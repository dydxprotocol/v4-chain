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

func TestBigPow10(t *testing.T) {
	tests := map[string]struct {
		exponent       uint64
		expectedResult *big.Int
	}{
		"Regular exponent": {
			exponent:       3,
			expectedResult: new(big.Int).SetUint64(1000),
		},
		"Zero exponent": {
			exponent:       0,
			expectedResult: new(big.Int).SetUint64(1),
		},
		"One exponent": {
			exponent:       1,
			expectedResult: new(big.Int).SetUint64(10),
		},
		"Power of 2": {
			exponent:       8,
			expectedResult: new(big.Int).SetUint64(100_000_000),
		},
		"Non-power of 2": {
			exponent:       6,
			expectedResult: new(big.Int).SetUint64(1_000_000),
		},
		"Greater than max uint64": {
			exponent:       20,
			expectedResult: big_testutil.MustFirst(new(big.Int).SetString("100000000000000000000", 10)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigPow10(tc.exponent)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigMulPow10(t *testing.T) {
	tests := map[string]struct {
		val            *big.Int
		exponent       int32
		expectedResult *big.Rat
	}{
		"exponent = 2": {
			val:            new(big.Int).SetUint64(12345678),
			exponent:       2,
			expectedResult: big.NewRat(1234567800, 1),
		},
		"exponent = 10": {
			val:            new(big.Int).SetUint64(12345678),
			exponent:       10,
			expectedResult: big.NewRat(123456780000000000, 1),
		},
		"exponent = 0": {
			val:            new(big.Int).SetUint64(12345678),
			exponent:       0,
			expectedResult: big.NewRat(12345678, 1),
		},
		"exponent = -1": {
			val:            new(big.Int).SetUint64(12345678),
			exponent:       -1,
			expectedResult: big.NewRat(12345678, 10),
		},
		"exponent = -3": {
			val:            new(big.Int).SetUint64(12345678),
			exponent:       -3,
			expectedResult: big.NewRat(12345678, 1000),
		},
		"exponent = -8": {
			val:            new(big.Int).SetUint64(12345678),
			exponent:       -8,
			expectedResult: big.NewRat(12345678, 100000000),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigMulPow10(tc.val, tc.exponent)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestRatPow10(t *testing.T) {
	tests := map[string]struct {
		exponent       int32
		expectedResult *big.Rat
	}{
		"Positive exponent": {
			exponent:       3,
			expectedResult: new(big.Rat).SetUint64(1000),
		},
		"Negative exponent": {
			exponent:       -3,
			expectedResult: new(big.Rat).SetFrac64(1, 1000),
		},
		"Zero exponent": {
			exponent:       0,
			expectedResult: new(big.Rat).SetUint64(1),
		},
		"One exponent": {
			exponent:       1,
			expectedResult: new(big.Rat).SetUint64(10),
		},
		"Negative one exponent": {
			exponent:       -1,
			expectedResult: new(big.Rat).SetFrac64(1, 10),
		},
		"Power of 2": {
			exponent:       8,
			expectedResult: new(big.Rat).SetUint64(100_000_000),
		},
		"Negative power of 2": {
			exponent:       -8,
			expectedResult: new(big.Rat).SetFrac64(1, 100_000_000),
		},
		"Non-power of 2": {
			exponent:       6,
			expectedResult: new(big.Rat).SetUint64(1_000_000),
		},
		"Negative non-power of 2": {
			exponent:       -6,
			expectedResult: new(big.Rat).SetFrac64(1, 1_000_000),
		},
		"Greater than max uint64": {
			exponent:       20,
			expectedResult: big_testutil.MustFirst(new(big.Rat).SetString("100000000000000000000")),
		},
		"Denom greater than max uint64": {
			exponent: -20,
			expectedResult: new(big.Rat).SetFrac(
				new(big.Int).SetInt64(1),
				big_testutil.MustFirst(
					new(big.Int).SetString("100000000000000000000", 10),
				),
			),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.RatPow10(tc.exponent)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestBigPow10AllValuesInMemo(t *testing.T) {
	exponentString := "1"
	for i := 0; i < 100; i++ {
		bigValue, ok := new(big.Int).SetString(exponentString, 0)

		require.True(t, ok)
		require.Equal(t, lib.BigPow10(uint64(i)), bigValue)

		exponentString = exponentString + "0"
	}
}

func TestBigPow10AllValueNotInMemo(t *testing.T) {
	exponentString := "1" + strings.Repeat("0", 110)
	bigValue, ok := new(big.Int).SetString(exponentString, 0)
	require.True(t, ok)
	require.Equal(t, lib.BigPow10(uint64(110)), bigValue)
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
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigRatMulPpm(tc.input, tc.ppm)
			require.Equal(t, tc.expectedResult, result)
		})
	}
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

func TestBigInt0(t *testing.T) {
	require.Equal(t, big.NewInt(0), lib.BigInt0())
}

func TestBigFloat0(t *testing.T) {
	require.Equal(t, big.NewFloat(0), lib.BigFloat0())
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

func TestBigRatRoundToNearestMultiple(t *testing.T) {
	tests := map[string]struct {
		value          *big.Rat
		base           uint32
		up             bool
		expectedResult uint64
	}{
		"Round 5 down to a multiple of 2": {
			value:          big.NewRat(5, 1),
			base:           2,
			up:             false,
			expectedResult: 4,
		},
		"Round 5 up to a multiple of 2": {
			value:          big.NewRat(5, 1),
			base:           2,
			up:             true,
			expectedResult: 6,
		},
		"Round 7 down to a multiple of 14": {
			value:          big.NewRat(7, 1),
			base:           14,
			up:             false,
			expectedResult: 0,
		},
		"Round 7 up to a multiple of 14": {
			value:          big.NewRat(7, 1),
			base:           14,
			up:             true,
			expectedResult: 14,
		},
		"Round 123 down to a multiple of 123": {
			value:          big.NewRat(123, 1),
			base:           123,
			up:             false,
			expectedResult: 123,
		},
		"Round 123 up to a multiple of 123": {
			value:          big.NewRat(123, 1),
			base:           123,
			up:             true,
			expectedResult: 123,
		},
		"Round 100/6 down to a multiple of 3": {
			value:          big.NewRat(100, 6),
			base:           3,
			up:             false,
			expectedResult: 15,
		},
		"Round 100/6 up to a multiple of 3": {
			value:          big.NewRat(100, 6),
			base:           3,
			up:             true,
			expectedResult: 18,
		},
		"Round 7/2 down to a multiple of 1": {
			value:          big.NewRat(7, 2),
			base:           1,
			up:             false,
			expectedResult: 3,
		},
		"Round 7/2 up to a multiple of 1": {
			value:          big.NewRat(7, 2),
			base:           1,
			up:             true,
			expectedResult: 4,
		},
		"Round 10 down to a multiple of 0": {
			value:          big.NewRat(10, 1),
			base:           0,
			up:             false,
			expectedResult: 0,
		},
		"Round 10 up to a multiple of 0": {
			value:          big.NewRat(10, 1),
			base:           0,
			up:             true,
			expectedResult: 0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.BigRatRoundToNearestMultiple(
				tc.value,
				tc.base,
				tc.up,
			)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
