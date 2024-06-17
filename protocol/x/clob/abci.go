package clob

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// PreBlocker executes all ABCI PreBlock logic respective to the clob module.
func PreBlocker(
	ctx sdk.Context,
	keeper types.ClobKeeper,
) {
	keeper.Initialize(ctx)
}

// BeginBlocker executes all ABCI BeginBlock logic respective to the clob module.
func BeginBlocker(
	ctx sdk.Context,
	keeper types.ClobKeeper,
) {
	ctx = log.AddPersistentTagsToLogger(ctx,
		log.Handler, log.BeginBlocker,
		log.BlockHeight, ctx.BlockHeight(),
	)

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
	ctx = log.AddPersistentTagsToLogger(ctx,
		log.Handler, log.EndBlocker,
		log.BlockHeight, ctx.BlockHeight(),
	)

	processProposerMatchesEvents := keeper.GetProcessProposerMatchesEvents(ctx)

	// Prune any fill amounts from state which are now past their `pruneableBlockHeight`.
	keeper.PruneStateFillAmountsForShortTermOrders(ctx)

	// Prune expired XOrders from state.
	keeper.ExpireXOrders(ctx)

	// Prune expired stateful orders completely from state.
	expiredStatefulOrderIds := keeper.RemoveExpiredStatefulOrdersTimeSlices(ctx, ctx.BlockTime())
	for _, orderId := range expiredStatefulOrderIds {
		// Remove the order fill amount from state.
		keeper.RemoveOrderFillAmount(ctx, orderId)

		// Delete the stateful order placement from state.
		keeper.DeleteLongTermOrderPlacement(ctx, orderId)

		// Emit an on-chain indexer event for Stateful Order Expiration.
		keeper.GetIndexerEventManager().AddBlockEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK,
			indexerevents.StatefulOrderEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewStatefulOrderRemovalEvent(
					orderId,
					indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_EXPIRED,
				),
			),
		)
		metrics.IncrCounterWithLabels(
			metrics.ClobExpiredStatefulOrders,
			1,
			orderId.GetOrderIdLabels()...,
		)
	}

	// Prune expired untriggered conditional orders from the in-memory UntriggeredConditionalOrders struct.
	keeper.PruneUntriggeredConditionalOrders(
		expiredStatefulOrderIds,
		processProposerMatchesEvents.PlacedStatefulCancellationOrderIds,
	)

	// Update the memstore with expired order ids.
	// These expired stateful order ids will be purged from the memclob in `Commit`.
	processProposerMatchesEvents.ExpiredStatefulOrderIds = expiredStatefulOrderIds

	// Before triggering conditional orders, add newly-placed conditional orders to the clob keeper's
	// in-memory UntriggeredConditionalOrders data structure to allow conditional orders to
	// trigger in the same block they are placed. Skip triggering orders which have been cancelled
	// or expired.
	// TODO(CLOB-773) Support conditional order replacements. Ensure replacements are de-duplicated.
	keeper.AddUntriggeredConditionalOrders(
		ctx,
		processProposerMatchesEvents.PlacedConditionalOrderIds,
		lib.UniqueSliceToSet(processProposerMatchesEvents.GetPlacedStatefulCancellationOrderIds()),
		lib.UniqueSliceToSet(expiredStatefulOrderIds),
	)

	// Trigger X Orders.
	err := keeper.TriggerXOrders(ctx)
	ctx.Logger().Error("TriggerXOrders", "err", err)

	// Poll out all triggered conditional orders from `UntriggeredConditionalOrders` and update state.
	triggeredConditionalOrderIds := keeper.MaybeTriggerConditionalOrders(ctx)
	// Update the memstore with conditional order ids triggered in the last block.
	// These triggered conditional orders will be placed in the `PrepareCheckState``.
	processProposerMatchesEvents.ConditionalOrderIdsTriggeredInLastBlock = triggeredConditionalOrderIds

	// Write the ProcessProposerMatchcesEvents with all the EndBlocker updates to state.
	keeper.MustSetProcessProposerMatchesEvents(
		ctx,
		processProposerMatchesEvents,
	)

	// Prune any rate limiting information that is no longer relevant.
	keeper.PruneRateLimits(ctx)

	// Emit relevant metrics at the end of every block.
	metrics.SetGauge(
		metrics.InsuranceFundBalance,
		metrics.GetMetricValueFromBigInt(keeper.GetCrossInsuranceFundBalance(ctx)),
	)
}

// PrepareCheckState executes all ABCI PrepareCheckState logic respective to the clob module.
func PrepareCheckState(
	ctx sdk.Context,
	keeper *keeper.Keeper,
) {
	ctx = log.AddPersistentTagsToLogger(ctx,
		log.Handler, log.PrepareCheckState,
		// Prepare check state is for the next block.
		log.BlockHeight, ctx.BlockHeight()+1,
	)

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
	localValidatorOperationsQueue, shortTermOrderTxBytes := keeper.MemClob.GetOperationsToReplay(ctx)

	log.DebugLog(ctx, "Clearing local operations queue",
		log.LocalValidatorOperationsQueue, types.GetInternalOperationsQueueTextString(localValidatorOperationsQueue),
	)

	keeper.MemClob.RemoveAndClearOperationsQueue(ctx, localValidatorOperationsQueue)

	// 2. Purge invalid state from the memclob.
	offchainUpdates := types.NewOffchainUpdates()
	offchainUpdates = keeper.MemClob.PurgeInvalidMemclobState(
		ctx,
		processProposerMatchesEvents.OrderIdsFilledInLastBlock,
		processProposerMatchesEvents.ExpiredStatefulOrderIds,
		processProposerMatchesEvents.PlacedStatefulCancellationOrderIds,
		processProposerMatchesEvents.RemovedStatefulOrderIds,
		offchainUpdates,
	)

	// 3. Place all stateful order placements included in the last block on the memclob.
	// Note telemetry is measured outside of the function call because `PlaceStatefulOrdersFromLastBlock`
	// is called within `PlaceConditionalOrdersTriggeredInLastBlock`.
	startPlaceLongTermOrders := time.Now()
	offchainUpdates = keeper.PlaceStatefulOrdersFromLastBlock(
		ctx,
		processProposerMatchesEvents.PlacedLongTermOrderIds,
		offchainUpdates,
	)
	telemetry.MeasureSince(
		startPlaceLongTermOrders,
		types.ModuleName,
		metrics.PlaceLongTermOrdersFromLastBlock,
		metrics.Latency,
	)
	telemetry.SetGauge(
		float32(len(processProposerMatchesEvents.PlacedLongTermOrderIds)),
		types.ModuleName,
		metrics.PlaceLongTermOrdersFromLastBlock,
		metrics.Count,
	)

	// 4. Place all conditional orders triggered in EndBlocker of last block on the memclob.
	offchainUpdates = keeper.PlaceConditionalOrdersTriggeredInLastBlock(
		ctx,
		processProposerMatchesEvents.ConditionalOrderIdsTriggeredInLastBlock,
		offchainUpdates,
	)

	// 5. Replay the local validator’s operations onto the book.
	replayUpdates := keeper.MemClob.ReplayOperations(
		ctx,
		localValidatorOperationsQueue,
		shortTermOrderTxBytes,
		offchainUpdates,
	)

	// TODO(CLOB-275): Do not gracefully handle panics in `PrepareCheckState`.
	if replayUpdates != nil {
		offchainUpdates = replayUpdates
	}

	// 6. Get all potentially liquidatable subaccount IDs and attempt to liquidate them.
	liquidatableSubaccountIds := keeper.DaemonLiquidationInfo.GetLiquidatableSubaccountIds()
	subaccountsToDeleverage, err := keeper.LiquidateSubaccountsAgainstOrderbook(ctx, liquidatableSubaccountIds)
	if err != nil {
		panic(err)
	}
	// Add subaccounts with open positions in final settlement markets to the slice of subaccounts/perps
	// to be deleveraged.
	subaccountsToDeleverage = append(
		subaccountsToDeleverage,
		keeper.GetSubaccountsWithPositionsInFinalSettlementMarkets(ctx)...,
	)

	// 7. Deleverage subaccounts.
	// TODO(CLOB-1052) - decouple steps 6 and 7 by using DaemonLiquidationInfo.NegativeTncSubaccounts
	// as the input for this function.
	if err := keeper.DeleverageSubaccounts(ctx, subaccountsToDeleverage); err != nil {
		panic(err)
	}

	// 8. Gate withdrawals by inserting a zero-fill deleveraging operation into the operations queue if any
	// of the negative TNC subaccounts still have negative TNC after liquidations and deleveraging steps.
	negativeTncSubaccountIds := keeper.DaemonLiquidationInfo.GetNegativeTncSubaccountIds()
	if err := keeper.GateWithdrawalsIfNegativeTncSubaccountSeen(ctx, negativeTncSubaccountIds); err != nil {
		panic(err)
	}

	// Send all off-chain Indexer events
	keeper.SendOffchainMessages(offchainUpdates, nil, metrics.SendPrepareCheckStateOffchainUpdates)

	newLocalValidatorOperationsQueue, _ := keeper.MemClob.GetOperationsToReplay(ctx)

	log.DebugLog(ctx, "Local operations queue after PrepareCheckState",
		log.NewLocalValidatorOperationsQueue,
		types.GetInternalOperationsQueueTextString(newLocalValidatorOperationsQueue),
	)

	// Initialize new GRPC streams with orderbook snapshots, if any.
	keeper.InitializeNewGrpcStreams(ctx)

	// Set per-orderbook gauges.
	keeper.MemClob.SetMemclobGauges(ctx)
}
