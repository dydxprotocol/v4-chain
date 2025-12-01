package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
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
		// Setup
		userStats        *stattypes.UserStats
		globalStats      *stattypes.GlobalStats
		stakingTiers     []*types.StakingTier
		userBondedTokens *big.Int

		// Input
		req *types.QueryUserStakingTierRequest

		// Expected
		expectedError    error
		expectedResponse *types.QueryUserStakingTierResponse
	}{
		"returns error for nil request": {
			req:           nil,
			expectedError: status.Error(codes.InvalidArgument, "invalid request"),
		},
		"returns error for empty address": {
			req: &types.QueryUserStakingTierRequest{
				Address: "",
			},
			expectedError: status.Error(codes.InvalidArgument, "invalid bech32 address"),
		},
		"valid user with no bonded tokens (0 discount)": {
			userStats: &stattypes.UserStats{
				TakerNotional: 10,
				MakerNotional: 10,
			},
			globalStats: &stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      100_000,
						},
					},
				},
			},
			userBondedTokens: nil,
			req: &types.QueryUserStakingTierRequest{
				Address: constants.BobAccAddress.String(),
			},
			expectedResponse: &types.QueryUserStakingTierResponse{
				FeeTierName:      "1",
				StakedBaseTokens: dtypes.NewInt(0),
				DiscountPpm:      0,
			},
		},
		"valid user with bonded tokens, no staking tiers": {
			userStats: &stattypes.UserStats{
				TakerNotional: 10,
				MakerNotional: 10,
			},
			globalStats: &stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			stakingTiers:     nil,
			userBondedTokens: big.NewInt(5000),
			req: &types.QueryUserStakingTierRequest{
				Address: constants.AliceAccAddress.String(),
			},
			expectedResponse: &types.QueryUserStakingTierResponse{
				FeeTierName:      "1",
				StakedBaseTokens: dtypes.NewInt(5000),
				DiscountPpm:      0,
			},
		},
		"valid user in first fee tier and qualifies for staking discount": {
			userStats: &stattypes.UserStats{
				TakerNotional: 10,
				MakerNotional: 10,
			},
			globalStats: &stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      100_000, // 10% discount
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(5000),
							FeeDiscountPpm:      200_000, // 20% discount
						},
					},
				},
			},
			userBondedTokens: big.NewInt(5000),
			req: &types.QueryUserStakingTierRequest{
				Address: constants.AliceAccAddress.String(),
			},
			expectedResponse: &types.QueryUserStakingTierResponse{
				FeeTierName:      "1",
				StakedBaseTokens: dtypes.NewInt(5000),
				DiscountPpm:      200_000,
			},
		},
		"valid user in second fee tier and qualifies for staking discount": {
			userStats: &stattypes.UserStats{
				TakerNotional: 1_000_000_000_000,
				MakerNotional: 150,
			},
			globalStats: &stattypes.GlobalStats{
				NotionalTraded: 10_000_000_000_000,
			},
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "2",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      150_000, // 15%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(2000),
							FeeDiscountPpm:      250_000, // 25%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(3000),
							FeeDiscountPpm:      350_000, // 35%
						},
					},
				},
			},
			userBondedTokens: big.NewInt(2000),
			req: &types.QueryUserStakingTierRequest{
				Address: constants.AliceAccAddress.String(),
			},
			expectedResponse: &types.QueryUserStakingTierResponse{
				FeeTierName:      "2",
				StakedBaseTokens: dtypes.NewInt(2000),
				DiscountPpm:      250_000,
			},
		},
		"valid user doesn't qualify for staking discount": {
			userStats: &stattypes.UserStats{
				TakerNotional: 10,
				MakerNotional: 10,
			},
			globalStats: &stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(10000),
							FeeDiscountPpm:      200_000,
						},
					},
				},
			},
			userBondedTokens: big.NewInt(500),
			req: &types.QueryUserStakingTierRequest{
				Address: constants.BobAccAddress.String(),
			},
			expectedResponse: &types.QueryUserStakingTierResponse{
				FeeTierName:      "1",
				StakedBaseTokens: dtypes.NewInt(500),
				DiscountPpm:      0,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper

			// Set up staking tiers
			if tc.stakingTiers != nil {
				err := k.SetStakingTiers(ctx, tc.stakingTiers)
				require.NoError(t, err)
			}

			// Set up user stats
			if tc.userStats != nil && tc.req != nil && tc.req.Address != "" {
				statsKeeper := tApp.App.StatsKeeper
				statsKeeper.SetUserStats(ctx, tc.req.Address, tc.userStats)
				statsKeeper.SetGlobalStats(ctx, tc.globalStats)
			}

			// Set up user bonded tokens
			if tc.req != nil && tc.req.Address != "" {
				statsKeeper := tApp.App.StatsKeeper
				bondedAmount := big.NewInt(0)
				if tc.userBondedTokens != nil {
					bondedAmount = tc.userBondedTokens
				}
				statsKeeper.UnsafeSetCachedStakedBaseTokens(ctx, tc.req.Address, &stattypes.CachedStakedBaseTokens{
					StakedBaseTokens: dtypes.NewIntFromBigInt(bondedAmount),
					CachedAt:         ctx.BlockTime().Unix(),
				})
			}

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
