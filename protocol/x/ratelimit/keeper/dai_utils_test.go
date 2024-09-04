package keeper

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDivideAmountBySDaiDecimals(t *testing.T) {
	tests := map[string]struct {
		x              *big.Int
		expectedResult *big.Int
	}{
		"Divide 0.": {
			x:              ConvertStringToBigIntWithPanicOnErr("0"),
			expectedResult: big.NewInt(0),
		},
		"Divide positive even amount.": {
			x:              ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000"),
			expectedResult: big.NewInt(1),
		},
		"Divide negative.": {
			x:              ConvertStringToBigIntWithPanicOnErr("-1000000000000000000000000000"),
			expectedResult: big.NewInt(-1),
		},
		"Divide positive uneven amount.": {
			x:              ConvertStringToBigIntWithPanicOnErr("1234567890123456789123456789"),
			expectedResult: big.NewInt(1),
		},
		"Divide negative uneven amount.": {
			x:              ConvertStringToBigIntWithPanicOnErr("-1234567890123456789123456789"),
			expectedResult: big.NewInt(-2),
		},
		"Divide large positive even amount.": {
			x:              ConvertStringToBigIntWithPanicOnErr("1000000000000000000000000000000000000000"),
			expectedResult: big.NewInt(1000000000000),
		},
		"Divide large negative even amount.": {
			x:              ConvertStringToBigIntWithPanicOnErr("-1000000000000000000000000000000000000000"),
			expectedResult: big.NewInt(-1000000000000),
		},
		"Divide large positive uneven amount.": {
			x:              ConvertStringToBigIntWithPanicOnErr("1234567890123456789123456789123456789123"),
			expectedResult: big.NewInt(1234567890123),
		},
		"Divide large negatibe uneven amount.": {
			x:              ConvertStringToBigIntWithPanicOnErr("-1234567890123456789123456789123456789123"),
			expectedResult: big.NewInt(-1234567890124),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotResult := divideAmountBySDaiDecimals(tc.x)
			require.Equal(
				t,
				0,
				tc.expectedResult.Cmp(gotResult),
				"divideAmountBySDaiDecimals value does not match the expected value. Expected Result %v. Got result %v.",
				tc.expectedResult,
				gotResult,
			)
		})
	}
}

func TestGetTenScaledBySDaiDecimals(t *testing.T) {
	expectedInt, ok := big.NewInt(0).SetString("1000000000000000000000000000", 10)
	if !ok {
		panic("Could not set up test")
	}
	require.Equal(t, expectedInt, getTenScaledBySDaiDecimals())
}

func TestDivideAndRoundUp_Success(t *testing.T) {
	tests := map[string]struct {
		x              *big.Int
		y              *big.Int
		expectedResult *big.Int
	}{
		"Divide positive number by positive number: Larger number divided evenly by smaller number.": {
			x:              big.NewInt(100),
			y:              big.NewInt(5),
			expectedResult: big.NewInt(20),
		},
		"Divide positive number by another positive number: Larger number divided unevenly by smaller number.": {
			x:              big.NewInt(100),
			y:              big.NewInt(3),
			expectedResult: big.NewInt(34),
		},
		"Divide positive number by positive number: Smaller number divided by larger number with result closer to larger whole number.": {
			x:              big.NewInt(5),
			y:              big.NewInt(6),
			expectedResult: big.NewInt(1),
		},
		"Divide positive number by positive number: Smaller number divided by larger number with result closer to smaller whole number.": {
			x:              big.NewInt(5),
			y:              big.NewInt(100),
			expectedResult: big.NewInt(1),
		},
		"Divide positive number by positive number: Divide by itself.": {
			x:              big.NewInt(100),
			y:              big.NewInt(100),
			expectedResult: big.NewInt(1),
		},
		"Divide positive number by positive number: Divide by one.": {
			x:              big.NewInt(100),
			y:              big.NewInt(1),
			expectedResult: big.NewInt(100),
		},
		"Divide positive number by positive number: Divide two big integers.": {
			x:              big.NewInt(1000000000000),
			y:              big.NewInt(987654321),
			expectedResult: big.NewInt(1013),
		},
		"Divide 0 by positive number.": {
			x:              big.NewInt(0),
			y:              big.NewInt(987654321),
			expectedResult: big.NewInt(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotResult, err := divideAndRoundUp(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, gotResult, "DivideAndRoundUp value does not match the expected value")
			require.Equal(t, err, nil, "Error should have been nil on success, but got non-nil.")
		})
	}
}

func TestDivideAndRoundUp_Failure(t *testing.T) {
	tests := map[string]struct {
		x              *big.Int
		y              *big.Int
		expectedResult *big.Int
		expectedErr    error
	}{
		"Divide positive number by 0.": {
			x:              big.NewInt(10000000),
			y:              big.NewInt(0),
			expectedResult: nil,
			expectedErr:    errors.New("division by zero"),
		},
		"Divide nil by 0.": {
			x:              nil,
			y:              big.NewInt(0),
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be nil"),
		},
		"Divide negative number by 0.": {
			x:              big.NewInt(-10000000),
			y:              big.NewInt(0),
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be negative"),
		},
		"One input is negative: x is negative.": {
			x:              big.NewInt(-10000000),
			y:              big.NewInt(10),
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be negative"),
		},
		"One input is negative: y is negative.": {
			x:              big.NewInt(10000000),
			y:              big.NewInt(-10),
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be negative"),
		},
		"Both input are negative.": {
			x:              big.NewInt(-20),
			y:              big.NewInt(-10),
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be negative"),
		},
		"One input is nil: x is nil.": {
			x:              nil,
			y:              big.NewInt(10),
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be nil"),
		},
		"One input is nil: y is nil.": {
			x:              big.NewInt(10),
			y:              nil,
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be nil"),
		},
		"Both inputs are nil.": {
			x:              nil,
			y:              nil,
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be nil"),
		},
		"x is nil, y is negative.": {
			x:              nil,
			y:              big.NewInt(-10),
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be nil"),
		},
		"y is nil, x is negative.": {
			x:              big.NewInt(-10),
			y:              nil,
			expectedResult: nil,
			expectedErr:    errors.New("input values cannot be nil"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotResult, err := divideAndRoundUp(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, gotResult, "Expected nil value on failure, but got non-nil.")
			require.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestConvertStringToBigInt(t *testing.T) {
	tests := map[string]struct {
		x              string
		expectedResult *big.Int
		expectErr      bool
	}{
		"Zero.": {
			x:              "0",
			expectedResult: big.NewInt(0),
			expectErr:      false,
		},
		"Basic positive example.": {
			x:              "21",
			expectedResult: big.NewInt(21),
			expectErr:      false,
		},
		"Basic negative example.": {
			x:              "-21",
			expectedResult: big.NewInt(-21),
			expectErr:      false,
		},
		"Large even positive example.": {
			x: "10000000000000000000000000000000000000",
			expectedResult: func() *big.Int {
				result, ok := new(big.Int).SetString("10000000000000000000000000000000000000", 10)
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large even negative example.": {
			x: "-10000000000000000000000000000000000000",
			expectedResult: func() *big.Int {
				result, ok := new(big.Int).SetString("-10000000000000000000000000000000000000", 10)
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large uneven positive example.": {
			x: "123456789123456789123456789123456789",
			expectedResult: func() *big.Int {
				result, ok := new(big.Int).SetString("123456789123456789123456789123456789", 10)
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large uneven negative example.": {
			x: "-123456789123456789123456789123456789",
			expectedResult: func() *big.Int {
				result, ok := new(big.Int).SetString("-123456789123456789123456789123456789", 10)
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Fails: empty string": {
			x:              "",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: non-numeric characters": {
			x:              "123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: underscores": {
			x:              "1_000_000",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: leading whitespace": {
			x:              " 123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: trailing whitespace": {
			x:              " 123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: only a sign 1": {
			x:              "+",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: only a sign 2": {
			x:              "-",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: multiple signs": {
			x:              "-+123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: sign not at strart": {
			x:              "123+",
			expectedResult: nil,
			expectErr:      true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotResult, err := ConvertStringToBigInt(tc.x)
			require.Equal(t, 0, tc.expectedResult.Cmp(gotResult), "Expected Result: %v. Got %v.", tc.expectedResult, gotResult)
			if tc.expectErr {
				require.ErrorContains(t, err, "Unable to convert the sDAI conversion rate to a big int")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConvertStringToBigIntWithPanicOnErr(t *testing.T) {
	tests := map[string]struct {
		x              string
		expectedResult *big.Int
		expectErr      bool
	}{
		"Zero.": {
			x:              "0",
			expectedResult: big.NewInt(0),
			expectErr:      false,
		},
		"Basic positive example.": {
			x:              "21",
			expectedResult: big.NewInt(21),
			expectErr:      false,
		},
		"Basic negative example.": {
			x:              "-21",
			expectedResult: big.NewInt(-21),
			expectErr:      false,
		},
		"Large even positive example.": {
			x: "10000000000000000000000000000000000000",
			expectedResult: func() *big.Int {
				result, ok := new(big.Int).SetString("10000000000000000000000000000000000000", 10)
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large even negative example.": {
			x: "-10000000000000000000000000000000000000",
			expectedResult: func() *big.Int {
				result, ok := new(big.Int).SetString("-10000000000000000000000000000000000000", 10)
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large uneven positive example.": {
			x: "123456789123456789123456789123456789",
			expectedResult: func() *big.Int {
				result, ok := new(big.Int).SetString("123456789123456789123456789123456789", 10)
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large uneven negative example.": {
			x: "-123456789123456789123456789123456789",
			expectedResult: func() *big.Int {
				result, ok := new(big.Int).SetString("-123456789123456789123456789123456789", 10)
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Fails: empty string": {
			x:              "",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: non-numeric characters": {
			x:              "123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: underscores": {
			x:              "1_000_000",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: leading whitespace": {
			x:              " 123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: trailing whitespace": {
			x:              " 123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: only a sign 1": {
			x:              "+",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: only a sign 2": {
			x:              "-",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: multiple signs": {
			x:              "-+123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: sign not at strart": {
			x:              "123+",
			expectedResult: nil,
			expectErr:      true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectErr {
				require.Panics(t, func() { ConvertStringToBigIntWithPanicOnErr(tc.x) })
			} else {
				gotResult := ConvertStringToBigIntWithPanicOnErr(tc.x)
				require.Equal(t, 0, tc.expectedResult.Cmp(gotResult), "Expected Result: %v. Got %v.", tc.expectedResult, gotResult)
			}
		})
	}
}

func TestConvertStringToBigRatWithPanicOnErr(t *testing.T) {
	tests := map[string]struct {
		x              string
		expectedResult *big.Rat
		expectErr      bool
	}{
		"Zero.": {
			x:              "0",
			expectedResult: big.NewRat(0, 1),
			expectErr:      false,
		},
		"Basic positive example.": {
			x:              "2",
			expectedResult: big.NewRat(2, 1),
			expectErr:      false,
		},
		"Basic negative example.": {
			x:              "-2",
			expectedResult: big.NewRat(-2, 1),
			expectErr:      false,
		},
		"Large even positive example.": {
			x: "10000000000000000000000000000000000000",
			expectedResult: func() *big.Rat {
				result, ok := new(big.Rat).SetString("10000000000000000000000000000000000000")
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large even negative example.": {
			x: "-10000000000000000000000000000000000000",
			expectedResult: func() *big.Rat {
				result, ok := new(big.Rat).SetString("-10000000000000000000000000000000000000")
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large uneven positive example.": {
			x: "10000000000000000000000000000000000000.1234566789",
			expectedResult: func() *big.Rat {
				result, ok := new(big.Rat).SetString("10000000000000000000000000000000000000.1234566789")
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Large uneven negative example.": {
			x: "-10000000000000000000000000000000000000.1234566789",
			expectedResult: func() *big.Rat {
				result, ok := new(big.Rat).SetString("-10000000000000000000000000000000000000.1234566789")
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Underscores": {
			x: "1_000_000",
			expectedResult: func() *big.Rat {
				result, ok := new(big.Rat).SetString("1000000")
				if !ok {
					panic("Failed to set up test")
				}
				return result
			}(),
			expectErr: false,
		},
		"Fails: empty string": {
			x:              "",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: non-numeric characters": {
			x:              "123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: leading whitespace": {
			x:              " 123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: trailing whitespace": {
			x:              " 123abc123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: only a sign 1": {
			x:              "+",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: only a sign 2": {
			x:              "-",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: multiple signs": {
			x:              "-+123",
			expectedResult: nil,
			expectErr:      true,
		},
		"Fails: sign not at strart": {
			x:              "123+",
			expectedResult: nil,
			expectErr:      true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectErr {
				require.Panics(t, func() { ConvertStringToBigRatWithPanicOnErr(tc.x) })
			} else {
				gotResult := ConvertStringToBigRatWithPanicOnErr(tc.x)
				require.Equal(t, 0, tc.expectedResult.Cmp(gotResult), "Expected Result: %v. Got %v.", tc.expectedResult, gotResult)
			}
		})
	}
}
