package e2e_test

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/config"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"

	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	sendingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb/testutil"
)

var (
	Clob_0 = testapp.MustGetClobPairsFromGenesis(testapp.DefaultGenesis())[0]
)

func TestParallelAnteHandler_ClobAndOther(t *testing.T) {
	// We concurrently send transfers and clob messages for the same accounts primarily relying on go's `-race` flag
	// to detect data races
	r := testutil.NewRand()
	simAccounts := simtypes.RandomAccounts(r, 10)
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
							&constants.TDai_Asset_500_000,
						},
					})
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *banktypes.GenesisState) {
				for _, simAccount := range simAccounts {
					genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
						Address: sdktypes.AccAddress(simAccount.PubKey.Address()).String(),
						Coins: sdktypes.NewCoins(sdktypes.NewInt64Coin(
							constants.TDai.Denom,
							constants.TDai_Asset_500_000.Quantums.BigInt().Int64(),
						)),
					})
				}
			},
		)
		return genesis
	}).WithNonDeterminismChecksEnabled(false).Build()
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

	accounts := make([]sdktypes.AccountI, len(simAccounts))
	for i, simAccount := range simAccounts {
		accounts[i] = tApp.App.AccountKeeper.GetAccount(ctx, simAccount.Address)
	}

	// We specifically use an atomic to ensure that we aren't providing any synchronization between the threads
	// maximizing any data races that could exist. The wait group is only used to synchronize the testing thread
	// when the other threads are done.
	blockLimitReached := atomic.Bool{}
	blockLimitReached.Store(false)
	wg := sync.WaitGroup{}

	// Start a block advancement thread.
	wg.Add(1)
	blockHeight := atomic.Uint64{}
	blockHeight.Store(2)

	go func() {
		defer wg.Done()
		defer func() {
			blockLimitReached.Store(true)
		}()
		for i := uint32(3); i < 50; i++ {
			tApp.AdvanceToBlock(i, testapp.AdvanceToBlockOptions{})
			blockHeight.Store(uint64(i))
		}
	}()

	// Start threads that will withdraw funds from each account.
	transferCounts := make([]atomic.Uint64, len(simAccounts))
	for i := 0; i < len(simAccounts); i++ {
		ii := i
		simAccount := simAccounts[i]
		account := accounts[i]

		wg.Add(1)
		go func() {
			defer wg.Done()

			for sequenceNumber := uint64(0); !blockLimitReached.Load(); sequenceNumber++ {
				checkTx, err := sims.GenSignedMockTx(
					rand.New(rand.NewSource(42)),
					tApp.App.TxConfig(),

					[]sdktypes.Msg{
						&sendingtypes.MsgWithdrawFromSubaccount{
							Sender: satypes.SubaccountId{
								Owner:  simAccount.Address.String(),
								Number: 0,
							},
							Recipient: simAccount.Address.String(),
							AssetId:   constants.TDai.Id,
							Quantums:  constants.TDai_Asset_1.Quantums.BigInt().Uint64(),
						},
					},
					constants.TestFeeCoins_5Cents,
					120_000,
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
				transferCounts[ii].Add(1)
			}
		}()
	}

	// Start threads with placing and cancelling orders over the same set of accounts.
	for i := 0; i < len(simAccounts); i++ {
		simAccount := simAccounts[i]
		account := accounts[i]

		wg.Add(1)
		go func() {
			defer wg.Done()

			for clientId := uint32(0); !blockLimitReached.Load(); clientId++ {
				orderId := clobtypes.OrderId{
					SubaccountId: satypes.SubaccountId{
						Owner:  sdktypes.MustBech32ifyAddressBytes(config.Bech32PrefixAccAddr, simAccount.PubKey.Address()),
						Number: 0,
					},
					ClientId:   clientId,
					ClobPairId: 0,
				}

				gtb := uint32(blockHeight.Load()) + 20
				order := testapp.MustScaleOrder(
					clobtypes.Order{
						OrderId:      orderId,
						Side:         clobtypes.Order_SIDE_BUY,
						Quantums:     1,
						Subticks:     10,
						GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: gtb},
					},
					Clob_0)

				checkTx, err := sims.GenSignedMockTx(
					rand.New(rand.NewSource(42)),
					tApp.App.TxConfig(),
					[]sdktypes.Msg{
						clobtypes.NewMsgPlaceOrder(order),
					},
					nil,
					0,
					ctx.ChainID(),
					[]uint64{account.GetAccountNumber()},
					[]uint64{0},
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

				checkTx, err = sims.GenSignedMockTx(
					rand.New(rand.NewSource(42)),
					tApp.App.TxConfig(),
					[]sdktypes.Msg{
						clobtypes.NewMsgCancelOrderShortTerm(orderId, gtb),
					},
					nil,
					0,
					ctx.ChainID(),
					[]uint64{account.GetAccountNumber()},
					[]uint64{0},
					simAccount.PrivKey,
				)
				if err != nil {
					panic(err)
				}
				bytes, err = tApp.App.TxConfig().TxEncoder()(checkTx)
				if err != nil {
					panic(err)
				}
				resp = tApp.CheckTx(abcitypes.RequestCheckTx{
					Tx:   bytes,
					Type: abcitypes.CheckTxType_New,
				})
				require.Conditionf(t, resp.IsOK, "Expected response to be ok: %+v", resp)
			}
		}()
	}

	wg.Wait()

	// Deliver the last of the transactions.
	ctx = tApp.AdvanceToBlock(50, testapp.AdvanceToBlockOptions{})

	// Verify the transfers occurred.
	for i := 0; i < len(simAccounts); i++ {
		account := accounts[i]
		subAccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, satypes.SubaccountId{
			Owner:  account.GetAddress().String(),
			Number: 0,
		})
		require.Equal(
			t,
			[]*satypes.AssetPosition{{
				AssetId: constants.TDai.Id,
				Quantums: dtypes.NewIntFromUint64(
					constants.TDai_Asset_500_000.Quantums.BigInt().Uint64() -
						transferCounts[i].Load()*constants.TDai_Asset_1.Quantums.BigInt().Uint64()),
			}},
			subAccount.AssetPositions,
		)
	}
}
