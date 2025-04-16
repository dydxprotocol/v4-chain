package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
	if orderId.IsLongTermOrder() || orderId.IsTwapSuborder() {
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

// ValidateLiquidationOrderAgainstProposedLiquidation performs stateless validation of a liquidation order
// against a proposed liquidation.
// An error is returned when
//   - The CLOB pair IDs of the order and proposed liquidation do not match.
//   - The perpetual IDs of the order and proposed liquidation do not match.
//   - The total size of the order and proposed liquidation do not match.
//   - The side of the order and proposed liquidation do not match.
func (k Keeper) ValidateLiquidationOrderAgainstProposedLiquidation(
	ctx sdk.Context,
	order *types.LiquidationOrder,
	proposedMatch *types.MatchPerpetualLiquidation,
) error {
	if order.GetClobPairId() != types.ClobPairId(proposedMatch.GetClobPairId()) {
		return errorsmod.Wrapf(
			types.ErrClobPairAndPerpetualDoNotMatch,
			"Order CLOB Pair ID: %v, Match CLOB Pair ID: %v",
			order.GetClobPairId(),
			proposedMatch.GetClobPairId(),
		)
	}

	if order.MustGetLiquidatedPerpetualId() != proposedMatch.GetPerpetualId() {
		return errorsmod.Wrapf(
			types.ErrClobPairAndPerpetualDoNotMatch,
			"Order Perpetual ID: %v, Match Perpetual ID: %v",
			order.MustGetLiquidatedPerpetualId(),
			proposedMatch.GetPerpetualId(),
		)
	}

	if order.GetBaseQuantums() != satypes.BaseQuantums(proposedMatch.TotalSize) {
		return errorsmod.Wrapf(
			types.ErrInvalidLiquidationOrderTotalSize,
			"Order Size: %v, Match Size: %v",
			order.GetBaseQuantums(),
			proposedMatch.TotalSize,
		)
	}

	if order.IsBuy() != proposedMatch.GetIsBuy() {
		return errorsmod.Wrapf(
			types.ErrInvalidLiquidationOrderSide,
			"Order Side: %v, Match Side: %v",
			order.IsBuy(),
			proposedMatch.GetIsBuy(),
		)
	}
	return nil
}
