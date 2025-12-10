package ante

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

// IsInternalMsg returns true if the given msg is an internal message.
func IsInternalMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// ------- CosmosSDK default modules
		// auth
		*auth.MsgUpdateParams,

		// bank
		*bank.MsgSetSendEnabled,
		*bank.MsgUpdateParams,

		// consensus
		*consensus.MsgUpdateParams,

		// crisis
		*crisis.MsgUpdateParams,

		// distribution
		*distribution.MsgCommunityPoolSpend,
		*distribution.MsgUpdateParams,

		// gov
		*gov.MsgExecLegacyContent,
		*gov.MsgUpdateParams,

		// slashing
		*slashing.MsgUpdateParams,

		// staking
		*staking.MsgSetProposers,
		*staking.MsgUpdateParams,

		// upgrade
		*upgrade.MsgCancelUpgrade,
		*upgrade.MsgSoftwareUpgrade,

		// ------- Custom modules
		// accountplus
		*accountplus.MsgSetActiveState,

		// blocktime
		*blocktime.MsgUpdateDowntimeParams,
		*blocktime.MsgUpdateSynchronyParams,
		// bridge
		*bridge.MsgCompleteBridge,
		*bridge.MsgUpdateEventParams,
		*bridge.MsgUpdateProposeParams,
		*bridge.MsgUpdateSafetyParams,

		// clob
		*clob.MsgCreateClobPair,
		*clob.MsgUpdateBlockRateLimitConfiguration,
		*clob.MsgUpdateClobPair,
		*clob.MsgUpdateEquityTierLimitConfiguration,
		*clob.MsgUpdateLiquidationsConfig,

		// delaymsg
		*delaymsg.MsgDelayMessage,

		// feetiers
		*feetiers.MsgSetStakingTiers,
		*feetiers.MsgUpdatePerpetualFeeParams,
		*feetiers.MsgSetMarketFeeDiscountParams,

		// govplus
		*govplus.MsgSlashValidator,

		// listing
		*listing.MsgSetMarketsHardCap,
		*listing.MsgSetListingVaultDepositParams,
		*listing.MsgUpgradeIsolatedPerpetualToCross,

		// perpetuals
		*perpetuals.MsgCreatePerpetual,
		*perpetuals.MsgSetLiquidityTier,
		*perpetuals.MsgUpdateParams,
		*perpetuals.MsgUpdatePerpetualParams,

		// prices
		*prices.MsgCreateOracleMarket,
		*prices.MsgUpdateMarketParam,

		// ratelimit
		*ratelimit.MsgSetLimitParams,
		*ratelimit.MsgSetLimitParamsResponse,

		// revshare
		*revshare.MsgSetMarketMapperRevenueShare,
		*revshare.MsgSetMarketMapperRevShareDetailsForMarket,
		*revshare.MsgUpdateUnconditionalRevShareConfig,
		*revshare.MsgSetOrderRouterRevShare,

		// rewards
		*rewards.MsgUpdateParams,

		// sending
		*sending.MsgSendFromModuleToAccount,
		*sending.MsgSendFromAccountToAccount,

		// stats
		*stats.MsgUpdateParams,

		// vault
		*vault.MsgUnlockShares,
		*vault.MsgUpdateOperatorParams,

		// vest
		*vest.MsgDeleteVestEntry,
		*vest.MsgSetVestEntry,

		// ibc
		*icahosttypes.MsgUpdateParams,
		*ibctransfer.MsgUpdateParams,
		*ibcclient.MsgUpdateParams,
		*ibcconn.MsgUpdateParams,

		// affiliates
		*affiliates.MsgUpdateAffiliateTiers,
		*affiliates.MsgUpdateAffiliateWhitelist,
		*affiliates.MsgUpdateAffiliateOverrides,
		*affiliates.MsgUpdateAffiliateParameters:

		return true

	default:
		return false
	}
}
