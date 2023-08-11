package types

import (
	"fmt"
	"sort"

	"github.com/dydxprotocol/v4/lib"
)

// Nonce is a nonce tied to an operation in the operation queue.
// It is used to sort the operations queue by nonce in ascending order.
type Nonce uint64

// OperationsToPropose is a struct encapsulating data required for determining the operations
// to propose in a block.
// TODO(DEC-1889): Fix hack where struct fields are public.
type OperationsToPropose struct {
	// This map will represent the nonce for all operations that will be included
	// in the next proposed block, where an operation is any of the following:
	// - An order placement that was part of a valid match (stateful or Short-Term).
	// - A pre-existing stateful order that was part of a valid match.
	// - A new stateful order placement that was successfully added to the orderbook.
	// - A valid stateful order cancellation.
	// - Order matches, specifically regular and liquidation matches.
	// - An operation that caused an order included in the operations to propose to be
	//   removed from the book. This could happen in the following instances:
	//   - Short-Term order cancellation of an order included in the operations to propose.
	//   - Order placement causing an order in the operations to propose to be removed from
	//     the book, for example an order placement causing a maker order in the
	//     operations to propose to become undercollateralized, or a new order placement that
	//     replaces a partially-matched maker order.
	NonceToOperationToPropose map[Nonce]Operation
	// This map represents a mapping from operation hash to the nonce for that operation.
	// The nonce is used for the relative ordering of all operations in
	// the operations to propose.
	OperationHashToNonce map[OperationHash]Nonce
	// A counter used for tracking the next available nonce to assign to new operations
	// that are inserted into `operationHashToNonce`.
	NextAvailableNonce Nonce
	// A set of order hashes that have gone through the add-to-orderbook collateralization check in
	// the last block.
	AddToOrderbookCollatCheckOrders map[OrderHash]bool
}

// NewOperationsToPropose returns a new instance of `OperationsToPropose`.
func NewOperationsToPropose() *OperationsToPropose {
	return &OperationsToPropose{
		NonceToOperationToPropose:       make(map[Nonce]Operation, 0),
		OperationHashToNonce:            make(map[OperationHash]Nonce),
		NextAvailableNonce:              0,
		AddToOrderbookCollatCheckOrders: make(map[OrderHash]bool, 0),
	}
}

// AssignNonceToOrder assigns a nonce to an order.
// This function panics if the order has already been assigned a nonce.
// TODO(DEC-1890): Update this code to verify stateful orders are not assigned nonces as
// pre-existing stateful orders and newly-placed stateful orders.
func (o *OperationsToPropose) AssignNonceToOrder(
	order Order,
	isPreexistingStatefulOrder bool,
) {
	operation := o.getOrderPlacementOperation(order, isPreexistingStatefulOrder)
	o.assignNonceToOperation(operation)
}

// MustAddToOrderbookCollatCheckOrders inserts an order into `AddToOrderbookCollatCheckOrders`
// to mark that the add-to-orderbook collateralization check was performed in the last block on
// this order.
func (o *OperationsToPropose) MustAddToOrderbookCollatCheckOrders(
	order Order,
) {
	orderHash := order.GetOrderHash()
	if _, exists := o.AddToOrderbookCollatCheckOrders[orderHash]; exists {
		panic(
			fmt.Sprintf(
				"MustAddToOrderbookCollatCheckOrders: order (%+v) already exists in "+
					"AddToOrderbookCollatCheckOrders",
				order.GetOrderTextString(),
			),
		)
	}

	o.AddToOrderbookCollatCheckOrders[orderHash] = true
}

// RemovePreexistingStatefulOrderPlacementNonce removes a pre-existing order nonce from the
// `operationsHashToNonce` data structure.
// This function will panic if the operation has no nonce, exists in the operations to propose,
// or if the order is not a stateful order.
func (o *OperationsToPropose) RemovePreexistingStatefulOrderPlacementNonce(
	order Order,
) {
	operation := o.getOrderPlacementOperation(order, true)
	o.mustRemoveOperationNonce(operation)
}

// RemoveOrderPlacementNonce removes an order nonce from the `operationsHashToNonce` data structure.
// This function will panic if the operation has no nonce or exists in the operations to propose.
func (o *OperationsToPropose) RemoveOrderPlacementNonce(
	order Order,
) {
	operation := o.getOrderPlacementOperation(order, false)
	o.mustRemoveOperationNonce(operation)
}

// mustRemoveOperationNonce removes an operation nonce from the `operationsHashToNonce` data structure.
// This function will panic if the operation has no nonce or exists in the operations to propose.
func (o *OperationsToPropose) mustRemoveOperationNonce(
	operation Operation,
) {
	operationHash := operation.GetOperationHash()
	operationNonce, exists := o.OperationHashToNonce[operationHash]
	if !exists {
		panic(
			fmt.Sprintf(
				"mustRemoveOperationNonce: nonce for operation (%+v) does not exist",
				operation.GetOperationTextString(),
			),
		)
	}

	// If the operation nonce exists in the operations to propose, then the nonce cannot be dropped.
	if _, existsInOperationsToPropose := o.NonceToOperationToPropose[operationNonce]; existsInOperationsToPropose {
		panic(
			fmt.Sprintf(
				"mustRemoveOperationNonce: operation (%+v) has nonce %d in operations to propose",
				operation.GetOperationTextString(),
				operationNonce,
			),
		)
	}

	delete(o.OperationHashToNonce, operationHash)
}

// getOrderPlacementOperation gets the correct operation depending on if the order is a pre-existing
// stateful order or not. This function will panic if the order is a Short-Term order and
// `isPreexistingStatefulOrder` is true.
func (o *OperationsToPropose) getOrderPlacementOperation(
	order Order,
	isPreexistingStatefulOrder bool,
) Operation {
	if isPreexistingStatefulOrder {
		order.MustBeStatefulOrder()
		return NewPreexistingStatefulOrderPlacementOperation(order)
	} else {
		return NewOrderPlacementOperation(order)
	}
}

// orderMustHaveNonce panics if the provided order does not have a nonce assigned.
func (o *OperationsToPropose) orderMustHaveNonce(order Order) {
	o.mustIsPreexistingStatefulOrder(order)
}

// mustIsPreexistingStatefulOrder returns true if the order is a pre-existing stateful order, and
// false otherwise. panics if the order does not have a nonce.
func (o *OperationsToPropose) mustIsPreexistingStatefulOrder(
	order Order,
) (isPreexistingStatefulOrder bool) {
	// Check if the maker order is an order placement operation.
	placementOp := NewOrderPlacementOperation(order)
	if order.IsShortTermOrder() {
		o.mustGetNonceFromOperation(placementOp)
		return false
	}

	// At this point the order must be a stateful order.
	// Check if the order is a newly-placed stateful order.
	_, nonceExists := o.OperationHashToNonce[placementOp.GetOperationHash()]
	if nonceExists {
		return false
	}

	// If the operation wasn't an order placement operation, check
	// if it's a pre-existing order placement operation.
	// If it's not, panic.
	preexistingOp := NewPreexistingStatefulOrderPlacementOperation(order)
	o.mustGetNonceFromOperation(preexistingOp)

	return true
}

// IsMakerOrderPreexistingStatefulOrder returns true if the maker order is a pre-existing stateful
// order, or false if it's a newly-placed stateful order.
// This method panics if the provided order is not a stateful order.
// This method panics if the order has not been assigned a nonce.
func (o *OperationsToPropose) IsMakerOrderPreexistingStatefulOrder(
	order Order,
) bool {
	order.MustBeStatefulOrder()
	return o.mustIsPreexistingStatefulOrder(order)
}

// assignNonceToOperation assigns a nonce to an operation.
// This function panics if the operation has already been assigned a nonce.
// TODO(DEC-1890): Update this code to verify stateful orders are not assigned nonces as
// pre-existing stateful orders and newly-placed stateful orders.
func (o *OperationsToPropose) assignNonceToOperation(operation Operation) {
	operationHash := operation.GetOperationHash()

	// Panic if the operation hash has already been assigned a nonce.
	if nonce, exists := o.OperationHashToNonce[operationHash]; exists {
		panic(
			fmt.Sprintf(
				"assignNonceToOperation: operation (%+v) has already been assigned nonce %d. "+
					"The current operations queue is: %s",
				operation.GetOperationTextString(),
				nonce,
				GetOperationsQueueTextString(o.GetOperationsQueue()),
			),
		)
	}

	// Assign the operation hash to the next available nonce and increment the next available nonce.
	o.OperationHashToNonce[operationHash] = o.NextAvailableNonce
	o.NextAvailableNonce++
}

// AddPreexistingStatefulOrderPlacementToOperationsQueue adds a preexisting stateful order placement
// (a stateful order placed in a previous block) into the operations to propose. The order passed in *must*
// already have a nonce assigned to it. Order placement operations will be inserted in the list of operations
// to propose in ascending nonce order. Function will panic if order hash is not present in
// `OperationHashToNonce` or if the order is not stateful.
func (o *OperationsToPropose) AddPreexistingStatefulOrderPlacementToOperationsQueue(
	order Order,
) {
	order.MustBeStatefulOrder()
	operation := NewPreexistingStatefulOrderPlacementOperation(order)
	o.insertOperationIntoOperationsToPropose(operation)
}

// AddOrderPlacementToOperationsQueue adds an order placement into the operations to propose.
// The order passed in *must* already have a nonce assigned to it. Order placement operations
// will be inserted in the list of operations according to ascending nonce order. Order placement
// can be either short term or stateful.
// Function will panic if order hash is not present in `OperationHashToNonce`.
func (o *OperationsToPropose) AddOrderPlacementToOperationsQueue(
	order Order,
) {
	operation := NewOrderPlacementOperation(order)
	o.insertOperationIntoOperationsToPropose(operation)
}

// AddOrderCancellationToOperationsQueue inserts an order cancellation operation into the
// operations to propose. Function panics if the order cancellation operation has a nonce.
func (o *OperationsToPropose) AddOrderCancellationToOperationsQueue(cancel MsgCancelOrder) {
	// Assign a nonce to the order cancellation operation.
	operation := NewOrderCancellationOperation(&cancel)
	o.assignNonceToOperation(operation)

	// Insert the order cancellation into the operations to propose.
	o.insertOperationIntoOperationsToPropose(operation)
}

// AddMatchToOperationsQueue creates a match operation, assigns a nonce to the match operation,
// and adds the match operation to the operations queue.
// The `takerMatchableOrder` can also be a liquidation order.
// Note that liquidation orders themselves should not / cannot be assigned nonces.
// This function will panic if the taker order or any of the maker orders are not present in
// `operationHashToNonce`.
func (o *OperationsToPropose) AddMatchToOperationsQueue(
	takerMatchableOrder MatchableOrder,
	makerFillsWithOrders []MakerFillWithOrder,
) {
	// If the order is not a liquidation order, it must have a nonce.
	if !takerMatchableOrder.IsLiquidation() {
		o.orderMustHaveNonce(takerMatchableOrder.MustGetOrder())
	}

	// Every maker order must have a nonce in the operations queue.
	for _, makerOrder := range makerFillsWithOrders {
		o.orderMustHaveNonce(makerOrder.Order)
	}

	// Create the match operation.
	matchOperation := NewMatchOperation(
		takerMatchableOrder,
		MakerFillsWithOrderToMakerFills(makerFillsWithOrders),
	)

	// Insert the match operation into the operations to propose.
	o.assignNonceToOperation(matchOperation)
	o.insertOperationIntoOperationsToPropose(matchOperation)
}

// IsOrderPlacementInOperationsQueue returns true if an order's corresponding order placement exists
// in the operations queue, otherwise returns false. This function is necessary for determining
// whether to add regular / liquidation order placements and cancellations to the operations queue,
// even if they didn't generate any matches.
// This function should not be called for pre-existing stateful orders, and will panic if the order
// placement operation does not have a nonce.
func (o *OperationsToPropose) IsOrderPlacementInOperationsQueue(
	order Order,
) bool {
	operation := o.getOrderPlacementOperation(order, false)
	return o.isOperationInOperationsQueue(operation)
}

// IsPreexistingStatefulOrderInOperationsQueue returns true if an order's corresponding order placement exists
// in the operations queue, otherwise returns false. This function is necessary for determining
// whether to add regular / liquidation order placements and cancellations to the operations queue,
// even if they didn't generate any matches.
// This function should only be called with pre-existing stateful orders, and will panic
// if called with a non-stateful order or the pre-existing stateful order placement operation does
// not have a nonce.
func (o *OperationsToPropose) IsPreexistingStatefulOrderInOperationsQueue(
	order Order,
) bool {
	operation := o.getOrderPlacementOperation(order, true)
	return o.isOperationInOperationsQueue(operation)
}

// DoesOperationHaveNonce returns a boolean indicating whether the provided operation has a nonce,
// true if the operation has a nonce and false if not.
// TODO(DEC-1948): Add unit test coverage to this method.
func (o *OperationsToPropose) DoesOperationHaveNonce(operation Operation) bool {
	_, exists := o.OperationHashToNonce[operation.GetOperationHash()]
	return exists
}

// isOperationInOperationsQueue returns whether the provided operation exists in the operations queue,
// true if it does and false if not.
// This function will panic if the provided operation does not have a nonce.
func (o *OperationsToPropose) isOperationInOperationsQueue(operation Operation) bool {
	nonce, nonceExists := o.OperationHashToNonce[operation.GetOperationHash()]
	if !nonceExists {
		panic(
			fmt.Sprintf(
				"isOperationInOperationsQueue: operation (%+v) has no nonce",
				operation.GetOperationTextString(),
			),
		)
	}

	_, operationExistsInOperationsQueue := o.NonceToOperationToPropose[nonce]
	return operationExistsInOperationsQueue
}

// mustGetNonceFromOperation fetches the assigned nonce for the provided operation from the
// `OperationHashToNonce` data structure. Function panics if the operation doesn't have
// a nonce.
func (o *OperationsToPropose) mustGetNonceFromOperation(operation Operation) Nonce {
	nonce, nonceExists := o.OperationHashToNonce[operation.GetOperationHash()]
	if !nonceExists {
		panic(
			fmt.Sprintf(
				"mustGetNonceFromOperation: operation (%+v) has no nonce",
				operation.GetOperationTextString(),
			),
		)
	}
	return nonce
}

// insertOperationIntoOperationsToPropose inserts an operation into the operations to propose with
// the nonce returned from `mustGetNonceFromOperation`. Function panics if the operation doesn't have
// a nonce.
func (o *OperationsToPropose) insertOperationIntoOperationsToPropose(operation Operation) {
	nonceToAdd := o.mustGetNonceFromOperation(operation)

	if existingOperation, exists := o.NonceToOperationToPropose[nonceToAdd]; exists {
		panic(
			fmt.Sprintf(
				"insertOperationIntoOperationsToPropose: an operation with nonce %d already exists "+
					"in the operations to propose. New operation: (%+v). Existing operation: (%+v).",
				nonceToAdd,
				operation.GetOperationTextString(),
				existingOperation.GetOperationTextString(),
			),
		)
	}

	o.NonceToOperationToPropose[nonceToAdd] = operation
}

// GetOperationsQueue returns the operations queue stored in `NonceToOperationToPropose`.
// The Operation slice that is returned is sorted in ascending order by Nonce.
func (o *OperationsToPropose) GetOperationsQueue() []Operation {
	nonces := lib.ConvertMapToSliceOfKeys(o.NonceToOperationToPropose)
	sort.Slice(nonces, func(i, j int) bool {
		return nonces[i] < nonces[j]
	})

	operations := make([]Operation, len(o.NonceToOperationToPropose))
	index := 0
	for _, nonce := range nonces {
		operations[index] = o.NonceToOperationToPropose[nonce]
		index += 1
	}
	return operations
}

// GetOrdersWithAddToOrderbookCollatCheck returns a sorted slice of order hashes of all orders that had the
// add-to-orderbook collateralization check performed in the last block.
// TODO(CLOB-513): Only include order hashes that were included in the proposed operations queue.
// TODO(CLOB-514): Add unit test coverage for this method.
func (o *OperationsToPropose) GetOrdersWithAddToOrderbookCollatCheck() []OrderHash {
	addToOrderbookCollatCheckOrderHashes := lib.ConvertMapToSliceOfKeys(o.AddToOrderbookCollatCheckOrders)
	sort.Slice(addToOrderbookCollatCheckOrderHashes, func(i, j int) bool {
		return string(addToOrderbookCollatCheckOrderHashes[i][:]) < string(addToOrderbookCollatCheckOrderHashes[j][:])
	})

	return addToOrderbookCollatCheckOrderHashes
}

// ClearOperationsQueue clears out all operations in the operations queue and all entries in the
// `NonceToOperationToPropose` and `AddToOrderbookCollatCheckOrders` maps. For each operation in
// `NonceToOperationToPropose`, the operation hash is removed from `OperationHashToNonce`. Note
// that this is *not* a full clear of the `OperationsToPropose` internal data structures because
// maker orders can rest across multiple blocks and require nonces if matched.
// Function panics if an operation has no corresponding `OperationHashToNonce` value or if there is a
// mismatch between an operation's nonce in `NonceToOperationToPropose` and `OperationHashToNonce`.
func (o *OperationsToPropose) ClearOperationsQueue() {
	for operationNonce, operation := range o.NonceToOperationToPropose {
		hash := operation.GetOperationHash()
		operationAssignedNonce, exists := o.OperationHashToNonce[hash]
		if !exists {
			panic(
				fmt.Sprintf(
					"ClearOperationsQueue: No nonce to remove for operation %+v",
					operation.GetOperationTextString(),
				),
			)
		}
		if operationNonce != operationAssignedNonce {
			panic(
				fmt.Sprintf(
					"ClearOperationsQueue: Mismatch between nonces for operation (%+v). "+
						"Assigned nonce: %d, nonce in operation queue: %d",
					operation.GetOperationTextString(),
					operationAssignedNonce,
					operationNonce,
				),
			)
		}
		delete(o.OperationHashToNonce, hash)
	}

	// Now that all operation nonce's are deleted from `OperationHashToNonce`, clear out the
	// operations to propose to effectively clear the operations queue.
	o.NonceToOperationToPropose = make(map[Nonce]Operation, 0)

	// Clear all entries from `AddToOrderbookCollatCheckOrders` to signify that all orders currently
	// resting on the book should not have the add-to-orderbook collateralization check performed in
	// in the next proposed block.
	o.AddToOrderbookCollatCheckOrders = make(map[OrderHash]bool, 0)
}
