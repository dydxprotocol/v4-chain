package clob_test

import (
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
)

func TestReduceOnlyOrders(t *testing.T) {
	tests := map[string]struct {
		subaccounts          []satypes.Subaccount
		ordersForFirstBlock  []clobtypes.Order
		ordersForSecondBlock []clobtypes.Order

		priceUpdateForFirstBlock  map[uint32]ve.VEPricePair
		priceUpdateForSecondBlock map[uint32]ve.VEPricePair

		crashingAppCheckTxNonDeterminismChecksDisabled bool

		expectedInTriggeredStateAfterBlock map[uint32]map[clobtypes.OrderId]bool

		expectedOrderOnMemClob  map[clobtypes.OrderId]bool
		expectedOrderFillAmount map[clobtypes.OrderId]uint64
		expectedSubaccounts     []satypes.Subaccount
	}{
		"IOC Reduce only order partially matches short term order same block, maker order fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{},

			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          false,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          100,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 100,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(95_000_550_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(100),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(504_997_500_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_900),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"IOC Reduce only order partially matches short term order second block, maker order fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},

			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          false,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy10_Price500000_GTB20.OrderId:          100,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 100,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(95_000_550_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(100),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(504_997_500_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_900),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"IOC Reduce only order partially matches short term order second block, maker order partially filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},

			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          true,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          150,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_IOC_RO.OrderId: 150,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(9_250_0825_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(150),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(507_496_250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_850),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"FOK Reduce only order fully matches short term order second block, maker order partially filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},

			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

			// Crashing app checks have to be disabled because the FOK order will not match
			// with an empty orderbook and fail to be placed.
			crashingAppCheckTxNonDeterminismChecksDisabled: true,

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          true,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_FOK_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId: 150,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(9_250_0825_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(150),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(507_496_250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_850),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"FOK Reduce only order fully matches short term order same block, maker order partially filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_100000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20,
					testapp.DefaultGenesis(),
				),
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			ordersForSecondBlock: []clobtypes.Order{},

			priceUpdateForFirstBlock:  map[uint32]ve.VEPricePair{},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

			// Crashing app checks don't need to be disabled since matches occur in same block.
			crashingAppCheckTxNonDeterminismChecksDisabled: false,

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId:          true,
				constants.Order_Alice_Num1_Id1_Clob0_Sell15_Price500000_GTB20_FOK_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy80_Price500000_GTB20.OrderId: 150,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(9_250_0825_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(150),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(507_496_250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(99_999_850),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"Conditional FOK Reduce only order fully matches short term order same block, maker order partially filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_500000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_FOK_RO,
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10,
			},
			ordersForSecondBlock: []clobtypes.Order{},

			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_FOK_RO.OrderId: true,
				},
				3: {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_FOK_RO.OrderId: true,
				},
			},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10.OrderId:                                true,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_FOK_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price500000_GTB10.OrderId:                                50_000_000,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_FOK_RO.OrderId: 50_000_000,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(250_027_500_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(50_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(749_875_000_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(50_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
		"Conditional IOC Reduce only order partially matches short term order same block, maker order fully filled": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_500000USD,
				constants.Alice_Num1_1BTC_Long_500_000USD,
			},
			ordersForFirstBlock: []clobtypes.Order{
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO,
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10,
			},
			ordersForSecondBlock: []clobtypes.Order{},
			priceUpdateForFirstBlock: map[uint32]ve.VEPricePair{
				0: {
					SpotPrice: 5_000_300_000,
					PnlPrice:  5_000_300_000,
				},
			},
			priceUpdateForSecondBlock: map[uint32]ve.VEPricePair{},

			expectedInTriggeredStateAfterBlock: map[uint32]map[clobtypes.OrderId]bool{
				2: {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: true,
				},
				3: {
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: true,
				},
			},

			expectedOrderOnMemClob: map[clobtypes.OrderId]bool{
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10.OrderId:                              false,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: false,
			},
			expectedOrderFillAmount: map[clobtypes.OrderId]uint64{
				constants.Order_Carl_Num0_Id0_Clob0_Buy025BTC_Price500000_GTB10.OrderId:                              25_000_000,
				constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO.OrderId: 25_000_000,
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(375_013_750_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(25_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(624_937_500_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(75_000_000),
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithCrashingAppCheckTxNonDeterminismChecksEnabled(!tc.crashingAppCheckTxNonDeterminismChecksDisabled).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
							genesisState.Perpetuals = []perptypes.Perpetual{
								constants.BtcUsd_20PercentInitial_10PercentMaintenance,
								constants.EthUsd_20PercentInitial_10PercentMaintenance,
							}
						},
					)
					return genesis
				}).Build()

			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)

			require.NoError(t, conversionErr)

			tApp.App.RatelimitKeeper.SetSDAIPrice(tApp.App.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.App.RatelimitKeeper.SetAssetYieldIndex(tApp.App.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.CrashingApp.RatelimitKeeper.SetSDAIPrice(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.CrashingApp.RatelimitKeeper.SetAssetYieldIndex(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.NoCheckTxApp.RatelimitKeeper.SetSDAIPrice(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.NoCheckTxApp.RatelimitKeeper.SetAssetYieldIndex(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.ParallelApp.RatelimitKeeper.SetSDAIPrice(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.ParallelApp.RatelimitKeeper.SetAssetYieldIndex(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			ctx := tApp.InitChain()

			// Create all orders.
			deliverTxsOverride := make([][]byte, 0)
			deliverTxsOverride = append(
				deliverTxsOverride,
				constants.ValidEmptyMsgProposedOperationsTxBytes,
			)

			for _, order := range tc.ordersForFirstBlock {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

					if order.IsStatefulOrder() {
						deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
					}
				}
			}

			// Add an empty premium vote.
			deliverTxsOverride = append(deliverTxsOverride, constants.EmptyMsgAddPremiumVotesTxBytes)

			// Add the price update.

			// txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
			// require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForFirstBlock))
			// priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
			// require.NoError(t, err)
			_, extCommitBz, err := vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdateForFirstBlock,
				"",
				tApp.GetHeader().Height,
			)
			require.NoError(t, err)
			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[2]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 3)
				}
			}

			// Create all orders.
			deliverTxsOverride = make([][]byte, 0)
			deliverTxsOverride = append(
				deliverTxsOverride,
				constants.ValidEmptyMsgProposedOperationsTxBytes,
			)

			// Place orders for second block
			for _, order := range tc.ordersForSecondBlock {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

					if order.IsStatefulOrder() {
						deliverTxsOverride = append(deliverTxsOverride, checkTx.Tx)
					}
				}
			}

			// Add an empty premium vote.
			deliverTxsOverride = append(deliverTxsOverride, constants.EmptyMsgAddPremiumVotesTxBytes)

			// Add the price update.

			// txBuilder := encoding.GetTestEncodingCfg().TxConfig.NewTxBuilder()
			// require.NoError(t, txBuilder.SetMsgs(tc.priceUpdateForSecondBlock))
			// priceUpdateTxBytes, err := encoding.GetTestEncodingCfg().TxConfig.TxEncoder()(txBuilder.GetTx())
			// require.NoError(t, err)

			_, extCommitBz, err = vetesting.GetInjectedExtendedCommitInfoForTestApp(
				&tApp.App.ConsumerKeeper,
				ctx,
				tc.priceUpdateForFirstBlock,
				"",
				tApp.GetHeader().Height+1,
			)
			require.NoError(t, err)
			deliverTxsOverride = append([][]byte{extCommitBz}, deliverTxsOverride...)

			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
				DeliverTxsOverride: deliverTxsOverride,
			})

			if expectedTriggeredOrders, ok := tc.expectedInTriggeredStateAfterBlock[3]; ok {
				for orderId, triggered := range expectedTriggeredOrders {
					require.Equal(t, triggered, tApp.App.ClobKeeper.IsConditionalOrderTriggered(ctx, orderId), "Block %d", 4)
				}
			}

			// Verify expectations.
			for orderId, exists := range tc.expectedOrderOnMemClob {
				_, existsOnMemclob := tApp.App.ClobKeeper.MemClob.GetOrder(ctx, orderId)
				require.Equal(t, exists, existsOnMemclob)
			}

			for orderId, expectedFillAmount := range tc.expectedOrderFillAmount {
				exists, fillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
				require.True(t, exists)
				require.Equal(t, expectedFillAmount, fillAmount.ToUint64())
			}

			for _, subaccount := range tc.expectedSubaccounts {
				actualSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *subaccount.Id)
				require.Equal(t, subaccount, actualSubaccount)
			}
		})
	}
}

func TestReduceOnlyOrderFailure(t *testing.T) {
	tests := map[string]struct {
		subaccounts []satypes.Subaccount
		orders      []clobtypes.Order
		errorMsg    []string
	}{
		"Zero perpetual position subaccount position cannot place sell RO order": {
			orders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Sell10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		"Zero perpetual position subaccount position cannot place buy RO order": {
			orders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob1_Buy10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrReduceOnlyWouldIncreasePositionSize.Error(),
			},
		},
		"FOK Reduce only order is placed but does not match immediately and is cancelled.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_1BTC_Short_100_000USD,
			},
			orders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.Order_Alice_Num1_Id1_Clob0_Buy10_Price15_GTB20_FOK_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				clobtypes.ErrFokOrderCouldNotBeFullyFilled.Error(),
			},
		},
		"Conditional FOK Reduce only order is placed successfully.": {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num1_1BTC_Short_100_000USD,
			},
			orders: []clobtypes.Order{
				testapp.MustScaleOrder(
					constants.ConditionalOrder_Alice_Num1_Id1_Clob0_Sell05BTC_Price500000_GTBT20_TP_50001_IOC_RO,
					testapp.DefaultGenesis(),
				),
			},
			errorMsg: []string{
				"",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				if len(tc.subaccounts) > 0 {
					testapp.UpdateGenesisDocWithAppStateForModule(
						&genesis,
						func(genesisState *satypes.GenesisState) {
							genesisState.Subaccounts = tc.subaccounts
						},
					)
				}
				return genesis
			}).Build()

			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)

			require.NoError(t, conversionErr)

			tApp.App.RatelimitKeeper.SetSDAIPrice(tApp.App.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.App.RatelimitKeeper.SetAssetYieldIndex(tApp.App.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.CrashingApp.RatelimitKeeper.SetSDAIPrice(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.CrashingApp.RatelimitKeeper.SetAssetYieldIndex(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.NoCheckTxApp.RatelimitKeeper.SetSDAIPrice(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.NoCheckTxApp.RatelimitKeeper.SetAssetYieldIndex(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.ParallelApp.RatelimitKeeper.SetSDAIPrice(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.ParallelApp.RatelimitKeeper.SetAssetYieldIndex(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			ctx := tApp.InitChain()

			for idx, order := range tc.orders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)

					if tc.errorMsg[idx] == "" {
						require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
					} else {
						require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
						require.Contains(
							t,
							resp.Log,
							tc.errorMsg[idx],
						)
					}
				}
			}
		})
	}
}

func TestReduceOnlyOrderReplacement(t *testing.T) {
	tests := map[string]struct {
		subaccounts        []satypes.Subaccount
		firstOrders        []clobtypes.Order
		secondOrders       []clobtypes.Order
		secondOrdersErrors []string

		expectedOrderFillAmounts map[uint32]map[clobtypes.OrderId]uint64
	}{
		`A regular order is partially filled. Replacement FOK RO order fails because it is immediate execution.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
				constants.Carl_Num0_100000USD,
			},

			firstOrders: []clobtypes.Order{
				// Regular order on the opposite side of the following replacement RO order.
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price500000_GTB20,
				// Partial match for the above order. 70 quantums are matched. Thus, current position size
				// for Alice is 60 quantums long - 70 quantums = 10 quantums short.
				constants.Order_Carl_Num0_Id0_Clob0_Buy70_Price500000_GTB10,
			},

			secondOrders: []clobtypes.Order{
				// Currently, IOC/FOK replacement orders for orders that are partially filled are not allowed.
				// Because the order was partially filled, this results in a ErrInvalidReplacement error.
				// If IOC/FOK replacement orders were allowed, this buy RO order should succeed because the
				// current position size is 10 quantums short, and buying reduces a short position.
				constants.Order_Alice_Num1_Id0_Clob0_Buy110_Price50000_GTB21_FOK_RO,
			},
			secondOrdersErrors: []string{
				clobtypes.ErrInvalidReplacement.Error(),
			},

			expectedOrderFillAmounts: map[uint32]map[clobtypes.OrderId]uint64{
				2: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price500000_GTB20.OrderId: 70,
					constants.Order_Carl_Num0_Id0_Clob0_Buy70_Price500000_GTB10.OrderId:    70,
				},
				3: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price500000_GTB20.OrderId: 70,
					constants.Order_Carl_Num0_Id0_Clob0_Buy70_Price500000_GTB10.OrderId:    70,
				},
			},
		},
		`A regular order is partially filled. Replacement IOC RO order fails because it is immediate execution.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
				constants.Carl_Num0_100000USD,
			},

			firstOrders: []clobtypes.Order{
				// Regular order on the opposite side of the following replacement RO order.
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price500000_GTB20,
				// Partial match for the above order. 70 quantums are matched. Thus, current position size
				// for Alice is 60 quantums long - 70 quantums = 10 quantums short.
				constants.Order_Carl_Num0_Id0_Clob0_Buy70_Price500000_GTB10,
			},

			secondOrders: []clobtypes.Order{
				// Currently, IOC/FOK replacement orders for orders that are partially filled are not allowed.
				// Because the order was partially filled, this results in a ErrInvalidReplacement error.
				// If IOC/FOK replacement orders were allowed, this buy RO order should succeed because the
				// current position size is 10 quantums short, and buying reduces a short position.
				constants.Order_Alice_Num1_Id0_Clob0_Buy110_Price50000_GTB21_IOC_RO,
			},
			secondOrdersErrors: []string{
				clobtypes.ErrInvalidReplacement.Error(),
			},

			expectedOrderFillAmounts: map[uint32]map[clobtypes.OrderId]uint64{
				2: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price500000_GTB20.OrderId: 70,
					constants.Order_Carl_Num0_Id0_Clob0_Buy70_Price500000_GTB10.OrderId:    70,
				},
				3: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price500000_GTB20.OrderId: 70,
					constants.Order_Carl_Num0_Id0_Clob0_Buy70_Price500000_GTB10.OrderId:    70,
				},
			},
		},
		`Position size is long. A regular order is placed but not filled. Replacement Sell FOK RO
		 reduces current long position size, but fails because it is resized to a smaller size due to
		 current position size.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
				constants.Carl_Num0_100000USD,
			},

			firstOrders: []clobtypes.Order{
				// Regular order on the opposite side of the following replacement RO order.
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price51000_GTB20,
			},

			secondOrders: []clobtypes.Order{
				// Full match for the below order.
				constants.Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB10,
				// The original order being replaced has no partial fills. This sell RO order should succeed because the
				// current position size is 60 quantums long, and selling reduces a long position. However,
				// the RO property resizes the order to the current position size (60) and thus the FOK order
				// cannot be fully filled and errors out.
				constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_FOK_RO,
			},
			secondOrdersErrors: []string{
				"",
				clobtypes.ErrFokOrderCouldNotBeFullyFilled.Error(),
			},

			expectedOrderFillAmounts: map[uint32]map[clobtypes.OrderId]uint64{
				2: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price51000_GTB20.OrderId: 0,
				},
				3: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_FOK_RO.OrderId: 0,
					constants.Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB10.OrderId:          0,
				},
			},
		},
		`Position size is long. A regular order is placed but not filled. Replacement Sell FOK RO
		reduces current long position size, and succeeds because subaccount has enough position
		to not resize the order smaller.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(110), // 110 quantums of BTC long
						},
					},
				},
				constants.Carl_Num0_100000USD,
			},

			firstOrders: []clobtypes.Order{
				// Regular order on the opposite side of the following replacement RO order.
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price51000_GTB20,
			},

			secondOrders: []clobtypes.Order{
				// Full match for the below order.
				constants.Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB10,
				// The original order being replaced has no partial fills. This sell RO order should succeed because the
				// current position size is 110 quantums long, and selling reduces a long position. RO property of the
				// order does not resize the order and order succeeds for full amount.
				constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_FOK_RO,
			},
			secondOrdersErrors: []string{
				"",
				"",
			},

			expectedOrderFillAmounts: map[uint32]map[clobtypes.OrderId]uint64{
				2: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price51000_GTB20.OrderId: 0,
				},
				3: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_FOK_RO.OrderId: 110,
					constants.Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB10.OrderId:          110,
				},
			},
		},
		`Position size is long. A regular order is placed but not filled. Replacement Sell IOC RO
		reduces current long position size, but is resized to a smaller size due to current position size
		and partially filled.`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_100_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(60), // 60 quantums of BTC long
						},
					},
				},
				constants.Carl_Num0_100000USD,
			},

			firstOrders: []clobtypes.Order{
				// Regular order on the opposite side of the following replacement RO order.
				constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price51000_GTB20,
			},

			secondOrders: []clobtypes.Order{
				// Full match for the below order.
				constants.Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB10,
				// The original order being replaced has no partial fills. This sell RO order should succeed because the
				// current position size is 60 quantums long, and selling reduces a long position.
				// The RO property resizes the order to the current position size (60) and thus the IOC order
				// is partially filled.
				constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_IOC_RO,
			},
			secondOrdersErrors: []string{
				"",
				"",
			},

			expectedOrderFillAmounts: map[uint32]map[clobtypes.OrderId]uint64{
				2: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell100_Price51000_GTB20.OrderId: 0,
				},
				3: {
					constants.Order_Alice_Num1_Id0_Clob0_Sell110_Price50000_GTB21_IOC_RO.OrderId: 60,
					constants.Order_Carl_Num0_Id0_Clob0_Buy110_Price50000_GTB10.OrderId:          60,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).
				WithGenesisDocFn(func() (genesis types.GenesisDoc) {
					genesis = testapp.DefaultGenesis()
					if len(tc.subaccounts) > 0 {
						testapp.UpdateGenesisDocWithAppStateForModule(
							&genesis,
							func(genesisState *satypes.GenesisState) {
								genesisState.Subaccounts = tc.subaccounts
							},
						)
					}
					return genesis
				}).
				WithCrashingAppCheckTxNonDeterminismChecksEnabled(false).
				Build()

			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)

			require.NoError(t, conversionErr)

			tApp.App.RatelimitKeeper.SetSDAIPrice(tApp.App.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.App.RatelimitKeeper.SetAssetYieldIndex(tApp.App.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.CrashingApp.RatelimitKeeper.SetSDAIPrice(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.CrashingApp.RatelimitKeeper.SetAssetYieldIndex(tApp.CrashingApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.NoCheckTxApp.RatelimitKeeper.SetSDAIPrice(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.NoCheckTxApp.RatelimitKeeper.SetAssetYieldIndex(tApp.NoCheckTxApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			tApp.ParallelApp.RatelimitKeeper.SetSDAIPrice(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), rate)
			tApp.ParallelApp.RatelimitKeeper.SetAssetYieldIndex(tApp.ParallelApp.NewUncachedContext(false, tmproto.Header{}), big.NewRat(1, 1))

			ctx := tApp.InitChain()

			// place first set of orders.
			for _, order := range tc.firstOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)
					require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			// Advance the block to persist matches.
			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// validate order fill amounts.
			if orderMap, exists := tc.expectedOrderFillAmounts[2]; exists {
				for orderId, expectedFillAmount := range orderMap {
					exists, actualFillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
					if expectedFillAmount == 0 {
						require.False(t, exists)
					} else {
						require.True(t, exists)
						require.Equal(t, expectedFillAmount, actualFillAmount.ToUint64())
					}
				}
			}

			// place second set of orders.
			for idx, order := range tc.secondOrders {
				for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(
					ctx,
					tApp.App,
					*clobtypes.NewMsgPlaceOrder(order),
				) {
					resp := tApp.CheckTx(checkTx)

					if tc.secondOrdersErrors[idx] == "" {
						require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
					} else {
						require.Conditionf(t, resp.IsErr, "Expected CheckTx to error. Response: %+v", resp)
						require.Contains(
							t,
							resp.Log,
							tc.secondOrdersErrors[idx],
						)
					}
				}
			}

			// Advance the block to persist matches.
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

			// validate order fill amounts.
			if orderMap, exists := tc.expectedOrderFillAmounts[3]; exists {
				for orderId, expectedFillAmount := range orderMap {
					exists, actualFillAmount, _ := tApp.App.ClobKeeper.GetOrderFillAmount(ctx, orderId)
					if expectedFillAmount == 0 {
						require.False(t, exists)
					} else {
						require.True(t, exists)
						require.Equal(t, expectedFillAmount, actualFillAmount.ToUint64())
					}
				}
			}
		})
	}
}
