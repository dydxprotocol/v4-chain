package keeper

import (
	"math/big"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

type SharesAfterWithdrawal struct {
	RedeemedShares  types.NumShares
	RemainingShares types.NumShares
	TotalShares     types.NumShares
}

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

func (k Keeper) RedeemShares(
	ctx sdk.Context,
	vaultId types.VaultId,
	owner string,
	quantumToWithdraw *big.Int,
) (SharesAfterWithdrawal, error) {

	// WATCH OUT FOR ANY ROUNDING. Should usually round up so that the vault gets the benefit.

	// There are two cases:
	// a) withdrawal amount + slippage <= owner vault equity
	//      = shares: redeem shares equal to the above amount
	//		= transfer: transfer withdrawal amount
	// b) withdrawal amount + slippage > owner vault equity
	//      = shares: redeem all shares
	//		= transfer: transfer withdrawal amount minus slippage that's over the owner vault equity
	// Note: if the withdrawal amount > owner vault equity,
	// then update the withdrawal amount to owner vault equity.
	// In that case, then it's the same as case b.

	// 0. get the owner's equity in the vault (owner vault equity)
	// 1. cap the withdrawal amount to the owner vault equity
	//      adjustedAmount = min(toWithdraw, owner vault equity)
	// 2. calculate the slippage to withdraw the amount in 1
	//      slippage = some formula
	// 3. calculate the shares to redeem
	//      sharesToRedeem = min(owner shares, shares(adjustedAmount + slippage))
	// 4. calculate the effective amount to withdraw
	//      adjWithSlippage = adjustedAmount + slippage
	//      if adjWithSlippage <= owner vault equity:
	//				then effectiveAmount = toWithdraw
	//      if adjWithSlippage > owner vault equity:
	//				then effectiveAmount = toWithdraw - (adjWithSlippage - owner vault equity)

	// The effective quantums to withdraw should be positive.
	if quantumToWithdraw.Sign() <= 0 {
		return SharesAfterWithdrawal{}, types.ErrInvalidWithdrawalAmount
	}

	// Get vault equity.
	vaultEquity, err := k.GetVaultEquity(ctx, vaultId)
	if err != nil {
		return SharesAfterWithdrawal{}, err
	}
	if vaultEquity.Sign() <= 0 {
		return SharesAfterWithdrawal{}, types.ErrNonPositiveEquity
	}
	// Check if the vault has enough equity.
	if vaultEquity.Cmp(quantumToWithdraw) < 0 {
		return SharesAfterWithdrawal{}, types.ErrInsufficientEquity
	}

	// Calculate shares to redeem.
	// Should be capped by the owner's shares in the vault.

	// Update the total and owner shares (remove entries if zero).

	return SharesAfterWithdrawal{}, nil
}
