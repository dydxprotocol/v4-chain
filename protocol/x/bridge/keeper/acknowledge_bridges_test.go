package keeper_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAcknowledgeBridges(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Bridge events to acknowledge.
		bridgeEvents []types.BridgeEvent
		// Whether bridging is disabled.
		bridgingDisabled bool
		// Error responses of mock delayMsgKeeper.
		delayMsgErrors []error

		/* --- Expectations --- */
		// Expected AcknowledgedEventInfo.
		expectedAEI types.BridgeEventInfo
		// Expected error.
		expectedError string
	}{
		"Success: no events": {
			bridgeEvents: []types.BridgeEvent{},
			expectedAEI: types.BridgeEventInfo{
				NextId:         0,
				EthBlockHeight: 0,
			},
		},
		"Success: 1 event": {
			bridgeEvents: []types.BridgeEvent{
				constants.BridgeEvent_Id55_Height15,
			},
			delayMsgErrors: []error{nil},
			expectedAEI: types.BridgeEventInfo{
				NextId:         56,
				EthBlockHeight: 15,
			},
		},
		"Success: 2 events": {
			bridgeEvents: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
			delayMsgErrors: []error{nil, nil},
			expectedAEI: types.BridgeEventInfo{
				NextId:         2,
				EthBlockHeight: 0,
			},
		},
		"Error: bridging disabled": {
			bridgeEvents: []types.BridgeEvent{
				constants.BridgeEvent_Id55_Height15,
			},
			delayMsgErrors:   []error{nil},
			bridgingDisabled: true,
			expectedError:    types.ErrBridgingDisabled.Error(),
		},
		"Error: 2 events, delaying second msg returns error": {
			bridgeEvents: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
			delayMsgErrors: []error{nil, errors.New("failed to delay message 1")},
			expectedError:  "failed to delay message 1",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize context, keeper, and mockDelayMsgKeeper.
			ctx, bridgeKeeper, _, _, _, _, mockDelayMsgKeeper := keepertest.BridgeKeepers(t)
			err := bridgeKeeper.UpdateSafetyParams(ctx, types.SafetyParams{
				IsDisabled:  tc.bridgingDisabled,
				DelayBlocks: bridgeKeeper.GetSafetyParams(ctx).DelayBlocks,
			})
			require.NoError(t, err)
			for i := range tc.bridgeEvents {
				mockDelayMsgKeeper.On(
					"DelayMessageByBlocks",
					ctx,
					mock.Anything,
					mock.Anything,
				).Return(uint32(i), tc.delayMsgErrors[i]).Once()
			}
			initialAei := bridgeKeeper.GetAcknowledgedEventInfo(ctx)

			// Invoke AcknowledgeBridges.
			err = bridgeKeeper.AcknowledgeBridges(ctx, tc.bridgeEvents)

			if tc.expectedError != "" {
				// Verify that error is as expected.
				require.ErrorContains(t, err, tc.expectedError)

				// Verify that AcknowledgedEventInfo was not updated.
				require.Equal(t, initialAei, bridgeKeeper.GetAcknowledgedEventInfo(ctx))

				if tc.bridgingDisabled {
					// Verify that no messages were delayed.
					mockDelayMsgKeeper.AssertNotCalled(t, "DelayMessageByBlocks")
				}
			} else {
				// Verify that AcknowledgeBridges returns no error.
				require.NoError(t, err)

				// Verify that AcknowledgedEventInfo is updated in state.
				aei := bridgeKeeper.GetAcknowledgedEventInfo(ctx)
				require.Equal(t, tc.expectedAEI, aei)

				// Assert mock expectations.
				mockDelayMsgKeeper.AssertExpectations(t)
			}
		})
	}
}

func TestGetAcknowledgeBridges(t *testing.T) {
	timeNow := time.Now()

	tests := map[string]struct {
		// Setup.
		blockTimestamp        time.Time
		eventTimestamp        time.Time
		proposeParams         types.ProposeParams
		bridgeEventsToAdd     []types.BridgeEvent
		acknowledgedEventInfo types.BridgeEventInfo
		bridgingDisabled      bool

		// Expectations.
		expectedMsg *types.MsgAcknowledgeBridges
	}{
		"Empty events due to probabilistic skipping": {
			blockTimestamp: timeNow,
			eventTimestamp: timeNow.Add(-time.Second * 2),
			proposeParams: types.ProposeParams{
				// 100% skip rate.
				SkipRatePpm: uint32(constants.OneMillion),
				// propose events recognized at least one second ago.
				ProposeDelayDuration: time.Second,
			},
			bridgeEventsToAdd: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
			expectedMsg: &types.MsgAcknowledgeBridges{
				Events: []types.BridgeEvent{},
			},
		},
		"Empty events due to deterministic skipping": {
			// Skip proposing bridge events as blockTimestamp <= timeNow - SkipIfBlockDelayedByDuration.
			blockTimestamp: timeNow.Add(-time.Second),
			proposeParams: types.ProposeParams{
				SkipRatePpm:                  0, // do not skip based on pseudo-randomness.
				SkipIfBlockDelayedByDuration: time.Second,
			},
			bridgeEventsToAdd: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
			expectedMsg: &types.MsgAcknowledgeBridges{
				Events: []types.BridgeEvent{},
			},
		},
		"More than MaxBridgesPerBlock events recognized": {
			blockTimestamp: timeNow,
			eventTimestamp: timeNow.Add(-time.Second * 2),
			proposeParams: types.ProposeParams{
				SkipRatePpm:                  0,           // do not skip based on pseudo-randomness.
				SkipIfBlockDelayedByDuration: time.Second, // do not skip based on time.
				MaxBridgesPerBlock:           3,           // propose up to 3 events per block.
				ProposeDelayDuration:         time.Second, // propose events recognized at least one second ago.
			},
			bridgeEventsToAdd: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id2_Height1,
				constants.BridgeEvent_Id3_Height3, // this event should not be proposed.
			},
			expectedMsg: &types.MsgAcknowledgeBridges{
				Events: []types.BridgeEvent{
					constants.BridgeEvent_Id0_Height0,
					constants.BridgeEvent_Id1_Height0,
					constants.BridgeEvent_Id2_Height1,
				},
			},
		},
		"MaxBridgesPerBlock events recognized": {
			blockTimestamp: timeNow,
			eventTimestamp: timeNow.Add(-time.Second * 2),
			proposeParams: types.ProposeParams{
				SkipRatePpm:                  0,           // do not skip based on pseudo-randomness.
				SkipIfBlockDelayedByDuration: time.Second, // do not skip based on time.
				MaxBridgesPerBlock:           3,           // propose up to 3 events per block.
				ProposeDelayDuration:         time.Second, // propose events recognized at least one second ago.
			},
			bridgeEventsToAdd: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id2_Height1,
			},
			expectedMsg: &types.MsgAcknowledgeBridges{
				Events: []types.BridgeEvent{
					constants.BridgeEvent_Id0_Height0,
					constants.BridgeEvent_Id1_Height0,
					constants.BridgeEvent_Id2_Height1,
				},
			},
		},
		"Fewer than MaxBridgesPerBlock events recognized": {
			blockTimestamp: timeNow,
			eventTimestamp: timeNow.Add(-time.Second * 2),
			proposeParams: types.ProposeParams{
				SkipRatePpm:                  0,           // do not skip based on pseudo-randomness.
				SkipIfBlockDelayedByDuration: time.Second, // do not skip based on time.
				MaxBridgesPerBlock:           3,           // propose up to 3 events per block.
				ProposeDelayDuration:         time.Second, // propose events recognized at least one second ago.
			},
			bridgeEventsToAdd: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
			expectedMsg: &types.MsgAcknowledgeBridges{
				Events: []types.BridgeEvent{
					constants.BridgeEvent_Id0_Height0,
					constants.BridgeEvent_Id1_Height0,
				},
			},
		},
		"Already acknowledged events are not proposed": {
			blockTimestamp: timeNow,
			eventTimestamp: timeNow.Add(-time.Second * 2),
			proposeParams: types.ProposeParams{
				SkipRatePpm:                  0,           // do not skip based on pseudo-randomness.
				SkipIfBlockDelayedByDuration: time.Second, // do not skip based on time.
				MaxBridgesPerBlock:           3,           // propose up to 3 events per block.
				ProposeDelayDuration:         time.Second, // propose events recognized at least one second ago.
			},
			bridgeEventsToAdd: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0, // this event should not be proposed.
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id2_Height1,
			},
			acknowledgedEventInfo: types.BridgeEventInfo{
				NextId:         1,
				EthBlockHeight: 0,
			},
			expectedMsg: &types.MsgAcknowledgeBridges{
				Events: []types.BridgeEvent{
					constants.BridgeEvent_Id1_Height0,
					constants.BridgeEvent_Id2_Height1,
				},
			},
		},
		"Events recognized at or after cutoff time are not proposed": {
			blockTimestamp: timeNow,
			eventTimestamp: timeNow.Add(-time.Second),
			proposeParams: types.ProposeParams{
				SkipRatePpm:                  0,           // do not skip based on pseudo-randomness.
				SkipIfBlockDelayedByDuration: time.Second, // do not skip based on time.
				MaxBridgesPerBlock:           3,           // propose up to 3 events per block.
				ProposeDelayDuration:         time.Second, // propose events recognized at least one second ago.
			},
			bridgeEventsToAdd: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
			},
			expectedMsg: &types.MsgAcknowledgeBridges{
				Events: []types.BridgeEvent{},
			},
		},
		"No event is proposed when bridging is disabled": {
			blockTimestamp: timeNow,
			eventTimestamp: timeNow.Add(-time.Second * 2),
			proposeParams: types.ProposeParams{
				SkipRatePpm:                  0,           // do not skip based on pseudo-randomness.
				SkipIfBlockDelayedByDuration: time.Second, // do not skip based on time.
				MaxBridgesPerBlock:           3,           // propose up to 3 events per block.
				ProposeDelayDuration:         time.Second, // propose events recognized at least one second ago.
			},
			bridgeEventsToAdd: []types.BridgeEvent{
				constants.BridgeEvent_Id0_Height0,
				constants.BridgeEvent_Id1_Height0,
				constants.BridgeEvent_Id2_Height1,
			},
			acknowledgedEventInfo: types.BridgeEventInfo{
				NextId:         0,
				EthBlockHeight: 0,
			},
			bridgingDisabled: true,
			expectedMsg: &types.MsgAcknowledgeBridges{
				Events: []types.BridgeEvent{},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper, bridgeEventManager, and mockTimeProvider.
			ctx, bridgeKeeper, _, mockTimeProvider, bridgeEventManager, _, _ := keepertest.BridgeKeepers(t)
			err := bridgeKeeper.SetAcknowledgedEventInfo(ctx, tc.acknowledgedEventInfo)
			require.NoError(t, err)
			err = bridgeKeeper.UpdateProposeParams(ctx, tc.proposeParams)
			require.NoError(t, err)
			err = bridgeKeeper.UpdateSafetyParams(ctx, types.SafetyParams{
				IsDisabled:  tc.bridgingDisabled,
				DelayBlocks: bridgeKeeper.GetSafetyParams(ctx).DelayBlocks,
			})
			require.NoError(t, err)
			mockTimeProvider.On("Now").Return(tc.eventTimestamp).Once()
			err = bridgeEventManager.AddBridgeEvents(tc.bridgeEventsToAdd)
			require.NoError(t, err)

			// Get MsgAcknowledgeBridges.
			mockTimeProvider.On("Now").Return(timeNow).Once()
			msg := bridgeKeeper.GetAcknowledgeBridges(ctx, tc.blockTimestamp)

			// Assert expected MsgAcknowledgeBridges.
			require.Equal(t, tc.expectedMsg, msg)
		})
	}
}
