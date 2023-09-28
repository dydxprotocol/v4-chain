package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensus "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	blocktime "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	bridge "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clob "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	delaymsg "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	feetiers "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perpetuals "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	rewards "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	sending "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	stats "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	vest "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
)

var (
	// InternalMsgSamplesAll are msgs that are used only used internally.
	InternalMsgSamplesAll = lib.MergeAllMapsMustHaveDistinctKeys(InternalMsgSamplesGovAuth)

	// InternalMsgSamplesGovAuth are msgs that are used only used internally.
	// GovAuth means that these messages must originate from the gov module and
	// signed by gov module account.
	InternalMsgSamplesGovAuth = map[string]sdk.Msg{
		// ------- CosmosSDK default modules
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
		"/cosmos.staking.v1beta1.MsgUpdateParams":         &staking.MsgUpdateParams{},
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse": nil,

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade":           &upgrade.MsgCancelUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse":   nil,
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade":         &upgrade.MsgSoftwareUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse": nil,

		// ------- Custom modules
		// blocktime
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParams":         &blocktime.MsgUpdateDowntimeParams{},
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse": nil,

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
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams":         &feetiers.MsgUpdatePerpetualFeeParams{},
		"/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse": nil,

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

		// rewards
		"/dydxprotocol.rewards.MsgUpdateParams":         &rewards.MsgUpdateParams{},
		"/dydxprotocol.rewards.MsgUpdateParamsResponse": nil,

		// sending
		"/dydxprotocol.sending.MsgSendFromModuleToAccount":         &sending.MsgSendFromModuleToAccount{},
		"/dydxprotocol.sending.MsgSendFromModuleToAccountResponse": nil,

		// stats
		"/dydxprotocol.stats.MsgUpdateParams":         &stats.MsgUpdateParams{},
		"/dydxprotocol.stats.MsgUpdateParamsResponse": nil,

		// vest
		"/dydxprotocol.vest.MsgSetVestEntry":            &vest.MsgSetVestEntry{},
		"/dydxprotocol.vest.MsgSetVestEntryResponse":    nil,
		"/dydxprotocol.vest.MsgDeleteVestEntry":         &vest.MsgDeleteVestEntry{},
		"/dydxprotocol.vest.MsgDeleteVestEntryResponse": nil,
	}
)
