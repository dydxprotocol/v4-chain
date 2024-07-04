package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	// "time"

	// errorsmod "cosmossdk.io/errors"
	// sdkmath "cosmossdk.io/math"

	// "github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	// testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	// big_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/big"
	// blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	// "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	// cometbfttypes "github.com/cometbft/cometbft/types"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	// banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
)

func TestDivideAndRoundUp_Success(t *testing.T) {
	tests := map[string]struct {
		x                	 *big.Int
		y				     *big.Int
		expectedResult		 *big.Int
	}{
		"Divide positive number by positive number: Larger number divided evenly by smaller number.": {
			x: big.NewInt(100),
			y: big.NewInt(5),
			expectedResult: big.NewInt(20), 
		},
		"Divide positive number by another positive number: Larger number divided unevenly by smaller number.": {
			x: big.NewInt(100),
			y: big.NewInt(3),
			expectedResult: big.NewInt(34), 
		},
		"Divide positive number by positive number: Smaller number divided by larger number with result closer to larger whole number.": {
			x: big.NewInt(5),
			y: big.NewInt(6),
			expectedResult: big.NewInt(1), 
		},
		"Divide positive number by positive number: Smaller number divided by larger number with result closer to smaller whole number.": {
			x: big.NewInt(5),
			y: big.NewInt(100),
			expectedResult: big.NewInt(1), 
		},
		"Divide positive number by positive number: Divide by itself.": {
			x: big.NewInt(100),
			y: big.NewInt(100),
			expectedResult: big.NewInt(1), 
		},
		"Divide positive number by positive number: Divide by one.": {
			x: big.NewInt(100),
			y: big.NewInt(1),
			expectedResult: big.NewInt(100), 
		},
		"Divide positive number by positive number: Divide two big integers.": {
			x: big.NewInt(1000000000000),
			y: big.NewInt(987654321),
			expectedResult: big.NewInt(1013), 
		},
		"Divide 0 by positive number.": {
			x: big.NewInt(0),
			y: big.NewInt(987654321),
			expectedResult: big.NewInt(0), 
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotResult, err := keeper.DivideAndRoundUp(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, gotResult, "DivideAndRoundUp value does not match the expected value")
			require.Equal(t, err, nil, "Error should have been nil on success, but got non-nil.")
		})
	}
}

func TestDivideAndRoundUp_Failure(t *testing.T) {
	tests := map[string]struct {
		x                	 *big.Int
		y				     *big.Int
		expectedResult		 *big.Int
		expectedErr			 error
	}{
		"Divide positive number by 0.": {
			x: big.NewInt(10000000),
			y: big.NewInt(0),
			expectedResult: nil,
			expectedErr: errors.New("division by zero"),
		},
		"Divide nil by 0.": {
			x: nil,
			y: big.NewInt(0),
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be nil"),
		},
		"Divide negative number by 0.": {
			x: big.NewInt(-10000000),
			y: big.NewInt(0),
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be negative"),
		},
		"One input is negative: x is negative.": {
			x: big.NewInt(-10000000),
			y: big.NewInt(10),
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be negative"),
		},
		"One input is negative: y is negative.": {
			x: big.NewInt(10000000),
			y: big.NewInt(-10),
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be negative"),
		},
		"Both input are negative.": {
			x: big.NewInt(-20),
			y: big.NewInt(-10),
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be negative"),
		},
		"One input is nil: x is nil.": {
			x: nil,
			y: big.NewInt(10),
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be nil"),
		},
		"One input is nil: y is nil.": {
			x: big.NewInt(10),
			y: nil,
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be nil"),
		},
		"Both inputs are nil.": {
			x: nil,
			y: nil,
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be nil"),
		},
		"x is nil, y is negative.": {
			x: nil,
			y: big.NewInt(-10),
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be nil"),
		},
		"y is nil, x is negative.": {
			x: big.NewInt(-10),
			y: nil,
			expectedResult: nil,
			expectedErr: errors.New("input values cannot be nil"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotResult, err := keeper.DivideAndRoundUp(tc.x, tc.y)
			require.Equal(t, tc.expectedResult, gotResult, "Expected nil value on failure, but got non-nil.")
			require.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}
