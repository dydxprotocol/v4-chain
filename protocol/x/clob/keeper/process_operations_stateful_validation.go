package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
			return order, sdkerrors.Wrapf(
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
			return order, sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful long term order id %+v does not exist in state.",
				orderId,
			)
		}
		return statefulOrderPlacement.Order, nil
	} else if orderId.IsConditionalOrder() {
		conditionalOrderPlacement, found := k.GetTriggeredConditionalOrderPlacement(ctx, orderId)
		if !found {
			return order, sdkerrors.Wrapf(
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

// StatefulValidateProposedOperationMatch performs stateful validation on a match orders object.
// The following validations are performed:
// - Validation on all match orders
// - Validation on all match liquidations
// - Validation that delevearging matches have a valid perpetual id
func (k Keeper) StatefulValidateProposedOperationMatch(
	ctx sdk.Context,
	clobMatch *types.ClobMatch,
	shortTermOrdersMap map[types.OrderId]types.Order,
) error {
	switch castedMatch := clobMatch.Match.(type) {
	case *types.ClobMatch_MatchOrders:
		matchOrder := castedMatch.MatchOrders
		if err := k.StatefulValidateProposedOperationMatchOrders(
			ctx,
			matchOrder,
			shortTermOrdersMap,
		); err != nil {
			return err
		}
	case *types.ClobMatch_MatchPerpetualLiquidation:
		matchLiquidation := castedMatch.MatchPerpetualLiquidation
		if err := k.StatefulValidateProposedOperationMatchLiquidation(
			ctx,
			matchLiquidation,
			shortTermOrdersMap,
		); err != nil {
			return err
		}
	case *types.ClobMatch_MatchPerpetualDeleveraging:
		perpId := castedMatch.MatchPerpetualDeleveraging.GetPerpetualId()
		_, err := k.perpetualsKeeper.GetPerpetual(ctx, perpId)
		if err != nil {
			return sdkerrors.Wrapf(
				types.ErrPerpetualDoesNotExist,
				"Perpetual id %+v does not exist in state.",
				perpId,
			)
		}
	default:
		panic(
			fmt.Sprintf(
				"StatefulValidateProposedOperationMatch: Unrecognized operation type for match: %+v",
				clobMatch,
			),
		)
	}
	return nil
}

// StatefulValidateProposedOperationMatchOrders performs stateful validation on a match orders object.
// The following validations are performed:
// - Validation on any short term orders.
// - Validation on all maker fills.
// - Validation that taker order cannot be post only.
func (k Keeper) StatefulValidateProposedOperationMatchOrders(
	ctx sdk.Context,
	matchOrder *types.MatchOrders,
	shortTermOrdersMap map[types.OrderId]types.Order,
) error {
	takerOrderId := matchOrder.GetTakerOrderId()
	// Fetch the taker order from either short term orders or state
	takerOrder, err := k.FetchOrderFromOrderId(ctx, takerOrderId, shortTermOrdersMap)
	if err != nil {
		return err
	}

	// Taker order cannot be post only.
	if takerOrder.GetTimeInForce() == types.Order_TIME_IN_FORCE_POST_ONLY {
		return sdkerrors.Wrapf(
			types.ErrInvalidMatchOrder,
			"Taker order %+v cannot be post only.",
			takerOrder.GetOrderTextString(),
		)
	}

	fills := matchOrder.GetFills()
	for _, fill := range fills {
		if err := k.StatefulValidateMakerFill(ctx, &fill, shortTermOrdersMap, &takerOrder); err != nil {
			return err
		}
	}
	return nil
}

// StatefulValidateProposedOperationMatchLiquidation performs stateful validation on a match liquidation.
// The following validations are performed:
// - Validation on the maker fills
// - Validation that clob pair id, perpetual id exists
func (k Keeper) StatefulValidateProposedOperationMatchLiquidation(
	ctx sdk.Context,
	matchLiquidation *types.MatchPerpetualLiquidation,
	shortTermOrdersMap map[types.OrderId]types.Order,
) error {
	fills := matchLiquidation.GetFills()
	for _, fill := range fills {
		if err := k.StatefulValidateMakerFill(ctx, &fill, shortTermOrdersMap, nil); err != nil {
			return err
		}
	}
	perpId := matchLiquidation.GetPerpetualId()
	_, err := k.perpetualsKeeper.GetPerpetual(ctx, perpId)
	if err != nil {
		return sdkerrors.Wrapf(
			types.ErrPerpetualDoesNotExist,
			"Perpetual id %+v does not exist in state.",
			perpId,
		)
	}
	clobPair := matchLiquidation.ClobPairId
	if _, found := k.GetClobPair(ctx, types.ClobPairId(clobPair)); !found {
		return sdkerrors.Wrapf(
			types.ErrInvalidClob,
			"Clob Pair id %+v does not exist in state.",
			clobPair,
		)
	}
	return nil
}

// StatefulValidateMakerFill performs stateful validation on a maker fill.
// The following validations are performed:
// - Validation on any short term orders
// - Validation that maker order cannot be FOK or IOC
// - Taker and Maker must be on opposite sides
func (k Keeper) StatefulValidateMakerFill(
	ctx sdk.Context,
	fill *types.MakerFill,
	shortTermOrdersMap map[types.OrderId]types.Order,
	takerOrder *types.Order,
) error {
	makerOrderId := fill.GetMakerOrderId()
	// Fetch the maker order from either short term orders or state
	makerOrder, err := k.FetchOrderFromOrderId(ctx, makerOrderId, shortTermOrdersMap)
	if err != nil {
		return err
	}

	// Orders must be on different sides of the book.
	if takerOrder != nil {
		if takerOrder.IsBuy() == makerOrder.IsBuy() {
			return sdkerrors.Wrapf(
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
		return sdkerrors.Wrapf(
			types.ErrInvalidMatchOrder,
			"Maker order %+v cannot be FOK or IOC.",
			makerOrder.GetOrderTextString(),
		)
	}
	return nil
}
