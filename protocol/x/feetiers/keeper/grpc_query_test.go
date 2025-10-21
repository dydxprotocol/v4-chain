package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
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

// TestPerMarketFeeDiscountParams tests the PerMarketFeeDiscountParams query handler
func TestPerMarketFeeDiscountParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set up a test fee discount params
	clobPairId := uint32(42)
	discountParams := types.PerMarketFeeDiscountParams{
		ClobPairId: clobPairId,
		StartTime:  time.Unix(1100, 0).UTC(),
		EndTime:    time.Unix(1200, 0).UTC(),
		ChargePpm:  500_000, // 50% discount
	}

	// Set current block time for validation
	ctx = ctx.WithBlockTime(time.Unix(1000, 0))
	err := k.SetPerMarketFeeDiscountParams(ctx, discountParams)
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		req *types.QueryPerMarketFeeDiscountParamsRequest
		res *types.QueryPerMarketFeeDiscountParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryPerMarketFeeDiscountParamsRequest{
				ClobPairId: clobPairId,
			},
			res: &types.QueryPerMarketFeeDiscountParamsResponse{
				Params: discountParams,
			},
			err: nil,
		},
		"Nil": {
			req: nil,
			res: nil,
			err: status.Error(codes.InvalidArgument, "invalid request"),
		},
		"Not Found": {
			req: &types.QueryPerMarketFeeDiscountParamsRequest{
				ClobPairId: 999, // non-existent CLOB pair ID
			},
			res: nil,
			err: status.Error(codes.NotFound, "fee discount not found for the specified market/CLOB pair"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			res, err := k.PerMarketFeeDiscountParams(ctx, tc.req)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

// TestAllMarketFeeDiscountParams tests the AllMarketFeeDiscountParams query handler
func TestAllMarketFeeDiscountParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Set current block time for validation
	ctx = ctx.WithBlockTime(time.Unix(1000, 0))

	// Set up multiple test fee discount params
	discountParams := []types.PerMarketFeeDiscountParams{
		{
			ClobPairId: 1,
			StartTime:  time.Unix(1100, 0).UTC(),
			EndTime:    time.Unix(1200, 0).UTC(),
			ChargePpm:  0, // 100% discount (free)
		},
		{
			ClobPairId: 2,
			StartTime:  time.Unix(1150, 0).UTC(),
			EndTime:    time.Unix(1250, 0).UTC(),
			ChargePpm:  500_000, // 50% discount
		},
		{
			ClobPairId: 3,
			StartTime:  time.Unix(1200, 0).UTC(),
			EndTime:    time.Unix(1300, 0).UTC(),
			ChargePpm:  750_000, // 25% discount
		},
	}

	// Store the fee discount params
	for _, params := range discountParams {
		err := k.SetPerMarketFeeDiscountParams(ctx, params)
		require.NoError(t, err)
	}

	for name, tc := range map[string]struct {
		req *types.QueryAllMarketFeeDiscountParamsRequest
		res *types.QueryAllMarketFeeDiscountParamsResponse
		err error
	}{
		"Success": {
			req: &types.QueryAllMarketFeeDiscountParamsRequest{},
			res: &types.QueryAllMarketFeeDiscountParamsResponse{
				Params: discountParams,
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
			res, err := k.AllMarketFeeDiscountParams(ctx, tc.req)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
				// We can't guarantee the order of the returned fee discount params, so we need to compare them differently
				require.Equal(t, len(tc.res.Params), len(res.Params))

				// Create a map to make comparison easier
				paramsMap := make(map[uint32]types.PerMarketFeeDiscountParams)
				for _, p := range res.Params {
					paramsMap[p.ClobPairId] = p
				}

				// Check that each expected params entry is in the result
				for _, expected := range tc.res.Params {
					actual, found := paramsMap[expected.ClobPairId]
					require.True(t, found)
					require.Equal(t, expected.ClobPairId, actual.ClobPairId)
					require.Equal(t, expected.StartTime, actual.StartTime)
					require.Equal(t, expected.EndTime, actual.EndTime)
					require.Equal(t, expected.ChargePpm, actual.ChargePpm)
				}
			}
		})
	}
}

// TestAllMarketFeeDiscountParamsEmpty tests the AllMarketFeeDiscountParams query handler with no params
func TestAllMarketFeeDiscountParamsEmpty(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.FeeTiersKeeper

	// Don't set any fee discount params - test empty response
	req := &types.QueryAllMarketFeeDiscountParamsRequest{}
	res, err := k.AllMarketFeeDiscountParams(ctx, req)

	// Should succeed with empty params list
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Empty(t, res.Params)
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
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      10000,
						},
					},
				},
				{
					FeeTierName: "2",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(500),
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
