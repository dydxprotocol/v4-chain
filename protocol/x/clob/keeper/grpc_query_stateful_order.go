package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StatefulOrder(
	c context.Context,
	req *types.QueryStatefulOrderRequest,
) (*types.QueryStatefulOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	val, found := k.GetLongTermOrderPlacement(
		ctx,
		req.OrderId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	res := &types.QueryStatefulOrderResponse{
		OrderPlacement: val,
	}

	// Get the fill amount
	_, fillAmount, _ := k.GetOrderFillAmount(ctx, req.OrderId)
	res.FillAmount = fillAmount.ToUint64()

	// Get triggered status for conditional orders
	if req.OrderId.IsConditionalOrder() {
		res.Triggered = k.IsConditionalOrderTriggered(ctx, req.OrderId)
	}

	return res, nil
}
