package isolated_pool_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// TestIsolatedPoolUnderflow_PrepareCheckStateContinues covers the same isolated-pool scenario as the
// short-term "closing a position at a profit fills when the pool is funded by the opener's deposit"
// case in x/clob/e2e/isolated_subaccount_orders_test.go, but with long-term (stateful) orders.
//
// A fresh opener (Alice) buys 1 ISO to open a position, funding the isolated collateral pool with
// her deposit. A closer (Bob) sells to fully close his isolated position at a high price, for a
// gross payout (~$19,995) that exceeds the isolated pool's seeded balance ($15,000) on its own.
//
// Short-term orders match immediately in CheckTx. Long-term orders are placed on the book and
// matched later: HandleMsgPlaceOrder (x/clob/keeper/msg_server_place_order.go) runs only in
// DeliverTx, and x/clob/ante.ClobDecorator merely records a stateful placement during CheckTx. The
// order is placed on the book (and can match) when PlaceStatefulOrdersFromLastBlock runs from
// x/clob/abci.go's PrepareCheckState, which BaseApp.Commit invokes directly. This test asserts the
// crossing match settles there, funded by the opener's deposit before the closer's payout is drawn
// from the pool.
func TestIsolatedPoolUnderflow_PrepareCheckStateContinues(t *testing.T) {
	orderQuantums := 1_000_000_000

	// Long-term analogs of PlaceOrder_Alice/Bob_Num0_Id0_Clob3_Buy/Sell_1ISO_HighPrice_GTB5 from
	// isolated_subaccount_orders_test.go, modeled on the
	// LongTermOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price5_GTBT5 fixture in that same file.
	LongTermOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_HighPrice_GTBT5 := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   3,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     uint64(orderQuantums),
		Subticks:     1_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
	}
	LongTermOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_HighPrice_GTBT5 := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Bob_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   3,
		},
		Side:         clobtypes.Order_SIDE_SELL,
		Quantums:     uint64(orderQuantums),
		Subticks:     1_000_000_000,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
	}

	// Alice opens an isolated position at a price far above the ISO-USD oracle price, paying nearly
	// her entire prior balance for it.
	Alice_Num0_IsolatedAfterHighPriceMatch := satypes.Subaccount{
		Id: &constants.Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			testutil.CreateSingleAssetPosition(
				uint32(0),
				// Match = 1e9 * 1e9 * 1e-8 = 10_000_000_000 quantums. Maker fee = -1_100_000 (a rebate).
				big.NewInt(1_100_000),
			),
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				uint32(3),
				big.NewInt(int64(orderQuantums)),
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}

	// Bob closes his isolated position at a profit: he receives more from the match than his
	// isolated collateral pool is seeded with, so the payout must be covered by Alice's deposit.
	Bob_Num0_ClosedAtProfitAfterMatch := satypes.Subaccount{
		Id: &constants.Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			testutil.CreateSingleAssetPosition(
				uint32(0),
				// Match = 10_000_000_000 quantums. Taker fee = 5_000_000.
				big.NewInt(19_995_000_000),
			),
		},
	}

	crossPoolAddr := satypes.ModuleAddress.String()
	isoPoolAddr := authtypes.NewModuleAddress(
		satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
	).String()

	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					constants.Alice_Num0_10_000USD,
					constants.Bob_Num0_1ISO_LONG_10_000USD,
				}
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
					constants.IsoUsd_IsolatedMarket,
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *banktypes.GenesisState) {
				collateralPoolBalances := map[string]int64{
					// $20,000 USDC.
					crossPoolAddr: 20_000_000_000,
					// $15,000 USDC: below Bob's ~$19,995 gross payout on close.
					isoPoolAddr: 15_000_000_000,
				}
				foundPools := make(map[string]struct{})
				for i, bal := range genesisState.Balances {
					usdcBal, exists := collateralPoolBalances[bal.Address]
					if exists {
						foundPools[bal.Address] = struct{}{}
						genesisState.Balances[i] = banktypes.Balance{
							Address: bal.Address,
							Coins: sdktypes.Coins{
								sdktypes.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(usdcBal)),
							},
						}
					}
				}
				for addr, bal := range collateralPoolBalances {
					if _, exists := foundPools[addr]; exists {
						continue
					}
					genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
						Address: addr,
						Coins: sdktypes.Coins{
							sdktypes.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(bal)),
						},
					})
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *clobtypes.GenesisState) {
				genesisState.ClobPairs = []clobtypes.ClobPair{
					constants.ClobPair_Btc,
					constants.ClobPair_Eth,
					constants.ClobPair_3_Iso,
				}
				genesisState.LiquidationsConfig = constants.LiquidationsConfig_FillablePrice_Max_Smmr
				genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	// CheckTx-place both long-term orders. Since HandleMsgPlaceOrder only runs in DeliverTx and
	// ClobDecorator does not place stateful orders on the book during CheckTx, neither order is
	// matched here: both simply rest as pending stateful order placements.
	for _, order := range []clobtypes.Order{
		LongTermOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_HighPrice_GTBT5,
		LongTermOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_HighPrice_GTBT5,
	} {
		for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
			ctx,
			tApp.App,
			*clobtypes.NewMsgPlaceOrder(order),
		) {
			resp := tApp.CheckTx(checkTx)
			require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
		}
	}

	// Advancing to block 2 delivers both stateful order placements (writing them to committed
	// state), commits the block, and -- as part of that same Commit call -- runs PrepareCheckState
	// for block 3, which replays last block's long-term order placements onto the memclob and
	// triggers the crossing match. The isolated pool is funded by the opener's deposit before the
	// closer's payout is drawn, so the match settles rather than underflowing the pool.
	require.NotPanics(t, func() {
		ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
	})

	for _, orderId := range []clobtypes.OrderId{
		LongTermOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_HighPrice_GTBT5.OrderId,
		LongTermOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_HighPrice_GTBT5.OrderId,
	} {
		exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
		require.True(t, exists)
		require.Equal(t, satypes.BaseQuantums(orderQuantums), fillAmount)
	}

	require.Equal(
		t,
		Alice_Num0_IsolatedAfterHighPriceMatch,
		tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0),
	)
	require.Equal(
		t,
		Bob_Num0_ClosedAtProfitAfterMatch,
		tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0),
	)

	require.Equal(
		t,
		// $20,000 - $10,000 (Alice's deposit, opened) + $19,995 (Bob's payout, closed)
		sdkmath.NewInt(29_995_000_000),
		tApp.App.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(crossPoolAddr), constants.Usdc.Denom).Amount,
	)
	require.Equal(
		t,
		// $15,000 + $10,000 (Alice's deposit) - $19,995 (Bob's payout) - $3,900 net fees.
		sdkmath.NewInt(5_001_100_000),
		tApp.App.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(isoPoolAddr), constants.Usdc.Denom).Amount,
	)
}
