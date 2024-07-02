package memclob

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	memclobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/memclob"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPurgeInvalidMemclobState(t *testing.T) {
	tests := map[string]struct {
		// State.
		placedOperations    []types.Operation
		newOrderFillAmounts map[types.OrderId]satypes.BaseQuantums
		existingCancels     []*types.MsgCancelOrder

		// Parameters.
		fullyFilledOrderIds      []types.OrderId
		expiredStatefulOrderIds  []types.OrderId
		canceledStatefulOrderIds []types.OrderId
		removedStatefulOrderIds  []types.OrderId

		// Expectations.
		expectedRemainingBids          []OrderWithRemainingSize
		expectedRemovedExpiredOrderIds []types.OrderId
		expectedRemainingAsks          []OrderWithRemainingSize
		expectedCancels                []*types.MsgCancelOrder
	}{
		`Empty memclob and no fully-filled or expired order IDs`: {
			placedOperations:    []types.Operation{},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{},

			fullyFilledOrderIds:     []types.OrderId{},
			expiredStatefulOrderIds: []types.OrderId{},

			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
		`Memclob has expired Short-Term orders and Short-Term order cancellations and they're removed`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15),
				clobtest.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10),
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id2_Clob0_Sell25_Price95_GTB10),
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num1_Id1_Clob1_Sell25_Price85_GTB10),
			},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{},
			existingCancels: []*types.MsgCancelOrder{
				&constants.CancelOrder_Alice_Num0_Id12_Clob0_GTB5,
				&constants.CancelOrder_Bob_Num0_Id2_Clob1_GTB5,
				&constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20,
				&constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30,
			},

			fullyFilledOrderIds:     []types.OrderId{},
			expiredStatefulOrderIds: []types.OrderId{},
			expectedCancels: []*types.MsgCancelOrder{
				&constants.CancelOrder_Alice_Num1_Id11_Clob1_GTB20,
				&constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB30,
			},
			expectedRemovedExpiredOrderIds: []types.OrderId{
				constants.Order_Bob_Num0_Id2_Clob0_Sell25_Price95_GTB10.OrderId,
				constants.Order_Bob_Num1_Id1_Clob1_Sell25_Price85_GTB10.OrderId,
			},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
					RemainingSize: 5,
				},
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					RemainingSize: 45,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
		`There are fully-filled and expired order IDs but the memclob is empty`: {
			placedOperations:    []types.Operation{},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{},

			fullyFilledOrderIds: []types.OrderId{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
				constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30.OrderId,
				constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO.OrderId,
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
			},
			expectedRemovedExpiredOrderIds: []types.OrderId{},
			expiredStatefulOrderIds: []types.OrderId{
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
			},
			expectedRemainingBids: []OrderWithRemainingSize{},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
		`There are fully-filled and expired orders but none exist on the memclob`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15),
				clobtest.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10),
			},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{},

			fullyFilledOrderIds: []types.OrderId{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
				constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30.OrderId,
				constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO.OrderId,
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
			},
			expectedRemovedExpiredOrderIds: []types.OrderId{},
			expiredStatefulOrderIds: []types.OrderId{
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
			},

			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32,
					RemainingSize: 35,
				},
				{
					Order:         constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
					RemainingSize: 5,
				},
				{
					Order:         constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10,
					RemainingSize: 45,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
		`There are fully-filled and expired orders and they're removed from the memclob`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32),
				clobtest.NewOrderPlacementOperation(constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15),
				clobtest.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10),
				clobtest.NewOrderPlacementOperation(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15),
			},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32.OrderId:           35,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId: 5,
			},

			fullyFilledOrderIds: []types.OrderId{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32.OrderId, // This order is on the memclob.
				constants.Order_Alice_Num1_Id13_Clob0_Buy50_Price50_GTB30.OrderId,
				constants.Order_Alice_Num0_Id1_Clob0_Buy15_Price10_GTB18_PO.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId, // This order is on the memclob.
			},
			expectedRemovedExpiredOrderIds: []types.OrderId{
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId, // This order is on the memclob.
			},
			expiredStatefulOrderIds: []types.OrderId{
				constants.LongTermOrder_Alice_Num1_Id1_Clob0_Sell25_Price30_GTBT10.OrderId,
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId, // This order is on the memclob.
			},

			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
		`An order is specified as fully-filled in the last block but it's not removed since it's
			only partially-filled on the validator's local memclob`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32),
			},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32.OrderId: 30,
			},

			fullyFilledOrderIds: []types.OrderId{
				constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32.OrderId, // This order is on the memclob.
			},

			expectedRemovedExpiredOrderIds: []types.OrderId{},
			expiredStatefulOrderIds:        []types.OrderId{},
			expectedRemainingBids: []OrderWithRemainingSize{
				{
					Order:         constants.Order_Bob_Num0_Id1_Clob0_Buy35_Price55_GTB32,
					RemainingSize: 5,
				},
			},
			expectedRemainingAsks: []OrderWithRemainingSize{},
		},
		`An order is canceled and removed from the memclob`: {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10),
			},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{},
			canceledStatefulOrderIds: []types.OrderId{
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
			},

			fullyFilledOrderIds: []types.OrderId{},

			expectedRemovedExpiredOrderIds: []types.OrderId{},
			expiredStatefulOrderIds:        []types.OrderId{},
			expectedRemainingBids:          []OrderWithRemainingSize{},
			expectedRemainingAsks:          []OrderWithRemainingSize{},
		},
		`An order is canceled, but it is not present on the memclob so it is a no-op`: {
			placedOperations:    []types.Operation{},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{},
			canceledStatefulOrderIds: []types.OrderId{
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
			},

			fullyFilledOrderIds: []types.OrderId{},

			expectedRemovedExpiredOrderIds: []types.OrderId{},
			expiredStatefulOrderIds:        []types.OrderId{},
			expectedRemainingBids:          []OrderWithRemainingSize{},
			expectedRemainingAsks:          []OrderWithRemainingSize{},
		},
		"An order in RemovedStatefulOrderIds is removed from the memclob": {
			placedOperations: []types.Operation{
				clobtest.NewOrderPlacementOperation(constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10),
			},
			newOrderFillAmounts: map[types.OrderId]satypes.BaseQuantums{},
			removedStatefulOrderIds: []types.OrderId{
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Buy45_Price10_GTBT10.OrderId,
			},

			fullyFilledOrderIds: []types.OrderId{},

			expectedRemovedExpiredOrderIds: []types.OrderId{},
			expiredStatefulOrderIds:        []types.OrderId{},
			expectedRemainingBids:          []OrderWithRemainingSize{},
			expectedRemainingAsks:          []OrderWithRemainingSize{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state.
			ctx, _, _ := sdktest.NewSdkContextWithMultistore()
			ctx = ctx.WithIsCheckTx(true)
			memclob := NewMemClobPriceTimePriority(true)
			mockMemClobKeeper := &mocks.MemClobKeeper{}
			memclob.SetClobKeeper(mockMemClobKeeper)
			mockMemClobKeeper.On("Logger", mock.Anything).Return(log.NewNopLogger()).Maybe()
			mockMemClobKeeper.On("SendOrderbookUpdates", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

			for _, operation := range tc.placedOperations {
				switch operation.Operation.(type) {
				case *types.Operation_ShortTermOrderPlacement:
					order := operation.GetShortTermOrderPlacement().Order
					orderId := order.OrderId
					// Mock out calls to GetOrderFillAmount during test setup.
					mockMemClobKeeper.On("GetOrderFillAmount", mock.Anything, orderId).Return(
						false,
						satypes.BaseQuantums(0),
						uint32(0),
					)
				}
			}

			// Create all unique orderbooks.
			createOrderbooks(
				t,
				ctx,
				memclob,
				3,
			)

			// Create all orders.
			applyOperationsToMemclob(
				t,
				ctx,
				memclob,
				tc.placedOperations,
				mockMemClobKeeper,
			)

			for _, operation := range tc.placedOperations {
				switch operation.Operation.(type) {
				case *types.Operation_ShortTermOrderPlacement:
					order := operation.GetShortTermOrderPlacement().Order
					orderId := order.OrderId
					// Mock out all remaining calls to GetOrderFillAmount, which is called in
					// `memclob.PurgeInvalidMemclobState` and during test assertions.
					fillAmount, exists := tc.newOrderFillAmounts[orderId]
					mockMemClobKeeper.On("GetOrderFillAmount", mock.Anything, orderId).Unset()
					mockMemClobKeeper.On("GetOrderFillAmount", mock.Anything, orderId).Return(exists, fillAmount, uint32(5))
				}
			}

			// Run the test.
			ctx = ctx.WithBlockHeight(10)
			offchainUpdates := types.NewOffchainUpdates()
			memclob.PurgeInvalidMemclobState(
				ctx,
				tc.fullyFilledOrderIds,
				tc.expiredStatefulOrderIds,
				tc.canceledStatefulOrderIds,
				tc.removedStatefulOrderIds,
				offchainUpdates,
			)

			// Verify that all removed orders have an associated off-chain removal.
			require.Equal(
				t,
				memclobtest.MessageCountOfType(offchainUpdates, types.RemoveMessageType),
				len(tc.expectedRemovedExpiredOrderIds),
			)

			for _, orderId := range tc.expectedRemovedExpiredOrderIds {
				require.True(
					t,
					memclobtest.HasMessage(offchainUpdates, orderId, types.RemoveMessageType),
				)
			}

			// Verify expecatations.
			AssertMemclobHasOrders(
				t,
				ctx,
				memclob,
				tc.expectedRemainingBids,
				tc.expectedRemainingAsks,
			)

			mockMemClobKeeper.AssertExpectations(t)

			//TODO(DEC-1936): Update these test assertions to verify nonces in `operationsToPropose`.
		})
	}
}

func TestPurgeInvalidMemclobState_DoesNotPanicWhenCalledWithDuplicateCanceledStatefulOrderIds(t *testing.T) {
	// Setup memclob state.
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	memclob := NewMemClobPriceTimePriority(true)
	mockMemClobKeeper := &mocks.MemClobKeeper{}
	memclob.SetClobKeeper(mockMemClobKeeper)
	memclob.CreateOrderbook(constants.ClobPair_Btc)
	mockMemClobKeeper.On("SendOrderbookUpdates", mock.Anything, mock.Anything).Return().Maybe()

	canceledStatefulOrderIds := []types.OrderId{
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
	}

	require.NotPanics(
		t,
		func() {
			memclob.PurgeInvalidMemclobState(
				ctx,
				[]types.OrderId{},
				[]types.OrderId{},
				canceledStatefulOrderIds,
				[]types.OrderId{},
				types.NewOffchainUpdates(),
			)
		},
	)
}

func TestPurgeInvalidMemclobState_PanicsWhenNonStatefulOrderIsCanceled(t *testing.T) {
	// Setup memclob state.
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	memclob := NewMemClobPriceTimePriority(true)
	mockMemClobKeeper := &mocks.MemClobKeeper{}
	memclob.SetClobKeeper(mockMemClobKeeper)
	memclob.CreateOrderbook(constants.ClobPair_Btc)
	mockMemClobKeeper.On("SendOrderbookUpdates", mock.Anything, mock.Anything).Return().Maybe()

	shortTermOrderId := constants.Order_Alice_Num0_Id0_Clob2_Buy5_Price10_GTB15.OrderId

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrderId,
		),
		func() {
			memclob.PurgeInvalidMemclobState(
				ctx,
				[]types.OrderId{},
				[]types.OrderId{},
				[]types.OrderId{shortTermOrderId},
				[]types.OrderId{},
				types.NewOffchainUpdates(),
			)
		},
	)
}

func TestPurgeInvalidMemclobState_DoesNotPanicWhenCalledWithDuplicateExpiredStatefulOrders(t *testing.T) {
	// Setup memclob state.
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)

	memclob := NewMemClobPriceTimePriority(true)
	mockMemClobKeeper := &mocks.MemClobKeeper{}
	memclob.SetClobKeeper(mockMemClobKeeper)
	memclob.CreateOrderbook(constants.ClobPair_Btc)
	mockMemClobKeeper.On("SendOrderbookUpdates", mock.Anything, mock.Anything).Return().Maybe()

	expiredStatefulOrderIds := []types.OrderId{
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
	}

	require.NotPanics(
		t,
		func() {
			memclob.PurgeInvalidMemclobState(
				ctx,
				[]types.OrderId{},
				expiredStatefulOrderIds,
				[]types.OrderId{},
				[]types.OrderId{},
				types.NewOffchainUpdates(),
			)
		},
	)
}

func TestPurgeInvalidMemclobState_PanicsWhenCalledWithShortTermExpiredStatefulOrders(t *testing.T) {
	// Setup memclob state.
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)

	memclob := NewMemClobPriceTimePriority(true)
	mockMemClobKeeper := &mocks.MemClobKeeper{}
	memclob.SetClobKeeper(mockMemClobKeeper)
	mockMemClobKeeper.On("SendOrderbookUpdates", mock.Anything, mock.Anything).Return().Maybe()

	shortTermOrderId := constants.Order_Alice_Num0_Id0_Clob2_Buy5_Price10_GTB15.OrderId

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrderId,
		),
		func() {
			memclob.PurgeInvalidMemclobState(
				ctx,
				[]types.OrderId{},
				[]types.OrderId{shortTermOrderId},
				[]types.OrderId{},
				[]types.OrderId{},
				types.NewOffchainUpdates(),
			)
		},
	)
}
