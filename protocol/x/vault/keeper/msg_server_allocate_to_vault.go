package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// AllocateToVault allocates funds from main vault to a vault.
func (k msgServer) AllocateToVault(
	goCtx context.Context,
	msg *types.MsgAllocateToVault,
) (*types.MsgAllocateToVaultResponse, error) {
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

	// Transfer from main vault to the specified vault.
	if err := k.Keeper.subaccountsKeeper.TransferFundsFromSubaccountToSubaccount(
		ctx,
		types.MegavaultMainSubaccount,
		*msg.VaultId.ToSubaccountId(),
		assetstypes.AssetUsdc.Id,
		msg.QuoteQuantums.BigInt(),
	); err != nil {
		return nil, err
	}

	return &types.MsgAllocateToVaultResponse{}, nil
}
