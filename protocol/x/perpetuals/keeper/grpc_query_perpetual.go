package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
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
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

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
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

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

func (k Keeper) NextPerpetualId(
	c context.Context,
	req *types.QueryNextPerpetualIdRequest,
) (
	*types.QueryNextPerpetualIdResponse,
	error,
) {
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	return &types.QueryNextPerpetualIdResponse{
		NextPerpetualId: k.GetNextPerpetualID(ctx),
	}, nil
}
