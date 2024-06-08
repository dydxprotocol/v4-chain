package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
)

// IsUnsupportedMsg returns true if the msg is unsupported by the app.
func IsUnsupportedMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// ICA Controller messages
		*icacontrollertypes.MsgUpdateParams,
		*icacontrollertypes.MsgSendTx,
		*icacontrollertypes.MsgRegisterInterchainAccount:
		return true
	}
	return false
}
