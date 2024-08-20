package keeper

import (
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
	if totalShares.NumShares.Sign() < 0 {
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

// GetOwnerShares gets owner shares for an owner in a vault.
func (k Keeper) GetOwnerShares(
	ctx sdk.Context,
	vaultId types.VaultId,
	owner string,
) (val types.NumShares, exists bool) {
	store := k.getVaultOwnerSharesStore(ctx, vaultId)

	b := store.Get([]byte(owner))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetOwnerShares sets owner shares for an owner in a vault.
func (k Keeper) SetOwnerShares(
	ctx sdk.Context,
	vaultId types.VaultId,
	owner string,
	ownerShares types.NumShares,
) error {
	if ownerShares.NumShares.Sign() < 0 {
		return types.ErrNegativeShares
	}

	b := k.cdc.MustMarshal(&ownerShares)
	store := k.getVaultOwnerSharesStore(ctx, vaultId)
	store.Set([]byte(owner), b)

	return nil
}

// getVaultOwnerSharesStore returns the store for owner shares of a given vault.
func (k Keeper) getVaultOwnerSharesStore(
	ctx sdk.Context,
	vaultId types.VaultId,
) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OwnerSharesKeyPrefix))
	return prefix.NewStore(store, vaultId.ToStateKeyPrefix())
}

// GetAllOwnerShares gets all owner shares of a given vault.
func (k Keeper) GetAllOwnerShares(
	ctx sdk.Context,
	vaultId types.VaultId,
) []*types.OwnerShare {
	allOwnerShares := []*types.OwnerShare{}
	ownerSharesStore := k.getVaultOwnerSharesStore(ctx, vaultId)
	ownerSharesIterator := storetypes.KVStorePrefixIterator(ownerSharesStore, []byte{})
	defer ownerSharesIterator.Close()
	for ; ownerSharesIterator.Valid(); ownerSharesIterator.Next() {
		owner := string(ownerSharesIterator.Key())
		var ownerShares types.NumShares
		k.cdc.MustUnmarshal(ownerSharesIterator.Value(), &ownerShares)
		allOwnerShares = append(allOwnerShares, &types.OwnerShare{
			Owner:  owner,
			Shares: &ownerShares,
		})
	}
	return allOwnerShares
}
