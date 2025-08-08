package keeper

import (
	"errors"
	"fmt"
	"sort"
	"time"

	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/store/prefix"
	cometbftlog "github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
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
	store := k.fetchStateStoresForOrder(ctx, order.OrderId)
	orderKey := order.OrderId.ToStateKey()

	// Write the `LongTermOrderPlacement` to state.
	store.Set(orderKey, longTermOrderPlacementBytes)

	if !found {
		// Increment the stateful order count.
		k.CheckAndIncrementStatefulOrderCount(ctx, order.OrderId)

		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.StatefulOrder, metrics.Count},
			1,
			[]metrics.Label{
				metrics.GetLabelForIntValue(metrics.ClobPairId, int(order.GetClobPairId())),
				metrics.GetLabelForBoolValue(metrics.Conditional, order.OrderId.IsConditionalOrder()),
			},
		)
	}
}

// GetTriggeredConditionalOrderPlacement gets an triggered conditional order placement from the store.
// Returns false if no triggered conditional order exists in store with `orderId`.
func (k Keeper) GetTriggeredConditionalOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.LongTermOrderPlacement, found bool) {
	store := k.GetTriggeredConditionalOrderPlacementStore(ctx)

	b := store.Get(orderId.ToStateKey())
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetUntriggeredConditionalOrderPlacement gets an untriggered conditional order placement from the store.
// Returns false if no untriggered conditional order exists in store with `orderId`.
func (k Keeper) GetUntriggeredConditionalOrderPlacement(
	ctx sdk.Context,
	orderId types.OrderId,
) (val types.LongTermOrderPlacement, found bool) {
	store := k.GetUntriggeredConditionalOrderPlacementStore(ctx)

	b := store.Get(orderId.ToStateKey())
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

	store := k.fetchStateStoresForOrder(ctx, orderId)

	b := store.Get(orderId.ToStateKey())
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

	store := k.fetchStateStoresForOrder(ctx, orderId)
	orderKey := orderId.ToStateKey()
	orderExists := store.Has(orderKey)

	// Delete the `StatefulOrderPlacement` from state.
	store.Delete(orderKey)

	// Set the count.
	if orderExists {
		k.CheckAndDecrementStatefulOrderCount(ctx, orderId)
	}

	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.StatefulOrderRemoved, metrics.Count},
		1,
		[]metrics.Label{
			metrics.GetLabelForIntValue(metrics.ClobPairId, int(orderId.GetClobPairId())),
			metrics.GetLabelForBoolValue(metrics.Conditional, orderId.IsConditionalOrder()),
		},
	)
}

// GetStatefulOrderIdExpirations gets a slice of stateful order IDs that expire at `goodTilBlockTime`,
// sorted by order ID state key.
func (k Keeper) GetStatefulOrderIdExpirations(ctx sdk.Context, goodTilBlockTime time.Time) (
	orderIds []types.OrderId,
) {
	return k.GetOrderIds(
		ctx,
		ctx.KVStore(k.storeKey),
		fmt.Sprintf(types.StatefulOrdersExpirationsKeyPrefix, sdk.FormatTimeString(goodTilBlockTime)),
	)
}

// RemoveStatefulOrderIdExpiration removes a stateful order id expiration for an order id and
// `goodTilBlockTime`.
func (k Keeper) RemoveStatefulOrderIdExpiration(
	ctx sdk.Context,
	goodTilBlockTime time.Time,
	orderId types.OrderId,
) {
	k.RemoveUnorderedOrderId(
		ctx,
		ctx.KVStore(k.storeKey),
		fmt.Sprintf(types.StatefulOrdersExpirationsKeyPrefix, sdk.FormatTimeString(goodTilBlockTime)),
		orderId,
	)
}

// AddStatefulOrderIdExpiration adds a stateful order id expiration for an order id and
// `goodTilBlockTime`.
func (k Keeper) AddStatefulOrderIdExpiration(
	ctx sdk.Context,
	goodTilBlockTime time.Time,
	orderId types.OrderId,
) {
	k.SetUnorderedOrderId(
		ctx,
		ctx.KVStore(k.storeKey),
		fmt.Sprintf(types.StatefulOrdersExpirationsKeyPrefix, sdk.FormatTimeString(goodTilBlockTime)),
		orderId,
	)
}

// RemoveExpiredStatefulOrders removes the stateful order id expirations up to `blockTime` and
// returns the removed order ids as a slice.
func (k Keeper) RemoveExpiredStatefulOrders(ctx sdk.Context, blockTime time.Time) (
	expiredOrderIds []types.OrderId,
) {
	expiredOrderIds = make([]types.OrderId, 0)
	store := ctx.KVStore(k.storeKey)
	it := store.Iterator(
		[]byte(fmt.Sprintf(types.StatefulOrdersExpirationsKeyPrefix, sdk.FormatTimeString(time.Time{}))),
		storetypes.PrefixEndBytes(
			[]byte(fmt.Sprintf(types.StatefulOrdersExpirationsKeyPrefix, sdk.FormatTimeString(blockTime))),
		),
	)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		var orderId types.OrderId
		k.cdc.MustUnmarshal(it.Value(), &orderId)
		expiredOrderIds = append(expiredOrderIds, orderId)
		store.Delete(it.Key())
	}
	return expiredOrderIds
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

	untriggeredConditionalOrderStore := k.GetUntriggeredConditionalOrderPlacementStore(ctx)

	bytes := untriggeredConditionalOrderStore.Get(orderId.ToStateKey())
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

	// Write the StatefulOrderPlacement to the Triggered state store.
	longTermOrderPlacementBytes := k.cdc.MustMarshal(&longTermOrderPlacement)
	triggeredConditionalOrderStore := k.GetTriggeredConditionalOrderPlacementStore(ctx)
	orderKey := orderId.ToStateKey()
	triggeredConditionalOrderStore.Set(orderKey, longTermOrderPlacementBytes)

	// Delete the `StatefulOrderPlacement` from Untriggered state store.
	untriggeredConditionalOrderStore.Delete(orderKey)
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

	order, exists := k.getOrderFromStore(ctx, orderId)
	if !exists {
		panic(fmt.Sprintf("MustRemoveStatefulOrder: order %v does not exist", orderId))
	}

	// TWAP orders are not maintained by these states because expiry and
	// fills are updated by the generated suborders.
	if !order.OrderId.IsTwapOrder() {
		goodTilBlockTime := order.MustGetUnixGoodTilBlockTime()
		k.RemoveStatefulOrderIdExpiration(ctx, goodTilBlockTime, orderId)
		// Remove the order fill amount from state.
		k.RemoveOrderFillAmount(ctx, orderId)
		// Delete the Stateful order placement from state.
		k.DeleteLongTermOrderPlacement(ctx, orderId)
	} else {
		k.DeleteTWAPOrderPlacement(ctx, orderId)
		// Cancelling a TWAP parent order will attempt to cancel the in-flight suborder.
		// Since a TWAP suborder is maintained as a normal stateful order, cancelling a
		// suborder follows the same flow as other stateful orders.
		suborderId := k.twapToSuborderId(order.OrderId)
		// GoodTilBlockTime is set to the current block time + 10 seconds.
		// This is because an in-flight suborder will have a good-til-block time of
		// the current block time + 3 seconds, so 10 seconds is a reasonable buffer.
		err := k.HandleMsgCancelOrder(ctx, &types.MsgCancelOrder{
			OrderId: suborderId,
			GoodTilOneof: &types.MsgCancelOrder_GoodTilBlockTime{
				GoodTilBlockTime: uint32(ctx.BlockTime().Unix() + 10),
			},
		})
		// We can expect a suborder to not exist during the window while it is waiting
		// to be triggered.
		if err != nil && !errors.Is(err, types.ErrStatefulOrderDoesNotExist) {
			panic(fmt.Sprintf("MustRemoveStatefulOrder: error cancelling twap order and suborder: %v", err))
		}
	}
}

// IsConditionalOrderTriggered checks if a given order ID is triggered or untriggered in state.
// Note: If the given order ID is neither in triggered or untriggered state, function will return false.
func (k Keeper) IsConditionalOrderTriggered(
	ctx sdk.Context,
	orderId types.OrderId,
) (triggered bool) {
	// If this is not a conditional order, panic.
	orderId.MustBeConditionalOrder()
	store := k.GetTriggeredConditionalOrderPlacementStore(ctx)
	orderKey := orderId.ToStateKey()
	return store.Has(orderKey)
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
func (k Keeper) getStatefulOrders(statefulOrderIterator dbm.Iterator) []types.Order {
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

// getAllOrdersIterator returns an iterator over all stateful orders, which includes all
// Long-Term orders, triggered and untriggered conditional orders.
func (k Keeper) getAllOrdersIterator(ctx sdk.Context) storetypes.Iterator {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.StatefulOrderKeyPrefix),
	)
	return storetypes.KVStorePrefixIterator(store, []byte{})
}

// getPlacedOrdersIterator returns an iterator over all placed orders, which includes all
// Long-Term orders and triggered conditional orders.
func (k Keeper) getPlacedOrdersIterator(ctx sdk.Context) storetypes.Iterator {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.PlacedStatefulOrderKeyPrefix),
	)
	return storetypes.KVStorePrefixIterator(store, []byte{})
}

// getUntriggeredConditionalOrdersIterator returns an iterator over all untriggered conditional
// orders.
func (k Keeper) getUntriggeredConditionalOrdersIterator(ctx sdk.Context) storetypes.Iterator {
	store := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.UntriggeredConditionalOrderKeyPrefix),
	)
	return storetypes.KVStorePrefixIterator(store, []byte{})
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

// Increment the stateful order count only once for TWAP orders.
// Suborders that are generated from a TWAP order are not counted
// towards the stateful order count because technically they are
// part of one parent order.
func (k Keeper) CheckAndIncrementStatefulOrderCount(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	if orderId.IsTwapSuborder() {
		return
	}
	subaccountId := orderId.SubaccountId
	count := k.GetStatefulOrderCount(ctx, subaccountId)
	k.SetStatefulOrderCount(ctx, subaccountId, count+1)
}

func (k Keeper) CheckAndDecrementStatefulOrderCount(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	if orderId.IsTwapSuborder() {
		return
	}

	subaccountId := orderId.SubaccountId
	count := k.GetStatefulOrderCount(ctx, subaccountId)
	if count == 0 {
		log.ErrorLog(ctx, "Stateful order count is zero but order is in the store. Underflow",
			"orderId", cometbftlog.NewLazySprintf("%+v", orderId),
		)
	} else {
		k.SetStatefulOrderCount(ctx, subaccountId, count-1)
	}
}
