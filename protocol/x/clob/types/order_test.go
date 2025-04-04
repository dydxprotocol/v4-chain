package types_test

import (
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

func TestOrder_GetBaseQuantums(t *testing.T) {
	order := types.Order{
		Quantums: 100,
	}

	quantums := order.GetBaseQuantums()
	require.Equal(t, satypes.BaseQuantums(100), quantums)
}

func TestOrder_GetBigQuantums(t *testing.T) {
	tests := map[string]struct {
		// Order parameters.
		quantums uint64
		side     bool

		// Expectations.
		expectedBigQuantums *big.Int
	}{
		"Buy order": {
			quantums:            10,
			side:                true,
			expectedBigQuantums: big.NewInt(10),
		},
		"Sell order": {
			quantums:            10,
			side:                true,
			expectedBigQuantums: big.NewInt(10),
		},
		"Quantums of 0 returns 0 for buy order": {
			quantums:            0,
			side:                true,
			expectedBigQuantums: big.NewInt(0),
		},
		"Quantums of 0 returns 0 for sell order": {
			quantums:            0,
			side:                false,
			expectedBigQuantums: big.NewInt(0),
		},
		"Max Uint64 buy order": {
			quantums:            math.MaxUint64,
			side:                true,
			expectedBigQuantums: new(big.Int).SetUint64(math.MaxUint64),
		},
		"Max Uint64 sell order": {
			quantums:            math.MaxUint64,
			side:                false,
			expectedBigQuantums: constants.BigNegMaxUint64(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			side := types.Order_SIDE_BUY
			if !tc.side {
				side = types.Order_SIDE_SELL
			}
			order := types.Order{
				Side:     side,
				Quantums: tc.quantums,
			}

			require.Equal(t, tc.expectedBigQuantums, order.GetBigQuantums())
		})
	}
}

func TestOrder_GetOrderSubticks(t *testing.T) {
	order := types.Order{
		Subticks: uint64(100),
	}

	subticks := order.GetOrderSubticks()
	require.Equal(t, types.Subticks(100), subticks)
}

func TestOrder_MustBeValidOrderSide(t *testing.T) {
	tests := map[string]struct {
		side        types.Order_Side
		shouldPanic bool
	}{
		"Buy side doesn't panic": {
			side:        types.Order_SIDE_BUY,
			shouldPanic: false,
		},
		"Sell side doesn't panic": {
			side:        types.Order_SIDE_SELL,
			shouldPanic: false,
		},
		"Invalid side panics": {
			side:        types.Order_SIDE_UNSPECIFIED,
			shouldPanic: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			order := types.Order{
				Side: tc.side,
			}
			if tc.shouldPanic {
				// Expect a panic
				require.Panics(t, func() { order.MustBeValidOrderSide() })
			} else {
				require.NotPanics(t, func() { order.MustBeValidOrderSide() })
			}
		})
	}
}

func TestOrder_IsBuy(t *testing.T) {
	tests := map[string]struct {
		side     types.Order_Side
		expected bool
	}{
		"Is BUY": {
			side:     types.Order_SIDE_BUY,
			expected: true,
		},
		"Is SELL": {
			side:     types.Order_SIDE_SELL,
			expected: false,
		},
		"Is UNSPECIFIED": {
			side:     types.Order_SIDE_UNSPECIFIED,
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			order := types.Order{
				Side: tc.side,
			}

			result := order.IsBuy()
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestOrder_GetOrderHash(t *testing.T) {
	tests := map[string]struct {
		order        types.Order
		expectedHash types.OrderHash
	}{
		"Can take SHA256 hash of an empty order": {
			order:        types.Order{},
			expectedHash: constants.OrderHash_Empty,
		},
		"Can take SHA256 hash of a regular order": {
			order:        constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
			expectedHash: constants.OrderHash_Alice_Number0_Id0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expectedHash, tc.order.GetOrderHash())
		})
	}
}

func TestOrder_MustCmpReplacementOrder(t *testing.T) {
	tests := map[string]struct {
		x types.Order
		y types.Order
		r int
	}{
		"Short-Term GTB 1": {
			x: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 2}},
			y: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1}},
			r: 1,
		},
		"Short-Term GTB -1": {
			x: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1}},
			y: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 2}},
			r: -1,
		},
		"Short-Term Hash -1": {
			x: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1}, Subticks: 1},
			y: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1}, Subticks: 2},
			r: -1,
		},
		"Short-Term Hash 1": {
			x: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1}, Subticks: 2},
			y: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1}, Subticks: 1},
			r: 1,
		},
		"Short-Term Equal": {
			x: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 2}},
			y: types.Order{GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 2}},
			r: 0,
		},
		"Long-Term GTBT 1": {
			x: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
			},
			y: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			r: 1,
		},
		"Long-Term GTBT -1": {
			x: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			y: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
			},
			r: -1,
		},
		"Long-Term Hash -1": {
			x: types.Order{
				Subticks:     1,
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			y: types.Order{
				Subticks:     2,
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			r: -1,
		},
		"Long-Term Hash 1": {
			x: types.Order{
				Subticks:     2,
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			y: types.Order{
				Subticks:     1,
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			r: 1,
		},
		"Long-Term Equal": {
			x: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
			},
			y: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
			},
			r: 0,
		},
		"Conditional GTBT 1": {
			x: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_Conditional},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
			},
			y: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_Conditional},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			r: 1,
		},
		"Conditional GTBT -1": {
			x: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_Conditional},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			y: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_Conditional},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
			},
			r: -1,
		},
		"Conditional Hash -1": {
			x: types.Order{
				Subticks:     2,
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_Conditional},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			y: types.Order{
				Subticks:     1,
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_Conditional},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			r: -1,
		},
		"Conditional Hash 1": {
			x: types.Order{
				Subticks:     1,
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_Conditional},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			y: types.Order{
				Subticks:     2,
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_Conditional},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 1},
			},
			r: 1,
		},
		"Conditional Equal": {
			x: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
			},
			y: types.Order{
				OrderId:      types.OrderId{OrderFlags: types.OrderIdFlags_LongTerm},
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
			},
			r: 0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.r, tc.x.MustCmpReplacementOrder(&tc.y))
		})
	}
}

func TestOrder_MustCmpReplacementOrder_PanicsWithNonReplacementOrder(t *testing.T) {
	order1 := types.Order{OrderId: types.OrderId{}}
	order2 := types.Order{OrderId: types.OrderId{ClientId: 1}}
	require.NotEqual(t, order1.OrderId, order2.OrderId)

	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustCmpReplacementOrder: order ID (%v) does not equal order ID (%v)",
			order1.OrderId,
			order2.OrderId,
		),
		func() {
			order1.MustCmpReplacementOrder(&order2)
		},
	)
}

func TestOrder_GetSubaccountId(t *testing.T) {
	expectedSubaccountId := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15.OrderId.SubaccountId
	order := types.Order{
		OrderId: types.OrderId{
			SubaccountId: expectedSubaccountId,
		},
	}

	subaccountId := order.GetSubaccountId()
	require.Equal(t, expectedSubaccountId, subaccountId)
}

func TestOrder_IsLiquidation(t *testing.T) {
	order := types.Order{}

	isLiquidation := order.IsLiquidation()
	require.False(t, isLiquidation)
}

func TestOrder_MustGetOrder(t *testing.T) {
	expectedOrder := &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15

	order := expectedOrder.MustGetOrder()
	require.Equal(t, *expectedOrder, order)
}

func TestOrder_MustGetLiquidatedPerpetualIdPanics(t *testing.T) {
	expectedOrder := &constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15

	require.PanicsWithValue(
		t,
		"MustGetLiquidatedPerpetualId: No liquidated perpetual on an Order type.",
		func() {
			expectedOrder.MustGetLiquidatedPerpetualId()
		},
	)
}

func TestOrder_IsReduceOnly(t *testing.T) {
	require.False(t, constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.IsReduceOnly())
	require.True(t, constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO.IsReduceOnly())
}

func TestOrder_RequiresImmediateExecution(t *testing.T) {
	require.False(t, constants.Order_Alice_Num0_Id1_Clob0_Sell10_Price15_GTB15.RequiresImmediateExecution())
	require.True(t, constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC.RequiresImmediateExecution())
}

func TestOrder_IsShortTermOrder(t *testing.T) {
	for i := 0; i < numOrderIdFlagsTestCases; i++ {
		orderFlags := uint32(i)
		order := types.Order{
			OrderId: types.OrderId{
				OrderFlags: orderFlags,
			},
		}

		expectedIsShortTermOrder := orderFlags == types.OrderIdFlags_ShortTerm
		require.Equal(t, expectedIsShortTermOrder, order.IsShortTermOrder(), "OrderFlag: %d", i)
	}
}

func TestOrder_IsStatefulOrder(t *testing.T) {
	for i := 0; i < numOrderIdFlagsTestCases; i++ {
		orderFlags := uint32(i)
		order := types.Order{
			OrderId: types.OrderId{
				OrderFlags: orderFlags,
			},
		}

		expectedIsStatefulOrder := orderFlags == types.OrderIdFlags_LongTerm ||
			orderFlags == types.OrderIdFlags_Conditional ||
			orderFlags == types.OrderIdFlags_Twap ||
			orderFlags == types.OrderIdFlags_TwapSuborder
		require.Equal(t, expectedIsStatefulOrder, order.IsStatefulOrder(), "OrderFlag: %d", i)
	}
}

func TestOrder_MustBeStatefulOrder(t *testing.T) {
	tests := map[string]struct {
		orderFlags  uint32
		shouldPanic bool
	}{
		"Long-Term order doesn't panic": {
			orderFlags:  types.OrderIdFlags_LongTerm,
			shouldPanic: false,
		},
		"Conditional order doesn't panic": {
			orderFlags:  types.OrderIdFlags_Conditional,
			shouldPanic: false,
		},
		"Short-Term order panics": {
			orderFlags:  types.OrderIdFlags_ShortTerm,
			shouldPanic: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			order := types.Order{
				OrderId: types.OrderId{OrderFlags: tc.orderFlags},
			}
			if tc.shouldPanic {
				require.Panics(t, func() { order.MustBeStatefulOrder() })
			} else {
				require.NotPanics(t, func() { order.MustBeStatefulOrder() })
			}
		})
	}
}

func TestOrder_GetUnixGoodTilBlockTime(t *testing.T) {
	shortTermOrder := &types.Order{
		OrderId: types.OrderId{
			OrderFlags: types.OrderIdFlags_ShortTerm,
		},
	}
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			types.OrderId{
				OrderFlags: shortTermOrder.OrderId.OrderFlags,
			},
		),
		func() {
			shortTermOrder.MustGetUnixGoodTilBlockTime()
		},
	)

	invalidOrder := &types.Order{
		OrderId: types.OrderId{
			OrderFlags: types.OrderIdFlags_LongTerm,
		},
	}
	require.PanicsWithError(
		t,
		fmt.Errorf(
			"MustGetUnixGoodTilBlockTime: order (%v) goodTilBlockTime is zero",
			invalidOrder,
		).Error(),
		func() {
			invalidOrder.MustGetUnixGoodTilBlockTime()
		},
	)

	require.Equal(
		t,
		time.Unix(15, 0),
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.MustGetUnixGoodTilBlockTime(),
	)

	require.Equal(
		t,
		time.Unix(15, 0),
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.MustGetUnixGoodTilBlockTime(),
	)
}

func TestOrder_MustGetOrderJson(t *testing.T) {
	expectedOrder := &constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15

	require.Equal(t,
		"order_id: <\n  subaccount_id: <\n    "+
			"owner: \"dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4\"\n  >\n  order_flags: 64\n>"+
			"\nside: SIDE_BUY\nquantums: 5\nsubticks: 10\ngood_til_block_time: 15\n",
		expectedOrder.GetOrderTextString(),
	)
}

func TestOrder_GetClientMetadata(t *testing.T) {
	order := types.Order{
		ClientMetadata: uint32(100),
	}

	clientMetadata := order.GetClientMetadata()
	require.Equal(t, uint32(100), clientMetadata)
}
