package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

// CompleteBridge finalizes a bridge by transferring coins to an address.
func (k msgServer) CompleteBridge(
	goCtx context.Context,
	msg *types.MsgCompleteBridge,
) (*types.MsgCompleteBridgeResponse, error) {
	// MsgCompleteBridge's authority should be bridge module.
	bridgeModuleAddressString := authtypes.NewModuleAddress(types.ModuleName).String()
	if bridgeModuleAddressString != msg.Authority {
		return nil, errors.Wrapf(
			types.ErrInvalidAuthority,
			"expected %s, got %s",
			bridgeModuleAddressString,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.CompleteBridge(ctx, msg.Event); err != nil {
		return nil, err
	}

	return &types.MsgCompleteBridgeResponse{}, nil
}
