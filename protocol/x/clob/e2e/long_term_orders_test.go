package clob_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/tmhash"

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
	tApp := testapp.NewTestAppBuilder(t).Build()
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
			tApp := testapp.NewTestAppBuilder(t).Build()

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
	tApp := testapp.NewTestAppBuilder(t).Build()
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
	blockAdvancement.AdvanceToBlock(ctx, 2, tApp, t)
}

func TestLongTermOrderExpires(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
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
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// subaccounts for indexer expectation assertions
	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)

	// order msgs
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
			Subticks:     500_000_000,    // 50k USDC / BTC, assuming QCE of -8
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		},
	)
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *aliceSubaccount.Id,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     20_000_000_000,
			Subticks:     500_000_000,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		},
	)
	PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     10_000_000_000,
			Subticks:     500_000_000,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
	)
	PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 1, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     10_000_000_000,
			Subticks:     500_000_000,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
	)
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *aliceSubaccount.Id,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     10_000_000_000,
			Subticks:     499_990_000,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
			TimeInForce:  clobtypes.Order_TIME_IN_FORCE_POST_ONLY,
		},
	)

	// CheckTx Txs needed for indexer expectation assertions
	CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
			),
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
	)
	CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5,
			),
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(&PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20),
		},
		&PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id1_Sell1_Price50000_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(&PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20),
		},
		&PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20,
	)
	CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO,
			),
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO,
	)

	type ordersAndExpectations struct {
		orderMsgs   []clobtypes.MsgPlaceOrder
		blockHeight uint32

		expectedOffchainMessagesCheckTx    []msgsender.Message
		expectedOffchainMessagesAfterBlock []msgsender.Message
		expectedOnchainMessagesAfterBlock  []msgsender.Message
	}

	tests := map[string]struct {
		// Long-term order to track
		order clobtypes.Order

		// Orders to place in each block and expectations to verify
		ordersAndExpectationsPerBlock []ordersAndExpectations

		// Expectations to verify at end of test
		orderShouldRestOnOrderbook bool
		expectedOrderFillAmount    uint64
		expectedSubaccounts        []satypes.Subaccount
	}{
		"Test placing an order": {
			order:                      LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
			orderShouldRestOnOrderbook: true,
			expectedOrderFillAmount:    0,
			expectedSubaccounts:        []satypes.Subaccount{aliceSubaccount},

			ordersAndExpectationsPerBlock: []ordersAndExpectations{
				{
					blockHeight: 2,
					orderMsgs: []clobtypes.MsgPlaceOrder{
						LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
					},
					// No offchain messages in CheckTx because stateful MsgPlaceOrder is not placed in CheckTx
					expectedOffchainMessagesCheckTx: []msgsender.Message{},
					// Order update message, note order place messages are skipped for stateful orders
					expectedOffchainMessagesAfterBlock: []msgsender.Message{
						off_chain_updates.MustCreateOrderUpdateMessage(
							nil,
							LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.OrderId,
							0,
						),
					},
					// Stateful order placement event is an onchain event
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 2,
							Time:   ctx.BlockTime(),
							Events: []*indexer_manager.IndexerTendermintEvent{
								{
									Subtype:             indexerevents.SubtypeStatefulOrder,
									OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
									EventIndex:          0,
									Version:             indexerevents.StatefulOrderEventVersion,
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewLongTermOrderPlacementEvent(
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
										),
									),
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
				// No matches generated, empty set of events
				{
					blockHeight: 3,
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 3,
							Time:   ctx.BlockTime(),
							Events: []*indexer_manager.IndexerTendermintEvent{},
						},
					)},
				},
			},
		},
		"Test matching an order fully as taker": {
			order:                      LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
			orderShouldRestOnOrderbook: false,
			expectedOrderFillAmount:    0, // order is fully-filled and removed from state
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: Clob_0.MustGetPerpetualId(),
							Quantums: dtypes.NewInt(int64(
								LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
							FundingIndex: dtypes.NewInt(0),
						},
					},
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId: 0,
							Quantums: dtypes.NewIntFromBigInt(
								new(big.Int).Sub(
									aliceSubaccount.GetUsdcPosition(),
									new(big.Int).SetInt64(
										50_000_000_000+25_000_000, // taker fee of .5%
									),
								),
							),
						},
					},
					MarginEnabled: true,
				},
				{
					Id: &constants.Bob_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: Clob_0.MustGetPerpetualId(),
							Quantums: dtypes.NewInt(-int64(
								LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
							FundingIndex: dtypes.NewInt(0),
						},
					},
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId: 0,
							Quantums: dtypes.NewIntFromBigInt(
								new(big.Int).Add(
									bobSubaccount.GetUsdcPosition(),
									new(big.Int).SetInt64(
										50_000_000_000+5_500_000, // maker rebate of .110%
									),
								),
							),
						},
					},
					MarginEnabled: true,
				},
			},
			ordersAndExpectationsPerBlock: []ordersAndExpectations{
				{
					blockHeight: 2,
					orderMsgs: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20,
						LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
					},
					// Short term order placement results in Create and Update with 0 fill amount
					expectedOffchainMessagesCheckTx: []msgsender.Message{
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
					expectedOffchainMessagesAfterBlock: []msgsender.Message{
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
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 2,
							Time:   ctx.BlockTime(),
							Events: []*indexer_manager.IndexerTendermintEvent{
								{
									Subtype:             indexerevents.SubtypeStatefulOrder,
									OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
									EventIndex:          0,
									Version:             indexerevents.StatefulOrderEventVersion,
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewLongTermOrderPlacementEvent(
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
										),
									),
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
				{
					blockHeight: 3,
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 3,
							Time:   ctx.BlockTime(),
							Events: []*indexer_manager.IndexerTendermintEvent{
								// taker subaccount state transition
								{
									Subtype: indexerevents.SubtypeSubaccountUpdate,
									DataBytes: indexer_manager.GetBytes(
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
													AssetId: 0,
													Quantums: dtypes.NewIntFromBigInt(
														new(big.Int).Sub(
															aliceSubaccount.GetUsdcPosition(),
															new(big.Int).SetInt64(
																50_000_000_000+25_000_000, // taker fee of .5%
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
									DataBytes: indexer_manager.GetBytes(
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
													AssetId: 0,
													Quantums: dtypes.NewIntFromBigInt(
														new(big.Int).Add(
															bobSubaccount.GetUsdcPosition(),
															new(big.Int).SetInt64(
																50_000_000_000+5_500_000, // maker rebate of .110%
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
									DataBytes: indexer_manager.GetBytes(
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
			},
		},
		"Test post-only order placed on the book": {
			order:                      LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO.Order,
			orderShouldRestOnOrderbook: true,
			expectedOrderFillAmount:    0,
			expectedSubaccounts:        []satypes.Subaccount{aliceSubaccount, bobSubaccount},

			ordersAndExpectationsPerBlock: []ordersAndExpectations{
				{
					blockHeight: 2,
					orderMsgs: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20,
						LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO,
					},
					expectedOffchainMessagesCheckTx: []msgsender.Message{
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
					expectedOffchainMessagesAfterBlock: []msgsender.Message{
						// post-only shouldn't match and will have 0 fill size in update message
						off_chain_updates.MustCreateOrderUpdateMessage(
							nil,
							LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO.Order.OrderId,
							0,
						),
					},
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 2,
							Time:   ctx.BlockTime(),
							Events: []*indexer_manager.IndexerTendermintEvent{
								{
									Subtype:             indexerevents.SubtypeStatefulOrder,
									OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
									EventIndex:          0,
									Version:             indexerevents.StatefulOrderEventVersion,
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewLongTermOrderPlacementEvent(
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO.Order,
										),
									),
								},
							},
							TxHashes: []string{
								string(lib.GetTxHash(
									CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO.Tx,
								)),
							},
						},
					)},
				},
				{
					blockHeight: 3,
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 3,
							Time:   ctx.BlockTime(),
						},
					)},
				},
			},
		},
		"Test matching an order partially as taker then fully as maker": {
			order:                      LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order,
			orderShouldRestOnOrderbook: false,
			// order is fully-filled and removed from state, resulting in zero fill amount in state
			expectedOrderFillAmount: 0,
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: Clob_0.MustGetPerpetualId(),
							Quantums: dtypes.NewInt(int64(
								LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.GetQuantums())),
							FundingIndex: dtypes.NewInt(0),
						},
					},
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId: 0,
							Quantums: dtypes.NewIntFromBigInt(
								new(big.Int).Sub(
									aliceSubaccount.GetUsdcPosition(),
									new(big.Int).SetInt64(
										50_000_000_000+25_000_000+ // taker fee of .5%
											50_000_000_000-5_500_000, // maker rebate of .110%
									),
								),
							),
						},
					},
					MarginEnabled: true,
				},
				{
					Id: &constants.Bob_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: Clob_0.MustGetPerpetualId(),
							Quantums: dtypes.NewInt(-int64(
								PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums() +
									PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums(),
							)),
							FundingIndex: dtypes.NewInt(0),
						},
					},
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId: 0,
							Quantums: dtypes.NewIntFromBigInt(
								new(big.Int).Add(
									bobSubaccount.GetUsdcPosition(),
									new(big.Int).SetInt64(
										50_000_000_000+5_500_000+ // maker rebate of .110% from first order
											50_000_000_000-25_000_000, // taker fee of .5% from second order
									),
								),
							),
						},
					},
					MarginEnabled: true,
				},
			},

			ordersAndExpectationsPerBlock: []ordersAndExpectations{
				{
					// Short term order placed in CheckTx, then long term order placed in PrepareCheckState after advancing block
					blockHeight: 2,
					orderMsgs: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20,
						LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5,
					},
					// Short term order placement results in Create and Update with 0 fill amount
					expectedOffchainMessagesCheckTx: []msgsender.Message{
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
					expectedOffchainMessagesAfterBlock: []msgsender.Message{
						// maker fully filled
						off_chain_updates.MustCreateOrderUpdateMessage(
							nil,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
						),
						// taker, partially filled
						off_chain_updates.MustCreateOrderUpdateMessage(
							nil,
							LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.OrderId,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
						),
					},
					// Stateful order placement
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 2,
							Time:   ctx.BlockTime(),
							Events: []*indexer_manager.IndexerTendermintEvent{
								{
									Subtype:             indexerevents.SubtypeStatefulOrder,
									OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
									EventIndex:          0,
									Version:             indexerevents.StatefulOrderEventVersion,
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewLongTermOrderPlacementEvent(
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order,
										),
									),
								},
							},
							TxHashes: []string{
								string(lib.GetTxHash(
									CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Tx,
								)),
							},
						},
					)},
				},
				{
					// Result of partial match of long-term taker order and short-term maker order
					blockHeight: 3,
					expectedOffchainMessagesAfterBlock: []msgsender.Message{
						// attempt to replay the stateful order in PrepareCheckState after advancing the block, fill amount
						// will stay constant
						off_chain_updates.MustCreateOrderUpdateMessage(
							nil,
							LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.OrderId,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
						),
					},
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 3,
							Time:   ctx.BlockTime(),
							Events: []*indexer_manager.IndexerTendermintEvent{
								// taker subaccount state transition
								{
									Subtype: indexerevents.SubtypeSubaccountUpdate,
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewSubaccountUpdateEvent(
											&constants.Alice_Num0,
											[]*satypes.PerpetualPosition{
												{
													PerpetualId: Clob_0.MustGetPerpetualId(),
													Quantums: dtypes.NewInt(int64(
														PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums())),
													FundingIndex: dtypes.NewInt(0),
												},
											},
											[]*satypes.AssetPosition{
												{
													AssetId: 0,
													Quantums: dtypes.NewIntFromBigInt(
														new(big.Int).Sub(
															aliceSubaccount.GetUsdcPosition(),
															new(big.Int).SetInt64(
																50_000_000_000+25_000_000, // taker fee of .5%
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
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewSubaccountUpdateEvent(
											&constants.Bob_Num0,
											[]*satypes.PerpetualPosition{
												{
													PerpetualId: Clob_0.MustGetPerpetualId(),
													Quantums: dtypes.NewInt(-int64(
														PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums())),
													FundingIndex: dtypes.NewInt(0),
												},
											},
											[]*satypes.AssetPosition{
												{
													AssetId: 0,
													Quantums: dtypes.NewIntFromBigInt(
														new(big.Int).Add(
															bobSubaccount.GetUsdcPosition(),
															new(big.Int).SetInt64(
																50_000_000_000+5_500_000, // maker rebate of .110%
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
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewOrderFillEvent(
											PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order,
											PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
											-5_500_000,
											25_000_000,
											PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
											PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
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
											&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order,
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
				{
					// Place another short term order from Bob to match the remaining size of the long-term order
					blockHeight: 4,
					orderMsgs: []clobtypes.MsgPlaceOrder{
						PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20,
					},
					// Short term order placement results in Create and Update with fully-filled amount for both
					// taker and maker
					expectedOffchainMessagesCheckTx: []msgsender.Message{
						off_chain_updates.MustCreateOrderPlaceMessage(
							nil,
							PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order,
						).AddHeader(msgsender.MessageHeader{
							Key:   msgsender.TransactionHashHeaderKey,
							Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id1_Sell1_Price50000_GTB20.Tx),
						}),
						off_chain_updates.MustCreateOrderUpdateMessage(
							nil,
							LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.OrderId,
							LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.GetBaseQuantums(),
						).AddHeader(msgsender.MessageHeader{
							Key:   msgsender.TransactionHashHeaderKey,
							Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id1_Sell1_Price50000_GTB20.Tx),
						}),
						off_chain_updates.MustCreateOrderUpdateMessage(
							nil,
							PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
							PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
						).AddHeader(msgsender.MessageHeader{
							Key:   msgsender.TransactionHashHeaderKey,
							Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id1_Sell1_Price50000_GTB20.Tx),
						}),
					},
					expectedOffchainMessagesAfterBlock: []msgsender.Message{},
					expectedOnchainMessagesAfterBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
						&indexer_manager.IndexerTendermintBlock{
							Height: 4,
							Time:   ctx.BlockTime(),
							Events: []*indexer_manager.IndexerTendermintEvent{
								// taker subaccount state transition
								{
									Subtype: indexerevents.SubtypeSubaccountUpdate,
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewSubaccountUpdateEvent(
											&constants.Bob_Num0,
											[]*satypes.PerpetualPosition{
												{
													PerpetualId: Clob_0.MustGetPerpetualId(),
													// perpetual position size should equal sum of base quantums of Bob's orders
													// because they are both fully filled
													Quantums: dtypes.NewInt(-int64(
														PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums() +
															PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums(),
													)),
													FundingIndex: dtypes.NewInt(0),
												},
											},
											[]*satypes.AssetPosition{
												{
													AssetId: 0,
													Quantums: dtypes.NewIntFromBigInt(
														new(big.Int).Add(
															bobSubaccount.GetUsdcPosition(),
															new(big.Int).SetInt64(
																50_000_000_000+5_500_000+ // maker rebate of .110% from first order
																	50_000_000_000-25_000_000, // taker fee of .5% from second order
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
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewSubaccountUpdateEvent(
											&constants.Alice_Num0,
											[]*satypes.PerpetualPosition{
												{
													PerpetualId: Clob_0.MustGetPerpetualId(),
													// Order was fully filled
													Quantums: dtypes.NewInt(int64(
														LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.GetQuantums())),
													FundingIndex: dtypes.NewInt(0),
												},
											},
											[]*satypes.AssetPosition{
												{
													AssetId: 0,
													Quantums: dtypes.NewIntFromBigInt(
														new(big.Int).Sub(
															aliceSubaccount.GetUsdcPosition(),
															new(big.Int).SetInt64(
																50_000_000_000+25_000_000+ // taker fee of .5% from first match
																	50_000_000_000-5_500_000, // maker rebate of .110% from second match
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
									DataBytes: indexer_manager.GetBytes(
										indexerevents.NewOrderFillEvent(
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order,
											PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order,
											PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
											-5_500_000,
											25_000_000,
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.GetBaseQuantums(),
											PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
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
												ShortTermOrderPlacement: CheckTx_PlaceOrder_Bob_Num0_Id1_Sell1_Price50000_GTB20.Tx,
											},
										},
										clobtestutils.NewMatchOperationRaw(
											&PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order,
											[]clobtypes.MakerFill{
												{
													FillAmount: PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.
														Order.GetBaseQuantums().ToUint64(),
													MakerOrderId: LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.OrderId,
												},
											},
										),
									},
								}))),
							},
						},
					)},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()
			ctx := tApp.InitChain()

			for _, ordersAndExpectations := range tc.ordersAndExpectationsPerBlock {
				// CheckTx
				for _, order := range ordersAndExpectations.orderMsgs {
					for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
						ctx,
						tApp.App,
						order,
					) {
						resp := tApp.CheckTx(checkTx)
						require.True(
							t,
							resp.IsOK(),
							"Expected CheckTx to succeed. Response: %+v, Block Height: %d",
							resp,
							ordersAndExpectations.blockHeight,
						)
					}
				}
				require.ElementsMatch(
					t,
					ordersAndExpectations.expectedOffchainMessagesCheckTx,
					msgSender.GetOffchainMessages(),
					"Block height: %d",
					ordersAndExpectations.blockHeight,
				)
				msgSender.Clear()

				// Block Processing
				ctx = tApp.AdvanceToBlock(ordersAndExpectations.blockHeight, testapp.AdvanceToBlockOptions{})
				require.ElementsMatch(
					t,
					ordersAndExpectations.expectedOnchainMessagesAfterBlock,
					msgSender.GetOnchainMessages(),
					"Block height: %d",
					ordersAndExpectations.blockHeight,
				)
				require.ElementsMatch(
					t,
					ordersAndExpectations.expectedOffchainMessagesAfterBlock,
					msgSender.GetOffchainMessages(),
					"Block height: %d",
					ordersAndExpectations.blockHeight,
				)
				msgSender.Clear()
			}

			// Verify orderbook
			_, found := tApp.App.ClobKeeper.MemClob.GetOrder(ctx, tc.order.OrderId)
			require.Equal(t, tc.orderShouldRestOnOrderbook, found)

			// Verify fill amount
			_, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, tc.order.OrderId)
			require.Equal(
				t,
				tc.expectedOrderFillAmount,
				fillAmount.ToUint64(),
				"Fill amount should be %d, not %d",
				tc.expectedOrderFillAmount,
				fillAmount,
			)

			// Verify subaccounts
			for _, expectedSubaccount := range tc.expectedSubaccounts {
				subaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id)
				require.Equal(t, expectedSubaccount, subaccount)
			}
		})
	}
}
