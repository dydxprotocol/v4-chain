package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// CompleteBridge finalizes a bridge by transferring coins to an address.
func (k msgServer) CompleteBridge(
	goCtx context.Context,
	msg *types.MsgCompleteBridge,
) (*types.MsgCompleteBridgeResponse, error) {
	if !k.Keeper.HasAuthority(msg.GetAuthority()) {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidAuthority,
			"message authority %s is not valid for sending complete bridge messages",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if err := k.Keeper.CompleteBridge(ctx, msg.Event); err != nil {
		return nil, err
	}

	return &types.MsgCompleteBridgeResponse{}, nil
}
