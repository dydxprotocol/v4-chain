package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	types "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/stretchr/testify/require"
)

var (
	// validAuthority is a valid bech32 address.
	validAuthority = constants.AliceAccAddress.String()
)

func TestMsgUpgradeIsolatedPerpetualToCross_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpgradeIsolatedPerpetualToCross
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 1,
			},
		},
		"Failure: Empty authority": {
			msg: types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   "",
				PerpetualId: 1,
			},
			expectedErr: "Authority is invalid",
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   "invalid",
				PerpetualId: 1,
			},
			expectedErr: "Authority is invalid",
		},
		"Failure: Invalid perpetual ID": {
			msg: types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 0,
			},
			expectedErr: "Invalid perpetual ID",
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
