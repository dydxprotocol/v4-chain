package types_test

import (
	"errors"
	"math"
	"math/big"
	"testing"

	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestFillAmountToQuoteQuantums(t *testing.T) {
	tests := map[string]struct {
		baseQuantums              satypes.BaseQuantums
		subticks                  types.Subticks
		quantumConversionExponent int32
		bigExpectedQuoteQuantums  *big.Int
	}{
		"Converts from base to quote quantums": {
			baseQuantums:              1,
			subticks:                  1,
			quantumConversionExponent: 1,
			bigExpectedQuoteQuantums:  big.NewInt(10),
		},
		"quantumConversionExponent is negative": {
			baseQuantums:              1_000,
			subticks:                  5_000_000,
			quantumConversionExponent: -4,
			bigExpectedQuoteQuantums:  big.NewInt(500000),
		},
		"quantumConversionExponent is zero": {
			baseQuantums:              5_000_000,
			subticks:                  1,
			quantumConversionExponent: 0,
			bigExpectedQuoteQuantums:  big.NewInt(5_000_000),
		},
		"Calculation rounds down and can return 0 quoteQuantums": {
			baseQuantums:              9,
			subticks:                  1,
			quantumConversionExponent: -1,
			bigExpectedQuoteQuantums:  big.NewInt(0),
		},
		"Calculation overflows": {
			baseQuantums:              math.MaxUint64,
			subticks:                  1,
			quantumConversionExponent: 6,
			bigExpectedQuoteQuantums:  big_testutil.MustFirst(new(big.Int).SetString("18446744073709551615000000", 10)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			quoteQuantums := types.FillAmountToQuoteQuantums(
				tc.subticks,
				tc.baseQuantums,
				tc.quantumConversionExponent,
			)
			if tc.bigExpectedQuoteQuantums.Cmp(quoteQuantums) != 0 {
				t.Fatalf(
					"%s: expectedQuoteQuantums: %s, quoteQuantums: %s",
					name,
					tc.bigExpectedQuoteQuantums.String(),
					quoteQuantums.String(),
				)
			}
		})
	}
}

func TestGetAveragePriceSubticks(t *testing.T) {
	tests := map[string]struct {
		bigQuoteQuantums          *big.Int
		bigBaseQuantums           *big.Int
		quantumConversionExponent int32
		bigRatExpectedSubticks    *big.Rat
		expectedErr               error
	}{
		"Converts from quote quantums to subticks": {
			bigQuoteQuantums:          big.NewInt(10),
			bigBaseQuantums:           big.NewInt(1),
			quantumConversionExponent: 1,
			bigRatExpectedSubticks:    big.NewRat(1, 1),
		},
		"quantumConversionExponent is negative": {
			bigQuoteQuantums:          big.NewInt(500_000),
			bigBaseQuantums:           big.NewInt(1_000),
			quantumConversionExponent: -4,
			bigRatExpectedSubticks:    big.NewRat(5_000_000, 1),
		},
		"quantumConversionExponent is zero": {
			bigQuoteQuantums:          big.NewInt(5_000_000),
			bigBaseQuantums:           big.NewInt(5_000_000),
			quantumConversionExponent: 0,
			bigRatExpectedSubticks:    big.NewRat(1, 1),
		},
		"Not divisble": {
			bigQuoteQuantums:          big.NewInt(24_001_111),
			bigBaseQuantums:           big.NewInt(5_000_000),
			quantumConversionExponent: -2,
			bigRatExpectedSubticks:    big.NewRat(24001111, 50_000),
		},
		"Panics if bigBaseQuantums is zero": {
			bigQuoteQuantums:          big.NewInt(100),
			bigBaseQuantums:           big.NewInt(0),
			quantumConversionExponent: 6,
			expectedErr:               errors.New("GetAveragePriceSubticks: bigBaseQuantums = 0"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedErr != nil {
				// Call EndBlocker.
				require.PanicsWithError(t,
					tc.expectedErr.Error(),
					func() {
						//nolint:errcheck
						types.GetAveragePriceSubticks(
							tc.bigQuoteQuantums,
							tc.bigBaseQuantums,
							tc.quantumConversionExponent,
						)
					})
				return
			}

			bigRatSubticks := types.GetAveragePriceSubticks(
				tc.bigQuoteQuantums,
				tc.bigBaseQuantums,
				tc.quantumConversionExponent,
			)
			require.Equal(t,
				tc.bigRatExpectedSubticks,
				bigRatSubticks,
			)
		})
	}
}
