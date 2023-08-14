package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

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

	ctx := sdk.UnwrapSDKContext(c)
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

	ctx := sdk.UnwrapSDKContext(c)
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

	ctx := sdk.UnwrapSDKContext(c)
	info := k.GetAllDowntimeInfo(ctx)
	return &types.QueryAllDowntimeInfoResponse{
		Info: info,
	}, nil
}
