package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatefulOrderCount returns the count of stateful orders for a subaccount.
func (k Keeper) StatefulOrderCount(
	c context.Context,
	req *types.QueryStatefulOrderCountRequest,
) (*types.QueryStatefulOrderCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	count := k.GetStatefulOrderCount(ctx, satypes.SubaccountId{
		Owner:  req.Owner,
		Number: req.SubaccountNumber,
	})
	return &types.QueryStatefulOrderCountResponse{
		Count: count,
	}, nil
}
