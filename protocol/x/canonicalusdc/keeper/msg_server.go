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
	if m.keeper.PendingCount(ctx) != 0 {
		if next.Mode == types.Mode_MODE_DISABLED {
			return nil, fmt.Errorf("%w: cannot disable with pending settlements", types.ErrInvalidControls)
		}
		if !current.SameRoute(next) {
			return nil, fmt.Errorf("%w: cannot change route identity with pending settlements", types.ErrInvalidControls)
		}
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
