package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Params processes a query request/response for the Params from state.
func (k Keeper) PerpetualFeeParams(
	c context.Context,
	req *types.QueryPerpetualFeeParamsRequest,
) (
	*types.QueryPerpetualFeeParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetPerpetualFeeParams(ctx)
	return &types.QueryPerpetualFeeParamsResponse{
		Params: params,
	}, nil
}

func (k Keeper) UserFeeTier(
	c context.Context,
	req *types.QueryUserFeeTierRequest,
) (
	*types.QueryUserFeeTierResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	index, tier := k.getUserFeeTier(ctx, req.User)
	return &types.QueryUserFeeTierResponse{
		Index: index,
		Tier:  tier,
	}, nil
}
