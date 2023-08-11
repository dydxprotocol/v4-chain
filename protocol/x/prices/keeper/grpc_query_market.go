package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4/x/prices/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllMarkets(
	c context.Context,
	req *types.QueryAllMarketsRequest,
) (*types.QueryAllMarketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var markets []types.Market
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	marketStore := prefix.NewStore(store, types.KeyPrefix(types.MarketKeyPrefix))

	pageRes, err := query.Paginate(marketStore, req.Pagination, func(key []byte, value []byte) error {
		var market types.Market
		if err := k.cdc.Unmarshal(value, &market); err != nil {
			return err
		}

		markets = append(markets, market)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllMarketsResponse{Market: markets, Pagination: pageRes}, nil
}

func (k Keeper) Market(c context.Context, req *types.QueryMarketRequest) (*types.QueryMarketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, err := k.GetMarket(
		ctx,
		req.Id,
	)
	if err != nil {
		if sdkerrors.IsOf(err, types.ErrMarketDoesNotExist) {
			return nil, status.Error(codes.NotFound, "not found")
		} else {
			return nil, status.Error(codes.Internal, "unknown error getting market")
		}
	}

	return &types.QueryMarketResponse{Market: val}, nil
}
