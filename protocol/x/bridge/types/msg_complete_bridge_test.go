package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgCompleteBridge_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgCompleteBridge
		expectedErr error
	}{
		"Success": {
			msg: *constants.TestMsg1,
		},
		"Failure: Empty authority": {
			msg: types.MsgCompleteBridge{
				Authority: "",
			},
			expectedErr: types.ErrInvalidAuthority,
		},
		"Failure: Not an address": {
			msg: types.MsgCompleteBridge{
				Authority: "invalid",
			},
			expectedErr: types.ErrInvalidAuthority,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
