package types

import (
	"fmt"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// NewShortTermOrderPlacementInternalOperation returns a new internal operation for placing
// a Short-Term order.
// This function will panic if it's called with a non Short-Term order.
func NewShortTermOrderPlacementInternalOperation(order Order) InternalOperation {
	order.OrderId.MustBeShortTermOrder()
	return InternalOperation{
		Operation: &InternalOperation_ShortTermOrderPlacement{
			ShortTermOrderPlacement: NewMsgPlaceOrder(order),
		},
	}
}

// NewPreexistingStatefulOrderPlacementInternalOperation returns a new internal operation for placing
// a stateful order.
// This function will panic if it's called with a non stateful order.
func NewPreexistingStatefulOrderPlacementInternalOperation(order Order) InternalOperation {
	order.OrderId.MustBeStatefulOrder()
	return InternalOperation{
		Operation: &InternalOperation_PreexistingStatefulOrder{
			PreexistingStatefulOrder: &order.OrderId,
		},
	}
}

// NewMatchOrdersInternalOperation returns a new operation for matching maker orders against a
// taker order.
// This function panics if there are zero maker fills.
func NewMatchOrdersInternalOperation(
	takerOrder Order,
	makerFills []MakerFill,
) InternalOperation {
	if len(makerFills) == 0 {
		panic(
			fmt.Sprintf(
				"NewMatchOrdersInternalOperation: cannot create a match orders "+
					"internal operation with no maker fills: %+v",
				takerOrder,
			),
		)
	}

	return InternalOperation{
		Operation: &InternalOperation_Match{
			Match: &ClobMatch{
				Match: &ClobMatch_MatchOrders{
					MatchOrders: &MatchOrders{
						TakerOrderId: takerOrder.OrderId,
						Fills:        makerFills,
					},
				},
			},
		},
	}
}

// NewMatchPerpetualLiquidationInternalOperation returns a new operation for matching maker orders
// against a perpetual liquidation order.
// This function panics if this is called with a non-liquidation order or there are zero maker fills.
func NewMatchPerpetualLiquidationInternalOperation(
	takerLiquidationOrder MatchableOrder,
	makerFills []MakerFill,
) InternalOperation {
	if !takerLiquidationOrder.IsLiquidation() {
		panic(
			fmt.Sprintf(
				"NewMatchPerpetualLiquidationInternalOperation: called with a non-liquidation order: %+v",
				takerLiquidationOrder,
			),
		)
	}

	if len(makerFills) == 0 {
		panic(
			fmt.Sprintf(
				"NewMatchPerpetualLiquidationInternalOperation: cannot create a match perpetual "+
					"liquidation internal operation with no maker fills: %+v",
				takerLiquidationOrder,
			),
		)
	}

	return InternalOperation{
		Operation: &InternalOperation_Match{
			Match: &ClobMatch{
				Match: &ClobMatch_MatchPerpetualLiquidation{
					MatchPerpetualLiquidation: &MatchPerpetualLiquidation{
						Liquidated:  takerLiquidationOrder.GetSubaccountId(),
						ClobPairId:  takerLiquidationOrder.GetClobPairId().ToUint32(),
						PerpetualId: takerLiquidationOrder.MustGetLiquidatedPerpetualId(),
						TotalSize:   takerLiquidationOrder.GetBaseQuantums().ToUint64(),
						IsBuy:       takerLiquidationOrder.IsBuy(),
						Fills:       makerFills,
					},
				},
			},
		},
	}
}

// NewMatchPerpetualDeleveragingInternalOperation returns a new operation for deleveraging liquidated subaccount's
// position against one or more offsetting subaccounts.
func NewMatchPerpetualDeleveragingInternalOperation(
	liquidatedSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	fills []MatchPerpetualDeleveraging_Fill,
	isFinalSettlement bool,
) InternalOperation {
	return InternalOperation{
		Operation: &InternalOperation_Match{
			Match: &ClobMatch{
				Match: &ClobMatch_MatchPerpetualDeleveraging{
					MatchPerpetualDeleveraging: &MatchPerpetualDeleveraging{
						Liquidated:        liquidatedSubaccountId,
						PerpetualId:       perpetualId,
						Fills:             fills,
						IsFinalSettlement: isFinalSettlement,
					},
				},
			},
		},
	}
}

// NewOrderRemovalInternalOperation returns a new operation for removing an order.
// This function panics if it's called with an order removal containing an OrderId
// for a non stateful order or the removal reason is unspecified.
func NewOrderRemovalInternalOperation(
	orderId OrderId,
	removalReason OrderRemoval_RemovalReason,
) InternalOperation {
	orderId.MustBeStatefulOrder()

	if removalReason == OrderRemoval_REMOVAL_REASON_UNSPECIFIED {
		panic("NewOrderRemovalInternalOperation: removal reason unspecified")
	}

	orderRemoval := OrderRemoval{
		OrderId:       orderId,
		RemovalReason: removalReason,
	}
	return InternalOperation{
		Operation: &InternalOperation_OrderRemoval{
			OrderRemoval: &orderRemoval,
		},
	}
}
