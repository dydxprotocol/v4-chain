package clob_test

import (
	"fmt"
	"sync"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4/app/config"
	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/testutil/rand"
	"gopkg.in/typ.v4/slices"

	"github.com/cometbft/cometbft/types"

	"github.com/dydxprotocol/v4/indexer"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/indexer/msgsender"
	"github.com/dydxprotocol/v4/indexer/off_chain_updates"
	"github.com/dydxprotocol/v4/lib"
	testapp "github.com/dydxprotocol/v4/testutil/app"
	"github.com/dydxprotocol/v4/testutil/constants"
	testmsgs "github.com/dydxprotocol/v4/testutil/msgs"
	testtx "github.com/dydxprotocol/v4/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	Clob_0                                       = MustGetClobPairsFromGenesis(testapp.DefaultGenesis())[0]
	PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	CancelOrder_Alice_Num0_Id0_GTB20 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			ClobPairId:   0,
		},
		20,
	)
	PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	CancelOrder_Bob_Num0_Id0_GTB20 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			ClobPairId:   0,
		},
		20,
	)

	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
		testapp.DefaultGenesis(),
	))
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
		testapp.DefaultGenesis(),
	))
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
		testapp.DefaultGenesis(),
	))
	ConditionalPlaceOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
		testapp.DefaultGenesis(),
	))
	ConditionalPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
		testapp.DefaultGenesis(),
	))
)

// We place 300 orders that match and 700 orders followed by their cancellations concurrently.
//
// This test heavily relies on golangs race detector to validate memory reads and writes are properly ordered.
func TestConcurrentMatchesAndCancels(t *testing.T) {
	r := rand.NewRand()
	simAccounts := simtypes.RandomAccounts(r, 1000)
	tApp := testapp.NewTestAppBuilder().WithTesting(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *auth.GenesisState) {
				for _, simAccount := range simAccounts {
					acct := &auth.BaseAccount{
						Address: sdktypes.AccAddress(simAccount.PubKey.Address()).String(),
						PubKey:  codectypes.UnsafePackAny(simAccount.PubKey),
					}
					genesisState.Accounts = append(genesisState.Accounts, codectypes.UnsafePackAny(acct))
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				for _, simAccount := range simAccounts {
					genesisState.Subaccounts = append(genesisState.Subaccounts, satypes.Subaccount{
						Id: &satypes.SubaccountId{
							Owner:  sdktypes.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, simAccount.PubKey.Address()),
							Number: 0,
						},
						AssetPositions: []*satypes.AssetPosition{
							&constants.Usdc_Asset_500_000,
						},
					})
				}
			},
		)
		return genesis
	}).Build()

	ctx := tApp.AdvanceToBlock(2)

	expectedFills := make([]clobtypes.Order, 300)
	expectedCancels := make([]clobtypes.Order, len(simAccounts)-len(expectedFills))
	checkTxsPerAccount := make([][]abcitypes.RequestCheckTx, len(simAccounts))
	for i, simAccount := range simAccounts {
		privKeySupplier := func(accAddress string) cryptotypes.PrivKey {
			expectedAccAddress := sdktypes.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, simAccount.PubKey.Address())
			if accAddress != expectedAccAddress {
				panic(fmt.Errorf("Unknown account, got %s, expected %s", accAddress, expectedAccAddress))
			}
			return simAccount.PrivKey
		}
		orderId := clobtypes.OrderId{
			SubaccountId: satypes.SubaccountId{
				Owner:  sdktypes.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, simAccount.PubKey.Address()),
				Number: 0,
			},
			ClientId:   0,
			ClobPairId: 0,
		}

		if i < len(expectedFills) {
			// 300 orders, 150 buys and 150 sells where there are 50 each of size 5, 10, and 15 accounting for a total
			// matched volume of 250 + 500 + 750 = 1500 quantums. We specifically use 5, 10 and 15 to ensure that we get
			// orders that are partial matches, full matches, and matches that cross multiple orders.
			checkTxsPerAccount[i] = make([]abcitypes.RequestCheckTx, 1)
			var side clobtypes.Order_Side
			var quantums uint64
			// We use 6 here since we want 3 different sizes (5/10/15) * 2 different sides (buy/sell)
			switch i % 6 {
			case 0:
				side = clobtypes.Order_SIDE_BUY
				quantums = 5
			case 1:
				side = clobtypes.Order_SIDE_BUY
				quantums = 10
			case 2:
				side = clobtypes.Order_SIDE_BUY
				quantums = 15
			case 3:
				side = clobtypes.Order_SIDE_SELL
				quantums = 5
			case 4:
				side = clobtypes.Order_SIDE_SELL
				quantums = 10
			case 5:
				side = clobtypes.Order_SIDE_SELL
				quantums = 15
			default:
				panic("Unimplemented case?")
			}
			expectedFills[i] = MustScaleOrder(clobtypes.Order{
				OrderId:      orderId,
				Side:         side,
				Quantums:     quantums,
				Subticks:     20,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
			}, Clob_0)
			msg := clobtypes.NewMsgPlaceOrder(expectedFills[i])
			checkTxsPerAccount[i][0] = testapp.MustMakeCheckTxWithPrivKeySupplier(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(msg),
				},
				privKeySupplier,
				msg,
			)
		} else {
			// The remainder are cancels for orders that would never match.
			checkTxsPerAccount[i] = make([]abcitypes.RequestCheckTx, 2)
			idx := i - len(expectedFills)

			// We use 2 here since we want orders that we will cancel on both sides (buy/sell)
			switch i % 2 {
			case 0:
				expectedCancels[idx] = MustScaleOrder(clobtypes.Order{
					OrderId:      orderId,
					Side:         clobtypes.Order_SIDE_BUY,
					Quantums:     1,
					Subticks:     10,
					GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
				},
					Clob_0)
			case 1:
				expectedCancels[idx] = MustScaleOrder(clobtypes.Order{
					OrderId:      orderId,
					Side:         clobtypes.Order_SIDE_SELL,
					Quantums:     1,
					Subticks:     30,
					GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
				},
					Clob_0)
			default:
				panic("Unimplemented case?")
			}
			placeOrderMsg := clobtypes.NewMsgPlaceOrder(expectedCancels[idx])
			checkTxsPerAccount[i][0] = testapp.MustMakeCheckTxWithPrivKeySupplier(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(placeOrderMsg),
				},
				privKeySupplier,
				placeOrderMsg,
			)
			cancelOrderMsg := clobtypes.NewMsgCancelOrderShortTerm(orderId, 20)
			checkTxsPerAccount[i][1] = testapp.MustMakeCheckTxWithPrivKeySupplier(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(cancelOrderMsg),
				},
				privKeySupplier,
				cancelOrderMsg,
			)
		}
	}

	// Shuffle the ordering of CheckTx calls to increase the randomness of the order of execution. Note
	// that the wait group and concurrent goroutine execution adds randomness as well because it will be
	// dependent on which goroutine wakeup order.
	slices.Shuffle(checkTxsPerAccount)

	var wgStart, wgFinish sync.WaitGroup
	wgStart.Add(len(checkTxsPerAccount))
	wgFinish.Add(len(checkTxsPerAccount))
	for i := 0; i < len(checkTxsPerAccount); i++ {
		checkTxs := checkTxsPerAccount[i]
		go func() {
			// Ensure that we unlock the wait group regardless of how this goroutine completes.
			defer wgFinish.Done()

			// Mark that we have started and then wait till everyone else starts to increase the amount of contention
			// and parallelization.
			wgStart.Done()
			wgStart.Wait()
			for _, checkTx := range checkTxs {
				resp := tApp.CheckTx(checkTx)
				require.True(t, resp.IsOK())
			}
		}()
	}

	// Wait till all the orders were placed and cancelled.
	wgFinish.Wait()

	// Advance the block and ensure that the appropriate orders were filled and cancelled.
	tApp.AdvanceToBlock(3)
	for _, expectedFill := range expectedFills {
		exists, amount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, expectedFill.OrderId)
		require.True(t, exists)
		require.Equal(t, expectedFill.Quantums, amount.ToUint64())
	}
	for _, expectedCancel := range expectedCancels {
		exists, amount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, expectedCancel.OrderId)
		require.False(t, exists)
		require.Equal(t, uint64(0), amount.ToUint64())
	}
}

func TestPlaceOrder(t *testing.T) {
	msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
	appOpts := map[string]interface{}{
		indexer.MsgSenderInstanceForTest: msgSender,
	}
	tAppBuilder := testapp.NewTestAppBuilder().WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts))
	tApp := tAppBuilder.Build()
	ctx := tApp.InitChain()

	aliceSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0)
	bobSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0)

	placeOrderAlice := &PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20
	CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(placeOrderAlice),
		},
		placeOrderAlice,
	)
	placeOrderBob := &PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20
	CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(placeOrderBob),
		},
		placeOrderBob,
	)

	tests := map[string]struct {
		orders                                  []clobtypes.MsgPlaceOrder
		expectedOrdersFilled                    []clobtypes.OrderId
		expectedOffchainMessagesAfterPlaceOrder []msgsender.Message
		expectedOnchainMessagesAfterPlaceOrder  []msgsender.Message
		expectedOffchainMessagesInNextBlock     []msgsender.Message
		expectedOnchainMessagesInNextBlock      []msgsender.Message
	}{
		"Test placing an order": {
			orders: []clobtypes.MsgPlaceOrder{PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 3,
					Time:   ctx.BlockTime(),
				})},
		},
		"Test matching an order fully": {
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20,
				PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20,
			},
			expectedOrdersFilled: []clobtypes.OrderId{
				PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
				PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.OrderId,
			},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 3,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(-int64(
												PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetQuantums())),
											FundingIndex: dtypes.NewInt(0),
										},
									},
									// Maker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										{
											AssetId:  lib.UsdcAssetId,
											Quantums: dtypes.NewIntFromBigInt(bobSubaccount.GetUsdcPosition()),
										},
									},
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
						},
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(int64(
												PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetQuantums())),
											FundingIndex: dtypes.NewInt(0),
										},
									},
									// Taker fees calculate to 0 so asset position doesn't change.
									[]*satypes.AssetPosition{
										{
											AssetId:  lib.UsdcAssetId,
											Quantums: dtypes.NewIntFromBigInt(aliceSubaccount.GetUsdcPosition()),
										},
									},
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
						},
					},
					TxHashes: []string{string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
						OperationsQueue: []clobtypes.Operation{
							clobtypes.NewOrderPlacementOperation(
								PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order),
							clobtypes.NewOrderPlacementOperation(
								PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order),
							clobtypes.NewMatchOperation(
								&PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order,
								[]clobtypes.MakerFill{
									{
										FillAmount: PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.
											Order.GetBaseQuantums().ToUint64(),
										MakerOrderId: PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.OrderId,
									},
								},
							),
						},
						AddToOrderbookCollatCheckOrderHashes: [][]byte{
							PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Order.GetOrderHash().ToBytes(),
							PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Order.GetOrderHash().ToBytes(),
						},
					},
					)))},
				})},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Reset for each iteration of the loop
			tApp.Reset()

			ctx = tApp.AdvanceToBlock(2)
			// Clear any messages produced prior to these checkTx calls.
			msgSender.Clear()
			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					require.True(t, tApp.CheckTx(checkTx).IsOK())
				}
			}
			require.ElementsMatch(
				t,
				tc.expectedOffchainMessagesAfterPlaceOrder,
				msgSender.GetOffchainMessages(),
			)
			require.ElementsMatch(t, tc.expectedOnchainMessagesAfterPlaceOrder, msgSender.GetOnchainMessages())

			// Clear the messages that we already matched prior to advancing to the next block.
			msgSender.Clear()
			ctx = tApp.AdvanceToBlock(3)
			require.ElementsMatch(t, tc.expectedOffchainMessagesInNextBlock, msgSender.GetOffchainMessages())
			require.ElementsMatch(t, tc.expectedOnchainMessagesInNextBlock, msgSender.GetOnchainMessages())
			for _, order := range tc.orders {
				if lib.ContainsValue(tc.expectedOrdersFilled, order.Order.OrderId) {
					exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(
						ctx,
						order.Order.OrderId,
					)

					require.True(t, exists)
					require.Equal(t, order.Order.GetBaseQuantums(), fillAmount)
				}
			}
		})
	}
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
							Msg:          &PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20,
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
				ctx := tApp.AdvanceToBlock(blockWithMessages.Block)

				for _, testSdkMsg := range blockWithMessages.Msgs {
					result := tApp.CheckTx(testapp.MustMakeCheckTx(
						ctx,
						tApp.App,
						testapp.MustMakeCheckTxOptions{
							AccAddressForSigning: testtx.MustGetSignerAddress(testSdkMsg.Msg),
						},
						testSdkMsg.Msg,
					)).IsOK()
					require.Equal(
						t,
						testSdkMsg.ExpectedIsOk,
						result,
					)
				}
			}

			ctx := tApp.AdvanceToBlockIfNecessary(tc.checkResultsBlock)
			require.True(t, tApp.App.ClobKeeper.HasSeenPlaceOrder(ctx, tc.checkCancelledPlaceOrder))
			exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(
				ctx,
				tc.checkCancelledPlaceOrder.Order.OrderId,
			)
			require.False(t, exists)
			require.Equal(t, uint64(0), fillAmount.ToUint64())
		})
	}
}

// MustScaleOrder scales clobtypes.Order and clobtypes.MsgPlaceorder based upon the clob information provided.
// Will panic if:
//   - OrderT is an unknown type.
//   - ClobPairSrcT is an unknown type.
//   - The clob pair id can't be used to look up the clob pair from the ClobPairSrcT.
func MustScaleOrder[
	OrderT clobtypes.Order | clobtypes.MsgPlaceOrder,
	ClobPairSrcT clobtypes.ClobPair | types.GenesisDoc](
	order OrderT,
	clobPairSrc ClobPairSrcT,
) OrderT {
	var msgPlaceOrder clobtypes.MsgPlaceOrder

	// Find the clob pair id based upon the type of order passed in.
	var clobPairId clobtypes.ClobPairId
	switch v := any(order).(type) {
	case clobtypes.MsgPlaceOrder:
		clobPairId = v.Order.GetClobPairId()
		msgPlaceOrder = v
	case clobtypes.Order:
		clobPairId = v.GetClobPairId()
		msgPlaceOrder = *clobtypes.NewMsgPlaceOrder(v)
	default:
		panic(fmt.Errorf("Unknown order type %T to get order", order))
	}

	// Find the clob pair based upon the clobPairSrc of the clob information passed in.
	var clobPair clobtypes.ClobPair
	switch v := any(clobPairSrc).(type) {
	case clobtypes.ClobPair:
		clobPair = v
	case types.GenesisDoc:
		clobPairs := MustGetClobPairsFromGenesis(v)
		if hasClobPair, ok := clobPairs[clobPairId]; ok {
			clobPair = hasClobPair
		} else {
			panic(fmt.Errorf("Clob not found in genesis doc for clob id %d", clobPairId))
		}
	default:
		panic(fmt.Errorf("Unknown source type %T to get clob pair", clobPairSrc))
	}

	// Scale the order based upon the quantums and subticks passed in.
	msgPlaceOrder.Order.Quantums = msgPlaceOrder.Order.Quantums * clobPair.MinOrderBaseQuantums
	msgPlaceOrder.Order.Subticks = msgPlaceOrder.Order.Subticks * uint64(clobPair.SubticksPerTick)

	// Return a type that matches what the user passed in for the order type.
	switch any(order).(type) {
	case clobtypes.MsgPlaceOrder:
		return any(msgPlaceOrder).(OrderT)
	case clobtypes.Order:
		return any(msgPlaceOrder.Order).(OrderT)
	default:
		panic(fmt.Errorf("Unable to convert to %T to %T", clobtypes.MsgPlaceOrder{}, order))
	}
}

// MustGetClobPairsFromGenesis unmarshals the initial genesis state and returns a map from clob pair id to clob pair.
func MustGetClobPairsFromGenesis(genesisDoc types.GenesisDoc) map[clobtypes.ClobPairId]clobtypes.ClobPair {
	var genesisState clobtypes.GenesisState
	testapp.UpdateGenesisDocWithAppStateForModule(&genesisDoc, func(genesisStatePtr *clobtypes.GenesisState) {
		genesisState = *genesisStatePtr
	})

	clobPairs := make(map[clobtypes.ClobPairId]clobtypes.ClobPair, len(genesisState.ClobPairs))
	for _, clobPair := range genesisState.ClobPairs {
		clobPairs[clobPair.GetClobPairId()] = clobPair
	}
	return clobPairs
}
