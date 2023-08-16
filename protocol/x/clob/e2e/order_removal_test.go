package clob_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestConditionalOrderRemoval(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order

		// Optional withdraw message for under-collateralized tests.
		withdrawal  *sendingtypes.MsgWithdrawFromSubaccount
		priceUpdate *prices.MsgUpdateMarketPrices

		// Optional short term order
		subsequentOrder *clobtypes.Order

		expectedOrderRemovals []bool
	}{
		"conditional post-only order crosses maker": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
				constants.ConditionalOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO_SL_15,
			},

			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 1_490_000),
				},
			},
			expectedOrderRemovals: []bool{
				false,
				true, // P0 order should be removed
			},
		},
		"conditional fill-or-kill order does not fully match and is removed": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_FOK,
			},

			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_400_000),
				},
			},
			expectedOrderRemovals: []bool{
				false,
				true, // non fully filled FOK order should be removed
			},
		},
		"conditional IOC order does not fully match and is removed": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell025BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
			},

			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_400_000),
				},
			},
			expectedOrderRemovals: []bool{
				false,
				true, // non fully filled IOC order should be removed
			},
		},
		"conditional self trade removes maker order": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
				constants.ConditionalOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10_SL_15,
			},

			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 1_490_000),
				},
			},
			expectedOrderRemovals: []bool{
				true, // Self trade removes the maker order.
				false,
			},
		},
		"fully filled maker orders triggered by conditional order are removed": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
				constants.ConditionalOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15_SL_15,
			},
			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 1_490_000),
				},
			},
			expectedOrderRemovals: []bool{
				true, // maker order fully filled
				false,
			},
		},
		"fully filled conditional taker orders are removed": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5_SL_15,
			},
			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 1_510_000),
				},
			},

			expectedOrderRemovals: []bool{
				false,
				true, // taker order fully filled
			},
		},
		"under-collateralized conditional taker during matching is removed": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_50003,
			},
			withdrawal: &sendingtypes.MsgWithdrawFromSubaccount{
				Sender:    constants.Dave_Num0,
				Recipient: constants.DaveAccAddress.String(),
				AssetId:   constants.Usdc.Id,
				Quantums:  10_000_000_000,
			},
			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_250_000),
				},
			},

			expectedOrderRemovals: []bool{
				false,
				true, // taker order fails collateralization check during matching
			},
		},
		"under-collateralized conditional taker when adding to book is removed": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Dave_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTBT10,
				// Does not cross with best bid.
				constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_50003,
			},
			withdrawal: &sendingtypes.MsgWithdrawFromSubaccount{
				Sender:    constants.Dave_Num0,
				Recipient: constants.DaveAccAddress.String(),
				AssetId:   constants.Usdc.Id,
				Quantums:  10_000_000_000,
			},
			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_250_000),
				},
			},

			expectedOrderRemovals: []bool{
				false,
				true, // taker order fails add-to-orderbook collateralization check
			},
		},
		"under-collateralized conditional maker is removed": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_500000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_SL_50003,
			},
			withdrawal: &sendingtypes.MsgWithdrawFromSubaccount{
				Sender:    constants.Dave_Num0,
				Recipient: constants.DaveAccAddress.String(),
				AssetId:   constants.Usdc.Id,
				Quantums:  500_000_000_000,
			},
			priceUpdate: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_250_000),
				},
			},

			subsequentOrder: &constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,

			expectedOrderRemovals: []bool{
				true, // maker is under-collateralized
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
							constants.ClobPair_Btc_No_Fee,
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
					require.True(t, tApp.CheckTx(checkTx).IsOK())
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

			// Make sure conditional orders are triggered.
			for _, order := range tc.orders {
				if order.IsConditionalOrder() {
					require.Equal(t, true, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, order.OrderId))
				}
			}

			// Do the optional withdraw.
			if tc.withdrawal != nil {
				CheckTx_MsgWithdrawFromSubaccount := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: testtx.MustGetSignerAddress(tc.withdrawal),
						Gas:                  100_000,
					},
					tc.withdrawal,
				)
				checkTxResp := tApp.CheckTx(CheckTx_MsgWithdrawFromSubaccount)
				require.True(t, checkTxResp.IsOK())
			}
			// Advance to the next block, persisting removals in operations queue to state.
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

			if tc.subsequentOrder != nil {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(*tc.subsequentOrder),
				) {
					require.True(t, tApp.CheckTx(checkTx).IsOK())
				}
			}

			// Advance to the next block, persisting removals in operations queue to state.
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})

			require.Equal(t, len(tc.orders), len(tc.expectedOrderRemovals))

			// Verify expectations.
			for idx, order := range tc.orders {
				_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
				require.Equal(t, tc.expectedOrderRemovals[idx], !found)
			}
		})
	}
}

func TestOrderRemoval(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		firstOrder  clobtypes.Order
		secondOrder clobtypes.Order

		// Optional withdraw message for under-collateralized tests.
		withdrawal *sendingtypes.MsgWithdrawFromSubaccount

		expectedFirstOrderRemoved  bool
		expectedSecondOrderRemoved bool
	}{
		"post-only order crosses maker": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			firstOrder:  constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			secondOrder: constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell10_Price10_GTBT10_PO,

			expectedFirstOrderRemoved:  false,
			expectedSecondOrderRemoved: true, // PO order should be removed.
		},
		"self trade removes maker order": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
			},
			firstOrder:  constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			secondOrder: constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10,

			expectedFirstOrderRemoved:  true, // Self trade removes the maker order.
			expectedSecondOrderRemoved: false,
		},
		"fully filled maker orders are removed": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			firstOrder:  constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			secondOrder: constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,

			expectedFirstOrderRemoved:  true, // maker order fully filled
			expectedSecondOrderRemoved: false,
		},
		"fully filled taker orders are removed": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			firstOrder:  constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
			secondOrder: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,

			expectedFirstOrderRemoved:  false,
			expectedSecondOrderRemoved: true, // taker order fully filled
		},
		"under-collateralized taker during matching is removed": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
			},
			firstOrder:  constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			secondOrder: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,

			withdrawal: &sendingtypes.MsgWithdrawFromSubaccount{
				Sender:    constants.Dave_Num0,
				Recipient: constants.DaveAccAddress.String(),
				AssetId:   constants.Usdc.Id,
				Quantums:  10_000_000_000,
			},

			expectedFirstOrderRemoved:  false,
			expectedSecondOrderRemoved: true, // taker order fails collateralization check during matching
		},
		"under-collateralized taker when adding to book is removed": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
			},
			firstOrder: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTBT10,
			// Does not cross with best bid.
			secondOrder: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,

			withdrawal: &sendingtypes.MsgWithdrawFromSubaccount{
				Sender:    constants.Dave_Num0,
				Recipient: constants.DaveAccAddress.String(),
				AssetId:   constants.Usdc.Id,
				Quantums:  10_000_000_000,
			},

			expectedFirstOrderRemoved:  false,
			expectedSecondOrderRemoved: true, // taker order fails add-to-orderbook collateralization check
		},
		"under-collateralized maker is removed": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
			},
			firstOrder:  constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			secondOrder: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,

			withdrawal: &sendingtypes.MsgWithdrawFromSubaccount{
				Sender:    constants.Carl_Num0,
				Recipient: constants.CarlAccAddress.String(),
				AssetId:   constants.Usdc.Id,
				Quantums:  10_000_000_000,
			},

			expectedFirstOrderRemoved:  true, // maker is under-collateralized
			expectedSecondOrderRemoved: false,
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
							constants.ClobPair_Btc_No_Fee,
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
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				tApp.App,
				*clobtypes.NewMsgPlaceOrder(tc.firstOrder),
			) {
				require.True(t, tApp.CheckTx(checkTx).IsOK())
			}
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				tApp.App,
				*clobtypes.NewMsgPlaceOrder(tc.secondOrder),
			) {
				require.True(t, tApp.CheckTx(checkTx).IsOK())
			}

			// Do the optional withdraw.
			if tc.withdrawal != nil {
				CheckTx_MsgWithdrawFromSubaccount := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: testtx.MustGetSignerAddress(tc.withdrawal),
						Gas:                  100_000,
					},
					tc.withdrawal,
				)
				checkTxResp := tApp.CheckTx(CheckTx_MsgWithdrawFromSubaccount)
				require.True(t, checkTxResp.IsOK())
			}

			// First block only persists stateful orders to state without matching them.
			// Therefore, both orders should be in state at this point.
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.firstOrder.OrderId)
			require.True(t, found)
			_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.secondOrder.OrderId)
			require.True(t, found)

			// Verify expectations.
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.firstOrder.OrderId)
			require.Equal(t, tc.expectedFirstOrderRemoved, !found)

			_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.secondOrder.OrderId)
			require.Equal(t, tc.expectedSecondOrderRemoved, !found)
		})
	}
}

func TestOrderRemoval_MultipleReplayOperationsDuringPrepareCheckState(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					constants.Alice_Num0_10_000USD,
					constants.Bob_Num0_10_000USD,
				}
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
					constants.ClobPair_Btc_No_Fee,
				}
				genesisState.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
			},
		)
		return genesis
	}).WithTesting(t).Build()
	ctx := tApp.InitChain()

	// Create a resting order for alice.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15_PO,
		),
	) {
		require.True(t, tApp.CheckTx(checkTx).IsOK())
	}
	// Partially match alice's order so that it's in the operations queue.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell5_Price10_GTBT10,
		),
	) {
		require.True(t, tApp.CheckTx(checkTx).IsOK())
	}
	// Now remove alice's order somehow. Self-trade in this case.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10,
		),
	) {
		require.True(t, tApp.CheckTx(checkTx).IsOK())
	}
	// Place another order to invalidate Alice's post only order.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell5_Price5_GTBT10,
		),
	) {
		require.True(t, tApp.CheckTx(checkTx).IsOK())
	}

	_ = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Local operations queue would be [placement(Alice_Order), ..., removal(Alice_Order)].
	// Let's say block proposer does not include these operations. Make sure we don't panic in this case.
	_ = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{},
	})
	_ = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
}
