package keeper

import (
	"math/big"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// This file contains all functions used to fetch deprecated state for the purpose of upgrade
// functions.

// Deprecated state keys
const (
	// Deprecated: For use by the v7.x upgrade handler
	// QuotingParamsKeyPrefix is the prefix to retrieve all QuotingParams.
	// QuotingParams store: vaultId VaultId -> QuotingParams.
	QuotingParamsKeyPrefix = "QuotingParams:"

	// Deprecated: For use by the v7.x upgrade handler
	// TotalSharesKeyPrefix is the prefix to retrieve all TotalShares.
	TotalSharesKeyPrefix = "TotalShares:"

	// Deprecated: For use by the v7.x upgrade handler
	// OwnerSharesKeyPrefix is the prefix to retrieve all OwnerShares.
	// OwnerShares store: vaultId VaultId -> owner string -> shares NumShares.
	OwnerSharesKeyPrefix = "OwnerShares:"
)

// v5.x state, used for v6.x upgrade.

// UnsafeGetParams returns `Params` in state.
// Used for v6.x upgrade handler.
func (k Keeper) UnsafeGetParams(
	ctx sdk.Context,
) (
	params types.QuotingParams,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte("Params"))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// UnsafeDeleteParams deletes `Params` in state.
// Used for v6.x upgrade handler.
func (k Keeper) UnsafeDeleteParams(
	ctx sdk.Context,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte("Params"))
}

// v6.x state, used for v7.x upgrade

// UnsafeSetQuotingParams sets quoting parameters for a given vault from state.
// Used for v7.x upgrade handler
func (k Keeper) UnsafeSetQuotingParams(
	ctx sdk.Context,
	vaultId types.VaultId,
	quotingParams types.QuotingParams,
) error {
	if err := quotingParams.Validate(); err != nil {
		return err
	}

	b := k.cdc.MustMarshal(&quotingParams)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(QuotingParamsKeyPrefix))
	store.Set(vaultId.ToStateKey(), b)

	return nil
}

// UnsafeGetQuotingParams returns quoting parameters for a given vault from state.
// Used for v7.x upgrade handler
func (k Keeper) UnsafeGetQuotingParams(
	ctx sdk.Context,
	vaultId types.VaultId,
) (
	quotingParams types.QuotingParams,
	exists bool,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(QuotingParamsKeyPrefix))

	b := store.Get(vaultId.ToStateKey())
	if b == nil {
		return quotingParams, false
	}

	k.cdc.MustUnmarshal(b, &quotingParams)
	return quotingParams, true
}

// UnsafeDeleteQuotingParams deletes quoting parameters for a given vault from state.
// Used for v7.x upgrade handler
func (k Keeper) UnsafeDeleteQuotingParams(
	ctx sdk.Context,
	vaultId types.VaultId,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(QuotingParamsKeyPrefix))
	b := store.Get(vaultId.ToStateKey())
	if b == nil {
		return
	}
	store.Delete(vaultId.ToStateKey())
}

// UnsafeGetAllVaultIds returns all vault ids from state using the deprecated total shares
// state.
func (k Keeper) UnsafeGetAllVaultIds(ctx sdk.Context) []types.VaultId {
	vaultIds := []types.VaultId{}
	totalSharesStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(TotalSharesKeyPrefix))
	totalSharesIterator := storetypes.KVStorePrefixIterator(totalSharesStore, []byte{})
	defer totalSharesIterator.Close()
	for ; totalSharesIterator.Valid(); totalSharesIterator.Next() {
		vaultId, err := types.GetVaultIdFromStateKey(totalSharesIterator.Key())
		if err != nil {
			panic(err)
		}
		vaultIds = append(vaultIds, *vaultId)
	}
	return vaultIds
}

// GetTotalShares gets TotalShares for a vault.
// Deprecated and used for v7.x upgrade handler
func (k Keeper) UnsafeGetTotalShares(
	ctx sdk.Context,
	vaultId types.VaultId,
) (val types.NumShares, exists bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(TotalSharesKeyPrefix))

	b := store.Get(vaultId.ToStateKey())
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// UnsafeGetAllOwnerShares gets all owner shares of a given vault.
// Deprecated and used for v7.x upgrade handler
func (k Keeper) UnsafeGetAllOwnerShares(
	ctx sdk.Context,
	vaultId types.VaultId,
) []*types.OwnerShare {
	allOwnerShares := []*types.OwnerShare{}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.OwnerSharesKeyPrefix))
	ownerSharesStore := prefix.NewStore(store, vaultId.ToStateKeyPrefix())
	ownerSharesIterator := storetypes.KVStorePrefixIterator(ownerSharesStore, []byte{})
	defer ownerSharesIterator.Close()
	for ; ownerSharesIterator.Valid(); ownerSharesIterator.Next() {
		owner := string(ownerSharesIterator.Key())
		var ownerShares types.NumShares
		k.cdc.MustUnmarshal(ownerSharesIterator.Value(), &ownerShares)
		allOwnerShares = append(allOwnerShares, &types.OwnerShare{
			Owner:  owner,
			Shares: ownerShares,
		})
	}
	return allOwnerShares
}

// UnsafeGetAllOwnerEquities returns equity that belongs to each owner across all vaults
// using the deprecated owner shares and total shares state.
// Deprecated and used for v7.x upgrade handler
func (k Keeper) UnsafeGetAllOwnerEquities(ctx sdk.Context) map[string]*big.Rat {
	ownerEquities := make(map[string]*big.Rat)
	totalSharesStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(TotalSharesKeyPrefix))
	totalSharesIterator := storetypes.KVStorePrefixIterator(totalSharesStore, []byte{})
	defer totalSharesIterator.Close()
	for ; totalSharesIterator.Valid(); totalSharesIterator.Next() {
		vaultId, err := types.GetVaultIdFromStateKey(totalSharesIterator.Key())
		if err != nil {
			panic(err)
		}

		var vaultTotalShares types.NumShares
		k.cdc.MustUnmarshal(totalSharesIterator.Value(), &vaultTotalShares)
		bigVaultTotalShares := vaultTotalShares.NumShares.BigInt()
		vaultEquity, err := k.GetVaultEquity(ctx, *vaultId)
		if err != nil {
			panic(err)
		}

		ownerShares := k.UnsafeGetAllOwnerShares(ctx, *vaultId)
		for _, ownerShare := range ownerShares {
			// owner equity in this vault = vault equity * owner shares / vault total shares
			ownerEquity := new(big.Rat).SetInt(vaultEquity)
			ownerEquity.Mul(
				ownerEquity,
				new(big.Rat).SetInt(ownerShare.Shares.NumShares.BigInt()),
			)
			ownerEquity.Quo(
				ownerEquity,
				new(big.Rat).SetInt(bigVaultTotalShares),
			)

			if e, ok := ownerEquities[ownerShare.Owner]; ok {
				ownerEquities[ownerShare.Owner] = e.Add(e, ownerEquity)
			} else {
				ownerEquities[ownerShare.Owner] = ownerEquity
			}
		}
	}

	return ownerEquities
}

// UnsafeDeleteVaultTotalShares deletes total shares of a given vault from state.
// Used for v7.x upgrade handler
func (k Keeper) UnsafeDeleteAllVaultTotalShares(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(TotalSharesKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}

// UnsafeDeleteVaultOwnerShares deletes all owner shares of a given vault from state.
// Used for v7.x upgrade handler
func (k Keeper) UnsafeDeleteAllVaultOwnerShares(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(OwnerSharesKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}
