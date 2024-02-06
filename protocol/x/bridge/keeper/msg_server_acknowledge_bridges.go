package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// AcknowledgeBridges acknowledges bridge events and sets them to complete
// at a later block.
func (k msgServer) AcknowledgeBridges(
	goCtx context.Context,
	msg *types.MsgAcknowledgeBridges,
) (*types.MsgAcknowledgeBridgesResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if err := k.Keeper.AcknowledgeBridges(ctx, msg.Events); err != nil {
		return nil, err
	}

	return &types.MsgAcknowledgeBridgesResponse{}, nil
}
