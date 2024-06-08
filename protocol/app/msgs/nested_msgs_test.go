package msgs_test

import (
	"sort"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/msgs"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestNestedMsgs_Key(t *testing.T) {
	expectedMsgs := []string{
		// authz
		"/cosmos.authz.v1beta1.MsgExec",
		"/cosmos.authz.v1beta1.MsgExecResponse",
	}
	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.NestedMsgSamples))
}

func TestNestedMsgs_Value(t *testing.T) {
	validateMsgValue(t, msgs.NestedMsgSamples)
}
