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
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse",

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
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams",
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse",

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

		// rewards
		"/dydxprotocol.rewards.MsgUpdateParams",
		"/dydxprotocol.rewards.MsgUpdateParamsResponse",

		// sending
		"/dydxprotocol.sending.MsgSendFromModuleToAccount",
		"/dydxprotocol.sending.MsgSendFromModuleToAccountResponse",

		// stats
		"/dydxprotocol.stats.MsgUpdateParams",
		"/dydxprotocol.stats.MsgUpdateParamsResponse",

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
