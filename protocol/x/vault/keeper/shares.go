package keeper

import (
	"math/big"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetTotalShares gets TotalShares for a vault.
func (k Keeper) GetTotalShares(
	ctx sdk.Context,
	vaultId types.VaultId,
) (val types.NumShares, exists bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.TotalSharesKeyPrefix))

	b := store.Get(vaultId.ToStateKey())
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetTotalShares sets TotalShares for a vault. Returns error if `totalShares` fails validation
// or is negative.
func (k Keeper) SetTotalShares(
	ctx sdk.Context,
	vaultId types.VaultId,
	totalShares types.NumShares,
) error {
	totalSharesRat, err := totalShares.ToBigRat()
	if err != nil {
		return err
	}
	if totalSharesRat.Sign() < 0 {
		return types.ErrNegativeShares
	}

	b := k.cdc.MustMarshal(&totalShares)
	totalSharesStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.TotalSharesKeyPrefix))
	totalSharesStore.Set(vaultId.ToStateKey(), b)

	// Emit metric on TotalShares.
	totalSharesFloat, _ := totalSharesRat.Float32()
	vaultId.SetGaugeWithLabels(
		metrics.TotalShares,
		totalSharesFloat,
	)

	return nil
}

// getTotalSharesIterator returns an iterator over all TotalShares.
func (k Keeper) getTotalSharesIterator(ctx sdk.Context) storetypes.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.TotalSharesKeyPrefix))

	return storetypes.KVStorePrefixIterator(store, []byte{})
}

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
	if quantumsToDeposit.Cmp(big.NewInt(0)) <= 0 {
		return types.ErrInvalidDepositAmount
	}
	// Get existing TotalShares of the vault.
	totalShares, exists := k.GetTotalShares(ctx, vaultId)
	var existingTotalShares *big.Rat
	var err error
	if exists {
		existingTotalShares, err = totalShares.ToBigRat()
		if err != nil {
			return err
		}
	} else {
		existingTotalShares = new(big.Rat).SetInt(big.NewInt(0))
	}
	// Calculate shares to mint.
	var sharesToMint *big.Rat
	if !exists || existingTotalShares.Sign() <= 0 {
		// Mint `quoteQuantums` number of shares.
		sharesToMint = new(big.Rat).SetInt(quantumsToDeposit)
		// Initialize existingTotalShares as 0.
		existingTotalShares = new(big.Rat).SetInt(big.NewInt(0))
	} else {
		// Get vault equity.
		equity, err := k.GetVaultEquity(ctx, vaultId)
		if err != nil {
			return err
		}
		// Don't mint shares if equity is non-positive.
		if equity.Cmp(big.NewInt(0)) <= 0 {
			return types.ErrNonPositiveEquity
		}
		// Mint `deposit (in quote quantums) * existing shares / vault equity (in quote quantums)`
		// number of shares.
		// For example:
		// - a vault currently has 5000 shares and 4000 equity (in quote quantums)
		// - each quote quantum is worth 5000 / 4000 = 1.25 shares
		// - a deposit of 1000 quote quantums should thus be given 1000 * 1.25 = 1250 shares
		sharesToMint = new(big.Rat).SetInt(quantumsToDeposit)
		sharesToMint = sharesToMint.Mul(sharesToMint, existingTotalShares)
		sharesToMint = sharesToMint.Quo(sharesToMint, new(big.Rat).SetInt(equity))
	}

	// Increase TotalShares of the vault.
	err = k.SetTotalShares(
		ctx,
		vaultId,
		types.BigRatToNumShares(
			existingTotalShares.Add(existingTotalShares, sharesToMint),
		),
	)
	if err != nil {
		return err
	}

	// TODO (TRA-170): Increase owner shares.

	return nil
}
