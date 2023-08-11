package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGetOperationHash(t *testing.T) {
	operation := types.Operation{}
	require.Equal(
		t,
		types.OperationHash(constants.OperationHash_Empty),
		operation.GetOperationHash(),
	)

	orderPlacementOperation := types.NewOrderPlacementOperation(
		constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
	)
	require.Equal(
		t,
		constants.OperationHash_Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
		orderPlacementOperation.GetOperationHash(),
	)

	cancellationOperation := types.NewOrderCancellationOperation(
		&constants.CancelOrder_Alice_Num1_Id13_Clob0_GTB25,
	)
	require.Equal(
		t,
		constants.OperationHash_CancelOrder_Alice_Num1_Id13_Clob0_GTB25,
		cancellationOperation.GetOperationHash(),
	)

	matchOperation := types.NewMatchOperation(
		&constants.LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC,
		[]types.MakerFill{},
	)
	require.Equal(
		t,
		constants.OperationHash_MatchLiquidationOrder_Alice_Num0_Sell20_Price25_BTC,
		matchOperation.GetOperationHash(),
	)

	operation1 := types.NewOrderPlacementOperation(
		types.Order{
			GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		},
	)
	operation2 := types.NewOrderPlacementOperation(
		types.Order{
			GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
		},
	)
	require.NotEqual(
		t,
		operation1.GetOperationHash(),
		operation2.GetOperationHash(),
	)
}
