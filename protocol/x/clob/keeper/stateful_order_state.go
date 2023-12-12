package keeper

import (
	"fmt"
	"sort"
	"time"

	gometrics "github.com/armon/go-metrics"
	db "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// TODO(CLOB-739) Rename all functions in this file to StatefulOrder instead of LongTermOrder

// SetLongTermOrderPlacement sets a stateful order in state, along with information about when
// it was placed. The placed order can either be a conditional order or a long term order.
// If the order is conditional, it will be placed into the Untriggered Conditional Orders state store.
// If it is a long term order, it will be placed in the Long Term Order state store.
// If the `OrderId` doesn't exist then the stateful order count is incremented.
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
	orderKey := order.OrderId.ToStateKey()

	// Write the `LongTermOrderPlacement` to state.
	store.Set(orderKey, longTermOrderPlacementBytes)

	// Write the `LongTermOrderPlacement` to memstore.
	memStore.Set(orderKey, longTermOrderPlacementBytes)

	if !found {
		// Increment the stateful order count.
		k.SetStatefulOrderCount(
			ctx,
			order.OrderId.SubaccountId,
			k.GetStatefulOrderCount(ctx, order.OrderId.SubaccountId)+1,
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

	b := memStore.Get(orderId.ToStateKey())
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

	b := memStore.Get(orderId.ToStateKey())
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

	b := memStore.Get(orderId.ToStateKey())
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// DeleteLongTermOrderPlacement deletes a long term order and the placement information from state
// and decrements the stateful order count if the `orderId` exists.
// This function is a no-op if no stateful order exists in state with `orderId`.
func (k Keeper) DeleteLongTermOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	// If this is a Short-Term order, panic.
	orderId.MustBeStatefulOrder()

	store, memStore := k.fetchStateStoresForOrder(ctx, orderId)

	// Note that since store reads/writes can cost gas we need to ensure that the number of operations is the
	// same regardless of whether the memstore has the order or not.
	count := k.GetStatefulOrderCount(ctx, orderId.SubaccountId)
	orderKey := orderId.ToStateKey()
	if memStore.Has(orderKey) {
		if count == 0 {
			k.Logger(ctx).Error(
				"Stateful order count is zero but order is in the memstore. Underflow",
				"orderId", log.NewLazySprintf("%+v", orderId),
			)
		} else {
			count--
		}
	}

	// Delete the `StatefulOrderPlacement` from state.
	store.Delete(orderKey)

	// Delete the `StatefulOrderPlacement` from memstore.
	memStore.Delete(orderKey)

	// Set the count.
	k.SetStatefulOrderCount(ctx, orderId.SubaccountId, count)

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
	statefulOrdersTimeSliceBytes := store.Get(sdk.FormatTimeBytes(goodTilBlockTime))

	// If there are no stateful orders that expire at this block, then return an empty slice.
	if statefulOrdersTimeSliceBytes == nil {
		return []types.OrderId{}
	}

	var longTermOrders types.StatefulOrderTimeSliceValue
	k.cdc.MustUnmarshal(statefulOrdersTimeSliceBytes, &longTermOrders)

	return longTermOrders.OrderIds
}

// GetNextStatefulOrderTransactionIndex returns the next stateful order block transaction index
// to be used, defaulting to zero if not set. It then increments the transaction index by one.
func (k Keeper) GetNextStatefulOrderTransactionIndex(ctx sdk.Context) uint32 {
	// Get the existing value
	nextTransactionIndexTransientStore := k.getTransientStore(ctx)
	b := nextTransactionIndexTransientStore.Get(
		[]byte(types.NextStatefulOrderBlockTransactionIndexKey),
	)
	index := gogotypes.UInt32Value{Value: 0}
	if b != nil {
		k.cdc.MustUnmarshal(b, &index)
	}
	oldValue := index.Value

	// Set the next stateful order transaction index to be one greater than the current transaction
	// index, to ensure that transaction indexes are monotonically increasing.
	index.Value = oldValue + 1
	nextTransactionIndexTransientStore.Set(
		[]byte(types.NextStatefulOrderBlockTransactionIndexKey),
		k.cdc.MustMarshal(&index),
	)
	return oldValue
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

	untriggeredConditionalOrderMemStore := k.GetUntriggeredConditionalOrderPlacementMemStore(ctx)
	untriggeredConditionalOrderStore := k.GetUntriggeredConditionalOrderPlacementStore(ctx)

	bytes := untriggeredConditionalOrderMemStore.Get(orderId.ToStateKey())
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
	orderKey := orderId.ToStateKey()
	triggeredConditionalOrderStore.Set(orderKey, longTermOrderPlacementBytes)
	triggeredConditionalOrderMemStore.Set(orderKey, longTermOrderPlacementBytes)

	// Delete the `StatefulOrderPlacement` from Untriggered state store/memstore.
	untriggeredConditionalOrderStore.Delete(orderKey)
	untriggeredConditionalOrderMemStore.Delete(orderKey)

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
		store.Delete(sdk.FormatTimeBytes(goodTilBlockTime))
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
	triggeredMemstore := k.GetTriggeredConditionalOrderPlacementMemStore(ctx)
	orderKey := orderId.ToStateKey()
	return triggeredMemstore.Has(orderKey)
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

// GetAllStatefulOrders iterates over all stateful order placements and returns a list
// of orders, ordered by ascending time priority.
func (k Keeper) GetAllStatefulOrders(ctx sdk.Context) []types.Order {
	return k.getStatefulOrders(k.getAllOrdersIterator(ctx))
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
		sdk.FormatTimeBytes(goodTilBlockTime),
		b,
	)
}

// getStatefulOrdersTimeSliceIterator returns an iterator over all stateful order time slice values
// from time 0 until `endTime`.
func (k Keeper) getStatefulOrdersTimeSliceIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	startKey := []byte(types.StatefulOrdersTimeSlicePrefix)
	endKey := append(
		startKey,
		sdk.InclusiveEndBytes(
			sdk.FormatTimeBytes(endTime),
		)...,
	)
	return store.Iterator(
		startKey,
		endKey,
	)
}

// getAllOrdersIterator returns an iterator over all stateful orders, which includes all
// Long-Term orders, triggered and untriggered conditional orders.
func (k Keeper) getAllOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.StatefulOrderKeyPrefix),
	)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// getPlacedOrdersIterator returns an iterator over all placed orders, which includes all
// Long-Term orders and triggered conditional orders.
func (k Keeper) getPlacedOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.PlacedStatefulOrderKeyPrefix),
	)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// getUntriggeredConditionalOrdersIterator returns an iterator over all untriggered conditional
// orders.
func (k Keeper) getUntriggeredConditionalOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.UntriggeredConditionalOrderKeyPrefix),
	)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// GetStatefulOrderCount gets a count of how many stateful orders are written to state for a subaccount.
func (k Keeper) GetStatefulOrderCount(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
) uint32 {
	store := k.GetStatefulOrderCountMemStore(ctx)

	b := store.Get(subaccountId.ToStateKey())
	result := gogotypes.UInt32Value{Value: 0}
	if b != nil {
		k.cdc.MustUnmarshal(b, &result)
	}
	return result.Value
}

// SetStatefulOrderCount sets a count of how many stateful orders are written to state.
func (k Keeper) SetStatefulOrderCount(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	count uint32,
) {
	store := k.GetStatefulOrderCountMemStore(ctx)

	if count == 0 {
		store.Delete(subaccountId.ToStateKey())
	} else {
		result := gogotypes.UInt32Value{Value: count}
		store.Set(
			subaccountId.ToStateKey(),
			k.cdc.MustMarshal(&result),
		)
	}
}
