package memclob

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4/testutil/constants"
	testutil_memclob "github.com/dydxprotocol/v4/testutil/memclob"
	sdktest "github.com/dydxprotocol/v4/testutil/sdk"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestShortTermCancelOrder_CancelAlreadyExists(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	memclob := NewMemClobPriceTimePriority(true)
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15

	// Create all unique orderbooks.
	createAllOrderbooksForOrders(
		t,
		ctx,
		memclob,
		[]types.Order{order},
	)

	offchainUpdates, err := memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, 100))
	require.NoError(t, err)
	testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, order.OrderId)
	_, err = memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, 99))
	require.Equal(t, types.ErrMemClobCancelAlreadyExists, err)
	_, err = memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, 100))
	require.Equal(t, types.ErrMemClobCancelAlreadyExists, err)
	offchainUpdates, err = memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(order.OrderId, 101))
	require.NoError(t, err)
	testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, order.OrderId)
}

func TestShortTermCancelOrder_OrdersTilBlockExceedsCancels(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	memClobKeeper := testutil_memclob.NewFakeMemClobKeeper()
	memclob := NewMemClobPriceTimePriority(true)
	memclob.SetClobKeeper(memClobKeeper)

	order := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	orderId := order.OrderId
	cancelGoodTilBlock1 := order.GetGoodTilBlock() - 2
	cancelGoodTilBlock2 := order.GetGoodTilBlock() - 1

	// Create all unique orderbooks.
	createAllOrderbooksForOrders(
		t,
		ctx,
		memclob,
		[]types.Order{order},
	)

	_, _, _, err := memclob.PlaceOrder(
		ctx,
		order,
		true,
	)
	require.NoError(t, err)

	// Cancel without error once.
	offchainUpdates, err := memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(orderId, cancelGoodTilBlock1))
	require.NoError(t, err)
	testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, order.OrderId)
	cancelBlock, isCanceled := memclob.cancels.get(orderId)
	require.True(t, isCanceled)
	require.Equal(t, cancelGoodTilBlock1, cancelBlock)

	// Cancel without error again.
	offchainUpdates, err = memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(orderId, cancelGoodTilBlock2))
	require.NoError(t, err)
	testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, order.OrderId)
	cancelBlock, isCanceled = memclob.cancels.get(orderId)
	require.True(t, isCanceled)
	require.Equal(t, cancelGoodTilBlock2, cancelBlock)

	// Order is still on the book.
	gottenOrder, found := memclob.GetOrder(ctx, orderId)
	require.True(t, found)
	require.Equal(t, order, gottenOrder)
}

func TestCancelOrder(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	tests := map[string]struct {
		// State.
		existingOrders  []types.Order
		existingCancels []*types.MsgCancelOrder

		// Parameters.
		order        types.Order
		goodTilBlock uint32

		// Expectations.
		expectedBestBid                              types.Subticks
		expectedBestAsk                              types.Subticks
		expectedTotalLevels                          int
		expectedTotalLevelQuantums                   uint64
		expectLevelToExist                           bool
		expectBlockExpirationsForOrdersToExist       bool
		expectSubaccountOpenClobOrdersForSideToExist bool
		expectSubaccountOpenClobOrdersToExist        bool
		expectedCanceledOrdersLen                    int
		expectedCancelOrderExpirationsLen            int
		expectedCancelOrderExpirationsForTilBlockLen int
	}{
		"Cancels an order where an existing cancel does not exist": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			order:        constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			goodTilBlock: 123,

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectedCanceledOrdersLen:                    1,
			expectedCancelOrderExpirationsLen:            1,
			expectedCancelOrderExpirationsForTilBlockLen: 1,
		},
		"Cancels an order when the order does not exist, and no previous cancel exists": {
			existingOrders: []types.Order{},
			order:          constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			goodTilBlock:   123,

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectedCanceledOrdersLen:                    1,
			expectedCancelOrderExpirationsLen:            1,
			expectedCancelOrderExpirationsForTilBlockLen: 1,
		},
		"Cancels an order where an existing cancel already exists with a lower `goodTilBlock`": {
			existingOrders: []types.Order{},
			existingCancels: []*types.MsgCancelOrder{
				{
					OrderId:      constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 123},
				},
			},
			order:        constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			goodTilBlock: 124,

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectedCanceledOrdersLen:                    1,
			expectedCancelOrderExpirationsLen:            1,
			expectedCancelOrderExpirationsForTilBlockLen: 1,
		},
		"Cancels an order where an existing cancel already exists for this `goodTilBlock`, but for a different order": {
			existingOrders: []types.Order{
				constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
				constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
				constants.Order_Bob_Num0_Id7_Clob0_Buy20_Price10000_GTB22,
			},
			existingCancels: []*types.MsgCancelOrder{
				{
					OrderId:      constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22.OrderId,
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 123},
				},
				{
					OrderId:      constants.Order_Bob_Num0_Id2_Clob1_Sell12_Price13_GTB20.OrderId,
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 120},
				},
			},
			order:        constants.Order_Bob_Num0_Id7_Clob0_Buy20_Price10000_GTB22,
			goodTilBlock: 123,

			expectedTotalLevels:                          2,
			expectedBestBid:                              1000,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
			expectedCanceledOrdersLen:                    3,
			expectedCancelOrderExpirationsLen:            2,
			expectedCancelOrderExpirationsForTilBlockLen: 2,
		},
		`Cancels an order where an existing cancel already exists, and an existing cancel
			already exists for this goodTilBlock for a different order`: {
			existingOrders: []types.Order{
				constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
				constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
			},
			existingCancels: []*types.MsgCancelOrder{
				{
					OrderId:      constants.Order_Bob_Num0_Id4_Clob1_Buy20_Price35_GTB22.OrderId,
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 123},
				},
				{
					OrderId:      constants.Order_Bob_Num0_Id7_Clob0_Buy20_Price10000_GTB22.OrderId,
					GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 122},
				},
			},
			order:        constants.Order_Bob_Num0_Id7_Clob0_Buy20_Price10000_GTB22,
			goodTilBlock: 123,

			expectedTotalLevels:                          2,
			expectedBestBid:                              1000,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
			expectedCanceledOrdersLen:                    2,
			expectedCancelOrderExpirationsLen:            1,
			expectedCancelOrderExpirationsForTilBlockLen: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the memclob state.
			memclob := NewMemClobPriceTimePriority(true)
			memclob.SetClobKeeper(testutil_memclob.NewFakeMemClobKeeper())

			// Create all unique orderbooks.
			createAllOrderbooksForOrders(
				t,
				ctx,
				memclob,
				append(tc.existingOrders, tc.order),
			)

			// Place all existing orders on the orderbook.
			for _, order := range tc.existingOrders {
				_, _, _, err := memclob.PlaceOrder(
					ctx,
					order,
					true,
				)
				require.NoError(t, err)
			}

			// Place all existing cancels on the orderbook.
			for _, cancel := range tc.existingCancels {
				offchainUpdates, err := memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(
					cancel.OrderId,
					cancel.GetGoodTilBlock(),
				))
				require.NoError(t, err)
				testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, cancel.OrderId)
			}

			// Run the test case.
			offchainUpdates, err := memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(tc.order.OrderId, tc.goodTilBlock))
			require.NoError(t, err)
			testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, tc.order.OrderId)

			// Verify that the memclob orderbook has the correct state.
			requireOrderDoesNotExistInMemclob(t, ctx, tc.order, memclob)
			for _, existingOrder := range tc.existingOrders {
				if existingOrder != tc.order {
					requireOrderExistsInMemclob(t, ctx, existingOrder, memclob)
				}
			}

			// Verify that the cancel now exists in the memclob.
			block, exists := memclob.cancels.get(tc.order.OrderId)
			require.Equal(t, tc.goodTilBlock, block)
			require.True(t, exists)
			require.True(t, memclob.cancels.expiryToOrderIds[tc.goodTilBlock][tc.order.OrderId])

			// Verify public method matches private method.
			publicBlock, publicExists := memclob.GetCancelOrder(ctx, tc.order.OrderId)
			require.Equal(t, block, publicBlock)
			require.Equal(t, exists, publicExists)

			// Verify that the cancel data structures meet expectations.
			require.Len(t, memclob.cancels.orderIdToExpiry, tc.expectedCanceledOrdersLen)
			require.Len(t, memclob.cancels.expiryToOrderIds, tc.expectedCancelOrderExpirationsLen)
			require.Len(t, memclob.cancels.expiryToOrderIds[tc.goodTilBlock], tc.expectedCancelOrderExpirationsForTilBlockLen)

			// Enforce various expectations around the in-memory data structures for
			// orders in the clob.
			assertOrderbookStateExpectations(
				t,
				memclob,
				tc.order,
				tc.expectedBestBid,
				tc.expectedBestAsk,
				tc.expectedTotalLevels,
				tc.expectLevelToExist,
				tc.expectBlockExpirationsForOrdersToExist,
				tc.expectSubaccountOpenClobOrdersForSideToExist,
				tc.expectSubaccountOpenClobOrdersToExist,
			)
		})
	}
}

func TestCancelOrder_Telemetry(t *testing.T) {
	m, err := telemetry.New(telemetry.Config{
		Enabled:        true,
		EnableHostname: false,
		ServiceName:    "test",
	})
	require.NoError(t, err)
	require.NotNil(t, m)

	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	memclob := NewMemClobPriceTimePriority(true)
	memclob.SetClobKeeper(testutil_memclob.NewFakeMemClobKeeper())

	orderOne := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	orderTwo := constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20

	// Create all unique orderbooks.
	createAllOrderbooksForOrders(
		t,
		ctx,
		memclob,
		[]types.Order{orderOne, orderTwo},
	)

	_, _, _, err = memclob.PlaceOrder(
		ctx,
		orderOne,
		true,
	)
	require.NoError(t, err)
	_, _, _, err = memclob.PlaceOrder(
		ctx,
		orderTwo,
		true,
	)
	require.NoError(t, err)

	// Cancel both orders.
	offchainUpdates, err := memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(orderOne.OrderId, 100))
	require.NoError(t, err)
	testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, orderOne.OrderId)
	offchainUpdates, err = memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(orderTwo.OrderId, 100))
	require.NoError(t, err)
	testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, orderTwo.OrderId)

	gr, err := m.Gather(telemetry.FormatText)
	require.NoError(t, err)
	require.Equal(t, "application/json", gr.ContentType)

	jsonMetrics := make(map[string]interface{})
	require.NoError(t, json.Unmarshal(gr.Metrics, &jsonMetrics))

	counters := jsonMetrics["Counters"].([]any)
	require.Condition(t, func() bool {
		for _, counter := range counters {
			if counter.(map[string]any)["Name"].(string) == "test.clob.cancel_order.removed_from_orderbook" &&
				counter.(map[string]any)["Count"].(float64) == 2.0 {
				return true
			}
		}
		return false
	})

	samples := jsonMetrics["Samples"].([]interface{})
	require.Condition(t, func() bool {
		for _, sample := range samples {
			if sample.(map[string]any)["Name"].(string) == "test.memclob.removed_from_orderbook.latency" &&
				sample.(map[string]any)["Count"].(float64) == 2.0 {
				return true
			}
		}
		return false
	})
}

func TestCancelOrder_AddToOperationsQueue(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedOperations       []types.Operation
		collateralizationCheck map[int]testutil_memclob.CollateralizationCheck

		// Expectations.
		expectedOperations       []types.Operation
		expectedOperationToNonce map[types.Operation]types.Nonce
		expectedErr              error
	}{
		`Stateful order cancellation added to operations queue if order has not been seen`: {
			placedOperations: []types.Operation{
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5,
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{},

			expectedOperations: []types.Operation{
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5,
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderCancellationOperation(&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5): 0,
			},
		},
		`Stateful order cancellation added to operations queue if order has been seen`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
				),
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5,
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
			},

			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
				),
				types.NewOrderCancellationOperation(
					&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5,
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5): 0,
				types.NewOrderCancellationOperation(&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT5):    1,
			},
		},
		`Can cancel a non-existent order`: {
			placedOperations: []types.Operation{
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.OrderId,
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{},

			expectedOperations:       []types.Operation{},
			expectedOperationToNonce: map[types.Operation]types.Nonce{},
		},
		`Can cancel a partially-matched order and the cancellation is added to the operations queue`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 30,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Alice_Num1: satypes.Success,
					},
				},
			},

			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25): 0,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20):  1,
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				): 2,
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				): 3,
			},
		},
		`Can cancel a partially-matched order multiple times and only the first cancellation is
					added to the operations queue`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 30,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Alice_Num1: satypes.Success,
					},
				},
			},

			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25): 0,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20):  1,
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				): 2,
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				): 3,
			},
		},
		`Canceling an unmatched order multiple times does not add a cancel to the operations queue`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 30,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
			},

			expectedOperations:       []types.Operation{},
			expectedOperationToNonce: map[types.Operation]types.Nonce{},
		},
		`Canceling a fully-matched order multiple times does not add a cancel to the operations queue`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				types.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 30,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 30,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 30,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
				2: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Bob_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          35,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Bob_Num0: satypes.Success,
					},
				},
			},

			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
				),
				types.NewMatchOperation(
					&constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   30,
						},
					},
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25): 0,
				types.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30):  1,
				types.NewMatchOperation(
					&constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   30,
						},
					},
				): 2,
			},
		},
		`Canceling multiple partially matched orders adds multiple cancels to the operations queue`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				// #1 partially-matched order.
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				// Cancel partially-matched order #1.
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				// #2 partially-matched order.
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34),
				// #3 partially-matched order.
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15),
				// Cancel partially-matched order #3.
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20.OrderId,
						constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20.GetGoodTilBlock(),
					),
				),
				// Cancel partially-matched order #2.
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.OrderId,
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           true,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Alice_Num1: satypes.Success,
					},
				},
				2: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
				3: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 25,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          15,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Bob_Num0: satypes.Success,
					},
				},
				4: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          15,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           true,
								Subticks:          15,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Alice_Num1: satypes.Success,
					},
				},
				5: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 10,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          45,
								ClobPairId:        1,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
				6: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          45,
								ClobPairId:        1,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          45,
								ClobPairId:        1,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Alice_Num1: satypes.Success,
					},
				},
			},

			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20.OrderId,
						constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20.GetGoodTilBlock(),
					),
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.OrderId,
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.GetGoodTilBlock(),
					),
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20):  0,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25): 1,
				types.NewMatchOperation(
					&constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				): 2,
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				): 3,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20): 4,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34):   5,
				types.NewMatchOperation(
					&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				): 6,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20): 7,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15):  8,
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				): 9,
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_User1_Num1_Id11_Clob1_GTB20.OrderId,
						constants.CancelOrder_User1_Num1_Id11_Clob1_GTB20.GetGoodTilBlock(),
					),
				): 10,
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_User1_Num0_Id10_Clob0_GTB20.OrderId,
						constants.CancelOrder_User1_Num0_Id10_Clob0_GTB20.GetGoodTilBlock(),
					),
				): 11,
			},
		},
		`Can cancel a partially-matched order, replace it, then re-cancel the order multiple times
					and only the first cancellation is added to the operations queue`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 30,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Alice_Num1: satypes.Success,
					},
				},
				2: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 45,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
			},

			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25): 0,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20):  1,
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				): 2,
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				): 3,
			},
		},
		`Can cancel a partially-matched order, replace it and it's partially-matched again, then
					re-cancel the order multiple times and only two cancellations are added to the
					operations queue`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				types.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				// Collateralization checks for first order placement.
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 30,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
					},
				},
				// Collateralization checks for second order placement.
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Alice_Num1: satypes.Success,
					},
				},
				// Collateralization checks for third order placement.
				2: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Bob_Num0: {
							{
								RemainingQuantums: 20,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Bob_Num0: satypes.Success,
					},
				},
				// Collateralization checks for fourth order placement.
				3: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 20,
								IsBuy:             true,
								IsTaker:           true,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 20,
								IsBuy:             false,
								IsTaker:           false,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
						constants.Bob_Num0:   satypes.Success,
					},
				},
				4: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 25,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
			},

			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				types.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   20,
						},
					},
				),
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25): 0,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20):  1,
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				): 2,
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				): 3,
				types.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22):   4,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30): 5,
				types.NewMatchOperation(
					&constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   20,
						},
					},
				): 6,
				types.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				): 7,
			},
		},
		`Can replace a partially-matched order then cancel it, and no cancellation is added to the
					operations queue`: {
			placedOperations: []types.Operation{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 30,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								IsTaker:           true,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num0: satypes.Success,
						constants.Alice_Num1: satypes.Success,
					},
				},
				2: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 45,
								IsBuy:             true,
								IsTaker:           false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
					},
					Result: map[satypes.SubaccountId]satypes.UpdateResult{
						constants.Alice_Num1: satypes.Success,
					},
				},
			},

			expectedOperations: []types.Operation{
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
				),
			},
			expectedOperationToNonce: map[types.Operation]types.Nonce{
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25): 0,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20):  1,
				types.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				): 2,
				types.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30): 3,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			memclob, _ := memclobOperationsTestSetup(
				t,
				ctx,
				tc.placedOperations,
				tc.collateralizationCheck,
				constants.GetStatePosition_ZeroPositionSize,
				[]types.StatefulOrderPlacement{},
			)

			assertMemclobHasOperations(
				t,
				ctx,
				memclob,
				tc.expectedOperations,
				tc.expectedOperationToNonce,
			)

			// TODO(DEC-1587): Verify the correct offchain update messages
			// were returned for order cancellations.
		})
	}
}
