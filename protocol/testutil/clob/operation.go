package clob

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// buildTxAndGetBytesFromMsg builds a transaction with a single message
// and returns that transaction's serialized bytes. It should be noted
// that this function does NOT sign the transaction. This is important
// because certain tests (end to end) require the transaction to be signed
// as transactions pass through CheckTx. Therefore this function should
// only be used for unit tests.
func buildTxAndGetBytesFromMsg(msg sdk.Msg) []byte {
	txConfig := constants.TestEncodingCfg.TxConfig
	tx := txConfig.NewTxBuilder()
	_ = tx.SetMsgs(msg)
	bytes, _ := txConfig.TxEncoder()(tx.GetTx())
	return bytes
}

// NewOrderPlacementOperation returns a new operation for placing an order.
func NewOrderPlacementOperation(order types.Order) types.Operation {
	// TODO Rename function and add assert short term order once tx validation is completed.
	return types.Operation{
		Operation: &types.Operation_ShortTermOrderPlacement{
			ShortTermOrderPlacement: types.NewMsgPlaceOrder(order.MustGetOrder()),
		},
	}
}

// NewPreexistingStatefulOrderPlacementOperation returns a new operation for placing a
// pre-existing stateful order.
// Note this function panics if called with a non-stateful order.
func NewPreexistingStatefulOrderPlacementOperation(order types.Order) types.Operation {
	order.MustBeStatefulOrder()

	orderId := order.GetOrderId()
	return types.Operation{
		Operation: &types.Operation_PreexistingStatefulOrder{
			PreexistingStatefulOrder: &orderId,
		},
	}
}

// NewMatchOperation returns a new operation for matching maker orders against a matchable order.
func NewMatchOperation(
	takerMatchableOrder types.MatchableOrder,
	makerFills []types.MakerFill,
) types.Operation {
	if takerMatchableOrder.IsLiquidation() {
		return types.Operation{
			Operation: &types.Operation_Match{
				Match: &types.ClobMatch{
					Match: &types.ClobMatch_MatchPerpetualLiquidation{
						MatchPerpetualLiquidation: &types.MatchPerpetualLiquidation{
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
		return types.Operation{
			Operation: &types.Operation_Match{
				Match: &types.ClobMatch{
					Match: &types.ClobMatch_MatchOrders{
						MatchOrders: &types.MatchOrders{
							TakerOrderId: order.GetOrderId(),
							Fills:        makerFills,
						},
					},
				},
			},
		}
	}
}

// NewMatchOperationFromPerpetualDeleveragingLiquidation returns a new match operation
// wrapping the `perpDeleveraging` object.
func NewMatchOperationFromPerpetualDeleveragingLiquidation(
	perpDeleveraging types.MatchPerpetualDeleveraging,
) types.Operation {
	return types.Operation{
		Operation: &types.Operation_Match{
			Match: &types.ClobMatch{
				Match: &types.ClobMatch_MatchPerpetualDeleveraging{
					MatchPerpetualDeleveraging: &perpDeleveraging,
				},
			},
		},
	}
}

// NewMatchOperationFromPerpetualLiquidation returns a new match operation
// wrapping the `perpLiquidation` object.
func NewMatchOperationFromPerpetualLiquidation(perpLiquidation types.MatchPerpetualLiquidation) types.Operation {
	return types.Operation{
		Operation: &types.Operation_Match{
			Match: &types.ClobMatch{
				Match: &types.ClobMatch_MatchPerpetualLiquidation{
					MatchPerpetualLiquidation: &perpLiquidation,
				},
			},
		},
	}
}

// NewDeleveragingMatchOperation returns a new match operation for deleveraging
// against a undercollateralized subaccount that has failed liquidation.
func NewDeleveragingMatchOperation(
	liquidatedSubaccountId satypes.SubaccountId,
	perpetualId uint32,
	fills []types.MatchPerpetualDeleveraging_Fill,
) types.Operation {
	return types.Operation{
		Operation: &types.Operation_Match{
			Match: &types.ClobMatch{
				Match: &types.ClobMatch_MatchPerpetualDeleveraging{
					MatchPerpetualDeleveraging: &types.MatchPerpetualDeleveraging{
						Liquidated:  liquidatedSubaccountId,
						PerpetualId: perpetualId,
						Fills:       fills,
					},
				},
			},
		},
	}
}

// NewOrderCancellationOperation returns a new operation for canceling an order.
func NewOrderCancellationOperation(msgCancelOrder *types.MsgCancelOrder) types.Operation {
	// TODO Rename function and add assert short term order once tx validation is completed.
	return types.Operation{
		Operation: &types.Operation_ShortTermOrderCancellation{
			ShortTermOrderCancellation: msgCancelOrder,
		},
	}
}

// NewShortTermOrderPlacementOperationRaw returns a new raw operation for placing an order.
func NewShortTermOrderPlacementOperationRaw(order types.Order) types.OperationRaw {
	// Create new tx that wraps the msg
	msg := types.NewMsgPlaceOrder(order.MustGetOrder())
	bytes := buildTxAndGetBytesFromMsg(msg)

	return types.OperationRaw{
		Operation: &types.OperationRaw_ShortTermOrderPlacement{
			ShortTermOrderPlacement: bytes,
		},
	}
}

// NewMatchOperationRaw returns a new raw operation for matching maker orders against a matchable order.
func NewMatchOperationRaw(
	takerMatchableOrder types.MatchableOrder,
	makerFills []types.MakerFill,
) types.OperationRaw {
	if takerMatchableOrder.IsLiquidation() {
		return types.OperationRaw{
			Operation: &types.OperationRaw_Match{
				Match: &types.ClobMatch{
					Match: &types.ClobMatch_MatchPerpetualLiquidation{
						MatchPerpetualLiquidation: &types.MatchPerpetualLiquidation{
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
		return types.OperationRaw{
			Operation: &types.OperationRaw_Match{
				Match: &types.ClobMatch{
					Match: &types.ClobMatch_MatchOrders{
						MatchOrders: &types.MatchOrders{
							TakerOrderId: order.GetOrderId(),
							Fills:        makerFills,
						},
					},
				},
			},
		}
	}
}

// NewMatchOperationRawFromPerpetualLiquidation returns a new raw match operation
// wrapping the `perpLiquidation` object.
func NewMatchOperationRawFromPerpetualLiquidation(
	perpLiquidation types.MatchPerpetualLiquidation,
) types.OperationRaw {
	return types.OperationRaw{
		Operation: &types.OperationRaw_Match{
			Match: &types.ClobMatch{
				Match: &types.ClobMatch_MatchPerpetualLiquidation{
					MatchPerpetualLiquidation: &perpLiquidation,
				},
			},
		},
	}
}

// NewMatchOperationRawFromPerpetualDeleveragingLiquidation returns a new raw match operation
// wrapping the `perpDeleveraging` object.
func NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
	perpDeleveraging types.MatchPerpetualDeleveraging,
) types.OperationRaw {
	return types.OperationRaw{
		Operation: &types.OperationRaw_Match{
			Match: &types.ClobMatch{
				Match: &types.ClobMatch_MatchPerpetualDeleveraging{
					MatchPerpetualDeleveraging: &perpDeleveraging,
				},
			},
		},
	}
}

// NewOrderRemovalOperationRaw returns a new raw order removal operation.
func NewOrderRemovalOperationRaw(
	orderId types.OrderId,
	reason types.OrderRemoval_RemovalReason,
) types.OperationRaw {
	return types.OperationRaw{
		Operation: &types.OperationRaw_OrderRemoval{
			OrderRemoval: &types.OrderRemoval{
				OrderId:       orderId,
				RemovalReason: reason,
			},
		},
	}
}
