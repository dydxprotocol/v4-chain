package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetTotalShares gets total shares.
func (k Keeper) GetTotalShares(
	ctx sdk.Context,
) (
	totalShares types.NumShares,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.TotalSharesKey))
	k.cdc.MustUnmarshal(b, &totalShares)
	return totalShares
}

// SetTotalShares sets total shares. Returns error if `totalShares` is negative.
func (k Keeper) SetTotalShares(
	ctx sdk.Context,
	totalShares types.NumShares,
) error {
	if totalShares.NumShares.Sign() < 0 {
		return types.ErrNegativeShares
	}

	b := k.cdc.MustMarshal(&totalShares)
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.TotalSharesKey), b)

	return nil
}

// GetOwnerShares gets owner shares for an owner.
func (k Keeper) GetOwnerShares(
	ctx sdk.Context,
	owner string,
) (val types.NumShares, exists bool) {
	store := k.getOwnerSharesStore(ctx)

	b := store.Get([]byte(owner))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// SetOwnerShares sets owner shares for an owner. Returns error if `ownerShares` is negative.
func (k Keeper) SetOwnerShares(
	ctx sdk.Context,
	owner string,
	ownerShares types.NumShares,
) error {
	if ownerShares.NumShares.Sign() < 0 {
		return types.ErrNegativeShares
	}

	b := k.cdc.MustMarshal(&ownerShares)
	store := k.getOwnerSharesStore(ctx)
	store.Set([]byte(owner), b)

	return nil
}

// getOwnerSharesStore returns the store for owner shares.
func (k Keeper) getOwnerSharesStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OwnerSharesKeyPrefix))
}

// GetAllOwnerShares gets all owner shares.
func (k Keeper) GetAllOwnerShares(ctx sdk.Context) []types.OwnerShare {
	allOwnerShares := []types.OwnerShare{}
	ownerSharesStore := k.getOwnerSharesStore(ctx)
	ownerSharesIterator := storetypes.KVStorePrefixIterator(ownerSharesStore, []byte{})
	defer ownerSharesIterator.Close()
	for ; ownerSharesIterator.Valid(); ownerSharesIterator.Next() {
		owner := string(ownerSharesIterator.Key())
		var ownerShares types.NumShares
		k.cdc.MustUnmarshal(ownerSharesIterator.Value(), &ownerShares)
		allOwnerShares = append(allOwnerShares, types.OwnerShare{
			Owner:  owner,
			Shares: ownerShares,
		})
	}
	return allOwnerShares
}
