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
		"/cosmos.gov.v1.MsgCancelProposal",
		"/cosmos.gov.v1.MsgCancelProposalResponse",
		"/cosmos.gov.v1beta1.MsgSubmitProposal",
		"/cosmos.gov.v1beta1.MsgSubmitProposalResponse",

		// cosmwasm messages
		"/cosmwasm.wasm.v1.AcceptedMessageKeysFilter",
		"/cosmwasm.wasm.v1.AcceptedMessagesFilter",
		"/cosmwasm.wasm.v1.AllowAllMessagesFilter",
		"/cosmwasm.wasm.v1.ClearAdminProposal",
		"/cosmwasm.wasm.v1.CombinedLimit",
		"/cosmwasm.wasm.v1.ContractExecutionAuthorization",
		"/cosmwasm.wasm.v1.ContractMigrationAuthorization",
		"/cosmwasm.wasm.v1.ExecuteContractProposal",
		"/cosmwasm.wasm.v1.InstantiateContract2Proposal",
		"/cosmwasm.wasm.v1.InstantiateContractProposal",
		"/cosmwasm.wasm.v1.MaxCallsLimit",
		"/cosmwasm.wasm.v1.MaxFundsLimit",
		"/cosmwasm.wasm.v1.MigrateContractProposal",
		"/cosmwasm.wasm.v1.MsgIBCCloseChannel",
		"/cosmwasm.wasm.v1.MsgIBCSend",
		"/cosmwasm.wasm.v1.PinCodesProposal",
		"/cosmwasm.wasm.v1.StoreAndInstantiateContractProposal",
		"/cosmwasm.wasm.v1.StoreCodeAuthorization",
		"/cosmwasm.wasm.v1.StoreCodeProposal",
		"/cosmwasm.wasm.v1.SudoContractProposal",
		"/cosmwasm.wasm.v1.UnpinCodesProposal",
		"/cosmwasm.wasm.v1.UpdateAdminProposal",
		"/cosmwasm.wasm.v1.UpdateInstantiateConfigProposal",

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
