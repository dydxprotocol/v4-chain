package types_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

// numOrderIdFlagsTestCases is set to 513 to verify that we run a test case where
// `OrderFlags` is greater than one byte (proto varints are encoded with 7 bits per byte).
const numOrderIdFlagsTestCases = 513

func TestToStateKey(t *testing.T) {
	// Success
	b, _ := constants.OrderId_Alice_Num0_ClientId0_Clob0.Marshal()
	require.Equal(t, b, constants.OrderId_Alice_Num0_ClientId0_Clob0.ToStateKey())

	// No panic case. MustMarshal() > Marshal() > MarshalToSizedBuffer() which never returns an error.
}

func TestIsShortTermOrder(t *testing.T) {
	for i := 0; i < numOrderIdFlagsTestCases; i++ {
		orderFlags := uint32(i)
		orderId := types.OrderId{
			OrderFlags: orderFlags,
		}

		expectedIsShortTermOrder := orderFlags == types.OrderIdFlags_ShortTerm
		require.Equal(t, expectedIsShortTermOrder, orderId.IsShortTermOrder(), "OrderFlag: %d", i)
	}
}

func TestIsLongTermOrder(t *testing.T) {
	for i := 0; i < numOrderIdFlagsTestCases; i++ {
		orderFlags := uint32(i)
		orderId := types.OrderId{
			OrderFlags: orderFlags,
		}

		expectedIsLongTermOrder := orderFlags == types.OrderIdFlags_LongTerm
		require.Equal(t, expectedIsLongTermOrder, orderId.IsLongTermOrder(), "OrderFlag: %d", i)
	}
}

func TestIsConditionalOrder(t *testing.T) {
	for i := 0; i < numOrderIdFlagsTestCases; i++ {
		orderFlags := uint32(i)
		orderId := types.OrderId{
			OrderFlags: orderFlags,
		}

		expectedIsConditionalOrder := orderFlags == types.OrderIdFlags_Conditional
		require.Equal(t, expectedIsConditionalOrder, orderId.IsConditionalOrder(), "OrderFlag: %d", i)
	}
}

func TestIsStatefulOrder(t *testing.T) {
	for i := 0; i < numOrderIdFlagsTestCases; i++ {
		orderFlags := uint32(i)
		orderId := types.OrderId{OrderFlags: orderFlags}

		expectedIsStatefulOrder := orderFlags == types.OrderIdFlags_LongTerm ||
			orderFlags == types.OrderIdFlags_Conditional ||
			orderFlags == types.OrderIdFlags_Twap ||
			orderFlags == types.OrderIdFlags_TwapSuborder
		require.Equal(t, expectedIsStatefulOrder, orderId.IsStatefulOrder(), "OrderFlag: %d", i)
	}
}

func TestSortOrders(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		orders         []types.OrderId
		expectedOrders []types.OrderId
	}{
		"sorts with different subaccount owners": {
			orders: []types.OrderId{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
				constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId,
			},
			expectedOrders: []types.OrderId{
				constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			},
		},
		"sorts with same subaccount owner different subaccount number": {
			orders: []types.OrderId{
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			},
			expectedOrders: []types.OrderId{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
				constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20.OrderId,
			},
		},
		"sorts with same subaccount owner and number, different client id": {
			orders: []types.OrderId{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			},
			expectedOrders: []types.OrderId{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
			},
		},
		"sorts with same subaccount owner, number and client id, different order flags": {
			orders: []types.OrderId{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			},
			expectedOrders: []types.OrderId{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			},
		},
		"sorts with same subaccount owner, number, client id, and order flags, different clob pair id": {
			orders: []types.OrderId{
				constants.Order_Alice_Num0_Id0_Clob2_Buy5_Price10_GTB15.OrderId,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
				constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15.OrderId,
			},
			expectedOrders: []types.OrderId{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
				constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15.OrderId,
				constants.Order_Alice_Num0_Id0_Clob2_Buy5_Price10_GTB15.OrderId,
			},
		},
		"sorts with same order": {
			orders: []types.OrderId{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
			},
			expectedOrders: []types.OrderId{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
			},
		},
		"sorts with one order": {
			orders: []types.OrderId{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
			},
			expectedOrders: []types.OrderId{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			orders := tc.orders
			sort.Sort(types.SortedOrders(orders))
			require.Equal(t, orders, tc.expectedOrders)
		})
	}
}

func TestMustSortAndHaveNoDuplicates_Panic(t *testing.T) {
	// duplicates should panic
	require.PanicsWithError(
		t,
		fmt.Sprintf(
			"cannot sort orders with duplicate order id %+v",
			constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
		),
		func() {
			orders := []types.OrderId{
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
				constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,
			}
			types.MustSortAndHaveNoDuplicates(orders)
		},
	)
}

func TestMustSortAndHaveNoDuplicates_Success(t *testing.T) {
	// without duplicates, should work fine
	orders := []types.OrderId{
		constants.Order_Alice_Num0_Id0_Clob2_Buy5_Price10_GTB15.OrderId,
		constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
		constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15.OrderId,
	}
	expectedOrders := []types.OrderId{
		constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
		constants.Order_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB15.OrderId,
		constants.Order_Alice_Num0_Id0_Clob2_Buy5_Price10_GTB15.OrderId,
	}
	types.MustSortAndHaveNoDuplicates(orders)
	require.Equal(t, orders, expectedOrders)
}

func TestOrderIdIsEqual(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		o1 types.OrderId
		o2 types.OrderId

		// Expectations.
		isEqual bool
	}{
		"Two order IDs are equal if they have the same owner, subaccount number, client ID, and order flags": {
			o1: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			o2: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,

			isEqual: true,
		},
		"Two order IDs are not equal if they don't have the same owner": {
			o1: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			o2: constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.OrderId,

			isEqual: false,
		},
		"Two order IDs are not equal if they don't have the same subaccount number": {
			o1: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			o2: constants.Order_Alice_Num1_Id0_Clob0_Sell10_Price15_GTB20.OrderId,

			isEqual: false,
		},
		"Two order IDs are not equal if they don't have the same client ID": {
			o1: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,
			o2: constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.OrderId,

			isEqual: false,
		},
		"Two order IDs are not equal if they don't have the same order flags": {
			o1: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
			o2: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId,

			isEqual: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.isEqual {
				require.Equal(t, tc.o1, tc.o2)
			} else {
				require.NotEqual(t, tc.o1, tc.o2)
			}
		})
	}
}

func TestOrderIdEquality(t *testing.T) {
	orderIds := make(map[types.OrderId]bool)
	o1 := types.OrderId{
		SubaccountId: satypes.SubaccountId{
			Owner:  "5",
			Number: 2,
		},
		ClientId:   2,
		OrderFlags: types.OrderIdFlags_LongTerm,
	}
	orderIds[o1] = true

	o2 := types.OrderId{
		SubaccountId: satypes.SubaccountId{
			Owner:  "5",
			Number: 2,
		},
		ClientId:   2,
		OrderFlags: types.OrderIdFlags_LongTerm,
	}

	require.Equal(t, o1, o2)

	_, exists := orderIds[o2]
	require.True(t, exists)
}
