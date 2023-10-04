package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateProposeParams_GetSigners(t *testing.T) {
	msg := types.MsgUpdateProposeParams{
		Authority: constants.CarlAccAddress.String(),
	}
	require.Equal(t, []sdk.AccAddress{constants.CarlAccAddress}, msg.GetSigners())
}

func TestMsgUpdateProposeParams_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateProposeParams
		expectedErr error
	}{
		"Success": {
			msg: types.MsgUpdateProposeParams{
				Authority: validAuthority,
				Params: types.ProposeParams{
					MaxBridgesPerBlock:           5,
					ProposeDelayDuration:         10_000,
					SkipRatePpm:                  800_000,
					SkipIfBlockDelayedByDuration: 5_000,
				},
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpdateProposeParams{
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
