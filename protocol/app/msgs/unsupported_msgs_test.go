package msgs_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestUnsupportedMsgSamples_Key(t *testing.T) {
	expectedMsgs := []string{
		"/cosmos.gov.v1beta1.MsgSubmitProposal",
		"/cosmos.gov.v1beta1.MsgSubmitProposalResponse",

		// ICA Controller messages
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount",
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTx",
	}

	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.UnsupportedMsgSamples))
}

func TestUnsupportedMsgSamples_Value(t *testing.T) {
	validateMsgValue(t, msgs.UnsupportedMsgSamples)
}
