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
		"/cosmwasm.wasm.v1.AcceptedMessageKeysFilter",
		"/cosmwasm.wasm.v1.AcceptedMessagesFilter",
		"/cosmwasm.wasm.v1.AllowAllMessagesFilter",
		"/cosmwasm.wasm.v1.ClearAdminProposal",
		"/cosmwasm.wasm.v1.CombinedLimit",
		"/cosmwasm.wasm.v1.ContractExecutionAuthorization",
		"/cosmwasm.wasm.v1.ContractMigrationAuthorization",
		"/cosmwasm.wasm.v1.MaxCallsLimit",
		"/cosmwasm.wasm.v1.MaxFundsLimit",
		"/cosmwasm.wasm.v1.MsgClearAdmin",
		"/cosmwasm.wasm.v1.MsgClearAdminResponse",
		"/cosmwasm.wasm.v1.MsgIBCCloseChannel",
		"/cosmwasm.wasm.v1.MsgIBCSend",
		"/cosmwasm.wasm.v1.MsgMigrateContract",
		"/cosmwasm.wasm.v1.MsgMigrateContractResponse",
		"/cosmwasm.wasm.v1.MsgPinCodes",
		"/cosmwasm.wasm.v1.MsgPinCodesResponse",
		"/cosmwasm.wasm.v1.MsgRemoveCodeUploadParamsAddresses",
		"/cosmwasm.wasm.v1.MsgRemoveCodeUploadParamsAddressesResponse",
		"/cosmwasm.wasm.v1.MsgSudoContract",
		"/cosmwasm.wasm.v1.MsgSudoContractResponse",
		"/cosmwasm.wasm.v1.MsgUnpinCodes",
		"/cosmwasm.wasm.v1.MsgUnpinCodesResponse",
		"/cosmwasm.wasm.v1.MsgUpdateAdmin",
		"/cosmwasm.wasm.v1.MsgUpdateAdminResponse",
		"/cosmwasm.wasm.v1.MsgUpdateContractLabel",
		"/cosmwasm.wasm.v1.MsgUpdateContractLabelResponse",
		"/cosmwasm.wasm.v1.MsgUpdateInstantiateConfig",
		"/cosmwasm.wasm.v1.MsgUpdateInstantiateConfigResponse",
		"/cosmwasm.wasm.v1.PinCodesProposal",
		"/cosmwasm.wasm.v1.StoreCodeAuthorization",
		"/cosmwasm.wasm.v1.SudoContractProposal",
		"/cosmwasm.wasm.v1.UnpinCodesProposal",

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
