package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var (
	// NestedMsgSamples are msgs that have can include other arbitrary messages inside.
	NestedMsgSamples = map[string]sdk.Msg{
		// authz
		"/cosmos.authz.v1beta1.MsgExec":         &authz.MsgExec{},
		"/cosmos.authz.v1beta1.MsgExecResponse": nil,
	}
)
