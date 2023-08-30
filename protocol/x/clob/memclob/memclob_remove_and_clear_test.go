package memclob

import (
	"testing"

	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	memclobtestutil "github.com/dydxprotocol/v4-chain/protocol/testutil/memclob"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestRemoveAndClearOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		placedOperations          []types.Operation
		preexistingStatefulOrders []types.LongTermOrderPlacement

		// Expectations.
		expectedRemainingBids []OrderWithRemainingSize
		expectedRemainingAsks []OrderWithRemainingSize
	}{
		"no orders were placed": {
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedRemainingBids: []OrderWithRemainingSize{},
		},
		"one short-term order is placed": {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					RemainingSize: 30,
				},
			},
		},
		"a pre-existing stateful order is placed": {
			placedOperations: []types.Operation{
				clobtest.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
			},
			preexistingStatefulOrders: []types.LongTermOrderPlacement{
				{
					Order: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					RemainingSize: 25,
				},
			},
		},
		"two short-term orders are placed and cross": {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30),
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedRemainingBids: []OrderWithRemainingSize{},
		},
		"two short-term orders are placed and cross, the partially filled order is canceled": {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30.OrderId,
						31,
					),
				),
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedRemainingBids: []OrderWithRemainingSize{},
		},
		"two short-term orders are placed and cross, a new short-term is placed and does not cross": {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20),
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20,
					RemainingSize: 25,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
		},
		"two short-term orders are placed and cross, a new stateful order is placed and does not cross": {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30),
				clobtest.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10),
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					RemainingSize: 25,
				},
			},
		},
		"two short-term orders are placed and cross, the previous taker order is replaced and crosses": {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell35_Price15_GTB25),
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedRemainingBids: []OrderWithRemainingSize{},
		},
		"two short-term orders are placed and cross, the previous taker order is replaced and does not cross": {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB34),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell35_Price15_GTB25),
			},
			expectedRemainingAsks: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id10_Clob0_Sell35_Price15_GTB25,
					RemainingSize: 25,
				},
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
		},
		"a pre-existing stateful order crosses a short-term order": {
			placedOperations: []types.Operation{
				clobtest.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20),
			},
			preexistingStatefulOrders: []types.LongTermOrderPlacement{
				{
					Order: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
			expectedRemainingBids: []OrderWithRemainingSize{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			ctx, _, _ := sdktest.NewSdkContextWithMultistore()
			ctx = ctx.WithIsCheckTx(true)

			// Setup memclob state and test expectations.
			memclob, _ := memclobOperationsTestSetupWithCustomCollatCheck(
				t,
				ctx,
				tc.placedOperations,
				memclobtestutil.AlwaysSuccessfulCollatCheckFn,
				constants.GetStatePosition_ZeroPositionSize,
				tc.preexistingStatefulOrders,
			)

			operations, _ := memclob.operationsToPropose.GetOperationsToReplay()
			memclob.RemoveAndClearOperationsQueue(
				ctx,
				operations,
			)

			AssertMemclobHasOrders(
				t,
				ctx,
				memclob,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
			)

			operations, shortTermTxBytes := memclob.operationsToPropose.GetOperationsToReplay()
			require.Empty(t, operations)
			require.Empty(t, shortTermTxBytes)
		})
	}
}
