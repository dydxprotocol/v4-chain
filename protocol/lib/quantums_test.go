package lib_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
)

func TestBaseToQuoteQuantums(t *testing.T) {
	tests := map[string]struct {
		baseQuantums                 *int256.Int
		baseCurrencyAtomicResolution int32
		priceValue                   uint64
		priceExponent                int32
		expectedQuoteQuantums        *int256.Int
	}{
		"Converts from base to quote quantums": {
			baseQuantums:                 int256.NewInt(1),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			expectedQuoteQuantums:        int256.NewInt(10),
		},
		"Correctly converts negative value": {
			baseQuantums:                 int256.NewInt(-100),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			expectedQuoteQuantums:        int256.NewInt(-1000),
		},
		"priceExponent is negative": {
			baseQuantums:                 int256.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                -5,
			expectedQuoteQuantums:        int256.NewInt(350),
		},
		"priceExponent is zero": {
			baseQuantums:                 int256.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                0,
			expectedQuoteQuantums:        int256.NewInt(35_000_000),
		},
		"baseCurrencyAtomicResolution is greater than 10^6": {
			baseQuantums:                 int256.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   7,
			priceExponent:                0,
			expectedQuoteQuantums:        int256.NewInt(350_000),
		},
		"baseCurrencyAtomicResolution is less than 10^6": {
			baseQuantums:                 int256.NewInt(5_000_000),
			baseCurrencyAtomicResolution: -4,
			priceValue:                   7,
			priceExponent:                1,
			expectedQuoteQuantums:        int256.NewInt(35_000_000_000),
		},
		"Calculation rounds down": {
			baseQuantums:                 int256.NewInt(9),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   1,
			priceExponent:                1,
			expectedQuoteQuantums:        int256.NewInt(0),
		},
		"Negative calculation rounds up": {
			baseQuantums:                 int256.NewInt(-9),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   1,
			priceExponent:                1,
			expectedQuoteQuantums:        int256.NewInt(0),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			quoteQuantums := lib.BaseToQuoteQuantums(
				tc.baseQuantums,
				tc.baseCurrencyAtomicResolution,
				tc.priceValue,
				tc.priceExponent,
			)
			if tc.expectedQuoteQuantums.Cmp(quoteQuantums) != 0 {
				t.Fatalf(
					"%s: expectedQuoteQuantums: %s, quoteQuantums: %s",
					name,
					tc.expectedQuoteQuantums.String(),
					quoteQuantums.String(),
				)
			}
		})
	}
}

func TestQuoteToBaseQuantums(t *testing.T) {
	tests := map[string]struct {
		quoteQuantums                *int256.Int
		baseCurrencyAtomicResolution int32
		priceValue                   uint64
		priceExponent                int32
		expectedBaseQuantums         *int256.Int
	}{
		"Converts from base to quote quantums": {
			quoteQuantums:                int256.NewInt(10),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			expectedBaseQuantums:         int256.NewInt(1),
		},
		"Correctly converts negative value": {
			quoteQuantums:                int256.NewInt(-1000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			expectedBaseQuantums:         int256.NewInt(-100),
		},
		"priceExponent is negative": {
			quoteQuantums:                int256.NewInt(350),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                -5,
			expectedBaseQuantums:         int256.NewInt(5_000_000),
		},
		"priceExponent is zero": {
			quoteQuantums:                int256.NewInt(35_000_000),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   7,
			priceExponent:                0,
			expectedBaseQuantums:         int256.NewInt(5_000_000),
		},
		"realistic values: 1 BTC at $29001": {
			quoteQuantums:                int256.NewInt(29_001_000_000), // $29_001
			baseCurrencyAtomicResolution: -10,
			priceValue:                   2_900_100_000,
			priceExponent:                -5,
			expectedBaseQuantums:         int256.NewInt(10_000_000_000),
		},
		"realistic values: 25.123 BTC at $29001": {
			quoteQuantums:                int256.NewInt(728_592_123_000), // $728_592.123
			baseCurrencyAtomicResolution: -10,
			priceValue:                   2_900_100_000,
			priceExponent:                -5,
			expectedBaseQuantums:         int256.NewInt(251_230_000_000),
		},
		"baseCurrencyAtomicResolution is greater than 10^6": {
			quoteQuantums:                int256.NewInt(350_000),
			baseCurrencyAtomicResolution: -8,
			priceValue:                   7,
			priceExponent:                0,
			expectedBaseQuantums:         int256.NewInt(5_000_000),
		},
		"baseCurrencyAtomicResolution is less than 10^6": {
			quoteQuantums:                int256.NewInt(35_000_000_000),
			baseCurrencyAtomicResolution: -4,
			priceValue:                   7,
			priceExponent:                1,
			expectedBaseQuantums:         int256.NewInt(5_000_000),
		},
		"Calculation rounds down": {
			quoteQuantums:                int256.NewInt(99),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			expectedBaseQuantums:         int256.NewInt(9),
		},
		"Negative calculation rounds up": {
			quoteQuantums:                int256.NewInt(-99),
			baseCurrencyAtomicResolution: -6,
			priceValue:                   1,
			priceExponent:                1,
			expectedBaseQuantums:         int256.NewInt(-9),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			baseQuantums := lib.QuoteToBaseQuantums(
				tc.quoteQuantums,
				tc.baseCurrencyAtomicResolution,
				tc.priceValue,
				tc.priceExponent,
			)
			if tc.expectedBaseQuantums.Cmp(baseQuantums) != 0 {
				t.Fatalf(
					"%s: expectedBaseQuantums: %s, baseQuantums: %s",
					name,
					tc.expectedBaseQuantums.String(),
					baseQuantums.String(),
				)
			}
		})
	}
}
