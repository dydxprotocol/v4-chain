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
)

// IsInternalMsg returns true if the given msg is an internal message.
func IsInternalMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// MsgUpdateParams
		*auth.MsgUpdateParams,
		*bank.MsgUpdateParams,
		*consensus.MsgUpdateParams,
		*crisis.MsgUpdateParams,
		*distribution.MsgUpdateParams,
		*gov.MsgUpdateParams,
		*slashing.MsgUpdateParams,
		*staking.MsgUpdateParams,

		// bank
		*bank.MsgSetSendEnabled,

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
