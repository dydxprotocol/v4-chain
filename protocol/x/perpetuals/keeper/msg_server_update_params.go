package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	// TODO(CORE-560): implement this method
	return nil, status.Errorf(codes.Unimplemented, "UpdateParams not implemented")
}
