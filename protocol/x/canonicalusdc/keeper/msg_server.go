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

func (m msgServer) ExecuteBackingSwap(
	goCtx context.Context,
	msg *types.MsgExecuteBackingSwap,
) (*types.MsgExecuteBackingSwapResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil execute backing swap message")
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	amount, err := types.ParsePositiveAmount(msg.Amount)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	sequence, err := m.keeper.ExecuteBackingSwap(ctx, msg.Controller, amount, msg.SourcePort, msg.TimeoutTimestamp)
	if err != nil {
		return nil, err
	}
	return &types.MsgExecuteBackingSwapResponse{Sequence: sequence}, nil
}

func (m msgServer) CancelBackingSwap(
	goCtx context.Context,
	msg *types.MsgCancelBackingSwap,
) (*types.MsgCancelBackingSwapResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil cancel backing swap message")
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	amount, err := types.ParsePositiveAmount(msg.Amount)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.keeper.CancelBackingSwap(ctx, msg.Controller, amount); err != nil {
		return nil, err
	}
	return &types.MsgCancelBackingSwapResponse{}, nil
}

func (m msgServer) UpdateControls(
	goCtx context.Context,
	msg *types.MsgUpdateControls,
) (*types.MsgUpdateControlsResponse, error) {
	if msg == nil || !m.keeper.HasAuthority(msg.Authority) {
		return nil, types.ErrUnauthorized
	}
	if err := msg.Controls.Validate(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.keeper.ValidateRouteBindings(ctx, msg.Controls); err != nil {
		return nil, err
	}
	current := m.keeper.GetControls(ctx)
	if m.keeper.PendingCount(ctx) != 0 {
		if msg.Controls.Mode == types.Mode_MODE_DISABLED {
			return nil, fmt.Errorf("%w: cannot disable with pending settlements", types.ErrInvalidControls)
		}
		if routeIdentityChanged(current, msg.Controls) {
			return nil, fmt.Errorf("%w: route identity cannot change with pending settlements", types.ErrInvalidControls)
		}
	}
	cacheCtx, write := ctx.CacheContext()
	if err := m.keeper.SetControls(cacheCtx, msg.Controls); err != nil {
		return nil, err
	}
	if err := m.keeper.SetLedger(cacheCtx, m.keeper.GetLedger(cacheCtx)); err != nil {
		return nil, err
	}
	if reason, broken := BackingInvariant(*m.keeper)(cacheCtx); broken {
		return nil, fmt.Errorf("%w: %s", types.ErrInvalidLedger, reason)
	}
	write()
	return &types.MsgUpdateControlsResponse{}, nil
}

func routeIdentityChanged(a, b types.Controls) bool {
	return a.LogicalDenom != b.LogicalDenom ||
		a.NobleChannel != b.NobleChannel ||
		a.InjectiveChannel != b.InjectiveChannel ||
		a.NoblePacketDenom != b.NoblePacketDenom ||
		a.InjectivePacketDenom != b.InjectivePacketDenom ||
		a.InjectiveDenom != b.InjectiveDenom
}

func (m msgServer) UpdateParticipants(
	goCtx context.Context,
	msg *types.MsgUpdateParticipants,
) (*types.MsgUpdateParticipantsResponse, error) {
	if msg == nil || !m.keeper.HasAuthority(msg.Authority) {
		return nil, types.ErrUnauthorized
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	participants, err := m.mergeParticipantState(ctx, msg.Participants)
	if err != nil {
		return nil, err
	}
	if err := m.keeper.ReplaceParticipants(ctx, participants); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParticipantsResponse{}, nil
}

func (m msgServer) mergeParticipantState(ctx sdk.Context, updates []types.Participant) ([]types.Participant, error) {
	seen := make(map[string]struct{}, len(updates))
	merged := make([]types.Participant, len(updates))
	for i := range updates {
		update := updates[i]
		if (update.Funded != "" && update.Funded != "0") ||
			(update.Released != "" && update.Released != "0") {
			return nil, fmt.Errorf("%w: accounting fields are read-only", types.ErrInvalidParticipant)
		}
		if _, duplicate := seen[update.Controller]; duplicate {
			return nil, fmt.Errorf("%w: duplicate controller", types.ErrInvalidParticipant)
		}
		seen[update.Controller] = struct{}{}
		if existing, found := m.keeper.GetParticipant(ctx, update.Controller); found {
			update.Funded = existing.Funded
			update.Released = existing.Released
		} else {
			update.Funded = "0"
			update.Released = "0"
		}
		merged[i] = update
	}
	for _, existing := range m.keeper.GetAllParticipants(ctx) {
		if _, retained := seen[existing.Controller]; retained {
			continue
		}
		funded, _ := types.ParseAmount(existing.Funded)
		if !funded.IsZero() || hasPendingController(m.keeper.GetAllPending(ctx), existing.Controller) {
			return nil, fmt.Errorf("%w: participant has funded or pending state", types.ErrInvalidParticipant)
		}
	}
	return merged, nil
}

func hasPendingController(pending []types.PendingSettlement, controller string) bool {
	for _, settlement := range pending {
		if settlement.Controller == controller {
			return true
		}
	}
	return false
}
