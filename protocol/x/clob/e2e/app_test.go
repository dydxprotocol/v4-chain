package clob_test

import (
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/rand"
	"gopkg.in/typ.v4/slices"

	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	epochtypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

var (
	Clob_0                                            = testapp.MustGetClobPairsFromGenesis(testapp.DefaultGenesis())[0]
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 5},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB23 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 23},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB24 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 24},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB47 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 47},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num1, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy6_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     6,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	// replacement of above order with smaller quantums
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB21 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 21},
		},
		testapp.DefaultGenesis(),
	))
	// replacement of order with larger quantums
	PlaceOrder_Alice_Num0_Id0_Clob0_Buy7_Price10_GTB21 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     7,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 21},
		},
		testapp.DefaultGenesis(),
	))
	// replacement of order on opposite side
	PlaceOrder_Alice_Num0_Id0_Clob0_Sell6_Price10_GTB21 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     6,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 21},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
	CancelOrder_Alice_Num0_Id0_Clob0_GTB47 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			ClobPairId:   0,
		},
		47,
	)
	CancelOrder_Alice_Num0_Id0_Clob0_GTB23 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			ClobPairId:   0,
		},
		23,
	)
	CancelOrder_Alice_Num1_Id0_Clob0_GTB20 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num1,
			ClientId:     0,
			ClobPairId:   0,
		},
		20,
	)
	CancelOrder_Alice_Num0_Id1_Clob0_GTB20 = *clobtypes.NewMsgCancelOrderShortTerm(
		clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			ClobPairId:   1,
		},
		20,
	)
	PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
		testapp.DefaultGenesis(),
	))
	PlaceOrder_Bob_Num0_Id0_Clob0_Sell4_Price10_GTB20 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     4,
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

	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB5,
		testapp.DefaultGenesis(),
	))
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT5,
		testapp.DefaultGenesis(),
	))
	LongTermPlaceOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob1_Buy5_Price10_GTBT5,
		testapp.DefaultGenesis(),
	))
	LongTermPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
		testapp.DefaultGenesis(),
	))
	LongTermPlaceOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTBT5 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.LongTermOrder_Alice_Num1_Id0_Clob0_Buy5_Price10_GTBT5,
		testapp.DefaultGenesis(),
	))
	ConditionalPlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
		testapp.DefaultGenesis(),
	))
	ConditionalPlaceOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15 = *clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
		constants.ConditionalOrder_Alice_Num1_Id0_Clob0_Sell5_Price10_GTB15,
		testapp.DefaultGenesis(),
	))

	BatchCancel_Alice_Num0_Clob0_1_2_3_GTB5 = *clobtypes.NewMsgBatchCancel(
		constants.Alice_Num0,
		[]clobtypes.OrderBatch{
			{
				ClobPairId: 0,
				ClientIds:  []uint32{1, 2, 3},
			},
		},
		5,
	)
	BatchCancel_Alice_Num0_Clob0_1_2_3_GTB47 = *clobtypes.NewMsgBatchCancel(
		constants.Alice_Num0,
		[]clobtypes.OrderBatch{
			{
				ClobPairId: 0,
				ClientIds:  []uint32{1, 2, 3},
			},
		},
		47,
	)
	BatchCancel_Alice_Num0_Clob0_1_2_3_GTB20 = *clobtypes.NewMsgBatchCancel(
		constants.Alice_Num0,
		[]clobtypes.OrderBatch{
			{
				ClobPairId: 0,
				ClientIds:  []uint32{1, 2, 3},
			},
		},
		20,
	)
	BatchCancel_Alice_Num0_Clob1_1_2_3_GTB20 = *clobtypes.NewMsgBatchCancel(
		constants.Alice_Num0,
		[]clobtypes.OrderBatch{
			{
				ClobPairId: 1,
				ClientIds:  []uint32{1, 2, 3},
			},
		},
		20,
	)
	BatchCancel_Alice_Num1_Clob0_1_2_3_GTB20 = *clobtypes.NewMsgBatchCancel(
		constants.Alice_Num1,
		[]clobtypes.OrderBatch{
			{
				ClobPairId: 0,
				ClientIds:  []uint32{1, 2, 3},
			},
		},
		20,
	)

	// Leverage update message constants
	UpdateLeverage_Alice_Num0_PerpId0_Lev5 = clobtypes.MsgUpdateLeverage{
		SubaccountId: &constants.Alice_Num0,
		ClobPairLeverage: []*clobtypes.LeverageEntry{
			{
				ClobPairId:   0,
				CustomImfPpm: 200_000,
			},
		},
	}
	UpdateLeverage_Alice_Num0_PerpId1_Lev10 = clobtypes.MsgUpdateLeverage{
		SubaccountId: &constants.Alice_Num0,
		ClobPairLeverage: []*clobtypes.LeverageEntry{
			{
				ClobPairId:   1,
				CustomImfPpm: 100_000,
			},
		},
	}
	UpdateLeverage_Alice_Num1_PerpId0_Lev4 = clobtypes.MsgUpdateLeverage{
		SubaccountId: &constants.Alice_Num1,
		ClobPairLeverage: []*clobtypes.LeverageEntry{
			{
				ClobPairId:   0,
				CustomImfPpm: 250_000,
			},
		},
	}
	UpdateLeverage_Bob_Num0_PerpId0_Lev5 = clobtypes.MsgUpdateLeverage{
		SubaccountId: &constants.Bob_Num0,
		ClobPairLeverage: []*clobtypes.LeverageEntry{
			{
				ClobPairId:   0,
				CustomImfPpm: 200_000,
			},
		},
	}
)

func TestHydrationInPreBlocker(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis tmtypes.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
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
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					constants.Carl_Num0_100000USD,
					constants.Dave_Num0_10000USD,
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *clobtypes.GenesisState) {
				genesisState.ClobPairs = []clobtypes.ClobPair{
					constants.ClobPair_Btc,
				}
				genesisState.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
			},
		)
		return genesis
	}).WithNonDeterminismChecksEnabled(false).Build()

	// Let's add some pre-existing orders to state.
	// Note that the order is not added to memclob.
	tApp.App.ClobKeeper.SetLongTermOrderPlacement(
		tApp.App.NewUncachedContext(false, tmproto.Header{}),
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
		1,
	)
	tApp.App.ClobKeeper.AddStatefulOrderIdExpiration(
		tApp.App.NewUncachedContext(false, tmproto.Header{}),
		time.Unix(50, 0),
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
	)

	// Advance one block so that pre blocker is called and clob is hydrated.
	_ = tApp.InitChain()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Order should exist in state
	_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(
		ctx,
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
	)
	require.True(t, found)

	// Order should be on the orderbook
	_, found = tApp.App.ClobKeeper.MemClob.GetOrder(
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
	)
	require.True(t, found)
}

func TestHydrationWithMatchPreBlocker(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis tmtypes.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
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
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					constants.Carl_Num0_100000USD,
					constants.Dave_Num0_500000USD,
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *clobtypes.GenesisState) {
				genesisState.ClobPairs = []clobtypes.ClobPair{
					constants.ClobPair_Btc,
				}
				genesisState.LiquidationsConfig = clobtypes.LiquidationsConfig_Default
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *feetierstypes.GenesisState) {
				genesisState.Params = constants.PerpetualFeeParamsNoFee
			},
		)
		return genesis
	}).WithNonDeterminismChecksEnabled(false).Build()

	// 1. Let's add some pre-existing orders to state before clob is initialized.
	tApp.App.ClobKeeper.SetLongTermOrderPlacement(
		tApp.App.NewUncachedContext(false, tmproto.Header{}),
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
		1,
	)
	tApp.App.ClobKeeper.AddStatefulOrderIdExpiration(
		tApp.App.NewUncachedContext(false, tmproto.Header{}),
		time.Unix(10, 0),
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
	)

	// Let's add a crossing order to state.
	tApp.App.ClobKeeper.SetLongTermOrderPlacement(
		tApp.App.NewUncachedContext(false, tmproto.Header{}),
		constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
		1,
	)
	tApp.App.ClobKeeper.AddStatefulOrderIdExpiration(
		tApp.App.NewUncachedContext(false, tmproto.Header{}),
		time.Unix(10, 0),
		constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
	)

	// 2. Advance one block so that pre blocker is called and clob is hydrated.
	ctx := tApp.InitChain()

	// Here, PreBlocker has been called and Carl's and Dave's orders are placed against the orderbook.
	// They should generate a match in the local operations queue, but their state changes should've been discarded since
	// preblocker happens during deliver state and context was cached with IsCheckTx set to true.

	// Make sure order still exists in state, with a fill amount of 0.
	// Note that `ctx` is the check tx context, so need to read from the uncached cms to make sure changes are discarded.
	uncachedCtx := tApp.App.NewUncachedContext(false, tmproto.Header{})
	_, found := tApp.App.ClobKeeper.GetLongTermOrderPlacement(
		uncachedCtx,
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
	)
	require.True(t, found)
	fillAmount := tApp.App.ClobKeeper.MemClob.GetOrderFilledAmount(
		uncachedCtx,
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
	)
	require.Equal(t, satypes.BaseQuantums(0), fillAmount)

	_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(
		uncachedCtx,
		constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
	)
	require.True(t, found)
	fillAmount = tApp.App.ClobKeeper.MemClob.GetOrderFilledAmount(
		uncachedCtx,
		constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
	)
	require.Equal(t, satypes.BaseQuantums(0), fillAmount)

	// Make sure orders are not on the orderbook.
	_, found = tApp.App.ClobKeeper.MemClob.GetOrder(
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
	)
	require.False(t, found)

	_, found = tApp.App.ClobKeeper.MemClob.GetOrder(
		constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
	)
	require.False(t, found)

	// Make sure match is in the operations queue.
	operations := tApp.App.ClobKeeper.MemClob.GetOperationsRaw(ctx)
	require.Len(t, operations, 1)

	// Advance to the next block to persist the matches.
	ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	// Order should not exist in state because they are filly filled.
	_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(
		ctx,
		constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
	)
	require.False(t, found)

	_, found = tApp.App.ClobKeeper.GetLongTermOrderPlacement(
		ctx,
		constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
	)
	require.False(t, found)

	// Carl and Dave's state should get updated accordingly.
	carl := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Carl_Num0)
	require.Equal(t, satypes.Subaccount{
		Id: &constants.Carl_Num0,
		AssetPositions: []*satypes.AssetPosition{
			testutil.CreateSingleAssetPosition(
				0,
				big.NewInt(100_000_000_000-50_000_000_000),
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
	}, carl)

	dave := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Dave_Num0)
	require.Equal(t, satypes.Subaccount{
		Id: &constants.Dave_Num0,
		AssetPositions: []*satypes.AssetPosition{
			testutil.CreateSingleAssetPosition(
				0,
				big.NewInt(500_000_000_000+50_000_000_000),
			),
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			testutil.CreateSinglePerpetualPosition(
				0,
				big.NewInt(-100_000_000),
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}, dave)

	require.Empty(t, tApp.App.ClobKeeper.MemClob.GetOperationsRaw(ctx))
}

// We place 300 orders that match and 700 orders followed by their cancellations concurrently.
//
// This test heavily relies on golangs race detector to validate memory reads and writes are properly ordered.
func TestConcurrentMatchesAndCancels(t *testing.T) {
	r := rand.NewRand()
	simAccounts := simtypes.RandomAccounts(r, 1000)
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis tmtypes.GenesisDoc) {
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
			expectedFills[i] = testapp.MustScaleOrder(clobtypes.Order{
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
					AccAddressForSigning: expectedFills[i].OrderId.SubaccountId.Owner,
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
				expectedCancels[idx] = testapp.MustScaleOrder(clobtypes.Order{
					OrderId:      orderId,
					Side:         clobtypes.Order_SIDE_BUY,
					Quantums:     1,
					Subticks:     10,
					GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
				},
					Clob_0)
			case 1:
				expectedCancels[idx] = testapp.MustScaleOrder(clobtypes.Order{
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
					AccAddressForSigning: orderId.SubaccountId.Owner,
				},
				privKeySupplier,
				placeOrderMsg,
			)
			cancelOrderMsg := clobtypes.NewMsgCancelOrderShortTerm(orderId, 20)
			checkTxsPerAccount[i][1] = testapp.MustMakeCheckTxWithPrivKeySupplier(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: orderId.SubaccountId.Owner,
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
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()
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

			proposal, err := tApp.PrepareProposal()
			require.NoError(t, err)
			proposal.Txs[0] = testtx.MustGetTxBytes(
				&clobtypes.MsgProposedOperations{
					OperationsQueue: operationsQueue,
				},
			)

			tApp.AdvanceToBlock(3,
				testapp.AdvanceToBlockOptions{
					RequestProcessProposalTxsOverride: proposal.Txs,
					ValidateFinalizeBlock: func(
						ctx sdktypes.Context,
						request abcitypes.RequestFinalizeBlock,
						response abcitypes.ResponseFinalizeBlock,
					) (haltchain bool) {
						txResult := response.TxResults[0]
						require.Condition(t, txResult.IsErr, "Expected DeliverTx to fail but passed %+v", response)
						require.Contains(t, txResult.Log, "invalid pubkey: MsgProposedOperations is invalid")
						return true
					},
				},
			)
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
			msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
			appOpts := map[string]interface{}{
				indexer.MsgSenderInstanceForTest: msgSender,
			}
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()
			tApp.InitChain()
			tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			proposal, err := tApp.PrepareProposal()
			require.NoError(t, err)
			proposal.Txs[0] = tc.proposedOperationsTx

			tApp.AdvanceToBlock(
				3,
				testapp.AdvanceToBlockOptions{
					RequestProcessProposalTxsOverride: proposal.Txs,
					ValidateFinalizeBlock: func(
						ctx sdktypes.Context,
						request abcitypes.RequestFinalizeBlock,
						response abcitypes.ResponseFinalizeBlock,
					) (haltchain bool) {
						txResult := response.TxResults[0]
						require.Condition(t, txResult.IsErr, "Expected DeliverTx to fail but passed %+v", response)
						require.Contains(t, txResult.Log, "Error: no signatures supplied: MsgProposedOperations is invalid")
						return true
					},
				},
			)
		})
	}
}

func TestStats(t *testing.T) {
	msgSender := msgsender.NewIndexerMessageSenderInMemoryCollector()
	appOpts := map[string]interface{}{
		indexer.MsgSenderInstanceForTest: msgSender,
	}
	tApp := testapp.NewTestAppBuilder(t).WithAppOptions(appOpts).Build()

	// Epochs start at block height 2.
	startTime := time.Unix(10, 0).UTC()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
		BlockTime: startTime,
	})

	aliceAddress := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Alice_Num0).Id.MustGetAccAddress().String()
	bobAddress := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, constants.Bob_Num0).Id.MustGetAccAddress().String()

	createAliceBuyOrder := func(clientId uint32) *clobtypes.MsgPlaceOrder {
		return clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
		return clobtypes.NewMsgPlaceOrder(testapp.MustScaleOrder(
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
