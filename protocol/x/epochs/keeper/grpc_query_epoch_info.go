package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) EpochInfoAll(
	c context.Context, req *types.QueryAllEpochInfoRequest) (*types.QueryEpochInfoAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var epochInfos []types.EpochInfo
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	epochInfoStore := k.getEpochInfoStore(ctx)

	pageRes, err := query.Paginate(epochInfoStore, req.Pagination, func(key []byte, value []byte) error {
		var epochInfo types.EpochInfo
		if err := k.cdc.Unmarshal(value, &epochInfo); err != nil {
			return err
		}

		epochInfos = append(epochInfos, epochInfo)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryEpochInfoAllResponse{EpochInfo: epochInfos, Pagination: pageRes}, nil
}

func (k Keeper) EpochInfo(
	c context.Context, req *types.QueryGetEpochInfoRequest) (*types.QueryEpochInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	val, found := k.GetEpochInfo(
		ctx,
		types.EpochInfoName(req.Name),
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryEpochInfoResponse{EpochInfo: val}, nil
}
