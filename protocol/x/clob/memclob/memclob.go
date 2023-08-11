package memclob

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/clob/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// Ensure that `memClobPriceTimePriority` struct properly implements
// the `MemClob` interface.
var _ types.MemClob = &MemClobPriceTimePriority{}

type MemClobPriceTimePriority struct {
	// ---- Fields for open orders ----
	// Struct for storing all open orders (including their expiries).
	openOrders *memclobOpenOrders

	// ---- Fields for canceled orders ----
	// Struct for storing order cancelations (including their expiries).
	cancels *memclobCancels

	// ---- Fields for pending fills ----
	pendingFills *memclobPendingFills

	// A reference to an expected clob keeper.
	clobKeeper types.MemClobKeeper

	// ---- Fields for replicating state ----
	// A map from a perpetual ID to a list of CLOB pair IDs that trade that perpetual. Used for
	// determining which CLOB to place a liquidation order on, when liquidating a perpetual.
	perpetualIdToClobPairId map[uint32][]types.ClobPairId

	// ---- Fields for determining if off-chain update messages should be generated ----
	generateOffchainUpdates bool

	// A reference to an interface that returns random bools.
	randomBooler types.RandomBooler
}

func NewMemClobPriceTimePriority(
	generateOffchainUpdates bool,
) *MemClobPriceTimePriority {
	return &MemClobPriceTimePriority{
		openOrders:              newMemclobOpenOrders(),
		cancels:                 newMemclobCancels(),
		pendingFills:            newMemclobPendingFills(),
		perpetualIdToClobPairId: make(map[uint32][]types.ClobPairId),
		generateOffchainUpdates: generateOffchainUpdates,
		randomBooler:            &types.RealRandomBooler{},
	}
}

// SetRandomBooler sets the RandomBooler implementation for the memclob.
func (m *MemClobPriceTimePriority) SetRandomBooler(randomBooler types.RandomBooler) {
	m.randomBooler = randomBooler
}

// SetClobKeeper sets the MemClobKeeper reference for this MemClob.
// This method is called after the MemClob struct is initialized.
// This reference is set with an explicit method call rather than during `NewMemClobPriceTimePriority`
// due to the bidirectional dependency between the Keeper and the MemClob.
func (m *MemClobPriceTimePriority) SetClobKeeper(clobKeeper types.MemClobKeeper) {
	m.clobKeeper = clobKeeper
}

// CancelOrder removes an order by `OrderId` (if it exists) from all order-related data structures
// in the memclob. Order cancellation can be stateful or short term.
// For short-term orders, CancelOrder adds (or updates) a cancel to the desired `goodTilBlock` in the memclob.
// If a cancel already exists for this order with a lower `goodTilBlock`, the cancel is updated to the
// new `goodTilBlock`.
//
// For short-term cancels, an error will be returned if any of the following conditions are true:
// - A cancel already exists for this order with the same or greater `GoodTilBlock`.
//
// If the order is removed from the orderbook, an off-chain update message will be generated.
func (m *MemClobPriceTimePriority) CancelOrder(
	ctx sdk.Context,
	msgCancelOrder *types.MsgCancelOrder,
) (offchainUpdates *types.OffchainUpdates, err error) {
	orderIdToCancel := msgCancelOrder.GetOrderId()

	// TODO(DEC-1949): Don't store order cancelations in state, and return an error if the associated
	// cancel does not exist.
	if orderIdToCancel.IsStatefulOrder() {
		// Always add the stateful order cancelation to the operations to propose.
		m.pendingFills.operationsToPropose.AddOrderCancellationToOperationsQueue(*msgCancelOrder)

		// Remove the stateful order from the book if it exists.
		// TODO(DEC-1934): Replay/re-add stateful orders that were canceled during `PrepareCheckState`.
		if levelOrder, exists := m.openOrders.orderIdToLevelOrder[orderIdToCancel]; exists {
			m.mustRemoveOrder(ctx, orderIdToCancel)

			if m.pendingFills.operationsToPropose.IsMakerOrderPreexistingStatefulOrder(levelOrder.Value.Order) {
				// If the cancelation is for a pre-existing stateful order that is not present in the operations to propose,
				// we can remove the nonce for the pre-existing stateful order.
				if !m.pendingFills.operationsToPropose.IsPreexistingStatefulOrderInOperationsQueue(levelOrder.Value.Order) {
					m.pendingFills.operationsToPropose.RemovePreexistingStatefulOrderPlacementNonce(levelOrder.Value.Order)
				}
			}
		}
	} else {
		// Retrieve the existing short-term cancel, if it exists.
		oldCancellationGoodTilBlock, cancelAlreadyExists := m.cancels.get(orderIdToCancel)
		goodTilBlock := msgCancelOrder.GetGoodTilBlock()

		// If the existing short-term cancel has the same or greater `goodTilBlock`, then there is
		// nothing for us to do. Return an error.
		if cancelAlreadyExists && oldCancellationGoodTilBlock >= goodTilBlock {
			return nil, types.ErrMemClobCancelAlreadyExists
		}

		// If there exists a resting order on the book with a `GoodTilBlock` not-greater-than that of
		// the short-term cancel, remove the order and add the order cancellation to the operations queue if necessary.
		// TODO(DEC-824): Perform correct cancellation validation of stateful orders.
		if levelOrder, orderExists := m.openOrders.orderIdToLevelOrder[orderIdToCancel]; orderExists &&
			goodTilBlock >= levelOrder.Value.Order.GetGoodTilBlock() {
			// If the canceled order exists in the `operationsToPropose`, then we need to add the cancelation
			// to the operations queue as well.
			if m.pendingFills.operationsToPropose.IsOrderPlacementInOperationsQueue(levelOrder.Value.Order) {
				m.pendingFills.operationsToPropose.AddOrderCancellationToOperationsQueue(*msgCancelOrder)
			} else {
				// If the cancelation is for an order that is not present in the operations to propose,
				// we can remove the nonce for the order.
				m.pendingFills.operationsToPropose.RemoveOrderPlacementNonce(levelOrder.Value.Order)
			}

			m.mustRemoveOrder(ctx, orderIdToCancel)
			telemetry.IncrCounter(1, types.ModuleName, metrics.CancelOrder, metrics.RemovedFromOrderBook)
		}

		// Remove the existing cancel, if any.
		if cancelAlreadyExists {
			m.cancels.remove(orderIdToCancel)
		}

		// Add the new order cancelation.
		m.cancels.addShortTermCancel(orderIdToCancel, goodTilBlock)
	}

	offchainUpdates = types.NewOffchainUpdates()
	if m.generateOffchainUpdates {
		if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
			ctx.Logger(),
			orderIdToCancel,
			off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_USER_CANCELED,
			off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
		); success {
			offchainUpdates.AddRemoveMessage(orderIdToCancel, message)
		}
	}

	return offchainUpdates, nil
}

// CreateOrderbook is used for updating memclob internal data structures to mark an orderbook as created.
// This function will panic if `clobPairId` already exists in any of the memclob's internal data structures.
func (m *MemClobPriceTimePriority) CreateOrderbook(
	ctx sdk.Context,
	clobPair types.ClobPair,
) {
	clobPairId := clobPair.GetClobPairId()
	subticksPerTick := clobPair.GetClobPairSubticksPerTick()
	minOrderBaseQuantums := clobPair.GetClobPairMinOrderBaseQuantums()

	// Create the in-memory orderbook for this `clobPairId`.
	m.openOrders.createOrderbook(ctx, clobPairId, subticksPerTick, minOrderBaseQuantums)

	// If this `ClobPair` is for a perpetual, add the `clobPairId` to the list of CLOB pair IDs
	// that facilitate trading of this perpetual.
	if perpetualClobMetadata := clobPair.GetPerpetualClobMetadata(); perpetualClobMetadata != nil {
		perpetualId := perpetualClobMetadata.PerpetualId
		clobPairIds, exists := m.perpetualIdToClobPairId[perpetualId]
		if !exists {
			clobPairIds = make([]types.ClobPairId, 0)
		}
		m.perpetualIdToClobPairId[perpetualId] = append(
			clobPairIds,
			clobPairId,
		)
	}
}

// GetOrder gets an order by ID and returns it.
func (m *MemClobPriceTimePriority) GetOrder(
	ctx sdk.Context,
	orderId types.OrderId,
) (order types.Order, found bool) {
	return m.openOrders.getOrder(ctx, orderId)
}

// GetCancelOrder returns the `tilBlock` expiry of an order cancelation and a bool indicating whether the expiry exists.
func (m *MemClobPriceTimePriority) GetCancelOrder(
	ctx sdk.Context,
	orderId types.OrderId,
) (tilBlock uint32, found bool) {
	return m.cancels.get(orderId)
}

// GetOrderFilledAmount returns the total filled amount of an order from state.
func (m *MemClobPriceTimePriority) GetOrderFilledAmount(
	ctx sdk.Context,
	orderId types.OrderId,
) satypes.BaseQuantums {
	exists, orderStateFilledAmount, _ := m.clobKeeper.GetOrderFillAmount(ctx, orderId)
	if !exists {
		orderStateFilledAmount = 0
	}

	return orderStateFilledAmount
}

// GetSubaccountOrders gets all of a subaccount's order on a specific CLOB and side.
// This function will panic if `side` is invalid or if the orderbook does not exist.
func (m *MemClobPriceTimePriority) GetSubaccountOrders(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
	subaccountId satypes.SubaccountId,
	side types.Order_Side,
) (openOrders []types.Order, err error) {
	return m.openOrders.getSubaccountOrders(
		ctx,
		clobPairId,
		subaccountId,
		side,
	)
}

// mustUpdateMemclobStateWithMatches updates the memclob state by applying matches to all bookkeeping data structures.
// Namely, it will perform the following operations:
//   - Update `orderHashToMatchableOrder` with all order hashes of newly matched orders, if they do not already exist.
//   - Update `orderIdToOrder` with all order hashes of newly matched orders, if they do not already exist,
//     excluding liquidation orders.
//   - Starting with the first matched order, add each new match to back of the `fills` and to the
//     `subaccountIdToFillIndexes`, for that respective subaccount.
//   - Update orderbook state by updating the filled amount of all matched maker orders, and removing them if fully
//     filled.
func (m *MemClobPriceTimePriority) mustUpdateMemclobStateWithMatches(
	ctx sdk.Context,
	takerOrder types.MatchableOrder,
	newPendingFills []types.PendingFill,
	newMakerFills []types.MakerFill,
	matchedOrderHashToOrder map[types.OrderHash]types.MatchableOrder,
	matchedMakerOrderIdToOrder map[types.OrderId]types.Order,
) (offchainUpdates *types.OffchainUpdates) {
	offchainUpdates = types.NewOffchainUpdates()

	// For each order, update `orderHashToMatchableOrder` and `orderIdToOrder`.
	// TODO(DEC-1924): Deprecate this below loop after CLOB refactor is completed.
	for orderHash, matchedOrder := range matchedOrderHashToOrder {
		// Update `orderHashToMatchableOrder`.
		if _, exists := m.pendingFills.orderHashToMatchableOrder[orderHash]; !exists {
			m.pendingFills.orderHashToMatchableOrder[orderHash] = matchedOrder
		}

		// If this is not a liquidation, update `orderIdToOrder`.
		if !matchedOrder.IsLiquidation() {
			order := matchedOrder.MustGetOrder()
			orderId := order.OrderId
			other, exists := m.pendingFills.orderIdToOrder[orderId]
			if exists && order.MustCmpReplacementOrder(&other) < 0 {
				// Shouldn't happen as the order should have already been replaced or rejected.
				panic("mustUpdateMemclobStateWithMatches: newly matched order is lesser than existing order")
			}
			m.pendingFills.orderIdToOrder[orderId] = order
		}
	}

	// Ensure each filled maker order has an order placement in `OperationsToPropose`.
	makerFillWithOrders := make([]types.MakerFillWithOrder, 0, len(newMakerFills))
	for _, newFill := range newMakerFills {
		matchedMakerOrder, exists := matchedMakerOrderIdToOrder[newFill.MakerOrderId]
		if !exists {
			panic(
				fmt.Sprintf(
					"mustUpdateMemclobStateWithMatches: matched maker order %+v does not exist in `matchedMakerOrderIdToOrder`",
					matchedMakerOrder,
				),
			)
		}

		makerFillWithOrders = append(
			makerFillWithOrders,
			types.MakerFillWithOrder{
				Order:     matchedMakerOrder,
				MakerFill: newFill,
			},
		)

		// Skip adding order placement in the operations to propose if it already exists.
		isMakerPreexistingStatefulOrder := matchedMakerOrder.IsStatefulOrder() &&
			m.pendingFills.operationsToPropose.IsMakerOrderPreexistingStatefulOrder(
				matchedMakerOrder,
			)
		if m.pendingFills.isMakerOrderInOperationsToPropose(
			ctx,
			matchedMakerOrder,
			isMakerPreexistingStatefulOrder,
		) {
			continue
		}

		// Add the maker order placement to the operations to propose.
		m.pendingFills.mustAddOrderToOperationsToPropose(
			ctx,
			matchedMakerOrder,
			isMakerPreexistingStatefulOrder,
		)
	}

	// If the taker order is not a liquidation, add the taker order placement to the operations queue.
	if !takerOrder.IsLiquidation() {
		taker := takerOrder.MustGetOrder()

		// Assign the taker order a nonce and add it to the operations to propose.
		isTakerPreexistingStatefulOrder := taker.IsStatefulOrder() &&
			m.clobKeeper.DoesStatefulOrderExistInState(ctx, taker)
		m.pendingFills.operationsToPropose.AssignNonceToOrder(taker, isTakerPreexistingStatefulOrder)
		m.pendingFills.mustAddOrderToOperationsToPropose(
			ctx,
			taker,
			isTakerPreexistingStatefulOrder,
		)
	}

	// Add the new matches to the operations queue.
	m.pendingFills.operationsToPropose.AddMatchToOperationsQueue(takerOrder, makerFillWithOrders)

	// Update the memclob fields for match bookkeeping with the new matches.

	// Define a data-structure for storing the total number of matched quantums for each subaccount
	// in the matching loop. This is used for reduce-only logic.
	subaccountTotalMatchedQuantums := make(map[satypes.SubaccountId]*big.Int)

	// TODO(DEC-1546): Refactor this to iterate over `newMakerFills`.
	for _, newFill := range newPendingFills {
		// Define shared variables.
		matchedQuantums := newFill.Quantums
		matchedMakerOrder := m.pendingFills.mustGetMatchedOrderByHash(ctx, newFill.MakerOrderHash)
		matchedTakerOrder := m.pendingFills.mustGetMatchedOrderByHash(ctx, newFill.TakerOrderHash)

		// Sanity checks.
		if matchedQuantums == 0 {
			panic(fmt.Sprintf(
				"mustUpdateMemclobStateWithMatches: Fill has 0 quantums. Fill %v and maker order %v",
				newFill,
				matchedMakerOrder,
			))
		}

		// Update `fills` with the new fill.
		m.pendingFills.fills = append(m.pendingFills.fills, newFill)

		// Update the orderbook state to reflect the maker order was matched.
		makerOrder := matchedMakerOrder.MustGetOrder()
		matchOffchainUpdates := m.mustUpdateOrderbookStateWithMatchedMakerOrder(
			ctx,
			makerOrder,
		)
		offchainUpdates.BulkUpdate(matchOffchainUpdates)

		// Update the total matched quantums for this matching loop stored in `subaccountTotalMatchedQuantums`.
		for _, order := range []types.MatchableOrder{
			matchedTakerOrder,
			matchedMakerOrder,
		} {
			bigTotalMatchedQuantums, exists := subaccountTotalMatchedQuantums[order.GetSubaccountId()]
			if !exists {
				bigTotalMatchedQuantums = big.NewInt(0)
			}

			bigMatchedQuantums := matchedQuantums.ToBigInt()
			if order.IsBuy() {
				bigTotalMatchedQuantums = bigTotalMatchedQuantums.Add(bigTotalMatchedQuantums, bigMatchedQuantums)
			} else {
				bigTotalMatchedQuantums = bigTotalMatchedQuantums.Sub(bigTotalMatchedQuantums, bigMatchedQuantums)
			}

			subaccountTotalMatchedQuantums[order.GetSubaccountId()] = bigTotalMatchedQuantums
		}
	}

	// Build a slice of all subaccounts which had matches this matching loop, and sort them for determinism.
	allSubaccounts := lib.ConvertMapToSliceOfKeys(subaccountTotalMatchedQuantums)
	sort.Sort(satypes.SortedSubaccountIds(allSubaccounts))

	// For each subaccount that had a match in the matching loop, determine whether we should cancel
	// open reduce-only orders for the subaccount. This occurs when the sign of the position size before matching
	// differs from the sign of the position size after matching.
	for _, subaccountId := range allSubaccounts {
		cancelledOffchainUpdates := m.maybeCancelReduceOnlyOrders(
			ctx,
			subaccountId,
			takerOrder.GetClobPairId(),
			subaccountTotalMatchedQuantums[subaccountId],
		)

		offchainUpdates.BulkUpdate(cancelledOffchainUpdates)
	}

	return offchainUpdates
}

// GetOperations fetches the operations to propose in the next block.
func (m *MemClobPriceTimePriority) GetOperations(ctx sdk.Context) (
	operationsQueue []types.Operation,
) {
	return m.pendingFills.operationsToPropose.GetOperationsQueue()
}

// GetOrdersWithAddToOrderbookCollatCheck fetches all order hashes that had
// the add-to-orderbook collateralization check performed in the last block.
func (m *MemClobPriceTimePriority) GetOrdersWithAddToOrderbookCollatCheck(ctx sdk.Context) (
	ordersWithAddToOrderbookCollatCheck []types.OrderHash,
) {
	return m.pendingFills.operationsToPropose.GetOrdersWithAddToOrderbookCollatCheck()
}

// PlaceOrder will perform the following operations:
// - Validate the order against memclob in-memory state.
// - If the newly placed order causes an overlap, match orders within that orderbook.
//   - Note that if any maker orders fail collateralization checks they will be removed, and if the taker order fails
//     collateralization checks then matching will stop.
//   - If there were any matches resulting from matching the taker order, memclob state will be updated accordingly.
//   - If the order has nonzero remaining size after it has been matched and passes collateralization checks, the order
//     will be added to the orderbook and other bookkeeping datastructures.
//
// This function will return the amount of optimistically filled size (in base quantums) of the order that was filled
// while attempting to match the taker order against the book, along with the status of the placed order.
// If order validation failed, no state in the memclob will be modified and an error will be returned.
func (m *MemClobPriceTimePriority) PlaceOrder(
	ctx sdk.Context,
	order types.Order,
	performAddToOrderbookCollatCheck bool,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	// Perform invariant checks that the orderbook is not crossed after `PlaceOrder` finishes execution.
	defer func() {
		orderbook := m.openOrders.mustGetOrderbook(ctx, order.GetClobPairId())
		bestBid, hasBid := m.openOrders.getBestOrderOnSide(
			orderbook,
			true, // isBuy
		)
		bestAsk, hasAsk := m.openOrders.getBestOrderOnSide(
			orderbook,
			false, // isBuy
		)
		if hasBid && hasAsk && bestBid.Value.Order.Subticks >= bestAsk.Value.Order.Subticks {
			panic(
				fmt.Sprintf(
					"PlaceOrder: orderbook ID %v is crossed. Best bid: (%+v), best ask: (%+v), placed order: (%+v)",
					order.GetClobPairId(),
					bestBid.Value.Order,
					bestAsk.Value.Order,
					order,
				),
			)
		}
	}()

	offchainUpdates = types.NewOffchainUpdates()

	// Validate the order and return an error if any validation fails.
	if err := m.validateNewOrder(ctx, order); err != nil {
		return 0, 0, offchainUpdates, err
	}

	if m.generateOffchainUpdates {
		if message, success := off_chain_updates.CreateOrderPlaceMessage(
			ctx.Logger(),
			order,
		); success {
			offchainUpdates.AddPlaceMessage(order.OrderId, message)
		}
	}

	// If we should perform the add-to-orderbook collateralization check for this order, mark that
	// we should perform the add-to-orderbook collateralization check on this order in the next
	// proposed block. Note this is necessary for the `DeliverTx` validation flow to know which
	// orders should not skip the add-to-orderbook collateralization check.
	if performAddToOrderbookCollatCheck {
		m.pendingFills.operationsToPropose.MustAddToOrderbookCollatCheckOrders(order)
	}

	// Attempt to match the order against the orderbook.
	takerOrderStatus, takerOffchainUpdates, err := m.matchOrder(ctx, &order)
	offchainUpdates.BulkUpdate(takerOffchainUpdates)

	if err != nil {
		if m.generateOffchainUpdates {
			// Send an off-chain update message indicating the order should be removed from the orderbook
			// on the Indexer.
			if message, success := off_chain_updates.CreateOrderRemoveMessage(
				ctx.Logger(),
				order.OrderId,
				takerOrderStatus.OrderStatus,
				err,
				off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(order.OrderId, message)
			}
		}

		return 0, 0, offchainUpdates, err
	}

	remainingSize := takerOrderStatus.RemainingQuantums
	orderSizeOptimisticallyFilledFromMatchingQuantums = takerOrderStatus.OrderOptimisticallyFilledQuantums

	// If the status of the taker order is not successful, do not attempt to add the order to the orderbook.
	if !takerOrderStatus.OrderStatus.IsSuccess() {
		if m.generateOffchainUpdates {
			// Send an off-chain update message indicating the order should be removed from the orderbook
			// on the Indexer.
			if message, success := off_chain_updates.CreateOrderRemoveMessage(
				ctx.Logger(),
				order.OrderId,
				takerOrderStatus.OrderStatus,
				nil,
				off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(order.OrderId, message)
			}
		}
		return orderSizeOptimisticallyFilledFromMatchingQuantums, takerOrderStatus.OrderStatus, offchainUpdates, nil
	}

	// If the order has no remaining size, we do not have to add the order to the orderbook and we can return early.
	if remainingSize == 0 {
		// If the status of the taker order after matching is success and the order has no remaining size, send an
		// off-chain message with the total filled size of the order equal to the size of the order.
		// This is needed to account for the case where an order was partially matched, rewound, then was fully matched
		// during uncrossing.
		if m.generateOffchainUpdates {
			if message, success := off_chain_updates.CreateOrderUpdateMessage(
				ctx.Logger(),
				order.OrderId,
				order.GetBaseQuantums(),
			); success {
				offchainUpdates.AddUpdateMessage(order.OrderId, message)
			}
		}
		return orderSizeOptimisticallyFilledFromMatchingQuantums, takerOrderStatus.OrderStatus, offchainUpdates, nil
	}

	// If this is an IOC order, cancel the remaining size since IOC orders cannot be maker orders.
	if order.RequiresImmediateExecution() {
		orderStatus := types.ImmediateOrCancelWouldRestOnBook
		if m.generateOffchainUpdates {
			// Send an off-chain update message indicating the order should be removed from the orderbook
			// on the Indexer.
			if message, success := off_chain_updates.CreateOrderRemoveMessage(
				ctx.Logger(),
				order.OrderId,
				orderStatus,
				nil,
				off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(order.OrderId, message)
			}
		}
		return orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, nil
	}

	// The taker order has unfilled size which will be added to the orderbook as a maker order.
	// Verify the maker order can be added to the orderbook by performing the add-to-orderbook collateralization
	// check if necessary.
	addOrderOrderStatus := types.Success
	if performAddToOrderbookCollatCheck {
		addOrderOrderStatus = m.addOrderToOrderbookCollateralizationCheck(
			ctx,
			order,
		)
	}

	// If the add order to orderbook collateralization check failed, we cannot add the order to the orderbook.
	if !addOrderOrderStatus.IsSuccess() {
		if m.generateOffchainUpdates {
			// Send an off-chain update message indicating the order should be removed from the orderbook
			// on the Indexer.
			if message, success := off_chain_updates.CreateOrderRemoveMessage(
				ctx.Logger(),
				order.OrderId,
				addOrderOrderStatus,
				nil,
				off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(order.OrderId, message)
			}
		}
		return orderSizeOptimisticallyFilledFromMatchingQuantums, addOrderOrderStatus, offchainUpdates, nil
	}

	// If the order hasn't been assigned a nonce, then assign a nonce to the new maker
	// order since it passed all validation.
	// Note this order could have been assigned a nonce if it was partially matched or caused order
	// removals of orders in the operations queue, meaning it was added to the operations to propose.
	isMakerPreexistingStatefulOrder := order.IsStatefulOrder() &&
		m.clobKeeper.DoesStatefulOrderExistInState(ctx, order)
	var operation types.Operation
	if isMakerPreexistingStatefulOrder {
		operation = types.NewPreexistingStatefulOrderPlacementOperation(order)
	} else {
		operation = types.NewOrderPlacementOperation(order)
	}

	if !m.pendingFills.operationsToPropose.DoesOperationHaveNonce(operation) {
		m.pendingFills.operationsToPropose.AssignNonceToOrder(
			order,
			isMakerPreexistingStatefulOrder,
		)
	}

	// Ensure newly-placed stateful orders are always added to the operations queue.
	// Note it could already be present in the operations queue if it was partially matched or
	// caused order removals of orders in the operations queue.
	if order.IsStatefulOrder() &&
		!isMakerPreexistingStatefulOrder &&
		!m.pendingFills.operationsToPropose.IsOrderPlacementInOperationsQueue(order) {
		m.pendingFills.mustAddOrderToOperationsToPropose(ctx, order, isMakerPreexistingStatefulOrder)
	}

	// Add the order to the orderbook and all other bookkeeping data structures.
	m.mustAddOrderToOrderbook(ctx, order, false)

	// If the taker order is added to the orderbook successfully, send an off-chain message with
	// the total filled size of the order (size of order - remaining size).
	if m.generateOffchainUpdates {
		if message, success := off_chain_updates.CreateOrderUpdateMessage(
			ctx.Logger(),
			order.OrderId,
			order.GetBaseQuantums()-remainingSize,
		); success {
			offchainUpdates.AddUpdateMessage(order.OrderId, message)
		}
	}

	// TODO(DEC-1347): Ensure emitted stats have tags for which ABCI callback was the caller.
	telemetry.IncrCounterWithLabels(
		[]string{types.ModuleName, metrics.PlaceOrder, metrics.AddedToOrderBook},
		1,
		order.GetOrderLabels(),
	)

	return orderSizeOptimisticallyFilledFromMatchingQuantums, types.Success, offchainUpdates, nil
}

// PlacePerpetualLiquidation matches an IOC liquidation order against the orderbook. Specifically,
// it will perform the following operations:
//   - If the liquidation order overlaps the orderbook, it will match orders within that orderbook
//     until there is no overlap.
//   - Note that if any maker orders fail collateralization checks they will be removed, and it won't
//     perform collateralization checks on the taker order.
//   - If there were any matches resulting from matching the liquidation order, memclob state will
//     be updated accordingly.
func (m *MemClobPriceTimePriority) PlacePerpetualLiquidation(
	ctx sdk.Context,
	order types.LiquidationOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	// Attempt to match the liquidation order against the orderbook.
	// TODO(DEC-1157): Update liquidations flow to send off-chain indexer messages.
	liquidationOrderStatus, offchainUpdates, err := m.matchOrder(ctx, &order)
	return liquidationOrderStatus.OrderOptimisticallyFilledQuantums,
		liquidationOrderStatus.OrderStatus,
		offchainUpdates,
		err
}

// matchOrder will match the provided `MatchableOrder` as a taker order against the respective orderbook.
// This function will return the status of the matched order, along with the new taker pending matches.
func (m *MemClobPriceTimePriority) matchOrder(
	ctx sdk.Context,
	order types.MatchableOrder,
) (
	orderStatus types.TakerOrderStatus,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	offchainUpdates = types.NewOffchainUpdates()

	// Attempt to match the order against the orderbook.
	newPendingFills,
		newMakerFills,
		matchedOrderHashToOrder,
		matchedMakerOrderIdToOrder,
		makerOrdersToRemove,
		takerOrderStatus := m.mustPerformTakerOrderMatching(
		ctx,
		order,
	)

	// If this is a replacement order, then ensure we remove the existing order from the orderbook.
	if !order.IsLiquidation() {
		orderId := order.MustGetOrder().OrderId
		if orderToBeReplaced, found := m.openOrders.getOrder(ctx, orderId); found {
			makerOrdersToRemove = append(makerOrdersToRemove, orderToBeReplaced)
		}
	}

	// If there are any removed maker orders that were included in the operations to propose in the next block,
	// the original taker order must be added to the operationsToPropose.
	// TODO(DEC-1957) After pre-existing stateful order removals are implemented, if there are any
	// pre-existing stateful orders that are removed, the taker order must be added to OTP.
	makerOrderInOperationsToProposeRemoved := false

	// For each maker order that should be removed, remove it from the orderbook and determine
	// whether any maker orders in `operationsToPropose` were removed. This is necessary for determining
	// whether order placements that didn't generate matches should be added to `operationsToPropose`.
	for _, makerOrder := range makerOrdersToRemove {
		// TODO(DEC-847): Update logic to properly remove long-term orders.
		makerOrderId := makerOrder.OrderId
		if m.generateOffchainUpdates {
			// If the taker order and the removed maker order are from the same subaccount, set
			// the reason to SELF_TRADE error, otherwise set the reason to be UNDERCOLLATERALIZED.
			// TODO(DEC-1409): Update this to support order replacements on indexer.
			reason := off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_UNDERCOLLATERALIZED
			if order.GetSubaccountId() == makerOrderId.SubaccountId {
				reason = off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_SELF_TRADE_ERROR
			}

			// Send an off-chain update message to the indexer to remove the maker orders that failed
			// collateralization checks from the orderbook.
			if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
				ctx.Logger(),
				makerOrderId,
				reason,
				off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(makerOrderId, message)
			}
		}

		// TODO(DEC-1945) Remove nonces for removed orders that aren't in the operation queue.
		m.mustRemoveOrder(ctx, makerOrderId)

		// To determine if we need to include a taker order placement in the `operationToPropose` even if it
		// doesn't match, check if any of the removed maker orders are in `operationsToPropose`. Note this is
		// necessary for the `DeliverTx` flow validation to work properly.
		isMakerPreexistingStatefulOrder := makerOrder.IsStatefulOrder() &&
			m.pendingFills.operationsToPropose.IsMakerOrderPreexistingStatefulOrder(
				makerOrder,
			)
		if m.pendingFills.isMakerOrderInOperationsToPropose(ctx, makerOrder, isMakerPreexistingStatefulOrder) {
			makerOrderInOperationsToProposeRemoved = true
		} else if isMakerPreexistingStatefulOrder {
			// TODO(CLOB-237): Handle stateful order removals properly here.
			m.pendingFills.operationsToPropose.RemovePreexistingStatefulOrderPlacementNonce(makerOrder)
		} else {
			m.pendingFills.operationsToPropose.RemoveOrderPlacementNonce(makerOrder)
		}
	}

	var matchingErr error

	// If this is a fill-or-kill order and it wasn't fully filled, set the matching error
	// so that the order is canceled. We check the taker order status since the order may
	// have been unsuccessful for other reasons which would be the true root cause
	// (e.g. failed a collateralization check).
	if !order.IsLiquidation() &&
		order.MustGetOrder().TimeInForce == types.Order_TIME_IN_FORCE_FILL_OR_KILL &&
		takerOrderStatus.RemainingQuantums > 0 {
		// FOK orders _must_ return an error here if they are not fully filled regardless
		// of the reason why they were not fully filled. If an error is not returned here, then
		// any partial matches that occurred during matching will be committed to state, and included
		// in the operations queue. This violates the invariant of FOK orders that they must be fully
		// filled or not filled at all.
		// TODO(CLOB-267): Create more granular error types here that indicate why the order was not
		// fully filled (i.e. undercollateralized, reduce only resized, etc).
		matchingErr = types.ErrFokOrderCouldNotBeFullyFilled
	}

	// If the order is post only and it's not the rewind step, then it cannot be filled.
	// Set the matching error so that the order is canceled.
	// TODO (DEC-998): Determine if allowing post-only orders to match in rewind step is valid.
	if len(newPendingFills) > 0 &&
		!order.IsLiquidation() &&
		order.MustGetOrder().TimeInForce == types.Order_TIME_IN_FORCE_POST_ONLY {
		matchingErr = types.ErrPostOnlyWouldCrossMakerOrder
	}

	takerGeneratedValidMatches := len(newMakerFills) > 0 && matchingErr == nil

	// If the match is valid and placing the taker order generated valid matches, update memclob state.
	if takerGeneratedValidMatches {
		matchOffchainUpdates := m.mustUpdateMemclobStateWithMatches(
			ctx,
			order,
			newPendingFills,
			newMakerFills,
			matchedOrderHashToOrder,
			matchedMakerOrderIdToOrder,
		)
		offchainUpdates.BulkUpdate(matchOffchainUpdates)
	}

	// If there were maker orders in the operations to propose in the next block that were removed and
	// the taker order was a post-only order that crossed the book and will be dropped, add the maker order
	// which crossed the post-only order to the operations queue to ensure the post-only order is dropped when
	// validating the block.
	// TODO(DEC-1542): Research downsides of including order placements that don't match but cause
	// order removals in the operations queue.
	if makerOrderInOperationsToProposeRemoved && errors.Is(matchingErr, types.ErrPostOnlyWouldCrossMakerOrder) {
		// Since the error is `ErrPostOnlyWouldCrossMakerOrder`, we know at least one element must
		// exist in `newMakerFills`.
		firstMatchedMakerOrderId := newMakerFills[0].MakerOrderId
		matchedMakerOrder := matchedMakerOrderIdToOrder[firstMatchedMakerOrderId]

		// Ensure the maker order is included in the operations to propose.
		isPreexistingStatefulOrder := matchedMakerOrder.IsStatefulOrder() &&
			m.pendingFills.operationsToPropose.IsMakerOrderPreexistingStatefulOrder(matchedMakerOrder)
		isMakerOrderInOperationsToPropose := m.pendingFills.isMakerOrderInOperationsToPropose(
			ctx,
			matchedMakerOrder,
			isPreexistingStatefulOrder,
		)
		if !isMakerOrderInOperationsToPropose {
			m.pendingFills.mustAddOrderToOperationsToPropose(
				ctx,
				matchedMakerOrder,
				isPreexistingStatefulOrder,
			)
		}
	}

	// If there were maker orders in the operations to propose in the next block that were removed and the
	// taker order was not placed in the operations queue, add the taker order placement to the operations
	// queue to ensure that the `DeliverTx` flow memclob validation works properly.
	// TODO(DEC-1542): Research downsides of including order placements that don't match but cause
	// order removals in the operations queue.
	if makerOrderInOperationsToProposeRemoved && !takerGeneratedValidMatches {
		if order.IsLiquidation() {
			m.pendingFills.operationsToPropose.AddMatchToOperationsQueue(order, []types.MakerFillWithOrder{})
		} else {
			order := order.MustGetOrder()

			// Assign the taker order a nonce and add it to the operations to propose.
			isTakerPreexistingStatefulorder := order.IsStatefulOrder() &&
				m.clobKeeper.DoesStatefulOrderExistInState(ctx, order)
			m.pendingFills.operationsToPropose.AssignNonceToOrder(order, isTakerPreexistingStatefulorder)
			m.pendingFills.mustAddOrderToOperationsToPropose(
				ctx,
				order,
				isTakerPreexistingStatefulorder,
			)
		}
	}

	return takerOrderStatus, offchainUpdates, matchingErr
}

// GetClobPairForPerpetual gets the first CLOB pair ID associated with the provided perpetual ID.
// It returns an error if there are no CLOB pair IDs associated with the perpetual ID.
func (m *MemClobPriceTimePriority) GetClobPairForPerpetual(
	ctx sdk.Context,
	perpetualId uint32,
) (
	clobPairId types.ClobPairId,
	err error,
) {
	clobPairIds, exists := m.perpetualIdToClobPairId[perpetualId]
	if !exists {
		return 0, sdkerrors.Wrapf(
			types.ErrNoClobPairForPerpetual,
			"Perpetual ID %d has no associated CLOB pairs",
			perpetualId,
		)
	}

	if len(clobPairIds) == 0 {
		panic("GetClobPairForPerpetual: Perpetual ID was created without a CLOB pair ID.")
	}

	return clobPairIds[0], nil
}

// ReplayOperations will replay the provided operations onto the memclob.
// This is used to replay all local operations from the local `operationsToPropose` from the previous block.
// The following operations are supported:
// - Short-Term orders.
// - Newly-placed stateful orders.
// - Pre-existing stateful orders.
// - Stateful cancelations.
// Note that match operations are no-op.
func (m *MemClobPriceTimePriority) ReplayOperations(
	ctx sdk.Context,
	localOperationsQueue []types.Operation,
	existingOffchainUpdates *types.OffchainUpdates,
	canceledStatefulOrderIds []types.OrderId,
) *types.OffchainUpdates {
	// Recover from any panics that occur during replay operations.
	// This could happen in cases where i.e. A subaccount balance overflowed
	// during a match. We don't want to halt the entire chain in this case.
	// TODO(CLOB-275): Do not gracefully handle panics in `PrepareCheckState`.
	defer func() {
		if r := recover(); r != nil {
			ctx.Logger().Error("panic in replay operations", "panic", r)
		}
	}()

	// Iterate over all provided operations.
	for _, operation := range localOperationsQueue {
		switch operation.Operation.(type) {
		// Replay all short-term and stateful order placements.
		case *types.Operation_OrderPlacement:
			order := operation.GetOrderPlacement().Order

			// TODO(CLOB-269): Do not replay orders that were canceled in the previous block.
			// This should no longer be necessary once transaction validation is implemented,
			// as the stateful order placement being rewound here should fail sequence number
			// validation.
			if lib.ContainsValue(canceledStatefulOrderIds, order.OrderId) {
				continue
			}

			// Branch the state to avoid writing to state on failed operations.
			placeOrderCtx, writeCache := ctx.CacheContext()

			// Note we use `clobKeeper.PlaceOrder` here to ensure the proper stateful validation is performed and
			// newly-placed stateful orders are written to state. In the future this will be important for sequence number
			// verification as well.
			// TODO(DEC-1755): Account for sequence number verification.
			// TODO(DEC-998): Research whether it's fine for two post-only orders to be matched. Currently they are dropped.
			msg := types.NewMsgPlaceOrder(order)
			orderSizeOptimisticallyFilledFromMatchingQuantums,
				orderStatus, placeOrderOffchainUpdates, err := m.clobKeeper.ReplayPlaceOrder(
				placeOrderCtx,
				msg,
			)

			ctx.Logger().Info(
				"Received new order",
				"orderHash",
				fmt.Sprintf("%X", order.GetOrderHash()),
				"msg",
				msg,
				"status",
				orderStatus,
				"orderSizeOptimisticallyFilledFromMatchingQuantums",
				orderSizeOptimisticallyFilledFromMatchingQuantums,
				"err",
				err,
				"block",
				ctx.BlockHeight(),
			)

			if err != nil {
				ctx.Logger().Debug(
					"ReplayOperations: PlaceOrder() returned an error.",
					"error",
					err,
					"operation",
					operation,
					"order",
					order,
				)

				// If the order is dropped while adding it to the book, return an off-chain order remove
				// message for the order.
				// Note: Currently, the error returned from placing the order determines whether an order
				// removal message is sent to the Indexer. This may change later on to be a check on whether
				// the order has an existing nonce.
				if m.generateOffchainUpdates && off_chain_updates.ShouldSendOrderRemovalOnReplay(err) {
					if message, success := off_chain_updates.CreateOrderRemoveMessageWithDefaultReason(
						ctx.Logger(),
						order.OrderId,
						orderStatus,
						err,
						off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
						off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_INTERNAL_ERROR,
					); success {
						existingOffchainUpdates.AddRemoveMessage(order.OrderId, message)
					}
				}
			} else {
				writeCache()

				if m.generateOffchainUpdates {
					// None of the `PlaceMessages` should be needed as the indexer
					// should have already learned about these orders.
					placeOrderOffchainUpdates.ClearPlaceMessages()
					existingOffchainUpdates.BulkUpdate(placeOrderOffchainUpdates)
				}
			}

		// Replay all pre-existing stateful order placements.
		case *types.Operation_PreexistingStatefulOrder:
			orderId := operation.GetPreexistingStatefulOrder()

			// TODO(DEC-1974): The `PreexistingStatefulOrder` operation
			// does not contain the order hash, so we cannot check if the
			// order is the same as the one in the book (rather than a replacement).
			// For consistency we should fix this, but currently it is not an issue as
			// replacements are not currently supported, and trying to place an older version
			// of a stateful order should fail due to the `GoodTilBlockTime`.
			statefulOrderPlacement, found := m.clobKeeper.GetStatefulOrderPlacement(ctx, *orderId)
			if !found {
				// It's possible that this order was canceled or expired in the last committed block.
				continue
			}

			// Branch the state to avoid writing to state on failed operations.
			placeOrderCtx, writeCache := ctx.CacheContext()

			// Note that we use `memclob.PlaceOrder` here, this will skip writing the stateful order placement to state.
			// TODO(DEC-998): Research whether it's fine for two post-only orders to be matched. Currently they are dropped.
			_, orderStatus, placeOrderOffchainUpdates, err := m.PlaceOrder(
				placeOrderCtx,
				statefulOrderPlacement.Order,
				false,
			)
			if err != nil {
				ctx.Logger().Debug(
					"ReplayOperations: PlaceOrder() returned an error for a pre-existing stateful order.",
					"error",
					err,
					"operation",
					operation,
					"statefulOrderPlacement",
					statefulOrderPlacement,
				)

				// If the stateful order is dropped while adding it to the book, return an off-chain order remove
				// message for the order.
				// Note: Currently, the error returned from placing the order determines whether an order
				// removal message is sent to the Indexer. This may change later on to be a check on whether
				// the order has an existing nonce.
				if m.generateOffchainUpdates && off_chain_updates.ShouldSendOrderRemovalOnReplay(err) {
					if message, success := off_chain_updates.CreateOrderRemoveMessageWithDefaultReason(
						ctx.Logger(),
						*orderId,
						orderStatus,
						err,
						off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
						off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_INTERNAL_ERROR,
					); success {
						existingOffchainUpdates.AddRemoveMessage(*orderId, message)
					}
				}
			} else {
				writeCache()

				if m.generateOffchainUpdates {
					// None of the `PlaceMessages` should be needed as the indexer
					// should have already learned about these orders.
					placeOrderOffchainUpdates.ClearPlaceMessages()
					existingOffchainUpdates.BulkUpdate(placeOrderOffchainUpdates)
				}
			}

		// Replay stateful order cancelations.
		case *types.Operation_OrderCancellation:
			orderIdToCancel := operation.GetOrderCancellation().GetOrderId()
			goodTilBlockTime := operation.GetOrderCancellation().GetGoodTilBlockTime()
			if orderIdToCancel.IsShortTermOrder() {
				// Skip replaying Short-Term cancellations since the
				// latest cancellation information is already
				// stored in the memclob.
				continue
			}

			// We don't need to fork the state for order cancelations as they do not leave the
			// store in a dirty state upon failure, however doing so would not hurt so we could consider
			// doing it for consistency in the future.
			err := m.clobKeeper.CheckTxCancelOrder(
				ctx,
				types.NewMsgCancelOrderStateful(orderIdToCancel, goodTilBlockTime),
			)
			if err != nil {
				ctx.Logger().Debug(
					"ReplayOperations: CancelOrder() returned an error.",
					"error",
					err,
					"operation",
					operation,
					"orderIdToCancel",
					orderIdToCancel,
					"goodTilBlockTime",
					goodTilBlockTime,
				)
			}
		// Matches are no-op.
		case *types.Operation_Match:
		default:
			panic(fmt.Sprintf("unknown operation type: %T", operation.Operation))
		}
	}

	return existingOffchainUpdates
}

// RemoveAndClearOperationsQueue is called during `Commit`/`PrepareCheckState`
// to clear and remove all orders that exist in the provided local validator's OTP (`operationsToPropose`).
// It performs the following steps:
// 1. Copy the operations queue.
// 2. Clear the OTP. Note that this also removes nonces for every operation in the OTP.
// 3. For each order placement operation in the copy, remove the order from the book if it exists.
func (m *MemClobPriceTimePriority) RemoveAndClearOperationsQueue(
	ctx sdk.Context,
	localValidatorOperationsQueue []types.Operation,
) {
	// Clear the OTP. This will also remove nonces for every operation in `operationsQueueCopy`.
	m.pendingFills.operationsToPropose.ClearOperationsQueue()

	// For each order placement operation in the copy, remove the order from the book
	// if it exists.
	for _, operation := range localValidatorOperationsQueue {
		switch operation.Operation.(type) {
		case *types.Operation_OrderPlacement:
			otpOrderId := operation.GetOrderPlacement().Order.OrderId
			otpOrderHash := operation.GetOrderPlacement().Order.GetOrderHash()

			// If the order exists in the book, remove it.
			existingOrder, found := m.openOrders.getOrder(ctx, otpOrderId)
			if found && existingOrder.GetOrderHash() == otpOrderHash {
				m.mustRemoveOrder(ctx, otpOrderId)
			}
		case *types.Operation_PreexistingStatefulOrder:
			otpOrderId := operation.GetPreexistingStatefulOrder()

			// TODO(DEC-1974): The `PreexistingStatefulOrder` operation
			// does not contain the order hash, so we cannot check if the
			// order is the same as the one in the book (rather than a replacement).
			// For consistency we should fix this, but currently it is not an issue as
			// we would expect the replacement to always be included in the
			// OTP, and therefore be removed in this loop as well.
			if m.openOrders.hasOrder(ctx, *otpOrderId) {
				m.mustRemoveOrder(ctx, *otpOrderId)
			}
		}
	}

	// TODO(DEC-1545): Remove all match data structures after operations queue refactor.
	m.pendingFills.fills = make([]types.PendingFill, 0)
	m.pendingFills.orderHashToMatchableOrder = make(map[types.OrderHash]types.MatchableOrder)
	m.pendingFills.orderIdToOrder = make(map[types.OrderId]types.Order)
}

// PurgeInvalidMemclobState will purge the following invalid state from the memclob:
// - Expired Short-Term order cancellations.
// - Expired Short-Term and stateful orders from the orderbook and remove their nonces from `OperationsToPropose`.
// - Fully-filled orders from the orderbook and remove their nonces from `OperationsToPropose`.
// - Canceled stateful orders and remove their nonces from `OperationsToPropose`.
func (m *MemClobPriceTimePriority) PurgeInvalidMemclobState(
	ctx sdk.Context,
	fullyFilledOrderIds []types.OrderId,
	expiredStatefulOrderIds []types.OrderId,
	canceledStatefulOrderIds []types.OrderId,
	existingOffchainUpdates *types.OffchainUpdates,
) *types.OffchainUpdates {
	blockHeight := lib.MustConvertIntegerToUint32(ctx.BlockHeight())

	// Remove all fully-filled order IDs from the memclob if they exist.
	for _, orderId := range fullyFilledOrderIds {
		m.RemoveOrderIfFilled(ctx, orderId)
	}

	// Remove all canceled stateful order IDs from the memclob if they exist.
	// If the slice has non-stateful order IDs or contains duplicates, panic.
	if lib.ContainsDuplicates(canceledStatefulOrderIds) {
		panic(
			fmt.Sprintf(
				"PurgeInvalidMemclobState: called with canceledStatefulOrderIds slice %v which contains duplicate order IDs",
				canceledStatefulOrderIds,
			),
		)
	}

	for _, statefulOrderId := range canceledStatefulOrderIds {
		statefulOrderId.MustBeStatefulOrder()

		// TODO(DEC-798/DEC-1279): Update this logic once we've determined how to rewind `MsgRemoveOrder` messages.
		// Currently stateful orders can be removed from the book due to things such as collateralization
		// check failures, self-trade errors, etc and will not be removed from state. Therefore it
		// is possible that when they are canceled they will not exist on the orderbook.
		if m.openOrders.hasOrder(ctx, statefulOrderId) {
			statefulOrder, _ := m.GetOrder(ctx, statefulOrderId)
			m.mustRemoveOrder(ctx, statefulOrderId)
			m.pendingFills.operationsToPropose.RemovePreexistingStatefulOrderPlacementNonce(statefulOrder)
		}
	}

	// Remove all expired stateful order IDs from the memclob if they exist.
	// If the slice has non-stateful order IDs or contains duplicates, panic.
	if lib.ContainsDuplicates(expiredStatefulOrderIds) {
		panic(
			fmt.Sprintf(
				"PurgeInvalidMemclobState: called with expiredStatefulOrderIds slice %v which contains duplicate order IDs",
				expiredStatefulOrderIds,
			),
		)
	}

	for _, statefulOrderId := range expiredStatefulOrderIds {
		statefulOrderId.MustBeStatefulOrder()

		// TODO(DEC-1800): Ensure correct indexer updates are returned for expired stateful orders.
		// TODO(DEC-798/DEC-1279): Update this logic once we've determined how to rewind `MsgRemoveOrder` messages.
		// Currently stateful orders can be removed from the book due to things such as collateralization
		// check failures, self-trade errors, etc and will not be removed from state. Therefore it
		// is possible that when they expire they will not exist on the orderbook.
		if m.openOrders.hasOrder(ctx, statefulOrderId) {
			statefulOrder, _ := m.GetOrder(ctx, statefulOrderId)
			m.mustRemoveOrder(ctx, statefulOrderId)
			m.pendingFills.operationsToPropose.RemovePreexistingStatefulOrderPlacementNonce(statefulOrder)

			if m.generateOffchainUpdates {
				// Send an off-chain update message indicating the stateful order should be removed from the
				// orderbook on the Indexer. As the order is expired, the status of the order is canceled
				// and not best-effort-canceled.
				if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
					ctx.Logger(),
					statefulOrderId,
					off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_EXPIRED,
					off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_CANCELED,
				); success {
					existingOffchainUpdates.AddRemoveMessage(statefulOrderId, message)
				}
			}
		}
	}

	// Remove all expired Short-Term order IDs from the memclob.
	if blockExpirations, beExists := m.openOrders.blockExpirationsForOrders[blockHeight]; beExists {
		for shortTermOrderId := range blockExpirations {
			if m.generateOffchainUpdates {
				// Send an off-chain update message indicating the order should be removed from the
				// orderbook on the Indexer. As the order is expired, the status of the order is canceled
				// and not best-effort-canceled.
				if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
					ctx.Logger(),
					shortTermOrderId,
					off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_EXPIRED,
					off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_CANCELED,
				); success {
					existingOffchainUpdates.AddRemoveMessage(shortTermOrderId, message)
				}
			}
			shortTermOrder, _ := m.GetOrder(ctx, shortTermOrderId)
			m.mustRemoveOrder(ctx, shortTermOrderId)
			m.pendingFills.operationsToPropose.RemoveOrderPlacementNonce(shortTermOrder)
		}
	}

	// Remove expired cancels. Note we don't have to remove a nonce for Short-Term order cancellations
	// since they will be removed in `ClearOperationsQueue`.
	m.cancels.removeAllAtBlock(blockHeight)

	return existingOffchainUpdates
}

// validateNewOrder will perform the following validation against the memclob's in-memory state to ensure the order
// can be placed (and if any condition is false, an error will be returned):
//   - The order is not canceled (with an equal-to-or-greater-than `GoodTilBlock` than the new order).
//   - If the order is replacing another order, then the new order's expiration must not be less than the
//     existing order's expiration.
//   - This subaccount has strictly less than `MaxSubaccountOrdersPerClobAndSide` open orders on the new order's
//     CLOB and side.
//
// Note that it does not perform collateralization checks since that will be done when matching the order (if the order
// overlaps the book) and when adding the order to the book (if the order has remaining size after matching).
//
// This function does not perform any order validation that doesn't depend on the memclob's in-memory state.
// Specifically, it will assume the following is true:
// - `Order.OrderId` is a valid `OrderId`.
// - The `Order.ClobPairId` references a valid `ClobPair`.
// - The order is not expired (`Order.GoodTilBlock >= currentBlock`).
// - The order expiration is within `ShortBlockWindow` (`Order.GoodTilBlock <= currentBlock + ShortBlockWindow`).
// - This order has not already been fully filled.
// - `Order.Side` is a valid side.
// - The order is a valid order for the referenced `ClobPair` (where `Order.ClobPairId == ClobPair.Id`). Specifically:
//   - `Order.Subticks` is a multiple of `ClobPair.SubticksPerTick`.
//   - `Order.Quantums` is a multiple of `ClobPair.MinOrderBaseQuantums`.
func (m *MemClobPriceTimePriority) validateNewOrder(
	ctx sdk.Context,
	order types.Order,
) (
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.PlaceOrder,
		metrics.Memclob,
		metrics.ValidateOrder,
		metrics.Latency,
	)
	orderId := order.OrderId

	if orderId.IsShortTermOrder() {
		// If the cancelation has an equal-to-or-greater `GoodTilBlock` than the new order, return an error.
		// If the cancelation has a lesser `GoodTilBlock` than the new order, we do not remove the cancelation.
		if cancelTilBlock, cancelExists := m.cancels.get(orderId); cancelExists && cancelTilBlock >= order.GetGoodTilBlock() {
			return sdkerrors.Wrapf(
				types.ErrOrderIsCanceled,
				"Order: %+v, Cancellation GoodTilBlock: %d",
				order,
				cancelTilBlock,
			)
		}
	}

	existingRestingOrder, restingOrderExists := m.openOrders.getOrder(ctx, orderId)
	existingMatchedOrder, matchedOrderExists := m.pendingFills.orderIdToOrder[orderId]

	// If an order with the same `OrderId` already exists on the orderbook (or was already matched),
	// then we must validate that the new order's `GoodTilBlock` is greater-in-value than the old order.
	// If greater, then it can be placed (replacing the old order if it was resting on the book).
	// If equal-or-lesser, then it is dropped.
	if restingOrderExists && existingRestingOrder.MustCmpReplacementOrder(&order) >= 0 {
		return types.ErrInvalidReplacement
	}

	if matchedOrderExists && existingMatchedOrder.MustCmpReplacementOrder(&order) >= 0 {
		return types.ErrInvalidReplacement
	}

	// If an order with this `OrderId` does not already exist resting on the book for the same side, then
	// we need to ensure that adding the new order would not cause the subaccount to exceed `MaxOpenOrdersPerClobAndSide`.
	// Note: The order could already be resting on the book for the same side if this order is a replacement.
	// Note: The order could already be resting on the book for a different side if this order is a replacement.
	doesOrderAlreadyExistForSide := restingOrderExists && existingRestingOrder.Side == order.Side
	if !doesOrderAlreadyExistForSide {
		existingSubaccountOrdersForClobAndSide, err := m.GetSubaccountOrders(
			ctx,
			order.GetClobPairId(),
			order.OrderId.SubaccountId,
			order.Side,
		)
		if err != nil {
			// This is an unexpected error and implies prior memclob order validation failed.
			panic(err)
		}

		// Verify that opening this order would not exceed the maximum amount of orders per CLOB and side.
		// This limit is enforced to limit the number of orders that are accounted for in the add to orderbook
		// collateralization check.
		if len(existingSubaccountOrdersForClobAndSide) >= types.MaxSubaccountOrdersPerClobAndSide {
			return sdkerrors.Wrapf(
				types.ErrOrderWouldExceedMaxOpenOrdersPerClobAndSide,
				"order id: %+v",
				order.GetOrderId(),
			)
		}
	}

	// If the order is a reduce-only order, we should ensure it does not increase the subaccount's
	// current position size. Note that we intentionally do not validate that the reduce-only order
	// does not change the subaccount's position _side_, and that will be validated if the order is matched.
	// TODO(DEC-1228): use `MustValidateReduceOnlyOrder` and move this to `PerformStatefulOrderValidation`.
	if order.IsReduceOnly() {
		existingPositionSize := m.clobKeeper.GetStatePosition(ctx, orderId.SubaccountId, order.GetClobPairId())
		orderSize := order.GetBigQuantums()

		// If the reduce-only order is not on the opposite side of the existing position size,
		// cancel the order by returning an error.
		if orderSize.Sign()*existingPositionSize.Sign() != -1 {
			return types.ErrReduceOnlyWouldIncreasePositionSize
		}
	}

	// Check if the order being replaced has at least `MinOrderBaseQuantums` of size remaining, otherwise the order
	// is considered fully filled and cannot be placed/replaced.
	orderbook := m.openOrders.mustGetOrderbook(ctx, order.GetClobPairId())
	remainingAmount, hasRemainingAmount := m.getOrderRemainingAmount(ctx, order)
	if !hasRemainingAmount || remainingAmount < orderbook.MinOrderBaseQuantums {
		return types.ErrOrderFullyFilled
	}

	// Check if the order has already been seen by determining if the associated operation already has a nonce.
	isPreexistingStatefulOrder := order.IsStatefulOrder() &&
		m.clobKeeper.DoesStatefulOrderExistInState(ctx, order)

	var existingOperation types.Operation
	if isPreexistingStatefulOrder {
		existingOperation = types.NewPreexistingStatefulOrderPlacementOperation(order)
	} else {
		existingOperation = types.NewOrderPlacementOperation(order)
	}

	// This order already has a nonce, which means we have already
	// received it. This could happen during `ReplayOperations`, or if a transaction
	// were to fall out of CometBFT's cache and re-submitted.
	// TODO(CLOB-243): Write PlaceOrder test for this case.
	if m.pendingFills.operationsToPropose.DoesOperationHaveNonce(existingOperation) {
		return types.ErrOrderReprocessed
	}

	return nil
}

// addOrderToOrderbookCollateralizationCheck will perform a collateralization check to verify that the subaccount would
// remain collateralized if the new maker order were to be fully filled.
// It returns the result of this collateralization check. If the collateralization check returns an
// error, it will return the collateralization check error so that it can be surfaced to the client.
//
// This function will assume that all prior order validation has passed, including the pre-requisite validation of
// `validateNewOrder` and the actual validation performed within `validateNewOrder`.
// Note that this is a loose check, mainly for the purposes of spam mitigation. We perform an additional
// collateralization check on orders when we attempt to match them.
func (m *MemClobPriceTimePriority) addOrderToOrderbookCollateralizationCheck(
	ctx sdk.Context,
	order types.Order,
) types.OrderStatus {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.PlaceOrder,
		metrics.Memclob,
		metrics.AddToOrderbookCollateralizationCheck,
		metrics.Latency,
	)

	orderId := order.OrderId
	subaccountId := orderId.SubaccountId

	// For the collateralization check, use the remaining amount of the order that is resting on the book.
	remainingAmount, hasRemainingAmount := m.getOrderRemainingAmount(ctx, order)
	if !hasRemainingAmount {
		panic(fmt.Sprintf("addOrderToOrderbookCollateralizationCheck: order has no remaining amount %v", order))
	}

	pendingOpenOrder := types.PendingOpenOrder{
		RemainingQuantums: remainingAmount,
		IsBuy:             order.IsBuy(),
		// This order will be added to the book as a maker order, so it cannot be a taker order.
		IsTaker:    false,
		Subticks:   order.GetOrderSubticks(),
		ClobPairId: order.GetClobPairId(),
	}

	// Temporarily construct the subaccountOpenOrders with a single PendingOpenOrder.
	subaccountOpenOrders := make(map[satypes.SubaccountId][]types.PendingOpenOrder)
	subaccountOpenOrders[subaccountId] = []types.PendingOpenOrder{pendingOpenOrder}

	// TODO(DEC-1896): AddOrderToOrderbookCollatCheck should accept a single PendingOpenOrder as a
	// parameter rather than the subaccountOpenOrders map.
	_, successPerSubaccountUpdate := m.clobKeeper.AddOrderToOrderbookCollatCheck(
		ctx,
		order.GetClobPairId(),
		subaccountOpenOrders,
	)

	return updateResultToOrderStatus(successPerSubaccountUpdate[subaccountId])
}

// mustAddOrderToOrderbook will add the order to the resting orderbook.
// This function will assume that all order validation has already been done.
// If `forceToFrontOfLevel` is true, places the order at the head of the level,
// otherwise places it at the tail.
func (m *MemClobPriceTimePriority) mustAddOrderToOrderbook(
	ctx sdk.Context,
	newOrder types.Order,
	forceToFrontOfLevel bool,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.PlaceOrder,
		metrics.Memclob,
		metrics.AddedToOrderBook,
		metrics.Latency,
	)

	// Ensure that the order is not fully-filled.
	if _, hasRemainingQuantums := m.getOrderRemainingAmount(ctx, newOrder); !hasRemainingQuantums {
		panic(fmt.Sprintf("mustAddOrderToOrderbook: order has no remaining amount %+v", newOrder))
	}

	m.openOrders.mustAddOrderToOrderbook(ctx, newOrder, forceToFrontOfLevel)
}

// mustPerformTakerOrderMatching performs matching using the provided taker order while the order
// overlaps the other side of the orderbook. It returns multiple variables used for representing the result
// of matching with the taker order, which are documented further below. Note that this function does not modify
// memclob state.
//
// This function will assume that all order validation has already been done through the `validateNewOrder` function.
// If `order.Side` is an invalid side or `order.ClobPairId` does not reference a valid CLOB, this function will panic.
func (m *MemClobPriceTimePriority) mustPerformTakerOrderMatching(
	ctx sdk.Context,
	newTakerOrder types.MatchableOrder,
) (
	// A list of pending fills to apply to the memclob's orderbook.
	newPendingFills []types.PendingFill,
	// A slice of new maker fills created from matching this taker order.
	newMakerFills []types.MakerFill,
	// A map of matched order hashes to the order.
	matchedOrderHashToOrder map[types.OrderHash]types.MatchableOrder,
	// A map of matched maker order IDs to the order.
	matchedMakerOrderIdToOrder map[types.OrderId]types.Order,
	// A list of maker orders that failed collateralization checks during matching and should be removed from the
	// orderbook.
	makerOrdersToRemove []types.Order,
	// The status of the taker order, specifically the remaining size, optimistically filled size, and the result of the
	// last collateralization check.
	// This is necessary for determining whether remaining taker order size can be added to the orderbook, and for
	// returning the optimistically filled size to the caller.
	takerOrderStatus types.TakerOrderStatus,
) {
	// Initialize return variables.
	newPendingFills = make([]types.PendingFill, 0)
	newMakerFills = make([]types.MakerFill, 0)
	matchedOrderHashToOrder = make(map[types.OrderHash]types.MatchableOrder)
	matchedMakerOrderIdToOrder = make(map[types.OrderId]types.Order)
	takerOrderStatus.OrderStatus = types.Success
	makerOrdersToRemove = make([]types.Order, 0)

	// Initialize variables used for traversing the orderbook.
	clobPairId := newTakerOrder.GetClobPairId()
	orderbook := m.openOrders.mustGetOrderbook(ctx, clobPairId)
	takerIsBuy := newTakerOrder.IsBuy()
	takerSubaccountId := newTakerOrder.GetSubaccountId()
	takerIsLiquidation := newTakerOrder.IsLiquidation()

	// Store the remaining size of the taker order to determine the filled amount of an order after
	// matching has ended.
	// If the order is a liquidation, then the remaining size is the full size of the order.
	// Else, this is a regular order and might already be partially matched, so we fetch the
	// remaining size of this order.
	var takerRemainingSize satypes.BaseQuantums
	if takerIsLiquidation {
		takerRemainingSize = newTakerOrder.GetBaseQuantums()
	} else {
		var takerHasRemainingSize bool
		takerRemainingSize, takerHasRemainingSize = m.getOrderRemainingAmount(
			ctx,
			newTakerOrder.MustGetOrder(),
		)
		if !takerHasRemainingSize {
			panic(fmt.Sprintf("mustPerformTakerOrderMatching: order has no remaining amount %v", newTakerOrder))
		}
	}
	takerRemainingSizeBeforeMatching := takerRemainingSize

	// Initialize variables used for tracking matches made during this matching cycle.
	var makerLevelOrder *types.LevelOrder
	var takerOrderHash types.OrderHash
	var takerOrderHashWasSet bool
	var bigTotalMatchedAmount *big.Int = big.NewInt(0)

	// Begin attempting to match orders. The below loop performs the following high-level operations, in order:
	// - Find the next best maker order if it exists. If not, stop matching.
	// - Check if the orderbook is crossed. If not, stop matching.
	// - Match the orders and check collateralization. If the maker order failed collateralization, it must be removed
	//   from the book and we can return to step 1 if taker order passed collateralization. If the taker order failed
	//   collateralization, stop matching.
	// - Update local bookkeeping variables with the new match. If the taker order is fully matched, stop matching.
	for {
		var foundMakerOrder bool
		// If the maker level order has not been initialized, then we are just starting matching and need to find the
		// best order on the opposite side.
		// Else, the maker order must have been fully matched (since the taker order has nonzero remaining size), and we
		// need to find the next best maker order.
		if makerLevelOrder == nil {
			makerLevelOrder, foundMakerOrder = m.openOrders.getBestOrderOnSide(orderbook, !takerIsBuy)
		} else {
			makerLevelOrder, foundMakerOrder = m.openOrders.findNextBestLevelOrder(ctx, makerLevelOrder)
		}

		// If no next best maker order was found, stop matching.
		if !foundMakerOrder {
			break
		}

		makerOrder := makerLevelOrder.Value

		// Check if the orderbook is crossed.
		var takerOrderCrossesMakerOrder bool
		if takerIsBuy {
			takerOrderCrossesMakerOrder = newTakerOrder.GetOrderSubticks() >= makerOrder.Order.GetOrderSubticks()
		} else {
			takerOrderCrossesMakerOrder = newTakerOrder.GetOrderSubticks() <= makerOrder.Order.GetOrderSubticks()
		}

		// If the taker order no longer crosses the maker order, stop matching.
		if !takerOrderCrossesMakerOrder {
			break
		}

		makerOrderId := makerOrder.Order.OrderId
		makerSubaccountId := makerOrderId.SubaccountId

		// If the taker order is replacing the maker order, skip this order and continue matching.
		// Note that the maker order will be removed from the orderbook after matching has completed.
		if !takerIsLiquidation && makerOrderId == newTakerOrder.MustGetOrder().OrderId {
			continue
		}

		// If the matched maker order does not have same order ID and is from the same subaccount
		// as the taker order, then we cannot match the orders. Cancel the maker order and continue matching.
		// TODO(DEC-1562): determine if we should handle self-trades differently.
		if makerSubaccountId == takerSubaccountId {
			makerOrdersToRemove = append(makerOrdersToRemove, makerOrder.Order)
			continue
		}

		makerRemainingSize, makerHasRemainingSize := m.getOrderRemainingAmount(ctx, makerOrder.Order)
		if !makerHasRemainingSize {
			panic(fmt.Sprintf("mustPerformTakerOrderMatching: maker order has no remaining amount %v", makerOrder.Order))
		}

		// The matched amount is the minimum of the remaining amount of both orders.
		var matchedAmount satypes.BaseQuantums
		if takerRemainingSize >= makerRemainingSize {
			matchedAmount = makerRemainingSize
		} else {
			matchedAmount = takerRemainingSize
		}

		// For each subaccount involved in the match, if the order is reduce-only we should verify
		// that the position side does not change or increase as a result of matching the order.
		if makerOrder.Order.IsReduceOnly() {
			currentPositionSize := m.clobKeeper.GetStatePosition(ctx, makerSubaccountId, clobPairId)
			resizedMatchAmount := m.pendingFills.resizeReduceOnlyMatchIfNecessary(
				ctx,
				makerSubaccountId,
				clobPairId,
				currentPositionSize,
				matchedAmount,
				!takerIsBuy,
			)

			// If the match size is zero, that indicates the maker order was a reduce-only order that
			// would have increased the maker's position size and we need to find the next best maker
			// order. This can happen if the maker has previous matches within this matching loop
			// that changed their position side, meaning all their resting reduce-only orders are invalid.
			if resizedMatchAmount == 0 {
				// TODO(DEC-1415): Revert this reduce-only bug patch.
				makerOrdersToRemove = append(makerOrdersToRemove, makerOrder.Order)
				continue
			}

			matchedAmount = resizedMatchAmount
		}

		if newTakerOrder.IsReduceOnly() {
			currentPositionSize := m.clobKeeper.GetStatePosition(ctx, takerSubaccountId, clobPairId)
			resizedMatchAmount := m.pendingFills.resizeReduceOnlyMatchIfNecessary(
				ctx,
				takerSubaccountId,
				clobPairId,
				currentPositionSize,
				matchedAmount,
				takerIsBuy,
			)

			// If the taker reduce-only order was resized to 0, that indicates the order is on the
			// same side as the taker's position side and this order should have failed validation.
			if resizedMatchAmount == 0 {
				panic("mustPerformTakerOrderMatching: taker reduce-only order resized to 0")
			}

			matchedAmount = resizedMatchAmount
		}

		subticks := makerOrder.Order.GetOrderSubticks()

		// Perform collateralization checks to verify the orders can be filled.
		matchWithOrders := types.MatchWithOrders{
			TakerOrder: newTakerOrder,
			MakerOrder: &makerOrder.Order,
			FillAmount: matchedAmount,
		}

		success, takerUpdateResult, makerUpdateResult, _, err := m.clobKeeper.ProcessSingleMatch(ctx, matchWithOrders)
		if err != nil {
			if errors.Is(err, types.ErrLiquidationExceedsSubaccountMaxInsuranceLost) {
				// Subaccount has reached max insurance lost block limit. Stop matching.
				telemetry.IncrCounter(1, types.ModuleName, metrics.SubaccountMaxInsuranceLost, metrics.Count)
				break
			}
			if errors.Is(err, types.ErrLiquidationExceedsSubaccountMaxNotionalLiquidated) {
				// Subaccount has reached max notional liquidated block limit. Stop matching.
				telemetry.IncrCounter(1, types.ModuleName, metrics.SubaccountMaxNotionalLiquidated, metrics.Count)
				break
			}
			if errors.Is(err, types.ErrInsuranceFundHasInsufficientFunds) {
				// Deleveraging is required. Stop matching.
				telemetry.IncrCounter(1, types.ModuleName, metrics.LiquidationRequiresDeleveraging, metrics.Count)
				takerOrderStatus.OrderStatus = types.LiquidationRequiresDeleveraging
				break
			}
			if !errors.Is(err, satypes.ErrFailedToUpdateSubaccounts) {
				ctx.Logger().Error("Unexpected error from `ProcessSingleMatch`", "error", err)
				panic(err)
			}
		}

		// If the collateralization check has failed, one or both of the taker or maker orders have failed the
		// collateralization check. Note if the taker is order is liquidation order, only the maker could
		// have failed collateralization checks. Therefore, we must perform the following conditional logic:
		// - If the maker order failed collateralization checks we need to remove it from the orderbook.
		// - If the taker order is not a liquidation order and failed collateralization checks, we
		//   need to stop matching.
		// - If the taker order is a liquidation order or passed collateralization checks, then we
		//   need to continue matching by attempting to find a new overlapping maker order.
		if !success {
			makerCollatOkay := updateResultToOrderStatus(makerUpdateResult).IsSuccess()
			takerCollatOkay := takerIsLiquidation ||
				updateResultToOrderStatus(takerUpdateResult).IsSuccess()

			// If the maker order failed collateralization checks, add the maker order ID to the
			// list of order IDs to be removed after matching has ended.
			if !makerCollatOkay {
				makerOrdersToRemove = append(makerOrdersToRemove, makerOrder.Order)
			}

			// If this is not a liquidation order and the taker order failed collateralization checks,
			// stop matching.
			if !takerCollatOkay {
				takerOrderStatus.OrderStatus = updateResultToOrderStatus(
					takerUpdateResult,
				)
				break
			}

			// The taker order is a liquidation or it passed collateralization checks, therefore we
			// can continue matching by attempting to find a new overlapping maker order.
			continue
		}

		// If the match was successful, and we're in `DeliverTx` mode, add the events for indexer which indicates
		// a match.
		if lib.IsDeliverTxMode(ctx) {
			// Send on-chain events
			if matchWithOrders.TakerOrder.IsLiquidation() {
				m.clobKeeper.GetIndexerEventManager().AddTxnEvent(
					ctx,
					indexerevents.SubtypeOrderFill,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewLiquidationOrderFillEvent(
							matchWithOrders.MakerOrder.MustGetOrder(),
							matchWithOrders.TakerOrder,
							matchWithOrders.FillAmount,
						),
					),
				)
			} else {
				m.clobKeeper.GetIndexerEventManager().AddTxnEvent(
					ctx,
					indexerevents.SubtypeOrderFill,
					indexer_manager.GetB64EncodedEventMessage(
						indexerevents.NewOrderFillEvent(
							matchWithOrders.MakerOrder.MustGetOrder(),
							matchWithOrders.TakerOrder.MustGetOrder(),
							matchWithOrders.FillAmount,
						),
					),
				)
			}
		}

		// The orders have matched successfully, and the state has been updated.
		// To mark the orders as matched, perform the following actions:
		// 1. Deduct `matchedAmount` from the taker order's remaining quantums, and add the matched
		//    amount to the total matched amount for this matching loop.
		// 2. Add the maker and taker order hash to the map of order hashes.
		// 3. Add the pending fill to the ordered slice of new pending maker fills.
		// 4. If the taker order is a reduce-only order and the user's position size is now zero, cancel any remaining
		//    size of the reduce-only order, and stop matching.

		// 1.
		takerRemainingSize -= matchedAmount

		if newTakerOrder.IsBuy() {
			bigTotalMatchedAmount.Add(bigTotalMatchedAmount, matchedAmount.ToBigInt())
		} else {
			bigTotalMatchedAmount.Sub(bigTotalMatchedAmount, matchedAmount.ToBigInt())
		}

		// 2.
		makerOrderHash := makerOrder.Order.GetOrderHash()
		matchedOrderHashToOrder[makerOrderHash] = &makerOrder.Order
		matchedMakerOrderIdToOrder[makerOrderId] = makerOrder.Order

		// Note that this if statement will only be entered once per every matching cycle. The taker order is hashed
		// inside the loop since we expect the ratio of placed to matched orders to be roughly 100:1, and therefore
		// we want to avoid hashing the taker order if it is not matched.
		if !takerOrderHashWasSet {
			takerOrderHash = newTakerOrder.GetOrderHash()
			matchedOrderHashToOrder[takerOrderHash] = newTakerOrder
			takerOrderHashWasSet = true
		}

		// 3.
		takerSide := types.Order_SIDE_BUY
		if !takerIsBuy {
			takerSide = types.Order_SIDE_SELL
		}
		fillType := types.Trade
		if takerIsLiquidation {
			fillType = types.PerpetualLiquidate
		}
		pendingFill := types.PendingFill{
			MakerOrderHash: makerOrderHash,
			TakerOrderHash: takerOrderHash,
			TakerSide:      takerSide,
			Subticks:       subticks,
			Quantums:       matchedAmount,
			Type:           fillType,
		}
		newPendingFills = append(newPendingFills, pendingFill)
		newMakerFills = append(newMakerFills, types.MakerFill{
			MakerOrderId: makerOrderId,
			FillAmount:   matchedAmount.ToUint64(),
		})

		// 4.
		if newTakerOrder.IsReduceOnly() && takerRemainingSize > 0 {
			takerStatePositionSize := m.clobKeeper.GetStatePosition(ctx, takerSubaccountId, clobPairId)
			if takerStatePositionSize.Sign() == 0 {
				// TODO(DEC-847): Update logic to properly remove long-term taker reduce-only orders.
				takerOrderStatus.OrderStatus = types.ReduceOnlyResized
				break
			}
		}

		// If the taker order was fully matched, stop matching.
		if takerRemainingSize == 0 {
			break
		}
	}

	// Update the remaining size of the taker order now that matching has ended.
	takerOrderStatus.RemainingQuantums = takerRemainingSize
	takerOrderStatus.OrderOptimisticallyFilledQuantums = takerRemainingSizeBeforeMatching - takerRemainingSize

	return newPendingFills,
		newMakerFills,
		matchedOrderHashToOrder,
		matchedMakerOrderIdToOrder,
		makerOrdersToRemove,
		takerOrderStatus
}

// mustRemoveOrder completely removes an order from all data structures for tracking
// open orders in the memclob. If the order does not exist, this method will panic.
// NOTE: `mustRemoveOrder` does _not_ remove cancels.
func (m *MemClobPriceTimePriority) mustRemoveOrder(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.Memclob,
		metrics.RemovedFromOrderBook,
		metrics.Latency,
	)

	// Verify that the order exists.
	levelOrder, exists := m.openOrders.orderIdToLevelOrder[orderId]
	if !exists {
		panic(fmt.Sprintf("mustRemoveOrder: order does not exist %v", orderId))
	}

	m.openOrders.mustRemoveOrder(ctx, levelOrder)
}

// mustUpdateOrderbookStateWithMatchedMakerOrder updates the orderbook with a matched maker order.
// If the maker order is fully filled, it removes it from the orderbook. Else, it updates the total number of quantums
// in the price level containing the maker order.
func (m *MemClobPriceTimePriority) mustUpdateOrderbookStateWithMatchedMakerOrder(
	ctx sdk.Context,
	makerOrder types.Order,
) *types.OffchainUpdates {
	offchainUpdates := types.NewOffchainUpdates()
	makerOrderBaseQuantums := makerOrder.GetBaseQuantums()
	newTotalFilledAmount := m.GetOrderFilledAmount(ctx, makerOrder.OrderId)

	// If the filled amount of the maker order is greater than the order size, panic to avoid silent failure.
	if newTotalFilledAmount > makerOrderBaseQuantums {
		panic("Total filled size of maker order greater than the order size")
	}

	// If the order is fully filled, remove it from the orderbook.
	// Else, update the total level quantums of the level the order is in.
	if newTotalFilledAmount == makerOrderBaseQuantums {
		makerOrderId := makerOrder.OrderId
		m.mustRemoveOrder(ctx, makerOrderId)
	}

	if m.generateOffchainUpdates {
		// Send an off-chain update message to the indexer to update the total filled size of the maker
		// order.
		if message, success := off_chain_updates.CreateOrderUpdateMessage(
			ctx.Logger(),
			makerOrder.OrderId,
			newTotalFilledAmount,
		); success {
			offchainUpdates.AddUpdateMessage(makerOrder.OrderId, message)
		}
	}

	return offchainUpdates
}

// updateResultToOrderStatus translates the result of a collateralization check into a resulting order status.
func updateResultToOrderStatus(updateResult satypes.UpdateResult) types.OrderStatus {
	if updateResult.IsSuccess() {
		return types.Success
	}

	if updateResult == satypes.UpdateCausedError {
		return types.InternalError
	}

	return types.Undercollateralized
}

// getOrderRemainingAmount returns the remaining amount of an order (its size minus its filled amount).
// It also returns a boolean indicating whether the remaining amount is positive (true) or not (false).
func (m *MemClobPriceTimePriority) getOrderRemainingAmount(
	ctx sdk.Context,
	order types.Order,
) (
	remainingAmount satypes.BaseQuantums,
	hasRemainingAmount bool,
) {
	totalFillAmount := m.GetOrderFilledAmount(ctx, order.OrderId)

	// Case: order is completely filled.
	if totalFillAmount >= order.GetBaseQuantums() {
		return 0, false
	}

	return order.GetBaseQuantums() - totalFillAmount, true
}

// RemoveOrderIfFilled removes an order from the orderbook if it currently fully filled in state.
func (m *MemClobPriceTimePriority) RemoveOrderIfFilled(
	ctx sdk.Context,
	orderId types.OrderId,
) {
	// Get LevelOrder.
	levelOrder, levelExists := m.openOrders.orderIdToLevelOrder[orderId]

	// If order is not on the book, return early.
	if !levelExists {
		return
	}

	// Get current fill amount for this order.
	exists, orderStateFillAmount, _ := m.clobKeeper.GetOrderFillAmount(ctx, orderId)

	// If there is no fill amount for this order, return early.
	if !exists {
		return
	}

	// Case: order is now completely filled and can be removed.
	order := levelOrder.Value.Order
	if orderStateFillAmount >= order.GetBaseQuantums() {
		m.openOrders.mustRemoveOrder(ctx, levelOrder)

		// Ensure the nonce is removed from this order.
		isPreexistingStatefulOrder := order.OrderId.IsStatefulOrder() &&
			m.clobKeeper.DoesStatefulOrderExistInState(ctx, order)
		if isPreexistingStatefulOrder {
			m.pendingFills.operationsToPropose.RemovePreexistingStatefulOrderPlacementNonce(order)
		} else {
			m.pendingFills.operationsToPropose.RemoveOrderPlacementNonce(order)
		}
	}
}

// maybeCancelReduceOnlyOrders cancels all open reduce-only orders on the CLOB pair if the new fill would change the
// position side of the subaccount.
func (m *MemClobPriceTimePriority) maybeCancelReduceOnlyOrders(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	clobPairId types.ClobPairId,
	totalBigMatchedQuantums *big.Int,
) (offchainUpdates *types.OffchainUpdates) {
	offchainUpdates = types.NewOffchainUpdates()
	// Get the new position size after matching.
	newPositionSize := m.clobKeeper.GetStatePosition(ctx, subaccountId, clobPairId)

	// Subtract the new match amount from the current position size. This should give us the position size before
	// matching occurred.
	previousPositionSize := new(big.Int).Sub(newPositionSize, totalBigMatchedQuantums)

	// If the subaccount's position size has changed sign, remove all open reduce-only orders.
	if newPositionSize.Sign() != previousPositionSize.Sign() {
		orderbook := m.openOrders.orderbooksMap[clobPairId]

		if openReduceOnlyOrders, exists := orderbook.SubaccountOpenReduceOnlyOrders[subaccountId]; exists {
			// Copy the list of open reduce-only orders.
			openReduceOnlyOrdersCopy := make([]types.OrderId, 0, len(openReduceOnlyOrders))
			for orderId := range openReduceOnlyOrders {
				openReduceOnlyOrdersCopy = append(openReduceOnlyOrdersCopy, orderId)
			}

			// Sort open reduce-only orders by ClientId so that order removal is deterministic. ClientId
			// can be used here since all orders are from the same subaccount, and there should be no
			// duplicates. Determinism is necessary as these removals are part of `DeliverTx`
			// which updates state.
			types.MustSortAndHaveNoDuplicates(openReduceOnlyOrdersCopy)

			// Remove each open reduce-only order from the memclob.
			for _, orderId := range openReduceOnlyOrdersCopy {
				// TODO(DEC-847): Update logic to properly remove long-term orders.
				m.mustRemoveOrder(ctx, orderId)
				if m.generateOffchainUpdates {
					if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
						ctx.Logger(),
						orderId,
						off_chain_updates.OrderRemove_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE,
						off_chain_updates.OrderRemove_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
					); success {
						offchainUpdates.AddRemoveMessage(orderId, message)
					}
				}
			}
		}
	}

	return offchainUpdates
}

// getImpactPriceSubticks returns the impact ask or bid price (in subticks), given the clob pair
// and orderbook. The bid (or ask) impact price is the average price a trader
// would receive if they sold (or bought) from the order book using `impactNotionalAmount`.
// Returns `hasEnoughLiquidity = false` if the orderbook doesn't have
// enough orders on the side to absorb the impact notional amount.
func (m *MemClobPriceTimePriority) getImpactPriceSubticks(
	ctx sdk.Context,
	clobPair types.ClobPair,
	orderbook *types.Orderbook,
	isBid bool,
	impactNotionalQuoteQuantums *big.Int,
) (
	impactPriceSubticks *big.Rat,
	hasEnoughLiquidity bool,
) {
	remainingImpactQuoteQuantums := new(big.Int).Set(impactNotionalQuoteQuantums)
	accumulatedBaseQuantums := new(big.Rat).SetInt64(0)

	makerLevelOrder, foundMakerOrder := m.openOrders.getBestOrderOnSide(orderbook, isBid)
	if impactNotionalQuoteQuantums.Sign() == 0 && foundMakerOrder {
		// Impact notional is zero, returns the price of the best order as impact price.
		return makerLevelOrder.Value.Order.GetOrderSubticks().ToBigRat(), true
	}

	for remainingImpactQuoteQuantums.Sign() > 0 && foundMakerOrder {
		makerOrder := makerLevelOrder.Value.Order
		makerRemainingSize, makerHasRemainingSize := m.getOrderRemainingAmount(ctx, makerOrder)
		if !makerHasRemainingSize {
			panic(fmt.Sprintf("getImpactPriceSubticks: maker order has no remaining amount (%+v)", makerOrder))
		}

		quoteQuantumsIfFullyMatched := types.FillAmountToQuoteQuantums(
			makerOrder.GetOrderSubticks(),
			makerRemainingSize,
			clobPair.QuantumConversionExponent,
		)

		if remainingImpactQuoteQuantums.Cmp(quoteQuantumsIfFullyMatched) > 0 {
			accumulatedBaseQuantums.Add(
				accumulatedBaseQuantums,
				new(big.Rat).SetUint64(makerRemainingSize.ToUint64()),
			)
		} else {
			lastFillFraction := new(big.Rat).SetFrac(
				remainingImpactQuoteQuantums,
				quoteQuantumsIfFullyMatched,
			)

			fractionalBaseQuantums := lastFillFraction.Mul(
				lastFillFraction,
				new(big.Rat).SetInt(makerRemainingSize.ToBigInt()),
			)

			accumulatedBaseQuantums.Add(
				accumulatedBaseQuantums,
				fractionalBaseQuantums,
			)
		}
		remainingImpactQuoteQuantums.Sub(
			remainingImpactQuoteQuantums,
			quoteQuantumsIfFullyMatched,
		)

		// The previous maker order must have been fully matched by the impact order (which has nonzero remaining
		// size), and we need to find the next best maker order.
		makerLevelOrder, foundMakerOrder = m.openOrders.findNextBestLevelOrder(ctx, makerLevelOrder)
	}

	if remainingImpactQuoteQuantums.Sign() > 0 {
		// Impact order was not fully matched, caused by insufficient liquidity.
		return nil, false
	}

	// Impact order was fully matched. Calculate average impact price.
	return types.GetAveragePriceSubticks(
		impactNotionalQuoteQuantums,
		new(big.Int).Div(
			accumulatedBaseQuantums.Num(),
			accumulatedBaseQuantums.Denom(),
		),
		clobPair.QuantumConversionExponent,
	), true
}

// GetPricePremium calculates the premium for a perpetual market, using the equation
// `P = (Max(0, Impact Bid - Index Price) - Max(0, Index Price - Impact Ask)) / Index Price`.
// This is equivalent to the following piece-wise function:
//
//		If Index < Impact Bid:
//	 		P = Impact Bid / Index - 1
//		If Impact Bid ≤ Index ≤ Impact Ask:
//			P = 0
//		If Impact Ask < Index:
//			P = Impact Ask / Index - 1
//
// `Impact Bid/Ask Price` is the average price at which the impact bid/ask order
// (with size of `ImpactNotionalQuoteQuantums`) is filled. If `ImpactNotionalQuoteQuantums`
// is zero, the `Best Bid/Ask Price` is used as `Impact Price`.
// Note that this implies that if there's not enough liquidity for both ask and bid,
// 0 premium is returned since Impact Bid = `0` and Impact Ask = `infinity`.
func (m *MemClobPriceTimePriority) GetPricePremium(
	ctx sdk.Context,
	clobPair types.ClobPair,
	params perptypes.GetPricePremiumParams,
) (
	premiumPpm int32,
	err error,
) {
	// Convert premium vote clamp to int32 (panics if underflows or overflows).
	maxPremiumPpm := lib.MustConvertBigIntToInt32(params.MaxAbsPremiumVotePpm)
	minPremiumPpm := -maxPremiumPpm

	// Check the `ClobPair` is a perpetual.
	if clobPair.GetPerpetualClobMetadata() == nil {
		return 0, sdkerrors.Wrapf(
			types.ErrPremiumWithNonPerpetualClobPair,
			"ClobPair ID: %d",
			clobPair.Id,
		)
	}
	orderbook := m.openOrders.mustGetOrderbook(ctx, clobPair.GetClobPairId())

	// Get index price represented in subticks.
	indexPriceSubticks := types.PriceToSubticks(
		params.Market,
		clobPair,
		params.BaseAtomicResolution,
		params.QuoteAtomicResolution,
	)

	// Check `indexPriceSubticks` is non-zero.
	if indexPriceSubticks.Sign() == 0 {
		return 0, sdkerrors.Wrapf(
			types.ErrZeroIndexPriceForPremiumCalculation,
			"market = %+v, clobPair = %+v, baseAtomicResolution = %d, quoteAtomicResolution = %d",
			params.Market,
			clobPair,
			params.BaseAtomicResolution,
			params.QuoteAtomicResolution,
		)
	}

	bestBid, hasBid := m.openOrders.getBestOrderOnSide(
		orderbook,
		true, // isBuy
	)
	bestAsk, hasAsk := m.openOrders.getBestOrderOnSide(
		orderbook,
		false, // isBuy
	)

	if !hasBid && !hasAsk {
		// Orderbook is empty.
		return 0, nil
	}

	if hasBid && hasAsk && bestBid.Value.Order.Subticks >= bestAsk.Value.Order.Subticks {
		panic(fmt.Sprintf(
			"GetPricePremium: crossing orderbook. ClobPairId = (%+v), bestBid = (%+v), bestAsk = (%+v)",
			clobPair.Id,
			bestBid.Value.Order,
			bestAsk.Value.Order,
		))
	}

	if hasBid && indexPriceSubticks.Cmp(
		new(big.Rat).SetInt(bestBid.Value.Order.GetOrderSubticks().ToBigInt()),
	) < 0 {
		// Index < Best Bid, need to calculate Impact Bid
		return m.getPricePremiumFromSide(
			ctx,
			clobPair,
			orderbook,
			true, // isBid
			params.ImpactNotionalQuoteQuantums,
			indexPriceSubticks,
			minPremiumPpm,
			maxPremiumPpm,
		), nil
	} else if hasAsk && indexPriceSubticks.Cmp(
		new(big.Rat).SetInt(bestAsk.Value.Order.GetOrderSubticks().ToBigInt()),
	) > 0 {
		// Best Ask < Index, need to calculate Impact Ask
		return m.getPricePremiumFromSide(
			ctx,
			clobPair,
			orderbook,
			false, // isBid
			params.ImpactNotionalQuoteQuantums,
			indexPriceSubticks,
			minPremiumPpm,
			maxPremiumPpm,
		), nil
	}

	// Impact Bid <= Best Bid <= Index <= Best Ask <= Impact Ask
	return 0, nil
}

// getPricePremiumFromSide returns the price premium given
// which side (bid/ask) the index price is on.
// `isBid == true` means Index < Best Bid; `isBid == false` means
// Index > Best Ask.
//
// The computed premium is non-zero if and only if one of the two
// cases below is true:
//
// Case 1: `isBid == true` and Impact Bid < Impact Ask < Index:
//
//	P = Impact Ask / Index - 1
//
// Case 2: `isBid == false` and Index < Impact Bid < Impact Ask:
//
//	P = Impact Bid / Index - 1
//
// Computed result is rounded towards zero.
func (m *MemClobPriceTimePriority) getPricePremiumFromSide(
	ctx sdk.Context,
	clobPair types.ClobPair,
	orderbook *types.Orderbook,
	isBid bool,
	impactNotionalQuoteQuantums *big.Int,
	indexPriceSubticks *big.Rat,
	minPremiumPpm int32,
	maxPremiumPpm int32,
) (
	premiumPpm int32,
) {
	impactPriceSubticks, hasEnoughLiquidity := m.getImpactPriceSubticks(
		ctx,
		clobPair,
		orderbook,
		isBid,
		impactNotionalQuoteQuantums,
	)

	if !hasEnoughLiquidity {
		// Impact Ask (Bid) is infinity (Zero), return 0 premium by definition.
		return 0
	}

	cmp := indexPriceSubticks.Cmp(impactPriceSubticks)
	if (!isBid && cmp <= 0) || (isBid && cmp >= 0) {
		// Best Ask < Index <= Impact Ask
		// or
		// Impact Bid <= Index < Best Bid
		return 0
	}

	// Calculate either of the following (in parts-per-million):
	//  Impact Ask / Index - 1
	// or
	//  Impact Bid / Index - 1
	result := new(big.Rat)
	result.Set(impactPriceSubticks).Mul(
		result, lib.BigRatOneMillion(),
	).Quo(
		result,
		indexPriceSubticks,
	).Sub(
		result,
		lib.BigRatOneMillion(),
	)

	// Round result towards zero.
	var resultRounded *big.Int
	if result.Sign() > 0 {
		resultRounded = lib.BigRatRound(result, false)
	} else {
		resultRounded = lib.BigRatRound(result, true)
	}

	return lib.BigInt32Clamp(
		resultRounded,
		minPremiumPpm,
		maxPremiumPpm,
	)
}
