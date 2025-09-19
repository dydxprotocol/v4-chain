package msgs_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestNormalMsgs_Key(t *testing.T) {
	expectedMsgs := []string{
		// auth
		"/cosmos.auth.v1beta1.BaseAccount",
		"/cosmos.auth.v1beta1.ModuleAccount",
		"/cosmos.auth.v1beta1.ModuleCredential",

		// authz
		"/cosmos.authz.v1beta1.GenericAuthorization",
		"/cosmos.authz.v1beta1.MsgGrant",
		"/cosmos.authz.v1beta1.MsgGrantResponse",
		"/cosmos.authz.v1beta1.MsgRevoke",
		"/cosmos.authz.v1beta1.MsgRevokeResponse",

		// bank
		"/cosmos.bank.v1beta1.MsgMultiSend",
		"/cosmos.bank.v1beta1.MsgMultiSendResponse",
		"/cosmos.bank.v1beta1.MsgSend",
		"/cosmos.bank.v1beta1.MsgSendResponse",
		"/cosmos.bank.v1beta1.SendAuthorization",
		"/cosmos.bank.v1beta1.Supply",

		// consensus

		// crypto
		"/cosmos.crypto.ed25519.PrivKey",
		"/cosmos.crypto.ed25519.PubKey",
		"/cosmos.crypto.multisig.LegacyAminoPubKey",
		"/cosmos.crypto.secp256k1.PrivKey",
		"/cosmos.crypto.secp256k1.PubKey",
		"/cosmos.crypto.secp256r1.PubKey",

		// distribution
		"/cosmos.distribution.v1beta1.CommunityPoolSpendProposal",
		"/cosmos.distribution.v1beta1.MsgDepositValidatorRewardsPool",
		"/cosmos.distribution.v1beta1.MsgDepositValidatorRewardsPoolResponse",
		"/cosmos.distribution.v1beta1.MsgFundCommunityPool",
		"/cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse",
		"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
		"/cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse",
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse",
		"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission",
		"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse",

		// evidence
		"/cosmos.evidence.v1beta1.Equivocation",
		"/cosmos.evidence.v1beta1.MsgSubmitEvidence",
		"/cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse",

		// feegrant
		"/cosmos.feegrant.v1beta1.AllowedMsgAllowance",
		"/cosmos.feegrant.v1beta1.BasicAllowance",
		"/cosmos.feegrant.v1beta1.MsgGrantAllowance",
		"/cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse",
		"/cosmos.feegrant.v1beta1.MsgPruneAllowances",
		"/cosmos.feegrant.v1beta1.MsgPruneAllowancesResponse",
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowance",
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse",
		"/cosmos.feegrant.v1beta1.PeriodicAllowance",

		// gov
		"/cosmos.gov.v1.MsgDeposit",
		"/cosmos.gov.v1.MsgDepositResponse",
		"/cosmos.gov.v1.MsgVote",
		"/cosmos.gov.v1.MsgVoteResponse",
		"/cosmos.gov.v1.MsgVoteWeighted",
		"/cosmos.gov.v1.MsgVoteWeightedResponse",
		"/cosmos.gov.v1beta1.MsgDeposit",
		"/cosmos.gov.v1beta1.MsgDepositResponse",
		"/cosmos.gov.v1beta1.MsgVote",
		"/cosmos.gov.v1beta1.MsgVoteResponse",
		"/cosmos.gov.v1beta1.MsgVoteWeighted",
		"/cosmos.gov.v1beta1.MsgVoteWeightedResponse",
		"/cosmos.gov.v1beta1.TextProposal",

		// params
		"/cosmos.params.v1beta1.ParameterChangeProposal",

		// slashing
		"/cosmos.slashing.v1beta1.MsgUnjail",
		"/cosmos.slashing.v1beta1.MsgUnjailResponse",

		// staking
		"/cosmos.staking.v1beta1.MsgBeginRedelegate",
		"/cosmos.staking.v1beta1.MsgBeginRedelegateResponse",
		"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation",
		"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegationResponse",
		"/cosmos.staking.v1beta1.MsgCreateValidator",
		"/cosmos.staking.v1beta1.MsgCreateValidatorResponse",
		"/cosmos.staking.v1beta1.MsgDelegate",
		"/cosmos.staking.v1beta1.MsgDelegateResponse",
		"/cosmos.staking.v1beta1.MsgEditValidator",
		"/cosmos.staking.v1beta1.MsgEditValidatorResponse",
		"/cosmos.staking.v1beta1.MsgUndelegate",
		"/cosmos.staking.v1beta1.MsgUndelegateResponse",
		"/cosmos.staking.v1beta1.StakeAuthorization",

		// tx
		"/cosmos.tx.v1beta1.Tx",

		// upgrade
		"/cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal",
		"/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal",

		// accountplus
		"/dydxprotocol.accountplus.MsgAddAuthenticator",
		"/dydxprotocol.accountplus.MsgAddAuthenticatorResponse",
		"/dydxprotocol.accountplus.MsgRemoveAuthenticator",
		"/dydxprotocol.accountplus.MsgRemoveAuthenticatorResponse",
		"/dydxprotocol.accountplus.TxExtension",

		// affiliates
		"/dydxprotocol.affiliates.MsgRegisterAffiliate",
		"/dydxprotocol.affiliates.MsgRegisterAffiliateResponse",

		// clob
		"/dydxprotocol.clob.MsgBatchCancel",
		"/dydxprotocol.clob.MsgBatchCancelResponse",
		"/dydxprotocol.clob.MsgCancelOrder",
		"/dydxprotocol.clob.MsgCancelOrderResponse",
		"/dydxprotocol.clob.MsgPlaceOrder",
		"/dydxprotocol.clob.MsgPlaceOrderResponse",
		"/dydxprotocol.clob.MsgUpdateLeverage",
		"/dydxprotocol.clob.MsgUpdateLeverageResponse",

		// listing
		"/dydxprotocol.listing.MsgCreateMarketPermissionless",
		"/dydxprotocol.listing.MsgCreateMarketPermissionlessResponse",

		// perpetuals

		// prices

		// sending
		"/dydxprotocol.sending.MsgCreateTransfer",
		"/dydxprotocol.sending.MsgCreateTransferResponse",
		"/dydxprotocol.sending.MsgDepositToSubaccount",
		"/dydxprotocol.sending.MsgDepositToSubaccountResponse",
		"/dydxprotocol.sending.MsgWithdrawFromSubaccount",
		"/dydxprotocol.sending.MsgWithdrawFromSubaccountResponse",

		// vault
		"/dydxprotocol.vault.MsgAllocateToVault",
		"/dydxprotocol.vault.MsgAllocateToVaultResponse",
		"/dydxprotocol.vault.MsgDepositToMegavault",
		"/dydxprotocol.vault.MsgDepositToMegavaultResponse",
		"/dydxprotocol.vault.MsgRetrieveFromVault",
		"/dydxprotocol.vault.MsgRetrieveFromVaultResponse",
		"/dydxprotocol.vault.MsgSetVaultParams",
		"/dydxprotocol.vault.MsgSetVaultParamsResponse",
		"/dydxprotocol.vault.MsgUpdateDefaultQuotingParams",
		"/dydxprotocol.vault.MsgUpdateDefaultQuotingParamsResponse",
		"/dydxprotocol.vault.MsgWithdrawFromMegavault",
		"/dydxprotocol.vault.MsgWithdrawFromMegavaultResponse",

		// ibc application module: ICA
		"/ibc.applications.interchain_accounts.v1.InterchainAccount",

		// ibc.applications
		"/ibc.applications.transfer.v1.MsgTransfer",
		"/ibc.applications.transfer.v1.MsgTransferResponse",
		"/ibc.applications.transfer.v1.TransferAuthorization",

		// ibc.core.channel
		"/ibc.core.channel.v1.Channel",
		"/ibc.core.channel.v1.Counterparty",
		"/ibc.core.channel.v1.MsgAcknowledgement",
		"/ibc.core.channel.v1.MsgAcknowledgementResponse",
		"/ibc.core.channel.v1.MsgChannelCloseConfirm",
		"/ibc.core.channel.v1.MsgChannelCloseConfirmResponse",
		"/ibc.core.channel.v1.MsgChannelCloseInit",
		"/ibc.core.channel.v1.MsgChannelCloseInitResponse",
		"/ibc.core.channel.v1.MsgChannelOpenAck",
		"/ibc.core.channel.v1.MsgChannelOpenAckResponse",
		"/ibc.core.channel.v1.MsgChannelOpenConfirm",
		"/ibc.core.channel.v1.MsgChannelOpenConfirmResponse",
		"/ibc.core.channel.v1.MsgChannelOpenInit",
		"/ibc.core.channel.v1.MsgChannelOpenInitResponse",
		"/ibc.core.channel.v1.MsgChannelOpenTry",
		"/ibc.core.channel.v1.MsgChannelOpenTryResponse",
		"/ibc.core.channel.v1.MsgRecvPacket",
		"/ibc.core.channel.v1.MsgRecvPacketResponse",
		"/ibc.core.channel.v1.MsgTimeout",
		"/ibc.core.channel.v1.MsgTimeoutOnClose",
		"/ibc.core.channel.v1.MsgTimeoutOnCloseResponse",
		"/ibc.core.channel.v1.MsgTimeoutResponse",
		"/ibc.core.channel.v1.Packet",

		// ibc.core.client
		"/ibc.core.client.v1.ClientUpdateProposal",
		"/ibc.core.client.v1.Height",
		"/ibc.core.client.v1.MsgCreateClient",
		"/ibc.core.client.v1.MsgCreateClientResponse",
		"/ibc.core.client.v1.MsgIBCSoftwareUpgrade",
		"/ibc.core.client.v1.MsgIBCSoftwareUpgradeResponse",
		"/ibc.core.client.v1.MsgRecoverClient",
		"/ibc.core.client.v1.MsgRecoverClientResponse",
		"/ibc.core.client.v1.MsgSubmitMisbehaviour",
		"/ibc.core.client.v1.MsgSubmitMisbehaviourResponse",
		"/ibc.core.client.v1.MsgUpdateClient",
		"/ibc.core.client.v1.MsgUpdateClientResponse",
		"/ibc.core.client.v1.MsgUpgradeClient",
		"/ibc.core.client.v1.MsgUpgradeClientResponse",
		"/ibc.core.client.v1.UpgradeProposal",

		// ibc.core.commitment
		"/ibc.core.commitment.v1.MerklePath",
		"/ibc.core.commitment.v1.MerklePrefix",
		"/ibc.core.commitment.v1.MerkleProof",
		"/ibc.core.commitment.v1.MerkleRoot",

		// ibc.core.connection
		"/ibc.core.connection.v1.ConnectionEnd",
		"/ibc.core.connection.v1.Counterparty",
		"/ibc.core.connection.v1.MsgConnectionOpenAck",
		"/ibc.core.connection.v1.MsgConnectionOpenAckResponse",
		"/ibc.core.connection.v1.MsgConnectionOpenConfirm",
		"/ibc.core.connection.v1.MsgConnectionOpenConfirmResponse",
		"/ibc.core.connection.v1.MsgConnectionOpenInit",
		"/ibc.core.connection.v1.MsgConnectionOpenInitResponse",
		"/ibc.core.connection.v1.MsgConnectionOpenTry",
		"/ibc.core.connection.v1.MsgConnectionOpenTryResponse",

		// ibc.lightclients
		"/ibc.lightclients.localhost.v2.ClientState",
		"/ibc.lightclients.tendermint.v1.ClientState",
		"/ibc.lightclients.tendermint.v1.ConsensusState",
		"/ibc.lightclients.tendermint.v1.Header",
		"/ibc.lightclients.tendermint.v1.Misbehaviour",

		// slinky marketmap messages
		"/slinky.marketmap.v1.MsgCreateMarkets",
		"/slinky.marketmap.v1.MsgCreateMarketsResponse",
		"/slinky.marketmap.v1.MsgParams",
		"/slinky.marketmap.v1.MsgParamsResponse",
		"/slinky.marketmap.v1.MsgRemoveMarketAuthorities",
		"/slinky.marketmap.v1.MsgRemoveMarketAuthoritiesResponse",
		"/slinky.marketmap.v1.MsgRemoveMarkets",
		"/slinky.marketmap.v1.MsgRemoveMarketsResponse",
		"/slinky.marketmap.v1.MsgUpdateMarkets",
		"/slinky.marketmap.v1.MsgUpdateMarketsResponse",
		"/slinky.marketmap.v1.MsgUpsertMarkets",
		"/slinky.marketmap.v1.MsgUpsertMarketsResponse",
	}

	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.NormalMsgs))
}

func TestNormalMsgs_Value(t *testing.T) {
	validateMsgValue(t, msgs.NormalMsgs)
}
