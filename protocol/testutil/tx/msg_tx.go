package tx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4/testutil/constants"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	sendingtypes "github.com/dydxprotocol/v4/x/sending/types"
)

// Returns the encoded msg as transaction. Will panic if encoding fails.
func MustGetTxBytes(msgs ...sdk.Msg) []byte {
	tx := constants.TestEncodingCfg.TxConfig.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}
	bz, err := constants.TestEncodingCfg.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		panic(err)
	}
	return bz
}

// Returns the account address that should sign the msg. Will panic if it is an unsupported message type.
func MustGetSignerAddress(msg sdk.Msg) string {
	switch v := any(msg).(type) {
	case *clobtypes.MsgPlaceOrder:
		return v.Order.OrderId.SubaccountId.Owner
	case *clobtypes.MsgCancelOrder:
		return v.OrderId.SubaccountId.Owner
	case *sendingtypes.MsgDepositToSubaccount:
		return v.Sender
	case *sendingtypes.MsgWithdrawFromSubaccount:
		return v.Sender.Owner
	default:
		panic(fmt.Errorf("Not a supported type %T", msg))
	}
}
