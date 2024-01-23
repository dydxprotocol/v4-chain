package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"

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
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	if err := k.Keeper.SetVestEntry(ctx, msg.Entry); err != nil {
		return nil, err
	}

	return &types.MsgSetVestEntryResponse{}, nil
}

func (k msgServer) DeleteVestEntry(
	goCtx context.Context,
	msg *types.MsgDeleteVestEntry,
) (*types.MsgDeleteVestEntryResponse, error) {
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	err := k.Keeper.DeleteVestEntry(ctx, msg.VesterAccount)
	if err != nil {
		return nil, err
	}

	return &types.MsgDeleteVestEntryResponse{}, nil
}
