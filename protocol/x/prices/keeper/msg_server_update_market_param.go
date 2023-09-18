package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) UpdateMarketParam(
	goCtx context.Context,
	msg *types.MsgUpdateMarketParam,
) (*types.MsgUpdateMarketParamResponse, error) {
	// TODO(CORE-564): implement this method
	return nil, status.Errorf(codes.Unimplemented, "UpdateMarketParam not implemented")
}
