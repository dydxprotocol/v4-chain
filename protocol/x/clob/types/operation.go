package types

import (
	"github.com/cosmos/gogoproto/proto"
)

// OperationHash is used to represent the SHA256 hash of an operation.
type OperationHash [32]byte

// GetOperationHash returns the SHA256 hash of this operation.
func (o *Operation) GetOperationHash() OperationHash {
	protoHash := GetHash(ProtobufHashable(o))
	return OperationHash(protoHash)
}

// NewOrderPlacementOperation returns a new operation for placing an order.
func NewOrderPlacementOperation(order Order) Operation {
	return Operation{
		Operation: &Operation_OrderPlacement{
			OrderPlacement: NewMsgPlaceOrder(order.MustGetOrder()),
		},
	}
}

// NewPreexistingStatefulOrderPlacementOperation returns a new operation for placing a
// pre-existing stateful order.
// Note this function panics if called with a non-stateful order.
func NewPreexistingStatefulOrderPlacementOperation(order Order) Operation {
	order.MustBeStatefulOrder()

	orderId := order.GetOrderId()
	return Operation{
		Operation: &Operation_PreexistingStatefulOrder{
			PreexistingStatefulOrder: &orderId,
		},
	}
}

// NewMatchOperation returns a new operation for matching maker orders against a matchable order.
func NewMatchOperation(
	takerMatchableOrder MatchableOrder,
	makerFills []MakerFill,
) Operation {
	if takerMatchableOrder.IsLiquidation() {
		return Operation{
			Operation: &Operation_Match{
				Match: &ClobMatch{
					Match: &ClobMatch_MatchPerpetualLiquidation{
						MatchPerpetualLiquidation: &MatchPerpetualLiquidation{
							Liquidated:  takerMatchableOrder.GetSubaccountId(),
							ClobPairId:  takerMatchableOrder.GetClobPairId().ToUint32(),
							PerpetualId: takerMatchableOrder.MustGetLiquidatedPerpetualId(),
							TotalSize:   takerMatchableOrder.GetBaseQuantums().ToUint64(),
							IsBuy:       takerMatchableOrder.IsBuy(),
							Fills:       makerFills,
						},
					},
				},
			},
		}
	} else {
		order := takerMatchableOrder.MustGetOrder()
		return Operation{
			Operation: &Operation_Match{
				Match: &ClobMatch{
					Match: &ClobMatch_MatchOrders{
						MatchOrders: &MatchOrders{
							TakerOrderId:   order.GetOrderId(),
							TakerOrderHash: order.GetOrderHash().ToBytes(),
							Fills:          makerFills,
						},
					},
				},
			},
		}
	}
}

// NewMatchOperationFromPerpetualDeleveragingLiquidation returns a new match operation
// wrapping the `perpDeleveraging` object.
func NewMatchOperationFromPerpetualDeleveragingLiquidation(perpDeleveraging MatchPerpetualDeleveraging) Operation {
	return Operation{
		Operation: &Operation_Match{
			Match: &ClobMatch{
				Match: &ClobMatch_MatchPerpetualDeleveraging{
					MatchPerpetualDeleveraging: &perpDeleveraging,
				},
			},
		},
	}
}

// NewMatchOperationFromPerpetualLiquidation returns a new match operation
// wrapping the `perpLiquidation` object.
func NewMatchOperationFromPerpetualLiquidation(perpLiquidation MatchPerpetualLiquidation) Operation {
	return Operation{
		Operation: &Operation_Match{
			Match: &ClobMatch{
				Match: &ClobMatch_MatchPerpetualLiquidation{
					MatchPerpetualLiquidation: &perpLiquidation,
				},
			},
		},
	}
}

// NewOrderCancellationOperation returns a new operation for canceling an order.
func NewOrderCancellationOperation(msgCancelOrder *MsgCancelOrder) Operation {
	return Operation{
		Operation: &Operation_OrderCancellation{
			OrderCancellation: msgCancelOrder,
		},
	}
}

// GetOperationTextString returns the text string representation of this operation.
// TODO(DEC-1772): Add method for encoding operation protos as JSON to make debugging easier.
func (o *Operation) GetOperationTextString() string {
	return proto.MarshalTextString(o)
}

// GetOperationsQueueString returns a string representation of the provided operations.
func GetOperationsQueueTextString(operations []Operation) string {
	var result string
	for _, operation := range operations {
		result += operation.GetOperationTextString() + "\n"
	}

	return result
}
