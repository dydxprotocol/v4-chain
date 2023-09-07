package lib_test

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	zero = uint64(0)
)

var (
	overflowError            = errors.New("value overflows uint64")
	underflowError           = errors.New("value underflows uint64")
	maxUint64PlusOneBigFloat = new(big.Float).Add(new(big.Float).SetUint64(1), lib.BigFloatMaxUint64())
)

func TestMustConvertIntegerToUint32_Int8(t *testing.T) {
	tests := map[string]struct {
		// parameters
		value int8

		// expectations
		expectedPanic bool
		expected      uint32
	}{
		"Convert 0 successfully": {
			value:    0,
			expected: 0,
		},
		"Convert max int8 successfully": {
			value:    math.MaxInt8,
			expected: math.MaxInt8,
		},
		"Convert -1 panic": {
			value:         -1,
			expectedPanic: true,
		},
		"Convert min int8 panics": {
			value:         math.MinInt8,
			expectedPanic: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic {
				require.Panics(t, func() { lib.MustConvertIntegerToUint32(tc.value) })
			} else {
				require.Equal(t, tc.expected, lib.MustConvertIntegerToUint32(tc.value))
			}
		})
	}
}

func TestMustConvertIntegerToUint32_Int32(t *testing.T) {
	tests := map[string]struct {
		// parameters
		value int32

		// expectations
		expectedPanic bool
		expected      uint32
	}{
		"Convert 0 successfully": {
			value:    0,
			expected: 0,
		},
		"Convert max int32 successfully": {
			value:    math.MaxInt32,
			expected: math.MaxInt32,
		},
		"Convert -1 panic": {
			value:         -1,
			expectedPanic: true,
		},
		"Convert min int32 panics": {
			value:         math.MinInt32,
			expectedPanic: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic {
				require.Panics(t, func() { lib.MustConvertIntegerToUint32(tc.value) })
			} else {
				require.Equal(t, tc.expected, lib.MustConvertIntegerToUint32(tc.value))
			}
		})
	}
}

func TestMustConvertIntegerToUint32_Int64(t *testing.T) {
	tests := map[string]struct {
		// parameters
		value int64

		// expectations
		expectedPanic bool
		expected      uint32
	}{
		"Convert 0 successfully": {
			value:    0,
			expected: 0,
		},
		"Convert max uint32 successfully": {
			value:    math.MaxUint32,
			expected: math.MaxUint32,
		},
		"Convert max uint32 + 1 panics": {
			value:         math.MaxUint32 + 1,
			expectedPanic: true,
		},
		"Convert max int64 panics": {
			value:         math.MaxInt64,
			expectedPanic: true,
		},
		"Convert -1 panics": {
			value:         -1,
			expectedPanic: true,
		},
		"Convert min int64 panics": {
			value:         math.MinInt64,
			expectedPanic: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic {
				require.Panics(t, func() { lib.MustConvertIntegerToUint32(tc.value) })
			} else {
				require.Equal(t, tc.expected, lib.MustConvertIntegerToUint32(tc.value))
			}
		})
	}
}

func TestMustConvertIntegerToUint32_Uint64(t *testing.T) {
	tests := map[string]struct {
		// parameters
		value uint64

		// expectations
		expectedPanic bool
		expected      uint32
	}{
		"Convert 0 successfully": {
			value:    0,
			expected: 0,
		},
		"Convert max uint32 successfully": {
			value:    math.MaxUint32,
			expected: math.MaxUint32,
		},
		"Convert max uint32 + 1 panics": {
			value:         math.MaxUint32 + 1,
			expectedPanic: true,
		},
		"Convert max uint64 panics": {
			value:         math.MaxUint64,
			expectedPanic: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanic {
				require.Panics(t, func() { lib.MustConvertIntegerToUint32(tc.value) })
			} else {
				require.Equal(t, tc.expected, lib.MustConvertIntegerToUint32(tc.value))
			}
		})
	}
}

func TestConvertBigFloatToUint64(t *testing.T) {
	tests := map[string]struct {
		// parameters
		bigFloatValue *big.Float

		// expectations
		expectedError  error
		expectedUint64 uint64
	}{
		"Convert 0 successfully": {
			bigFloatValue:  new(big.Float).SetInt64(0),
			expectedError:  nil,
			expectedUint64: uint64(0),
		},
		"Convert max uint64 successfully": {
			bigFloatValue:  lib.BigFloatMaxUint64(),
			expectedError:  nil,
			expectedUint64: math.MaxUint64,
		},
		"Convert successfully": {
			bigFloatValue:  new(big.Float).SetInt64(200),
			expectedError:  nil,
			expectedUint64: uint64(200),
		},
		"Convert successfully and rounds down 0.2": {
			bigFloatValue:  new(big.Float).SetFloat64(200.2),
			expectedError:  nil,
			expectedUint64: uint64(200),
		},
		"Convert successfully and rounds down 0.7": {
			bigFloatValue:  new(big.Float).SetFloat64(200.7),
			expectedError:  nil,
			expectedUint64: uint64(200),
		},
		"Convert and overflow": {
			bigFloatValue:  maxUint64PlusOneBigFloat,
			expectedError:  overflowError,
			expectedUint64: zero,
		},
		"Convert and underflow": {
			bigFloatValue:  new(big.Float).SetInt64(-10),
			expectedError:  underflowError,
			expectedUint64: zero,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := lib.ConvertBigFloatToUint64(tc.bigFloatValue)
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}

			require.Equal(t, tc.expectedUint64, result)
		})
	}
}

func TestConvertStringSliceToBigFloatSlice(t *testing.T) {
	tests := map[string]struct {
		// parameters
		stringSlice []string

		// expectations
		expectedError         error
		expectedBigFloatSlice []*big.Float
	}{
		"Convert successfully": {
			stringSlice:   []string{"100"},
			expectedError: nil,
			expectedBigFloatSlice: []*big.Float{
				new(big.Float).SetUint64(100),
			},
		},
		"Convert empty string returns an error": {
			stringSlice:   []string{""},
			expectedError: fmt.Errorf("invalid, value is not a number: %v", ""),
		},
		"Multiple values and one returns an error": {
			stringSlice: []string{
				"100.0001",
				"300.02",
				"5",
				"50000.001.001", // Invalid
			},
			expectedError:         fmt.Errorf("invalid, value is not a number: %v", "50000.001.001"),
			expectedBigFloatSlice: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := lib.ConvertStringSliceToBigFloatSlice(tc.stringSlice)
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}

			require.Equal(t, tc.expectedBigFloatSlice, result)
		})
	}
}

func TestConvertStringSliceToBigFloatSlice_MultipleValues(t *testing.T) {
	stringValues := []string{
		"-100.0001",
		"300.02",
		"5",
		"50000",
	}

	result, err := lib.ConvertStringSliceToBigFloatSlice(stringValues)

	expectedResult := make([]*big.Float, 0, len(stringValues))
	for _, val := range stringValues {
		v, success := new(big.Float).SetString(val)
		require.True(t, success)
		expectedResult = append(expectedResult, v)
	}

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestConvertBigFloatSliceToUint64Slice(t *testing.T) {
	tests := map[string]struct {
		// parameters
		bigFloatSlice []*big.Float

		// expectations
		expectedError       error
		expectedUint64Slice []uint64
	}{
		"Convert successfully": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(100.0001),
			},
			expectedUint64Slice: []uint64{
				uint64(100),
			},
		},
		"Convert successfully for 0 and max value": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(0),
				lib.BigFloatMaxUint64(),
			},
			expectedUint64Slice: []uint64{
				uint64(0),
				math.MaxUint64,
			},
		},
		"Convert successfully with multiple values": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(100.0001),
				new(big.Float).SetFloat64(300.02),
				new(big.Float).SetFloat64(5),
				new(big.Float).SetFloat64(50000),
				lib.BigFloatMaxUint64(),
			},
			expectedUint64Slice: []uint64{
				uint64(100),
				uint64(300),
				uint64(5),
				uint64(50000),
				math.MaxUint64,
			},
		},
		"Multiple values and one errors": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(100.0001),
				new(big.Float).SetFloat64(300.02),
				new(big.Float).SetFloat64(5),
				new(big.Float).SetFloat64(-1), // Invalid
				lib.BigFloatMaxUint64(),
			},
			expectedError: underflowError,
		},
		"value overflows": {
			bigFloatSlice: []*big.Float{
				maxUint64PlusOneBigFloat,
			},
			expectedError: overflowError,
		},
		"value underflows": {
			bigFloatSlice: []*big.Float{
				new(big.Float).SetFloat64(-100),
			},
			expectedError: underflowError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := lib.ConvertBigFloatSliceToUint64Slice(tc.bigFloatSlice)
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}

			require.Equal(t, tc.expectedUint64Slice, result)
		})
	}
}
