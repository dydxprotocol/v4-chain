package clob_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/rand"
	"gopkg.in/typ.v4/slices"

	"github.com/cometbft/cometbft/types"

	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	epochtypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	Clob_0                                            = MustGetClobPairsFromGenesis(testapp.DefaultGenesis())[0]
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 5},
		},
		testapp.DefaultGenesis(),
	))
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
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB27 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 27},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num1, ClientId: 0, ClobPairId: 0},
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
	CancelOrder_Alice_Num0_Id0_Clob0_GTB27 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			ClobPairId:   0,
		},
		27,
	)
	CancelOrder_Alice_Num1_Id0_Clob0_GTB20 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num1,
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
	LongTermPlaceOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTBT5 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.LongTermOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTBT5,
		testapp.DefaultGenesis(),
	))
	ConditionalPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
		testapp.DefaultGenesis(),
	))
	ConditionalPlaceOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15 = *clobtypes.NewMsgPlaceOrder(MustScaleOrder(
		constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
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
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(msg),
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
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(placeOrderMsg),
				},
				privKeySupplier,
				placeOrderMsg,
			)
			cancelOrderMsg := clobtypes.NewMsgCancelOrderShortTerm(orderId, 20)
			checkTxsPerAccount[i][1] = testapp.MustMakeCheckTxWithPrivKeySupplier(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: testtx.MustGetOnlySignerAddress(cancelOrderMsg),
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
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
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

func TestFailsDeliverTxWithIncorrectlySignedPlaceOrderTx(t *testing.T) {
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

func TestFailsDeliverTxWithUnsignedTransactions(t *testing.T) {
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
			resp := tApp.CheckTx(checkTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
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
			resp := tApp.CheckTx(checkTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
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
			resp := tApp.CheckTx(checkTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
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
	msgPlaceOrder.Order.Quantums = msgPlaceOrder.Order.Quantums * clobPair.StepBaseQuantums
	msgPlaceOrder.Order.Subticks = msgPlaceOrder.Order.Subticks * uint64(clobPair.SubticksPerTick)
	msgPlaceOrder.Order.ConditionalOrderTriggerSubticks = msgPlaceOrder.Order.ConditionalOrderTriggerSubticks *
		uint64(clobPair.SubticksPerTick)

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
