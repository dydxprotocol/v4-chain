package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetTradePricesForPerpetual gets the maximum and minimum traded prices for a perpetual for the
// current block.
// These prices are intended to be used for improved conditional order triggering in the EndBlocker.
func (k Keeper) GetTradePricesForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
) (
	minTradePriceSubticks types.Subticks,
	maxTradePriceSubticks types.Subticks,
	found bool,
) {
	minTradePriceSubticks, found = k.getMinTradePriceForPerpetual(ctx, perpetualId)
	if !found {
		return 0, 0, false
	}

	maxTradePriceSubticks, found = k.getMaxTradePriceForPerpetual(ctx, perpetualId)
	return minTradePriceSubticks, maxTradePriceSubticks, found
}

// getMinTradePriceForPerpetual gets the min trade price for a perpetual.
func (k Keeper) getMinTradePriceForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
) (subticks types.Subticks, found bool) {
	store := k.GetMinTradePriceStore(ctx)

	b := store.Get(lib.Uint32ToKey(perpetualId))
	if b == nil {
		return 0, false
	}

	result := gogotypes.UInt64Value{Value: 0}
	k.cdc.MustUnmarshal(b, &result)
	return types.Subticks(result.Value), true
}

// getMaxTradePriceForPerpetual gets the max trade price for a perpetual.
func (k Keeper) getMaxTradePriceForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
) (subticks types.Subticks, found bool) {
	store := k.GetMaxTradePriceStore(ctx)

	b := store.Get(lib.Uint32ToKey(perpetualId))
	if b == nil {
		return 0, false
	}

	result := gogotypes.UInt64Value{Value: 0}
	k.cdc.MustUnmarshal(b, &result)
	return types.Subticks(result.Value), true
}

// SetTradePricesForPerpetual sets the maximum and minimum traded prices for a perpetual for the current
// block.
// Note that this method updates the transient store and is meant to be called in `DeliverTx` when
// matches are persisted to state.
func (k Keeper) SetTradePricesForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
	subticks types.Subticks,
) {
	minTradePriceSubticks, maxTradePriceSubticks, found := k.GetTradePricesForPerpetual(ctx, perpetualId)
	key := lib.Uint32ToKey(perpetualId)
	value := gogotypes.UInt64Value{Value: subticks.ToUint64()}
	b := k.cdc.MustMarshal(&value)

	// Update the min price if not previously saved or if the new price is lower.
	if !found || subticks < minTradePriceSubticks {
		store := k.GetMinTradePriceStore(ctx)
		store.Set(key, b)
	}

	// Update the max price if not previously saved or if the new price is higher.
	if !found || subticks > maxTradePriceSubticks {
		store := k.GetMaxTradePriceStore(ctx)
		store.Set(key, b)
	}
}
