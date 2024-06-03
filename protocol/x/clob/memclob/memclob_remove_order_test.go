package memclob

import (
	"math"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil_memclob "github.com/dydxprotocol/v4-chain/protocol/testutil/memclob"
	sdktest "github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRemoveOrder_PanicsIfNotExists(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	memclob := NewMemClobPriceTimePriority(false)

	order1 := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15

	require.Panics(t, func() {
		memclob.mustRemoveOrder(ctx, order1.OrderId)
	})
}

func TestRemoveOrderIfFilled(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		existingOrders []types.Order

		// Parameters.
		order               types.Order
		orderFilledQuantums satypes.BaseQuantums
		isReplacement       bool

		// Expectations.
		expectedBestBid                              types.Subticks
		expectedBestAsk                              types.Subticks
		expectedTotalLevels                          int
		expectedTotalLevelQuantums                   uint64
		expectLevelToExist                           bool
		expectBlockExpirationsForOrdersToExist       bool
		expectSubaccountOpenClobOrdersForSideToExist bool
		expectSubaccountOpenClobOrdersToExist        bool
		expectOrderToExist                           bool
	}{
		"Removes bid from otherwise empty book": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			order:               constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.GetBaseQuantums(),

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Removes bid from otherwise empty book if overfilled": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			order:               constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.GetBaseQuantums() * 2,

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Does not remove bid if not fully filled": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			order:               constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.GetBaseQuantums() - 1,

			expectedTotalLevels:                          1,
			expectedBestBid:                              10,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
			expectOrderToExist:                           true,
		},
		"Removes ask from otherwise empty book": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},
			order:               constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.GetBaseQuantums(),

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Removes ask from otherwise empty book if overfilled": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},
			order:               constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.GetBaseQuantums() * 2,

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Does not remove ask if not fully filled": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},
			order:               constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.GetBaseQuantums() - 1,

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              15,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
			expectOrderToExist:                           true,
		},
		"Removes higher time-priority bid from book at same level": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
			},
			order:               constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.GetBaseQuantums(),

			expectedTotalLevels:                          1,
			expectedBestBid:                              10,
			expectedBestAsk:                              math.MaxUint64,
			expectedTotalLevelQuantums:                   20,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Removes replaced order that has been fully-filled at the new size": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num1_Id10_Clob0_Buy15_Price30_GTB33,
				constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB34,
			},
			order:               constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB34,
			orderFilledQuantums: constants.Order_Alice_Num1_Id10_Clob0_Buy10_Price30_GTB34.GetBaseQuantums(),
			isReplacement:       true,

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Removes higher time-priority ask from book at same level": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
			},
			order:               constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
			orderFilledQuantums: constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20.GetBaseQuantums(),

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              15,
			expectedTotalLevelQuantums:                   10,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Removes lower time-priority bid from book at same level": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
			},
			order:               constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
			orderFilledQuantums: constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22.GetBaseQuantums(),

			expectedTotalLevels:                          1,
			expectedBestBid:                              10,
			expectedBestAsk:                              math.MaxUint64,
			expectedTotalLevelQuantums:                   5,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Removes lower time-priority ask from book at same level": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
			},
			order:               constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
			orderFilledQuantums: constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20.GetBaseQuantums(),

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              15,
			expectedTotalLevelQuantums:                   25,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Removes best price level bid": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
			},
			order:               constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
			orderFilledQuantums: constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20.GetBaseQuantums(),

			expectedTotalLevels:                          1,
			expectedBestBid:                              10,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
			expectOrderToExist:                           false,
		},
		"Removes best price level ask": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
			},
			order:               constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
			orderFilledQuantums: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22.GetBaseQuantums(),

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              15,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
			expectOrderToExist:                           false,
		},
		"Removes worst price level bid": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
			},
			order:               constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.GetBaseQuantums(),

			expectedTotalLevels:                          1,
			expectedBestBid:                              5,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
			expectOrderToExist:                           false,
		},
		"Removes worst price level ask": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
			},
			order:               constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
			orderFilledQuantums: constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20.GetBaseQuantums(),

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              10,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Removes subaccountOpenClobOrders for one side": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			order:               constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			orderFilledQuantums: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.GetBaseQuantums(),

			expectedTotalLevels:                          0,
			expectedBestBid:                              10,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        true,
			expectOrderToExist:                           false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the memclob state.
			memClobKeeper := mocks.MemClobKeeper{}
			memclob := NewMemClobPriceTimePriority(false)
			memclob.SetClobKeeper(&memClobKeeper)

			memClobKeeper.On("ValidateSubaccountEquityTierLimitForNewOrder", mock.Anything, mock.Anything).Return(nil)
			memClobKeeper.On("SendOrderbookUpdates", mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

			// Set initial fill amount to `0` for all orders.
			initialCall := memClobKeeper.On("GetOrderFillAmount", mock.Anything, mock.Anything).
				Return(true, satypes.BaseQuantums(0), uint32(0))

			// Create all unique orderbooks.
			createAllOrderbooksForOrders(
				t,
				ctx,
				memclob,
				append(tc.existingOrders, tc.order),
			)

			// Place all existing orders on the orderbook
			for _, order := range tc.existingOrders {
				_, _, _, err := memclob.PlaceOrder(ctx, order)
				require.NoError(t, err)
			}

			// Unset the initial mock call to overwrite it with `orderFilledQuantums` for the provided `order`.
			initialCall.Unset()

			// Set the order to be filled based on `tc.orderFilledQuantums`.
			memClobKeeper.On("GetOrderFillAmount", mock.Anything, tc.order.OrderId).
				Return(true, satypes.BaseQuantums(tc.orderFilledQuantums), uint32(0))

			// Run the test case.
			memclob.RemoveOrderIfFilled(ctx, tc.order.OrderId)
			for _, existingOrder := range tc.existingOrders {
				if !tc.isReplacement && (existingOrder != tc.order || tc.expectOrderToExist) {
					requireOrderExistsInMemclob(t, ctx, existingOrder, memclob)
				} else {
					requireOrderDoesNotExistInMemclob(t, ctx, tc.order, memclob)
				}
			}

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

func TestRemoveOrder(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()
	ctx = ctx.WithIsCheckTx(true)
	tests := map[string]struct {
		// State.
		existingOrders []types.Order

		// Parameters.
		order types.Order

		// Expectations.
		expectedBestBid                              types.Subticks
		expectedBestAsk                              types.Subticks
		expectedTotalLevels                          int
		expectedTotalLevelQuantums                   uint64
		expectLevelToExist                           bool
		expectBlockExpirationsForOrdersToExist       bool
		expectSubaccountOpenClobOrdersForSideToExist bool
		expectSubaccountOpenClobOrdersToExist        bool
	}{
		"Removes bid from otherwise empty book": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},
			order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Removes ask from otherwise empty book": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			},
			order: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,

			expectedTotalLevels:                          0,
			expectedBestBid:                              0,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Removes higher time-priority bid from book at same level": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
			},
			order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,

			expectedTotalLevels:                          1,
			expectedBestBid:                              10,
			expectedBestAsk:                              math.MaxUint64,
			expectedTotalLevelQuantums:                   20,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Removes higher time-priority ask from book at same level": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
			},
			order: constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              15,
			expectedTotalLevelQuantums:                   10,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Removes lower time-priority bid from book at same level": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,
			},
			order: constants.Order_Bob_Num0_Id5_Clob0_Buy20_Price10_GTB22,

			expectedTotalLevels:                          1,
			expectedBestBid:                              10,
			expectedBestAsk:                              math.MaxUint64,
			expectedTotalLevelQuantums:                   5,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Removes lower time-priority ask from book at same level": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,
			},
			order: constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20,

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              15,
			expectedTotalLevelQuantums:                   25,
			expectLevelToExist:                           true,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Removes best price level bid": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
			},
			order: constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,

			expectedTotalLevels:                          1,
			expectedBestBid:                              10,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
		},
		"Removes best price level ask": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
			},
			order: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              15,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Removes worst price level bid": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Alice_Num0_Id6_Clob0_Buy25_Price5_GTB20,
			},
			order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,

			expectedTotalLevels:                          1,
			expectedBestBid:                              5,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
		},
		"Removes worst price level ask": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,
				constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
			},
			order: constants.Order_Alice_Num0_Id7_Clob0_Sell25_Price15_GTB20,

			expectedTotalLevels:                          1,
			expectedBestBid:                              0,
			expectedBestAsk:                              10,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       false,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        false,
		},
		"Resorts to iterating over all bid levels after `totalLevels` of iteration": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
				constants.Order_Bob_Num0_Id6_Clob0_Buy20_Price1000_GTB22,
				constants.Order_Bob_Num0_Id7_Clob0_Buy20_Price10000_GTB22,
			},
			order: constants.Order_Bob_Num0_Id7_Clob0_Buy20_Price10000_GTB22,

			expectedTotalLevels:                          2,
			expectedBestBid:                              1000,
			expectedBestAsk:                              math.MaxUint64,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
		},
		"Resorts to iterating over all ask levels after `totalLevels` of iteration": {
			existingOrders: []types.Order{
				constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,
				constants.Order_Bob_Num0_Id9_Clob0_Sell20_Price1000,
				constants.Order_Bob_Num0_Id10_Clob0_Sell20_Price10000,
			},
			order: constants.Order_Bob_Num0_Id8_Clob0_Sell20_Price10_GTB22,

			expectedTotalLevels:                          2,
			expectedBestBid:                              0,
			expectedBestAsk:                              1000,
			expectLevelToExist:                           false,
			expectBlockExpirationsForOrdersToExist:       true,
			expectSubaccountOpenClobOrdersForSideToExist: true,
			expectSubaccountOpenClobOrdersToExist:        true,
		},
		"Removes subaccountOpenClobOrders for one side": {
			existingOrders: []types.Order{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			},

			order:                                  constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15,
			expectedTotalLevels:                    0,
			expectedBestBid:                        10,
			expectedBestAsk:                        math.MaxUint64,
			expectLevelToExist:                     false,
			expectBlockExpirationsForOrdersToExist: true,
			expectSubaccountOpenClobOrdersForSideToExist: false,
			expectSubaccountOpenClobOrdersToExist:        true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the memclob state.
			memClobKeeper := testutil_memclob.NewFakeMemClobKeeper()
			memclob := NewMemClobPriceTimePriority(false)
			memclob.SetClobKeeper(memClobKeeper)

			// Create all unique orderbooks.
			createAllOrderbooksForOrders(
				t,
				ctx,
				memclob,
				append(tc.existingOrders, tc.order),
			)

			// Place all existing orders on the orderbook
			for _, order := range tc.existingOrders {
				_, _, _, err := memclob.PlaceOrder(ctx, order)
				require.NoError(t, err)
			}

			// Run the test case.
			memclob.mustRemoveOrder(ctx, tc.order.OrderId)
			requireOrderDoesNotExistInMemclob(t, ctx, tc.order, memclob)
			for _, existingOrder := range tc.existingOrders {
				if existingOrder != tc.order {
					requireOrderExistsInMemclob(t, ctx, existingOrder, memclob)
				}
			}

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
