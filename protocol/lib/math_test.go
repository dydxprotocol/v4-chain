package lib_test

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestUint64LinearInterpolate(t *testing.T) {
	tests := map[string]struct {
		v0          uint64
		v1          uint64
		cPpm        uint32
		expected    uint64
		expectedErr error
	}{
		"lerp near uint64 max": {
			v0:       math.MaxUint64 - 1,
			v1:       math.MaxUint64,
			cPpm:     1_000_000,
			expected: math.MaxUint64,
		},
		"lerp near uint64 max with rounding": {
			v0:       math.MaxUint64 - 1,
			v1:       math.MaxUint64,
			cPpm:     500_000,
			expected: math.MaxUint64 - 1,
		},
		"lerp with same inputs": {
			v0:       math.MaxUint64,
			v1:       math.MaxUint64,
			cPpm:     1_000_000,
			expected: math.MaxUint64,
		},
		"lerp with << inputs": {
			v0:       2_000_000,
			v1:       3_000_000,
			cPpm:     300_000,
			expected: 2_300_000,
		},
		"lerp with << inputs, v1 < v0": {
			v0:       3_000_000,
			v1:       2_000_000,
			cPpm:     700_000,
			expected: 2_300_000,
		},
		"lerp with invalid inputs, cPpm > 1_000_000": {
			v0:          3_000_000,
			v1:          2_000_000,
			cPpm:        1_500_000,
			expectedErr: fmt.Errorf("uint64 interpolation requires 0 <= cPpm <= 1_000_000, but received cPpm value of 1500000"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := lib.Uint64LinearInterpolate(tc.v0, tc.v1, tc.cPpm)
			if tc.expectedErr == nil {
				require.Nil(t, err)
				require.Equal(t, tc.expected, result)
			} else {
				require.Zero(t, result)
				require.EqualError(t, tc.expectedErr, err.Error())
			}
		})
	}
}

func TestAddToUint32(t *testing.T) {
	tests := map[string]struct {
		a           int64
		b           uint32
		expected    int64
		expectedErr error
	}{
		"success: smallest possible output": {
			a:        math.MinInt64,
			b:        0,
			expected: math.MinInt64,
		},
		"a + b overflows int64: << b": {
			a:           math.MaxInt64,
			b:           1,
			expectedErr: fmt.Errorf("int64 overflow: %d + %d", math.MaxInt64, 1),
		},
		"a + b overflows int64: smallest possible a": {
			a:           math.MaxInt64 - math.MaxUint32 + 1,
			b:           math.MaxUint32,
			expectedErr: fmt.Errorf("int64 overflow: %d + %d", math.MaxInt64-math.MaxUint32+1, math.MaxUint32),
		},
		"success": {
			a:        math.MaxUint32,
			b:        1 << 20,
			expected: math.MaxUint32 + 1<<20,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual, actualError := lib.AddUint32(tc.a, tc.b)
			if tc.expectedErr == nil {
				require.Nil(t, actualError)
				require.Equal(t, tc.expected, actual)
			} else {
				require.EqualError(t, actualError, tc.expectedErr.Error())
				require.Zero(t, actual)
			}
		})
	}
}

func TestMustDivideUint32RoundUp(t *testing.T) {
	tests := map[string]struct {
		x              uint32
		y              uint32
		expectedResult uint32
	}{
		"y divides x": {
			x:              3600,
			y:              60,
			expectedResult: 60,
		},
		"y doesn't divide x, round up": {
			x:              3601,
			y:              60,
			expectedResult: 61,
		},
		"x == 0": {
			x:              0,
			y:              60,
			expectedResult: 0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.MustDivideUint32RoundUp(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGenericMinFloat64(t *testing.T) {
	tests := map[string]struct {
		x              float64
		y              float64
		expectedResult float64
	}{
		"x is smaller": {
			x:              float64(32.213),
			y:              float64(32.214),
			expectedResult: float64(32.213),
		},
		"y is smaller": {
			x:              float64(.00001),
			y:              float64(-9123.31241),
			expectedResult: float64(-9123.31241),
		},
		"x == y": {
			x:              float64(931.2322),
			y:              float64(931.2322),
			expectedResult: float64(931.2322),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.Min(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGenericMinUInt64(t *testing.T) {
	tests := map[string]struct {
		x              uint64
		y              uint64
		expectedResult uint64
	}{
		"x is smaller": {
			x:              uint64(32),
			y:              uint64(33),
			expectedResult: uint64(32),
		},
		"y is smaller": {
			x:              uint64(1),
			y:              uint64(0),
			expectedResult: uint64(0),
		},
		"x == y": {
			x:              uint64(931),
			y:              uint64(931),
			expectedResult: uint64(931),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.Min(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGenericMaxUInt64(t *testing.T) {
	tests := map[string]struct {
		x              uint64
		y              uint64
		expectedResult uint64
	}{
		"x is smaller": {
			x:              uint64(32),
			y:              uint64(33),
			expectedResult: uint64(33),
		},
		"y is smaller": {
			x:              uint64(1),
			y:              uint64(0),
			expectedResult: uint64(1),
		},
		"x == y": {
			x:              uint64(931),
			y:              uint64(931),
			expectedResult: uint64(931),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.Max(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGenericMaxFloat64(t *testing.T) {
	tests := map[string]struct {
		x              float64
		y              float64
		expectedResult float64
	}{
		"x is smaller": {
			x:              float64(32.213),
			y:              float64(32.214),
			expectedResult: float64(32.214),
		},
		"y is smaller": {
			x:              float64(.00001),
			y:              float64(-9123.31241),
			expectedResult: float64(.00001),
		},
		"x == y": {
			x:              float64(931.2322),
			y:              float64(931.2322),
			expectedResult: float64(931.2322),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.Max(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestInt64MulPpm(t *testing.T) {
	tests := map[string]struct {
		x              int64
		ppm            uint32
		expectedResult int64
		shouldPanic    bool
		expectedError  error
	}{
		"60 * 25% = 15": {
			x:              60,
			ppm:            250_000, // 25%
			expectedResult: 15,
		},
		"60 * 10% = 6": {
			x:              60,
			ppm:            100_000, // 10%
			expectedResult: 6,
		},
		"69 * 10% rounds down to 6 (round towards negative infinity)": {
			x:              69,
			ppm:            100_000, // 10%
			expectedResult: 6,
		},
		"-61 * 10% rounds down to -7 (round towards negative infinity)": {
			x:              -61,
			ppm:            100_000, // 10%
			expectedResult: -7,
		},
		"overflow causes panic": {
			x:             math.MaxInt64,
			ppm:           1_000_001,
			shouldPanic:   true,
			expectedError: fmt.Errorf("IntMulPpm (int = 9223372036854775807, ppm = 1000001) results in integer overflow"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldPanic {
				require.PanicsWithError(
					t,
					tc.expectedError.Error(),
					func() {
						lib.Int64MulPpm(tc.x, tc.ppm)
					},
				)
				return
			}

			result := lib.Int64MulPpm(tc.x, tc.ppm)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestUint64MulPpm(t *testing.T) {
	tests := map[string]struct {
		x              uint64
		ppm            uint32
		expectedResult uint64
		shouldPanic    bool
		expectedError  error
	}{
		"60 * 25% = 15": {
			x:              60,
			ppm:            250_000, // 25%
			expectedResult: 15,
		},
		"60 * 10% = 6": {
			x:              60,
			ppm:            100_000, // 10%
			expectedResult: 6,
		},
		"61 * 10% rounds down to 6": {
			x:              61,
			ppm:            100_000, // 10%
			expectedResult: 6,
		},
		"overflow causes panic": {
			x:             math.MaxUint64,
			ppm:           1_000_001,
			shouldPanic:   true,
			expectedError: fmt.Errorf("UintMulPpm (uint = 18446744073709551615, ppm = 1000001) results in integer overflow"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldPanic {
				require.PanicsWithError(
					t,
					tc.expectedError.Error(),
					func() {
						lib.Uint64MulPpm(tc.x, tc.ppm)
					},
				)
				return
			}

			result := lib.Uint64MulPpm(tc.x, tc.ppm)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestAvgInt32(t *testing.T) {
	tests := map[string]struct {
		nums           []int32
		expectedResult int32
		shouldPanic    bool
		expectedError  error
	}{
		"Array 1": {
			nums:           []int32{20, 20, 20, 20, 30},
			expectedResult: 22,
		},
		"Array 2": {
			nums:           []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			expectedResult: 6,
		},
		"Rounds to zero": {
			nums:           []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 13},
			expectedResult: 6,
		},
		"Rounds to zero for negative numbers": {
			nums:           []int32{-1, -2, -3, -4, -5, -6, -7, -8, -9, -10, -13},
			expectedResult: -6,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldPanic {
				require.PanicsWithError(
					t,
					tc.expectedError.Error(),
					func() {
						lib.AvgInt32(tc.nums)
					},
				)
				return
			}

			result := lib.AvgInt32(tc.nums)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestAbsInt32(t *testing.T) {
	tests := map[string]struct {
		num            int32
		expectedResult uint32
	}{
		"Negative number": {
			num:            -7,
			expectedResult: 7,
		},
		"Positive number": {
			num:            7,
			expectedResult: 7,
		},
		"Zero": {
			num:            0,
			expectedResult: 0,
		},
		"Min integer": {
			num:            math.MinInt32,
			expectedResult: 2_147_483_648,
		},
		"Max integer": {
			num:            math.MaxInt32,
			expectedResult: 2_147_483_647,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.AbsInt32(tc.num)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestAbsInt64(t *testing.T) {
	tests := map[string]struct {
		num            int64
		expectedResult uint64
	}{
		"Negative number": {
			num:            -7,
			expectedResult: 7,
		},
		"Positive number": {
			num:            7,
			expectedResult: 7,
		},
		"Zero": {
			num:            0,
			expectedResult: 0,
		},
		"Min int64": {
			num:            math.MinInt64,
			expectedResult: 9_223_372_036_854_775_808,
		},
		"Max int64": {
			num:            math.MaxInt64,
			expectedResult: 9_223_372_036_854_775_807,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.AbsInt64(tc.num)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestAbsDiffUint64(t *testing.T) {
	tests := map[string]struct {
		x              uint64
		y              uint64
		expectedResult uint64
	}{
		"x < y": {
			x:              10,
			y:              11,
			expectedResult: 1,
		},
		"x == y": {
			x:              10,
			y:              10,
			expectedResult: 0,
		},
		"x > y": {
			x:              11,
			y:              10,
			expectedResult: 1,
		},
		"Zero": {
			x:              0,
			y:              0,
			expectedResult: 0,
		},
		"Max uint64": {
			x:              math.MaxUint64,
			y:              0,
			expectedResult: math.MaxUint64,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := lib.AbsDiffUint64(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestChangeRateUint64(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		originalV uint64
		newV      uint64

		// Expected.
		expectedResult float32
		expectedError  error
	}{
		"Zero": {
			originalV:      1000,
			newV:           1000, // change = 0%
			expectedResult: 0.0,
		},
		"Positive": {
			originalV:      1,
			newV:           1001, // change = 100,000%
			expectedResult: 1000.0,
		},
		"Positive: truncated decimals": {
			originalV:      1_000_000_000,
			newV:           1_987_654_321, // change = 98.7654321%
			expectedResult: 0.9876543,     // float32 precision is 8 decimal, so trucates.
		},
		"Negative": {
			originalV:      10000,
			newV:           1, // change = -99.99%
			expectedResult: -0.9999,
		},
		"Negative: truncated decimals": {
			originalV:      1_000_000_000,
			newV:           123_456_789, // change = 87.6543211%
			expectedResult: -0.8765432,  // float32 precision is 8 decimal, so trucates.
		},
		"Invalid: divide by zero return error": {
			originalV:     0, // divide by zero is not allowed
			newV:          1_000_000_000,
			expectedError: errors.New("original value cannot be zero since we cannot divide by zero"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := lib.ChangeRateUint64(tc.originalV, tc.newV)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestMustGetMedian_Failure(t *testing.T) {
	require.PanicsWithError(t,
		"input cannot be empty",
		func() {
			lib.MustGetMedian([]int32{})
		},
	)
}

func TestMustGetMedian_Success(t *testing.T) {
	require.Equal(t,
		int32(5),
		lib.MustGetMedian([]int32{8, 1, -5, 100, -50, 59}),
	)
}

func TestMedian_Int32(t *testing.T) {
	tests := map[string]struct {
		input          []int32
		expectedResult int32
		expectedError  bool
	}{
		"Empty input causes error": {
			input:          []int32{},
			expectedResult: 0,
			expectedError:  true,
		},
		"Odd number input": {
			input:          []int32{2, 0, 1, 3, 4},
			expectedResult: 2,
			expectedError:  false,
		},
		"Even number input, middle numbers are both odd": {
			input:          []int32{5, 11, 1, 3, 12, 50}, // median is (5+11)/2=8
			expectedResult: 8,
			expectedError:  false,
		},
		"Even number input, middle numbers are odd and even": {
			input:          []int32{5, 12, 1, 3, 12, 50}, // median is (5+12)/2=8.5
			expectedResult: 9,
			expectedError:  false,
		},
		"Even number input, middle numbers are even and odd": {
			input:          []int32{6, 11, 1, 3, 12, 50}, // median is (6+11)/2=8.5
			expectedResult: 9,
			expectedError:  false,
		},
		"Even number input, middle numbers are both even": {
			input:          []int32{6, 12, 1, 3, 12, 50}, // median is (6+12)/2=9
			expectedResult: 9,
			expectedError:  false,
		},
		"All negative, even number input, middle numbers are both odd": {
			input:          []int32{-6, -12, -1, -3, -12, -50}, // median is (-6-12)/2=-9
			expectedResult: -9,
			expectedError:  false,
		},
		"All negative, even number input,  middle numbers are even and odd": {
			input:          []int32{-5, -12, -1, -3, -12, -50}, // median is (-12-5)/2=-8.5
			expectedResult: -9,
			expectedError:  false,
		},
		"All negative, even number input, middle numbers are odd and even": {
			input:          []int32{-6, -11, -1, -3, -12, -50}, // median is (-11-6)/2=8.5
			expectedResult: -9,
			expectedError:  false,
		},
		"All negative, even number input, middle numbers are both even": {
			input:          []int32{-6, -12, -1, -3, -12, -50}, // median is (-12-6)/2=-9
			expectedResult: -9,
			expectedError:  false,
		},
		"Mixed signs, odd number input": {
			input:          []int32{6, 12, 1, -12, -50},
			expectedResult: 1,
			expectedError:  false,
		},
		"Mixed signs, even number input, middle numbers both non-negative": {
			input:          []int32{-12, 20, 100, 1, 0, -50},
			expectedResult: 1, // median is (0+1)/2=1
			expectedError:  false,
		},
		"Mixed signs, odd number input, middle numbers both non-positive": {
			input:          []int32{-12, 20, 100, -9, 0, -50},
			expectedResult: -5, // median is (-9+0)/2=-5
			expectedError:  false,
		},
		"Mixed signs, odd number input, middle numbers different signs, both odd": {
			input:          []int32{-12, 20, 100, -9, 1, -50},
			expectedResult: -4, // median is (-9+1)/2=-4
			expectedError:  false,
		},
		"Mixed signs, odd number input, middle numbers different signs, both even": {
			input:          []int32{-12, 6, 100, -2, 20, -50},
			expectedResult: 2, // median is (-2+6)/2=2
			expectedError:  false,
		},
		"Mixed signs, odd number input, middle numbers different signs, even and odd, negative median": {
			input:          []int32{-12, 5, 100, -10, 20, -50},
			expectedResult: -3, // median is (-10+5)/2=-3
			expectedError:  false,
		},
		"Mixed signs, odd number input, middle numbers different signs, even and odd, positive median": {
			input:          []int32{-12, 5, 100, -4, 20, -50},
			expectedResult: 1, // median is (-4+5)/2=1
			expectedError:  false,
		},
		"Mixed signs, odd number input, middle numbers different signs, odd and even, negative median": {
			input:          []int32{-12, 2, 100, -9, 20, -50},
			expectedResult: -4, // median is (-9+2)/2=-4
			expectedError:  false,
		},
		"Mixed signs, odd number input, middle numbers different signs, odd and even, positive median": {
			input:          []int32{-12, 16, 100, -9, 20, -50},
			expectedResult: 4, // median is (-9+16)/2=4
			expectedError:  false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := lib.Median(tc.input)
			require.Equal(t, tc.expectedResult, result)
			if tc.expectedError {
				require.EqualError(t, err, "input cannot be empty")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMedian_Uint64(t *testing.T) {
	tests := map[string]struct {
		input          []uint64
		expectedResult uint64
		expectedError  bool
	}{
		"Empty input causes error": {
			input:          []uint64{},
			expectedResult: 0,
			expectedError:  true,
		},
		"Odd number input": {
			input:          []uint64{2, 0, 1, 3, 4},
			expectedResult: 2,
			expectedError:  false,
		},
		"Even number input, middle numbers are both odd": {
			input:          []uint64{5, 11, 1, 3, 12, 50}, // median is (5+11)/2=8
			expectedResult: 8,
			expectedError:  false,
		},
		"Even number input, middle numbers are odd and even": {
			input:          []uint64{5, 12, 1, 3, 12, 50}, // median is (5+12)/2=8.5
			expectedResult: 9,
			expectedError:  false,
		},
		"Even number input, middle numbers are even and odd": {
			input:          []uint64{6, 11, 1, 3, 12, 50}, // median is (6+11)/2=8.5
			expectedResult: 9,
			expectedError:  false,
		},
		"Even number input, middle numbers are both even": {
			input:          []uint64{6, 12, 1, 3, 12, 50}, // median is (6+12)/2=9
			expectedResult: 9,
			expectedError:  false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := lib.Median(tc.input)
			require.Equal(t, tc.expectedResult, result)
			if tc.expectedError {
				require.EqualError(t, err, "input cannot be empty")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
