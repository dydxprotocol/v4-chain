package lib_test

import (
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"

	"github.com/stretchr/testify/require"
)

func BenchmarkBigI(b *testing.B) {
	var result *big.Int
	for i := 0; i < b.N; i++ {
		result = lib.BigI(int64(i))
	}
	require.Equal(b, result, result)
}

func BenchmarkBigU(b *testing.B) {
	var result *big.Int
	for i := 0; i < b.N; i++ {
		result = lib.BigU(uint32(i))
	}
	require.Equal(b, result, result)
}

func TestBigI(t *testing.T) {
	require.Equal(t, big.NewInt(-123), lib.BigI(int(-123)))
	require.Equal(t, big.NewInt(-123), lib.BigI(int32(-123)))
	require.Equal(t, big.NewInt(-123), lib.BigI(int64(-123)))
	require.Equal(t, big.NewInt(math.MaxInt64), lib.BigI(math.MaxInt64))
}

func TestBigU(t *testing.T) {
	require.Equal(t, big.NewInt(123), lib.BigU(uint(123)))
	require.Equal(t, big.NewInt(123), lib.BigU(uint32(123)))
	require.Equal(t, big.NewInt(123), lib.BigU(uint64(123)))
	require.Equal(t, new(big.Int).SetUint64(math.MaxUint64), lib.BigU(uint64(math.MaxUint64)))
}

func BenchmarkBigMulPpm_RoundDown(b *testing.B) {
	val := big.NewInt(543_211)
	ppm := big.NewInt(876_543)
	var result *big.Int
	for i := 0; i < b.N; i++ {
		result = lib.BigMulPpm(val, ppm, false)
	}
	require.Equal(b, big.NewInt(476147), result)
}

func BenchmarkBigMulPpm_RoundUp(b *testing.B) {
	val := big.NewInt(543_211)
	ppm := big.NewInt(876_543)
	var result *big.Int
	for i := 0; i < b.N; i++ {
		result = lib.BigMulPpm(val, ppm, true)
	}
	require.Equal(b, big.NewInt(476148), result)
}

func TestBigMulPpm(t *testing.T) {
	tests := map[string]struct {
		val            *big.Int
		ppm            *big.Int
		roundUp        bool
		expectedResult *big.Int
	}{
		"Positive round down": {
			val:            big.NewInt(543_211),
			ppm:            big.NewInt(876_543),
			roundUp:        false,
			expectedResult: big.NewInt(476147),
		},
		"Negative round down": {
			val:            big.NewInt(-543_211),
			ppm:            big.NewInt(876_543),
			roundUp:        false,
			expectedResult: big.NewInt(-476148),
		},
		"Positive round up": {
			val:            big.NewInt(543_211),
			ppm:            big.NewInt(876_543),
			roundUp:        true,
			expectedResult: big.NewInt(476148),
		},
		"Negative round up": {
			val:            big.NewInt(-543_211),
			ppm:            big.NewInt(876_543),
			roundUp:        true,
			expectedResult: big.NewInt(-476147),
		},
		"Zero val": {
			val:            big.NewInt(0),
			ppm:            big.NewInt(876_543),
			roundUp:        true,
			expectedResult: big.NewInt(0),
		},
		"Zero ppm": {
			val:            big.NewInt(543_211),
			ppm:            big.NewInt(0),
			roundUp:        true,
			expectedResult: big.NewInt(0),
		},
		"Zero val and ppm": {
			val:            big.NewInt(0),
			ppm:            big.NewInt(0),
			roundUp:        true,
			expectedResult: big.NewInt(0),
		},
		"Negative val": {
			val:            big.NewInt(-543_211),
			ppm:            big.NewInt(876_543),
			roundUp:        true,
			expectedResult: big.NewInt(-476147),
		},
		"Negative ppm": {
			val:            big.NewInt(543_211),
			ppm:            big.NewInt(-876_543),
			roundUp:        true,
			expectedResult: big.NewInt(-476147),
		},
		"Negative val and ppm": {
			val:            big.NewInt(-543_211),
			ppm:            big.NewInt(-876_543),
			roundUp:        true,
			expectedResult: big.NewInt(476148),
		},
		"Greater than max int64": {
			val:            big_testutil.MustFirst(new(big.Int).SetString("1000000000000000000000000", 10)),
			ppm:            big.NewInt(10_000),
			roundUp:        true,
			expectedResult: big_testutil.MustFirst(new(big.Int).SetString("10000000000000000000000", 10)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigMulPpm(tc.val, tc.ppm, tc.roundUp)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigPow10(t *testing.T) {
	tests := map[string]struct {
		exponent        int64
		expectedValue   *big.Int
		expectedInverse bool
	}{
		"0":   {0, big.NewInt(1), false},
		"1":   {1, big.NewInt(10), false},
		"2":   {2, big.NewInt(100), false},
		"3":   {3, big.NewInt(1000), false},
		"4":   {4, big.NewInt(10000), false},
		"5":   {5, big.NewInt(100000), false},
		"20":  {20, big_testutil.MustFirst(new(big.Int).SetString("100000000000000000000", 10)), false},
		"-1":  {-1, big.NewInt(10), true},
		"-2":  {-2, big.NewInt(100), true},
		"-3":  {-3, big.NewInt(1000), true},
		"-4":  {-4, big.NewInt(10000), true},
		"-5":  {-5, big.NewInt(100000), true},
		"-20": {-20, big_testutil.MustFirst(new(big.Int).SetString("100000000000000000000", 10)), true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, inverse := lib.BigPow10(tc.exponent)
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedInverse, inverse)
		})
	}
}

func TestBigPow10AllValuesInMemo(t *testing.T) {
	exponentString := "1"
	for i := 0; i < 100; i++ {
		expected, ok := new(big.Int).SetString(exponentString, 10)

		require.True(t, ok)
		result, _ := lib.BigPow10(i)
		require.Equal(t, expected, result)

		exponentString = exponentString + "0"
	}
}

func TestBigPow10AllValueNotInMemo(t *testing.T) {
	exponentString := "1" + strings.Repeat("0", 110)
	expected, ok := new(big.Int).SetString(exponentString, 10)
	require.True(t, ok)
	result, _ := lib.BigPow10(110)
	require.Equal(t, expected, result)
}

func TestBigIntMulPpm(t *testing.T) {
	tests := map[string]struct {
		input          *big.Int
		ppm            uint32
		expectedResult *big.Int
	}{
		"Ppm of 5": {
			input:          big.NewInt(1_000_000),
			ppm:            5,
			expectedResult: big.NewInt(5),
		},
		"Ppm of 10^6": {
			input:          big.NewInt(1_000_000),
			ppm:            1_000_000,
			expectedResult: big.NewInt(1_000_000),
		},
		"Ppm of 0": {
			input:          big.NewInt(1_000_000),
			ppm:            0,
			expectedResult: big.NewInt(0),
		},
		"Ppm over 10^6": {
			input:          big.NewInt(1_000_000),
			ppm:            1_000_000_000,
			expectedResult: big.NewInt(1_000_000_000),
		},
		"Ppm of max uint32": {
			input:          big.NewInt(1_000_000),
			ppm:            math.MaxUint32,
			expectedResult: big.NewInt(4_294_967_295),
		},
		"Positive rounding towards negative infinity": {
			input:          big.NewInt(3),
			ppm:            500_000,
			expectedResult: big.NewInt(1), // 3 * .5 = 1.5, rounds down to 1
		},
		"Negative rounding towards negative infinity": {
			input:          big.NewInt(-3),
			ppm:            500_000,
			expectedResult: big.NewInt(-2), // -3 * .5 = -1.5, rounds down to -2
		},
		"0 input": {
			input:          big.NewInt(0),
			ppm:            1_000_000,
			expectedResult: big.NewInt(0),
		},
		"Input greater than max uint": {
			input:          big_testutil.MustFirst(new(big.Int).SetString("1000000000000000000000000", 10)),
			ppm:            10_000,
			expectedResult: big_testutil.MustFirst(new(big.Int).SetString("10000000000000000000000", 10)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigIntMulPpm(tc.input, tc.ppm)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigMin(t *testing.T) {
	tests := map[string]struct {
		a              *big.Int
		b              *big.Int
		expectedResult *big.Int
	}{
		"a is smaller than b": {
			a:              big.NewInt(5),
			b:              big.NewInt(6),
			expectedResult: big.NewInt(5),
		},
		"b is smaller than a": {
			a:              big.NewInt(7),
			b:              big.NewInt(4),
			expectedResult: big.NewInt(4),
		},
		"a is equal to b": {
			a:              big.NewInt(8),
			b:              big.NewInt(8),
			expectedResult: big.NewInt(8),
		},
		"a and b are negative, a is less than b": {
			a:              big.NewInt(-8),
			b:              big.NewInt(-7),
			expectedResult: big.NewInt(-8),
		},
		"a and b are negative, b is less than a": {
			a:              big.NewInt(-9),
			b:              big.NewInt(-10),
			expectedResult: big.NewInt(-10),
		},
		"a is positive, b is negative, and abs(a) is less than abs(b)": {
			a:              big.NewInt(4),
			b:              big.NewInt(-7),
			expectedResult: big.NewInt(-7),
		},
		"a is positive, b is negative, and abs(a) is greater than abs(b)": {
			a:              big.NewInt(7),
			b:              big.NewInt(-4),
			expectedResult: big.NewInt(-4),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigMin(tc.a, tc.b)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigRatMin(t *testing.T) {
	tests := map[string]struct {
		a        *big.Rat
		b        *big.Rat
		expected *big.Rat
	}{
		"a is smaller than b": {
			a:        big.NewRat(5, 2),
			b:        big.NewRat(6, 2),
			expected: big.NewRat(5, 2),
		},
		"b is smaller than a": {
			a:        big.NewRat(7, 1),
			b:        big.NewRat(4, 1),
			expected: big.NewRat(4, 1),
		},
		"a is equal to b": {
			a:        big.NewRat(8, 7),
			b:        big.NewRat(8, 7),
			expected: big.NewRat(8, 7),
		},
		"a and b are negative, a is less than b": {
			a:        big.NewRat(-8, 3),
			b:        big.NewRat(-7, 3),
			expected: big.NewRat(-8, 3),
		},
		"a and b are negative, b is less than a": {
			a:        big.NewRat(-9, 5),
			b:        big.NewRat(-10, 5),
			expected: big.NewRat(-10, 5),
		},
		"a is positive, b is negative, and abs(a) is less than abs(b)": {
			a:        big.NewRat(4, 3),
			b:        big.NewRat(-7, 2),
			expected: big.NewRat(-7, 2),
		},
		"a is positive, b is negative, and abs(a) is greater than abs(b)": {
			a:        big.NewRat(7, 2),
			b:        big.NewRat(-4, 3),
			expected: big.NewRat(-4, 3),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigRatMin(tc.a, tc.b)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestBigMax(t *testing.T) {
	tests := map[string]struct {
		a              *big.Int
		b              *big.Int
		expectedResult *big.Int
	}{
		"a is smaller than b": {
			a:              big.NewInt(5),
			b:              big.NewInt(6),
			expectedResult: big.NewInt(6),
		},
		"b is smaller than a": {
			a:              big.NewInt(7),
			b:              big.NewInt(4),
			expectedResult: big.NewInt(7),
		},
		"a is equal to b": {
			a:              big.NewInt(8),
			b:              big.NewInt(8),
			expectedResult: big.NewInt(8),
		},
		"a and b are negative, a is less than b": {
			a:              big.NewInt(-8),
			b:              big.NewInt(-7),
			expectedResult: big.NewInt(-7),
		},
		"a and b are negative, b is less than a": {
			a:              big.NewInt(-9),
			b:              big.NewInt(-10),
			expectedResult: big.NewInt(-9),
		},
		"a is positive, b is negative, and abs(a) is less than abs(b)": {
			a:              big.NewInt(4),
			b:              big.NewInt(-7),
			expectedResult: big.NewInt(4),
		},
		"a is positive, b is negative, and abs(a) is greater than abs(b)": {
			a:              big.NewInt(7),
			b:              big.NewInt(-4),
			expectedResult: big.NewInt(7),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigMax(tc.a, tc.b)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigRatMulPpm(t *testing.T) {
	tests := map[string]struct {
		input          *big.Rat
		ppm            uint32
		expectedResult *big.Rat
	}{
		"Ppm of 5": {
			input:          big.NewRat(int64(lib.OneMillion), 1),
			ppm:            5,
			expectedResult: big.NewRat(5, 1),
		},
		"Ppm of 10^6": {
			input:          big.NewRat(int64(lib.OneMillion), 1),
			ppm:            lib.OneMillion,
			expectedResult: big.NewRat(int64(lib.OneMillion), 1),
		},
		"Ppm of 0": {
			input:          big.NewRat(int64(lib.OneMillion), 1),
			ppm:            0,
			expectedResult: big.NewRat(0, 1),
		},
		"Ppm over 10^6": {
			input:          big.NewRat(int64(lib.OneMillion), 1),
			ppm:            lib.OneMillion * 1000,
			expectedResult: big.NewRat(int64(lib.OneMillion*1000), 1),
		},
		"Ppm of max uint32": {
			input:          big.NewRat(int64(lib.OneMillion), 1),
			ppm:            math.MaxUint32,
			expectedResult: big.NewRat(4_294_967_295, 1),
		},
		"0 input": {
			input:          big.NewRat(0, 1),
			ppm:            lib.OneMillion,
			expectedResult: big.NewRat(0, 1),
		},
		"Input greater than max uint": {
			input:          big_testutil.MustFirst(new(big.Rat).SetString("1000000000000000000000000")),
			ppm:            10_000,
			expectedResult: big_testutil.MustFirst(new(big.Rat).SetString("10000000000000000000000")),
		},
		"Positive with a common divisor": {
			input:          big.NewRat(3, 1),
			ppm:            500_000,
			expectedResult: big.NewRat(3, 2),
		},
		"Negative with a common divisor": {
			input:          big.NewRat(-3, 1),
			ppm:            500_000,
			expectedResult: big.NewRat(-3, 2),
		},
		"Positive with no common divisor": {
			input:          big.NewRat(117, 61),
			ppm:            419,
			expectedResult: big.NewRat(49023, 61000000),
		},
		"Negative with no common divisor": {
			input:          big.NewRat(-117, 61),
			ppm:            419,
			expectedResult: big.NewRat(-49023, 61000000),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigRatMulPpm(tc.input, tc.ppm)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func BenchmarkBigRatMulPpm(b *testing.B) {
	input := big.NewRat(1_000_000, 1_234_567)
	ppm := uint32(5)
	var result *big.Rat
	for i := 0; i < b.N; i++ {
		result = lib.BigRatMulPpm(input, ppm)
	}
	require.Equal(b, big.NewRat(5, 1_234_567), result)
}

func TestBigRatClamp(t *testing.T) {
	tests := map[string]struct {
		input          *big.Rat
		lower          *big.Rat
		upper          *big.Rat
		expectedResult *big.Rat
	}{
		"Input is returned when input is between lower and upper bound": {
			input:          big.NewRat(2, 1),
			lower:          big.NewRat(1, 1),
			upper:          big.NewRat(3, 1),
			expectedResult: big.NewRat(2, 1),
		},
		"Input is returned when input is at lower bound": {
			input:          big.NewRat(2, 1),
			lower:          big.NewRat(2, 1),
			upper:          big.NewRat(3, 1),
			expectedResult: big.NewRat(2, 1),
		},
		"Input is returned when input is at upper bound": {
			input:          big.NewRat(2, 1),
			lower:          big.NewRat(1, 1),
			upper:          big.NewRat(2, 1),
			expectedResult: big.NewRat(2, 1),
		},
		"Lower is returned when input is below lower bound": {
			input:          big.NewRat(-1, 1),
			lower:          big.NewRat(0, 1),
			upper:          big.NewRat(5, 1),
			expectedResult: big.NewRat(0, 1),
		},
		"Upper is returned when input is above upper bound": {
			input:          big.NewRat(100, 1),
			lower:          big.NewRat(1, 1),
			upper:          big.NewRat(5, 1),
			expectedResult: big.NewRat(5, 1),
		},
		"Lower is returned when lower > upper && n < lower": {
			input:          big.NewRat(4, 1),
			lower:          big.NewRat(5, 1),
			upper:          big.NewRat(3, 1),
			expectedResult: big.NewRat(5, 1),
		},
		"Upper is returned when lower > upper && n >= lower": {
			input:          big.NewRat(4, 1),
			lower:          big.NewRat(4, 1),
			upper:          big.NewRat(3, 1),
			expectedResult: big.NewRat(3, 1),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigRatClamp(tc.input, tc.lower, tc.upper)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigIntClamp(t *testing.T) {
	tests := map[string]struct {
		input          *big.Int
		lower          *big.Int
		upper          *big.Int
		expectedResult *big.Int
	}{
		"Input is returned when input is between lower and upper bound": {
			input:          big.NewInt(2),
			lower:          big.NewInt(-1),
			upper:          big.NewInt(3),
			expectedResult: big.NewInt(2),
		},
		"Input is returned when input is at lower bound": {
			input:          big.NewInt(-2),
			lower:          big.NewInt(-2),
			upper:          big.NewInt(3),
			expectedResult: big.NewInt(-2),
		},
		"Input is returned when input is at upper bound": {
			input:          big.NewInt(4),
			lower:          big.NewInt(-3),
			upper:          big.NewInt(4),
			expectedResult: big.NewInt(4),
		},
		"Lower is returned when input is below lower bound": {
			input:          big.NewInt(-1),
			lower:          big.NewInt(0),
			upper:          big.NewInt(5),
			expectedResult: big.NewInt(0),
		},
		"Upper is returned when input is above upper bound": {
			input:          big.NewInt(100),
			lower:          big.NewInt(-2),
			upper:          big.NewInt(5),
			expectedResult: big.NewInt(5),
		},
		"Lower is returned when lower > upper && n < lower": {
			input:          big.NewInt(1),
			lower:          big.NewInt(2),
			upper:          big.NewInt(-4),
			expectedResult: big.NewInt(2),
		},
		"Upper is returned when lower > upper && n >= lower": {
			input:          big.NewInt(4),
			lower:          big.NewInt(4),
			upper:          big.NewInt(-2),
			expectedResult: big.NewInt(-2),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigIntClamp(tc.input, tc.lower, tc.upper)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func BenchmarkBigDivCeil(b *testing.B) {
	numerator := big.NewInt(10)
	denominator := big.NewInt(3)
	var result *big.Int
	for i := 0; i < b.N; i++ {
		result = lib.BigDivCeil(numerator, denominator)
	}
	require.Equal(b, big.NewInt(4), result)
}

func TestBigDivCeil(t *testing.T) {
	tests := map[string]struct {
		numerator      *big.Int
		denominator    *big.Int
		expectedResult *big.Int
	}{
		"Divides evenly": {
			numerator:      big.NewInt(10),
			denominator:    big.NewInt(5),
			expectedResult: big.NewInt(2),
		},
		"Doesn't divide evenly": {
			numerator:      big.NewInt(10),
			denominator:    big.NewInt(3),
			expectedResult: big.NewInt(4),
		},
		"Negative numerator": {
			numerator:      big.NewInt(-10),
			denominator:    big.NewInt(3),
			expectedResult: big.NewInt(-3),
		},
		"Negative numerator 2": {
			numerator:      big.NewInt(-1),
			denominator:    big.NewInt(2),
			expectedResult: big.NewInt(0),
		},
		"Negative denominator": {
			numerator:      big.NewInt(10),
			denominator:    big.NewInt(-3),
			expectedResult: big.NewInt(-3),
		},
		"Negative denominator 2": {
			numerator:      big.NewInt(1),
			denominator:    big.NewInt(-2),
			expectedResult: big.NewInt(0),
		},
		"Negative numerator and denominator": {
			numerator:      big.NewInt(-10),
			denominator:    big.NewInt(-3),
			expectedResult: big.NewInt(4),
		},
		"Negative numerator and denominator 2": {
			numerator:      big.NewInt(-1),
			denominator:    big.NewInt(-2),
			expectedResult: big.NewInt(1),
		},
		"Zero numerator": {
			numerator:      big.NewInt(0),
			denominator:    big.NewInt(3),
			expectedResult: big.NewInt(0),
		},
		"Zero denominator": {
			numerator:      big.NewInt(10),
			denominator:    big.NewInt(0),
			expectedResult: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Panics if the expected result is nil
			if tc.expectedResult == nil {
				require.Panics(t, func() {
					lib.BigDivCeil(tc.numerator, tc.denominator)
				})
				return
			}
			// Otherwise test the result
			result := lib.BigDivCeil(tc.numerator, tc.denominator)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigRatRound(t *testing.T) {
	tests := map[string]struct {
		input          *big.Rat
		roundUp        bool
		expectedResult *big.Int
	}{
		"Input is unrounded if it is zero and we round up": {
			input:          big.NewRat(0, 1),
			roundUp:        true,
			expectedResult: big.NewInt(0),
		},
		"Input is unrounded if it is zero and we round down": {
			input:          big.NewRat(0, 1),
			roundUp:        false,
			expectedResult: big.NewInt(0),
		},
		"Input is unrounded if it is an int and we round up": {
			input:          big.NewRat(2, 1),
			roundUp:        true,
			expectedResult: big.NewInt(2),
		},
		"Input is unrounded if it is an int and we round down": {
			input:          big.NewRat(2, 1),
			roundUp:        false,
			expectedResult: big.NewInt(2),
		},
		"Input is unrounded if it isn't normalized, it is an int and we round up": {
			input:          big.NewRat(21, 3),
			roundUp:        true,
			expectedResult: big.NewInt(7),
		},
		"Input is unrounded if it isn't normalized, it is an int and we round down": {
			input:          big.NewRat(21, 3),
			roundUp:        false,
			expectedResult: big.NewInt(7),
		},
		"Input is rounded up if we round up": {
			input:          big.NewRat(5, 4),
			roundUp:        true,
			expectedResult: big.NewInt(2),
		},
		"Input is rounded up if it isn't normalized and we round up": {
			input:          big.NewRat(10, 4),
			roundUp:        true,
			expectedResult: big.NewInt(3),
		},
		"Input is rounded down if rational number isn't normalized and we round down": {
			input:          big.NewRat(10, 4),
			roundUp:        false,
			expectedResult: big.NewInt(2),
		},
		"Input is rounded down if we round down": {
			input:          big.NewRat(5, 4),
			roundUp:        false,
			expectedResult: big.NewInt(1),
		},
		"Input is rounded down if input is negative and we round down": {
			input:          big.NewRat(-22, 7),
			roundUp:        false,
			expectedResult: big.NewInt(-4),
		},
		"Input is rounded up if input is negative and we round up": {
			input:          big.NewRat(-22, 7),
			roundUp:        true,
			expectedResult: big.NewInt(-3),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigRatRound(tc.input, tc.roundUp)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigIntRoundToMultiple(t *testing.T) {
	tests := map[string]struct {
		n              *big.Int
		multiple       *big.Int
		roundUp        bool
		expectedResult *big.Int
		shouldPanic    bool
	}{
		"Input is unrounded if n % multiple == 0 and we round up": {
			n:              big.NewInt(5),
			multiple:       big.NewInt(5),
			roundUp:        true,
			expectedResult: big.NewInt(5),
		},
		"Input is unrounded if n % multiple == 0 and we round down": {
			n:              big.NewInt(-5),
			multiple:       big.NewInt(5),
			roundUp:        false,
			expectedResult: big.NewInt(-5),
		},
		"Input is rounded up if n % multiple != 0 and we round up": {
			n:              big.NewInt(7),
			multiple:       big.NewInt(3),
			roundUp:        true,
			expectedResult: big.NewInt(9),
		},
		"Input is rounded down if n % multiple != 0 and we round down": {
			n:              big.NewInt(7),
			multiple:       big.NewInt(3),
			roundUp:        false,
			expectedResult: big.NewInt(6),
		},
		"Input is rounded up if n is negative, n % multiple != 0 and we round up": {
			n:              big.NewInt(-7),
			multiple:       big.NewInt(3),
			roundUp:        true,
			expectedResult: big.NewInt(-6),
		},
		"Input is rounded down if n is negative, n % multiple != 0 and we round down": {
			n:              big.NewInt(-7),
			multiple:       big.NewInt(3),
			roundUp:        false,
			expectedResult: big.NewInt(-9),
		},
		"Panics if multiple is zero": {
			n:        big.NewInt(-7),
			multiple: big.NewInt(0),
			roundUp:  false,

			shouldPanic: true,
		},
		"Panics if multiple is negative": {
			n:        big.NewInt(7),
			multiple: big.NewInt(-1),
			roundUp:  true,

			shouldPanic: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldPanic {
				require.Panics(t, func() {
					//nolint:errcheck
					lib.BigIntRoundToMultiple(
						tc.n,
						tc.multiple,
						tc.roundUp,
					)
				})
				return
			}
			result := lib.BigIntRoundToMultiple(
				tc.n,
				tc.multiple,
				tc.roundUp,
			)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigFloatMaxUint64(t *testing.T) {
	require.Equal(t, new(big.Float).SetUint64(math.MaxUint64), lib.BigFloatMaxUint64())
}

func TestBigRat0(t *testing.T) {
	require.Equal(t, big.NewRat(0, 1), lib.BigRat0())
}

func TestBigRat1(t *testing.T) {
	require.Equal(t, big.NewRat(1, 1), lib.BigRat1())
}

func TestBigUint64Clamp(t *testing.T) {
	tests := map[string]struct {
		input          *big.Int
		lower          uint64
		upper          uint64
		expectedResult uint64
	}{
		"Input is returned when input is between lower and upper bound": {
			input:          big.NewInt(2),
			lower:          1,
			upper:          3,
			expectedResult: 2,
		},
		"Input is returned when input is at lower bound": {
			input:          big.NewInt(2),
			lower:          2,
			upper:          3,
			expectedResult: 2,
		},
		"Input is returned when input is at upper bound": {
			input:          big.NewInt(2),
			lower:          1,
			upper:          2,
			expectedResult: 2,
		},
		"Lower is returned when input is below lower bound": {
			input:          big.NewInt(0),
			lower:          1,
			upper:          5,
			expectedResult: 1,
		},
		"Upper is returned when input is above upper bound": {
			input:          big.NewInt(100),
			lower:          1,
			upper:          5,
			expectedResult: 5,
		},
		"Lower is returned when lower > upper && n < lower": {
			input:          big.NewInt(4),
			lower:          5,
			upper:          3,
			expectedResult: 5,
		},
		"Upper is returned when lower > upper && n >= lower": {
			input:          big.NewInt(4),
			lower:          4,
			upper:          3,
			expectedResult: 3,
		},
		"Upper is returned when input is above max uint64": {
			input: big_testutil.MustFirst(
				new(big.Int).SetString("100000000000000000000", 10),
			),
			lower:          1,
			upper:          math.MaxUint64,
			expectedResult: math.MaxUint64,
		},
		"Lower is returned when input is negative": {
			input: big_testutil.MustFirst(
				new(big.Int).SetString("-100000000000000000000", 10),
			),
			lower:          1,
			upper:          20,
			expectedResult: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigUint64Clamp(tc.input, tc.lower, tc.upper)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigInt32Clamp(t *testing.T) {
	tests := map[string]struct {
		input          *big.Int
		lower          int32
		upper          int32
		expectedResult int32
	}{
		"Input is returned when input is between lower and upper bound": {
			input:          big.NewInt(11),
			lower:          -100,
			upper:          123,
			expectedResult: 11,
		},
		"Input is returned when input is at lower bound": {
			input:          big.NewInt(-100),
			lower:          -100,
			upper:          123,
			expectedResult: -100,
		},
		"Input is returned when input is at upper bound": {
			input:          big.NewInt(123),
			lower:          -100,
			upper:          123,
			expectedResult: 123,
		},
		"Lower is returned when input is below lower bound": {
			input:          big.NewInt(-100),
			lower:          -10,
			upper:          5,
			expectedResult: -10,
		},
		"Upper is returned when input is above upper bound": {
			input:          big.NewInt(100),
			lower:          -5,
			upper:          5,
			expectedResult: 5,
		},
		"Lower is returned when lower > upper && n < lower": {
			input:          big.NewInt(4),
			lower:          5,
			upper:          -3,
			expectedResult: 5,
		},
		"Upper is returned when lower > upper && n >= lower": {
			input:          big.NewInt(4),
			lower:          -4,
			upper:          3,
			expectedResult: 3,
		},
		"Upper is returned when input is above max int32": {
			input: big_testutil.MustFirst(
				new(big.Int).SetString("100000000000000000000", 10),
			),
			lower:          -100,
			upper:          100,
			expectedResult: 100,
		},
		"Lower is returned when input is below min int32": {
			input: big_testutil.MustFirst(
				new(big.Int).SetString("-100000000000000000000", 10),
			),
			lower:          -100,
			upper:          100,
			expectedResult: -100,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigInt32Clamp(tc.input, tc.lower, tc.upper)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestMustConvertBigIntToInt32(t *testing.T) {
	tests := map[string]struct {
		input          *big.Int
		expectedResult int32
		shouldPanic    bool
	}{
		"Returns true when input is between MinInt32 and MaxInt32": {
			input:          big.NewInt(11),
			expectedResult: int32(11),
		},
		"Returns true when input is MinInt32": {
			input:          big.NewInt(math.MinInt32),
			expectedResult: int32(math.MinInt32),
		},
		"Returns true when input is MaxInt32": {
			input:          big.NewInt(math.MaxInt32),
			expectedResult: int32(math.MaxInt32),
		},
		"Returns false when input is larger than MaxInt32": {
			input:       big.NewInt(math.MaxInt32 + 1),
			shouldPanic: true,
		},
		"Returns false when input is smaller than MinInt32": {
			input:       big.NewInt(math.MinInt32 - 1),
			shouldPanic: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldPanic {
				require.Panics(t, func() {
					//nolint:errcheck
					lib.MustConvertBigIntToInt32(tc.input)
				})
				return
			}
			result := lib.MustConvertBigIntToInt32(tc.input)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
