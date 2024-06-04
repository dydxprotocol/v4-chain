package msgs

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
)

var (
	// AllTypeMessages is a list of all messages and types that are used in the app.
	// This list comes from the app's `InterfaceRegistry`.
	AllTypeMessages = map[string]struct{}{
		// auth
		"/cosmos.auth.v1beta1.BaseAccount":      {},
		"/cosmos.auth.v1beta1.ModuleAccount":    {},
		"/cosmos.auth.v1beta1.ModuleCredential": {},
		"/cosmos.auth.v1beta1.MsgUpdateParams":  {},

		// authz
		"/cosmos.authz.v1beta1.GenericAuthorization": {},
		"/cosmos.authz.v1beta1.MsgExec":              {},
		"/cosmos.authz.v1beta1.MsgExecResponse":      {},
		"/cosmos.authz.v1beta1.MsgGrant":             {},
		"/cosmos.authz.v1beta1.MsgGrantResponse":     {},
		"/cosmos.authz.v1beta1.MsgRevoke":            {},
		"/cosmos.authz.v1beta1.MsgRevokeResponse":    {},

		// bank
		"/cosmos.bank.v1beta1.MsgMultiSend":              {},
		"/cosmos.bank.v1beta1.MsgMultiSendResponse":      {},
		"/cosmos.bank.v1beta1.MsgSend":                   {},
		"/cosmos.bank.v1beta1.MsgSendResponse":           {},
		"/cosmos.bank.v1beta1.MsgSetSendEnabled":         {},
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse": {},
		"/cosmos.bank.v1beta1.MsgUpdateParams":           {},
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse":   {},
		"/cosmos.bank.v1beta1.SendAuthorization":         {},
		"/cosmos.bank.v1beta1.Supply":                    {},

		// consensus
		"/cosmos.consensus.v1.MsgUpdateParams":         {},
		"/cosmos.consensus.v1.MsgUpdateParamsResponse": {},

		// crisis
		"/cosmos.crisis.v1beta1.MsgUpdateParams":            {},
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse":    {},
		"/cosmos.crisis.v1beta1.MsgVerifyInvariant":         {},
		"/cosmos.crisis.v1beta1.MsgVerifyInvariantResponse": {},

		// crypto
		"/cosmos.crypto.ed25519.PrivKey":            {},
		"/cosmos.crypto.ed25519.PubKey":             {},
		"/cosmos.crypto.multisig.LegacyAminoPubKey": {},
		"/cosmos.crypto.secp256k1.PrivKey":          {},
		"/cosmos.crypto.secp256k1.PubKey":           {},
		"/cosmos.crypto.secp256r1.PubKey":           {},

		// evidence
		"/cosmos.evidence.v1beta1.Equivocation":              {},
		"/cosmos.evidence.v1beta1.MsgSubmitEvidence":         {},
		"/cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse": {},

		// feegrant
		"/cosmos.feegrant.v1beta1.AllowedMsgAllowance":        {},
		"/cosmos.feegrant.v1beta1.BasicAllowance":             {},
		"/cosmos.feegrant.v1beta1.MsgGrantAllowance":          {},
		"/cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse":  {},
		"/cosmos.feegrant.v1beta1.MsgPruneAllowances":         {},
		"/cosmos.feegrant.v1beta1.MsgPruneAllowancesResponse": {},
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowance":         {},
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse": {},
		"/cosmos.feegrant.v1beta1.PeriodicAllowance":          {},

		// params
		"/cosmos.params.v1beta1.ParameterChangeProposal": {},

		// slashing
		"/cosmos.slashing.v1beta1.MsgUnjail":               {},
		"/cosmos.slashing.v1beta1.MsgUnjailResponse":       {},
		"/cosmos.slashing.v1beta1.MsgUpdateParams":         {},
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse": {},

		// tx
		"/cosmos.tx.v1beta1.Tx": {},

		// upgrade
		"/cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal": {},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade":              {},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse":      {},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade":            {},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse":    {},
		"/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal":       {},

		// blocktime
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParams":         {},
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse": {},

		// clob
		"/dydxprotocol.clob.MsgCancelOrder":                                {},
		"/dydxprotocol.clob.MsgCancelOrderResponse":                        {},
		"/dydxprotocol.clob.MsgCreateClobPair":                             {},
		"/dydxprotocol.clob.MsgCreateClobPairResponse":                     {},
		"/dydxprotocol.clob.MsgPlaceOrder":                                 {},
		"/dydxprotocol.clob.MsgPlaceOrderResponse":                         {},
		"/dydxprotocol.clob.MsgProposedOperations":                         {},
		"/dydxprotocol.clob.MsgProposedOperationsResponse":                 {},
		"/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration":          {},
		"/dydxprotocol.clob.MsgUpdateBlockRateLimitConfigurationResponse":  {},
		"/dydxprotocol.clob.MsgUpdateClobPair":                             {},
		"/dydxprotocol.clob.MsgUpdateClobPairResponse":                     {},
		"/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration":         {},
		"/dydxprotocol.clob.MsgUpdateEquityTierLimitConfigurationResponse": {},
		"/dydxprotocol.clob.MsgUpdateLiquidationsConfig":                   {},
		"/dydxprotocol.clob.MsgUpdateLiquidationsConfigResponse":           {},

		// delaymsg
		"/dydxprotocol.delaymsg.MsgDelayMessage":         {},
		"/dydxprotocol.delaymsg.MsgDelayMessageResponse": {},

		// feetiers
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams":         {},
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse": {},

		// perpetuals
		"/dydxprotocol.perpetuals.MsgAddPremiumVotes":               {},
		"/dydxprotocol.perpetuals.MsgAddPremiumVotesResponse":       {},
		"/dydxprotocol.perpetuals.MsgCreatePerpetual":               {},
		"/dydxprotocol.perpetuals.MsgCreatePerpetualResponse":       {},
		"/dydxprotocol.perpetuals.MsgSetLiquidityTier":              {},
		"/dydxprotocol.perpetuals.MsgSetLiquidityTierResponse":      {},
		"/dydxprotocol.perpetuals.MsgUpdateParams":                  {},
		"/dydxprotocol.perpetuals.MsgUpdateParamsResponse":          {},
		"/dydxprotocol.perpetuals.MsgUpdatePerpetualParams":         {},
		"/dydxprotocol.perpetuals.MsgUpdatePerpetualParamsResponse": {},

		// prices
		"/dydxprotocol.prices.MsgCreateOracleMarket":         {},
		"/dydxprotocol.prices.MsgCreateOracleMarketResponse": {},
		"/dydxprotocol.prices.MsgUpdateMarketPrices":         {},
		"/dydxprotocol.prices.MsgUpdateMarketPricesResponse": {},
		"/dydxprotocol.prices.MsgUpdateMarketParam":          {},
		"/dydxprotocol.prices.MsgUpdateMarketParamResponse":  {},

		// ratelimit
		"/dydxprotocol.ratelimit.MsgSetLimitParams":         {},
		"/dydxprotocol.ratelimit.MsgSetLimitParamsResponse": {},

		// sending
		"/dydxprotocol.sending.MsgCreateTransfer":                  {},
		"/dydxprotocol.sending.MsgCreateTransferResponse":          {},
		"/dydxprotocol.sending.MsgDepositToSubaccount":             {},
		"/dydxprotocol.sending.MsgDepositToSubaccountResponse":     {},
		"/dydxprotocol.sending.MsgWithdrawFromSubaccount":          {},
		"/dydxprotocol.sending.MsgWithdrawFromSubaccountResponse":  {},
		"/dydxprotocol.sending.MsgSendFromModuleToAccount":         {},
		"/dydxprotocol.sending.MsgSendFromModuleToAccountResponse": {},

		// stats
		"/dydxprotocol.stats.MsgUpdateParams":         {},
		"/dydxprotocol.stats.MsgUpdateParamsResponse": {},

		// vest
		"/dydxprotocol.vest.MsgSetVestEntry":            {},
		"/dydxprotocol.vest.MsgSetVestEntryResponse":    {},
		"/dydxprotocol.vest.MsgDeleteVestEntry":         {},
		"/dydxprotocol.vest.MsgDeleteVestEntryResponse": {},

		// rewards
		"/dydxprotocol.rewards.MsgUpdateParams":         {},
		"/dydxprotocol.rewards.MsgUpdateParamsResponse": {},

		// ibc.applications
		"/ibc.applications.transfer.v1.MsgTransfer":             {},
		"/ibc.applications.transfer.v1.MsgTransferResponse":     {},
		"/ibc.applications.transfer.v1.MsgUpdateParams":         {},
		"/ibc.applications.transfer.v1.MsgUpdateParamsResponse": {},
		"/ibc.applications.transfer.v1.TransferAuthorization":   {},

		// ibc.core.channel
		"/ibc.core.channel.v1.Channel":                          {},
		"/ibc.core.channel.v1.Counterparty":                     {},
		"/ibc.core.channel.v1.MsgAcknowledgement":               {},
		"/ibc.core.channel.v1.MsgAcknowledgementResponse":       {},
		"/ibc.core.channel.v1.MsgChannelCloseConfirm":           {},
		"/ibc.core.channel.v1.MsgChannelCloseConfirmResponse":   {},
		"/ibc.core.channel.v1.MsgChannelCloseInit":              {},
		"/ibc.core.channel.v1.MsgChannelCloseInitResponse":      {},
		"/ibc.core.channel.v1.MsgChannelOpenAck":                {},
		"/ibc.core.channel.v1.MsgChannelOpenAckResponse":        {},
		"/ibc.core.channel.v1.MsgChannelOpenConfirm":            {},
		"/ibc.core.channel.v1.MsgChannelOpenConfirmResponse":    {},
		"/ibc.core.channel.v1.MsgChannelOpenInit":               {},
		"/ibc.core.channel.v1.MsgChannelOpenInitResponse":       {},
		"/ibc.core.channel.v1.MsgChannelOpenTry":                {},
		"/ibc.core.channel.v1.MsgChannelOpenTryResponse":        {},
		"/ibc.core.channel.v1.MsgRecvPacket":                    {},
		"/ibc.core.channel.v1.MsgRecvPacketResponse":            {},
		"/ibc.core.channel.v1.MsgTimeout":                       {},
		"/ibc.core.channel.v1.MsgTimeoutOnClose":                {},
		"/ibc.core.channel.v1.MsgTimeoutOnCloseResponse":        {},
		"/ibc.core.channel.v1.MsgTimeoutResponse":               {},
		"/ibc.core.channel.v1.Packet":                           {},
		"/ibc.core.channel.v1.MsgChannelUpgradeAck":             {},
		"/ibc.core.channel.v1.MsgChannelUpgradeAckResponse":     {},
		"/ibc.core.channel.v1.MsgChannelUpgradeCancel":          {},
		"/ibc.core.channel.v1.MsgChannelUpgradeCancelResponse":  {},
		"/ibc.core.channel.v1.MsgChannelUpgradeConfirm":         {},
		"/ibc.core.channel.v1.MsgChannelUpgradeConfirmResponse": {},
		"/ibc.core.channel.v1.MsgChannelUpgradeInit":            {},
		"/ibc.core.channel.v1.MsgChannelUpgradeInitResponse":    {},
		"/ibc.core.channel.v1.MsgChannelUpgradeOpen":            {},
		"/ibc.core.channel.v1.MsgChannelUpgradeOpenResponse":    {},
		"/ibc.core.channel.v1.MsgChannelUpgradeTimeout":         {},
		"/ibc.core.channel.v1.MsgChannelUpgradeTimeoutResponse": {},
		"/ibc.core.channel.v1.MsgChannelUpgradeTry":             {},
		"/ibc.core.channel.v1.MsgChannelUpgradeTryResponse":     {},
		"/ibc.core.channel.v1.MsgPruneAcknowledgements":         {},
		"/ibc.core.channel.v1.MsgPruneAcknowledgementsResponse": {},
		"/ibc.core.channel.v1.MsgUpdateParams":                  {},
		"/ibc.core.channel.v1.MsgUpdateParamsResponse":          {},

		// ibc.core.client
		"/ibc.core.client.v1.ClientUpdateProposal":          {},
		"/ibc.core.client.v1.Height":                        {},
		"/ibc.core.client.v1.MsgCreateClient":               {},
		"/ibc.core.client.v1.MsgCreateClientResponse":       {},
		"/ibc.core.client.v1.MsgIBCSoftwareUpgrade":         {},
		"/ibc.core.client.v1.MsgIBCSoftwareUpgradeResponse": {},
		"/ibc.core.client.v1.MsgRecoverClient":              {},
		"/ibc.core.client.v1.MsgRecoverClientResponse":      {},
		"/ibc.core.client.v1.MsgSubmitMisbehaviour":         {},
		"/ibc.core.client.v1.MsgSubmitMisbehaviourResponse": {},
		"/ibc.core.client.v1.MsgUpdateClient":               {},
		"/ibc.core.client.v1.MsgUpdateClientResponse":       {},
		"/ibc.core.client.v1.MsgUpgradeClient":              {},
		"/ibc.core.client.v1.MsgUpgradeClientResponse":      {},
		"/ibc.core.client.v1.MsgUpdateParams":               {},
		"/ibc.core.client.v1.MsgUpdateParamsResponse":       {},
		"/ibc.core.client.v1.UpgradeProposal":               {},

		// ibc.core.commitment
		"/ibc.core.commitment.v1.MerklePath":   {},
		"/ibc.core.commitment.v1.MerklePrefix": {},
		"/ibc.core.commitment.v1.MerkleProof":  {},
		"/ibc.core.commitment.v1.MerkleRoot":   {},

		// ibc.core.connection
		"/ibc.core.connection.v1.ConnectionEnd":                    {},
		"/ibc.core.connection.v1.Counterparty":                     {},
		"/ibc.core.connection.v1.MsgConnectionOpenAck":             {},
		"/ibc.core.connection.v1.MsgConnectionOpenAckResponse":     {},
		"/ibc.core.connection.v1.MsgConnectionOpenConfirm":         {},
		"/ibc.core.connection.v1.MsgConnectionOpenConfirmResponse": {},
		"/ibc.core.connection.v1.MsgConnectionOpenInit":            {},
		"/ibc.core.connection.v1.MsgConnectionOpenInitResponse":    {},
		"/ibc.core.connection.v1.MsgConnectionOpenTry":             {},
		"/ibc.core.connection.v1.MsgConnectionOpenTryResponse":     {},
		"/ibc.core.connection.v1.MsgUpdateParams":                  {},
		"/ibc.core.connection.v1.MsgUpdateParamsResponse":          {},

		// ibc.lightclients
		"/ibc.lightclients.localhost.v2.ClientState":     {},
		"/ibc.lightclients.tendermint.v1.ClientState":    {},
		"/ibc.lightclients.tendermint.v1.ConsensusState": {},
		"/ibc.lightclients.tendermint.v1.Header":         {},
		"/ibc.lightclients.tendermint.v1.Misbehaviour":   {},

		// ica messages
		// Note: the `interchain_accounts.controller` messages are not actually used by the app,
		// since ICA Controller Keeper is initialized as nil.
		// However, since the ica.AppModuleBasic{} needs to be passed to basic_mananger as a whole, these messages
		// registered in the interface registry.
		"/ibc.applications.interchain_accounts.v1.InterchainAccount":                               {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTx":                            {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgSendTxResponse":                    {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccount":         {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgRegisterInterchainAccountResponse": {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParams":                      {},
		"/ibc.applications.interchain_accounts.controller.v1.MsgUpdateParamsResponse":              {},
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams":                            {},
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse":                    {},
		"/ibc.applications.interchain_accounts.host.v1.MsgModuleQuerySafe":			    {},
		"/ibc.applications.interchain_accounts.host.v1.MsgModuleQuerySafeResponse":            	    {},
	}

	// DisallowMsgs are messages that cannot be externally submitted.
	DisallowMsgs = lib.MergeAllMapsMustHaveDistinctKeys(
		AppInjectedMsgSamples,
		InternalMsgSamplesAll,
		NestedMsgSamples,
		UnsupportedMsgSamples,
	)

	// AllowMsgs are messages that can be externally submitted.
	AllowMsgs = NormalMsgs
)
