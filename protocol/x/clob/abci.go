package clob

import (
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	liquidationtypes "github.com/dydxprotocol/v4/daemons/server/types/liquidations"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// BeginBlocker executes all ABCI BeginBlock logic respective to the clob module.
func BeginBlocker(
	ctx sdk.Context,
	keeper types.ClobKeeper,
) {
	// Initialize the set of process proposer match events for the next block effectively
	// removing any events that occurred in the last block.
	keeper.MustSetProcessProposerMatchesEvents(
		ctx,
		types.ProcessProposerMatchesEvents{
			BlockHeight: lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
		},
	)
}

// EndBlocker executes all ABCI EndBlock logic respective to the clob module.
func EndBlocker(
	ctx sdk.Context,
	keeper keeper.Keeper,
) {
	processProposerMatchesEvents := keeper.GetProcessProposerMatchesEvents(ctx)
	previousBlockStatefulOrderCancellations := processProposerMatchesEvents.GetPlacedStatefulCancellations()

	// Create a set of placed stateful cancellations by `OrderId`.
	placedStatefulCancellationOrderIds := lib.SliceToSet(previousBlockStatefulOrderCancellations)
	// Create a set of removed stateful orders by `OrderId`.
	removedStatefulOrderIds := lib.SliceToSet(processProposerMatchesEvents.RemovedStatefulOrderIds)

	// Retrieve the fill amounts for all orders which were filled in the last block, and populate
	// the `offchainUpdates` with updates for the fill amounts.
	offchainUpdates := types.NewOffchainUpdates()
	for _, orderId := range processProposerMatchesEvents.OrdersIdsFilledInLastBlock {
		// Skip sending order updates for orders that have been cancelled / removed since they have been
		// already removed from state.
		_, cancelled := placedStatefulCancellationOrderIds[orderId]
		_, removed := removedStatefulOrderIds[orderId]
		if cancelled || removed {
			continue
		}

		exists, fillAmount, _ := keeper.GetOrderFillAmount(ctx, orderId)
		if !exists {
			ctx.Logger().Error(
				fmt.Sprintf(
					"EndBlocker: order fill amount does not exist in state for Indexer event for orderId %v",
					orderId,
				),
			)
			continue
		}

		if message, success := off_chain_updates.CreateOrderUpdateMessage(
			ctx.Logger(),
			orderId,
			fillAmount,
		); success {
			offchainUpdates.AddUpdateMessage(orderId, message)
		}
	}

	// Prune any fill amounts from state which are now past their `pruneableBlockHeight`.
	keeper.PruneStateFillAmountsForShortTermOrders(ctx)

	// Update the block time of the previously committed block.
	keeper.SetBlockTimeForLastCommittedBlock(ctx)

	// Prune expired stateful orders completely from state.
	expiredStatefulOrderIds := keeper.RemoveExpiredStatefulOrdersTimeSlices(ctx, ctx.BlockTime())
	for _, orderId := range expiredStatefulOrderIds {
		// Remove the order fill amount from state.
		keeper.RemoveOrderFillAmount(ctx, orderId)

		// Delete the stateful order placement from state.
		keeper.DeleteLongTermOrderPlacement(ctx, orderId)

		// Emit an on-chain indexer event for Stateful Order Expiration.
		keeper.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexer_manager.GetB64EncodedEventMessage(
				indexerevents.NewStatefulOrderExpirationEvent(
					orderId,
				),
			),
		)
		telemetry.IncrCounter(1, types.ModuleName, metrics.Expired, metrics.StatefulOrderRemoved, metrics.Count)
	}

	// Update the memstore with expired order ids.
	// These expired stateful order ids will be purged from the memclob in `Commit`.
	processProposerMatchesEvents.ExpiredStatefulOrderIds = expiredStatefulOrderIds
	keeper.MustSetProcessProposerMatchesEvents(
		ctx,
		processProposerMatchesEvents,
	)

	// Send all off-chain Indexer updates with these new fill amounts.
	keeper.SendOffchainMessages(offchainUpdates, nil, metrics.SendPrepareCheckStateOffchainUpdates)

	// Prune any rate limiting information that is no longer relevant.
	keeper.PruneRateLimits(ctx)
}

// PrepareCheckState executes all ABCI PrepareCheckState logic respective to the clob module.
func PrepareCheckState(
	ctx sdk.Context,
	keeper keeper.Keeper,
	memClob types.MemClob,
	liquidatableSubaccountIds *liquidationtypes.LiquidatableSubaccountIds,
) {
	// Get the events generated from processing the matches in the latest block.
	processProposerMatchesEvents := keeper.GetProcessProposerMatchesEvents(ctx)
	if ctx.BlockHeight() != int64(processProposerMatchesEvents.BlockHeight) {
		panic(
			fmt.Errorf(
				"block height %d for ProcessProposerMatchesEvents does not equal current block height %d",
				processProposerMatchesEvents.BlockHeight,
				ctx.BlockHeight(),
			),
		)
	}

	// 1. Remove all operations in the local validators operations queue from the memclob.
	localValidatorOperationsQueue, shortTermOrderTxBytes := memClob.GetOperationsToReplay(ctx)
	ctx.Logger().Debug(
		"Clearing local operations queue",
		"localValidatorOperationsQueue",
		types.GetInternalOperationsQueueTextString(localValidatorOperationsQueue),
		"block",
		ctx.BlockHeight(),
	)

	memClob.RemoveAndClearOperationsQueue(ctx, localValidatorOperationsQueue)

	// 2. Purge invalid state from the memclob.
	offchainUpdates := types.NewOffchainUpdates()
	offchainUpdates = memClob.PurgeInvalidMemclobState(
		ctx,
		processProposerMatchesEvents.OrdersIdsFilledInLastBlock,
		processProposerMatchesEvents.ExpiredStatefulOrderIds,
		processProposerMatchesEvents.PlacedStatefulCancellations,
		processProposerMatchesEvents.RemovedStatefulOrderIds,
		offchainUpdates,
	)

	// 3. Place all stateful order placements included in the last block on the memclob.
	offchainUpdates = keeper.PlaceStatefulOrdersFromLastBlock(
		ctx,
		processProposerMatchesEvents.PlacedStatefulOrders,
		offchainUpdates,
	)

	// 4. Replay the local validator’s operations onto the book.
	replayUpdates := memClob.ReplayOperations(
		ctx,
		localValidatorOperationsQueue,
		shortTermOrderTxBytes,
		offchainUpdates,
	)

	// TODO(CLOB-275): Do not gracefully handle panics in `PrepareCheckState`.
	if replayUpdates != nil {
		offchainUpdates = replayUpdates
	}

	// 5. Get all potentially liquidatable subaccount IDs and attempt to liquidate them.
	subaccountIds := liquidatableSubaccountIds.GetSubaccountIds()

	telemetry.ModuleSetGauge(
		types.ModuleName,
		float32(len(subaccountIds)),
		metrics.Liquidations,
		metrics.LiquidatableSubaccountIds,
		metrics.Count,
	)

	// Get the liquidation order for each subaccount.
	liquidationOrders := make([]types.LiquidationOrder, 0)
	for _, subaccountId := range subaccountIds {
		// If attempting to liquidate a subaccount returns an error, panic.
		liquidationOrder, err := keeper.MaybeGetLiquidationOrder(ctx, subaccountId)
		if err != nil {
			panic(err)
		}

		if liquidationOrder != nil {
			liquidationOrders = append(liquidationOrders, *liquidationOrder)
		}
	}

	// Sort liquidation orders by clob pair id, then by fillable price, then by order hash.
	sort.Sort(types.SortedLiquidationOrders(liquidationOrders))

	// Attempt to place each liquidation order and perform deleveraging if necessary.
	for _, liquidationOrder := range liquidationOrders {
		if _, _, err := keeper.PlacePerpetualLiquidation(ctx, liquidationOrder); err != nil {
			ctx.Logger().Error(
				fmt.Sprintf(
					"Failed to liquidate subaccount. Liquidation Order: (%+v). Err: %v",
					liquidationOrder,
					err,
				),
			)
			panic(err)
		}
	}

	// Send all off-chain Indexer events
	keeper.SendOffchainMessages(offchainUpdates, nil, metrics.SendPrepareCheckStateOffchainUpdates)

	newLocalValidatorOperationsQueue, _ := memClob.GetOperationsToReplay(ctx)
	ctx.Logger().Debug(
		"Local operations queue after PrepareCheckState",
		"newLocalValidatorOperationsQueue",
		types.GetInternalOperationsQueueTextString(newLocalValidatorOperationsQueue),
		"block",
		ctx.BlockHeight(),
	)

	// Set per-orderbook gauges.
	memClob.SetMemclobGauges(ctx)
}
