package keeper_test

import (
	"testing"

	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestStatefulOrderTest(t *testing.T) {
	tApp := testApp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	tApp.App.ClobKeeper.SetLongTermOrderPlacement(
		ctx,
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
		1,
	)

	tApp.App.ClobKeeper.SetOrderFillAmount(
		ctx,
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
		5,
		123456789,
	)

	res, err := tApp.App.ClobKeeper.StatefulOrder(
		ctx,
		&types.QueryStatefulOrderRequest{
			OrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
		},
	)

	require.NoError(t, err)
	require.Equal(
		t,
		&types.QueryStatefulOrderResponse{
			OrderPlacement: types.LongTermOrderPlacement{
				Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
				PlacementIndex: types.TransactionOrdering{
					BlockHeight:      1,
					TransactionIndex: 0,
				},
			},
			FillAmount: 5,
			Triggered:  false,
		},
		res,
	)
}
