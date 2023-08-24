package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k Keeper) StatefulValidateProposedOperations(
	ctx sdk.Context,
	operations []types.InternalOperation,
) error {
	// Collect all the short-term orders placed for subsequent lookups.
	placedShortTermOrders := make(map[types.OrderId]types.Order, 0)

	for _, operation := range operations {
		switch castedOperation := operation.Operation.(type) {
		case *types.InternalOperation_Match:
			clobMatch := castedOperation.Match
			if err := k.StatefulValidateProposedOperationMatch(ctx, clobMatch, placedShortTermOrders); err != nil {
				return err
			}
		case *types.InternalOperation_ShortTermOrderPlacement:
			order := castedOperation.ShortTermOrderPlacement.GetOrder()
			placedShortTermOrders[order.GetOrderId()] = order
		case *types.InternalOperation_OrderRemoval:
			// Order removals are always for stateful orders that must exist.
			orderId := castedOperation.OrderRemoval.OrderId
			_, found := k.GetLongTermOrderPlacement(ctx, orderId)
			if !found {
				return sdkerrors.Wrapf(
					types.ErrStatefulOrderDoesNotExist,
					"Stateful order id %+v does not exist in state.",
					orderId,
				)
			}
		case *types.InternalOperation_PreexistingStatefulOrder:
			// When we fetch operations to propose, preexisting stateful orders are not included
			// in the operations queue.
			panic(
				fmt.Sprintf(
					"StatefulProposedOperationsValidation: Preexisting Stateful Orders should not exist in operations queue: %+v",
					castedOperation.PreexistingStatefulOrder,
				),
			)
		default:
			panic(
				fmt.Sprintf(
					"StatefulProposedOperationsValidation: Unrecognized operation type for operation: %+v",
					operation.GetInternalOperationTextString(),
				),
			)
		}
	}
	return nil
}

// FetchOrderFromOrderId is a helper function that fetches a order from and order id.
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

func (k Keeper) StatefulValidateProposedOperationMatch(
	ctx sdk.Context,
	clobMatch *types.ClobMatch,
	shortTermOrdersMap map[types.OrderId]types.Order,
) error {
	switch castedMatch := clobMatch.Match.(type) {
	case *types.ClobMatch_MatchOrders:
		matchOrder := castedMatch.MatchOrders
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
				takerOrder,
			)
		}

		fills := matchOrder.GetFills()
		for _, fill := range fills {
			makerOrderId := fill.GetMakerOrderId()
			// Fetch the maker order from either short term orders or state
			makerOrder, err := k.FetchOrderFromOrderId(ctx, makerOrderId, shortTermOrdersMap)
			if err != nil {
				return err
			}

			// Maker order cannot be FOK or IOC.
			if makerOrder.GetTimeInForce() == types.Order_TIME_IN_FORCE_FILL_OR_KILL ||
				makerOrder.GetTimeInForce() == types.Order_TIME_IN_FORCE_IOC {
				return sdkerrors.Wrapf(
					types.ErrInvalidMatchOrder,
					"Maker order %+v cannot be FOK or IOC.",
					makerOrder,
				)
			}
		}

	case *types.ClobMatch_MatchPerpetualLiquidation:
		matchLiquidation := castedMatch.MatchPerpetualLiquidation
		fills := matchLiquidation.GetFills()
		for _, fill := range fills {
			makerOrderId := fill.GetMakerOrderId()
			if makerOrderId.IsStatefulOrder() {
				_, found := k.GetLongTermOrderPlacement(ctx, makerOrderId)
				if !found {
					return sdkerrors.Wrapf(
						types.ErrStatefulOrderDoesNotExist,
						"Stateful order id %+v does not exist in state.",
						makerOrderId,
					)
				}
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
				types.ErrPerpetualDoesNotExist,
				"Clob Pair id %+v does not exist in state.",
				clobPair,
			)
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
