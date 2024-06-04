package memclob

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/cosmos/cosmos-sdk/telemetry"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil_memclob "github.com/dydxprotocol/v4-chain/protocol/testutil/memclob"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestShortTermCancelOrder_CancelAlreadyExists(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	memclob := NewMemClobPriceTimePriority(true)
	memclob.SetClobKeeper(testutil_memclob.NewFakeMemClobKeeper())
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
	ctx = ctx.WithIsCheckTx(true)
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

	_, _, _, err := memclob.PlaceOrder(ctx, order)
	require.NoError(t, err)

	// Cancel without error once.
	offchainUpdates, err := memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(orderId, cancelGoodTilBlock1))
	require.NoError(t, err)
	testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, order.OrderId)
	cancelBlock, isCanceled := memclob.GetCancelOrder(orderId)
	require.True(t, isCanceled)
	require.Equal(t, cancelGoodTilBlock1, cancelBlock)

	// Cancel without error again.
	offchainUpdates, err = memclob.CancelOrder(ctx, types.NewMsgCancelOrderShortTerm(orderId, cancelGoodTilBlock2))
	require.NoError(t, err)
	testutil_memclob.RequireCancelOrderOffchainUpdate(t, ctx, offchainUpdates, order.OrderId)
	cancelBlock, isCanceled = memclob.GetCancelOrder(orderId)
	require.True(t, isCanceled)
	require.Equal(t, cancelGoodTilBlock2, cancelBlock)

	// Order is still on the book.
	gottenOrder, found := memclob.GetOrder(orderId)
	require.True(t, found)
	require.Equal(t, order, gottenOrder)
}

func TestCancelOrder_PanicsOnStatefulOrder(t *testing.T) {
	memclob := NewMemClobPriceTimePriority(true)
	orderId := constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)

	expectedError := fmt.Sprintf(
		"MustBeShortTermOrder: called with stateful order ID (%+v)",
		orderId,
	)

	memclob.CreateOrderbook(constants.ClobPair_Btc)
	require.PanicsWithValue(t, expectedError, func() {
		//nolint:errcheck
		memclob.CancelOrder(ctx, types.NewMsgCancelOrderStateful(orderId, 100))
	})
}

func TestCancelOrder(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
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
				constants.Order_Bob_Num0_Id8_Clob1_Sell20_Price10_GTB22,
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
			expectedCanceledOrdersLen:                    1,
			expectedCancelOrderExpirationsLen:            1,
			expectedCancelOrderExpirationsForTilBlockLen: 1,
		},
		`Cancels an order where an existing cancel already exists, and an existing cancel
			already exists for this goodTilBlock for a different order`: {
			existingOrders: []types.Order{
				constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
				constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
				constants.Order_Bob_Num0_Id8_Clob1_Sell20_Price10_GTB22,
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
			expectedCanceledOrdersLen:                    1,
			expectedCancelOrderExpirationsLen:            1,
			expectedCancelOrderExpirationsForTilBlockLen: 1,
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
				_, _, _, err := memclob.PlaceOrder(ctx, order)
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
			orderbook := memclob.mustGetOrderbook(tc.order.GetClobPairId())
			block, exists := orderbook.getCancel(tc.order.OrderId)
			require.Equal(t, tc.goodTilBlock, block)
			require.True(t, exists)
			require.True(t, orderbook.cancelExpiryToOrderIds[tc.goodTilBlock][tc.order.OrderId])

			// Verify public method matches private method.
			publicBlock, publicExists := memclob.GetCancelOrder(tc.order.OrderId)
			require.Equal(t, block, publicBlock)
			require.Equal(t, exists, publicExists)

			// Verify that the cancel data structures meet expectations.
			require.Len(t, orderbook.orderIdToCancelExpiry, tc.expectedCanceledOrdersLen)
			require.Len(t, orderbook.cancelExpiryToOrderIds, tc.expectedCancelOrderExpirationsLen)
			require.Len(t, orderbook.cancelExpiryToOrderIds[tc.goodTilBlock], tc.expectedCancelOrderExpirationsForTilBlockLen)

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
	ctx = ctx.WithIsCheckTx(true)
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

	_, _, _, err = memclob.PlaceOrder(ctx, orderOne)
	require.NoError(t, err)
	_, _, _, err = memclob.PlaceOrder(ctx, orderTwo)
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
			if counter.(map[string]any)["Name"].(string) == "test.clob.cancel_short_term_order.removed_from_orderbook" &&
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

func TestCancelOrder_OperationsQueue(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		placedOperations       []types.Operation
		collateralizationCheck map[int]testutil_memclob.CollateralizationCheck

		// Expectations.
		expectedOperations         []types.Operation
		expectedInternalOperations []types.InternalOperation
		expectedErr                error
	}{
		`Can cancel a non-existent order`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.OrderId,
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{},

			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`Can cancel a partially-matched order`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
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
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Can cancel a partially-matched order multiple times`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.GetGoodTilBlock(),
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
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Canceling an unmatched order multiple times`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.GetGoodTilBlock(),
					),
				),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{},

			expectedOperations:         []types.Operation{},
			expectedInternalOperations: []types.InternalOperation{},
		},
		`Canceling a fully-matched order multiple times`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderCancellationOperation(
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
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 30,
								IsBuy:             false,
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
			},

			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   30,
						},
					},
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Bob_Num0_Id13_Clob0_Sell35_Price35_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   30,
						},
					},
				),
			},
		},
		`Can cancel multiple partially matched orders`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				// #1 partially-matched order.
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				// Cancel partially-matched order #1.
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				// #2 partially-matched order.
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34),
				// #3 partially-matched order.
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15),
				// Cancel partially-matched order #3.
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20.OrderId,
						constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20.GetGoodTilBlock(),
					),
				),
				// Cancel partially-matched order #2.
				clobtest.NewOrderCancellationOperation(
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
								Subticks:          10,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								Subticks:          15,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				2: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								Subticks:          45,
								ClobPairId:        1,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20.OrderId,
						constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.OrderId,
						constants.CancelOrder_Alice_Num0_Id10_Clob0_GTB20.GetGoodTilBlock(),
					),
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB34,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num0_Id10_Clob0_Sell25_Price15_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id11_Clob1_Buy10_Price45_GTB20.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Can cancel a partially-matched order, replace it, then re-cancel the order multiple times`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.GetGoodTilBlock(),
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
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
		},
		`Can cancel a partially-matched order, replace it and it's partially-matched again, then
					re-cancel the order multiple times`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB35.GetGoodTilBlock(),
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
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				1: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num1: {
							{
								RemainingQuantums: 20,
								IsBuy:             true,
								Subticks:          10,
								ClobPairId:        0,
							},
						},
						constants.Bob_Num0: {
							{
								RemainingQuantums: 20,
								IsBuy:             false,
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
			},

			expectedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25.GetGoodTilBlock(),
					),
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   20,
						},
					},
				),
				clobtest.NewOrderCancellationOperation(
					types.NewMsgCancelOrderShortTerm(
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.OrderId,
						constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30.GetGoodTilBlock(),
					),
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.OrderId,
							FillAmount:   20,
						},
					},
				),
			},
		},
		`Can replace a partially-matched order then cancel it`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30),
			},
			collateralizationCheck: map[int]testutil_memclob.CollateralizationCheck{
				0: {
					CollatCheck: map[satypes.SubaccountId][]types.PendingOpenOrder{
						constants.Alice_Num0: {
							{
								RemainingQuantums: 5,
								IsBuy:             false,
								Subticks:          50,
								ClobPairId:        0,
							},
						},
						constants.Alice_Num1: {
							{
								RemainingQuantums: 5,
								IsBuy:             true,
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
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				clobtest.NewMatchOperation(
					&constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
				clobtest.NewOrderPlacementOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30,
				),
			},
			expectedInternalOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Alice_Num1_Id13_Clob0_Buy30_Price50_GTB25.OrderId,
							FillAmount:   5,
						},
					},
				),
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
				[]types.LongTermOrderPlacement{},
			)

			assertMemclobHasOperations(
				t,
				ctx,
				memclob,
				tc.expectedOperations,
				tc.expectedInternalOperations,
			)

			// TODO(DEC-1587): Verify the correct offchain update messages
			// were returned for order cancellations.
		})
	}
}
