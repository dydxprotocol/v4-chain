package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// ValidateWithdrawFromVault performs stateful validations on `MsgWithdrawFromVault`.
func (k Keeper) ValidateWithdrawFromVault(
	ctx sdk.Context,
	msgWithdraw *types.MsgWithdrawFromVault,
) error {
	// 1. Vault exists.
	if _, vaultExists := k.GetTotalShares(ctx, *msgWithdraw.GetVaultId()); !vaultExists {
		return errors.Wrapf(types.ErrVaultNotFound, "VaultId: %v", *msgWithdraw.GetVaultId())
	}

	// 2. Subaccount has shares in the vault.
	_, sharesExists := k.GetOwnerShares(ctx, *msgWithdraw.GetVaultId(), msgWithdraw.SubaccountId.GetOwner())
	if !sharesExists {
		return errors.Wrapf(
			types.ErrOwnerShareNotFound,
			"VaultId: %v, Owner: %v",
			*msgWithdraw.GetVaultId(),
			msgWithdraw.SubaccountId.GetOwner(),
		)
	}
	return nil
}
