package ante

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
	blocktime "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	bridge "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	delaymsg "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	feetiers "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	rewards "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	stats "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	vest "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
)

// IsInternalMsg returns true if the given msg is an internal message.
func IsInternalMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// MsgUpdateParams
		*auth.MsgUpdateParams,
		*blocktime.MsgUpdateDowntimeParams,
		*bank.MsgUpdateParams,
		*consensus.MsgUpdateParams,
		*crisis.MsgUpdateParams,
		*feetiers.MsgUpdatePerpetualFeeParams,
		*distribution.MsgUpdateParams,
		*gov.MsgUpdateParams,
		*slashing.MsgUpdateParams,
		*staking.MsgUpdateParams,
		*stats.MsgUpdateParams,
		*rewards.MsgUpdateParams,
		*vest.MsgSetVestEntry,
		*vest.MsgDeleteVestEntry,

		// bank
		*bank.MsgSetSendEnabled,

		// bridge
		*bridge.MsgCompleteBridge,
		*bridge.MsgUpdateEventParams,
		*bridge.MsgUpdateProposeParams,
		*bridge.MsgUpdateSafetyParams,

		// delaymsg
		*delaymsg.MsgDelayMessage,

		// distribution
		*distribution.MsgCommunityPoolSpend,

		// gov
		*gov.MsgExecLegacyContent,

		// upgrade
		*upgrade.MsgCancelUpgrade,
		*upgrade.MsgSoftwareUpgrade:

		return true

	default:
		return false
	}
}
