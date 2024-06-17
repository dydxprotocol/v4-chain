package prices_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	"github.com/stretchr/testify/require"
)

func TestMustHumanPriceToMarketPrice(t *testing.T) {
	tests := []struct {
		humanPrice  string
		exponent    int32
		expected    uint64
		expectPanic bool
	}{
		{"20000", -5, 2000_000_000, false},
		{"12345.67", -3, 12_345_670, false},
		{"1.123", -8, 112_300_000, false},
		{"0.00000001", -8, 1, false},
		{"1", -10, 10_000_000_000, false},
		{"0.0000000001", -10, 1, false},
		{"500", 2, 5, false},
		{"500", 0, 500, false},
		{"abc", -8, 0, true}, // Invalid humanPrice
	}

	for _, test := range tests {
		if test.expectPanic {
			require.Panics(t, func() {
				prices.MustHumanPriceToMarketPrice(test.humanPrice, test.exponent)
			}, "For humanPrice %s and exponent %d, expected a panic", test.humanPrice, test.exponent)
		} else {
			result := prices.MustHumanPriceToMarketPrice(test.humanPrice, test.exponent)
			require.Equal(t, test.expected, result, "expected = %v, result = %v for humanPrice %s and exponent %d",
				test.expected, result, test.humanPrice, test.exponent)
		}
	}
}
