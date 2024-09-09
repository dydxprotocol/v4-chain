package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/sample"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMsgClaimYieldForSubaccount_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgClaimYieldForSubaccount
		err  error
	}{
		{
			name: "Invalid subaccount owner",
			msg: types.MsgClaimYieldForSubaccount{
				Id: &types.SubaccountId{
					Owner:  "invalid_owner",
					Number: uint32(0),
				},
			},
			err: types.ErrInvalidSubaccountIdOwner,
		},
		{
			name: "Invalid subaccount number",
			msg: types.MsgClaimYieldForSubaccount{
				Id: &types.SubaccountId{
					Owner:  sample.AccAddress(),
					Number: uint32(999_999),
				},
			},
			err: types.ErrInvalidSubaccountIdNumber,
		},
		{
			name: "Valid address 1",
			msg: types.MsgClaimYieldForSubaccount{
				Id: &constants.Carl_Num0,
			},
		},
		{
			name: "Valid address 2",
			msg: types.MsgClaimYieldForSubaccount{
				Id: &constants.Dave_Num0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
