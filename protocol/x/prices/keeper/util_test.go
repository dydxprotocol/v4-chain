package keeper

import (
	"errors"
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
)

const (
	maxUint64 = uint64(18_446_744_073_709_551_615) // 2 ^ 64 - 1
	maxUint32 = uint32(4_294_967_295)              // 2 ^ 32 - 1
)

func TestGetMinPriceChangeAmountForMarket(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		market types.Market

		// Expected.
		expectedResult uint64
		expectedPanic  error
	}{
		"Valid": {
			market: types.Market{
				Price:             uint64(123_000),
				MinPriceChangePpm: uint32(1_000), // 0.1%
			},
			expectedResult: 123,
		},
		"Valid: discards decimal": {
			market: types.Market{
				Price:             uint64(1_234),
				MinPriceChangePpm: uint32(1_000), // 0.1%
			},
			expectedResult: 1,
		},
		"Zero": {
			market: types.Market{
				Price:             uint64(0),
				MinPriceChangePpm: uint32(1_000), // 0.1%
			},
			expectedResult: 0,
		},
		"Result exceeds max uint64": {
			market: types.Market{
				Price:             math.MaxUint64,
				MinPriceChangePpm: uint32(1_000_001), // must be <= 1,000,000
			},
			expectedPanic: errors.New(
				"getMinPriceChangeAmountForMarket: min price change amount is greater than max uint64 value",
			)},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic != nil {
				require.PanicsWithError(
					t,
					tc.expectedPanic.Error(),
					func() { getMinPriceChangeAmountForMarket(tc.market) })
				return
			}

			result := getMinPriceChangeAmountForMarket(tc.market)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestIsTowardsIndexPrice(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		oldPrice   uint64
		newPrice   uint64
		indexPrice uint64

		// Expected.
		expectedResult bool
	}{
		"Towards: curr < new < idx": {
			oldPrice:       1,
			newPrice:       2,
			indexPrice:     3,
			expectedResult: true,
		},
		"Towards: idx < new < curr": {
			oldPrice:       3,
			newPrice:       2,
			indexPrice:     1,
			expectedResult: true,
		},
		"Towards: curr == new < idx": {
			oldPrice:       1,
			newPrice:       1,
			indexPrice:     2,
			expectedResult: true,
		},
		"Towards: idx < curr == new": {
			oldPrice:       2,
			newPrice:       2,
			indexPrice:     1,
			expectedResult: true,
		},
		"Towards: curr < new == idx": {
			oldPrice:       1,
			newPrice:       2,
			indexPrice:     2,
			expectedResult: true,
		},
		"Towards: new == idx < curr": {
			oldPrice:       2,
			newPrice:       1,
			indexPrice:     1,
			expectedResult: true,
		},
		"Towards: new == idx == curr": {
			oldPrice:       1,
			newPrice:       1,
			indexPrice:     1,
			expectedResult: true,
		},
		"Not Towards: new < curr < idx": {
			oldPrice:       2,
			newPrice:       1,
			indexPrice:     3,
			expectedResult: false,
		},
		"Not Towards: new < idx < curr": {
			oldPrice:       3,
			newPrice:       1,
			indexPrice:     2,
			expectedResult: false,
		},
		"Not Towards: new < idx == curr": {
			oldPrice:       2,
			newPrice:       1,
			indexPrice:     2,
			expectedResult: false,
		},
		"Not Towards: curr < idx < new": {
			oldPrice:       1,
			newPrice:       3,
			indexPrice:     2,
			expectedResult: false,
		},
		"Not Towards: idx < curr < new": {
			oldPrice:       2,
			newPrice:       3,
			indexPrice:     1,
			expectedResult: false,
		},
		"Not Towards: curr == idx < new": {
			oldPrice:       1,
			newPrice:       2,
			indexPrice:     1,
			expectedResult: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := isTowardsIndexPrice(tc.oldPrice, tc.newPrice, tc.indexPrice)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestIsCrossingIndexPrice(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		oldPrice   uint64
		newPrice   uint64
		indexPrice uint64

		// Expected.
		expectedResult bool
	}{
		"Crossing: curr < index < new": {
			oldPrice:       1,
			newPrice:       3,
			indexPrice:     2,
			expectedResult: true,
		},
		"Crossing: new < index < curr": {
			oldPrice:       3,
			newPrice:       1,
			indexPrice:     2,
			expectedResult: true,
		},
		"Not Crossing: curr < new < index": {
			oldPrice:       1,
			newPrice:       2,
			indexPrice:     3,
			expectedResult: false,
		},
		"Not Crossing: new < curr < index": {
			oldPrice:       2,
			newPrice:       1,
			indexPrice:     3,
			expectedResult: false,
		},
		"Not Crossing: curr < new = index": {
			oldPrice:       1,
			newPrice:       2,
			indexPrice:     2,
			expectedResult: false,
		},
		"Not Crossing: new < curr = index": {
			oldPrice:       2,
			newPrice:       1,
			indexPrice:     2,
			expectedResult: false,
		},
		"Not Crossing: new = curr = index": {
			oldPrice:       1,
			newPrice:       1,
			indexPrice:     1,
			expectedResult: false,
		},
		"Not Crossing: index < curr < new": {
			oldPrice:       2,
			newPrice:       3,
			indexPrice:     1,
			expectedResult: false,
		},
		"Not Crossing: index < new < curr": {
			oldPrice:       3,
			newPrice:       2,
			indexPrice:     1,
			expectedResult: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := isCrossingIndexPrice(tc.oldPrice, tc.newPrice, tc.indexPrice)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestComputeTickSizePpm(t *testing.T) {
	tests := map[string]struct {
		oldPrice          uint64
		minPriceChangePpm uint32
		expected          *big.Int
	}{
		"non-overflow case": {
			oldPrice:          1_000_000_000_000,
			minPriceChangePpm: 50,
			expected:          new(big.Int).SetUint64(1_000_000_000_000 * 50),
		},
		"overflow case": {
			oldPrice:          maxUint64,
			minPriceChangePpm: maxUint32,
			expected: new(big.Int).Mul(
				new(big.Int).SetUint64(maxUint64),
				new(big.Int).SetUint64(uint64(maxUint32))),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := computeTickSizePpm(tc.oldPrice, tc.minPriceChangePpm)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestPriceDeltaIsWithinOneTick(t *testing.T) {
	tests := map[string]struct {
		priceDelta  *big.Int
		tickSizePpm *big.Int
		expected    bool
	}{
		"Within: less than one tick": {
			priceDelta:  new(big.Int).SetUint64(500_000),
			tickSizePpm: new(big.Int).SetUint64(600_000_000_000),
			expected:    true,
		},
		"Within: exactly one tick": {
			priceDelta:  new(big.Int).SetUint64(500_000),
			tickSizePpm: new(big.Int).SetUint64(500_000_000_000),
			expected:    true,
		},
		"Not within: greater than one tick": {
			priceDelta:  new(big.Int).SetUint64(500_000),
			tickSizePpm: new(big.Int).SetUint64(400_000_000_000),
			expected:    false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, priceDeltaIsWithinOneTick(tc.priceDelta, tc.tickSizePpm))
		})
	}
}

func TestNewPriceMeetSqrtCondition(t *testing.T) {
	tests := map[string]struct {
		oldDelta    *big.Int
		newDelta    *big.Int
		tickSizePpm *big.Int
		expected    bool
	}{
		"Meets condition: new_ticks < sqrt(old_ticks)": {
			newDelta:    new(big.Int).SetUint64(4999),
			oldDelta:    new(big.Int).SetUint64(500000),
			tickSizePpm: new(big.Int).SetUint64(50000000),
			expected:    true,
		},
		"Meets condition: new_ticks == sqrt(old_ticks)": {
			newDelta:    new(big.Int).SetUint64(5000),
			oldDelta:    new(big.Int).SetUint64(500000),
			tickSizePpm: new(big.Int).SetUint64(50000000),
			expected:    true,
		},
		"Does not meet condition: new_ticks > sqrt(old_ticks)": {
			newDelta:    new(big.Int).SetUint64(5001),
			oldDelta:    new(big.Int).SetUint64(500000),
			tickSizePpm: new(big.Int).SetUint64(50000000),
			expected:    false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, newPriceMeetsSqrtCondition(tc.oldDelta, tc.newDelta, tc.tickSizePpm))
		})
	}
}

func TestMaximumAllowedPriceDelta(t *testing.T) {
	tests := map[string]struct {
		oldDelta           *big.Int
		tickSizePpm        *big.Int
		expectedMaxAllowed *big.Int
	}{
		"no precision loss": {
			oldDelta:           new(big.Int).SetUint64(500_000),
			tickSizePpm:        new(big.Int).SetUint64(500_000),
			expectedMaxAllowed: new(big.Int).SetUint64(500),
		},
		"precision loss from division, sqrt": {
			oldDelta:           new(big.Int).SetUint64(512_345),
			tickSizePpm:        new(big.Int).SetUint64(567_899),
			expectedMaxAllowed: new(big.Int).SetUint64(539),
		},
		"precision loss from sqrt": {
			oldDelta:           new(big.Int).SetUint64(512_000),
			tickSizePpm:        new(big.Int).SetUint64(567_000),
			expectedMaxAllowed: new(big.Int).SetUint64(538),
		},
		"tiny currency change": {
			oldDelta:           new(big.Int).SetUint64(50),
			tickSizePpm:        new(big.Int).SetUint64(567_000),
			expectedMaxAllowed: new(big.Int).SetUint64(5),
		},
		"no error: unrealistically small tick size": {
			oldDelta:           new(big.Int).SetUint64(512_000),
			tickSizePpm:        new(big.Int).SetUint64(50),
			expectedMaxAllowed: new(big.Int).SetUint64(5),
		},
		"no error: zero inputs": {
			oldDelta:           new(big.Int).SetUint64(0),
			tickSizePpm:        new(big.Int).SetUint64(0),
			expectedMaxAllowed: new(big.Int).SetUint64(0),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedMaxAllowed, maximumAllowedPriceDelta(tc.oldDelta, tc.tickSizePpm))
		})
	}
}
