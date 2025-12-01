package msgs_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestInternalMsgSamples_All_Key(t *testing.T) {
	expectedAllInternalMsgs := lib.MergeAllMapsMustHaveDistinctKeys(msgs.InternalMsgSamplesGovAuth)
	require.Equal(t, expectedAllInternalMsgs, msgs.InternalMsgSamplesAll)
}

func TestInternalMsgSamples_All_Value(t *testing.T) {
	validateMsgValue(t, msgs.InternalMsgSamplesAll)
}

func TestInternalMsgSamples_Gov_Key(t *testing.T) {
	expectedMsgs := []string{
		// auth
		"/cosmos.auth.v1beta1.MsgUpdateParams",

		// bank
		"/cosmos.bank.v1beta1.MsgSetSendEnabled",
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse",
		"/cosmos.bank.v1beta1.MsgUpdateParams",
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse",

		// consensus
		"/cosmos.consensus.v1.MsgUpdateParams",
		"/cosmos.consensus.v1.MsgUpdateParamsResponse",

		// crisis
		"/cosmos.crisis.v1beta1.MsgUpdateParams",
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse",

		// distribution
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend",
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse",
		"/cosmos.distribution.v1beta1.MsgUpdateParams",
		"/cosmos.distribution.v1beta1.MsgUpdateParamsResponse",

		// gov
		"/cosmos.gov.v1.MsgExecLegacyContent",
		"/cosmos.gov.v1.MsgExecLegacyContentResponse",
		"/cosmos.gov.v1.MsgUpdateParams",
		"/cosmos.gov.v1.MsgUpdateParamsResponse",

		// slashing
		"/cosmos.slashing.v1beta1.MsgUpdateParams",
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse",

		// staking
		"/cosmos.staking.v1beta1.MsgSetProposers",
		"/cosmos.staking.v1beta1.MsgSetProposersResponse",
		"/cosmos.staking.v1beta1.MsgUpdateParams",
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse",

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade",
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse",

		// accountplus
		"/dydxprotocol.accountplus.MsgSetActiveState",
		"/dydxprotocol.accountplus.MsgSetActiveStateResponse",

		// affiliates
		"/dydxprotocol.affiliates.MsgUpdateAffiliateOverrides",
		"/dydxprotocol.affiliates.MsgUpdateAffiliateOverridesResponse",
		"/dydxprotocol.affiliates.MsgUpdateAffiliateParameters",
		"/dydxprotocol.affiliates.MsgUpdateAffiliateParametersResponse",
		"/dydxprotocol.affiliates.MsgUpdateAffiliateTiers",
		"/dydxprotocol.affiliates.MsgUpdateAffiliateTiersResponse",
		"/dydxprotocol.affiliates.MsgUpdateAffiliateWhitelist",
		"/dydxprotocol.affiliates.MsgUpdateAffiliateWhitelistResponse",

		// blocktime
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse",
		"/dydxprotocol.blocktime.MsgUpdateSynchronyParams",
		"/dydxprotocol.blocktime.MsgUpdateSynchronyParamsResponse",

		// bridge
		"/dydxprotocol.bridge.MsgCompleteBridge",
		"/dydxprotocol.bridge.MsgCompleteBridgeResponse",
		"/dydxprotocol.bridge.MsgUpdateEventParams",
		"/dydxprotocol.bridge.MsgUpdateEventParamsResponse",
		"/dydxprotocol.bridge.MsgUpdateProposeParams",
		"/dydxprotocol.bridge.MsgUpdateProposeParamsResponse",
		"/dydxprotocol.bridge.MsgUpdateSafetyParams",
		"/dydxprotocol.bridge.MsgUpdateSafetyParamsResponse",

		// clob
		"/dydxprotocol.clob.MsgCreateClobPair",
		"/dydxprotocol.clob.MsgCreateClobPairResponse",
		"/dydxprotocol.clob.MsgUpdateBlockRateLimitConfiguration",
		"/dydxprotocol.clob.MsgUpdateBlockRateLimitConfigurationResponse",
		"/dydxprotocol.clob.MsgUpdateClobPair",
		"/dydxprotocol.clob.MsgUpdateClobPairResponse",
		"/dydxprotocol.clob.MsgUpdateEquityTierLimitConfiguration",
		"/dydxprotocol.clob.MsgUpdateEquityTierLimitConfigurationResponse",
		"/dydxprotocol.clob.MsgUpdateLiquidationsConfig",
		"/dydxprotocol.clob.MsgUpdateLiquidationsConfigResponse",

		// delaymsg
		"/dydxprotocol.delaymsg.MsgDelayMessage",
		"/dydxprotocol.delaymsg.MsgDelayMessageResponse",

		// feetiers
		"/dydxprotocol.feetiers.MsgSetMarketFeeDiscountParams",
		"/dydxprotocol.feetiers.MsgSetMarketFeeDiscountParamsResponse",
		"/dydxprotocol.feetiers.MsgSetStakingTiers",
		"/dydxprotocol.feetiers.MsgSetStakingTiersResponse",
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams",
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse",

		// govplus
		"/dydxprotocol.govplus.MsgSlashValidator",
		"/dydxprotocol.govplus.MsgSlashValidatorResponse",

		// listing
		"/dydxprotocol.listing.MsgSetListingVaultDepositParams",
		"/dydxprotocol.listing.MsgSetListingVaultDepositParamsResponse",
		"/dydxprotocol.listing.MsgSetMarketsHardCap",
		"/dydxprotocol.listing.MsgSetMarketsHardCapResponse",
		"/dydxprotocol.listing.MsgUpgradeIsolatedPerpetualToCross",
		"/dydxprotocol.listing.MsgUpgradeIsolatedPerpetualToCrossResponse",

		// perpeutals
		"/dydxprotocol.perpetuals.MsgCreatePerpetual",
		"/dydxprotocol.perpetuals.MsgCreatePerpetualResponse",
		"/dydxprotocol.perpetuals.MsgSetLiquidityTier",
		"/dydxprotocol.perpetuals.MsgSetLiquidityTierResponse",
		"/dydxprotocol.perpetuals.MsgUpdateParams",
		"/dydxprotocol.perpetuals.MsgUpdateParamsResponse",
		"/dydxprotocol.perpetuals.MsgUpdatePerpetualParams",
		"/dydxprotocol.perpetuals.MsgUpdatePerpetualParamsResponse",

		// prices
		"/dydxprotocol.prices.MsgCreateOracleMarket",
		"/dydxprotocol.prices.MsgCreateOracleMarketResponse",
		"/dydxprotocol.prices.MsgUpdateMarketParam",
		"/dydxprotocol.prices.MsgUpdateMarketParamResponse",

		// ratelimit
		"/dydxprotocol.ratelimit.MsgSetLimitParams",
		"/dydxprotocol.ratelimit.MsgSetLimitParamsResponse",

		// revshare
		"/dydxprotocol.revshare.MsgSetMarketMapperRevShareDetailsForMarket",
		"/dydxprotocol.revshare.MsgSetMarketMapperRevShareDetailsForMarketResponse",
		"/dydxprotocol.revshare.MsgSetMarketMapperRevenueShare",
		"/dydxprotocol.revshare.MsgSetMarketMapperRevenueShareResponse",
		"/dydxprotocol.revshare.MsgSetOrderRouterRevShare",
		"/dydxprotocol.revshare.MsgSetOrderRouterRevShareResponse",
		"/dydxprotocol.revshare.MsgUpdateUnconditionalRevShareConfig",
		"/dydxprotocol.revshare.MsgUpdateUnconditionalRevShareConfigResponse",

		// rewards
		"/dydxprotocol.rewards.MsgUpdateParams",
		"/dydxprotocol.rewards.MsgUpdateParamsResponse",

		// sending
		"/dydxprotocol.sending.MsgSendFromAccountToAccount",
		"/dydxprotocol.sending.MsgSendFromAccountToAccountResponse",
		"/dydxprotocol.sending.MsgSendFromModuleToAccount",
		"/dydxprotocol.sending.MsgSendFromModuleToAccountResponse",

		// stats
		"/dydxprotocol.stats.MsgUpdateParams",
		"/dydxprotocol.stats.MsgUpdateParamsResponse",

		// vault
		"/dydxprotocol.vault.MsgUnlockShares",
		"/dydxprotocol.vault.MsgUnlockSharesResponse",
		"/dydxprotocol.vault.MsgUpdateOperatorParams",
		"/dydxprotocol.vault.MsgUpdateOperatorParamsResponse",

		// vest
		"/dydxprotocol.vest.MsgDeleteVestEntry",
		"/dydxprotocol.vest.MsgDeleteVestEntryResponse",
		"/dydxprotocol.vest.MsgSetVestEntry",
		"/dydxprotocol.vest.MsgSetVestEntryResponse",

		// ibc
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams",
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse",
		"/ibc.applications.transfer.v1.MsgUpdateParams",
		"/ibc.applications.transfer.v1.MsgUpdateParamsResponse",
		"/ibc.core.client.v1.MsgUpdateParams",
		"/ibc.core.client.v1.MsgUpdateParamsResponse",
		"/ibc.core.connection.v1.MsgUpdateParams",
		"/ibc.core.connection.v1.MsgUpdateParamsResponse",
	}

	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.InternalMsgSamplesGovAuth))
}

func TestInternalMsgSamples_Gov_Value(t *testing.T) {
	validateMsgValue(t, msgs.InternalMsgSamplesGovAuth)
}
