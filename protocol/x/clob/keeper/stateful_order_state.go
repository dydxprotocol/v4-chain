package keeper

import (
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// SetStatefulOrderPlacement sets a stateful order in state, along with information about when
// it was placed. Note the following:
// - If a stateful order placement already exists in state with `order.OrderId`, this function will overwrite it.
// - The `TransactionIndex` field will be set to the next unused transaction index for this block.
func (k Keeper) SetStatefulOrderPlacement(
	ctx sdk.Context,
	order types.Order,
	blockHeight uint32,
) {
	// If this is a Short-Term order, panic.
	order.MustBeStatefulOrder()

	// Get the next stateful order block transaction index, defaulting to zero if not set.
	// Note that the transaction index will always be overwritten at the end of this method.
	nextTransactionIndexTransientStore := k.getTransientStore(ctx)
	nextStatefulOrderTransactionIndexBytes := nextTransactionIndexTransientStore.Get(
		types.KeyPrefix(
			types.NextStatefulOrderBlockTransactionIndexKey,
		),
	)
	nextStatefulOrderTransactionIndex := uint32(0)
	if nextStatefulOrderTransactionIndexBytes != nil {
		nextStatefulOrderTransactionIndex = lib.BytesToUint32(nextStatefulOrderTransactionIndexBytes)
	}

	orderIdBytes := types.OrderIdKey(order.OrderId)
	statefulOrderPlacement := types.StatefulOrderPlacement{
		Order:            order,
		BlockHeight:      blockHeight,
		TransactionIndex: nextStatefulOrderTransactionIndex,
	}
	statefulOrderPlacementBytes := k.cdc.MustMarshal(&statefulOrderPlacement)

	// Write the `StatefulOrderPlacement` to state.
	store := k.getStatefulOrderPlacementStore(ctx)
	store.Set(orderIdBytes, statefulOrderPlacementBytes)

	// Write the `StatefulOrderPlacement` to memstore.
	memStore := k.getStatefulOrderPlacementMemStore(ctx)
	memStore.Set(orderIdBytes, statefulOrderPlacementBytes)

	// Set the next stateful order transaction index to be one greater than the current transaction
	// index, to ensure that transaction indexes are monotonically increasing.
	nextTransactionIndexTransientStore.Set(
		types.KeyPrefix(types.NextStatefulOrderBlockTransactionIndexKey),
		lib.Uint32ToBytes(nextStatefulOrderTransactionIndex+1),
	)
}

// GetStatefulOrderPlacement gets a stateful order and the placement information from state.
// Returns false if no stateful order exists in state with `orderId`.
func (k Keeper) GetStatefulOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.StatefulOrderPlacement, found bool) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	// Get the `StatefulOrderPlacement` from state.
	memStore := k.getStatefulOrderPlacementMemStore(ctx)

	b := memStore.Get(types.OrderIdKey(orderId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// DeleteStatefulOrderPlacement deletes a stateful order and the placement information from state.
// This function is a no-op if no stateful order exists in state with `orderId`.
func (k Keeper) DeleteStatefulOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	orderIdBytes := types.OrderIdKey(orderId)

	// Delete the `StatefulOrderPlacement` from state.
	store := k.getStatefulOrderPlacementStore(ctx)
	store.Delete(orderIdBytes)

	// Delete the `StatefulOrderPlacement` from memstore.
	memStore := k.getStatefulOrderPlacementMemStore(ctx)
	memStore.Delete(orderIdBytes)
}

// GetStatefulOrdersTimeSlice gets a slice of stateful order IDs that expire at `goodTilBlockTime`,
// sorted by order ID.
func (k Keeper) GetStatefulOrdersTimeSlice(ctx sdk.Context, goodTilBlockTime time.Time) (
	orderIds []types.OrderId,
) {
	store := k.getStatefulOrdersTimeSliceStore(ctx)
	statefulOrdersTimeSliceBytes := store.Get(types.GetTimeSliceKey(goodTilBlockTime))

	// If there are no stateful orders that expire at this block, then return an empty slice.
	if statefulOrdersTimeSliceBytes == nil {
		return []types.OrderId{}
	}

	var longTermOrders types.StatefulOrderTimeSliceValue
	k.cdc.MustUnmarshal(statefulOrdersTimeSliceBytes, &longTermOrders)

	return longTermOrders.OrderIds
}

// MustAddOrderToStatefulOrdersTimeSlice adds a new `OrderId` to an existing time slice, or creates a new time slice
// containing the `OrderId` and writes it to state. It first sorts all order IDs before writing them
// to state to avoid non-determinism issues.
func (k Keeper) MustAddOrderToStatefulOrdersTimeSlice(
	ctx sdk.Context,
	goodTilBlockTime time.Time,
	orderId types.OrderId,
) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	longTermOrdersExpiringAtTime := k.GetStatefulOrdersTimeSlice(ctx, goodTilBlockTime)

	// Panic if this order ID is already written to state.
	for _, foundOrderId := range longTermOrdersExpiringAtTime {
		if orderId == foundOrderId {
			panic(
				fmt.Sprintf(
					"MustAddOrderToStatefulOrdersTimeSlice: order ID %v is already contained in state for time %v",
					orderId,
					goodTilBlockTime,
				),
			)
		}
	}

	longTermOrdersExpiringAtTime = append(longTermOrdersExpiringAtTime, orderId)

	k.setStatefulOrdersTimeSliceInState(ctx, goodTilBlockTime, longTermOrdersExpiringAtTime)
}

// MustRemoveStatefulOrder removes an order by `OrderId` from an existing time slice.
// If the time slice is empty after removing the `OrderId`, then the time slice is pruned from state.
// For the `OrderId` which is removed, this method also calls `DeleteStatefulOrderPlacement`
// and `RemoveOrderFillAmount`.
func (k Keeper) MustRemoveStatefulOrder(
	ctx sdk.Context,
	goodTilBlockTime time.Time,
	orderId types.OrderId,
) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	longTermOrdersExpiringAtTime := k.GetStatefulOrdersTimeSlice(ctx, goodTilBlockTime)
	updatedStatefulOrdersExpiringAtTime := make([]types.OrderId, 0, len(longTermOrdersExpiringAtTime))

	// Loop through all order IDs and remove any that equal `orderId`.
	for _, longTermOrderId := range longTermOrdersExpiringAtTime {
		if longTermOrderId != orderId {
			updatedStatefulOrdersExpiringAtTime = append(updatedStatefulOrdersExpiringAtTime, longTermOrderId)
		}
	}

	// Panic if the length of the new list is not one less, since that indicates no element was removed.
	if len(longTermOrdersExpiringAtTime) != len(updatedStatefulOrdersExpiringAtTime)+1 {
		panic(
			fmt.Sprintf(
				"MustRemoveStatefulOrder: order ID %v is not in state for time %v",
				orderId,
				goodTilBlockTime,
			),
		)
	}

	// If `updatedStatefulOrdersExpiringAtTime` is empty, remove the key prefix from state.
	// Else, set the updated list of order IDs in state.
	if len(updatedStatefulOrdersExpiringAtTime) == 0 {
		store := k.getStatefulOrdersTimeSliceStore(ctx)
		store.Delete(types.GetTimeSliceKey(goodTilBlockTime))
	} else {
		k.setStatefulOrdersTimeSliceInState(ctx, goodTilBlockTime, updatedStatefulOrdersExpiringAtTime)
	}

	// Remove the order fill amount from state.
	k.RemoveOrderFillAmount(ctx, orderId)

	// Delete the stateful order placement from state.
	k.DeleteStatefulOrderPlacement(ctx, orderId)
}

// RemoveExpiredStatefulOrdersTimeSlices iterates all time slices from 0 until the time specified by
// `blockTime` (inclusive) and removes the time slices from state. It returns all order IDs that were removed.
func (k Keeper) RemoveExpiredStatefulOrdersTimeSlices(ctx sdk.Context, blockTime time.Time) (
	expiredOrderIds []types.OrderId,
) {
	longTermOrderIterator := k.getStatefulOrdersTimeSliceIterator(ctx, blockTime)

	defer longTermOrderIterator.Close()

	expiredOrderIds = make([]types.OrderId, 0)
	store := ctx.KVStore(k.storeKey)

	// Delete all orders from state that expire before or at `blockTime`.
	for ; longTermOrderIterator.Valid(); longTermOrderIterator.Next() {
		longTermOrderTimeSlice := types.StatefulOrderTimeSliceValue{}
		value := longTermOrderIterator.Value()
		k.cdc.MustUnmarshal(value, &longTermOrderTimeSlice)
		expiredOrderIds = append(expiredOrderIds, longTermOrderTimeSlice.OrderIds...)

		store.Delete(longTermOrderIterator.Key())
	}

	return expiredOrderIds
}

// SetBlockTimeForLastCommittedBlock writes the block time of the previously committed block
// to state. This is necessary for consensus validation of stateful orders,
// since `order.GoodTIlBlockTime` is always validated against the previous block's timestamp
// and cannot be validated against the current block's timestamp.
// Note that this function overwrites the current value and does not validate `ctx.BlockTime()`
// against the current value in state.
func (k Keeper) SetBlockTimeForLastCommittedBlock(
	ctx sdk.Context,
) {
	blockTime := ctx.BlockTime()
	if blockTime.IsZero() {
		panic("Block-time is zero")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(
		types.KeyPrefix(types.LastCommittedBlockTimeKey),
		sdk.FormatTimeBytes(blockTime),
	)
}

// MustGetBlockTimeForLastCommittedBlock returns the block time of the previously commited block.
// Panics if the previously committed block time is not found.
func (k Keeper) MustGetBlockTimeForLastCommittedBlock(
	ctx sdk.Context,
) (
	blockTime time.Time,
) {
	store := ctx.KVStore(k.storeKey)
	time, err := sdk.ParseTimeBytes(
		store.Get(types.KeyPrefix(types.LastCommittedBlockTimeKey)),
	)

	if err != nil {
		panic("Failed to get the block time of the previously committed block")
	}
	return time.UTC()
}

// GetAllStatefulOrders iterates over all stateful order placements and returns a list
// of orders, ordered by ascending time priority.
func (k Keeper) GetAllStatefulOrders(ctx sdk.Context) []types.Order {
	statefulOrderPlacementIterator := k.getStatefulOrderPlacementIterator(ctx)

	defer statefulOrderPlacementIterator.Close()

	statefulOrderPlacements := make([]types.StatefulOrderPlacement, 0)

	// Get all stateful order placements from state in any order.
	for ; statefulOrderPlacementIterator.Valid(); statefulOrderPlacementIterator.Next() {
		statefulOrderPlacement := types.StatefulOrderPlacement{}
		value := statefulOrderPlacementIterator.Value()
		k.cdc.MustUnmarshal(value, &statefulOrderPlacement)
		statefulOrderPlacements = append(statefulOrderPlacements, statefulOrderPlacement)
	}

	// Sort all stateful order placements in ascending time priority and return the orders.
	sort.Sort(types.SortedStatefulOrderPlacements(statefulOrderPlacements))
	sortedOrders := make([]types.Order, 0, len(statefulOrderPlacements))
	for _, orderPlacement := range statefulOrderPlacements {
		sortedOrders = append(sortedOrders, orderPlacement.Order)
	}

	return sortedOrders
}

// DoesStatefulOrderExistInState returns true if the stateful order exists in state, false if not.
// It checks the order hashes of the order stored in state to determine equality.
// Note this function panics if called with a Short-Term order.
func (k Keeper) DoesStatefulOrderExistInState(
	ctx sdk.Context,
	order types.Order,
) bool {
	// If this is a Short-Term order, panic.
	order.MustBeStatefulOrder()

	orderInState, found := k.GetStatefulOrderPlacement(ctx, order.OrderId)
	return found && order.GetOrderHash() == orderInState.Order.GetOrderHash()
}

// setStatefulOrdersTimeSliceInState sets a sorted list of order IDs in state at a `goodTilBlockTime`.
// This function automatically sorts the order IDs before writing them to state.
func (k Keeper) setStatefulOrdersTimeSliceInState(
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
	store := k.getStatefulOrdersTimeSliceStore(ctx)
	store.Set(
		types.GetTimeSliceKey(
			goodTilBlockTime,
		),
		b,
	)
}

// getStatefulOrdersTimeSliceIterator returns an iterator over all stateful order time slice values
// from time 0 until `endTime`.
func (k Keeper) getStatefulOrdersTimeSliceIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	startKey :=
		types.KeyPrefix(types.StatefulOrdersTimeSlicePrefix)
	endKey := append(
		startKey,
		sdk.InclusiveEndBytes(
			types.GetTimeSliceKey(
				endTime,
			),
		)...,
	)
	return store.Iterator(
		startKey,
		endKey,
	)
}

// getStatefulOrderPlacementIterator returns an iterator over all stateful orders.
func (k Keeper) getStatefulOrderPlacementIterator(ctx sdk.Context) sdk.Iterator {
	store := k.getStatefulOrderPlacementStore(ctx)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
