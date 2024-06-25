package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) StreamOrderbookUpdates(
	req *types.StreamOrderbookUpdatesRequest,
	stream types.Query_StreamOrderbookUpdatesServer,
) error {
	err := k.GetGrpcStreamingManager().Subscribe(*req, stream)
	if err != nil {
		return err
	}

	// Keep this scope alive because once this scope exits - the stream is closed
	ctx := stream.Context()
	<-ctx.Done()
	return nil
}

// revertOptimisticFills sends OrderUpdates that effectively reverts the optimistic fills.
func (k Keeper) revertOptimisticFills(ctx sdk.Context) {
	// Send absolute fill amounts from local to the grpc stream.
	// This must be sent out to account for checkState being discarded and deliverState being used.
	localValidatorOperationsQueue, _ := k.MemClob.GetOperationsToReplay(ctx)
	orderIdsFromLocalOpsQueue := fetchOrdersInvolvedInOpQueue(
		localValidatorOperationsQueue,
	)
	allUpdates := types.NewOffchainUpdates()

	for orderId := range orderIdsFromLocalOpsQueue {
		orderbookUpdate := k.MemClob.GetOrderbookUpdatesForOrderUpdate(ctx, orderId)
		allUpdates.Append(orderbookUpdate)
	}
	k.SendOrderbookUpdates(ctx, allUpdates)
}

// streamOnchainFills sends updates resulting from onchain fills.
func (k Keeper) streamOnchainFills(ctx sdk.Context) {
	// Get onchain stream events stored in transient store.
	events := k.GetIndexerEventManager().GetOnchainStreamEvents(ctx)
	streamOrderbookFills := make([]types.StreamOrderbookFill, len(events))
	// Unmarshall the onchain stream events.
	// TODO(CT-940): Avoid unnecessary marshalling and unmarshalling from using transient store,
	// currently not a huge drawback given the low amount of onchain fills.
	for i, event := range events {
		var streamOrderbookFill types.StreamOrderbookFill
		if err := proto.Unmarshal(event.Event.DataBytes, &streamOrderbookFill); err != nil {
			panic(err)
		}
		streamOrderbookFills[i] = streamOrderbookFill
	}
	k.SendOrderbookFillUpdates(
		ctx,
		streamOrderbookFills,
	)
}

// StreamEventsAfterBlockFinalized streams all events resulted from block processing.
// It should be called after the block is finalized.
func (k Keeper) StreamEventsAfterBlockFinalized(
	ctx sdk.Context,
) {
	if streamingManager := k.GetGrpcStreamingManager(); !streamingManager.Enabled() {
		// Return early if grpc streams are not enabled.
		return
	}
	k.revertOptimisticFills(ctx)
	k.streamOnchainFills(ctx)
}
