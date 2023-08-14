package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()

	logger := tApp.App.FeeTiersKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

func TestGetPerpetualFeePpm(t *testing.T) {
	tests := map[string]struct {
		UserStats           *stattypes.UserStats
		GlobalStats         *stattypes.GlobalStats
		expectedTakerFeePpm int32
		expectedMakerFeePpm int32
	}{
		"first tier": {
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
		"increased tier": {
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
		"partial requirements doesn't increase tier": {
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
		"top tier": {
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
			ctx := tApp.InitChain()
			user := "alice"
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
			statsKeeper.SetUserStats(ctx, user, tc.UserStats)
			statsKeeper.SetGlobalStats(ctx, tc.GlobalStats)

			require.Equal(t, tc.expectedTakerFeePpm, k.GetPerpetualFeePpm(ctx, user, true))
			require.Equal(t, tc.expectedMakerFeePpm, k.GetPerpetualFeePpm(ctx, user, false))
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
			tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
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
