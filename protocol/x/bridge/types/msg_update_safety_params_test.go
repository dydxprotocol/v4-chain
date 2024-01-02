package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateSafetyParams_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateSafetyParams
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpdateSafetyParams{
				Authority: validAuthority,
				Params: types.SafetyParams{
					IsDisabled:  false,
					DelayBlocks: 500,
				},
			},
		},
		"Failure: empty authority": {
			msg: types.MsgUpdateSafetyParams{
				Authority: "",
			},
			expectedErr: types.ErrInvalidAuthority.Error(),
		},
		"Failure: invalid authority": {
			msg: types.MsgUpdateSafetyParams{
				Authority: "dydx1abc",
			},
			expectedErr: types.ErrInvalidAuthority.Error(),
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
