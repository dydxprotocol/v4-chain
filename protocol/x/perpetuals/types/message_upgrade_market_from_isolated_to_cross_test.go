package types_test

import (
	"testing"

	types "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpgradeMarketFromIsolatedToCross_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpgradeMarketFromIsolatedToCross
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpgradeMarketFromIsolatedToCross{
				Authority:   validAuthority,
				PerpetualId: 1,
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpgradeMarketFromIsolatedToCross{
				Authority: "",
			},
			expectedErr: "Authority is invalid",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
