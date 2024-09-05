package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	affiliateskeeper "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
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
		UserStats           *stattypes.UserStats
		GlobalStats         *stattypes.GlobalStats
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			tApp.App.VaultKeeper.AddVaultToAddressStore(ctx, constants.Vault_Clob0)
			k := tApp.App.FeeTiersKeeper
			err := k.SetPerpetualFeeParams(
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

			statsKeeper := tApp.App.StatsKeeper
			statsKeeper.SetUserStats(ctx, tc.user, tc.UserStats)
			statsKeeper.SetGlobalStats(ctx, tc.GlobalStats)

			require.Equal(t, tc.expectedTakerFeePpm, k.GetPerpetualFeePpm(ctx, tc.user, true))
			require.Equal(t, tc.expectedMakerFeePpm, k.GetPerpetualFeePpm(ctx, tc.user, false))
		})
	}
}

func TestGetPerpetualFeePpm_Referral(t *testing.T) {
	tests := map[string]struct {
		expectedTakerFeePpm int32
		setup               func(ctx sdk.Context, affiliatesKeeper *affiliateskeeper.Keeper)
	}{
		"regular user, first tier, no referral": {
			expectedTakerFeePpm: 10,
			setup:               func(ctx sdk.Context, affiliatesKeeper *affiliateskeeper.Keeper) {},
		},
		"regular user, referral": {
			expectedTakerFeePpm: 30,
			setup: func(ctx sdk.Context, affiliatesKeeper *affiliateskeeper.Keeper) {
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

			statsKeeper := tApp.App.StatsKeeper
			statsKeeper.SetUserStats(ctx, constants.AliceAccAddress.String(), &stattypes.UserStats{
				TakerNotional: 10,
				MakerNotional: 10,
			})
			statsKeeper.SetGlobalStats(ctx, &stattypes.GlobalStats{
				NotionalTraded: 10_000,
			})

			affiliatesKeeper := tApp.App.AffiliatesKeeper
			if tc.setup != nil {
				tc.setup(ctx, &affiliatesKeeper)
			}

			require.Equal(t, tc.expectedTakerFeePpm,
				k.GetPerpetualFeePpm(ctx, constants.AliceAccAddress.String(), true))
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
