package keeper

import (
	"fmt"

	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/clob/types"
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
	UntriggeredLTEOraclePriceConditionalOrders []types.Order
	// All untriggered take profit sell orders and stop loss buy orders sorted by time priority.
	// These orders will be triggered when the oracle price goes greater than or equal to the trigger price.
	// This array functions like a min heap.
	UntriggeredGTEOraclePriceConditionalOrders []types.Order
}

func NewUntriggeredConditionalOrders() *UntriggeredConditionalOrders {
	return &UntriggeredConditionalOrders{
		UntriggeredLTEOraclePriceConditionalOrders: make([]types.Order, 0),
		UntriggeredGTEOraclePriceConditionalOrders: make([]types.Order, 0),
	}
}

// AddUntriggeredConditionalOrder adds an untriggered conditional order to the UntriggeredConditionalOrders
// data structure. It will panic if the order is not a conditional order.
func (untriggeredOrders *UntriggeredConditionalOrders) AddUntriggeredConditionalOrder(order types.Order) {
	order.MustBeConditionalOrder()

	if order.IsTakeProfitOrder() {
		if order.IsBuy() {
			untriggeredOrders.UntriggeredLTEOraclePriceConditionalOrders = append(
				untriggeredOrders.UntriggeredLTEOraclePriceConditionalOrders,
				order,
			)
		} else {
			untriggeredOrders.UntriggeredGTEOraclePriceConditionalOrders = append(
				untriggeredOrders.UntriggeredGTEOraclePriceConditionalOrders,
				order,
			)
		}
	}

	if order.IsStopLossOrder() {
		if order.IsBuy() {
			untriggeredOrders.UntriggeredGTEOraclePriceConditionalOrders = append(
				untriggeredOrders.UntriggeredGTEOraclePriceConditionalOrders,
				order,
			)
		} else {
			untriggeredOrders.UntriggeredLTEOraclePriceConditionalOrders = append(
				untriggeredOrders.UntriggeredLTEOraclePriceConditionalOrders,
				order,
			)
		}
	}
}

// RemoveExpiredUntriggeredConditionalOrders removes a list of expired order ids from the `UntriggeredConditionalOrders`
// data structure.
func (untriggeredOrders *UntriggeredConditionalOrders) RemoveExpiredUntriggeredConditionalOrders(
	expiredOrders []types.OrderId,
) {
	expiredOrderSet := lib.SliceToSet(expiredOrders)

	newUntriggeredLTEOraclePriceConditionalOrders := make([]types.Order, 0)
	for _, order := range untriggeredOrders.UntriggeredLTEOraclePriceConditionalOrders {
		if _, exists := expiredOrderSet[order.OrderId]; !exists {
			newUntriggeredLTEOraclePriceConditionalOrders = append(newUntriggeredLTEOraclePriceConditionalOrders, order)
		}
	}
	untriggeredOrders.UntriggeredLTEOraclePriceConditionalOrders = newUntriggeredLTEOraclePriceConditionalOrders

	newUntriggeredGTEOraclePriceConditionalOrders := make([]types.Order, 0)
	for _, order := range untriggeredOrders.UntriggeredGTEOraclePriceConditionalOrders {
		if _, exists := expiredOrderSet[order.OrderId]; !exists {
			newUntriggeredGTEOraclePriceConditionalOrders = append(newUntriggeredGTEOraclePriceConditionalOrders, order)
		}
	}
	untriggeredOrders.UntriggeredGTEOraclePriceConditionalOrders = newUntriggeredGTEOraclePriceConditionalOrders
}

// PollTriggeredConditionalOrders removes all triggered conditional orders from the
// `UntriggeredConditionalOrders` struct given a new oracle price for a clobPairId. It returns
// a list of orders that were triggered. This is only called in EndBlocker.
func (untriggeredOrders *UntriggeredConditionalOrders) PollTriggeredConditionalOrders(
	currentSubticks types.Subticks,
) []types.Order {
	triggeredLTEOrders := make([]types.Order, 0)
	// For the lte array, find all orders that are triggered when oracle price goes lower
	// than or equal to the trigger price.
	newUntriggeredLTEOraclePriceConditionalOrders := make([]types.Order, 0)
	for _, order := range untriggeredOrders.UntriggeredLTEOraclePriceConditionalOrders {
		if order.CanTrigger(currentSubticks) {
			triggeredLTEOrders = append(triggeredLTEOrders, order)
		} else {
			newUntriggeredLTEOraclePriceConditionalOrders = append(
				newUntriggeredLTEOraclePriceConditionalOrders,
				order,
			)
		}
	}
	untriggeredOrders.UntriggeredLTEOraclePriceConditionalOrders = newUntriggeredLTEOraclePriceConditionalOrders

	triggeredGTEOrders := make([]types.Order, 0)
	// For the gte array, find all orders that are triggered when oracle price goes greater
	// than or equal to the trigger price.
	newUntriggeredGTEOraclePriceConditionalOrders := make([]types.Order, 0)
	for _, order := range untriggeredOrders.UntriggeredGTEOraclePriceConditionalOrders {
		if order.CanTrigger(currentSubticks) {
			triggeredGTEOrders = append(triggeredGTEOrders, order)
		} else {
			newUntriggeredGTEOraclePriceConditionalOrders = append(
				newUntriggeredGTEOraclePriceConditionalOrders,
				order,
			)
		}
	}
	untriggeredOrders.UntriggeredGTEOraclePriceConditionalOrders = newUntriggeredGTEOraclePriceConditionalOrders

	if len(triggeredGTEOrders) > 0 && len(triggeredLTEOrders) > 0 {
		panic(
			fmt.Errorf(
				"PollTriggeredConditionalOrders: orders triggered from both lte and gte trigger arrays. "+
					"gte orders: %+v, lte orders: %+v, oracle price: %+v subticks",
				triggeredGTEOrders,
				triggeredLTEOrders,
				currentSubticks,
			),
		)
	}
	if len(triggeredGTEOrders) > 0 {
		return triggeredGTEOrders
	}
	return triggeredLTEOrders
}
