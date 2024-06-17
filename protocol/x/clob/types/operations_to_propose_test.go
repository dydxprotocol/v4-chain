package types_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestClearOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		// State.
		operationsQueue              []types.InternalOperation
		orderHashesInOperationsQueue map[types.OrderHash]bool
		shortTermOrderHashToTxBytes  map[types.OrderHash][]byte
	}{
		"Empty operations queue": {
			operationsQueue:              []types.InternalOperation{},
			orderHashesInOperationsQueue: map[types.OrderHash]bool{},
			shortTermOrderHashToTxBytes:  map[types.OrderHash][]byte{},
		},
		"Operations queue and order hashes in operations queue gets cleared": {
			operationsQueue: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20,
					[]types.MakerFill{
						{
							FillAmount:   10,
							MakerOrderId: constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.GetOrderId(),
						},
					},
				),
			},
			orderHashesInOperationsQueue: map[types.OrderHash]bool{
				constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20.GetOrderHash(): true,
				constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.GetOrderHash():   true,
			},
			shortTermOrderHashToTxBytes: map[types.OrderHash][]byte{
				constants.Order_Alice_Num0_Id5_Clob1_Sell25_Price15_GTB20.GetOrderHash(): {4, 0, 8},
				constants.Order_Bob_Num0_Id0_Clob1_Sell10_Price15_GTB20.GetOrderHash():   {8, 8, 8},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			otp.OperationsQueue = tc.operationsQueue
			otp.OrderHashesInOperationsQueue = tc.orderHashesInOperationsQueue
			otp.ShortTermOrderHashToTxBytes = tc.shortTermOrderHashToTxBytes

			otp.ClearOperationsQueue()
			// Verify expectations of state.
			require.Empty(t, otp.OperationsQueue)
			require.Empty(t, otp.OrderHashesInOperationsQueue)
			require.Equal(t, tc.shortTermOrderHashToTxBytes, otp.ShortTermOrderHashToTxBytes)
		})
	}
}

func TestMustAddShortTermOrderTxBytes(t *testing.T) {
	shortTermOrder1 := constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC
	shortTermOrder2 := constants.Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000
	shortTermOrder3 := constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000
	shortTermOrder4 := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	shortTermOrder5 := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16

	otp := types.NewOperationsToPropose()

	// Add `shortTermOrder1`, `shortTermOrder2`, and `shortTermOrder4` to `ShortTermOrderHashToTxBytes`
	// and verify they're present in `ShortTermOrderHashToTxBytes` but not in the operations queue.
	for _, order := range []types.Order{
		shortTermOrder1,
		shortTermOrder2,
		shortTermOrder4,
	} {
		orderHash := order.GetOrderHash()
		otp.MustAddShortTermOrderTxBytes(order, orderHash.ToBytes())
		_, exists := otp.ShortTermOrderHashToTxBytes[orderHash]
		require.True(t, exists)
		require.False(t, otp.IsOrderPlacementInOperationsQueue(order))
	}

	// Verify `shortTermOrder3` and `shortTermOrder5` are not present in `ShortTermOrderHashToTxBytes`
	// and not present in the operations queue.
	for _, orderNotInOpQueue := range []types.Order{
		shortTermOrder3,
		shortTermOrder5,
	} {
		orderHash := orderNotInOpQueue.GetOrderHash()
		_, exists := otp.ShortTermOrderHashToTxBytes[orderHash]
		require.False(t, exists)
		require.False(t, otp.IsOrderPlacementInOperationsQueue(orderNotInOpQueue))
	}

	// Add `shortTermOrder3` and `shortTermOrder5`to the queue and
	// verify they're present in `ShortTermOrderHashToTxBytes` but not in the operations queue.
	for _, order := range []types.Order{
		shortTermOrder3,
		shortTermOrder5,
	} {
		orderHash := order.GetOrderHash()
		otp.MustAddShortTermOrderTxBytes(order, orderHash.ToBytes())
		_, exists := otp.ShortTermOrderHashToTxBytes[orderHash]
		require.True(t, exists)
		require.False(t, otp.IsOrderPlacementInOperationsQueue(order))
	}

	// Verify all Short-Term orders can now be found in `ShortTermOrderHashToTxBytes` but are
	// still not in the operations queue.
	for _, order := range []types.Order{
		shortTermOrder1,
		shortTermOrder2,
		shortTermOrder3,
		shortTermOrder4,
		shortTermOrder5,
	} {
		orderHash := order.GetOrderHash()
		_, exists := otp.ShortTermOrderHashToTxBytes[orderHash]
		require.True(t, exists)
		require.False(t, otp.IsOrderPlacementInOperationsQueue(order))
	}
}

func TestMustAddShortTermOrderTxBytes_PanicsOnStatefulOrder(t *testing.T) {
	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20
	otp := types.NewOperationsToPropose()
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeShortTermOrder: called with stateful order ID (%+v)",
			order.OrderId,
		),
		func() {
			otp.MustAddShortTermOrderTxBytes(order, []byte{})
		},
	)
}

func TestMustAddShortTermOrderTxBytes_PanicsOnOrderInShortTermOrderHashToTxBytes(t *testing.T) {
	order := constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000
	otp := types.NewOperationsToPropose()
	otp.MustAddShortTermOrderTxBytes(
		order,
		[]byte{0, 1},
	)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustAddShortTermOrderTxBytes: Order (%s) already exists in `ShortTermOrderHashToTxBytes`.",
			order.GetOrderTextString(),
		),
		func() {
			otp.MustAddShortTermOrderTxBytes(order, []byte{0, 1})
		},
	)
}

func TestMustAddShortTermOrderPlacementToOperationsQueue(t *testing.T) {
	shortTermOrder1 := constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC
	shortTermOrder2 := constants.Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000
	shortTermOrder3 := constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000
	shortTermOrder4 := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	shortTermOrder5 := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16

	otp := types.NewOperationsToPropose()

	// Add `shortTermOrder1`, `shortTermOrder2`, and `shortTermOrder4` to the operations queue
	// and verify they're present in the operations queue.
	for _, order := range []types.Order{
		shortTermOrder1,
		shortTermOrder2,
		shortTermOrder4,
	} {
		otp.MustAddShortTermOrderTxBytes(order, order.GetOrderHash().ToBytes())
		otp.MustAddShortTermOrderPlacementToOperationsQueue(order)
		require.True(t, otp.IsOrderPlacementInOperationsQueue(order))
	}

	// Verify `shortTermOrder3` and `shortTermOrder5` are not present in the operations queue.
	for _, orderNotInOpQueue := range []types.Order{
		shortTermOrder3,
		shortTermOrder5,
	} {
		require.False(t, otp.IsOrderPlacementInOperationsQueue(orderNotInOpQueue))
	}

	// Add `shortTermOrder3` and `shortTermOrder5`to the operations queue and verify they're now
	// in the operations queue.
	for _, order := range []types.Order{
		shortTermOrder3,
		shortTermOrder5,
	} {
		otp.MustAddShortTermOrderTxBytes(order, order.GetOrderHash().ToBytes())
		otp.MustAddShortTermOrderPlacementToOperationsQueue(order)
		require.True(t, otp.IsOrderPlacementInOperationsQueue(order))
	}

	// Verify all Short-Term orders can now be found in `ShortTermOrderHashToTxBytes` but are
	// still not in the operations queue.
	for _, order := range []types.Order{
		shortTermOrder1,
		shortTermOrder2,
		shortTermOrder3,
		shortTermOrder4,
		shortTermOrder5,
	} {
		require.True(t, otp.IsOrderPlacementInOperationsQueue(order))
	}
}

func TestMustAddShortTermOrderPlacementToOperationsQueue_PanicsOnOrderInOrderHashesInOperationsQueue(t *testing.T) {
	order := constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000
	otp := types.NewOperationsToPropose()
	otp.OrderHashesInOperationsQueue[order.GetOrderHash()] = true
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustAddShortTermOrderPlacementToOperationsQueue: Order (%s) already exists in "+
				"`OrderHashesInOperationsQueue`.",
			order.GetOrderTextString(),
		),
		func() {
			otp.MustAddShortTermOrderPlacementToOperationsQueue(order)
		},
	)
}

func TestMustAddShortTermOrderPlacementToOperationsQueue_PanicsOnOrderNotInShortTermOrderHashToTxBytes(t *testing.T) {
	order := constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000
	otp := types.NewOperationsToPropose()
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustAddShortTermOrderPlacementToOperationsQueue: Order (%s) does not exist in "+
				"`ShortTermOrderHashToTxBytes`.",
			order.GetOrderTextString(),
		),
		func() {
			otp.MustAddShortTermOrderPlacementToOperationsQueue(order)
		},
	)
}

func TestRemoveShortTermOrderTxBytes(t *testing.T) {
	shortTermOrder1 := constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC
	shortTermOrder2 := constants.Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000
	shortTermOrder3 := constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000
	shortTermOrder4 := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	shortTermOrder5 := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16

	otp := types.NewOperationsToPropose()

	// Add all short term orders to ShortTermOrderHashToTxBytes.
	for _, order := range []types.Order{
		shortTermOrder1,
		shortTermOrder2,
		shortTermOrder3,
		shortTermOrder4,
		shortTermOrder5,
	} {
		otp.MustAddShortTermOrderTxBytes(order, order.GetOrderHash().ToBytes())
	}

	// Verify all Short-Term order hashes can now be found in `ShortTermOrderHashToTxBytes`,
	// and remove them.
	for _, order := range []types.Order{
		shortTermOrder1,
		shortTermOrder2,
		shortTermOrder3,
		shortTermOrder4,
		shortTermOrder5,
	} {
		require.Contains(t, otp.ShortTermOrderHashToTxBytes, order.GetOrderHash())
		otp.RemoveShortTermOrderTxBytes(order)
	}

	// Verify all order hashes have been removed from `ShortTermOrderHashToTxBytes`.
	require.Empty(t, otp.ShortTermOrderHashToTxBytes)
}

func TestRemoveShortTermOrderTxBytes_PanicsOnStatefulOrder(t *testing.T) {
	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20
	otp := types.NewOperationsToPropose()
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeShortTermOrder: called with stateful order ID (%+v)",
			order.OrderId,
		),
		func() {
			otp.RemoveShortTermOrderTxBytes(order)
		},
	)
}

func TestRemoveShortTermOrderTxBytes_PanicsOnOrderHashInOperationsQueue(t *testing.T) {
	order := constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000
	otp := types.NewOperationsToPropose()
	otp.OrderHashesInOperationsQueue[order.GetOrderHash()] = true
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"RemoveShortTermOrderTxBytes: Order (%s) exists in `OrderHashesInOperationsQueue`.",
			order.GetOrderTextString(),
		),
		func() {
			otp.RemoveShortTermOrderTxBytes(order)
		},
	)
}

func TestRemoveShortTermOrderTxBytes_PanicsOnOrderNotInShortTermOrderHashToTxBytes(t *testing.T) {
	order := constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000
	otp := types.NewOperationsToPropose()
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"RemoveShortTermOrderTxBytes: Order (%s) does not exist in `ShortTermOrderHashToTxBytes`.",
			order.GetOrderTextString(),
		),
		func() {
			otp.RemoveShortTermOrderTxBytes(order)
		},
	)
}

func TestMustAddStatefulOrderPlacementToOperationsQueue(t *testing.T) {
	statefulOrder1 := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15
	statefulOrder2 := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20
	statefulOrder3 := constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25
	statefulOrder4 := constants.LongTermOrder_Bob_Num1_Id3_Clob0_Buy10_Price40_GTBT10
	statefulOrder5 := constants.LongTermOrder_Alice_Num1_Id1_Clob0_Buy02BTC_Price10_GTB15

	otp := types.NewOperationsToPropose()

	// Add `statefulOrder1`, `statefulOrder2`, and `statefulOrder4` to the operations queue and
	// verify they're present.
	for _, order := range []types.Order{
		statefulOrder1,
		statefulOrder2,
		statefulOrder4,
	} {
		otp.MustAddStatefulOrderPlacementToOperationsQueue(order)
		require.True(t, otp.IsOrderPlacementInOperationsQueue(order))
	}

	// Verify `statefulOrder3` and `statefulOrder5` are not present in the operations queue.
	for _, orderNotInOpQueue := range []types.Order{
		statefulOrder3,
		statefulOrder5,
	} {
		require.False(t, otp.IsOrderPlacementInOperationsQueue(orderNotInOpQueue))
	}

	// Verify `statefulOrder3` and `statefulOrder5` can be added to the operations queue.
	for _, order := range []types.Order{
		statefulOrder3,
		statefulOrder5,
	} {
		otp.MustAddStatefulOrderPlacementToOperationsQueue(order)
	}

	// Verify all stateful orders can now be found in the operations queue.
	for _, order := range []types.Order{
		statefulOrder1,
		statefulOrder2,
		statefulOrder3,
		statefulOrder4,
		statefulOrder5,
	} {
		require.True(t, otp.IsOrderPlacementInOperationsQueue(order))
	}
}

func TestMustAddStatefulOrderPlacementToOperationsQueue_PanicsOnShortTermOrder(t *testing.T) {
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
	otp := types.NewOperationsToPropose()
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			order.OrderId,
		),
		func() {
			otp.MustAddStatefulOrderPlacementToOperationsQueue(order)
		},
	)
}

func TestMustAddStatefulOrderPlacementToOperationsQueue_PanicsOnOrderInOrderHashesInOperationsQueue(t *testing.T) {
	order := constants.LongTermOrder_Alice_Num1_Id2_Clob0_Sell02BTC_Price10_GTB15
	otp := types.NewOperationsToPropose()
	otp.MustAddStatefulOrderPlacementToOperationsQueue(order)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustAddStatefulOrderPlacementToOperationsQueue: Order (%s) already exists in "+
				"`OrderHashesInOperationsQueue`.",
			order.GetOrderTextString(),
		),
		func() {
			otp.MustAddStatefulOrderPlacementToOperationsQueue(order)
		},
	)
}

func TestMustAddMatchToOperationsQueue(t *testing.T) {
	shortTermOrder1 := constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC
	shortTermOrder2 := constants.Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000
	shortTermOrder3 := constants.Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price500000_GTB10
	longTermOrder1 := constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10
	longTermOrder2 := constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10
	liquidationOrder1 := constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000

	// Verify you can add a match containing `shortTermOrder1`, `shortTermOrder2`, and `shortTermOrder3`
	// to the operations queue.
	otp := types.NewOperationsToPropose()
	for _, shortTermOrder := range []types.Order{
		shortTermOrder1,
		shortTermOrder2,
		shortTermOrder3,
	} {
		orderHash := shortTermOrder.GetOrderHash()
		otp.MustAddShortTermOrderTxBytes(shortTermOrder, orderHash.ToBytes())
		otp.MustAddShortTermOrderPlacementToOperationsQueue(shortTermOrder)
	}
	otp.MustAddMatchToOperationsQueue(
		&shortTermOrder3,
		[]types.MakerFillWithOrder{
			{
				Order: shortTermOrder1,
				MakerFill: types.MakerFill{
					FillAmount:   50_000_000,
					MakerOrderId: shortTermOrder1.OrderId,
				},
			},
			{
				Order: shortTermOrder2,
				MakerFill: types.MakerFill{
					FillAmount:   1_000,
					MakerOrderId: shortTermOrder2.OrderId,
				},
			},
		},
	)

	// Verify you can add a match containing `shortTermOrder3` and `longTermOrder2` to the
	// operations queue.
	otp.MustAddStatefulOrderPlacementToOperationsQueue(longTermOrder2)
	otp.MustAddMatchToOperationsQueue(
		&longTermOrder2,
		[]types.MakerFillWithOrder{
			{
				Order: shortTermOrder3,
				MakerFill: types.MakerFill{
					FillAmount:   50_000_000,
					MakerOrderId: shortTermOrder3.OrderId,
				},
			},
		},
	)

	// Verify you can add a match containing `longTermOrder1` and `longTermOrder2` to the
	// operations queue.
	otp.MustAddStatefulOrderPlacementToOperationsQueue(longTermOrder1)
	otp.MustAddMatchToOperationsQueue(
		&longTermOrder1,
		[]types.MakerFillWithOrder{
			{
				Order: longTermOrder2,
				MakerFill: types.MakerFill{
					FillAmount:   50_000_000,
					MakerOrderId: longTermOrder2.OrderId,
				},
			},
		},
	)

	// Verify you can add a match containing `shortTermOrder3` and `liquidationOrder1` to the
	// operations queue.
	otp.MustAddMatchToOperationsQueue(
		&liquidationOrder1,
		[]types.MakerFillWithOrder{
			{
				Order: shortTermOrder3,
				MakerFill: types.MakerFill{
					FillAmount:   50_000_000,
					MakerOrderId: shortTermOrder3.OrderId,
				},
			},
		},
	)
}

func TestMustAddMatchToOperationsQueue_PanicsOnTakerOrderNotInOperationsQueue(t *testing.T) {
	takerOrder := constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10
	makerOrder := constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10
	otp := types.NewOperationsToPropose()
	otp.MustAddStatefulOrderPlacementToOperationsQueue(makerOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustAddMatchToOperationsQueue: Order (%s) does not exist in "+
				"`OrderHashesInOperationsQueue`.",
			takerOrder.GetOrderTextString(),
		),
		func() {
			otp.MustAddMatchToOperationsQueue(
				&takerOrder,
				[]types.MakerFillWithOrder{
					{
						Order: makerOrder,
						MakerFill: types.MakerFill{
							FillAmount:   25,
							MakerOrderId: makerOrder.OrderId,
						},
					},
				},
			)
		},
	)
}

func TestMustAddMatchToOperationsQueue_PanicsOnMakerOrderNotInOperationsQueue(t *testing.T) {
	takerOrder := constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10
	makerOrder := constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10
	otp := types.NewOperationsToPropose()
	otp.MustAddStatefulOrderPlacementToOperationsQueue(takerOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustAddMatchToOperationsQueue: Order (%s) does not exist in "+
				"`OrderHashesInOperationsQueue`.",
			makerOrder.GetOrderTextString(),
		),
		func() {
			otp.MustAddMatchToOperationsQueue(
				&takerOrder,
				[]types.MakerFillWithOrder{
					{
						Order: makerOrder,
						MakerFill: types.MakerFill{
							FillAmount:   25,
							MakerOrderId: makerOrder.OrderId,
						},
					},
				},
			)
		},
	)
}

func TestMustAddMatchToOperationsQueue_PanicsOnRegularTakerOrderWithZeroFills(t *testing.T) {
	takerOrder := constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10
	otp := types.NewOperationsToPropose()
	otp.MustAddStatefulOrderPlacementToOperationsQueue(takerOrder)
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"NewMatchOrdersInternalOperation: cannot create a match orders "+
				"internal operation with no maker fills: %+v",
			takerOrder,
		),
		func() {
			otp.MustAddMatchToOperationsQueue(
				&takerOrder,
				[]types.MakerFillWithOrder{},
			)
		},
	)
}

func TestMustAddMatchToOperationsQueue_PanicsOnLiquidationOrderWithZeroFills(t *testing.T) {
	liquidationTakerOrder := constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000
	otp := types.NewOperationsToPropose()
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"NewMatchPerpetualLiquidationInternalOperation: cannot create a match perpetual "+
				"liquidation internal operation with no maker fills: %+v",
			&liquidationTakerOrder,
		),
		func() {
			otp.MustAddMatchToOperationsQueue(
				&liquidationTakerOrder,
				[]types.MakerFillWithOrder{},
			)
		},
	)
}

func TestMustAddDeleveraingToOperationsQueue(t *testing.T) {
	otp := types.NewOperationsToPropose()

	otp.MustAddDeleveragingToOperationsQueue(
		constants.Alice_Num0,
		0,
		[]types.MatchPerpetualDeleveraging_Fill{
			{
				OffsettingSubaccountId: constants.Bob_Num0,
				FillAmount:             5,
			},
		},
		false,
	)
}

func TestMustAddDeleveragingToOperationsQueue_Panics(t *testing.T) {
	tests := map[string]struct {
		liquidatedSubaccountId satypes.SubaccountId
		perpetualId            uint32
		fills                  []types.MatchPerpetualDeleveraging_Fill

		expectedPanic string
	}{
		"number of fills is zero": {
			liquidatedSubaccountId: constants.Alice_Num0,
			perpetualId:            0,
			fills:                  []types.MatchPerpetualDeleveraging_Fill{},

			expectedPanic: fmt.Sprintf(
				"MustAddDeleveragingToOperationsQueue: number of fills is zero. "+
					"liquidatedSubaccountId = (%+v), perpetualId = (%d)",
				constants.Alice_Num0,
				0,
			),
		},
		"fill amount is zero": {
			liquidatedSubaccountId: constants.Alice_Num0,
			perpetualId:            0,
			fills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Bob_Num0,
					FillAmount:             0,
				},
			},

			expectedPanic: fmt.Sprintf(
				"MustAddDeleveragingToOperationsQueue: fill amount is zero. "+
					"liquidatedSubaccountId = (%+v), perpetualId = (%d), fill = (%+v)",
				constants.Alice_Num0,
				0,
				types.MatchPerpetualDeleveraging_Fill{
					OffsettingSubaccountId: constants.Bob_Num0,
					FillAmount:             0,
				},
			),
		},
		"offsetting is the same as liquidated": {
			liquidatedSubaccountId: constants.Alice_Num0,
			perpetualId:            0,
			fills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Alice_Num0,
					FillAmount:             5,
				},
			},

			expectedPanic: fmt.Sprintf(
				"MustAddDeleveragingToOperationsQueue: offsetting subaccount is the same as liquidated subaccount. "+
					"liquidatedSubaccountId = (%+v), perpetualId = (%d), fill = (%+v)",
				constants.Alice_Num0,
				0,
				types.MatchPerpetualDeleveraging_Fill{
					OffsettingSubaccountId: constants.Alice_Num0,
					FillAmount:             5,
				},
			),
		},
		"duplicated subaccount ids": {
			liquidatedSubaccountId: constants.Alice_Num0,
			perpetualId:            0,
			fills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Bob_Num0,
					FillAmount:             5,
				},
				{
					OffsettingSubaccountId: constants.Bob_Num0,
					FillAmount:             5,
				},
			},

			expectedPanic: fmt.Sprintf(
				"MustAddDeleveragingToOperationsQueue: duplicated subaccount ids. "+
					"liquidatedSubaccountId = (%+v), perpetualId = (%d), fill = (%+v)",
				constants.Alice_Num0,
				0,
				types.MatchPerpetualDeleveraging_Fill{
					OffsettingSubaccountId: constants.Bob_Num0,
					FillAmount:             5,
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			otp := types.NewOperationsToPropose()
			require.PanicsWithValue(
				t,
				tc.expectedPanic,
				func() {
					otp.MustAddDeleveragingToOperationsQueue(
						tc.liquidatedSubaccountId,
						tc.perpetualId,
						tc.fills,
						false,
					)
				},
			)
		})
	}
}

func TestMustAddOrderRemovalToOperationsQueue(t *testing.T) {
	tests := map[string]struct {
		orderId       types.OrderId
		removalReason types.OrderRemoval_RemovalReason
	}{
		"order removal for stateful order": {
			orderId:       constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0,
			removalReason: types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			otp := types.NewOperationsToPropose()

			otp.MustAddOrderRemovalToOperationsQueue(
				tc.orderId,
				tc.removalReason,
			)
		})
	}
}

func TestMustAddOrderRemovalToOperationsQueue_Panics(t *testing.T) {
	tests := map[string]struct {
		orderId       types.OrderId
		removalReason types.OrderRemoval_RemovalReason
		expectedPanic string
	}{
		"order removal reason unspecified": {
			orderId:       constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0,
			removalReason: types.OrderRemoval_REMOVAL_REASON_UNSPECIFIED,
			expectedPanic: "MustAddOrderRemovalToOperationsQueue: removal reason unspecified",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			otp := types.NewOperationsToPropose()
			require.PanicsWithValue(
				t,
				tc.expectedPanic,
				func() {
					otp.MustAddOrderRemovalToOperationsQueue(
						tc.orderId,
						tc.removalReason,
					)
				},
			)
		})
	}
}

func TestGetOperationsToReplay_PanicsOnNonexistentShortTermOrderHashToTxBytesOrder(t *testing.T) {
	otp := types.NewOperationsToPropose()
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
	operation := types.NewShortTermOrderPlacementInternalOperation(order)
	otp.OperationsQueue = []types.InternalOperation{
		operation,
	}
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"GetOperationsToReplay: Short-Term order (%s) does not exist in "+
				"`ShortTermOrderHashToTxBytes`.",
			order.GetOrderTextString(),
		),
		func() {
			otp.GetOperationsToReplay()
		},
	)
}

func TestGetOperationsToReplay_Success(t *testing.T) {
	tests := map[string]struct {
		// Params.
		setup func(otp *types.OperationsToPropose)

		// Expectations.
		expectedOperations          []types.InternalOperation
		expectedShortTermOrderBytes map[types.OrderHash][]byte
	}{
		"Can get short term orders to replay": {
			setup: func(otp *types.OperationsToPropose) {
				shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
				otp.MustAddShortTermOrderTxBytes(shortTermOrder, shortTermOrder.GetOrderHash().ToBytes())
				otp.MustAddShortTermOrderPlacementToOperationsQueue(shortTermOrder)
			},
			expectedOperations: []types.InternalOperation{
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
				),
			},
			expectedShortTermOrderBytes: map[types.OrderHash][]byte{
				constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.GetOrderHash(): constants.
					Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.GetOrderHash().ToBytes(),
			},
		},
		"Can get long term orders to replay": {
			setup: func(otp *types.OperationsToPropose) {
				longTermOrder := constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25
				otp.MustAddStatefulOrderPlacementToOperationsQueue(longTermOrder)
			},
			expectedOperations: []types.InternalOperation{
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25,
				),
			},
			expectedShortTermOrderBytes: map[types.OrderHash][]byte{},
		},
		"Can get order matches to replay": {
			setup: func(otp *types.OperationsToPropose) {
				takerOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
				makerOrder := constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15
				makerFillsWithOrders := []types.MakerFillWithOrder{
					{
						MakerFill: types.MakerFill{
							FillAmount:   5,
							MakerOrderId: makerOrder.OrderId,
						},
						Order: makerOrder,
					},
				}
				otp.MustAddStatefulOrderPlacementToOperationsQueue(makerOrder)
				otp.MustAddShortTermOrderTxBytes(takerOrder, []byte{4, 0, 8})
				otp.MustAddShortTermOrderPlacementToOperationsQueue(takerOrder)
				otp.MustAddMatchToOperationsQueue(&takerOrder, makerFillsWithOrders)
			},
			expectedOperations: []types.InternalOperation{
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15,
				),
				types.NewShortTermOrderPlacementInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
				),
				types.NewMatchOrdersInternalOperation(
					constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16,
					[]types.MakerFill{
						{
							MakerOrderId: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
							FillAmount:   5,
						},
					},
				),
			},
			expectedShortTermOrderBytes: map[types.OrderHash][]byte{
				constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.GetOrderHash(): {4, 0, 8},
			},
		},
		"Can get liquidations matches to replay": {
			setup: func(otp *types.OperationsToPropose) {
				liquidationOrder := constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC
				makerOrder := constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15
				otp.MustAddStatefulOrderPlacementToOperationsQueue(makerOrder)
				otp.MustAddMatchToOperationsQueue(
					&liquidationOrder,
					[]types.MakerFillWithOrder{
						{
							MakerFill: types.MakerFill{
								FillAmount:   5,
								MakerOrderId: makerOrder.OrderId,
							},
							Order: makerOrder,
						},
					},
				)
			},
			expectedOperations: []types.InternalOperation{
				types.NewPreexistingStatefulOrderPlacementInternalOperation(
					constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15,
				),
				types.NewMatchPerpetualLiquidationInternalOperation(
					&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
					[]types.MakerFill{
						{
							FillAmount:   5,
							MakerOrderId: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
						},
					},
				),
			},
			expectedShortTermOrderBytes: map[types.OrderHash][]byte{},
		},
		"Can get deleveraging matches to replay": {
			setup: func(otp *types.OperationsToPropose) {
				otp.MustAddDeleveragingToOperationsQueue(
					constants.Alice_Num0,
					0,
					[]types.MatchPerpetualDeleveraging_Fill{
						{
							OffsettingSubaccountId: constants.Bob_Num0,
							FillAmount:             10,
						},
					},
					false,
				)
			},
			expectedOperations: []types.InternalOperation{
				{
					Operation: &types.InternalOperation_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchPerpetualDeleveraging{
								MatchPerpetualDeleveraging: &types.MatchPerpetualDeleveraging{
									Liquidated:  constants.Alice_Num0,
									PerpetualId: 0,
									Fills: []types.MatchPerpetualDeleveraging_Fill{
										{
											OffsettingSubaccountId: constants.Bob_Num0,
											FillAmount:             10,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedShortTermOrderBytes: map[types.OrderHash][]byte{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			tc.setup(otp)

			// Verify expectations.
			operation, shortTermOrdersTxBytes := otp.GetOperationsToReplay()
			require.Equal(t, tc.expectedOperations, operation)
			require.Equal(t, tc.expectedShortTermOrderBytes, shortTermOrdersTxBytes)
		})
	}
}

func TestGetOperationsToPropose_PanicsOnNonexistentShortTermOrderHashToTxBytesOrder(t *testing.T) {
	otp := types.NewOperationsToPropose()
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
	operation := types.NewShortTermOrderPlacementInternalOperation(order)
	otp.OperationsQueue = []types.InternalOperation{
		operation,
	}
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"GetOperationsToPropose: Order (%s) does not exist in "+
				"`ShortTermOrderHashToTxBytes`.",
			order.GetOrderTextString(),
		),
		func() {
			otp.GetOperationsToPropose()
		},
	)
}

func TestGetOperationsToPropose_Success(t *testing.T) {
	tests := map[string]struct {
		// Params.
		setup func(otp *types.OperationsToPropose)

		// Expectations.
		expectedOperations []types.OperationRaw
	}{
		"Short term orders are included in operations to propose": {
			setup: func(otp *types.OperationsToPropose) {
				shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
				// Dummy bytes for testing.
				otp.MustAddShortTermOrderTxBytes(shortTermOrder, []byte{4, 0, 8})
				otp.MustAddShortTermOrderPlacementToOperationsQueue(shortTermOrder)
			},
			expectedOperations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_ShortTermOrderPlacement{
						ShortTermOrderPlacement: []byte{4, 0, 8},
					},
				},
			},
		},
		"Stateful orders do not get included in operations to propose": {
			setup: func(otp *types.OperationsToPropose) {
				statefulOrder := constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25
				otp.MustAddStatefulOrderPlacementToOperationsQueue(statefulOrder)
			},
			expectedOperations: []types.OperationRaw{},
		},
		"Order matches are included in operations to propose": {
			setup: func(otp *types.OperationsToPropose) {
				takerOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
				makerOrder := constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15
				makerFillsWithOrders := []types.MakerFillWithOrder{
					{
						MakerFill: types.MakerFill{
							FillAmount:   5,
							MakerOrderId: makerOrder.OrderId,
						},
						Order: makerOrder,
					},
				}
				// This is not included in operations to propose.
				otp.MustAddStatefulOrderPlacementToOperationsQueue(makerOrder)
				otp.MustAddShortTermOrderTxBytes(takerOrder, []byte{4, 0, 8})
				otp.MustAddShortTermOrderPlacementToOperationsQueue(takerOrder)
				otp.MustAddMatchToOperationsQueue(&takerOrder, makerFillsWithOrders)
			},
			expectedOperations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_ShortTermOrderPlacement{
						ShortTermOrderPlacement: []byte{4, 0, 8},
					},
				},
				{
					Operation: &types.OperationRaw_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchOrders{
								MatchOrders: &types.MatchOrders{
									TakerOrderId: constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16.OrderId,
									Fills: []types.MakerFill{
										{
											FillAmount:   5,
											MakerOrderId: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"Liquidation matches are included in operations to propose": {
			setup: func(otp *types.OperationsToPropose) {
				liquidationOrder := constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC
				makerOrder := constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15
				// This is not included in operations to propose.
				otp.MustAddStatefulOrderPlacementToOperationsQueue(makerOrder)
				otp.MustAddMatchToOperationsQueue(
					&liquidationOrder,
					[]types.MakerFillWithOrder{
						{
							MakerFill: types.MakerFill{
								FillAmount:   5,
								MakerOrderId: makerOrder.OrderId,
							},
							Order: makerOrder,
						},
					},
				)
			},
			expectedOperations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchPerpetualLiquidation{
								MatchPerpetualLiquidation: &types.MatchPerpetualLiquidation{
									Liquidated:  constants.Alice_Num0,
									ClobPairId:  constants.ClobPair_Btc.Id,
									PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().GetPerpetualId(),
									TotalSize:   20,
									IsBuy:       false,
									Fills: []types.MakerFill{
										{
											FillAmount:   5,
											MakerOrderId: constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15.OrderId,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"Deleveraging matches are included in operations to propose": {
			setup: func(otp *types.OperationsToPropose) {
				otp.MustAddDeleveragingToOperationsQueue(
					constants.Alice_Num0,
					0,
					[]types.MatchPerpetualDeleveraging_Fill{
						{
							OffsettingSubaccountId: constants.Bob_Num0,
							FillAmount:             10,
						},
					},
					false,
				)
			},
			expectedOperations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_Match{
						Match: &types.ClobMatch{
							Match: &types.ClobMatch_MatchPerpetualDeleveraging{
								MatchPerpetualDeleveraging: &types.MatchPerpetualDeleveraging{
									Liquidated:  constants.Alice_Num0,
									PerpetualId: 0,
									Fills: []types.MatchPerpetualDeleveraging_Fill{
										{
											OffsettingSubaccountId: constants.Bob_Num0,
											FillAmount:             10,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"Order Removals are included in operations to propose": {
			setup: func(otp *types.OperationsToPropose) {
				otp.MustAddOrderRemovalToOperationsQueue(
					constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0,
					types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
				)
			},
			expectedOperations: []types.OperationRaw{
				{
					Operation: &types.OperationRaw_OrderRemoval{
						OrderRemoval: &types.OrderRemoval{
							OrderId:       constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0,
							RemovalReason: types.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup the test.
			otp := types.NewOperationsToPropose()
			tc.setup(otp)

			// Verify expectations.
			require.Equal(t, tc.expectedOperations, otp.GetOperationsToPropose())
		})
	}
}
