package clob_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
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

		crashingAppCheckTxNonDeterminismChecksDisabled bool

		expectedInTriggeredStateAfterBlock map[uint32]map[clobtypes.OrderId]bool

		expectedOrderOnMemClob  map[clobtypes.OrderId]bool
		expectedOrderFillAmount map[clobtypes.OrderId]uint64
		expectedSubaccounts     []satypes.Subaccount
	}{
		"IOC Reduce only order partially matches short term order same block, maker order fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
				testapp.MustScaleOrder(
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
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(95_000_550_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(100),
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
							big.NewInt(504_997_500_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(99_999_900),
							big.NewInt(0),
							big.NewInt(0),
						),
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
				testapp.MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
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
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(95_000_550_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(100),
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
							big.NewInt(504_997_500_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(99_999_900),
							big.NewInt(0),
							big.NewInt(0),
						),
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
				testapp.MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
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
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(9_250_0825_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(150),
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
							big.NewInt(507_496_250_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(99_999_850),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		"Conditional IOC Reduce only order partially matches short term order same block, maker order fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_500000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO,
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10,
			},
			ordersForSecondBlock: []clobtypes.Order{},

			priceUpdateForFirstBlock: &prices.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*prices.MsgUpdateMarketPrices_MarketPrice{
					prices.NewMarketPriceUpdate(0, 5_000_300_000),
				},
			},
			priceUpdateForSecondBlock: &prices.MsgUpdateMarketPrices{},

			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: true,
				},
				3: {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: true,
				},
			},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10.OrderId:                              false,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10.OrderId:                              25_000_000,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: 25_000_000,
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
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(75_000_000),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithCrashingAppCheckTxNonDeterminismChecksEnabled(!tc.crashingAppCheckTxNonDeterminismChecksDisabled).
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
				txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
				require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForFirstBlock))
				priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
				require.NoError(t, err)
				deliverTxsOverride = append(deliverTxsOverride, priceUpdateTxBytes)
			}

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

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
				_, existsOnMemclob := tApp.App.ClobKeeper.MemClob.GetOrder(orderId)
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
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		"Zero perpetual position subaccount position cannot place buy RO order": {
			orders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Buy10_Price15_GTB20_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		"Regular Reduce only order fails because disabled": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_1BTC_Short_100_000USD,
			},
			orders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell10_Price15_GTB20_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyDisabled.Error(),
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

func TestClosePositionOrder(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	CheckTx_PlaceOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Carl_Num0.Owner,
		},
		clobtypes.NewMsgPlaceOrder(
			constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
		),
	)
	CheckTx_PlaceOrder_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Carl_Num0.Owner,
		},
		clobtypes.NewMsgPlaceOrder(
			constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10,
		),
	)
	CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num1.Owner,
		},
		clobtypes.NewMsgPlaceOrder(
			constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO,
		),
	)

	tests := map[string]struct {
		subaccounts          []satypes.Subaccount
		orders               []clobtypes.Order
		matchIncludedInBlock bool

		expectedOrderOnMemClob   map[clobtypes.OrderId]bool
		expectedOrderFillAmount  map[clobtypes.OrderId]uint64
		expectedSubaccounts      []satypes.Subaccount
		expectedOffchainMessages []msgsender.Message
	}{
		"Close position order (IOC reduce-only) fully filled, maker order fully filled, match in block": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				// Initialize Alice subaccount 1 with a 1 BTC long position.
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			orders: []clobtypes.Order{
				// 1. an order from Carl that buys 1 BTC at price 50_000
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				// 2. an order from Alice that closes their position by selling 1 BTC at price 50_000
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO,
			},
			matchIncludedInBlock: true,

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId:          false,
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId:          100_000_000,
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId: 100_000_000,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							// 100_000 usdc - 50_000 (from buying 1 btc) + 5.5 usdc maker fee
							big.NewInt(50_005_500_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(100_000_000),
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
							// 500_000 usdc + 50_000 (from selling 1 btc) - 25 usdc taker fee
							big.NewInt(549_975_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{},
				},
			},
			expectedOffchainMessages: []msgsender.Message{
				// 1. Order place of Carl's order
				// 2. Order update of Carl's order with 0 fill amount
				// 3. Order place of Alice's order
				// 4. Order update that Carl's order is fully filled
				// 5. Order update that Alice's order is fully filled
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId,
					100_000_000,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId,
					100_000_000,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
			},
		},
		"Close position order (IOC reduce-only) fully filled, maker order partially filled, match in block": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				// Initialize Alice subaccount 1 with a 1 BTC long position.
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			orders: []clobtypes.Order{
				// 1. an order from Carl that buys 2 BTC at price 50_000
				constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10,
				// 2. an order from Alice that closes their position by selling 1 BTC at price 50_000
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO,
			},
			matchIncludedInBlock: true,

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10.OrderId:          true,
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10.OrderId:          100_000_000,
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId: 100_000_000,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							// 100_000 usdc - 50_000 (from buying 1 btc) + 5.5 usdc maker fee
							big.NewInt(50_005_500_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(100_000_000),
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
							// 500_000 usdc + 50_000 (from selling 1 btc) - 25 usdc taker fee
							big.NewInt(549_975_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{},
				},
			},
			expectedOffchainMessages: []msgsender.Message{
				// 1. Order place of Carl's order
				// 2. Order update of Carl's order with 0 fill amount
				// 3. Order place of Alice's order
				// 4. Order update that Carl's order is partially filled
				// 5. Order update that Alice's order is fully filled
				// 6. Order update that Carl's order is partially filled (due to operations being replayed)
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10.OrderId,
					100_000_000,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId,
					100_000_000,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price50000_GTB10.OrderId,
					100_000_000,
				),
			},
		},
		"Close position order (IOC reduce-only) fully filled, maker order fully filled, match not in block": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				// Initialize Alice subaccount 1 with a 1 BTC long position.
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			orders: []clobtypes.Order{
				// 1. an order from Carl that buys 1 BTC at price 50_000
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				// 2. an order from Alice that closes their position by selling 1 BTC at price 50_000
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO,
			},
			matchIncludedInBlock: false,

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId:          false,
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId:          100_000_000,
				constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId: 100_000_000,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							// 100_000 usdc - 50_000 (from buying 1 btc) + 5.5 usdc maker fee
							big.NewInt(50_005_500_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(100_000_000),
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
							// 500_000 usdc + 50_000 (from selling 1 btc) - 25 usdc taker fee
							big.NewInt(549_975_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{},
				},
			},
			expectedOffchainMessages: []msgsender.Message{
				// 1. Order place of Carl's order
				// 2. Order update of Carl's order with 0 fill amount
				// 3. Order place of Alice's order
				// 4. Order update that Carl's order is fully filled
				// 5. Order update that Alice's order is fully filled
				// 6. Order update that Carl's order is fully filled (due to operations being replayed as match wasn't in block)
				// 7. Order update that Alice's order is fully filled (due to operations being replayed as match wasn't in block)
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId,
					100_000_000,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId,
					100_000_000,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId,
					100_000_000,
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					ctx,
					constants.Order_Alice_Num1_Id0_Clob0_Sell1BTC_Price50000_GTB20_IOC_RO.OrderId,
					100_000_000,
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()

			tApp := testapp.NewTestAppBuilder(t).
				WithAppOptions(map[string]interface{}{
					indexer.MsgSenderInstanceForTest: msgSender,
				}).
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

			// Place orders.
			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			options := testapp.AdvanceToBlockOptions{}
			if !tc.matchIncludedInBlock {
				options.DeliverTxsOverride = [][]byte{}
			}
			ctx = tApp.AdvanceToBlock(2, options)

			// Verify expectations.
			for orderId, exists := range tc.expectedOrderOnMemClob {
				_, existsOnMemclob := tApp.App.ClobKeeper.MemClob.GetOrder(orderId)
				require.Equal(t, exists, existsOnMemclob)
			}

			for orderId, expectedFillAmount := range tc.expectedOrderFillAmount {
				exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				require.True(t, exists)
				require.Equal(t, expectedFillAmount, fillAmount.ToUint64())
			}

			require.ElementsMatch(
				t,
				tc.expectedOffchainMessages,
				msgSender.GetOffchainMessages(),
			)
		})
	}
}
