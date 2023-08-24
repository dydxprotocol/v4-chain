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
				ClobPairId: 0,
				Status: types.ClobPair_STATUS_ACTIVE,
			},
		},
		"invalid": {
			msg: types.MsgSetClobPairStatus{
				ClobPairId: 0,
				Status: types.ClobPair_STATUS_UNSPECIFIED,
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

// func TestMsgSetClobPairStatus_GetSigners(t *testing.T) {
// 	tests := map[string]struct {
// 		msg    types.MsgSetClobPairStatus
// 		signers []sdk.AccAddress
// 	}{
// 		"valid": {
// 			msg: types.MsgSetClobPairStatus{
// 				ClobPairId: 0,
// 				Status: types.ClobPair_STATUS_ACTIVE,
// 			},
// 			signers: make([]sdk.AccAddress, 0),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			signers := tc.msg.GetSigners()
// 			require.Equal(t, tc.signers, signers)
// 		})
// 	}
// }