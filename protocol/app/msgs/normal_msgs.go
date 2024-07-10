package msgs

import (
	evidence "cosmossdk.io/x/evidence/types"
	feegrant "cosmossdk.io/x/feegrant"
	wasm "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcconn "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibccore "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clob "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	sending "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	vault "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

var (
	// NormalMsgs are messages that can be submitted by external users.
	NormalMsgs = lib.MergeAllMapsMustHaveDistinctKeys(NormalMsgsDefault, NormalMsgsDydxCustom)

	// Default modules
	NormalMsgsDefault = map[string]sdk.Msg{
		// auth
		"/cosmos.auth.v1beta1.BaseAccount":      nil,
		"/cosmos.auth.v1beta1.ModuleAccount":    nil,
		"/cosmos.auth.v1beta1.ModuleCredential": nil,

		// authz
		"/cosmos.authz.v1beta1.GenericAuthorization": nil,
		"/cosmos.authz.v1beta1.MsgGrant":             &authz.MsgGrant{},
		"/cosmos.authz.v1beta1.MsgGrantResponse":     nil,
		"/cosmos.authz.v1beta1.MsgRevoke":            &authz.MsgRevoke{},
		"/cosmos.authz.v1beta1.MsgRevokeResponse":    nil,

		// bank
		"/cosmos.bank.v1beta1.MsgMultiSend":         &bank.MsgMultiSend{},
		"/cosmos.bank.v1beta1.MsgMultiSendResponse": nil,
		"/cosmos.bank.v1beta1.MsgSend":              &bank.MsgSend{},
		"/cosmos.bank.v1beta1.MsgSendResponse":      nil,
		"/cosmos.bank.v1beta1.SendAuthorization":    nil,
		"/cosmos.bank.v1beta1.Supply":               nil,

		// consensus

		// crisis
		"/cosmos.crisis.v1beta1.MsgVerifyInvariant":         &crisis.MsgVerifyInvariant{},
		"/cosmos.crisis.v1beta1.MsgVerifyInvariantResponse": nil,

		// crypto
		"/cosmos.crypto.ed25519.PrivKey":            nil,
		"/cosmos.crypto.ed25519.PubKey":             nil,
		"/cosmos.crypto.multisig.LegacyAminoPubKey": nil,
		"/cosmos.crypto.secp256k1.PrivKey":          nil,
		"/cosmos.crypto.secp256k1.PubKey":           nil,
		"/cosmos.crypto.secp256r1.PubKey":           nil,

		// distribution
		"/cosmos.distribution.v1beta1.CommunityPoolSpendProposal":             nil,
		"/cosmos.distribution.v1beta1.MsgDepositValidatorRewardsPool":         &distr.MsgDepositValidatorRewardsPool{},
		"/cosmos.distribution.v1beta1.MsgDepositValidatorRewardsPoolResponse": nil,
		"/cosmos.distribution.v1beta1.MsgFundCommunityPool":                   &distr.MsgFundCommunityPool{},
		"/cosmos.distribution.v1beta1.MsgFundCommunityPoolResponse":           nil,
		"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress":                  &distr.MsgSetWithdrawAddress{},
		"/cosmos.distribution.v1beta1.MsgSetWithdrawAddressResponse":          nil,
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":             &distr.MsgWithdrawDelegatorReward{},
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorRewardResponse":     nil,
		"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission":         &distr.MsgWithdrawValidatorCommission{},
		"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommissionResponse": nil,

		// evidence
		"/cosmos.evidence.v1beta1.Equivocation":              nil,
		"/cosmos.evidence.v1beta1.MsgSubmitEvidence":         &evidence.MsgSubmitEvidence{},
		"/cosmos.evidence.v1beta1.MsgSubmitEvidenceResponse": nil,

		// feegrant
		"/cosmos.feegrant.v1beta1.AllowedMsgAllowance":        nil,
		"/cosmos.feegrant.v1beta1.BasicAllowance":             nil,
		"/cosmos.feegrant.v1beta1.MsgGrantAllowance":          &feegrant.MsgGrantAllowance{},
		"/cosmos.feegrant.v1beta1.MsgGrantAllowanceResponse":  nil,
		"/cosmos.feegrant.v1beta1.MsgPruneAllowances":         &feegrant.MsgPruneAllowances{},
		"/cosmos.feegrant.v1beta1.MsgPruneAllowancesResponse": nil,
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowance":         &feegrant.MsgRevokeAllowance{},
		"/cosmos.feegrant.v1beta1.MsgRevokeAllowanceResponse": nil,
		"/cosmos.feegrant.v1beta1.PeriodicAllowance":          nil,

		// gov
		"/cosmos.gov.v1.MsgDeposit":                   &gov.MsgDeposit{},
		"/cosmos.gov.v1.MsgDepositResponse":           nil,
		"/cosmos.gov.v1.MsgVote":                      &gov.MsgVote{},
		"/cosmos.gov.v1.MsgVoteResponse":              nil,
		"/cosmos.gov.v1.MsgVoteWeighted":              &gov.MsgVoteWeighted{},
		"/cosmos.gov.v1.MsgVoteWeightedResponse":      nil,
		"/cosmos.gov.v1beta1.MsgDeposit":              &govbeta.MsgDeposit{},
		"/cosmos.gov.v1beta1.MsgDepositResponse":      nil,
		"/cosmos.gov.v1beta1.MsgVote":                 &govbeta.MsgVote{},
		"/cosmos.gov.v1beta1.MsgVoteResponse":         nil,
		"/cosmos.gov.v1beta1.MsgVoteWeighted":         &govbeta.MsgVoteWeighted{},
		"/cosmos.gov.v1beta1.MsgVoteWeightedResponse": nil,
		"/cosmos.gov.v1beta1.TextProposal":            nil,

		// params
		"/cosmos.params.v1beta1.ParameterChangeProposal": nil,

		// slashing
		"/cosmos.slashing.v1beta1.MsgUnjail":         &slashing.MsgUnjail{},
		"/cosmos.slashing.v1beta1.MsgUnjailResponse": nil,

		// staking
		"/cosmos.staking.v1beta1.MsgBeginRedelegate":                   &staking.MsgBeginRedelegate{},
		"/cosmos.staking.v1beta1.MsgBeginRedelegateResponse":           nil,
		"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegation":         &staking.MsgCancelUnbondingDelegation{},
		"/cosmos.staking.v1beta1.MsgCancelUnbondingDelegationResponse": nil,
		"/cosmos.staking.v1beta1.MsgCreateValidator":                   &staking.MsgCreateValidator{},
		"/cosmos.staking.v1beta1.MsgCreateValidatorResponse":           nil,
		"/cosmos.staking.v1beta1.MsgDelegate":                          &staking.MsgDelegate{},
		"/cosmos.staking.v1beta1.MsgDelegateResponse":                  nil,
		"/cosmos.staking.v1beta1.MsgEditValidator":                     &staking.MsgEditValidator{},
		"/cosmos.staking.v1beta1.MsgEditValidatorResponse":             nil,
		"/cosmos.staking.v1beta1.MsgUndelegate":                        &staking.MsgUndelegate{},
		"/cosmos.staking.v1beta1.MsgUndelegateResponse":                nil,
		"/cosmos.staking.v1beta1.StakeAuthorization":                   nil,

		// tx
		"/cosmos.tx.v1beta1.Tx": nil,

		// upgrade
		"/cosmos.upgrade.v1beta1.CancelSoftwareUpgradeProposal": nil,
		"/cosmos.upgrade.v1beta1.SoftwareUpgradeProposal":       nil,

		// ibc.applications
		"/ibc.applications.transfer.v1.MsgTransfer":           &ibctransfer.MsgTransfer{},
		"/ibc.applications.transfer.v1.MsgTransferResponse":   nil,
		"/ibc.applications.transfer.v1.TransferAuthorization": nil,

		// ibc.core.channel
		"/ibc.core.channel.v1.Channel":                        nil,
		"/ibc.core.channel.v1.Counterparty":                   nil,
		"/ibc.core.channel.v1.MsgAcknowledgement":             &ibccore.MsgAcknowledgement{},
		"/ibc.core.channel.v1.MsgAcknowledgementResponse":     nil,
		"/ibc.core.channel.v1.MsgChannelCloseConfirm":         &ibccore.MsgChannelCloseConfirm{},
		"/ibc.core.channel.v1.MsgChannelCloseConfirmResponse": nil,
		"/ibc.core.channel.v1.MsgChannelCloseInit":            &ibccore.MsgChannelCloseInit{},
		"/ibc.core.channel.v1.MsgChannelCloseInitResponse":    nil,
		"/ibc.core.channel.v1.MsgChannelOpenAck":              &ibccore.MsgChannelOpenAck{},
		"/ibc.core.channel.v1.MsgChannelOpenAckResponse":      nil,
		"/ibc.core.channel.v1.MsgChannelOpenConfirm":          &ibccore.MsgChannelOpenConfirm{},
		"/ibc.core.channel.v1.MsgChannelOpenConfirmResponse":  nil,
		"/ibc.core.channel.v1.MsgChannelOpenInit":             &ibccore.MsgChannelOpenInit{},
		"/ibc.core.channel.v1.MsgChannelOpenInitResponse":     nil,
		"/ibc.core.channel.v1.MsgChannelOpenTry":              &ibccore.MsgChannelOpenTry{},
		"/ibc.core.channel.v1.MsgChannelOpenTryResponse":      nil,
		"/ibc.core.channel.v1.MsgRecvPacket":                  &ibccore.MsgRecvPacket{},
		"/ibc.core.channel.v1.MsgRecvPacketResponse":          nil,
		"/ibc.core.channel.v1.MsgTimeout":                     &ibccore.MsgTimeout{},
		"/ibc.core.channel.v1.MsgTimeoutOnClose":              &ibccore.MsgTimeoutOnClose{},
		"/ibc.core.channel.v1.MsgTimeoutOnCloseResponse":      nil,
		"/ibc.core.channel.v1.MsgTimeoutResponse":             nil,
		"/ibc.core.channel.v1.Packet":                         nil,

		// ibc.core.client
		"/ibc.core.client.v1.ClientUpdateProposal":          nil,
		"/ibc.core.client.v1.Height":                        nil,
		"/ibc.core.client.v1.MsgCreateClient":               &ibcclient.MsgCreateClient{},
		"/ibc.core.client.v1.MsgCreateClientResponse":       nil,
		"/ibc.core.client.v1.MsgIBCSoftwareUpgrade":         &ibcclient.MsgIBCSoftwareUpgrade{},
		"/ibc.core.client.v1.MsgIBCSoftwareUpgradeResponse": nil,
		"/ibc.core.client.v1.MsgRecoverClient":              &ibcclient.MsgRecoverClient{},
		"/ibc.core.client.v1.MsgRecoverClientResponse":      nil,
		"/ibc.core.client.v1.MsgSubmitMisbehaviour":         &ibcclient.MsgSubmitMisbehaviour{}, //nolint:staticcheck
		"/ibc.core.client.v1.MsgSubmitMisbehaviourResponse": nil,
		// TODO(CORE-851): Move MsgUpdateClient and MsgUpgradeClient to unsupported_msgs once v4.0.0 upgrade has
		// been completed and Cosmos 0.50 performs well.
		"/ibc.core.client.v1.MsgUpdateClient":          &ibcclient.MsgUpdateClient{},
		"/ibc.core.client.v1.MsgUpdateClientResponse":  nil,
		"/ibc.core.client.v1.MsgUpgradeClient":         &ibcclient.MsgUpgradeClient{},
		"/ibc.core.client.v1.MsgUpgradeClientResponse": nil,
		"/ibc.core.client.v1.UpgradeProposal":          nil,

		// ibc.core.commitment
		"/ibc.core.commitment.v1.MerklePath":   nil,
		"/ibc.core.commitment.v1.MerklePrefix": nil,
		"/ibc.core.commitment.v1.MerkleProof":  nil,
		"/ibc.core.commitment.v1.MerkleRoot":   nil,

		// ibc.core.connection
		"/ibc.core.connection.v1.ConnectionEnd":                    nil,
		"/ibc.core.connection.v1.Counterparty":                     nil,
		"/ibc.core.connection.v1.MsgConnectionOpenAck":             &ibcconn.MsgConnectionOpenAck{},
		"/ibc.core.connection.v1.MsgConnectionOpenAckResponse":     nil,
		"/ibc.core.connection.v1.MsgConnectionOpenConfirm":         &ibcconn.MsgConnectionOpenConfirm{},
		"/ibc.core.connection.v1.MsgConnectionOpenConfirmResponse": nil,
		"/ibc.core.connection.v1.MsgConnectionOpenInit":            &ibcconn.MsgConnectionOpenInit{},
		"/ibc.core.connection.v1.MsgConnectionOpenInitResponse":    nil,
		"/ibc.core.connection.v1.MsgConnectionOpenTry":             &ibcconn.MsgConnectionOpenTry{},
		"/ibc.core.connection.v1.MsgConnectionOpenTryResponse":     nil,

		// ibc.lightclients
		"/ibc.lightclients.localhost.v2.ClientState":     nil,
		"/ibc.lightclients.tendermint.v1.ClientState":    nil,
		"/ibc.lightclients.tendermint.v1.ConsensusState": nil,
		"/ibc.lightclients.tendermint.v1.Header":         nil,
		"/ibc.lightclients.tendermint.v1.Misbehaviour":   nil,

		// ica
		"/ibc.applications.interchain_accounts.v1.InterchainAccount": nil,

		// wasm
		// TODO(OTE-461): Audit and remove unnecessary messages
		"/cosmwasm.wasm.v1.AcceptedMessageKeysFilter":                  nil,
		"/cosmwasm.wasm.v1.AcceptedMessagesFilter":                     nil,
		"/cosmwasm.wasm.v1.AllowAllMessagesFilter":                     nil,
		"/cosmwasm.wasm.v1.ClearAdminProposal":                         nil,
		"/cosmwasm.wasm.v1.CombinedLimit":                              nil,
		"/cosmwasm.wasm.v1.ContractExecutionAuthorization":             nil,
		"/cosmwasm.wasm.v1.ContractMigrationAuthorization":             nil,
		"/cosmwasm.wasm.v1.ExecuteContractProposal":                    nil,
		"/cosmwasm.wasm.v1.InstantiateContract2Proposal":               nil,
		"/cosmwasm.wasm.v1.InstantiateContractProposal":                nil,
		"/cosmwasm.wasm.v1.MaxCallsLimit":                              nil,
		"/cosmwasm.wasm.v1.MaxFundsLimit":                              nil,
		"/cosmwasm.wasm.v1.MigrateContractProposal":                    nil,
		"/cosmwasm.wasm.v1.MsgAddCodeUploadParamsAddresses":            &wasm.MsgAddCodeUploadParamsAddresses{},
		"/cosmwasm.wasm.v1.MsgAddCodeUploadParamsAddressesResponse":    nil,
		"/cosmwasm.wasm.v1.MsgClearAdmin":                              &wasm.MsgClearAdmin{},
		"/cosmwasm.wasm.v1.MsgClearAdminResponse":                      nil,
		"/cosmwasm.wasm.v1.MsgExecuteContract":                         &wasm.MsgExecuteContract{},
		"/cosmwasm.wasm.v1.MsgExecuteContractResponse":                 nil,
		"/cosmwasm.wasm.v1.MsgIBCCloseChannel":                         &wasm.MsgIBCCloseChannel{},
		"/cosmwasm.wasm.v1.MsgIBCSend":                                 &wasm.MsgIBCSend{},
		"/cosmwasm.wasm.v1.MsgInstantiateContract":                     &wasm.MsgInstantiateContract{},
		"/cosmwasm.wasm.v1.MsgInstantiateContract2":                    &wasm.MsgInstantiateContract2{},
		"/cosmwasm.wasm.v1.MsgInstantiateContract2Response":            nil,
		"/cosmwasm.wasm.v1.MsgInstantiateContractResponse":             nil,
		"/cosmwasm.wasm.v1.MsgMigrateContract":                         &wasm.MsgMigrateContract{},
		"/cosmwasm.wasm.v1.MsgMigrateContractResponse":                 nil,
		"/cosmwasm.wasm.v1.MsgPinCodes":                                &wasm.MsgPinCodes{},
		"/cosmwasm.wasm.v1.MsgPinCodesResponse":                        nil,
		"/cosmwasm.wasm.v1.MsgRemoveCodeUploadParamsAddresses":         &wasm.MsgRemoveCodeUploadParamsAddresses{},
		"/cosmwasm.wasm.v1.MsgRemoveCodeUploadParamsAddressesResponse": nil,
		"/cosmwasm.wasm.v1.MsgStoreAndInstantiateContract":             &wasm.MsgStoreAndInstantiateContract{},
		"/cosmwasm.wasm.v1.MsgStoreAndInstantiateContractResponse":     nil,
		"/cosmwasm.wasm.v1.MsgStoreAndMigrateContract":                 &wasm.MsgStoreAndMigrateContract{},
		"/cosmwasm.wasm.v1.MsgStoreAndMigrateContractResponse":         nil,
		"/cosmwasm.wasm.v1.MsgStoreCode":                               &wasm.MsgStoreCode{},
		"/cosmwasm.wasm.v1.MsgStoreCodeResponse":                       nil,
		"/cosmwasm.wasm.v1.MsgSudoContract":                            &wasm.MsgSudoContract{},
		"/cosmwasm.wasm.v1.MsgSudoContractResponse":                    nil,
		"/cosmwasm.wasm.v1.MsgUnpinCodes":                              &wasm.MsgUnpinCodes{},
		"/cosmwasm.wasm.v1.MsgUnpinCodesResponse":                      nil,
		"/cosmwasm.wasm.v1.MsgUpdateAdmin":                             &wasm.MsgUpdateAdmin{},
		"/cosmwasm.wasm.v1.MsgUpdateAdminResponse":                     nil,
		"/cosmwasm.wasm.v1.MsgUpdateContractLabel":                     &wasm.MsgUpdateContractLabel{},
		"/cosmwasm.wasm.v1.MsgUpdateContractLabelResponse":             nil,
		"/cosmwasm.wasm.v1.MsgUpdateInstantiateConfig":                 &wasm.MsgUpdateInstantiateConfig{},
		"/cosmwasm.wasm.v1.MsgUpdateInstantiateConfigResponse":         nil,
		"/cosmwasm.wasm.v1.MsgUpdateParams":                            &wasm.MsgUpdateParams{},
		"/cosmwasm.wasm.v1.MsgUpdateParamsResponse":                    nil,
		"/cosmwasm.wasm.v1.PinCodesProposal":                           nil,
		"/cosmwasm.wasm.v1.StoreAndInstantiateContractProposal":        nil,
		"/cosmwasm.wasm.v1.StoreCodeAuthorization":                     nil,
		"/cosmwasm.wasm.v1.StoreCodeProposal":                          nil,
		"/cosmwasm.wasm.v1.SudoContractProposal":                       nil,
		"/cosmwasm.wasm.v1.UnpinCodesProposal":                         nil,
		"/cosmwasm.wasm.v1.UpdateAdminProposal":                        nil,
		"/cosmwasm.wasm.v1.UpdateInstantiateConfigProposal":            nil,
	}

	// Custom modules
	NormalMsgsDydxCustom = map[string]sdk.Msg{
		// clob
		"/dydxprotocol.clob.MsgBatchCancel":         &clob.MsgBatchCancel{},
		"/dydxprotocol.clob.MsgBatchCancelResponse": nil,
		"/dydxprotocol.clob.MsgCancelOrder":         &clob.MsgCancelOrder{},
		"/dydxprotocol.clob.MsgCancelOrderResponse": nil,
		"/dydxprotocol.clob.MsgPlaceOrder":          &clob.MsgPlaceOrder{},
		"/dydxprotocol.clob.MsgPlaceOrderResponse":  nil,

		// perpetuals

		// prices

		// sending
		"/dydxprotocol.sending.MsgCreateTransfer":                 &sending.MsgCreateTransfer{},
		"/dydxprotocol.sending.MsgCreateTransferResponse":         nil,
		"/dydxprotocol.sending.MsgDepositToSubaccount":            &sending.MsgDepositToSubaccount{},
		"/dydxprotocol.sending.MsgDepositToSubaccountResponse":    nil,
		"/dydxprotocol.sending.MsgWithdrawFromSubaccount":         &sending.MsgWithdrawFromSubaccount{},
		"/dydxprotocol.sending.MsgWithdrawFromSubaccountResponse": nil,

		// vault
		"/dydxprotocol.vault.MsgDepositToVault":            &vault.MsgDepositToVault{},
		"/dydxprotocol.vault.MsgDepositToVaultResponse":    nil,
		"/dydxprotocol.vault.MsgWithdrawFromVault":         &vault.MsgWithdrawFromVault{},
		"/dydxprotocol.vault.MsgWithdrawFromVaultResponse": nil,
	}
)
