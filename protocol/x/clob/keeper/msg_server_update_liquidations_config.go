package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateLiquidationsConfig updates the liquidation config in state.
func (k msgServer) UpdateLiquidationsConfig(
	goCtx context.Context,
	msg *types.MsgUpdateLiquidationsConfig,
) (*types.MsgUpdateLiquidationsConfigResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.UpdateLiquidationsConfig(ctx, msg.LiquidationsConfig); err != nil {
		return nil, err
	}
	return &types.MsgUpdateLiquidationsConfigResponse{}, nil
}
