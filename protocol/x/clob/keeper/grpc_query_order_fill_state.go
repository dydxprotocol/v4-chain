package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) OrderFillStates(
	c context.Context,
	req *types.QueryOrderFillStatesRequest,
) (*types.QueryOrderFillStatesResponse, error) {
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	var fstates []*types.OrderIdAndFillState

	orderIdsAndFillStates := k.GetAllOrderFillStates(ctx)
	for _, idFillState := range orderIdsAndFillStates {
		fstates = append(fstates, &types.OrderIdAndFillState{
			OrderId:   &idFillState.OrderId,
			FillState: &idFillState.OrderFillState,
		})
	}

	return &types.QueryOrderFillStatesResponse{OrderIdAndFillState: fstates}, nil
}
