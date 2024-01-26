package util_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/util"
	"github.com/stretchr/testify/require"
)

func TestCalculateNewCapacityList(t *testing.T) {
	tests := map[string]struct {
		bigTvl               *big.Int
		limiterCapacityList  []types.LimiterCapacity
		expectedCapacityList []dtypes.SerializableInt
		timeSinceLastBlock   time.Duration
	}{
		"Prev capacity equals baseline": {
			bigTvl: big.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
			limiterCapacityList: []types.LimiterCapacity{
				{
					Limiter: types.Limiter{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					Capacity: dtypes.NewInt(250_000_000_000), // 250k tokens, which equals baseline
				},
				{
					Limiter: types.Limiter{
						Period:          86_400 * time.Second,
						BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
						BaselineTvlPpm:  100_000,                          // 10%
					},
					Capacity: dtypes.NewInt(2_500_000_000_000), // 2.5M tokens, which equals baseline
				},
			},
			timeSinceLastBlock: time.Second,
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(250_000_000_000),   // 250k tokens
				dtypes.NewInt(2_500_000_000_000), // 2.5M tokens
			},
		},
		"Prev capacity < baseline": {
			bigTvl: big.NewInt(25_000_000_000_000), // 25M token (assuming 6 decimals)
			limiterCapacityList: []types.LimiterCapacity{
				{
					// baseline = 25M * 1% = 250k tokens
					Limiter: types.Limiter{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					Capacity: dtypes.NewInt(99_000_000_000), // 99k tokens, < baseline (250k)
				},
				{
					// baseline = 25M * 10% = 2.5M tokens
					Limiter: types.Limiter{
						Period:          86_400 * time.Second,
						BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
						BaselineTvlPpm:  100_000,                          // 10%
					},
					Capacity: dtypes.NewInt(990_000_000_000), // 0.99M tokens, < baseline (2.5M)
				},
			},
			timeSinceLastBlock: time.Second + 90*time.Millisecond, // 1.09 second
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(99_075_694_444),  // recovered by 1.09/3600 * 250k = 75.694444 tokens
				dtypes.NewInt(990_031_539_351), // recovered by 1.09/86400 * 2.5M = 31.539 tokens
			},
		},
		"prev capacity < baseline, 18 decimals": {
			bigTvl: big_testutil.Int64MulPow10(25, 24), // 25M tokens
			limiterCapacityList: []types.LimiterCapacity{
				{
					// baseline = 25M * 1% = 250k tokens
					Limiter: types.Limiter{
						Period: 3_600 * time.Second,
						BaselineMinimum: dtypes.NewIntFromBigInt(
							big_testutil.Int64MulPow10(100_000, 18), // 100k tokens(assuming 18 decimals)
						),
						BaselineTvlPpm: 10_000, // 1%
					},
					Capacity: dtypes.NewIntFromBigInt(
						big_testutil.Int64MulPow10(99_000, 18),
					), // 99k tokens < baseline (250k)
				},
				{
					// baseline = 25M * 10% = 2.5M tokens
					Limiter: types.Limiter{
						Period: 86_400 * time.Second,
						BaselineMinimum: dtypes.NewIntFromBigInt(
							big_testutil.Int64MulPow10(1_000_000, 18), // 1M tokens(assuming 18 decimals)
						),
						BaselineTvlPpm: 100_000, // 10%
					},
					Capacity: dtypes.NewIntFromBigInt(
						big_testutil.Int64MulPow10(990_000, 18),
					), // 0.99M tokens, < baseline (2.5M)
				},
			},
			timeSinceLastBlock: time.Second,
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
			limiterCapacityList: []types.LimiterCapacity{
				{
					// baseline = baseline minimum = 100k tokens
					Limiter: types.Limiter{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					Capacity: dtypes.NewInt(0),
				},
				{
					// baseline = baseline minimum = 1M tokens
					Limiter: types.Limiter{
						Period:          86_400 * time.Second,
						BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
						BaselineTvlPpm:  100_000,                          // 10%
					},
					Capacity: dtypes.NewInt(0),
				},
			},
			timeSinceLastBlock: time.Second + 150*time.Millisecond, // 1.15 second
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(31_944_444), // recovered by 1.15/3600 * 100k ~= 31.94
				dtypes.NewInt(13_310_185), // recovered by 1.15/86400 * 1M ~= 13.31
			},
		},
		"Prev capacity = 0, capacity_diff rounds down": {
			bigTvl: big.NewInt(1_000_000_000_000), // 1M token (assuming 6 decimals)
			limiterCapacityList: []types.LimiterCapacity{
				{
					// baseline = baseline minimum = 100k tokens
					Limiter: types.Limiter{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					Capacity: dtypes.NewInt(0),
				},
			},
			timeSinceLastBlock: 12 * time.Second, // 12 second
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(333_333_333), // recovered by 12/3600 * 100k ~= 333.333
			},
		},
		"Prev capacity = 2 * baseline, capacity_diff rounds down": {
			bigTvl: big.NewInt(1_000_000_000_000), // 1M token (assuming 6 decimals)
			limiterCapacityList: []types.LimiterCapacity{
				{
					// baseline = baseline minimum = 100k tokens
					Limiter: types.Limiter{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					Capacity: dtypes.NewInt(200_000_000_000),
				},
			},
			timeSinceLastBlock: 12 * time.Second, // 12 second
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(199_666_666_667), // recovered by 12/3600 * 100k ~= 333.333
			},
		},
		"baseline < prev capacity < 2 * baseline": {
			bigTvl: big.NewInt(20_000_000_000_000), // 20M token (assuming 6 decimals)
			limiterCapacityList: []types.LimiterCapacity{
				{
					// baseline = 200k tokens
					Limiter: types.Limiter{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					Capacity: dtypes.NewInt(329_000_000_000),
				},
				{
					// baseline = 2M tokens
					Limiter: types.Limiter{
						Period:          86_400 * time.Second,
						BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
						BaselineTvlPpm:  100_000,                          // 10%
					},
					Capacity: dtypes.NewInt(3_500_000_000_000),
				},
			},
			timeSinceLastBlock: time.Second + 150*time.Millisecond, // 1.15 second
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(328_936_111_112),   // recovered by 1.15/3600 * 200k ~= 63.89
				dtypes.NewInt(3_499_973_379_630), // recovered by 1.15/86400 * 2M ~=  26.62
			},
		},
		"prev capacity > 2 * baseline + capacity < baseline": {
			bigTvl: big.NewInt(20_000_000_000_000), // 20M token (assuming 6 decimals)
			limiterCapacityList: []types.LimiterCapacity{
				{
					// baseline = 200k tokens
					Limiter: types.Limiter{
						Period:          3_600 * time.Second,
						BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k tokens (assuming 6 decimals)
						BaselineTvlPpm:  10_000,                         // 1%
					},
					Capacity: dtypes.NewInt(629_000_000_000),
				},
				{
					// baseline = 2M tokens
					Limiter: types.Limiter{
						Period:          86_400 * time.Second,
						BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1M tokens (assuming 6 decimals)
						BaselineTvlPpm:  100_000,                          // 10%
					},
					Capacity: dtypes.NewInt(1_200_000_000_000),
				},
			},
			timeSinceLastBlock: time.Second + 150*time.Millisecond, // 1.15 second
			expectedCapacityList: []dtypes.SerializableInt{
				dtypes.NewInt(628_862_958_334),   // recovered by 1.15/3600 * (629k - 200k) ~= 137.04
				dtypes.NewInt(1_200_026_620_370), //  recovered by 1.15/86400 * 2M ~= 26.62
			},
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			newCapacityList := util.CalculateNewCapacityList(
				tc.bigTvl,
				tc.limiterCapacityList,
				tc.timeSinceLastBlock,
			)

			require.Equal(t,
				tc.expectedCapacityList,
				newCapacityList,
			)
		})
	}
}
