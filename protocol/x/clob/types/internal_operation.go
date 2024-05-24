package types

func NewOrderPlacementInternalOperation(order Order) InternalOperation {
	return InternalOperation{
		Operation: &InternalOperation_OrderPlacement{
			OrderPlacement: NewMsgPlaceOrder(order),
		},
	}
}

// TODO: Operation for liquidation order and deleveraging?

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
