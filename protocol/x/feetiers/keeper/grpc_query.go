package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Params processes a query request/response for the Params from state.
func (k Keeper) PerpetualFeeParams(
	c context.Context,
	req *types.QueryPerpetualFeeParamsRequest,
) (
	*types.QueryPerpetualFeeParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetPerpetualFeeParams(ctx)
	return &types.QueryPerpetualFeeParamsResponse{
		Params: params,
	}, nil
}

func (k Keeper) UserFeeTier(
	c context.Context,
	req *types.QueryUserFeeTierRequest,
) (
	*types.QueryUserFeeTierResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	if _, err := sdk.AccAddressFromBech32(req.User); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bech32 address")
	}
	index, tier := k.getUserFeeTier(ctx, req.User)
	return &types.QueryUserFeeTierResponse{
		Index: index,
		Tier:  tier,
	}, nil
}

// PerMarketFeeDiscountParams processes a query for fee discount parameters for a specific market/CLOB pair.
func (k Keeper) PerMarketFeeDiscountParams(
	c context.Context,
	req *types.QueryPerMarketFeeDiscountParamsRequest,
) (
	*types.QueryPerMarketFeeDiscountParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params, err := k.GetPerMarketFeeDiscountParams(ctx, req.ClobPairId)
	if err != nil {
		if errors.Is(err, types.ErrMarketFeeDiscountNotFound) {
			return nil, status.Error(codes.NotFound, "fee discount not found for the specified market/CLOB pair")
		}
		return nil, status.Errorf(codes.Internal, "failed to get per-market fee discount: %v", err)
	}

	return &types.QueryPerMarketFeeDiscountParamsResponse{
		Params: params,
	}, nil
}

// AllMarketFeeDiscountParams processes a query for all market fee discount parameters.
func (k Keeper) AllMarketFeeDiscountParams(
	c context.Context,
	req *types.QueryAllMarketFeeDiscountParamsRequest,
) (
	*types.QueryAllMarketFeeDiscountParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetAllMarketFeeDiscountParams(ctx)

	return &types.QueryAllMarketFeeDiscountParamsResponse{
		Params: params,
	}, nil
}
