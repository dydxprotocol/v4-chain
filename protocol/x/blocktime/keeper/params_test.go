package keeper_test

import (
	"testing"
	"time"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/stretchr/testify/require"
)

func TestGetDowntimeParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.BlockTimeKeeper

	require.Equal(t, types.DefaultGenesis().Params, k.GetDowntimeParams(ctx))
}

func TestSetDowntimeParams_Success(t *testing.T) {
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

	params := types.DowntimeParams{
		Durations: []time.Duration{
			5 * time.Second,
			20 * time.Second,
			25 * time.Second,
			35 * time.Second,
			45 * time.Second,
		},
	}
	expectedAllDowntimeInfo := &types.AllDowntimeInfo{
		Infos: []*types.AllDowntimeInfo_DowntimeInfo{
			{
				Duration: 5 * time.Second,
				BlockInfo: types.BlockInfo{
					Height:    40,
					Timestamp: time.Unix(400, 0).UTC(),
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
				Duration: 25 * time.Second,
				BlockInfo: types.BlockInfo{
					Height:    20,
					Timestamp: time.Unix(200, 0).UTC(),
				},
			},
			{
				Duration: 35 * time.Second,
				BlockInfo: types.BlockInfo{
					Height:    10,
					Timestamp: time.Unix(100, 0).UTC(),
				},
			},
			{
				Duration: 45 * time.Second,
				BlockInfo: types.BlockInfo{
					Height:    1,
					Timestamp: time.Unix(10, 0).UTC(),
				},
			},
		},
	}
	require.NoError(t, params.Validate())

	require.NoError(t, k.SetDowntimeParams(ctx, params))
	require.Equal(t, params, k.GetDowntimeParams(ctx))
	require.Equal(t, expectedAllDowntimeInfo, k.GetAllDowntimeInfo(ctx))
}
