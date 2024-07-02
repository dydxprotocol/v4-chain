package keeper

import (
	"math/big"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// MintShares mints shares of a vault for `owner` based on `quantumsToDeposit` by:
// 1. Increasing total shares of the vault.
// 2. Increasing owner shares of the vault for given `owner`.
func (k Keeper) MintShares(
	ctx sdk.Context,
	vaultId types.VaultId,
	owner string,
	quantumsToDeposit *big.Int,
) error {
	// Quantums to deposit should be positive.
	if quantumsToDeposit.Sign() <= 0 {
		return types.ErrInvalidDepositAmount
	}
	// Get existing TotalShares of the vault.
	totalShares, exists := k.GetTotalShares(ctx, vaultId)
	existingTotalShares := totalShares.NumShares.BigInt()
	// Calculate shares to mint.
	var sharesToMint *big.Int
	if !exists || existingTotalShares.Sign() <= 0 {
		// Mint `quoteQuantums` number of shares.
		sharesToMint = new(big.Int).Set(quantumsToDeposit)
		// Initialize existingTotalShares as 0.
		existingTotalShares = big.NewInt(0)
	} else {
		// Get vault equity.
		equity, err := k.GetVaultEquity(ctx, vaultId)
		if err != nil {
			return err
		}
		// Don't mint shares if equity is non-positive.
		// This shouldn't happen, but check as a defense in depth.
		if equity.Sign() <= 0 {
			return errors.Wrapf(types.ErrNonPositiveEquity, "VaultId: %v, Equity: %v", vaultId, equity)
		}
		// Mint `deposit (in quote quantums) * existing shares / vault equity (in quote quantums)`
		// number of shares.
		// For example:
		// - a vault currently has 5000 shares and 4000 equity (in quote quantums)
		// - each quote quantum is worth 5000 / 4000 = 1.25 shares
		// - a deposit of 1000 quote quantums should thus be given 1000 * 1.25 = 1250 shares
		sharesToMint = new(big.Int).Set(quantumsToDeposit)
		sharesToMint = sharesToMint.Mul(sharesToMint, existingTotalShares)
		sharesToMint = sharesToMint.Quo(sharesToMint, equity)

		// Return error if `sharesToMint` is rounded down to 0.
		if sharesToMint.Sign() == 0 {
			return errors.Wrapf(
				types.ErrZeroSharesToMint,
				"VaultId: %v, Equity: %v, Deposit: %v, TotalShares: %v, SharesToMint: %v",
				vaultId,
				equity,
				quantumsToDeposit,
				existingTotalShares,
				sharesToMint,
			)
		}
	}

	// Increase TotalShares of the vault.
	err := k.SetTotalShares(
		ctx,
		vaultId,
		types.BigIntToNumShares(
			existingTotalShares.Add(existingTotalShares, sharesToMint),
		),
	)
	if err != nil {
		return err
	}

	// Increase owner shares in the vault.
	ownerShares, exists := k.GetOwnerShares(ctx, vaultId, owner)
	if !exists {
		// Set owner shares to be sharesToMint.
		err := k.SetOwnerShares(
			ctx,
			vaultId,
			owner,
			types.BigIntToNumShares(sharesToMint),
		)
		if err != nil {
			return err
		}
	} else {
		// Increase existing owner shares by sharesToMint.
		existingOwnerShares := ownerShares.NumShares.BigInt()
		err = k.SetOwnerShares(
			ctx,
			vaultId,
			owner,
			types.BigIntToNumShares(
				existingOwnerShares.Add(existingOwnerShares, sharesToMint),
			),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
