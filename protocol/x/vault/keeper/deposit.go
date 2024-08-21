package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// MintShares mints shares for `owner` based on `quantumsToDeposit` by:
// 1. Increasing total shares.
// 2. Increasing owner shares for given `owner`.
func (k Keeper) MintShares(
	ctx sdk.Context,
	owner string,
	quantumsToDeposit *big.Int,
) error {
	// Quantums to deposit should be positive.
	if quantumsToDeposit.Sign() <= 0 {
		return types.ErrInvalidDepositAmount
	}
	// Get existing TotalShares of the vault.
	existingTotalShares := k.GetTotalShares(ctx).NumShares.BigInt()
	// Calculate shares to mint.
	var sharesToMint *big.Int
	if existingTotalShares.Sign() <= 0 {
		// Mint `quoteQuantums` number of shares.
		sharesToMint = new(big.Int).Set(quantumsToDeposit)
	} else {
		// Get megavault equity.
		equity, err := k.GetMegavaultEquity(ctx)
		if err != nil {
			return err
		}
		// Don't mint shares if equity is non-positive.
		if equity.Sign() <= 0 {
			return types.ErrNonPositiveEquity
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
			return types.ErrZeroSharesToMint
		}
	}

	// Increase total shares.
	err := k.SetTotalShares(
		ctx,
		types.BigIntToNumShares(
			existingTotalShares.Add(existingTotalShares, sharesToMint),
		),
	)
	if err != nil {
		return err
	}

	// Increase owner shares.
	ownerShares, exists := k.GetOwnerShares(ctx, owner)
	if !exists {
		// Set owner shares to be sharesToMint.
		err := k.SetOwnerShares(
			ctx,
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
