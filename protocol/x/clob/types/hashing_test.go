package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

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
