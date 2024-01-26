package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
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

func (k msgServer) SetLimitParams(
	ctx context.Context,
	msg *types.MsgSetLimitParams,
) (*types.MsgSetLimitParamsResponse, error) {
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	// msg.LimitParams.Validate() is called in `Keeper.SetLimitParams`
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.Keeper.SetLimitParams(sdkCtx, msg.LimitParams); err != nil {
		return nil, err
	}

	return &types.MsgSetLimitParamsResponse{}, nil
}
