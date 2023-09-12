package memclob

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/pkg/errors"
	"github.com/zyedidia/generic/list"
)

// memclobOpenOrders is a utility struct used for storing orders (within an orderbook) and order expirations.
type memclobOpenOrders struct {
	// Holds every `Orderbook` by ID of the CLOB.
	orderbooksMap map[types.ClobPairId]*types.Orderbook
	// Map from order IDs to their `LevelOrder` reference contained in the
	// orderbook, necessary for O(1) order removal from the orderbook.
	orderIdToLevelOrder map[types.OrderId]*types.LevelOrder
	// Map from block number to a set of all orders that expire at this block
	// (with each order keyed by `OrderId`). Necessary for O(1) order removal
	// from the orderbook when expiring orders in the EndBlocker.
	blockExpirationsForOrders map[uint32]map[types.OrderId]bool
}

// newMemclobOpenOrders returns a new `memclobOpenOrders`.
func newMemclobOpenOrders() *memclobOpenOrders {
	return &memclobOpenOrders{
		orderbooksMap:             make(map[types.ClobPairId]*types.Orderbook),
		orderIdToLevelOrder:       make(map[types.OrderId]*types.LevelOrder),
		blockExpirationsForOrders: make(map[uint32]map[types.OrderId]bool),
	}
}

// mustGetOrderbook returns the orderbook for the given clobPairId. Panics if the orderbook cannot be found.
func (m *memclobOpenOrders) mustGetOrderbook(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
) *types.Orderbook {
	orderbook, exists := m.orderbooksMap[clobPairId]
	if !exists {
		panic(fmt.Sprintf("No orderbook exists with id %d", clobPairId))
	}
	return orderbook
}

// findNextBestLevelOrder is a helper method for finding the next best level order in the orderbook.
// Note that it does not modify any state in the memclob.
// It will return the next best level order along with a boolean indicating whether the next best level order exists.
func (m *memclobOpenOrders) findNextBestLevelOrder(
	ctx sdk.Context,
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
	orderbook := m.mustGetOrderbook(ctx, order.GetClobPairId())
	isBuy := order.IsBuy()

	nextBestSubticks, foundOrder := m.findNextBestSubticks(ctx, subticks, orderbook, isBuy)
	if !foundOrder {
		return nil, false
	}

	nextBestLevelOrder, foundOrder = m.getFirstOrderAtSideAndSubticks(
		orderbook,
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
func (m *memclobOpenOrders) findNextBestSubticks(
	ctx sdk.Context,
	startingTicks types.Subticks,
	orderbook *types.Orderbook,
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
	levels := orderbook.GetSide(isBuy)
	numLevels := len(levels)
	orderbookSubticksPerTick := types.Subticks(orderbook.SubticksPerTick)

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

// getBestOrderOnSide returns a reference to the best order on the passed in side of the book, along with a boolean
// indicating whether such an order exists at that subtick.
func (m *memclobOpenOrders) getBestOrderOnSide(
	orderbook *types.Orderbook,
	isBuy bool,
) (
	bestOrder *types.LevelOrder,
	found bool,
) {
	var bestSubticks types.Subticks
	if isBuy {
		bestSubticks = orderbook.BestBid
	} else {
		bestSubticks = orderbook.BestAsk
	}

	return m.getFirstOrderAtSideAndSubticks(orderbook, isBuy, bestSubticks)
}

// getFirstOrderAtSideAndSubticks returns the first level order at a specific side and subticks for the passed in
// orderbook, along with a boolean indicating whether any order was found. If there are no orders on that side of
// the orderbook, the returned boolean will be `false` to indicate this.
func (m *memclobOpenOrders) getFirstOrderAtSideAndSubticks(
	orderbook *types.Orderbook,
	isBuy bool,
	subticks types.Subticks,
) (
	firstOrder *types.LevelOrder,
	found bool,
) {
	levels := orderbook.GetSide(isBuy)
	levelOrders, exists := levels[subticks]
	// If no orders at this price level exist, return false.
	if !exists {
		return nil, false
	}

	return levelOrders.LevelOrders.Front, true
}

// hasOrder returns true if the order by ID exists on the book.
func (m *memclobOpenOrders) hasOrder(
	ctx sdk.Context,
	orderId types.OrderId,
) bool {
	_, exists := m.orderIdToLevelOrder[orderId]
	return exists
}

// getOrder gets an order by ID and returns it.
func (m *memclobOpenOrders) getOrder(
	ctx sdk.Context,
	orderId types.OrderId,
) (order types.Order, found bool) {
	levelOrder, exists := m.orderIdToLevelOrder[orderId]
	if !exists {
		return types.Order{}, exists
	}

	return levelOrder.Value.Order, true
}

// getSubaccountOrders gets all of a subaccount's order on a specific CLOB and side.
// This function will panic if `side` is invalid or if the orderbook does not exist.
func (m *memclobOpenOrders) getSubaccountOrders(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
	subaccountId satypes.SubaccountId,
	side types.Order_Side,
) (openOrders []types.Order, err error) {
	if side == types.Order_SIDE_UNSPECIFIED {
		return openOrders, errors.WithStack(types.ErrInvalidOrderSide)
	}

	orderbook, exists := m.orderbooksMap[clobPairId]
	// If the CLOB doesn't exist, then `clobPairId` is invalid.
	if !exists {
		return openOrders, errorsmod.Wrapf(
			types.ErrInvalidClob,
			"Invalid ClobPair ID: %d",
			clobPairId,
		)
	}

	openClobOrdersForSubaccount, exists := orderbook.SubaccountOpenClobOrders[subaccountId]
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
		order, found := m.getOrder(ctx, orderId)
		if !found {
			panic("Open subaccount order does not exist in memclob")
		}
		openOrders[i] = order
		i++
	}

	return openOrders, nil
}

// createOrderbook is used for updating memclob internal data structures to mark an orderbook as created.
// This function will panic if `clobPairId` already exists in any of the memclob's internal data structures.
func (m *memclobOpenOrders) createOrderbook(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
	subticksPerTick types.SubticksPerTick,
	minOrderBaseQuantums satypes.BaseQuantums,
) {
	if _, exists := m.orderbooksMap[clobPairId]; exists {
		panic(fmt.Sprintf("Orderbook for ClobPair ID %d already exists", clobPairId))
	}

	if subticksPerTick == 0 {
		panic("subticksPerTick must be greater than zero")
	}

	if minOrderBaseQuantums == 0 {
		panic("minOrderBaseQuantums must be greater than zero")
	}

	m.orderbooksMap[clobPairId] = &types.Orderbook{
		Asks:                           make(map[types.Subticks]*types.Level),
		BestAsk:                        math.MaxUint64,
		BestBid:                        0,
		Bids:                           make(map[types.Subticks]*types.Level),
		MinOrderBaseQuantums:           minOrderBaseQuantums,
		SubaccountOpenClobOrders:       make(map[satypes.SubaccountId]map[types.Order_Side]map[types.OrderId]bool),
		SubticksPerTick:                subticksPerTick,
		SubaccountOpenReduceOnlyOrders: make(map[satypes.SubaccountId]map[types.OrderId]bool),
	}
}

// mustAddShortTermOrderToBlockExpirationsForOrders is a function used for providing a simple interface
// for adding Short-Term orders to the `blockExpirationsForOrders` data structure. It will add an
// order to the set of orders expiring at that block, and if necessary will create any intermediate maps
// that do not already exist.
// This function assumes that it will only be called with Short-Term orders that have passed order validation
// in the CLOB keeper and the `validateNewOrder` function.
func (m *memclobOpenOrders) mustAddShortTermOrderToBlockExpirationsForOrders(
	ctx sdk.Context,
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
	ordersExpiringAtBlock, exists := m.blockExpirationsForOrders[order.GetGoodTilBlock()]
	if !exists {
		ordersExpiringAtBlock = make(map[types.OrderId]bool)
		m.blockExpirationsForOrders[order.GetGoodTilBlock()] = ordersExpiringAtBlock
	}
	ordersExpiringAtBlock[order.OrderId] = true
}

// mustAddOrderToSubaccountOrders is a function used for providing a simple interface for adding orders to the
// `SubaccountOpenClobOrders` data structure on an orderbook. It will add an order to a subaccount's currently
// open orders, and if necessary will create any intermediate maps that do not already exist.
// This function assumes that it will only be called with orders that have passed order validation
// in the CLOB keeper and the `validateNewOrder` function.
// If `order.Side` is an invalid side or `order.ClobPairId` does not reference a valid CLOB, this function will panic.
func (m *memclobOpenOrders) mustAddOrderToSubaccountOrders(
	ctx sdk.Context,
	order types.Order,
) {
	// If `ClobPairId` does not reference a valid CLOB, panic.
	orderbook, exists := m.orderbooksMap[order.GetClobPairId()]
	if !exists {
		panic(types.ErrInvalidClob)
	}

	// Create the map containing all of a subaccount's open orders on this CLOB,
	// if it doesn't exist.
	subaccountId := order.OrderId.SubaccountId
	subaccountOpenClobOrders, exists := orderbook.SubaccountOpenClobOrders[subaccountId]
	if !exists {
		subaccountOpenClobOrders = make(map[types.Order_Side]map[types.OrderId]bool)
		orderbook.SubaccountOpenClobOrders[subaccountId] = subaccountOpenClobOrders
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
func (m *memclobOpenOrders) mustAddOrderToOrderbook(
	ctx sdk.Context,
	newOrder types.Order,
	forceToFrontOfLevel bool,
) {
	// Verify that the order has a valid side, and panic if that's not the case.
	newOrder.MustBeValidOrderSide()

	// Initialize variables used for traversing the orderbook.
	orderbook := m.mustGetOrderbook(ctx, newOrder.GetClobPairId())
	isBuy := newOrder.IsBuy()
	clobOrder := types.ClobOrder{
		Order: newOrder,
	}

	// Add the order to the proper side of the orderbook and update orderbook state as necessary.
	if isBuy {
		// Update the `orderbook.BestBid` if the new order is a bid and has a higher price.
		if orderbook.BestBid < newOrder.GetOrderSubticks() {
			orderbook.BestBid = newOrder.GetOrderSubticks()
		}
	} else {
		// Update the `orderbook.BestAsk` if the new order is an ask and has a lower price.
		if orderbook.BestAsk > newOrder.GetOrderSubticks() {
			orderbook.BestAsk = newOrder.GetOrderSubticks()
		}
	}

	// Create the price level if it doesn't exist (it will only exist if there is at least one order at this price level).
	orders := orderbook.GetSide(isBuy)
	level, exists := orders[newOrder.GetOrderSubticks()]
	if !exists {
		level = &types.Level{
			LevelOrders: *list.New[types.ClobOrder](),
		}
		orders[newOrder.GetOrderSubticks()] = level
	}

	// Verify that the order ID is not already present in the `orderIdToLevelOrder` mapping. If not,
	// this means a replacement order was not properly validated and we should panic.
	if levelOrder, exists := m.orderIdToLevelOrder[newOrder.OrderId]; exists {
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
		m.orderIdToLevelOrder[newOrder.OrderId] = level.LevelOrders.Front
	} else {
		level.LevelOrders.PushBack(clobOrder)
		m.orderIdToLevelOrder[newOrder.OrderId] = level.LevelOrders.Back
	}

	// Add the order to the subaccount's currently open orders.
	m.mustAddOrderToSubaccountOrders(ctx, newOrder)

	// Increment the total number of open orders for this orderbook.
	orderbook.TotalOpenOrders++

	// If the order is a Short-Term order, add the order to the order block expirations map.
	if newOrder.IsShortTermOrder() {
		m.mustAddShortTermOrderToBlockExpirationsForOrders(ctx, newOrder)
	}

	// If the order is reduce-only, add it to the open reduce-only orders for this subaccount.
	if newOrder.IsReduceOnly() {
		openReduceOnlyOrders, exists := orderbook.SubaccountOpenReduceOnlyOrders[newOrder.OrderId.SubaccountId]
		if !exists {
			openReduceOnlyOrders = make(map[types.OrderId]bool)
		}
		openReduceOnlyOrders[newOrder.OrderId] = true
		orderbook.SubaccountOpenReduceOnlyOrders[newOrder.OrderId.SubaccountId] = openReduceOnlyOrders
	}
}

// mustRemoveOrder completely removes an order from all data structures for tracking
// open orders in the memclob. If the order does not exist, this method will panic.
// NOTE: `mustRemoveOrder` does _not_ remove cancels.
// TODO(DEC-847): Remove stateful orders properly.
func (m *memclobOpenOrders) mustRemoveOrder(
	ctx sdk.Context,
	levelOrder *types.LevelOrder,
) {
	// Define variables related to this order for more succinct reference.
	order := levelOrder.Value.Order
	orderId := order.OrderId
	subaccountId := orderId.SubaccountId
	clobPairId := order.GetClobPairId()
	side := order.Side
	subticks := order.GetOrderSubticks()
	isBuy := order.IsBuy()

	// Delete this order from various data structures.
	// If this is a Short-Term order, remove the order from `blockExpirationsForOrders[goodTilBlock]`.
	if order.OrderId.IsShortTermOrder() {
		goodTilBlock := order.GetGoodTilBlock()
		delete(m.blockExpirationsForOrders[goodTilBlock], orderId)
		if len(m.blockExpirationsForOrders[goodTilBlock]) == 0 {
			delete(m.blockExpirationsForOrders, goodTilBlock)
		}
	}

	delete(m.orderbooksMap[clobPairId].SubaccountOpenClobOrders[subaccountId][side], orderId)
	if len(m.orderbooksMap[clobPairId].SubaccountOpenClobOrders[subaccountId][side]) == 0 {
		delete(m.orderbooksMap[clobPairId].SubaccountOpenClobOrders[subaccountId], side)

		if len(m.orderbooksMap[clobPairId].SubaccountOpenClobOrders[subaccountId]) == 0 {
			delete(m.orderbooksMap[clobPairId].SubaccountOpenClobOrders, subaccountId)
		}
	}

	delete(m.orderIdToLevelOrder, orderId)

	// If this is a reduce-only order, remove it from the open reduce-only orders for
	// this subaccount. If the subaccount has no more open reduce-only orders, delete the inner map.
	// TODO(DEC-847): Remove stateful reduce-only orders properly.
	if order.IsReduceOnly() {
		delete(m.orderbooksMap[clobPairId].SubaccountOpenReduceOnlyOrders[subaccountId], orderId)
		if len(m.orderbooksMap[clobPairId].SubaccountOpenReduceOnlyOrders[subaccountId]) == 0 {
			delete(m.orderbooksMap[clobPairId].SubaccountOpenReduceOnlyOrders, subaccountId)
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

	orderbook := m.mustGetOrderbook(ctx, clobPairId)
	levels := orderbook.GetSide(isBuy)
	level, levelExists := levels[subticks]
	if !levelExists {
		panic("mustRemoveOrder: Level does not exist for order")
	}

	// Remove the order from the level.
	level.LevelOrders.Remove(levelOrder)

	// Decrement the total number of orders for this orderbook.
	orderbook.TotalOpenOrders--

	// Edge case: If this removed order was the last order in its
	// level, we need to remove the level altogether.
	if level.LevelOrders.Front == nil {
		delete(levels, subticks)
		if isBuy {
			if subticks == orderbook.BestBid {
				// Edge case: If this removed level represented the best price level for this side
				// of the book, we need to find the next best price level.
				nextBestSubticks, _ := m.findNextBestSubticks(ctx, subticks, orderbook, true)
				orderbook.BestBid = nextBestSubticks
			}
		} else {
			if subticks == orderbook.BestAsk {
				// Edge case: If this removed level represented the best price level for this side
				// of the book, we need to find the next best price level.
				nextBestSubticks, _ := m.findNextBestSubticks(ctx, subticks, orderbook, false)
				orderbook.BestAsk = nextBestSubticks
			}
		}
	}
}
