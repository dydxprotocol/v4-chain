package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	affiliateskeeper "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	feetierskeeper "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	statskeeper "github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	logger := tApp.App.FeeTiersKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

func TestGetPerpetualFeePpm(t *testing.T) {
	tests := map[string]struct {
		user                string
<<<<<<< HEAD
		UserStats           *stattypes.UserStats
		GlobalStats         *stattypes.GlobalStats
=======
		userStats           *stattypes.UserStats
		globalStats         *stattypes.GlobalStats
		setupFeeDiscount    bool
		discountParams      types.PerMarketFeeDiscountParams
		setupTime           *time.Time
		blockTime           time.Time
		clobPairId          uint32
		stakingTiers        []*types.StakingTier
		userBondedTokens    *big.Int
>>>>>>> c667de27 (consider staking tiers when calculating fees (#3195))
		expectedTakerFeePpm int32
		expectedMakerFeePpm int32
	}{
		"regular user, first tier": {
			"alice",
			&stattypes.UserStats{
				TakerNotional: 10,
				MakerNotional: 10,
			},
			&stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			10,
			1,
		},
		"regular user, increased tier": {
			"alice",
			&stattypes.UserStats{
				TakerNotional: 1_000,
				MakerNotional: 150,
			},
			&stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			20,
			2,
		},
		"regular user, partial requirements doesn't increase tier": {
			"alice",
			&stattypes.UserStats{
				TakerNotional: 1_000,
				MakerNotional: 1_000_000_000,
			},
			&stattypes.GlobalStats{
				NotionalTraded: 10_000_000_000,
			},
			20,
			2,
		},
		"regular user, top tier": {
			"alice",
			&stattypes.UserStats{
				TakerNotional: 1_000,
				MakerNotional: 1_000_000_000,
			},
			&stattypes.GlobalStats{
				NotionalTraded: 2_000_000_000,
			},
			30,
			3,
		},
		"vault is top tier regardless of stats": {
			constants.Vault_Clob0.ToModuleAccountAddress(),
			&stattypes.UserStats{
				TakerNotional: 10,
				MakerNotional: 10,
			},
			&stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			30,
			3,
		},
		{
			name: "staking discount applies to positive maker fees in tier 3",
			user: constants.AliceAccAddress.String(),
			userStats: &stattypes.UserStats{
				TakerNotional: 1_000,
				MakerNotional: 1_000_000_000,
			},
			globalStats: &stattypes.GlobalStats{
				NotionalTraded: 2_000_000_000,
			},
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "3",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      500_000, // 50% discount
						},
					},
				},
			},
			userBondedTokens:    big.NewInt(1000),
			blockTime:           time.Now().UTC(),
			clobPairId:          1,
			expectedTakerFeePpm: 15, // 30 * 0.5 = 15
			expectedMakerFeePpm: 1,  // 3 * 0.5 = 1.5, rounded to 1
		},
		{
			name: "staking discount combined with per-market discount",
			user: constants.AliceAccAddress.String(),
			userStats: &stattypes.UserStats{
				TakerNotional: 10,
				MakerNotional: 10,
			},
			globalStats: &stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			setupFeeDiscount: true,
			discountParams: types.PerMarketFeeDiscountParams{
				ClobPairId: 1,
				StartTime:  time.Unix(1000, 0).UTC(),
				EndTime:    time.Unix(3000, 0).UTC(),
				ChargePpm:  500_000, // 50% per-market discount
			},
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      200_000, // 20% staking discount
						},
					},
				},
			},
			userBondedTokens:    big.NewInt(1000),
			blockTime:           time.Unix(2000, 0).UTC(),
			clobPairId:          1,
			expectedTakerFeePpm: 4, // 10 * 0.5 (market) * 0.8 (staking) = 4
			expectedMakerFeePpm: 0, // 1 * 0.5 * 0.8 = 0.4, rounded to 0
		},
		{
			name: "higher tier with 20% staking discount",
			user: constants.AliceAccAddress.String(),
			userStats: &stattypes.UserStats{
				TakerNotional: 1_000,
				MakerNotional: 150,
			},
			globalStats: &stattypes.GlobalStats{
				NotionalTraded: 10_000,
			},
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "2",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(5000),
							FeeDiscountPpm:      200_000, // 20% discount
						},
					},
				},
			},
			userBondedTokens:    big.NewInt(5000),
			blockTime:           time.Now().UTC(),
			clobPairId:          1,
			expectedTakerFeePpm: 16, // 20 * 0.8 = 16
			expectedMakerFeePpm: 1,  // 2 * 0.8 = 1.6, rounded to 1
		},
		{
			name: "100% staking discount",
			user: constants.AliceAccAddress.String(),
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
							FeeDiscountPpm:      1_000_000, // 100% discount (free)
						},
					},
				},
			},
			userBondedTokens:    big.NewInt(10000),
			blockTime:           time.Now().UTC(),
			clobPairId:          1,
			expectedTakerFeePpm: 0, // 10 * 0 = 0
			expectedMakerFeePpm: 0, // 1 * 0 = 0
		},
		{
			name: "doesn't qualify for staking discount",
			user: constants.BobAccAddress.String(),
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
							FeeDiscountPpm:      200_000, // 20% discount
						},
					},
				},
			},
			userBondedTokens:    big.NewInt(500), // Not enough to qualify
			blockTime:           time.Now().UTC(),
			clobPairId:          1,
			expectedTakerFeePpm: 10, // No discount applied
			expectedMakerFeePpm: 1,  // No discount applied
		},
		{
			name: "user with no bonded tokens",
			user: constants.CarlAccAddress.String(),
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
							FeeDiscountPpm:      200_000, // 20% discount
						},
					},
				},
			},
			userBondedTokens:    big.NewInt(0),
			blockTime:           time.Now().UTC(),
			clobPairId:          1,
			expectedTakerFeePpm: 10, // No discount applied
			expectedMakerFeePpm: 1,  // No discount applied
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			tApp.App.VaultKeeper.AddVaultToAddressStore(ctx, constants.Vault_Clob0)
			k := tApp.App.FeeTiersKeeper
			err := tApp.App.AffiliatesKeeper.UpdateAffiliateParameters(ctx, &affiliatetypes.MsgUpdateAffiliateParameters{
				AffiliateParameters: affiliatetypes.AffiliateParameters{
					RefereeMinimumFeeTierIdx: 2,
				},
			})
			require.NoError(t, err)
			err = k.SetPerpetualFeeParams(
				ctx,
				types.PerpetualFeeParams{
					Tiers: []*types.PerpetualFeeTier{
						{
							Name:        "1",
							TakerFeePpm: 10,
							MakerFeePpm: 1,
						},
						{
							Name:                      "2",
							AbsoluteVolumeRequirement: 1_000,
							TakerFeePpm:               20,
							MakerFeePpm:               2,
						},
						{
							Name:                           "3",
							AbsoluteVolumeRequirement:      1_000_000_000,
							MakerVolumeShareRequirementPpm: 500_000,
							TakerFeePpm:                    30,
							MakerFeePpm:                    3,
						},
					},
				},
			)
			require.NoError(t, err)

<<<<<<< HEAD
=======
			// Setup fee discount if needed
			if tc.setupFeeDiscount {
				// Use setupCtx with the appropriate time for setting up the discount
				err = k.SetPerMarketFeeDiscountParams(setupCtx, tc.discountParams)
				require.NoError(t, err)

				// Verify fee discount was set correctly
				params, err := k.GetPerMarketFeeDiscountParams(ctx, tc.discountParams.ClobPairId)
				require.NoError(t, err)
				require.Equal(t, tc.discountParams.ClobPairId, params.ClobPairId)
				require.Equal(t, tc.discountParams.StartTime, params.StartTime)
				require.Equal(t, tc.discountParams.EndTime, params.EndTime)
				require.Equal(t, tc.discountParams.ChargePpm, params.ChargePpm)
			}

			// Set up staking tiers if needed
			if tc.stakingTiers != nil {
				err = k.SetStakingTiers(ctx, tc.stakingTiers)
				require.NoError(t, err)
			}

			// Set up stats
>>>>>>> c667de27 (consider staking tiers when calculating fees (#3195))
			statsKeeper := tApp.App.StatsKeeper
			statsKeeper.SetUserStats(ctx, tc.user, tc.UserStats)
			statsKeeper.SetGlobalStats(ctx, tc.GlobalStats)

<<<<<<< HEAD
			require.Equal(t, tc.expectedTakerFeePpm, k.GetPerpetualFeePpm(ctx, tc.user, true, 2))
			require.Equal(t, tc.expectedMakerFeePpm, k.GetPerpetualFeePpm(ctx, tc.user, false, 2))
=======
			// Set up user bonded tokens
			if tc.userBondedTokens != nil {
				statsKeeper.UnsafeSetCachedStakedAmount(ctx, tc.user, &stattypes.CachedStakeAmount{
					StakedAmount: dtypes.NewIntFromBigInt(tc.userBondedTokens),
					CachedAt:     ctx.BlockTime().Unix(),
				})
			}

			require.Equal(t, tc.expectedTakerFeePpm, k.GetPerpetualFeePpm(ctx, tc.user, true, 2, tc.clobPairId))
			require.Equal(t, tc.expectedMakerFeePpm, k.GetPerpetualFeePpm(ctx, tc.user, false, 2, tc.clobPairId))
>>>>>>> c667de27 (consider staking tiers when calculating fees (#3195))
		})
	}
}

func TestGetPerpetualFeePpm_Referral(t *testing.T) {
	testFeePerpetualParams := types.PerpetualFeeParams{
		Tiers: []*types.PerpetualFeeTier{
			{
				Name:        "1",
				TakerFeePpm: 10,
				MakerFeePpm: 1,
			},
			{
				Name:                      "2",
				AbsoluteVolumeRequirement: 1_000,
				TakerFeePpm:               20,
				MakerFeePpm:               2,
			},
			{
				Name:                           "3",
				AbsoluteVolumeRequirement:      1_000_000_000,
				MakerVolumeShareRequirementPpm: 500_000,
				TakerFeePpm:                    30,
				MakerFeePpm:                    3,
			},
			{
				Name:                           "4",
				AbsoluteVolumeRequirement:      2_000_000_000,
				MakerVolumeShareRequirementPpm: 600_000,
				TakerFeePpm:                    40,
				MakerFeePpm:                    4,
			},
			{
				Name:                           "5",
				AbsoluteVolumeRequirement:      5_000_000_000,
				MakerVolumeShareRequirementPpm: 700_000,
				TakerFeePpm:                    50,
				MakerFeePpm:                    5,
			},
		},
	}
	tests := map[string]struct {
		expectedTakerFeePpm int32
		setup               func(ctx sdk.Context, affiliatesKeeper *affiliateskeeper.Keeper, statsKeeper *statskeeper.Keeper)
	}{
		"regular user, first tier, no referral": {
			expectedTakerFeePpm: 10,
			setup: func(ctx sdk.Context, affiliatesKeeper *affiliateskeeper.Keeper, statsKeeper *statskeeper.Keeper) {
				statsKeeper.SetUserStats(ctx, constants.AliceAccAddress.String(), &stattypes.UserStats{
					TakerNotional: 10,
					MakerNotional: 10,
				})
				statsKeeper.SetGlobalStats(ctx, &stattypes.GlobalStats{
					NotionalTraded: 10_000,
				})
			},
		},
		"regular user, referral": {
			expectedTakerFeePpm: 30,
			setup: func(ctx sdk.Context, affiliatesKeeper *affiliateskeeper.Keeper, statsKeeper *statskeeper.Keeper) {
				statsKeeper.SetUserStats(ctx, constants.AliceAccAddress.String(), &stattypes.UserStats{
					TakerNotional: 10,
					MakerNotional: 10,
				})
				statsKeeper.SetGlobalStats(ctx, &stattypes.GlobalStats{
					NotionalTraded: 10_000,
				})

				err := affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(), constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		"regular user, referral, already in tier 3": {
			expectedTakerFeePpm: 30,
			setup: func(ctx sdk.Context, affiliatesKeeper *affiliateskeeper.Keeper, statsKeeper *statskeeper.Keeper) {
				statsKeeper.SetUserStats(ctx, constants.AliceAccAddress.String(), &stattypes.UserStats{
					TakerNotional: 10,
					MakerNotional: 1_000_000_000,
				})
				statsKeeper.SetGlobalStats(ctx, &stattypes.GlobalStats{
					NotionalTraded: 1_000_000_000,
				})

				err := affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(), constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		"regular user, referral, above tier 3": {
			expectedTakerFeePpm: 40,
			setup: func(ctx sdk.Context, affiliatesKeeper *affiliateskeeper.Keeper, statsKeeper *statskeeper.Keeper) {
				statsKeeper.SetUserStats(ctx, constants.AliceAccAddress.String(), &stattypes.UserStats{
					TakerNotional: 10,
					MakerNotional: 2_000_000_001,
				})
				statsKeeper.SetGlobalStats(ctx, &stattypes.GlobalStats{
					NotionalTraded: 3_000_000_000,
				})

				err := affiliatesKeeper.RegisterAffiliate(ctx, constants.AliceAccAddress.String(), constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper
			err := k.SetPerpetualFeeParams(
				ctx,
				testFeePerpetualParams,
			)
			require.NoError(t, err)

			statsKeeper := tApp.App.StatsKeeper
			affiliatesKeeper := tApp.App.AffiliatesKeeper

			err = affiliatesKeeper.UpdateAffiliateParameters(ctx, &affiliatetypes.MsgUpdateAffiliateParameters{
				AffiliateParameters: affiliatetypes.AffiliateParameters{
					RefereeMinimumFeeTierIdx: 2,
				},
			})
			require.NoError(t, err)

			// common setup
			err = affiliatesKeeper.UpdateAffiliateTiers(ctx, affiliatetypes.DefaultAffiliateTiers)
			require.NoError(t, err)

			if tc.setup != nil {
				tc.setup(ctx, &affiliatesKeeper, &statsKeeper)
			}

			require.Equal(t, tc.expectedTakerFeePpm,
				k.GetPerpetualFeePpm(ctx, constants.AliceAccAddress.String(), true, 2))
		})
	}
}

func TestGetMaxMakerRebate(t *testing.T) {
	tests := map[string]struct {
		expectedLowestMakerFee int32
		feeTiers               []*types.PerpetualFeeTier
	}{
		"all positive maker fee": {
			feeTiers: []*types.PerpetualFeeTier{
				{
					Name:        "1",
					TakerFeePpm: 200,
					MakerFeePpm: 100,
				},
				{
					Name:        "2",
					TakerFeePpm: 100,
					MakerFeePpm: 50,
				},
			},
			expectedLowestMakerFee: 50,
		},
		"includes 0 maker fee": {
			feeTiers: []*types.PerpetualFeeTier{
				{
					Name:        "1",
					TakerFeePpm: 200,
					MakerFeePpm: 100,
				},
				{
					Name:        "2",
					TakerFeePpm: 100,
					MakerFeePpm: 50,
				},
				{
					Name:        "3",
					TakerFeePpm: 75,
					MakerFeePpm: 0,
				},
			},
			expectedLowestMakerFee: 0,
		},
		"includes negative maker fee": {
			feeTiers: []*types.PerpetualFeeTier{
				{
					Name:        "1",
					TakerFeePpm: 200,
					MakerFeePpm: 100,
				},
				{
					Name:        "2",
					TakerFeePpm: 100,
					MakerFeePpm: 50,
				},
				{
					Name:        "3",
					TakerFeePpm: 75,
					MakerFeePpm: 0,
				},
				{
					Name:        "3",
					TakerFeePpm: 60,
					MakerFeePpm: -10,
				},
				{
					Name:        "4",
					TakerFeePpm: 50,
					MakerFeePpm: -20,
				},
			},
			expectedLowestMakerFee: -20,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper
			err := k.SetPerpetualFeeParams(
				ctx,
				types.PerpetualFeeParams{
					Tiers: tc.feeTiers,
				},
			)
			require.NoError(t, err)

			require.Equal(t, tc.expectedLowestMakerFee, k.GetLowestMakerFee(ctx))
		})
	}
}

func TestGetAffiliateRefereeLowestTakerFee(t *testing.T) {
	tests := map[string]struct {
		expectedLowestTakerFee int32
		feeTiers               types.PerpetualFeeParams
	}{
		"tiers are ordered by absolute volume requirement": {
			feeTiers:               types.StandardParams(),
			expectedLowestTakerFee: 350,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper
			err := k.SetPerpetualFeeParams(
				ctx,
				tc.feeTiers,
			)
			require.NoError(t, err)

			feeTiers := k.GetPerpetualFeeParams(ctx).Tiers
			require.Equal(t, tc.expectedLowestTakerFee, feetierskeeper.GetAffiliateRefereeLowestTakerFeeFromTiers(feeTiers))
		})
	}
}
