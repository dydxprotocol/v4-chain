package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSetClobPairStatus_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg           types.MsgSetClobPairStatus
		expectedError error
	}{
		"valid": {
			msg: types.MsgSetClobPairStatus{
				ClobPairId:     0,
				ClobPairStatus: int32(types.ClobPair_STATUS_ACTIVE),
			},
		},
		"invalid unsupported status": {
			msg: types.MsgSetClobPairStatus{
				ClobPairId:     0,
				ClobPairStatus: int32(types.ClobPair_STATUS_UNSPECIFIED),
			},
			expectedError: types.ErrInvalidMsgSetClobPairStatus,
		},
		"invalid negative out of bounds status": {
			msg: types.MsgSetClobPairStatus{
				ClobPairId:     0,
				ClobPairStatus: -1,
			},
			expectedError: types.ErrInvalidMsgSetClobPairStatus,
		},
		"invalid positive out of bounds status": {
			msg: types.MsgSetClobPairStatus{
				ClobPairId:     0,
				ClobPairStatus: 100,
			},
			expectedError: types.ErrInvalidMsgSetClobPairStatus,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
