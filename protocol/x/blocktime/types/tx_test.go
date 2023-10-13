package types_test

import (
	"testing"
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/stretchr/testify/require"
)

var (
	validAuthority = constants.AliceAccAddress.String()
)

func TestMsgUpdateDowntimeParams_GetSigners(t *testing.T) {
	msg := types.MsgUpdateDowntimeParams{
		Authority: validAuthority,
	}
	require.Equal(t, []sdk.AccAddress{constants.AliceAccAddress}, msg.GetSigners())
}

func TestMsgUpdateDowntimeParams_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateDowntimeParams
		expectedErr error
	}{
		"Success": {
			msg: types.MsgUpdateDowntimeParams{
				Authority: validAuthority,
				Params:    types.DowntimeParams{},
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpdateDowntimeParams{
				Authority: "", // invalid
			},
			expectedErr: types.ErrInvalidAuthority,
		},
		"Failure: Invalid params": {
			msg: types.MsgUpdateDowntimeParams{
				Authority: validAuthority,
				Params: types.DowntimeParams{
					Durations: []time.Duration{
						5 * time.Second,
						1 * time.Second,
					},
				},
			},
			expectedErr: types.ErrUnorderedDurations,
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
