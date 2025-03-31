package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TwapOrder(
	c context.Context,
	req *types.QueryStatefulOrderRequest,
) (*types.QueryTWAPOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	// Get the TWAP order placement
	twapOrder, found := k.GetTwapOrderPlacement(
		ctx,
		req.OrderId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	// Get the TWAP trigger placements
	triggerPlacements, found_suborders := k.GetTwapTriggerPlacements(
		ctx,
		req.OrderId,
	)

	if !found || !found_suborders {
		return nil, status.Error(codes.NotFound, "not found")
	}

	res := &types.QueryTWAPOrderResponse{
		TwapOrderPlacement:    twapOrder,
		TwapTriggerPlacement:  triggerPlacements,
	}

	return res, nil
}
