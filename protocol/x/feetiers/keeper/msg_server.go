package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"gopkg.in/typ.v4/maps"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
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

func (k msgServer) UpdatePerpetualFeeParams(
	goCtx context.Context,
	msg *types.MsgUpdatePerpetualFeeParams,
) (*types.MsgUpdatePerpetualFeeParamsResponse, error) {
	if _, ok := k.GetAuthorities()[msg.Authority]; !ok {
		return nil, sdkerrors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected one of %s, got %s",
			maps.Keys(k.GetAuthorities()),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetPerpetualFeeParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdatePerpetualFeeParamsResponse{}, nil
}
