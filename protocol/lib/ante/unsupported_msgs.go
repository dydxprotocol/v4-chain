package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
)

// IsUnsupportedMsg returns true if the msg is unsupported by the app.
func IsUnsupportedMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// Some gov msgs
		*gov.MsgCancelProposal,
		*govbeta.MsgSubmitProposal,
		// ICA Controller messages
		*icacontrollertypes.MsgUpdateParams,
		*icacontrollertypes.MsgSendTx,
		*icacontrollertypes.MsgRegisterInterchainAccount:
		return true
	}
	return false
}
