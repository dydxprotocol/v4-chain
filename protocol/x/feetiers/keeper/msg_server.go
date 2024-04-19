package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
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

func (k msgServer) UpdatePerpetualFeeParams(
	goCtx context.Context,
	msg *types.MsgUpdatePerpetualFeeParams,
) (*types.MsgUpdatePerpetualFeeParamsResponse, error) {
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	if err := k.SetPerpetualFeeParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdatePerpetualFeeParamsResponse{}, nil
}
