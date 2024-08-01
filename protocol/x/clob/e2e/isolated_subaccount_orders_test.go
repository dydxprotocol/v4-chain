package clob_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	sdkmath "cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func TestIsolatedSubaccountOrders(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	orderQuantums := 1_000_000_000
	PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 3},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     uint64(orderQuantums),
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 5},
		})
	PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5_IOC := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 3},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     uint64(orderQuantums),
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 5},
			TimeInForce:  clobtypes.Order_TIME_IN_FORCE_IOC,
		})
	PlaceOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_Price10_GTB5 := *clobtypes.NewMsgPlaceOrder(
		clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 0, ClobPairId: 3},
			Side:         clobtypes.Order_SIDE_SELL,
			Quantums:     uint64(orderQuantums),
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 5},
		})
	LongTermOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price5_GTBT5 := clobtypes.Order{
		OrderId: clobtypes.OrderId{
			SubaccountId: constants.Alice_Num0,
			ClientId:     0,
			OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
			ClobPairId:   3,
		},
		Side:         clobtypes.Order_SIDE_BUY,
		Quantums:     uint64(orderQuantums),
		Subticks:     10,
		GoodTilOneof: &clobtypes.Order_GoodTilBlockTime{GoodTilBlockTime: 5},
	}

	// Alice holds a long position after the match.
	Alice_Num0_IsolatedAfterMatch := satypes.Subaccount{
		Id: &constants.Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			// USDC asset position.
			testutil.CreateSingleAssetPosition(
				uint32(0),
				// Match = 10e9 * 10e-8 * 10 = 100 quantums. Fees = 0.
				// Alice is buying, subtract match quantums from asset position.
				big.NewInt(10_000_000_000-100),
			),
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			// Isolated perpetual position.
			testutil.CreateSinglePerpetualPosition(
				uint32(3),
				big.NewInt(int64(orderQuantums)),
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}

	// Bob holds a short position after the match.
	Bob_Num0_IsolatedAfterMatch := satypes.Subaccount{
		Id: &constants.Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			// USDC asset position.
			testutil.CreateSingleAssetPosition(
				uint32(0),
				// Match = 10e9 * 10e-8 * 10 = 100 quantums. Fees = 1 quantum.
				// Bob is selling, add match quantums from asset position.
				big.NewInt(10_000_000_000+99),
			),
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			// Isolated perpetual position.
			testutil.CreateSinglePerpetualPosition(
				uint32(3),
				big.NewInt(-1*int64(orderQuantums)),
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}

	// Alice holds a larger long position after the match.
	Alice_Num0_MoreIsolatedAfterMatch := satypes.Subaccount{
		Id: &constants.Alice_Num0,
		AssetPositions: []*satypes.AssetPosition{
			// USDC asset position.
			testutil.CreateSingleAssetPosition(
				uint32(0),
				// Match = 10e9 * 10e-8 * 10 = 100 quantums. Fees = 0.
				// Alice is buying, subtract match quantums from asset position.
				big.NewInt(10_000_000_000-100),
			),
		},
		PerpetualPositions: []*satypes.PerpetualPosition{
			// Isolated perpetual position.
			testutil.CreateSinglePerpetualPosition(
				uint32(3),
				// Alice buys 1 more ISO,
				big.NewInt(2*int64(orderQuantums)),
				big.NewInt(0),
				big.NewInt(0),
			),
		},
	}

	// Bob closes an isolated long position.
	Bob_Num0_CrossAfterMatch := satypes.Subaccount{
		Id: &constants.Bob_Num0,
		AssetPositions: []*satypes.AssetPosition{
			// USDC asset position.
			testutil.CreateSingleAssetPosition(
				uint32(0),
				// Match = 10e9 * 10e-8 * 10 = 100 quantums. Fees = 1.
				// Bob is selling, add match quantums from asset position.
				big.NewInt(10_000_000_000+99),
			),
		},
	}

	tests := map[string]struct {
		// Initial state
		subaccounts            []satypes.Subaccount
		perpetuals             []perptypes.Perpetual
		clobPairs              []clobtypes.ClobPair
		collateralPoolBalances map[string]int64

		// Test params
		orders []clobtypes.MsgPlaceOrder

		// Expectation
		expectedOrdersFilled           []clobtypes.OrderId
		expectedSubaccounts            []satypes.Subaccount
		expectedCollateralPoolBalances map[string]int64
		expectedErrMsg                 string
		expectedErrCode                uint32
	}{
		"Isolated subaccount will not have matches for cross-market orders": {
			subaccounts: []satypes.Subaccount{
				// Alice subaccount is isolated to ISO perpetual market with id 3.
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20, // this order is invalid, so a match won't happen
				PlaceOrder_Bob_Num0_Id0_Clob0_Sell5_Price10_GTB20,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000, // $10,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // $10,000 USDC
			},
			// No orders filled.
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			// No changes as no match should have happened.
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000,
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000,
			},
		},
		"Cross subaccount (with cross position) will not have matches for isolated-market orders": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1BTC_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5, // this order is invalid, so a match won't happen
				PlaceOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_Price10_GTB5,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000, // $10,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // $10,000 USDC
			},
			// No orders filled.
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1BTC_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			// No changes as no match should have happened.
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000,
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000,
			},
		},
		`Empty subaccount becomes isolated if an order matches for an isolated market, collateral balances
		move to isolated collateral pools`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			// Orders should match.
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5,
				PlaceOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_Price10_GTB5,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 5_000_000_000, // $5,000 USDC
			},
			expectedOrdersFilled: []clobtypes.OrderId{
				PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5.Order.OrderId,
				PlaceOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_Price10_GTB5.Order.OrderId,
			},
			expectedSubaccounts: []satypes.Subaccount{
				Alice_Num0_IsolatedAfterMatch,
				Bob_Num0_IsolatedAfterMatch,
			},
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000, // $30,000 USDC - $10,000 USDC - $10,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 24_999_999_999, // $5,000 USDC + $10,000 USDC + $10,000 USDC - fee (1 quote quantum)
			},
		},
		"Isolated subaccount closing position moves collateral back to cross collateral pool": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_1ISO_LONG_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			// Orders should match.
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5,
				PlaceOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_Price10_GTB5,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000, // $10,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 30_000_000_000, // $30,000 USDC
			},
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				Alice_Num0_MoreIsolatedAfterMatch,
				Bob_Num0_CrossAfterMatch,
			},
			expectedCollateralPoolBalances: map[string]int64{
				// $10,000 USDC + $10,000 USDC + (match) 100 quote quantums - fee (1 quote quantum)
				satypes.ModuleAddress.String(): 20_000_000_099,
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
					// $30,000 USDC - $10,000 USDC - (match) 100 quote quantums
				).String(): 19_999_999_900,
			},
		},
		`Empty subaccount will not have matches for isolated market orders if cross collateral pool does not 
		have enough balance`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5,
				PlaceOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_Price10_GTB5,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 1_000_000_000, // $1,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 30_000_000_000, // $30,000 USDC
			},
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			// No changes
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 1_000_000_000, // $1,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 30_000_000_000, // $30,000 USDC
			},
			expectedErrCode: sdkerrors.ErrPanic.ABCICode(),
			expectedErrMsg:  "insufficient funds",
		},
		`Isolated subaccount will not have matches to close isolated perpetual position if isolated collateral pool does not 
		have enough balance`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_1ISO_LONG_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			// Orders should match.
			orders: []clobtypes.MsgPlaceOrder{
				PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5,
				PlaceOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_Price10_GTB5,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_1ISO_LONG_10_000USD,
			},
			// No changes
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedErrCode: sdkerrors.ErrPanic.ABCICode(),
			expectedErrMsg:  "insufficient funds",
		},
		`Subaccount with isolated perpetual position fails to place stateful order for cross perpetual`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			orders: []clobtypes.MsgPlaceOrder{
				*clobtypes.NewMsgPlaceOrder(constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy100_Price10_GTBT15),
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
			},
			// No changes
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedErrCode: clobtypes.ErrWouldViolateIsolatedSubaccountConstraints.ABCICode(),
			expectedErrMsg:  "Order would violate isolated subaccount constraints.",
		},
		`Subaccount with cross perpetual position fails to place stateful order for isolated perpetual`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1BTC_LONG_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			orders: []clobtypes.MsgPlaceOrder{
				*clobtypes.NewMsgPlaceOrder(LongTermOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price5_GTBT5),
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1BTC_LONG_10_000USD,
			},
			// No changes
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedErrCode: clobtypes.ErrWouldViolateIsolatedSubaccountConstraints.ABCICode(),
			expectedErrMsg:  "Order would violate isolated subaccount constraints.",
		},
		`Subaccount with isolated perpetual position fails to place IOC order for cross perpetual`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			orders: []clobtypes.MsgPlaceOrder{
				// Liquidity to match the IOC order
				*clobtypes.NewMsgPlaceOrder(constants.Order_Bob_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
				*clobtypes.NewMsgPlaceOrder(constants.Order_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10_IOC),
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			// No changes
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedErrCode: clobtypes.ErrWouldViolateIsolatedSubaccountConstraints.ABCICode(),
			expectedErrMsg:  "Order would violate isolated subaccount constraints.",
		},
		`Subaccount with cross perpetual position fails to place IOC order for isolated perpetual`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1BTC_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			clobPairs: []clobtypes.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
				constants.ClobPair_3_Iso,
			},
			orders: []clobtypes.MsgPlaceOrder{
				// Liquidity to match the IOC order
				PlaceOrder_Bob_Num0_Id0_Clob3_Sell_1ISO_Price10_GTB5,
				PlaceOrder_Alice_Num0_Id0_Clob3_Buy_1ISO_Price10_GTB5_IOC,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedOrdersFilled: []clobtypes.OrderId{},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1BTC_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			// No changes
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 30_000_000_000, // $30,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 1_000_000_000, // $1,000 USDC
			},
			expectedErrCode: clobtypes.ErrWouldViolateIsolatedSubaccountConstraints.ABCICode(),
			expectedErrMsg:  "Order would violate isolated subaccount constraints.",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = constants.LiquidityTiers
						genesisState.Perpetuals = tc.perpetuals
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						// If the collateral pool address is already in bank genesis state, update it.
						foundPools := make(map[string]struct{})
						for i, bal := range genesisState.Balances {
							usdcBal, exists := tc.collateralPoolBalances[bal.Address]
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
						// If the collateral pool address is not in the bank genesis state, add it.
						for addr, bal := range tc.collateralPoolBalances {
							_, exists := foundPools[addr]
							if exists {
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
						genesisState.ClobPairs = tc.clobPairs
						genesisState.LiquidationsConfig = constants.LiquidationsConfig_FillablePrice_Max_Smmr
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				return genesis
			}).Build()
			ctx = tApp.InitChain()

			for i, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, order) {
					resp := tApp.CheckTx(checkTx)
					// Error should only be returned for the second order, as it results in a match.
					if tc.expectedErrMsg == "" || (i != len(tc.orders)-1) {
						require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
					} else {
						require.False(t, resp.IsOK())
						require.Equal(t, resp.Code, tc.expectedErrCode)
						require.Contains(t, resp.Log, tc.expectedErrMsg)
					}
				}
			}

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			for _, order := range tc.orders {
				if slices.Contains(tc.expectedOrdersFilled, order.Order.OrderId) {
					exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(
						ctx,
						order.Order.OrderId,
					)

					require.True(t, exists)
					require.Equal(t, order.Order.GetBaseQuantums(), fillAmount)
				}
			}

			for _, expectedSubaccount := range tc.expectedSubaccounts {
				require.Equal(
					t,
					expectedSubaccount,
					tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id),
				)
			}

			for addr, expectedBalance := range tc.expectedCollateralPoolBalances {
				require.Equal(
					t,
					sdkmath.NewIntFromBigInt(big.NewInt(expectedBalance)),
					tApp.App.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(addr), constants.Usdc.Denom).Amount,
				)
			}
		})
	}
}
