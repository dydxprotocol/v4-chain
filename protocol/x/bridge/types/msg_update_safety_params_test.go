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
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpdateSafetyParams{
				Authority: "test",
				Params: types.SafetyParams{
					IsDisabled:  false,
					DelayBlocks: 500,
				},
			},
		},
		"Failure: Empty authority": {
			msg: types.MsgUpdateSafetyParams{
				Authority: "",
			},
			expectedErr: "authority cannot be empty",
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
