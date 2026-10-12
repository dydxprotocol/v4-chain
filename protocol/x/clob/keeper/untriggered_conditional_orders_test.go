package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddUntriggeredConditionalOrder(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		conditionalOrdersToAdd []types.Order

		// Expectations.
		expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice []types.Order
		expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice []types.Order
		expectedNumberOfMatches                               uint32
	}{
		"Can add a stop loss buy to the GTE array": {
			conditionalOrdersToAdd: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
			},

			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
			},
			expectedNumberOfMatches: 1,
		},
		"Can add a take profit sell to the GTE array": {
			conditionalOrdersToAdd: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_TakeProfit20,
			},

			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_TakeProfit20,
			},
			expectedNumberOfMatches: 1,
		},
		"Can add a take profit buy to the LTE array": {
			conditionalOrdersToAdd: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
			},

			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
			},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
			expectedNumberOfMatches:                               1,
		},
		"Can add a stop loss sell to the LTE array": {
			conditionalOrdersToAdd: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss20,
			},

			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss20,
			},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
			expectedNumberOfMatches:                               1,
		},
		"Can add multiple conditional orders to both heaps": {
			conditionalOrdersToAdd: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
			},

			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
			},
			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
			},
			expectedNumberOfMatches: 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			untriggeredConditionalOrders := keeper.NewUntriggeredConditionalOrders()

			for _, order := range tc.conditionalOrdersToAdd {
				untriggeredConditionalOrders.AddUntriggeredConditionalOrder(order)
			}

			require.Equal(
				t,
				tc.expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice,
				untriggeredConditionalOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice,
			)
			require.Equal(
				t,
				tc.expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice,
				untriggeredConditionalOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice,
			)
		})
	}
}

func TestAddUntriggeredConditionalOrder_NonConditionalOrder(t *testing.T) {
	untriggeredConditionalOrders := keeper.NewUntriggeredConditionalOrders()
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeConditionalOrder: called with non-conditional order ID (%+v)",
			&constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
		),
		func() {
			untriggeredConditionalOrders.AddUntriggeredConditionalOrder(
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
			)
		},
	)
}

func TestOrganizeUntriggeredConditionalOrdersFromState(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		conditionalOrdersFromState []types.Order

		// Expectations.
		expectedUntriggeredConditionalOrders map[types.ClobPairId]*keeper.UntriggeredConditionalOrders
	}{
		"Only GTE orders, one ClobPair": {
			conditionalOrdersFromState: []types.Order{
				// GTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20,
			},
			expectedUntriggeredConditionalOrders: map[types.ClobPairId]*keeper.UntriggeredConditionalOrders{
				0: {
					OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{},
					OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
						constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
						constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20,
					},
				},
			},
		},
		"Both GTE and LTE orders, one ClobPair": {
			conditionalOrdersFromState: []types.Order{
				// GTE
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				// LTE
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20,
			},
			expectedUntriggeredConditionalOrders: map[types.ClobPairId]*keeper.UntriggeredConditionalOrders{
				0: {
					OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20,
					},
					OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					},
				},
			},
		},
		"Multiple ClobPair + both LTE and GTE orders": {
			conditionalOrdersFromState: []types.Order{
				// GTE, ClobPair 1
				constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit20,
				// GTE, ClobPair 0
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				// LTE, ClobPair 0
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20,
				// GTE, ClobPair 0
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				// LTE, ClobPair 1
				constants.ConditionalOrder_Alice_Num0_Id1_Clob1_Buy5_Price10_GTBT15_StopLoss20,
			},
			expectedUntriggeredConditionalOrders: map[types.ClobPairId]*keeper.UntriggeredConditionalOrders{
				0: {
					OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20,
					},
					OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
						constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
					},
				},
				1: {
					OrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT15_TakeProfit20,
					},
					OrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
						constants.ConditionalOrder_Alice_Num0_Id1_Clob1_Buy5_Price10_GTBT15_StopLoss20,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := keeper.OrganizeUntriggeredConditionalOrdersFromState(tc.conditionalOrdersFromState)

			require.Equal(
				t,
				tc.expectedUntriggeredConditionalOrders,
				got,
			)
		})
	}
}

func TestRemoveUntriggeredConditionalOrders(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		conditionalOrdersToAdd      []types.Order
		conditionalOrderIdsToExpire []types.OrderId

		// Expectations.
		expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice []types.Order
		expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice []types.Order
	}{
		"Removes multiple expired order from GTE array": {
			conditionalOrdersToAdd: []types.Order{
				// GTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_StopLoss20,
			},
			conditionalOrderIdsToExpire: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_StopLoss20.OrderId,
			},

			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_StopLoss20,
			},
		},
		"Removes multiple expired order from LTE array": {
			conditionalOrdersToAdd: []types.Order{
				// LTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price10_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20,
			},
			conditionalOrderIdsToExpire: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price10_GTBT15_StopLoss20.OrderId,
			},

			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20,
			},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
		},
		"Full clear of both GTE and LTE orders": {
			conditionalOrdersToAdd: []types.Order{
				// GTE
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
				// LTE
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20,
			},
			conditionalOrderIdsToExpire: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price10_GTBT15_TakeProfit20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
			},

			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			untriggeredConditionalOrders := keeper.NewUntriggeredConditionalOrders()

			for _, order := range tc.conditionalOrdersToAdd {
				untriggeredConditionalOrders.AddUntriggeredConditionalOrder(order)
			}

			untriggeredConditionalOrders.RemoveUntriggeredConditionalOrders(tc.conditionalOrderIdsToExpire)

			require.Equal(
				t,
				tc.expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice,
				untriggeredConditionalOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice,
			)
			require.Equal(
				t,
				tc.expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice,
				untriggeredConditionalOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice,
			)
		})
	}
}

func TestPollTriggeredConditionalOrders(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		conditionalOrdersToAdd []types.Order
		clobPairId             types.ClobPairId
		currentSubticks        *big.Rat

		// Expectations.
		expectedTriggeredOrderIds                             []types.OrderId
		expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice []types.Order
		expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice []types.Order
	}{
		"No conditional orders triggered": {
			conditionalOrdersToAdd: []types.Order{
				// GTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,

				// LTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
			},
			currentSubticks:           big.NewRat(15, 1),
			expectedTriggeredOrderIds: []types.OrderId{},
			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
			},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,
			},
		},
		"Trigger GTE orders": {
			conditionalOrdersToAdd: []types.Order{
				// GTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,

				// LTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
			},
			currentSubticks: big.NewRat(20, 1),

			expectedTriggeredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20.OrderId,
			},
			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
			},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,
			},
		},
		"Trigger LTE orders": {
			conditionalOrdersToAdd: []types.Order{
				// GTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,

				// LTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
			},
			currentSubticks: big.NewRat(10, 1),

			expectedTriggeredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10.OrderId,
			},
			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
			},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,
			},
		},
		"Trigger all LTE orders": {
			conditionalOrdersToAdd: []types.Order{
				// GTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,

				// LTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
			},
			currentSubticks: big.NewRat(0, 1),

			expectedTriggeredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10.OrderId,
			},
			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,
			},
		},
		"Trigger all GTE orders": {
			conditionalOrdersToAdd: []types.Order{
				// GTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25,

				// LTE orders
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
			},
			currentSubticks: big.NewRat(50, 1),

			expectedTriggeredOrderIds: []types.OrderId{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price25_GTBT15_StopLoss25.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Sell20_Price20_GTBT15_TakeProfit20.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Buy25_Price25_GTBT15_StopLoss25.OrderId,
			},
			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Buy15_Price10_GTBT15_TakeProfit5,
				constants.ConditionalOrder_Alice_Num0_Id2_Clob0_Buy20_Price10_GTBT15_TakeProfit10,
				constants.ConditionalOrder_Alice_Num0_Id3_Clob0_Sell25_Price10_GTBT15_StopLoss10,
			},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
		},
		"Pessimistically rounds and doesn't trigger GTE as a result": {
			conditionalOrdersToAdd: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
			},
			currentSubticks:           big.NewRat(39, 2), // 19.5 will round down to 19 and not trigger
			expectedTriggeredOrderIds: []types.OrderId{},
			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price20_GTBT15_StopLoss20,
			},
		},
		"Pessimistically rounds and doesn't trigger LTE as a result": {
			conditionalOrdersToAdd: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
			},
			currentSubticks:           big.NewRat(21, 2), // 10.5 will round up to 11 and not trigger
			expectedTriggeredOrderIds: []types.OrderId{},
			expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice: []types.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_TakeProfit10,
			},
			expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice: []types.Order{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			untriggeredConditionalOrders := keeper.NewUntriggeredConditionalOrders()

			for _, order := range tc.conditionalOrdersToAdd {
				untriggeredConditionalOrders.AddUntriggeredConditionalOrder(order)
			}

			triggeredOrderIds := untriggeredConditionalOrders.PollTriggeredConditionalOrders(
				tc.currentSubticks,
			)

			require.Equal(
				t,
				tc.expectedTriggeredOrderIds,
				triggeredOrderIds,
			)

			require.Equal(
				t,
				tc.expectedOrdersToTriggerWhenOraclePriceGTETriggerPrice,
				untriggeredConditionalOrders.OrdersToTriggerWhenOraclePriceGTETriggerPrice,
			)
			require.Equal(
				t,
				tc.expectedOrdersToTriggerWhenOraclePriceLTETriggerPrice,
				untriggeredConditionalOrders.OrdersToTriggerWhenOraclePriceLTETriggerPrice,
			)
		})
	}
}

// A market the oracle hasn't priced yet must not panic EndBlocker; its conditional orders stay untriggered
// until it has a real price.
func TestMaybeTriggerConditionalOrders_SkipsZeroPriceMarket(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	mockIndexerEventManager.On("Enabled").Return(false)
	mockIndexerEventManager.On(
		"AddTxnEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

	keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

	// A freshly listed market with no oracle price yet.
	marketParam := pricestypes.MarketParam{
		Id:                 10,
		Pair:               "ZEROPRICE-USD",
		Exponent:           -6,
		MinExchanges:       1,
		MinPriceChangePpm:  50,
		ExchangeConfigJson: `{"test_config_placeholder":{}}`,
	}
	_, err := keepertest.CreateTestMarket(
		t, ks.Ctx, ks.PricesKeeper, marketParam,
		pricestypes.MarketPrice{Id: marketParam.Id, Exponent: marketParam.Exponent, Price: 0},
	)
	require.NoError(t, err)

	const perpetualId = uint32(10)
	perpetual, err := ks.PerpetualsKeeper.CreatePerpetual(
		ks.Ctx, perpetualId, "ZEROPRICE-USD", marketParam.Id, int32(-8), int32(0),
		uint32(0), perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
	)
	require.NoError(t, err)

	clobPair, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
		ks.Ctx, 10, perpetual.Params.Id, satypes.BaseQuantums(10), -8, 1, types.ClobPair_STATUS_ACTIVE,
	)
	require.NoError(t, err)

	order := types.Order{
		OrderId: types.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   types.OrderIdFlags_Conditional,
			ClobPairId:   clobPair.Id,
		},
		Side:                            types.Order_SIDE_BUY,
		Quantums:                        5,
		Subticks:                        10,
		GoodTilOneof:                    &types.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
		ConditionType:                   types.Order_CONDITION_TYPE_STOP_LOSS,
		ConditionalOrderTriggerSubticks: 1,
	}
	// A stop-loss sell would trigger against a zero price.
	sellOrder := order
	sellOrder.OrderId.ClientId = 1
	sellOrder.Side = types.Order_SIDE_SELL
	for i, o := range []types.Order{order, sellOrder} {
		ks.ClobKeeper.GetUntriggeredConditionalOrderPlacementStore(ks.Ctx).Set(
			o.OrderId.ToStateKey(),
			ks.Cdc.MustMarshal(&types.LongTermOrderPlacement{
				Order:          o,
				PlacementIndex: types.TransactionOrdering{TransactionIndex: uint32(i)},
			}),
		)
	}

	require.NotPanics(t, func() {
		triggered := ks.ClobKeeper.MaybeTriggerConditionalOrders(ks.Ctx)
		require.Empty(t, triggered, "zero-price market's conditional orders must not be triggered")
	})
	for _, o := range []types.Order{order, sellOrder} {
		_, found := ks.ClobKeeper.GetUntriggeredConditionalOrderPlacement(ks.Ctx, o.OrderId)
		require.True(t, found, "order must remain untriggered, not dropped")
	}

	// Once the oracle catches up with a real price, the same order triggers normally.
	require.NoError(t, ks.PricesKeeper.UpdateMarketPrices(ks.Ctx, []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
		{MarketId: marketParam.Id, Price: 1_000_000},
	}))
	triggered := ks.ClobKeeper.MaybeTriggerConditionalOrders(ks.Ctx)
	require.Equal(t, []types.OrderId{order.OrderId}, triggered)
}
