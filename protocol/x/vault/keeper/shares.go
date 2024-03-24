package keeper

import (
	"math/big"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
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

// SetTotalShares sets TotalShares for a vault. Returns error if `totalShares` is negative.
func (k Keeper) SetTotalShares(
	ctx sdk.Context,
	vaultId types.VaultId,
	totalShares types.NumShares,
) error {
	if totalShares.NumShares.Cmp(dtypes.NewInt(0)) == -1 {
		return types.ErrNegativeShares
	}

	b := k.cdc.MustMarshal(&totalShares)
	totalSharesStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.TotalSharesKeyPrefix))
	totalSharesStore.Set(vaultId.ToStateKey(), b)

	// Emit metric on TotalShares.
	vaultId.SetGaugeWithLabels(
		metrics.TotalShares,
		float32(totalShares.NumShares.BigInt().Uint64()),
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
	totalShares, exists := k.GetTotalShares(ctx, vaultId)
	existingTotalShares := totalShares.NumShares.BigInt()
	sharesToMint := big.NewInt(0)
	if !exists || existingTotalShares.Cmp(big.NewInt(0)) < 1 {
		// Mint `quoteQuantums` number of shares.
		sharesToMint = quantumsToDeposit
		// Initialize existingTotalShares as 0.
		existingTotalShares = big.NewInt(0)
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
		// For shares calculation to be accurate, scale everything by denominator and mint
		// numerator number of shares.
		// For example:
		// - a vault currently has 1 share and 1000 equity (in quote quantums)
		// - a deposit of 1 quote quantum is worth 1 * 1 / 1000 shares
		// - to be accurate, we credit 1 share and scale existing shares by 1000
		newShares := new(big.Rat).SetInt(quantumsToDeposit)
		newShares = newShares.Mul(newShares, new(big.Rat).SetInt(existingTotalShares))
		newShares = newShares.Quo(newShares, new(big.Rat).SetInt(equity))
		existingTotalShares = existingTotalShares.Mul(existingTotalShares, newShares.Denom())
		sharesToMint = newShares.Num()
	}

	// Increase TotalShares of the vault.
	err := k.SetTotalShares(
		ctx,
		vaultId,
		types.NumShares{
			NumShares: dtypes.NewIntFromBigInt(
				sharesToMint.Add(sharesToMint, existingTotalShares),
			),
		},
	)
	if err != nil {
		return err
	}

	// TODO (TRA-170): Increase owner shares.

	return nil
}
