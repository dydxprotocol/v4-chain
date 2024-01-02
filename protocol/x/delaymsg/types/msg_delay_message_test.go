package types_test

import (
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	AcceptedAuthority = types.ModuleAddress
)

func TestMsgDelayMessage_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		mdm         types.MsgDelayMessage
		expectedErr error
	}{
		"Failure: invalid authority": {
			mdm:         types.MsgDelayMessage{},
			expectedErr: types.ErrInvalidAuthority,
		},
		"Failure: nil message": {
			mdm: types.MsgDelayMessage{
				Authority: AcceptedAuthority.String(),
			},
			expectedErr: types.ErrMsgIsNil,
		},
		"Success": {
			mdm: types.MsgDelayMessage{
				Authority: AcceptedAuthority.String(),
				Msg:       encoding.EncodeMessageToAny(t, constants.TestMsg1),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.mdm.ValidateBasic()
			if test.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, test.expectedErr.Error())
			}
		})
	}
}

func TestMsgDelayMessage_GetMessage(t *testing.T) {
	tests := map[string]struct {
		mdm         types.MsgDelayMessage
		expectedMsg sdk.Msg
		expectedErr error
	}{
		"Failure: nil message": {
			mdm:         types.MsgDelayMessage{},
			expectedErr: types.ErrMsgIsNil,
		},
		"Failure: uncached message": {
			mdm: types.MsgDelayMessage{
				Msg: &codectypes.Any{},
			},
			expectedErr: fmt.Errorf("any cached value is nil, delayed messages must be correctly packed any values"),
		},
		"Failure: encoded message is not a sdk.Msg": {
			mdm: types.MsgDelayMessage{
				Msg: codectypes.UnsafePackAny("not an sdk.Msg"),
			},
			expectedErr: fmt.Errorf("cached value is not a sdk.Msg"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg, err := test.mdm.GetMessage()
			if test.expectedErr == nil {
				require.NoError(t, err)
				require.Equal(t, test.expectedMsg, msg)
			} else {
				require.ErrorContains(t, err, test.expectedErr.Error())
				require.Nil(t, msg)
			}
		})
	}
}
