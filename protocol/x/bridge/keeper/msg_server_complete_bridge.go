package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// CompleteBridge finalizes a bridge by minting coins to an address.
func (k msgServer) CompleteBridge(
	goCtx context.Context,
	msg *types.MsgCompleteBridge,
) (*types.MsgCompleteBridgeResponse, error) {
	// MsgCompleteBridge's authority should be bridge module.
	if k.Keeper.GetBridgeAuthority() != msg.Authority {
		return nil, errors.Wrapf(
			types.ErrInvalidAuthority,
			"expected %s, got %s",
			k.Keeper.GetBridgeAuthority(),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.CompleteBridge(ctx, msg.Event); err != nil {
		return nil, err
	}

	return &types.MsgCompleteBridgeResponse{}, nil
}
