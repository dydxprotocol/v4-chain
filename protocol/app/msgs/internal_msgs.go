package msgs

import (
	upgrade "cosmossdk.io/x/upgrade/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	blocktime "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	clob "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	delaymsg "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	feetiers "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	perpetuals "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	prices "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimit "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	rewards "github.com/StreamFinance-Protocol/stream-chain/protocol/x/rewards/types"
	sending "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	stats "github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/types"
	vest "github.com/StreamFinance-Protocol/stream-chain/protocol/x/vest/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensus "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcconn "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
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

		// slashing
		"/cosmos.slashing.v1beta1.MsgUpdateParams":         &slashing.MsgUpdateParams{},
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse": nil,

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
		// blocktime
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParams":         &blocktime.MsgUpdateDowntimeParams{},
		"/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse": nil,

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

		// ratelimit
		"/dydxprotocol.ratelimit.MsgSetLimitParams":         &ratelimit.MsgSetLimitParams{},
		"/dydxprotocol.ratelimit.MsgSetLimitParamsResponse": nil,

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
