package keeper

import (
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// TriggerDirectionLTE is the direction byte for orders that trigger when the oracle price is
// less than or equal to the order's trigger price (take-profit buys + stop-loss sells).
const TriggerDirectionLTE byte = 0x00

// TriggerDirectionGTE is the direction byte for orders that trigger when the oracle price is
// greater than or equal to the order's trigger price (take-profit sells + stop-loss buys).
const TriggerDirectionGTE byte = 0x01

// conditionalOrderTriggerDirection returns the direction byte for the given order by replicating
// the exact classification used by AddUntriggeredConditionalOrder / CanTrigger:
//   - LTE (0x00): take-profit BUY or stop-loss SELL (triggered when oracle price ≤ trigger price)
//   - GTE (0x01): take-profit SELL or stop-loss BUY  (triggered when oracle price ≥ trigger price)
func conditionalOrderTriggerDirection(order types.Order) byte {
	if order.IsTakeProfitOrder() && order.IsBuy() ||
		order.IsStopLossOrder() && !order.IsBuy() {
		return TriggerDirectionLTE
	}
	return TriggerDirectionGTE
}

// conditionalOrderTriggerDirectionForId returns the direction byte given an orderId and the
// order's trigger metadata. This is the inverse-read helper used during delete (when we only have
// the placement record, not the full order struct).  The caller must pass the fully-loaded Order.
// This is just an alias kept for clarity at call sites.
func conditionalOrderTriggerDirectionForOrder(order types.Order) byte {
	return conditionalOrderTriggerDirection(order)
}

// encodeTriggerPriceIndexKey builds the composite key for the trigger-price secondary index.
// The key layout (all within the ConditionalOrderTriggerPriceIndexKeyPrefix prefix store) is:
//
//	<clobPairId:4 big-endian> <directionByte:1> <triggerSubticks:8 big-endian> <orderId state key>
//
// The big-endian encoding of both clobPairId and triggerSubticks ensures that a forward byte scan
// over the prefix store is equivalent to ascending numeric order, enabling efficient range queries.
func encodeTriggerPriceIndexKey(
	clobPairId uint32,
	direction byte,
	triggerSubticks uint64,
	orderId types.OrderId,
) []byte {
	orderKey := orderId.ToStateKey()
	// 4 bytes clobPairId + 1 byte direction + 8 bytes triggerSubticks + len(orderKey)
	key := make([]byte, 4+1+8+len(orderKey))
	binary.BigEndian.PutUint32(key[0:4], clobPairId)
	key[4] = direction
	binary.BigEndian.PutUint64(key[5:13], triggerSubticks)
	copy(key[13:], orderKey)
	return key
}

// GetConditionalOrderTriggerPriceIndexStore returns the prefix store for the trigger-price index.
func (k Keeper) GetConditionalOrderTriggerPriceIndexStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.ConditionalOrderTriggerPriceIndexKeyPrefix),
	)
}

// addConditionalOrderToTriggerPriceIndex inserts an entry into the trigger-price secondary index
// for the given untriggered conditional order.  The value stored is empty (the key itself
// encodes all necessary lookup information).
//
// This must be called whenever a conditional order is placed into the untriggered store.
func (k Keeper) addConditionalOrderToTriggerPriceIndex(ctx sdk.Context, order types.Order) {
	direction := conditionalOrderTriggerDirection(order)
	key := encodeTriggerPriceIndexKey(
		uint32(order.GetClobPairId()),
		direction,
		order.ConditionalOrderTriggerSubticks,
		order.OrderId,
	)
	store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
	store.Set(key, []byte{})
}

// removeConditionalOrderFromTriggerPriceIndex removes the entry for the given order from the
// trigger-price secondary index.  This is a no-op if the key does not exist (safe for all
// delete/cancel/trigger/expire call sites).
func (k Keeper) removeConditionalOrderFromTriggerPriceIndex(ctx sdk.Context, order types.Order) {
	direction := conditionalOrderTriggerDirectionForOrder(order)
	key := encodeTriggerPriceIndexKey(
		uint32(order.GetClobPairId()),
		direction,
		order.ConditionalOrderTriggerSubticks,
		order.OrderId,
	)
	store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
	store.Delete(key)
}

// BackfillConditionalOrderTriggerPriceIndex populates the trigger-price secondary index for every
// untriggered conditional order that does not yet have an entry.  It is called exactly once at the
// governance flag-enable transition (OFF->ON) to backfill orders that were resting before the
// config was enabled and were therefore never written into the index by the placement hooks.
//
// DETERMINISM: GetAllUntriggeredConditionalOrders returns orders sorted by ascending time
// priority (SortedLongTermOrderPlacements), so the iteration order is identical on every node.
// addConditionalOrderToTriggerPriceIndex is a pure KV store.Set — it is idempotent and safe to
// call even when an entry already exists.
//
// This must NOT be called again after the upgrade block because the per-placement hooks
// (SetLongTermOrderPlacement) maintain the index from that point forward.  InitMemStore
// rehydrates NumUcond / NumUcondSa: counters independently from the untriggered store, so no
// double-counting occurs.
func (k Keeper) BackfillConditionalOrderTriggerPriceIndex(ctx sdk.Context) {
	// Reconcile: clear any existing index entries, then repopulate from the authoritative
	// untriggered conditional order store. Deterministic (fixed iteration order) and safe to run
	// at the mitigation's enable transition regardless of prior index state (e.g. after a
	// disabled window during which the per-placement hooks did not maintain the index).
	indexStore := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
	staleKeys := make([][]byte, 0)
	it := indexStore.Iterator(nil, nil)
	for ; it.Valid(); it.Next() {
		staleKeys = append(staleKeys, it.Key())
	}
	it.Close()
	for _, key := range staleKeys {
		indexStore.Delete(key)
	}

	orders := k.GetAllUntriggeredConditionalOrders(ctx)
	for _, order := range orders {
		k.addConditionalOrderToTriggerPriceIndex(ctx, order)
	}
	ctx.Logger().Info(
		"BackfillConditionalOrderTriggerPriceIndex: reconciled trigger-price index",
		"cleared", len(staleKeys),
		"count", len(orders),
	)
}

// IterateCrossedConditionalOrders iterates over orders in the trigger-price index that are
// crossed by the given price for a specific clobPairId and direction.
//
//   - LTE direction (TriggerDirectionLTE): an order is crossed when
//     oracle_price ≤ order.triggerSubticks, i.e. all keys with triggerSubticks ≥ priceSubticks.
//
//   - GTE direction (TriggerDirectionGTE): an order is crossed when
//     oracle_price ≥ order.triggerSubticks, i.e. all keys with triggerSubticks ≤ priceSubticks.
//
// The callback fn receives each crossed orderId.  Iteration stops early if fn returns false.
// priceSubticks must already be the pessimistic (rounded) value, matching PollTriggeredConditionalOrders.
func (k Keeper) IterateCrossedConditionalOrders(
	ctx sdk.Context,
	clobPairId uint32,
	direction byte,
	priceSubticks uint64,
	fn func(orderId types.OrderId) bool,
) {
	store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)

	// Key layout in the prefix store: <clobPairId:4><dir:1><subticks:8><orderId:N>
	// Big-endian encoding ensures lexicographic order == numeric order for all fixed-width fields.

	var startKey, endKey []byte

	switch direction {
	case TriggerDirectionLTE:
		// Crossed orders: triggerSubticks ∈ [priceSubticks, MaxUint64].
		// Start: first key at or after (clobPairId, LTE, priceSubticks).
		startKey = make([]byte, 4+1+8)
		binary.BigEndian.PutUint32(startKey[0:4], clobPairId)
		startKey[4] = TriggerDirectionLTE
		binary.BigEndian.PutUint64(startKey[5:13], priceSubticks)
		// End: first key in GTE bucket for same clobPairId (exclusive).
		endKey = make([]byte, 4+1)
		binary.BigEndian.PutUint32(endKey[0:4], clobPairId)
		endKey[4] = TriggerDirectionGTE

	case TriggerDirectionGTE:
		// Crossed orders: triggerSubticks ∈ [0, priceSubticks].
		// Start: first key in GTE bucket for this clobPairId.
		startKey = make([]byte, 4+1)
		binary.BigEndian.PutUint32(startKey[0:4], clobPairId)
		startKey[4] = TriggerDirectionGTE
		// End: first key after (clobPairId, GTE, priceSubticks) — exclusive.
		if priceSubticks < ^uint64(0) {
			endKey = make([]byte, 4+1+8)
			binary.BigEndian.PutUint32(endKey[0:4], clobPairId)
			endKey[4] = TriggerDirectionGTE
			binary.BigEndian.PutUint64(endKey[5:13], priceSubticks+1)
		} else {
			// MaxUint64: include all GTE keys; end at the next clobPairId + LTE bucket.
			// Use PrefixEndBytes on the GTE prefix to get the exclusive upper bound.
			endBase := make([]byte, 4+1)
			binary.BigEndian.PutUint32(endBase[0:4], clobPairId)
			endBase[4] = TriggerDirectionGTE
			endKey = storetypes.PrefixEndBytes(endBase)
		}
	}

	it := store.Iterator(startKey, endKey)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		// Key is relative to the ConditionalOrderTriggerPriceIndexKeyPrefix prefix store.
		// Layout: <clobPairId:4><dir:1><subticks:8><orderId:N>
		rawKey := it.Key()
		if len(rawKey) <= 4+1+8 {
			// Malformed key — skip.
			continue
		}
		orderKeyBytes := rawKey[4+1+8:]

		var orderId types.OrderId
		k.cdc.MustUnmarshal(orderKeyBytes, &orderId)

		if !fn(orderId) {
			break
		}
	}
}
