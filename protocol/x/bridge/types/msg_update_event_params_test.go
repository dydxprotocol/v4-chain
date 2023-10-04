package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

var (
	// validAuthority is a valid bech32 address string.
	validAuthority = constants.AliceAccAddress.String()
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
		expectedErr error
	}{
		"Success": {
			msg: types.MsgUpdateEventParams{
				Authority: validAuthority,
				Params: types.EventParams{
					Denom:      "test-denom",
					EthChainId: 0,
					EthAddress: "test",
				},
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpdateEventParams{
				Authority: "",
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
