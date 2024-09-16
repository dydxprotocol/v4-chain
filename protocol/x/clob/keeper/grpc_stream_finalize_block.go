package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// Returns the order updates needed to update the fill amount for the orders
// from local ops queue, according to the latest onchain state (after FinalizeBlock).
// Effectively reverts the optimistic fill amounts removed from the CheckTx to DeliverTx state transition.
func (k Keeper) getUpdatesToSyncLocalOpsQueue(
	ctx sdk.Context,
) *types.OffchainUpdates {
	localValidatorOperationsQueue, _ := k.MemClob.GetOperationsToReplay(ctx)
	fetchOrdersInvolvedInOpQueue(localValidatorOperationsQueue)
	orderIdsFromLocal := fetchOrdersInvolvedInOpQueue(
		localValidatorOperationsQueue,
	)
	allUpdates := types.NewOffchainUpdates()
	for orderId := range orderIdsFromLocal {
		orderbookUpdate := k.MemClob.GetOrderbookUpdatesForOrderUpdate(ctx, orderId)
		allUpdates.Append(orderbookUpdate)
	}
	return allUpdates
}

// Grpc Streaming logic after consensus agrees on a block.
// - Stream all events staged during `FinalizeBlock`.
// - Stream orderbook updates to sync fills in local ops queue.
func (k Keeper) StreamBatchUpdatesAfterFinalizeBlock(
	ctx sdk.Context,
) {
	defer telemetry.MeasureSince(
		time.Now(),
		types.ModuleName,
		metrics.StreamBatchUpdatesAfterFinalizeBlock,
		metrics.Latency,
	)
	orderBookUpdatesToSyncLocalOpsQueue := k.getUpdatesToSyncLocalOpsQueue(ctx)

	k.GetFullNodeStreamingManager().StreamBatchUpdatesAfterFinalizeBlock(
		ctx,
		orderBookUpdatesToSyncLocalOpsQueue,
		k.PerpetualIdToClobPairId,
	)
}
