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
		leveragePpm   *big.Int
		expected      *big.Int
	}{
		"Zero skew factor and leverage": {
			skewFactorPpm: 0,
			leveragePpm:   big.NewInt(0),
			expected:      big.NewInt(0),
		},
		"Non-zero skew factor, zero leverage": {
			skewFactorPpm: 1_000_000,
			leveragePpm:   big.NewInt(0),
			expected:      big.NewInt(0),
		},
		"Zero skew factor, non-zero leverage": {
			skewFactorPpm: 0,
			leveragePpm:   big.NewInt(1_000_000),
			expected:      big.NewInt(0),
		},
		"Small skew factor and small positive leverage": {
			skewFactorPpm: 500_000,             // 0.5
			leveragePpm:   big.NewInt(800_000), // 0.8
			// 0.5 * 0.8^2 + 0.5^2 * 0.8^3 / 3 ~= 0.362666
			// round up to 0.362667
			expected: big.NewInt(362_667),
		},
		"Small skew factor and small negative leverage": {
			skewFactorPpm: 500_000,              // 0.5
			leveragePpm:   big.NewInt(-800_000), // -0.8
			// 0.5 * (-0.8)^2 + 0.5^2 * (-0.8)^3 / 3 ~= 0.277333
			// round up to 0.277334
			expected: big.NewInt(277_334),
		},
		"Large skew factor and large positive leverage": {
			skewFactorPpm: 5_000_000,             // 5
			leveragePpm:   big.NewInt(8_700_000), // 8.7
			// 5 * (8.7)^2 + 5^2 * (8.7)^3 / 3 = 5865.975
			expected: big.NewInt(5_865_975_000),
		},
		"Large skew factor and large negative leverage": {
			skewFactorPpm: 5_000_000,              // 5
			leveragePpm:   big.NewInt(-8_700_000), // -8.7
			// 5 * (-8.7)^2 + 5^2 * (-8.7)^3 / 3 = -5109.075
			expected: big.NewInt(-5_109_075_000),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := vault.SkewAntiderivativePpm(tc.skewFactorPpm, tc.leveragePpm)
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
