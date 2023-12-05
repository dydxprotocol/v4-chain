package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllLiquidityTiers(
	c context.Context,
	req *types.QueryAllLiquidityTiersRequest,
) (*types.QueryAllLiquidityTiersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	liquidityTiers := k.GetAllLiquidityTiers(ctx)

	return &types.QueryAllLiquidityTiersResponse{LiquidityTiers: liquidityTiers}, nil
}
