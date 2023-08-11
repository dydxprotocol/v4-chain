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
	"github.com/dydxprotocol/v4/lib/maps"
)

var (
	// InternalMsgSamplesAll are msgs that are used only used internally.
	InternalMsgSamplesAll = maps.MergeAllMapsMustHaveDistinctKeys(InternalMsgSamplesGovAuth)

	// InternalMsgSamplesGovAuth are msgs that are used only used internally.
	// GovAuth means that these messages must originate from the gov module and
	// signed by gov module account.
	InternalMsgSamplesGovAuth = map[string]sdk.Msg{
		// MsgUpdateParams
		"/cosmos.auth.v1beta1.MsgUpdateParams":                 &auth.MsgUpdateParams{},
		"/cosmos.bank.v1beta1.MsgUpdateParams":                 &bank.MsgUpdateParams{},
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse":         nil,
		"/cosmos.consensus.v1.MsgUpdateParams":                 &consensus.MsgUpdateParams{},
		"/cosmos.consensus.v1.MsgUpdateParamsResponse":         nil,
		"/cosmos.crisis.v1beta1.MsgUpdateParams":               &crisis.MsgUpdateParams{},
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse":       nil,
		"/cosmos.distribution.v1beta1.MsgUpdateParams":         &distribution.MsgUpdateParams{},
		"/cosmos.distribution.v1beta1.MsgUpdateParamsResponse": nil,
		"/cosmos.gov.v1.MsgUpdateParams":                       &gov.MsgUpdateParams{},
		"/cosmos.gov.v1.MsgUpdateParamsResponse":               nil,
		"/cosmos.slashing.v1beta1.MsgUpdateParams":             &slashing.MsgUpdateParams{},
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse":     nil,
		"/cosmos.staking.v1beta1.MsgUpdateParams":              &staking.MsgUpdateParams{},
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse":      nil,

		// bank
		"/cosmos.bank.v1beta1.MsgSetSendEnabled":         &bank.MsgSetSendEnabled{},
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse": nil,

		// distribution
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend":         &distribution.MsgCommunityPoolSpend{},
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse": nil,

		// gov
		"/cosmos.gov.v1.MsgExecLegacyContent":         &gov.MsgExecLegacyContent{},
		"/cosmos.gov.v1.MsgExecLegacyContentResponse": nil,

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade":           &upgrade.MsgCancelUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse":   nil,
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade":         &upgrade.MsgSoftwareUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse": nil,
	}
)
