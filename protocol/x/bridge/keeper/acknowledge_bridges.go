package keeper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

// GetAcknowledgeBridges returns a `MsgAcknowledgeBridges` for recognized but not-yet-acknowledged
// bridge events, up to a maximum number of `ProposeParams.MaxBridgesPerBlock`.
func (k Keeper) GetAcknowledgeBridges(
	ctx sdk.Context,
	blockTimestamp time.Time,
) (msg *types.MsgAcknowledgeBridges) {
	wallClock := k.bridgeEventManager.GetNow()
	proposeParams := k.GetProposeParams(ctx)

	// In order to ensure an upper-bound on liveness issues in the case that +1/3 of validators cannot
	// properly get logs from an Ethereum node, skip proposing bridge events if any of the following:
	// - rand.Uint32(1_000_000) < ProposeParams.skip_rate_ppm
	// - blockTimestamp ≤ wallClock - ProposeParams.skip_if_block_delayed_by_duration
	if uint32(rand.Intn(int(lib.OneMillion))) < proposeParams.SkipRatePpm ||
		!blockTimestamp.After(wallClock.Add(-proposeParams.SkipIfBlockDelayedByDuration)) {
		return &types.MsgAcknowledgeBridges{
			Events: []types.BridgeEvent{},
		}
	}

	// Measure latency if not skipping proposing bridge events.
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.GetAcknowledgeBridges,
		metrics.Latency,
	)
	acknowledgedEventInfo := k.GetAcknowledgedEventInfo(ctx)
	recognizedCutoffTime := wallClock.Add(-proposeParams.ProposeDelayDuration)
	events := make([]types.BridgeEvent, 0)
	for i := uint32(0); i < proposeParams.MaxBridgesPerBlock; i++ {
		// 1. Try to retrieve recognized event with id `NextId + i` from BridgeEventManager.
		eventToAcknowledge, eventRecognizedAt, found := k.bridgeEventManager.GetBridgeEventById(
			acknowledgedEventInfo.NextId + i)
		// Stop looking for events with higher IDs if event with current ID is not found.
		// This assumes that recognized events are assigned IDs that increment by 1 each time.
		if !found {
			break
		}

		// 2. Append the new event if it is recognized before the cutoff time.
		if eventRecognizedAt.Before(recognizedCutoffTime) {
			events = append(events, eventToAcknowledge)
		} else {
			// Stop looking for events with higher IDs if event with current ID is not old enough.
			// This assumes that events with lower IDs are recognized before events with higher IDs.
			break
		}
	}

	return &types.MsgAcknowledgeBridges{
		Events: events,
	}
}

// AcknowledgeBridges acknowledges a list of bridge events.
func (k Keeper) AcknowledgeBridges(
	ctx sdk.Context,
	bridgeEvents []types.BridgeEvent,
) (err error) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.AcknowledgeBridges,
		metrics.Latency,
	)

	if len(bridgeEvents) == 0 {
		return nil
	}

	// For each bridge event, delay a `MsgCompleteBridge` to be executed `safetyParams.DelayBlocks`
	// blocks in the future. Panic if fails to delay any of the messages.
	safetyParams := k.GetSafetyParams(ctx)
	delayMsgModuleAccAddrString := authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String()
	for _, bridgeEvent := range bridgeEvents {
		// delaymsg module should be the authority for completing bridges.
		msgCompleteBridge := types.MsgCompleteBridge{
			Authority: delayMsgModuleAccAddrString,
			Event:     bridgeEvent,
		}
		_, err := k.delayMsgKeeper.DelayMessageByBlocks(
			ctx,
			&msgCompleteBridge,
			safetyParams.DelayBlocks,
		)
		if err != nil {
			panic(
				fmt.Sprintf(
					"failed to delay completing bridge: %s",
					err.Error(),
				),
			)
		}
	}

	// Update `AcknowledgedEventInfo` in state.
	lastBridgeEvent := bridgeEvents[len(bridgeEvents)-1]
	if err = k.SetAcknowledgedEventInfo(ctx, types.BridgeEventInfo{
		NextId:         lastBridgeEvent.GetId() + 1,
		EthBlockHeight: lastBridgeEvent.GetEthBlockHeight(),
	}); err != nil {
		return err
	}

	return nil
}
