package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllStatefulOrders(
	c context.Context,
	req *types.QueryAllStatefulOrdersRequest,
) (
	*types.QueryAllStatefulOrdersResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	statefulOrders := k.GetAllStatefulOrders(ctx)

	return &types.QueryAllStatefulOrdersResponse{
		StatefulOrders: statefulOrders,
	}, nil
}

func (k Keeper) StatefulOrderCount(
	c context.Context,
	req *types.QueryStatefulOrderCountRequest,
) (
	*types.QueryStatefulOrderCountResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	count := k.GetStatefulOrderCount(ctx, *req.SubaccountId)

	return &types.QueryStatefulOrderCountResponse{
		Count: count,
	}, nil
}
