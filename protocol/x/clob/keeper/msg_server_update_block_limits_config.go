package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UpdateBlockLimitsConfig updates the block limits config in state.
func (k msgServer) UpdateBlockLimitsConfig(
	goCtx context.Context,
	msg *types.MsgUpdateBlockLimitsConfig,
) (resp *types.MsgUpdateBlockLimitsConfigResponse, err error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	if err := k.Keeper.UpdateBlockLimitsConfig(ctx, msg.BlockLimitsConfig); err != nil {
		return nil, err
	}
	return &types.MsgUpdateBlockLimitsConfigResponse{}, nil
}
