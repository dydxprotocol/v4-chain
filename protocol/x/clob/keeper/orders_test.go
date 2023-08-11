package keeper_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	indexer_manager "github.com/dydxprotocol/v4/indexer/indexer_manager"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/mocks"
	clobtest "github.com/dydxprotocol/v4/testutil/clob"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/testutil/proto"
	"github.com/dydxprotocol/v4/testutil/tracer"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/dydxprotocol/v4/x/perpetuals"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	"github.com/dydxprotocol/v4/x/prices"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs          []types.ClobPair
		existingOrders []types.Order

		// Parameters.
		order types.Order

		// Expectations.
		expectedMultiStoreWrites []string
		expectedOrderStatus      types.OrderStatus
		expectedFilledSize       satypes.BaseQuantums
		expectedErr              error
		expectedSeenPlaceOrder   bool
	}{
		"Can place an order on the orderbook closing a position": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},

			order: constants.Order_Carl_Num0_Id1_Clob0_Buy1BTC_Price49999,

			expectedOrderStatus:    types.Success,
			expectedFilledSize:     0,
			expectedSeenPlaceOrder: true,
		},
		"Can place an order on the orderbook in a different market than their current perpetual position": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},

			order: constants.Order_Carl_Num0_Id2_Clob1_Buy10ETH_Price3000,

			expectedOrderStatus:    types.Success,
			expectedFilledSize:     0,
			expectedSeenPlaceOrder: true,
		},
		"Can place an order and the order is fully matched": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			existingOrders: []types.Order{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
			},

			order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,

			expectedOrderStatus:    types.Success,
			expectedFilledSize:     constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.GetBaseQuantums(),
			expectedSeenPlaceOrder: true,
			expectedMultiStoreWrites: []string{
				// Update taker subaccount
				"Subaccount/value/\n+" + constants.Carl_Num0.Owner + "/",
				// Indexer event
				indexer_manager.IndexerEventsKey,
				// Update maker subaccount
				"Subaccount/value/\n+" + constants.Dave_Num0.Owner + "/",
				// Indexer event
				indexer_manager.IndexerEventsKey,
				// Update prunable block height for taker fill amount
				"BlockHeightToPotentiallyPrunableOrders/value",
				// Update taker order fill amount
				"OrderAmount/value",
				// Update taker order fill amount in memStore
				"OrderAmount/value",
				// Update prunable block height for maker fill amount
				"BlockHeightToPotentiallyPrunableOrders/value",
				// Update maker order fill amount
				"OrderAmount/value",
				// Update maker order fill amount in memStore
				"OrderAmount/value",
			},
		},
		"Cannot place an order on the orderbook if the account would be undercollateralized": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},

			order: constants.Order_Carl_Num0_Id3_Clob1_Buy1ETH_Price3000,

			expectedOrderStatus:    types.Undercollateralized,
			expectedFilledSize:     0,
			expectedSeenPlaceOrder: true,
		},
		"Can place an order on the orderbook if the subaccount is right at the initial margin ratio": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},

			order: constants.Order_Carl_Num0_Id0_Clob0_Buy10QtBTC_Price10000QuoteQt,

			expectedOrderStatus:    types.Success,
			expectedFilledSize:     0,
			expectedSeenPlaceOrder: true,
		},
		"Cannot place an order on the orderbook if the account would be undercollateralized due to fees paid": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{
				// Exact same set-up as the previous test, except the clob pair has fees.
				constants.ClobPair_Btc,
			},

			order: constants.Order_Carl_Num0_Id0_Clob0_Buy10QtBTC_Price10000QuoteQt,

			expectedOrderStatus:    types.Undercollateralized,
			expectedFilledSize:     0,
			expectedSeenPlaceOrder: true,
		},
		"Cannot open an order if it doesn't reference a valid CLOB": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			order: constants.Order_Carl_Num0_Id3_Clob1_Buy1ETH_Price3000,

			expectedErr: types.ErrInvalidClob,
		},
		"Cannot open an order if the subticks are invalid": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num1, ClientId: 1,
				},
				Subticks: 2,
			},

			expectedErr: types.ErrInvalidPlaceOrder,
		},
		"Cannot open an order that is smaller than the minimum base quantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num1, ClientId: 1,
				},
				Quantums: 1,
			},

			expectedErr:            types.ErrInvalidPlaceOrder,
			expectedSeenPlaceOrder: false,
		},
		"Cannot open an order that is not divisible by the step size base quantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			order: types.Order{
				OrderId:  types.OrderId{},
				Quantums: 11,
			},

			expectedErr: types.ErrInvalidPlaceOrder,
		},
		"Cannot open an order with a GoodTilBlock in the past": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			order: types.Order{
				OrderId:      types.OrderId{},
				Quantums:     10,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 0},
			},

			expectedErr: types.ErrHeightExceedsGoodTilBlock,
		},
		"Cannot open an order with a GoodTilBlock greater than ShortBlockWindow blocks in the future": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			order: types.Order{
				OrderId:      types.OrderId{},
				Quantums:     10,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 2 + types.ShortBlockWindow},
			},

			expectedErr: types.ErrGoodTilBlockExceedsShortBlockWindow,
		},
		// This is a regression test for an issue whereby orders that had been previously matched were being checked for
		// collateralization as if the subticks of the order were `0`. This resulted in always using `0`
		// `bigFillQuoteQuantums` for the order when performing collateralization checks during `PlaceOrder`.
		// This meant that previous buy orders in the match queue could only ever increase collateralization
		// of the subaccount.
		// Context: https://dydx-team.slack.com/archives/C03SLFHC3L7/p1668105457456389
		`Regression: New order should be undercollateralized when adding to the orderbook when previous fills make it
			undercollateralized`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num1_500USD,
				constants.Carl_Num0_10000USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders: []types.Order{
				// The maker subaccount places an order which is a maker order to buy $500 worth of BTC.
				// The subaccount has a balance of $500 worth of USDC, and the perpetual has a 100% margin requirement.
				// This order does not match, and is placed on the book as a maker order.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1kQtBTC_Price50000,
				// The taker subaccount places an order which fully fills the previous order.
				constants.Order_Carl_Num0_Id0_Clob0_Sell1kQtBTC_Price50000,
			},
			// The maker subaccount places a second order identical to the first.
			// This should fail, because the maker subaccount currently has a balance of $0 USDC, and a perpetual of size
			// 0.01 BTC ($500), and the perpetual has a 100% margin requirement.
			order:                  constants.Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000,
			expectedOrderStatus:    types.Undercollateralized,
			expectedSeenPlaceOrder: true,
		},
		`Regression: New order should be undercollateralized when matching when previous fills make it
				undercollateralized`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num1_500USD,
				constants.Carl_Num0_10000USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders: []types.Order{
				// The maker subaccount places an order which is a maker order to buy $500 worth of BTC.
				// The subaccount has a balance of $500 worth of USDC, and the perpetual has a 100% margin requirement.
				// This order does not match, and is placed on the book as a maker order.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1kQtBTC_Price50000,
				// The taker subaccount places an order which fully fills the previous order.
				constants.Order_Carl_Num0_Id0_Clob0_Sell1kQtBTC_Price50000,
				// Match queue is now empty.
				// The subaccount from the above order now places an order which is added to the book.
				constants.Order_Carl_Num0_Id1_Clob0_Sell1kQtBTC_Price50000,
			},
			// The maker subaccount places a second order identical to the first.
			// This should fail, because the maker during matching, because subaccount currently has a balance of $0 USDC,
			// and a perpetual of size 0.01 BTC ($500), and the perpetual has a 100% margin requirement.
			order:                  constants.Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000,
			expectedOrderStatus:    types.Undercollateralized,
			expectedSeenPlaceOrder: true,
		},
		`New order should be undercollateralized when matching when previous fills make it undercollateralized when using
				maker orders subticks, but would be collateralized if using taker order subticks`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num1_500USD,
				constants.Carl_Num0_10000USD,
			},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders: []types.Order{
				// Alice places sell order for $500 worth of BTC.
				constants.Order_Carl_Num0_Id0_Clob0_Sell1kQtBTC_Price50000,
				// Bob places a buy order to buy $500 worth of BTC at $600.
				// This completely fills Alice's order at the maker price of $500.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1kQtBTC_Price60000,
				// Orderbook is now empty.
				// Alice places another sell order for $500 worth of BTC which now rests on the book.
				constants.Order_Carl_Num0_Id1_Clob0_Sell1kQtBTC_Price50000,
			},
			// Bob places a second order at a lower price than his first.
			// This should succeed, because Bob currently has a balance of $0 USDC,
			// and a perpetual of size 0.01 BTC ($500), and the perpetual has a 50% margin requirement.
			// This should bring Bob's balance to -500$ USD, and 0.02 BTC which is exactly at the 50% margin requirement.
			// If the Bob's previous order was viewed as being filled at the taker subticks (60_000_000_000) instead when
			// checking collateralization, the match would fail.
			order:                  constants.Order_Carl_Num1_Id1_Clob0_Buy1kQtBTC_Price50000,
			expectedOrderStatus:    types.Success,
			expectedFilledSize:     1_000_000,
			expectedSeenPlaceOrder: true,
			expectedMultiStoreWrites: []string{
				// Update taker subaccount
				"Subaccount/value/\n+" + constants.Carl_Num1.Owner,
				indexer_manager.IndexerEventsKey,
				// Update maker subaccount
				"Subaccount/value/\n+" + constants.Carl_Num0.Owner,
				indexer_manager.IndexerEventsKey,
				// Update prunable block height for taker fill amount
				"BlockHeightToPotentiallyPrunableOrders/value",
				// Update taker order fill amount
				"OrderAmount/value",
				// Update taker order fill amount in memStore
				"OrderAmount/value",
				// Update prunable block height for maker fill amount
				"BlockHeightToPotentiallyPrunableOrders/value",
				// Update maker order fill amount
				"OrderAmount/value",
				// Update maker order fill amount in memStore
				"OrderAmount/value",
			},
		},
		// This is a regression test for an issue whereby orders that had been previously matched were being checked for
		// collateralization as if the CLOB pair ID of the order was `0`. This resulted in always using `0`
		// the quantum conversion exponent of the first CLOB when performing collateralization checks during `PlaceOrder`.
		// This lead to an incorrect amount of quote quantums being returned for all previously matched orders
		// that weren't placed on the first CLOB. If firstClobPair.QuantumConversionExponent >
		// expectedClobPair.QuantumConversionExponent, then sellers receive more quote quantums and buyers are charged
		// more. Vice versa if firstClobPair.QuantumConversionExponent < expectedClobPair.QuantumConversionExponent.
		// Context: https://github.com/dydxprotocol/v4/pull/562#discussion_r1024319468
		`Regression: New order should be fully collateralized when matching with previous fills
				because the correct quantum conversion exponent was used`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_660USD,
				constants.Dave_Num0_10000USD,
			},
			clobs: []types.ClobPair{
				{
					Id: 0,
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:                    types.ClobPair_STATUS_ACTIVE,
					StepBaseQuantums:          10,
					SubticksPerTick:           100,
					MinOrderBaseQuantums:      10,
					QuantumConversionExponent: -1,
				},
				constants.ClobPair_Eth_No_Fee,
			},
			existingOrders: []types.Order{
				// Alice places buy order for $3,000 worth of ETH.
				constants.Order_Carl_Num0_Id3_Clob1_Buy1ETH_Price3000,
				// Bob places sell order for $3,000 worth of ETH.
				// This completely fills Alice's order at the maker price of $3,000.
				constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000,
			},
			// Bob places a second buy order at a the same price as the first.
			// This should bring his total initial margin requirement to $3,300 * 20%, or $660
			// which is exactly equal to the quote balance.
			// However, if the quantum conversion exponent of the first perpetual CLOB pair
			// is used for the collateralization check, it will over-report the amount of quote
			// quantums required to open the previous buy order and fail.
			order:                  constants.Order_Carl_Num0_Id4_Clob1_Buy01ETH_Price3000,
			expectedOrderStatus:    types.Success,
			expectedSeenPlaceOrder: true,
		},
		`Subaccount cannot place maker buy order for 1 BTC at 5 subticks with 0 collateral`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_0USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders:           []types.Order{},
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price5subticks_GTB10,
			expectedOrderStatus:      types.Undercollateralized,
			expectedFilledSize:       0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{},
		},
		`Subaccount cannot place maker sell order for 1 BTC at 500,000 with 0 collateral`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_0USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders:           []types.Order{},
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price500000_GTB10,
			expectedOrderStatus:      types.Undercollateralized,
			expectedFilledSize:       0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{},
		},
		// <grouped tests: pessimistic value collateralization check -- BUY>
		// The following 3 tests are a group to test the pessimistic value used for the collateralization check.
		// 1. The first should have a lower asset value in its subaccount. (undercollateralized)
		// 2. The second should have a buy price above the oracle price of 50,000. (undercollateralized)
		// 3. The third should have the order in common with #1 and the subaccount in common with #2 and should succeed.
		`Subaccount cannot place buy order due to a failed collateralization check with the oracle price but would
				pass if using the maker price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_10000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders:           []types.Order{},
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price5subticks_GTB10,
			expectedOrderStatus:      types.Undercollateralized,
			expectedFilledSize:       0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{},
		},
		`Subaccount cannot place buy order due to a failed collateralization check with its maker price but would
				pass if using the oracle price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_100000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders:           []types.Order{},
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10,
			expectedOrderStatus:      types.Undercollateralized,
			expectedFilledSize:       0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{},
		},
		`Subaccount placed buy order passes collateralization check when using the oracle price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_100000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders:           []types.Order{},
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price5subticks_GTB10,
			expectedOrderStatus:      types.Success,
			expectedFilledSize:       0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{},
		},
		// <end of grouped tests: pessimistic value collateralization check -- BUY>
		// <grouped tests: pessimistic value collateralization check -- SELL>
		// The following 3 tests are a group to test the pessimistic value used for the collateralization check.
		// 1. The first should have a lower asset value in its subaccount. (undercollateralized)
		// 2. The second should have a sell price below the oracle price of 50,000 subticks. (undercollateralized)
		// 3. The third should have the order in common with #1 and the subaccount in common with #2 and should succeed.
		`Subaccount cannot place sell order due to a failed collateralization check with the oracle price but would
				pass if using the maker price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_10000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders:           []types.Order{},
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price500000_GTB10,
			expectedOrderStatus:      types.Undercollateralized,
			expectedFilledSize:       0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{},
		},
		`Subaccount cannot place sell order due to a failed collateralization check with its maker price but would
				pass if using the oracle price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_50000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders:           []types.Order{},
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price5000_GTB10,
			expectedOrderStatus:      types.Undercollateralized,
			expectedFilledSize:       0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{},
		},
		`Subaccount placed sell order passes collateralization check when using the oracle price`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_50PercentInitial_40PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{constants.Carl_Num0_50000USD},
			clobs: []types.ClobPair{
				constants.ClobPair_Btc_No_Fee,
			},
			existingOrders:           []types.Order{},
			order:                    constants.Order_Carl_Num0_Id0_Clob0_Sell1BTC_Price500000_GTB10,
			expectedOrderStatus:      types.Success,
			expectedFilledSize:       0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{},
		},
		// <end of grouped tests: pessimistic value collateralization check -- SELL>
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)

			ctx,
				clobKeeper,
				pricesKeeper,
				assetsKeeper,
				perpetualsKeeper,
				subaccountsKeeper,
				_,
				_ := keepertest.ClobKeepers(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ctx = ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, pricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up USDC asset in assets module.
			err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := perpetualsKeeper.CreatePerpetual(
					ctx,
					p.Ticker,
					p.MarketId,
					p.AtomicResolution,
					p.DefaultFundingPpm,
					p.LiquidityTier,
				)
				require.NoError(t, err)
			}

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				subaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobs {
				_, err = clobKeeper.CreatePerpetualClobPair(
					ctx,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
					clobPair.MakerFeePpm,
					clobPair.TakerFeePpm,
				)
				require.NoError(t, err)
			}

			// Create all existing orders.
			for _, order := range tc.existingOrders {
				_, _, err := clobKeeper.CheckTxPlaceOrder(ctx, &types.MsgPlaceOrder{Order: order})
				require.NoError(t, err)
			}

			// Run the test.
			traceDecoder := &tracer.TraceDecoder{}
			ctx.MultiStore().SetTracer(traceDecoder)

			msg := &types.MsgPlaceOrder{Order: tc.order}
			orderSizeOptimisticallyFilledFromMatching,
				orderStatus,
				err := clobKeeper.CheckTxPlaceOrder(ctx, msg)

			// Verify test expectations.
			require.ErrorIs(t, err, tc.expectedErr)
			if err == nil {
				require.Equal(t, tc.expectedOrderStatus, orderStatus)
				require.Equal(t, tc.expectedFilledSize, orderSizeOptimisticallyFilledFromMatching)
			}

			traceDecoder.RequireKeyPrefixWrittenInSequence(t, tc.expectedMultiStoreWrites)

			require.Equal(
				t,
				tc.expectedSeenPlaceOrder,
				clobKeeper.HasSeenPlaceOrder(ctx, *msg),
			)
		})
	}
}

// TODO(DEC-1648): Add test cases for additional order placement scenarios.
func TestPlaceOrder_LongTerm(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs          []types.ClobPair
		existingOrders []types.Order

		// Parameters.
		order types.Order

		// Expectations.
		expectedMultiStoreWrites []string
		expectedOrderStatus      types.OrderStatus
		expectedFilledSize       satypes.BaseQuantums
		expectedErr              error
		expectedSeenPlaceOrder   bool
		expectedTransactionIndex uint32
	}{
		"Can place a stateful order on the orderbook": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},

			order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,

			expectedOrderStatus:      types.Success,
			expectedFilledSize:       0,
			expectedTransactionIndex: 0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{
				// Write the stateful order to state and memStore.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId.Marshal(),
					)),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				"StatefulOrdersTimeSlice/value/1970-01-01T00:00:10.000000000",
			},
		},
		`Can place multiple stateful orders on the orderbook, and the newly-placed stateful order
			matches and is written to state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
			},

			order: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,

			expectedOrderStatus:      types.Success,
			expectedFilledSize:       100_000_000,
			expectedTransactionIndex: 2,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{
				// Update taker subaccount.
				fmt.Sprintf(
					"Subaccount/value/%v/",
					string(proto.MustFirst(constants.Carl_Num0.Marshal())),
				),
				indexer_manager.IndexerEventsKey,
				// Update maker subaccount.
				fmt.Sprintf(
					"Subaccount/value/%v/",
					string(proto.MustFirst(constants.Dave_Num0.Marshal())),
				),
				indexer_manager.IndexerEventsKey,
				// Update taker order fill amount to state and memStore.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId.Marshal())),
				),
				// Update maker order fill amount to state and memStore.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId.Marshal())),
				),
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId.Marshal())),
				),
				// Write the taker stateful order placement to memStore and state.
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId.Marshal(),
					)),
				),
				fmt.Sprintf(
					"StatefulOrderPlacement/value/%v",
					string(proto.MustFirst(
						constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId.Marshal(),
					)),
				),
				"NextStatefulOrderBlockTransactionIndex/value",
				"StatefulOrdersTimeSlice/value/1970-01-01T00:00:10.000000000",
			},
		},
		`Can place a stateful post-only order that crosses the book and the order placement isn't
			written to state`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Alice_Num1_10_000USD,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20,
			},

			order: constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO,

			expectedOrderStatus:      types.Success,
			expectedErr:              types.ErrPostOnlyWouldCrossMakerOrder,
			expectedFilledSize:       0,
			expectedTransactionIndex: 0,
			expectedSeenPlaceOrder:   true,
			expectedMultiStoreWrites: []string{
				// Update taker subaccount.
				fmt.Sprintf(
					"Subaccount/value/%v/",
					string(proto.MustFirst(constants.Alice_Num0.Marshal())),
				),
				indexer_manager.IndexerEventsKey,
				// Update maker subaccount.
				fmt.Sprintf(
					"Subaccount/value/%v/",
					string(proto.MustFirst(constants.Alice_Num1.Marshal())),
				),
				indexer_manager.IndexerEventsKey,
				// Update taker order fill amount.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO.OrderId.Marshal())),
				),
				// Update taker order fill amount in memStore.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25_PO.OrderId.Marshal())),
				),
				// Update prunable block height for maker fill amount.
				fmt.Sprintf(
					"BlockHeightToPotentiallyPrunableOrders/value/%v",
					string(types.BlockHeightToPotentiallyPrunableOrdersKey(
						constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.GetGoodTilBlock()+types.ShortBlockWindow),
					),
				),
				// Update maker order fill amount
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId.Marshal())),
				),
				// Update maker order fill amount in memStore.
				fmt.Sprintf(
					"OrderAmount/value/%v",
					string(proto.MustFirst(constants.Order_Alice_Num1_Clob0_Id4_Buy10_Price45_GTB20.OrderId.Marshal())),
				),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)

			ctx,
				clobKeeper,
				pricesKeeper,
				assetsKeeper,
				perpetualsKeeper,
				subaccountsKeeper,
				_,
				_ := keepertest.ClobKeepers(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ctx = ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, pricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up USDC asset in assets module.
			err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := perpetualsKeeper.CreatePerpetual(
					ctx,
					p.Ticker,
					p.MarketId,
					p.AtomicResolution,
					p.DefaultFundingPpm,
					p.LiquidityTier,
				)
				require.NoError(t, err)
			}

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				subaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobs {
				_, err := clobKeeper.CreatePerpetualClobPair(
					ctx,
					clobPair.GetPerpetualClobMetadata().PerpetualId,
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
					clobPair.MakerFeePpm,
					clobPair.TakerFeePpm,
				)
				require.NoError(t, err)
			}

			// Set the block height and last committed block time.
			blockHeight := uint32(2)
			ctx = ctx.WithBlockHeight(int64(blockHeight)).WithBlockTime(time.Unix(5, 0))
			clobKeeper.SetBlockTimeForLastCommittedBlock(ctx)

			// Create all existing orders.
			for _, order := range tc.existingOrders {
				_, orderStatus, err := clobKeeper.CheckTxPlaceOrder(ctx, &types.MsgPlaceOrder{Order: order})
				require.NoError(t, err)
				require.True(t, orderStatus.IsSuccess())
			}

			// Verify the order that will be placed is a Long-Term order.
			require.True(t, tc.order.OrderId.IsLongTermOrder())

			// Run the test.
			traceDecoder := &tracer.TraceDecoder{}
			ctx.MultiStore().SetTracer(traceDecoder)

			msg := &types.MsgPlaceOrder{Order: tc.order}
			orderSizeOptimisticallyFilledFromMatching,
				orderStatus,
				err := clobKeeper.CheckTxPlaceOrder(ctx, msg)

			// Verify test expectations.
			require.ErrorIs(t, err, tc.expectedErr)
			statefulOrderPlacement, found := clobKeeper.GetStatefulOrderPlacement(ctx, tc.order.OrderId)
			statefulOrderIds := clobKeeper.GetStatefulOrdersTimeSlice(ctx, tc.order.MustGetUnixGoodTilBlockTime())
			if err == nil {
				require.Equal(t, tc.expectedOrderStatus, orderStatus)
				require.Equal(t, tc.expectedFilledSize, orderSizeOptimisticallyFilledFromMatching)
				require.Equal(t, tc.order, statefulOrderPlacement.Order)
				// The block height is incremented by 1 because the order placement is written to state in `CheckTx`.
				require.Equal(t, blockHeight+1, statefulOrderPlacement.BlockHeight)
				require.Equal(t, tc.expectedTransactionIndex, statefulOrderPlacement.TransactionIndex)
				require.Contains(
					t,
					statefulOrderIds,
					tc.order.OrderId,
				)
			} else {
				require.False(t, found)
				require.NotContains(
					t,
					statefulOrderIds,
					tc.order.OrderId,
				)
			}

			traceDecoder.RequireKeyPrefixWrittenInSequence(t, tc.expectedMultiStoreWrites)

			require.Equal(
				t,
				tc.expectedSeenPlaceOrder,
				clobKeeper.HasSeenPlaceOrder(ctx, *msg),
			)
		})
	}
}
func TestPlaceOrder_SendOffchainMessages(t *testing.T) {
	indexerEventManager := &mocks.IndexerEventManager{}
	for _, message := range constants.TestOffchainMessages {
		indexerEventManager.On(
			"SendOffchainData",
			message.AddHeader(constants.TestTxHashHeader),
		).Return().Once()
	}

	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()

	ctx, keeper, pricesKeeper, _, perpetualsKeeper, _, _, _ :=
		keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, indexerEventManager)
	prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	ctx = ctx.WithIsCheckTx(true)

	memClob.On("CreateOrderbook", ctx, constants.ClobPair_Btc).Return()
	_, err := keeper.CreatePerpetualClobPair(
		ctx,
		clobtest.MustPerpetualId(constants.ClobPair_Btc),
		satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
		satypes.BaseQuantums(constants.ClobPair_Btc.MinOrderBaseQuantums),
		constants.ClobPair_Btc.QuantumConversionExponent,
		constants.ClobPair_Btc.SubticksPerTick,
		constants.ClobPair_Btc.Status,
		constants.ClobPair_Btc.MakerFeePpm,
		constants.ClobPair_Btc.TakerFeePpm,
	)
	require.NoError(t, err)

	order := constants.Order_Carl_Num0_Id5_Clob0_Buy2BTC_Price50000
	msgPlaceOrder := &types.MsgPlaceOrder{Order: order}
	memClob.On("PlaceOrder", ctx, order, true).
		Return(order.GetBaseQuantums(), types.OrderStatus(0), constants.TestOffchainUpdates, nil)

	_, _, err = keeper.CheckTxPlaceOrder(ctx, msgPlaceOrder)
	require.NoError(t, err)
	indexerEventManager.AssertNumberOfCalls(t, "SendOffchainData", len(constants.TestOffchainMessages))
	indexerEventManager.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestPerformStatefulOrderValidation_PreExistingStatefulOrder(t *testing.T) {
	// Setup keeper state.
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ctx,
		clobKeeper,
		pricesKeeper,
		_,
		perpetualsKeeper,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	memClob.On("CreateOrderbook", ctx, constants.ClobPair_Btc).Return()
	_, err := clobKeeper.CreatePerpetualClobPair(
		ctx,
		clobtest.MustPerpetualId(constants.ClobPair_Btc),
		satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
		satypes.BaseQuantums(constants.ClobPair_Btc.MinOrderBaseQuantums),
		constants.ClobPair_Btc.QuantumConversionExponent,
		constants.ClobPair_Btc.SubticksPerTick,
		constants.ClobPair_Btc.Status,
		constants.ClobPair_Btc.MakerFeePpm,
		constants.ClobPair_Btc.TakerFeePpm,
	)
	require.NoError(t, err)
	ctx = ctx.WithBlockHeight(int64(100)).WithBlockTime(time.Unix(5, 0))
	clobKeeper.SetBlockTimeForLastCommittedBlock(ctx)
	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15

	// TODO(CLOB-249): Re-implement test once sanity check is added.
	// // Run the test if the preexisting order is not in state. Expected panic.
	// require.PanicsWithError(
	// 	t,
	// 	fmt.Errorf(
	// 		"PerformStatefulOrderValidation: Called for preExistingStatefulOrder "+
	// 			"%+v but does not exist in state",
	// 		order,
	// 	).Error(),
	// 	func() {
	// 		_ = clobKeeper.PerformStatefulOrderValidation(ctx, &order, 10, true)
	// 	},
	// )

	// Run the test if the preexisting order is in state. Expected no panic.
	clobKeeper.SetStatefulOrderPlacement(ctx, order, 10)
	err = clobKeeper.PerformStatefulOrderValidation(ctx, &order, 10, true)
	require.NoError(t, err)
}

func TestPerformStatefulOrderValidation(t *testing.T) {
	blockHeight := uint32(5)

	tests := map[string]struct {
		setupState  func(ctx sdk.Context, k *keeper.Keeper)
		order       types.Order
		expectedErr string
	}{
		"Succeeds with a GoodTilBlock of blockHeight": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight},
			},
		},
		"Succeeds with a GoodTilBlock of blockHeight + ShortBlockWindow": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + types.ShortBlockWindow},
			},
		},
		"Fails with invalid ClobPairId": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(1),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + 5},
			},
			expectedErr: types.ErrInvalidClob.Error(),
		},
		"Fails if Subticks is not a multiple of SubticksPerTick": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     77,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + 5},
			},
			expectedErr: "must be a multiple of the ClobPair's SubticksPerTick",
		},
		"Fails if Quantums < MinOrderBaseQuantums": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     24,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + 5},
			},
			expectedErr: "must be greater than the ClobPair's MinOrderBaseQuantums",
		},
		"Fails if Quantums is not a multiple of StepBaseQuantums": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     599,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + 5},
			},
			expectedErr: "must be a multiple of the ClobPair's StepBaseQuantums",
		},
		"Fails if GoodTilBlock is in the past": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight - 1},
			},
			expectedErr: types.ErrHeightExceedsGoodTilBlock.Error(),
		},
		"Fails if GoodTilBlock Exceeds ShortBlockWindow": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					ClobPairId:   uint32(0),
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: blockHeight + types.ShortBlockWindow + 1},
			},
			expectedErr: types.ErrGoodTilBlockExceedsShortBlockWindow.Error(),
		},
		"Long-term: Fails if GoodTilBlockTime is less than or equal to the block time of the previous block": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 4,
				},
			},
			expectedErr: types.ErrTimeExceedsGoodTilBlockTime.Error(),
		},
		"Long-term: Fails if GoodTilBlockTime Exceeds StatefulOrderTimeWindow": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: lib.MustConvertIntegerToUint32(
						time.Unix(5, 0).Add(types.StatefulOrderTimeWindow).Unix() + 1,
					),
				},
			},
			expectedErr: types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow.Error(),
		},
		`Long-term: Returns error when order already exists in state`: {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetStatefulOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     600,
						Subticks:     78,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
		`Long-term: Returns error when order with same order id but lower priority exists in state`: {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetStatefulOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_LongTerm,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     600,
						Subticks:     78,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_LongTerm,
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
		"Conditional: Fails if GoodTilBlockTime is in the past": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_Conditional,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: 4,
				},
			},
			expectedErr: types.ErrTimeExceedsGoodTilBlockTime.Error(),
		},
		"Conditional: Fails if GoodTilBlockTime Exceeds StatefulOrderTimeWindow": {
			order: types.Order{
				OrderId: types.OrderId{
					ClientId:     0,
					SubaccountId: constants.Alice_Num0,
					OrderFlags:   types.OrderIdFlags_Conditional,
					ClobPairId:   uint32(0),
				},
				Side:     types.Order_SIDE_BUY,
				Quantums: 600,
				Subticks: 78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{
					GoodTilBlockTime: lib.MustConvertIntegerToUint32(
						time.Unix(5, 0).Add(types.StatefulOrderTimeWindow).Unix() + 1,
					),
				},
			},
			expectedErr: types.ErrGoodTilBlockTimeExceedsStatefulOrderTimeWindow.Error(),
		},
		`Conditional: Returns error when order already exists in state`: {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetStatefulOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_Conditional,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     600,
						Subticks:     78,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_Conditional,
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
		`Conditional: Returns error when order with same order id but lower priority exists in state`: {
			setupState: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetStatefulOrderPlacement(
					ctx,
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Alice_Num0,
							ClientId:     0,
							OrderFlags:   types.OrderIdFlags_Conditional,
						},
						Side:         types.Order_SIDE_BUY,
						Quantums:     600,
						Subticks:     78,
						GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
					},
					lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				)
			},
			order: types.Order{
				OrderId: types.OrderId{
					SubaccountId: constants.Alice_Num0,
					ClientId:     0,
					OrderFlags:   types.OrderIdFlags_Conditional,
				},
				Side:         types.Order_SIDE_BUY,
				Quantums:     600,
				Subticks:     78,
				GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 15},
			},
			expectedErr: types.ErrStatefulOrderAlreadyExists.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("CreateOrderbook", mock.Anything, mock.Anything, mock.Anything)
			memClob.On("SetClobKeeper", mock.Anything).Return()

			ctx, keeper, pricesKeeper, _, perpetualsKeeper, _, _, _ := keepertest.ClobKeepers(
				t,
				memClob,
				&mocks.BankKeeper{},
				&mocks.IndexerEventManager{},
			)
			ctx = ctx.WithBlockTime(time.Unix(5, 0))
			prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			clobPair := types.ClobPair{
				Metadata: &types.ClobPair_PerpetualClobMetadata{
					PerpetualClobMetadata: &types.PerpetualClobMetadata{
						PerpetualId: 0,
					},
				},
				Status:               types.ClobPair_STATUS_ACTIVE,
				StepBaseQuantums:     12,
				SubticksPerTick:      39,
				MinOrderBaseQuantums: 204,
			}

			_, err := keeper.CreatePerpetualClobPair(
				ctx,
				clobtest.MustPerpetualId(clobPair),
				satypes.BaseQuantums(clobPair.StepBaseQuantums),
				satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
				clobPair.QuantumConversionExponent,
				clobPair.SubticksPerTick,
				clobPair.Status,
				clobPair.MakerFeePpm,
				clobPair.TakerFeePpm,
			)
			require.NoError(t, err)
			keeper.SetBlockTimeForLastCommittedBlock(ctx)

			if tc.setupState != nil {
				tc.setupState(ctx, keeper)
			}

			err = keeper.PerformStatefulOrderValidation(ctx, &tc.order, blockHeight, false)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetStatePosition_Success(t *testing.T) {
	tests := map[string]struct {
		// Subaccount state.
		subaccount *satypes.Subaccount
		// CLOB state.
		clob types.ClobPair

		// Parameters.
		subaccountId satypes.SubaccountId
		clobPairId   types.ClobPairId

		// Expectations.
		expectedPositionSize *big.Int
	}{
		`Can fetch the position size of a long position`: {
			subaccount: &constants.Dave_Num0_1BTC_Long_50000USD,

			subaccountId: constants.Dave_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Btc.Id),

			expectedPositionSize: constants.Dave_Num0_1BTC_Long_50000USD.PerpetualPositions[0].GetBigQuantums(),
		},
		`Can fetch the position size from multiple positions`: {
			subaccount: &constants.Dave_Num0_1BTC_Long_1ETH_Long_46000USD_Short,

			subaccountId: constants.Dave_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Eth.Id),

			expectedPositionSize: constants.Dave_Num0_1BTC_Long_1ETH_Long_46000USD_Short.PerpetualPositions[1].GetBigQuantums(),
		},
		`Can fetch the position size of a short position`: {
			subaccount: &constants.Carl_Num0_1BTC_Short,

			subaccountId: constants.Carl_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Btc.Id),

			expectedPositionSize: constants.Carl_Num0_1BTC_Short.PerpetualPositions[0].GetBigQuantums(),
		},
		`Fetching a non-existent subaccount returns 0`: {
			subaccountId: constants.Carl_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Btc.Id),

			expectedPositionSize: big.NewInt(0),
		},
		`Fetching a subaccount that doesn't have an open position corresponding to CLOB returns 0`: {
			subaccount: &constants.Dave_Num0_1BTC_Long_50000USD,

			subaccountId: constants.Dave_Num0,
			clobPairId:   types.ClobPairId(constants.ClobPair_Eth.Id),

			expectedPositionSize: big.NewInt(0),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ctx,
				clobKeeper,
				pricesKeeper,
				_,
				perpetualsKeeper,
				subaccountsKeeper,
				_,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Create subaccount if it's specified.
			if tc.subaccount != nil {
				subaccountsKeeper.SetSubaccount(ctx, *tc.subaccount)
			}
			prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			// Create CLOB pairs.
			clobPairs := []types.ClobPair{constants.ClobPair_Btc, constants.ClobPair_Eth}
			for _, cp := range clobPairs {
				_, err := clobKeeper.CreatePerpetualClobPair(
					ctx,
					clobtest.MustPerpetualId(cp),
					satypes.BaseQuantums(cp.StepBaseQuantums),
					satypes.BaseQuantums(cp.MinOrderBaseQuantums),
					cp.QuantumConversionExponent,
					cp.SubticksPerTick,
					cp.Status,
					cp.MakerFeePpm,
					cp.TakerFeePpm,
				)
				require.NoError(t, err)
			}

			// Run the test and verify expectations.
			positionSizeBig := clobKeeper.GetStatePosition(ctx, tc.subaccountId, tc.clobPairId)

			require.Equal(t, tc.expectedPositionSize, positionSizeBig)
		})
	}
}

func TestGetStatePosition_PanicsOnInvalidClob(t *testing.T) {
	// Setup keeper state.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ctx,
		clobKeeper,
		pricesKeeper,
		_,
		perpetualsKeeper,
		_,
		_,
		_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
	prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	// Run the test and verify expectations.
	clobPairId := types.ClobPairId(constants.ClobPair_Eth.Id)
	require.PanicsWithValue(
		t,
		fmt.Sprintf("GetStatePosition: CLOB pair %d not found", clobPairId),
		func() {
			clobKeeper.GetStatePosition(ctx, constants.Alice_Num0, clobPairId)
		},
	)
}

// TODO(DEC-1535): Uncomment this test when we can create spot CLOB pairs.
// func TestGetStatePosition_PanicsOnAssetClob(t *testing.T) {
// 	// Setup keeper state.
// 	MemClob := memclob.NewMemClobPriceTimePriority(false)
// 	ctx,
// 		clobKeeper,
// 		pricesKeeper,
// 		_,
// 		perpetualsKeeper,
// 		_,
// 		_,
// 		_ := keepertest.ClobKeepers(t, MemClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
// 	prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
// 	perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

// 	// Create CLOB pair.
// 	clobPair := constants.ClobPair_Asset
// 	clobKeeper.CreatePerpetualClobPair(
// 		ctx,
// 		clobPair.Metadata.(*types.ClobPair_PerpetualClobMetadata),
// 		satypes.BaseQuantums(clobPair.StepBaseQuantums),
// 		satypes.BaseQuantums(clobPair.MinOrderBaseQuantums),
// 		clobPair.QuantumConversionExponent,
// 		clobPair.SubticksPerTick,
// 		clobPair.Status,
// 		clobPair.MakerFeePpm,
// 		clobPair.TakerFeePpm,
// 	)

// 	// Run the test and verify expectations.
// 	require.PanicsWithError(
// 		t,
// 		sdkerrors.Wrap(
// 			types.ErrAssetOrdersNotImplemented,
// 			"GetStatePosition: Reduce-only orders for assets not implemented",
// 		).Error(),
// 		func() {
// 			clobKeeper.GetStatePosition(ctx, constants.Alice_Num0, clobPair.Id)
// 		},
// 	)
// }

func TestInitStatefulOrdersInMemClob(t *testing.T) {
	tests := map[string]struct {
		// CLOB module return values.
		statefulOrders       []types.Order
		orderPlacementErrors []error
	}{
		`Can initialize no stateful order in the memclob with no errors`: {
			statefulOrders:       []types.Order{},
			orderPlacementErrors: []error{},
		},
		`Can initialize one stateful order in the memclob with no errors`: {
			statefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
			orderPlacementErrors: []error{
				nil,
			},
		},
		`Can initialize multiple stateful orders in the memclob with no errors`: {
			statefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
			orderPlacementErrors: []error{
				nil,
				nil,
				nil,
			},
		},
		`Can initialize multiple stateful orders in the memclob with errors`: {
			statefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
			orderPlacementErrors: []error{
				nil,
				types.ErrInvalidStatefulOrderGoodTilBlockTime,
				types.ErrTimeExceedsGoodTilBlockTime,
			},
		},
		`Can initialize multiple stateful orders in the memclob where each order throws an error`: {
			statefulOrders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Alice_Num0_Id2_Clob0_Sell65_Price10_GTBT25,
			},
			orderPlacementErrors: []error{
				types.ErrInvalidStatefulOrderGoodTilBlockTime,
				types.ErrInvalidStatefulOrderGoodTilBlockTime,
				types.ErrTimeExceedsGoodTilBlockTime,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup state.
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()

			require.Len(t, tc.statefulOrders, len(tc.orderPlacementErrors))

			indexerEventManager := &mocks.IndexerEventManager{}

			ctx, keeper, pricesKeeper, _, perpetualsKeeper, _, _, _ :=
				keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, indexerEventManager)
			prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			// Create CLOB pair.
			memClob.On("CreateOrderbook", mock.Anything, constants.ClobPair_Btc).Return()
			_, err := keeper.CreatePerpetualClobPair(
				ctx,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				satypes.BaseQuantums(constants.ClobPair_Btc.MinOrderBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
				constants.ClobPair_Btc.MakerFeePpm,
				constants.ClobPair_Btc.TakerFeePpm,
			)
			require.NoError(t, err)

			// Create each stateful order placement in state and properly mock the MemClob call.
			for i, order := range tc.statefulOrders {
				require.True(t, order.IsStatefulOrder())

				keeper.SetStatefulOrderPlacement(ctx, order, uint32(i))
				orderPlacementErr := tc.orderPlacementErrors[i]
				memClob.On("PlaceOrder", mock.Anything, order, false).Return(
					satypes.BaseQuantums(0),
					types.Success,
					constants.TestOffchainUpdates,
					orderPlacementErr,
				).Once()

				for _, message := range constants.TestOffchainMessages {
					indexerEventManager.On("SendOffchainData", message).Return().Once()
				}
			}

			// Run the test and verify expectations.
			keeper.InitStatefulOrdersInMemClob(ctx)
			indexerEventManager.AssertExpectations(t)
			indexerEventManager.AssertNumberOfCalls(
				t,
				"SendOffchainData",
				len(constants.TestOffchainMessages)*len(tc.statefulOrders),
			)
			memClob.AssertExpectations(t)
		})
	}
}

func TestPlaceStatefulOrdersFromLastBlock(t *testing.T) {
	tests := map[string]struct {
		orders []types.Order

		expectedOrderPlacementCalls []types.Order
	}{
		"empty stateful orders": {
			orders: []types.Order{},

			expectedOrderPlacementCalls: []types.Order{},
		},
		"places stateful orders from last block": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
			},
			expectedOrderPlacementCalls: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
			},
		},
		"does not place orders with GTBT equal to block time": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"does not place orders with GTBT less than block time": {
			orders: []types.Order{
				{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						OrderFlags:   types.OrderIdFlags_LongTerm,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     5,
					Subticks:     10,
					GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
				},
			},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"does not place orders with invalid clob pair id": {
			orders: []types.Order{
				{
					OrderId:      constants.InvalidClobPairId_Long_Term_Order,
					Side:         types.Order_SIDE_BUY,
					Quantums:     5,
					Subticks:     10,
					GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 2},
				},
			},
			expectedOrderPlacementCalls: []types.Order{},
		},
		"does not place orders further than StatefulOrderTimeWindow in the future": {
			orders: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				{
					OrderId: types.OrderId{
						SubaccountId: constants.Alice_Num0,
						ClientId:     0,
						OrderFlags:   types.OrderIdFlags_LongTerm,
						ClobPairId:   0,
					},
					Side:     types.Order_SIDE_BUY,
					Quantums: 5,
					Subticks: 10,
					GoodTilOneof: &types.Order_GoodTilBlockTime{
						GoodTilBlockTime: lib.MustConvertIntegerToUint32(
							5 + int64(types.StatefulOrderTimeWindow.Seconds()+1),
						),
					},
				},
			},
			expectedOrderPlacementCalls: []types.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup state.
			memClob := &mocks.MemClob{}

			memClob.On("SetClobKeeper", mock.Anything).Return()

			ctx, keeper, pricesKeeper, _, perpetualsKeeper, _, _, _ :=
				keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, indexer_manager.NewIndexerEventManagerNoop())
			prices.InitGenesis(ctx, *pricesKeeper, constants.Prices_DefaultGenesisState)
			perpetuals.InitGenesis(ctx, *perpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

			ctx = ctx.WithBlockHeight(int64(100)).WithBlockTime(time.Unix(5, 0))
			ctx = ctx.WithIsCheckTx(true)
			keeper.SetBlockTimeForLastCommittedBlock(ctx)

			// Create CLOB pair.
			memClob.On("CreateOrderbook", mock.Anything, constants.ClobPair_Btc).Return()
			_, err := keeper.CreatePerpetualClobPair(
				ctx,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				satypes.BaseQuantums(constants.ClobPair_Btc.MinOrderBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
				constants.ClobPair_Btc.MakerFeePpm,
				constants.ClobPair_Btc.TakerFeePpm,
			)
			require.NoError(t, err)

			// Create each stateful order placement in state
			for i, order := range tc.orders {
				require.True(t, order.IsStatefulOrder())

				keeper.SetStatefulOrderPlacement(ctx, order, uint32(i))
			}

			// Assert expected order placement memclob calls.
			for _, order := range tc.expectedOrderPlacementCalls {
				memClob.On("PlaceOrder", mock.Anything, order, false).Return(
					satypes.BaseQuantums(0),
					types.Success,
					constants.TestOffchainUpdates,
					nil,
				).Once()
			}

			// Run the test and verify expectations.
			offchainUpdates := types.NewOffchainUpdates()
			keeper.PlaceStatefulOrdersFromLastBlock(ctx, tc.orders, offchainUpdates)

			// Verify that all removed orders have an associated off-chain update.
			orderMap := make(map[types.OrderId]bool)
			for _, order := range tc.orders {
				orderMap[order.OrderId] = true
			}

			removedOrders := lib.FilterSlice(tc.expectedOrderPlacementCalls, func(order types.Order) bool {
				return !orderMap[order.OrderId]
			})

			for _, order := range removedOrders {
				_, exists := offchainUpdates.RemoveMessages[order.OrderId]
				require.True(t, exists)
			}

			memClob.AssertExpectations(t)
		})
	}
}
