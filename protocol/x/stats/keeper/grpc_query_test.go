package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
)

func TestParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithNonDeterminismChecksEnabled(false).Build()
	ctx := tApp.InitChain()
	k := tApp.App.StatsKeeper

	for name, tc := range map[string]struct {
		req *types.QueryParamsRequest
		res *types.QueryParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryParamsRequest{},
			res: &types.QueryParamsResponse{
				Params: types.DefaultGenesis().Params,
			},
			err: nil,
		},
		//"Nil": {
		//	req: nil,
		//	res: nil,
		//	err: status.Error(codes.InvalidArgument, "invalid request"),
		//},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.Params(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestStatsMetadata(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.StatsKeeper
	statsMetadata := &types.StatsMetadata{
		TrailingEpoch: 10,
	}
	k.SetStatsMetadata(ctx, statsMetadata)

	for name, tc := range map[string]struct {
		req *types.QueryStatsMetadataRequest
		res *types.QueryStatsMetadataResponse
		err error
	}{
		"Success": {
			req: &types.QueryStatsMetadataRequest{},
			res: &types.QueryStatsMetadataResponse{
				Metadata: statsMetadata,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.StatsMetadata(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestGlobalStats(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.StatsKeeper
	globalStats := &types.GlobalStats{
		NotionalTraded: 10,
	}
	k.SetGlobalStats(ctx, globalStats)

	for name, tc := range map[string]struct {
		req *types.QueryGlobalStatsRequest
		res *types.QueryGlobalStatsResponse
		err error
	}{
		"Success": {
			req: &types.QueryGlobalStatsRequest{},
			res: &types.QueryGlobalStatsResponse{
				Stats: globalStats,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.GlobalStats(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestUserStats(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.StatsKeeper
	user := "alice"
	userStats := &types.UserStats{
		TakerNotional: 10,
		MakerNotional: 10,
	}
	k.SetUserStats(ctx, user, userStats)

	for name, tc := range map[string]struct {
		req *types.QueryUserStatsRequest
		res *types.QueryUserStatsResponse
		err error
	}{
		"Success": {
			req: &types.QueryUserStatsRequest{
				User: user,
			},
			res: &types.QueryUserStatsResponse{
				Stats: userStats,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.UserStats(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
