package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGetPreviousBlockStatefulOrderCancellations(t *testing.T) {
	tests := map[string]struct {
		operations []types.Operation

		expectedOrderCancellations []types.OrderId
	}{
		"empty operations": {
			operations:                 []types.Operation{},
			expectedOrderCancellations: []types.OrderId{},
		},
		"one stateful order placement": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
				),
			},
			expectedOrderCancellations: []types.OrderId{},
		},
		"one stateful order cancellation": {
			operations: []types.Operation{
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5,
				),
			},
			expectedOrderCancellations: []types.OrderId{
				constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5.GetOrderId(),
			},
		},
		"one short term place order": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				),
			},
			expectedOrderCancellations: []types.OrderId{},
		},
		"one stateful order cancellation of a partially filled stateful order, not replaced": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num1_Id1_Clob0_GTBT_20,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
			},
			expectedOrderCancellations: []types.OrderId{
				constants.CancelLongTermOrder_Alice_Num1_Id1_Clob0_GTBT_20.GetOrderId(),
			},
		},
		"one stateful order cancellation of a partially filled stateful order, re-placed": {
			operations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num1_Id1_Clob0_GTBT_20,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
			},
			expectedOrderCancellations: []types.OrderId{},
		},
		"multiple stateful order cancellations, some being replaced, some cancellations": {
			operations: []types.Operation{
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num1_Id1_Clob0_GTBT_20,
				),
				types.NewOrderCancellationOperation(
					&constants.CancelConditionalOrder_Alice_Num1_Id0_Clob0_GTBT15,
				),
				// Short term order.
				types.NewOrderCancellationOperation(
					&constants.CancelOrder_User1_Num0_Id12_Clob0_GTB5,
				),
				// Valid cancellation
				types.NewOrderCancellationOperation(
					&constants.CancelConditionalOrder_Alice_Num1_Id0_Clob1_GTBT15,
				),
				types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id9_Clob1_Buy15_Price45_GTB19.OrderId,
							FillAmount:   10,
						},
					},
				),
				types.NewOrderPlacementOperation(
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15,
				),
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
				),
				// Valid cancellation
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Bob_Num0_Id0_Clob0_GTBT5,
				),
				types.NewOrderCancellationOperation(
					&constants.CancelOrder_Bob_Num1_Id11_Clob1_GTB20,
				),
			},
			expectedOrderCancellations: []types.OrderId{
				constants.CancelLongTermOrder_Bob_Num0_Id0_Clob0_GTBT5.GetOrderId(),
				constants.CancelConditionalOrder_Alice_Num1_Id0_Clob1_GTBT15.GetOrderId(),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			events := types.ProcessProposerMatchesEvents{
				OperationsProposedInLastBlock: tc.operations,
			}
			orderCancellations := events.GetPreviousBlockStatefulOrderCancellations()
			require.Equal(t, tc.expectedOrderCancellations, orderCancellations)
		})
	}
}
