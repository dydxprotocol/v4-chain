package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
)

// IsUnsupportedMsg returns true if the msg is unsupported by the app.
func IsUnsupportedMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// ICA Controller messages
		*icacontrollertypes.MsgSendTx,
		*icacontrollertypes.MsgRegisterInterchainAccount,
		// ------- CosmosSDK default modules
		// gov
		*govbeta.MsgSubmitProposal:
		return true
	}
	return false
}
