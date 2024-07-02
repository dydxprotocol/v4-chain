package funding_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/funding"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestGetFundingIndexDelta(t *testing.T) {
	testCases := map[string]struct {
		perp                 types.Perpetual
		marketPrice          pricestypes.MarketPrice
		big8hrFundingRatePpm *big.Int
		timeSinceLastFunding uint32
		expected             *big.Int
	}{
		"Positive Funding Rate (rounds towards zero)": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithAtomicResolution(-12),
			),
			marketPrice: pricestypes.MarketPrice{
				Id:       0,
				Exponent: 0,
				Price:    1_000,
			},
			big8hrFundingRatePpm: big.NewInt(1_001_999),
			timeSinceLastFunding: 8 * 60 * 60,
			expected:             big.NewInt(1_001),
		},
		"Negative Funding Rate (rounds towards zero)": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithAtomicResolution(-12),
			),
			marketPrice: pricestypes.MarketPrice{
				Id:       0,
				Exponent: 0,
				Price:    1_000,
			},
			big8hrFundingRatePpm: big.NewInt(-1_001_999),
			timeSinceLastFunding: 8 * 60 * 60,
			expected:             big.NewInt(-1_001),
		},
		"Varied parameters (1)": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithAtomicResolution(-4),
			),
			marketPrice: pricestypes.MarketPrice{
				Id:       0,
				Exponent: 2,
				Price:    1_000,
			},
			big8hrFundingRatePpm: big.NewInt(-1_001_999),
			timeSinceLastFunding: 8 * 60 * 60 / 2,
			expected:             big.NewInt(-5_009_995_000_000),
		},
		"Varied parameters (2)": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithAtomicResolution(0),
			),
			marketPrice: pricestypes.MarketPrice{
				Id:       0,
				Exponent: -6,
				Price:    1_000,
			},
			big8hrFundingRatePpm: big.NewInt(-1_001_999),
			timeSinceLastFunding: 8 * 60 * 60 / 8,
			expected:             big.NewInt(-125_249_875),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := funding.GetFundingIndexDelta(
				tc.perp,
				tc.marketPrice,
				tc.big8hrFundingRatePpm,
				tc.timeSinceLastFunding,
			)

			require.Equal(t, tc.expected, result)
		})
	}
}
