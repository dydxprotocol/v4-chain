package big_test

import (
	"math/big"
	"testing"

	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/stretchr/testify/require"
)

func TestInt64MulPow10(t *testing.T) {
	tests := map[string]struct {
		val            int64
		exponent       uint64
		expectedResult string
	}{
		"Regular value and exponent": {
			val:            215,
			exponent:       4,
			expectedResult: "2150000",
		},
		"Zero value": {
			val:            0,
			exponent:       3,
			expectedResult: "0",
		},
		"Zero exponent": {
			val:            2,
			exponent:       0,
			expectedResult: "2",
		},
		"(-2) * 1e3": {
			val:            -2,
			exponent:       3,
			expectedResult: "-2000",
		},
		"123456789 * 1e10": {
			val:            123456789,
			exponent:       10,
			expectedResult: "1234567890000000000",
		},
		"87654321 * 1e18": {
			val:            87654321,
			exponent:       18,
			expectedResult: "87654321000000000000000000",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := big_testutil.Int64MulPow10(tc.val, tc.exponent)
			bigExpected, valid := new(big.Int).SetString(tc.expectedResult, 10)
			require.True(t, valid)
			require.Equal(t, bigExpected, result)
		})
	}
}
