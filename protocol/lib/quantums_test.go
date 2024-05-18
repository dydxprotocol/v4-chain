package lib_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
)

func TestBaseToQuoteQuantums(t *testing.T) {
	tests := map[string]struct {
		bigBaseQuantums              *big.Int
		baseCurrencyAtomicResolution int32
		priceValue                   uint64
		priceExponent                int32
		bigExpectedQuoteQuantums     *big.Int
	}{
		"Converts from base to quote quantums": {
			bigBaseQuantums:              big.NewInt(1),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     big.NewInt(10),
		},
		"Correctly converts negative value": {
			bigBaseQuantums:              big.NewInt(-100),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     big.NewInt(-1000),
		},
		"priceExponent is negative": {
			bigBaseQuantums:              big.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                -5,
			bigExpectedQuoteQuantums:     big.NewInt(350),
		},
		"priceExponent is zero": {
			bigBaseQuantums:              big.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                0,
			bigExpectedQuoteQuantums:     big.NewInt(35_000_000),
		},
		"baseCurrencyAtomicResolution is greater than 10^6": {
			bigBaseQuantums:              big.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   7,
			priceExponent:                0,
			bigExpectedQuoteQuantums:     big.NewInt(350_000),
		},
		"baseCurrencyAtomicResolution is less than 10^6": {
			bigBaseQuantums:              big.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -4,
			priceValue:                   7,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     big.NewInt(35_000_000_000),
		},
		"Calculation rounds down": {
			bigBaseQuantums:              big.NewInt(9),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     big.NewInt(0),
		},
		"Negative calculation rounds up": {
			bigBaseQuantums:              big.NewInt(-9),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     big.NewInt(0),
		},
		"Calculation overflows": {
			bigBaseQuantums:              new(big.Int).SetUint64(math.MaxUint64),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   2,
			priceExponent:                0,
			bigExpectedQuoteQuantums:     big_testutil.MustFirst(new(big.Int).SetString("36893488147419103230", 10)),
		},
		"Calculation underflows": {
			bigBaseQuantums:              big_testutil.MustFirst(new(big.Int).SetString("-18446744073709551615", 10)),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   2,
			priceExponent:                0,
			bigExpectedQuoteQuantums:     big_testutil.MustFirst(new(big.Int).SetString("-36893488147419103230", 10)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			quoteQuantums := lib.BaseToQuoteQuantums(
				tc.bigBaseQuantums,
				tc.baseCurrencyAtomicResolution,
				tc.priceValue,
				tc.priceExponent,
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

func TestBaseToQuoteQuantumsInt256(t *testing.T) {
	tests := map[string]struct {
		baseQuantums                 *int256.Int
		baseCurrencyAtomicResolution int32
		priceValue                   uint64
		priceExponent                int32
		bigExpectedQuoteQuantums     *int256.Int
	}{
		"Converts from base to quote quantums": {
			baseQuantums:                 int256.NewInt(1),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     int256.NewInt(10),
		},
		"Correctly converts negative value": {
			baseQuantums:                 int256.NewInt(-100),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     int256.NewInt(-1000),
		},
		"priceExponent is negative": {
			baseQuantums:                 int256.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                -5,
			bigExpectedQuoteQuantums:     int256.NewInt(350),
		},
		"priceExponent is zero": {
			baseQuantums:                 int256.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                0,
			bigExpectedQuoteQuantums:     int256.NewInt(35_000_000),
		},
		"baseCurrencyAtomicResolution is greater than 10^6": {
			baseQuantums:                 int256.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   7,
			priceExponent:                0,
			bigExpectedQuoteQuantums:     int256.NewInt(350_000),
		},
		"baseCurrencyAtomicResolution is less than 10^6": {
			baseQuantums:                 int256.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -4,
			priceValue:                   7,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     int256.NewInt(35_000_000_000),
		},
		"Calculation rounds down": {
			baseQuantums:                 int256.NewInt(9),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     int256.NewInt(0),
		},
		"Negative calculation rounds up": {
			baseQuantums:                 int256.NewInt(-9),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedQuoteQuantums:     int256.NewInt(0),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			quoteQuantums := lib.BaseToQuoteQuantumsInt256(
				tc.baseQuantums,
				tc.baseCurrencyAtomicResolution,
				tc.priceValue,
				tc.priceExponent,
			)
			if !tc.bigExpectedQuoteQuantums.Eq(quoteQuantums) {
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

func TestQuoteToBaseQuantums(t *testing.T) {
	tests := map[string]struct {
		bigQuoteQuantums             *big.Int
		baseCurrencyAtomicResolution int32
		priceValue                   uint64
		priceExponent                int32
		bigExpectedBaseQuantums      *big.Int
	}{
		"Converts from base to quote quantums": {
			bigQuoteQuantums:             big.NewInt(10),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedBaseQuantums:      big.NewInt(1),
		},
		"Correctly converts negative value": {
			bigQuoteQuantums:             big.NewInt(-1000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedBaseQuantums:      big.NewInt(-100),
		},
		"priceExponent is negative": {
			bigQuoteQuantums:             big.NewInt(350),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                -5,
			bigExpectedBaseQuantums:      big.NewInt(5_000_000),
		},
		"priceExponent is zero": {
			bigQuoteQuantums:             big.NewInt(35_000_000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                0,
			bigExpectedBaseQuantums:      big.NewInt(5_000_000),
		},
		"realistic values: 1 BTC at $29001": {
			bigQuoteQuantums:             big.NewInt(29_001_000_000), // $29_001
			baseCurrencyAtomicResolution: -10,
			priceValue:                   2_900_100_000,
			priceExponent:                -5,
			bigExpectedBaseQuantums:      big.NewInt(10_000_000_000),
		},
		"realistic values: 25.123 BTC at $29001": {
			bigQuoteQuantums:             big.NewInt(728_592_123_000), // $728_592.123
			baseCurrencyAtomicResolution: -10,
			priceValue:                   2_900_100_000,
			priceExponent:                -5,
			bigExpectedBaseQuantums:      big.NewInt(251_230_000_000),
		},
		"baseCurrencyAtomicResolution is greater than 10^6": {
			bigQuoteQuantums:             big.NewInt(350_000),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   7,
			priceExponent:                0,
			bigExpectedBaseQuantums:      big.NewInt(5_000_000),
		},
		"baseCurrencyAtomicResolution is less than 10^6": {
			bigQuoteQuantums:             big.NewInt(35_000_000_000),
			baseCurrencyAtomicResolution: -4,
			priceValue:                   7,
			priceExponent:                1,
			bigExpectedBaseQuantums:      big.NewInt(5_000_000),
		},
		"Calculation rounds down": {
			bigQuoteQuantums:             big.NewInt(99),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedBaseQuantums:      big.NewInt(9),
		},
		"Negative calculation rounds up": {
			bigQuoteQuantums:             big.NewInt(-99),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			bigExpectedBaseQuantums:      big.NewInt(-9),
		},
		"Calculation overflows": {
			bigQuoteQuantums:             new(big.Int).SetUint64(math.MaxUint64),
			baseCurrencyAtomicResolution: -7,
			priceValue:                   1,
			priceExponent:                0,
			bigExpectedBaseQuantums:      big_testutil.MustFirst(new(big.Int).SetString("184467440737095516150", 10)),
		},
		"Calculation underflows": {
			bigQuoteQuantums:             big_testutil.MustFirst(new(big.Int).SetString("-18446744073709551615", 10)),
			baseCurrencyAtomicResolution: -7,
			priceValue:                   1,
			priceExponent:                0,
			bigExpectedBaseQuantums:      big_testutil.MustFirst(new(big.Int).SetString("-184467440737095516150", 10)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			baseQuantums := lib.QuoteToBaseQuantums(
				tc.bigQuoteQuantums,
				tc.baseCurrencyAtomicResolution,
				tc.priceValue,
				tc.priceExponent,
			)
			if tc.bigExpectedBaseQuantums.Cmp(baseQuantums) != 0 {
				t.Fatalf(
					"%s: expectedBaseQuantums: %s, baseQuantums: %s",
					name,
					tc.bigExpectedBaseQuantums.String(),
					baseQuantums.String(),
				)
			}
		})
	}
}
