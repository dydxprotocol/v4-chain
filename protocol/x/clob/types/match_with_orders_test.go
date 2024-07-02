package types_test

import (
	"errors"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestPerformStatelessMatchOrdersValidation(t *testing.T) {
	tests := map[string]struct {
		makerOrder types.MatchableOrder
		takerOrder types.MatchableOrder
		fillAmount uint64

		expectedError error
	}{
		"Stateless match validation: match constitutes a self-trade": {
			makerOrder: &types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
				Side:         types.Order_SIDE_BUY,
				Quantums:     100_000_000,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
			},
			takerOrder: &types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
				Side:         types.Order_SIDE_SELL,
				Quantums:     100_000_000,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
			},
			fillAmount:    50_000_000, // .5 BTC
			expectedError: errors.New("Match constitutes a self-trade"),
		},
		"Stateless match validation: fillAmount must be greater than 0": {
			makerOrder:    &constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			takerOrder:    &constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
			fillAmount:    0,
			expectedError: errors.New("fillAmount must be greater than 0"),
		},
		"Stateless match validation: clobPairIds do not match": {
			makerOrder: &types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: constants.ClobPair_Eth.Id},
				Side:         types.Order_SIDE_BUY,
				Quantums:     1000,
				Subticks:     1000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
			},
			takerOrder:    &constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
			fillAmount:    100_000_000, // 1 BTC
			expectedError: errors.New("ClobPairIds do not match"),
		},
		"Stateless match validation: matches are on the same side of the book": {
			makerOrder: &types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
				Side:         types.Order_SIDE_BUY,
				Quantums:     10,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 6},
			},
			takerOrder: &types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
				Side:         types.Order_SIDE_BUY,
				Quantums:     10,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 6},
			},
			fillAmount:    100_000_000,
			expectedError: errors.New("Orders are not on opposing sides of the book in match"),
		},
		"Stateless match validation: orders dont cross with maker buy order": {
			makerOrder:    &constants.Order_Carl_Num0_Id1_Clob0_Buy1BTC_Price49999,
			takerOrder:    &constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
			fillAmount:    100_000_000, // 1 BTC
			expectedError: errors.New("Orders do not cross in match"),
		},
		"Stateless match validation: orders dont cross with taker buy order": {
			makerOrder:    &constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
			takerOrder:    &constants.Order_Carl_Num0_Id1_Clob0_Buy1BTC_Price49999,
			fillAmount:    100_000_000, // 1 BTC
			expectedError: errors.New("Orders do not cross in match"),
		},
		"Stateless match validation: minimum initial order quantums exceeds fill amount": {
			makerOrder:    &constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
			takerOrder:    &constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			fillAmount:    200_000_000, // 2 BTC. Too big!
			expectedError: errors.New("Minimum initial order quantums exceeds fill amount"),
		},
		"Stateless match validation: maker order is a liquidation order": {
			makerOrder:    &constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
			takerOrder:    &constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			fillAmount:    100_000_000, // 1 BTC.
			expectedError: errors.New("Liquidation order cannot be matched as a maker order"),
		},
		"Stateless match validation: maker order is an IOC order": {
			makerOrder: &types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 0, ClobPairId: 0},
				Side:         types.Order_SIDE_SELL,
				Quantums:     100_000_000,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				TimeInForce:  types.Order_TIME_IN_FORCE_IOC,
			},
			takerOrder:    &constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			fillAmount:    100_000_000, // 1 BTC.
			expectedError: errors.New("IOC order cannot be matched as a maker order"),
		},
		"Stateless match validation: taker order is an IOC order": {
			makerOrder: &constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			takerOrder: &types.Order{
				OrderId:      types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 0, ClobPairId: 0},
				Side:         types.Order_SIDE_SELL,
				Quantums:     100_000_000,
				Subticks:     50_000_000_000,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
				TimeInForce:  types.Order_TIME_IN_FORCE_IOC,
			},
			fillAmount: 100_000_000, // 1 BTC.
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			matchedOrder := types.MatchWithOrders{
				MakerOrder: tc.makerOrder,
				TakerOrder: tc.takerOrder,
				FillAmount: satypes.BaseQuantums(tc.fillAmount),
			}
			err := matchedOrder.Validate()
			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
