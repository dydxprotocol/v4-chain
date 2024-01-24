package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PremiumVotes(
	c context.Context,
	req *types.QueryPremiumVotesRequest,
) (*types.QueryPremiumVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	premiumVotes := k.GetPremiumVotes(ctx)

	return &types.QueryPremiumVotesResponse{PremiumVotes: premiumVotes}, nil
}

func (k Keeper) PremiumSamples(
	c context.Context,
	req *types.QueryPremiumSamplesRequest,
) (*types.QueryPremiumSamplesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	premiumSamples := k.GetPremiumSamples(ctx)

	return &types.QueryPremiumSamplesResponse{PremiumSamples: premiumSamples}, nil
}
