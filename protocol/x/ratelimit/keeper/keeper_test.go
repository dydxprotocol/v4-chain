package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"math/big"
	"testing"

	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

const (
	testDenom    = "ibc/xxx"
	testAddress1 = "dydx16h7p7f4dysrgtzptxx2gtpt5d8t834g9dj830z"
	testAddress2 = "dydx168pjt8rkru35239fsqvz7rzgeclakp49zx3aum"
	testAddress3 = "dydx1fjg6zp6vv8t9wvy4lps03r5l4g7tkjw9wvmh70"
)

func TestSetGetDenomCapacity(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	capacityList := []dtypes.SerializableInt{
		dtypes.NewInt(123_456_789),
		dtypes.NewInt(500_000_000),
	}
	denomCapacity := types.DenomCapacity{
		Denom:        testDenom,
		CapacityList: capacityList,
	}

	// Test SetDenomCapacity
	k.SetDenomCapacity(ctx, denomCapacity)

	// Test GetDenomCapacity
	gotDenomCapacity := k.GetDenomCapacity(ctx, testDenom)
	require.Equal(t, denomCapacity, gotDenomCapacity, "retrieved DenomCapacity does not match the set value")

	k.SetDenomCapacity(ctx, types.DenomCapacity{
		Denom:        testDenom,
		CapacityList: []dtypes.SerializableInt{}, // Empty list, results in deletion of the key.
	})

	// Check that the key is deleted under `DenomCapacity` storage.
	require.Equal(t,
		types.DenomCapacity{
			Denom:        testDenom,
			CapacityList: nil,
		},
		k.GetDenomCapacity(ctx, testDenom),
		"retrieved LimitParams do not match the set value",
	)
}

func TestSetGetLimitParams_Success(t *testing.T) {
	// Run tests.
	tests := map[string]struct {
		denom                string
		balances             []banktypes.Balance
		limiters             []types.Limiter
		expectedCapacityList []dtypes.SerializableInt
	}{
		"0 TVL, capactiy correctly initialized as minimum baseline": {
			denom: testDenom,
			limiters: []types.Limiter{
				{
					PeriodSec:       3_600,
					BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
					BaselineTvlPpm:  10_000,                         // 1%
				},
				{
					PeriodSec:       86_400,
					BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1m tokens (assuming 6 decimals)
					BaselineTvlPpm:  100_000,                          // 10%
				},
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(100_000_000_000),
				dtypes.NewInt(1_000_000_000_000),
			},
		},
		"50m TVL, capactiy correctly initialized to 1% and 10%": {
			denom: testDenom,
			limiters: []types.Limiter{
				{
					PeriodSec:       3_600,
					BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
					BaselineTvlPpm:  10_000,                         // 1%
				},
				{
					PeriodSec:       86_400,
					BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1m tokens (assuming 6 decimals)
					BaselineTvlPpm:  100_000,                          // 10%
				},
			},
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token
						},
					},
				},
				{
					Address: testAddress2,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token
						},
					},
				},
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(500_000_000_000),   // 500k tokens (1% of 50m)
				dtypes.NewInt(5_000_000_000_000), // 5m tokens (10% of 50m)
			},
		},
		"50m TVL, capactiy correctly initialized to 5% and 20% (rounds down)": {
			denom: testDenom,
			limiters: []types.Limiter{
				{
					PeriodSec:       3_600,
					BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
					BaselineTvlPpm:  50_000,                         // 5%
				},
				{
					PeriodSec:       86_400,
					BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1m tokens (assuming 6 decimals)
					BaselineTvlPpm:  200_000,                          // 20%
				},
			},
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_123_450_000_000), // 25M token
						},
					},
				},
				{
					Address: testAddress2,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token
						},
					},
				},
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(2_506_172_500_000),  // ~2.5M tokens (5% of ~50m)
				dtypes.NewInt(10_024_690_000_000), // ~5m tokens (20% of 50m)
			},
		},
		"50m TVL, capactiy correctly initialized to 10% and 100%": {
			denom: testDenom,
			limiters: []types.Limiter{
				{
					PeriodSec:       3_600,
					BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
					BaselineTvlPpm:  100_000,                        // 10%
				},
				{
					PeriodSec:       86_400,
					BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1m tokens (assuming 6 decimals)
					BaselineTvlPpm:  1_000_000,                        // 100%
				},
			},
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token
						},
					},
				},
				{
					Address: testAddress2,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token
						},
					},
				},
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(5_000_000_000_000),  // 5m tokens (10% of 50m)
				dtypes.NewInt(50_000_000_000_000), // 50m tokens (100% of 50m)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set up treasury account balance in genesis state
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						genesisState.Balances = append(genesisState.Balances, tc.balances...)
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			limitParams := types.LimitParams{
				Denom:    tc.denom,
				Limiters: tc.limiters,
			}

			// Test SetLimitParams
			k.SetLimitParams(ctx, limitParams)

			// Test GetLimitParams
			gotLimitParams := k.GetLimitParams(ctx, tc.denom)
			require.Equal(t, limitParams, gotLimitParams, "retrieved LimitParams do not match the set value")

			// Query for `DenomCapacity` of `testDenom`.
			gotDenomCapacity := k.GetDenomCapacity(ctx, tc.denom)
			// Expected `DenomCapacity` is initialized such that each capacity is equal to the baseline.
			expectedDenomCapacity := types.DenomCapacity{
				Denom:        tc.denom,
				CapacityList: tc.expectedCapacityList,
			}
			require.Equal(t, expectedDenomCapacity, gotDenomCapacity, "retrieved DenomCapacity does not match the set value")

			// Set empty `LimitParams` for `testDenom`.
			k.SetLimitParams(ctx, types.LimitParams{
				Denom:    tc.denom,
				Limiters: []types.Limiter{}, // Empty list, results in deletion of the key.
			})

			// Check that the key is deleted under `LimitParams` storage.
			require.Equal(t,
				types.LimitParams{
					Denom:    tc.denom,
					Limiters: nil,
				},
				k.GetLimitParams(ctx, tc.denom),
				"retrieved LimitParams do not match the set value")

			// Check that the key is deleted under `DenomCapacity` storage.
			require.Equal(t,
				types.DenomCapacity{
					Denom:        tc.denom,
					CapacityList: nil,
				},
				k.GetDenomCapacity(ctx, tc.denom),
				"retrieved LimitParams do not match the set value")
		})
	}
}

func TestGetBaseline(t *testing.T) {
	tests := map[string]struct {
		denom            string
		balances         []banktypes.Balance
		limiter          types.Limiter
		expectedBaseline *big.Int
	}{
		"max(1% of TVL, 100k token), TVL = 5M token": {
			denom: testDenom,
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(1_000_000_000_000), // 1M token
						},
					},
				},
				{
					Address: testAddress2,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(4_000_000_000_000), // 4M token
						},
					},
				},
			},
			limiter: types.Limiter{
				PeriodSec:       3_600,
				BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k token
				BaselineTvlPpm:  10_000,                         // 1%
			},
			expectedBaseline: big.NewInt(100_000_000_000), // 100k token (baseline minimum)
		},
		"max(1% of TVL, 100k token), TVL = 15M token": {
			denom: testDenom,
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(1_000_000_000_000), // 1M token
						},
					},
				},
				{
					Address: testAddress2,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(4_000_000_000_000), // 4M token
						},
					},
				},
				{
					Address: testAddress3,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(10_000_000_000_000), // 10M token
						},
					},
				},
			},
			limiter: types.Limiter{
				PeriodSec:       3_600,
				BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k token
				BaselineTvlPpm:  10_000,                         // 1%
			},
			expectedBaseline: big.NewInt(150_000_000_000), // 150k token (1% of 15m)
		},
		"max(1% of TVL, 100k token), TVL = ~15M token, rounds down": {
			denom: testDenom,
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(1_000_000_000_000), // 1M token
						},
					},
				},
				{
					Address: testAddress2,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(4_000_123_456_777), // ~4M token
						},
					},
				},
				{
					Address: testAddress3,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(10_200_000_000_000), // ~10M token
						},
					},
				},
			},
			limiter: types.Limiter{
				PeriodSec:       3_600,
				BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k token
				BaselineTvlPpm:  10_000,                         // 1%
			},
			expectedBaseline: big.NewInt(152_001_234_567), // ~152k token (1% of 15.2m)
		},
		"max(10% of TVL, 1 million), TVL = 20M token": {
			denom: testDenom,
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(6_000_000_000_000), // 6M token
						},
					},
				},
				{
					Address: testAddress2,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(4_000_000_000_000), // 4M token
						},
					},
				},
				{
					Address: testAddress3,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(10_000_000_000_000), // 10M token
						},
					},
				},
			},
			limiter: types.Limiter{
				PeriodSec:       3_600,
				BaselineMinimum: dtypes.NewInt(100_000_000_000), // 1m token
				BaselineTvlPpm:  100_000,                        // 10%
			},
			expectedBaseline: big.NewInt(2_000_000_000_000), // 2m token (10% of 20m)
		},
		"max(10% of TVL, 1 million), TVL = 8M token": {
			denom: testDenom,
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(2_000_000_000_000), // 2M token
						},
					},
				},
				{
					Address: testAddress2,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(4_000_000_000_000), // 4M token
						},
					},
				},
				{
					Address: testAddress3,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(2_000_000_000_000), // 2M token
						},
					},
				},
			},
			limiter: types.Limiter{
				PeriodSec:       3_600,
				BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1m token
				BaselineTvlPpm:  100_000,                          // 10%
			},
			expectedBaseline: big.NewInt(1_000_000_000_000), // 1m token (baseline minimum)
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set up treasury account balance in genesis state
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						genesisState.Balances = append(genesisState.Balances, tc.balances...)
					},
				)
				return genesis
			}).Build()

			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			gotBaseline := k.GetBaseline(ctx, tc.denom, tc.limiter)

			require.Equal(t, tc.expectedBaseline, gotBaseline, "retrieved baseline does not match the expected value")
		})
	}
}
