package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) SetLiquidityTier(
	goCtx context.Context,
	msg *types.MsgSetLiquidityTier,
) (*types.MsgSetLiquidityTierResponse, error) {
	// TODO(CORE-563): implement this method
	return nil, status.Errorf(codes.Unimplemented, "SetLiquidityTier not implemented")
}
