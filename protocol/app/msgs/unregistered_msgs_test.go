package msgs_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestUnregisteredMsgs_Key(t *testing.T) {
	expectedMsgs := []string{
		// authz
		"/cosmos.authz.v1.MsgExec",
		"/cosmos.authz.v1.MsgExecResponse",

		// group
		"/cosmos.group.v1.MsgSubmitProposal",
		"/cosmos.group.v1.MsgSubmitProposalResponse",
		"/cosmos.group.v1beta1.MsgSubmitProposal",
		"/cosmos.group.v1beta1.MsgSubmitProposalResponse",
	}
	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.UnregisteredMsgs))
}

func TestUnregisteredMsgs_Value(t *testing.T) {
	for _, msg := range msgs.UnregisteredMsgs {
		require.Equal(t, struct{}{}, msg)
	}
}
