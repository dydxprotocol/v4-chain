package keeper_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestQueryParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

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
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
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

func TestQueryRewardShare(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RewardsKeeper

	rewardShare := types.RewardShare{
		Address: TestAddress1,
		Weight:  dtypes.NewInt(12_345_678),
	}

	err := k.SetRewardShare(ctx, rewardShare)
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		req *types.QueryRewardShareRequest
		res *types.QueryRewardShareResponse
		err error
	}{
		"Success": {
			req: &types.QueryRewardShareRequest{
				Address: rewardShare.Address,
			},
			res: &types.QueryRewardShareResponse{
				RewardShare: rewardShare,
			},
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.RewardShare(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
