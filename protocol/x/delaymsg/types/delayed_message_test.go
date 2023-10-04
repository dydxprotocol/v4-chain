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

func TestDelayedMessage_Validate(t *testing.T) {
	tests := map[string]struct {
		dm          types.DelayedMessage
		expectedErr error
	}{
		"Failure: nil message": {
			dm:          types.DelayedMessage{},
			expectedErr: types.ErrMsgIsNil,
		},
		"Success": {
			dm: types.DelayedMessage{
				Msg: encoding.EncodeMessageToAny(t, constants.TestMsg1),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.dm.Validate()
			if test.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, test.expectedErr.Error())
			}
		})
	}
}

func TestDelayedMessage_GetMessage(t *testing.T) {
	tests := map[string]struct {
		dm          types.DelayedMessage
		expectedMsg sdk.Msg
		expectedErr error
	}{
		"Failure: nil message": {
			dm:          types.DelayedMessage{},
			expectedErr: types.ErrMsgIsNil,
		},
		"Failure: uncached message": {
			dm: types.DelayedMessage{
				Msg: &codectypes.Any{},
			},
			expectedErr: fmt.Errorf("any cached value is nil, delayed messages must be correctly packed any values"),
		},
		"Failure: encoded message is not a sdk.Msg": {
			dm: types.DelayedMessage{
				Msg: codectypes.UnsafePackAny("not an sdk.Msg"),
			},
			expectedErr: fmt.Errorf("cached value is not a sdk.Msg"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			msg, err := test.dm.GetMessage()
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
