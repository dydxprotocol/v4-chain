package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// StreamStagedEventsAfterFinalizeBlock streams all events emitted during `FinalizeBlock`
// (PreBlocker + EndBlocker + DeliverTx + EndBlocker).
// It should be called after the consensus agrees on a block (e.g. Precommitter).
func (k Keeper) StreamStagedEventsAfterFinalizeBlock(
	ctx sdk.Context,
) {
	defer telemetry.MeasureSince(
		time.Now(),
		types.ModuleName,
		metrics.StreamStagedEventsAfterFinalizeBlock,
		metrics.Latency,
	)

	// Get onchain stream events stored in transient store.
	stagedEvents := k.GetFullNodeStreamingManager().GetStagedFinalizeBlockEvents(ctx)

	finalizedFillUpdates := []types.StreamOrderbookFill{}
	finalizedSubaccountUpdates := []satypes.StreamSubaccountUpdate{}

	for _, stagedEvent := range stagedEvents {
		switch event := stagedEvent.Event.(type) {
		case *types.StagedFinalizeBlockEvent_OrderFill:
			finalizedFillUpdates = append(finalizedFillUpdates, *event.OrderFill)
		case *types.StagedFinalizeBlockEvent_SubaccountUpdate:
			finalizedSubaccountUpdates = append(finalizedSubaccountUpdates, *event.SubaccountUpdate)
		}
	}

	k.SendOrderbookFillUpdates(
		ctx,
		finalizedFillUpdates,
	)

	k.GetFullNodeStreamingManager().SendFinalizedSubaccountUpdates(
		finalizedSubaccountUpdates,
		uint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)
}
