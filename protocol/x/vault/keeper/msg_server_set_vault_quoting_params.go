package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// SetVaultQuotingParams sets the quoting parameters of a specific vault.
func (k msgServer) SetVaultQuotingParams(
	goCtx context.Context,
	msg *types.MsgSetVaultQuotingParams,
) (*types.MsgSetVaultQuotingParamsResponse, error) {
	// Check if authority is valid.
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	// Validate quoting parameters.
	if err := msg.QuotingParams.Validate(); err != nil {
		return nil, err
	}

	// Set quoting parameters for specified vault.
	if err := k.Keeper.SetVaultQuotingParams(ctx, msg.VaultId, msg.QuotingParams); err != nil {
		return nil, err
	}

	return &types.MsgSetVaultQuotingParamsResponse{}, nil
}
