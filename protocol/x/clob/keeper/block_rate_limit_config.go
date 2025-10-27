package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetBlockRateLimitConfiguration gets the block rate limit configuration from state.
func (k Keeper) GetBlockRateLimitConfiguration(
	ctx sdk.Context,
) (config types.BlockRateLimitConfiguration) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.BlockRateLimitConfigKey))

	// The block rate limit configuration should be set in state by the genesis logic.
	// If it's not found, then that indicates it was never set in state, which is invalid.
	if b == nil {
		panic("GetBlockRateLimitConfiguration: BlockRateLimitConfig was never set in state")
	}

	k.cdc.MustUnmarshal(b, &config)

	return config
}

// InitalizeBlockRateLimitFromStateIfExists initializes the `placeOrderRateLimiter`, `cancelOrderRateLimiter`,
// and `updateLeverageRateLimiter` from state. Should be invoked during application start and before CLOB genesis.
func (k *Keeper) InitalizeBlockRateLimitFromStateIfExists(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.BlockRateLimitConfigKey))

	if b == nil {
		return
	}

	var config types.BlockRateLimitConfiguration
	k.cdc.MustUnmarshal(b, &config)

	k.placeCancelOrderRateLimiter = rate_limit.NewPlaceCancelOrderRateLimiter(config)
	k.updateLeverageRateLimiter = rate_limit.NewUpdateLeverageRateLimiter(config)
}

// InitializeBlockRateLimit initializes the block rate limit configuration in state and uses
// the configuration to initialize the `placeOrderRateLimiter`, `cancelOrderRateLimiter`, and
// `updateLeverageRateLimiter`. This function should only be called from CLOB genesis or when a
// block rate limit configuration change is accepted via governance.
//
// Note that any previously tracked rates will be reset.
func (k *Keeper) InitializeBlockRateLimit(
	ctx sdk.Context,
	config types.BlockRateLimitConfiguration,
) error {
	// Validate the block rate limit config before writing it to state.
	if err := config.Validate(); err != nil {
		return err
	}

	// Write the rate limit configuration to state.
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&config)
	store.Set([]byte(types.BlockRateLimitConfigKey), b)

	k.placeCancelOrderRateLimiter = rate_limit.NewPlaceCancelOrderRateLimiter(config)
	k.updateLeverageRateLimiter = rate_limit.NewUpdateLeverageRateLimiter(config)

	return nil
}
