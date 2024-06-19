package lib_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestIsValidStateTransitionForUndercollateralizedSubaccount_ZeroMarginRequirements(t *testing.T) {
	tests := map[string]struct {
		oldNC  *big.Int
		oldIMR *big.Int
		oldMMR *big.Int
		newNC  *big.Int
		newMMR *big.Int

		expectedResult types.UpdateResult
	}{
		// Tests when current margin requirement is zero and margin requirement increases.
		"fails when MMR increases and TNC decreases - negative TNC": {
			oldNC:          big.NewInt(-1),
			oldIMR:         big.NewInt(0),
			oldMMR:         big.NewInt(0),
			newNC:          big.NewInt(-2),
			newMMR:         big.NewInt(1),
			expectedResult: types.StillUndercollateralized,
		},
		"fails when MMR increases and TNC stays the same - negative TNC": {
			oldNC:          big.NewInt(-1),
			oldIMR:         big.NewInt(0),
			oldMMR:         big.NewInt(0),
			newNC:          big.NewInt(-1),
			newMMR:         big.NewInt(1),
			expectedResult: types.StillUndercollateralized,
		},
		"fails when MMR increases and TNC increases - negative TNC": {
			oldNC:          big.NewInt(-1),
			oldIMR:         big.NewInt(0),
			oldMMR:         big.NewInt(0),
			newNC:          big.NewInt(100),
			newMMR:         big.NewInt(1),
			expectedResult: types.StillUndercollateralized,
		},
		// Tests when both margin requirements are zero.
		"fails when both new and old MMR are zero and TNC stays the same": {
			oldNC:          big.NewInt(-1),
			oldIMR:         big.NewInt(0),
			oldMMR:         big.NewInt(0),
			newNC:          big.NewInt(-1),
			newMMR:         big.NewInt(0),
			expectedResult: types.StillUndercollateralized,
		},
		"fails when both new and old MMR are zero and TNC decrease from negative to negative": {
			oldNC:          big.NewInt(-1),
			oldIMR:         big.NewInt(0),
			oldMMR:         big.NewInt(0),
			newNC:          big.NewInt(-2),
			newMMR:         big.NewInt(0),
			expectedResult: types.StillUndercollateralized,
		},
		"succeeds when both new and old MMR are zero and TNC increases": {
			oldNC:          big.NewInt(-2),
			oldIMR:         big.NewInt(0),
			oldMMR:         big.NewInt(0),
			newNC:          big.NewInt(-1),
			newMMR:         big.NewInt(0),
			expectedResult: types.Success,
		},
		// Tests when new margin requirement is zero.
		"fails when MMR decreased to zero, and TNC increases but is still negative": {
			oldNC:          big.NewInt(-2),
			oldIMR:         big.NewInt(1),
			oldMMR:         big.NewInt(1),
			newNC:          big.NewInt(-1),
			newMMR:         big.NewInt(0),
			expectedResult: types.StillUndercollateralized,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedResult,
				lib.IsValidStateTransitionForUndercollateralizedSubaccount(
					margin.Risk{
						NC:  tc.oldNC,
						IMR: tc.oldIMR,
						MMR: tc.oldMMR,
					},
					margin.Risk{
						NC:  tc.newNC,
						MMR: tc.newMMR,
					},
				),
			)
		})
	}
}
