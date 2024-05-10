package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) PrunableOrders(
	c context.Context,
	req *types.QueryPrunableOrdersRequest,
) (*types.QueryPrunableOrdersResponse, error) {
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	var orderIds []*types.OrderId

	potentiallyPrunableOrdersStore := k.GetPruneableOrdersStore(ctx, req.BlockHeight)
	it := potentiallyPrunableOrdersStore.Iterator(nil, nil)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var orderId types.OrderId
		k.cdc.MustUnmarshal(it.Value(), &orderId)
		orderIds = append(orderIds, &orderId)
	}

	return &types.QueryPrunableOrdersResponse{PrunableOrders: orderIds}, nil
}
