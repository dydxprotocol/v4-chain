package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func (k msgServer) SetMarketsHardCap(
	goCtx context.Context,
	msg *types.MsgSetMarketsHardCap,
) (*types.MsgSetMarketsHardCapResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	// Check if the sender has the authority to set the hard cap
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	// Set the hard cap for listed markets
	err := k.Keeper.SetMarketsHardCap(ctx, msg.HardCapForMarkets)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetMarketsHardCapResponse{}, nil
}
