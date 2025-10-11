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

func TestStakingTiers(t *testing.T) {
	tests := map[string]struct {
		// Setup
		stakingTiers []*types.StakingTier

		// Input
		req *types.QueryStakingTiersRequest

		// Expected
		expectedError      error
		expectedTiersCount int
	}{
		"returns empty when nothing set": {
			stakingTiers:       []*types.StakingTier{},
			req:                &types.QueryStakingTiersRequest{},
			expectedError:      nil,
			expectedTiersCount: 0,
		},
		"returns staking tiers correctly": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: "100",
							FeeDiscountPpm:      10000,
						},
					},
				},
				{
					FeeTierName: "2",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: "500",
							FeeDiscountPpm:      25000,
						},
					},
				},
			},
			req:                &types.QueryStakingTiersRequest{},
			expectedError:      nil,
			expectedTiersCount: 2,
		},
		"returns error for nil request": {
			stakingTiers:  []*types.StakingTier{},
			req:           nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper

			// Set staking tiers
			err := k.SetStakingTiers(ctx, tc.stakingTiers)
			require.NoError(t, err)

			// Verify query
			resp, err := k.StakingTiers(ctx, tc.req)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Len(t, resp.StakingTiers, tc.expectedTiersCount)
				require.ElementsMatch(t, tc.stakingTiers, resp.StakingTiers)
			}
		})
	}
}

func TestUserStakingTier(t *testing.T) {
	tests := map[string]struct {
		// Input
		req *types.QueryUserStakingTierRequest

		// Expected
		expectedError    error
		expectedResponse *types.QueryUserStakingTierResponse
	}{
		// TODO: add valid test cases after implementation
		"returns error for nil request": {
			req:           nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid request"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper

			// Verify query
			resp, err := k.UserStakingTier(ctx, tc.req)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.expectedResponse, resp)
			}
		})
	}
}
