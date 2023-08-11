package types_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestAssignNonceToOrder(t *testing.T) {
	// OrderToAssignNonce is a helper struct for performing test assertions on the `OperationHashToNonce`
	// data structure.
	type OrderToAssignNonce struct {
		order                      types.Order
		isPreexistingStatefulOrder bool
	}

	tests := map[string]struct {
		// State.
		ordersToAssignNonces []OrderToAssignNonce

		// Expectations.
		expectedOrderNonces []struct {
			orderToAssignNonce OrderToAssignNonce
			nonce              types.Nonce
		}
	}{
		"Can assign zero orders a nonce": {
			ordersToAssignNonces: []OrderToAssignNonce{},

			expectedOrderNonces: []struct {
				orderToAssignNonce OrderToAssignNonce
				nonce              types.Nonce
			}{},
		},
		"Can assign a single Short-Term order a nonce": {
			ordersToAssignNonces: []OrderToAssignNonce{
				{
					order:                      constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					isPreexistingStatefulOrder: false,
				},
			},

			expectedOrderNonces: []struct {
				orderToAssignNonce OrderToAssignNonce
				nonce              types.Nonce
			}{
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
						isPreexistingStatefulOrder: false,
					},
					nonce: 0,
				},
			},
		},
		"Can assign a single new stateful order a nonce": {
			ordersToAssignNonces: []OrderToAssignNonce{
				{
					order:                      constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					isPreexistingStatefulOrder: false,
				},
			},

			expectedOrderNonces: []struct {
				orderToAssignNonce OrderToAssignNonce
				nonce              types.Nonce
			}{
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
						isPreexistingStatefulOrder: false,
					},
					nonce: 0,
				},
			},
		},
		"Can assign multiple Short-Term, new stateful, and pre-existing stateful orders a nonce": {
			ordersToAssignNonces: []OrderToAssignNonce{
				{
					order:                      constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					isPreexistingStatefulOrder: false,
				},
				{
					order:                      constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
					isPreexistingStatefulOrder: false,
				},
				{
					order:                      constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
					isPreexistingStatefulOrder: true,
				},
				{
					order:                      constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
					isPreexistingStatefulOrder: false,
				},
				{
					order:                      constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
					isPreexistingStatefulOrder: false,
				},
				{
					order:                      constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
					isPreexistingStatefulOrder: false,
				},
				{
					order:                      constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					isPreexistingStatefulOrder: true,
				},
				{
					order:                      constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
					isPreexistingStatefulOrder: false,
				},
				{
					order:                      constants.LongTermOrder_Alice_Num1_Id2_Clob0_Buy10_Price40_GTBT10,
					isPreexistingStatefulOrder: false,
				},
			},

			expectedOrderNonces: []struct {
				orderToAssignNonce OrderToAssignNonce
				nonce              types.Nonce
			}{
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
						isPreexistingStatefulOrder: false,
					},
					nonce: 0,
				},
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.Order_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
						isPreexistingStatefulOrder: false,
					},
					nonce: 1,
				},
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
						isPreexistingStatefulOrder: true,
					},
					nonce: 2,
				},
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
						isPreexistingStatefulOrder: false,
					},
					nonce: 3,
				},
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.Order_Alice_Num1_Id12_Clob0_Sell20_Price5_GTB25,
						isPreexistingStatefulOrder: false,
					},
					nonce: 4,
				},
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.Order_Alice_Num0_Id3_Clob1_Sell5_Price10_GTB15,
						isPreexistingStatefulOrder: false,
					},
					nonce: 5,
				},
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
						isPreexistingStatefulOrder: true,
					},
					nonce: 6,
				},
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
						isPreexistingStatefulOrder: false,
					},
					nonce: 7,
				},
				{
					orderToAssignNonce: OrderToAssignNonce{
						order:                      constants.LongTermOrder_Alice_Num1_Id2_Clob0_Buy10_Price40_GTBT10,
						isPreexistingStatefulOrder: false,
					},
					nonce: 8,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			for _, order := range tc.ordersToAssignNonces {
				otp.AssignNonceToOrder(order.order, order.isPreexistingStatefulOrder)
			}

			// Verify expectations.
			for _, expectedOrderNonce := range tc.expectedOrderNonces {
				var operation types.Operation
				order := expectedOrderNonce.orderToAssignNonce.order
				if expectedOrderNonce.orderToAssignNonce.isPreexistingStatefulOrder {
					operation = types.NewPreexistingStatefulOrderPlacementOperation(order)
				} else {
					operation = types.NewOrderPlacementOperation(order)
				}
				operationHash := operation.GetOperationHash()

				// Verify the order is present and the nonce is correct.
				nonce, ok := otp.OperationHashToNonce[operationHash]
				if !ok {
					t.Errorf(
						"Expected order %+v to have nonce %d, but it was not found. Pre-existing: %t",
						expectedOrderNonce.orderToAssignNonce,
						expectedOrderNonce.nonce,
						expectedOrderNonce.orderToAssignNonce.isPreexistingStatefulOrder,
					)
				}
				if nonce != expectedOrderNonce.nonce {
					t.Errorf(
						"Expected order %v to have nonce %d, but it had nonce %d. Pre-existing: %t",
						expectedOrderNonce.orderToAssignNonce,
						expectedOrderNonce.nonce,
						nonce,
						expectedOrderNonce.orderToAssignNonce.isPreexistingStatefulOrder,
					)
				}
			}

			// Verify the number of expected order nonces is the same as the number of entries
			// in the `OperationHashToNonce` map.
			require.Equal(t, len(tc.expectedOrderNonces), len(otp.OperationHashToNonce))

			// Verify the next available nonce is correct.
			require.Equal(t, types.Nonce(len(tc.expectedOrderNonces)), otp.NextAvailableNonce)
		})
	}
}

func TestAssignNonceToOrder_PanicsOnDuplicate(t *testing.T) {
	otp := types.NewOperationsToPropose()

	// Assign a nonce to the same Short-Term twice.
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	shortTermOrderPlacementOperation := types.NewOrderPlacementOperation(shortTermOrder)
	otp.AssignNonceToOrder(shortTermOrder, false)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"assignNonceToOperation: operation (%+v) has already been assigned nonce 0. "+
				"The current operations queue is: %s",
			shortTermOrderPlacementOperation.GetOperationTextString(),
			types.GetOperationsQueueTextString(otp.GetOperationsQueue()),
		),
		func() {
			otp.AssignNonceToOrder(shortTermOrder, false)
		},
	)

	// Assign a nonce to the same newly-placed stateful order twice.
	statefulOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20
	statefulOrderPlacementOperation := types.NewOrderPlacementOperation(statefulOrder)
	otp.AssignNonceToOrder(statefulOrder, false)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"assignNonceToOperation: operation (%+v) has already been assigned nonce 1. "+
				"The current operations queue is: %s",
			statefulOrderPlacementOperation.GetOperationTextString(),
			types.GetOperationsQueueTextString(otp.GetOperationsQueue()),
		),
		func() {
			otp.AssignNonceToOrder(statefulOrder, false)
		},
	)

	// Assign a nonce to the same pre-existing stateful order twice.
	preexistingStatefulOrder := constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25
	preexistingStatefulOrderPlacementOperation := types.NewPreexistingStatefulOrderPlacementOperation(
		preexistingStatefulOrder,
	)
	otp.AssignNonceToOrder(preexistingStatefulOrder, true)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"assignNonceToOperation: operation (%+v) has already been assigned nonce 2. "+
				"The current operations queue is: %s",
			preexistingStatefulOrderPlacementOperation.GetOperationTextString(),
			types.GetOperationsQueueTextString(otp.GetOperationsQueue()),
		),
		func() {
			otp.AssignNonceToOrder(preexistingStatefulOrder, true)
		},
	)
}

func TestAssignNonceToOrder_PanicsOnPreexistingShortTermOrder(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrder.OrderId,
		),
		func() {
			otp.AssignNonceToOrder(shortTermOrder, true)
		},
	)
}

func TestAddOrderCancellationToOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		cancels []types.MsgCancelOrder

		// Expectations.
		expectedOperationsQueue map[types.Nonce]types.Operation
	}{
		"Starts with an empty operations queue": {
			cancels: []types.MsgCancelOrder{},

			expectedOperationsQueue: map[types.Nonce]types.Operation{},
		},
		"Can add a cancellation to the operations queue": {
			cancels: []types.MsgCancelOrder{
				constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
			},

			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
			},
		},
		"Can assign multiple cancellations to the operations queue": {
			cancels: []types.MsgCancelOrder{
				constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25,
				constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35,
				constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5,
				constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20,
			},

			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
				1: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25,
				),
				2: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35,
				),
				3: types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5,
				),
				4: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20,
				),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			for _, cancel := range tc.cancels {
				otp.AddOrderCancellationToOperationsQueue(cancel)
			}

			// Verify expectations.
			require.Equal(t, tc.expectedOperationsQueue, otp.NonceToOperationToPropose)

			// Verify the next available nonce is correct.
			require.Equal(t, types.Nonce(len(tc.expectedOperationsQueue)), otp.NextAvailableNonce)
		})
	}
}

func TestAddOrderCancellationToOperationsQueue_PanicsOnDuplicateCancel(t *testing.T) {
	otp := types.NewOperationsToPropose()
	cancel := constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20
	cancelOperation := types.NewOrderCancellationOperation(&cancel)
	otp.AddOrderCancellationToOperationsQueue(cancel)
	// TODO(DEC-1772): Add assertion on panic error message.
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"assignNonceToOperation: operation (%+v) has already been assigned nonce 0. "+
				"The current operations queue is: %s",
			cancelOperation.GetOperationTextString(),
			types.GetOperationsQueueTextString(otp.GetOperationsQueue()),
		),
		func() {
			otp.AddOrderCancellationToOperationsQueue(cancel)
		},
	)
}

func TestAddPreexistingStatefulOrderPlacementToOperationsQueue_PanicsOnShortTermOrder(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	shortTermOrderPlacementOperation := types.NewOrderPlacementOperation(shortTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrderPlacementOperation.GetOrderPlacement().GetOrder().OrderId,
		),
		func() {
			otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(shortTermOrder)
		},
	)
}

func TestAddPreexistingStatefulOrderPlacementToOperationsQueue_PanicsIfNonexistentNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5
	longTermOrderPlacementOperation := types.NewPreexistingStatefulOrderPlacementOperation(longTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			longTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(longTermOrder)
		},
	)
}

func TestAddPreexistingStatefulOrderPlacementToOperationsQueue_PanicsOnDuplicatePlacement(t *testing.T) {
	otp := types.NewOperationsToPropose()
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5
	// Assign nonce to order.
	otp.AssignNonceToOrder(longTermOrder, true)
	// First addition to the operations queue.
	otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(longTermOrder)
	longTermOrderPlacementOperation := types.NewPreexistingStatefulOrderPlacementOperation(longTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"insertOperationIntoOperationsToPropose: an operation with nonce %d already exists "+
				"in the operations to propose. New operation: (%+v). Existing operation: (%+v).",
			0,
			longTermOrderPlacementOperation.GetOperationTextString(),
			longTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(longTermOrder)
		},
	)
}

func TestAddOrderPlacementToOperationsQueue_PanicsIfNonexistentNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	shortTermOrderPlacementOperation := types.NewOrderPlacementOperation(shortTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			shortTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.AddOrderPlacementToOperationsQueue(shortTermOrder)
		},
	)
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5
	longTermOrderPlacementOperation := types.NewOrderPlacementOperation(longTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			longTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.AddOrderPlacementToOperationsQueue(longTermOrder)
		},
	)
}

func TestAddOrderPlacementToOperationsQueue_PanicsOnDuplicatePlacement(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	// Assign nonce to order.
	otp.AssignNonceToOrder(shortTermOrder, false)
	// First addition to the operations queue.
	otp.AddOrderPlacementToOperationsQueue(shortTermOrder)
	shortTermOrderPlacementOperation := types.NewOrderPlacementOperation(shortTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"insertOperationIntoOperationsToPropose: an operation with nonce %d already exists "+
				"in the operations to propose. New operation: (%+v). Existing operation: (%+v).",
			0,
			shortTermOrderPlacementOperation.GetOperationTextString(),
			shortTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.AddOrderPlacementToOperationsQueue(shortTermOrder)
		},
	)
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5
	// Assign nonce to order.
	otp.AssignNonceToOrder(longTermOrder, false)
	// First addition to the operations queue.
	otp.AddOrderPlacementToOperationsQueue(longTermOrder)
	longTermOrderPlacementOperation := types.NewOrderPlacementOperation(longTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"insertOperationIntoOperationsToPropose: an operation with nonce %d already exists "+
				"in the operations to propose. New operation: (%+v). Existing operation: (%+v).",
			1,
			longTermOrderPlacementOperation.GetOperationTextString(),
			longTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.AddOrderPlacementToOperationsQueue(longTermOrder)
		},
	)
}

func TestAddOrderPlacementToOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		orders []types.Order

		// Expectations.
		expectedOperationsQueue map[types.Nonce]types.Operation
	}{
		"Starts with an empty operations queue": {
			orders: []types.Order{},

			expectedOperationsQueue: map[types.Nonce]types.Operation{},
		},
		"Can add an order to the operations queue": {
			orders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
			},
		},
		"Can assign multiple orders to the operations queue": {
			orders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
				constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				constants.Order_Bob_Num0_Id9_Clob0_Sell20_Price1000,
			},

			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				1: types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
				),
				2: types.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
				),
				3: types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				),
				4: types.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id9_Clob0_Sell20_Price1000,
				),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			for _, order := range tc.orders {
				// Generate a nonce for the order then add it to operations queue.
				otp.AssignNonceToOrder(order, false)
				otp.AddOrderPlacementToOperationsQueue(order)

				if order.IsStatefulOrder() {
					require.False(t, otp.IsMakerOrderPreexistingStatefulOrder(order))
				}
			}

			// Verify expectations.
			require.Equal(t, tc.expectedOperationsQueue, otp.NonceToOperationToPropose)

			// Verify the next available nonce is correct.
			require.Equal(t, types.Nonce(len(tc.expectedOperationsQueue)), otp.NextAvailableNonce)
		})
	}
}

func TestAddPreexistingStatefulOrderPlacementToOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		orders []types.Order

		// Expectations.
		expectedOperationsQueue map[types.Nonce]types.Operation
	}{
		"Starts with an empty operations queue": {
			orders: []types.Order{},

			expectedOperationsQueue: map[types.Nonce]types.Operation{},
		},
		"Can add an order to the operations queue": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},

			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
			},
		},
		"Can assign multiple orders to the operations queue": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
			},

			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				1: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				),
				2: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				),
				3: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				),
				4: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			for _, order := range tc.orders {
				// Generate a nonce for the order then add it to operations queue.
				otp.AssignNonceToOrder(order, true)
				otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(order)
				require.True(t, otp.IsMakerOrderPreexistingStatefulOrder(order))
			}

			// Verify expectations.
			require.Equal(t, tc.expectedOperationsQueue, otp.NonceToOperationToPropose)

			// Verify the next available nonce is correct.
			require.Equal(t, types.Nonce(len(tc.expectedOperationsQueue)), otp.NextAvailableNonce)
		})
	}
}

func TestIsMakerOrderPreExistingStatefulOrder_PanicsOnShortTermOrder(t *testing.T) {
	otp := types.NewOperationsToPropose()

	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrder.OrderId,
		),
		func() {
			otp.IsMakerOrderPreexistingStatefulOrder(shortTermOrder)
		},
	)
}

func TestIsMakerOrderPreExistingStatefulOrder_PanicsIfNonExistentNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()

	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5
	longTermOrderPlacementOperation := types.NewOrderPlacementOperation(longTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			longTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.AddOrderPlacementToOperationsQueue(longTermOrder)
		},
	)
}

func TestIsMakerOrderPreExistingStatefulOrder(t *testing.T) {
	// Returns true with a pre-existing stateful order placement operation.
	otp := types.NewOperationsToPropose()
	preexistingOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	otp.AssignNonceToOrder(preexistingOrder, true)
	otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(
		preexistingOrder,
	)
	result := otp.IsMakerOrderPreexistingStatefulOrder(preexistingOrder)
	require.True(t, result)

	// Returns false with a new stateful order placement operation.
	newOrder := constants.LongTermOrder_Bob_Num0_Id2_Clob0_Buy15_Price5_GTBT10
	otp.AssignNonceToOrder(newOrder, false)
	otp.AddOrderPlacementToOperationsQueue(
		newOrder,
	)
	result = otp.IsMakerOrderPreexistingStatefulOrder(newOrder)
	require.False(t, result)
}

func TestRemoveOrderPlacementNonce(t *testing.T) {
	tests := map[string]struct {
		// State.
		ordersToAssignNonces []types.Order

		// Parameters.
		ordersToRemoveNonces []types.Order

		// Expectations.
		expectedOrderNonces map[types.Nonce]types.Order
	}{
		"Can assign and remove a nonce from an order": {
			ordersToAssignNonces: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			ordersToRemoveNonces: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			expectedOrderNonces: map[types.Nonce]types.Order{},
		},
		"Can assign a nonce to multiple orders and remove a nonce from only one order": {
			ordersToAssignNonces: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price35_GTB32,
			},

			ordersToRemoveNonces: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			expectedOrderNonces: map[types.Nonce]types.Order{
				1: constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price35_GTB32,
			},
		},
		"Can assign a nonce to multiple orders and remove a nonce from all of them": {
			ordersToAssignNonces: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price35_GTB32,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
			},

			ordersToRemoveNonces: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.Order_Bob_Num0_Id12_Clob0_Sell20_Price35_GTB32,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
			},

			expectedOrderNonces: map[types.Nonce]types.Order{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			for _, order := range tc.ordersToAssignNonces {
				otp.AssignNonceToOrder(order, false)
			}

			for _, order := range tc.ordersToRemoveNonces {
				otp.RemoveOrderPlacementNonce(order)
			}

			// Verify expectations.
			expectedOperationsHashToNonce := make(map[types.OperationHash]types.Nonce)
			for nonce, order := range tc.expectedOrderNonces {
				operation := types.NewOrderPlacementOperation(order)
				expectedOperationsHashToNonce[operation.GetOperationHash()] = nonce
			}
			require.Equal(t, expectedOperationsHashToNonce, otp.OperationHashToNonce)

			// Verify the next available nonce is correct.
			require.Equal(t, types.Nonce(len(tc.ordersToAssignNonces)), otp.NextAvailableNonce)
		})
	}
}

func TestRemovePreexistingStatefulOrderPlacementNonce(t *testing.T) {
	tests := map[string]struct {
		// State.
		ordersToAssignNonces []types.Order

		// Parameters.
		ordersToRemoveNonces []types.Order

		// Expectations.
		expectedOrderNonces map[types.Nonce]types.Order
	}{
		"Can assign and remove a nonce from an order": {
			ordersToAssignNonces: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
			},

			ordersToRemoveNonces: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
			},

			expectedOrderNonces: map[types.Nonce]types.Order{},
		},
		"Can assign a nonce to multiple orders and remove a nonce from only one order": {
			ordersToAssignNonces: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
			},

			ordersToRemoveNonces: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
			},

			expectedOrderNonces: map[types.Nonce]types.Order{
				1: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
			},
		},
		"Can assign a nonce to multiple orders and remove a nonce from all of them": {
			ordersToAssignNonces: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
			},

			ordersToRemoveNonces: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
			},

			expectedOrderNonces: map[types.Nonce]types.Order{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			for _, order := range tc.ordersToAssignNonces {
				otp.AssignNonceToOrder(order, true)
			}

			for _, order := range tc.ordersToRemoveNonces {
				otp.RemovePreexistingStatefulOrderPlacementNonce(order)
			}

			// Verify expectations.
			expectedOperationsHashToNonce := make(map[types.OperationHash]types.Nonce)
			for nonce, order := range tc.expectedOrderNonces {
				operation := types.NewPreexistingStatefulOrderPlacementOperation(order)
				expectedOperationsHashToNonce[operation.GetOperationHash()] = nonce
			}
			require.Equal(t, expectedOperationsHashToNonce, otp.OperationHashToNonce)

			// Verify the next available nonce is correct.
			require.Equal(t, types.Nonce(len(tc.ordersToAssignNonces)), otp.NextAvailableNonce)
		})
	}
}

func TestRemoveOrderPlacementNonce_PanicsOnNoNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	shortTermOrderPlacementOperation := types.NewOrderPlacementOperation(shortTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustRemoveOperationNonce: nonce for operation (%+v) does not exist",
			shortTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.RemoveOrderPlacementNonce(shortTermOrder)
		},
	)
}

func TestRemovePreexistingStatefulOrderPlacementNonce_PanicsOnNoNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20
	longTermOrderPlacementOperation := types.NewPreexistingStatefulOrderPlacementOperation(longTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustRemoveOperationNonce: nonce for operation (%+v) does not exist",
			longTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.RemovePreexistingStatefulOrderPlacementNonce(longTermOrder)
		},
	)
}

func TestRemoveOrderPlacementNonce_PanicsIfInOperationsToPropose(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	otp.AssignNonceToOrder(shortTermOrder, false)
	otp.AddOrderPlacementToOperationsQueue(shortTermOrder)
	shortTermOrderPlacementOperation := types.NewOrderPlacementOperation(shortTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustRemoveOperationNonce: operation (%+v) has nonce 0 in operations to propose",
			shortTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.RemoveOrderPlacementNonce(shortTermOrder)
		},
	)
}

func TestRemovePreexistingStatefulOrderPlacementNonce_PanicsIfInOperationsToPropose(t *testing.T) {
	otp := types.NewOperationsToPropose()
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25
	otp.AssignNonceToOrder(longTermOrder, true)
	otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(longTermOrder)
	longTermOrderPlacementOperation := types.NewPreexistingStatefulOrderPlacementOperation(longTermOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustRemoveOperationNonce: operation (%+v) has nonce 0 in operations to propose",
			longTermOrderPlacementOperation.GetOperationTextString(),
		),
		func() {
			otp.RemovePreexistingStatefulOrderPlacementNonce(longTermOrder)
		},
	)
}

func TestRemovePreexistingStatefulOrderPlacementNonce_PanicsWithShortTermOrder(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrder.OrderId,
		),
		func() {
			otp.RemovePreexistingStatefulOrderPlacementNonce(shortTermOrder)
		},
	)
}

func TestIsOrderPlacementInOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		orders []types.Order
	}{
		"Orders are not in an empty operations queue": {
			orders: []types.Order{},
		},
		"Can add an order to the operations queue and it's present": {
			orders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
		},
		"Can assign multiple orders to the operations queue and they're all present": {
			orders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
				constants.Order_Bob_Num0_Id3_Clob0_Sell20_Price10_GTB20_RO,
				constants.Order_Carl_Num0_Id4_Clob1_Buy01ETH_Price3000,
				constants.Order_Bob_Num0_Id9_Clob0_Sell20_Price1000,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			for _, order := range tc.orders {
				// Generate a nonce for the order then add it to operations queue.
				otp.AssignNonceToOrder(order, false)
				otp.AddOrderPlacementToOperationsQueue(order)
			}

			ordersNotInOperationsQueue := []types.Order{
				constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
				constants.Order_Alice_Num0_Id4_Clob2_Buy25_Price5_GTB20,
			}
			for _, order := range ordersNotInOperationsQueue {
				otp.AssignNonceToOrder(order, false)
			}

			// Verify orders added to the operations queue are present.
			for _, order := range tc.orders {
				require.True(t, otp.IsOrderPlacementInOperationsQueue(order))
			}

			// Verify orders not added to the operations queue are not present.
			for _, order := range ordersNotInOperationsQueue {
				require.False(t, otp.IsOrderPlacementInOperationsQueue(order))
			}
		})
	}
}

func TestIsPreexistingStatefulOrderPlacementInOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		orders []types.Order
	}{
		"Orders are not in an empty operations queue": {
			orders: []types.Order{},
		},
		"Can add an order to the operations queue and it's present": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
			},
		},
		"Can assign multiple orders to the operations queue and they're all present": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				constants.LongTermOrder_Alice_Num1_Id2_Clob0_Buy10_Price40_GTBT10,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			for _, order := range tc.orders {
				// Generate a nonce for the order then add it to operations queue.
				otp.AssignNonceToOrder(order, true)
				otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(order)
				require.True(t, otp.IsMakerOrderPreexistingStatefulOrder(order))
			}

			ordersNotInOperationsQueue := []types.Order{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10,
			}
			for _, order := range ordersNotInOperationsQueue {
				otp.AssignNonceToOrder(order, true)
			}

			// Verify orders added to the operations queue are present.
			for _, order := range tc.orders {
				require.True(t, otp.IsPreexistingStatefulOrderInOperationsQueue(order))
			}

			// Verify orders not added to the operations queue are not present.
			for _, order := range ordersNotInOperationsQueue {
				require.False(t, otp.IsPreexistingStatefulOrderInOperationsQueue(order))
			}
		})
	}
}

func TestIsOrderPlacementInOperationsQueue_PanicsOnNoNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30
	operation := types.NewOrderPlacementOperation(shortTermOrder)

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"isOperationInOperationsQueue: operation (%+v) has no nonce",
			operation.GetOperationTextString(),
		),
		func() {
			otp.IsOrderPlacementInOperationsQueue(shortTermOrder)
		},
	)
}

func TestIsPreexistingStatefulOrderPlacementInOperationsQueue_PanicsOnNoNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	longTermOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	operation := types.NewPreexistingStatefulOrderPlacementOperation(longTermOrder)

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"isOperationInOperationsQueue: operation (%+v) has no nonce",
			operation.GetOperationTextString(),
		),
		func() {
			otp.IsPreexistingStatefulOrderInOperationsQueue(longTermOrder)
		},
	)
}

func TestIsPreexistingStatefulOrderPlacementInOperationsQueue_PanicsOnShortTermOrder(t *testing.T) {
	otp := types.NewOperationsToPropose()
	shortTermOrder := constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrder.OrderId,
		),
		func() {
			otp.IsPreexistingStatefulOrderInOperationsQueue(shortTermOrder)
		},
	)
}

func TestGetOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		nonceToOperationToPropose map[types.Nonce]types.Operation

		// Expectations.
		expectedOperations []types.Operation
	}{
		"Empty nonce map": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{},

			expectedOperations: []types.Operation{},
		},
		"Nonce map single entry, zero index": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{
				0: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
			},
			expectedOperations: []types.Operation{
				types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
			},
		},
		"Nonce map single entry, nonzero index": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{
				5: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
			},
			expectedOperations: []types.Operation{
				types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
			},
		},
		"Nonce map multiple entry, no nonces missing": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{
				1: types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				0: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
				2: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				),
			},
			expectedOperations: []types.Operation{
				types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				),
			},
		},
		"Nonce map multiple entry, missing nonces": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{
				7130: types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				1998: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				),
				8675309: types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
				408: types.NewMatchOperation(
					&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					[]types.MakerFill{},
				),
			},
			expectedOperations: []types.Operation{
				types.NewMatchOperation(
					&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					[]types.MakerFill{},
				),
				types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				types.NewOrderCancellationOperation(
					&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
				),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			otp.NonceToOperationToPropose = tc.nonceToOperationToPropose

			operations := otp.GetOperationsQueue()
			// Verify expectations.
			require.Equal(t, tc.expectedOperations, operations)
		})
	}
}

func TestClearOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		nonceToOperationToPropose map[types.Nonce]types.Operation
		operationHashToNonce      map[types.OperationHash]types.Nonce
		nextAvailableNonce        types.Nonce

		// Expectations.
		expectedOperationHashToNonce map[types.OperationHash]types.Nonce
		expectedNextAvailableNonce   types.Nonce
	}{
		"Empty nonce map with no preassigned nonces": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{},
			operationHashToNonce:      map[types.OperationHash]types.Nonce{},
			nextAvailableNonce:        2,

			expectedOperationHashToNonce: map[types.OperationHash]types.Nonce{},
			expectedNextAvailableNonce:   2,
		},
		"Empty nonce map with preassigned nonces": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{},
			operationHashToNonce: convertListOfOperationsToOperationHashToNonce(
				map[types.Operation]types.Nonce{
					types.NewPreexistingStatefulOrderPlacementOperation(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					): types.Nonce(0),
				},
			),
			nextAvailableNonce: 2,

			expectedOperationHashToNonce: convertListOfOperationsToOperationHashToNonce(
				map[types.Operation]types.Nonce{
					types.NewPreexistingStatefulOrderPlacementOperation(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					): types.Nonce(0),
				},
			),
			expectedNextAvailableNonce: 2,
		},
		"nonce map clear results in empty operationHashToNonce": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{
				types.Nonce(0): types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				),
			},
			operationHashToNonce: convertListOfOperationsToOperationHashToNonce(
				map[types.Operation]types.Nonce{
					types.NewPreexistingStatefulOrderPlacementOperation(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					): types.Nonce(0),
				},
			),
			nextAvailableNonce: 2,

			expectedOperationHashToNonce: convertListOfOperationsToOperationHashToNonce(
				map[types.Operation]types.Nonce{},
			),
			expectedNextAvailableNonce: 2,
		},
		"nonce map clear results in non-empty operationHashToNonce": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{
				types.Nonce(0): types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				),
			},
			operationHashToNonce: convertListOfOperationsToOperationHashToNonce(
				map[types.Operation]types.Nonce{
					types.NewOrderCancellationOperation(
						&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
					): types.Nonce(1),
					types.NewPreexistingStatefulOrderPlacementOperation(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					): types.Nonce(0),
				},
			),
			nextAvailableNonce: 2,

			expectedOperationHashToNonce: convertListOfOperationsToOperationHashToNonce(
				map[types.Operation]types.Nonce{
					types.NewOrderCancellationOperation(
						&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
					): types.Nonce(1),
				},
			),
			expectedNextAvailableNonce: 2,
		},
		"multiple operation nonce map clear results in non-empty operationHashToNonce": {
			nonceToOperationToPropose: map[types.Nonce]types.Operation{
				types.Nonce(0): types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				),
				types.Nonce(3399): types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				),
				types.Nonce(408): types.NewMatchOperation(
					&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					[]types.MakerFill{},
				),
			},
			operationHashToNonce: convertListOfOperationsToOperationHashToNonce(
				map[types.Operation]types.Nonce{
					// Extra preexisting operation.
					types.NewOrderCancellationOperation(
						&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
					): types.Nonce(1),
					types.NewPreexistingStatefulOrderPlacementOperation(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					): types.Nonce(0),
					types.NewOrderPlacementOperation(
						constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					): types.Nonce(3399),
					types.NewMatchOperation(
						&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
						[]types.MakerFill{},
					): types.Nonce(408),
					// Extra preexisting operation.
					types.NewOrderPlacementOperation(
						constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
					): types.Nonce(12),
				},
			),
			nextAvailableNonce: 500,

			expectedOperationHashToNonce: convertListOfOperationsToOperationHashToNonce(
				map[types.Operation]types.Nonce{
					types.NewOrderCancellationOperation(
						&constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20,
					): types.Nonce(1),
					types.NewOrderPlacementOperation(
						constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
					): types.Nonce(12),
				},
			),
			expectedNextAvailableNonce: 500,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			otp.NonceToOperationToPropose = tc.nonceToOperationToPropose
			otp.OperationHashToNonce = tc.operationHashToNonce
			otp.NextAvailableNonce = tc.nextAvailableNonce

			otp.ClearOperationsQueue()
			// Verify expectations of state.
			require.Empty(t, otp.NonceToOperationToPropose)
			require.Equal(t, tc.expectedOperationHashToNonce, otp.OperationHashToNonce)
			require.Equal(t, tc.expectedNextAvailableNonce, otp.NextAvailableNonce)
			require.Empty(t, otp.GetOperationsQueue())
		})
	}
}

func convertListOfOperationsToOperationHashToNonce(
	operations map[types.Operation]types.Nonce,
) map[types.OperationHash]types.Nonce {
	operationHashToNonces := make(map[types.OperationHash]types.Nonce, len(operations))
	for operation, nonce := range operations {
		hash := operation.GetOperationHash()
		operationHashToNonces[hash] = nonce
	}
	return operationHashToNonces
}

func TestClearOperationsQueue_PanicsOnNoNonceToRemove(t *testing.T) {
	otp := types.NewOperationsToPropose()
	operation := types.NewMatchOperation(
		&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
		[]types.MakerFill{},
	)
	otp.NonceToOperationToPropose = map[types.Nonce]types.Operation{
		types.Nonce(408): operation,
	}
	otp.OperationHashToNonce = map[types.OperationHash]types.Nonce{}
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"ClearOperationsQueue: No nonce to remove for operation %+v",
			operation.GetOperationTextString(),
		),
		func() {
			otp.ClearOperationsQueue()
		},
	)
}

func TestClearOperationsQueue_PanicsOnMismatchedNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	operation := types.NewMatchOperation(
		&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
		[]types.MakerFill{},
	)
	otp.NonceToOperationToPropose = map[types.Nonce]types.Operation{
		types.Nonce(408): operation,
	}
	otp.OperationHashToNonce = map[types.OperationHash]types.Nonce{
		operation.GetOperationHash(): types.Nonce(680),
	}
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"ClearOperationsQueue: Mismatch between nonces for operation (%+v). "+
				"Assigned nonce: %d, nonce in operation queue: %d",
			operation.GetOperationTextString(),
			types.Nonce(680),
			types.Nonce(408),
		),
		func() {
			otp.ClearOperationsQueue()
		},
	)
}

func TestAddMatchToOperationsQueue_PanicsOnStatefulTakerNoNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	takerOrder := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	makerFillsWithOrders := []types.MakerFillWithOrder{
		{
			Order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			MakerFill: types.MakerFill{
				FillAmount: 1,
			},
		},
	}

	nonexistentTakerOp := types.NewPreexistingStatefulOrderPlacementOperation(takerOrder)

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			nonexistentTakerOp.GetOperationTextString(),
		),
		func() {
			otp.AddMatchToOperationsQueue(&takerOrder, makerFillsWithOrders)
		},
	)
}

func TestAddMatchToOperationsQueue_PanicsOnShortTermTakerNoNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	takerOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	makerFillsWithOrders := []types.MakerFillWithOrder{
		{
			Order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			MakerFill: types.MakerFill{
				FillAmount: 1,
			},
		},
	}

	nonExistentTakerOp := types.NewOrderPlacementOperation(takerOrder)

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			nonExistentTakerOp.GetOperationTextString(),
		),
		func() {
			otp.AddMatchToOperationsQueue(&takerOrder, makerFillsWithOrders)
		},
	)
}

func TestAddMatchToOperationsQueue_PanicsOnStatefulMakerNoNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	validTakerOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	validMakerOrder := constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20

	makerFillsWithOrders := []types.MakerFillWithOrder{
		{
			Order: validMakerOrder,
			MakerFill: types.MakerFill{
				FillAmount: 1,
			},
		},
		{
			Order: constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10,
			MakerFill: types.MakerFill{
				FillAmount: 1,
			},
		},
	}

	otp.AssignNonceToOrder(validMakerOrder, false)
	otp.AddOrderPlacementToOperationsQueue(validMakerOrder)
	otp.AssignNonceToOrder(validTakerOrder, false)
	otp.AddOrderPlacementToOperationsQueue(validTakerOrder)

	nonexistentMakerOp := types.NewPreexistingStatefulOrderPlacementOperation(makerFillsWithOrders[1].Order)

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			nonexistentMakerOp.GetOperationTextString(),
		),
		func() {
			otp.AddMatchToOperationsQueue(&validTakerOrder, makerFillsWithOrders)
		},
	)
}

func TestAddMatchToOperationsQueue_PanicsOnMakerNoNonce(t *testing.T) {
	otp := types.NewOperationsToPropose()
	takerOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	makerFillsWithOrders := []types.MakerFillWithOrder{
		{
			Order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			MakerFill: types.MakerFill{
				FillAmount: 1,
			},
		},
	}
	otp.AssignNonceToOrder(takerOrder, false)
	otp.AddOrderPlacementToOperationsQueue(takerOrder)

	nonexistentMakerOp := types.NewOrderPlacementOperation(makerFillsWithOrders[0].Order)

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"mustGetNonceFromOperation: operation (%+v) has no nonce",
			nonexistentMakerOp.GetOperationTextString(),
		),
		func() {
			otp.AddMatchToOperationsQueue(&takerOrder, makerFillsWithOrders)
		},
	)
}

func TestAddMatchToOperationsQueue_Success(t *testing.T) {
	tests := map[string]struct {
		// State.
		takerOrder                 types.MatchableOrder
		makerFillsWithOrders       []types.MakerFillWithOrder
		isPreexistingStatefulOrder map[types.OrderId]bool

		// Expectations.
		expectedOperationsQueue map[types.Nonce]types.Operation
	}{
		"Match a single pre-existing stateful taker order with Short-Term maker orders": {
			takerOrder: &constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			makerFillsWithOrders: []types.MakerFillWithOrder{
				{
					Order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					MakerFill: types.MakerFill{
						FillAmount:   5,
						MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
					},
				},
			},
			isPreexistingStatefulOrder: map[types.OrderId]bool{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId: true,
			},
			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				1: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				),
				2: types.NewMatchOperation(
					&constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
						},
					},
				),
			},
		},
		"Match a liquidation order with one pre-existing stateful order and one newly-placed stateful order": {
			takerOrder: &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
			makerFillsWithOrders: []types.MakerFillWithOrder{
				{
					Order: constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10,
					MakerFill: types.MakerFill{
						FillAmount:   1,
						MakerOrderId: constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId,
					},
				},
				{
					Order: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
					MakerFill: types.MakerFill{
						FillAmount:   1,
						MakerOrderId: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO.OrderId,
					},
				},
			},
			isPreexistingStatefulOrder: map[types.OrderId]bool{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO.OrderId: true,
			},
			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderPlacementOperation(constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10),
				1: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,
				),
				2: types.NewMatchOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					[]types.MakerFill{
						{
							FillAmount:   1,
							MakerOrderId: constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId,
						},
						{
							FillAmount:   1,
							MakerOrderId: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO.OrderId,
						},
					},
				),
			},
		},
		"Match a liquidation order with one newly-placed stateful order and one pre-existing stateful order": {
			takerOrder: &constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
			makerFillsWithOrders: []types.MakerFillWithOrder{
				{
					Order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					MakerFill: types.MakerFill{
						FillAmount:   5,
						MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
					},
				},
			},
			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				1: types.NewMatchOperation(
					&constants.LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
						},
					},
				),
			},
		},
		"Match a single Short-Term taker order with Short-Term maker orders": {
			takerOrder: &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			makerFillsWithOrders: []types.MakerFillWithOrder{
				{
					Order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					MakerFill: types.MakerFill{
						FillAmount:   5,
						MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
					},
				},
			},
			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				1: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
				2: types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
						},
					},
				),
			},
		},
		"Match a single Short-Term taker order with one Short-Term maker order and one pre-existing Stateful maker order": {
			takerOrder: &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			makerFillsWithOrders: []types.MakerFillWithOrder{
				{
					Order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					MakerFill: types.MakerFill{
						FillAmount:   1,
						MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
					},
				},
				{
					Order: constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10,
					MakerFill: types.MakerFill{
						FillAmount:   1,
						MakerOrderId: constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId,
					},
				},
			},
			isPreexistingStatefulOrder: map[types.OrderId]bool{
				constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId: true,
			},
			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				1: types.NewPreexistingStatefulOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10,
				),
				2: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
				3: types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							FillAmount:   1,
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
						},
						{
							FillAmount:   1,
							MakerOrderId: constants.LongTermOrder_Alice_Num1_Id0_Clob0_Sell15_Price5_GTBT10.OrderId,
						},
					},
				),
			},
		},
		"Match a single Short-Term taker order with multiple Short-Term maker orders": {
			takerOrder: &constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			makerFillsWithOrders: []types.MakerFillWithOrder{
				{
					Order: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					MakerFill: types.MakerFill{
						FillAmount:   1,
						MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
					},
				},
				{
					Order: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					MakerFill: types.MakerFill{
						FillAmount:   1,
						MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
					},
				},
			},
			expectedOperationsQueue: map[types.Nonce]types.Operation{
				0: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				1: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15),
				2: types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15),
				3: types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					[]types.MakerFill{
						{
							FillAmount:   1,
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
						},
						{
							FillAmount:   1,
							MakerOrderId: constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15.OrderId,
						},
					},
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()

			// Insert all maker order placements
			for _, makerFillWithOrder := range tc.makerFillsWithOrders {
				order := makerFillWithOrder.Order
				isPreexistingStatefulOrder := tc.isPreexistingStatefulOrder[order.OrderId]
				otp.AssignNonceToOrder(order, isPreexistingStatefulOrder)
				if isPreexistingStatefulOrder {
					otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(order)
				} else {
					otp.AddOrderPlacementToOperationsQueue(order)
				}
			}

			// Insert taker order placement if it's not a liquidation
			if !tc.takerOrder.IsLiquidation() {
				order := tc.takerOrder.MustGetOrder()
				isPreexistingStatefulOrder := tc.isPreexistingStatefulOrder[order.OrderId]
				otp.AssignNonceToOrder(order, isPreexistingStatefulOrder)
				if isPreexistingStatefulOrder {
					otp.AddPreexistingStatefulOrderPlacementToOperationsQueue(order)
				} else {
					otp.AddOrderPlacementToOperationsQueue(order)
				}
			}

			// Insert match operation
			otp.AddMatchToOperationsQueue(tc.takerOrder, tc.makerFillsWithOrders)

			// Verify expectations.
			require.Equal(t, tc.expectedOperationsQueue, otp.NonceToOperationToPropose)

			// Verify the next available nonce is correct.
			require.Equal(t, types.Nonce(len(tc.expectedOperationsQueue)), otp.NextAvailableNonce)
		})
	}
}
