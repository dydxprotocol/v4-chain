package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestNewClobMatchFromMatchOrders(t *testing.T) {
	msgMatchOrder := &types.MatchOrders{
		TakerOrderId: constants.OrderId_Alice_Num0_ClientId0_Clob0,
		Fills: []types.MakerFill{
			{
				MakerOrderId: constants.OrderId_Bob_Num0_ClientId0_Clob0,
				FillAmount:   30,
			},
		},
	}

	msgClobMatch := types.NewClobMatchFromMatchOrders(msgMatchOrder)

	require.Equal(t, msgMatchOrder, msgClobMatch.GetMatchOrders())
}

func TestNewClobMatchFromMatchPerpetualLiquidation(t *testing.T) {
	msgMatchPerpetualLiquidation := &types.MatchPerpetualLiquidation{
		ClobPairId:  5,
		IsBuy:       true,
		TotalSize:   30,
		Liquidated:  constants.Alice_Num1,
		PerpetualId: 1,
		Fills: []types.MakerFill{
			{
				MakerOrderId: constants.OrderId_Alice_Num0_ClientId1_Clob0,
				FillAmount:   10,
			},
			{
				MakerOrderId: constants.OrderId_Alice_Num0_ClientId2_Clob0,
				FillAmount:   20,
			},
		},
	}

	msgClobMatch := types.NewClobMatchFromMatchPerpetualLiquidation(msgMatchPerpetualLiquidation)

	require.Equal(t, msgMatchPerpetualLiquidation, msgClobMatch.GetMatchPerpetualLiquidation())
}
