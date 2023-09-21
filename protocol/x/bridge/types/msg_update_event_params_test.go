package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateEventParams_GetSigners(t *testing.T) {
	msg := types.MsgUpdateEventParams{
		Authority: constants.CarlAccAddress.String(),
	}
	require.Equal(t, []sdk.AccAddress{constants.CarlAccAddress}, msg.GetSigners())
}

func TestMsgUpdateEventParams_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateEventParams
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpdateEventParams{
				Authority: "test",
				Params: types.EventParams{
					Denom:      "test-denom",
					EthChainId: 0,
					EthAddress: "test",
				},
			},
		},
		"Failure: Empty authority": {
			msg: types.MsgUpdateEventParams{
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
