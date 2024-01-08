package clob_test

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestReduceOnlyOrders(t *testing.T) {
	tests := map[string]struct {
		subaccounts          []satypes.Subaccount
		ordersForFirstBlock  []clobtypes.Order
		ordersForSecondBlock []clobtypes.Order

		priceUpdateForFirstBlock  *prices.MsgUpdateMarketPrices
		priceUpdateForSecondBlock *prices.MsgUpdateMarketPrices

		crashingAppCheckTxNonDeterminsmChecksDisabled bool

		expectedInTriggeredStateAfterBlock map[uint32]map[clobtypes.OrderId]bool

		expectedOrderOnMemClob  map[clobtypes.OrderId]bool
		expectedOrderFillAmount map[clobtypes.OrderId]uint64
		expectedSubaccounts     []satypes.Subaccount
	}{
		/*
			"IOC Reduce only order partially matches short term order same block, maker order fully filled": {
				subaccounts: []satypes.Subaccount{
					constants.Carl_Num0_100000USD,
					constants.Alice_Num1_1BTC_Long_500_000USD,
				},
				ordersForFirstBlock: []clobtypes.Order{
					MustScaleOrder(
						constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20,
						testapp.DefaultGenesis(),
					),
					MustScaleOrder(
						constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
						testapp.DefaultGenesis(),
					),
				},
				ordersForSecondBlock: []clobtypes.Order{},

				expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          false,
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
				},
				expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          100,
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 100,
				},
				expectedSubaccounts: []satypes.Subaccount{
					{
						Id: &constants.Carl_Num0,
						AssetPositions: []*satypes.AssetPosition{
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(95_000_550_000),
							},
						},
						PerpetualPositions: []*satypes.PerpetualPosition{
							{
								PerpetualId:  0,
								Quantums:     dtypes.NewInt(100),
								FundingIndex: dtypes.NewInt(0),
							},
						},
					},
					{
						Id: &constants.Alice_Num1,
						AssetPositions: []*satypes.AssetPosition{
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(504_997_500_000),
							},
						},
						PerpetualPositions: []*satypes.PerpetualPosition{
							{
								PerpetualId:  0,
								Quantums:     dtypes.NewInt(99_999_900),
								FundingIndex: dtypes.NewInt(0),
							},
						},
					},
				},
			},
			"IOC Reduce only order partially matches short term order second block, maker order fully filled": {
				subaccounts: []satypes.Subaccount{
					constants.Carl_Num0_100000USD,
					constants.Alice_Num1_1BTC_Long_500_000USD,
				},
				ordersForFirstBlock: []clobtypes.Order{
					MustScaleOrder(
						constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20,
						testapp.DefaultGenesis(),
					),
				},
				ordersForSecondBlock: []clobtypes.Order{
					MustScaleOrder(
						constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
						testapp.DefaultGenesis(),
					),
				},

				expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          false,
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
				},
				expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          100,
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 100,
				},
				expectedSubaccounts: []satypes.Subaccount{
					{
						Id: &constants.Carl_Num0,
						AssetPositions: []*satypes.AssetPosition{
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(95_000_550_000),
							},
						},
						PerpetualPositions: []*satypes.PerpetualPosition{
							{
								PerpetualId:  0,
								Quantums:     dtypes.NewInt(100),
								FundingIndex: dtypes.NewInt(0),
							},
						},
					},
					{
						Id: &constants.Alice_Num1,
						AssetPositions: []*satypes.AssetPosition{
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(504_997_500_000),
							},
						},
						PerpetualPositions: []*satypes.PerpetualPosition{
							{
								PerpetualId:  0,
								Quantums:     dtypes.NewInt(99_999_900),
								FundingIndex: dtypes.NewInt(0),
							},
						},
					},
				},
			},
			"IOC Reduce only order partially matches short term order second block, maker order partially filled": {
				subaccounts: []satypes.Subaccount{
					constants.Carl_Num0_100000USD,
					constants.Alice_Num1_1BTC_Long_500_000USD,
				},
				ordersForFirstBlock: []clobtypes.Order{
					MustScaleOrder(
						constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20,
						testapp.DefaultGenesis(),
					),
				},
				ordersForSecondBlock: []clobtypes.Order{
					MustScaleOrder(
						constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
						testapp.DefaultGenesis(),
					),
				},

				expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          true,
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
				},
				expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          150,
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 150,
				},
				expectedSubaccounts: []satypes.Subaccount{
					{
						Id: &constants.Carl_Num0,
						AssetPositions: []*satypes.AssetPosition{
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(9_250_0825_000),
							},
						},
						PerpetualPositions: []*satypes.PerpetualPosition{
							{
								PerpetualId:  0,
								Quantums:     dtypes.NewInt(150),
								FundingIndex: dtypes.NewInt(0),
							},
						},
					},
					{
						Id: &constants.Alice_Num1,
						AssetPositions: []*satypes.AssetPosition{
							{
								AssetId:  0,
								Quantums: dtypes.NewInt(507_496_250_000),
							},
						},
						PerpetualPositions: []*satypes.PerpetualPosition{
							{
								PerpetualId:  0,
								Quantums:     dtypes.NewInt(99_999_850),
								FundingIndex: dtypes.NewInt(0),
							},
						},
					},
				},
			},
		*/
		"FOK Reduce only order fully matches short term order second block, maker order partially filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},

			// Crashing app checks have to be disabled because the FOK order will not match
			// with an empty orderbook and fail to be placed.
			crashingAppCheckTxNonDeterminsmChecksDisabled: true,

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          true,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_FOK_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId: 150,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(9_250_0825_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(150),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(507_496_250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_850),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
			},
		},
		"FOK Reduce only order fully matches short term order same block, maker order partially filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{},

			// Crashing app checks don't need to be disabled since matches occur in same block.
			crashingAppCheckTxNonDeterminsmChecksDisabled: false,

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          true,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_FOK_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId: 150,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(9_250_0825_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(150),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(507_496_250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_850),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
			},
		},
		"Conditional FOK Reduce only order fully matches short term order same block, maker order partially filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
				satypes.Subaccount{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_500_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(100_000_000), // +1 BTC
						},
					},
				},
			},
			// ordersForFirstBlock: []clobtypes.Order{
			// 	// Conditional order will be placed and triggered in same block.
			// 	MustScaleOrder(
			// 		constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO,
			// 		testapp.DefaultGenesis(),
			// 	),
			// 	// This short term order will match the conditional order after it is triggered.
			// 	MustScaleOrder(
			// 		constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10,
			// 		testapp.DefaultGenesis(),
			// 	),
			// },
			ordersForFirstBlock: []clobtypes.Order{
				// constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO,
				// constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10,
				constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999,
			},
			ordersForSecondBlock: []clobtypes.Order{},

			// priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{

			// 	MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
			// 		MustScaleMarketPriceUpdate(
			// 			*prices.NewMarketPriceUpdate(0, 5_000_300_000),
			// 			testapp.DefaultGenesis(),
			// 		),
			// 	},
			// },
			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					// prices.NewMarketPriceUpdate(0, 5_000_300_000),
					prices.NewMarketPriceUpdate(0, 4_999_700_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_TP_49999.OrderId: true,
					// constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: true,
				},
			},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10.OrderId: true,
				// constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId: 150,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(9_250_0825_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(150),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(507_496_250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_850),
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithCrashingAppCheckTxNonDeterminismChecksEnabled(!tc.crashingAppCheckTxNonDeterminsmChecksDisabled).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *satypes.GenesisState) {
							genesisState.Subaccounts = tc.subaccounts
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

			for _, order := range tc.ordersForFirstBlock {
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
			if tc.priceUpdateForFirstBlock != nil {
				fmt.Println("price updated")
				txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
				require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForFirstBlock))
				priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
				require.NoError(t, err)
				deliverTxsOverride = append(deliverTxsOverride, priceUpdateTxBytes)
			}

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			prices := tApp.App.PricesKeeper.GetAllMarketPrices(ctx)
			fmt.Println("the prices", prices)

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[2]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 2)
				}
			}

			// Create all orders.
			deliverTxsOverride = make([][]byte, 0)
			deliverTxsOverride = append(
				deliverTxsOverride,
				constants.ValidEmptyMsgProposedOperationsTxBytes,
			)

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

			// Add an empty premium vote.
			deliverTxsOverride = append(deliverTxsOverride, constants.EmptyMsgAddPremiumVotesTxBytes)

			// Add the price update.
			if tc.priceUpdateForSecondBlock != nil {
				txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
				require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForSecondBlock))
				priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
				require.NoError(t, err)
				deliverTxsOverride = append(deliverTxsOverride, priceUpdateTxBytes)
			}

			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[3]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 3)
				}
			}

			// Verify expectations.
			for orderId, exists := range tc.expectedOrderOnMemClob {
				_, existsOnMemclob := tApp.App.ClobKeeper.MemClob.GetOrder(ctx, orderId)
				require.Equal(t, exists, existsOnMemclob)
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

func TestReduceOnlyOrderFailure(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order
		errorMsg    []string
	}{
		"Zero perpetual position subaccount position cannot place sell RO order": {
			orders: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		"Zero perpetual position subaccount position cannot place buy RO order": {
			orders: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Buy10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		"Reduce only order fails to replace non-reduce only order": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_1BTC_Short_100_000USD,
			},
			orders: []clobtypes.Order{
				// non reduce only order
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20,
					testapp.DefaultGenesis(),
				),
				// reduce only replacement fails
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Buy10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				"",
				clobtypes.ErrReduceOnlyOrderReplacement.Error(),
			},
		},
		"FOK Reduce only order is placed but does not match immediately and is cancelled.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_1BTC_Short_100_000USD,
			},
			orders: []clobtypes.Order{
				MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Buy10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrFokOrderCouldNotBeFullyFilled.Error(),
			},
		},
		"Conditional FOK Reduce only order is placed successfully.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_1BTC_Short_100_000USD,
			},
			orders: []clobtypes.Order{
				MustScaleOrder(
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				"",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				if len(tc.subaccounts) > 0 {
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *satypes.GenesisState) {
							genesisState.Subaccounts = tc.subaccounts
						},
					)
				}
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			for idx, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)

					if tc.errorMsg[idx] == "" {
						require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
					} else {
						require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
						require.Contains(
							t,
							resp.Log,
							tc.errorMsg[idx],
						)
					}
				}
			}
		})
	}
}
