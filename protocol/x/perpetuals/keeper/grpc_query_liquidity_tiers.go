package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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
	var liquidityTiers []types.LiquidityTier

	store := ctx.KVStore(k.storeKey)
	liquidityTierStore := prefix.NewStore(store, []byte(types.LiquidityTierKeyPrefix))

	pageRes, err := query.Paginate(liquidityTierStore, req.Pagination, func(key []byte, value []byte) error {
		var liquidityTier types.LiquidityTier
		if err := k.cdc.Unmarshal(value, &liquidityTier); err != nil {
			return err
		}

		liquidityTiers = append(liquidityTiers, liquidityTier)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllLiquidityTiersResponse{
		LiquidityTiers: liquidityTiers,
		Pagination:     pageRes,
	}, nil
}
