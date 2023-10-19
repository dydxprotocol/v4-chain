package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllMarketPrices(
	c context.Context,
	req *types.QueryAllMarketPricesRequest,
) (*types.QueryAllMarketPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var marketPrices []types.MarketPrice
	ctx := sdk.UnwrapSDKContext(c)

	marketPriceStore := k.getMarketPriceStore(ctx)

	pageRes, err := query.Paginate(marketPriceStore, req.Pagination, func(key []byte, value []byte) error {
		var marketPrice types.MarketPrice
		if err := k.cdc.Unmarshal(value, &marketPrice); err != nil {
			return err
		}

		marketPrices = append(marketPrices, marketPrice)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllMarketPricesResponse{MarketPrices: marketPrices, Pagination: pageRes}, nil
}

func (k Keeper) MarketPrice(
	c context.Context,
	req *types.QueryMarketPriceRequest,
) (
	*types.QueryMarketPriceResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, err := k.GetMarketPrice(
		ctx,
		req.Id,
	)
	if err != nil {
		if errorsmod.IsOf(err, types.ErrMarketPriceDoesNotExist) {
			return nil, status.Error(codes.NotFound, "not found")
		} else {
			return nil, status.Error(codes.Internal, "unknown error getting market price")
		}
	}

	return &types.QueryMarketPriceResponse{MarketPrice: val}, nil
}

func (k Keeper) AllMarketParams(
	c context.Context,
	req *types.QueryAllMarketParamsRequest,
) (*types.QueryAllMarketParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var marketParams []types.MarketParam
	ctx := sdk.UnwrapSDKContext(c)

	marketParamStore := k.getMarketParamStore(ctx)

	pageRes, err := query.Paginate(marketParamStore, req.Pagination, func(key []byte, value []byte) error {
		var marketParam types.MarketParam
		if err := k.cdc.Unmarshal(value, &marketParam); err != nil {
			return err
		}

		marketParams = append(marketParams, marketParam)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllMarketParamsResponse{MarketParams: marketParams, Pagination: pageRes}, nil
}

func (k Keeper) MarketParam(
	c context.Context,
	req *types.QueryMarketParamRequest,
) (
	*types.QueryMarketParamResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, exists := k.GetMarketParam(
		ctx,
		req.Id,
	)
	if !exists {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryMarketParamResponse{MarketParam: val}, nil
}
