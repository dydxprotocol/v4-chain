package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	"github.com/stretchr/testify/require"
)

func TestSubaccountLiquidationInfo_HasPerpetualBeenLiquidatedForSubaccount(t *testing.T) {
	tests := map[string]struct {
		// State.
		liquidatedPerpetuals []uint32

		// Parameters.
		perpetualId uint32

		// Expectations.
		expectedHasBeenLiquidated bool
	}{
		"no liquidated perpetuals doesn't return an error": {
			perpetualId: 5,
		},
		"new liquidated perpetual doesn't return an error": {
			liquidatedPerpetuals: []uint32{1, 2, 3, 4},

			perpetualId: 5,
		},
		"perpetual that has already been liquidated returns an error": {
			liquidatedPerpetuals: []uint32{1, 4, 3, 2},

			perpetualId: 3,

			expectedHasBeenLiquidated: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			subaccountLiquidationInfo := types.SubaccountLiquidationInfo{
				PerpetualsLiquidated: tc.liquidatedPerpetuals,
			}

			hasBeenLiquidated := subaccountLiquidationInfo.HasPerpetualBeenLiquidatedForSubaccount(tc.perpetualId)
			require.Equal(t, tc.expectedHasBeenLiquidated, hasBeenLiquidated)
		})
	}
}
