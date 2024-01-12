package keeper_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	logger := tApp.App.BlockTimeKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

func TestUpdateAllDowntimeInfo(t *testing.T) {
	tests := map[string]struct {
		previousBlockInfo       *types.BlockInfo
		currentTimestamp        time.Time
		inputAllDowntimeInfo    *types.AllDowntimeInfo
		expectedAllDowntimeInfo *types.AllDowntimeInfo
	}{
		"update downtime": {
			&types.BlockInfo{
				Height:    10,
				Timestamp: time.Unix(100, 0).UTC(),
			},
			time.Unix(105, 0).UTC(),
			&types.AllDowntimeInfo{
				Infos: []*types.AllDowntimeInfo_DowntimeInfo{
					{
						Duration: time.Second,
						BlockInfo: types.BlockInfo{
							Height:    1,
							Timestamp: time.Unix(10, 0).UTC(),
						},
					},
					{
						Duration: 5 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    1,
							Timestamp: time.Unix(10, 0).UTC(),
						},
					},
					{
						Duration: 10 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    1,
							Timestamp: time.Unix(10, 0).UTC(),
						},
					},
				},
			},
			&types.AllDowntimeInfo{
				Infos: []*types.AllDowntimeInfo_DowntimeInfo{
					{
						Duration: time.Second,
						BlockInfo: types.BlockInfo{
							Height:    11,
							Timestamp: time.Unix(105, 0).UTC(),
						},
					},
					{
						Duration: 5 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    11,
							Timestamp: time.Unix(105, 0).UTC(),
						},
					},
					{
						Duration: 10 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    1,
							Timestamp: time.Unix(10, 0).UTC(),
						},
					},
				},
			},
		},
		"no update": {
			&types.BlockInfo{
				Height:    1,
				Timestamp: time.Unix(100, 0).UTC(),
			},
			time.Unix(102, 0).UTC(),
			&types.AllDowntimeInfo{
				Infos: []*types.AllDowntimeInfo_DowntimeInfo{
					{
						Duration: 5 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    2,
							Timestamp: time.Unix(10, 0).UTC(),
						},
					},
					{
						Duration: 10 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    1,
							Timestamp: time.Unix(5, 0).UTC(),
						},
					},
				},
			},
			&types.AllDowntimeInfo{
				Infos: []*types.AllDowntimeInfo_DowntimeInfo{
					{
						Duration: 5 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    2,
							Timestamp: time.Unix(10, 0).UTC(),
						},
					},
					{
						Duration: 10 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    1,
							Timestamp: time.Unix(5, 0).UTC(),
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			tApp.InitChain()
			ctx := tApp.AdvanceToBlock(
				tc.previousBlockInfo.Height+1,
				testapp.AdvanceToBlockOptions{
					BlockTime: tc.currentTimestamp,
				},
			)
			k := tApp.App.BlockTimeKeeper
			k.SetPreviousBlockInfo(ctx, tc.previousBlockInfo)
			k.SetAllDowntimeInfo(ctx, tc.inputAllDowntimeInfo)

			k.UpdateAllDowntimeInfo(ctx)
			actual := k.GetAllDowntimeInfo(ctx)
			require.Equal(t, tc.expectedAllDowntimeInfo, actual)
		})
	}
}

func TestGetDowntimeInfoFor(t *testing.T) {
	tests := map[string]struct {
		duration             time.Duration
		expectedDowntimeInfo types.AllDowntimeInfo_DowntimeInfo
	}{
		"smaller than all durations": {
			duration: 5 * time.Second,
			expectedDowntimeInfo: types.AllDowntimeInfo_DowntimeInfo{
				Duration: 0,
				BlockInfo: types.BlockInfo{
					Height:    40,
					Timestamp: time.Unix(400, 0).UTC(),
				},
			},
		},
		"equal to duration": {
			duration: 20 * time.Second,
			expectedDowntimeInfo: types.AllDowntimeInfo_DowntimeInfo{
				Duration: 20 * time.Second,
				BlockInfo: types.BlockInfo{
					Height:    20,
					Timestamp: time.Unix(200, 0).UTC(),
				},
			},
		},
		"not equal to duration": {
			duration: 25 * time.Second,
			expectedDowntimeInfo: types.AllDowntimeInfo_DowntimeInfo{
				Duration: 20 * time.Second,
				BlockInfo: types.BlockInfo{
					Height:    20,
					Timestamp: time.Unix(200, 0).UTC(),
				},
			},
		},
		"greater than all durations": {
			duration: 45 * time.Second,
			expectedDowntimeInfo: types.AllDowntimeInfo_DowntimeInfo{
				Duration: 40 * time.Second,
				BlockInfo: types.BlockInfo{
					Height:    1,
					Timestamp: time.Unix(10, 0).UTC(),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			tApp.InitChain()
			ctx := tApp.AdvanceToBlock(
				40,
				testapp.AdvanceToBlockOptions{
					BlockTime: time.Unix(400, 0).UTC(),
				},
			)
			k := tApp.App.BlockTimeKeeper
			k.SetAllDowntimeInfo(ctx, &types.AllDowntimeInfo{
				Infos: []*types.AllDowntimeInfo_DowntimeInfo{
					{
						Duration: 10 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    30,
							Timestamp: time.Unix(300, 0).UTC(),
						},
					},
					{
						Duration: 20 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    20,
							Timestamp: time.Unix(200, 0).UTC(),
						},
					},
					{
						Duration: 30 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    10,
							Timestamp: time.Unix(100, 0).UTC(),
						},
					},
					{
						Duration: 40 * time.Second,
						BlockInfo: types.BlockInfo{
							Height:    1,
							Timestamp: time.Unix(10, 0).UTC(),
						},
					},
				},
			})

			actual := k.GetDowntimeInfoFor(ctx, tc.expectedDowntimeInfo.Duration)
			require.Equal(t, tc.expectedDowntimeInfo, actual)
		})
	}
}

func TestGetTimeSinceLastBlock(t *testing.T) {
	testPrevBlockHeight := uint32(5)
	tests := map[string]struct {
		prevBlockTime              time.Time
		currBlockTime              time.Time
		expectedTimeSinceLastBlock time.Duration
	}{
		"2 sec": {
			prevBlockTime:              time.Unix(100, 0).UTC(),
			currBlockTime:              time.Unix(102, 0).UTC(),
			expectedTimeSinceLastBlock: time.Second * 2,
		},
		"Realistic values": {
			prevBlockTime:              time.Unix(1_704_827_023, 123_000_000).UTC(),
			currBlockTime:              time.Unix(1_704_827_024, 518_000_000).UTC(),
			expectedTimeSinceLastBlock: time.Second*1 + time.Nanosecond*395_000_000,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			tApp.InitChain()

			ctx := tApp.AdvanceToBlock(
				testPrevBlockHeight,
				testapp.AdvanceToBlockOptions{
					BlockTime: tc.prevBlockTime,
				},
			)

			k := tApp.App.BlockTimeKeeper

			actual := k.GetTimeSinceLastBlock(ctx.WithBlockTime(tc.currBlockTime))
			require.Equal(
				t,
				tc.expectedTimeSinceLastBlock,
				actual,
			)
		})
	}
}
