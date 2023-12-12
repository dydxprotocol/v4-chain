package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ListLimitParams(
	ctx context.Context,
	req *types.ListLimitParamsRequest,
) (*types.ListLimitParamsResponse, error) {
	// 	// TODO(CORE-823): implement query for `x/ratelimit`
	return nil, status.Errorf(codes.Unimplemented, "method ListLimitParams not implemented")
}
func (k Keeper) CapacityByDenom(
	ctx context.Context,
	req *types.QueryCapacityByDenomRequest,
) (*types.QueryCapacityByDenomResponse, error) {
	// 	// TODO(CORE-823): implement query for `x/ratelimit`
	return nil, status.Errorf(codes.Unimplemented, "method CapacityByDenom not implemented")
}
