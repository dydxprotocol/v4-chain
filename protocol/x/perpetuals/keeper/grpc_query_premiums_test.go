package keeper

import (
	"context"
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

}
