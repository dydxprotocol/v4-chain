package keeper

import (
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
)

// v5.x state, used for v6.x upgrade.

// Deprecated: Only used to get vault params as they were in v5.x.
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

// Deprecated: Only used to delete params as they were in v5.x.
// UnsafeDeleteParams deletes `Params` in state.
// Used for v6.x upgrade handler.
func (k Keeper) UnsafeDeleteParams(
	ctx sdk.Context,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte("Params"))
}

// v6.x state, used for v7.x upgrade

// Deprecated: Only used to set quoting params as they were in v6.x.
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

// Deprecated: Only used to get quoting params as they were in v6.x.
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

// Deprecated: Only used to set quoting params as they were in v6.x.
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

// Deprecated: Only used to fetch vault ids as they were in v6.x
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
