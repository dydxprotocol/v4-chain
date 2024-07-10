package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// GetParams returns `Params` in state.
func (k Keeper) GetParams(
	ctx sdk.Context,
) (
	params types.Params,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.ParamsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// SetParams updates `Params` in state.
// Returns an error if validation fails.
func (k Keeper) SetParams(
	ctx sdk.Context,
	params types.Params,
) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.ParamsKey), b)

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
