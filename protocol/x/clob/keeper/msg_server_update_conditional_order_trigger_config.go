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

	// Stateless bounds: reject out-of-range budgets / caps. Enforced here
	// (the guaranteed execution point for a governance message) in addition to being registered as
	// the message's ValidateBasic.
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Stateful enable gate. The disabled->enabled transition backfills
	// the trigger-price index over the ENTIRE resting untriggered set in one consensus block
	// (unbudgeted O(N) — see BackfillConditionalOrderTriggerPriceIndex). Because admission caps are
	// off while disabled, that set can be arbitrarily large. Refuse to enable when the live count
	// exceeds the size the backfill can safely rebuild in one block, or exceeds the global admission
	// cap that will apply once enabled (which would otherwise freeze all new placements immediately).
	// Operators must drain (cancel / let expire) the resting set below the ceiling before enabling.
	if msg.Enabled && !k.Keeper.IsConditionalOrderTriggerConfigEnabled(ctx) {
		liveCount := k.Keeper.GetUntriggeredConditionalOrderCountGlobal(ctx)

		effectiveGlobalCap := msg.MaxUntriggeredConditionalOrdersGlobal
		if effectiveGlobalCap == 0 {
			// Mirror the setter's zero-normalization so the gate reasons about the cap that will
			// actually apply once enabled.
			effectiveGlobalCap = MaxUntriggeredConditionalOrdersGlobal
		}

		limit := uint32(MaxBackfillCardinality)
		if effectiveGlobalCap < limit {
			limit = effectiveGlobalCap
		}

		if liveCount > limit {
			return nil, errorsmod.Wrapf(
				types.ErrUntriggeredSetTooLargeToEnable,
				"live untriggered conditional order count %d exceeds the enable ceiling %d "+
					"(backfill bound %d, effective global cap %d); drain the resting set before enabling",
				liveCount,
				limit,
				MaxBackfillCardinality,
				effectiveGlobalCap,
			)
		}
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
