package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AreSubaccountsLiquidatable(
	c context.Context,
	req *types.AreSubaccountsLiquidatableRequest,
) (
	*types.AreSubaccountsLiquidatableResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	results := make([]types.AreSubaccountsLiquidatableResponse_Result, len(req.SubaccountIds))
	for i, subaccountId := range req.SubaccountIds {
		isLiquidatable, err := k.IsLiquidatable(ctx, subaccountId)

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		results[i] = types.AreSubaccountsLiquidatableResponse_Result{
			SubaccountId:   subaccountId,
			IsLiquidatable: isLiquidatable,
		}
	}
	return &types.AreSubaccountsLiquidatableResponse{
		Results: results,
	}, nil
}
