package msgs_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestNestedMsgs_Key(t *testing.T) {
	expectedMsgs := []string{
		// gov
		"/cosmos.gov.v1.MsgSubmitProposal",
		"/cosmos.gov.v1.MsgSubmitProposalResponse",
	}
	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.NestedMsgSamples))
}

func TestNestedMsgs_Value(t *testing.T) {
	validateSampleMsgValue(t, msgs.NestedMsgSamples)
}
