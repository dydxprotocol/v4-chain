package e2e_test

import (
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb/testutil"
	"math/rand"
	"sync/atomic"
	"testing"
)

var (
	Clob_0 = testapp.MustGetClobPairsFromGenesis(testapp.DefaultGenesis())[0]
)

func TestParallelAnteHandler_Other(t *testing.T) {
	r := testutil.NewRand()
	simAccounts := simtypes.RandomAccounts(r, 20)
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
	ctx := tApp.InitChain()

	accounts := make([]sdktypes.AccountI, len(simAccounts))
	for i, simAccount := range simAccounts {
		accounts[i] = tApp.App.AccountKeeper.GetAccount(ctx, simAccount.Address)
	}

	transferCount := atomic.Uint64{}

	for i := 0; i < len(simAccounts)/2; i++ {
		simAccount := simAccounts[i]
		account := accounts[i]

		go func() {
			for sequenceNumber := uint64(0); ; i++ {
				checkTx, err := sims.GenSignedMockTx(
					rand.New(rand.NewSource(42)),
					tApp.App.TxConfig(),
					[]sdktypes.Msg{&sendingtypes.MsgCreateTransfer{
						Transfer: &sendingtypes.Transfer{
							Sender: satypes.SubaccountId{
								Owner:  simAccount.Address.String(),
								Number: 0,
							},
							Recipient: satypes.SubaccountId{
								Owner:  simAccounts[(i+1)%len(simAccounts)].Address.String(),
								Number: 0,
							},
							AssetId: assettypes.AssetUsdc.Id,
							Amount:  constants.Usdc_Asset_500.Quantums.BigInt().Uint64(),
						},
					}},
					constants.TestFeeCoins_5Cents,
					100_000,
					ctx.ChainID(),
					[]uint64{account.GetAccountNumber()},
					[]uint64{sequenceNumber},
					simAccount.PrivKey,
				)
				if err != nil {
					panic(err)
				}
				bytes, err := tApp.App.TxConfig().TxEncoder()(checkTx)
				if err != nil {
					panic(err)
				}
				resp := tApp.CheckTx(abcitypes.RequestCheckTx{
					Tx:   bytes,
					Type: abcitypes.CheckTxType_New,
				})
				require.Conditionf(t, resp.IsOK, "Expected response to be ok: %+v", resp)
				transferCount.Add(1)
			}
		}()
	}
}

//
//func TestParallelAnteHandler_Clob(t *testing.T) {
//	r := testutilrand.NewRand()
//	simAccounts := simtypes.RandomAccounts(r, 10)
//	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
//		genesis = testapp.DefaultGenesis()
//		testapp.UpdateGenesisDocWithAppStateForModule(
//			&genesis,
//			func(genesisState *auth.GenesisState) {
//				for _, simAccount := range simAccounts {
//					acct := &auth.BaseAccount{
//						Address: sdktypes.AccAddress(simAccount.PubKey.Address()).String(),
//						PubKey:  codectypes.UnsafePackAny(simAccount.PubKey),
//					}
//					genesisState.Accounts = append(genesisState.Accounts, codectypes.UnsafePackAny(acct))
//				}
//			},
//		)
//		testapp.UpdateGenesisDocWithAppStateForModule(
//			&genesis,
//			func(genesisState *satypes.GenesisState) {
//				for _, simAccount := range simAccounts {
//					genesisState.Subaccounts = append(genesisState.Subaccounts, satypes.Subaccount{
//						Id: &satypes.SubaccountId{
//							Owner:  sdktypes.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, simAccount.PubKey.Address()),
//							Number: 0,
//						},
//						AssetPositions: []*satypes.AssetPosition{
//							&constants.Usdc_Asset_500_000,
//						},
//					})
//				}
//			},
//		)
//		return genesis
//	}).Build()
//
//	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
//
//	// In parallel:
//	//   - advance the block
//	//   - get the app/version
//	//   - load block info from store directly
//	//   - perform a custom gRPC query to get the block info
//
//	// We specifically use an atomic to ensure that we aren't providing any synchronization between the threads
//	// maximizing any data races that could exist. The wait group is only used to synchronize the testing thread
//	// when the other 4 threads are done.
//	blockLimitReached := atomic.Bool{}
//	blockLimitReached.Store(false)
//	wg := sync.WaitGroup{}
//	wg.Add(4)
//
//	go func() {
//		defer wg.Done()
//		defer func() {
//			blockLimitReached.Store(true)
//		}()
//		for i := uint32(2); i < 50; i++ {
//			tApp.AdvanceToBlock(i, testapp.AdvanceToBlockOptions{})
//		}
//	}()
//
//	for i, simAccount := range simAccounts {
//		privKeySupplier := func(accAddress string) cryptotypes.PrivKey {
//			expectedAccAddress := sdktypes.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, simAccount.PubKey.Address())
//			if accAddress != expectedAccAddress {
//				panic(fmt.Errorf("Unknown account, got %s, expected %s", accAddress, expectedAccAddress))
//			}
//			return simAccount.PrivKey
//		}
//
//		orderId := clobtypes.OrderId{
//			SubaccountId: satypes.SubaccountId{
//				Owner:  sdktypes.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, simAccount.PubKey.Address()),
//				Number: 0,
//			},
//			ClientId:   0,
//			ClobPairId: 0,
//		}
//
//		// Place matching orders
//		go func() {
//
//		}()
//
//		// Place and cancel orders
//		go func() {
//			var order clobtypes.Order
//			// We use 2 here since we want orders that we will cancel on both sides (buy/sell)
//			switch i % 2 {
//			case 0:
//				order = testapp.MustScaleOrder(clobtypes.Order{
//					OrderId:      orderId,
//					Side:         clobtypes.Order_SIDE_BUY,
//					Quantums:     1,
//					Subticks:     10,
//					GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
//				},
//					Clob_0)
//			case 1:
//				order = testapp.MustScaleOrder(clobtypes.Order{
//					OrderId:      orderId,
//					Side:         clobtypes.Order_SIDE_SELL,
//					Quantums:     1,
//					Subticks:     30,
//					GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
//				},
//					Clob_0)
//			default:
//				panic("Unimplemented case?")
//			}
//			placeOrderMsg := clobtypes.NewMsgPlaceOrder(order)
//			tApp.CheckTx(testapp.MustMakeCheckTxWithPrivKeySupplier(
//				ctx,
//				tApp.App,
//				testapp.MustMakeCheckTxOptions{
//					AccAddressForSigning: orderId.SubaccountId.Owner,
//				},
//				privKeySupplier,
//				placeOrderMsg,
//			))
//			cancelOrderMsg := clobtypes.NewMsgCancelOrderShortTerm(orderId, 20)
//			tApp.CheckTx(testapp.MustMakeCheckTxWithPrivKeySupplier(
//				ctx,
//				tApp.App,
//				testapp.MustMakeCheckTxOptions{
//					AccAddressForSigning: orderId.SubaccountId.Owner,
//				},
//				privKeySupplier,
//				cancelOrderMsg,
//			))
//		}()
//
//		//
//		//if i < len(expectedFills) {
//		//	// 300 orders, 150 buys and 150 sells where there are 50 each of size 5, 10, and 15 accounting for a total
//		//	// matched volume of 250 + 500 + 750 = 1500 quantums. We specifically use 5, 10 and 15 to ensure that we get
//		//	// orders that are partial matches, full matches, and matches that cross multiple orders.
//		//	checkTxsPerAccount[i] = make([]abcitypes.RequestCheckTx, 1)
//		//	var side clobtypes.Order_Side
//		//	var quantums uint64
//		//	// We use 6 here since we want 3 different sizes (5/10/15) * 2 different sides (buy/sell)
//		//	switch i % 6 {
//		//	case 0:
//		//		side = clobtypes.Order_SIDE_BUY
//		//		quantums = 5
//		//	case 1:
//		//		side = clobtypes.Order_SIDE_BUY
//		//		quantums = 10
//		//	case 2:
//		//		side = clobtypes.Order_SIDE_BUY
//		//		quantums = 15
//		//	case 3:
//		//		side = clobtypes.Order_SIDE_SELL
//		//		quantums = 5
//		//	case 4:
//		//		side = clobtypes.Order_SIDE_SELL
//		//		quantums = 10
//		//	case 5:
//		//		side = clobtypes.Order_SIDE_SELL
//		//		quantums = 15
//		//	default:
//		//		panic("Unimplemented case?")
//		//	}
//		//	expectedFills[i] = testapp.MustScaleOrder(clobtypes.Order{
//		//		OrderId:      orderId,
//		//		Side:         side,
//		//		Quantums:     quantums,
//		//		Subticks:     20,
//		//		GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
//		//	}, Clob_0)
//		//	msg := clobtypes.NewMsgPlaceOrder(expectedFills[i])
//		//	checkTxsPerAccount[i][0] = testapp.MustMakeCheckTxWithPrivKeySupplier(
//		//		ctx,
//		//		tApp.App,
//		//		testapp.MustMakeCheckTxOptions{
//		//			AccAddressForSigning: expectedFills[i].OrderId.SubaccountId.Owner,
//		//		},
//		//		privKeySupplier,
//		//		msg,
//		//	)
//		//} else {
//		//	// The remainder are cancels for orders that would never match.
//		//	checkTxsPerAccount[i] = make([]abcitypes.RequestCheckTx, 2)
//		//	idx := i - len(expectedFills)
//		//
//		//
//		//
//		//	}
//		//	placeOrderMsg := clobtypes.NewMsgPlaceOrder(expectedCancels[idx])
//		//	checkTxsPerAccount[i][0] = testapp.MustMakeCheckTxWithPrivKeySupplier(
//		//		ctx,
//		//		tApp.App,
//		//		testapp.MustMakeCheckTxOptions{
//		//			AccAddressForSigning: orderId.SubaccountId.Owner,
//		//		},
//		//		privKeySupplier,
//		//		placeOrderMsg,
//		//	)
//		//	cancelOrderMsg := clobtypes.NewMsgCancelOrderShortTerm(orderId, 20)
//		//	checkTxsPerAccount[i][1] = testapp.MustMakeCheckTxWithPrivKeySupplier(
//		//		ctx,
//		//		tApp.App,
//		//		testapp.MustMakeCheckTxOptions{
//		//			AccAddressForSigning: orderId.SubaccountId.Owner,
//		//		},
//		//		privKeySupplier,
//		//		cancelOrderMsg,
//		//	)
//		//}
//	}
//
//	//// Shuffle the ordering of CheckTx calls to increase the randomness of the order of execution. Note
//	//// that the wait group and concurrent goroutine execution adds randomness as well because it will be
//	//// dependent on which goroutine wakeup order.
//	//slices.Shuffle(checkTxsPerAccount)
//	//
//	//var wgStart, wgFinish sync.WaitGroup
//	//wgStart.Add(len(checkTxsPerAccount))
//	//wgFinish.Add(len(checkTxsPerAccount))
//	//for i := 0; i < len(checkTxsPerAccount); i++ {
//	//	checkTxs := checkTxsPerAccount[i]
//	//	go func() {
//	//		// Ensure that we unlock the wait group regardless of how this goroutine completes.
//	//		defer wgFinish.Done()
//	//
//	//		// Mark that we have started and then wait till everyone else starts to increase the amount of contention
//	//		// and parallelization.
//	//		wgStart.Done()
//	//		wgStart.Wait()
//	//		for _, checkTx := range checkTxs {
//	//			resp := tApp.CheckTx(checkTx)
//	//			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
//	//		}
//	//	}()
//	//}
//	//
//	//// Wait till all the orders were placed and cancelled.
//	//wgFinish.Wait()
//	//
//	//// Advance the block and ensure that the appropriate orders were filled and cancelled.
//	//tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})
//	//for _, expectedFill := range expectedFills {
//	//	exists, amount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, expectedFill.OrderId)
//	//	require.True(t, exists)
//	//	require.Equal(t, expectedFill.Quantums, amount.ToUint64())
//	//}
//	//for _, expectedCancel := range expectedCancels {
//	//	exists, amount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, expectedCancel.OrderId)
//	//	require.False(t, exists)
//	//	require.Equal(t, uint64(0), amount.ToUint64())
//	//}
//}
