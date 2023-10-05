package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	errorlib "github.com/dydxprotocol/v4-chain/protocol/lib/error"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateLiquidationsConfig updates the liquidation config in state.
func (k msgServer) UpdateLiquidationsConfig(
	goCtx context.Context,
	msg *types.MsgUpdateLiquidationsConfig,
) (resp *types.MsgUpdateLiquidationsConfigResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	defer func() {
		if err != nil {
			errorlib.LogErrorWithBlockHeight(k.Keeper.Logger(ctx), err, ctx.BlockHeight(), metrics.DeliverTx)
		}
	}()

	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	if err := k.Keeper.UpdateLiquidationsConfig(ctx, msg.LiquidationsConfig); err != nil {
		return nil, err
	}
	return &types.MsgUpdateLiquidationsConfigResponse{}, nil
}
