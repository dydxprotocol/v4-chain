package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	premiumVotes := k.GetPremiumVotes(sdk.UnwrapSDKContext(c))

	return &types.QueryPremiumVotesResponse{PremiumVotes: premiumVotes}, nil
}

func (k Keeper) PremiumSamples(
	c context.Context,
	req *types.QueryPremiumSamplesRequest,
) (*types.QueryPremiumSamplesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	premiumSamples := k.GetPremiumSamples(sdk.UnwrapSDKContext(c))

	return &types.QueryPremiumSamplesResponse{PremiumSamples: premiumSamples}, nil
}
