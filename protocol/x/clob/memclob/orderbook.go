package memclob

import (
	"fmt"
	"math"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/pkg/errors"
	"github.com/zyedidia/generic/list"
)

// Orderbook holds the bids and asks for a specific `ClobPairId`.
type Orderbook struct {
	// Defines the tick size of the orderbook by defining how many subticks
	// are in one tick. That is, the subticks of any valid order must be a
	// multiple of this value. Generally this value should start `>= 100` to
	// allow room for decreasing it. This field is stored in state as part of a
	// `ClobPair`, but must be made available to the in-memory `Orderbook` in
	// order to efficiently remove orders from the orderbook. See the `removeOrder`
	// implementation for more information.
	SubticksPerTick types.SubticksPerTick
	// Map of price level (in subticks) to buy orders contained at that level.
	Bids map[types.Subticks]*types.Level
	// Map of price level (in subticks) to sell orders contained at that level.
	Asks map[types.Subticks]*types.Level
	// The highest bid on this orderbook, in subticks. 0 if no bids exist.
	BestBid types.Subticks
	// The lowest ask on this orderbook, in subticks. math.MaxUint64 if no asks exist.
	BestAsk types.Subticks
	// Contains all open orders on this CLOB for a given subaccount and side.
	// Used for fetching open orders for the add to orderbook collateralization
	// check for a subaccount.
	SubaccountOpenClobOrders map[satypes.SubaccountId]map[types.Order_Side]map[types.OrderId]bool
	// Minimum size of an order on the CLOB, in base quantums.
	MinOrderBaseQuantums satypes.BaseQuantums
	// Contains all open reduce-only orders on this CLOB from each subaccount. Used for tracking
	// which open reduce-only orders should be canceled when a position changes sides.
	SubaccountOpenReduceOnlyOrders map[satypes.SubaccountId]map[types.OrderId]bool
	// TotalOpenOrders tracks the total number of open orders in an orderbook for observability purposes.
	TotalOpenOrders uint
	// Map from order IDs to their `LevelOrder` reference contained in the
	// orderbook, necessary for O(1) order removal from the orderbook.
	orderIdToLevelOrder map[types.OrderId]*types.LevelOrder
	// Map from block number to a set of all orders that expire at this block
	// (with each order keyed by `OrderId`). Necessary for O(1) order removal
	// from the orderbook when expiring orders in the EndBlocker.
	blockExpirationsForOrders map[uint32]map[types.OrderId]bool
	// A map of all known canceled order IDs mapped to their expiry block.
	orderIdToCancelExpiry map[types.OrderId]uint32
	// A map from a block height to a set of all canceled order IDs that expire at the block.
	cancelExpiryToOrderIds map[uint32]map[types.OrderId]bool
}

// GetSide returns the Bid-side levels if `isBuy == true` otherwise, returns the Ask-side levels.
func (ob *Orderbook) GetSide(isBuy bool) map[types.Subticks]*types.Level {
	if isBuy {
		return ob.Bids
	}
	return ob.Asks
}

// GetMidPrice returns the mid price of the orderbook and whether or not it exists.
func (ob *Orderbook) GetMidPrice() (
	midPrice types.Subticks,
	exists bool,
) {
	if ob.BestBid == 0 || ob.BestAsk == math.MaxUint64 {
		return 0, false
	}
	return ob.BestBid + (ob.BestAsk-ob.BestBid)/2, true
}

// findNextBestLevelOrder is a helper method for finding the next best level order in the orderbook.
// Note that it does not modify any state in the orderbook.
// It will return the next best level order along with a boolean indicating whether the next best level order exists.
func (ob *Orderbook) findNextBestLevelOrder(
	levelOrder *types.LevelOrder,
) (
	nextBestLevelOrder *types.LevelOrder,
	foundOrder bool,
) {
	// Use the next order in the level if it exists.
	if levelOrder.Next != nil {
		return levelOrder.Next, true
	}

	// No more orders in this level exist. Attempt to get the first order at the next best level, if it exists.
	order := levelOrder.Value.Order
	subticks := order.GetOrderSubticks()
	isBuy := order.IsBuy()

	nextBestSubticks, foundOrder := ob.findNextBestSubticks(subticks, isBuy)
	if !foundOrder {
		return nil, false
	}

	nextBestLevelOrder, foundOrder = ob.getFirstOrderAtSideAndSubticks(
		isBuy,
		nextBestSubticks,
	)
	return nextBestLevelOrder, foundOrder
}

// findNextBestSubticks finds the next best price level subticks for a given orderbook and side.
// It naively scans the book until it finds a populated level. If no populated level has been found after
// N scans (where N is the number of populated levels),
// then simply take the set of all populated levels, and iterates over all of them to get the result.
// It returns the next best subtick for the given orderbook and side, and a boolean to indicate whether it exists
// (false if no "next best" subticks exist on that side of the book).
func (ob *Orderbook) findNextBestSubticks(
	startingTicks types.Subticks,
	isBuy bool,
) (
	nextBestSubtick types.Subticks,
	found bool,
) {
	var curSubticks types.Subticks = startingTicks
	var curLevel *types.Level

	// Search iteratively through the map of all levels from the current subticks
	// in multiples of `SubticksPerTick`. If the number of iterations exceeds the
	// number of levels on a given side of the book, then we will fallback on iterating
	// through all levels on this side of the book to find the best remaining level. We
	// do this so the worst case runtime is `O(n)`, where `n` is the number of prices on this side.
	//
	// Note that `curSubticks` can not overflow because:
	// 1. We know there is another price level on the book (`numLevels` is not zero).
	// 2. We are always moving closer to the next best price level on each iteration.
	levels := ob.GetSide(isBuy)
	numLevels := len(levels)
	orderbookSubticksPerTick := types.Subticks(ob.SubticksPerTick)

	for i := 0; i < numLevels; i++ {
		if isBuy {
			curSubticks -= orderbookSubticksPerTick
		} else {
			curSubticks += orderbookSubticksPerTick
		}

		curLevel = levels[curSubticks]
		if curLevel != nil {
			return curSubticks, true
		}
	}

	// The level was not found after `numLevels` of iteration so we fallback on iterating
	// through all levels on this side of the book to find the best remaining level. We
	// do this so the worst case runtime is `O(n)`, where `n` is the number of prices on this side.
	// Set `nextBestSubtick` to the worst possible value for each side.
	if isBuy {
		nextBestSubtick = 0
	} else {
		nextBestSubtick = math.MaxUint64
	}

	// If a subtick is found with a better price than `nextBestSubtick`, then it is set as the new `nextBestSubtick`.
	for subtick := range levels {
		// If the current subtick is better than the best seen subtick and worse than the starting subtick, mark it
		// as the new best seen subtick.
		if isBuy && nextBestSubtick < subtick && startingTicks > subtick {
			nextBestSubtick = subtick
			found = true
		} else if !isBuy && nextBestSubtick > subtick && startingTicks < subtick {
			nextBestSubtick = subtick
			found = true
		}
	}

	return nextBestSubtick, found
}

// getFirstOrderAtSideAndSubticks returns the first level order at a specific side and subticks for the passed in
// orderbook, along with a boolean indicating whether any order was found. If there are no orders on that side of
// the orderbook, the returned boolean will be `false` to indicate this.
func (ob *Orderbook) getFirstOrderAtSideAndSubticks(
	isBuy bool,
	subticks types.Subticks,
) (
	firstOrder *types.LevelOrder,
	found bool,
) {
	levels := ob.GetSide(isBuy)
	levelOrders, exists := levels[subticks]
	// If no orders at this price level exist, return false.
	if !exists {
		return nil, false
	}

	return levelOrders.LevelOrders.Front, true
}

// getBestOrderOnSide returns a reference to the best order on the passed in side of the book, along with a boolean
// indicating whether such an order exists at that subtick.
func (ob *Orderbook) getBestOrderOnSide(
	isBuy bool,
) (
	bestOrder *types.LevelOrder,
	found bool,
) {
	var bestSubticks types.Subticks
	if isBuy {
		bestSubticks = ob.BestBid
	} else {
		bestSubticks = ob.BestAsk
	}

	return ob.getFirstOrderAtSideAndSubticks(isBuy, bestSubticks)
}

// hasOrder returns true if the order by ID exists on the book.
func (ob *Orderbook) hasOrder(
	orderId types.OrderId,
) bool {
	_, exists := ob.orderIdToLevelOrder[orderId]
	return exists
}

// getOrder gets an order by ID and returns it.
func (ob *Orderbook) getOrder(
	orderId types.OrderId,
) (order types.Order, found bool) {
	levelOrder, exists := ob.orderIdToLevelOrder[orderId]
	if !exists {
		return types.Order{}, exists
	}

	return levelOrder.Value.Order, true
}

// getSubaccountOrders gets all of a subaccount's order on a specific CLOB and side.
// This function will panic if `side` is invalid or if the orderbook does not exist.
func (ob *Orderbook) getSubaccountOrders(
	subaccountId satypes.SubaccountId,
	side types.Order_Side,
) (openOrders []types.Order, err error) {
	if side == types.Order_SIDE_UNSPECIFIED {
		return openOrders, errors.WithStack(types.ErrInvalidOrderSide)
	}

	openClobOrdersForSubaccount, exists := ob.SubaccountOpenClobOrders[subaccountId]
	if !exists {
		return openOrders, nil
	}

	openClobOrdersForSubaccountAndSide, exists := openClobOrdersForSubaccount[side]
	if !exists {
		return openOrders, nil
	}

	// For each order ID, get the corresponding order.
	openOrders = make([]types.Order, len(openClobOrdersForSubaccountAndSide))
	i := 0
	for orderId := range openClobOrdersForSubaccountAndSide {
		order, found := ob.getOrder(orderId)
		if !found {
			panic("Open subaccount order does not exist in memclob")
		}
		openOrders[i] = order
		i++
	}

	return openOrders, nil
}

// mustAddShortTermOrderToBlockExpirationsForOrders is a function used for providing a simple interface
// for adding Short-Term orders to the `blockExpirationsForOrders` data structure. It will add an
// order to the set of orders expiring at that block, and if necessary will create any intermediate maps
// that do not already exist.
// This function assumes that it will only be called with Short-Term orders that have passed order validation
// in the CLOB keeper and the `validateNewOrder` function.
func (ob *Orderbook) mustAddShortTermOrderToBlockExpirationsForOrders(
	order types.Order,
) {
	if !order.OrderId.IsShortTermOrder() {
		panic(
			fmt.Sprintf(
				"mustAddShortTermOrderToBlockExpirationsForOrders: order ID %v is not a Short-Term order",
				order.OrderId,
			),
		)
	}

	// Create the map containing the set of orders expiring at this block, if it doesn't already exist.
	ordersExpiringAtBlock, exists := ob.blockExpirationsForOrders[order.GetGoodTilBlock()]
	if !exists {
		ordersExpiringAtBlock = make(map[types.OrderId]bool)
		ob.blockExpirationsForOrders[order.GetGoodTilBlock()] = ordersExpiringAtBlock
	}
	ordersExpiringAtBlock[order.OrderId] = true
}

// mustAddOrderToSubaccountOrders is a function used for providing a simple interface for adding orders to the
// `SubaccountOpenClobOrders` data structure on an orderbook. It will add an order to a subaccount's currently
// open orders, and if necessary will create any intermediate maps that do not already exist.
// This function assumes that it will only be called with orders that have passed order validation
// in the CLOB keeper and the `validateNewOrder` function.
func (ob *Orderbook) mustAddOrderToSubaccountOrders(
	order types.Order,
) {
	// Create the map containing all of a subaccount's open orders on this CLOB,
	// if it doesn't exist.
	subaccountId := order.OrderId.SubaccountId
	subaccountOpenClobOrders, exists := ob.SubaccountOpenClobOrders[subaccountId]
	if !exists {
		subaccountOpenClobOrders = make(map[types.Order_Side]map[types.OrderId]bool)
		ob.SubaccountOpenClobOrders[subaccountId] = subaccountOpenClobOrders
	}

	// Create the map containing all of a subaccount's open orders on this CLOB
	// on this side, if it doesn't exist.
	subaccountOpenClobOrdersSide, exists := subaccountOpenClobOrders[order.Side]
	if !exists {
		subaccountOpenClobOrdersSide = make(map[types.OrderId]bool)
		subaccountOpenClobOrders[order.Side] = subaccountOpenClobOrdersSide
	}

	// Add the order to the subaccount's open orders on this CLOB and side.
	subaccountOpenClobOrdersSide[order.OrderId] = true
}

// mustAddOrderToOrderbook will add the order to the resting orderbook.
// This function will assume that all order validation has already been done.
// If `forceToFrontOfLevel` is true, places the order at the head of the level,
// otherwise places it at the tail.
func (ob *Orderbook) mustAddOrderToOrderbook(
	newOrder types.Order,
	forceToFrontOfLevel bool,
) {
	// Verify that the order has a valid side, and panic if that's not the case.
	newOrder.MustBeValidOrderSide()

	isBuy := newOrder.IsBuy()
	clobOrder := types.ClobOrder{
		Order: newOrder,
	}

	// Add the order to the proper side of the orderbook and update orderbook state as necessary.
	if isBuy {
		// Update the `orderbook.BestBid` if the new order is a bid and has a higher price.
		if ob.BestBid < newOrder.GetOrderSubticks() {
			ob.BestBid = newOrder.GetOrderSubticks()
		}
	} else {
		// Update the `orderbook.BestAsk` if the new order is an ask and has a lower price.
		if ob.BestAsk > newOrder.GetOrderSubticks() {
			ob.BestAsk = newOrder.GetOrderSubticks()
		}
	}

	// Create the price level if it doesn't exist (it will only exist if there is at least one order at this price level).
	orders := ob.GetSide(isBuy)
	level, exists := orders[newOrder.GetOrderSubticks()]
	if !exists {
		level = &types.Level{
			LevelOrders: *list.New[types.ClobOrder](),
		}
		orders[newOrder.GetOrderSubticks()] = level
	}

	// Verify that the order ID is not already present in the `orderIdToLevelOrder` mapping. If not,
	// this means a replacement order was not properly validated and we should panic.
	if levelOrder, exists := ob.orderIdToLevelOrder[newOrder.OrderId]; exists {
		panic(
			fmt.Sprintf(
				"mustAddOrderToOrderbook: order (%+v) should be removed from the orderbook before replacing with order (%+v)",
				levelOrder.Value.Order,
				newOrder,
			),
		)
	}

	// Add the order to:
	// - The price level (either front or back)
	// - The orderIdToLevelOrder mapping
	if forceToFrontOfLevel {
		level.LevelOrders.PushFront(clobOrder)
		ob.orderIdToLevelOrder[newOrder.OrderId] = level.LevelOrders.Front
	} else {
		level.LevelOrders.PushBack(clobOrder)
		ob.orderIdToLevelOrder[newOrder.OrderId] = level.LevelOrders.Back
	}

	// Add the order to the subaccount's currently open orders.
	ob.mustAddOrderToSubaccountOrders(newOrder)

	// Increment the total number of open orders for this orderbook.
	ob.TotalOpenOrders++

	// If the order is a Short-Term order, add the order to the order block expirations map.
	if newOrder.IsShortTermOrder() {
		ob.mustAddShortTermOrderToBlockExpirationsForOrders(newOrder)
	}

	// If the order is reduce-only, add it to the open reduce-only orders for this subaccount.
	if newOrder.IsReduceOnly() {
		openReduceOnlyOrders, exists := ob.SubaccountOpenReduceOnlyOrders[newOrder.OrderId.SubaccountId]
		if !exists {
			openReduceOnlyOrders = make(map[types.OrderId]bool)
		}
		openReduceOnlyOrders[newOrder.OrderId] = true
		ob.SubaccountOpenReduceOnlyOrders[newOrder.OrderId.SubaccountId] = openReduceOnlyOrders
	}
}

// mustRemoveOrder completely removes an order from all data structures for tracking
// open orders in the memclob. If the order does not exist, this method will panic.
// NOTE: `mustRemoveOrder` does _not_ remove cancels.
// TODO(DEC-847): Remove stateful orders properly.
func (ob *Orderbook) mustRemoveOrder(
	levelOrder *types.LevelOrder,
) {
	// Define variables related to this order for more succinct reference.
	order := levelOrder.Value.Order
	orderId := order.OrderId
	subaccountId := orderId.SubaccountId
	side := order.Side
	subticks := order.GetOrderSubticks()
	isBuy := order.IsBuy()

	// Delete this order from various data structures.
	// If this is a Short-Term order, remove the order from `blockExpirationsForOrders[goodTilBlock]`.
	if order.OrderId.IsShortTermOrder() {
		goodTilBlock := order.GetGoodTilBlock()
		delete(ob.blockExpirationsForOrders[goodTilBlock], orderId)
		if len(ob.blockExpirationsForOrders[goodTilBlock]) == 0 {
			delete(ob.blockExpirationsForOrders, goodTilBlock)
		}
	}

	delete(ob.SubaccountOpenClobOrders[subaccountId][side], orderId)
	if len(ob.SubaccountOpenClobOrders[subaccountId][side]) == 0 {
		delete(ob.SubaccountOpenClobOrders[subaccountId], side)

		if len(ob.SubaccountOpenClobOrders[subaccountId]) == 0 {
			delete(ob.SubaccountOpenClobOrders, subaccountId)
		}
	}

	delete(ob.orderIdToLevelOrder, orderId)

	// If this is a reduce-only order, remove it from the open reduce-only orders for
	// this subaccount. If the subaccount has no more open reduce-only orders, delete the inner map.
	// TODO(DEC-847): Remove stateful reduce-only orders properly.
	if order.IsReduceOnly() {
		delete(ob.SubaccountOpenReduceOnlyOrders[subaccountId], orderId)
		if len(ob.SubaccountOpenReduceOnlyOrders[subaccountId]) == 0 {
			delete(ob.SubaccountOpenReduceOnlyOrders, subaccountId)
		}
	}

	// Next we'll remove the order from the orderbook. This process involves the following steps:
	// 1. Fetch the relevant orderbook and side for this order via the `orderbooksMap`.
	// 2. Remove the `levelOrder` from the `LevelOrders` linked list within the `Level`.
	// 3. If `LevelOrders.Front` is not `nil` (i.e. There exist other orders at this level), then we are done.
	// Otherwise, continue.
	// 4. If there exist no other orders at this level, then delete the `Level` from the appropriate side of the book.
	// 5. If the price of the deleted order was not the best price on its side of the book, then we are done.
	// Otherwise, continue.
	// 6. If the price of the deleted order _was_ the best price on its side of the book, then we need to find a new best
	// price. While the current `Level` is `nil`, attempt to find the next best price by iterating from the previous best
	// price in increments of `SubticksPerTick` toward the worst possible price (decrementing if bids, incrementing if
	// asks).
	// 7. If a Level is found in fewer iterations than there are total levels, then this is the next best level and we
	// are done. Otherwise, continue.
	// 8. If after `totalLevels` of iteration, we still have not found the next best level, then we fall back to
	// iterating over every level on the appropriate side to find the next best price.

	levels := ob.GetSide(isBuy)
	level, levelExists := levels[subticks]
	if !levelExists {
		panic("mustRemoveOrder: Level does not exist for order")
	}

	// Remove the order from the level.
	level.LevelOrders.Remove(levelOrder)

	// Decrement the total number of orders for this orderbook.
	ob.TotalOpenOrders--

	// Edge case: If this removed order was the last order in its
	// level, we need to remove the level altogether.
	if level.LevelOrders.Front == nil {
		delete(levels, subticks)
		if isBuy {
			if subticks == ob.BestBid {
				// Edge case: If this removed level represented the best price level for this side
				// of the book, we need to find the next best price level.
				nextBestSubticks, _ := ob.findNextBestSubticks(subticks, true)
				ob.BestBid = nextBestSubticks
			}
		} else {
			if subticks == ob.BestAsk {
				// Edge case: If this removed level represented the best price level for this side
				// of the book, we need to find the next best price level.
				nextBestSubticks, _ := ob.findNextBestSubticks(subticks, false)
				ob.BestAsk = nextBestSubticks
			}
		}
	}
}

// getCancel returns the `tilBlock` expiry of an order cancelation and a bool indicating whether the expiry exists.
func (ob *Orderbook) getCancel(
	orderId types.OrderId,
) (
	tilBlock uint32,
	exists bool,
) {
	tilBlock, exists = ob.orderIdToCancelExpiry[orderId]
	return tilBlock, exists
}

// addShortTermCancel adds a cancel expiring at block `tilBlock`.
// Panics if the cancel already exists.
func (ob *Orderbook) addShortTermCancel(
	orderId types.OrderId,
	tilBlock uint32,
) {
	orderId.MustBeShortTermOrder()
	// Add the `orderId` to the `orderIdToExpiry` map, panicing if it already exists.
	if _, exists := ob.orderIdToCancelExpiry[orderId]; exists {
		panic(fmt.Sprintf(
			"mustAddCancel: orderId %+v already exists in orderIdToExpiry",
			orderId,
		))
	}
	ob.orderIdToCancelExpiry[orderId] = tilBlock

	// Fetch a reference to the `expiryToOrderIds` for this `tilBlock`, creating it if it does not already exist.
	orderIdsInBlock, exists := ob.cancelExpiryToOrderIds[tilBlock]
	if !exists {
		orderIdsInBlock = make(map[types.OrderId]bool)
		ob.cancelExpiryToOrderIds[tilBlock] = orderIdsInBlock
	} else if _, exists = orderIdsInBlock[orderId]; exists {
		panic(fmt.Sprintf(
			"memclobCancels#add: orderId %+v already exists in expiryToOrderIds[%d]",
			orderId,
			tilBlock,
		))
	}

	// Set the `OrderId` in the `expiryToOrderIds` data structure for the new cancel's `tilBlock`.
	orderIdsInBlock[orderId] = true
}

// mustRemoveCancel removes a cancel. Panics if the `orderId` is not found.
func (ob *Orderbook) mustRemoveCancel(
	orderId types.OrderId,
) {
	// Panic if the `orderId` does not exist in the `orderIdToExpiry` map.
	goodTilBlock, exists := ob.orderIdToCancelExpiry[orderId]
	if !exists {
		panic(fmt.Sprintf(
			"memclobCancels#remove: orderId %+v does not exist in orderIdToExpiry",
			orderId,
		))
	}

	// Panic if the `orderId` does not exist in the appropriate submap of `expiryToOrderIds`.
	expiryToOrderIdsForBlock, exists := ob.cancelExpiryToOrderIds[goodTilBlock]
	if !exists {
		panic(fmt.Sprintf(
			"memclobCancels#remove: %d does not exist in expiryToOrderIds",
			goodTilBlock,
		))
	}
	if _, exists = expiryToOrderIdsForBlock[orderId]; !exists {
		panic(fmt.Sprintf(
			"memclobCancels#remove: orderId %+v does not exist in expiryToOrderIds[%d]",
			orderId,
			goodTilBlock,
		))
	}

	// Delete the `orderId` from the `orderIdToExpiry` map.
	delete(ob.orderIdToCancelExpiry, orderId)

	// Delete the `orderId` from the `expiryToOrderIds` submap.
	// If this is the last order in the submap, delete the submap.
	if len(expiryToOrderIdsForBlock) == 1 {
		delete(ob.cancelExpiryToOrderIds, goodTilBlock)
	} else {
		delete(expiryToOrderIdsForBlock, orderId)
	}
}

// removeAllCancelsAtBlock iterates through and removes all cancels that expire at a certain `block`.
func (ob *Orderbook) removeAllCancelsAtBlock(
	block uint32,
) {
	orderIds, exists := ob.cancelExpiryToOrderIds[block]

	// If map entry does not exist, return early.
	if !exists {
		return
	}

	// Remove all ids. This also removes `block` from the `expiryToOrderIds` map.
	for id := range orderIds {
		ob.mustRemoveCancel(id)
	}
}
