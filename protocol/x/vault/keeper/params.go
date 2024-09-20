package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
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

	if vaultParams.Status == types.VaultStatus_VAULT_STATUS_DEACTIVATED {
		vaultEquity, err := k.GetVaultEquity(ctx, vaultId)
		if err != nil {
			return err
		}
		if vaultEquity.Sign() > 0 {
			return types.ErrDeactivatePositiveEquityVault
		}
	}

	b := k.cdc.MustMarshal(&vaultParams)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultParamsKeyPrefix))
	store.Set(vaultId.ToStateKey(), b)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeUpsertVault,
		indexerevents.UpsertVaultEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewUpsertVaultEvent(
				vaultId.ToModuleAccountAddress(),
				vaultId.Number,
				vaultParams.Status,
			),
		),
	)

	return nil
}

// getVaultParamsIterator returns an iterator over all VaultParams.
func (k Keeper) getVaultParamsIterator(ctx sdk.Context) storetypes.Iterator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.VaultParamsKeyPrefix))
	return storetypes.KVStorePrefixIterator(store, []byte{})
}

// GetVaultAndQuotingParams returns vault params and quoting parameters for a given vault.
// Quoting parameters is
// - `VaultParams.QuotingParams` if set
// - `DefaultQuotingParams` otherwise
// `exists` is false if `VaultParams` does not exist for the given vault.
func (k Keeper) GetVaultAndQuotingParams(
	ctx sdk.Context,
	vaultId types.VaultId,
) (
	vaultParams types.VaultParams,
	quotingParams types.QuotingParams,
	exists bool,
) {
	vaultParams, exists = k.GetVaultParams(ctx, vaultId)
	if !exists {
		return vaultParams, quotingParams, false
	}
	if vaultParams.QuotingParams == nil {
		return vaultParams, k.GetDefaultQuotingParams(ctx), true
	} else {
		return vaultParams, *vaultParams.QuotingParams, true
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

// GetOperatorParams returns `OperatorParams` in state.
func (k Keeper) GetOperatorParams(
	ctx sdk.Context,
) (
	params types.OperatorParams,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.OperatorParamsKey))
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// SetOperatorParams sets `OperatorParams` in state.
// Returns an error if validation fails.
func (k Keeper) SetOperatorParams(
	ctx sdk.Context,
	params types.OperatorParams,
) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.OperatorParamsKey), b)

	return nil
}
