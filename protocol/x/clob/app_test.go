package clob_test

import (
	"bytes"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4/app/config"
	"github.com/dydxprotocol/v4/daemons/liquidation/api"
	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/rand"
	"gopkg.in/typ.v4/slices"

	"github.com/cometbft/cometbft/types"

	"github.com/dydxprotocol/v4/indexer"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"github.com/dydxprotocol/v4/indexer/msgsender"
	"github.com/dydxprotocol/v4/indexer/off_chain_updates"
	testapp "github.com/dydxprotocol/v4/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4/testutil/clob"
	"github.com/dydxprotocol/v4/testutil/constants"
	testmsgs "github.com/dydxprotocol/v4/testutil/msgs"
	testtx "github.com/dydxprotocol/v4/testutil/tx"
	assettypes "github.com/dydxprotocol/v4/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	epochtypes "github.com/dydxprotocol/v4/x/epochs/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4/x/sending/types"
	stattypes "github.com/dydxprotocol/v4/x/stats/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	Clob_0                                             = MustGetClobPairsFromGenesis(testapp.DefaultGenesis())[0]
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     6,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 1},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	CancelOrder_Alice_Num0_Id0_Clob0_GTB5 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			ClobPairId:   0,
		},
		5,
	)
	CancelOrder_Alice_Num0_Id0_Clob1_GTB5 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			ClobPairId:   1,
		},
		5,
	)
	CancelOrder_Alice_Num0_Id0_Clob0_GTB20 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			ClobPairId:   0,
		},
		20,
	)
	PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	CancelOrder_Bob_Num0_Id0_Clob0_GTB20 = *clobtypes.NewMsgCancelOrderShortTerm(
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
	LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
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
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
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

	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

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
	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
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

	CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20),
		},
		&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
	)
	CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(&PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20),
		},
		&PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
	)
	CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20 := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20),
		},
		&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
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
			orders: []clobtypes.MsgPlaceOrder{PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
					Time:   ctx.BlockTime(),
				})},
		},
		"Test matching an order fully": {
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			},
			expectedOrdersFilled: []clobtypes.OrderId{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
			},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
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
												PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
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
												PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
						},
					},
					TxHashes: []string{string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
						OperationsQueue: []clobtypes.OperationRaw{
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Alice_Num0_Id0_Buy5_Price10_GTB20.Tx,
								},
							},
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx,
								},
							},
							clobtestutils.NewMatchOperationRaw(
								&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
								[]clobtypes.MakerFill{
									{
										FillAmount: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.
											Order.GetBaseQuantums().ToUint64(),
										MakerOrderId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order.OrderId,
									},
								},
							),
						},
					},
					)))},
				})},
		},
		"Test matching an order partially, maker order remains on book": {
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			},
			expectedOrdersFilled: []clobtypes.OrderId{
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
			},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
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
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
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
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
						},
					},
					TxHashes: []string{string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
						OperationsQueue: []clobtypes.OperationRaw{
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx,
								},
							},
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx,
								},
							},
							clobtestutils.NewMatchOperationRaw(
								&PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
								[]clobtypes.MakerFill{
									{
										FillAmount: PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.
											Order.GetBaseQuantums().ToUint64(),
										MakerOrderId: PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
									},
								},
							),
						},
					},
					)))},
				})},
		},
		"Test matching an order partially, taker order remains on book": {
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
			},
			expectedOrdersFilled: []clobtypes.OrderId{
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
			},
			expectedOffchainMessagesAfterPlaceOrder: []msgsender.Message{
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					0,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderPlaceMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				).AddHeader(msgsender.MessageHeader{
					Key:   msgsender.TransactionHashHeaderKey,
					Value: tmhash.Sum(CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx),
				}),
			},
			expectedOffchainMessagesInNextBlock: []msgsender.Message{
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
				off_chain_updates.MustCreateOrderUpdateMessage(
					nil,
					PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order.OrderId,
					PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
				),
			},
			expectedOnchainMessagesInNextBlock: []msgsender.Message{indexer_manager.CreateIndexerBlockEventMessage(
				&indexer_manager.IndexerTendermintBlock{
					Height: 2,
					Time:   ctx.BlockTime(),
					Events: []*indexer_manager.IndexerTendermintEvent{
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Alice_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          0,
						},
						{
							Subtype: indexerevents.SubtypeSubaccountUpdate,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewSubaccountUpdateEvent(
									&constants.Bob_Num0,
									[]*satypes.PerpetualPosition{
										{
											PerpetualId: Clob_0.MustGetPerpetualId(),
											Quantums: dtypes.NewInt(-int64(
												PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetQuantums())),
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
									nil, // no funding payments
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          1,
						},
						{
							Subtype: indexerevents.SubtypeOrderFill,
							Data: indexer_manager.GetB64EncodedEventMessage(
								indexerevents.NewOrderFillEvent(
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order,
									PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
									PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.GetBaseQuantums(),
									0, // Fees are 0 due to lost precision
									0,
								),
							),
							OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{},
							EventIndex:          2,
						},
					},
					TxHashes: []string{string(lib.GetTxHash(testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
						OperationsQueue: []clobtypes.OperationRaw{
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Bob_Num0_Id0_Sell5_Price10_GTB20.Tx,
								},
							},
							{
								Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
									ShortTermOrderPlacement: CheckTx_PlaceOrder_Alice_Num0_Id0_Buy6_Price10_GTB20.Tx,
								},
							},
							clobtestutils.NewMatchOperationRaw(
								&PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20.Order,
								[]clobtypes.MakerFill{
									{
										FillAmount: PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.
											Order.GetBaseQuantums().ToUint64(),
										MakerOrderId: PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20.Order.OrderId,
									},
								},
							),
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

			ctx = tApp.InitChain()
			// Clear any messages produced prior to these checkTx calls.
			msgSender.Clear()
			for _, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					checkTxResp := tApp.CheckTx(checkTx)
					require.True(t, checkTxResp.IsOK(), fmt.Sprintf("CheckTx failed: %v\n", checkTxResp))
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
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
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

func TestPlacePerpetualLiquidation_Deleveraging(t *testing.T) {
	tests := map[string]struct {
		// State.
		subaccounts []satypes.Subaccount

		// Parameters.
		placedMatchableOrders     []clobtypes.MatchableOrder
		liquidatableSubaccountIds []satypes.SubaccountId
		liquidationConfig         clobtypes.LiquidationsConfig

		// Expectations.
		expectedSubaccounts []satypes.Subaccount
	}{
		`Can place a liquidation order that is fully filled and does not require deleveraging`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10, // Order at $50,000
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_499_000_000 - 50_000_000_000 - 250_000_000),
						},
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(100_000_000_000), // $100,000
						},
					},
				},
			},
		},
		`Can place a liquidation order that is partially filled and does not require deleveraging`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// First order at $50,000
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				// Second order at $60,000, which does not cross the liquidation order
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price60000_GTB10,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_499_000_000 - 12_500_000_000 - 62_500_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(-75_000_000), // -0.75 BTC
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 12_500_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(75_000_000), // 0.75 BTC
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
			},
		},
		`Can place a liquidation order that is unfilled and full position size is deleveraged`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $50,499, and closing at $50,500
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 50_499_000_000),
						},
					},
				},
			},
		},
		`Can place a liquidation order that is partially-filled and remaining size is deleveraged`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// First order at $50,498, Carl pays $0.25 to the insurance fund.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11,
				// Carl's bankruptcy price to close 0.75 BTC short is $50,499, and closing at $50,500
				// would require $0.75 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 50_499_000_000 - 250_000),
						},
					},
				},
			},
		},
		`Can place a liquidation order that is unfilled and cannot be deleveraged due to
		non-overlapping bankruptcy prices`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing at $50,000
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails.
				// Dave's bankruptcy price to close 1 BTC long is $50,000, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
			},
		},
		`Can place a liquidation order that is partially-filled and cannot be deleveraged due to
		non-overlapping bankruptcy prices`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				constants.Dave_Num1_025BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Dave_Num1_Id0_Clob0_Sell025BTC_Price49999_GTB10,
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing 0.75 BTC at $50,000
				// would require $0.75 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails for remaining amount.
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(49_999_000_000 - 12_499_750_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(-75_000_000), // -0.75 BTC
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				// Dave's bankruptcy price to close 1 BTC long is $50,000, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				// Dave_Num0 does not change since deleveraging against this subaccount failed.
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				{
					Id: &constants.Dave_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 12_499_750_000),
						},
					},
				},
			},
		},
		`Can place a liquidation order that is unfilled, then only a portion of the remaining size can
		deleveraged due to non-overlapping bankruptcy prices with some subaccounts`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				constants.Dave_Num1_05BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing 0.75 BTC at $50,000
				// would require $0.75 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(49_999_000_000 - 24_999_500_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							// Deleveraging fails for remaining amount.
							Quantums:     dtypes.NewInt(-50_000_000), // -0.5 BTC
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				// Dave_Num0 does not change since deleveraging against this subaccount failed.
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				{
					Id: &constants.Dave_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 24_999_500_000),
						},
					},
				},
			},
		},
		`Can place a liquidation order that is partially-filled, then deleveraged for only a
		portion of the remaining size due to non-overlapping bankruptcy prices with some subaccounts`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				constants.Dave_Num1_05BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Dave_Num1_Id0_Clob0_Sell025BTC_Price49999_GTB10,
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing 0.75 BTC at $50,000
				// would require $0.75 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails for remaining amount.
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId: 0,
							// 0.25 BTC closed by liquidation order, 0.25 BTC closed by deleveraging.
							Quantums: dtypes.NewInt(49_999_000_000 - 12_499_750_000 - 12_499_750_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(-50_000_000), // -0.5 BTC
							FundingIndex: dtypes.NewInt(0),
						},
					},
				},
				// Dave_Num0 does not change since deleveraging against this subaccount failed.
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				{
					Id: &constants.Dave_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId: 0,
							// 0.25 BTC closed by liquidation order, 0.25 BTC closed by deleveraging.
							Quantums: dtypes.NewInt(50_000_000_000 + 12_499_750_000 + 12_499_750_000),
						},
					},
				},
			},
		},
		`Deleveraging takes precedence - can place a liquidation order that would fail due to exceeding 
		subaccount limit and full position size is deleveraged`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $50,499, and closing at $50,500
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Max_Smmr,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: clobtypes.SubaccountBlockLimits{
					MaxNotionalLiquidated:    math.MaxUint64,
					MaxQuantumsInsuranceLost: 1,
				},
			},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 50_499_000_000),
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *assettypes.GenesisState) {
						genesisState.Assets = []assettypes.Asset{
							*constants.Usdc,
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
							constants.EthUsd_20PercentInitial_10PercentMaintenance,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.ClobPairs = []clobtypes.ClobPair{
							constants.ClobPair_Btc_No_Fee,
							constants.ClobPair_Eth_No_Fee,
						}
						genesisState.LiquidationsConfig = tc.liquidationConfig
					},
				)
				return genesis
			}).WithTesting(t).Build()

			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Create all existing orders.
			existingOrderMsgs := make([]clobtypes.MsgPlaceOrder, len(tc.placedMatchableOrders))
			for i, matchableOrder := range tc.placedMatchableOrders {
				existingOrderMsgs[i] = clobtypes.MsgPlaceOrder{Order: matchableOrder.MustGetOrder()}
			}
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, existingOrderMsgs...) {
				require.True(t, tApp.CheckTx(checkTx).IsOK())
			}

			_, err := tApp.App.Server.LiquidateSubaccounts(ctx, &api.LiquidateSubaccountsRequest{
				SubaccountIds: tc.liquidatableSubaccountIds,
			})
			require.NoError(t, err)

			// Verify test expectations.
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			for _, expectedSubaccount := range tc.expectedSubaccounts {
				require.Equal(
					t,
					expectedSubaccount,
					tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id),
				)
			}
		})
	}
}

func TestProcessProposalFailsDeliverTxWithIncorrectlySignedPlaceOrderTx(t *testing.T) {
	tests := map[string]struct {
		accAddressForSigning string
		msg                  sdktypes.Msg
	}{
		// these orders are from Alice, but are instead signed by Bob
		"Signed order placement with incorrect signer": {
			accAddressForSigning: constants.BobAccAddress.String(),
			msg:                  constants.Msg_PlaceOrder,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					require.ErrorContains(
						t,
						r.(error),
						"invalid pubkey: MsgProposedOperations is invalid",
					)
				} else {
					t.Error("Expected panic")
				}
			}()
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tAppBuilder := testapp.NewTestAppBuilder().WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts))
			tApp := tAppBuilder.Build()
			tApp.InitChain()
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			signedTransaction := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{AccAddressForSigning: tc.accAddressForSigning},
				tc.msg,
			).Tx

			operationsQueue := make([]clobtypes.OperationRaw, 0)
			switch tc.msg.(type) {
			case *clobtypes.MsgPlaceOrder:
				operationsQueue = append(
					operationsQueue,
					clobtypes.OperationRaw{
						Operation: &clobtypes.OperationRaw_ShortTermOrderPlacement{
							ShortTermOrderPlacement: signedTransaction,
						},
					},
				)
			default:
				require.Fail(t, "Invalid operation type: %+v", tc.msg)
			}

			proposal := tApp.PrepareProposal()
			proposal.Txs[0] = testtx.MustGetTxBytes(
				&clobtypes.MsgProposedOperations{
					OperationsQueue: operationsQueue,
				},
			)

			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{RequestProcessProposalTxsOverride: proposal.Txs})
		})
	}
}

func TestProcessProposalFailsDeliverTxWithUnsignedTransactions(t *testing.T) {
	tests := map[string]struct {
		proposedOperationsTx []byte
	}{
		"Unsigned order placement": {
			proposedOperationsTx: testtx.MustGetTxBytes(&clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{
					clobtestutils.NewShortTermOrderPlacementOperationRaw(
						PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20.Order,
					),
				},
			}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					require.ErrorContains(
						t,
						r.(error),
						"Error: no signatures supplied: MsgProposedOperations is invalid",
					)
				} else {
					t.Error("Expected panic")
				}
			}()

			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tAppBuilder := testapp.NewTestAppBuilder().WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts))
			tApp := tAppBuilder.Build()
			tApp.InitChain()
			tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			proposal := tApp.PrepareProposal()
			proposal.Txs[0] = tc.proposedOperationsTx

			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{RequestProcessProposalTxsOverride: proposal.Txs})
		})
	}
}

func TestRateLimitingOrders_RateLimitsAreEnforced(t *testing.T) {
	tests := map[string]struct {
		blockRateLimitConifg clobtypes.BlockRateLimitConfiguration
		firstMsg             sdktypes.Msg
		secondMsg            sdktypes.Msg
	}{
		"Short term orders": {
			blockRateLimitConifg: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersPerMarketPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			secondMsg: &PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
		},
		"Stateful orders": {
			blockRateLimitConifg: clobtypes.BlockRateLimitConfiguration{
				MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
			secondMsg: &LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
		},
		"Short term order cancellations": {
			blockRateLimitConifg: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrderCancellationsPerMarketPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMsg:  &CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
			secondMsg: &CancelOrder_Alice_Num0_Id0_Clob0_GTB20,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.BlockRateLimitConfig = tc.blockRateLimitConifg
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			firstCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(tc.firstMsg),
				},
				tc.firstMsg,
			)
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			// First order should be allowed.
			require.True(t, tApp.CheckTx(firstCheckTx).IsOK())

			secondCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(tc.secondMsg),
				},
				tc.secondMsg,
			)
			// Rate limit is 1 over two block, second attempt should be blocked.
			resp := tApp.CheckTx(secondCheckTx)
			require.True(t, resp.IsErr())
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "Rate of 2 exceeds configured block rate limit")

			// Rate limit of 1 over two blocks should still apply, total should be 3 now (2 in block 2, 1 in block 3).
			tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.True(t, resp.IsErr())
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "Rate of 3 exceeds configured block rate limit")

			// Rate limit of 1 over two blocks should still apply, total should be 2 now (1 in block 3, 1 in block 4).
			tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.True(t, resp.IsErr())
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "Rate of 2 exceeds configured block rate limit")

			// Advancing two blocks should make the total count 0 now and the msg should be accepted.
			tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{})
			resp = tApp.CheckTx(secondCheckTx)
			require.True(t, resp.IsOK())
		})
	}
}

func TestRateLimitingOrders_ShortTermOrderRateLimitsArePerMarket(t *testing.T) {
	tests := map[string]struct {
		blockRateLimitConifg clobtypes.BlockRateLimitConfiguration
		firstMarketMsg       sdktypes.Msg
		secondMarketMsg      sdktypes.Msg
		firstMarketSecondMsg sdktypes.Msg
	}{
		"Short term orders": {
			blockRateLimitConifg: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrdersPerMarketPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMarketMsg:       &PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
			secondMarketMsg:      &PlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB20,
			firstMarketSecondMsg: &PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20,
		},
		"Short term order cancellations": {
			blockRateLimitConifg: clobtypes.BlockRateLimitConfiguration{
				MaxShortTermOrderCancellationsPerMarketPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
					{
						NumBlocks: 2,
						Limit:     1,
					},
				},
			},
			firstMarketMsg:       &CancelOrder_Alice_Num0_Id0_Clob0_GTB5,
			secondMarketMsg:      &CancelOrder_Alice_Num0_Id0_Clob1_GTB5,
			firstMarketSecondMsg: &CancelOrder_Alice_Num0_Id0_Clob0_GTB20,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						genesisState.BlockRateLimitConfig = tc.blockRateLimitConifg
					},
				)
				return genesis
			}).WithTesting(t).Build()
			ctx := tApp.InitChain()

			firstMarketCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(tc.firstMarketMsg),
				},
				tc.firstMarketMsg,
			)
			secondMarketCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(tc.secondMarketMsg),
				},
				tc.secondMarketMsg,
			)
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			// First order for each market should be allowed.
			require.True(t, tApp.CheckTx(firstMarketCheckTx).IsOK())
			require.True(t, tApp.CheckTx(secondMarketCheckTx).IsOK())

			firstMarketSecondCheckTx := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetSignerAddress(tc.firstMarketSecondMsg),
				},
				tc.firstMarketSecondMsg,
			)
			// Rate limit is 1 over two block, second attempt should be blocked.
			resp := tApp.CheckTx(firstMarketSecondCheckTx)
			require.True(t, resp.IsErr())
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
			require.Contains(t, resp.Log, "Rate of 2 exceeds configured block rate limit")
		})
	}
}

func TestCancellationAndMatchInTheSameBlock_Regression(t *testing.T) {
	msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
	appOpts := map[string]interface{}{
		indexer.MsgSenderInstanceForTest: msgSender,
	}
	tAppBuilder := testapp.NewTestAppBuilder().WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts))
	tApp := tAppBuilder.Build()

	LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20 := *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0, ClobPairId: 0,
				OrderFlags: clobtypes.OrderIdFlags_LongTerm,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 20},
		},
		testapp.DefaultGenesis(),
	))
	LCancelOrder_Alice_Num0_Id0_Clob0_GTBT20 := *clobtypes.NewMsgCancelOrderStateful(
		LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.Order.OrderId,
		20,
	)

	PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price10_GTB20 := *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     1,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Bob_Num0_Id0_Clob0_Sell7_Price10_GTB20 := *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     7,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))

	tApp.InitChain()
	ctx := tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
	for _, msg := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		LPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
	) {
		resp := tApp.CheckTx(msg)
		require.True(t, resp.IsOK())
	}
	ctx = tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
	for _, msg := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		PlaceOrder_Bob_Num0_Id0_Clob0_Sell1_Price10_GTB20,
	) {
		resp := tApp.CheckTx(msg)
		require.True(t, resp.IsOK())
	}
	ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{})
	for _, msg := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		LCancelOrder_Alice_Num0_Id0_Clob0_GTBT20,
	) {
		resp := tApp.CheckTx(msg)
		require.True(t, resp.IsOK())
	}
	for _, msg := range testapp.MustMakeCheckTxsWithClobMsg(
		ctx,
		tApp.App,
		PlaceOrder_Bob_Num0_Id0_Clob0_Sell7_Price10_GTB20,
	) {
		resp := tApp.CheckTx(msg)
		require.True(t, resp.IsOK())
	}
	// We shouldn't be overfilling orders and the line below shouldn't panic.
	_ = tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{
		ValidateDeliverTxs: func(_ sdktypes.Context, _ abcitypes.RequestDeliverTx, _ abcitypes.ResponseDeliverTx) bool {
			// Don't halt the chain since it's expected that the order will be removed after getting fully filled,
			// so the subsequent cancellation will be invalid.
			return false
		},
	})
}

func TestRateLimitingOrders_StatefulOrderRateLimitsAreAcrossMarkets(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *clobtypes.GenesisState) {
				genesisState.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{
					MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 2,
							Limit:     1,
						},
					},
				}
			},
		)
		return genesis
	}).WithTesting(t).Build()
	ctx := tApp.InitChain()

	firstMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5),
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
	)

	// Second order should not be allowed in 2nd block and allowed in 4th block.
	secondMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5),
			AccSequenceNumberForSigning: 2,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
	)

	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	// First order should be allowed and second should be rejected.
	require.True(t, tApp.CheckTx(firstMarketCheckTx).IsOK())
	resp := tApp.CheckTx(secondMarketCheckTx)
	require.True(t, resp.IsErr())
	require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), resp.Code)
	require.Contains(t, resp.Log, "Rate of 2 exceeds configured block rate limit")

	// Retrying in the 4th block should succeed since the rate limits should have been pruned.
	tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
	require.True(t, tApp.CheckTx(secondMarketCheckTx).IsOK())
}

func TestRateLimitingOrders_StatefulOrdersDuringDeliverTxAreRateLimited(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *clobtypes.GenesisState) {
				genesisState.BlockRateLimitConfig = clobtypes.BlockRateLimitConfiguration{
					MaxStatefulOrdersPerNBlocks: []clobtypes.MaxPerNBlocksRateLimit{
						{
							NumBlocks: 2,
							Limit:     1,
						},
					},
				}
			},
		)
		return genesis
	}).WithTesting(t).Build()
	ctx := tApp.InitChain()

	tApp.CheckTx(testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5),
			AccSequenceNumberForSigning: 2,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
	))

	firstMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5),
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
	)
	secondMarketCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{
			AccAddressForSigning: testtx.MustGetSignerAddress(
				&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5),
			AccSequenceNumberForSigning: 2,
		},
		&LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
	)

	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{firstMarketCheckTx.Tx, secondMarketCheckTx.Tx},
		ValidateDeliverTxs: func(
			context sdktypes.Context,
			request abcitypes.RequestDeliverTx,
			response abcitypes.ResponseDeliverTx,
		) (haltChain bool) {
			if bytes.Equal(request.Tx, firstMarketCheckTx.Tx) {
				require.True(t, response.IsOK())
			} else {
				require.True(t, response.IsErr())
				require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), response.Code)
				require.Contains(t, response.Log, "Rate of 2 exceeds configured block rate limit")
			}
			return false
		},
	})

	// Advance to block 3 which should cause the delivered stateful order to still be rejected from block 2.
	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{secondMarketCheckTx.Tx},
		ValidateDeliverTxs: func(
			ctx sdktypes.Context,
			request abcitypes.RequestDeliverTx,
			response abcitypes.ResponseDeliverTx,
		) (haltchain bool) {
			require.True(t, response.IsErr())
			require.Equal(t, clobtypes.ErrBlockRateLimitExceeded.ABCICode(), response.Code)
			require.Contains(t, response.Log, "Rate of 3 exceeds configured block rate limit")
			return false
		},
	})

	// Advance to block 4 should clear out the delivered transactions in 2 and 3 allowing them to be
	// delivered in block 5.
	tApp.AdvanceToBlock(4, testapp.AdvanceToBlockOptions{})
	tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
		DeliverTxsOverride: [][]byte{secondMarketCheckTx.Tx},
	})
}

func TestOrderRemoval(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		msgs        []sdktypes.Msg
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
		// TODO(CLOB-735): Investigate whether we still need to skip add-to-book collat check for pre-existing
		// stateful orders.
		// "under-collateralized taker when adding to book is removed": {
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_10000USD,
		// 		constants.Dave_Num0_10000USD,
		// 	},
		// 	firstOrder:  constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTBT10,
		// 	// Does not cross with best bid.
		// 	secondOrder: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,

		// 	withdrawal: &sendingtypes.MsgWithdrawFromSubaccount{
		// 		Sender:    constants.Dave_Num0,
		// 		Recipient: constants.DaveAccAddress.String(),
		// 		AssetId:   constants.Usdc.Id,
		// 		Quantums:  10_000_000_000,
		// 	},

		// 	expectedFirstOrderRemoved:  false,
		// 	expectedSecondOrderRemoved: true, // taker order fails add-to-orderbook collateralization check
		// },
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

func TestStats(t *testing.T) {
	msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
	appOpts := map[string]interface{}{
		indexer.MsgSenderInstanceForTest: msgSender,
	}
	tAppBuilder := testapp.NewTestAppBuilder().WithAppCreatorFn(testapp.DefaultTestAppCreatorFn(appOpts))
	tApp := tAppBuilder.Build()

	// Epochs start at block height 2.
	startTime := time.Unix(10, 0).UTC()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: startTime,
	})

	aliceAddress := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0).Id.MustGetAccAddress().String()
	bobAddress := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0).Id.MustGetAccAddress().String()

	createAliceBuyOrder := func(clientId uint32) *clobtypes.MsgPlaceOrder {
		return clobtypes.NewMsgPlaceOrder(MustScaleOrder(
			clobtypes.Order{
				OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: clientId, ClobPairId: 0},
				Side:         clobtypes.Order_SIDE_BUY,
				Quantums:     5000,
				Subticks:     1000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
			},
			testapp.DefaultGenesis(),
		))
	}
	createBobSellOrder := func(clientId uint32) *clobtypes.MsgPlaceOrder {
		return clobtypes.NewMsgPlaceOrder(MustScaleOrder(
			clobtypes.Order{
				OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: clientId, ClobPairId: 0},
				Side:         clobtypes.Order_SIDE_SELL,
				Quantums:     5000,
				Subticks:     1000,
				GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
			},
			testapp.DefaultGenesis(),
		))
	}

	// Check that UserStats and GlobalStats reflect the orders filled
	requireStatsEqual := func(expectedNotional uint64) {
		require.Equal(t, &stattypes.UserStats{
			TakerNotional: 0,
			MakerNotional: expectedNotional,
		}, tApp.App.StatsKeeper.GetUserStats(ctx, aliceAddress))
		require.Equal(t, &stattypes.UserStats{
			TakerNotional: expectedNotional,
			MakerNotional: 0,
		}, tApp.App.StatsKeeper.GetUserStats(ctx, bobAddress))
		require.Equal(t, &stattypes.GlobalStats{
			NotionalTraded: expectedNotional,
		}, tApp.App.StatsKeeper.GetGlobalStats(ctx))
	}

	// Check that the correct epoch stats exist
	requireEpochStatsEqual := func(epoch uint32, expectedNotional uint64) {
		require.Equal(t, &stattypes.EpochStats{
			EpochEndTime: time.Unix(0, 0).
				Add((time.Duration((epoch + 1) * epochtypes.StatsEpochDuration)) * time.Second).
				UTC(),
			// Alice's address happens to be after Bob's alphabetically
			Stats: []*stattypes.EpochStats_UserWithStats{
				{
					User: bobAddress,
					Stats: &stattypes.UserStats{
						TakerNotional: expectedNotional,
						MakerNotional: 0,
					},
				},
				{
					User: aliceAddress,
					Stats: &stattypes.UserStats{
						TakerNotional: 0,
						MakerNotional: expectedNotional,
					},
				},
			},
		}, tApp.App.StatsKeeper.GetEpochStatsOrNil(ctx, epoch))
	}

	// Multiple orders per block
	orderMsgs := []*clobtypes.MsgPlaceOrder{
		createAliceBuyOrder(0),
		createBobSellOrder(0),
		createAliceBuyOrder(1),
		createBobSellOrder(1),
	}
	for _, order := range orderMsgs {
		for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *order) {
			require.True(t, tApp.CheckTx(checkTx).IsOK())
		}
	}
	currTime := startTime
	ctx = tApp.AdvanceToBlock(5, testapp.AdvanceToBlockOptions{
		BlockTime: currTime,
	})
	requireStatsEqual(10000)
	requireEpochStatsEqual(0, 10000)

	orderMsgs = []*clobtypes.MsgPlaceOrder{
		createAliceBuyOrder(2),
		createBobSellOrder(2),
	}
	for _, order := range orderMsgs {
		for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *order) {
			require.True(t, tApp.CheckTx(checkTx).IsOK())
		}
	}
	// Don't advance the epoch, so these stats are on the same epoch as the previous block
	currTime = time.Unix(0, 0).Add((time.Duration(epochtypes.StatsEpochDuration) - 1) * time.Second)
	ctx = tApp.AdvanceToBlock(6, testapp.AdvanceToBlockOptions{
		BlockTime: currTime,
	})
	requireStatsEqual(15000)
	requireEpochStatsEqual(0, 15000)

	// Advance epoch without adding stats
	currTime = currTime.Add(time.Duration(epochtypes.StatsEpochDuration) * time.Second)
	ctx = tApp.AdvanceToBlock(7, testapp.AdvanceToBlockOptions{
		BlockTime: currTime,
	})
	requireStatsEqual(15000)

	orderMsgs = []*clobtypes.MsgPlaceOrder{
		createAliceBuyOrder(3),
		createBobSellOrder(3),
	}
	for _, order := range orderMsgs {
		for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, *order) {
			require.True(t, tApp.CheckTx(checkTx).IsOK())
		}
	}
	currTime = currTime.Add(time.Duration(epochtypes.StatsEpochDuration) * time.Second)
	ctx = tApp.AdvanceToBlock(8, testapp.AdvanceToBlockOptions{
		BlockTime: currTime,
	})
	requireStatsEqual(20000)
	requireEpochStatsEqual(2, 5000)

	// Advance the window one epoch at a time and check to make sure stats fall out of the window
	windowDuration := tApp.App.StatsKeeper.GetWindowDuration(ctx)
	currTime = time.Unix(0, 0).Add(time.Duration(windowDuration)).Add(time.Second).UTC()
	ctx = tApp.AdvanceToBlock(10, testapp.AdvanceToBlockOptions{
		BlockTime: currTime,
	})
	requireStatsEqual(20000)

	currTime = currTime.Add(time.Duration(epochtypes.StatsEpochDuration) * time.Second)
	ctx = tApp.AdvanceToBlock(11, testapp.AdvanceToBlockOptions{
		BlockTime: currTime,
	})
	requireStatsEqual(5000)

	// Recall that we made an epoch without any fills
	currTime = currTime.Add(time.Duration(epochtypes.StatsEpochDuration) * time.Second)
	ctx = tApp.AdvanceToBlock(12, testapp.AdvanceToBlockOptions{
		BlockTime: currTime,
	})
	requireStatsEqual(5000)

	currTime = currTime.Add(time.Duration(epochtypes.StatsEpochDuration) * time.Second)
	ctx = tApp.AdvanceToBlock(13, testapp.AdvanceToBlockOptions{
		BlockTime: currTime,
	})
	requireStatsEqual(0)
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
