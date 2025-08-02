package clob_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"

	"github.com/cometbft/cometbft/crypto/tmhash"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
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
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
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

// TestCancelFullyFilledStatefulOrderInSameBlockItIsFilled tests the scenario where an honest block proposer
// may propose a stateful cancellation which fails because the order was fully filled in the same block.
func TestCancelFullyFilledStatefulOrderInSameBlockItIsFilled(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Place order
	result := tApp.CheckTx(testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
	))
	require.True(t, result.IsOK(), "Expected CheckTx to succeed. Response: %+v", result)
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Place order which fully matches the first order
	result = tApp.CheckTx(testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Bob_Num0.Owner,
		},
		&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
	))
	require.True(t, result.IsOK(), "Expected CheckTx to succeed. Response: %+v", result)

	// Place cancellation
	cancellationTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
		},
		&constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
	)
	result = tApp.CheckTx(cancellationTx)
	require.True(t, result.IsOK(), "Expected CheckTx to succeed. Response: %+v", result)

	// DeliverTx should fail for cancellation tx
	ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		ValidateFinalizeBlock: func(
			ctx sdktypes.Context,
			request abcitypes.RequestFinalizeBlock,
			response abcitypes.ResponseFinalizeBlock,
		) (haltChain bool) {
			for txIndex, execResult := range response.TxResults {
				if txIndex == 1 {
					require.True(t, execResult.IsErr())
					require.Equal(t, clobtypes.ErrStatefulOrderCancellationFailedForAlreadyRemovedOrder.ABCICode(), execResult.Code)
				} else {
					require.True(t, execResult.IsOK(), "Expected DeliverTx to succeed. Response log: %+v", execResult.Log)
				}
			}

			return false
		},
	})

	_, exists := tApp.App.ClobKeeper.GetLongTermOrderPlacement(
		ctx,
		LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.Order.OrderId,
	)
	require.False(t, exists)
	exists, _, _ = tApp.App.ClobKeeper.GetOrderFillAmount(
		ctx,
		LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.Order.OrderId,
	)
	require.False(t, exists)
}

func TestCancelStatefulOrder(t *testing.T) {
	type checkResults struct {
		orderId       clobtypes.OrderId
		existsInState bool
	}

	tests := map[string]struct {
		blockWithMessages []testmsgs.TestBlockWithMsgs
		expectations      checkResults
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

			expectations: checkResults{
				orderId:       LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.Order.OrderId,
				existsInState: false,
			},
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

			expectations: checkResults{
				orderId:       LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.Order.OrderId,
				existsInState: false,
			},
		},
		"Test stateful order is cancelled when placed and then partially matched and cancelled in next block": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
							ExpectedIsOk: true,
						},
					},
				},
				{
					Block: 3,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &PlaceOrder_Bob_Num0_Id0_Clob0_Sell4_Price10_GTB20,
							ExpectedIsOk: true,
						},
						{
							Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
							ExpectedIsOk: true,
						},
					},
				},
			},
			expectations: checkResults{
				orderId:       LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.Order.OrderId,
				existsInState: false,
			},
		},
		"Test stateful order is placed when placed, cancelled, then re-placed with the same order id": {
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

			expectations: checkResults{
				orderId:       LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5.Order.OrderId,
				existsInState: true,
			},
		},
		"Test stateful order cancel for non existent order fails": {
			blockWithMessages: []testmsgs.TestBlockWithMsgs{
				{
					Block: 2,
					Msgs: []testmsgs.TestSdkMsg{
						{
							Msg:          &constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15,
							ExpectedIsOk: false,
						},
					},
				},
			},

			expectations: checkResults{
				orderId:       constants.CancelLongTermOrder_Alice_Num0_Id0_Clob0_GTBT15.OrderId,
				existsInState: false,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()

			for _, blockWithMessages := range tc.blockWithMessages {
				for _, testSdkMsg := range blockWithMessages.Msgs {
					result := tApp.CheckTx(testapp.MustMakeCheckTx(
						ctx,
						tApp.App,
						testapp.MustMakeCheckTxOptions{
							AccAddressForSigning: testtx.MustGetOnlySignerAddress(tApp.App.AppCodec(), testSdkMsg.Msg),
						},
						testSdkMsg.Msg,
					))
					require.Conditionf(t, func() bool {
						return result.IsOK() == testSdkMsg.ExpectedIsOk
					}, "Expected CheckTx to succeed. Response: %+v", result)
				}

				ctx = tApp.AdvanceToBlock(blockWithMessages.Block, testapp.AdvanceToBlockOptions{})
			}

			_, exists := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, tc.expectations.orderId)
			require.Equal(t, exists, tc.expectations.existsInState)
			exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(
				ctx,
				tc.expectations.orderId,
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
			clobtypes.Order{
				OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 1},
				Side:         clobtypes.Order_SIDE_BUY,
				Quantums:     20,
				Subticks:     15,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
				TimeInForce:  clobtypes.Order_TIME_IN_FORCE_FILL_OR_KILL,
			},
		),
	) {
		resp := tApp.CheckTx(checkTx)
		require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
		require.Contains(t, resp.Log, clobtypes.ErrDeprecatedField.Error())
	}

	// Reject long-term IOC/FOK in DeliverTx
	blockAdvancement := testapp.BlockAdvancementWithErrors{
		BlockAdvancement: testapp.BlockAdvancement{
			StatefulOrders: []clobtypes.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_IOC,
				{
					OrderId: clobtypes.OrderId{
						SubaccountId: constants.Carl_Num0,
						ClientId:     0,
						OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
						ClobPairId:   0,
					},
					Side:         clobtypes.Order_SIDE_SELL,
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
					TimeInForce:  clobtypes.Order_TIME_IN_FORCE_FILL_OR_KILL,
				},
			},
		},
		ExpectedDeliverTxErrors: map[int]string{
			0: clobtypes.ErrLongTermOrdersCannotRequireImmediateExecution.Error(),
			1: clobtypes.ErrDeprecatedField.Error(),
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
			testapp.MustScaleOrder(order, testapp.DefaultGenesis()),
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
			AccAddressForSigning: constants.Alice_Num0.Owner,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
	)
	CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Bob_Num0.Owner,
		},
		&PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id1_Sell1_Price50000_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Bob_Num0.Owner,
		},
		&PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20,
	)
	CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy10_Price49999_GTBT15_PO := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: constants.Alice_Num0.Owner,
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
							ctx,
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
						testutil.CreateSinglePerpetualPosition(
							Clob_0.MustGetPerpetualId(),
							big.NewInt(int64(
								LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							new(big.Int).Sub(
								aliceSubaccount.GetUsdcPosition(),
								new(big.Int).SetInt64(
									50_000_000_000+25_000_000, // taker fee of .5%
								),
							),
						),
					},
					MarginEnabled: true,
				},
				{
					Id: &constants.Bob_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							Clob_0.MustGetPerpetualId(),
							big.NewInt(-int64(
								LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							new(big.Int).Add(
								bobSubaccount.GetUsdcPosition(),
								new(big.Int).SetInt64(
									50_000_000_000+5_500_000, // maker rebate of .110%
								),
							),
						),
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
							ctx,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
						).AddHeader(msgsender.MessageHeader{
							Key:   msgsender.TransactionHashHeaderKey,
							Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20.Tx),
						}),
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
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
							ctx,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
						),
						// taker
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
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
												testutil.CreateSinglePerpetualPosition(
													Clob_0.MustGetPerpetualId(),
													big.NewInt(int64(
														LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
													big.NewInt(0),
													big.NewInt(0),
												),
											},
											[]*satypes.AssetPosition{
												testutil.CreateSingleAssetPosition(
													0,
													new(big.Int).Sub(
														aliceSubaccount.GetUsdcPosition(),
														new(big.Int).SetInt64(
															50_000_000_000+25_000_000, // taker fee of .5%
														),
													),
												),
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
												testutil.CreateSinglePerpetualPosition(
													Clob_0.MustGetPerpetualId(),
													big.NewInt(-int64(
														LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
													big.NewInt(0),
													big.NewInt(0),
												),
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
											0,
											0,
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
											big.NewInt(0),
											0,
											0,
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
							ctx,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
						).AddHeader(msgsender.MessageHeader{
							Key:   msgsender.TransactionHashHeaderKey,
							Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20.Tx),
						}),
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
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
							ctx,
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
						testutil.CreateSinglePerpetualPosition(
							Clob_0.MustGetPerpetualId(),
							big.NewInt(int64(
								LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.GetQuantums())),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							new(big.Int).Sub(
								aliceSubaccount.GetUsdcPosition(),
								new(big.Int).SetInt64(
									50_000_000_000+25_000_000+ // taker fee of .5%
										50_000_000_000-5_500_000, // maker rebate of .110%
								),
							),
						),
					},
					MarginEnabled: true,
				},
				{
					Id: &constants.Bob_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							Clob_0.MustGetPerpetualId(),
							big.NewInt(-int64(
								PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums()+
									PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums(),
							)),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							new(big.Int).Add(
								bobSubaccount.GetUsdcPosition(),
								new(big.Int).SetInt64(
									50_000_000_000+5_500_000+ // maker rebate of .110% from first order
										50_000_000_000-25_000_000, // taker fee of .5% from second order
								),
							),
						),
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
							ctx,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
						).AddHeader(msgsender.MessageHeader{
							Key:   msgsender.TransactionHashHeaderKey,
							Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell1_Price50000_GTB20.Tx),
						}),
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
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
							ctx,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
							PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
						),
						// taker, partially filled
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
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
							ctx,
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
												testutil.CreateSinglePerpetualPosition(
													Clob_0.MustGetPerpetualId(),
													big.NewInt(int64(
														PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums())),
													big.NewInt(0),
													big.NewInt(0),
												),
											},
											[]*satypes.AssetPosition{
												testutil.CreateSingleAssetPosition(
													0,
													new(big.Int).Sub(
														aliceSubaccount.GetUsdcPosition(),
														new(big.Int).SetInt64(
															50_000_000_000+25_000_000, // taker fee of .5%
														),
													),
												),
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
												testutil.CreateSinglePerpetualPosition(
													Clob_0.MustGetPerpetualId(),
													big.NewInt(-int64(
														PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums())),
													big.NewInt(0),
													big.NewInt(0),
												),
											},
											[]*satypes.AssetPosition{
												testutil.CreateSingleAssetPosition(
													0,
													new(big.Int).Add(
														bobSubaccount.GetUsdcPosition(),
														new(big.Int).SetInt64(
															50_000_000_000+5_500_000, // maker rebate of .110%
														),
													),
												),
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
											0,
											0,
											PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
											PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
											big.NewInt(0),
											0,
											0,
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
							ctx,
							PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order,
						).AddHeader(msgsender.MessageHeader{
							Key:   msgsender.TransactionHashHeaderKey,
							Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id1_Sell1_Price50000_GTB20.Tx),
						}),
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
							LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.OrderId,
							LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.GetBaseQuantums(),
						).AddHeader(msgsender.MessageHeader{
							Key:   msgsender.TransactionHashHeaderKey,
							Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id1_Sell1_Price50000_GTB20.Tx),
						}),
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
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
												testutil.CreateSinglePerpetualPosition(
													Clob_0.MustGetPerpetualId(),
													// perpetual position size should equal sum of base quantums of Bob's orders
													// because they are both fully filled
													big.NewInt(-int64(
														PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums()+
															PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.GetQuantums(),
													)),
													big.NewInt(0),
													big.NewInt(0),
												),
											},
											[]*satypes.AssetPosition{
												testutil.CreateSingleAssetPosition(
													0,
													new(big.Int).Add(
														bobSubaccount.GetUsdcPosition(),
														new(big.Int).SetInt64(
															50_000_000_000+5_500_000+ // maker rebate of .110% from first order
																50_000_000_000-25_000_000, // taker fee of .5% from second order
														),
													),
												),
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
												testutil.CreateSinglePerpetualPosition(
													Clob_0.MustGetPerpetualId(),
													// Order was fully filled
													big.NewInt(int64(
														LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.GetQuantums())),
													big.NewInt(0),
													big.NewInt(0),
												),
											},
											[]*satypes.AssetPosition{
												testutil.CreateSingleAssetPosition(
													0,
													new(big.Int).Sub(
														aliceSubaccount.GetUsdcPosition(),
														new(big.Int).SetInt64(
															50_000_000_000+25_000_000+ // taker fee of .5% from first match
																50_000_000_000-5_500_000, // maker rebate of .110% from second match
														),
													),
												),
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
											0,
											0,
											LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy2_Price50000_GTBT5.Order.GetBaseQuantums(),
											PlaceOrder_Bob_Num0_Id1_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
											big.NewInt(0),
											0,
											0,
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

				messages := msgSender.GetOnchainMessages()
				fmt.Println("Onchain messages", messages)
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
			_, found := tApp.App.ClobKeeper.MemClob.GetOrder(tc.order.OrderId)
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

func TestRegression_InvalidTimeInForce(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).
		// Disable non-determinism checks since we mutate keeper state directly.
		WithNonDeterminismChecksEnabled(false).
		Build()
	ctx := tApp.InitChain()

	// subaccounts for indexer expectation assertions
	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)

	// order msgs
	Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5 := *clobtypes.NewMsgPlaceOrder(
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
	// CheckTx Txs needed for indexer expectation assertions
	CheckTx_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetOnlySignerAddress(
				tApp.App.AppCodec(),
				&Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
			),
		},
		&Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
	)

	// Pre-existing order with invalid time in force.
	LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Bob_Num0,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   0,
			},
			Side:     clobtypes.Order_SIDE_SELL,
			Quantums: 10_000_000_000,
			Subticks: 500_000_000,
			// Invalid time in force
			TimeInForce:  clobtypes.Order_TimeInForce(uint32(999)),
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		},
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
		"Test matching an order fully as taker against order with invalid time in force": {
			order:                      Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
			orderShouldRestOnOrderbook: false,
			expectedOrderFillAmount:    0, // order is fully-filled and removed from state
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							Clob_0.MustGetPerpetualId(),
							big.NewInt(int64(
								Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							new(big.Int).Sub(
								aliceSubaccount.GetUsdcPosition(),
								new(big.Int).SetInt64(
									50_000_000_000+25_000_000, // taker fee of .5%
								),
							),
						),
					},
					MarginEnabled: true,
				},
				{
					Id: &constants.Bob_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							Clob_0.MustGetPerpetualId(),
							big.NewInt(-int64(
								Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
							big.NewInt(0),
							big.NewInt(0),
						),
					},
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							new(big.Int).Add(
								bobSubaccount.GetUsdcPosition(),
								new(big.Int).SetInt64(
									50_000_000_000+5_500_000, // maker rebate of .110%
								),
							),
						),
					},
					MarginEnabled: true,
				},
			},
			ordersAndExpectationsPerBlock: []ordersAndExpectations{
				{
					blockHeight: 2,
					orderMsgs: []clobtypes.MsgPlaceOrder{
						Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
					},
					// Short term order placement results in Create and Update with 0 fill amount
					expectedOffchainMessagesCheckTx: []msgsender.Message{},
					// Short term order update for fill amount, stateful order update for fill amount
					// Note there are no headers because these events are generated in PrepareCheckState
					expectedOffchainMessagesAfterBlock: []msgsender.Message{
						// maker
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
							LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
							LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.GetBaseQuantums(),
						),
						// taker
						off_chain_updates.MustCreateOrderUpdateMessage(
							ctx,
							Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.OrderId,
							Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
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
											Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
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
												testutil.CreateSinglePerpetualPosition(
													Clob_0.MustGetPerpetualId(),
													big.NewInt(int64(
														Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
													big.NewInt(0),
													big.NewInt(0),
												),
											},
											[]*satypes.AssetPosition{
												testutil.CreateSingleAssetPosition(
													0,
													new(big.Int).Sub(
														aliceSubaccount.GetUsdcPosition(),
														new(big.Int).SetInt64(
															50_000_000_000+25_000_000, // taker fee of .5%
														),
													),
												),
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
												testutil.CreateSinglePerpetualPosition(
													Clob_0.MustGetPerpetualId(),
													big.NewInt(-int64(
														Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetQuantums())),
													big.NewInt(0),
													big.NewInt(0),
												),
											},
											[]*satypes.AssetPosition{
												testutil.CreateSingleAssetPosition(
													0,
													new(big.Int).Add(
														bobSubaccount.GetUsdcPosition(),
														new(big.Int).SetInt64(
															50_000_000_000+5_500_000, // maker rebate of .110%
														),
													),
												),
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
											LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
											Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
											Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
											-5_500_000,
											25_000_000,
											0,
											0,
											Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
											Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.GetBaseQuantums(),
											big.NewInt(0),
											0,
											0,
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
										clobtestutils.NewMatchOperationRaw(
											&Invalid_TIF_LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order,
											[]clobtypes.MakerFill{
												{
													FillAmount: LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.
														Order.GetBaseQuantums().ToUint64(),
													MakerOrderId: LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
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
			tApp := testapp.NewTestAppBuilder(t).
				// Disable non-determinism checks since we mutate keeper state directly.
				WithNonDeterminismChecksEnabled(false).
				WithAppOptions(appOpts).Build()
			ctx := tApp.InitChain()

			// Add the order with invalid time in force to state and orderbook.
			tApp.App.ClobKeeper.SetLongTermOrderPlacement(
				tApp.App.NewUncachedContext(true, tmproto.Header{}),
				LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
				1,
			)
			tApp.App.ClobKeeper.AddStatefulOrderIdExpiration(
				tApp.App.NewUncachedContext(true, tmproto.Header{}),
				time.Unix(5, 0),
				LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order.OrderId,
			)
			_, _, _, err := tApp.App.ClobKeeper.MemClob.PlaceOrder(
				tApp.App.NewUncachedContext(true, tmproto.Header{}),
				LongTermPlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price50000_GTB20.Order,
			)
			require.NoError(t, err)

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
			_, found := tApp.App.ClobKeeper.MemClob.GetOrder(tc.order.OrderId)
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

func TestMultiplePlaceOrdersInSingleTransaction(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	// Get subaccounts for indexer expectation assertions
	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	aliceSubaccount1 := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num1)

	// Create multiple orders for Alice
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *aliceSubaccount.Id,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     10_000_000_000, // 1 BTC
			Subticks:     500_000_000,    // 50k USDC / BTC
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		},
	)

	LongTermPlaceOrder_Alice_Num0_Id1_Clob0_Buy1_Price50000_GTBT5 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *aliceSubaccount.Id,
				ClientId:     1,
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     10_000_000_000, // 1 BTC
			Subticks:     500_000_000,    // 50k USDC / BTC
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		},
	)

	// Long term orders for Alice's subaccount 1
	LongTermPlaceOrder_Alice_Num1_Id0_Clob0_Buy1_Price50000_GTBT5 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *aliceSubaccount1.Id,
				ClientId:     0,
				OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     10_000_000_000, // 1 BTC
			Subticks:     500_000_000,    // 50k USDC / BTC
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
		},
	)

	// Create a short-term order for Alice
	ShortTermPlaceOrder_Alice_Num0_Id2_Clob0_Buy1_Price50000_GTB20 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: *aliceSubaccount.Id,
				ClientId:     2,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     10_000_000_000, // 1 BTC
			Subticks:     500_000_000,    // 50k USDC / BTC
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
	)

	// Create a cancel order for Alice's first long-term order
	CancelOrder_Alice_Num0_Id0_Clob0 := *clobtypes.NewMsgCancelOrderStateful(
		LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.OrderId,
		5, // Cancel at block 5
	)

	// Create a transfer message from Alice's subaccount 0 to Alice's subaccount 1
	Transfer_Alice_Num0_To_Alice_Num1 := *sendingtypes.NewMsgCreateTransfer(
		&sendingtypes.Transfer{
			Sender:    constants.Alice_Num0,
			Recipient: constants.Alice_Num1,
			AssetId:   assettypes.AssetUsdc.Id,
			Amount:    100_000_000_000, // 100k USDC
		},
	)

	// Create test cases
	tests := map[string]struct {
		msgs                     []sdktypes.Msg
		expectedOrdersInMemclob  map[clobtypes.OrderId]bool
		expectedOrderFillAmounts map[clobtypes.OrderId]uint64
		expectedSubaccounts      []satypes.Subaccount
		expectedErr              string
	}{
		"Successfully places multiple long-term orders in single transaction": {
			msgs: []sdktypes.Msg{
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
				&LongTermPlaceOrder_Alice_Num0_Id1_Clob0_Buy1_Price50000_GTBT5,
			},
			expectedOrdersInMemclob: map[clobtypes.OrderId]bool{
				LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.OrderId: true,
				LongTermPlaceOrder_Alice_Num0_Id1_Clob0_Buy1_Price50000_GTBT5.Order.OrderId: true,
			},
			expectedSubaccounts: []satypes.Subaccount{aliceSubaccount},
		},
		"Fails when mixing long-term and short-term orders": {
			msgs: []sdktypes.Msg{
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
				&ShortTermPlaceOrder_Alice_Num0_Id2_Clob0_Buy1_Price50000_GTB20,
			},
			expectedErr: "a transaction containing short term MsgCancelOrder or MsgPlaceOrder " +
				"may not contain more than one message",
		},
		"Successfully places and cancels order in single transaction": {
			msgs: []sdktypes.Msg{
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5,
				&CancelOrder_Alice_Num0_Id0_Clob0,
			},
			expectedOrdersInMemclob: map[clobtypes.OrderId]bool{
				LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy1_Price50000_GTBT5.Order.OrderId: false,
			},
			expectedSubaccounts: []satypes.Subaccount{aliceSubaccount},
		},
		"Successfully transfers and places order in single transaction": {
			msgs: []sdktypes.Msg{
				&Transfer_Alice_Num0_To_Alice_Num1,
				&LongTermPlaceOrder_Alice_Num1_Id0_Clob0_Buy1_Price50000_GTBT5,
			},
			expectedOrdersInMemclob: map[clobtypes.OrderId]bool{
				LongTermPlaceOrder_Alice_Num1_Id0_Clob0_Buy1_Price50000_GTBT5.Order.OrderId: true,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							new(big.Int).Add(
								aliceSubaccount1.GetUsdcPosition(),
								new(big.Int).SetInt64(100_000_000_000), // 100k USDC transfer
							),
						),
					},
					MarginEnabled: false,
				},
			},
		},
		"Fails when mixing transfer and short-term orders": {
			msgs: []sdktypes.Msg{
				&Transfer_Alice_Num0_To_Alice_Num1,
				&ShortTermPlaceOrder_Alice_Num0_Id2_Clob0_Buy1_Price50000_GTB20,
			},
			expectedErr: "a transaction containing short term MsgCancelOrder or MsgPlaceOrder " +
				"may not contain more than one message",
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
				appOpts := map[string]interface{}{
					indexer.MsgSenderInstanceForTest: msgSender,
				}
				tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()
				ctx := tApp.InitChain()

				// Create a transaction with multiple messages
				tx := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: constants.AliceAccAddress.String(),
						Gas:                  100_000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					tc.msgs...,
				)

				// Check that CheckTx succeeds
				checkTxResp := tApp.CheckTx(tx)
				if tc.expectedErr != "" {
					require.Conditionf(t, checkTxResp.IsErr, "Expected CheckTx to error. Response: %+v", checkTxResp)
					require.Contains(t, checkTxResp.Log, tc.expectedErr)
					return
				}
				require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

				// Advance to block 3
				ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

				// Verify orders in memclob
				for orderId, shouldExist := range tc.expectedOrdersInMemclob {
					_, exists := tApp.App.ClobKeeper.GetLongTermOrderPlacement(ctx, orderId)
					require.Equal(t, shouldExist, exists)
				}

				// Verify subaccounts
				for _, expectedSubaccount := range tc.expectedSubaccounts {
					subaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id)
					require.Equal(t, expectedSubaccount, subaccount)
				}
			},
		)
	}
}
