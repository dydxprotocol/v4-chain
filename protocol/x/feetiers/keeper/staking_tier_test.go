package keeper_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetStakingTiers(t *testing.T) {
	tests := map[string]struct {
		// Input
		initialTiers []*types.StakingTier
		newTiers     []*types.StakingTier
	}{
		"Set tiers twice": {
			initialTiers: []*types.StakingTier{
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
							MinStakedBaseTokens: dtypes.NewInt(200),
							FeeDiscountPpm:      20000,
						},
					},
				},
			},
			newTiers: []*types.StakingTier{
				{
					FeeTierName: "3",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(500),
							FeeDiscountPpm:      50000,
						},
					},
				},
			},
		},
		"Set tiers and then set to empty": {
			initialTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      10000,
						},
					},
				},
			},
			newTiers: []*types.StakingTier{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper

			// Set initial tiers
			err := k.SetStakingTiers(ctx, tc.initialTiers)
			require.NoError(t, err)

			// Verify all tiers
			allTiers := k.GetAllStakingTiers(ctx)
			require.ElementsMatch(t, tc.initialTiers, allTiers)

			// Set new tiers
			err = k.SetStakingTiers(ctx, tc.newTiers)
			require.NoError(t, err)

			// Verify all tiers
			allTiers = k.GetAllStakingTiers(ctx)
			require.ElementsMatch(t, tc.newTiers, allTiers)

			// Verify each tier
			for _, expectedTier := range tc.newTiers {
				tier, found := k.GetStakingTier(ctx, expectedTier.FeeTierName)
				require.True(t, found)
				require.Equal(t, expectedTier, tier)
			}

			// Verify old tiers (that don't exist in new tiers) no longer exist
			for _, oldTier := range tc.initialTiers {
				existsInNew := false
				for _, newTier := range tc.newTiers {
					if newTier.FeeTierName == oldTier.FeeTierName {
						existsInNew = true
						break
					}
				}

				if !existsInNew {
					_, found := k.GetStakingTier(ctx, oldTier.FeeTierName)
					require.False(t, found, "old tier %s should not be accessible", oldTier.FeeTierName)
				}
			}
		})
	}
}

func TestSetStakingTiers_ValidationError(t *testing.T) {
	tests := map[string]struct {
		// Input
		stakingTiers []*types.StakingTier

		// Expected
		expectedError string
	}{
		"empty fee tier name": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      10000,
						},
					},
				},
			},
			expectedError: "fee tier name cannot be empty",
		},
		"duplicate fee tier names": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels:      []*types.StakingLevel{},
				},
				{
					FeeTierName: "1",
					Levels:      []*types.StakingLevel{},
				},
			},
			expectedError: "duplicate staking tier for fee tier: 1",
		},
		"negative min staked tokens": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(-100),
							FeeDiscountPpm:      10000,
						},
					},
				},
			},
			expectedError: "min staked tokens cannot be negative for tier 1 level 0",
		},
		"levels in decreasing order": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      10000,
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(999),
							FeeDiscountPpm:      20000,
						},
					},
				},
			},
			expectedError: "staking levels must be in increasing order for tier 1",
		},
		"levels with equal amounts": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      10000,
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      20000,
						},
					},
				},
			},
			expectedError: "staking levels must be in increasing order for tier 1",
		},
		"discount exceeds 100%": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      1_000_001,
						},
					},
				},
			},
			expectedError: "fee discount cannot exceed 100% for tier 1 level 0",
		},
		"fails with non-existent fee tier": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "7777",
					Levels:      []*types.StakingLevel{},
				},
			},
			expectedError: "fee tier 7777 does not exist",
		},
		"fails when one of multiple tiers doesn't exist": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels:      []*types.StakingLevel{},
				},
				{
					FeeTierName: "7777",
					Levels:      []*types.StakingLevel{},
				},
			},
			expectedError: "fee tier 7777 does not exist",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper

			err := k.SetStakingTiers(ctx, tc.stakingTiers)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

func TestGetStakingTier_NotFound(t *testing.T) {
	tests := map[string]struct {
		// Setup
		initialTiers []*types.StakingTier

		// Input
		queryTier string

		// Expected
		expectedFound bool
	}{
		"returns false for non-existent tier": {
			initialTiers:  []*types.StakingTier{},
			queryTier:     "1",
			expectedFound: false,
		},
		"returns false for non-existent tier when store has data": {
			initialTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels:      []*types.StakingLevel{},
				},
			},
			queryTier:     "999",
			expectedFound: false,
		},
		"returns false for empty tier name": {
			initialTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels:      []*types.StakingLevel{},
				},
			},
			queryTier:     "",
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.FeeTiersKeeper

			// Setup initial tiers
			err := k.SetStakingTiers(ctx, tc.initialTiers)
			require.NoError(t, err)

			// Query tier
			tier, found := k.GetStakingTier(ctx, tc.queryTier)
			require.Equal(t, tc.expectedFound, found)
			if !tc.expectedFound {
				require.Nil(t, tier)
			}
		})
	}
}

func TestGetStakingDiscountPpm(t *testing.T) {
	tests := map[string]struct {
		// Setup
		stakingTiers []*types.StakingTier

		// Input
		feeTierName      string
		stakedBaseTokens *big.Int

		// Expected
		expectedDiscountPpm uint32
	}{
		"No staking tier configured": {
			stakingTiers:        []*types.StakingTier{},
			feeTierName:         "1",
			stakedBaseTokens:    big.NewInt(1000),
			expectedDiscountPpm: 0,
		},
		"Staking tier exists but no levels": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels:      []*types.StakingLevel{},
				},
			},
			feeTierName:         "1",
			stakedBaseTokens:    big.NewInt(1000),
			expectedDiscountPpm: 0,
		},
		"User has zero staked tokens": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      50000, // 5%
						},
					},
				},
			},
			feeTierName:         "1",
			stakedBaseTokens:    big.NewInt(0),
			expectedDiscountPpm: 0,
		},
		"User qualifies for first level": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "2",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      50000, // 5%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      100000, // 10%
						},
					},
				},
			},
			feeTierName:         "2",
			stakedBaseTokens:    big.NewInt(999),
			expectedDiscountPpm: 50000,
		},
		"User qualifies for middle level": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "7",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      50000, // 5%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(500),
							FeeDiscountPpm:      75000, // 7.5%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      100000, // 10%
						},
					},
				},
			},
			feeTierName:         "7",
			stakedBaseTokens:    big.NewInt(500),
			expectedDiscountPpm: 75000,
		},
		"Two staking tiers": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "3",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      50000, // 5%
						},
					},
				},
				{
					FeeTierName: "9",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(200),
							FeeDiscountPpm:      75000, // 7.5%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(500),
							FeeDiscountPpm:      100000, // 10%
						},
					},
				},
			},
			feeTierName:         "9",
			stakedBaseTokens:    big.NewInt(499),
			expectedDiscountPpm: 75000,
		},
		"Maximum discount - 100%": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "1",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      50000, // 5%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(10000),
							FeeDiscountPpm:      1000000, // 100%
						},
					},
				},
			},
			feeTierName:         "1",
			stakedBaseTokens:    big.NewInt(10000),
			expectedDiscountPpm: 1000000,
		},
		"A tier with five levels": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "6",
					Levels: []*types.StakingLevel{
						{
							MinStakedBaseTokens: dtypes.NewInt(100),
							FeeDiscountPpm:      20000, // 2%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(500),
							FeeDiscountPpm:      40000, // 4%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(1000),
							FeeDiscountPpm:      60000, // 6%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(5000),
							FeeDiscountPpm:      80000, // 8%
						},
						{
							MinStakedBaseTokens: dtypes.NewInt(10000),
							FeeDiscountPpm:      100000, // 10%
						},
					},
				},
			},
			feeTierName:         "6",
			stakedBaseTokens:    big.NewInt(7500),
			expectedDiscountPpm: 80000,
		},
		"Large staked amount": {
			stakingTiers: []*types.StakingTier{
				{
					FeeTierName: "5",
					Levels: []*types.StakingLevel{
						{
							// 1e24
							MinStakedBaseTokens: dtypes.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(24), nil)),
							FeeDiscountPpm:      50000, // 5%
						},
						{
							// 1e28
							MinStakedBaseTokens: dtypes.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(28), nil)),
							FeeDiscountPpm:      110000, // 11%
						},
						{
							// 1e33
							MinStakedBaseTokens: dtypes.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(33), nil)),
							FeeDiscountPpm:      200000, // 20%
						},
					},
				},
			},
			feeTierName:         "5",
			stakedBaseTokens:    new(big.Int).Exp(big.NewInt(10), big.NewInt(30), nil), // 1e30
			expectedDiscountPpm: 110000,
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

			// Verify discount
			discountPpm := k.GetStakingDiscountPpm(ctx, tc.feeTierName, tc.stakedBaseTokens)
			require.Equal(t, tc.expectedDiscountPpm, discountPpm)
		})
	}
}
