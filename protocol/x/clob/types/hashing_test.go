package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/stretchr/testify/require"
)

// TODO(DEC-1654) Add more robust tests for different fields
func TestHashing_Operation_Equal(t *testing.T) {
	tests := map[string]struct {
		operationOne types.Operation
		operationTwo types.Operation
	}{
		"place order": {
			operationOne: types.NewOrderPlacementOperation(types.Order{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				Side:         types.Order_SIDE_BUY,
				Quantums:     100_000_000,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
			}),
			operationTwo: types.NewOrderPlacementOperation(types.Order{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				Side:         types.Order_SIDE_BUY,
				Quantums:     100_000_000,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10}, // ptrs to diff vals
			}),
		},
		"cancel order": {
			operationOne: types.NewOrderCancellationOperation(&types.MsgCancelOrder{
				OrderId:      constants.InvalidSubaccountIdNumber_OrderId,
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 1},
			}),
			operationTwo: types.NewOrderCancellationOperation(&types.MsgCancelOrder{
				OrderId:      constants.InvalidSubaccountIdNumber_OrderId,
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 1},
			}),
		},
		"cancel order order id structs": {
			operationOne: types.NewOrderCancellationOperation(&types.MsgCancelOrder{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					ClobPairId:   0,
				},
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 1},
			}),
			operationTwo: types.NewOrderCancellationOperation(&types.MsgCancelOrder{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 1},
			}),
		},
		"match order": {
			operationOne: types.NewMatchOperation(
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				[]types.MakerFill{
					{
						FillAmount:   100_000_000, // 1 BTC
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
					},
				}),
			operationTwo: types.NewMatchOperation(
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				[]types.MakerFill{
					{
						FillAmount: 100_000_000, // 1 BTC
						// should be the same
						MakerOrderId: types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 0, ClobPairId: 0},
					},
				},
			),
		},
		"liquidation order": {
			operationOne: types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
				Liquidated:  constants.Carl_Num0,
				ClobPairId:  0,
				PerpetualId: 0,
				TotalSize:   1,
				IsBuy:       true,
				Fills: []types.MakerFill{
					{
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						FillAmount:   0,
					},
				},
			}),
			operationTwo: types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
				Liquidated:  constants.Carl_Num0,
				ClobPairId:  0,
				PerpetualId: 0,
				TotalSize:   1,
				IsBuy:       true,
				Fills: []types.MakerFill{
					{
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						FillAmount:   0,
					},
				},
			}),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.True(t, types.IsEqual(&tc.operationOne, &tc.operationTwo))
		})
	}
	operationOne := types.NewOrderPlacementOperation(types.Order{
		OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
		Side:         types.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
	})
	operationTwo := types.NewOrderPlacementOperation(types.Order{
		OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
		Side:         types.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
	})
	// direct pointer comparison no address
	require.True(t, types.IsEqual(&operationOne, &operationTwo))
}

func TestHashing_Operation_NotEqual(t *testing.T) {
	tests := map[string]struct {
		operationOne types.Operation
		operationTwo types.Operation
	}{
		"place order quantums": {
			operationOne: types.NewOrderPlacementOperation(types.Order{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				Side:         types.Order_SIDE_BUY,
				Quantums:     100_000_000,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
			}),
			operationTwo: types.NewOrderPlacementOperation(types.Order{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				Side:         types.Order_SIDE_SELL,
				Quantums:     100_000, // different
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
			}),
		},
		"place order clobPairId": {
			operationOne: types.NewOrderPlacementOperation(types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 1},
				Side:         types.Order_SIDE_BUY,
				Quantums:     100_000_000,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
			}),
			operationTwo: types.NewOrderPlacementOperation(types.Order{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				Side:         types.Order_SIDE_SELL,
				Quantums:     100_000, // different
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
			}),
		},
		"cancel order goodTilBlock": {
			operationOne: types.NewOrderCancellationOperation(&types.MsgCancelOrder{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 1},
			}),
			operationTwo: types.NewOrderCancellationOperation(&types.MsgCancelOrder{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 10}, // different
			}),
		},
		"cancel order order id structs": {
			operationOne: types.NewOrderCancellationOperation(&types.MsgCancelOrder{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					ClobPairId:   1, // different
				},
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 1},
			}),
			operationTwo: types.NewOrderCancellationOperation(&types.MsgCancelOrder{
				OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
				GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: 1},
			}),
		},
		"match order fill amount": {
			operationOne: types.NewMatchOperation(
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				[]types.MakerFill{
					{
						FillAmount:   100_000_000,
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
					},
				},
			),
			operationTwo: types.NewMatchOperation(
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				[]types.MakerFill{
					{
						FillAmount:   100,
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
					},
				},
			),
		},
		"match order fill array size": {
			operationOne: types.NewMatchOperation(
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				[]types.MakerFill{
					{
						FillAmount:   100_000_000,
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
					},
				},
			),
			operationTwo: types.NewMatchOperation(
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				[]types.MakerFill{}, // empty
			),
		},
		"match order and liquidation order": {
			operationOne: types.NewMatchOperation(
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				[]types.MakerFill{
					{
						FillAmount:   100_000_000,
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
					},
				},
			),
			operationTwo: types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
				Liquidated:  constants.Carl_Num0,
				ClobPairId:  0,
				PerpetualId: 0,
				TotalSize:   1,
				IsBuy:       true,
				Fills: []types.MakerFill{
					{
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						FillAmount:   0,
					},
				},
			}),
		},
		"liquidation order fill amount": {
			operationOne: types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
				Liquidated:  constants.Carl_Num0,
				ClobPairId:  0,
				PerpetualId: 0,
				TotalSize:   1,
				IsBuy:       true,
				Fills: []types.MakerFill{
					{
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						FillAmount:   0,
					},
				},
			}),
			operationTwo: types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
				Liquidated:  constants.Carl_Num0,
				ClobPairId:  0,
				PerpetualId: 0,
				TotalSize:   1,
				IsBuy:       true,
				Fills: []types.MakerFill{
					{
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						FillAmount:   100, // wrong
					},
				},
			}),
		},
		"liquidation order fill array size": {
			operationOne: types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
				Liquidated:  constants.Carl_Num0,
				ClobPairId:  0,
				PerpetualId: 0,
				TotalSize:   1,
				IsBuy:       true,
				Fills: []types.MakerFill{
					{
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						FillAmount:   0,
					},
				},
			}),
			operationTwo: types.NewMatchOperationFromPerpetualLiquidation(types.MatchPerpetualLiquidation{
				Liquidated:  constants.Carl_Num0,
				ClobPairId:  0,
				PerpetualId: 0,
				TotalSize:   1,
				IsBuy:       true,
				Fills: []types.MakerFill{
					{
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						FillAmount:   100,
					},
					{
						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.GetOrderId(),
						FillAmount:   100,
					},
				},
			}),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.False(t, types.IsEqual(&tc.operationOne, &tc.operationTwo))
		})
	}
}

func Test_Hashing_Order_IsEqual(t *testing.T) {
	order1 := types.Order{
		OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
		Side:         types.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	order2 := types.Order{
		OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
		Side:         types.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10}, // pointers are different here
	}
	require.False(t, order1.GoodTilOneof == order2.GoodTilOneof)
	require.True(t, types.IsEqual(&order1, &order2))

	order1 = types.Order{
		OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
		Side:         types.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
	}
	order2 = types.Order{
		OrderId:      constants.OrderId_Alice_Num0_ClientId0_Clob0,
		Side:         types.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000_000,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 100}, // diff gtb
	}
	require.False(t, order1.GoodTilOneof == order2.GoodTilOneof)
	require.False(t, types.IsEqual(&order1, &order2))
}

func Test_IsEqualPanicOnNil(t *testing.T) {
	require.Panics(
		t,
		func() {
			var ptrOne *types.Order
			var ptrTwo *types.Order
			types.IsEqual(ptrOne, ptrTwo)
		},
	)
	require.Panics(
		t,
		func() {
			types.IsEqual(&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10, nil)
		},
	)
}
