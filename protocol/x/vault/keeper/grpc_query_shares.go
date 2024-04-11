package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func (k Keeper) OwnerShares(
	c context.Context,
	req *types.QueryOwnerSharesRequest,
) (*types.QueryOwnerSharesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	vaultId := types.VaultId{
		Type:   req.Type,
		Number: req.Number,
	}
	_, exists := k.GetTotalShares(ctx, vaultId)
	if !exists {
		return nil, status.Error(codes.NotFound, "vault not found")
	}

	var ownerShares []*types.OwnerShare
	ownerSharesStore := k.getVaultOwnerSharesStore(ctx, vaultId)
	pageRes, err := query.Paginate(ownerSharesStore, req.Pagination, func(key []byte, value []byte) error {
		owner := string(key)

		var shares types.NumShares
		k.cdc.MustUnmarshal(value, &shares)

		ownerShares = append(ownerShares, &types.OwnerShare{
			Owner:  owner,
			Shares: &shares,
		})

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOwnerSharesResponse{
		OwnerShares: ownerShares,
		Pagination:  pageRes,
	}, nil
}
