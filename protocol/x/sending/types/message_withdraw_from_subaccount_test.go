package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMsgWithdrawFromSubaccount_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg types.MsgWithdrawFromSubaccount
		err error
	}{
		"Valid - withdraw from same account's subaccount": {
			msg: constants.MsgWithdrawFromSubaccount_Alice_Num0_To_Alice_500,
		},
		"Valid - withdraw from different account's subaccount": {
			msg: constants.MsgWithdrawFromSubaccount_Carl_Num0_To_Alice_750,
		},
		"Invalid sender owner": {
			msg: types.MsgWithdrawFromSubaccount{
				Sender: satypes.SubaccountId{
					Owner:  "invalid_owner",
					Number: uint32(0),
				},
				Recipient: constants.AliceAccAddress.String(),
			},
			err: satypes.ErrInvalidSubaccountIdOwner,
		},
		"Invalid recipient address": {
			msg: types.MsgWithdrawFromSubaccount{
				Sender:    constants.Alice_Num0,
				Recipient: "invalid_address",
			},
			err: types.ErrInvalidAccountAddress,
		},
		"Non-TDai asset transfer not supported": {
			msg: types.MsgWithdrawFromSubaccount{
				Sender:    constants.Alice_Num0,
				Recipient: constants.AliceAccAddress.String(),
				AssetId:   uint32(1),
				Quantums:  uint64(100),
			},
			err: types.ErrNonTDaiAssetTransferNotImplemented,
		},
		"Invalid quantums": {
			msg: types.MsgWithdrawFromSubaccount{
				Sender:    constants.Alice_Num0,
				Recipient: constants.AliceAccAddress.String(),
				AssetId:   uint32(0),
				Quantums:  uint64(0),
			},
			err: types.ErrInvalidTransferAmount,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
