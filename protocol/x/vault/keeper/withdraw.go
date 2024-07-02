package keeper

import (
	"math/big"

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

	ownerShares, sharesExists := k.GetOwnerShares(ctx, *msgWithdraw.GetVaultId(), msgWithdraw.SubaccountId.GetOwner())

	// 2. Subaccount has shares in the vault.
	if !sharesExists {
		return errors.Wrapf(
			types.ErrOwnerShareNotFound,
			"VaultId: %v, Owner: %v",
			*msgWithdraw.GetVaultId(),
			msgWithdraw.SubaccountId.GetOwner(),
		)
	}

	// 3. Shares to withdraw cannot be greater than the owner shares.
	if ownerShares.NumShares.BigInt().Cmp(msgWithdraw.Shares.NumShares.BigInt()) < 0 {
		return errors.Wrapf(
			types.ErrInvalidWithdrawalAmount,
			"VaultId: %v, Owner: %v, OwnerShares: %v, WithdrawalShares: %v",
			*msgWithdraw.GetVaultId(),
			msgWithdraw.SubaccountId.GetOwner(),
			ownerShares.NumShares.BigInt(),
			msgWithdraw.Shares.NumShares.BigInt(),
		)
	}

	return nil
}

func (k Keeper) RedeemShares(
	ctx sdk.Context,
	vaultId types.VaultId,
	owner string,
	sharesToRedeem types.NumShares,
) (quantumsToWithdraw *big.Int, err error) {
	// **** Validate ****
	// Validate that `sharesToRedeem` is positive.
	sharesToRedeemInt := sharesToRedeem.NumShares.BigInt()
	if sharesToRedeemInt.Sign() <= 0 {
		return nil, errors.Wrapf(
			types.ErrInvalidWithdrawalAmount,
			"Shares to withdraw must be positive",
		)
	}
	totalShares, vaultExists := k.GetTotalShares(ctx, vaultId)
	totalShareInt := totalShares.NumShares.BigInt()
	ownerShares, ownerSharesExist := k.GetOwnerShares(ctx, vaultId, owner)
	ownerSharesInt := ownerShares.NumShares.BigInt()
	// Validate that vault and owner shares exist.
	if !vaultExists || !ownerSharesExist {
		return nil, errors.Wrapf(
			types.ErrOwnerShareNotFound,
			"VaultId: %v, Owner: %v",
			vaultId,
			owner,
		)
	}
	// Validate that owner shares and total shares are positive.
	if ownerSharesInt.Sign() <= 0 || totalShareInt.Sign() <= 0 {
		return nil, errors.Wrapf(
			types.ErrNegativeShares,
			"Owner/Total shares must be positive. OwnerShares: %v, TotalShares: %v",
			ownerSharesInt,
			totalShareInt,
		)
	}
	// Validate that
	// 1. `sharesToRedeem` are less than or equal to the owner shares and
	// 2. owner shares are less than or equal to the total vault shares.
	if ownerSharesInt.Cmp(sharesToRedeemInt) < 0 ||
		totalShareInt.Cmp(ownerSharesInt) < 0 {
		return nil, errors.Wrapf(
			types.ErrInvalidWithdrawalAmount,
			"VaultId: %v, TotalShares: %v, OwnerShares: %v, WithdrawalShares: %v",
			vaultId,
			totalShareInt,
			ownerSharesInt,
			sharesToRedeem.NumShares.BigInt(),
		)
	}

	// **** Calculate Withdrawal Amount ****
	// Calculate the amount of quantums corresponding to the shares to withdraw
	equity, err := k.GetVaultEquity(ctx, vaultId)
	if err != nil {
		return nil, err
	}
	// Don't redeem shares if equity is non-positive.
	// This shouldn't happen, but check as a defense in depth.
	if equity.Sign() <= 0 {
		return nil, errors.Wrapf(types.ErrNonPositiveEquity, "VaultId: %v, Equity: %v", vaultId, equity)
	}
	quantumsToWithdraw = new(big.Int).Set(equity)
	quantumsToWithdraw = quantumsToWithdraw.Mul(quantumsToWithdraw, ownerShares.NumShares.BigInt())
	// Note: `Quo` discards the remainder, so this keeps the remainder staying in vault.
	quantumsToWithdraw = quantumsToWithdraw.Quo(quantumsToWithdraw, totalShares.NumShares.BigInt())

	// **** Update Shares ****
	// Decrease TotalShares of the vault by the amount of redeemed shares.
	err = k.SetTotalShares(
		ctx,
		vaultId,
		types.BigIntToNumShares(totalShareInt.Sub(totalShareInt, sharesToRedeemInt)),
	)
	if err != nil {
		return nil, err
	}
	// Decrease owner shares in the vault by the amount of redeemed shares.
	err = k.SetOwnerShares(
		ctx,
		vaultId,
		owner,
		types.BigIntToNumShares(ownerSharesInt.Sub(ownerSharesInt, sharesToRedeemInt)),
	)
	if err != nil {
		return nil, err
	}

	// **** Cleanup ****
	// Note: do not remove vault total shares that are zero, because the quoting strategy
	// may depend on the total shares of the vault.
	// TODO: clean up zero ownershares

	return quantumsToWithdraw, nil
}
