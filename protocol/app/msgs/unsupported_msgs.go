package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var (
	// UnsupportedMsgSamples are msgs that are registered with the app, but are not supported.
	UnsupportedMsgSamples = map[string]sdk.Msg{
		// gov
		"/cosmos.gov.v1beta1.MsgSubmitProposal":         &govbeta.MsgSubmitProposal{},
		"/cosmos.gov.v1beta1.MsgSubmitProposalResponse": nil,
	}
)
