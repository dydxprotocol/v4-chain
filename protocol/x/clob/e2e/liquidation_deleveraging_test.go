package clob_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"

	"github.com/cometbft/cometbft/types"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestLiquidationConfig(t *testing.T) {
	tests := map[string]struct {
		// State.
		subaccounts                   []satypes.Subaccount
		marketIdToOraclePriceOverride map[uint32]uint64

		// Parameters.
		placedMatchableOrders     []clobtypes.MatchableOrder
		liquidatableSubaccountIds []satypes.SubaccountId

		// Configuration.
		liquidationConfig clobtypes.LiquidationsConfig
		liquidityTiers    []perptypes.LiquidityTier
		perpetuals        []perptypes.Perpetual
		clobPairs         []clobtypes.ClobPair

		// Expectations.
		expectedSubaccounts []satypes.Subaccount
	}{
		`Liquidating short respects position block limit - MinPositionNotionalLiquidated`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10, // Order at $50,000
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: clobtypes.PositionBlockLimits{
					// 1% of $50,000 is $500 so $500 worth of BTC should get liquidated.
					// However, this is smaller than the minimum position notional liquidated of $100,000,
					// so the entire position should get liquidated.
					MinPositionNotionalLiquidated:   100_000_000_000, // $100,000
					MaxPositionPortionLiquidatedPpm: 10_000,          // 1%
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_499_000_000-50_000_000_000-250_000_000),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(100_000_000_000), // $100,000
						),
					},
				},
			},
		},
		`Liquidatiing long respects position block limit - MinPositionNotionalLiquidated`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_49501USD_Short,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10, // Order at $50,000
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Dave_Num0},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: clobtypes.PositionBlockLimits{
					// 1% of $50,000 is $500 so $500 worth of BTC should get liquidated.
					// However, this is smaller than the minimum position notional liquidated of $100,000,
					// so the entire position should get liquidated.
					MinPositionNotionalLiquidated:   100_000_000_000, // $100,000
					MaxPositionPortionLiquidatedPpm: 10_000,          // 1%
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(100_000_000_000-50_000_000_000),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(-49_501_000_000+50_000_000_000-250_000_000),
						),
					},
				},
			},
		},
		`Liquidatiing short respects position block limit - MaxPositionPortionLiquidatedPpm`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10, // Order at $50,000
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: clobtypes.PositionBlockLimits{
					// 10% of $50,000 is $5,000 so $5,000 worth of BTC should get liquidated.
					MinPositionNotionalLiquidated:   100_000_000, // $1,000
					MaxPositionPortionLiquidatedPpm: 100_000,     // 10%
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_499_000_000-5_000_000_000-25_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-90_000_000), // -0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(55_000_000_000), // $55,000
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(90_000_000), // 0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		`Liquidatiing long respects position block limit - MaxPositionPortionLiquidatedPpm`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_49501USD_Short,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10, // Order at $50,000
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Dave_Num0},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: clobtypes.PositionBlockLimits{
					// 10% of $50,000 is $5,000 so $5,000 worth of BTC should get liquidated.
					MinPositionNotionalLiquidated:   100_000_000, // $1,000
					MaxPositionPortionLiquidatedPpm: 100_000,     // 10%
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(100_000_000_000-5_000_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-90_000_000), // -0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(-49_501_000_000+5_000_000_000-25_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(90_000_000), // 0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		`Liquidating short respects subaccount block limit - MaxNotionalLiquidated`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				&constants.Order_Dave_Num0_Id2_Clob0_Sell1BTC_Price49500_GTB10, // Order at $49,500
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_Default,
				SubaccountBlockLimits: clobtypes.SubaccountBlockLimits{
					// Subaccount may only liquidate $5,000 per block.
					MaxNotionalLiquidated:    5_000_000_000,
					MaxQuantumsInsuranceLost: 100_000_000_000_000,
				},
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_499_000_000-4_950_000_000-24_750_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-90_000_000), // -0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(54_950_000_000), // $54,950
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(90_000_000), // 0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		`Liquidating long respects subaccount block limit - MaxNotionalLiquidated`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_49501USD_Short,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Maker order at $50,500.
				// This maker order is specifically chosen to be above the oracle price, to ensure that
				// block limits use the notional value of the position (oracle price),
				// and not the notional liquidated (match price)
				&constants.Order_Carl_Num0_Id2_Clob0_Buy1BTC_Price50500_GTB10, // Order at $50,500
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Dave_Num0},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_Default,
				SubaccountBlockLimits: clobtypes.SubaccountBlockLimits{
					// Subaccount may only liquidate $5,000 per block.
					MaxNotionalLiquidated:    5_000_000_000,
					MaxQuantumsInsuranceLost: 100_000_000_000_000,
				},
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(100_000_000_000-5_050_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-90_000_000), // -0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(-49_501_000_000+5_050_000_000-25_250_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(90_000_000), // 0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		`Liquidating short respects subaccount block limit - MaxQuantumsInsuranceLost`: {
			subaccounts: []satypes.Subaccount{
				// Carl_Num0 is irrelevant to the test, but is used to seed the insurance fund.
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Carl_Num1_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// This order is irrelevant to the test, but is used to seed the insurance fund.
				&constants.Order_Dave_Num0_Id2_Clob0_Sell1BTC_Price49500_GTB10, // Order at $49,500

				// Bankruptcy price is $50,499, and closing at $50,500 would require $1 from the insurance fund.
				// First order would transfer $0.1 from the insurance fund and would succeed.
				// Second order would require $0.9 from the insurance fund and would fail since subaccounts
				// may only lose $0.5 per block.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell01BTC_Price50500_GTB10, // Order at $50,500
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,  // Order at $50,500
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num1, constants.Carl_Num0},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Max_Smmr,
				PositionBlockLimits:  constants.PositionBlockLimits_Default,
				SubaccountBlockLimits: clobtypes.SubaccountBlockLimits{
					// Subaccount may only lose $0.5 per block.
					MaxNotionalLiquidated:    100_000_000_000_000,
					MaxQuantumsInsuranceLost: 500_000,
				},
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num1,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							// $0.1 from insurance fund
							big.NewInt(50_499_000_000-5_050_000_000+100_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-90_000_000), // -0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+49_500_000_000+5_050_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-10_000_000), // -0.1 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		`Liquidating long respects subaccount block limit - MaxQuantumsInsuranceLost`: {
			subaccounts: []satypes.Subaccount{
				// Carl_Num0 is irrelevant to the test, but is used to seed the insurance fund.
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_49501USD_Short,
				constants.Dave_Num1_1BTC_Long_49501USD_Short,
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// This order is irrelevant to the test, but is used to seed the insurance fund.
				&constants.Order_Carl_Num0_Id2_Clob0_Buy1BTC_Price50500_GTB10, // Order at $50,500

				// Bankruptcy price is $49,501, and closing at $49,500 would require $1 from the insurance fund.
				// First order would transfer $0.1 from the insurance fund and would succeed.
				// Second order would require $0.9 from the insurance fund and would fail since subaccounts
				// may only lose $0.5 per block.
				&constants.Order_Carl_Num0_Id1_Clob0_Buy01BTC_Price49500_GTB10, // Order at $49,500
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,  // Order at $49,500
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Dave_Num0, constants.Dave_Num1},
			liquidationConfig: clobtypes.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Max_Smmr,
				PositionBlockLimits:  constants.PositionBlockLimits_Default,
				SubaccountBlockLimits: clobtypes.SubaccountBlockLimits{
					// Subaccount may only lose $0.5 per block.
					MaxNotionalLiquidated:    100_000_000_000_000,
					MaxQuantumsInsuranceLost: 500_000,
				},
			},

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(100_000_000_000-50_500_000_000-4_950_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(10_000_000), // 0.1 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(-49_501_000_000+4_950_000_000+100_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(90_000_000), // 0.9 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
						// Set oracle prices in the genesis.
						pricesGenesis := constants.TestPricesGenesisState

						// Make a copy of the MarketPrices slice to avoid modifying by reference.
						marketPricesCopy := make([]prices.MarketPrice, len(pricesGenesis.MarketPrices))
						copy(marketPricesCopy, pricesGenesis.MarketPrices)

						for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {
							exponent, exists := constants.TestMarketIdsToExponents[marketId]
							require.True(t, exists)

							marketPricesCopy[marketId] = prices.MarketPrice{
								Id:       marketId,
								Price:    oraclePrice,
								Exponent: exponent,
							}
						}

						pricesGenesis.MarketPrices = marketPricesCopy
						*genesisState = pricesGenesis
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = tc.liquidityTiers
						genesisState.Perpetuals = tc.perpetuals
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
						genesisState.ClobPairs = tc.clobPairs
						genesisState.LiquidationsConfig = tc.liquidationConfig
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).Build()

			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Create all existing orders.
			existingOrderMsgs := make([]clobtypes.MsgPlaceOrder, len(tc.placedMatchableOrders))
			for i, matchableOrder := range tc.placedMatchableOrders {
				existingOrderMsgs[i] = clobtypes.MsgPlaceOrder{Order: matchableOrder.MustGetOrder()}
			}
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, existingOrderMsgs...) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			_, err := tApp.App.Server.LiquidateSubaccounts(ctx, &api.LiquidateSubaccountsRequest{
				LiquidatableSubaccountIds: tc.liquidatableSubaccountIds,
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

func TestPlacePerpetualLiquidation_Deleveraging(t *testing.T) {
	tests := map[string]struct {
		// State.
		subaccounts                   []satypes.Subaccount
		marketIdToOraclePriceOverride map[uint32]uint64

		// Parameters.
		placedMatchableOrders     []clobtypes.MatchableOrder
		liquidatableSubaccountIds []satypes.SubaccountId

		// Configuration.
		liquidationConfig clobtypes.LiquidationsConfig
		liquidityTiers    []perptypes.LiquidityTier
		perpetuals        []perptypes.Perpetual
		clobPairs         []clobtypes.ClobPair

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

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_499_000_000-50_000_000_000-250_000_000),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(100_000_000_000), // $100,000
						),
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

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_499_000_000-12_500_000_000-62_500_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-75_000_000), // -0.75 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+12_500_000_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(75_000_000), // 0.75 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
			},
		},
		`Can place a liquidation order that is unfilled and full position size is deleveraged`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},

			placedMatchableOrders: []clobtypes.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $50,499, and closing at $50,500
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{constants.Carl_Num0},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+50_499_000_000),
						),
					},
				},
			},
		},
		`Can place a liquidation order that is partially-filled and deleveraging is skipped`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
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

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_499_000_000-(50_498_000_000/4)-250_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-75_000_000), // -0.75 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+(50_498_000_000/4)),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(75_000_000), // 0.75 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
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

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

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

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails for remaining amount.
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(49_999_000_000-12_499_750_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-75_000_000), // -0.75 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
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

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(49_999_000_000-24_999_500_000),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							// Deleveraging fails for remaining amount.
							big.NewInt(-50_000_000), // -0.5 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
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
		`Deleveraging takes precedence - can place a liquidation order that would fail due to exceeding
		subaccount limit and full position size is deleveraged`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
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

			liquidityTiers: constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+50_499_000_000),
						),
					},
				},
			},
		},
		`Deleveraging occurs at bankruptcy price for negative TNC subaccount with open position in final settlement market`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},
			// Account should be deleveraged regardless of whether or not the liquidations engine returns this subaccount
			// in the list of liquidatable subaccounts. Pass empty list to confirm this.
			liquidatableSubaccountIds: []satypes.SubaccountId{},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,
			liquidityTiers:            constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc_Final_Settlement},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+50_499_000_000),
						),
					},
				},
			},
		},
		`Deleveraging occurs at oracle price for non-negative TNC subaccounts 
			with open positions in final settlement market`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_100000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			liquidatableSubaccountIds: []satypes.SubaccountId{},
			liquidationConfig:         constants.LiquidationsConfig_FillablePrice_Max_Smmr,
			liquidityTiers:            constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1,
			},
			clobPairs: []clobtypes.ClobPair{constants.ClobPair_Btc_Final_Settlement},

			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(100_000_000_000-50_000_000_000),
						),
					},
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+50_000_000_000),
						),
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
						// Set oracle prices in the genesis.
						pricesGenesis := constants.TestPricesGenesisState

						// Make a copy of the MarketPrices slice to avoid modifying by reference.
						marketPricesCopy := make([]prices.MarketPrice, len(pricesGenesis.MarketPrices))
						copy(marketPricesCopy, pricesGenesis.MarketPrices)

						for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {
							exponent, exists := constants.TestMarketIdsToExponents[marketId]
							require.True(t, exists)

							marketPricesCopy[marketId] = prices.MarketPrice{
								Id:       marketId,
								Price:    oraclePrice,
								Exponent: exponent,
							}
						}

						pricesGenesis.MarketPrices = marketPricesCopy
						*genesisState = pricesGenesis
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = tc.liquidityTiers
						genesisState.Perpetuals = tc.perpetuals
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
						genesisState.ClobPairs = tc.clobPairs
						genesisState.LiquidationsConfig = tc.liquidationConfig
						genesisState.EquityTierLimitConfig = clobtypes.EquityTierLimitConfiguration{}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).Build()

			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Create all existing orders.
			existingOrderMsgs := make([]clobtypes.MsgPlaceOrder, len(tc.placedMatchableOrders))
			for i, matchableOrder := range tc.placedMatchableOrders {
				existingOrderMsgs[i] = clobtypes.MsgPlaceOrder{Order: matchableOrder.MustGetOrder()}
			}
			for _, checkTx := range testapp.MustMakeCheckTxsWithClobMsg(ctx, tApp.App, existingOrderMsgs...) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			_, err := tApp.App.Server.LiquidateSubaccounts(ctx, &api.LiquidateSubaccountsRequest{
				LiquidatableSubaccountIds:  tc.liquidatableSubaccountIds,
				SubaccountOpenPositionInfo: clobtest.GetOpenPositionsFromSubaccounts(tc.subaccounts),
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
