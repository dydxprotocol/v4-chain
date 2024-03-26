package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
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

	return nil
}

// getTotalSharesIterator returns an iterator over all TotalShares.
func (k Keeper) getTotalSharesIterator(ctx sdk.Context) storetypes.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.TotalSharesKeyPrefix))

	return storetypes.KVStorePrefixIterator(store, []byte{})
}
