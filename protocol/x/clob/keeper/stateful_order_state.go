package keeper

import (
	"fmt"
	"sort"
	"time"

	gometrics "github.com/armon/go-metrics"
	db "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// TODO(CLOB-739) Rename all functions in this file to StatefulOrder instead of LongTermOrder

// SetLongTermOrderPlacement sets a stateful order in state, along with information about when
// it was placed. The placed order can either be a conditional order or a long term order.
// If the order is conditional, it will be placed into the Untriggered Conditional Orders state store.
// If it is a long term order, it will be placed in the Long Term Order state store.
// If the `OrderId` doesn't exist then the `to be committed` stateful order count is incremented.
// Note the following:
// - If a stateful order placement already exists in state with `order.OrderId`, this function will overwrite it.
// - The `TransactionIndex` field will be set to the next unused transaction index for this block.
// - Triggered conditional orders should use the `TriggerConditionalOrder` write path.
func (k Keeper) SetLongTermOrderPlacement(
	ctx sdk.Context,
	order types.Order,
	blockHeight uint32,
) {
	// If this is a Short-Term order, panic.
	order.MustBeStatefulOrder()

	_, found := k.GetLongTermOrderPlacement(ctx, order.GetOrderId())

	// Get the next stateful order block transaction index, defaulting to zero if not set.
	// Note that the transaction index will always be overwritten at the end of this method.
	nextStatefulOrderTransactionIndex := k.GetNextStatefulOrderTransactionIndex(ctx)

	orderIdBytes := types.OrderIdKey(order.OrderId)
	longTermOrderPlacement := types.LongTermOrderPlacement{
		Order: order,
		PlacementIndex: types.TransactionOrdering{
			BlockHeight:      blockHeight,
			TransactionIndex: nextStatefulOrderTransactionIndex,
		},
	}
	longTermOrderPlacementBytes := k.cdc.MustMarshal(&longTermOrderPlacement)

	// For setting long term order placements, always set conditional orders to the untriggered state store.
	store, memStore := k.fetchStateStoresForOrder(ctx, order.OrderId)
	// Write the `LongTermOrderPlacement` to state.
	store.Set(orderIdBytes, longTermOrderPlacementBytes)

	// Write the `LongTermOrderPlacement` to memstore.
	memStore.Set(orderIdBytes, longTermOrderPlacementBytes)

	if !found {
		// Increment the `to be committed` stateful order count.
		k.SetToBeCommittedStatefulOrderCount(
			ctx,
			order.OrderId,
			k.GetToBeCommittedStatefulOrderCount(ctx, order.OrderId)+1,
		)

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.StatefulOrder, metrics.Count},
			1,
			[]gometrics.Label{
				metrics.GetLabelForIntValue(metrics.ClobPairId, int(order.GetClobPairId())),
				metrics.GetLabelForBoolValue(metrics.Conditional, order.OrderId.IsConditionalOrder()),
			},
		)
	}
}

// GetTriggeredConditionalOrderPlacement gets an triggered conditional order placement from the memstore.
// Returns false if no triggered conditional order exists in memstore with `orderId`.
func (k Keeper) GetTriggeredConditionalOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.LongTermOrderPlacement, found bool) {
	memStore := k.GetTriggeredConditionalOrderPlacementMemStore(ctx)

	b := memStore.Get(types.OrderIdKey(orderId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetUntriggeredConditionalOrderPlacement gets an untriggered conditional order placement from the memstore.
// Returns false if no untriggered conditional order exists in memstore with `orderId`.
func (k Keeper) GetUntriggeredConditionalOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.LongTermOrderPlacement, found bool) {
	memStore := k.GetUntriggeredConditionalOrderPlacementMemStore(ctx)

	b := memStore.Get(types.OrderIdKey(orderId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetLongTermOrderPlacement gets a long term order and the placement information from state.
// OrderId can be conditional or long term.
// Returns false if no stateful order exists in state with `orderId`.
func (k Keeper) GetLongTermOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.LongTermOrderPlacement, found bool) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	_, memStore := k.fetchStateStoresForOrder(ctx, orderId)

	b := memStore.Get(types.OrderIdKey(orderId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// DeleteLongTermOrderPlacement deletes a long term order and the placement information from state.
// This function is a no-op if no stateful order exists in state with `orderId`.
func (k Keeper) DeleteLongTermOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	orderIdBytes := types.OrderIdKey(orderId)

	store, memStore := k.fetchStateStoresForOrder(ctx, orderId)

	// Delete the `StatefulOrderPlacement` from state.
	store.Delete(orderIdBytes)

	// Delete the `StatefulOrderPlacement` from memstore.
	memStore.Delete(orderIdBytes)

	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.StatefulOrderRemoved, metrics.Count},
		1,
		[]gometrics.Label{
			metrics.GetLabelForIntValue(metrics.ClobPairId, int(orderId.GetClobPairId())),
			metrics.GetLabelForBoolValue(metrics.Conditional, orderId.IsConditionalOrder()),
		},
	)
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

// GetNextStatefulOrderTransactionIndex returns the next stateful order block transaction index
// to be used, defeaulting to zero if not set. It then increments the transaction index by one.
func (k Keeper) GetNextStatefulOrderTransactionIndex(ctx sdk.Context) (
	nextStatefulOrderTransactionIndex uint32,
) {
	nextTransactionIndexTransientStore := k.getTransientStore(ctx)
	nextStatefulOrderTransactionIndexBytes := nextTransactionIndexTransientStore.Get(
		types.KeyPrefix(
			types.NextStatefulOrderBlockTransactionIndexKey,
		),
	)
	nextStatefulOrderTransactionIndex = uint32(0)
	if nextStatefulOrderTransactionIndexBytes != nil {
		nextStatefulOrderTransactionIndex = lib.BytesToUint32(nextStatefulOrderTransactionIndexBytes)
	}
	// Set the next stateful order transaction index to be one greater than the current transaction
	// index, to ensure that transaction indexes are monotonically increasing.
	nextTransactionIndexTransientStore.Set(
		types.KeyPrefix(types.NextStatefulOrderBlockTransactionIndexKey),
		lib.Uint32ToBytes(nextStatefulOrderTransactionIndex+1),
	)
	return nextStatefulOrderTransactionIndex
}

// MustTriggerConditionalOrder triggers an untriggered conditional order. The conditional order must
// already exist in untriggered state, or else this function will panic. The LongTermOrderPlacement object
// will be removed from the untriggered state store and placed in the triggered state store.
// TODO(CLOB-746) define private R/W methods for conditional orders
func (k Keeper) MustTriggerConditionalOrder(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	// If this is not a conditional order, panic.
	orderId.MustBeConditionalOrder()

	blockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight())
	orderIdBytes := types.OrderIdKey(orderId)

	untriggeredConditionalOrderMemStore := k.GetUntriggeredConditionalOrderPlacementMemStore(ctx)
	untriggeredConditionalOrderStore := k.GetUntriggeredConditionalOrderPlacementStore(ctx)

	bytes := untriggeredConditionalOrderMemStore.Get(types.OrderIdKey(orderId))
	if bytes == nil {
		panic(
			fmt.Sprintf(
				"MustTriggerConditionalOrder: conditional order Id does not exist in Untriggered state: %+v",
				orderId,
			),
		)
	}
	var longTermOrderPlacement types.LongTermOrderPlacement
	k.cdc.MustUnmarshal(bytes, &longTermOrderPlacement)

	nextStatefulOrderTransactionIndex := k.GetNextStatefulOrderTransactionIndex(ctx)

	// Set the triggered block height and transaction index.
	longTermOrderPlacement.PlacementIndex = types.TransactionOrdering{
		BlockHeight:      blockHeight,
		TransactionIndex: nextStatefulOrderTransactionIndex,
	}

	// Write the StatefulOrderPlacement to the Triggered state store/memstore.
	longTermOrderPlacementBytes := k.cdc.MustMarshal(&longTermOrderPlacement)
	triggeredConditionalOrderMemStore := k.GetTriggeredConditionalOrderPlacementMemStore(ctx)
	triggeredConditionalOrderStore := k.GetTriggeredConditionalOrderPlacementStore(ctx)
	triggeredConditionalOrderStore.Set(orderIdBytes, longTermOrderPlacementBytes)
	triggeredConditionalOrderMemStore.Set(orderIdBytes, longTermOrderPlacementBytes)

	// Delete the `StatefulOrderPlacement` from Untriggered state store/memstore.
	untriggeredConditionalOrderStore.Delete(orderIdBytes)
	untriggeredConditionalOrderMemStore.Delete(orderIdBytes)

	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.ConditionalOrderTriggered, metrics.Count},
		1,
		append(
			orderId.GetOrderIdLabels(),
			metrics.GetLabelForIntValue(metrics.ClobPairId, int(orderId.GetClobPairId())),
		),
	)
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

// MustRemoveStatefulOrder removes an order by `OrderId` from an existing time slice. If the time slice is empty
// after removing the `OrderId`, then the time slice is pruned from state. For the `OrderId` which is removed,
// this method also calls `DeleteStatefulOrderPlacement` to remove the order placement from state.
func (k Keeper) MustRemoveStatefulOrder(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	longTermOrderPlacement, exists := k.GetLongTermOrderPlacement(ctx, orderId)
	if !exists {
		panic(fmt.Sprintf("MustRemoveStatefulOrder: order %v does not exist", orderId))
	}

	goodTilBlockTime := longTermOrderPlacement.Order.MustGetUnixGoodTilBlockTime()
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

	// Delete the Stateful order placement from state.
	k.DeleteLongTermOrderPlacement(ctx, orderId)
}

// IsConditionalOrderTriggered checks if a given order ID is triggered or untriggered in state.
// Note: If the given order ID is neither in triggered or untriggered state, function will return false.
func (k Keeper) IsConditionalOrderTriggered(
	ctx sdk.Context,
	orderId types.OrderId,
) (triggered bool) {
	// If this is not a conditional order, panic.
	orderId.MustBeConditionalOrder()
	orderIdBytes := types.OrderIdKey(orderId)
	triggeredMemstore := k.GetTriggeredConditionalOrderPlacementMemStore(ctx)
	return triggeredMemstore.Has(orderIdBytes)
}

// RemoveExpiredStatefulOrdersTimeSlices iterates all time slices from 0 until the time specified by
// `blockTime` (inclusive) and removes the time slices from state. It returns all order IDs that were removed.
func (k Keeper) RemoveExpiredStatefulOrdersTimeSlices(ctx sdk.Context, blockTime time.Time) (
	expiredOrderIds []types.OrderId,
) {
	statefulOrderPlacementIterator := k.getStatefulOrdersTimeSliceIterator(ctx, blockTime)

	defer statefulOrderPlacementIterator.Close()

	expiredOrderIds = make([]types.OrderId, 0)
	store := ctx.KVStore(k.storeKey)

	// Delete all orders from state that expire before or at `blockTime`.
	for ; statefulOrderPlacementIterator.Valid(); statefulOrderPlacementIterator.Next() {
		statefulOrderTimeSlice := types.StatefulOrderTimeSliceValue{}
		value := statefulOrderPlacementIterator.Value()
		k.cdc.MustUnmarshal(value, &statefulOrderTimeSlice)
		expiredOrderIds = append(expiredOrderIds, statefulOrderTimeSlice.OrderIds...)

		store.Delete(statefulOrderPlacementIterator.Key())
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

// GetAllPlacedStatefulOrders iterates over all stateful order placements and returns a list
// of orders, ordered by ascending time priority. Note that this only returns placed orders,
// and therefore will not return untriggered conditional orders.
func (k Keeper) GetAllPlacedStatefulOrders(ctx sdk.Context) []types.Order {
	return k.getStatefulOrders(k.getPlacedOrdersIterator(ctx))
}

// GetAllUntriggeredConditionalOrders iterates over all untriggered conditional order placements
// and returns a list of untriggered conditional orders, ordered by ascending time priority.
func (k Keeper) GetAllUntriggeredConditionalOrders(ctx sdk.Context) []types.Order {
	return k.getStatefulOrders(k.getUntriggeredConditionalOrdersIterator(ctx))
}

// getStatefulOrders takes an iterator and iterates over all stateful order placements in state.
// It returns a list of stateful order placements ordered by ascending time priority. Note this
// function handles closing the iterator.
func (k Keeper) getStatefulOrders(statefulOrderIterator db.Iterator) []types.Order {
	defer statefulOrderIterator.Close()

	statefulOrderPlacements := make([]types.LongTermOrderPlacement, 0)

	// Get all stateful order placements from state in any order.
	for ; statefulOrderIterator.Valid(); statefulOrderIterator.Next() {
		statefulOrderPlacement := types.LongTermOrderPlacement{}
		value := statefulOrderIterator.Value()
		k.cdc.MustUnmarshal(value, &statefulOrderPlacement)
		statefulOrderPlacements = append(statefulOrderPlacements, statefulOrderPlacement)
	}

	// Sort all stateful order placements in ascending time priority and return the orders.
	sort.Sort(types.SortedLongTermOrderPlacements(statefulOrderPlacements))
	sortedOrders := make([]types.Order, 0, len(statefulOrderPlacements))
	for _, orderPlacement := range statefulOrderPlacements {
		sortedOrders = append(sortedOrders, orderPlacement.Order)
	}

	return sortedOrders
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

// getPlacedOrdersIterator returns an iterator over all placed orders, which includes all
// Long-Term orders and triggered conditional orders.
func (k Keeper) getPlacedOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.KeyPrefix(types.PlacedStatefulOrderKeyPrefix),
	)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// getUntriggeredConditionalOrdersIterator returns an iterator over all untriggered conditional
// orders.
func (k Keeper) getUntriggeredConditionalOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.KeyPrefix(types.UntriggeredConditionalOrderKeyPrefix),
	)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
