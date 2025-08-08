package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
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

		// vault
		// MsgSetVaultQuotingParams is deprecated since v6.x and replaced by MsgSetVaultParams.
		// nolint:staticcheck
		"/dydxprotocol.vault.MsgSetVaultQuotingParams": &vaulttypes.MsgSetVaultQuotingParams{},
		// MsgUpdateParams is deprecated since v6.x and replaced by MsgUpdateDefaultQuotingParams.
		// nolint:staticcheck
		"/dydxprotocol.vault.MsgUpdateParams": &vaulttypes.MsgUpdateParams{},

		// Disable MsgVerifyInvariant in the crisis module, since:
		// 1. We currently do not rely on crisis module for any invariant assertion.
		// 2. MsgVerifyInvariant can potentially be abused to consume massive compute.
		"/cosmos.crisis.v1beta1.MsgVerifyInvariant":         &crisis.MsgVerifyInvariant{},
		"/cosmos.crisis.v1beta1.MsgVerifyInvariantResponse": nil,
	}
)
