package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// SetVaultParams sets the parameters of a specific vault.
func (k msgServer) SetVaultParams(
	goCtx context.Context,
	msg *types.MsgSetVaultParams,
) (*types.MsgSetVaultParamsResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	operator := k.GetOperatorParams(ctx).Operator

	// Check if authority is valid (must be a module authority or operator).
	if !k.HasAuthority(msg.Authority) && msg.Authority != operator {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidAuthority,
			"invalid authority %s",
			msg.Authority,
		)
	}

	// Validate parameters.
	if err := msg.VaultParams.Validate(); err != nil {
		return nil, err
	}

	// Set parameters for specified vault.
	if err := k.Keeper.SetVaultParams(ctx, msg.VaultId, msg.VaultParams); err != nil {
		return nil, err
	}

	return &types.MsgSetVaultParamsResponse{}, nil
}
