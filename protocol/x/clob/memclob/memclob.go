package memclob

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	"fmt"
	"math/big"
	"runtime/debug"
	"sort"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	indexershared "github.com/dydxprotocol/v4-chain/protocol/indexer/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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

	// OperationsToPropose struct for proposing operations in the next block.
	operationsToPropose types.OperationsToPropose

	// A reference to an expected clob keeper.
	clobKeeper types.MemClobKeeper

	// ---- Fields for determining if off-chain update messages should be generated ----
	generateOffchainUpdates bool
}

type OrderWithRemovalReason struct {
	Order         types.Order
	RemovalReason types.OrderRemoval_RemovalReason
}

func NewMemClobPriceTimePriority(
	generateOffchainUpdates bool,
) *MemClobPriceTimePriority {
	return &MemClobPriceTimePriority{
		openOrders:              newMemclobOpenOrders(),
		cancels:                 newMemclobCancels(),
		operationsToPropose:     *types.NewOperationsToPropose(),
		generateOffchainUpdates: generateOffchainUpdates,
	}
}

// SetClobKeeper sets the MemClobKeeper reference for this MemClob.
// This method is called after the MemClob struct is initialized.
// This reference is set with an explicit method call rather than during `NewMemClobPriceTimePriority`
// due to the bidirectional dependency between the Keeper and the MemClob.
func (m *MemClobPriceTimePriority) SetClobKeeper(clobKeeper types.MemClobKeeper) {
	m.clobKeeper = clobKeeper
}

// CancelOrder removes a Short-Term order by `OrderId` (if it exists) from all order-related data structures
// in the memclob. This method manages only Short-Term cancellations. For stateful cancellations, see
// `msg_server_cancel_orders.go`.

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

	// Stateful orders are not expected here.
	orderIdToCancel.MustBeShortTermOrder()

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
		m.mustRemoveOrder(ctx, orderIdToCancel)

		telemetry.IncrCounter(1, types.ModuleName, metrics.CancelShortTermOrder, metrics.RemovedFromOrderBook)
	}

	// Remove the existing cancel, if any.
	if cancelAlreadyExists {
		m.cancels.remove(orderIdToCancel)
	}

	// Add the new order cancelation.
	m.cancels.addShortTermCancel(orderIdToCancel, goodTilBlock)

	offchainUpdates = types.NewOffchainUpdates()
	if m.generateOffchainUpdates {
		if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
			m.clobKeeper.Logger(ctx),
			orderIdToCancel,
			indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_USER_CANCELED,
			off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
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
}

// CountSubaccountOrders will count all open orders for a given subaccount that match the provided filter.
func (m *MemClobPriceTimePriority) CountSubaccountOrders(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	filter func(types.OrderId) bool,
) (count uint32) {
	for _, openOrdersPerClob := range m.openOrders.orderbooksMap {
		for _, openOrdersPerClobAndSide := range openOrdersPerClob.SubaccountOpenClobOrders[subaccountId] {
			for orderId := range openOrdersPerClobAndSide {
				if filter(orderId) {
					count++
				}
			}
		}
	}
	return count
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
//   - Append all newly-matched orders to the operations queue, along with all new matches.
//   - Update orderbook state by updating the filled amount of all matched maker orders, and removing them if fully
//     filled.
func (m *MemClobPriceTimePriority) mustUpdateMemclobStateWithMatches(
	ctx sdk.Context,
	takerOrder types.MatchableOrder,
	newMakerFills []types.MakerFill,
	matchedOrderHashToOrder map[types.OrderHash]types.MatchableOrder,
	matchedMakerOrderIdToOrder map[types.OrderId]types.Order,
) (offchainUpdates *types.OffchainUpdates) {
	offchainUpdates = types.NewOffchainUpdates()

	// For each order, update `orderHashToMatchableOrder` and `orderIdToOrder`.
	for _, matchedOrder := range matchedOrderHashToOrder {
		// If this is not a liquidation, update `orderIdToOrder`.
		if !matchedOrder.IsLiquidation() {
			order := matchedOrder.MustGetOrder()
			orderId := order.OrderId
			other, exists := m.operationsToPropose.MatchedOrderIdToOrder[orderId]
			if exists && order.MustCmpReplacementOrder(&other) < 0 {
				// Shouldn't happen as the order should have already been replaced or rejected.
				panic(
					"mustUpdateMemclobStateWithMatches: newly matched order is lesser than existing order " +
						"Newly matched order %v, existing order %v",
				)
			}
			m.operationsToPropose.MatchedOrderIdToOrder[orderId] = order
		}
	}

	// Define a data-structure for storing the total number of matched quantums for each subaccount
	// in the matching loop. This is used for reduce-only logic.
	subaccountTotalMatchedQuantums := make(map[satypes.SubaccountId]*big.Int)
	// Ensure each filled maker order has an order placement in `OperationsToPropose`.
	makerFillWithOrders := make([]types.MakerFillWithOrder, 0, len(newMakerFills))
	for _, newFill := range newMakerFills {
		matchedMakerOrder, exists := matchedMakerOrderIdToOrder[newFill.MakerOrderId]
		if !exists {
			panic(
				fmt.Sprintf(
					"mustUpdateMemclobStateWithMatches: matched maker order %s does not exist in `matchedMakerOrderIdToOrder`",
					matchedMakerOrder.GetOrderTextString(),
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

		// Skip adding order placement in the operations queue if it already exists.
		if !m.operationsToPropose.IsOrderPlacementInOperationsQueue(
			matchedMakerOrder,
		) {
			// Add the maker order placement to the operations queue.
			if matchedMakerOrder.IsStatefulOrder() {
				m.operationsToPropose.MustAddStatefulOrderPlacementToOperationsQueue(
					matchedMakerOrder,
				)
			} else {
				m.operationsToPropose.MustAddShortTermOrderPlacementToOperationsQueue(
					matchedMakerOrder,
				)
			}
		}

		// Update the memclob fields for match bookkeeping with the new matches.
		matchedQuantums := satypes.BaseQuantums(newFill.GetFillAmount())

		// Sanity checks.
		if matchedQuantums == 0 {
			panic(fmt.Sprintf(
				"mustUpdateMemclobStateWithMatches: Fill has 0 quantums. Fill %v and maker order %v",
				newFill,
				matchedMakerOrder,
			))
		}

		// Update the orderbook state to reflect the maker order was matched.
		matchOffchainUpdates := m.mustUpdateOrderbookStateWithMatchedMakerOrder(
			ctx,
			matchedMakerOrder,
		)
		offchainUpdates.Append(matchOffchainUpdates)

		// Update the total matched quantums for this matching loop stored in `subaccountTotalMatchedQuantums`.
		for _, order := range []types.MatchableOrder{
			matchedOrderHashToOrder[matchedMakerOrder.GetOrderHash()],
			takerOrder,
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

	// If the taker order is not a liquidation, add the taker order placement to the operations queue.
	if !takerOrder.IsLiquidation() {
		taker := takerOrder.MustGetOrder()

		// Add the taker order placement to the operations queue.
		if taker.IsStatefulOrder() {
			m.operationsToPropose.MustAddStatefulOrderPlacementToOperationsQueue(
				taker,
			)
		} else {
			m.operationsToPropose.MustAddShortTermOrderTxBytes(
				taker,
				ctx.TxBytes(),
			)
			m.operationsToPropose.MustAddShortTermOrderPlacementToOperationsQueue(
				taker,
			)
		}
	}

	// Add the new matches to the operations queue.
	m.operationsToPropose.MustAddMatchToOperationsQueue(takerOrder, makerFillWithOrders)

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

		offchainUpdates.Append(cancelledOffchainUpdates)
	}

	return offchainUpdates
}

// GetOperationsRaw fetches the operations to propose in the next block in raw format
// for placement into MsgProposedOperations.
func (m *MemClobPriceTimePriority) GetOperationsRaw(ctx sdk.Context) (
	operationsQueue []types.OperationRaw,
) {
	return m.operationsToPropose.GetOperationsToPropose()
}

// GetOperationsToReplay fetches the operations to replay in `PrepareCheckState`.
func (m *MemClobPriceTimePriority) GetOperationsToReplay(ctx sdk.Context) (
	[]types.InternalOperation,
	map[types.OrderHash][]byte,
) {
	return m.operationsToPropose.GetOperationsToReplay()
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
		// If this is a replacement order, then ensure we send the appropriate removal message.
		if !order.IsLiquidation() {
			orderId := order.OrderId
			if _, found := m.openOrders.getOrder(ctx, orderId); found {
				if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
					m.clobKeeper.Logger(ctx),
					orderId,
					indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_REPLACED,
					off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
				); success {
					offchainUpdates.AddRemoveMessage(orderId, message)
				}
			}
		}
		if message, success := off_chain_updates.CreateOrderPlaceMessage(
			m.clobKeeper.Logger(ctx),
			order,
		); success {
			offchainUpdates.AddPlaceMessage(order.OrderId, message)
		}
	}

	// Attempt to match the order against the orderbook.
	takerOrderStatus, takerOffchainUpdates, _, err := m.matchOrder(ctx, &order)
	offchainUpdates.Append(takerOffchainUpdates)

	if err != nil {
		if order.IsStatefulOrder() {
			var removalReason types.OrderRemoval_RemovalReason

			if errors.Is(err, types.ErrFokOrderCouldNotBeFullyFilled) {
				if !order.IsConditionalOrder() {
					panic(
						fmt.Sprintf(
							"PlaceOrder: stateful FOK order must be conditional. Order %+v",
							order,
						),
					)
				}
				removalReason = types.OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED
			} else if errors.Is(err, types.ErrPostOnlyWouldCrossMakerOrder) {
				removalReason = types.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER
			}

			if !m.operationsToPropose.IsOrderRemovalInOperationsQueue(order.OrderId) {
				m.operationsToPropose.MustAddOrderRemovalToOperationsQueue(
					order.OrderId,
					removalReason,
				)
			}
		}

		if m.generateOffchainUpdates {
			// Send an off-chain update message indicating the order should be removed from the orderbook
			// on the Indexer.
			if message, success := off_chain_updates.CreateOrderRemoveMessage(
				m.clobKeeper.Logger(ctx),
				order.OrderId,
				takerOrderStatus.OrderStatus,
				err,
				off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
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
				m.clobKeeper.Logger(ctx),
				order.OrderId,
				takerOrderStatus.OrderStatus,
				nil,
				off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(order.OrderId, message)
			}
		}
		// If stateful taker order fails collateralization checks while matching, add Order Removal
		// to operations queue to forcefully remove the order from state.
		if takerOrderStatus.OrderStatus == types.Undercollateralized && order.IsStatefulOrder() {
			if !m.operationsToPropose.IsOrderRemovalInOperationsQueue(order.OrderId) {
				m.operationsToPropose.MustAddOrderRemovalToOperationsQueue(
					order.OrderId,
					types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
				)
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
				m.clobKeeper.Logger(ctx),
				order.OrderId,
				order.GetBaseQuantums(),
			); success {
				offchainUpdates.AddUpdateMessage(order.OrderId, message)
			}
		}
		return orderSizeOptimisticallyFilledFromMatchingQuantums, takerOrderStatus.OrderStatus, offchainUpdates, nil
	}

	// If this is an IOC order, cancel the remaining size since IOC orders cannot be maker orders.
	if order.GetTimeInForce() == types.Order_TIME_IN_FORCE_IOC {
		orderStatus := types.ImmediateOrCancelWouldRestOnBook
		if m.generateOffchainUpdates {
			// Send an off-chain update message indicating the order should be removed from the orderbook
			// on the Indexer.
			if message, success := off_chain_updates.CreateOrderRemoveMessage(
				m.clobKeeper.Logger(ctx),
				order.OrderId,
				orderStatus,
				nil,
				off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(order.OrderId, message)
			}
		}

		// long-term orders cannot use IOC, so we know this stateful order
		// is conditional. Remove the conditional order.
		if order.IsStatefulOrder() && !m.operationsToPropose.IsOrderRemovalInOperationsQueue(order.OrderId) {
			m.operationsToPropose.MustAddOrderRemovalToOperationsQueue(
				order.OrderId,
				types.OrderRemoval_REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK,
			)
		}
		return orderSizeOptimisticallyFilledFromMatchingQuantums, orderStatus, offchainUpdates, nil
	}

	// The taker order has unfilled size which will be added to the orderbook as a maker order.
	// Verify the maker order can be added to the orderbook by performing the add-to-orderbook
	// collateralization check.
	addOrderOrderStatus := m.addOrderToOrderbookCollateralizationCheck(
		ctx,
		order,
	)

	// If the add order to orderbook collateralization check failed, we cannot add the order to the orderbook.
	if !addOrderOrderStatus.IsSuccess() {
		if m.generateOffchainUpdates {
			// Send an off-chain update message indicating the order should be removed from the orderbook
			// on the Indexer.
			if message, success := off_chain_updates.CreateOrderRemoveMessage(
				m.clobKeeper.Logger(ctx),
				order.OrderId,
				addOrderOrderStatus,
				nil,
				off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(order.OrderId, message)
			}
		}

		// remove stateful orders which fail collateralization check while being added to orderbook
		if order.IsStatefulOrder() && !m.operationsToPropose.IsOrderRemovalInOperationsQueue(order.OrderId) {
			m.operationsToPropose.MustAddOrderRemovalToOperationsQueue(
				order.OrderId,
				types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
			)
		}
		return orderSizeOptimisticallyFilledFromMatchingQuantums, addOrderOrderStatus, offchainUpdates, nil
	}

	// If this is a Short-Term order and it's not in the operations queue, add the TX bytes to the
	// operations to propose.
	if order.IsShortTermOrder() &&
		!m.operationsToPropose.IsOrderPlacementInOperationsQueue(order) {
		m.operationsToPropose.MustAddShortTermOrderTxBytes(
			order,
			ctx.TxBytes(),
		)
	}

	// Add the order to the orderbook and all other bookkeeping data structures.
	m.mustAddOrderToOrderbook(ctx, order, false)

	// If the taker order is added to the orderbook successfully, send an off-chain message with
	// the total filled size of the order (size of order - remaining size).
	if m.generateOffchainUpdates {
		if message, success := off_chain_updates.CreateOrderUpdateMessage(
			m.clobKeeper.Logger(ctx),
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
//
// TODO(CLOB-852): Separate out deleveraging flow from liquidations flow.
func (m *MemClobPriceTimePriority) PlacePerpetualLiquidation(
	ctx sdk.Context,
	liquidationOrder types.LiquidationOrder,
) (
	orderSizeOptimisticallyFilledFromMatchingQuantums satypes.BaseQuantums,
	orderStatus types.OrderStatus,
	offchainUpdates *types.OffchainUpdates,
	err error,
) {
	// Attempt to match the liquidation order against the orderbook.
	// TODO(DEC-1157): Update liquidations flow to send off-chain indexer messages.
	liquidationOrderStatus, offchainUpdates, _, err := m.matchOrder(ctx, &liquidationOrder)
	if err != nil {
		return 0, 0, nil, err
	}

	// Skip checking if the account needs to be deleveraged if the liquidation order was partially
	// or fully-filled.
	if liquidationOrderStatus.OrderOptimisticallyFilledQuantums > 0 {
		return liquidationOrderStatus.OrderOptimisticallyFilledQuantums,
			liquidationOrderStatus.OrderStatus,
			offchainUpdates,
			err
	}

	canPerformDeleveraging, deleverageErr := m.clobKeeper.CanDeleverageSubaccount(ctx, liquidationOrder.GetSubaccountId())
	if deleverageErr != nil {
		return 0, 0, offchainUpdates, deleverageErr
	}

	// Early return to skip deleveraging if the subaccount can't be deleveraged.
	if !canPerformDeleveraging {
		return liquidationOrderStatus.OrderOptimisticallyFilledQuantums,
			liquidationOrderStatus.OrderStatus,
			offchainUpdates,
			err
	}

	// Deleverage the full liquidation order size from the subaccount's position size.
	deltaQuantums := liquidationOrder.GetBaseQuantums().ToBigInt()
	if !liquidationOrder.IsBuy() {
		deltaQuantums = deltaQuantums.Neg(deltaQuantums)
	}

	fills, _ := m.clobKeeper.OffsetSubaccountPerpetualPosition(
		ctx,
		liquidationOrder.GetSubaccountId(),
		liquidationOrder.MustGetLiquidatedPerpetualId(),
		deltaQuantums,
	)

	if len(fills) > 0 {
		m.operationsToPropose.MustAddDeleveragingToOperationsQueue(
			liquidationOrder.GetSubaccountId(),
			liquidationOrder.MustGetLiquidatedPerpetualId(),
			fills,
		)
	}

	return liquidationOrderStatus.OrderOptimisticallyFilledQuantums,
		liquidationOrderStatus.OrderStatus,
		offchainUpdates,
		err
}

// matchOrder will match the provided `MatchableOrder` as a taker order against the respective orderbook.
// This function will return the status of the matched order, along with the new taker pending matches.
// If order matching results in any error, all state updates wil be discarded.
func (m *MemClobPriceTimePriority) matchOrder(
	ctx sdk.Context,
	order types.MatchableOrder,
) (
	orderStatus types.TakerOrderStatus,
	offchainUpdates *types.OffchainUpdates,
	makerOrdersToRemove []OrderWithRemovalReason,
	err error,
) {
	offchainUpdates = types.NewOffchainUpdates()

	// // Branch the state. State will be wrote to only if matching does not return an error.
	branchedContext, writeCache := ctx.CacheContext()

	// Attempt to match the order against the orderbook.
	newMakerFills,
		matchedOrderHashToOrder,
		matchedMakerOrderIdToOrder,
		makerOrdersToRemove,
		takerOrderStatus := m.mustPerformTakerOrderMatching(
		branchedContext,
		order,
	)

	// If this is a replacement order, then ensure we remove the existing order from the orderbook.
	if !order.IsLiquidation() {
		orderId := order.MustGetOrder().OrderId
		if orderToBeReplaced, found := m.openOrders.getOrder(branchedContext, orderId); found {
			makerOrdersToRemove = append(makerOrdersToRemove, OrderWithRemovalReason{Order: orderToBeReplaced})
		}
	}

	// For each maker order that should be removed, remove it from the orderbook and emit off-chain
	// updates for the indexer.
	for _, makerOrderWithRemovalReason := range makerOrdersToRemove {
		// TODO(DEC-847): Update logic to properly remove long-term orders.
		makerOrderId := makerOrderWithRemovalReason.Order.OrderId
		// TODO(CLOB-669): Move logic outside of `memclob.go` by returning a slice of removed orders.
		// If the order is a replacement order, a message was already added above the place message.
		if m.generateOffchainUpdates && (order.IsLiquidation() || makerOrderId != order.MustGetOrder().OrderId) {
			// If the taker order and the removed maker order are from the same subaccount, set
			// the reason to SELF_TRADE error, otherwise set the reason to be UNDERCOLLATERALIZED.
			// TODO(DEC-1409): Update this to support order replacements on indexer.
			reason := indexershared.ConvertOrderRemovalReasonToIndexerOrderRemovalReason(
				makerOrderWithRemovalReason.RemovalReason,
			)
			if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
				branchedContext.Logger(),
				makerOrderId,
				reason,
				off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
			); success {
				offchainUpdates.AddRemoveMessage(makerOrderId, message)
			}
		}

		m.mustRemoveOrder(branchedContext, makerOrderId)
		if makerOrderId.IsStatefulOrder() && !m.operationsToPropose.IsOrderRemovalInOperationsQueue(makerOrderId) {
			m.operationsToPropose.MustAddOrderRemovalToOperationsQueue(
				makerOrderId,
				makerOrderWithRemovalReason.RemovalReason,
			)
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
	// TODO(DEC-998): Determine if allowing post-only orders to match in rewind step is valid.
	if len(newMakerFills) > 0 &&
		!order.IsLiquidation() &&
		order.MustGetOrder().TimeInForce == types.Order_TIME_IN_FORCE_POST_ONLY {
		matchingErr = types.ErrPostOnlyWouldCrossMakerOrder
	}

	// If the match is valid and placing the taker order generated valid matches, update memclob state.
	takerGeneratedValidMatches := len(newMakerFills) > 0 && matchingErr == nil
	if takerGeneratedValidMatches {
		matchOffchainUpdates := m.mustUpdateMemclobStateWithMatches(
			branchedContext,
			order,
			newMakerFills,
			matchedOrderHashToOrder,
			matchedMakerOrderIdToOrder,
		)
		offchainUpdates.Append(matchOffchainUpdates)
		writeCache()
	}

	return takerOrderStatus, offchainUpdates, makerOrdersToRemove, matchingErr
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
	localOperations []types.InternalOperation,
	shortTermOrderTxBytes map[types.OrderHash][]byte,
	existingOffchainUpdates *types.OffchainUpdates,
) *types.OffchainUpdates {
	// Recover from any panics that occur during replay operations.
	// This could happen in cases where i.e. A subaccount balance overflowed
	// during a match. We don't want to halt the entire chain in this case.
	// TODO(CLOB-275): Do not gracefully handle panics in `PrepareCheckState`.
	defer func() {
		if r := recover(); r != nil {
			stackTrace := string(debug.Stack())
			m.clobKeeper.Logger(ctx).Error("panic in replay operations", "panic", r, "stackTrace", stackTrace)
		}
	}()

	// Iterate over all provided operations.
	for _, operation := range localOperations {
		switch operation.Operation.(type) {
		// Replay all short-term and stateful order placements.
		case *types.InternalOperation_ShortTermOrderPlacement:
			order := operation.GetShortTermOrderPlacement().Order

			// Set underlying tx bytes so OperationsToPropose may access it and
			// store the tx bytes on OperationHashToTxBytes data structure
			shortTermOrderTxBytes, exists := shortTermOrderTxBytes[order.GetOrderHash()]
			if !exists || len(shortTermOrderTxBytes) == 0 {
				panic(
					fmt.Sprintf(
						"ReplayOperations: Short-Term order TX bytes not found for order %s",
						order.GetOrderTextString(),
					),
				)
			}
			ctx = ctx.WithTxBytes(shortTermOrderTxBytes)

			// Note we use `clobKeeper.PlaceOrder` here to ensure the proper stateful validation is performed and
			// newly-placed stateful orders are written to state. In the future this will be important for sequence number
			// verification as well.
			// TODO(DEC-1755): Account for sequence number verification.
			// TODO(DEC-998): Research whether it's fine for two post-only orders to be matched. Currently they are dropped.
			msg := types.NewMsgPlaceOrder(order)
			orderSizeOptimisticallyFilledFromMatchingQuantums,
				orderStatus, placeOrderOffchainUpdates, err := m.clobKeeper.ReplayPlaceOrder(
				ctx,
				msg,
			)

			m.clobKeeper.Logger(ctx).Debug(
				"Received new order",
				"orderHash",
				log.NewLazySprintf("%X", order.GetOrderHash()),
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

			existingOffchainUpdates = m.GenerateOffchainUpdatesForReplayPlaceOrder(
				ctx,
				err,
				operation,
				order,
				orderStatus,
				placeOrderOffchainUpdates,
				existingOffchainUpdates,
			)

		// Replay all pre-existing stateful order placements.
		case *types.InternalOperation_PreexistingStatefulOrder:
			orderId := operation.GetPreexistingStatefulOrder()

			// TODO(DEC-1974): The `PreexistingStatefulOrder` operation
			// does not contain the order hash, so we cannot check if the
			// order is the same as the one in the book (rather than a replacement).
			// For consistency we should fix this, but currently it is not an issue as
			// replacements are not currently supported, and trying to place an older version
			// of a stateful order should fail due to the `GoodTilBlockTime`.
			statefulOrderPlacement, found := m.clobKeeper.GetLongTermOrderPlacement(ctx, *orderId)
			if !found {
				// It's possible that this order was canceled or expired in the last committed block.
				continue
			}

			// Note that we use `memclob.PlaceOrder` here, this will skip writing the stateful order placement to state.
			// TODO(DEC-998): Research whether it's fine for two post-only orders to be matched. Currently they are dropped.
			_, orderStatus, placeOrderOffchainUpdates, err := m.PlaceOrder(
				ctx,
				statefulOrderPlacement.Order,
			)
			existingOffchainUpdates = m.GenerateOffchainUpdatesForReplayPlaceOrder(
				ctx,
				err,
				operation,
				statefulOrderPlacement.Order,
				orderStatus,
				placeOrderOffchainUpdates,
				existingOffchainUpdates,
			)
		// Matches are a no-op.
		case *types.InternalOperation_Match:
		case *types.InternalOperation_OrderRemoval:
			// Re-place orders which were not removed by the previous block proposer to give them a "second chance".
			// It is possible that this placement or a subsequent match operation will
			// cause the Order Removal to be generated once again.
			orderId := operation.GetOrderRemoval().OrderId
			statefulOrderPlacement, found := m.clobKeeper.GetLongTermOrderPlacement(ctx, orderId)

			// if not in state anymore, this means it was removed in the previous block. No-op.
			if !found {
				continue
			}

			_, orderStatus, placeOrderOffchainUpdates, err := m.PlaceOrder(
				ctx,
				statefulOrderPlacement.Order,
			)
			existingOffchainUpdates = m.GenerateOffchainUpdatesForReplayPlaceOrder(
				ctx,
				err,
				operation,
				statefulOrderPlacement.Order,
				orderStatus,
				placeOrderOffchainUpdates,
				existingOffchainUpdates,
			)
		default:
			panic(fmt.Sprintf("unknown operation type: %T", operation.Operation))
		}
	}

	existingOffchainUpdates.CondenseMessagesForReplay()
	return existingOffchainUpdates
}

// GenerateOffchainUpdatesForReplayPlaceOrder is a helper function intended to be used in ReplayOperations.
// It takes the results of a PlaceOrder function call, emits the according logs, and appends offchain updates for
// the replay operation to the existingOffchainUpdates object.
func (m *MemClobPriceTimePriority) GenerateOffchainUpdatesForReplayPlaceOrder(
	ctx sdk.Context,
	err error,
	operation types.InternalOperation,
	order types.Order,
	orderStatus types.OrderStatus,
	placeOrderOffchainUpdates *types.OffchainUpdates,
	existingOffchainUpdates *types.OffchainUpdates,
) *types.OffchainUpdates {
	orderId := order.OrderId
	if err != nil {
		var loggerString string
		switch operation.Operation.(type) {
		case *types.InternalOperation_ShortTermOrderPlacement:
			loggerString = "ReplayOperations: PlaceOrder() returned an error"
		case *types.InternalOperation_PreexistingStatefulOrder:
			loggerString = "ReplayOperations: PlaceOrder() returned an error for a pre-existing stateful order."
		case *types.InternalOperation_OrderRemoval:
			loggerString = "ReplayOperations: PlaceOrder() returned an error for a removed stateful order which was re-placed."
		}
		m.clobKeeper.Logger(ctx).Debug(
			loggerString,
			"error",
			err,
			"operation",
			operation,
			"order",
			order,
		)

		// If the order is dropped while adding it to the book, return an off-chain order remove
		// message for the order.
		if m.generateOffchainUpdates && off_chain_updates.ShouldSendOrderRemovalOnReplay(err) {
			if message, success := off_chain_updates.CreateOrderRemoveMessageWithDefaultReason(
				m.clobKeeper.Logger(ctx),
				orderId,
				orderStatus,
				err,
				off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
				indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_INTERNAL_ERROR,
			); success {
				existingOffchainUpdates.AddRemoveMessage(orderId, message)
			}
		}
	} else if m.generateOffchainUpdates {
		existingOffchainUpdates.Append(placeOrderOffchainUpdates)
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
	localValidatorOperationsQueue []types.InternalOperation,
) {
	// Clear the OTP. This will also remove nonces for every operation in `operationsQueueCopy`.
	m.operationsToPropose.ClearOperationsQueue()

	// For each order placement operation in the copy, remove the order from the book
	// if it exists.
	for _, operation := range localValidatorOperationsQueue {
		switch operation.Operation.(type) {
		case *types.InternalOperation_ShortTermOrderPlacement:
			otpOrderId := operation.GetShortTermOrderPlacement().Order.OrderId
			otpOrderHash := operation.GetShortTermOrderPlacement().Order.GetOrderHash()

			// If the order exists in the book, remove it.
			// Else if the order is a Short-Term order, since it's no longer on the book or operations
			// queue we should remove the order hash from `ShortTermOrderTxBytes`.
			existingOrder, found := m.openOrders.getOrder(ctx, otpOrderId)
			if found && existingOrder.GetOrderHash() == otpOrderHash {
				m.mustRemoveOrder(ctx, otpOrderId)
			} else if otpOrderId.IsShortTermOrder() {
				order := operation.GetShortTermOrderPlacement().Order
				if _, exists := m.operationsToPropose.
					ShortTermOrderHashToTxBytes[order.GetOrderHash()]; !exists {
					panic(
						fmt.Sprintf(
							"RemoveAndClearOperationsQueue: No TxBytes to remove for Short-Term order %+v",
							order.GetOrderTextString(),
						),
					)
				}
				m.operationsToPropose.RemoveShortTermOrderTxBytes(order)
			}
		case *types.InternalOperation_PreexistingStatefulOrder:
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
}

// PurgeInvalidMemclobState will purge the following invalid state from the memclob:
// - Expired Short-Term order cancellations.
// - Expired Short-Term and stateful orders.
// - Fully-filled orders.
// - Canceled stateful orders.
// - Forcefully removed stateful orders.
func (m *MemClobPriceTimePriority) PurgeInvalidMemclobState(
	ctx sdk.Context,
	fullyFilledOrderIds []types.OrderId,
	expiredStatefulOrderIds []types.OrderId,
	canceledStatefulOrderIds []types.OrderId,
	removedStatefulOrderIds []types.OrderId,
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
			m.mustRemoveOrder(ctx, statefulOrderId)
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
			m.mustRemoveOrder(ctx, statefulOrderId)

			if m.generateOffchainUpdates {
				// Send an off-chain update message indicating the stateful order should be removed from the
				// orderbook on the Indexer. As the order is expired, the status of the order is canceled
				// and not best-effort-canceled.
				if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
					m.clobKeeper.Logger(ctx),
					statefulOrderId,
					indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_EXPIRED,
					off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_CANCELED,
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
					m.clobKeeper.Logger(ctx),
					shortTermOrderId,
					indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_EXPIRED,
					off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_CANCELED,
				); success {
					existingOffchainUpdates.AddRemoveMessage(shortTermOrderId, message)
				}
			}

			m.mustRemoveOrder(ctx, shortTermOrderId)
		}
	}

	// Remove all forcefully removed stateful order IDs from the memclob if they exist.
	// Indexer events are sent during DeliverTx and therefore do not need to be sent here.
	for _, statefulOrderId := range removedStatefulOrderIds {
		statefulOrderId.MustBeStatefulOrder()

		if m.openOrders.hasOrder(ctx, statefulOrderId) {
			m.mustRemoveOrder(ctx, statefulOrderId)
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
//   - This subaccount has strictly less open orders than the equity tier limit the subaccount qualifies for.
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
//   - `Order.Quantums` is greater than orderbook's MinOrderBaseQuantums (equal to `ClobPair.StepBaseQuantums`)
//   - `Order.Quantums` is a multiple of `ClobPair.StepBaseQuantums`.
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
			return errorsmod.Wrapf(
				types.ErrOrderIsCanceled,
				"Order: %+v, Cancellation GoodTilBlock: %d",
				order,
				cancelTilBlock,
			)
		}
	}

	existingRestingOrder, restingOrderExists := m.openOrders.getOrder(ctx, orderId)
	existingMatchedOrder, matchedOrderExists := m.operationsToPropose.MatchedOrderIdToOrder[orderId]

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
	// we need to ensure that adding the new order would not exceed any equity tier limits.
	// Note: The order could already be resting on the book for the same side if this order is a replacement.
	// Note: The order could already be resting on the book for a different side if this order is a replacement.
	doesOrderAlreadyExistForSide := restingOrderExists && existingRestingOrder.Side == order.Side
	if !doesOrderAlreadyExistForSide {
		if err := m.clobKeeper.ValidateSubaccountEquityTierLimitForNewOrder(ctx, order); err != nil {
			return err
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
		return errorsmod.Wrapf(
			types.ErrOrderFullyFilled,
			"Order remaining amount is less than `MinOrderBaseQuantums`. Remaining amount: %d. Order: %+v",
			remainingAmount,
			order.GetOrderTextString(),
		)
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
	// A slice of new maker fills created from matching this taker order.
	newMakerFills []types.MakerFill,
	// A map of matched order hashes to the order.
	matchedOrderHashToOrder map[types.OrderHash]types.MatchableOrder,
	// A map of matched maker order IDs to the order.
	matchedMakerOrderIdToOrder map[types.OrderId]types.Order,
	// A list of maker orders that failed collateralization checks during matching and should be removed from the
	// orderbook.
	makerOrdersToRemove []OrderWithRemovalReason,
	// The status of the taker order, specifically the remaining size, optimistically filled size, and the result of the
	// last collateralization check.
	// This is necessary for determining whether remaining taker order size can be added to the orderbook, and for
	// returning the optimistically filled size to the caller.
	takerOrderStatus types.TakerOrderStatus,
) {
	// Initialize return variables.
	newMakerFills = make([]types.MakerFill, 0)
	matchedOrderHashToOrder = make(map[types.OrderHash]types.MatchableOrder)
	matchedMakerOrderIdToOrder = make(map[types.OrderId]types.Order)
	takerOrderStatus.OrderStatus = types.Success
	makerOrdersToRemove = make([]OrderWithRemovalReason, 0)

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
			makerOrdersToRemove = append(
				makerOrdersToRemove,
				OrderWithRemovalReason{
					Order:         makerOrder.Order,
					RemovalReason: types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
				},
			)
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
			resizedMatchAmount := m.resizeReduceOnlyMatchIfNecessary(
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
				makerOrdersToRemove = append(
					makerOrdersToRemove,
					OrderWithRemovalReason{
						Order:         makerOrder.Order,
						RemovalReason: types.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY,
					},
				)
				continue
			}

			matchedAmount = resizedMatchAmount
		}

		if newTakerOrder.IsReduceOnly() {
			currentPositionSize := m.clobKeeper.GetStatePosition(ctx, takerSubaccountId, clobPairId)
			resizedMatchAmount := m.resizeReduceOnlyMatchIfNecessary(
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

		// Perform collateralization checks to verify the orders can be filled.
		matchWithOrders := types.MatchWithOrders{
			TakerOrder: newTakerOrder,
			MakerOrder: &makerOrder.Order,
			FillAmount: matchedAmount,
		}

		success, takerUpdateResult, makerUpdateResult, _, err := m.clobKeeper.ProcessSingleMatch(ctx, &matchWithOrders)
		if err != nil && !errors.Is(err, satypes.ErrFailedToUpdateSubaccounts) {
			if errors.Is(err, types.ErrLiquidationExceedsSubaccountMaxInsuranceLost) {
				// Subaccount has reached max insurance lost block limit. Stop matching.
				telemetry.IncrCounter(1, types.ModuleName, metrics.SubaccountMaxInsuranceLost, metrics.Count)
				takerOrderStatus.OrderStatus = types.LiquidationExceededSubaccountMaxInsuranceLost
				break
			}
			if errors.Is(err, types.ErrLiquidationExceedsSubaccountMaxNotionalLiquidated) {
				// Subaccount has reached max notional liquidated block limit. Stop matching.
				telemetry.IncrCounter(1, types.ModuleName, metrics.SubaccountMaxNotionalLiquidated, metrics.Count)
				takerOrderStatus.OrderStatus = types.LiquidationExceededSubaccountMaxNotionalLiquidated
				break
			}
			if errors.Is(err, types.ErrInsuranceFundHasInsufficientFunds) {
				// Deleveraging is required. Stop matching.
				telemetry.IncrCounter(1, types.ModuleName, metrics.LiquidationRequiresDeleveraging, metrics.Count)
				takerOrderStatus.OrderStatus = types.LiquidationRequiresDeleveraging
				break
			}

			// Panic since this is an unknown error.
			m.clobKeeper.Logger(ctx).Error(
				"Unexpected error from `ProcessSingleMatch`",
				"error",
				err,
				"matchWithOrders",
				matchWithOrders,
			)
			panic(err)
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
				makerOrdersToRemove = append(
					makerOrdersToRemove,
					OrderWithRemovalReason{
						Order:         makerOrder.Order,
						RemovalReason: types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
					},
				)
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
		newMakerFills = append(newMakerFills, types.MakerFill{
			MakerOrderId: makerOrderId,
			FillAmount:   matchedAmount.ToUint64(),
		})

		// 4.
		if newTakerOrder.IsReduceOnly() && takerRemainingSize > 0 {
			takerStatePositionSize := m.clobKeeper.GetStatePosition(ctx, takerSubaccountId, clobPairId)
			if takerStatePositionSize.Sign() == 0 {
				// TODO(DEC-847): Update logic to properly remove stateful taker reduce-only orders.
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

	return newMakerFills,
		matchedOrderHashToOrder,
		matchedMakerOrderIdToOrder,
		makerOrdersToRemove,
		takerOrderStatus
}

// SetMemclobGauges sets gauges for each orderbook and the operations queue based on current memclob state.
// This is used only for observability purposes.
func (m *MemClobPriceTimePriority) SetMemclobGauges(
	ctx sdk.Context,
) {
	// Set gauges for each orderbook.
	for clobPairId, orderbook := range m.openOrders.orderbooksMap {
		// Set gauge for total open orders on each orderbook.
		telemetry.SetGaugeWithLabels(
			[]string{
				types.ModuleName,
				metrics.TotalOrdersClobPair,
			},
			float32(orderbook.TotalOpenOrders),
			[]gometrics.Label{
				metrics.GetLabelForIntValue(metrics.ClobPairId, int(clobPairId)),
			},
		)

		// Set gauge for best bid on each orderbook.
		telemetry.SetGaugeWithLabels(
			[]string{
				types.ModuleName,
				metrics.BestBidClobPair,
			},
			float32(orderbook.BestBid),
			[]gometrics.Label{
				metrics.GetLabelForIntValue(metrics.ClobPairId, int(clobPairId)),
			},
		)

		// Set gauge for best ask on each orderbook.
		telemetry.SetGaugeWithLabels(
			[]string{
				types.ModuleName,
				metrics.BestAskClobPair,
			},
			float32(orderbook.BestAsk),
			[]gometrics.Label{
				metrics.GetLabelForIntValue(metrics.ClobPairId, int(clobPairId)),
			},
		)
	}

	// Set gauges for the operations queue.

	telemetry.SetGauge(
		float32(len(m.operationsToPropose.OperationsQueue)),
		types.ModuleName,
		metrics.OperationsQueueLength,
	)

	telemetry.SetGauge(
		float32(len(m.operationsToPropose.OrderHashesInOperationsQueue)),
		types.ModuleName,
		metrics.NumMatchedOrdersInOperationsQueue,
	)

	telemetry.SetGauge(
		float32(len(m.operationsToPropose.ShortTermOrderHashToTxBytes)),
		types.ModuleName,
		metrics.NumShortTermOrderTxBytes,
	)
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

	// If this is a Short-Term order and it's not in the operations queue, then remove it from
	// `ShortTermOrderTxBytes`.
	order := levelOrder.Value.Order
	if order.IsShortTermOrder() &&
		!m.operationsToPropose.IsOrderPlacementInOperationsQueue(order) {
		m.operationsToPropose.RemoveShortTermOrderTxBytes(order)
	}
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
	// Note we shouldn't remove Short-Term order hashes from `ShortTermOrderTxBytes` here since
	// the order was matched.
	if newTotalFilledAmount == makerOrderBaseQuantums {
		makerOrderId := makerOrder.OrderId
		m.mustRemoveOrder(ctx, makerOrderId)
	}

	if m.generateOffchainUpdates {
		// Send an off-chain update message to the indexer to update the total filled size of the maker
		// order.
		if message, success := off_chain_updates.CreateOrderUpdateMessage(
			m.clobKeeper.Logger(ctx),
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
		m.mustRemoveOrder(ctx, order.OrderId)
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
				// TODO(DEC-847): Update logic to properly remove stateful orders.
				m.mustRemoveOrder(ctx, orderId)
				if orderId.IsStatefulOrder() && !m.operationsToPropose.IsOrderRemovalInOperationsQueue(orderId) {
					m.operationsToPropose.MustAddOrderRemovalToOperationsQueue(
						orderId,
						types.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY,
					)
				}
				if m.generateOffchainUpdates {
					if message, success := off_chain_updates.CreateOrderRemoveMessageWithReason(
						m.clobKeeper.Logger(ctx),
						orderId,
						indexershared.OrderRemovalReason_ORDER_REMOVAL_REASON_REDUCE_ONLY_RESIZE,
						off_chain_updates.OrderRemoveV1_ORDER_REMOVAL_STATUS_BEST_EFFORT_CANCELED,
					); success {
						offchainUpdates.AddRemoveMessage(orderId, message)
					}
				}
			}
		}
	}

	return offchainUpdates
}

// GetMidPrice returns the mid price of the orderbook for the given clob pair
// and whether or not it exists.
func (m *MemClobPriceTimePriority) GetMidPrice(
	ctx sdk.Context,
	clobPairId types.ClobPairId,
) (
	subticks types.Subticks,
	exists bool,
) {
	subticks, exists = m.openOrders.orderbooksMap[clobPairId].GetMidPrice()
	if !exists {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, metrics.MissingMidPrice, metrics.Count},
			1,
			[]gometrics.Label{
				metrics.GetLabelForIntValue(
					metrics.ClobPairId,
					int(clobPairId.ToUint32()),
				),
			},
		)
	}
	return subticks, exists
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
//		If Impact Bid ≤ Index ≤ Impact Ask:
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
		return 0, errorsmod.Wrapf(
			types.ErrPremiumWithNonPerpetualClobPair,
			"ClobPair ID: %d",
			clobPair.Id,
		)
	}
	orderbook := m.openOrders.mustGetOrderbook(ctx, clobPair.GetClobPairId())

	// Get index price represented in subticks.
	indexPriceSubticks := types.PriceToSubticks(
		params.MarketPrice,
		clobPair,
		params.BaseAtomicResolution,
		params.QuoteAtomicResolution,
	)

	// Check `indexPriceSubticks` is non-zero.
	if indexPriceSubticks.Sign() == 0 {
		return 0, errorsmod.Wrapf(
			types.ErrZeroIndexPriceForPremiumCalculation,
			"market = %+v, clobPair = %+v, baseAtomicResolution = %d, quoteAtomicResolution = %d",
			params.MarketPrice,
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

// resizeReduceOnlyMatchIfNecessary resizes a reduce-only match if it would change or increase
// the position side of the subaccount, and returns the resized match.
func (m *MemClobPriceTimePriority) resizeReduceOnlyMatchIfNecessary(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	clobPairId types.ClobPairId,
	currentPositionSize *big.Int,
	newlyMatchedAmount satypes.BaseQuantums,
	isBuy bool,
) satypes.BaseQuantums {
	// Get the signed size of the new match.
	newMatchSize := newlyMatchedAmount.ToBigInt()
	if !isBuy {
		newMatchSize.Neg(newMatchSize)
	}

	// If the match is not on the opposite side of the position, then the match is invalid.
	// Note that this can occur for reduce-only maker orders if the maker subaccount's position side
	// changes during the matching loop, and this should never happen for taker orders.
	if currentPositionSize.Sign()*newMatchSize.Sign() != -1 {
		return satypes.BaseQuantums(0)
	}

	// The match is on the opposite side of the position. Return the minimum of the match size and
	// position size to ensure the new match does not change the subaccount's position side.
	absPositionSize := new(big.Int).Abs(currentPositionSize)
	absNewMatchSize := new(big.Int).Abs(newMatchSize)
	maxMatchSize := lib.BigMin(absPositionSize, absNewMatchSize)
	return satypes.BaseQuantums(maxMatchSize.Uint64())
}
