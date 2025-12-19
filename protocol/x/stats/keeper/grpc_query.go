package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// Params processes a query request/response for the Params from state.
func (k Keeper) Params(
	c context.Context,
	req *types.QueryParamsRequest,
) (
	*types.QueryParamsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (k Keeper) StatsMetadata(
	c context.Context,
	req *types.QueryStatsMetadataRequest,
) (
	*types.QueryStatsMetadataResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	statsMetadata := k.GetStatsMetadata(ctx)
	return &types.QueryStatsMetadataResponse{
		Metadata: statsMetadata,
	}, nil
}

func (k Keeper) GlobalStats(
	c context.Context,
	req *types.QueryGlobalStatsRequest,
) (
	*types.QueryGlobalStatsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	globalStats := k.GetGlobalStats(ctx)
	return &types.QueryGlobalStatsResponse{
		Stats: globalStats,
	}, nil
}

func (k Keeper) UserStats(
	c context.Context,
	req *types.QueryUserStatsRequest,
) (
	*types.QueryUserStatsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	userStats := k.GetUserStats(ctx, req.User)
	return &types.QueryUserStatsResponse{
		Stats: userStats,
	}, nil
}

func (k Keeper) EpochStats(
	c context.Context,
	req *types.QueryEpochStatsRequest,
) (
	*types.QueryEpochStatsResponse,
	error,
) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	epochStats := k.GetEpochStatsOrNil(ctx, req.Epoch)

	// Return empty stats if epoch not found
	if epochStats == nil {
		epochStats = &types.EpochStats{
			Stats: []*types.EpochStats_UserWithStats{},
		}
	}

	return &types.QueryEpochStatsResponse{
		Stats: epochStats,
	}, nil
}
