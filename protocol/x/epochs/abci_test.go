package epochs_test

import (
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestBeginBlocker(t *testing.T) {
	tests := map[string]struct {
		epochInfosToCreate []types.EpochInfo
		nextBlockTimeSec   int64
		nextBlockHeight    int64
		expectedEpochInfos []types.EpochInfo
		expectedEvents     []sdk.Event
	}{
		"initialize `funding-tick` and `funding-sample`": {
			epochInfosToCreate: []types.EpochInfo{
				{
					Name:                string(types.FundingSampleEpochInfoName),
					Duration:            60,
					NextTick:            30,
					FastForwardNextTick: true,
					IsInitialized:       false,
				},
				{
					Name:                string(types.FundingTickEpochInfoName),
					Duration:            3600,
					NextTick:            0,
					FastForwardNextTick: true,
					IsInitialized:       false,
				},
			},
			nextBlockTimeSec: 1800000000,
			nextBlockHeight:  2,
			expectedEpochInfos: []types.EpochInfo{
				{
					Name:                   string(types.FundingSampleEpochInfoName),
					Duration:               60,
					NextTick:               1800000030,
					CurrentEpoch:           0,
					CurrentEpochStartBlock: 0,
					IsInitialized:          true,
					FastForwardNextTick:    true,
				},
				{
					Name:                   string(types.FundingTickEpochInfoName),
					Duration:               3600,
					NextTick:               1800003600,
					CurrentEpoch:           0,
					CurrentEpochStartBlock: 0,
					IsInitialized:          true,
					FastForwardNextTick:    true,
				},
			},
		},
		"initialized, nextTick not reached, same epoch": {
			epochInfosToCreate: []types.EpochInfo{
				{
					Name:          "name",
					Duration:      60,
					NextTick:      1800000060,
					IsInitialized: true,
				},
			},
			nextBlockTimeSec: 1800000055,
			nextBlockHeight:  50,
			expectedEpochInfos: []types.EpochInfo{
				{
					Name:                   "name",
					Duration:               60,
					NextTick:               1800000060,
					CurrentEpoch:           0,
					CurrentEpochStartBlock: 0,
					IsInitialized:          true,
				},
			},
		},
		"initialized, nextTick reached": {
			epochInfosToCreate: []types.EpochInfo{
				{
					Name:          "name",
					Duration:      60,
					NextTick:      1800000060,
					IsInitialized: true,
				},
			},
			nextBlockTimeSec: 1800000075,
			nextBlockHeight:  65,
			expectedEpochInfos: []types.EpochInfo{
				{
					Name:                   "name",
					Duration:               60,
					NextTick:               1800000120,
					CurrentEpoch:           1,
					CurrentEpochStartBlock: 65,
					IsInitialized:          true,
				},
			},
			expectedEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeNewEpoch,
					sdk.NewAttribute(types.AttributeKeyEpochInfoName, "name"),
					sdk.NewAttribute(types.AttributeKeyEpochNumber, "1"),
					sdk.NewAttribute(types.AttributeKeyEpochStartTickTime, "1800000060"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlockTime, "1800000075"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlock, "65"),
				),
			},
		},
		"two different epochs, both reached nextTick": {
			epochInfosToCreate: []types.EpochInfo{
				{
					Name:          "name",
					Duration:      60,
					IsInitialized: true,
					NextTick:      1800000060,
				},
				{
					Name:          "name2",
					Duration:      60,
					NextTick:      1800000000,
					IsInitialized: true,
				},
			},
			nextBlockTimeSec: 1800000075,
			nextBlockHeight:  65,
			expectedEpochInfos: []types.EpochInfo{
				{
					Name:                   "name",
					Duration:               60,
					NextTick:               1800000120,
					CurrentEpoch:           1,
					IsInitialized:          true,
					CurrentEpochStartBlock: 65,
				},
				{
					Name:                   "name2",
					Duration:               60,
					NextTick:               1800000060, // Only catch up one epoch
					CurrentEpoch:           1,
					IsInitialized:          true,
					CurrentEpochStartBlock: 65,
				},
			},
			expectedEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeNewEpoch,
					sdk.NewAttribute(types.AttributeKeyEpochInfoName, "name"),
					sdk.NewAttribute(types.AttributeKeyEpochNumber, "1"),
					sdk.NewAttribute(types.AttributeKeyEpochStartTickTime, "1800000060"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlockTime, "1800000075"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlock, "65"),
				),
				sdk.NewEvent(
					types.EventTypeNewEpoch,
					sdk.NewAttribute(types.AttributeKeyEpochInfoName, "name2"),
					sdk.NewAttribute(types.AttributeKeyEpochNumber, "1"),
					sdk.NewAttribute(types.AttributeKeyEpochStartTickTime, "1800000000"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlockTime, "1800000075"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlock, "65"),
				),
			},
		},
		"downtime recovery - only catch up one epoch every block": {
			epochInfosToCreate: []types.EpochInfo{
				{
					Name:          "name",
					Duration:      60,
					NextTick:      1800000060,
					IsInitialized: true,
				},
			},
			nextBlockTimeSec: 1800000601, // BeginBlocker() called after 600 seconds
			nextBlockHeight:  65,
			expectedEpochInfos: []types.EpochInfo{
				{
					Name:                   "name",
					Duration:               60,
					NextTick:               1800000120, // Only catch up one epoch
					CurrentEpoch:           1,
					IsInitialized:          true,
					CurrentEpochStartBlock: 65,
				},
			},
			expectedEvents: []sdk.Event{
				sdk.NewEvent(
					types.EventTypeNewEpoch,
					sdk.NewAttribute(types.AttributeKeyEpochInfoName, "name"),
					sdk.NewAttribute(types.AttributeKeyEpochNumber, "1"),
					sdk.NewAttribute(types.AttributeKeyEpochStartTickTime, "1800000060"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlockTime, "1800000601"),
					sdk.NewAttribute(types.AttributeKeyEpochStartBlock, "65"),
				),
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			initBlockHeight := int64(0)
			initBlockTimeSec := int64(1800000001)

			ctx, keeper, _ := keepertest.EpochsKeeper(t)
			initblockTime := time.Unix(initBlockTimeSec, 0)
			createCtx := ctx.WithBlockTime(initblockTime).WithBlockHeight(initBlockHeight)
			for _, epoch := range tc.epochInfosToCreate {
				err := keeper.CreateEpochInfo(createCtx, epoch)
				require.NoError(t, err)
			}

			nextBlockTime := time.Unix(tc.nextBlockTimeSec, 0)
			nextCtx := ctx.WithBlockTime(nextBlockTime).WithBlockHeight(tc.nextBlockHeight)
			epochs.BeginBlocker(nextCtx, *keeper)
			require.Equal(t,
				tc.expectedEpochInfos,
				keeper.GetAllEpochInfo(nextCtx),
			)
			require.ElementsMatch(t,
				sdk.Events(tc.expectedEvents).ToABCIEvents(),
				ctx.EventManager().ABCIEvents(),
			)
		})
	}
}
