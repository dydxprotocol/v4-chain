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
	goCtx context.Context,
	req *types.ListLimitParamsRequest,
) (*types.ListLimitParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.ListLimitParamsResponse{
		LimitParamsList: k.GetAllLimitParams(ctx),
	}, nil
}

func (k Keeper) CapacityByDenom(
	goCtx context.Context,
	req *types.QueryCapacityByDenomRequest,
) (*types.QueryCapacityByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := sdk.ValidateDenom(req.Denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	limiterCapacityList, err := k.GetLimiterCapacityListForDenom(ctx, req.Denom)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryCapacityByDenomResponse{
		LimiterCapacityList: limiterCapacityList,
	}, nil
}
