package msgs

import (
	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
)

var (
	// UnsupportedMsgSamples are msgs that are registered with the app, but are not supported.
	UnsupportedMsgSamples = map[string]sdk.Msg{

		// gov
		// MsgCancelProposal is not allowed by protocol, due to it's potential for abuse.
		"/cosmos.gov.v1.MsgCancelProposal":         &gov.MsgCancelProposal{},
		"/cosmos.gov.v1.MsgCancelProposalResponse": nil,
		// These are deprecated/legacy msgs that we should not support.
		"/cosmos.gov.v1beta1.MsgSubmitProposal":         &govbeta.MsgSubmitProposal{},
		"/cosmos.gov.v1beta1.MsgSubmitProposalResponse": nil,

		// cosmwasm
		"/cosmwasm.wasm.v1.AcceptedMessageKeysFilter":                  nil,
		"/cosmwasm.wasm.v1.AcceptedMessagesFilter":                     nil,
		"/cosmwasm.wasm.v1.AllowAllMessagesFilter":                     nil,
		"/cosmwasm.wasm.v1.ClearAdminProposal":                         nil,
		"/cosmwasm.wasm.v1.CombinedLimit":                              nil,
		"/cosmwasm.wasm.v1.MaxCallsLimit":                              nil,
		"/cosmwasm.wasm.v1.MaxFundsLimit":                              nil,
		"/cosmwasm.wasm.v1.MsgClearAdmin":                              &wasm.MsgClearAdmin{},
		"/cosmwasm.wasm.v1.MsgClearAdminResponse":                      nil,
		"/cosmwasm.wasm.v1.MsgIBCCloseChannel":                         &wasm.MsgIBCCloseChannel{},
		"/cosmwasm.wasm.v1.MsgIBCSend":                                 &wasm.MsgIBCSend{},
		"/cosmwasm.wasm.v1.MsgPinCodes":                                &wasm.MsgPinCodes{},
		"/cosmwasm.wasm.v1.MsgPinCodesResponse":                        nil,
		"/cosmwasm.wasm.v1.MsgRemoveCodeUploadParamsAddresses":         &wasm.MsgRemoveCodeUploadParamsAddresses{},
		"/cosmwasm.wasm.v1.MsgRemoveCodeUploadParamsAddressesResponse": nil,
		"/cosmwasm.wasm.v1.MsgSudoContract":                            &wasm.MsgSudoContract{},
		"/cosmwasm.wasm.v1.MsgSudoContractResponse":                    nil,
		"/cosmwasm.wasm.v1.MsgUnpinCodes":                              &wasm.MsgUnpinCodes{},
		"/cosmwasm.wasm.v1.MsgUnpinCodesResponse":                      nil,
		"/cosmwasm.wasm.v1.MsgUpdateAdmin":                             &wasm.MsgUpdateAdmin{},
		"/cosmwasm.wasm.v1.MsgUpdateAdminResponse":                     nil,
		"/cosmwasm.wasm.v1.UnpinCodesProposal":                         nil,
		"/cosmwasm.wasm.v1.PinCodesProposal":                           nil,
		"/cosmwasm.wasm.v1.SudoContractProposal":                       nil,
		"/cosmwasm.wasm.v1.MsgUpdateContractLabel":                     &wasm.MsgUpdateContractLabel{},
		"/cosmwasm.wasm.v1.MsgUpdateContractLabelResponse":             nil,
		"/cosmwasm.wasm.v1.MsgUpdateInstantiateConfig":                 &wasm.MsgUpdateInstantiateConfig{},
		"/cosmwasm.wasm.v1.MsgUpdateInstantiateConfigResponse":         nil,
		"/cosmwasm.wasm.v1.MsgMigrateContract":                         &wasm.MsgMigrateContract{},
		"/cosmwasm.wasm.v1.MsgMigrateContractResponse":                 nil,
		"/cosmwasm.wasm.v1.ContractExecutionAuthorization":             nil,
		"/cosmwasm.wasm.v1.ContractMigrationAuthorization":             nil,
		"/cosmwasm.wasm.v1.StoreCodeAuthorization":                     nil,

		// ICA Controller messages - these are not used since ICA Controller is disabled.
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount": &icacontrollertypes.
			MsgRegisterInterchainAccount{},
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccountResponse": nil,
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTx": &icacontrollertypes.
			MsgSendTx{},
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTxResponse": nil,
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParams": &icacontrollertypes.
			MsgUpdateParams{},
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParamsResponse": nil,
	}
)
