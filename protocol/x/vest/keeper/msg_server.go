package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"

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
	if k.GetAuthority() != msg.Authority {
		return nil, sdkerrors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			k.GetAuthority(),
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
	if k.GetAuthority() != msg.Authority {
		return nil, sdkerrors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			k.GetAuthority(),
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
