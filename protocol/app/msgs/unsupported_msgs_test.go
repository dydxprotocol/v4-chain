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
		// crisis
		"/cosmos.crisis.v1beta1.MsgVerifyInvariant",
		"/cosmos.crisis.v1beta1.MsgVerifyInvariantResponse",

		"/cosmos.gov.v1.MsgCancelProposal",
		"/cosmos.gov.v1.MsgCancelProposalResponse",
		"/cosmos.gov.v1beta1.MsgSubmitProposal",
		"/cosmos.gov.v1beta1.MsgSubmitProposalResponse",

		"/dydxprotocol.vault.MsgSetVaultQuotingParams",
		"/dydxprotocol.vault.MsgUpdateParams",

		// ICA Controller messages
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount",
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccountResponse",
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTx",
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTxResponse",
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParams",
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParamsResponse",
	}

	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.UnsupportedMsgSamples))
}

func TestUnsupportedMsgSamples_Value(t *testing.T) {
	validateMsgValue(t, msgs.UnsupportedMsgSamples)
}
