package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// FetchOrderFromOrderId is a helper function that fetches a order from an order id.
// If the order id is a short term order, the map will be used to populate the order.
// If the order id is a long term order, it will be fetched from state.
// If the order Id is a conditional order, it will be fetched from triggered conditional order state.
func (k Keeper) FetchOrderFromOrderId(
	ctx sdk.Context,
	orderId types.OrderId,
	shortTermOrdersMap map[types.OrderId]types.Order,
) (order types.Order, err error) {
	// In the case of short term orders, fetch from the orders map.
	// It should always exist in the orders map because short term order placement operations
	// should precede operations with order ids that reference them.
	if orderId.IsShortTermOrder() {
		order, exists := shortTermOrdersMap[orderId]
		if !exists {
			return order, errorsmod.Wrapf(
				types.ErrInvalidMatchOrder,
				"Failed fetching short term order id %+v from previous operations in operations queue",
				orderId,
			)
		}
		return order, nil
	}

	// For stateful orders, fetch from state. Do not fetch from untriggered conditional orders.
	if orderId.IsLongTermOrder() {
		statefulOrderPlacement, found := k.GetLongTermOrderPlacement(ctx, orderId)
		if !found {
			return order, errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful long term order id %+v does not exist in state.",
				orderId,
			)
		}
		return statefulOrderPlacement.Order, nil
	} else if orderId.IsConditionalOrder() {
		conditionalOrderPlacement, found := k.GetTriggeredConditionalOrderPlacement(ctx, orderId)
		if !found {
			return order, errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				orderId,
			)
		}
		return conditionalOrderPlacement.Order, nil
	}

	panic(
		fmt.Sprintf("FetchOrderFromOrderId: unknown order type. order id: %+v", orderId),
	)
}

// MustFetchOrderFromOrderId fetches an Order object given an orderId. If it is a short term order,
// `ordersMap` will be used to populate the order. If it is a stateful order, read from state.
// Note that this function is meant to be used for operation processing during DeliverTx and does not
// fetch untriggered conditional orders.
//
// Function will panic if for any reason, the order cannot be searched up.
func (k Keeper) MustFetchOrderFromOrderId(
	ctx sdk.Context,
	orderId types.OrderId,
	ordersMap map[types.OrderId]types.Order,
) types.Order {
	order, err := k.FetchOrderFromOrderId(ctx, orderId, ordersMap)
	if err != nil {
		panic(err)
	}
	return order
}

// StatefulValidateMakerFill performs stateful validation on a maker fill.
// Additionally, it returns the maker order referenced in the fill.
// The following validations are performed:
// - Validation on any short term orders
// - Validation that maker order cannot be FOK or IOC
// - Taker and Maker must be on opposite sides
func (k Keeper) StatefulValidateMakerFill(
	ctx sdk.Context,
	fill *types.MakerFill,
	shortTermOrdersMap map[types.OrderId]types.Order,
	takerOrder *types.Order,
) (makerOrder types.Order, err error) {
	makerOrderId := fill.GetMakerOrderId()
	// Fetch the maker order from either short term orders or state
	makerOrder, err = k.FetchOrderFromOrderId(ctx, makerOrderId, shortTermOrdersMap)
	if err != nil {
		return makerOrder, err
	}

	// Orders must be on different sides of the book.
	if takerOrder != nil {
		if takerOrder.IsBuy() == makerOrder.IsBuy() {
			return makerOrder, errorsmod.Wrapf(
				types.ErrInvalidMatchOrder,
				"Taker Order %+v and Maker order %+v are not on opposing sides of the book",
				takerOrder.GetOrderTextString(),
				makerOrder.GetOrderTextString(),
			)
		}
	}

	// Maker order cannot be FOK or IOC.
	if makerOrder.GetTimeInForce() == types.Order_TIME_IN_FORCE_FILL_OR_KILL ||
		makerOrder.GetTimeInForce() == types.Order_TIME_IN_FORCE_IOC {
		return makerOrder, errorsmod.Wrapf(
			types.ErrInvalidMatchOrder,
			"Maker order %+v cannot be FOK or IOC.",
			makerOrder.GetOrderTextString(),
		)
	}
	return makerOrder, nil
}
