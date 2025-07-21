package types

import (
	fmt "fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgProposedOperations = "proposed_operations"

var _ sdk.Msg = &MsgProposedOperations{}

// Stateless validation for MsgProposedOperations is located in ValidateAndTransformRawOperations.
func (msg *MsgProposedOperations) ValidateBasic() error {
	// Go through the operations one by one to validate them, updating state as necessary.
	for _, rawOperation := range msg.GetOperationsQueue() {
		switch operation := rawOperation.Operation.(type) {
		case
			*OperationRaw_Match,
			*OperationRaw_ShortTermOrderPlacement:
			// no-op, stateless validation is done in ValidateAndTransformRawOperations
		case *OperationRaw_OrderRemoval:
			orderId := operation.OrderRemoval.GetOrderId()
			if orderId.IsShortTermOrder() {
				return errorsmod.Wrapf(
					ErrInvalidMsgProposedOperations,
					"order removal is not allowed for short-term orders: %v",
					orderId,
				)
			}

			switch operation.OrderRemoval.RemovalReason {
			case OrderRemoval_REMOVAL_REASON_UNSPECIFIED:
				return errorsmod.Wrapf(
					ErrInvalidMsgProposedOperations,
					"order removal reason must be specified: %v",
					orderId,
				)
			case OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY:
				// Reduce-only order removals are now supported.
				// Validation will be performed in PersistOrderRemovalToState.
			}

		default:
			return errorsmod.Wrapf(
				ErrInvalidMsgProposedOperations,
				"operation queue type not implemented yet for raw operation %v",
				rawOperation,
			)
		}
	}
	return nil
}

// ValidateAndTransformRawOperations performs stateless validation on the proposed operation queue
// and transforms the input []OperationRaw into []InternalOperation.
// Validations differ based on operation types. We are able to supply a TxDecoder and AnteHandler
// to this function. These are needed to decode OperationRaw tx bytes and to validate that
// the operations' transactions were signed correctly.
func ValidateAndTransformRawOperations(
	ctx sdk.Context,
	rawOperations []OperationRaw,
	decoder sdk.TxDecoder,
	anteHandler sdk.AnteHandler,
) ([]InternalOperation, error) {
	operations := make([]InternalOperation, 0, len(rawOperations))

	validator := operationsQueueValidator{
		ordersPlacedInBlock: make(map[OrderId]Order, 0),
	}

	// Go through the operations one by one to validate them, updating state as necessary.
	for _, rawOperation := range rawOperations {
		var err error
		operation := &InternalOperation{}
		switch rawOperation.Operation.(type) {
		case *OperationRaw_Match:
			match := rawOperation.GetMatch()
			if err = validator.validateMatchOperation(match); err != nil {
				return nil, err
			}
			operation.Operation = &InternalOperation_Match{
				Match: match,
			}
		case *OperationRaw_ShortTermOrderPlacement:
			operation, err = decodeOperationRawShortTermOrderPlacementBytes(
				ctx,
				rawOperation.GetShortTermOrderPlacement(),
				decoder,
				anteHandler,
			)
			if err != nil {
				return nil, err
			}
			if err = validator.validateShortTermOrderPlacementOperation(
				operation.GetShortTermOrderPlacement(),
			); err != nil {
				return nil, err
			}
		case *OperationRaw_OrderRemoval:
			orderRemoval := rawOperation.GetOrderRemoval()
			if err := orderRemoval.OrderId.Validate(); err != nil {
				return nil, err
			}
			// Order removal reason fully filled is only used by indexer and should not be
			// placed in the operations queue.
			if orderRemoval.RemovalReason == OrderRemoval_REMOVAL_REASON_UNSPECIFIED ||
				orderRemoval.RemovalReason == OrderRemoval_REMOVAL_REASON_FULLY_FILLED {
				return nil, errorsmod.Wrapf(
					ErrInvalidOrderRemoval,
					"Invalid order removal reason: %+v",
					orderRemoval.RemovalReason,
				)
			}
			operation.Operation = &InternalOperation_OrderRemoval{
				OrderRemoval: rawOperation.GetOrderRemoval(),
			}
		default:
			return nil, fmt.Errorf("operation queue type not implemented yet for raw operation %v", rawOperation)
		}

		operations = append(operations, *operation)
	}
	return operations, nil
}

// operationsQueueValidator encapsulates all the previous context we need to validate sequential
// operations in the operations queue since the order of operations matters.
type operationsQueueValidator struct {
	// All the previous orders placed in this block (short and long term).
	// This field is used when ensuring short term OrderIds references an order in the last block.
	// ordersPlacedInBlock stores the most recently placed order.
	// It tracks orders placed via `OrderPlacement` operations.
	ordersPlacedInBlock map[OrderId]Order
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
			return errorsmod.Wrapf(
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

// validateShortTermOrderPlacementOperation performs stateless validation on an order placement.
// It also populates the validator object with the order.
// This validation does not perform any state reads, or memclob reads.
//
// The following validation occurs in this method:
//
//   - ValidateBasic for OrderPlacement message
//   - Orders placed in the same block with same OrderId must not be the same.
func (validator *operationsQueueValidator) validateShortTermOrderPlacementOperation(
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
			return errorsmod.Wrapf(
				ErrInvalidPlaceOrder,
				"Duplicate Order %s",
				order.GetOrderTextString(),
			)
		}
		// Replacement Orders have a higher priority than the previously placed order that it replaces.
		if prevOrder.MustCmpReplacementOrder(&order) != -1 {
			return errorsmod.Wrapf(
				ErrInvalidReplacement,
				"Replacement order is not higher priority. order: %s, prevOrder: %s",
				order.GetOrderTextString(),
				prevOrder.GetOrderTextString(),
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
//   - Match has at least one fill.
//   - For all fills, The fill amount is not zero.
//   - For all fills, maker order ids must be previously placed in an operation.
//   - Taker order id must be previously placed in an operation.
//   - There are no duplicate MakerOrderIds in fills.
func (validator *operationsQueueValidator) validateMatchOrdersOperation(
	matchOrders *MatchOrders,
) error {
	fills := matchOrders.GetFills()
	if len(fills) == 0 {
		return errorsmod.Wrapf(
			ErrInvalidMatchOrder,
			"Match has no fills: %+v",
			matchOrders,
		)
	}

	makerOrderIdSet := make(map[OrderId]struct{}, len(fills))
	takerOrderId := matchOrders.GetTakerOrderId()

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
			return errorsmod.Wrapf(
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
//   - Liquidation match has at least one fill.
//   - For all fills, maker order ids must be previously placed in an operation.
//   - The sum of all fill_amount entries in the list of fills is not greater than the total size.
func (validator *operationsQueueValidator) validateMatchPerpetualLiquidationOperation(
	liquidationMatch *MatchPerpetualLiquidation,
) error {
	fills := liquidationMatch.GetFills()
	if len(fills) == 0 {
		return errorsmod.Wrapf(
			ErrInvalidMatchOrder,
			"Liquidation match has no fills: %+v",
			liquidationMatch,
		)
	}

	// Make sure the total size greater than zero.
	totalSize := liquidationMatch.GetTotalSize()
	if totalSize == 0 {
		return errorsmod.Wrapf(
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

	if err := liquidationMatch.Liquidated.Validate(); err != nil {
		return err
	}

	if bigQuantumsFilled.Cmp(new(big.Int).SetUint64(totalSize)) == 1 {
		return errorsmod.Wrapf(
			ErrTotalFillAmountExceedsOrderSize,
			"Total fill size: %v match total size: %v",
			bigQuantumsFilled,
			totalSize,
		)
	}

	return nil
}

// verifyOrderPlacementInOperationsQueue is a pure function. For the referenced order, it checks:
//   - If the order is a short-term order, that it is included in the operations queue for this block.
//
// If this condition isn't met, an `ErrOrderPlacementNotInOperationsQueue` is returned.
func (validator *operationsQueueValidator) verifyOrderPlacementInOperationsQueue(orderId OrderId) error {
	if orderId.IsShortTermOrder() {
		if _, prevPlaced := validator.ordersPlacedInBlock[orderId]; !prevPlaced {
			return errorsmod.Wrapf(ErrOrderPlacementNotInOperationsQueue, "short term orderId: %v", orderId)
		}
	}

	return nil
}
