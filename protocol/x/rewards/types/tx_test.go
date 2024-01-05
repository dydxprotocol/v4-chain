package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	validAuthority = constants.BobAccAddress.String()
)

func TestValidateBasic(t *testing.T) {
	test := map[string]struct {
		msg         types.MsgUpdateParams
		expectedErr error
	}{
		"Success": {
			msg: types.MsgUpdateParams{
				Authority: validAuthority,
				Params: types.Params{
					TreasuryAccount: "treasury_account",
					Denom:           "denom",
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
					TreasuryAccount: "", // invalid - empty
				},
			},
			expectedErr: types.ErrInvalidTreasuryAccount,
		},
	}
	for name, tc := range test {
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
