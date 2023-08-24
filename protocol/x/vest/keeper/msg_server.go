package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"gopkg.in/typ.v4/maps"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetVestEntry(
	goCtx context.Context,
	msg *types.MsgSetVestEntry,
) (*types.MsgSetVestEntryResponse, error) {
	if _, ok := k.GetAuthorities()[msg.Authority]; !ok {
		return nil, sdkerrors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected one of %s, got %s",
			maps.Keys(k.GetAuthorities()),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.Keeper.SetVestEntry(ctx, msg.Entry); err != nil {
		return nil, err
	}

	return &types.MsgSetVestEntryResponse{}, nil
}

func (k msgServer) DeleteVestEntry(
	goCtx context.Context,
	msg *types.MsgDeleteVestEntry,
) (*types.MsgDeleteVestEntryResponse, error) {
	if _, ok := k.GetAuthorities()[msg.Authority]; !ok {
		return nil, sdkerrors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected one of %s, got %s",
			maps.Keys(k.GetAuthorities()),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.Keeper.DeleteVestEntry(ctx, msg.VesterAccount)
	if err != nil {
		return nil, err
	}

	return &types.MsgDeleteVestEntryResponse{}, nil
}
