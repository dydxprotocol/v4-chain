package msgs_test

import (
	"sort"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/msgs"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
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
		"/cosmos.staking.v1beta1.MsgUpdateParams",
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse",

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade",
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse",

		// blocktime
		"/klyraprotocol.blocktime.MsgUpdateDowntimeParams",
		"/klyraprotocol.blocktime.MsgUpdateDowntimeParamsResponse",

		// clob
		"/klyraprotocol.clob.MsgCreateClobPair",
		"/klyraprotocol.clob.MsgCreateClobPairResponse",
		"/klyraprotocol.clob.MsgUpdateBlockRateLimitConfiguration",
		"/klyraprotocol.clob.MsgUpdateBlockRateLimitConfigurationResponse",
		"/klyraprotocol.clob.MsgUpdateClobPair",
		"/klyraprotocol.clob.MsgUpdateClobPairResponse",
		"/klyraprotocol.clob.MsgUpdateEquityTierLimitConfiguration",
		"/klyraprotocol.clob.MsgUpdateEquityTierLimitConfigurationResponse",
		"/klyraprotocol.clob.MsgUpdateLiquidationsConfig",
		"/klyraprotocol.clob.MsgUpdateLiquidationsConfigResponse",

		// delaymsg
		"/klyraprotocol.delaymsg.MsgDelayMessage",
		"/klyraprotocol.delaymsg.MsgDelayMessageResponse",

		// feetiers
		"/klyraprotocol.feetiers.MsgUpdatePerpetualFeeParams",
		"/klyraprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse",

		// govplus
		"/klyraprotocol.govplus.MsgSlashValidator",
		"/klyraprotocol.govplus.MsgSlashValidatorResponse",

		// perpeutals
		"/klyraprotocol.perpetuals.MsgCreatePerpetual",
		"/klyraprotocol.perpetuals.MsgCreatePerpetualResponse",
		"/klyraprotocol.perpetuals.MsgSetLiquidityTier",
		"/klyraprotocol.perpetuals.MsgSetLiquidityTierResponse",
		"/klyraprotocol.perpetuals.MsgUpdateParams",
		"/klyraprotocol.perpetuals.MsgUpdateParamsResponse",
		"/klyraprotocol.perpetuals.MsgUpdatePerpetualParams",
		"/klyraprotocol.perpetuals.MsgUpdatePerpetualParamsResponse",

		// prices
		"/klyraprotocol.prices.MsgCreateOracleMarket",
		"/klyraprotocol.prices.MsgCreateOracleMarketResponse",
		"/klyraprotocol.prices.MsgUpdateMarketParam",
		"/klyraprotocol.prices.MsgUpdateMarketParamResponse",

		// ratelimit
		"/klyraprotocol.ratelimit.MsgSetLimitParams",
		"/klyraprotocol.ratelimit.MsgSetLimitParamsResponse",

		// rewards
		"/klyraprotocol.rewards.MsgUpdateParams",
		"/klyraprotocol.rewards.MsgUpdateParamsResponse",

		// sending
		"/klyraprotocol.sending.MsgSendFromModuleToAccount",
		"/klyraprotocol.sending.MsgSendFromModuleToAccountResponse",

		// stats
		"/klyraprotocol.stats.MsgUpdateParams",
		"/klyraprotocol.stats.MsgUpdateParamsResponse",

		// vest
		"/klyraprotocol.vest.MsgDeleteVestEntry",
		"/klyraprotocol.vest.MsgDeleteVestEntryResponse",
		"/klyraprotocol.vest.MsgSetVestEntry",
		"/klyraprotocol.vest.MsgSetVestEntryResponse",

		// ibc
		"/ibc.applications.interchain_accounts.host.v1.MsgModuleQuerySafe",
		"/ibc.applications.interchain_accounts.host.v1.MsgModuleQuerySafeResponse",
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
