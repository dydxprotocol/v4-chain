package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

var (
	// NestedMsgSamples are msgs that have can include other arbitrary messages inside.
	NestedMsgSamples = map[string]sdk.Msg{
		// authz
		"/cosmos.authz.v1beta1.MsgExec":         &authz.MsgExec{},
		"/cosmos.authz.v1beta1.MsgExecResponse": nil,

		// gov
		"/cosmos.gov.v1.MsgSubmitProposal":         &gov.MsgSubmitProposal{},
		"/cosmos.gov.v1.MsgSubmitProposalResponse": nil,
	}
)
