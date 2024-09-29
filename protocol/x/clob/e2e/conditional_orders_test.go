package clob_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"

	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	feetiertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	prices "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"

	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestConditionalOrder(t *testing.T) {
	tests := map[string]struct {
		subaccounts          []satypes.Subaccount
		orders               []clobtypes.Order
		ordersForSecondBlock []clobtypes.Order

		priceUpdateForFirstBlock  map[uint32]ve.VEPricePair
		priceUpdateForSecondBlock map[uint32]ve.VEPricePair

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
			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},

			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_000_000,
					PnlPrice:  5_000_000_000,
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},
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
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(10_000_000_000 - 12_500_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(25_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"TakeProfit/Sell FOK conditional order can place, trigger, not match, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK,
			},
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_10_000,
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"TakeProfit/Sell FOK conditional order can place, trigger, fully match, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_PO,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK,
			},
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(10_000_000_000 + 25_000_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(-50_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},
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
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(10_000_000_000 - 25_000_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(50_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_600_000,
					PnlPrice:  4_999_600_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},
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
						&constants.TDai_Asset_10_000,
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_600_000,
					PnlPrice:  4_999_600_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},
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
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(10_000_000_000 - 12_500_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(25_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"TakeProfit/Buy post-only conditional order can place, trigger, cross, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_TP_49999_PO,
			},
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_600_000,
					PnlPrice:  4_999_600_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},
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
						&constants.TDai_Asset_10_000,
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
						&constants.TDai_Asset_1,
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
		`Conditional FOK order that would violate isolated subaccount constraints can be placed, trigger, 
		fail isolated subaccount checks and get removed from state`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Dave_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_FOK,
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
				},
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_FOK.OrderId: false,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_FOK.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_FOK.OrderId: true},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999_FOK.OrderId: false},
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

			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)

			require.NoError(t, conversionErr)

			tApp.App.RatelimitKeeper.SetSDAIPrice(tApp.App.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.App.RatelimitKeeper.SetAssetYieldIndex(tApp.App.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.CrashingApp.RatelimitKeeper.SetSDAIPrice(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.CrashingApp.RatelimitKeeper.SetAssetYieldIndex(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.NoCheckTxApp.RatelimitKeeper.SetSDAIPrice(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.NoCheckTxApp.RatelimitKeeper.SetAssetYieldIndex(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.ParallelApp.RatelimitKeeper.SetSDAIPrice(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.ParallelApp.RatelimitKeeper.SetAssetYieldIndex(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

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

			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdateForFirstBlock,
				"",
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			// ve info has to be first in block
			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

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

			extCommitInfo, _, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdateForSecondBlock,
				"",
				tApp.GetHeader().Height-1, // prepare proposal gets called for block 2
			)
			require.NoError(t, err)

			// Advance to the next block with new price updates.
			deliverTxsOverride = [][]byte{tApp.GetProposedOperationsTx(
				extCommitInfo,
			)}

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

			_, extCommitBz, err = vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdateForSecondBlock,
				"",
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

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
				// Trigger price is $48,700.
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_48700,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $48,500.
				// This price can trigger the conditional order if unbounded.
				// The bounded price is $50,000 - $50,000 * 2.5% = $48,750, which would not trigger the conditional order.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price48500_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_48700.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_48700.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_48700.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_48700.OrderId: false},
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
				// Trigger price is $51,300.
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_51300,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $51,500.
				// This price can trigger the conditional order if unbounded.
				// The bounded price is $50,000 + $50,000 * 2.5% = $51,250, which would not trigger the conditional order.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price51500_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_51300.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_51300.OrderId: false},
				3: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_51300.OrderId: false},
				4: {constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_51300.OrderId: false},
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
				// Trigger price is $51,300.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_51300,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $51,500.
				// This price can trigger the conditional order if unbounded.
				// The bounded price is $50,000 + $50,000 * 2.5% = $51,250, which would not trigger the conditional order.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price51500_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_51300.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_51300.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_51300.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_TP_51300.OrderId: false},
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
				// Trigger price is $48,700.
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_48700,
			},
			ordersForSecondBlock: []clobtypes.Order{
				// Create a match with price $48,500.
				// This price can trigger the conditional order if unbounded.
				// The bounded price is $50,000 - $50,000 * 2.5% = $48,750, which would not trigger the conditional order.
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price48500_GTB10,
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
			},

			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_48700.OrderId: true,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_48700.OrderId: false},
				3: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_48700.OrderId: false},
				4: {constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_48700.OrderId: false},
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
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(10_000_000_000 - 12_500_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(25_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"TakeProfit/Sell FOK conditional order can place, trigger, not match, and be removed from state": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Carl_Num1_100000USD,
				constants.Dave_Num1_500000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				// Create a match with price $50,003.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50003_GTB10,
				constants.Order_Dave_Num1_Id0_Clob0_Sell1BTC_Price49997_GTB10,
				// Place the conditional order.
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_10_000,
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"TakeProfit/Sell FOK conditional order can place, trigger, fully match, and be removed from state": {
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
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_PO,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK,
			},
			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: true},
				3: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false},
				4: {constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false},
			},
			expectedExistInState: map[clobtypes.OrderId]bool{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell05BTC_Price50000_GTBT10_TP_50003_FOK.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(10_000_000_000 + 25_000_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(-50_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(10_000_000_000 - 25_000_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(50_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
						&constants.TDai_Asset_10_000,
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(10_000_000_000 - 12_500_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(25_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
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
						&constants.TDai_Asset_10_000,
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
									Id:        0,
									Exponent:  constants.BtcUsdExponent,
									SpotPrice: constants.FiveBillion, // $50,000 == 1 BTC
									PnlPrice:  constants.FiveBillion, // $50,000 == 1 BTC
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

			rate := sdaiservertypes.TestSDAIEventRequest.ConversionRate

			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				map[uint32]ve.VEPricePair{},
				rate,
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: [][]byte{extCommitBz},
			})

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

			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

			// First block should persist stateful orders to state.
			for _, order := range tc.ordersForFirstBlock {
				if order.IsStatefulOrder() {
					_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
					require.True(t, found)
				}
			}

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[2]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 3)
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

			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[3]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 4)
				}
			}

			// Advance to the next block so that matches are proposed and persisted.
			ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{})
			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[4]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 5)
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

		priceUpdate                       map[uint32]ve.VEPricePair
		expectedConditionalOrderTriggered bool
	}{
		"untriggered conditional order is cancelled": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_SL_50005,
			},

			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
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

			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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

			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
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

			rate := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				map[uint32]ve.VEPricePair{},
				rate,
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: [][]byte{extCommitBz},
			})

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

			// // Add the price update.
			_, extCommitBz, err = vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdate,
				"",
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

			// Advance to the next block, updating the price.
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
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
		firstPriceUpdate              map[uint32]ve.VEPricePair
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
			firstPriceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			firstPriceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			firstPriceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			firstPriceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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
			firstPriceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 4_999_700_000,
					PnlPrice:  4_999_700_000,
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

			rate := sdaiservertypes.TestSDAIEventRequest.ConversionRate

			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				map[uint32]ve.VEPricePair{},
				rate,
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: [][]byte{extCommitBz},
			})

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
			_, extCommitBz, err = vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.firstPriceUpdate,
				"",
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)
			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

			// Advance to the next block, updating the price.
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
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
