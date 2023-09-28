package keeper

import (
	"context"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllPerpetuals(
	c context.Context,
	req *types.QueryAllPerpetualsRequest,
) (*types.QueryAllPerpetualsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var perpetuals []types.Perpetual
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	perpetualStore := prefix.NewStore(store, []byte(types.PerpetualKeyPrefix))

	pageRes, err := query.Paginate(perpetualStore, req.Pagination, func(key []byte, value []byte) error {
		var perpetual types.Perpetual
		if err := k.cdc.Unmarshal(value, &perpetual); err != nil {
			return err
		}

		perpetuals = append(perpetuals, perpetual)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPerpetualsResponse{Perpetual: perpetuals, Pagination: pageRes}, nil
}

func (k Keeper) Perpetual(c context.Context, req *types.QueryPerpetualRequest) (*types.QueryPerpetualResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, err := k.GetPerpetual(
		ctx,
		req.Id,
	)
	if err != nil {
		if errors.Is(err, types.ErrPerpetualDoesNotExist) {
			return nil,
				status.Error(
					codes.NotFound,
					fmt.Sprintf(
						"Perpetual id %+v not found.",
						req.Id,
					),
				)
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryPerpetualResponse{Perpetual: val}, nil
}
