package types

import (
	fmt "fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4/lib"
)

const TypeMsgProposedOperations = "proposed_operations"

var _ sdk.Msg = &MsgProposedOperations{}

func (msg *MsgProposedOperations) GetSigners() []sdk.AccAddress {
	// Return empty slice because app-injected msg is not expected to be signed.
	return []sdk.AccAddress{}
}

// GetAddToOrderbookCollatCheckOrderHashesSet returns a set of order hashes from the
// list of order hashes in the `MsgProposedOperations` message. Note this function will
// panic if there are any duplicates.
func (msg *MsgProposedOperations) GetAddToOrderbookCollatCheckOrderHashesSet() map[OrderHash]bool {
	bytesSlice := lib.MapSlice(
		msg.AddToOrderbookCollatCheckOrderHashes,
		lib.BytesSliceToBytes32,
	)

	addToOrderbookCollatCheckOrderHashesSet := make(
		map[OrderHash]bool,
		len(bytesSlice),
	)
	for _, b := range bytesSlice {
		orderHash := OrderHash(b)
		if _, exists := addToOrderbookCollatCheckOrderHashesSet[orderHash]; exists {
			panic(
				fmt.Sprintf(
					"GetAddToOrderbookCollatCheckOrderHashesSet: duplicate order hash in AddToOrderbookCollatCheckOrderHashes: %+v",
					msg.AddToOrderbookCollatCheckOrderHashes,
				),
			)
		}
		addToOrderbookCollatCheckOrderHashesSet[orderHash] = true
	}

	return addToOrderbookCollatCheckOrderHashesSet
}

// operationsQueueValidator encapsulates all the previous context we need to validate sequential
// operations in the operations queue since the order of operations matters.
type operationsQueueValidator struct {

	// All the previous orders placed in this block (short and long term).
	// This field is used when ensuring short term OrderIds references an order in the last block.
	// ordersPlacedInBlock stores the most recently placed order.
	// It tracks orders placed via `OrderPlacement` operations.
	ordersPlacedInBlock map[OrderId]Order
	// preExistingStatefulOrders is a set of pre-existing stateful orders placed in a previous block.
	// It records orders placed via `PreExistingStatefulOrder` operations.
	preExistingStatefulOrders map[OrderId]struct{}
}

// ValidateBasic performs stateless validation on the proposed operation queue.
// Validations differ based on operation types.
// TODO(CLOB-510): Validate the `add_to_orderbook_collat_check_order_hashes` field.
func (msg *MsgProposedOperations) ValidateBasic() error {
	operations := msg.GetOperationsQueue()

	validator := operationsQueueValidator{
		ordersPlacedInBlock:       make(map[OrderId]Order, 0),
		preExistingStatefulOrders: make(map[OrderId]struct{}, 0),
	}

	// Go through the operations one by one to validate them, updating state as necessary.
	for _, operation := range operations {
		var err error
		switch operation.Operation.(type) {
		case *Operation_Match:
			match := operation.GetMatch()
			err = validator.validateMatchOperation(match)
		case *Operation_OrderPlacement:
			orderPlacement := operation.GetOrderPlacement()
			err = validator.validateOrderPlacementOperation(orderPlacement)
		case *Operation_OrderCancellation:
			orderCancellation := operation.GetOrderCancellation()
			err = validator.validateOrderCancellationOperation(orderCancellation)
		case *Operation_PreexistingStatefulOrder:
			preExistingStatefulOrder := operation.GetPreexistingStatefulOrder()
			err = validator.validatePreExistingOrderOperation(preExistingStatefulOrder)
		default:
			err = fmt.Errorf("operation Queue type not implemented yet for operation %v", operation)
		}
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidMsgProposedOperations, "Error: %+v", err)
		}
	}
	return nil
}

// validatePreExistingOrderOperation performs stateless validation on a PreExisting Stateful Order operation.
// It also populates the validator object with the order.
// This validation does not perform any state reads, or memclob reads.
//
// The following validation occurs in this method:
//
//   - Validate for Order Id
//   - Validate that the OrderId is stateful.
//   - Validate that there are no duplicate PreExisting Stateful Order OrderIds.
func (validator *operationsQueueValidator) validatePreExistingOrderOperation(orderId *OrderId) error {
	if err := orderId.Validate(); err != nil {
		return err
	}
	if !orderId.IsStatefulOrder() {
		return sdkerrors.Wrapf(
			ErrInvalidOrderFlag,
			"Invalid Preexisting Order Operation: OrderId %+v is not stateful.",
			*orderId,
		)
	}

	if _, exists := validator.preExistingStatefulOrders[*orderId]; exists {
		return sdkerrors.Wrapf(
			ErrInvalidMsgProposedOperations,
			"Duplicate Pre Existing Order Operation: OrderId %+v",
			*orderId,
		)
	}

	// Record the pre existing order in validator.
	validator.preExistingStatefulOrders[*orderId] = struct{}{}

	return nil
}

// validateMatchOperation unwraps the match message and performs validation for each match type.
func (validator *operationsQueueValidator) validateMatchOperation(match *ClobMatch) error {
	switch match.Match.(type) {
	case *ClobMatch_MatchOrders:
		matchOrders := match.GetMatchOrders()
		if err := validator.validateMatchOrdersOperation(matchOrders); err != nil {
			return err
		}
	case *ClobMatch_MatchPerpetualLiquidation:
		matchPerpetualLiquidation := match.GetMatchPerpetualLiquidation()
		if err := validator.validateMatchPerpetualLiquidationOperation(matchPerpetualLiquidation); err != nil {
			return err
		}
	case *ClobMatch_MatchPerpetualDeleveraging:
		matchPerpetualDeleveraging := match.GetMatchPerpetualDeleveraging()
		if err := matchPerpetualDeleveraging.Validate(); err != nil {
			return sdkerrors.Wrapf(
				err,
				"match: %+v",
				matchPerpetualDeleveraging,
			)
		}
	default:
		panic("Unsupported Clob Match type")
	}
	return nil
}

// validateOrderCancellationOperation performs stateless validation on an order cancellation.
// It also removes the order from the validator object.
// This validation does not perform any state reads, or memclob reads.
//
// The following validation occurs in this method:
//
//   - ValidateBasic for OrderCancellation message
//   - Short term order cancellations must reference an OrderId that was previously placed.
func (validator *operationsQueueValidator) validateOrderCancellationOperation(
	orderCancellation *MsgCancelOrder,
) error {
	// Order cancellation msg has to pass its own validation.
	if err := orderCancellation.ValidateBasic(); err != nil {
		return err
	}
	orderCancellationOrderId := orderCancellation.GetOrderId()

	// A short-term order cancellation must reference an OrderId that was previously placed.
	if orderCancellationOrderId.IsShortTermOrder() {
		err := validator.verifyOrderPlacementInOperationsQueue(
			orderCancellationOrderId,
		)
		if err != nil {
			return err
		}
	}

	// Delete the order Id from the validator order maps.
	delete(validator.ordersPlacedInBlock, orderCancellationOrderId)
	delete(validator.preExistingStatefulOrders, orderCancellationOrderId)

	return nil
}

// validateOrderPlacementOperation performs stateless validation on an order placement.
// It also populates the validator object with the order.
// This validation does not perform any state reads, or memclob reads.
//
// The following validation occurs in this method:
//
//   - ValidateBasic for OrderPlacement message
//   - Orders placed in the same block with same OrderId must not be the same.
func (validator *operationsQueueValidator) validateOrderPlacementOperation(
	orderPlacement *MsgPlaceOrder,
) error {
	// Order placement msg has to pass its own validation.
	if err := orderPlacement.ValidateBasic(); err != nil {
		return err
	}

	order := orderPlacement.GetOrder()
	orderId := order.GetOrderId()

	// For orders with the same orderId placed within this block, verify replacement order priority.
	if prevOrder, placedPreviously := validator.ordersPlacedInBlock[orderId]; placedPreviously {
		// No duplicate order placements allowed.
		if prevOrder.MustCmpReplacementOrder(&order) == 0 {
			return sdkerrors.Wrapf(
				ErrInvalidPlaceOrder,
				"Duplicate Order %+v",
				order,
			)
		}
		// Replacement Orders have a higher priority than the previously placed order that it replaces.
		// All short term replacement orders should be checked here. Note that for long term orders,
		// this check only takes effect if the order being replaced is in the same block.
		if prevOrder.MustCmpReplacementOrder(&order) != -1 {
			return sdkerrors.Wrapf(
				ErrInvalidReplacement,
				"Replacement order is not higher priority. order: %+v, prevOrder: %+v",
				order,
				prevOrder,
			)
		}
	}

	// Record the placed order in validator.
	validator.ordersPlacedInBlock[orderId] = order

	return nil
}

// validateMatchOrdersOperation performs stateless validation on an match orders.
// This validation does not perform any state reads, or memclob reads.
//
// The following validation occurs in this method:
//   - For all fills, The fill amount is not zero.
//   - For all fills, maker order ids must be previously placed in an operation.
//   - Taker order id must be previously placed in an operation.
//   - There are no duplicate MakerOrderIds in fills.
func (validator *operationsQueueValidator) validateMatchOrdersOperation(
	matchOrders *MatchOrders,
) error {
	fills := matchOrders.GetFills()
	makerOrderIdSet := make(map[OrderId]struct{}, len(fills))
	takerOrderId := matchOrders.GetTakerOrderId()
	takerOrderHash := matchOrders.GetTakerOrderHash()

	if len(takerOrderHash) != 32 {
		return sdkerrors.Wrapf(
			ErrInvalidMatchOrder,
			"taker order %+v has invalid order hash of length %d, expected 32",
			takerOrderId,
			len(takerOrderHash),
		)
	}

	for _, fill := range fills {
		// Fill amount must be greater than zero.
		if fill.GetFillAmount() == 0 {
			return ErrFillAmountIsZero
		}
		makerOrderId := fill.GetMakerOrderId()

		if err := makerOrderId.Validate(); err != nil {
			return err
		}

		// No duplicate maker order IDs in fills.
		if _, exists := makerOrderIdSet[makerOrderId]; exists {
			return sdkerrors.Wrapf(
				ErrInvalidMatchOrder,
				"duplicate Maker OrderId in a MatchOrder's fills, maker: %+v, taker %+v",
				makerOrderId,
				takerOrderId,
			)
		}

		// Maker order id must be previously placed in an operation.
		if err := validator.verifyOrderPlacementInOperationsQueue(
			makerOrderId,
		); err != nil {
			return err
		}
		makerOrderIdSet[makerOrderId] = struct{}{}
	}

	if err := takerOrderId.Validate(); err != nil {
		return err
	}

	// Taker order id must be previously placed in an operation.
	if err := validator.verifyOrderPlacementInOperationsQueue(takerOrderId); err != nil {
		return err
	}

	return nil
}

// validateMatchPerpetualLiquidationOperation performs stateless validation on a liquidation match.
// This validation does not perform any state reads, or memclob reads.
//
// The following validation occurs in this method:
//   - The total size of liquidation order is not zero.
//   - FillAmounts is not zero.
//   - For all fills, maker order ids must be previously placed in an operation.
//   - The sum of all fill_amount entries in the list of fills is not greater than the total size.
func (validator *operationsQueueValidator) validateMatchPerpetualLiquidationOperation(
	liquidationMatch *MatchPerpetualLiquidation,
) error {
	fills := liquidationMatch.GetFills()
	totalSize := liquidationMatch.GetTotalSize()

	// Make sure the total size greater than zero.
	if liquidationMatch.GetTotalSize() == 0 {
		return sdkerrors.Wrapf(
			ErrInvalidLiquidationOrderTotalSize,
			"Liquidation match total size is zero. match: %+v",
			liquidationMatch,
		)
	}

	// Make sure the sum of all fill_amount entries in the list of fills does not exceed the total size.
	// Get the total quantums filled for this liquidation order.
	bigQuantumsFilled := new(big.Int)

	for _, fill := range fills {
		fillAmt := fill.GetFillAmount()
		// Fill amount cannot be zero.
		if fillAmt == 0 {
			return ErrFillAmountIsZero
		}
		bigQuantumsFilled.Add(bigQuantumsFilled, new(big.Int).SetUint64(fill.FillAmount))

		fillMakerOrderId := fill.GetMakerOrderId()

		if err := fillMakerOrderId.Validate(); err != nil {
			return err
		}

		// Maker order id must be previously placed in an operation.
		if err := validator.verifyOrderPlacementInOperationsQueue(
			fillMakerOrderId,
		); err != nil {
			return err
		}
	}
	if bigQuantumsFilled.Cmp(new(big.Int).SetUint64(totalSize)) == 1 {
		return sdkerrors.Wrapf(
			ErrTotalFillAmountExceedsOrderSize,
			"Total fill size: %v match total size: %v",
			bigQuantumsFilled,
			totalSize,
		)
	}

	return nil
}

// verifyOrderPlacementInOperationsQueue is a pure function. For the referenced order, it checks:
//   - If the order is a stateful or short-term order included in the operations queue for this block.
//   - If the order is a pre-existing stateful order.
//
// If neither of these conditions are met, an `ErrOrderPlacementNotInOperationsQueue` is returned.
func (validator *operationsQueueValidator) verifyOrderPlacementInOperationsQueue(orderId OrderId) error {
	// First, check if the order is placed in `ordersPlacedInBlock`.
	if _, prevPlaced := validator.ordersPlacedInBlock[orderId]; prevPlaced {
		return nil
	}
	// If it is a short term order, we can return early because short term orders
	// cannot exist in `preExistingStatefulOrders`.
	if orderId.IsShortTermOrder() {
		return sdkerrors.Wrapf(ErrOrderPlacementNotInOperationsQueue, "short term orderId: %v", orderId)
	}
	// Second, check if the stateful order exists in `preExistingStatefulOrders`.
	if _, prevPlaced := validator.preExistingStatefulOrders[orderId]; prevPlaced {
		return nil
	}
	return sdkerrors.Wrapf(ErrOrderPlacementNotInOperationsQueue, "stateful orderId: %v", orderId)
}
