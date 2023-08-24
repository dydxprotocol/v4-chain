package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k msgServer) SetClobPairStatus(
	goCtx context.Context,
	msg *types.MsgSetClobPairStatus,
) (*types.MsgSetClobPairStatusResponse, error) {
	// TODO (CLOB-807): implement keeper SetClobPairStatus
	return &types.MsgSetClobPairStatusResponse{}, nil
}
