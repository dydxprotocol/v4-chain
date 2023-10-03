package clob_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cometbft/cometbft/types"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testmsgs "github.com/dydxprotocol/v4-chain/protocol/testutil/msgs"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder_StatefulCancelFollowedByPlaceInSameBlockErrorsInCheckTx(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()

	// Place the order.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// We should accept the cancellation since the order exists in state.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgCancelOrderStateful(
			LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.Order.OrderId,
			30,
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
	}
	// We should reject this order since there is already an uncommitted cancellation which
	// we would reject during `DeliverTx`.
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
		require.Contains(t, resp.Log, "An uncommitted stateful order cancellation with this OrderId already exists")
	}

	// Advancing to the next block should succeed and have the order be cancelled and a new one not being placed.
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
	orders := tApp.App.ClobKeeper.GetAllPlacedStatefulOrders(ctx)
	require.NotContains(t, orders, LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.Order.OrderId)
}

func TestCancelStatefulOrder(t *testing.T) {
	tests := map[string]struct {
		blockWithMessages        []testmsgs.TestBlockWithMsgs
		checkCancelledPlaceOrder clobtypes.MsgPlaceOrder
		checkResultsBlock        uint32
	}{
		"Test stateful order is cancelled when placed and cancelled in the same block": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
							ExpectedIsOk: true,
						},
						{
							Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
							ExpectedIsOk: true,
						},
					},
				},
			},

			checkCancelledPlaceOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			checkResultsBlock:        4,
		},
		"Test stateful order is cancelled when placed then cancelled in a future block": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{{
						Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
						ExpectedIsOk: true,
					}},
				},
				{
					Block: 3,
					Msgs: []testmsgs.TestSdkMsg{{
						Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
						ExpectedIsOk: true,
					}},
				},
			},

			checkCancelledPlaceOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			checkResultsBlock:        4,
		},
		"Test stateful order is cancelled when placed, matched, and cancelled in the same block": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
							ExpectedIsOk: true,
						},
						{
							Msg:          &PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
							ExpectedIsOk: true,
						},
						{
							Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
							ExpectedIsOk: true,
						},
					},
				},
			},

			checkCancelledPlaceOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			checkResultsBlock:        4,
		},
		"Test stateful order is cancelled when placed, cancelled, then re-placed with the same order id": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
							ExpectedIsOk: true,
						},
						{
							Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
							ExpectedIsOk: true,
						},
					},
				},
				{
					Block: 3,
					Msgs: []testmsgs.TestSdkMsg{{
						Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
						ExpectedIsOk: true,
					}},
				},
			},

			checkCancelledPlaceOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			checkResultsBlock:        4,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()

			for _, blockWithMessages := range tc.blockWithMessages {
				ctx := tApp.AdvanceToBlock(blockWithMessages.Block, testapp.AdvanceToBlockOptions{})

				for _, testSdkMsg := range blockWithMessages.Msgs {
					result := tApp.CheckTx(testapp.MustMakeCheckTx(
						ctx,
						tApp.App,
						testapp.MustMakeCheckTxOptions{
							AccAddressForSigning: testtx.MustGetOnlySignerAddress(testSdkMsg.Msg),
						},
						testSdkMsg.Msg,
					))
					require.Conditionf(t, func() bool {
						return result.IsOK() == testSdkMsg.ExpectedIsOk
					}, "Expected CheckTx to succeed. Response: %+v", result)
				}
			}

			ctx := tApp.AdvanceToBlock(tc.checkResultsBlock, testapp.AdvanceToBlockOptions{})
			exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(
				ctx,
				tc.checkCancelledPlaceOrder.Order.OrderId,
			)
			require.False(t, exists)
			require.Equal(t, uint64(0), fillAmount.ToUint64())
		})
	}
}

func TestImmediateExecutionLongTermOrders(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()

	// Reject long-term IOC in CheckTx
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_IOC,
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
		require.Contains(t, resp.Log, clobtypes.ErrLongTermOrdersCannotRequireImmediateExecution.Error())
	}

	// Reject long-term FOK in CheckTx
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_FOK,
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
		require.Contains(t, resp.Log, clobtypes.ErrLongTermOrdersCannotRequireImmediateExecution.Error())
	}

	// Reject long-term IOC/FOK in DeliverTx
	blockAdvancement := testapp.BlockAdvancementWithErrors{
		BlockAdvancement: testapp.BlockAdvancement{
			StatefulOrders: []clobtypes.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_IOC,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_FOK,
			},
		},
		ExpectedDeliverTxErrors: map[int]string{
			0: clobtypes.ErrLongTermOrdersCannotRequireImmediateExecution.Error(),
			1: clobtypes.ErrLongTermOrdersCannotRequireImmediateExecution.Error(),
		},
	}
	blockAdvancement.AdvanceToBlock(ctx, 2, &tApp, t)
}

func TestLongTermOrderExpires(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()

	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15
	for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		*clobtypes.NewMsgPlaceOrder(
			MustScaleOrder(order, testapp.DefaultGenesis()),
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.True(t, resp.IsOK(), "Expected CheckTx to succed, but failed: %+v", resp.Log)
	}

	// block time zero, not expired
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
	require.True(t, found, "Order is not expired and should still be in state")

	// block time ten, still not expired
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(10, 0).UTC(),
	})
	_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
	require.True(t, found, "Order is not expired and should still be in state")

	// block time fifteen, expired
	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{
		BlockTime: time.Unix(15, 0).UTC(),
	})
	_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, order.OrderId)
	require.False(t, found, "Order is expired and should not be in state")
}

func TestPlaceLongTermOrder(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().Build()
	ctx := tApp.InitChain()

	// subaccounts for indexer expectation assertions
	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)

	// orders
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *aliceSubaccount.Id,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     10_000_000_000, // 1 BTC, assuming atomic resolution of -10
			Subticks:     500_000_000, // 50k USDC / BTC, assuming QCE of -8
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		},
	)
	PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     10_000_000_000, // 1 BTC, assuming atomic resolution of -10
			Subticks:     500_000_000, // 50k USDC / BTC, assuming QCE of -8
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
	)

	// CheckTx Txs needed for indexer expectation assertions
	CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5),
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(&PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20),
		},
		&PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20,
	)

	tests := map[string]struct {
		makerOrders []clobtypes.MsgPlaceOrder
		takerOrder  clobtypes.MsgPlaceOrder

		expectedOffchainMessagesAfterCheckTx []msgsender.Message
		expectedOffchainMessagesInFirstBlock []msgsender.Message
		expectedOnchainMessagesInFirstBlock []msgsender.Message
		// expectedOffchainMessagesInSecondBlock []msgsender.Message
		expectedOnchainMessagesInSecondBlock []msgsender.Message
		// expectedFillAmount int
	} {
		"Test placing an order": {
			takerOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
			// No offchain messages in CheckTx because stateful MsgPlaceOrder is not placed in CheckTx
			expectedOffchainMessagesAfterCheckTx: []msgsender.Message{},
			// Order update message, note order place messages are skipped for stateful orders
			expectedOffchainMessagesInFirstBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.OrderId,
					0,
				),
			},
			// Stateful order placement event is an onchain event
			expectedOnchainMessagesInFirstBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						{
							Subtype: indexerevents.SubtypeStatefulOrder,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewLongTermOrderPlacementEvent(
									LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.StatefulOrderEventVersion,
						},
					},
					TxHashes: []string{
						string(lib.GetTxHash(
							CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Tx,
						)),
					},
				},
			)},
		},
		"Test matching an order fully as taker": {
			makerOrders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20,
			},
			takerOrder: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
			// Short term order placement results in Create and Update with 0 fill amount
			expectedOffchainMessagesAfterCheckTx: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20.Tx),
				}),
			},
			// Short term order update for fill amount, stateful order update for fill amount
			// Note there are no headers because these events are generated in PrepareCheckState
			expectedOffchainMessagesInFirstBlock: []msgsender.Message{
				// maker
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
				),
				// taker
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.OrderId,
					LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
				),
			},
			// Stateful order placement
			expectedOnchainMessagesInFirstBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						{
							Subtype: indexerevents.SubtypeStatefulOrder,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewLongTermOrderPlacementEvent(
									LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.StatefulOrderEventVersion,
						},
					},
					TxHashes: []string{
						string(lib.GetTxHash(
							CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Tx,
						)),
					},
				},
			)},
			expectedOnchainMessagesInSecondBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 3,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						// taker subaccount state transition
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(int64(
												LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
											FundingIndex: dtypes.NewInt(0),
										},
									},
									[]*satypes.AssetPosition{
										{
											AssetId:  lib.UsdcAssetId,
											Quantums: dtypes.NewIntFromBigInt(
												new(big.Int).Sub(
													aliceSubaccount.GetUsdcPosition(),
													new(big.Int).SetInt64(
														50_000_000_000 + 25_000_000, // taker fee of .5%
													),
												),
											),
										},
									},
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
							Version:             indexerevents.SubaccountUpdateEventVersion,
						},
						// maker subaccount state transition
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(-int64(
												LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
											FundingIndex: dtypes.NewInt(0),
										},
									},
									[]*satypes.AssetPosition{
										{
											AssetId:  lib.UsdcAssetId,
											Quantums: dtypes.NewIntFromBigInt(
												new(big.Int).Add(
													bobSubaccount.GetUsdcPosition(),
													new(big.Int).SetInt64(
														50_000_000_000 + 5_500_000, // maker rebate of .110%
													),
												),
											),
										},
									},
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
							Version:             indexerevents.SubaccountUpdateEventVersion,
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
									LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
									LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
									-5_500_000,
									25_000_000,
									LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
									LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
							Version:             indexerevents.OrderFillEventVersion,
						},
					},
					TxHashes: []string{
						string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
							OperationsQueue: []clobtypes.OperationRaw{
								{
									Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
										ShortTermOrderPlacement: CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20.Tx,
									},
								},
								clobtestutils.NewMatchOperationRaw(
									&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
									[]clobtypes.MakerFill{
										{
											FillAmount: PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.
												Order.GetBaseQuantums().ToUint64(),
											MakerOrderId: PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
										},
									},
								),
							},
						}))),
					},
				},
			)},
		},
		// "Test matching an order partially, maker order remains on book": {

		// },
		// "Test matching an order partially, taker order placed on book": {

		// },
		// "Test matching an order partially as taker then fully as maker": {

		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// testapp.UpdateGenesisDocWithAppStateForModule(
				// 	&genesis,
				// 	func(genesisState *satypes.GenesisState) {
				// 		genesisState.Subaccounts = tc.subaccounts
				// 	},
				// )
				return genesis
			}).WithTesting(t).WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts)).Build()
			ctx := tApp.InitChain()

			// Place makers
			for _, order := range tc.makerOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					order,
				) {
					resp := tApp.CheckTx(checkTx)
					require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}

			// Place taker
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
				ctx,
				tApp.App,
				tc.takerOrder,
			) {
				resp := tApp.CheckTx(checkTx)
				require.True(t, resp.IsOK(), "Expected CheckTx to succeed. Response: %+v", resp)
			}

			require.ElementsMatch(
				t,
				tc.expectedOffchainMessagesAfterCheckTx,
				msgSender.GetOffchainMessages(),
			)
			msgSender.Clear()

			// places short term makers on the book and writes long term orders to state
			tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			require.ElementsMatch(
				t,
				tc.expectedOnchainMessagesInFirstBlock,
				msgSender.GetOnchainMessages(),
			)
			// orders are immediately placed in PrepareCheckState, triggering offchain messages
			require.ElementsMatch(
				t,
				tc.expectedOffchainMessagesInFirstBlock,
				msgSender.GetOffchainMessages(),
			)
			msgSender.Clear()

			// matches generated in PrepareCheckState are proposed
			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			require.ElementsMatch(
				t,
				tc.expectedOnchainMessagesInSecondBlock,
				msgSender.GetOnchainMessages(),
			)

			// require.True(t, false)

			// // Verify expectations
			// // IOC orders should not have remaining size placed as makers
			// _, found := tApp.App.ClobKeeper.MemClob.GetOrder(ctx, order.OrderId)
			// require.True(t, found, "Partially filled order should be on the orderbook")

			// // Fill amount should be 50
			// _, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, order.OrderId)
			// require.Equal(t, 50, fillAmount, "Fill amount should be 50, not %d", fillAmount)
		})
	}
}
