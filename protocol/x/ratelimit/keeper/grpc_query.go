package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ListLimitParams(
	ctx context.Context,
	req *types.ListLimitParamsRequest,
) (*types.ListLimitParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return &types.ListLimitParamsResponse{
		LimitParamsList: k.GetAllLimitParams(sdkCtx),
	}, nil
}

func (k Keeper) CapacityByDenom(
	ctx context.Context,
	req *types.QueryCapacityByDenomRequest,
) (*types.QueryCapacityByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	limiterCapacityList, err := k.GetLimiterCapacityListForDenom(sdkCtx, req.Denom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCapacityByDenomResponse{
		LimiterCapacityList: limiterCapacityList,
	}, nil
}
