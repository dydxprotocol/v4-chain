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

// FeeDiscountCampaignParams processes a query for fee discount campaign parameters for a specific CLOB pair.
func (k Keeper) FeeDiscountCampaignParams(
	c context.Context,
	req *types.QueryFeeDiscountCampaignParamsRequest,
) (
	*types.QueryFeeDiscountCampaignParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	params, err := k.GetFeeDiscountCampaignParams(ctx, req.ClobPairId)
	if err != nil {
		if errors.Is(err, types.ErrFeeDiscountCampaignNotFound) {
			return nil, status.Error(codes.NotFound, "fee discount campaign not found for the specified CLOB pair")
		}
		return nil, status.Errorf(codes.Internal, "failed to get fee discount campaign: %v", err)
	}

	return &types.QueryFeeDiscountCampaignParamsResponse{
		Params: params,
	}, nil
}

// AllFeeDiscountCampaignParams processes a query for all fee discount campaign parameters.
func (k Keeper) AllFeeDiscountCampaignParams(
	c context.Context,
	req *types.QueryAllFeeDiscountCampaignParamsRequest,
) (
	*types.QueryAllFeeDiscountCampaignParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetAllFeeDiscountCampaignParams(ctx)

	return &types.QueryAllFeeDiscountCampaignParamsResponse{
		Params: params,
	}, nil
}
