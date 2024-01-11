package util_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	ratelimitutil "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/util"
	"github.com/stretchr/testify/require"
)

func TestGetBaseline(t *testing.T) {
	tests := map[string]struct {
		supply           *big.Int
		limiter          types.Limiter
		expectedBaseline *big.Int
	}{
		"max(1% of TVL, 100k token), TVL = 5M token": {
			supply: big.NewInt(5_000_000_000_000), // 5M token
			limiter: types.Limiter{
				Period:          3_600 * time.Second,
				BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k token
				BaselineTvlPpm:  10_000,                         // 1%
			},
			expectedBaseline: big.NewInt(100_000_000_000), // 100k token (baseline minimum)
		},
		"max(1% of TVL, 100k token), TVL = 15M token": {
			supply: big.NewInt(15_000_000_000_000), // 10M token
			limiter: types.Limiter{
				Period:          3_600 * time.Second,
				BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k token
				BaselineTvlPpm:  10_000,                         // 1%
			},
			expectedBaseline: big.NewInt(150_000_000_000), // 150k token (1% of 15m)
		},
		"max(1% of TVL, 100k token), TVL = ~15M token, rounds down": {
			supply: big.NewInt(15_200_123_456_777),
			limiter: types.Limiter{
				Period:          3_600 * time.Second,
				BaselineMinimum: dtypes.NewInt(100_000_000_000), // 100k token
				BaselineTvlPpm:  10_000,                         // 1%
			},
			expectedBaseline: big.NewInt(152_001_234_567), // ~152k token (1% of 15.2m)
		},
		"max(10% of TVL, 1 million), TVL = 20M token": {
			supply: big.NewInt(20_000_000_000_000), // 20M token,
			limiter: types.Limiter{
				Period:          3_600 * time.Second,
				BaselineMinimum: dtypes.NewInt(100_000_000_000), // 1m token
				BaselineTvlPpm:  100_000,                        // 10%
			},
			expectedBaseline: big.NewInt(2_000_000_000_000), // 2m token (10% of 20m)
		},
		"max(10% of TVL, 1 million), TVL = 8M token": {
			supply: big.NewInt(8_000_000_000_000), // 2m token (10% of 20m)
			limiter: types.Limiter{
				Period:          3_600 * time.Second,
				BaselineMinimum: dtypes.NewInt(1_000_000_000_000), // 1m token
				BaselineTvlPpm:  100_000,                          // 10%
			},
			expectedBaseline: big.NewInt(1_000_000_000_000), // 1m token (baseline minimum)
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotBaseline := ratelimitutil.GetBaseline(tc.supply, tc.limiter)

			require.Equal(t, tc.expectedBaseline, gotBaseline, "retrieved baseline does not match the expected value")
		})
	}
}
