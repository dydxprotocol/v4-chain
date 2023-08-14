package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// CompleteBridge finalizes a bridge by minting coins to an address.
func (k msgServer) CompleteBridge(
	goCtx context.Context,
	msg *types.MsgCompleteBridge,
) (*types.MsgCompleteBridgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.CompleteBridge(ctx, msg.Event); err != nil {
		return nil, err
	}

	return &types.MsgCompleteBridgeResponse{}, nil
}
