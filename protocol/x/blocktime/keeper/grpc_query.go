package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) SynchronyParams(
	c context.Context,
	req *types.QuerySynchronyParamsRequest,
) (
	*types.QuerySynchronyParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetSynchronyParams(ctx)
	return &types.QuerySynchronyParamsResponse{
		Params: params,
	}, nil
}

func (k Keeper) DowntimeParams(
	c context.Context,
	req *types.QueryDowntimeParamsRequest,
) (
	*types.QueryDowntimeParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetDowntimeParams(ctx)
	return &types.QueryDowntimeParamsResponse{
		Params: params,
	}, nil
}

func (k Keeper) PreviousBlockInfo(
	c context.Context,
	req *types.QueryPreviousBlockInfoRequest,
) (
	*types.QueryPreviousBlockInfoResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	info := k.GetPreviousBlockInfo(ctx)
	return &types.QueryPreviousBlockInfoResponse{
		Info: &info,
	}, nil
}

func (k Keeper) AllDowntimeInfo(
	c context.Context,
	req *types.QueryAllDowntimeInfoRequest,
) (
	*types.QueryAllDowntimeInfoResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	info := k.GetAllDowntimeInfo(ctx)
	return &types.QueryAllDowntimeInfoResponse{
		Info: info,
	}, nil
}
