package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

func TestParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	for name, tc := range map[string]struct {
		req *types.QueryPerpetualFeeParamsRequest
		res *types.QueryPerpetualFeeParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryPerpetualFeeParamsRequest{},
			res: &types.QueryPerpetualFeeParamsResponse{
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
			res, err := k.PerpetualFeeParams(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestUserFeeTier(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	for name, tc := range map[string]struct {
		req *types.QueryUserFeeTierRequest
		res *types.QueryUserFeeTierResponse
		err error
	}{
		"Success": {
			req: &types.QueryUserFeeTierRequest{
				User: "dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4",
			},
			res: &types.QueryUserFeeTierResponse{
				Index: 0,
				Tier: &types.PerpetualFeeTier{
					Name:                           "1",
					AbsoluteVolumeRequirement:      0,
					TotalVolumeShareRequirementPpm: 0,
					MakerVolumeShareRequirementPpm: 0,
					MakerFeePpm:                    -110,
					TakerFeePpm:                    500,
				},
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
		"Malformed address": {
			req: &types.QueryUserFeeTierRequest{
				User: "alice",
			},
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid bech32 address"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.UserFeeTier(ctx, tc.req)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}
