package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// UntriggeredConditionalOrders is an in-memory struct stored on the clob Keeper.
// It is intended to efficiently store placed conditional orders and poll out triggered
// conditional orders on oracle price changes for a given ClobPairId.
// All orders contained in this data structure are placed conditional orders with the same
// ClobPairId and are untriggered, unexpired, and uncancelled.
// Note that we are using a Order list for the initial implementation, but for
// optimal runtime a an AVL-tree backed priority queue would work.
// TODO(CLOB-717) Change list to use priority queue.
type UntriggeredConditionalOrders struct {
	// All untriggered take profit buy orders and stop loss sell orders sorted by time priority.
	// These orders will be triggered when the oracle price goes lower than or equal to the trigger price.
	// This array functions like a max heap.
	OrdersToTriggerWhenOraclePriceLTETriggerPrice []types.Order

	// All untriggered take profit sell orders and stop loss buy orders sorted by time priority.
	// These orders will be triggered when the oracle price goes greater than or equal to the trigger price.
	// This array functions like a min heap.
	OrdersToTriggerWhenOraclePriceGTETriggerPrice []types.Order
}

func (k Keeper) NewUntriggeredConditionalOrders() *UntriggeredConditionalOrders {
	return NewUntriggeredConditionalOrders()
}

func NewUntriggeredConditionalOrders() *UntriggeredConditionalOrders {
	return &UntriggeredConditionalOrders{
		OrdersToTriggerWhenOraclePriceLTETriggerPrice: make([]types.Order, 0),
		OrdersToTriggerWhenOraclePriceGTETriggerPrice: make([]types.Order, 0),
	}
}

// IsEmpty returns true if the UntriggeredConditionalOrders' order slices are both empty.
func (untriggeredOrders *UntriggeredConditionalOrders) IsEmpty() bool {
	return len(untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice) == 0 &&
		len(untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice) == 0
}

// AddUntriggeredConditionalOrders takes in a list of newly-placed conditional order ids and adds them
// to the in-memory UntriggeredConditionalOrders struct, filtering out orders that have been cancelled
// or expired in the last block. This function is used in EndBlocker and on application startup.
func (k Keeper) AddUntriggeredConditionalOrders(
	ctx sdk.Context,
	placedConditionalOrderIds []types.OrderId,
	placedStatefulCancellationOrderIds map[types.OrderId]struct{},
	expiredStatefulOrderIdsSet map[types.OrderId]struct{},
) {
	for _, orderId := range placedConditionalOrderIds {
		_, isCancelled := placedStatefulCancellationOrderIds[orderId]
		_, isExpired := expiredStatefulOrderIdsSet[orderId]
		if isCancelled || isExpired {
			continue
		}

		orderPlacement, exists := k.GetUntriggeredConditionalOrderPlacement(ctx, orderId)
		if !exists {
			panic(
				fmt.Sprintf(
					"AddUntriggeredConditionalOrders: order placement does not exist in state for untriggered "+
						"conditional order id, OrderId %+v.",
					orderId,
				),
			)
		}

		clobPairId := types.ClobPairId(orderId.GetClobPairId())
		untriggeredConditionalOrders, exists := k.UntriggeredConditionalOrders[clobPairId]
		if !exists {
			untriggeredConditionalOrders = k.NewUntriggeredConditionalOrders()
			k.UntriggeredConditionalOrders[clobPairId] = untriggeredConditionalOrders
		}
		untriggeredConditionalOrders.AddUntriggeredConditionalOrder(orderPlacement.GetOrder())
	}
}

// AddUntriggeredConditionalOrder adds an untriggered conditional order to the UntriggeredConditionalOrders
// data structure. It will panic if the order is not a conditional order.
func (untriggeredOrders *UntriggeredConditionalOrders) AddUntriggeredConditionalOrder(order types.Order) {
	order.MustBeConditionalOrder()

	if order.IsTakeProfitOrder() {
		if order.IsBuy() {
			untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice = append(
				untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice,
				order,
			)
		} else {
			untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice = append(
				untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice,
				order,
			)
		}
	}

	if order.IsStopLossOrder() {
		if order.IsBuy() {
			untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice = append(
				untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice,
				order,
			)
		} else {
			untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice = append(
				untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice,
				order,
			)
		}
	}
}

// PruneUntriggeredConditionalOrders takes in lists of expired and cancelled stateful order ids and removes
// all respective orders from the in-memory `UntriggeredConditionalOrders` data structure. This data structure
// stores untriggered orders in a map of ClobPairId -> []Order, so we first group orders by ClobPairId and then
// call `UntriggeredConditionalOrders.RemoveExpiredUntriggeredConditionalOrders` on each ClobPairId.
func (k Keeper) PruneUntriggeredConditionalOrders(
	expiredStatefulOrderIds []types.OrderId,
	cancelledStatefulOrderIds []types.OrderId,
) {
	// Merge lists of order ids.
	orderIdsToPrune := lib.UniqueSliceToSet(expiredStatefulOrderIds)
	for _, orderId := range cancelledStatefulOrderIds {
		if _, exists := orderIdsToPrune[orderId]; exists {
			panic(
				fmt.Sprintf(
					"PruneUntriggeredConditionalOrders: duplicate order id %+v in expired and "+
						"cancelled order lists", orderId,
				),
			)
		}
		orderIdsToPrune[orderId] = struct{}{}
	}

	prunableUntriggeredConditionalOrderIdsByClobPair := make(map[types.ClobPairId][]types.OrderId)
	for orderId := range orderIdsToPrune {
		// If the order id is conditional, add to prunable list of untriggered order ids.
		// Triggered conditional orders will be effectively ignored during removal as they are not part of
		// UntriggeredConditionalOrders anymore. No need to filter out here, we can avoid memstore reads.
		if orderId.IsConditionalOrder() {
			clobPairId := types.ClobPairId(orderId.GetClobPairId())
			if _, exists := prunableUntriggeredConditionalOrderIdsByClobPair[clobPairId]; !exists {
				prunableUntriggeredConditionalOrderIdsByClobPair[clobPairId] = []types.OrderId{}
			}

			prunableUntriggeredConditionalOrderIdsByClobPair[clobPairId] = append(
				prunableUntriggeredConditionalOrderIdsByClobPair[clobPairId],
				orderId,
			)
		}
	}

	for clobPairId := range prunableUntriggeredConditionalOrderIdsByClobPair {
		if untriggeredConditionalOrders, exists := k.UntriggeredConditionalOrders[clobPairId]; exists {
			untriggeredConditionalOrders.RemoveUntriggeredConditionalOrders(
				prunableUntriggeredConditionalOrderIdsByClobPair[clobPairId],
			)
			if untriggeredConditionalOrders.IsEmpty() {
				delete(k.UntriggeredConditionalOrders, clobPairId)
			}
		}
	}
}

// RemoveUntriggeredConditionalOrders removes a list of order ids from the `UntriggeredConditionalOrders`
// data structure. This function will panic if the order ids contained involve more than one ClobPairId.
func (untriggeredOrders *UntriggeredConditionalOrders) RemoveUntriggeredConditionalOrders(
	orderIdsToRemove []types.OrderId,
) {
	if len(orderIdsToRemove) == 0 {
		return
	}

	// all orders should have the same ClobPairId
	clobPairId := types.ClobPairId(orderIdsToRemove[0].GetClobPairId())
	for _, orderId := range orderIdsToRemove {
		orderClobPairId := types.ClobPairId(orderId.GetClobPairId())
		if types.ClobPairId(orderId.GetClobPairId()) != clobPairId {
			panic(
				fmt.Sprintf(
					"RemoveExpiredUntriggeredConditionalOrders: expected all orders to have the same ClobPairId. "+
						"Got %v and %v.",
					clobPairId,
					orderClobPairId,
				),
			)
		}
	}

	orderIdsToRemoveSet := lib.UniqueSliceToSet(orderIdsToRemove)

	newOrdersToTriggerWhenOraclePriceLTETriggerPrice := make([]types.Order, 0)
	for _, order := range untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice {
		if _, exists := orderIdsToRemoveSet[order.OrderId]; !exists {
			newOrdersToTriggerWhenOraclePriceLTETriggerPrice = append(newOrdersToTriggerWhenOraclePriceLTETriggerPrice, order)
		}
	}
	untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice = newOrdersToTriggerWhenOraclePriceLTETriggerPrice

	newOrdersToTriggerWhenOraclePriceGTETriggerPrice := make([]types.Order, 0)
	for _, order := range untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice {
		if _, exists := orderIdsToRemoveSet[order.OrderId]; !exists {
			newOrdersToTriggerWhenOraclePriceGTETriggerPrice = append(newOrdersToTriggerWhenOraclePriceGTETriggerPrice, order)
		}
	}
	untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice = newOrdersToTriggerWhenOraclePriceGTETriggerPrice
}

// PollTriggeredConditionalOrders removes all triggered conditional orders from the
// `UntriggeredConditionalOrders` struct given a new oracle price for a clobPairId. It returns
// a list of order ids that were triggered. This is only called in EndBlocker. We round up to the nearest
// subtick int for LTE and down to the nearest subtick int for GTE conditions. This is pessimistic rounding,
// as we want to trigger orders only when we are sure they are triggerable.
func (untriggeredOrders *UntriggeredConditionalOrders) PollTriggeredConditionalOrders(
	oraclePriceSubticksRat *big.Rat,
) []types.OrderId {
	triggeredOrderIds := make([]types.OrderId, 0)
	pessimisticLTESubticks := types.Subticks(lib.BigRatRound(oraclePriceSubticksRat, true).Uint64())
	// For the lte array, find all orders that are triggered when oracle price goes lower
	// than or equal to the trigger price.

	newOrdersToTriggerWhenOraclePriceLTETriggerPrice := make([]types.Order, 0)
	for _, order := range untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice {
		if order.CanTrigger(pessimisticLTESubticks) {
			triggeredOrderIds = append(triggeredOrderIds, order.OrderId)
		} else {
			newOrdersToTriggerWhenOraclePriceLTETriggerPrice = append(
				newOrdersToTriggerWhenOraclePriceLTETriggerPrice,
				order,
			)
		}
	}
	untriggeredOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice = newOrdersToTriggerWhenOraclePriceLTETriggerPrice

	pessimisticGTESubticks := types.Subticks(lib.BigRatRound(oraclePriceSubticksRat, false).Uint64())
	// For the gte array, find all orders that are triggered when oracle price goes greater
	// than or equal to the trigger price.
	newOrdersToTriggerWhenOraclePriceGTETriggerPrice := make([]types.Order, 0)
	for _, order := range untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice {
		if order.CanTrigger(pessimisticGTESubticks) {
			triggeredOrderIds = append(triggeredOrderIds, order.OrderId)
		} else {
			newOrdersToTriggerWhenOraclePriceGTETriggerPrice = append(
				newOrdersToTriggerWhenOraclePriceGTETriggerPrice,
				order,
			)
		}
	}
	untriggeredOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice = newOrdersToTriggerWhenOraclePriceGTETriggerPrice

	return triggeredOrderIds
}

// MaybeTriggerConditionalOrders queries the prices module for price updates and triggers
// any conditional orders in `UntriggeredConditionalOrders` that can be triggered. For each triggered
// order, it takes the stateful order placement stored in Untriggered state and moves it to Triggered state.
// A conditional order trigger event is emitted for each triggered order.
// Function returns a sorted list of conditional order ids that were triggered, intended to be written
// to `ProcessProposerMatchesEvents.ConditionalOrderIdsTriggeredInLastBlock`.
// This function is called in EndBlocker.
func (k Keeper) MaybeTriggerConditionalOrders(ctx sdk.Context) (triggeredConditionalOrderIds []types.OrderId) {
	triggeredConditionalOrderIds = make([]types.OrderId, 0)
	// Sort the keys for the untriggered conditional orders struct. We need to trigger
	// the conditional orders in an ordered way to have deterministic state writes.
	sortedKeys := lib.GetSortedKeys[types.SortedClobPairId](k.UntriggeredConditionalOrders)

	// For all clob pair ids in UntriggeredConditionalOrders, fetch the updated
	// oracle price and poll out triggered conditional orders.
	for _, clobPairId := range sortedKeys {
		untriggeredConditionalOrders := k.UntriggeredConditionalOrders[clobPairId]
		clobPair, found := k.GetClobPair(ctx, clobPairId)
		if !found {
			panic(
				fmt.Errorf(
					"EndBlocker: untriggeredConditionalOrders failed to find clobPairId %+v",
					clobPairId,
				),
			)
		}
		currentOraclePriceSubticksRat := k.GetOraclePriceSubticksRat(ctx, clobPair)
		triggeredOrderIds := untriggeredConditionalOrders.PollTriggeredConditionalOrders(
			currentOraclePriceSubticksRat,
		)
		triggeredConditionalOrderIds = append(triggeredConditionalOrderIds, triggeredOrderIds...)
		// Set the modified untriggeredConditionalOrders back on the keeper field.
		k.UntriggeredConditionalOrders[clobPairId] = untriggeredConditionalOrders
	}

	// State write - move the conditional order placement in state from untriggered to triggered state.
	// Emit an event for each triggered conditional order.
	for _, triggeredConditionalOrderId := range triggeredConditionalOrderIds {
		k.MustTriggerConditionalOrder(
			ctx,
			triggeredConditionalOrderId,
		)
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeStatefulOrder,
			indexerevents.StatefulOrderEventVersion,
			indexer_manager.GetBytes(
				indexerevents.NewConditionalOrderTriggeredEvent(
					triggeredConditionalOrderId,
				),
			),
		)
	}
	return triggeredConditionalOrderIds
}
