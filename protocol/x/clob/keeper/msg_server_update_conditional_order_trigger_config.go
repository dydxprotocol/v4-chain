package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateConditionalOrderTriggerConfig sets the conditional-order EndBlocker work conditional-order trigger mitigation
// configuration. It is authority-gated (governance), which is the on-chain activation lever for
// the mitigation: the binary ships with the config disabled (legacy behavior, rolling-deploy
// safe), and a governance proposal invoking this message flips Enabled and sets the per-block
// trigger/removal budgets and admission caps at a coordinated height. Enabling deterministically
// backfills the trigger-price index (handled inside SetConditionalOrderTriggerConfig).
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

	// Normalizes zero-valued numeric fields to their defaults and performs the
	// disabled->enabled index backfill (inside SetConditionalOrderTriggerConfig).
	k.Keeper.SetConditionalOrderTriggerConfigParams(
		ctx,
		msg.Enabled,
		msg.MaxTriggersPerBlock,
		msg.MaxRemovalsPerBlock,
		msg.MaxUntriggeredConditionalOrdersGlobal,
		msg.MaxUntriggeredConditionalOrdersPerSubaccount,
	)

	return &types.MsgUpdateConditionalOrderTriggerConfigResponse{}, nil
}
