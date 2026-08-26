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

// conditionalOrderTriggerDirectionForOrder is an alias kept for clarity at delete call sites, where
// we only have the placement record's Order. The caller must pass the fully-loaded Order.
func conditionalOrderTriggerDirectionForOrder(order types.Order) byte {
	return conditionalOrderTriggerDirection(order)
}

// triggerIndexKeySequenceLen is the width of the placement-sequence tie-breaker embedded in the
// trigger-price index key (packed block height + transaction index; see triggerIndexSequenceKey).
const triggerIndexKeySequenceLen = 8

// triggerIndexSequenceKey encodes an order's placement ordering (block height, transaction index)
// into an 8-byte tie-breaker that follows triggerSubticks in the index key, so that among orders
// with EQUAL trigger subticks the OLDER order (earlier placement) is triggered first under a
// per-block budget. This replaces the previous tie-break on the raw orderId, which embeds the
// client-chosen ClientId and thus let a later order jump an older equal-price order. Placement
// ordering is protocol-assigned (GetNextStatefulOrderTransactionIndex) and cannot be influenced
// by the client.
//
// The value is direction-aware so that oldest-first holds under BOTH iteration directions used by
// IterateCrossedConditionalOrders:
//   - LTE is scanned forward, so the raw (ascending) sequence already yields oldest-first.
//   - GTE is scanned in reverse (nearest-crossing = highest subticks first), so the bitwise
//     complement is stored; reversing a descending complement yields an ascending real sequence —
//     still oldest-first.
func triggerIndexSequenceKey(direction byte, placement types.TransactionOrdering) uint64 {
	rawSeq := (uint64(placement.BlockHeight) << 32) | uint64(placement.TransactionIndex)
	if direction == TriggerDirectionGTE {
		return ^rawSeq
	}
	return rawSeq
}

// encodeTriggerPriceIndexKey builds the composite key for the trigger-price secondary index.
// The key layout (all within the ConditionalOrderTriggerPriceIndexKeyPrefix prefix store) is:
//
//	<clobPairId:4 big-endian> <directionByte:1> <triggerSubticks:8 big-endian>
//	  <sequenceKey:8 big-endian> <orderId state key>
//
// The big-endian encoding of clobPairId, triggerSubticks, and sequenceKey ensures that a forward
// byte scan over the prefix store is equivalent to ascending numeric order, enabling efficient
// range queries. The sequenceKey (see triggerIndexSequenceKey) sits between triggerSubticks and the
// orderId so that orders with equal subticks are ordered by placement time, not by client-chosen
// orderId bytes.
func encodeTriggerPriceIndexKey(
	clobPairId uint32,
	direction byte,
	triggerSubticks uint64,
	sequenceKey uint64,
	orderId types.OrderId,
) []byte {
	orderKey := orderId.ToStateKey()
	// 4 clobPairId + 1 direction + 8 triggerSubticks + 8 sequenceKey + len(orderKey)
	key := make([]byte, 4+1+8+triggerIndexKeySequenceLen+len(orderKey))
	binary.BigEndian.PutUint32(key[0:4], clobPairId)
	key[4] = direction
	binary.BigEndian.PutUint64(key[5:13], triggerSubticks)
	binary.BigEndian.PutUint64(key[13:21], sequenceKey)
	copy(key[21:], orderKey)
	return key
}

// GetConditionalOrderTriggerPriceIndexStore returns the prefix store for the trigger-price index.
func (k Keeper) GetConditionalOrderTriggerPriceIndexStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(
		ctx.KVStore(k.storeKey),
		[]byte(types.ConditionalOrderTriggerPriceIndexKeyPrefix),
	)
}

// clobPairHasTriggerIndexEntries reports, in O(1), whether the trigger-price index holds any
// untriggered conditional order for the given clobPairId. The bounded EndBlocker trigger path uses
// it to skip pairs with nothing to trigger — visiting only pairs that have resting orders and
// avoiding an oracle-price fetch (and its zero-price panic) for empty pairs.
func (k Keeper) clobPairHasTriggerIndexEntries(ctx sdk.Context, clobPairId uint32) bool {
	store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
	prefixKey := make([]byte, 4)
	binary.BigEndian.PutUint32(prefixKey, clobPairId)
	it := store.Iterator(prefixKey, storetypes.PrefixEndBytes(prefixKey))
	defer it.Close()
	return it.Valid()
}

// addConditionalOrderToTriggerPriceIndex inserts an entry into the trigger-price secondary index
// for the given untriggered conditional order.  The value stored is empty (the key itself
// encodes all necessary lookup information).
//
// This must be called whenever a conditional order is placed into the untriggered store.
func (k Keeper) addConditionalOrderToTriggerPriceIndex(
	ctx sdk.Context,
	order types.Order,
	placementIndex types.TransactionOrdering,
) {
	direction := conditionalOrderTriggerDirection(order)
	key := encodeTriggerPriceIndexKey(
		uint32(order.GetClobPairId()),
		direction,
		order.ConditionalOrderTriggerSubticks,
		triggerIndexSequenceKey(direction, placementIndex),
		order.OrderId,
	)
	store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
	store.Set(key, []byte{})
}

// removeConditionalOrderFromTriggerPriceIndex removes the entry for the given order from the
// trigger-price secondary index.  This is a no-op if the key does not exist (safe for all
// delete/cancel/trigger/expire call sites).
func (k Keeper) removeConditionalOrderFromTriggerPriceIndex(
	ctx sdk.Context,
	order types.Order,
	placementIndex types.TransactionOrdering,
) {
	direction := conditionalOrderTriggerDirectionForOrder(order)
	key := encodeTriggerPriceIndexKey(
		uint32(order.GetClobPairId()),
		direction,
		order.ConditionalOrderTriggerSubticks,
		triggerIndexSequenceKey(direction, placementIndex),
		order.OrderId,
	)
	store := k.GetConditionalOrderTriggerPriceIndexStore(ctx)
	store.Delete(key)
}

// BuildConditionalOrderTriggerPriceIndex builds the trigger-price secondary index from the full
// set of resting untriggered conditional orders. It is a one-shot migration run by the
// state-breaking upgrade handler at a coordinated height, so it may perform an unbounded scan: the
// cost is paid once during the upgrade block, never per block. After this runs the index is
// authoritative and is maintained incrementally by the placement / cancel / trigger / expiry hooks.
func (k Keeper) BuildConditionalOrderTriggerPriceIndex(ctx sdk.Context) {
	source := k.GetUntriggeredConditionalOrderPlacementStore(ctx)
	it := source.Iterator(nil, nil)
	defer it.Close()
	indexed := 0
	for ; it.Valid(); it.Next() {
		var placement types.LongTermOrderPlacement
		k.cdc.MustUnmarshal(it.Value(), &placement)
		k.addConditionalOrderToTriggerPriceIndex(ctx, placement.Order, placement.PlacementIndex)
		indexed++
	}
	ctx.Logger().Info("conditional-order trigger-price index built", "indexed", indexed)
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
	k.iterateCrossedConditionalOrders(
		ctx,
		clobPairId,
		direction,
		priceSubticks,
		func(orderId types.OrderId, _ []byte) bool {
			return fn(orderId)
		},
	)
}

// iterateCrossedConditionalOrders is the internal form of
// IterateCrossedConditionalOrders. It also exposes the raw secondary-index key so callers can
// remove stale entries without reconstructing placement metadata that may no longer exist.
func (k Keeper) iterateCrossedConditionalOrders(
	ctx sdk.Context,
	clobPairId uint32,
	direction byte,
	priceSubticks uint64,
	fn func(orderId types.OrderId, indexKey []byte) bool,
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

	// Visit nearest-crossing orders first so that, under a per-block trigger budget, the orders
	// closest to the current price are triggered before farther ones:
	//   - LTE crosses at triggerSubticks ≥ price; nearest = lowest such subticks = ASCENDING from
	//     priceSubticks, i.e. a forward iterator.
	//   - GTE crosses at triggerSubticks ≤ price; nearest = highest such subticks = DESCENDING from
	//     priceSubticks, i.e. a reverse iterator over the same [0, price] range.
	// Both directions produce a deterministic, node-identical order.
	var it storetypes.Iterator
	if direction == TriggerDirectionGTE {
		it = store.ReverseIterator(startKey, endKey)
	} else {
		it = store.Iterator(startKey, endKey)
	}
	defer it.Close()

	for ; it.Valid(); it.Next() {
		// Key is relative to the ConditionalOrderTriggerPriceIndexKeyPrefix prefix store.
		// Layout: <clobPairId:4><dir:1><subticks:8><sequenceKey:8><orderId:N>
		rawKey := it.Key()
		if len(rawKey) <= 4+1+8+triggerIndexKeySequenceLen {
			// Malformed key — skip.
			continue
		}
		orderKeyBytes := rawKey[4+1+8+triggerIndexKeySequenceLen:]

		var orderId types.OrderId
		k.cdc.MustUnmarshal(orderKeyBytes, &orderId)

		if !fn(orderId, rawKey) {
			break
		}
	}
}

// removeConditionalOrderTriggerPriceIndexEntry removes a raw key yielded by
// iterateCrossedConditionalOrders. This is used to clean up expired or orphaned index entries
// while keeping consensus-state expiry deletion in the bounded expiry queue.
func (k Keeper) removeConditionalOrderTriggerPriceIndexEntry(ctx sdk.Context, indexKey []byte) {
	k.GetConditionalOrderTriggerPriceIndexStore(ctx).Delete(indexKey)
}
