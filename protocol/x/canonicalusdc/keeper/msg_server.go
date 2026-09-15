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
	if err := msg.Controls.Validate(); err != nil {
		return nil, err
	}
	if err := m.keeper.ValidateRouteBindings(ctx, msg.Controls); err != nil {
		return nil, err
	}
	current := m.keeper.GetControls(ctx)
	if m.keeper.PendingCount(ctx) != 0 {
		if msg.Controls.Mode == types.Mode_MODE_DISABLED {
			return nil, fmt.Errorf("%w: cannot disable with pending settlements", types.ErrInvalidControls)
		}
		if current.NobleChannel != msg.Controls.NobleChannel ||
			current.InjectiveChannel != msg.Controls.InjectiveChannel ||
			current.NoblePacketDenom != msg.Controls.NoblePacketDenom ||
			current.InjectivePacketDenom != msg.Controls.InjectivePacketDenom ||
			current.InjectiveDenom != msg.Controls.InjectiveDenom ||
			current.LogicalDenom != msg.Controls.LogicalDenom {
			return nil, fmt.Errorf("%w: cannot change route identity with pending settlements", types.ErrInvalidControls)
		}
	}
	cacheCtx, write := ctx.CacheContext()
	if err := m.keeper.SetControls(cacheCtx, msg.Controls); err != nil {
		return nil, err
	}
	if current.Mode == types.Mode_MODE_DISABLED && msg.Controls.Mode != types.Mode_MODE_DISABLED {
		if err := m.keeper.snapshotNobleBacking(cacheCtx); err != nil {
			return nil, err
		}
	} else if err := m.keeper.SetLedger(cacheCtx, m.keeper.GetLedger(cacheCtx)); err != nil {
		return nil, err
	}
	if reason, broken := BackingInvariant(*m.keeper)(cacheCtx); broken {
		return nil, fmt.Errorf("%w: %s", types.ErrInvalidLedger, reason)
	}
	write()
	return &types.MsgUpdateControlsResponse{}, nil
}
