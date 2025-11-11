package msgs

import (
	upgrade "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensus "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcconn "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	accountplus "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	affiliates "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	blocktime "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	bridge "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clob "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	delaymsg "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	feetiers "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	govplus "github.com/dydxprotocol/v4-chain/protocol/x/govplus/types"
	listing "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	perpetuals "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	ratelimit "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	revshare "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	rewards "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	sending "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	stats "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	vault "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	vest "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
)

var (
	// InternalMsgSamplesAll are msgs that are used only used internally.
	InternalMsgSamplesAll = lib.MergeAllMapsMustHaveDistinctKeys(InternalMsgSamplesGovAuth)

	// InternalMsgSamplesGovAuth are msgs that are used only used internally.
	// GovAuth means that these messages must originate from the gov module and
	// signed by gov module account.
	// InternalMsgSamplesAll are msgs that are used only used internally.
	InternalMsgSamplesGovAuth = lib.MergeAllMapsMustHaveDistinctKeys(
		InternalMsgSamplesDefault,
		InternalMsgSamplesDydxCustom,
	)

	// CosmosSDK default modules
	InternalMsgSamplesDefault = map[string]sdk.Msg{
		// auth
		"/cosmos.auth.v1beta1.MsgUpdateParams": &auth.MsgUpdateParams{},

		// bank
		"/cosmos.bank.v1beta1.MsgSetSendEnabled":         &bank.MsgSetSendEnabled{},
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse": nil,
		"/cosmos.bank.v1beta1.MsgUpdateParams":           &bank.MsgUpdateParams{},
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse":   nil,

		// consensus
		"/cosmos.consensus.v1.MsgUpdateParams":         &consensus.MsgUpdateParams{},
		"/cosmos.consensus.v1.MsgUpdateParamsResponse": nil,

		// crisis
		"/cosmos.crisis.v1beta1.MsgUpdateParams":         &crisis.MsgUpdateParams{},
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse": nil,

		// distribution
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend":         &distribution.MsgCommunityPoolSpend{},
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse": nil,
		"/cosmos.distribution.v1beta1.MsgUpdateParams":               &distribution.MsgUpdateParams{},
		"/cosmos.distribution.v1beta1.MsgUpdateParamsResponse":       nil,

		// gov
		"/cosmos.gov.v1.MsgExecLegacyContent":         &gov.MsgExecLegacyContent{},
		"/cosmos.gov.v1.MsgExecLegacyContentResponse": nil,
		"/cosmos.gov.v1.MsgUpdateParams":              &gov.MsgUpdateParams{},
		"/cosmos.gov.v1.MsgUpdateParamsResponse":      nil,

		// slashing
		"/cosmos.slashing.v1beta1.MsgUpdateParams":         &slashing.MsgUpdateParams{},
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse": nil,

		// staking
		"/cosmos.staking.v1beta1.MsgSetProposers":         &staking.MsgSetProposers{},
		"/cosmos.staking.v1beta1.MsgSetProposersResponse": nil,
		"/cosmos.staking.v1beta1.MsgUpdateParams":         &staking.MsgUpdateParams{},
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse": nil,

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade":           &upgrade.MsgCancelUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse":   nil,
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade":         &upgrade.MsgSoftwareUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse": nil,

		// ibc
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams":         &icahosttypes.MsgUpdateParams{},
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse": nil,
		"/ibc.applications.transfer.v1.MsgUpdateParams":                         &ibctransfer.MsgUpdateParams{},
		"/ibc.applications.transfer.v1.MsgUpdateParamsResponse":                 nil,
		"/ibc.core.client.v1.MsgUpdateParams":                                   &ibcclient.MsgUpdateParams{},
		"/ibc.core.client.v1.MsgUpdateParamsResponse":                           nil,
		"/ibc.core.connection.v1.MsgUpdateParams":                               &ibcconn.MsgUpdateParams{},
		"/ibc.core.connection.v1.MsgUpdateParamsResponse":                       nil,
	}

	// Custom modules
	InternalMsgSamplesDydxCustom = map[string]sdk.Msg{
		// affiliates
		"/dydxprotocol.affiliates.MsgUpdateAffiliateTiers":              &affiliates.MsgUpdateAffiliateTiers{},
		"/dydxprotocol.affiliates.MsgUpdateAffiliateTiersResponse":      nil,
		"/dydxprotocol.affiliates.MsgUpdateAffiliateWhitelist":          &affiliates.MsgUpdateAffiliateWhitelist{},
		"/dydxprotocol.affiliates.MsgUpdateAffiliateWhitelistResponse":  nil,
		"/dydxprotocol.affiliates.MsgUpdateAffiliateParameters":         &affiliates.MsgUpdateAffiliateParameters{},
		"/dydxprotocol.affiliates.MsgUpdateAffiliateParametersResponse": nil,
		"/dydxprotocol.affiliates.MsgUpdateAffiliateOverrides":          &affiliates.MsgUpdateAffiliateOverrides{},
		"/dydxprotocol.affiliates.MsgUpdateAffiliateOverridesResponse":  nil,

		// accountplus
		"/dydxprotocol.accountplus.MsgSetActiveState":         &accountplus.MsgSetActiveState{},
		"/dydxprotocol.accountplus.MsgSetActiveStateResponse": nil,

		// blocktime
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParams":          &blocktime.MsgUpdateDowntimeParams{},
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse":  nil,
		"/dydxprotocol.blocktime.MsgUpdateSynchronyParams":         &blocktime.MsgUpdateSynchronyParams{},
		"/dydxprotocol.blocktime.MsgUpdateSynchronyParamsResponse": nil,

		// bridge
		"/dydxprotocol.bridge.MsgCompleteBridge":              &bridge.MsgCompleteBridge{},
		"/dydxprotocol.bridge.MsgCompleteBridgeResponse":      nil,
		"/dydxprotocol.bridge.MsgUpdateEventParams":           &bridge.MsgUpdateEventParams{},
		"/dydxprotocol.bridge.MsgUpdateEventParamsResponse":   nil,
		"/dydxprotocol.bridge.MsgUpdateProposeParams":         &bridge.MsgUpdateProposeParams{},
		"/dydxprotocol.bridge.MsgUpdateProposeParamsResponse": nil,
		"/dydxprotocol.bridge.MsgUpdateSafetyParams":          &bridge.MsgUpdateSafetyParams{},
		"/dydxprotocol.bridge.MsgUpdateSafetyParamsResponse":  nil,

		// clob
		"/dydxprotocol.clob.MsgCreateClobPair":                             &clob.MsgCreateClobPair{},
		"/dydxprotocol.clob.MsgCreateClobPairResponse":                     nil,
		"/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration":          &clob.MsgUpdateBlockRateLimitConfiguration{},
		"/dydxprotocol.clob.MsgUpdateBlockRateLimitConfigurationResponse":  nil,
		"/dydxprotocol.clob.MsgUpdateClobPair":                             &clob.MsgUpdateClobPair{},
		"/dydxprotocol.clob.MsgUpdateClobPairResponse":                     nil,
		"/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration":         &clob.MsgUpdateEquityTierLimitConfiguration{},
		"/dydxprotocol.clob.MsgUpdateEquityTierLimitConfigurationResponse": nil,
		"/dydxprotocol.clob.MsgUpdateLiquidationsConfig":                   &clob.MsgUpdateLiquidationsConfig{},
		"/dydxprotocol.clob.MsgUpdateLiquidationsConfigResponse":           nil,

		// delaymsg
		"/dydxprotocol.delaymsg.MsgDelayMessage":         &delaymsg.MsgDelayMessage{},
		"/dydxprotocol.delaymsg.MsgDelayMessageResponse": nil,

		// feetiers
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams":           &feetiers.MsgUpdatePerpetualFeeParams{},
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse":   nil,
		"/dydxprotocol.feetiers.MsgSetMarketFeeDiscountParams":         &feetiers.MsgSetMarketFeeDiscountParams{},
		"/dydxprotocol.feetiers.MsgSetMarketFeeDiscountParamsResponse": nil,
		"/dydxprotocol.feetiers.MsgSetStakingTiers":                    &feetiers.MsgSetStakingTiers{},
		"/dydxprotocol.feetiers.MsgSetStakingTiersResponse":            nil,

		// govplus
		"/dydxprotocol.govplus.MsgSlashValidator":         &govplus.MsgSlashValidator{},
		"/dydxprotocol.govplus.MsgSlashValidatorResponse": nil,

		// listing
		"/dydxprotocol.listing.MsgSetMarketsHardCap":                       &listing.MsgSetMarketsHardCap{},
		"/dydxprotocol.listing.MsgSetMarketsHardCapResponse":               nil,
		"/dydxprotocol.listing.MsgSetListingVaultDepositParams":            &listing.MsgSetListingVaultDepositParams{},
		"/dydxprotocol.listing.MsgSetListingVaultDepositParamsResponse":    nil,
		"/dydxprotocol.listing.MsgUpgradeIsolatedPerpetualToCross":         &listing.MsgUpgradeIsolatedPerpetualToCross{},
		"/dydxprotocol.listing.MsgUpgradeIsolatedPerpetualToCrossResponse": nil,

		// perpetuals
		"/dydxprotocol.perpetuals.MsgCreatePerpetual":               &perpetuals.MsgCreatePerpetual{},
		"/dydxprotocol.perpetuals.MsgCreatePerpetualResponse":       nil,
		"/dydxprotocol.perpetuals.MsgSetLiquidityTier":              &perpetuals.MsgSetLiquidityTier{},
		"/dydxprotocol.perpetuals.MsgSetLiquidityTierResponse":      nil,
		"/dydxprotocol.perpetuals.MsgUpdateParams":                  &perpetuals.MsgUpdateParams{},
		"/dydxprotocol.perpetuals.MsgUpdateParamsResponse":          nil,
		"/dydxprotocol.perpetuals.MsgUpdatePerpetualParams":         &perpetuals.MsgUpdatePerpetualParams{},
		"/dydxprotocol.perpetuals.MsgUpdatePerpetualParamsResponse": nil,

		// prices
		"/dydxprotocol.prices.MsgCreateOracleMarket":         &prices.MsgCreateOracleMarket{},
		"/dydxprotocol.prices.MsgCreateOracleMarketResponse": nil,
		"/dydxprotocol.prices.MsgUpdateMarketParam":          &prices.MsgUpdateMarketParam{},
		"/dydxprotocol.prices.MsgUpdateMarketParamResponse":  nil,

		// ratelimit
		"/dydxprotocol.ratelimit.MsgSetLimitParams":         &ratelimit.MsgSetLimitParams{},
		"/dydxprotocol.ratelimit.MsgSetLimitParamsResponse": nil,

		// revshare
		"/dydxprotocol.revshare.MsgSetMarketMapperRevShareDetailsForMarket":         &revshare.MsgSetMarketMapperRevShareDetailsForMarket{}, //nolint:lll
		"/dydxprotocol.revshare.MsgSetMarketMapperRevShareDetailsForMarketResponse": nil,
		"/dydxprotocol.revshare.MsgSetMarketMapperRevenueShare":                     &revshare.MsgSetMarketMapperRevenueShare{}, //nolint:lll
		"/dydxprotocol.revshare.MsgSetMarketMapperRevenueShareResponse":             nil,
		"/dydxprotocol.revshare.MsgSetOrderRouterRevShare":                          &revshare.MsgSetOrderRouterRevShare{}, //nolint:lll
		"/dydxprotocol.revshare.MsgSetOrderRouterRevShareResponse":                  nil,
		"/dydxprotocol.revshare.MsgUpdateUnconditionalRevShareConfig":               &revshare.MsgUpdateUnconditionalRevShareConfig{}, //nolint:lll
		"/dydxprotocol.revshare.MsgUpdateUnconditionalRevShareConfigResponse":       nil,

		// rewards
		"/dydxprotocol.rewards.MsgUpdateParams":         &rewards.MsgUpdateParams{},
		"/dydxprotocol.rewards.MsgUpdateParamsResponse": nil,

		// sending
		"/dydxprotocol.sending.MsgSendFromModuleToAccount":          &sending.MsgSendFromModuleToAccount{},
		"/dydxprotocol.sending.MsgSendFromModuleToAccountResponse":  nil,
		"/dydxprotocol.sending.MsgSendFromAccountToAccount":         &sending.MsgSendFromAccountToAccount{},
		"/dydxprotocol.sending.MsgSendFromAccountToAccountResponse": nil,

		// stats
		"/dydxprotocol.stats.MsgUpdateParams":         &stats.MsgUpdateParams{},
		"/dydxprotocol.stats.MsgUpdateParamsResponse": nil,

		// vault
		"/dydxprotocol.vault.MsgUnlockShares":                 &vault.MsgUnlockShares{},
		"/dydxprotocol.vault.MsgUnlockSharesResponse":         nil,
		"/dydxprotocol.vault.MsgUpdateOperatorParams":         &vault.MsgUpdateOperatorParams{},
		"/dydxprotocol.vault.MsgUpdateOperatorParamsResponse": nil,

		// vest
		"/dydxprotocol.vest.MsgSetVestEntry":            &vest.MsgSetVestEntry{},
		"/dydxprotocol.vest.MsgSetVestEntryResponse":    nil,
		"/dydxprotocol.vest.MsgDeleteVestEntry":         &vest.MsgDeleteVestEntry{},
		"/dydxprotocol.vest.MsgDeleteVestEntryResponse": nil,
	}
)
