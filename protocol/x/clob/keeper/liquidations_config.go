package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetLiquidationsConfig gets the liquidations config from state.
func (k Keeper) GetLiquidationsConfig(
	ctx sdk.Context,
) (config types.LiquidationsConfig) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.LiquidationsConfigKey))

	// The liquidations config should be set in state by the genesis logic.
	// If it's not found, then that indicates it was never set in state, which is invalid.
	if b == nil {
		panic("getLiquidationsConfig: LiquidationsConfig was never set in state")
	}

	k.cdc.MustUnmarshal(b, &config)

	return config
}

// setLiquidationsConfig sets the passed-in liquidations config in state.
// It returns an error if the provided liquidations config fails validation.
func (k Keeper) setLiquidationsConfig(
	ctx sdk.Context,
	config types.LiquidationsConfig,
) error {
	// Validate the liquidations config before writing it to state.
	if err := config.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&config)
	store.Set([]byte(types.LiquidationsConfigKey), b)

	return nil
}

// UpdateLiquidationsConfig updates the liquidations config in state.
func (k Keeper) UpdateLiquidationsConfig(
	ctx sdk.Context,
	config types.LiquidationsConfig,
) error {
	// Write the liquidations config to state.
	if err := k.setLiquidationsConfig(ctx, config); err != nil {
		return err
	}

	return nil
}

// InitializeLiquidationsConfig initializes the liquidations config in state.
// This function should only be called from the CLOB genesis.
func (k Keeper) InitializeLiquidationsConfig(
	ctx sdk.Context,
	config types.LiquidationsConfig,
) error {
	// Write the liquidations config to state.
	if err := k.setLiquidationsConfig(ctx, config); err != nil {
		return err
	}

	return nil
}
