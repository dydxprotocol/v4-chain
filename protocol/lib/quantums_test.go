package lib_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/stretchr/testify/require"
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

func BenchmarkBaseToQuoteQuantums(b *testing.B) {
	value := big.NewInt(1234123412341)
	baseCurrencyAtomicResolution := int32(-8)
	priceValue := uint64(18446744073709551610)
	priceExponent1 := int32(-10)
	priceExponent2 := int32(10)
	var result1 *big.Int
	var result2 *big.Int

	for i := 0; i < b.N; i++ {
		result1 = lib.BaseToQuoteQuantums(
			value,
			baseCurrencyAtomicResolution,
			priceValue,
			priceExponent1,
		)
		result2 = lib.BaseToQuoteQuantums(
			value,
			baseCurrencyAtomicResolution,
			priceValue,
			priceExponent2,
		)
	}

	expected1, _ := new(big.Int).SetString("22765558742827551059", 10)
	require.Equal(b, expected1, result1)
	expected2, _ := new(big.Int).SetString("2276555874282755105905825041901000000000", 10)
	require.Equal(b, expected2, result2)
}

func BenchmarkQuoteToBaseQuantums(b *testing.B) {
	value, _ := new(big.Int).SetString("18446744073709551610", 10)
	baseCurrencyAtomicResolution := int32(-8)
	priceValue := uint64(1234123412341)
	priceExponent1 := int32(-10)
	priceExponent2 := int32(6)
	var result1 *big.Int
	var result2 *big.Int

	for i := 0; i < b.N; i++ {
		result1 = lib.QuoteToBaseQuantums(
			value,
			baseCurrencyAtomicResolution,
			priceValue,
			priceExponent1,
		)
		result2 = lib.QuoteToBaseQuantums(
			value,
			baseCurrencyAtomicResolution,
			priceValue,
			priceExponent2,
		)
	}

	expected1, _ := new(big.Int).SetString("14947244245790664347", 10)
	require.Equal(b, expected1, result1)
	expected2, _ := new(big.Int).SetString("1494", 10)
	require.Equal(b, expected2, result2)
}

func TestBigRatRoundToMultiple(t *testing.T) {
	tests := map[string]struct {
		value          *big.Rat
		multiple       *big.Int
		roundUp        bool
		expectedResult *big.Int
	}{
		// The fractional part must survive to the rounding decision. In these cases the
		// floor of the value already sits on a multiple, so discarding it before rounding
		// leaves the result a full multiple below the value it was asked to round up past.
		"Rounds up when the value's floor is already a multiple": {
			value:          big.NewRat(2001, 2), // 1000.5
			multiple:       big.NewInt(10),
			roundUp:        true,
			expectedResult: big.NewInt(1010),
		},
		"Rounds up to the next integer when the multiple is one": {
			value:          big.NewRat(21, 2), // 10.5
			multiple:       big.NewInt(1),
			roundUp:        true,
			expectedResult: big.NewInt(11),
		},
		"Rounds a value below the multiple up to the multiple": {
			value:          big.NewRat(1, 2), // 0.5
			multiple:       big.NewInt(10),
			roundUp:        true,
			expectedResult: big.NewInt(10),
		},
		"Rounds up when the value's floor is not a multiple": {
			value:          big.NewRat(2011, 2), // 1005.5
			multiple:       big.NewInt(10),
			roundUp:        true,
			expectedResult: big.NewInt(1010),
		},
		"Rounds down": {
			value:          big.NewRat(2011, 2), // 1005.5
			multiple:       big.NewInt(10),
			roundUp:        false,
			expectedResult: big.NewInt(1000),
		},
		"Leaves an exact multiple unchanged when rounding up": {
			value:          big.NewRat(1000, 1),
			multiple:       big.NewInt(10),
			roundUp:        true,
			expectedResult: big.NewInt(1000),
		},
		"Leaves an exact multiple unchanged when rounding down": {
			value:          big.NewRat(1000, 1),
			multiple:       big.NewInt(10),
			roundUp:        false,
			expectedResult: big.NewInt(1000),
		},
		"Rounds an integer value up to the next multiple": {
			value:          big.NewRat(7, 1),
			multiple:       big.NewInt(3),
			roundUp:        true,
			expectedResult: big.NewInt(9),
		},
		"Rounds an integer value down to the previous multiple": {
			value:          big.NewRat(7, 1),
			multiple:       big.NewInt(3),
			roundUp:        false,
			expectedResult: big.NewInt(6),
		},
		// Rounding up means towards positive infinity, so a negative value rounds
		// towards zero.
		"Rounds a negative value up towards zero": {
			value:          big.NewRat(-2001, 2), // -1000.5
			multiple:       big.NewInt(10),
			roundUp:        true,
			expectedResult: big.NewInt(-1000),
		},
		"Rounds a negative value down away from zero": {
			value:          big.NewRat(-2001, 2), // -1000.5
			multiple:       big.NewInt(10),
			roundUp:        false,
			expectedResult: big.NewInt(-1010),
		},
		"Rounds a negative value down when the multiple is one": {
			value:          big.NewRat(-21, 2), // -10.5
			multiple:       big.NewInt(1),
			roundUp:        false,
			expectedResult: big.NewInt(-11),
		},
		"Zero is unchanged": {
			value:          big.NewRat(0, 1),
			multiple:       big.NewInt(10),
			roundUp:        true,
			expectedResult: big.NewInt(0),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigRatRoundToMultiple(tc.value, tc.multiple, tc.roundUp)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
