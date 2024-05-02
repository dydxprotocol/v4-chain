package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
)

var (
	// UnsupportedMsgSamples are msgs that are registered with the app, but are not supported.
	UnsupportedMsgSamples = map[string]sdk.Msg{

		// ICA Controller messages - these are not used since ICA Controller is disabled.
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount": &icacontrollertypes.
			MsgRegisterInterchainAccount{},
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccountResponse": nil,
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTx": &icacontrollertypes.
			MsgSendTx{},
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTxResponse": nil,
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParams": &icacontrollertypes.
			MsgUpdateParams{},
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParamsResponse": nil,
	}
)
