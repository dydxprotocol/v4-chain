package types

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// OperationsToPropose is a struct encapsulating data required for determining the operations
// to propose in a block.
// TODO(DEC-1889): Fix hack where struct fields are public.
type OperationsToPropose struct {
	// An ordered list of operations to propose in a block.
	OperationsQueue []InternalOperation
	// A set of order hashes where the order placement has already been included in the operations queue.
	OrderHashesInOperationsQueue map[OrderHash]bool
	// A map of Short-Term order hashes to the raw transaction bytes of the order placement.
	// This is used in `GetOperationsQueueRaw` for returning a slice of `OperationRaw` for
	// the purposes of constructing `MsgProposedOperations`.
	ShortTermOrderHashToTxBytes map[OrderHash][]byte
	// A map from order ID to the orders themselves for each order that
	// was matched. Note: there may be multiple distinct orders with the same
	// ID that are matched. In that case, only the "greatest" of any such orders
	// is maintained in this map.
	MatchedOrderIdToOrder map[OrderId]Order
	// A set of order ids where the order removal has already been included in the operations queue.
	OrderRemovalsInOperationsQueue map[OrderId]bool
}

// NewOperationsToPropose returns a new instance of `OperationsToPropose`.
func NewOperationsToPropose() *OperationsToPropose {
	return &OperationsToPropose{
		OperationsQueue:                make([]InternalOperation, 0),
		OrderHashesInOperationsQueue:   make(map[OrderHash]bool),
		ShortTermOrderHashToTxBytes:    make(map[OrderHash][]byte),
		MatchedOrderIdToOrder:          make(map[OrderId]Order),
		OrderRemovalsInOperationsQueue: make(map[OrderId]bool),
	}
}

// ClearOperationsQueue clears out all operations in the operations queue and all entries in the
// `OrderHashesInOperationsQueue` set. Note that we don't clear `ShortTermOrderHashToTxBytes` because
// it is updated in the `mustRemoveOrder` and `RemoveAndClearOperationsQueue` functions.
func (o *OperationsToPropose) ClearOperationsQueue() {
	o.OperationsQueue = make([]InternalOperation, 0)
	o.OrderHashesInOperationsQueue = make(map[OrderHash]bool, 0)
	o.MatchedOrderIdToOrder = make(map[OrderId]Order)
	o.OrderRemovalsInOperationsQueue = make(map[OrderId]bool, 0)
}

// MustAddShortTermOrderTxBytes adds the provided Short-Term order hash and TX bytes into
// `ShortTermOrderHashToTxBytes`.
// This function will panic if the provided order is not a Short-Term order or this order already
// exists in `ShortTermOrderHashToTxBytes`.
func (o *OperationsToPropose) MustAddShortTermOrderTxBytes(
	order Order,
	txBytes []byte,
) {
	order.OrderId.MustBeShortTermOrder()

	if len(txBytes) == 0 {
		panic(
			fmt.Sprintf(
				"MustAddShortTermOrderTxBytes: provided TX bytes are empty for order (%s).",
				order.GetOrderTextString(),
			),
		)
	}

	orderHash := order.GetOrderHash()
	if _, exists := o.ShortTermOrderHashToTxBytes[orderHash]; exists {
		panic(
			fmt.Sprintf(
				"MustAddShortTermOrderTxBytes: Order (%s) already exists in `ShortTermOrderHashToTxBytes`.",
				order.GetOrderTextString(),
			),
		)
	}

	o.ShortTermOrderHashToTxBytes[orderHash] = txBytes
}

// MustAddShortTermOrderPlacementToOperationsQueue adds a Short-Term order placement operation to the
// operations queue.
// This function will panic if the order is not a Short-Term order, the order already exists in
// `OrderHashesInOperationsQueue`, or the order does not exist in `ShortTermOrderHashToTxBytes`.
func (o *OperationsToPropose) MustAddShortTermOrderPlacementToOperationsQueue(
	order Order,
) {
	order.OrderId.MustBeShortTermOrder()

	orderHash := order.GetOrderHash()
	if _, exists := o.OrderHashesInOperationsQueue[orderHash]; exists {
		panic(
			fmt.Sprintf(
				"MustAddShortTermOrderPlacementToOperationsQueue: Order (%s) already exists in "+
					"`OrderHashesInOperationsQueue`.",
				order.GetOrderTextString(),
			),
		)
	}

	if _, exists := o.ShortTermOrderHashToTxBytes[orderHash]; !exists {
		panic(
			fmt.Sprintf(
				"MustAddShortTermOrderPlacementToOperationsQueue: Order (%s) does not exist in "+
					"`ShortTermOrderHashToTxBytes`.",
				order.GetOrderTextString(),
			),
		)
	}

	o.OrderHashesInOperationsQueue[orderHash] = true
	o.OperationsQueue = append(o.OperationsQueue, NewShortTermOrderPlacementInternalOperation(order))
}

// RemoveShortTermOrderTxBytes removes a short term order from `ShortTermOrderHashToTxBytes`.
// This function will panic for any of the following:
// - the order is not a short term order.
// - the order hash is present in `OrderHashesInOperationsQueue`
// - the order hash is not present in `ShortTermOrderHashToTxBytes`
func (o *OperationsToPropose) RemoveShortTermOrderTxBytes(
	order Order,
) {
	order.OrderId.MustBeShortTermOrder()

	orderHash := order.GetOrderHash()
	if _, exists := o.OrderHashesInOperationsQueue[orderHash]; exists {
		panic(
			fmt.Sprintf(
				"RemoveShortTermOrderTxBytes: Order (%s) exists in "+
					"`OrderHashesInOperationsQueue`.",
				order.GetOrderTextString(),
			),
		)
	}

	if _, exists := o.ShortTermOrderHashToTxBytes[orderHash]; !exists {
		panic(
			fmt.Sprintf(
				"RemoveShortTermOrderTxBytes: Order (%s) does not exist in "+
					"`ShortTermOrderHashToTxBytes`.",
				order.GetOrderTextString(),
			),
		)
	}

	delete(o.ShortTermOrderHashToTxBytes, orderHash)
}

// MustAddStatefulOrderPlacementToOperationsQueue adds a stateful order placement operation to the
// operations queue.
// This function will panic if this stateful order already exists in `OrderHashesInOperationsQueue`,
// or the provided order is not a stateful order.
func (o *OperationsToPropose) MustAddStatefulOrderPlacementToOperationsQueue(
	order Order,
) {
	order.OrderId.MustBeStatefulOrder()

	orderHash := order.GetOrderHash()
	if _, exists := o.OrderHashesInOperationsQueue[orderHash]; exists {
		panic(
			fmt.Sprintf(
				"MustAddStatefulOrderPlacementToOperationsQueue: Order (%s) already exists in "+
					"`OrderHashesInOperationsQueue`.",
				order.GetOrderTextString(),
			),
		)
	}

	o.OrderHashesInOperationsQueue[orderHash] = true
	o.OperationsQueue = append(o.OperationsQueue, NewPreexistingStatefulOrderPlacementInternalOperation(order))
}

// MustAddMatchToOperationsQueue adds a match operation to the
// operations queue.
// This function will panic if any of the following conditions are true:
//   - The taker order is not a liquidation order and the number of maker fills is zero.
//   - The taker order hash is not present in `orderHashesInOperationsQueue`.
//   - There exist filled maker orders where the order hash is not present in
//     `orderHashesInOperationsQueue`.
func (o *OperationsToPropose) MustAddMatchToOperationsQueue(
	takerMatchableOrder MatchableOrder,
	makerFillsWithOrders []MakerFillWithOrder,
) InternalOperation {
	makerFills := lib.MapSlice(
		makerFillsWithOrders,
		func(mfwo MakerFillWithOrder) MakerFill {
			return mfwo.MakerFill
		},
	)

	// Add all maker orders to `ordersMustBeInOpQueue` so we can verify that they are present in
	// the operations queue.
	ordersMustBeInOpQueue := lib.MapSlice(
		makerFillsWithOrders,
		func(mfwo MakerFillWithOrder) Order {
			return mfwo.Order
		},
	)

	// If the order is a liquidation order, create a liquidation match.
	// Else the order is not a liquidation order, create a regular match and verify the taker order
	// is in the operations queue.
	var matchOperation InternalOperation
	if takerMatchableOrder.IsLiquidation() {
		matchOperation = NewMatchPerpetualLiquidationInternalOperation(
			takerMatchableOrder,
			makerFills,
		)
	} else {
		ordersMustBeInOpQueue = append(ordersMustBeInOpQueue, takerMatchableOrder.MustGetOrder())
		matchOperation = NewMatchOrdersInternalOperation(
			takerMatchableOrder.MustGetOrder(),
			makerFills,
		)
	}

	// Ensure each order that should be in the operations queue is in the operations queue.
	for _, order := range ordersMustBeInOpQueue {
		if !o.IsOrderPlacementInOperationsQueue(order) {
			panic(
				fmt.Sprintf(
					"MustAddMatchToOperationsQueue: Order (%s) does not exist in "+
						"`OrderHashesInOperationsQueue`.",
					order.GetOrderTextString(),
				),
			)
		}
	}

	o.OperationsQueue = append(
		o.OperationsQueue,
		matchOperation,
	)
	return matchOperation
}

// AddZeroFillDeleveragingToOperationsQueue adds a zero-fill deleveraging match operation to the
// operations queue.
func (o *OperationsToPropose) AddZeroFillDeleveragingToOperationsQueue(
	liquidatedSubaccountId satypes.SubaccountId,
	perpetualId uint32,
) {
	o.OperationsQueue = append(
		o.OperationsQueue,
		NewMatchPerpetualDeleveragingInternalOperation(
			liquidatedSubaccountId,
			perpetualId,
			[]MatchPerpetualDeleveraging_Fill{},
			false,
		),
	)
}

// MustAddDeleveragingToOperationsQueue adds a deleveraging match operation to the
// operations queue.
// This function will panic if:
//   - The number of maker fills is zero.
//   - The fill amount is zero for any of the maker fills.
//   - Maker fills contain duplicated subaccount IDs.
//   - Maker fills contain liquidated subaccount ID.
func (o *OperationsToPropose) MustAddDeleveragingToOperationsQueue(
	liquidatedSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	fills []MatchPerpetualDeleveraging_Fill,
	isFinalSettlement bool,
) {
	if len(fills) == 0 {
		panic(
			fmt.Sprintf(
				"MustAddDeleveragingToOperationsQueue: number of fills is zero. "+
					"liquidatedSubaccountId = (%+v), perpetualId = (%d)",
				liquidatedSubaccountId,
				perpetualId,
			),
		)
	}

	seenSubaccountIds := make(map[satypes.SubaccountId]bool)
	for _, fill := range fills {
		if fill.FillAmount == 0 {
			panic(
				fmt.Sprintf(
					"MustAddDeleveragingToOperationsQueue: fill amount is zero. "+
						"liquidatedSubaccountId = (%+v), perpetualId = (%d), fill = (%+v)",
					liquidatedSubaccountId,
					perpetualId,
					fill,
				),
			)
		}

		if fill.OffsettingSubaccountId == liquidatedSubaccountId {
			panic(
				fmt.Sprintf(
					"MustAddDeleveragingToOperationsQueue: offsetting subaccount is the same as liquidated subaccount. "+
						"liquidatedSubaccountId = (%+v), perpetualId = (%d), fill = (%+v)",
					liquidatedSubaccountId,
					perpetualId,
					fill,
				),
			)
		}

		if _, ok := seenSubaccountIds[fill.OffsettingSubaccountId]; ok {
			panic(
				fmt.Sprintf(
					"MustAddDeleveragingToOperationsQueue: duplicated subaccount ids. "+
						"liquidatedSubaccountId = (%+v), perpetualId = (%d), fill = (%+v)",
					liquidatedSubaccountId,
					perpetualId,
					fill,
				),
			)
		}

		seenSubaccountIds[fill.OffsettingSubaccountId] = true
	}

	o.OperationsQueue = append(
		o.OperationsQueue,
		NewMatchPerpetualDeleveragingInternalOperation(
			liquidatedSubaccountId,
			perpetualId,
			fills,
			isFinalSettlement,
		),
	)
}

// MustAddOrderRemovalToOperationsQueue adds an order removal operation to the
// operations queue. This function will panic if given an unspecified removal reason.
func (o *OperationsToPropose) MustAddOrderRemovalToOperationsQueue(
	orderId OrderId,
	removalReason OrderRemoval_RemovalReason,
) {
	// unspecified removal reason is not allowed
	if removalReason == OrderRemoval_REMOVAL_REASON_UNSPECIFIED {
		panic("MustAddOrderRemovalToOperationsQueue: removal reason unspecified")
	}

	if _, exists := o.OrderRemovalsInOperationsQueue[orderId]; exists {
		panic("MustAddOrderRemovalToOperationsQueue: order removal already exists in operations queue")
	}

	o.OperationsQueue = append(
		o.OperationsQueue,
		NewOrderRemovalInternalOperation(orderId, removalReason),
	)
	o.OrderRemovalsInOperationsQueue[orderId] = true
}

// IsOrderRemovalInOperationsQueue returns true if the provided order ID is included in
// `OrderRemovalsInOperationsQueue`, false if not.
func (o *OperationsToPropose) IsOrderRemovalInOperationsQueue(
	orderId OrderId,
) bool {
	_, exists := o.OrderRemovalsInOperationsQueue[orderId]
	return exists
}

// IsOrderPlacementInOperationsQueue returns true if the provided order hash is included in
// `orderHashesInOperationsQueue`, false if not. This function should be used for determining if a
// Short-Term or stateful order placement should be added to the operations queue.
func (o *OperationsToPropose) IsOrderPlacementInOperationsQueue(
	order Order,
) bool {
	orderHash := order.GetOrderHash()

	_, exists := o.OrderHashesInOperationsQueue[orderHash]
	return exists
}

// GetOperationsToReplay returns all operations in the operations queue and a map of all Short-Term
// order hashes to their TX bytes.
// Note the returned operations include pre-existing stateful order placements, since those
// operations are only used when replaying a local validator’s operations queue.
// This function will panic if any of the Short-Term order placement operations do not have an
// entry in ShortTermOrderHashToTxBytes, since that is necessary for constructing the list
// of OperationWithTxBytes.
func (o *OperationsToPropose) GetOperationsToReplay() (
	[]InternalOperation,
	map[OrderHash][]byte,
) {
	operations := make([]InternalOperation, 0, len(o.OperationsQueue))
	shortTermOrderTxBytesMap := make(map[OrderHash][]byte)

	for _, operation := range o.OperationsQueue {
		operations = append(operations, operation)

		if shortTermOrder := operation.GetShortTermOrderPlacement(); shortTermOrder != nil {
			orderHash := shortTermOrder.Order.GetOrderHash()
			shortTermOrderBytes, exists := o.ShortTermOrderHashToTxBytes[orderHash]
			if !exists {
				panic(
					fmt.Sprintf(
						"GetOperationsToReplay: Short-Term order (%s) does not exist in "+
							"`ShortTermOrderHashToTxBytes`.",
						shortTermOrder.Order.GetOrderTextString(),
					),
				)
			} else if len(shortTermOrderBytes) == 0 {
				panic(
					fmt.Sprintf(
						"GetOperationsToReplay: Short-Term order (%s) is assigned to an empty byte "+
							"array in `ShortTermOrderHashToTxBytes`.",
						shortTermOrder.Order.GetOrderTextString(),
					),
				)
			}

			// Copy to avoid accidentally modifying underlying TX bytes.
			shortTermOrderTxBytesMapCopy := make([]byte, len(shortTermOrderBytes))
			copy(shortTermOrderTxBytesMapCopy, shortTermOrderBytes)
			shortTermOrderTxBytesMap[orderHash] = shortTermOrderTxBytesMapCopy
			if len(shortTermOrderTxBytesMapCopy) == 0 {
				panic("GetOperationsToReplay: Short-Term order TX bytes are empty.")
			}
		}
	}

	return operations, shortTermOrderTxBytesMap
}

// GetOperationsToPropose returns a slice of OperationRaw.
// Note this function returns all operations in the operations queue *except* pre-existing
// stateful order placements, since those operations are only used when replaying a local
// validator’s operations queue. They do not need to be proposed again.
// This function will panic if any of the Short-Term order placement operations do not have an
// entry in ShortTermOrderHashToTxBytes, since that is necessary for constructing the list
// of OperationRaw.
func (o *OperationsToPropose) GetOperationsToPropose() []OperationRaw {
	operationRaws := make([]OperationRaw, 0)

	for _, operation := range o.OperationsQueue {
		switch operation := operation.Operation.(type) {
		case *InternalOperation_Match:
			operationRaws = append(operationRaws, OperationRaw{
				Operation: &OperationRaw_Match{
					Match: &ClobMatch{
						Match: operation.Match.Match,
					},
				},
			})
		case *InternalOperation_ShortTermOrderPlacement:
			order := operation.ShortTermOrderPlacement.GetOrder()
			operationBytes, exists := o.ShortTermOrderHashToTxBytes[order.GetOrderHash()]
			if !exists {
				panic(
					fmt.Sprintf(
						"GetOperationsToPropose: Order (%s) does not exist in "+
							"`ShortTermOrderHashToTxBytes`.",
						order.GetOrderTextString(),
					),
				)
			}
			operationRaws = append(operationRaws, OperationRaw{
				Operation: &OperationRaw_ShortTermOrderPlacement{
					ShortTermOrderPlacement: operationBytes,
				},
			})
		case *InternalOperation_PreexistingStatefulOrder:
		case *InternalOperation_OrderRemoval:
			operationRaws = append(operationRaws, OperationRaw{
				Operation: &OperationRaw_OrderRemoval{
					OrderRemoval: operation.OrderRemoval,
				},
			})
		default:
			panic(fmt.Sprintf("GetOperationsToReplay: Unrecognized operation: %+v", operation))
		}
	}

	return operationRaws
}

// MustGetShortTermOrderTxBytes returns the `ShortTermOrderHashToTxBytes` for a short term order.
// This function will panic for any of the following:
// - the order is not a short term order.
// - the order hash is not present in `ShortTermOrderHashToTxBytes`
func (o *OperationsToPropose) MustGetShortTermOrderTxBytes(
	order Order,
) (txBytes []byte) {
	order.OrderId.MustBeShortTermOrder()

	orderHash := order.GetOrderHash()
	bytes, exists := o.ShortTermOrderHashToTxBytes[orderHash]
	if !exists {
		panic(
			fmt.Sprintf(
				"MustGetShortTermOrderTxBytes: Order (%s) does not exist in "+
					"`ShortTermOrderHashToTxBytes`.",
				order.GetOrderTextString(),
			),
		)
	}

	return bytes
}
