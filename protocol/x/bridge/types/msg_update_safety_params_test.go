package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateSafetyParams_GetSigners(t *testing.T) {
	msg := types.MsgUpdateProposeParams{
		Authority: constants.BobAccAddress.String(),
	}
	require.Equal(t, []sdk.AccAddress{constants.BobAccAddress}, msg.GetSigners())
}

func TestMsgUpdateSafetyParams_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateSafetyParams
		expectedErr error
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
		"Failure: Invalid authority": {
			msg:         types.MsgUpdateSafetyParams{},
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
