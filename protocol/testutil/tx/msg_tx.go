package tx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
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
func MustGetOnlySignerAddress(msg sdk.Msg) string {
	if len(msg.GetSigners()) == 0 {
		panic(fmt.Errorf("msg does not have designated signer: %T", msg))
	}
	if len(msg.GetSigners()) > 1 {
		panic(fmt.Errorf("not supported - msg has multiple signers: %T", msg))
	}
	return msg.GetSigners()[0].String()
}
