package vault_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/lib/vault"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func TestSkewAntiderivativePpm(t *testing.T) {
	tests := map[string]struct {
		skewFactorPpm uint32
		leverage      *big.Rat
		expected      *big.Rat
	}{
		"Zero skew factor and leverage": {
			skewFactorPpm: 0,
			leverage:      big.NewRat(0, 1),
			expected:      big.NewRat(0, 1),
		},
		"Non-zero skew factor, zero leverage": {
			skewFactorPpm: 1_000_000,
			leverage:      big.NewRat(0, 1),
			expected:      big.NewRat(0, 1),
		},
		"Zero skew factor, non-zero leverage": {
			skewFactorPpm: 0,
			leverage:      big.NewRat(1_000_000, 1),
			expected:      big.NewRat(0, 1),
		},
		"Small skew factor and small positive leverage": {
			skewFactorPpm: 500_000,          // 0.5
			leverage:      big.NewRat(4, 5), // 0.8
			// 0.5 * 0.8^2 + 0.5^2 * 0.8^3 / 3 = 136/375
			expected: big.NewRat(136, 375),
		},
		"Small skew factor and small negative leverage": {
			skewFactorPpm: 500_000,           // 0.5
			leverage:      big.NewRat(-4, 5), // -0.8
			// 0.5 * (-0.8)^2 + 0.5^2 * (-0.8)^3 / 3 = 104/375
			expected: big.NewRat(104, 375),
		},
		"Large skew factor and large positive leverage": {
			skewFactorPpm: 5_000_000,          // 5
			leverage:      big.NewRat(87, 10), // 8.7
			// 5 * (8.7)^2 + 5^2 * (8.7)^3 / 3 = 234639/40
			expected: big.NewRat(234_639, 40),
		},
		"Large skew factor and large negative leverage": {
			skewFactorPpm: 5_000_000,           // 5
			leverage:      big.NewRat(-87, 10), // -8.7
			// 5 * (-8.7)^2 + 5^2 * (-8.7)^3 / 3 = -204363/40
			expected: big.NewRat(-204_363, 40),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := vault.SkewAntiderivative(tc.skewFactorPpm, tc.leverage)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestSpreadPpm(t *testing.T) {
	tests := map[string]struct {
		quotingParams *types.QuotingParams
		marketParam   *pricestypes.MarketParam
		expected      uint32
	}{
		"SpreadMinPpm > SpreadBufferPpm + MinPriceChangePpm": {
			quotingParams: &types.QuotingParams{
				SpreadMinPpm:    1000,
				SpreadBufferPpm: 200,
			},
			marketParam: &pricestypes.MarketParam{
				MinPriceChangePpm: 500,
			},
			expected: 1000,
		},
		"SpreadMinPpm < SpreadBufferPpm + MinPriceChangePpm": {
			quotingParams: &types.QuotingParams{
				SpreadMinPpm:    1000,
				SpreadBufferPpm: 600,
			},
			marketParam: &pricestypes.MarketParam{
				MinPriceChangePpm: 500,
			},
			expected: 1100,
		},
		"SpreadMinPpm = SpreadBufferPpm + MinPriceChangePpm": {
			quotingParams: &types.QuotingParams{
				SpreadMinPpm:    1000,
				SpreadBufferPpm: 400,
			},
			marketParam: &pricestypes.MarketParam{
				MinPriceChangePpm: 600,
			},
			expected: 1000,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expected,
				vault.SpreadPpm(tc.quotingParams, tc.marketParam),
			)
		})
	}
}
