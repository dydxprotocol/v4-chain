package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/typ.v4/maps"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
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

func (k msgServer) UpdateDowntimeParams(
	goCtx context.Context,
	msg *types.MsgUpdateDowntimeParams,
) (*types.MsgUpdateDowntimeParamsResponse, error) {
	if _, ok := k.GetAuthorities()[msg.Authority]; !ok {
		return nil, sdkerrors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected one of %s, got %s",
			maps.Keys(k.GetAuthorities()),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetDowntimeParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateDowntimeParamsResponse{}, nil
}

func (k msgServer) IsDelayedBlock(
	goCtx context.Context,
	msg *types.MsgIsDelayedBlock,
) (*types.MsgIsDelayedBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "IsDelayedBlock not implemented")
}
