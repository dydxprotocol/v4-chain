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
							Number: uint32(9999),
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
		"long-term: cannot be FOK": {
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
					TimeInForce: Order_TIME_IN_FORCE_FILL_OR_KILL,
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
		"success with fill-or-kill order": {
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
		},
		"reduce-only disabled": {
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
					ReduceOnly:   true,
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
