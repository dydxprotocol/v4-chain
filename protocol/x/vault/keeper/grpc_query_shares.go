package keeper

import (
	"context"
	"math/big"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
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

	ownerShares, exists := k.GetOwnerShares(ctx, req.Address)
	if !exists {
		return nil, status.Error(codes.NotFound, "owner not found")
	}
	if ownerShares.NumShares.Sign() < 0 {
		return nil, status.Error(codes.Internal, "owner has negative shares")
	} else if ownerShares.NumShares.Sign() == 0 {
		return &types.QueryMegavaultOwnerSharesResponse{
			Address: req.Address,
			Shares:  ownerShares,
		}, nil
	}

	shareUnlocks, exists := k.GetOwnerShareUnlocks(ctx, req.Address)
	totalLockedShares := big.NewInt(0)
	if exists {
		totalLockedShares = shareUnlocks.GetTotalLockedShares()
	}

	megavaultEquity, err := k.GetMegavaultEquity(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	totalShares := k.GetTotalShares(ctx)

	bigOwnerShares := ownerShares.NumShares.BigInt()
	ownerEquity := new(big.Int).Mul(bigOwnerShares, megavaultEquity)
	ownerEquity.Quo(ownerEquity, totalShares.NumShares.BigInt())

	ownerWithdrawableEquity := new(big.Int).Sub(bigOwnerShares, totalLockedShares)
	ownerWithdrawableEquity.Mul(ownerWithdrawableEquity, ownerEquity)
	ownerWithdrawableEquity.Quo(ownerWithdrawableEquity, bigOwnerShares)

	return &types.QueryMegavaultOwnerSharesResponse{
		Address:            req.Address,
		Shares:             ownerShares,
		ShareUnlocks:       shareUnlocks.ShareUnlocks,
		Equity:             dtypes.NewIntFromBigInt(ownerEquity),
		WithdrawableEquity: dtypes.NewIntFromBigInt(ownerWithdrawableEquity),
	}, nil
}

func (k Keeper) MegavaultAllOwnerShares(
	c context.Context,
	req *types.QueryMegavaultAllOwnerSharesRequest,
) (*types.QueryMegavaultAllOwnerSharesResponse, error) {
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

	return &types.QueryMegavaultAllOwnerSharesResponse{
		OwnerShares: ownerShares,
		Pagination:  pageRes,
	}, nil
}
