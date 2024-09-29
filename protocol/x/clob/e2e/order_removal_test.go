package clob_test

import (
	"math/big"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	clobtestutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/clob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	testtx "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/tx"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	feetiertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	prices "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	sendingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestConditionalOrderRemoval(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order

		// Optional withdraw message for under-collateralized tests.
		withdrawal  *sendingtypes.MsgWithdrawFromSubaccount
		priceUpdate map[uint32]ve.VEPricePair

		// Optional short term order
		subsequentOrder *clobtypes.Order

		expectedOrderRemovals []bool

		disableNonDeterminismChecks bool
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

			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 1_490_000,
					PnlPrice:  1_490_000,
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

			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
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

			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			expectedOrderRemovals: []bool{
				true,
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

			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 1_490_000,
					PnlPrice:  1_490_000,
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
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 1_490_000,
					PnlPrice:  1_490_000,
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
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 1_510_000,
					PnlPrice:  1_510_000,
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
				AssetId:   constants.TDai.Id,
				Quantums:  10_000_000_000,
			},
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_250_000,
					PnlPrice:  5_000_250_000,
				},
			},

			expectedOrderRemovals: []bool{
				false,
				true, // taker order fails collateralization check during matching
			},
			// TODO(CORE-858): Re-enable determinism checks once non-determinism issue is found and resolved.
			disableNonDeterminismChecks: true,
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
				AssetId:   constants.TDai.Id,
				Quantums:  10_000_000_000,
			},
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_250_000,
					PnlPrice:  5_000_250_000,
				},
			},

			expectedOrderRemovals: []bool{
				false,
				true, // taker order fails add-to-orderbook collateralization check
			},
			// TODO(CORE-858): Re-enable determinism checks once non-determinism issue is found and resolved.
			disableNonDeterminismChecks: true,
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
				AssetId:   constants.TDai.Id,
				Quantums:  500_000_000_000,
			},
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_250_000,
					PnlPrice:  5_000_250_000,
				},
			},

			subsequentOrder: &constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,

			expectedOrderRemovals: []bool{
				true, // maker is under-collateralized
			},
			// TODO(CORE-858): Re-enable determinism checks once non-determinism issue is found and resolved.
			disableNonDeterminismChecks: true,
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
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).WithNonDeterminismChecksEnabled(!tc.disableNonDeterminismChecks).Build()
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
						AccAddressForSigning: tc.withdrawal.Sender.Owner,
						Gas:                  120_000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					tc.withdrawal,
				)
				checkTxResp := tApp.CheckTx(CheckTx_MsgWithdrawFromSubaccount)
				require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)
			}
			// Advance to the next block, persisting removals in operations queue to state.
			ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})

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
			ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{})

			require.Equal(t, len(tc.orders), len(tc.expectedOrderRemovals))

			// Verify expectations.
			for idx, order := range tc.orders {
				_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
				require.Equal(t, tc.expectedOrderRemovals[idx], !found)
			}
		})
	}
}

func TestOrderRemoval_Invalid(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order

		// Optional withdraw message for under-collateralized tests.
		withdrawal  *sendingtypes.MsgWithdrawFromSubaccount
		priceUpdate map[uint32]ve.VEPricePair

		// Optional field to override MsgProposedOperations to inject invalid order removals
		msgProposedOperations *clobtypes.MsgProposedOperations

		expectedErr string
	}{
		// TODO(CLOB-877): re-enable these tests.
		// "invalid proposal: undercollateralized order removal invalid for fully-filled order": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Alice_Num0_10_000USD,
		// 		constants.Bob_Num0_10_000USD,
		// 	},
		// 	orders: []clobtypes.Order{
		// 		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
		// 		constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
		// 	},
		// priceUpdate: map[uint32]uint64{},
		// 	msgProposedOperations: &clobtypes.MsgProposedOperations{
		// 		OperationsQueue: []clobtypes.OperationRaw{
		// 			clobtestutils.NewMatchOperationRaw(
		// 				&constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell50_Price10_GTBT15,
		// 				[]clobtypes.MakerFill{
		// 					{
		// 						FillAmount:   5,
		// 						MakerOrderId: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.OrderId,
		// 					},
		// 				},
		// 			),
		// 			clobtestutils.NewOrderRemovalOperationRaw(
		// 				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.OrderId,
		// 				clobtypes.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
		// 			),
		// 		},
		// 	},
		// 	expectedErr: "Order is fully filled",
		// },
		// "invalid proposal: order for well collateralized account cannot be removed": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_10000USD,
		// 	},
		// 	orders: []clobtypes.Order{
		// 		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
		// 	},
		// priceUpdate: map[uint32]uint64{},
		// 	msgProposedOperations: &clobtypes.MsgProposedOperations{
		// 		OperationsQueue: []clobtypes.OperationRaw{
		// 			clobtestutils.NewOrderRemovalOperationRaw(
		// 				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
		// 				clobtypes.OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED,
		// 			),
		// 		},
		// 	},
		// 	expectedErr: "Order passes collateralization check",
		// },
		// Re-enable when reduce-only orders are re-enabled.
		// "invalid proposal: valid reduce-only order cannot be removed": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_1BTC_Short,
		// 	},
		// 	orders: []clobtypes.Order{
		// 		constants.LongTermOrder_Carl_Num0_Id2_Clob0_Buy10_Price35_GTB20_RO,
		// 	},
		// priceUpdate: map[uint32]uint64{},
		// 	msgProposedOperations: &clobtypes.MsgProposedOperations{
		// 		OperationsQueue: []clobtypes.OperationRaw{
		// 			clobtestutils.NewOrderRemovalOperationRaw(
		// 				constants.LongTermOrder_Carl_Num0_Id2_Clob0_Buy10_Price35_GTB20_RO.OrderId,
		// 				clobtypes.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY,
		// 			),
		// 		},
		// 	},
		// 	expectedErr: "Order fill must increase position size or change side",
		// },
		// "invalid proposal: non reduce-only order may not be removed with reduce-only reason": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_1BTC_Short,
		// 	},
		// 	orders: []clobtypes.Order{
		// 		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
		// 	},
		// priceUpdate: map[uint32]uint64{},
		// 	msgProposedOperations: &clobtypes.MsgProposedOperations{
		// 		OperationsQueue: []clobtypes.OperationRaw{
		// 			clobtestutils.NewOrderRemovalOperationRaw(
		// 				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
		// 				clobtypes.OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY,
		// 			),
		// 		},
		// 	},
		// 	expectedErrType: clobtypes.ErrInvalidOrderRemoval,
		// 	expectedErr:     "Order must be reduce only",
		// },
		"invalid proposal: conditional fok order cannot be removed when untriggered": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_FOK,
			},
			priceUpdate: map[uint32]ve.VEPricePair{},
			msgProposedOperations: &clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewOrderRemovalOperationRaw(
						constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_FOK.OrderId,
						clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED,
					),
				},
			},
			expectedErr: "does not exist in triggered conditional state",
		},
		"invalid proposal: conditional fok order removal is for non fok order": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
			},
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			msgProposedOperations: &clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewOrderRemovalOperationRaw(
						constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId,
						clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED,
					),
				},
			},
			expectedErr: "Order is not fill-or-kill",
		},
		"invalid proposal: conditional fok order cannot be removed when fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_FOK,
			},
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			msgProposedOperations: &clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewMatchOperationRaw(
						&constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_FOK,
						[]clobtypes.MakerFill{
							{
								FillAmount:   50_000_000,
								MakerOrderId: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							},
						},
					),
					clobtestutils.NewOrderRemovalOperationRaw(
						constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_FOK.OrderId,
						clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED,
					),
				},
			},
			expectedErr: "Fill-or-kill order is fully filled",
		},

		"invalid proposal: conditional ioc order cannot be removed when untriggered": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
			},
			priceUpdate: map[uint32]ve.VEPricePair{},
			msgProposedOperations: &clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewOrderRemovalOperationRaw(
						constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId,
						clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK,
					),
				},
			},
			expectedErr: "does not exist in triggered conditional state",
		},
		"invalid proposal: conditional ioc order removal is for non ioc order": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_FOK,
			},
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			msgProposedOperations: &clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewOrderRemovalOperationRaw(
						constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_FOK.OrderId,
						clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK,
					),
				},
			},
			expectedErr: "Order is not immediate-or-cancel",
		},
		"invalid proposal: conditional ioc order cannot be removed when fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_10000USD,
				constants.Dave_Num0_10000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
			},
			priceUpdate: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_400_000,
					PnlPrice:  5_000_400_000,
				},
			},
			msgProposedOperations: &clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewMatchOperationRaw(
						&constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC,
						[]clobtypes.MakerFill{
							{
								FillAmount:   50_000_000,
								MakerOrderId: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							},
						},
					),
					clobtestutils.NewOrderRemovalOperationRaw(
						constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTBT10_SL_50003_IOC.OrderId,
						clobtypes.OrderRemoval_REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK,
					),
				},
			},
			expectedErr: "Immediate-or-cancel order is fully filled",
		},
		"invalid proposal: post-only removal reason used for non post-only order": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_100_000USD,
			},
			orders: []clobtypes.Order{
				constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15,
			},
			priceUpdate: map[uint32]ve.VEPricePair{},
			msgProposedOperations: &clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewOrderRemovalOperationRaw(
						constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15.OrderId,
						clobtypes.OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER,
					),
				},
			},
			expectedErr: "Order is not post-only",
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
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
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

			// Create all orders and add to deliverTxsOverride
			deliverTxsOverride := make([][]byte, 0)
			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
					deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
				}
			}

			// Add the price update to deliverTxsOverride
			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				map[uint32]ve.VEPricePair{},
				"",
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)

			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

			// Advance to the next block, updating the price.
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})
			// Make sure stateful orders are in state.
			for _, order := range tc.orders {
				_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
				require.True(t, found)
			}

			_, extCommitBz, err = vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdate,
				"",
				2,
			)
			require.NoError(t, err)

			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: [][]byte{extCommitBz},
			})

			_, extCommitBz, err = vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				map[uint32]ve.VEPricePair{},
				"",
				3,
			)
			require.NoError(t, err)

			// Next block will have invalid Order Removals injected in proposal.
			tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: [][]byte{extCommitBz, testtx.MustGetTxBytes(tc.msgProposedOperations)},
				ValidateFinalizeBlock: func(
					ctx sdktypes.Context,
					request abcitypes.RequestFinalizeBlock,
					response abcitypes.ResponseFinalizeBlock,
				) (haltchain bool) {
					execResult := response.TxResults[1]
					require.True(t, execResult.IsErr())
					require.Equal(t, clobtypes.ErrInvalidOrderRemoval.ABCICode(), execResult.Code)
					require.Contains(t, execResult.Log, tc.expectedErr)
					return false
				},
			})
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

		disableNonDeterminismChecks bool
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
				AssetId:   constants.TDai.Id,
				Quantums:  10_000_000_000,
			},

			expectedFirstOrderRemoved:  false,
			expectedSecondOrderRemoved: true, // taker order fails collateralization check during matching
			// TODO(CORE-858): Re-enable determinism checks once non-determinism issue is found and resolved.
			disableNonDeterminismChecks: true,
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
				AssetId:   constants.TDai.Id,
				Quantums:  10_000_000_000,
			},

			expectedFirstOrderRemoved:  false,
			expectedSecondOrderRemoved: true, // taker order fails add-to-orderbook collateralization check
			// TODO(CORE-858): Re-enable determinism checks once non-determinism issue is found and resolved.
			disableNonDeterminismChecks: true,
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
				AssetId:   constants.TDai.Id,
				Quantums:  10_000_000_000,
			},

			expectedFirstOrderRemoved:  true, // maker is under-collateralized
			expectedSecondOrderRemoved: false,
			// TODO(CORE-858): Re-enable determinism checks once non-determinism issue is found and resolved.
			disableNonDeterminismChecks: true,
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
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).WithNonDeterminismChecksEnabled(!tc.disableNonDeterminismChecks).Build()

			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)

			require.NoError(t, conversionErr)

			tApp.App.RatelimitKeeper.SetSDAIPrice(tApp.App.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.App.RatelimitKeeper.SetAssetYieldIndex(tApp.App.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			if !(tApp.CrashingApp == nil) {
				tApp.CrashingApp.RatelimitKeeper.SetSDAIPrice(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), rate)
				tApp.CrashingApp.RatelimitKeeper.SetAssetYieldIndex(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))
			}

			if !(tApp.NoCheckTxApp == nil) {
				tApp.NoCheckTxApp.RatelimitKeeper.SetSDAIPrice(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), rate)
				tApp.NoCheckTxApp.RatelimitKeeper.SetAssetYieldIndex(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))
			}

			if !(tApp.ParallelApp == nil) {
				tApp.ParallelApp.RatelimitKeeper.SetSDAIPrice(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), rate)
				tApp.ParallelApp.RatelimitKeeper.SetAssetYieldIndex(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))
			}

			ctx := tApp.InitChain()

			// Create all orders.
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				tApp.App,
				*clobtypes.NewMsgPlaceOrder(tc.firstOrder),
			) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				tApp.App,
				*clobtypes.NewMsgPlaceOrder(tc.secondOrder),
			) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			// Do the optional withdraw.
			if tc.withdrawal != nil {
				CheckTx_MsgWithdrawFromSubaccount := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: tc.withdrawal.Sender.Owner,
						Gas:                  120_000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					tc.withdrawal,
				)
				checkTxResp := tApp.CheckTx(CheckTx_MsgWithdrawFromSubaccount)
				require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)
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
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
					constants.ClobPair_Btc,
				}
				genesisState.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
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

	// Create a resting order for alice.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15_PO,
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	// Partially match alice's order so that it's in the operations queue.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Bob_Num0_Id1_Clob0_Sell5_Price10_GTBT10,
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	// Now remove alice's order somehow. Self-trade in this case.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Alice_Num0_Id1_Clob0_Sell20_Price10_GTBT10,
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	// Place another order to invalidate Alice's post only order.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Bob_Num0_Id0_Clob0_Sell5_Price5_GTBT10,
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}

	_ = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Local operations queue would be [placement(Alice_Order), ..., removal(Alice_Order)].
	// Let's say block proposer does not include these operations. Make sure we don't panic in this case.
	_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
		&tApp.App.ConsumerKeeper,
		ctx,
		map[uint32]ve.VEPricePair{},
		"",
		tApp.GetHeader().Height,
	)
	require.NoError(t, err)
	_ = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{extCommitBz},
	})
	_ = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
}
