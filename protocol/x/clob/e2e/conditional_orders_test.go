package clob_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestConditionalOrder(t *testing.T) {
	tests := map[string]struct {
		subaccounts          []satypes.Subaccount
		orders               []clobtypes.Order
		ordersForSecondBlock []clobtypes.Order

		priceUpdateForFirstBlock  *prices.MsgUpdateMarketPrices
		priceUpdateForSecondBlock *prices.MsgUpdateMarketPrices

		expectedInTriggeredStateAfterBlock map[uint32]map[clobtypes.OrderId]bool

		// these expectations are asserted after all blocks are processed
		expectedExistInState    map[clobtypes.OrderId]bool
		expectedOrderFillAmount map[clobtypes.OrderId]uint64
		expectedSubaccounts     []satypes.Subaccount
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
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
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
			},
		},
		"StopLoss/Buy IOC conditional order can place, trigger, partially match, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_400_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(10_000_000_000-12_500_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(25_000_000),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		"StopLoss/Buy IOC conditional order can place, trigger, fully match, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_400_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(10_000_000_000-25_000_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(50_000_000),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		"TakeProfit/Buy post-only conditional order can place, trigger, not cross, and stay in state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id1_Clob0_Sell025BTC_Price50001_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_600_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_10_000,
					},
				},
			},
		},
		"TakeProfit/Buy post-only conditional order can place, trigger, not cross, and partially fill in a later block": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id1_Clob0_Sell025BTC_Price50001_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO,
			},
			ordersForSecondBlock: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_600_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: 25_000_000,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(10_000_000_000-12_500_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(25_000_000),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		"TakeProfit/Buy post-only conditional order can place, trigger, cross, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10_PO,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 4_999_600_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_10_000,
					},
				},
			},
		},
		"Undercollateralized conditional order can be placed": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			priceUpdateForFirstBlock:  &prices.MsgUpdateMarketPrices{},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
			},
		},
		"Undercollateralized conditional order can be placed, trigger, fail collat check and get removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1USD,
				constants.Dave_Num0_10000USD,
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
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_1,
					},
				},
			},
		},
		`Conditional order that would violate isolated subaccount constraints can be placed, trigger, 
		fail isolated subaccount checks and get removed from state`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Dave_Num0_10000USD,
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
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
			},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
							constants.IsoUsd_IsolatedMarket,
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
			}).Build()
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

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[2]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 2)
				}
			}

			// Advance to the next block with new price updates.
			deliverTxsOverride = [][]byte{tApp.GetProposedOperationsTx()}
			txBuilder = encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
			require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForSecondBlock))
			priceUpdateTxBytes, err = encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
			require.NoError(t, err)
			deliverTxsOverride = append(deliverTxsOverride, priceUpdateTxBytes)

			// Place orders for second block
			for _, order := range tc.ordersForSecondBlock {
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

			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[3]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 3)
				}
			}

			// Advance to the next block so that matches are proposed and persisted.
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[4]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 4)
				}
			}

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

			for _, subaccount := range tc.expectedSubaccounts {
				actualSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount.Id)
				require.Equal(t, subaccount, actualSubaccount)
			}
		})
	}
}

func TestConditionalOrder_TriggeringUsingMatchedPrice(t *testing.T) {
	tests := map[string]struct {
		subaccounts          []satypes.Subaccount
		ordersForFirstBlock  []clobtypes.Order
		ordersForSecondBlock []clobtypes.Order

		expectedInTriggeredStateAfterBlock map[uint32]map[clobtypes.OrderId]bool

		// these expectations are asserted after all blocks are processed
		expectedExistInState    map[clobtypes.OrderId]bool
		expectedOrderFillAmount map[clobtypes.OrderId]uint64
		expectedSubaccounts     []satypes.Subaccount
	}{
		"TakeProfit/Buy conditional order is placed and not triggered by matched price": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
			},
		},
		"TakeProfit/Buy conditional order is placed and not triggered by matched price (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Trigger price is $49,700.
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49700,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,500.
				// This price can trigger the conditional order if unbounded.
				// The bounded price is $50,000 - $50,000 * 0.5% = $49,750, which would not trigger the conditional order.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49500_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49995.OrderId: false},
			},
		},
		"StopLoss/Buy conditional order is placed and not triggered by matched price": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
			},
		},
		"StopLoss/Buy conditional order is placed and not triggered by matched price (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Trigger price is $50,300.
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50300,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,500.
				// This price can trigger the conditional order if unbounded.
				// The bounded price is $50,000 + $50,000 * 0.5% = $50,250, which would not trigger the conditional order.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50500_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005.OrderId: false},
			},
		},
		"TakeProfit/Sell conditional order is placed and not triggered by matched price": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
			},
		},
		"TakeProfit/Sell conditional order is placed and not triggered by matched price (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Trigger price is $50,300.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50300,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,500.
				// This price can trigger the conditional order if unbounded.
				// The bounded price is $50,000 + $50,000 * 0.5% = $50,250, which would not trigger the conditional order.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50500_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50005.OrderId: false},
			},
		},
		"StopLoss/Sell conditional order is placed and not triggered by matched price": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
			},
		},
		"StopLoss/Sell conditional order is placed and not triggered by matched price (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Trigger price is $49,700.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49700,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,500.
				// This price can trigger the conditional order if unbounded.
				// The bounded price is $50,000 - $50,000 * 0.5% = $49,750, which would not trigger the conditional order.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49500_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49995.OrderId: false},
			},
		},
		"TakeProfit/Buy conditional order is placed and triggered immediately by matched price": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
			},
		},
		"TakeProfit/Buy conditional order is placed and triggered immediately by matched price (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $49,500.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49500_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
			},
		},
		"TakeProfit/Buy conditional order is placed and triggered in later blocks": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
			},
		},
		"TakeProfit/Buy conditional order is placed and triggered in later blocks (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,500.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49500_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
			},
		},
		"StopLoss/Buy conditional order is placed and triggered immediately by matched price": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
			},
		},
		"StopLoss/Buy conditional order is placed and triggered immediately by matched price (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $50,500.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50500_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
			},
		},
		"StopLoss/Buy conditional order is placed and triggered in later blocks": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
			},
		},
		"StopLoss/Buy conditional order is placed and triggered in later blocks (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,500.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50500_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
			},
		},
		"TakeProfit/Sell conditional order is placed and triggered immediately by matched price": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
			},
		},
		"TakeProfit/Sell conditional order is placed and triggered immediately by matched price (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $50,500.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50500_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
			},
		},
		"TakeProfit/Sell conditional order is placed and triggered in later blocks": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
			},
		},
		"TakeProfit/Sell conditional order is placed and triggered in later blocks (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,500.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50500_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
			},
		},
		"StopLoss/Sell conditional order is placed and triggered immediately by matched price": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
			},
		},
		"StopLoss/Sell conditional order is placed and triggered immediately by matched price (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $49,500.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49500_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
			},
		},
		"StopLoss/Sell conditional order is placed and triggered in later blocks": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
			},
		},
		"StopLoss/Sell conditional order is placed and triggered in later blocks (bounded)": {
			subaccounts: []satypes.Subaccount{
				constants.Bob_Num0_100_000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,500.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49500_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
			},
		},
		"TakeProfit/Buy conditional order is placed, triggered, and partially matched": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
				constants.Dave_Num0_500000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the order that would match against the conditional order.
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true},
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
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the order that would match against the conditional order.
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50001.OrderId: true},
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
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the order that would match against the conditional order.
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_50001.OrderId: true},
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
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the order that would match against the conditional order.
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: true},
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000.OrderId:                          25_000_000,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_49999.OrderId: 25_000_000,
			},
		},
		"StopLoss/Buy IOC conditional order can place, trigger, partially match, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the conditional order and the order that would match against it.
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(10_000_000_000-12_500_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(25_000_000),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		"StopLoss/Buy IOC conditional order can place, trigger, fully match, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the conditional order and the order that would match against it.
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(10_000_000_000-25_000_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(50_000_000),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		"TakeProfit/Buy post-only conditional order can place, trigger, not cross, and stay in state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the conditional order and the order that would match against it.
				constants.LongTermOrder_Dave_Num0_Id1_Clob0_Sell025BTC_Price50001_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_10_000,
					},
				},
			},
		},
		"TakeProfit/Buy post-only conditional order can place, trigger, not cross, and partially fill in a later block": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the conditional order.
				constants.LongTermOrder_Dave_Num0_Id1_Clob0_Sell025BTC_Price50001_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO,
			},
			ordersForSecondBlock: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: 25_000_000,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(10_000_000_000-12_500_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(25_000_000),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		"TakeProfit/Buy post-only conditional order can place, trigger, cross, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $49,997.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				// Place the conditional order.
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10_PO,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_10_000,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
						*genesisState = prices.GenesisState{
							MarketParams: []prices.MarketParam{
								{
									Id:                 0,
									Pair:               constants.BtcUsdPair,
									Exponent:           constants.BtcUsdExponent,
									MinExchanges:       1,
									MinPriceChangePpm:  1_000,
									ExchangeConfigJson: constants.TestMarketExchangeConfigs[exchange_config.MARKET_BTC_USD],
								},
							},

							MarketPrices: []prices.MarketPrice{
								{
									Id:       0,
									Exponent: constants.BtcUsdExponent,
									Price:    constants.FiveBillion, // $50,000 == 1 BTC
								},
							},
						}
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
							{
								Id: 0,
								Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
									PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
										PerpetualId: 0,
									},
								},
								StepBaseQuantums:          1,
								SubticksPerTick:           1,
								QuantumConversionExponent: -8,
								Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
							},
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
			}).Build()
			ctx := tApp.InitChain()

			// Create all orders.
			for _, order := range tc.ordersForFirstBlock {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// First block should persist stateful orders to state.
			for _, order := range tc.ordersForFirstBlock {
				if order.IsStatefulOrder() {
					_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
					require.True(t, found)
				}
			}

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[2]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 2)
				}
			}

			// Place orders for second block
			for _, order := range tc.ordersForSecondBlock {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[3]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 3)
				}
			}

			// Advance to the next block so that matches are proposed and persisted.
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[4]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 4)
				}
			}

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

			for _, subaccount := range tc.expectedSubaccounts {
				actualSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount.Id)
				require.Equal(t, subaccount, actualSubaccount)
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
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
			}).Build()
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
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
			}).Build()
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

func TestConditionalIOCReduceOnlyOrders(t *testing.T) {
	tests := map[string]struct {
		subaccounts                        []satypes.Subaccount
		orders                             []clobtypes.Order
		priceUpdateForFirstBlock           *prices.MsgUpdateMarketPrices
		expectedInTriggeredStateAfterBlock map[uint32]map[clobtypes.OrderId]bool
		expectedExistInState               map[clobtypes.OrderId]bool
		expectedOrderOnMemClob             map[clobtypes.OrderId]bool
		expectedOrderFillAmount            map[clobtypes.OrderId]uint64
		expectedSubaccounts                []satypes.Subaccount
	}{
		"Conditional IOC reduce-only order closes position and gets resized to zero": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_500000USD,
				// Alice has a long position of exactly 0.25 BTC
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(500_000_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(25_000_000), // 0.25 BTC long
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
			orders: []clobtypes.Order{
				// Carl buys 0.25 BTC, which will exactly close Alice's position
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10,
				// Alice tries to sell 0.5 BTC reduce-only but only has 0.25 BTC position
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO,
			},
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000), // Trigger the conditional order
				},
			},

			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: true,
				},
				3: {
					// Should no longer be triggered after removal
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: false,
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				// Should be removed from state after resize to zero
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: false,
			},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: false,
				// Fully filled
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10.OrderId: false,
			},

			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10.OrderId: 25_000_000,
			},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(375_013_750_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(25_000_000),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(624_937_500_000),
						),
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *satypes.GenesisState) {
							genesisState.Subaccounts = tc.subaccounts
						},
					)
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *perptypes.GenesisState) {
							genesisState.Params = constants.PerpetualsGenesisParams
							genesisState.LiquidityTiers = constants.LiquidityTiers
							genesisState.Perpetuals = []perptypes.Perpetual{
								constants.BtcUsd_20PercentInitial_10PercentMaintenance,
								constants.EthUsd_20PercentInitial_10PercentMaintenance,
							}
						},
					)
					return genesis
				}).Build()
			ctx := tApp.InitChain()

			// Create all orders
			deliverTxsOverride := make([][]byte, 0)
			deliverTxsOverride = append(deliverTxsOverride, constants.ValidEmptyMsgProposedOperationsTxBytes)

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

			if tc.priceUpdateForFirstBlock != nil {
				txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
				require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForFirstBlock))
				priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
				require.NoError(t, err)
				deliverTxsOverride = append(deliverTxsOverride, priceUpdateTxBytes)
			}

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			// Verify conditional order triggering for block 2
			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[2]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 2)
				}
			}

			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

			// Verify conditional order triggering for block 3
			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[3]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 3)
				}
			}

			// Verify expectations
			for orderId, exists := range tc.expectedOrderOnMemClob {
				_, existsOnMemclob := tApp.App.ClobKeeper.MemClob.GetOrder(orderId)
				require.Equal(
					t,
					exists,
					existsOnMemclob,
					"Order %v expected on memclob: %v, actual: %v",
					orderId,
					exists,
					existsOnMemclob,
				)
			}

			for orderId, expectedFillAmount := range tc.expectedOrderFillAmount {
				exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				if expectedFillAmount > 0 {
					require.True(t, exists)
					require.Equal(t, expectedFillAmount, fillAmount.ToUint64())
				}
			}

			// Verify orders are removed from state
			for orderId, expectedExists := range tc.expectedExistInState {
				_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, orderId)
				require.Equal(
					t,
					expectedExists,
					found,
					"Order %v expected exists in state: %v, actual: %v",
					orderId,
					expectedExists,
					found,
				)
			}

			for _, subaccount := range tc.expectedSubaccounts {
				actualSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount.Id)
				require.Equal(t, subaccount, actualSubaccount)
			}
		})
	}
}
