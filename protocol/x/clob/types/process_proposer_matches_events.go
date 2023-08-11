package types

// GetPreviousBlockStatefulOrderCancellations returns a list of stateful OrderIds that were cancelled
// in the previous block.
func (events ProcessProposerMatchesEvents) GetPreviousBlockStatefulOrderCancellations() []OrderId {
	operations := events.OperationsProposedInLastBlock

	statefulOrderCancellations := make([]OrderId, 0)
	statefulOrderPlacements := make(map[OrderId]struct{}, 0)

	// Iterate through operations in reverse, collect all stateful order placement orderIds.
	// If we see a stateful order cancellation that's not followed by a stateful order placement
	// with the same orderId, we consider it a stateful order cancelled in the last block.
	for index := len(operations) - 1; index >= 0; index-- {
		operation := operations[index]
		if orderPlacement := operation.GetOrderPlacement(); orderPlacement != nil {
			order := orderPlacement.GetOrder()
			if order.IsStatefulOrder() {
				statefulOrderPlacements[order.GetOrderId()] = struct{}{}
			}
		}
		if orderCancellation := operation.GetOrderCancellation(); orderCancellation != nil {
			orderIdToCancel := orderCancellation.GetOrderId()
			if orderIdToCancel.IsStatefulOrder() {
				_, exists := statefulOrderPlacements[orderIdToCancel]
				if !exists {
					statefulOrderCancellations = append(statefulOrderCancellations, orderIdToCancel)
				}
			}
		}
	}
	return statefulOrderCancellations
}
