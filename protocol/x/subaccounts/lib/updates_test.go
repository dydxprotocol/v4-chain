package lib_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestIsValidStateTransitionForUndercollateralizedSubaccount_ZeroMarginRequirements(t *testing.T) {
	tests := map[string]struct {
		bigCurNetCollateral     *big.Int
		bigCurInitialMargin     *big.Int
		bigCurMaintenanceMargin *big.Int
		bigNewNetCollateral     *big.Int
		bigNewMaintenanceMargin *big.Int

		expectedResult types.UpdateResult
	}{
		// Tests when current margin requirement is zero and margin requirement increases.
		"fails when MMR increases and TNC decreases - negative TNC": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-2),
			bigNewMaintenanceMargin: big.NewInt(1),
			expectedResult:          types.StillUndercollateralized,
		},
		"fails when MMR increases and TNC stays the same - negative TNC": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-1),
			bigNewMaintenanceMargin: big.NewInt(1),
			expectedResult:          types.StillUndercollateralized,
		},
		"fails when MMR increases and TNC increases - negative TNC": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(100),
			bigNewMaintenanceMargin: big.NewInt(1),
			expectedResult:          types.StillUndercollateralized,
		},
		// Tests when both margin requirements are zero.
		"fails when both new and old MMR are zero and TNC stays the same": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-1),
			bigNewMaintenanceMargin: big.NewInt(0),
			expectedResult:          types.StillUndercollateralized,
		},
		"fails when both new and old MMR are zero and TNC decrease from negative to negative": {
			bigCurNetCollateral:     big.NewInt(-1),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-2),
			bigNewMaintenanceMargin: big.NewInt(0),
			expectedResult:          types.StillUndercollateralized,
		},
		"succeeds when both new and old MMR are zero and TNC increases": {
			bigCurNetCollateral:     big.NewInt(-2),
			bigCurInitialMargin:     big.NewInt(0),
			bigCurMaintenanceMargin: big.NewInt(0),
			bigNewNetCollateral:     big.NewInt(-1),
			bigNewMaintenanceMargin: big.NewInt(0),
			expectedResult:          types.Success,
		},
		// Tests when new margin requirement is zero.
		"fails when MMR decreased to zero, and TNC increases but is still negative": {
			bigCurNetCollateral:     big.NewInt(-2),
			bigCurInitialMargin:     big.NewInt(1),
			bigCurMaintenanceMargin: big.NewInt(1),
			bigNewNetCollateral:     big.NewInt(-1),
			bigNewMaintenanceMargin: big.NewInt(0),
			expectedResult:          types.StillUndercollateralized,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expectedResult,
				lib.IsValidStateTransitionForUndercollateralizedSubaccount(
					tc.bigCurNetCollateral,
					tc.bigCurInitialMargin,
					tc.bigCurMaintenanceMargin,
					tc.bigNewNetCollateral,
					tc.bigNewMaintenanceMargin,
				),
			)
		})
	}
}
