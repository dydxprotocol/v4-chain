package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// SetVaultParams sets the parameters of a specific vault.
func (k msgServer) SetVaultParams(
	goCtx context.Context,
	msg *types.MsgSetVaultParams,
) (*types.MsgSetVaultParamsResponse, error) {
	// Check if authority is valid.
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

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
