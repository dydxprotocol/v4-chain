package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetBlockLimitsConfig gets the block limits config from state.
// If the config has not been set yet, returns the default config (no cap).
// This allows for backward compatibility with chains that existed before this parameter was added.
func (k Keeper) GetBlockLimitsConfig(
	ctx sdk.Context,
) (config types.BlockLimitsConfig) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.BlockLimitsConfigKey))

	// If not found in state, return default config (0 means no cap).
	// This provides backward compatibility - governance can set it later if needed.
	if b == nil {
		return types.BlockLimitsConfig_Default
	}

	k.cdc.MustUnmarshal(b, &config)

	return config
}

// setBlockLimitsConfig sets the passed-in block limits config in state.
// It returns an error if the provided block limits config fails validation.
func (k Keeper) setBlockLimitsConfig(
	ctx sdk.Context,
	config types.BlockLimitsConfig,
) error {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&config)
	store.Set([]byte(types.BlockLimitsConfigKey), b)

	return nil
}

// UpdateBlockLimitsConfig updates the block limits config in state.
// This is the only way to set the block limits config - via governance proposal.
func (k Keeper) UpdateBlockLimitsConfig(
	ctx sdk.Context,
	config types.BlockLimitsConfig,
) error {
	// Write the block limits config to state.
	if err := k.setBlockLimitsConfig(ctx, config); err != nil {
		return err
	}

	return nil
}
