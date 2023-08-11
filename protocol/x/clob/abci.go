package clob

import (
	"fmt"

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
	previousBlockStatefulOrderCancellations := lib.SliceToSet(
		processProposerMatchesEvents.GetPreviousBlockStatefulOrderCancellations(),
	)

	// Retrieve the fill amounts for all orders which were filled in the last block, and populate
	// the `offchainUpdates` with updates for the fill amounts.
	offchainUpdates := types.NewOffchainUpdates()
	for _, orderId := range processProposerMatchesEvents.OrdersIdsFilledInLastBlock {
		// Skip sending order updates for orders that have been cancelled since they have been
		// already removed from state.
		if _, cancelled := previousBlockStatefulOrderCancellations[orderId]; cancelled {
			continue
		}

		exists, fillAmount, _ := keeper.GetOrderFillAmount(ctx, orderId)
		if !exists {
			ctx.Logger().Error(
				fmt.Sprintf(
					"PrepareCheckState: order fill amount does not exist in state for Indexer event for orderId %v",
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

	// Prune expired seen place orders from the in-memory map.
	keeper.PruneExpiredSeenPlaceOrders(ctx, lib.MustConvertIntegerToUint32(ctx.BlockHeight()))

	// Update the block time of the previously committed block.
	keeper.SetBlockTimeForLastCommittedBlock(ctx)

	// Prune expired stateful orders completely from state.
	expiredStatefulOrderIds := keeper.RemoveExpiredStatefulOrdersTimeSlices(ctx, ctx.BlockTime())
	for _, orderId := range expiredStatefulOrderIds {
		// Remove the order fill amount from state.
		keeper.RemoveOrderFillAmount(ctx, orderId)

		// Delete the stateful order placement from state.
		keeper.DeleteStatefulOrderPlacement(ctx, orderId)

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
	localValidatorOperationsQueue := memClob.GetOperations(ctx)
	ctx.Logger().Info(
		"Clearing local operations queue",
		"localValidatorOperationsQueue",
		types.GetOperationsQueueTextString(localValidatorOperationsQueue),
		"block",
		ctx.BlockHeight(),
	)

	memClob.RemoveAndClearOperationsQueue(ctx, localValidatorOperationsQueue)

	// 2. Purge invalid state from the memclob.
	offchainUpdates := types.NewOffchainUpdates()
	previousBlockStatefulOrderCancellations := processProposerMatchesEvents.GetPreviousBlockStatefulOrderCancellations()
	offchainUpdates = memClob.PurgeInvalidMemclobState(
		ctx,
		processProposerMatchesEvents.OrdersIdsFilledInLastBlock,
		processProposerMatchesEvents.ExpiredStatefulOrderIds,
		previousBlockStatefulOrderCancellations,
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
		offchainUpdates,
		previousBlockStatefulOrderCancellations,
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

	for _, subaccountId := range subaccountIds {
		// If attempting to liquidate a subaccount returns an error, panic.
		if err := keeper.MaybeLiquidateSubaccount(ctx, subaccountId); err != nil {
			panic(err)
		}
	}

	// Send all off-chain Indexer events
	keeper.SendOffchainMessages(offchainUpdates, nil, metrics.SendPrepareCheckStateOffchainUpdates)

	newLocalValidatorOperationsQueue := memClob.GetOperations(ctx)
	ctx.Logger().Info(
		"Local operations queue after PrepareCheckState",
		"newLocalValidatorOperationsQueue",
		types.GetOperationsQueueTextString(newLocalValidatorOperationsQueue),
		"block",
		ctx.BlockHeight(),
	)
}
