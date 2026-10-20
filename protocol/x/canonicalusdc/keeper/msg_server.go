package keeper

import (
	"context"
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper *Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) UpdateControls(
	goCtx context.Context,
	msg *types.MsgUpdateControls,
) (*types.MsgUpdateControlsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg == nil || !m.keeper.HasAuthority(msg.Authority) {
		return nil, types.ErrUnauthorized
	}
	next := msg.Controls
	if err := next.Validate(); err != nil {
		return nil, err
	}
	if err := m.keeper.ValidateRouteBindings(ctx, next); err != nil {
		return nil, err
	}
	current := m.keeper.GetControls(ctx)
	if err := m.keeper.validateTransition(ctx, current, next); err != nil {
		return nil, err
	}
	cacheCtx, write := ctx.CacheContext()
	if err := m.keeper.SetControls(cacheCtx, next); err != nil {
		return nil, err
	}
	if next.Mode == types.Mode_MODE_GRADUAL {
		if err := m.keeper.enterGradual(cacheCtx, current.Mode); err != nil {
			return nil, err
		}
	}
	write()
	return &types.MsgUpdateControlsResponse{}, nil
}

// enterGradual snapshots Noble backing when GRADUAL is entered, verifies the
// backing invariant at that point, and requires the Noble route to stay open
// while any Noble backing or Noble settlement remains.
func (k Keeper) enterGradual(ctx sdk.Context, previous types.Mode) error {
	if previous != types.Mode_MODE_GRADUAL {
		if err := k.snapshotNobleBacking(ctx); err != nil {
			return err
		}
		if reason, broken := BackingInvariant(k)(ctx); broken {
			return fmt.Errorf("%w: %s", types.ErrInvalidLedger, reason)
		}
	}
	if k.GetControls(ctx).NobleWithdrawalsEnabled {
		return nil
	}
	amounts, err := parseBackingAmounts(k.GetLedger(ctx))
	if err != nil {
		return err
	}
	if !amounts.noble.IsZero() || !amounts.pendingNoble.IsZero() {
		return types.ErrNobleBackingRemains
	}
	return nil
}

// validateTransition keeps the Injective reserve and in-flight settlements on
// the route that created them. While any is outstanding, routing cannot be
// disabled (DISABLED would let Injective-backed supply leave through the legacy
// Noble route) and route identity cannot change (the reserve would be stranded
// in the old physical denom). The pending cap cannot drop below the live
// count, or exported genesis would fail validation.
func (k Keeper) validateTransition(ctx sdk.Context, current, next types.Controls) error {
	amounts, err := parseBackingAmounts(k.GetLedger(ctx))
	if err != nil {
		return err
	}
	pendingCount := k.PendingCount(ctx)
	outstanding := pendingCount != 0 || !amounts.injective.IsZero() ||
		!amounts.pendingInjective.IsZero() || !amounts.pendingNoble.IsZero()
	if outstanding && next.Mode == types.Mode_MODE_DISABLED {
		return fmt.Errorf("%w: cannot disable while Injective backing or settlements are outstanding",
			types.ErrInvalidControls)
	}
	// DISABLED binds no route, and it cannot be entered with anything outstanding.
	if outstanding && current.Mode != types.Mode_MODE_DISABLED && !current.SameRoute(next) {
		return fmt.Errorf("%w: cannot change route identity while Injective backing or settlements are outstanding",
			types.ErrInvalidControls)
	}
	if next.Mode != types.Mode_MODE_DISABLED && next.MaxPendingSettlements < pendingCount {
		return fmt.Errorf("%w: max pending settlements below current pending count %d",
			types.ErrInvalidControls, pendingCount)
	}
	return nil
}
