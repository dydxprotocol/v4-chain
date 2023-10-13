package clob_test

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestConditionalOrder(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order

		priceUpdateForFirstBlock  *prices.MsgUpdateMarketPrices
		priceUpdateForSecondBlock *prices.MsgUpdateMarketPrices

		expectedExistInState    map[clobtypes.OrderId]bool
		expectedTriggered       map[clobtypes.OrderId]bool
		expectedOrderFillAmount map[clobtypes.OrderId]uint64
	}{
		"TakeProfit/Buy conditional order is placed but not triggered (no price update)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			priceUpdateForFirstBlock:  &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false,
			},
		},
		"StopLoss/Buy conditional order is placed but not triggered (no price update)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			},
			priceUpdateForFirstBlock:  &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false,
			},
		},
		"TakeProfit/Sell conditional order is placed but not triggered (no price update)": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},
			priceUpdateForFirstBlock:  &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false,
			},
		},
		"StopLoss/Sell conditional order is placed but not triggered (no price update)": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
			},
			priceUpdateForFirstBlock:  &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false,
			},
		},
		"TakeProfit/Buy conditional order is placed and not triggered by price update": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false,
			},
		},
		"StopLoss/Buy conditional order is placed and not triggered by price update": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false,
			},
		},
		"TakeProfit/Sell conditional order is placed and not triggered by price update": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false,
			},
		},
		"StopLoss/Sell conditional order is placed and not triggered by price update": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false,
			},
		},
		"TakeProfit/Buy conditional order is placed and triggered": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
		},
		"TakeProfit/Buy conditional order is placed and triggered in later blocks": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
		},
		"StopLoss/Buy conditional order is placed and triggered": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
		},
		"StopLoss/Buy conditional order is placed and triggered in later blocks": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
		},
		"TakeProfit/Sell conditional order is placed and triggered": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
		},
		"TakeProfit/Sell conditional order is placed and triggered in later blocks": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
		},
		"StopLoss/Sell conditional order is placed and triggered": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
		},
		"StopLoss/Sell conditional order is placed and triggered in later blocks": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
		},
		"TakeProfit/Buy conditional order is placed, triggered, and partially matched": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId:                    25_000_000,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: 25_000_000,
			},
		},
		"StopLoss/Buy conditional order is placed, triggered, and partially matched": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId:                    25_000_000,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: 25_000_000,
			},
		},
		"TakeProfit/Sell conditional order is placed, triggered, and partially matched": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num0_100000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000.OrderId:                          25_000_000,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: 25_000_000,
			},
		},
		"StopLoss/Sell conditional order is placed, triggered, and partially matched": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num0_100000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000.OrderId:                          25_000_000,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: 25_000_000,
			},
		},
		"Triggered conditional orders can not be untriggered": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_000_000),
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						*genesisState = constants.TestPricesGenesisState
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{
							constants.BtcUsd_20PercentInitial_10PercentMaintenance,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{
							constants.ClobPair_Btc,
						}
						genesisState.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			// Create all orders.
			deliverTxsOverride := make([][]byte, 0)
			deliverTxsOverride = append(
				deliverTxsOverride,
				constants.ValidEmptyMsgProposedOperationsTxBytes,
			)

			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

					if order.IsStatefulOrder() {
						deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
					}
				}
			}

			// Add an empty premium vote.
			deliverTxsOverride = append(deliverTxsOverride, constants.EmptyMsgAddPremiumVotesTxBytes)

			// Add the price update.
			txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
			require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForFirstBlock))
			priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
			require.NoError(t, err)

			deliverTxsOverride = append(deliverTxsOverride, priceUpdateTxBytes)

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			// First block should persist stateful orders to state.
			for _, order := range tc.orders {
				if order.IsStatefulOrder() {
					_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
					require.True(t, found)
				}
			}

			// Advance to the next block with new price updates.
			txBuilder = encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
			require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForSecondBlock))
			priceUpdateTxBytes, err = encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
			require.NoError(t, err)

			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: [][]byte{priceUpdateTxBytes},
			})

			// Make sure conditional orders are triggered.
			for orderId, triggered := range tc.expectedTriggered {
				require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId))
			}

			// Advance to the next block so that matches are proposed and persisted.
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})

			// Verify expectations.
			for orderId, exists := range tc.expectedExistInState {
				_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, orderId)
				require.Equal(t, exists, found)
			}

			for orderId, expectedFillAmount := range tc.expectedOrderFillAmount {
				exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				require.True(t, exists)
				require.Equal(t, expectedFillAmount, fillAmount.ToUint64())
			}
		})
	}
}

func TestConditionalOrderCancellation(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order

		priceUpdate                       *prices.MsgUpdateMarketPrices
		expectedConditionalOrderTriggered bool
	}{
		"untriggered conditional order is cancelled": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005,
			},

			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},
			expectedConditionalOrderTriggered: false,
		},
		"triggered conditional order is cancelled": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},

			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			expectedConditionalOrderTriggered: true,
		},
		"triggered and matched conditional order is cancelled": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},

			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},
			expectedConditionalOrderTriggered: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						*genesisState = constants.TestPricesGenesisState
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{
							constants.BtcUsd_20PercentInitial_10PercentMaintenance,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{
							constants.ClobPair_Btc,
						}
						genesisState.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			// Create all orders.
			deliverTxsOverride := make([][]byte, 0)
			deliverTxsOverride = append(
				deliverTxsOverride,
				constants.ValidEmptyMsgProposedOperationsTxBytes,
			)

			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

					if order.IsStatefulOrder() {
						deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
					}
				}
			}

			// Add an empty premium vote.
			deliverTxsOverride = append(deliverTxsOverride, constants.EmptyMsgAddPremiumVotesTxBytes)

			// Add the price update.
			txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
			require.NoError(t, txBuilder.SetMsgs(tc.priceUpdate))
			priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
			require.NoError(t, err)

			deliverTxsOverride = append(deliverTxsOverride, priceUpdateTxBytes)

			// Advance to the next block, updating the price.
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			// Verify placed conditional order exists
			_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.orders[0].OrderId)
			require.Equal(t, true, found)

			// Verify placed conditional order was triggered
			isTriggered := tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, tc.orders[0].OrderId)
			require.Equal(t, tc.expectedConditionalOrderTriggered, isTriggered)

			// If there was a short term order, assert order fill amount was created.
			if len(tc.orders) > 1 {
				exists, _, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, tc.orders[0].OrderId)
				require.Equal(t, true, exists)
				exists, _, _ = tApp.App.ClobKeeper.GetOrderFillAmount(ctx, tc.orders[0].OrderId)
				require.Equal(t, true, exists)
			}

			// Cancel the previously placed conditional order.
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				tApp.App,
				*clobtypes.NewMsgCancelOrderStateful(
					tc.orders[0].OrderId,
					lib.MustConvertIntegerToUint32(
						time.Unix(ctx.BlockTime().Unix(), 0).Add(clobtypes.StatefulOrderTimeWindow).Unix(),
					),
				),
			) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			// Advance to the next block, cancelling the previously placed conditional order.
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})

			// Verify conditional order cancellation cleared state.
			exists, _, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, tc.orders[0].OrderId)
			require.Equal(t, false, exists)

			_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.orders[0].OrderId)
			require.Equal(t, false, found)

			isTriggered = tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, tc.orders[0].OrderId)
			require.Equal(t, false, isTriggered)
		})
	}
}

func TestConditionalOrderExpiration(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order

		firstBlockTime                time.Time
		firstPriceUpdate              *prices.MsgUpdateMarketPrices
		firstExistInStateExpectations map[clobtypes.OrderId]bool
		firstExpectedTriggered        map[clobtypes.OrderId]bool

		secondBlockTime                time.Time
		secondExistInStateExpectations map[clobtypes.OrderId]bool
		secondExpectedTriggered        map[clobtypes.OrderId]bool
		expectedOrderFillsExist        map[clobtypes.OrderId]bool
	}{
		"untriggered conditional order that doesn't expire": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			// Expires at unix time 10.
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995,
			},
			firstBlockTime: time.Unix(5, 0).UTC(),
			// Does not trigger above conditional order.
			firstPriceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			firstExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
			},
			firstExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false,
			},

			// Does not expire above conditional order.
			secondBlockTime: time.Unix(9, 0).UTC(),
			secondExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
			},
			secondExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false,
			},
		},
		"untriggered conditional order that expires": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
			},
			// Expires at unix time 10.
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995,
			},
			firstBlockTime: time.Unix(5, 0).UTC(),
			// Does not trigger above conditional order.
			firstPriceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			firstExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
			},
			firstExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false,
			},

			// Expires above conditional order.
			secondBlockTime: time.Unix(11, 0).UTC(),
			secondExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false,
			},
			secondExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false,
			},
		},
		"triggered conditional order that doesn't expire": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			// Expires at unix time 10.
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			firstBlockTime: time.Unix(5, 0).UTC(),
			// Triggers above conditional order.
			firstPriceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			firstExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			firstExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},

			// Does not expire above conditional order.
			secondBlockTime: time.Unix(9, 0).UTC(),
			secondExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			secondExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
		},
		"triggered conditional order that expires": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			// Expires at unix time 10.
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			firstBlockTime: time.Unix(5, 0).UTC(),
			// Triggers above conditional order.
			firstPriceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			firstExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			firstExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},

			// Expires above conditional order.
			secondBlockTime: time.Unix(11, 0).UTC(),
			secondExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false,
			},
			secondExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false,
			},
		},
		"triggered conditional order partially matches, then expires": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num0_100000USD,
			},
			orders: []clobtypes.Order{
				// Expires at unix time 10.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
			},
			firstBlockTime: time.Unix(5, 0).UTC(),
			// Triggers above conditional order.
			firstPriceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			firstExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			firstExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},

			// Expires above conditional order.
			secondBlockTime: time.Unix(11, 0).UTC(),
			secondExistInStateExpectations: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false,
			},
			secondExpectedTriggered: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false,
			},
			expectedOrderFillsExist: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000.OrderId:                          true,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *prices.GenesisState) {
						*genesisState = constants.TestPricesGenesisState
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = []perptypes.Perpetual{
							constants.BtcUsd_20PercentInitial_10PercentMaintenance,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{
							constants.ClobPair_Btc,
						}
						genesisState.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			// Create all orders.
			deliverTxsOverride := make([][]byte, 0)
			deliverTxsOverride = append(
				deliverTxsOverride,
				constants.ValidEmptyMsgProposedOperationsTxBytes,
			)

			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
					if order.IsStatefulOrder() {
						deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
					}
				}
			}

			// Add an empty premium vote.
			deliverTxsOverride = append(deliverTxsOverride, constants.EmptyMsgAddPremiumVotesTxBytes)

			// Add the price update.
			txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
			require.NoError(t, txBuilder.SetMsgs(tc.firstPriceUpdate))
			priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
			require.NoError(t, err)

			deliverTxsOverride = append(deliverTxsOverride, priceUpdateTxBytes)

			// Advance to the next block, updating the price.
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				BlockTime:          tc.firstBlockTime,
				DeliverTxsOverride: deliverTxsOverride,
			})

			// Verify first test expectations of state.
			for orderId, exists := range tc.firstExistInStateExpectations {
				_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, orderId)
				require.Equal(t, exists, found)
			}
			// Make sure conditional orders are triggered.
			for orderId, triggered := range tc.firstExpectedTriggered {
				require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId))
			}

			// Advance to the next block, expiring the order if the test case calls for it.
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
				BlockTime: tc.secondBlockTime,
			})
			// Verify fill amounts gets pruned for expired matched conditional orders.
			for orderId, expectedExists := range tc.expectedOrderFillsExist {
				exists, _, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				require.Equal(t, expectedExists, exists)
			}

			// Verify second test expectations of state.
			for orderId, exists := range tc.secondExistInStateExpectations {
				_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, orderId)
				require.Equal(t, exists, found)
			}
			// Make sure conditional orders are triggered.
			for orderId, triggered := range tc.secondExpectedTriggered {
				require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId))
			}
		})
	}
}
