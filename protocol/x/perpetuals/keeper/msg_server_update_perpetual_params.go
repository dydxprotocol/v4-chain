package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) UpdatePerpetualParams(
	goCtx context.Context,
	msg *types.MsgUpdatePerpetualParams,
) (*types.MsgUpdatePerpetualParamsResponse, error) {
	// TODO(CORE-562): implement this method
	return nil, status.Errorf(codes.Unimplemented, "UpdatePerpetualParams not implemented")
}
