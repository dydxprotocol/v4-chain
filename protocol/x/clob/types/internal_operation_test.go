package types_test

import (
	fmt "fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestNewShortTermOrderPlacementInternalOperation(t *testing.T) {
	order := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
	operation := types.InternalOperation{
		Operation: &types.InternalOperation_ShortTermOrderPlacement{
			ShortTermOrderPlacement: types.NewMsgPlaceOrder(order),
		},
	}
	require.Equal(
		t,
		operation,
		types.NewShortTermOrderPlacementInternalOperation(order),
	)
}

func TestNewShortTermOrderPlacementInternalOperation_PanicsOnStatefulOrder(t *testing.T) {
	statefulOrder := constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeShortTermOrder: called with stateful order ID (%+v)",
			statefulOrder.OrderId,
		),
		func() {
			types.NewShortTermOrderPlacementInternalOperation(statefulOrder)
		},
	)
}

func TestNewPreexistingStatefulOrderPlacementInternalOperation(t *testing.T) {
	order := constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10
	operation := types.InternalOperation{
		Operation: &types.InternalOperation_PreexistingStatefulOrder{
			PreexistingStatefulOrder: &order.OrderId,
		},
	}
	require.Equal(
		t,
		operation,
		types.NewPreexistingStatefulOrderPlacementInternalOperation(order),
	)
}

func TestNewPreexistingStatefulOrderPlacementInternalOperation_PanicsOnShortTermOrder(t *testing.T) {
	shortTermOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrder.OrderId,
		),
		func() {
			types.NewPreexistingStatefulOrderPlacementInternalOperation(shortTermOrder)
		},
	)
}

func TestNewMatchOrdersInternalOperation(t *testing.T) {
	takerOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
	makerOrder := constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15
	makerFills := []types.MakerFill{
		{
			FillAmount:   5,
			MakerOrderId: makerOrder.OrderId,
		},
	}
	operation := types.InternalOperation{
		Operation: &types.InternalOperation_Match{
			Match: &types.ClobMatch{
				Match: &types.ClobMatch_MatchOrders{
					MatchOrders: &types.MatchOrders{
						TakerOrderId: takerOrder.GetOrderId(),
						Fills:        makerFills,
					},
				},
			},
		},
	}
	require.Equal(
		t,
		operation,
		types.NewMatchOrdersInternalOperation(
			takerOrder,
			makerFills,
		),
	)
}

func TestNewMatchOrdersInternalOperation_PanicsOnZeroFills(t *testing.T) {
	takerOrder := constants.Order_Alice_Num0_Id0_Clob0_Buy10_Price10_GTB16
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"NewMatchOrdersInternalOperation: cannot create a match orders "+
				"internal operation with no maker fills: %+v",
			takerOrder,
		),
		func() {
			types.NewMatchOrdersInternalOperation(
				takerOrder,
				[]types.MakerFill{},
			)
		},
	)
}

func TestNewMatchPerpetualLiquidationInternalOperation(t *testing.T) {
	makerOrder := constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTBT15_StopLoss15
	makerFills := []types.MakerFill{
		{
			FillAmount:   5,
			MakerOrderId: makerOrder.OrderId,
		},
	}

	liquidationTakerOrder := constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC
	liqOperation := types.InternalOperation{
		Operation: &types.InternalOperation_Match{
			Match: &types.ClobMatch{
				Match: &types.ClobMatch_MatchPerpetualLiquidation{
					MatchPerpetualLiquidation: &types.MatchPerpetualLiquidation{
						Liquidated:  liquidationTakerOrder.GetSubaccountId(),
						ClobPairId:  liquidationTakerOrder.GetClobPairId().ToUint32(),
						PerpetualId: liquidationTakerOrder.MustGetLiquidatedPerpetualId(),
						TotalSize:   liquidationTakerOrder.GetBaseQuantums().ToUint64(),
						IsBuy:       liquidationTakerOrder.IsBuy(),
						Fills:       makerFills,
					},
				},
			},
		},
	}
	require.Equal(
		t,
		liqOperation,
		types.NewMatchPerpetualLiquidationInternalOperation(
			&liquidationTakerOrder,
			makerFills,
		),
	)
}

func TestNewMatchPerpetualLiquidationInternalOperation_PanicsOnZeroFills(t *testing.T) {
	liquidationTakerOrder := constants.LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"NewMatchPerpetualLiquidationInternalOperation: cannot create a match perpetual "+
				"liquidation internal operation with no maker fills: %+v",
			&liquidationTakerOrder,
		),
		func() {
			types.NewMatchPerpetualLiquidationInternalOperation(
				&liquidationTakerOrder,
				[]types.MakerFill{},
			)
		},
	)
}

func TestNewMatchPerpetualLiquidationInternalOperation_PanicsOnNonLiquidationOrder(t *testing.T) {
	takerOrder := constants.Order_Alice_Num0_Id0_Clob0_Sell5_Price10_GTB20
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"NewMatchPerpetualLiquidationInternalOperation: called with a non-liquidation order: %+v",
			&takerOrder,
		),
		func() {
			types.NewMatchPerpetualLiquidationInternalOperation(
				&takerOrder,
				[]types.MakerFill{},
			)
		},
	)
}

func TestNewOrderRemovalInternalOperation_PanicsOnShortTermOrderId(t *testing.T) {
	shortTermOrderId := constants.OrderId_Alice_Num0_ClientId0_Clob0
	require.PanicsWithValue(
		t,
		fmt.Sprintf(
			"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
			shortTermOrderId,
		),
		func() {
			types.NewOrderRemovalInternalOperation(
				shortTermOrderId,
				types.OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE,
			)
		},
	)
}

func TestNewOrderRemovalInternalOperation_PanicsOnUnspecifiedRemovalReason(t *testing.T) {
	require.PanicsWithValue(
		t,
		"NewOrderRemovalInternalOperation: removal reason unspecified",
		func() {
			types.NewOrderRemovalInternalOperation(
				constants.LongTermOrderId_Alice_Num0_ClientId0_Clob0,
				types.OrderRemoval_REMOVAL_REASON_UNSPECIFIED,
			)
		},
	)
}
