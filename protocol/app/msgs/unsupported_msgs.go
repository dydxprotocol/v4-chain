package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
)

var (
	// UnsupportedMsgSamples are msgs that are registered with the app, but are not supported.
	UnsupportedMsgSamples = map[string]sdk.Msg{
		// gov
		"/cosmos.gov.v1beta1.MsgSubmitProposal":         &govbeta.MsgSubmitProposal{},
		"/cosmos.gov.v1beta1.MsgSubmitProposalResponse": nil,

		// ICA Controller messages - these are not used since ICA Controller is disabled.
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTx": &icacontrollertypes.MsgSendTx{},
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount": &icacontrollertypes.
			MsgRegisterInterchainAccount{},
	}
)
