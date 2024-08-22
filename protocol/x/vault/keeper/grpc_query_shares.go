package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func (k Keeper) MegavaultTotalShares(
	c context.Context,
	_ *types.QueryMegavaultTotalSharesRequest,
) (*types.QueryMegavaultTotalSharesResponse, error) {
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	totalShares := k.GetTotalShares(ctx)

	return &types.QueryMegavaultTotalSharesResponse{
		TotalShares: &totalShares,
	}, nil
}

func (k Keeper) MegavaultOwnerShares(
	c context.Context,
	req *types.QueryMegavaultOwnerSharesRequest,
) (*types.QueryMegavaultOwnerSharesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	var ownerShares []*types.OwnerShare
	ownerSharesStore := k.getOwnerSharesStore(ctx)
	pageRes, err := query.Paginate(ownerSharesStore, req.Pagination, func(key []byte, value []byte) error {
		owner := string(key)

		var shares types.NumShares
		k.cdc.MustUnmarshal(value, &shares)

		ownerShares = append(ownerShares, &types.OwnerShare{
			Owner:  owner,
			Shares: shares,
		})

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMegavaultOwnerSharesResponse{
		OwnerShares: ownerShares,
		Pagination:  pageRes,
	}, nil
}
