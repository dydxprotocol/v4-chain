package msgs_test

import (
	"sort"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/msgs"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestUnsupportedMsgSamples_Key(t *testing.T) {
	expectedMsgs := []string{

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
