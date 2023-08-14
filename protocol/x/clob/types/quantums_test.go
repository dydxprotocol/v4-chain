package types_test

import (
	"errors"
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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

func TestNotionalToCoinAmount(t *testing.T) {
	tests := map[string]struct {
		notionalQuoteQuantums    *big.Int
		denomExp                 int32
		marketPrice              pricestypes.MarketPrice
		bigRatExpectedCoinAmount *big.Rat
	}{
		"$9.5 notional, ATOM price at $9.5, get amount in `uatom` (exp = -6)": {
			notionalQuoteQuantums: big.NewInt(9_500_000),
			marketPrice: pricestypes.MarketPrice{
				Price:    95_000,
				Exponent: -4,
			},
			denomExp:                 -6,
			bigRatExpectedCoinAmount: big.NewRat(1_000_000, 1),
		},
		"$4.75 notional, ATOM price at $9.5, get amount in `uatom` (exp = -6)": {
			notionalQuoteQuantums: big.NewInt(4_750_000),
			marketPrice: pricestypes.MarketPrice{
				Price:    95_000,
				Exponent: -4,
			},
			denomExp:                 -6,
			bigRatExpectedCoinAmount: big.NewRat(500_000, 1),
		},
		"$10.5 notional, ETH price at $2000, get amount in `gwei` (exp = -9)": {
			notionalQuoteQuantums: big.NewInt(10_500_000),
			marketPrice: pricestypes.MarketPrice{
				Price:    20_000_000_000,
				Exponent: -7,
			},
			denomExp:                 -9,
			bigRatExpectedCoinAmount: big.NewRat(5_250_000, 1),
		},
		"$1000 notional, ETH price at $2001.57, get amount in `gwei` (exp = -9)": {
			notionalQuoteQuantums: big.NewInt(1_000_000_000),
			marketPrice: pricestypes.MarketPrice{
				Price:    20_015_700_000,
				Exponent: -7,
			},
			denomExp:                 -9,
			bigRatExpectedCoinAmount: big.NewRat(100_000_000_000_000, 200157), // 499607807.871 Gwei, or 0.499607807871 Eth.
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bigRatCoinAmount := types.NotionalToCoinAmount(
				tc.notionalQuoteQuantums,
				lib.QuoteCurrencyAtomicResolution,
				tc.denomExp,
				tc.marketPrice,
			)
			require.Equal(t,
				tc.bigRatExpectedCoinAmount,
				bigRatCoinAmount,
			)
		})
	}
}
