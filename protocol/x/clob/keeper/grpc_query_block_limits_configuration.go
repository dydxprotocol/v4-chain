package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BlockLimitsConfiguration returns the block limits configuration.
func (k Keeper) BlockLimitsConfiguration(
	c context.Context,
	req *types.QueryBlockLimitsConfigurationRequest,
) (*types.QueryBlockLimitsConfigurationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := lib.UnwrapSDKContext(c, types.ModuleName)
	blockLimitsConfig := k.GetBlockLimitsConfig(ctx)
	return &types.QueryBlockLimitsConfigurationResponse{
		BlockLimitsConfig: blockLimitsConfig,
	}, nil
}
