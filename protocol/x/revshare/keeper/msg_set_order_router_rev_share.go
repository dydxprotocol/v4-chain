package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func (k msgServer) SetOrderRouterRevShare(
	goCtx context.Context,
	msg *types.MsgSetOrderRouterRevShare,
) (*types.MsgSetOrderRouterRevShareResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if sender is authorized to set revenue share
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	revShare := msg.OrderRouterRevShare
	if err := k.Keeper.SetOrderRouterRevShare(
		ctx, revShare.Address, revShare.SharePpm,
	); err != nil {
		return nil, err
	}

	return &types.MsgSetOrderRouterRevShareResponse{}, nil
}
