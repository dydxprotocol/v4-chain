package msgs

import (
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/dydxprotocol/v4/testutil/encoding"
)

func init() {
	testEncodingCfg := encoding.GetTestEncodingCfg()
	testTxBuilder := testEncodingCfg.TxConfig.NewTxBuilder()

	_ = testTxBuilder.SetMsgs(GovBetaMsgSubmitProposal)
	GovBetaMsgSubmitProposalTxBytes, _ = testEncodingCfg.TxConfig.TxEncoder()(testTxBuilder.GetTx())
}

var (
	GovBetaMsgSubmitProposal        = &govbeta.MsgSubmitProposal{}
	GovBetaMsgSubmitProposalTxBytes []byte
)
