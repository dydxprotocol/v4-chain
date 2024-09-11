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

	telemetry.SetGauge(
		float32(len(stagedEvents)),
		types.ModuleName,
		metrics.GrpcStagedAllFinalizeBlockUpdates,
		metrics.Count,
	)

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

	telemetry.SetGauge(
		float32(len(finalizedFillUpdates)),
		types.ModuleName,
		metrics.GrpcStagedFillFinalizeBlockUpdates,
		metrics.Count,
	)

	k.GetFullNodeStreamingManager().SendFinalizedSubaccountUpdates(
		finalizedSubaccountUpdates,
		uint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)

	telemetry.SetGauge(
		float32(len(finalizedSubaccountUpdates)),
		types.ModuleName,
		metrics.GrpcStagedSubaccountFinalizeBlockUpdates,
		metrics.Count,
	)
}
