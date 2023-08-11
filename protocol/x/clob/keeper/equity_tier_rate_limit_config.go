package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// GetEquityTierLimitConfiguration gets the equity tier limit configuration from state.
func (k Keeper) GetEquityTierLimitConfiguration(
	ctx sdk.Context,
) (config types.EquityTierLimitConfiguration) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(
		types.KeyPrefix(
			types.EquityTierLimitConfigKey,
		),
	)

	// The equity tier limit configuration should be set in state by the genesis logic.
	// If it's not found, then that indicates it was never set in state, which is invalid.
	if b == nil {
		panic("GetEquityTierLimitConfiguration: EquityTierLimitConfiguration was never set in state")
	}

	k.cdc.MustUnmarshal(b, &config)

	return config
}

// InitializeEquityTierLimit initializes the equity tier limit configuration in state.
// This function should only be called from CLOB genesis.
func (k *Keeper) InitializeEquityTierLimit(
	ctx sdk.Context,
	config types.EquityTierLimitConfiguration,
) error {
	// Validate the equity tier limit config before writing it to state.
	if err := config.Validate(); err != nil {
		return err
	}

	// Write the rate limit configuration to state.
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&config)
	store.Set(
		types.KeyPrefix(
			types.EquityTierLimitConfigKey,
		),
		b,
	)

	return nil
}
