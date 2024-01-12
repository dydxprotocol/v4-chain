package util_test

import (
	"math/big"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/util"
	"github.com/stretchr/testify/require"
)

func TestCalculateNewCapacityList(t *testing.T) {
	testDenom := "testDenom"
	tests := map[string]struct {
		bigTvl               *big.Int
		limitParams          types.LimitParams
		prevCapacityList     []dtypes.SerializableInt
		expectedCapacityList []dtypes.SerializableInt
		timeSinceLastBlock   time.Duration
		expectedErr          error
	}{
		"Prev capacity equals baseline": {
			bigTvl: big.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
			limitParams: types.LimitParams{
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
			timeSinceLastBlock: time.Second,
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(250_000_000_000),   // 250k tokens, which equals baseline
				dtypes.NewInt(2_500_000_000_000), // 2.5M tokens, which equals baseline
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(250_000_000_000),   // 250k tokens
				dtypes.NewInt(2_500_000_000_000), // 2.5M tokens
			},
		},
		"Prev capacity < baseline": {
			bigTvl: big.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
			limitParams: types.LimitParams{
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
			timeSinceLastBlock: time.Second + 90*time.Millisecond, // 1.09 second
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(99_000_000_000),  // 99k tokens, < baseline (250k)
				dtypes.NewInt(990_000_000_000), // 0.99M tokens, < baseline (2.5M)
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(99_075_694_444),  // recovered by 1.09/3600 * 250k = 75.694444 tokens
				dtypes.NewInt(990_031_539_351), // recovered by 1.09/86400 * 2.5M = 31.539 tokens
			},
		},
		"prev capacity < baseline, 18 decimals": {
			bigTvl: big_testutil.Int64MulPow10(25, 24), // 25M tokens
			limitParams: types.LimitParams{
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
			timeSinceLastBlock: time.Second,
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewIntFromBigInt(
					big_testutil.Int64MulPow10(99_000, 18),
				), // 99k tokens < baseline (250k)
				dtypes.NewIntFromBigInt(
					big_testutil.Int64MulPow10(990_000, 18),
				), // 0.99M tokens, < baseline (2.5M)
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewIntFromBigInt(
					big_testutil.MustFirst(new(big.Int).SetString("99069444444444444444444", 10)),
				), // recovered by 1/3600 * 250k ~= 69.4444 tokens
				dtypes.NewIntFromBigInt(
					big_testutil.MustFirst(new(big.Int).SetString("990028935185185185185185", 10)),
				), // recovered by 1/86400 * 2.5M ~= 28.9351 tokens
			},
		},
		"Prev capacity = 0": {
			bigTvl: big.NewInt(1_000_000_000_000), // 1M token (assuming 6 decimals)
			limitParams: types.LimitParams{
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
			timeSinceLastBlock: time.Second + 150*time.Millisecond, // 1.15 second
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(0),
				dtypes.NewInt(0),
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(31_944_444), // recovered by 1.15/3600 * 100k ~= 31.94
				dtypes.NewInt(13_310_185), // recovered by 1.15/86400 * 1M ~= 13.31
			},
		},
		"Prev capacity = 0, capacity_diff rounds down": {
			bigTvl: big.NewInt(1_000_000_000_000), // 1M token (assuming 6 decimals)
			limitParams: types.LimitParams{
				Denom: testDenom,
				Limiters: []types.Limiter{
					// baseline = baseline minimum = 100k tokens
					{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
				},
			},
			timeSinceLastBlock: 12 * time.Second, // 12 second
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(0),
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(333_333_333), // recovered by 12/3600 * 100k ~= 333.333
			},
		},
		"Prev capacity = 2 * baseline, capacity_diff rounds down": {
			bigTvl: big.NewInt(1_000_000_000_000), // 1M token (assuming 6 decimals)
			limitParams: types.LimitParams{
				Denom: testDenom,
				Limiters: []types.Limiter{
					// baseline = baseline minimum = 100k tokens
					{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
				},
			},
			timeSinceLastBlock: 12 * time.Second, // 12 second
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(200_000_000_000),
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(199_666_666_667), // recovered by 12/3600 * 100k ~= 333.333
			},
		},
		"baseline < prev capacity < 2 * baseline": {
			bigTvl: big.NewInt(20_000_000_000_000), // 20M token (assuming 6 decimals)
			limitParams: types.LimitParams{
				Denom: testDenom,
				Limiters: []types.Limiter{
					// baseline = 200k tokens
					{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					// baseline = 2M tokens
					{
						Period:          86_400 * time.Second,
						BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
						BaselineTvlPpm:  100_000,                          // 10%
					},
				},
			},
			timeSinceLastBlock: time.Second + 150*time.Millisecond, // 1.15 second
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(329_000_000_000),
				dtypes.NewInt(3_500_000_000_000),
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(328_936_111_112),   // recovered by 1.15/3600 * 200k ~= 63.89
				dtypes.NewInt(3_499_973_379_630), // recovered by 1.15/86400 * 2M ~=  26.62
			},
		},
		"prev capacity > 2 * baseline + capacity < baseline": {
			bigTvl: big.NewInt(20_000_000_000_000), // 20M token (assuming 6 decimals)
			limitParams: types.LimitParams{
				Denom: testDenom,
				Limiters: []types.Limiter{
					// baseline = 200k tokens
					{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					// baseline = 2M tokens
					{
						Period:          86_400 * time.Second,
						BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
						BaselineTvlPpm:  100_000,                          // 10%
					},
				},
			},
			timeSinceLastBlock: time.Second + 150*time.Millisecond, // 1.15 second
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(629_000_000_000),   // > 2 * baseline
				dtypes.NewInt(1_200_000_000_000), // < baseline
			},
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(628_862_958_334),   // recovered by 1.15/3600 * (629k - 200k) ~= 137.04
				dtypes.NewInt(1_200_026_620_370), //  recovered by 1.15/86400 * 2M ~= 26.62
			},
		},
		"Error: len(capacityList) != len(limiters)": {
			bigTvl: big.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
			limitParams: types.LimitParams{
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
			timeSinceLastBlock: time.Second + 90*time.Millisecond, // 1.09 second
			prevCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(99_000_000_000),
				dtypes.NewInt(990_000_000_000),
				dtypes.NewInt(0),
			},
			expectedErr: errorsmod.Wrapf(
				types.ErrMismatchedCapacityLimitersLength,
				"denom = %v, len(limiters) = %v, len(prevCapacityList) = %v",
				testDenom,
				2,
				3,
			),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			newCapacityList, err := util.CalculateNewCapacityList(
				tc.bigTvl,
				tc.limitParams,
				tc.prevCapacityList,
				tc.timeSinceLastBlock,
			)

			if tc.expectedErr != nil {
				require.Error(t, tc.expectedErr, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t,
				tc.expectedCapacityList,
				newCapacityList,
			)
		})
	}
}
