package types

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMsgPlaceOrder_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg MsgPlaceOrder
		err error
	}{
		"invalid subaccountId owner": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  "invalid_owner",
							Number: uint32(0),
						},
					},
				},
			},
			err: satypes.ErrInvalidSubaccountIdOwner,
		},
		"invalid subaccountId number": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(999_999),
						},
					},
				},
			},
			err: satypes.ErrInvalidSubaccountIdNumber,
		},
		"invalid side": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					Side: Order_Side(uint32(999)),
				},
			},
			err: ErrInvalidOrderSide,
		},
		"invalid time in force": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					Side:        Order_SIDE_BUY,
					TimeInForce: Order_TimeInForce(uint32(999)),
				},
			},
			err: ErrInvalidTimeInForce,
		},
		"invalid condition type": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					Side:          Order_SIDE_BUY,
					ConditionType: Order_ConditionType(uint32(999)),
				},
			},
			err: ErrInvalidConditionType,
		},
		"unspecified side": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
				},
			},
			err: ErrInvalidOrderSide,
		},
		"zero quantums": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					Side: Order_SIDE_BUY,
				},
			},
			err: ErrInvalidOrderQuantums,
		},
		"zero GoodTilBlock": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
				},
			},
			err: ErrInvalidOrderGoodTilBlock,
		},
		"order flag: invalid order flag": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: uint32(999),
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
				},
			},
			err: ErrInvalidOrderFlag,
		},
		"long-term: zero GoodTilBlockTime": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
				},
			},
			err: ErrInvalidStatefulOrderGoodTilBlockTime,
		},
		"conditional: zero GoodTilBlockTime": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Conditional,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
				},
			},
			err: ErrInvalidStatefulOrderGoodTilBlockTime,
		},
		"long-term: cannot be IOC": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					TimeInForce: Order_TIME_IN_FORCE_IOC,
				},
			},
			err: ErrLongTermOrdersCannotRequireImmediateExecution,
		},
		"zero subticks": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					Side:         Order_SIDE_BUY,
					Quantums:     uint64(42),
					GoodTilOneof: &Order_GoodTilBlock{GoodTilBlock: uint32(100)},
				},
			},
			err: ErrInvalidOrderSubticks,
		},
		"success": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					Side:         Order_SIDE_BUY,
					Quantums:     uint64(42),
					GoodTilOneof: &Order_GoodTilBlock{GoodTilBlock: uint32(100)},
					Subticks:     uint64(10),
				},
			},
		},
		"fill-or-kill orders are deprecated": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
					},
					Side:         Order_SIDE_BUY,
					Quantums:     uint64(42),
					GoodTilOneof: &Order_GoodTilBlock{GoodTilBlock: uint32(100)},
					Subticks:     uint64(10),
					TimeInForce:  Order_TIME_IN_FORCE_FILL_OR_KILL,
					ReduceOnly:   false,
				},
			},
			err: ErrDeprecatedField,
		},
		"short-term IOC reduce-only success": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_ShortTerm,
					},
					Side:         Order_SIDE_BUY,
					Quantums:     uint64(42),
					GoodTilOneof: &Order_GoodTilBlock{GoodTilBlock: uint32(100)},
					Subticks:     uint64(10),
					TimeInForce:  Order_TIME_IN_FORCE_IOC,
					ReduceOnly:   true,
				},
			},
		},
		"long-term order reduce-only disabled": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:         Order_SIDE_BUY,
					Quantums:     uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{GoodTilBlockTime: uint32(100)},
					Subticks:     uint64(10),
					TimeInForce:  Order_TIME_IN_FORCE_UNSPECIFIED,
					ReduceOnly:   true,
				},
			},
			err: ErrReduceOnlyDisabled,
		},
		"long-term reduce-only disabled": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Conditional,
					},
					Side:                            Order_SIDE_BUY,
					Quantums:                        uint64(42),
					GoodTilOneof:                    &Order_GoodTilBlockTime{GoodTilBlockTime: uint32(100)},
					Subticks:                        uint64(10),
					TimeInForce:                     Order_TIME_IN_FORCE_UNSPECIFIED,
					ConditionType:                   Order_CONDITION_TYPE_STOP_LOSS,
					ConditionalOrderTriggerSubticks: uint64(8),
					ReduceOnly:                      true,
				},
			},
			err: ErrReduceOnlyDisabled,
		},
		"conditional: valid": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Conditional,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					ConditionType:                   Order_CONDITION_TYPE_TAKE_PROFIT,
					ConditionalOrderTriggerSubticks: uint64(10),
				},
			},
		},
		"conditional: unspecified condition type": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Conditional,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
				},
			},
			err: ErrInvalidConditionType,
		},
		"conditional: zero ConditionalOrderTriggerSubticks": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Conditional,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					ConditionType: Order_CONDITION_TYPE_TAKE_PROFIT,
				},
			},
			err: ErrInvalidConditionalOrderTriggerSubticks,
		},
		"non-conditional: specified condition type": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					ConditionType: Order_CONDITION_TYPE_TAKE_PROFIT,
				},
			},
			err: ErrInvalidConditionType,
		},
		"non-conditional: greater than zero ConditionalOrderTriggerSubticks": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					ConditionalOrderTriggerSubticks: uint64(10),
				},
			},
			err: ErrInvalidConditionalOrderTriggerSubticks,
		},
		"twap: missing config": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Twap,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
				},
			},
			err: ErrInvalidPlaceOrder,
		},
		"twap: interval too small": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Twap,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					TwapParameters: &TwapParameters{
						Interval: MinTwapOrderInterval - 1,
						Duration: 300,
					},
				},
			},
			err: ErrInvalidPlaceOrder,
		},
		"twap: interval too large": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Twap,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					TwapParameters: &TwapParameters{
						Interval: MaxTwapOrderInterval + 1,
						Duration: 300,
					},
				},
			},
			err: ErrInvalidPlaceOrder,
		},
		"twap: duration too small": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Twap,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					TwapParameters: &TwapParameters{
						Interval:       60,
						Duration:       MinTwapOrderDuration - 1,
						PriceTolerance: 1000,
					},
				},
			},
			err: ErrInvalidPlaceOrder,
		},
		"twap: duration too large": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Twap,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					TwapParameters: &TwapParameters{
						Interval:       60,
						Duration:       MaxTwapOrderDuration + 1,
						PriceTolerance: 1000,
					},
				},
			},
			err: ErrInvalidPlaceOrder,
		},
		"twap: duration not multiple of interval": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Twap,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					TwapParameters: &TwapParameters{
						Interval:       60,
						Duration:       301, // Not divisible by 60
						PriceTolerance: 1000,
					},
				},
			},
			err: ErrInvalidPlaceOrder,
		},
		"twap: price tolerance too high": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Twap,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					TwapParameters: &TwapParameters{
						Interval:       60,
						Duration:       300,
						PriceTolerance: MaxTwapOrderPriceTolerance,
					},
				},
			},
			err: ErrInvalidPlaceOrder,
		},
		"twap: valid config": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_Twap,
					},
					Side:     Order_SIDE_BUY,
					Subticks: uint64(10),
					Quantums: uint64(42),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					TwapParameters: &TwapParameters{
						Interval:       60,
						Duration:       300,
						PriceTolerance: 1000,
					},
				},
			},
		},
		"invalid builder address": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
					Subticks: uint64(10),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					BuilderCodeParameters: &BuilderCodeParameters{
						BuilderAddress: "invalid_builder_address",
						FeePpm:         1000,
					},
				},
			},
			err: ErrInvalidBuilderCode,
		},
		"builder code parameters: zero fee ppm": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
					Subticks: uint64(10),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					BuilderCodeParameters: &BuilderCodeParameters{
						BuilderAddress: sample.AccAddress(),
						FeePpm:         0,
					},
				},
			},
			err: ErrInvalidBuilderCode,
		},
		"builder code parameters: above 10_000 fee ppm": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
					Subticks: uint64(10),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					BuilderCodeParameters: &BuilderCodeParameters{
						BuilderAddress: sample.AccAddress(),
						FeePpm:         MaxBuilderCodeFeePpm + 1,
					},
				},
			},
			err: ErrInvalidBuilderCode,
		},
		"valid builder code parameters": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
					Subticks: uint64(10),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					BuilderCodeParameters: &BuilderCodeParameters{
						BuilderAddress: sample.AccAddress(),
						FeePpm:         MaxBuilderCodeFeePpm,
					},
				},
			},
		},
		"invalid order router address": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
					Subticks: uint64(10),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					OrderRouterAddress: "invalid_order_router_address",
				},
			},
			err: ErrInvalidOrderRouterAddress,
		},
		"valid order router address": {
			msg: MsgPlaceOrder{
				Order: Order{
					OrderId: OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner:  sample.AccAddress(),
							Number: uint32(0),
						},
						OrderFlags: OrderIdFlags_LongTerm,
					},
					Side:     Order_SIDE_BUY,
					Quantums: uint64(42),
					Subticks: uint64(10),
					GoodTilOneof: &Order_GoodTilBlockTime{
						GoodTilBlockTime: uint32(100),
					},
					OrderRouterAddress: sample.AccAddress(),
				},
			},
			err: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
