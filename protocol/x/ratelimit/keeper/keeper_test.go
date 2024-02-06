package keeper_test

import (
	"math/big"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

const (
	testDenom    = "ibc/xxx"
	testDenom2   = "testdenom2"
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
					Period:          3_600 * time.Second,
					BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
					BaselineTvlPpm:  10_000,                         // 1%
				},
				{
					Period:          86_400 * time.Second,
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
					Period:          3_600 * time.Second,
					BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
					BaselineTvlPpm:  10_000,                         // 1%
				},
				{
					Period:          86_400 * time.Second,
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
					Period:          3_600 * time.Second,
					BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
					BaselineTvlPpm:  50_000,                         // 5%
				},
				{
					Period:          86_400 * time.Second,
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
		"50m TVL, capactiy correctly initialized to 10% and 99%": {
			denom: testDenom,
			limiters: []types.Limiter{
				{
					Period:          3_600 * time.Second,
					BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
					BaselineTvlPpm:  100_000,                        // 10%
				},
				{
					Period:          86_400 * time.Second,
					BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1m tokens (assuming 6 decimals)
					BaselineTvlPpm:  990_000,                          // 99%
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
				dtypes.NewInt(49_500_000_000_000), // 49.5m tokens (99% of 50m)
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
			err := k.SetLimitParams(ctx, limitParams)
			require.NoError(t, err)

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
			err = k.SetLimitParams(ctx, types.LimitParams{
				Denom:    tc.denom,
				Limiters: []types.Limiter{}, // Empty list, results in deletion of the key.
			})
			require.NoError(t, err)

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

func TestUpdateAllCapacitiesEndBlocker(t *testing.T) {
	tests := map[string]struct {
		balances                  []banktypes.Balance // For initializing the current supply
		limitParamsList           []types.LimitParams
		prevBlockTime             time.Time
		blockTime                 time.Time
		initDenomCapacityList     []types.DenomCapacity
		expectedDenomCapacityList []types.DenomCapacity
	}{
		"One denom, prev capacity equals baseline": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = 25M * 1% = 250k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
						// baseline = 25M * 10% = 2.5M tokens
						{
							Period:          86_400 * time.Second,
							BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
							BaselineTvlPpm:  100_000,                          // 10%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1001, 0).In(time.UTC),
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(250_000_000_000),   // 250k tokens, which equals baseline
						dtypes.NewInt(2_500_000_000_000), // 2.5M tokens, which equals baseline
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(250_000_000_000),   // 250k tokens
						dtypes.NewInt(2_500_000_000_000), // 2.5M tokens
					},
				},
			},
		},
		"One denom, prev capacity < baseline": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = 25M * 1% = 250k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
						// baseline = 25M * 10% = 2.5M tokens
						{
							Period:          86_400 * time.Second,
							BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
							BaselineTvlPpm:  100_000,                          // 10%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1001, 0).In(time.UTC),
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(99_000_000_000),  // 99k tokens, < baseline (250k)
						dtypes.NewInt(990_000_000_000), // 0.99M tokens, < baseline (2.5M)
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(99_069_444_444),  // recovered by 1/3600 * 250k = 69.4444 tokens
						dtypes.NewInt(990_028_935_185), // recovered by 1/86400 * 2.5M = 28.9351 tokens
					},
				},
			},
		},
		"One denom, prev capacity < baseline, 18 decimals": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom: testDenom,
							Amount: sdkmath.NewIntFromBigInt(
								big_testutil.Int64MulPow10(25, 24), // 25M tokens (assuming 18 decimals)
							),
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = 25M * 1% = 250k tokens
						{
							Period: 3_600 * time.Second,
							BaselineMinimum: dtypes.NewIntFromBigInt(
								big_testutil.Int64MulPow10(100_000, 18), // 100k tokens(assuming 18 decimals)
							),
							BaselineTvlPpm: 10_000, // 1%
						},
						// baseline = 25M * 10% = 2.5M tokens
						{
							Period: 86_400 * time.Second,
							BaselineMinimum: dtypes.NewIntFromBigInt(
								big_testutil.Int64MulPow10(1_000_000, 18), // 1M tokens(assuming 18 decimals)
							),
							BaselineTvlPpm: 100_000, // 10%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1001, 0).In(time.UTC),
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewIntFromBigInt(
							big_testutil.Int64MulPow10(99_000, 18),
						), // 99k tokens < baseline (250k)
						dtypes.NewIntFromBigInt(
							big_testutil.Int64MulPow10(990_000, 18),
						), // 0.99M tokens, < baseline (2.5M)
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("99069444444444444444444", 10)),
						), // recovered by 1/3600 * 250k ~= 69.4444 tokens
						dtypes.NewIntFromBigInt(
							big_testutil.MustFirst(new(big.Int).SetString("990028935185185185185185", 10)),
						), // recovered by 1/86400 * 2.5M ~= 28.9351 tokens
					},
				},
			},
		},
		"One denom, prev capacity = 0": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(1_000_000_000_000), // 1M token (assuming 6 decimals)
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = baseline minimum = 100k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
						// baseline = baseline minimum = 1M tokens
						{
							Period:          86_400 * time.Second,
							BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
							BaselineTvlPpm:  100_000,                          // 10%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1001, 0).In(time.UTC),
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(0), // 0 Capacity
						dtypes.NewInt(0), // 0 Capacity
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(27_777_777), // recovered by 1/3600 * 100k ~= 27.7778 tokens
						dtypes.NewInt(11_574_074), // recovered by 1/86400 * 1M ~= 11.5741 tokens
					},
				},
			},
		},
		"One denom, baseline < prev capacity < 2 * baseline": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(20_000_000_000_000), // 20M token (assuming 6 decimals)
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = 20M * 1% = 200k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
						// baseline = 20M * 10% = 2M tokens
						{
							Period:          86_400 * time.Second,
							BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
							BaselineTvlPpm:  100_000,                          // 10%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1001, 0).In(time.UTC),
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(329_000_000_000),
						dtypes.NewInt(3_500_000_000_000),
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(328_944_444_445),   // recovered by 1/3600 * 200k ~= 55.5555555556
						dtypes.NewInt(3_499_976_851_852), // recovered by 1/86400 * 2M ~= 23.1481481482
					},
				},
			},
		},
		"One denom, prev capacity > 2 * baseline": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(20_000_000_000_000), // 20M token (assuming 6 decimals)
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = 20M * 1% = 200k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
						// baseline = 20M * 10% = 2M tokens
						{
							Period:          86_400 * time.Second,
							BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
							BaselineTvlPpm:  100_000,                          // 10%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1001, 0).In(time.UTC),
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(529_000_000_000),   // 529k tokens > 2 * baseline (200k)
						dtypes.NewInt(4_500_000_000_000), // 4.5M tokens > 2 * baseline (2)
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(528_908_611_112),   // recovered by 1/3600 * (529k - 200k) ~= 91.389
						dtypes.NewInt(4_499_971_064_815), // recovered by 1/86400 * (4.5M - 2M) ~= 28.935
					},
				},
			},
		},
		"Two denoms, mix of values from above cases": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(20_000_000_000_000), // 20M token (assuming 6 decimals)
						},
						{
							Denom:  testDenom2,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 20M token (assuming 6 decimals)
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = 20M * 1% = 200k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
						// baseline = 20M * 10% = 2M tokens
						{
							Period:          86_400 * time.Second,
							BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
							BaselineTvlPpm:  100_000,                          // 10%
						},
					},
				},
				{
					Denom: testDenom2,
					Limiters: []types.Limiter{
						// baseline = 25M * 1% = 250k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
						// baseline = 25M * 10% = 2.5M tokens
						{
							Period:          86_400 * time.Second,
							BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
							BaselineTvlPpm:  100_000,                          // 10%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1001, 0).In(time.UTC),
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(529_000_000_000),   // 529k tokens > 2 * baseline (200k)
						dtypes.NewInt(4_500_000_000_000), // 4.5M tokens > 2 * baseline (2)
					},
				},
				{
					Denom: testDenom2,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(99_000_000_000),  // 99k tokens, < baseline (250k)
						dtypes.NewInt(990_000_000_000), // 0.99M tokens, < baseline (2.5M)
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(528_908_611_112),   // recovered by 1/3600 * (529k - 200k) ~= 91.389
						dtypes.NewInt(4_499_971_064_815), // recovered by 1/86400 * (4.5M - 2M) ~= 28.935
					},
				},
				{
					Denom: testDenom2,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(99_069_444_444),  // recovered by 1/3600 * 250k = 69.4444 tokens
						dtypes.NewInt(990_028_935_185), // recovered by 1/86400 * 2.5M = 28.9351 tokens
					},
				},
			},
		},
		"(Error) one denom, current block time = prev block time, no changes applied": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = 25M * 1% = 250k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
						// baseline = 25M * 10% = 2.5M tokens
						{
							Period:          86_400 * time.Second,
							BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
							BaselineTvlPpm:  100_000,                          // 10%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1000, 0).In(time.UTC), // same as prev block time (should not happen in practice)
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(99_000_000_000),  // 99k tokens, < baseline (250k)
						dtypes.NewInt(990_000_000_000), // 0.99M tokens, < baseline (2.5M)
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(99_000_000_000),  // 99k tokens (unchanged)
						dtypes.NewInt(990_000_000_000), // 0.99M tokens (unchanged)
					},
				},
			},
		},
		"(Error) one denom, len(limiters) != len(capacityList)": {
			balances: []banktypes.Balance{
				{
					Address: testAddress1,
					Coins: sdk.Coins{
						{
							Denom:  testDenom,
							Amount: sdkmath.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
						},
					},
				},
			},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						// baseline = 25M * 1% = 250k tokens
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
							BaselineTvlPpm:  10_000,                         // 1%
						},
					},
				},
			},
			prevBlockTime: time.Unix(1000, 0).In(time.UTC),
			blockTime:     time.Unix(1001, 0).In(time.UTC), // same as prev block time (should not happen in practice)
			initDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(99_000_000_000),  // 99k tokens, < baseline (250k)
						dtypes.NewInt(990_000_000_000), // 0.99M tokens, < baseline (2.5M)
					},
				},
			},
			expectedDenomCapacityList: []types.DenomCapacity{
				{
					Denom: testDenom,
					CapacityList: []dtypes.SerializableInt{
						dtypes.NewInt(99_000_000_000),  // 99k tokens (unchanged)
						dtypes.NewInt(990_000_000_000), // 0.99M tokens (unchanged)
					},
				},
			},
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

			// Set previous block time
			tApp.App.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Timestamp: tc.prevBlockTime,
			})

			k := tApp.App.RatelimitKeeper

			// Initialize limit params
			for _, limitParams := range tc.limitParamsList {
				err := k.SetLimitParams(ctx, limitParams)
				require.NoError(t, err)
			}

			// Initialize denom capacity
			for _, denomCapacity := range tc.initDenomCapacityList {
				k.SetDenomCapacity(ctx, denomCapacity)
			}

			// Run the function being tested
			k.UpdateAllCapacitiesEndBlocker(ctx.WithBlockTime(tc.blockTime))

			// Check results
			for _, expectedDenomCapacity := range tc.expectedDenomCapacityList {
				gotDenomCapacity := k.GetDenomCapacity(ctx, expectedDenomCapacity.Denom)
				require.Equal(t,
					expectedDenomCapacity,
					gotDenomCapacity,
					"expected denom capacity: %+v, got: %+v",
					expectedDenomCapacity,
					gotDenomCapacity,
				)
			}
		})
	}
}

func TestGetAllLimitParams(t *testing.T) {
	denom4LimitParams := types.LimitParams{
		Denom: "denom4",
		Limiters: []types.Limiter{
			{
				Period:          3_600 * time.Second,
				BaselineMinimum: dtypes.NewInt(100_000_000_000),
				BaselineTvlPpm:  10_000,
			},
		},
	}
	denom3LimitParams := types.LimitParams{
		Denom: "denom3",
		Limiters: []types.Limiter{
			{
				Period:          24 * time.Hour,
				BaselineMinimum: dtypes.NewInt(1_000_000_000_000),
				BaselineTvlPpm:  100_000,
			},
		},
	}
	denom1LimitParams := types.LimitParams{
		Denom: "denom1",
		Limiters: []types.Limiter{
			{
				Period:          12 * time.Hour,
				BaselineMinimum: dtypes.NewInt(123_456_789_000),
				BaselineTvlPpm:  10_000,
			},
		},
	}
	denom2LimitParams := types.LimitParams{
		Denom: "denom2",
		Limiters: []types.Limiter{
			{
				Period:          72_000 * time.Second,
				BaselineMinimum: dtypes.NewInt(100_000_000_000),
				BaselineTvlPpm:  10_000,
			},
		},
	}
	testLimitParamsList := []types.LimitParams{
		denom4LimitParams,
		denom3LimitParams,
		denom1LimitParams,
		denom2LimitParams,
	}

	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *types.GenesisState) {
				// Empty genesis params for test
				genesisState.LimitParamsList = []types.LimitParams{}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	// Initialize limit params
	for _, limitParams := range testLimitParamsList {
		err := k.SetLimitParams(ctx, limitParams)
		require.NoError(t, err)
	}

	allLimitParams := k.GetAllLimitParams(ctx)

	expectedLimitParamsList := []types.LimitParams{
		denom1LimitParams,
		denom2LimitParams,
		denom3LimitParams,
		denom4LimitParams,
	}

	require.Equal(t,
		expectedLimitParamsList,
		allLimitParams,
	)
}

func TestProcessWithdrawal(t *testing.T) {
	tests := map[string]struct {
		balances              []banktypes.Balance // Use baseline minimum
		limitParamsList       []types.LimitParams
		withdrawDenom         string
		withdrawAmount        *big.Int
		expectedDenomCapacity types.DenomCapacity
		expectedErr           error
	}{
		"has limit params, withdrawal amount < capacity, succeeds": {
			balances: []banktypes.Balance{},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000),
							BaselineTvlPpm:  10_000,
						},
					},
				},
			},
			withdrawDenom:  testDenom,
			withdrawAmount: big.NewInt(98_760_000), // < baseline capacity
			expectedDenomCapacity: types.DenomCapacity{
				Denom: testDenom,
				CapacityList: []dtypes.SerializableInt{
					dtypes.NewInt(1_240_000), // 100_000_000 - 98_760_000
				},
			},
			expectedErr: nil,
		},
		"no limit params, succeeds": {
			balances: []banktypes.Balance{},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000),
							BaselineTvlPpm:  10_000,
						},
					},
				},
			},
			withdrawDenom:  testDenom2,
			withdrawAmount: big.NewInt(98_760_000),
			expectedDenomCapacity: types.DenomCapacity{
				Denom: testDenom,
				CapacityList: []dtypes.SerializableInt{
					dtypes.NewInt(100_000_000), // unchanged
				},
			},
			expectedErr: nil,
		},
		"has limit params, withdrawal amount > capacity, rate limited": {
			balances: []banktypes.Balance{},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000),
							BaselineTvlPpm:  10_000,
						},
					},
				},
			},
			withdrawDenom:  testDenom,
			withdrawAmount: big.NewInt(105_000_000), // < baseline capacity
			expectedDenomCapacity: types.DenomCapacity{
				Denom: testDenom,
				CapacityList: []dtypes.SerializableInt{
					dtypes.NewInt(100_000_000), // unchanged
				},
			},
			expectedErr: errorsmod.Wrapf(
				types.ErrWithdrawalExceedsCapacity,
				"denom = %v, capacity(index: %v) = %v, amount = %v",
				testDenom,
				0,
				big.NewInt(100_000_000),
				big.NewInt(105_000_000),
			),
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

			// Initialize limit params
			for _, limitParams := range tc.limitParamsList {
				err := k.SetLimitParams(ctx, limitParams)
				require.NoError(t, err)
			}

			// Run the function being tested
			err := k.ProcessWithdrawal(
				ctx,
				tc.withdrawDenom,
				tc.withdrawAmount,
			)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			gotDenomCapacity := k.GetDenomCapacity(ctx, tc.expectedDenomCapacity.Denom)
			require.Equal(
				t,
				tc.expectedDenomCapacity,
				gotDenomCapacity,
			)
		})
	}
}

func TestIncrementCapacitiesForDenom(t *testing.T) {
	tests := map[string]struct {
		balances              []banktypes.Balance // Use baseline minimum
		limitParamsList       []types.LimitParams
		incrementDenom        string
		incrementAmount       *big.Int
		expectedDenomCapacity types.DenomCapacity
		expectedErr           error
	}{
		"has limit params": {
			balances: []banktypes.Balance{},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000),
							BaselineTvlPpm:  10_000,
						},
					},
				},
			},
			incrementDenom:  testDenom,
			incrementAmount: big.NewInt(98_760_000), // < baseline capacity
			expectedDenomCapacity: types.DenomCapacity{
				Denom: testDenom,
				CapacityList: []dtypes.SerializableInt{
					dtypes.NewInt(198_760_000), // 100_000_000 + 98_760_000
				},
			},
			expectedErr: nil,
		},
		"no limit params": {
			balances: []banktypes.Balance{},
			limitParamsList: []types.LimitParams{
				{
					Denom: testDenom,
					Limiters: []types.Limiter{
						{
							Period:          3_600 * time.Second,
							BaselineMinimum: dtypes.NewInt(100_000_000),
							BaselineTvlPpm:  10_000,
						},
					},
				},
			},
			incrementDenom:  testDenom2,
			incrementAmount: big.NewInt(98_760_000),
			expectedDenomCapacity: types.DenomCapacity{
				Denom: testDenom,
				CapacityList: []dtypes.SerializableInt{
					dtypes.NewInt(100_000_000), // unchanged
				},
			},
			expectedErr: nil,
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

			// Initialize limit params
			for _, limitParams := range tc.limitParamsList {
				err := k.SetLimitParams(ctx, limitParams)
				require.NoError(t, err)
			}

			// Run the function being tested
			k.IncrementCapacitiesForDenom(
				ctx,
				tc.incrementDenom,
				tc.incrementAmount,
			)

			gotDenomCapacity := k.GetDenomCapacity(ctx, tc.expectedDenomCapacity.Denom)
			require.Equal(
				t,
				tc.expectedDenomCapacity,
				gotDenomCapacity,
			)
		})
	}
}
