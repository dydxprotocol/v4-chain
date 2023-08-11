package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/clob/types"
)

// AddSeenPlaceOrder adds a PlaceOrder to the in-memory store. Currently only supports short-term
// orders with a GoodTilBlock.
func (k Keeper) AddSeenPlaceOrder(
	ctx sdk.Context,
	placeOrder types.MsgPlaceOrder,
) {
	// Get relevant fields.
	orderId := placeOrder.Order.OrderId
	var goodTilBlock types.Order_GoodTilBlock

	// Only Short Term Orders (with GoodTilBlock) are currently supported.
	// TODO(DEC-1207): Support Long-term orders.
	switch goodTil := placeOrder.Order.GoodTilOneof.(type) {
	case *types.Order_GoodTilBlock:
		goodTilBlock = types.Order_GoodTilBlock{GoodTilBlock: goodTil.GoodTilBlock}
	default:
		return
	}

	// Log an info and return early if the GoodTilBlock is lower than the current one.
	// TODO(DEC-1925): Investigate if this error is valid, and convert back to error log if so.
	curGoodTilBlock, ok := k.seenPlaceOrderIds[orderId]
	if ok && goodTilBlock.GoodTilBlock < curGoodTilBlock.GoodTilBlock {
		k.Logger(ctx).Info(fmt.Sprintf(
			"AddSeenPlaceOrder: new seen orderId %v was found with new GoodTilBlock (%v) "+
				"which is lower than the current one (%v)",
			orderId,
			goodTilBlock.GoodTilBlock,
			curGoodTilBlock.GoodTilBlock,
		))
		return
	}

	k.seenPlaceOrderIds[orderId] = goodTilBlock

	if _, ok := k.seenGoodTilBlocks[goodTilBlock]; !ok {
		k.seenGoodTilBlocks[goodTilBlock] = make(map[types.OrderId]bool)
	}
	k.seenGoodTilBlocks[goodTilBlock][orderId] = true
}

// Prune place orders from the data structure that have expired (e.g. GoodTilBlock is older than
// currentBlockHeight - ShortBlockWindow). Currently only supports short-term orders with a GoodTilBlock.
func (k Keeper) PruneExpiredSeenPlaceOrders(
	ctx sdk.Context,
	goodTilBlockHeightToPrune uint32,
) {
	// TODO(DEC-1207): Support Long-term orders.
	goodTilBlockToPrune := types.Order_GoodTilBlock{GoodTilBlock: goodTilBlockHeightToPrune}

	// Get seenPlaceOrderIds that will be pruned.
	seenPlaceOrderIdsToPotentiallyPrune, ok := k.seenGoodTilBlocks[goodTilBlockToPrune]
	if !ok {
		// Nothing to prune.
		return
	}

	// Prune all seenPlaceOrders (unless they have an updated goodTilBlock)
	for orderId := range seenPlaceOrderIdsToPotentiallyPrune {
		orderGoodTilBlock, ok := k.seenPlaceOrderIds[orderId]
		if !ok {
			k.Logger(ctx).Error(fmt.Sprintf(
				"PruneExpiredSeenPlaceOrders: orderId %v was found in seenGoodTilBlocks "+
					"but not in seenPlaceOrderIds for block height %d",
				orderId,
				goodTilBlockHeightToPrune,
			))
		}

		// order's latest goodTilBlock matches the block we are pruning.
		if orderGoodTilBlock == goodTilBlockToPrune {
			delete(k.seenPlaceOrderIds, orderId)
		}
	}

	// Prune seenGoodTilBlock
	delete(k.seenGoodTilBlocks, goodTilBlockToPrune)
}

// HasSeenPlaceOrder checks if a MsgPlaceOrder has been seen with the same (or higher) GoodTilBlock.
func (k Keeper) HasSeenPlaceOrder(
	ctx sdk.Context,
	placeOrder types.MsgPlaceOrder,
) bool {
	// Get relevant fields.
	orderId := placeOrder.Order.OrderId
	var goodTilBlock types.Order_GoodTilBlock

	// Only Short Term Orders (with GoodTilBlock) are currently supported.
	// TODO(DEC-1207): Support Long-term orders.
	switch goodTil := placeOrder.Order.GoodTilOneof.(type) {
	case *types.Order_GoodTilBlock:
		goodTilBlock = types.Order_GoodTilBlock{GoodTilBlock: goodTil.GoodTilBlock}
	case *types.Order_GoodTilBlockTime:
		// Since long-term orders are not supported, mark them as seen to not produce false negatives.
		return true
	default:
		// Invalid place orders (e.g. no GoodTilOneof) are never seen.
		return false
	}

	seenLatestGoodTilBlock, ok := k.seenPlaceOrderIds[orderId]
	return ok && seenLatestGoodTilBlock.GoodTilBlock >= goodTilBlock.GoodTilBlock
}
