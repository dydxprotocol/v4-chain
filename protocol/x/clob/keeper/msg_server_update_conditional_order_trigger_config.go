package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateConditionalOrderTriggerConfig sets the governance-tunable budgets and caps for the
// conditional-order trigger mitigation. It is authority-gated (governance). The mitigation itself
// is always active (activated by the state-breaking upgrade that builds the trigger-price index);
// this message only retunes the per-block trigger/removal budgets and the admission caps.
func (k msgServer) UpdateConditionalOrderTriggerConfig(
	goCtx context.Context,
	msg *types.MsgUpdateConditionalOrderTriggerConfig,
) (resp *types.MsgUpdateConditionalOrderTriggerConfigResponse, err error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	// Stateless bounds: reject out-of-range budgets / caps. Enforced here
	// (the guaranteed execution point for a governance message) in addition to being registered as
	// the message's ValidateBasic.
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Normalizes zero-valued numeric fields to sane package defaults.
	k.Keeper.SetConditionalOrderTriggerConfigParams(
		ctx,
		msg.MaxTriggersPerBlock,
		msg.MaxRemovalsPerBlock,
		msg.MaxUntriggeredConditionalOrdersGlobal,
		msg.MaxUntriggeredConditionalOrdersPerSubaccount,
	)

	return &types.MsgUpdateConditionalOrderTriggerConfigResponse{}, nil
}
