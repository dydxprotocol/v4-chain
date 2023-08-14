package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMsgDepositToSubaccount_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg types.MsgDepositToSubaccount
		err error
	}{
		"Valid - deposit to same account's subaccount": {
			msg: constants.MsgDepositToSubaccount_Alice_To_Alice_Num0_500,
		},
		"Valid - deposit to different account's subaccount": {
			msg: constants.MsgDepositToSubaccount_Alice_To_Carl_Num0_750,
		},
		"Invalid sender address": {
			msg: types.MsgDepositToSubaccount{
				Sender:    "invalid_address",
				Recipient: constants.Alice_Num0,
			},
			err: types.ErrInvalidAccountAddress,
		},
		"Invalid recipient owner": {
			msg: types.MsgDepositToSubaccount{
				Sender: constants.AliceAccAddress.String(),
				Recipient: satypes.SubaccountId{
					Owner:  "invalid_owner",
					Number: uint32(0),
				},
			},
			err: satypes.ErrInvalidSubaccountIdOwner,
		},
		"Non-USDC asset transfer not supported": {
			msg: types.MsgDepositToSubaccount{
				Sender:    constants.AliceAccAddress.String(),
				Recipient: constants.Alice_Num0,
				AssetId:   uint32(1),
				Quantums:  uint64(100),
			},
			err: types.ErrNonUsdcAssetTransferNotImplemented,
		},
		"Invalid quantums": {
			msg: types.MsgDepositToSubaccount{
				Sender:    constants.AliceAccAddress.String(),
				Recipient: constants.Alice_Num0,
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
