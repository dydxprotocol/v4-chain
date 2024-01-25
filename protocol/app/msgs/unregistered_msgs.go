package msgs

var (
	// UnregisteredMsgs are msgs that should not be registered with the app.
	UnregisteredMsgs = map[string]struct{}{
		// authz
		"/cosmos.authz.v1.MsgExec":         {},
		"/cosmos.authz.v1.MsgExecResponse": {},

		// group
		"/cosmos.group.v1.MsgSubmitProposal":              {},
		"/cosmos.group.v1.MsgSubmitProposalResponse":      {},
		"/cosmos.group.v1beta1.MsgSubmitProposal":         {},
		"/cosmos.group.v1beta1.MsgSubmitProposalResponse": {},
	}
)
