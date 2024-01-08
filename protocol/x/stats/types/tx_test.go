package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	validAuthority = constants.BobAccAddress.String()
)

func TestValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpdateParams
		expectedErr error
	}{
		"Success": {
			msg: types.MsgUpdateParams{
				Authority: validAuthority,
				Params: types.Params{
					WindowDuration: 1 * time.Second,
				},
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpdateParams{
				Authority: "", // invalid - empty
			},
			expectedErr: types.ErrInvalidAuthority,
		},
		"Failure: Invalid params": {
			msg: types.MsgUpdateParams{
				Authority: validAuthority,
				Params: types.Params{
					WindowDuration: 0, // invalid - zero
				},
			},
			expectedErr: types.ErrNonpositiveDuration,
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
