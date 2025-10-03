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

// FeeHolidayParams processes a query for fee holiday parameters for a specific CLOB pair.
func (k Keeper) FeeHolidayParams(
	c context.Context,
	req *types.QueryFeeHolidayParamsRequest,
) (
	*types.QueryFeeHolidayParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	params, err := k.GetFeeHolidayParams(ctx, req.ClobPairId)
	if err != nil {
		if errors.Is(err, types.ErrFeeHolidayNotFound) {
			return nil, status.Error(codes.NotFound, "fee holiday not found for the specified CLOB pair")
		}
		return nil, status.Errorf(codes.Internal, "failed to get fee holiday: %v", err)
	}

	return &types.QueryFeeHolidayParamsResponse{
		Params: params,
	}, nil
}

// AllFeeHolidayParams processes a query for all fee holiday parameters.
func (k Keeper) AllFeeHolidayParams(
	c context.Context,
	req *types.QueryAllFeeHolidayParamsRequest,
) (
	*types.QueryAllFeeHolidayParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetAllFeeHolidayParams(ctx)

	return &types.QueryAllFeeHolidayParamsResponse{
		Params: params,
	}, nil
}
