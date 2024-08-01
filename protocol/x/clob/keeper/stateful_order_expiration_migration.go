package keeper

import (
	"fmt"
	"time"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UnsafeMigrateOrderExpirationState migrates order expiration state from slices based on time to
// individual keys.
// Deprecated: Only intended for use in the v5.2 upgrade handler.
func (k Keeper) UnsafeMigrateOrderExpirationState(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	prefixStore := prefix.NewStore(
		store,
		[]byte(types.LegacyStatefulOrdersTimeSlicePrefix), //nolint:staticcheck
	)
	it := prefixStore.Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		time, err := sdk.ParseTimeBytes(it.Key())
		if err != nil {
			panic(fmt.Sprintf("migration failed due to malformed time: %s", it.Key()))
		}
		var orders types.StatefulOrderTimeSliceValue
		k.cdc.MustUnmarshal(it.Value(), &orders)
		for _, orderId := range orders.OrderIds {
			k.AddStatefulOrderIdExpiration(ctx, time, orderId)
		}
		prefixStore.Delete(it.Key())
	}
}

// LegacySetStatefulOrdersTimeSliceInState sets a sorted list of order IDs in state at a `goodTilBlockTime`.
// This function automatically sorts the order IDs before writing them to state.
// Deprecated: Only intended for testing MigrateOrderExpirationState.
func (k Keeper) LegacySetStatefulOrdersTimeSliceInState(
	ctx sdk.Context,
	goodTilBlockTime time.Time,
	orderIds []types.OrderId,
) {
	// Sort the order IDs.
	types.MustSortAndHaveNoDuplicates(orderIds)

	statefulOrderPlacement := types.StatefulOrderTimeSliceValue{
		OrderIds: orderIds,
	}
	b := k.cdc.MustMarshal(&statefulOrderPlacement)
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.LegacyStatefulOrdersTimeSlicePrefix), //nolint:staticcheck
	)
	store.Set(
		sdk.FormatTimeBytes(goodTilBlockTime),
		b,
	)
}
