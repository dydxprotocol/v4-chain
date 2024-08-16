package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetDefaultQuotingParams returns `DefaultQuotingParams` in state.
func (k Keeper) GetDefaultQuotingParams(
	ctx sdk.Context,
) (
	params types.QuotingParams,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.DefaultQuotingParamsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// SetDefaultQuotingParams updates `DefaultQuotingParams` in state.
// Returns an error if validation fails.
func (k Keeper) SetDefaultQuotingParams(
	ctx sdk.Context,
	params types.QuotingParams,
) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.DefaultQuotingParamsKey), b)

	return nil
}

// GetVaultParams returns `VaultParams` in state for a given vault.
func (k Keeper) GetVaultParams(
	ctx sdk.Context,
	vaultId types.VaultId,
) (
	vaultParams types.VaultParams,
	exists bool,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultParamsKeyPrefix))

	b := store.Get(vaultId.ToStateKey())
	if b == nil {
		return vaultParams, false
	}

	k.cdc.MustUnmarshal(b, &vaultParams)
	return vaultParams, true
}

// SetVaultParams sets `VaultParams` in state for a given vault.
// Returns an error if validation fails.
func (k Keeper) SetVaultParams(
	ctx sdk.Context,
	vaultId types.VaultId,
	vaultParams types.VaultParams,
) error {
	if err := vaultParams.Validate(); err != nil {
		return err
	}

	b := k.cdc.MustMarshal(&vaultParams)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultParamsKeyPrefix))
	store.Set(vaultId.ToStateKey(), b)

	return nil
}

// getVaultParamsIterator returns an iterator over all VaultParams.
func (k Keeper) getVaultParamsIterator(ctx sdk.Context) storetypes.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultParamsKeyPrefix))
	return storetypes.KVStorePrefixIterator(store, []byte{})
}

// GetVaultQuotingParams returns quoting parameters for a given vault, which is
// - `VaultParams.QuotingParams` if set
// - `DefaultQuotingParams` otherwise
// `exists` is false if `VaultParams` does not exist for the given vault.
func (k Keeper) GetVaultQuotingParams(
	ctx sdk.Context,
	vaultId types.VaultId,
) (
	params types.QuotingParams,
	exists bool,
) {
	vaultParams, exists := k.GetVaultParams(ctx, vaultId)
	if !exists {
		return params, false
	}
	if vaultParams.QuotingParams == nil {
		return k.GetDefaultQuotingParams(ctx), true
	} else {
		return *vaultParams.QuotingParams, true
	}
}

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
