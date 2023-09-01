package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateLiquidationsConfig updates the liquidation config in state.
func (k msgServer) UpdateLiquidationsConfig(
	goCtx context.Context,
	msg *types.MsgUpdateLiquidationsConfig,
) (*types.MsgUpdateLiquidationsConfigResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, sdkerrors.Wrapf(
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