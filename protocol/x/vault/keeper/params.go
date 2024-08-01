package keeper

import (
	"cosmossdk.io/store/prefix"
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
	params *types.QuotingParams,
) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(params)
	store.Set([]byte(types.DefaultQuotingParamsKey), b)

	return nil
}

// GetVaultQuotingParams returns `QuotingParams` in state for a given vault, if it exists,
// and otherwise, returns module-wide `DefaultQuotingParams`.
func (k Keeper) GetVaultQuotingParams(
	ctx sdk.Context,
	vaultId types.VaultId,
) (
	quotingParams types.QuotingParams,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.QuotingParamsKeyPrefix))

	b := store.Get(vaultId.ToStateKey())
	if b == nil {
		return k.GetDefaultQuotingParams(ctx)
	}

	k.cdc.MustUnmarshal(b, &quotingParams)
	return quotingParams
}

// SetVaultQuotingParams sets `QuotingParams` in state for a given vault.
// Returns an error if validation fails.
func (k Keeper) SetVaultQuotingParams(
	ctx sdk.Context,
	vaultId types.VaultId,
	qoutingParams types.QuotingParams,
) error {
	if err := qoutingParams.Validate(); err != nil {
		return err
	}

	b := k.cdc.MustMarshal(&qoutingParams)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.QuotingParamsKeyPrefix))
	store.Set(vaultId.ToStateKey(), b)

	return nil
}
