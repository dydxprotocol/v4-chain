package msgs_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/msgs"
	"github.com/dydxprotocol/v4-chain/protocol/lib/maps"
	"github.com/stretchr/testify/require"
)

func TestInternalMsgSamples_All_Key(t *testing.T) {
	expectedAllInternalMsgs := maps.MergeAllMapsMustHaveDistinctKeys(msgs.InternalMsgSamplesGovAuth)
	require.Equal(t, expectedAllInternalMsgs, msgs.InternalMsgSamplesAll)
}

func TestInternalMsgSamples_All_Value(t *testing.T) {
	validateSampleMsgValue(t, msgs.InternalMsgSamplesAll)
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
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParams",
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse",

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

		// delaymsg
		"/dydxprotocol.delaymsg.MsgDelayMessage",
		"/dydxprotocol.delaymsg.MsgDelayMessageResponse",

		// feetiers
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams",
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse",

		// perpeutals
		"/dydxprotocol.perpetuals.MsgCreatePerpetual",
		"/dydxprotocol.perpetuals.MsgCreatePerpetualResponse",

		// prices
		"/dydxprotocol.prices.MsgCreateOracleMarket",
		"/dydxprotocol.prices.MsgCreateOracleMarketResponse",

		// rewards
		"/dydxprotocol.rewards.MsgUpdateParams",
		"/dydxprotocol.rewards.MsgUpdateParamsResponse",

		// stats
		"/dydxprotocol.stats.MsgUpdateParams",
		"/dydxprotocol.stats.MsgUpdateParamsResponse",

		// vest
		"/dydxprotocol.vest.MsgDeleteVestEntry",
		"/dydxprotocol.vest.MsgDeleteVestEntryResponse",
		"/dydxprotocol.vest.MsgSetVestEntry",
		"/dydxprotocol.vest.MsgSetVestEntryResponse",
	}

	require.Equal(t, expectedMsgs, maps.GetSortedKeys(msgs.InternalMsgSamplesGovAuth))
}
