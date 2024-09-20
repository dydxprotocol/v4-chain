package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// RetrieveFromVault retrieves funds from a vault to main vault.
func (k msgServer) RetrieveFromVault(
	goCtx context.Context,
	msg *types.MsgRetrieveFromVault,
) (*types.MsgRetrieveFromVaultResponse, error) {
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

	// Check if vault exists.
	if _, exists := k.Keeper.GetVaultParams(ctx, msg.VaultId); !exists {
		return nil, types.ErrVaultParamsNotFound
	}

	// Transfer from specified vault to main vault.
	if err := k.Keeper.sendingKeeper.ProcessTransfer(
		ctx,
		&sendingtypes.Transfer{
			Sender:    *msg.VaultId.ToSubaccountId(),
			Recipient: types.MegavaultMainSubaccount,
			AssetId:   assetstypes.AssetUsdc.Id,
			Amount:    msg.QuoteQuantums.BigInt().Uint64(), // validated to be positive above.
		},
	); err != nil {
		return nil, err
	}

	return &types.MsgRetrieveFromVaultResponse{}, nil
}
