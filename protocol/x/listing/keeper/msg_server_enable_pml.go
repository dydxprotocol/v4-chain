package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func (k msgServer) EnablePermissionlessMarketListing(
	goCtx context.Context,
	msg *types.MsgEnablePermissionlessMarketListing,
) (*types.MsgEnablePermissionlessMarketListingResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	// Check if the sender has the authority to enable permissionless market listing.
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	// Set the permissionless listing enable flag
	err := k.Keeper.SetPermissionlessListingEnable(ctx, msg.EnablePermissionlessMarketListing)
	if err != nil {
		return nil, err
	}

	return &types.MsgEnablePermissionlessMarketListingResponse{}, nil
}
