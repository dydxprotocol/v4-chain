package keeper_test

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	big_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/big"
	clobtest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/clob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	perptest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/perpetuals"
	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/heap"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/memclob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	feetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlacePerpetualLiquidation(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs     []types.ClobPair
		feeParams feetypes.PerpetualFeeParams

		existingOrders []types.Order

		// Parameters.
		order types.LiquidationOrder

		// Expectations.
		expectedPlacedOrders  []*types.MsgPlaceOrder
		expectedMatchedOrders []*types.ClobMatch
		// Expected remaining OI after test.
		// The test initializes each perp with default open interest of 1 full coin.
		expectedOpenInterests map[uint32]*big.Int
	}{
		`Can place a liquidation that doesn't match any maker orders`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			order: constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,

			expectedPlacedOrders:  []*types.MsgPlaceOrder{},
			expectedMatchedOrders: []*types.ClobMatch{},
			expectedOpenInterests: map[uint32]*big.Int{
				constants.BtcUsd_SmallMarginRequirement.Params.Id: big.NewInt(100_000_000), // unchanged
			},
		},
		`Can place a liquidation that matches maker orders`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			},

			order: constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,

			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   100_000_000,
						Liquidated:  constants.Dave_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			expectedOpenInterests: map[uint32]*big.Int{
				constants.BtcUsd_SmallMarginRequirement.Params.Id: new(big.Int), // fully liquidated
			},
		},
		`Can place a liquidation that matches maker orders and removes undercollateralized ones`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				// Note this order will be removed when matching.
				constants.Order_Carl_Num1_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			},

			order: constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,

			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   100_000_000,
						Liquidated:  constants.Dave_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			expectedOpenInterests: map[uint32]*big.Int{
				constants.BtcUsd_SmallMarginRequirement.Params.Id: new(big.Int), // fully liquidated
			},
		},
		`Can place a liquidation that matches maker orders with maker rebates and empty fee collector`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParamsMakerRebate,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			},

			order: constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,

			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   100_000_000,
						Liquidated:  constants.Dave_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			expectedOpenInterests: map[uint32]*big.Int{
				constants.BtcUsd_SmallMarginRequirement.Params.Id: new(big.Int), // fully liquidated
			},
		},
		`Can place a liquidation that matches maker orders with maker rebates`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParamsMakerRebate,
			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
			},

			order: constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,

			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   100_000_000,
						Liquidated:  constants.Dave_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			expectedOpenInterests: map[uint32]*big.Int{
				constants.BtcUsd_SmallMarginRequirement.Params.Id: new(big.Int), // fully liquidated
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				satypes.ModuleAddress,
				authtypes.NewModuleAddress(authtypes.FeeCollectorName),
				mock.Anything,
			).Return(nil)
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				authtypes.NewModuleAddress(satypes.ModuleName),
				perptypes.InsuranceFundModuleAddress,
				mock.Anything,
			).Return(nil)
			// Fee collector does not have any funds.
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				authtypes.NewModuleAddress(authtypes.FeeCollectorName),
				satypes.ModuleAddress,
				mock.Anything,
			).Return(sdkerrors.ErrInsufficientFunds)
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				authtypes.NewModuleAddress(satypes.LiquidityFeeModuleAddress),
				mock.Anything,
			).Return(nil)
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			// Give the insurance fund a 1M TDai balance.
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				perptypes.InsuranceFundModuleAddress,
				constants.TDai.Denom,
			).Return(
				sdk.NewCoin(
					constants.TDai.Denom,
					sdkmath.NewIntFromBigInt(big.NewInt(1_000_000_000_000)),
				),
			)
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			ctx := ks.Ctx.WithIsCheckTx(true)
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, tc.feeParams))

			// Set up tDAI asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobs {
				_, err = ks.ClobKeeper.CreatePerpetualClobPair(
					ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			// Initialize the liquidations config.
			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default),
			)

			// Create all existing orders.
			for _, order := range tc.existingOrders {
				_, _, err := ks.ClobKeeper.PlaceShortTermOrder(ctx, &types.MsgPlaceOrder{Order: order})
				require.NoError(t, err)
			}

			// Run the test.
			_, _, err = ks.ClobKeeper.PlacePerpetualLiquidation(ctx, tc.order)
			require.NoError(t, err)

			for _, perp := range tc.perpetuals {
				if expectedOI, exists := tc.expectedOpenInterests[perp.Params.Id]; exists {
					gotPerp, err := ks.PerpetualsKeeper.GetPerpetual(ks.Ctx, perp.Params.Id)
					require.NoError(t, err)
					require.Zero(t,
						expectedOI.Cmp(gotPerp.OpenInterest.BigInt()),
						"expected open interest %s, got %s",
						expectedOI.String(),
						gotPerp.OpenInterest.String(),
					)
				}
			}

			// Verify test expectations.
			// TODO(DEC-1979): Refactor these tests to support the operations queue refactor.
			// placedOrders, matchedOrders := memClob.GetPendingFills(ctx)

			// require.Equal(t, tc.expectedPlacedOrders, placedOrders, "Placed orders lists are not equal")
			// require.Equal(t, tc.expectedMatchedOrders, matchedOrders, "Matched orders lists are not equal")
		})
	}
}

func TestPlacePerpetualLiquidation_validateLiquidationAgainstClobPairStatus(t *testing.T) {
	tests := map[string]struct {
		status types.ClobPair_Status

		expectedError error
	}{
		"Cannot liquidate in initializing state": {
			status: types.ClobPair_STATUS_INITIALIZING,

			expectedError: types.ErrLiquidationConflictsWithClobPairStatus,
		},
		"Can liquidate in active state": {
			status: types.ClobPair_STATUS_ACTIVE,
		},
		"Cannot liquidate in final settlement state": {
			status: types.ClobPair_STATUS_FINAL_SETTLEMENT,

			expectedError: types.ErrLiquidationConflictsWithClobPairStatus,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))
			ctx := ks.Ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			perpetuals := []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			}
			for _, p := range perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			clobPair := constants.ClobPair_Btc
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ctx,
				clobPair.Id,
				clobtest.MustPerpetualId(clobPair),
				satypes.BaseQuantums(clobPair.StepBaseQuantums),
				clobPair.QuantumConversionExponent,
				clobPair.SubticksPerTick,
				tc.status,
			)
			require.NoError(t, err)

			_, _, err = ks.ClobKeeper.PlacePerpetualLiquidation(
				ctx,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
			)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			}
		})
	}
}

func TestPlacePerpetualLiquidation_PreexistingLiquidation(t *testing.T) {
	tests := map[string]struct {
		// State.
		subaccounts         []satypes.Subaccount
		setupMockBankKeeper func(m *mocks.BankKeeper)

		// Parameters.
		liquidationConfig     types.LiquidationsConfig
		placedMatchableOrders []types.MatchableOrder
		order                 types.LiquidationOrder

		// Expectations.
		panics                            bool
		expectedError                     error
		expectedFilledSize                satypes.BaseQuantums
		expectedOrderStatus               types.OrderStatus
		expectedPlacedOrders              []*types.MsgPlaceOrder
		expectedMatchedOrders             []*types.ClobMatch
		expectedSubaccountLiquidationInfo map[satypes.SubaccountId]types.SubaccountLiquidationInfo
		expectedLiquidationDeltaPerBlock  map[uint32]*big.Int
	}{
		`PlacePerpetualLiquidation succeeds with pre-existing liquidations in the block`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(54_999_000_000), // $54,999
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
						},
						{
							PerpetualId: 1,
							Quantums:    dtypes.NewInt(-1_000_000_000), // -1 ETH
						},
					},
					AssetYieldIndex: "1/1",
				},
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000,
				&constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000,
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,

			expectedOrderStatus: types.Success,
			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000,
				},
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Eth.Id,
						IsBuy:       true,
						TotalSize:   1_000_000_000,
						Liquidated:  constants.Carl_Num0,
						PerpetualId: constants.ClobPair_Eth.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   1_000_000_000,
							},
						},
					},
				),
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       true,
						TotalSize:   100_000_000,
						Liquidated:  constants.Carl_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{1, 0},
				},
				constants.Dave_Num0: {},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(-265000000),
				1: big.NewInt(-265000000),
			},
		},
		`PlacePerpetualLiquidation considers pre-existing liquidations and stops before exceeding
		max insurance fund lost per block`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(53_000_000_000), // $53,000
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
						},
						{
							PerpetualId: 1,
							Quantums:    dtypes.NewInt(-1_000_000_000), // -1 ETH
						},
					},
					AssetYieldIndex: "1/1",
				},
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			liquidationConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm:             5_000,
				ValidatorFeePpm:                 0,
				LiquidityFeePpm:                 0,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(50_000_000),
			},
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id4_Clob1_Sell1ETH_Price3030,
				&constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3030,
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,

			// Only matches one order since matching both orders would exceed `MaxQuantumsInsuranceLost`.
			expectedOrderStatus: types.LiquidationExceededSubaccountMaxInsuranceLost,
			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id4_Clob1_Sell1ETH_Price3030,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Eth.Id,
						IsBuy:       true,
						TotalSize:   1_000_000_000,
						Liquidated:  constants.Carl_Num0,
						PerpetualId: constants.ClobPair_Eth.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   1_000_000_000,
							},
						},
					},
				),
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{1, 0},
				},
				constants.Dave_Num0: {},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(30_000_000),
				1: big.NewInt(30_000_000),
			},
		},
		`PlacePerpetualLiquidation matches some order and stops before exceeding max insurance lost per block`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			liquidationConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm:             5_000,
				ValidatorFeePpm:                 0,
				LiquidityFeePpm:                 0,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(500_000),
			},
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
			},
			// Overall insurance lost when liquidating at $50,500 is $1.
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,

			// Only matches one order since matching both orders would exceed `MaxQuantumsInsuranceLost`.
			expectedOrderStatus: types.LiquidationExceededSubaccountMaxInsuranceLost,
			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       true,
						TotalSize:   100_000_000,
						Liquidated:  constants.Carl_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   25_000_000,
							},
						},
					},
				),
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
				constants.Dave_Num0: {},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(250_000),
			},
		},
		`Liquidation buy order does not generate a match when deleveraging is required`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
					constants.TDai.Denom,
				).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("TDAI", sdkmath.NewIntFromUint64(0))) // Insurance fund is empty.
			},

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Expected insurance fund delta is $-1.

			// Does not generate a match since insurance fund does not have enough to cover the losses.
			expectedOrderStatus:   types.LiquidationRequiresDeleveraging,
			expectedPlacedOrders:  []*types.MsgPlaceOrder{},
			expectedMatchedOrders: []*types.ClobMatch{},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
				constants.Dave_Num0: {},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
		},
		`Liquidation sell order does not generate a match when deleveraging is required`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_49501USD_Short,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
					constants.TDai.Denom,
				).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("TDAI", sdkmath.NewIntFromUint64(0))) // Insurance fund is empty.
			},

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
			},
			order: constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price49500, // Expected insurance fund delta is $-1.

			// Does not generate a match since insurance fund does not have enough to cover the losses.
			expectedOrderStatus:   types.LiquidationRequiresDeleveraging,
			expectedPlacedOrders:  []*types.MsgPlaceOrder{},
			expectedMatchedOrders: []*types.ClobMatch{},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
		},
		`Liquidation buy order matches with some orders and stops when insurance fund is empty`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
					constants.TDai.Denom,
				).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(
					// Insurance fund has $0.99 initially.
					sdk.NewCoin("TDAI", sdkmath.NewIntFromUint64(990_000)),
				).Once()
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(
					// Insurance fund has $0.74 after covering the loss of the first match.
					sdk.NewCoin("TDAI", sdkmath.NewIntFromUint64(740_000)),
				).Twice()
			},

			liquidationConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm:             5_000,
				ValidatorFeePpm:                 200_000,
				LiquidityFeePpm:                 800_000,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
			},
			// Overall insurance fund delta when liquidating at $50,500 is -$1.
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,

			// Matches the first order since insurance fund balance has enough to cover the losses (-$0.25).
			// Does not match the second order since insurance fund delta is -$0.75 and insurance fund balance
			// is $0.74 which is not enough to cover the loss, and therefore deleveraging is required.
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       true,
						TotalSize:   100_000_000,
						Liquidated:  constants.Carl_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   25_000_000,
							},
						},
					},
				),
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
				constants.Dave_Num0: {},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(250_000),
			},
		},
		`Liquidation sell order matches with some orders and stops when deleveraging is required`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_49501USD_Short,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
					constants.TDai.Denom,
				).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(
					// Insurance fund has $0.99 initially.
					sdk.NewCoin("TDAI", sdkmath.NewIntFromUint64(990_000)),
				).Once()
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(
					// Insurance fund has $0.74 after covering the loss of the first match.
					sdk.NewCoin("TDAI", sdkmath.NewIntFromUint64(740_000)),
				).Once()
			},

			liquidationConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm:             5_000,
				ValidatorFeePpm:                 200_000,
				LiquidityFeePpm:                 800_000,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price49500,
				&constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
			},
			// Overall insurance fund delta when liquidating at $50,500 is -$1.
			order: constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price49500,

			// Matches the first order since insurance fund balance has enough to cover the losses (-$0.25).
			// Does not match the second order since insurance fund delta is -$0.75 and insurance fund balance
			// is $0.74 which is not enough to cover the loss, and therefore deleveraging is required.
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price49500,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   100_000_000,
						Liquidated:  constants.Dave_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   25_000_000,
							},
						},
					},
				),
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(250_000),
			},
		},
		`PlacePerpetualLiquidation panics when trying to liquidate the same perpetual in a block`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(54_999_000_000), // $54,999
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
						},
						{
							PerpetualId: 1,
							Quantums:    dtypes.NewInt(-2_000_000_000), // -2 ETH
						},
					},
					AssetYieldIndex: "1/1",
				},
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000,
				&constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000,
				&constants.Order_Dave_Num0_Id4_Clob1_Sell1ETH_Price3000,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000,

			expectedError: errorsmod.Wrapf(
				types.ErrSubaccountHasLiquidatedPerpetual,
				"Subaccount %v and perpetual %v have already been liquidated within the last block",
				constants.Carl_Num0,
				1,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			memclob := memclob.NewMemClobPriceTimePriority(false)

			bankKeeper := &mocks.BankKeeper{}
			if tc.setupMockBankKeeper != nil {
				tc.setupMockBankKeeper(bankKeeper)
			} else {
				bankKeeper.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bankKeeper.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)
				bankKeeper.On(
					"GetBalance",
					mock.Anything,
					authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
					constants.TDai.Denom,
				).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))
				bankKeeper.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("TDAI", sdkmath.NewIntFromUint64(math.MaxUint64)))

			}

			mockIndexerEventManager := &mocks.IndexerEventManager{}
			mockIndexerEventManager.On("Enabled").Return(false)
			ks := keepertest.NewClobKeepersTestContext(t, memclob, bankKeeper, mockIndexerEventManager)
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			ctx := ks.Ctx.WithIsCheckTx(true)

			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParams))

			// Set up tDAI asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			testPerps := []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_100PercentMarginRequirement,
			}
			for _, perpetual := range testPerps {
				_, err = ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					perpetual.Params.Id,
					perpetual.Params.Ticker,
					perpetual.Params.MarketId,
					perpetual.Params.AtomicResolution,
					perpetual.Params.DefaultFundingPpm,
					perpetual.Params.LiquidityTier,
					perpetual.Params.MarketType,
					perpetual.Params.DangerIndexPpm,
					perpetual.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					perpetual.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				testPerps,
			)

			for _, s := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, s)
			}
			// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
			// the indexer event manager to expect these events.
			mockIndexerEventManager.On("AddTxnEvent",
				ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						0,
						0,
						constants.BtcUsd_100PercentMarginRequirement.Params.Ticker,
						constants.BtcUsd_100PercentMarginRequirement.Params.MarketId,
						constants.ClobPair_Btc.Status,
						constants.ClobPair_Btc.QuantumConversionExponent,
						constants.BtcUsd_100PercentMarginRequirement.Params.AtomicResolution,
						constants.ClobPair_Btc.SubticksPerTick,
						constants.ClobPair_Btc.StepBaseQuantums,
						constants.BtcUsd_100PercentMarginRequirement.Params.LiquidityTier,
						constants.BtcUsd_100PercentMarginRequirement.Params.MarketType,
						constants.BtcUsd_100PercentMarginRequirement.Params.DangerIndexPpm,
						fmt.Sprintf("%d", constants.BtcUsd_100PercentMarginRequirement.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock),
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)
			mockIndexerEventManager.On("AddTxnEvent",
				ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						1,
						1,
						constants.EthUsd_100PercentMarginRequirement.Params.Ticker,
						constants.EthUsd_100PercentMarginRequirement.Params.MarketId,
						constants.ClobPair_Eth.Status,
						constants.ClobPair_Eth.QuantumConversionExponent,
						constants.EthUsd_100PercentMarginRequirement.Params.AtomicResolution,
						constants.ClobPair_Eth.SubticksPerTick,
						constants.ClobPair_Eth.StepBaseQuantums,
						constants.EthUsd_100PercentMarginRequirement.Params.LiquidityTier,
						constants.EthUsd_100PercentMarginRequirement.Params.MarketType,
						constants.EthUsd_100PercentMarginRequirement.Params.DangerIndexPpm,
						fmt.Sprintf("%d", constants.EthUsd_100PercentMarginRequirement.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock),
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ctx,
				constants.ClobPair_Eth.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Eth),
				satypes.BaseQuantums(constants.ClobPair_Eth.StepBaseQuantums),
				constants.ClobPair_Eth.QuantumConversionExponent,
				constants.ClobPair_Eth.SubticksPerTick,
				constants.ClobPair_Eth.Status,
			)
			require.NoError(t, err)

			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, tc.liquidationConfig),
			)

			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Timestamp: time.Unix(5, 0),
			})

			// Place all existing orders on the orderbook.
			for _, matchableOrder := range tc.placedMatchableOrders {
				// If the order is a liquidation order, place the liquidation.
				// Else, assume it's a regular order and place it.
				if liquidationOrder, ok := matchableOrder.(*types.LiquidationOrder); ok {
					_, _, err := ks.ClobKeeper.PlacePerpetualLiquidation(
						ctx,
						*liquidationOrder,
					)
					require.NoError(t, err)
				} else {
					order := matchableOrder.MustGetOrder()
					_, _, err := ks.ClobKeeper.PlaceShortTermOrder(ctx, &types.MsgPlaceOrder{Order: order.MustGetOrder()})
					require.NoError(t, err)
				}
			}

			// Run the test case and verify expectations.
			if tc.expectedError != nil {
				require.PanicsWithError(
					t,
					tc.expectedError.Error(),
					func() {
						_, _, _ = ks.ClobKeeper.PlacePerpetualLiquidation(ctx, tc.order)
					},
				)
			} else {
				_, orderStatus, err := ks.ClobKeeper.PlacePerpetualLiquidation(ctx, tc.order)
				require.NoError(t, err)
				require.Equal(t, tc.expectedOrderStatus, orderStatus)

				for subaccountId, liquidationInfo := range tc.expectedSubaccountLiquidationInfo {
					require.Equal(
						t,
						liquidationInfo,
						ks.ClobKeeper.GetSubaccountLiquidationInfo(ctx, subaccountId),
					)
				}

				for perpetualId, expectedLiquidationDeltaPerBlock := range tc.expectedLiquidationDeltaPerBlock {
					liquidationDeltaPerBlock, err := ks.ClobKeeper.GetCumulativeInsuranceFundDelta(ctx, perpetualId)
					require.NoError(t, err)
					require.Equal(
						t,
						expectedLiquidationDeltaPerBlock,
						liquidationDeltaPerBlock,
					)
				}

				// Verify test expectations.
				// TODO(DEC-1979): Refactor these tests to support the operations queue refactor.
				// placedOrders, matchedOrders := memclob.GetPendingFills(ctx)

				// require.Equal(t, tc.expectedPlacedOrders, placedOrders, "Placed orders lists are not equal")
				// require.Equal(t, tc.expectedMatchedOrders, matchedOrders, "Matched orders lists are not equal")
			}
		})
	}
}

func TestGetFillablePrice(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		perpetualId uint32

		// Perpetual state.
		perpetuals []perptypes.Perpetual

		// Subaccount state.
		assetPositions     []*satypes.AssetPosition
		perpetualPositions []*satypes.PerpetualPosition

		// Liquidation config.
		liquidationConfig *types.LiquidationsConfig

		// Expectations.
		expectedFillablePrice *big.Rat
		expectedError         error
	}{
		`Can calculate fillable price for a subaccount with one long position that is slightly
		below maintenance margin requirements`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// $49,999 = (49,999 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,999.9 notional sell order.
			expectedFillablePrice: big.NewRat(49_999, 100),
		},
		`Can calculate fillable price for a subaccount with one long position when bankruptcyAdjustmentPpm is 2_000_000`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm: 5_000,
				ValidatorFeePpm:     200_000,
				LiquidityFeePpm:     800_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           2_000_000,
					SpreadToMaintenanceMarginRatioPpm: 100_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},
			// $49,998 = (49,998 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,999.8 notional sell order.
			expectedFillablePrice: big.NewRat(49_998, 100),
		},
		`Can calculate fillable price for a subaccount with one long position when 
		spreadToMaintenanceMarginRatioPpm is 200_000`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm: 5_000,
				ValidatorFeePpm:     200_000,
				LiquidityFeePpm:     800_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           lib.OneMillion,
					SpreadToMaintenanceMarginRatioPpm: 200_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},
			// $49,998 = (49,998 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,999.8 notional sell order.
			expectedFillablePrice: big.NewRat(49_998, 100),
		},
		`Can calculate fillable price for a subaccount with one short position that is slightly
		below maintenance margin requirements`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_499),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// $50,001 = (50,001 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC short with a $5,000.1 notional buy order.
			expectedFillablePrice: big.NewRat(50_001, 100),
		},
		`Can calculate fillable price for a subaccount with one short position when bankruptcyAdjustmentPpm is 2_000_000`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_499),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm: 5_000,
				ValidatorFeePpm:     200_000,
				LiquidityFeePpm:     800_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           2_000_000,
					SpreadToMaintenanceMarginRatioPpm: 100_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			// $50,002 = (50,002 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC short with a $5,000.2 notional buy order.
			expectedFillablePrice: big.NewRat(50_002, 100),
		},
		`Can calculate fillable price for a subaccount with one short position when 
		SpreadToMaintenanceMarginRatioPpm is 200_000`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_499),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm: 5_000,
				ValidatorFeePpm:     200_000,
				LiquidityFeePpm:     800_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           lib.OneMillion,
					SpreadToMaintenanceMarginRatioPpm: 200_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			// $50,002 = (50,002 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC short with a $5,000.2 notional buy order.
			expectedFillablePrice: big.NewRat(50_002, 100),
		},
		"Can calculate fillable price for a subaccount with one long position at the bankruptcy price": {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_000),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// $49,500 = (495 / 1) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,950 notional sell order.
			expectedFillablePrice: big.NewRat(495, 1),
		},
		`Can calculate fillable price for a subaccount with one long position at the bankruptcy price
		where we are liquidating half of the position`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_000),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// $49,500 = (495 / 1) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,950 notional sell order.
			// Note that even though we are closing half of the position, the fillable price is the same as
			// if we were closing the full position because it's calculated based on the position size.
			expectedFillablePrice: big.NewRat(495, 1),
		},
		"Can calculate fillable price for a subaccount with one short position at the bankruptcy price": {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_000),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// $50,500 = (505 / 1) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC short with a $5,050 notional buy order.
			expectedFillablePrice: big.NewRat(505, 1),
		},
		"Can calculate fillable price for a subaccount with one long position below the bankruptcy price": {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_500),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// $49,500 = (495 / 1) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,950 notional sell order.
			expectedFillablePrice: big.NewRat(495, 1),
		},
		"Can calculate fillable price for a subaccount with one short position below the bankruptcy price": {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_500),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// $50,500 = (505 / 1) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC short with a $5,050 notional buy order.
			expectedFillablePrice: big.NewRat(505, 1),
		},
		"Can calculate fillable price for a subaccount with multiple long positions": {
			perpetualId: 1,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -490),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_FourThousandthsBTCLong,
				&constants.PerpetualPosition_OneTenthEthLong,
			},

			// $2976 = (372 / 125) subticks * QuoteCurrencyAtomicResolution / BaseCurrencyAtomicResolution.
			// This means we should close our 0.1 ETH long for $2,976 dollars.
			expectedFillablePrice: big.NewRat(372, 125),
		},
		`Can calculate fillable price when bankruptcyAdjustmentPpm is max uint32`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm: 5_000,
				ValidatorFeePpm:     200_000,
				LiquidityFeePpm:     800_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           math.MaxUint32,
					SpreadToMaintenanceMarginRatioPpm: 100_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			// $49,500 = (495 / 1) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,950 notional sell order.
			expectedFillablePrice: big.NewRat(495, 1),
		},
		`Can calculate fillable price when SpreadTomaintenanceMarginRatioPpm is 1`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm: 5_000,
				ValidatorFeePpm:     200_000,
				LiquidityFeePpm:     800_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           lib.OneMillion,
					SpreadToMaintenanceMarginRatioPpm: 1,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			expectedFillablePrice: big.NewRat(4_999_999_999, 10_000_000),
		},
		`Can calculate fillable price when SpreadTomaintenanceMarginRatioPpm is one million`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm: 5_000,
				ValidatorFeePpm:     200_000,
				LiquidityFeePpm:     800_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           lib.OneMillion,
					SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			// $49,990 = (49990 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,999 notional sell order.
			expectedFillablePrice: big.NewRat(49_990, 100),
		},
		`Returns error when subaccount does not have an open position for perpetual id`: {
			perpetualId: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{},

			expectedError: types.ErrInvalidPerpetualPositionSizeDelta,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				mock.Anything,
				constants.TDai.Denom,
			).Return(
				sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int))),
			)

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, &mocks.IndexerEventManager{})
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			// Initialize the liquidations config.
			if tc.liquidationConfig != nil {
				require.NoError(t,
					ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, *tc.liquidationConfig),
				)
			} else {
				require.NoError(t,
					ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, types.LiquidationsConfig_Default),
				)
			}

			// Create the tdai asset
			_, err := ks.AssetsKeeper.CreateAsset(ks.Ctx, constants.TDai.Id, constants.TDai.Symbol, constants.TDai.Denom, constants.TDai.DenomExponent, constants.TDai.HasMarket, constants.TDai.MarketId, constants.TDai.AtomicResolution, constants.TDai.AssetYieldIndex)
			require.NoError(t, err)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			// Create the subaccount.
			subaccount := satypes.Subaccount{
				Id: &satypes.SubaccountId{
					Owner:  "liquidations_test",
					Number: 0,
				},
				AssetPositions:     tc.assetPositions,
				PerpetualPositions: tc.perpetualPositions,
			}
			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, subaccount)

			fillablePrice, err := ks.ClobKeeper.GetFillablePrice(
				ks.Ctx,
				*subaccount.Id,
				tc.perpetualId,
			)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedFillablePrice, fillablePrice)
			}
		})
	}
}

func TestPlacePerpetualLiquidation_Deleveraging(t *testing.T) {
	tests := map[string]struct {
		// State.
		subaccounts                   []satypes.Subaccount
		insuranceFundBalance          uint64
		marketIdToOraclePriceOverride map[uint32]uint64

		// Parameters.
		liquidationConfig     types.LiquidationsConfig
		placedMatchableOrders []types.MatchableOrder
		order                 types.LiquidationOrder

		// Expectations.
		expectedFilledSize                satypes.BaseQuantums
		expectedOrderStatus               types.OrderStatus
		expectedSubaccountLiquidationInfo map[satypes.SubaccountId]types.SubaccountLiquidationInfo
		expectedLiquidationDeltaPerBlock  map[uint32]*big.Int
		expectedSubaccounts               []satypes.Subaccount
		expectedOperationsQueue           []types.OperationRaw
	}{
		`Can place a liquidation order that is fully filled and does not require deleveraging`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			insuranceFundBalance: 0,

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10, // Order at $50,000
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(100_000_000),
			expectedOrderStatus: types.Success,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(-250_000_000),
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_499_000_000 - 50_000_000_000 - 250_000_000),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(100_000_000_000), // $100,000
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10),
				clobtest.NewMatchOperationRaw(
					&constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderId(),
							FillAmount:   100_000_000,
						},
					},
				),
			},
		},
		`Can place a liquidation order that is partially filled and does not require deleveraging`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			insuranceFundBalance: 0,

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				// First order at $50,000
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				// Second order at $60,000, which does not cross the liquidation order
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price60000_GTB10,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(25_000_000),
			expectedOrderStatus: types.Success,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(-62_500_000),
			},
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
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
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
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.GetOrderId(),
							FillAmount:   25_000_000,
						},
					},
				),
			},
		},
		`Can place a liquidation order that is unfilled and full position size is deleveraged`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			insuranceFundBalance: 0,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC
			},

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $50,499, and closing at $50,500
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(0),
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id:              &constants.Carl_Num0,
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 50_499_000_000),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num0,
								FillAmount:             100_000_000,
							},
						},
					},
				),
			},
		},
		`Can place a liquidation order that is partially-filled and it's not deleveraged`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			insuranceFundBalance: 0,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC.
			},

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				// First order at $50,498, Carl pays $0.25 to the insurance fund.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11,
				// Carl's bankruptcy price to close 0.75 BTC short is $50,499, and closing at $50,500
				// would require $0.75 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(25_000_000),
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(-250000),
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_499_000_000 - (50_498_000_000 / 4) - 250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(-75_000_000), // -0.75 BTC
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + (50_498_000_000 / 4)),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(75_000_000), // 0.75 BTC
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11.GetOrderId(),
							FillAmount:   25_000_000,
						},
					},
				),
			},
		},
		`Can place a liquidation order that is unfilled and cannot be deleveraged due to
			non-overlapping bankruptcy prices`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
			},
			insuranceFundBalance: 0,

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing at $50,000
				// would require $1 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(0),
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails.
				// Dave's bankruptcy price to close 1 BTC long is $50,000, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				{
					Id:                 constants.Carl_Num0_1BTC_Short_49999USD.Id,
					AssetPositions:     constants.Carl_Num0_1BTC_Short_49999USD.AssetPositions,
					PerpetualPositions: constants.Carl_Num0_1BTC_Short_49999USD.PerpetualPositions,
					MarginEnabled:      constants.Carl_Num0_1BTC_Short_49999USD.MarginEnabled,
					AssetYieldIndex:    big.NewRat(1, 1).String(),
				},
				{
					Id:                 constants.Dave_Num0_1BTC_Long_50000USD_Short.Id,
					AssetPositions:     constants.Dave_Num0_1BTC_Long_50000USD_Short.AssetPositions,
					PerpetualPositions: constants.Dave_Num0_1BTC_Long_50000USD_Short.PerpetualPositions,
					MarginEnabled:      constants.Dave_Num0_1BTC_Long_50000USD_Short.MarginEnabled,
					AssetYieldIndex:    big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{},
		},
		`Can place a liquidation order that is partially-filled and cannot be deleveraged due to
			non-overlapping bankruptcy prices`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				constants.Dave_Num1_025BTC_Long_50000USD,
			},
			insuranceFundBalance: 0,

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num1_Id0_Clob0_Sell025BTC_Price49999_GTB10,
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing 0.75 BTC at $50,000
				// would require $0.75 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(25_000_000),
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
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
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				// Dave's bankruptcy price to close 1 BTC long is $50,000, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				// Dave_Num0 does not change since deleveraging against this subaccount failed.
				{
					Id:                 constants.Dave_Num0_1BTC_Long_50000USD_Short.Id,
					AssetPositions:     constants.Dave_Num0_1BTC_Long_50000USD_Short.AssetPositions,
					PerpetualPositions: constants.Dave_Num0_1BTC_Long_50000USD_Short.PerpetualPositions,
					MarginEnabled:      constants.Dave_Num0_1BTC_Long_50000USD_Short.MarginEnabled,
					AssetYieldIndex:    big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Dave_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 12_499_750_000),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num1_Id0_Clob0_Sell025BTC_Price49999_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num1_Id0_Clob0_Sell025BTC_Price49999_GTB10.GetOrderId(),
							FillAmount:   25_000_000,
						},
					},
				),
			},
		},
		`Can place a liquidation order that is unfilled, then only a portion of the remaining size can
			deleveraged due to non-overlapping bankruptcy prices with some subaccounts`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
				constants.Dave_Num1_05BTC_Long_50000USD,
			},
			insuranceFundBalance: 0,

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				// Carl's bankruptcy price to close 1 BTC short is $49,999, and closing 0.75 BTC at $50,000
				// would require $0.75 from the insurance fund. Since the insurance fund is empty,
				// deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(0),
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(0),
			},
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
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				// Dave_Num0 does not change since deleveraging against this subaccount failed.
				{
					Id:                 constants.Dave_Num0_1BTC_Long_50000USD_Short.Id,
					AssetPositions:     constants.Dave_Num0_1BTC_Long_50000USD_Short.AssetPositions,
					PerpetualPositions: constants.Dave_Num0_1BTC_Long_50000USD_Short.PerpetualPositions,
					MarginEnabled:      constants.Dave_Num0_1BTC_Long_50000USD_Short.MarginEnabled,
					AssetYieldIndex:    big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Dave_Num1,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 24_999_500_000),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualDeleveragingLiquidation(
					types.MatchPerpetualDeleveraging{
						Liquidated:  constants.Carl_Num0,
						PerpetualId: 0,
						Fills: []types.MatchPerpetualDeleveraging_Fill{
							{
								OffsettingSubaccountId: constants.Dave_Num1,
								FillAmount:             50_000_000,
							},
						},
					},
				),
			},
		},
		`Partially matched but fails due to insufficient insurance fund balance and deleveraging is skipped -
			negative TNC`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			insuranceFundBalance: 740_000, // $0.74
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC.
			},

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				// First order at $50,498, Carl pays $0.25 to the insurance fund.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11,
				// Carl's bankruptcy price to close 0.75 BTC short is $50,499, and closing at $50,500
				// would require $0.75 from the insurance fund. The insurance fund balance cannot
				// cover this loss so deleveraging is required to close this position.
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(25_000_000),
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(-250_000),
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_499_000_000 - (50_498_000_000 / 4) - 250_000),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(-75_000_000), // -0.75 BTC
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + (50_498_000_000 / 4)),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(75_000_000), // 0.75 BTC
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11.GetOrderId(),
							FillAmount:   25_000_000,
						},
					},
				),
			},
		},
		`Partially matched deleveraging is skipped - negative TNC`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			insuranceFundBalance: 750_000, // $0.75
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_050_000_000, // $50,500 / BTC.
			},

			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				// First order at $50,498, Carl pays $0.25 to the insurance fund.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11,
				// Carl's bankruptcy price to close 0.75 BTC short is $50,499, and closing at $50,500
				// would require $0.75 from the insurance fund. The insurance fund balance can
				// cover this loss so the liquidation succeeds.
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500, // Liquidation order at $50,500

			expectedFilledSize:  satypes.BaseQuantums(100_000_000),
			expectedOrderStatus: types.Success,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(500_000),
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id:              &constants.Carl_Num0,
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 50_499_500_000),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11,
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50498_GTB11.GetOrderId(),
							FillAmount:   25_000_000,
						},
						{
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.GetOrderId(),
							FillAmount:   75_000_000,
						},
					},
				),
			},
		},
		`Can place a liquidation order that is partially-filled and subaccount becomes non-liquidatable -
			deleveraging is skipped`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			insuranceFundBalance: 0,
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			placedMatchableOrders: []types.MatchableOrder{
				// First order at $50,000, matching against this order will make Carl's TNC >= MMR.
				// Insurance fund fee will be maxed out. DeltaQuoteQuantums = .25 BTC * $50k/BTC = $12,500.
				// Current InsuranceFundFeePpm = 5000.
				// Fee = $12,500 * 5000 / 1,000,000 = $62.5.
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				// Carl's account has now become well-collateralized, however the liquidation order
				// is for the full position size and will try to match against this high price order.
				// Because of the high price, the insurance fund will be needed to cover the loss but the insurance
				// fund is empty so we must deleverage. Our deleveraging algorithm will verify that the account is
				// still liquidatable before actually doing any deleveraging. For this test, we expect that the
				// account is no longer liquidatable and so deleveraging will be skipped.
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price60000_GTB10,
			},
			// Liquidation order at $60,000, setting high fill price to allow the order to
			// attempt to fill against the high price order above.
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price60000,

			expectedFilledSize:  satypes.BaseQuantums(25_000_000),
			expectedOrderStatus: types.LiquidationRequiresDeleveraging,
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated: []uint32{0},
				},
			},
			expectedLiquidationDeltaPerBlock: map[uint32]*big.Int{
				0: big.NewInt(-62_500_000),
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId: 0,
							Quantums: dtypes.NewInt(
								54_999_000_000 - 50_000_000_000/4 -
									lib.BigIntMulPpm(
										big.NewInt(50_000_000_000/4),
										constants.LiquidationsConfig_No_Limit.InsuranceFundFeePpm,
									).Int64(),
							),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(-75_000_000), // -0.75 BTC
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(50_000_000_000 + 50_000_000_000/4),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  0,
							Quantums:     dtypes.NewInt(75_000_000), // 0.75 BTC
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetYieldIndex: big.NewRat(1, 1).String(),
				},
			},
			expectedOperationsQueue: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price60000,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.GetOrderId(),
							FillAmount:   25_000_000,
						},
					},
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup memclob state and test expectations.
			memclob := memclob.NewMemClobPriceTimePriority(false)

			bankKeeper := &mocks.BankKeeper{}
			bankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			bankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			bankKeeper.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))
			bankKeeper.On(
				"GetBalance",
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(sdk.NewCoin("TDAI", sdkmath.NewIntFromUint64(tc.insuranceFundBalance))).Twice()

			mockIndexerEventManager := &mocks.IndexerEventManager{}
			mockIndexerEventManager.On("Enabled").Return(false)
			ks := keepertest.NewClobKeepersTestContext(t, memclob, bankKeeper, mockIndexerEventManager)
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			ctx := ks.Ctx.WithIsCheckTx(true)

			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParamsNoFee))

			// Set up TDAI asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			perpetuals := []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			}
			for _, perpetual := range perpetuals {
				_, err = ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					perpetual.Params.Id,
					perpetual.Params.Ticker,
					perpetual.Params.MarketId,
					perpetual.Params.AtomicResolution,
					perpetual.Params.DefaultFundingPpm,
					perpetual.Params.LiquidityTier,
					perpetual.Params.MarketType,
					perpetual.Params.DangerIndexPpm,
					perpetual.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					perpetual.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				perpetuals,
			)

			for _, s := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, s)
			}

			ks.ClobKeeper.DaemonDeleveragingInfo.UpdateSubaccountsWithPositions(
				clobtest.GetOpenPositionsFromSubaccounts(tc.subaccounts),
			)

			for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {
				err := ks.PricesKeeper.UpdateSpotAndPnlMarketPrices(
					ctx,
					&pricestypes.MarketPriceUpdate{
						MarketId:  marketId,
						SpotPrice: oraclePrice,
						PnlPrice:  oraclePrice,
					},
				)
				require.NoError(t, err)
			}

			// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
			// the indexer event manager to expect these events.
			for i, clobPair := range []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth_No_Fee,
			} {
				mockIndexerEventManager.On("AddTxnEvent",
					ctx,
					indexerevents.SubtypePerpetualMarket,
					indexerevents.PerpetualMarketEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewPerpetualMarketCreateEvent(
							uint32(i),
							uint32(i),
							perpetuals[i].Params.Ticker,
							perpetuals[i].Params.MarketId,
							clobPair.Status,
							clobPair.QuantumConversionExponent,
							perpetuals[i].Params.AtomicResolution,
							clobPair.SubticksPerTick,
							clobPair.StepBaseQuantums,
							perpetuals[i].Params.LiquidityTier,
							perpetuals[i].Params.MarketType,
							perpetuals[i].Params.DangerIndexPpm,
							fmt.Sprintf("%d", perpetuals[i].Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock),
						),
					),
				).Once().Return()
				_, err = ks.ClobKeeper.CreatePerpetualClobPair(
					ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, tc.liquidationConfig),
			)

			ks.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Timestamp: time.Unix(5, 0),
			})

			// Place all existing orders on the orderbook.
			for _, matchableOrder := range tc.placedMatchableOrders {
				require.False(t, matchableOrder.IsLiquidation())

				order := matchableOrder.MustGetOrder()

				// Get raw tx bytes for this short term order placement and set on context
				// so bytes are properly stored in OperationsToPropose
				shortTermOrderPlacement := clobtest.NewShortTermOrderPlacementOperationRaw(order)
				bytes := shortTermOrderPlacement.GetShortTermOrderPlacement()
				tempCtx := ctx.WithTxBytes(bytes)
				_, orderStatus, err := ks.ClobKeeper.PlaceShortTermOrder(tempCtx, &types.MsgPlaceOrder{Order: order.MustGetOrder()})
				require.NoError(t, err)
				require.Equal(t, types.Success, orderStatus)
			}

			// Run the test case and verify expectations.
			actualFillAmount, orderStatus, err := ks.ClobKeeper.PlacePerpetualLiquidation(ctx, tc.order)
			require.NoError(t, err)

			require.Equal(t, tc.expectedOrderStatus, orderStatus)
			require.Equal(t, tc.expectedFilledSize, actualFillAmount)

			for subaccountId, liquidationInfo := range tc.expectedSubaccountLiquidationInfo {
				require.Equal(
					t,
					liquidationInfo,
					ks.ClobKeeper.GetSubaccountLiquidationInfo(ctx, subaccountId),
				)
			}

			for perpetualId, expectedLiquidationDeltaPerBlock := range tc.expectedLiquidationDeltaPerBlock {
				liquidationDeltaPerBlock, err := ks.ClobKeeper.GetCumulativeInsuranceFundDelta(ctx, perpetualId)
				require.NoError(t, err)
				require.Equal(
					t,
					expectedLiquidationDeltaPerBlock,
					liquidationDeltaPerBlock,
				)
			}

			if tc.expectedFilledSize == 0 {
				// Bankruptcy price in DeleveragingEvent is not exposed by API. It is also
				// being tested in other e2e tests. So we don't test it here.
				mockIndexerEventManager.On("AddTxnEvent",
					mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return()
				_, err = ks.ClobKeeper.MaybeDeleverageSubaccount(
					ctx,
					tc.order.GetSubaccountId(),
					tc.order.MustGetLiquidatedPerpetualId(),
				)
				require.NoError(t, err)
			}

			for _, expectedSubaccount := range tc.expectedSubaccounts {
				require.Equal(t, expectedSubaccount, ks.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.GetId()))
			}

			require.Equal(
				t,
				tc.expectedOperationsQueue,
				ks.ClobKeeper.GetOperations(ctx).GetOperationsQueue(),
			)
		})
	}
}

func TestPlacePerpetualLiquidation_SendOffchainMessages(t *testing.T) {
	indexerEventManager := &mocks.IndexerEventManager{}
	for _, message := range constants.TestOffchainMessages {
		indexerEventManager.On("SendOffchainData", message).Once().Return()
	}

	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()

	bankMock := &mocks.BankKeeper{}
	bankMock.On(
		"GetBalance",
		mock.Anything,
		authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
		constants.TDai.Denom,
	).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

	ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, indexerEventManager)
	ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

	ctx := ks.Ctx.WithTxBytes(constants.TestTxBytes)
	// CheckTx mode set correctly
	ctx = ctx.WithIsCheckTx(true)
	prices.InitGenesis(ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	memClob.On("CreateOrderbook", ctx, constants.ClobPair_Btc).Return()
	// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
	// the indexer event manager to expect these events.
	indexerEventManager.On("AddTxnEvent",
		ctx,
		indexerevents.SubtypePerpetualMarket,
		indexerevents.PerpetualMarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewPerpetualMarketCreateEvent(
				0,
				0,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.Ticker,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketId,
				constants.ClobPair_Btc.Status,
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.AtomicResolution,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.StepBaseQuantums,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.LiquidityTier,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.MarketType,
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.DangerIndexPpm,
				fmt.Sprintf("%d", constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock),
			),
		),
	).Once().Return()
	_, err := ks.ClobKeeper.CreatePerpetualClobPair(
		ctx,
		constants.ClobPair_Btc.Id,
		clobtest.MustPerpetualId(constants.ClobPair_Btc),
		satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
		constants.ClobPair_Btc.QuantumConversionExponent,
		constants.ClobPair_Btc.SubticksPerTick,
		constants.ClobPair_Btc.Status,
	)
	require.NoError(t, err)

	order := constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000
	memClob.On("PlacePerpetualLiquidation", ctx, order).Return(
		satypes.BaseQuantums(100_000_000),
		types.Success,
		constants.TestOffchainUpdates,
		nil,
	)

	_, _, err = ks.ClobKeeper.PlacePerpetualLiquidation(ctx, order)
	require.NoError(t, err)

	indexerEventManager.AssertNumberOfCalls(t, "SendOffchainData", len(constants.TestOffchainMessages))
	indexerEventManager.AssertExpectations(t)
	memClob.AssertExpectations(t)
}

func TestIsLiquidatable(t *testing.T) {
	tests := map[string]struct {
		// State.
		perpetuals []perptypes.Perpetual

		// Subaccount state.
		assetPositions     []*satypes.AssetPosition
		perpetualPositions []*satypes.PerpetualPosition

		// Expectations.
		expectedIsLiquidatable bool
	}{
		"Subaccount with no open positions but positive net collateral is not liquidatable": {
			expectedIsLiquidatable: false,
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 1),
			),
		},
		"Subaccount with no open positions but negative net collateral is not liquidatable": {
			expectedIsLiquidatable: false,
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -1),
			),
		},
		"Subaccount at initial margin requirements is not liquidatable": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			perpetualPositions: []*satypes.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
				},
			},
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_000),
			),
			expectedIsLiquidatable: false,
		},
		"Subaccount below initial but at maintenance margin requirements is not liquidatable": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			perpetualPositions: []*satypes.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
				},
			},
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_500),
			),
			expectedIsLiquidatable: false,
		},
		"Subaccount below maintenance margin requirements is liquidatable": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			perpetualPositions: []*satypes.PerpetualPosition{
				{
					PerpetualId: uint32(0),
					Quantums:    dtypes.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
				},
			},
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			expectedIsLiquidatable: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Set up TDAI asset in assets module.
			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			// Create the subaccount.
			subaccount := satypes.Subaccount{
				Id: &satypes.SubaccountId{
					Owner:  "liquidations_test",
					Number: 0,
				},
				AssetPositions:     tc.assetPositions,
				PerpetualPositions: tc.perpetualPositions,
			}
			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, subaccount)
			isLiquidatable, err := ks.ClobKeeper.IsLiquidatable(ks.Ctx, *subaccount.Id)

			// Note that there should never be errors when passing the empty update.
			require.NoError(t, err)
			require.Equal(t, tc.expectedIsLiquidatable, isLiquidatable)
		})
	}
}

func TestGetBankruptcyPriceInQuoteQuantums(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		perpetualId   uint32
		deltaQuantums int64

		// Perpetual state.
		perpetuals []perptypes.Perpetual

		// Subaccount state.
		assetPositions     []*satypes.AssetPosition
		perpetualPositions []*satypes.PerpetualPosition

		// Expectations.
		expectedBankruptcyPriceQuoteQuantums *big.Int
		expectedError                        error
	}{
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one long position that is slightly below maintenance margin requirements`: {
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// 4,501,000,000 quote quantums = $4,501. This means if 0.1 BTC can't be sold for at
			// least $4,501 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(4_501_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one short position that is slightly below maintenance margin requirements`: {
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_499),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// -5,499,000,000 quote quantums = -$5,499. This means if 0.1 BTC can't be bought for
			// at most $5,499 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(-5_499_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one long position that is at the bankruptcy price`: {
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_000),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// 5,000,000,000 quote quantums = $5,000. This means if 0.1 BTC can't be sold for at
			// least $5,000 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(5_000_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is partially closing
		one long position that is at the bankruptcy price`: {
			perpetualId:   0,
			deltaQuantums: -5_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_000),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// 2,500,000,000 quote quantums = $2,500. This means if 0.1 BTC can't be sold for at
			// least $2,500 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(2_500_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is partially closing
		one short position that is at the bankruptcy price`: {
			perpetualId:   0,
			deltaQuantums: 5_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_000),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// -2,500,000,000 quote quantums = -$2,500. This means if 0.1 BTC can't be bought for at
			// most $2,500 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(-2_500_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one short position that is at the bankruptcy price`: {
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_000),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// -5,000,000,000 quote quantums = -$5,000. This means if 0.1 BTC can't be bought for at
			// most $5,000 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(-5_000_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one long position that is below the bankruptcy price`: {
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// 5,100,000,000 quote quantums = $5,100. This means if 0.1 BTC can't be sold for at
			// least $5,100 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(5_100_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one short position that is below the bankruptcy price`: {
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// -4,900,000,000 quote quantums = -$4,900. This means if 0.1 BTC can't be bought for at
			// most $4,900 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(-4_900_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one long position and has multiple long positions`: {
			perpetualId:   1,
			deltaQuantums: -100_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -490),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_FourThousandthsBTCLong,
				&constants.PerpetualPosition_OneTenthEthLong,
			},

			// 294,000,000 quote quantums = $294. This means if 0.1 ETH can't be sold for at
			// least $294 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(294_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one short position and has multiple short positions`: {
			perpetualId:   1,
			deltaQuantums: 100_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 510),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_FourThousandthsBTCShort,
				&constants.PerpetualPosition_OneTenthEthShort,
			},

			// -306,000,000 quote quantums = -$306. This means if 0.1 ETH can't be bought for at
			// most $306 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(-306_000_000),
		},
		`Can calculate bankruptcy price in quote quantums for a subaccount that is fully closing
		one short position and has a long and short position`: {
			perpetualId:   1,
			deltaQuantums: 100_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 110),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_FourThousandthsBTCLong,
				&constants.PerpetualPosition_OneTenthEthShort,
			},

			// -306,000,000 quote quantums = -$306. This means if 0.1 ETH can't be bought for at
			// most $306 then the subaccount will be bankrupt when this position is closed.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(-306_000_000),
		},
		`Rounds up bankruptcy price in quote quantums for a subaccount that is partially closing
		one long position that is slightly below maintenance margin requirements`: {
			perpetualId:   0,
			deltaQuantums: -21_347,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -13),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},

			// 2,776 quote quantums = $0.002776. This means if 0.00021347 BTC can't be sold for
			// at least $0.002776 then the subaccount will be bankrupt when this position is closed.
			// Note that the result is rounded up from 2,775.11 quote quantums.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(2_776),
		},
		`Rounds up bankruptcy price in quote quantums for a subaccount that is partially closing
		one short position that is below the bankruptcy price`: {
			perpetualId:   0,
			deltaQuantums: 21_347,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 13),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCShort,
			},

			// -2,775 quote quantums = $0.002775. This means if 0.00021347 BTC can't be bought for
			// at most $0.002775 then the subaccount will be bankrupt when this position is closed.
			// Note that the result is rounded down from 2,775.11 quote quantums.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(-2_775),
		},
		`Account with a long position that cannot be liquidated at a loss has a negative
		bankruptcy price in quote quantums`: {
			perpetualId:   0,
			deltaQuantums: -100_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			// Note that if quote balance is positive for longs, this indicates that the subaccount's
			// quote balance exceeds the notional value of their long position.
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},

			// -1,000,000 quote quantums = -$1,000,000. This means if 1 BTC can't be sold for
			// at least -$1,000,000 then the subaccount will be bankrupt when this position is closed.
			// Note this is not possible since it's impossible to sell a position for less than 0 dollars.
			expectedBankruptcyPriceQuoteQuantums: big.NewInt(-1_000_000),
		},
		`Returns error when deltaQuantums is zero`: {
			perpetualId:   0,
			deltaQuantums: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},

			expectedError: types.ErrInvalidPerpetualPositionSizeDelta,
		},
		`Returns error when subaccount does not have an open position for perpetual id`: {
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{},

			expectedError: types.ErrInvalidPerpetualPositionSizeDelta,
		},
		`Returns error when delta quantums and perpetual position have the same sign`: {
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},

			expectedError: types.ErrInvalidPerpetualPositionSizeDelta,
		},
		`Returns error when abs delta quantums is greater than position size`: {
			perpetualId:   0,
			deltaQuantums: -100_000_001,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},

			expectedError: types.ErrInvalidPerpetualPositionSizeDelta,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			require.NoError(t, keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper))

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			// Create the subaccount.
			subaccountId := satypes.SubaccountId{
				Owner:  "liquidations_test",
				Number: 0,
			}
			subaccount := satypes.Subaccount{
				Id:                 &subaccountId,
				AssetPositions:     tc.assetPositions,
				PerpetualPositions: tc.perpetualPositions,
			}
			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, subaccount)

			bankruptcyPriceInQuoteQuantums, err := ks.ClobKeeper.GetBankruptcyPriceInQuoteQuantums(
				ks.Ctx,
				*subaccount.Id,
				tc.perpetualId,
				big.NewInt(tc.deltaQuantums),
			)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedBankruptcyPriceQuoteQuantums, bankruptcyPriceInQuoteQuantums)

				// Verify that the returned delta quote quantums can pass `CanUpdateSubaccounts` function.
				success, _, err := ks.SubaccountsKeeper.CanUpdateSubaccounts(
					ks.Ctx,
					[]satypes.Update{
						{
							SubaccountId: subaccountId,
							AssetUpdates: keepertest.CreateTDaiAssetUpdate(bankruptcyPriceInQuoteQuantums),
							PerpetualUpdates: []satypes.PerpetualUpdate{
								{
									PerpetualId:      tc.perpetualId,
									BigQuantumsDelta: big.NewInt(tc.deltaQuantums),
								},
							},
						},
					},
					satypes.CollatCheck,
				)

				require.True(t, success)
				require.NoError(t, err)
			}
		})
	}
}

func TestGetLiquidationInsuranceFundFeeAndRemainingAvailableCollateral(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		perpetualId uint32
		isBuy       bool
		fillAmount  uint64
		subticks    types.Subticks

		liquidationConfig *types.LiquidationsConfig

		// Perpetual and subaccount state.
		perpetuals []perptypes.Perpetual

		// Subaccount state.
		assetPositions     []*satypes.AssetPosition
		perpetualPositions []*satypes.PerpetualPosition

		// Expectations.
		expectedLiquidationInsuranceFundDeltaBig *big.Int
		expectedRemainingQuoteQuantumsBig        *big.Int
		expectedError                            error
	}{
		`Fully closing one long position above the bankruptcy price and pays max liquidation fee`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    56_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, 5,610,000,000 quote quantums.
			// abs(5,610,000,000) * 0.5% max liquidation fee < 5,610,000,000 - 5,100,000,000, so the max
			// liquidation fee is returned.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(28_050_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(481_950_000),
		},
		`Fully closing one long position above the bankruptcy price pays max liquidation fee 
		when InsuranceFundFeePpm is 25_000`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    56_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm:             25_000,
				ValidatorFeePpm:                 200_000,
				LiquidityFeePpm:                 800_000,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, 5,610,000,000 quote quantums.
			// abs(5,610,000,000) * 2.5% max liquidation fee < 5,610,000,000 - 5,100,000,000, so the max
			// liquidation fee is returned.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(140_250_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(369_750_000),
		},
		`Fully closing one long position above the bankruptcy price pays less than max liquidation fee 
		when InsuranceFundFeePpm is one million`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    56_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm:             1_000_000,
				ValidatorFeePpm:                 200_000,
				LiquidityFeePpm:                 800_000,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, 5,610,000,000 quote quantums.
			// abs(5,610,000,000) * 100% max liquidation fee > 5,610,000,000 - 5,100,000,000, so all
			// of the leftover collateral is transferred to the insurance fund.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(510_000_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(0),
		},
		`Fully closing one short position above the bankruptcy price and pays max liquidation fee`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    44_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, -4,410,000,000 quote quantums.
			// abs(-4,410,000,000) * 0.5% max liquidation fee < -4,900,000,000 - -4,410,000,000, so
			// the max liquidation fee is returned.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(22_050_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(467_950_000),
		},
		`Fully closing one short position above the bankruptcy price and pays max liquidation fee
		when InsuranceFundFeePpm is 25_000`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    44_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},
			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm:             25_000,
				ValidatorFeePpm:                 200_000,
				LiquidityFeePpm:                 800_000,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, -4,410,000,000 quote quantums.
			// abs(-4,410,000,000) * 2.5% max liquidation fee < -4,900,000,000 - -4,410,000,000, so
			// the max liquidation fee is returned.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(110_250_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(379_750_000),
		},
		`Fully closing one short position above the bankruptcy price and pays less than max liquidation fee
		when InsuranceFundFeePpm is one million`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    44_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},
			liquidationConfig: &types.LiquidationsConfig{
				InsuranceFundFeePpm:             1_000_000,
				ValidatorFeePpm:                 200_000,
				LiquidityFeePpm:                 800_000,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, -4,410,000,000 quote quantums.
			// abs(-4,410,000,000) * 100% max liquidation fee > -4,900,000,000 - -4,410,000,000, so all
			// of the leftover collateral is transferred to the insurance fund.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(490_000_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(0),
		},
		`Fully closing one long position above the bankruptcy price and pays less than max
		liquidation fee`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    51_051_000_000, // 0.1% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 0.1% above bankruptcy price, 5,105,100,000 quote quantums.
			// 5,105,100,000 * 0.5% max liquidation fee > 5,105,100,000 - 5,100,000,000, so all
			// of the leftover collateral is transferred to the insurance fund.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(5_100_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(0),
		},
		`Fully closing one short position above the bankruptcy price and pays less than max
		liquidation fee`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    48_951_000_000, // 0.1% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 0.1% above bankruptcy price, -4,895,100,000 quote quantums.
			// -4,895,100,000 * 0.5% max liquidation fee < -4,895,100,000 - -4,900,000,000, so all
			// of the leftover collateral is transferred to the insurance fund.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(4_900_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(0),
		},
		`Fully closing one long position at the bankruptcy price and the delta is 0`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    51_000_000_000, // 0% above bankruptcy price (equal).

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 0% above bankruptcy price, 5,100,000,000 quote quantums.
			// 5,100,000,000 * 0.5% max liquidation fee > 5,100,000,000 - 5,100,000,000, so all
			// of the leftover collateral (which is zero) is transferred to the insurance fund.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(0),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(0),
		},
		`Fully closing one short position above the bankruptcy price and the delta is 0`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    49_000_000_000, // 0% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 0.1% above bankruptcy price, -4,900,000,000 quote quantums.
			// -4,900,000,000 * 0.5% max liquidation fee < -4,900,000,000 - -4,900,000,000, so all
			// of the leftover collateral (which is zero) is transferred to the insurance fund.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(0),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(0),
		},
		`Fully closing one long position below the bankruptcy price and the insurance fund must
		cover the loss`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    50_490_000_000, // 1% below bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 1% below the bankruptcy price, 5,049,000,000 quote quantums.
			// 5,049,000,000 - 5,100,000,000 < 0, so the insurance fund must cover the losses.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(-51_000_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(0),
		},
		`If fully closing one short position below the bankruptcy price the insurance fund must
		cover the loss`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    49_490_000_000, // 1% below bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 1% below the bankruptcy price, -4,949,000,000 quote quantums.
			// -4,949,000,000 - -4,900,000,000 < 0, so the insurance fund msut cover the losses.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(-49_000_000),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(0),
		},
		"Returns error when delta quantums is zero": {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  0,
			subticks:    50_000_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},
			expectedError: types.ErrInvalidQuantumsForInsuranceFundDeltaCalculation,
		},
		"Succeeds when delta quote quantums is zero": {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000, // 0.1 BTC delta.
			subticks:    1,          // Quote quantums for 0.1 BTC is 1/10, rounded to zero.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: keepertest.CreateTDaiAssetPosition(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Insurance fund delta before applying position limit is 0 - -4,900,000,000 = 4,900,000,000.
			// abs(0) * 0.5% max liquidation fee < 4,900,000,000, so overall delta is zero.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(0),
			expectedRemainingQuoteQuantumsBig:        big.NewInt(4_900_000_000),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, mockIndexerEventManager)
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Set up TDAI asset in assets module.
			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			// Create clob pair.
			mockIndexerEventManager.On("AddTxnEvent",
				ks.Ctx,
				indexerevents.SubtypePerpetualMarket,
				indexerevents.PerpetualMarketEventVersion,
				indexer_manager.GetBytes(
					indexerevents.NewPerpetualMarketCreateEvent(
						0,
						0,
						tc.perpetuals[0].Params.Ticker,
						tc.perpetuals[0].Params.MarketId,
						constants.ClobPair_Btc.Status,
						constants.ClobPair_Btc.QuantumConversionExponent,
						tc.perpetuals[0].Params.AtomicResolution,
						constants.ClobPair_Btc.SubticksPerTick,
						constants.ClobPair_Btc.StepBaseQuantums,
						tc.perpetuals[0].Params.LiquidityTier,
						tc.perpetuals[0].Params.MarketType,
						tc.perpetuals[0].Params.DangerIndexPpm,
						fmt.Sprintf("%d", tc.perpetuals[0].Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock),
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPair(
				ks.Ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)

			// Create the subaccount.
			subaccount := satypes.Subaccount{
				Id: &satypes.SubaccountId{
					Owner:  "liquidations_test",
					Number: 0,
				},
				AssetPositions:     tc.assetPositions,
				PerpetualPositions: tc.perpetualPositions,
			}
			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, subaccount)

			// Initialize the liquidations config.
			if tc.liquidationConfig != nil {
				require.NoError(
					t,
					ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, *tc.liquidationConfig),
				)
			} else {
				require.NoError(
					t,
					ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, types.LiquidationsConfig_Default),
				)
			}

			// Run the test and verify expectations.
			remainingQuoteQuantumsBig, liquidationInsuranceFundDeltaBig, err := ks.ClobKeeper.GetLiquidationInsuranceFundFeeAndRemainingAvailableCollateral(
				ks.Ctx,
				*subaccount.Id,
				tc.perpetualId,
				tc.isBuy,
				tc.fillAmount,
				tc.subticks,
			)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(
					t,
					tc.expectedLiquidationInsuranceFundDeltaBig.Int64(),
					liquidationInsuranceFundDeltaBig.Int64(),
				)
				require.Equal(
					t,
					tc.expectedRemainingQuoteQuantumsBig.Int64(),
					remainingQuoteQuantumsBig.Int64(),
				)
			}
		})
	}
}

func TestConvertLiquidationPriceToSubticks(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		liquidationPrice  *big.Rat
		isLiquidatingLong bool
		clobPair          types.ClobPair

		// Expectations.
		expectedSubticks types.Subticks
	}{
		`Converts liquidation price to subticks for liquidating a BTC long position`: {
			liquidationPrice: big.NewRat(
				int64(constants.FiveBillion),
				1,
			),
			isLiquidatingLong: true,
			clobPair:          constants.ClobPair_Btc,

			expectedSubticks: 500_000_000_000_000_000,
		},
		`Converts liquidation price to subticks for liquidating a BTC short position`: {
			liquidationPrice: big.NewRat(
				int64(constants.FiveBillion),
				1,
			),
			isLiquidatingLong: false,
			clobPair:          constants.ClobPair_Btc,

			expectedSubticks: 500_000_000_000_000_000,
		},
		`Converts liquidation price to subticks for liquidating a long position and rounds up`: {
			liquidationPrice: big.NewRat(
				7,
				1,
			),
			isLiquidatingLong: true,
			clobPair: types.ClobPair{
				SubticksPerTick:           100,
				QuantumConversionExponent: 1,
			},

			expectedSubticks: 100,
		},
		`Converts liquidation price to subticks for liquidating a short position and rounds down`: {
			liquidationPrice: big.NewRat(
				197,
				1,
			),
			isLiquidatingLong: true,
			clobPair: types.ClobPair{
				SubticksPerTick:           100,
				QuantumConversionExponent: 1,
			},

			expectedSubticks: 100,
		},
		`Converts liquidation price to subticks for liquidating a short position and rounds down, but
		the result is lower bounded at SubticksPerTick`: {
			liquidationPrice: big.NewRat(
				7,
				1,
			),
			isLiquidatingLong: true,
			clobPair: types.ClobPair{
				SubticksPerTick:           100,
				QuantumConversionExponent: 1,
			},

			expectedSubticks: 100,
		},
		`Converts zero liquidation price to subticks for liquidating a short position and rounds down,
		but the result is lower bounded at SubticksPerTick`: {
			liquidationPrice: big.NewRat(
				0,
				1,
			),
			isLiquidatingLong: true,
			clobPair: types.ClobPair{
				SubticksPerTick:           100,
				QuantumConversionExponent: 1,
			},

			expectedSubticks: 100,
		},
		`Converts liquidation price to subticks for liquidating a long position and rounds up, but
		the result is upper bounded at the max Uint64 that is most aligned with SubticksPerTick`: {
			liquidationPrice: big_testutil.MustFirst(
				new(big.Rat).SetString("10000000000000000000000"),
			),
			isLiquidatingLong: true,
			clobPair: types.ClobPair{
				SubticksPerTick:           100,
				QuantumConversionExponent: 1,
			},

			expectedSubticks: 18_446_744_073_709_551_600,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Run the test.
			subticks := ks.ClobKeeper.ConvertLiquidationPriceToSubticks(
				ks.Ctx,
				tc.liquidationPrice,
				tc.isLiquidatingLong,
				tc.clobPair,
			)
			require.Equal(
				t,
				tc.expectedSubticks.ToBigInt().String(),
				subticks.ToBigInt().String(),
			)
		})
	}
}

func TestConvertLiquidationPriceToSubticks_PanicsOnNegativeLiquidationPrice(t *testing.T) {
	// Setup keeper state.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	// Run the test.
	require.Panics(t, func() {
		ks.ClobKeeper.ConvertLiquidationPriceToSubticks(
			ks.Ctx,
			big.NewRat(-1, 1),
			false,
			constants.ClobPair_Btc,
		)
	})
}

func TestGetBestPerpetualPositionToLiquidate(t *testing.T) {
	tests := map[string]struct {
		// Subaccount state.
		perpetualPositions []*satypes.PerpetualPosition
		// Perpetual state.
		perpetuals []perptypes.Perpetual
		// Clob state.
		liquidationConfig types.LiquidationsConfig
		// CLOB pair state.
		clobPairs []types.ClobPair

		// Expectations.
		expectedClobPair types.ClobPair
		expectedQuantums *big.Int
	}{
		`Full position size is returned when subaccount has one perpetual long position`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: constants.LiquidationsConfig_No_Limit,

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).Neg(
				constants.PerpetualPosition_OneTenthBTCLong.GetBigQuantums(),
			),
		},
		`full position is returned when position size is smaller than StepBaseQuantums`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				{
					PerpetualId: 0,
					Quantums:    dtypes.NewInt(5),
				},
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm:             5_000,
				ValidatorFeePpm:                 0,
				LiquidityFeePpm:                 0,
				FillablePriceConfig:             constants.FillablePriceConfig_Default,
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
			},

			clobPairs: []types.ClobPair{
				// StepBaseQuantums is 10.
				constants.ClobPair_Btc3,
			},

			expectedClobPair: constants.ClobPair_Btc3,
			expectedQuantums: new(big.Int).SetInt64(-5),
		},
		`Full position size is returned when subaccount has one perpetual short position`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCShort,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: constants.LiquidationsConfig_No_Limit,

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).Neg(
				constants.PerpetualPosition_OneBTCShort.GetBigQuantums(),
			),
		},
		`Full position size of max uint64 of perpetual and CLOB pair are returned when subaccount
		has one long perpetual position at max position size`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_MaxUint64EthLong,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: constants.LiquidationsConfig_No_Limit,

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},

			expectedClobPair: constants.ClobPair_Eth,
			expectedQuantums: new(big.Int).Neg(
				new(big.Int).SetUint64(18446744073709551615),
			),
		},
		`Full position size of negated max uint64 of perpetual and CLOB pair are returned when
		subaccount has one short perpetual position at max position size`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_MaxUint64EthShort,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: constants.LiquidationsConfig_No_Limit,

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},

			expectedClobPair: constants.ClobPair_Eth,
			expectedQuantums: new(big.Int).Neg(
				big_testutil.MustFirst(
					new(big.Int).SetString("-18446744073709551615", 10),
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			// Create the subaccount.
			subaccount := satypes.Subaccount{
				Id: &satypes.SubaccountId{
					Owner:  "liquidations_test",
					Number: 0,
				},
				PerpetualPositions: tc.perpetualPositions,
			}
			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, subaccount)

			// Create the CLOB pairs and store the expected CLOB pair.
			for i, clobPair := range tc.clobPairs {
				perpetualId := clobtest.MustPerpetualId(clobPair)
				// PerpetualMarketCreateEvents are emitted when initializing the genesis state, so we need to mock
				// the indexer event manager to expect these events.
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexerevents.PerpetualMarketEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewPerpetualMarketCreateEvent(
							perpetualId,
							uint32(i),
							tc.perpetuals[perpetualId].Params.Ticker,
							tc.perpetuals[perpetualId].Params.MarketId,
							clobPair.Status,
							clobPair.QuantumConversionExponent,
							tc.perpetuals[perpetualId].Params.AtomicResolution,
							clobPair.SubticksPerTick,
							clobPair.StepBaseQuantums,
							tc.perpetuals[perpetualId].Params.LiquidityTier,
							tc.perpetuals[perpetualId].Params.MarketType,
							tc.perpetuals[perpetualId].Params.DangerIndexPpm,
							fmt.Sprintf("%d", tc.perpetuals[perpetualId].Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock),
						),
					),
				).Once().Return()
				_, err := ks.ClobKeeper.CreatePerpetualClobPair(
					ks.Ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}
			// Initialize the liquidations config.
			err := ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, tc.liquidationConfig)
			require.NoError(t, err)

			perpetualId, err := ks.ClobKeeper.GetBestPerpetualPositionToLiquidate(
				ks.Ctx,
				*subaccount.Id,
			)
			require.NoError(t, err)

			deltaQuantums, err := ks.ClobKeeper.GetNegativePositionSize(
				ks.Ctx,
				*subaccount.Id,
				perpetualId,
			)
			require.NoError(t, err)
			require.Equal(t, tc.expectedQuantums, deltaQuantums)

			expectedPerpetualId, err := tc.expectedClobPair.GetPerpetualId()
			require.NoError(t, err)
			require.Equal(
				t,
				expectedPerpetualId,
				perpetualId,
			)
		})
	}
}

func TestMaybeGetLiquidationOrder(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs          []types.ClobPair
		existingOrders []types.Order

		// Parameters.
		liquidatableSubaccount satypes.SubaccountId
		setupState             func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext)

		// Expectations.
		expectedErr           error
		expectedPlacedOrders  []*types.MsgPlaceOrder
		expectedMatchedOrders []*types.ClobMatch
	}{
		`Subaccount liquidation matches no maker orders`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
			},

			liquidatableSubaccount: constants.Dave_Num0,

			expectedPlacedOrders:  []*types.MsgPlaceOrder{},
			expectedMatchedOrders: []*types.ClobMatch{},
		},
		`Subaccount liquidation matches maker orders`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
				constants.Order_Carl_Num0_Id4_Clob0_Buy05BTC_Price40000,
			},

			liquidatableSubaccount: constants.Dave_Num0,

			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				},
				{
					Order: constants.Order_Carl_Num0_Id3_Clob0_Buy025BTC_Price50000,
				},
			},
			expectedMatchedOrders: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						ClobPairId:  constants.ClobPair_Btc.Id,
						IsBuy:       false,
						TotalSize:   100_000_000,
						Liquidated:  constants.Dave_Num0,
						PerpetualId: constants.ClobPair_Btc.GetPerpetualClobMetadata().PerpetualId,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   50_000_000,
							},
							{
								MakerOrderId: types.OrderId{},
								FillAmount:   25_000_000,
							},
						},
					},
				),
			},
		},
		`Does not place liquidation order if subaccount has no perpetual positions to liquidate`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num1_Short_500USD,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs:          []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{},

			liquidatableSubaccount: constants.Dave_Num0,
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(ctx, constants.Dave_Num0, 0)
			},

			expectedErr:           types.ErrNoPerpetualPositionsToLiquidate,
			expectedPlacedOrders:  []*types.MsgPlaceOrder{},
			expectedMatchedOrders: []*types.ClobMatch{},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))
			mockBankKeeper.On(
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			// Give the insurance fund a 1M TDAI balance.
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				perptypes.InsuranceFundModuleAddress,
				constants.TDai.Denom,
			).Return(
				sdk.NewCoin(
					constants.TDai.Denom,
					sdkmath.NewIntFromBigInt(big.NewInt(1_000_000_000_000)),
				),
			)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))
			ctx := ks.Ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParams))

			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobs {
				_, err = ks.ClobKeeper.CreatePerpetualClobPair(
					ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			// Initialize the liquidations config.
			err = ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default)
			require.NoError(t, err)

			if tc.setupState != nil {
				tc.setupState(ctx, ks)
			}

			// Create all existing orders.
			for _, order := range tc.existingOrders {
				_, _, err := ks.ClobKeeper.PlaceShortTermOrder(ctx, &types.MsgPlaceOrder{Order: order})
				require.NoError(t, err)
			}

			// Run the test.
			liquidationOrder, err := ks.ClobKeeper.MaybeGetLiquidationOrder(ctx, tc.liquidatableSubaccount)

			// Verify test expectations.
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, liquidationOrder)
				_, _, err := ks.ClobKeeper.PlacePerpetualLiquidation(ctx, *liquidationOrder)
				require.NoError(t, err)

				// TODO(DEC-1979): Refactor these tests to support the operations queue refactor.
				// placedOrders, matchedOrders := memClob.GetPendingFills(ctx)
				// require.Equal(t, tc.expectedPlacedOrders, placedOrders, "Placed orders lists are not equal")
				// require.Equal(t, tc.expectedMatchedOrders, matchedOrders, "Matched orders lists are not equal")
			}
		})
	}
}

func TestGetNextSubaccountToLiquidate(t *testing.T) {
	tests := map[string]struct {
		// Inputs
		subaccountIds                 []heap.LiquidationPriority
		isolatedPositionsPriorityHeap []heap.LiquidationPriority
		numIsolatedLiquidations       int

		// Expected outputs
		expectedSubaccountId                  satypes.SubaccountId
		expectedNumIsolated                   int
		expectedIsolatedPositionsPriorityHeap *heap.LiquidationPriorityHeap
		expectedSubaccountIds                 *heap.LiquidationPriorityHeap
	}{
		"returns nil when both heaps are empty": {
			subaccountIds:                 []heap.LiquidationPriority{},
			isolatedPositionsPriorityHeap: []heap.LiquidationPriority{},
			numIsolatedLiquidations:       0,

			expectedSubaccountId:                  satypes.SubaccountId{},
			expectedNumIsolated:                   0,
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
		},
		"returns from subaccountIds when available": {
			subaccountIds: []heap.LiquidationPriority{
				{SubaccountId: constants.Alice_Num0, Priority: big.NewFloat(100)},
			},
			isolatedPositionsPriorityHeap: []heap.LiquidationPriority{},
			numIsolatedLiquidations:       0,

			expectedSubaccountId:                  constants.Alice_Num0,
			expectedNumIsolated:                   0,
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
		},
		"switches to isolated positions when subaccountIds is empty": {
			subaccountIds: []heap.LiquidationPriority{},
			isolatedPositionsPriorityHeap: []heap.LiquidationPriority{
				{SubaccountId: constants.Bob_Num0, Priority: big.NewFloat(200)},
			},
			numIsolatedLiquidations: 0,

			expectedSubaccountId:                  constants.Bob_Num0,
			expectedNumIsolated:                   -1000000,
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
		},
		"returns from subaccountIds when multiple subaccounts are available": {
			subaccountIds: []heap.LiquidationPriority{
				{SubaccountId: constants.Bob_Num0, Priority: big.NewFloat(100)},
				{SubaccountId: constants.Alice_Num0, Priority: big.NewFloat(50)},
			},
			isolatedPositionsPriorityHeap: []heap.LiquidationPriority{},
			numIsolatedLiquidations:       0,

			expectedSubaccountId:                  constants.Alice_Num0,
			expectedNumIsolated:                   0,
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountIds: &heap.LiquidationPriorityHeap{
				{SubaccountId: constants.Bob_Num0, Priority: big.NewFloat(100)},
			},
		},
		"returns from subaccountIds when subaccount exists in both normal and isolated": {
			subaccountIds: []heap.LiquidationPriority{
				{SubaccountId: constants.Alice_Num0, Priority: big.NewFloat(100)},
			},
			isolatedPositionsPriorityHeap: []heap.LiquidationPriority{
				{SubaccountId: constants.Alice_Num0, Priority: big.NewFloat(200)},
			},
			numIsolatedLiquidations: 0,

			expectedSubaccountId: constants.Alice_Num0,
			expectedNumIsolated:  0,
			expectedIsolatedPositionsPriorityHeap: &heap.LiquidationPriorityHeap{
				{SubaccountId: constants.Alice_Num0, Priority: big.NewFloat(200)},
			},
			expectedSubaccountIds: heap.NewLiquidationPriorityHeap(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			subaccountIds := heap.NewLiquidationPriorityHeap()
			for _, priority := range tc.subaccountIds {
				subaccountIds.AddSubaccount(priority.SubaccountId, priority.Priority)
			}

			isolatedPositionsPriorityHeap := heap.NewLiquidationPriorityHeap()
			for _, priority := range tc.isolatedPositionsPriorityHeap {
				isolatedPositionsPriorityHeap.AddSubaccount(priority.SubaccountId, priority.Priority)
			}

			// Call the function.
			_, subaccountId := ks.ClobKeeper.GetNextSubaccountToLiquidate(
				ks.Ctx,
				subaccountIds,
				isolatedPositionsPriorityHeap,
				&tc.numIsolatedLiquidations,
			)

			// Check the results.
			if subaccountId == nil {
				require.Equal(t, tc.expectedSubaccountId, satypes.SubaccountId{})
			} else {
				require.Equal(t, tc.expectedSubaccountId, subaccountId.SubaccountId)
			}
			require.Equal(t, tc.expectedNumIsolated, tc.numIsolatedLiquidations)
			require.Equal(t, tc.expectedIsolatedPositionsPriorityHeap.Len(), isolatedPositionsPriorityHeap.Len())
			require.Equal(t, tc.expectedSubaccountIds.Len(), subaccountIds.Len())
		})
	}
}

func TestGetHealth(t *testing.T) {
	tests := map[string]struct {
		netCollateral     *big.Int
		maintenanceMargin *big.Int
		expectedHealth    *big.Float
	}{
		"negative net collateral returns 0": {
			netCollateral:     big.NewInt(-100),
			maintenanceMargin: big.NewInt(50),
			expectedHealth:    big.NewFloat(0),
		},
		"zero maintenance margin returns max float64": {
			netCollateral:     big.NewInt(100),
			maintenanceMargin: big.NewInt(0),
			expectedHealth:    big.NewFloat(math.MaxFloat64),
		},
		"negative maintenance margin returns max float64": {
			netCollateral:     big.NewInt(100),
			maintenanceMargin: big.NewInt(-50),
			expectedHealth:    big.NewFloat(math.MaxFloat64),
		},
		"normal case - health less than 1": {
			netCollateral:     big.NewInt(50),
			maintenanceMargin: big.NewInt(100),
			expectedHealth:    big.NewFloat(0.5),
		},
		"normal case - health equal to 1": {
			netCollateral:     big.NewInt(100),
			maintenanceMargin: big.NewInt(100),
			expectedHealth:    big.NewFloat(1),
		},
		"normal case - health greater than 1": {
			netCollateral:     big.NewInt(150),
			maintenanceMargin: big.NewInt(100),
			expectedHealth:    big.NewFloat(1.5),
		},
		"large numbers": {
			netCollateral:     new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), // 10^18
			maintenanceMargin: new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil), // 10^15
			expectedHealth:    big.NewFloat(1000),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := keeper.GetHealth(tc.netCollateral, tc.maintenanceMargin)

			// Compare the result with the expected value
			if result.Cmp(tc.expectedHealth) != 0 {
				t.Errorf("Expected health %v, but got %v", tc.expectedHealth, result)
			}
		})
	}
}

func TestCalculateLiquidationPriority(t *testing.T) {
	tests := map[string]struct {
		totalNetCollateral        *big.Int
		totalMaintenanceMargin    *big.Int
		weightedMaintenanceMargin *big.Int
		expectedPriority          *big.Float
	}{
		"zero weighted maintenance margin returns max float64": {
			totalNetCollateral:        big.NewInt(100),
			totalMaintenanceMargin:    big.NewInt(50),
			weightedMaintenanceMargin: big.NewInt(0),
			expectedPriority:          big.NewFloat(math.MaxFloat64),
		},
		"negative weighted maintenance margin returns max float64": {
			totalNetCollateral:        big.NewInt(100),
			totalMaintenanceMargin:    big.NewInt(50),
			weightedMaintenanceMargin: big.NewInt(-10),
			expectedPriority:          big.NewFloat(math.MaxFloat64),
		},
		"normal case - health less than 1": {
			totalNetCollateral:        big.NewInt(50),
			totalMaintenanceMargin:    big.NewInt(100),
			weightedMaintenanceMargin: big.NewInt(200),
			expectedPriority:          big.NewFloat(0.0025), // (50/100) / 200
		},
		"normal case - health equal to 1": {
			totalNetCollateral:        big.NewInt(100),
			totalMaintenanceMargin:    big.NewInt(100),
			weightedMaintenanceMargin: big.NewInt(200),
			expectedPriority:          big.NewFloat(0.005), // (100/100) / 200
		},
		"normal case - health greater than 1": {
			totalNetCollateral:        big.NewInt(150),
			totalMaintenanceMargin:    big.NewInt(100),
			weightedMaintenanceMargin: big.NewInt(200),
			expectedPriority:          big.NewFloat(0.0075), // (150/100) / 200
		},
		"large numbers": {
			totalNetCollateral:        new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), // 10^18
			totalMaintenanceMargin:    new(big.Int).Exp(big.NewInt(10), big.NewInt(15), nil), // 10^15
			weightedMaintenanceMargin: new(big.Int).Exp(big.NewInt(10), big.NewInt(16), nil), // 10^16
			expectedPriority:          new(big.Float).SetFloat64(1e-13),                      // (10^18/10^15) / 10^16 = 1000 / 10^16 = 10^-13
		},
		"negative net collateral": {
			totalNetCollateral:        big.NewInt(-100),
			totalMaintenanceMargin:    big.NewInt(50),
			weightedMaintenanceMargin: big.NewInt(200),
			expectedPriority:          big.NewFloat(0), // (0/50) / 200 = 0
		},
		"zero maintenance margin": {
			totalNetCollateral:        big.NewInt(100),
			totalMaintenanceMargin:    big.NewInt(0),
			weightedMaintenanceMargin: big.NewInt(1),
			expectedPriority:          big.NewFloat(math.MaxFloat64), // MaxFloat64
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := keeper.CalculateLiquidationPriority(tc.totalNetCollateral, tc.totalMaintenanceMargin, tc.weightedMaintenanceMargin)

			// Compare the result with the expected value
			if !almostEqual(result, tc.expectedPriority, 0.000001) {
				t.Errorf("Expected priority %v, but got %v", tc.expectedPriority, result)
			}
		})
	}
}

// almostEqual compares two big.Float values with a given epsilon for floating-point comparison
func almostEqual(a, b *big.Float, epsilon float64) bool {
	diff := new(big.Float).Sub(a, b)
	return diff.Abs(diff).Cmp(big.NewFloat(epsilon)) < 0
}

func TestGetMostAggressivePrice(t *testing.T) {
	tests := map[string]struct {
		bankruptcyPrice *big.Rat
		fillablePrice   *big.Rat
		isLong          bool
		expectedPrice   *big.Rat
	}{
		"long position - bankruptcy price lower": {
			bankruptcyPrice: big.NewRat(90, 1),
			fillablePrice:   big.NewRat(100, 1),
			isLong:          true,
			expectedPrice:   big.NewRat(90, 1),
		},
		"long position - fillable price lower": {
			bankruptcyPrice: big.NewRat(110, 1),
			fillablePrice:   big.NewRat(100, 1),
			isLong:          true,
			expectedPrice:   big.NewRat(100, 1),
		},
		"long position - prices equal": {
			bankruptcyPrice: big.NewRat(100, 1),
			fillablePrice:   big.NewRat(100, 1),
			isLong:          true,
			expectedPrice:   big.NewRat(100, 1),
		},
		"short position - bankruptcy price higher": {
			bankruptcyPrice: big.NewRat(110, 1),
			fillablePrice:   big.NewRat(100, 1),
			isLong:          false,
			expectedPrice:   big.NewRat(110, 1),
		},
		"short position - fillable price higher": {
			bankruptcyPrice: big.NewRat(90, 1),
			fillablePrice:   big.NewRat(100, 1),
			isLong:          false,
			expectedPrice:   big.NewRat(100, 1),
		},
		"short position - prices equal": {
			bankruptcyPrice: big.NewRat(100, 1),
			fillablePrice:   big.NewRat(100, 1),
			isLong:          false,
			expectedPrice:   big.NewRat(100, 1),
		},
		"fractional prices - long position": {
			bankruptcyPrice: big.NewRat(9999, 100),
			fillablePrice:   big.NewRat(10001, 100),
			isLong:          true,
			expectedPrice:   big.NewRat(9999, 100),
		},
		"fractional prices - short position": {
			bankruptcyPrice: big.NewRat(10001, 100),
			fillablePrice:   big.NewRat(9999, 100),
			isLong:          false,
			expectedPrice:   big.NewRat(10001, 100),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := keeper.GetMostAggressivePrice(tc.bankruptcyPrice, tc.fillablePrice, tc.isLong)

			if result.Cmp(tc.expectedPrice) != 0 {
				t.Errorf("Expected price %v, but got %v", tc.expectedPrice, result)
			}
		})
	}
}

func TestRemovePerpetualPosition(t *testing.T) {
	tests := map[string]struct {
		initialPositions    []*satypes.PerpetualPosition
		perpetualIdToRemove uint32
		expectedPositions   []*satypes.PerpetualPosition
	}{
		"remove middle position": {
			initialPositions: []*satypes.PerpetualPosition{
				{PerpetualId: 0, Quantums: dtypes.NewInt(100)},
				{PerpetualId: 1, Quantums: dtypes.NewInt(200)},
				{PerpetualId: 2, Quantums: dtypes.NewInt(300)},
			},
			perpetualIdToRemove: 1,
			expectedPositions: []*satypes.PerpetualPosition{
				{PerpetualId: 0, Quantums: dtypes.NewInt(100)},
				{PerpetualId: 2, Quantums: dtypes.NewInt(300)},
			},
		},
		"remove first position": {
			initialPositions: []*satypes.PerpetualPosition{
				{PerpetualId: 0, Quantums: dtypes.NewInt(100)},
				{PerpetualId: 1, Quantums: dtypes.NewInt(200)},
			},
			perpetualIdToRemove: 0,
			expectedPositions: []*satypes.PerpetualPosition{
				{PerpetualId: 1, Quantums: dtypes.NewInt(200)},
			},
		},
		"remove last position": {
			initialPositions: []*satypes.PerpetualPosition{
				{PerpetualId: 0, Quantums: dtypes.NewInt(100)},
				{PerpetualId: 1, Quantums: dtypes.NewInt(200)},
			},
			perpetualIdToRemove: 1,
			expectedPositions: []*satypes.PerpetualPosition{
				{PerpetualId: 0, Quantums: dtypes.NewInt(100)},
			},
		},
		"remove non-existent position": {
			initialPositions: []*satypes.PerpetualPosition{
				{PerpetualId: 0, Quantums: dtypes.NewInt(100)},
				{PerpetualId: 1, Quantums: dtypes.NewInt(200)},
			},
			perpetualIdToRemove: 2,
			expectedPositions: []*satypes.PerpetualPosition{
				{PerpetualId: 0, Quantums: dtypes.NewInt(100)},
				{PerpetualId: 1, Quantums: dtypes.NewInt(200)},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			subaccount := &satypes.Subaccount{
				PerpetualPositions: tc.initialPositions,
			}
			keeper.RemovePerpetualPosition(subaccount, tc.perpetualIdToRemove)
			require.Equal(t, tc.expectedPositions, subaccount.PerpetualPositions)
		})
	}
}

func TestUpdateTDaiPosition(t *testing.T) {
	tests := map[string]struct {
		subaccount         satypes.Subaccount
		quantumsDelta      *big.Int
		expectedSubaccount satypes.Subaccount
		expectedError      bool
	}{
		"increase TDai position": {
			subaccount: satypes.Subaccount{
				AssetPositions: []*satypes.AssetPosition{
					{AssetId: 0, Quantums: dtypes.NewInt(1000)},
				},
			},
			quantumsDelta: big.NewInt(500),
			expectedSubaccount: satypes.Subaccount{
				AssetPositions: []*satypes.AssetPosition{
					{AssetId: 0, Quantums: dtypes.NewInt(1500)},
				},
			},
		},
		"decrease TDai position": {
			subaccount: satypes.Subaccount{
				AssetPositions: []*satypes.AssetPosition{
					{AssetId: 0, Quantums: dtypes.NewInt(1000)},
				},
			},
			quantumsDelta: big.NewInt(-300),
			expectedSubaccount: satypes.Subaccount{
				AssetPositions: []*satypes.AssetPosition{
					{AssetId: 0, Quantums: dtypes.NewInt(700)},
				},
			},
		},
		"TDai position goes to zero": {
			subaccount: satypes.Subaccount{
				AssetPositions: []*satypes.AssetPosition{
					{AssetId: 0, Quantums: dtypes.NewInt(1000)},
				},
			},
			quantumsDelta: big.NewInt(-1000),
			expectedSubaccount: satypes.Subaccount{
				AssetPositions: []*satypes.AssetPosition{},
			},
		},
		"error: first asset is not TDai": {
			subaccount: satypes.Subaccount{
				AssetPositions: []*satypes.AssetPosition{
					{AssetId: 1, Quantums: dtypes.NewInt(1000)},
				},
			},
			quantumsDelta: big.NewInt(500),
			expectedError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := keeper.UpdateTDaiPosition(&tc.subaccount, tc.quantumsDelta)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedSubaccount, tc.subaccount)
			}
		})
	}
}

func TestLiquidateSubaccountsAgainstOrderbookInternal(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccounts []satypes.Subaccount
		// CLOB state.
		clobs     []types.ClobPair
		feeParams feetypes.PerpetualFeeParams

		existingOrders []types.Order

		MaxLiquidationAttemptsPerBlock         uint32
		MaxIsolatedLiquidationAttemptsPerBlock uint32

		subaccountIds                 *heap.LiquidationPriorityHeap
		isolatedPositionsPriorityHeap *heap.LiquidationPriorityHeap

		expectedSubaccountsToDeleverage       []heap.SubaccountToDeleverage
		expectedSubaccountIds                 *heap.LiquidationPriorityHeap
		expectedIsolatedPositionsPriorityHeap *heap.LiquidationPriorityHeap
		expectedError                         error
		expectPanic                           bool
		ignorePriorityOnSubaccountIds         bool
	}{
		`Can place a liquidation that doesn't match any maker orders`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_49500USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price45000_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         2,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage: []heap.SubaccountToDeleverage{
				{
					SubaccountId: constants.Dave_Num0,
					PerpetualId:  0,
				},
			},
		},
		`Can place a liquidation that matches a maker order`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         2,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage:       nil,
		},
		`Chooses the correct order to liquidate when there are multiple`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_DangerIndex,
				constants.EthUsd_20PercentInitial_10PercentMaintenance_DangerIndex,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc, constants.ClobPair_Eth},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         2,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage: []heap.SubaccountToDeleverage{
				{
					SubaccountId: constants.Dave_Num0,
					PerpetualId:  1,
				},
			},
		},
		`Reinsert subaccount that is still liquidatable after liquidating eth position`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_DangerIndex,
				constants.EthUsd_20PercentInitial_10PercentMaintenance_DangerIndex,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc, constants.ClobPair_Eth},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
				constants.Order_Carl_Num0_Id0_Clob0_Buy_SmallETH_Price3000_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         1,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
			},
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage:       nil,
			ignorePriorityOnSubaccountIds:         true,
		},
		`Too many orders to liquidate, one get deleveraged`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
				constants.Dave_Num1_1BTC_Long_46000USD_Short,
				constants.Dave_Num2_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
				{
					SubaccountId: constants.Dave_Num1,
					Priority:     big.NewFloat(1),
				},
				{
					SubaccountId: constants.Dave_Num2,
					Priority:     big.NewFloat(2),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         2,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num2,
					Priority:     big.NewFloat(2),
				},
			},
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage: []heap.SubaccountToDeleverage{
				{
					SubaccountId: constants.Dave_Num1,
					PerpetualId:  0,
				},
			},
		},
		`Can place a liquidation for an isolated perpetualthat matches a maker order`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_Isolated,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         2,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage:       nil,
		},
		`Can only place one liquidation for an isolated perpetual per block`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_Isolated,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
				constants.Dave_Num1_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price49500_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
				{
					SubaccountId: constants.Dave_Num1,
					Priority:     big.NewFloat(1),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         1,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num1,
					Priority:     big.NewFloat(1),
				},
			},
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage:       nil,
		},
		`Can place two liquidations for an isolated perpetual per block`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_Isolated,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
				constants.Dave_Num1_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price49500_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
				{
					SubaccountId: constants.Dave_Num1,
					Priority:     big.NewFloat(1),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         2,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage:       nil,
		},
		`Can place one isolated and one normal liquidation`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_Isolated,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
				constants.Dave_Num1_1ETH_Long_2900USD_Short,
				constants.Dave_Num2_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc, constants.ClobPair_Eth},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price49500_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
				{
					SubaccountId: constants.Dave_Num2,
					Priority:     big.NewFloat(1),
				},
				{
					SubaccountId: constants.Dave_Num1,
					Priority:     big.NewFloat(2),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         2,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds: heap.NewLiquidationPriorityHeap(),
			expectedIsolatedPositionsPriorityHeap: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num2,
					Priority:     big.NewFloat(1),
				},
			},
			expectedSubaccountsToDeleverage: []heap.SubaccountToDeleverage{
				{
					SubaccountId: constants.Dave_Num1,
					PerpetualId:  1,
				},
			},
		},
		`Can place both isolated and one normal liquidation`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_Isolated,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_46000USD_Short,
				constants.Dave_Num1_1ETH_Long_2900USD_Short,
				constants.Dave_Num2_1BTC_Long_46000USD_Short,
			},
			clobs:     []types.ClobPair{constants.ClobPair_Btc, constants.ClobPair_Eth},
			feeParams: constants.PerpetualFeeParams,

			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id0_Clob0_Buy2BTC_Price49500_GTB10,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Dave_Num0,
					Priority:     big.NewFloat(0),
				},
				{
					SubaccountId: constants.Dave_Num2,
					Priority:     big.NewFloat(1),
				},
				{
					SubaccountId: constants.Dave_Num1,
					Priority:     big.NewFloat(2),
				},
			},
			isolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),

			MaxLiquidationAttemptsPerBlock:         3,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,

			expectedSubaccountIds:                 heap.NewLiquidationPriorityHeap(),
			expectedIsolatedPositionsPriorityHeap: heap.NewLiquidationPriorityHeap(),
			expectedSubaccountsToDeleverage: []heap.SubaccountToDeleverage{
				{
					SubaccountId: constants.Dave_Num1,
					PerpetualId:  1,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				satypes.ModuleAddress,
				authtypes.NewModuleAddress(authtypes.FeeCollectorName),
				mock.Anything,
			).Return(nil)
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				authtypes.NewModuleAddress(satypes.ModuleName),
				perptypes.InsuranceFundModuleAddress,
				mock.Anything,
			).Return(nil)
			// Fee collector does not have any funds.
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				authtypes.NewModuleAddress(authtypes.FeeCollectorName),
				satypes.ModuleAddress,
				mock.Anything,
			).Return(sdkerrors.ErrInsufficientFunds)
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				authtypes.NewModuleAddress(satypes.LiquidityFeeModuleAddress),
				mock.Anything,
			).Return(nil)
			mockBankKeeper.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			// Give the insurance fund a 1M TDai balance.
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				perptypes.InsuranceFundModuleAddress,
				constants.TDai.Denom,
			).Return(
				sdk.NewCoin(
					constants.TDai.Denom,
					sdkmath.NewIntFromBigInt(big.NewInt(1_000_000_000_000)),
				),
			)
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				mock.Anything,
				constants.TDai.Denom,
			).Return(
				sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int))),
			)

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			ctx := ks.Ctx.WithIsCheckTx(true)
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, tc.feeParams))

			// Set up TDai asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Create all CLOBs.
			for _, clobPair := range tc.clobs {
				_, err = ks.ClobKeeper.CreatePerpetualClobPair(
					ctx,
					clobPair.Id,
					clobtest.MustPerpetualId(clobPair),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			// Initialize the liquidations config.
			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default),
			)

			ks.ClobKeeper.Flags.MaxLiquidationAttemptsPerBlock = tc.MaxLiquidationAttemptsPerBlock
			ks.ClobKeeper.Flags.MaxIsolatedLiquidationAttemptsPerBlock = tc.MaxIsolatedLiquidationAttemptsPerBlock

			// Create all existing orders.
			for _, order := range tc.existingOrders {
				_, _, err := ks.ClobKeeper.PlaceShortTermOrder(ctx, &types.MsgPlaceOrder{Order: order})
				require.NoError(t, err)
			}

			if tc.expectPanic {
				require.Panics(t, func() {
					_, err := ks.ClobKeeper.LiquidateSubaccountsAgainstOrderbookInternal(ctx, tc.subaccountIds, tc.isolatedPositionsPriorityHeap)
					require.Error(t, err)
				})
				return
			}
			subaccountsToDeleverage, err := ks.ClobKeeper.LiquidateSubaccountsAgainstOrderbookInternal(ctx, tc.subaccountIds, tc.isolatedPositionsPriorityHeap)
			if tc.expectedError != nil {
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				if tc.ignorePriorityOnSubaccountIds {
					require.Equal(t, tc.expectedSubaccountIds.Len(), tc.subaccountIds.Len(), "Heap lengths should match")
					for i := 0; i < tc.expectedSubaccountIds.Len(); i++ {
						expected := tc.expectedSubaccountIds.PopLowestPriority()
						actual := tc.subaccountIds.PopLowestPriority()
						require.Equal(t, expected.SubaccountId, actual.SubaccountId, "SubaccountIds should match")
					}
				} else {
					require.Equal(t, tc.expectedSubaccountIds, tc.subaccountIds)
				}
				require.Equal(t, tc.expectedIsolatedPositionsPriorityHeap, tc.isolatedPositionsPriorityHeap)
				require.Equal(t, tc.expectedSubaccountsToDeleverage, subaccountsToDeleverage)
			}

		})
	}
}

func TestGetBestPerpetualPositionToLiquidateMultiplePositions(t *testing.T) {
	tests := map[string]struct {
		// Perpetuals state.
		perpetuals []perptypes.Perpetual
		// Subaccount state.
		subaccount satypes.Subaccount

		previouslyLiquidatedPerpetuals []uint32

		expectedPerpetualId uint32
		expectedError       bool
		expectedExactError  string
	}{
		`Expect ETH position to be liquidated first as BTC position is negligeable`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_DangerIndex,
				constants.EthUsd_20PercentInitial_10PercentMaintenance_DangerIndex,
			},
			subaccount: constants.Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short,

			expectedPerpetualId: 1,
		},
		"Expect BTC position to be liquidated when it's the only position": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_DangerIndex,
			},
			subaccount:          constants.Carl_Num0_1BTC_Short,
			expectedPerpetualId: 0,
		},
		"Expect error when all positions have been previously liquidated": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_DangerIndex,
				constants.EthUsd_20PercentInitial_10PercentMaintenance_DangerIndex,
			},
			subaccount:                     constants.Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short,
			previouslyLiquidatedPerpetuals: []uint32{0, 1},
			expectedError:                  true,
			expectedExactError:             "Subaccount has no perpetual positions to liquidate",
		},
		"Expect second position when first has been previously liquidated": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_DangerIndex,
				constants.EthUsd_20PercentInitial_10PercentMaintenance_DangerIndex,
			},
			subaccount:                     constants.Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short,
			previouslyLiquidatedPerpetuals: []uint32{1},
			expectedPerpetualId:            0,
		},
		"Expect error when subaccount has no perpetual positions": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_DangerIndex,
			},
			subaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: []*satypes.AssetPosition{
					{
						AssetId:  0,
						Quantums: dtypes.NewInt(1_000_000_000), // 1,000 TDai
					},
				},
			},
			expectedError:      true,
			expectedExactError: "Subaccount has no perpetual positions to liquidate",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				mock.Anything,
				constants.TDai.Denom,
			).Return(
				sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int))),
			)

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			ctx := ks.Ctx.WithIsCheckTx(true)
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			// Set up TDai asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			// Create all subaccounts.
			ks.SubaccountsKeeper.SetSubaccount(ctx, tc.subaccount)

			// Initialize the liquidations config.
			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default),
			)

			for _, perpId := range tc.previouslyLiquidatedPerpetuals {
				ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(ctx, *tc.subaccount.Id, perpId)
			}

			perpetualId, err := ks.ClobKeeper.GetBestPerpetualPositionToLiquidate(ctx, *tc.subaccount.Id)
			if tc.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedExactError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedPerpetualId, perpetualId)
			}
		})
	}
}

func TestEnsurePerpetualNotAlreadyLiquidated(t *testing.T) {
	tests := map[string]struct {
		perpetuals                     []perptypes.Perpetual
		subaccount                     satypes.Subaccount
		previouslyLiquidatedPerpetuals []uint32
		perpetualIdToCheck             uint32
		expectedError                  error
	}{
		"perpetual not liquidated": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccount:                     constants.Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short,
			previouslyLiquidatedPerpetuals: []uint32{},
			perpetualIdToCheck:             0, // BTC
			expectedError:                  nil,
		},
		"perpetual already liquidated": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccount:                     constants.Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short,
			previouslyLiquidatedPerpetuals: []uint32{0}, // BTC already liquidated
			perpetualIdToCheck:             0,           // BTC
			expectedError:                  types.ErrSubaccountHasLiquidatedPerpetual,
		},
		"different perpetual liquidated": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccount:                     constants.Dave_Num0_TinyBTC_Long_1ETH_Long_2900USD_Short,
			previouslyLiquidatedPerpetuals: []uint32{1}, // ETH already liquidated
			perpetualIdToCheck:             0,           // BTC
			expectedError:                  nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())

			ctx := ks.Ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			// Set up TDai asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			// Create the subaccount.
			ks.SubaccountsKeeper.SetSubaccount(ctx, tc.subaccount)

			// Initialize the liquidations config.
			require.NoError(t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default),
			)

			// Set up previously liquidated perpetuals.
			for _, perpId := range tc.previouslyLiquidatedPerpetuals {
				ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(ctx, *tc.subaccount.Id, perpId)
			}

			// Run the test.
			err = ks.ClobKeeper.EnsurePerpetualNotAlreadyLiquidated(ctx, *tc.subaccount.Id, tc.perpetualIdToCheck)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckInsuranceFundLimits(t *testing.T) {
	tests := map[string]struct {
		perpetuals         []perptypes.Perpetual
		liquidationsConfig types.LiquidationsConfig
		insuranceFundDelta *big.Int
		perpetualId        uint32
		expectedError      error
		expectPanic        bool
	}{
		"success - insurance fund delta within limits": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			liquidationsConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm: 10_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           10_000_000,
					SpreadToMaintenanceMarginRatioPpm: 10_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000),
			},
			insuranceFundDelta: big.NewInt(-500_000),
			perpetualId:        0,
			expectedError:      nil,
		},
		"failure - insurance fund delta exceeds remaining limit": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			liquidationsConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm: 10_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           10_000_000,
					SpreadToMaintenanceMarginRatioPpm: 10_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000),
			},
			insuranceFundDelta: big.NewInt(-1_100_000),
			perpetualId:        0,
			expectedError:      types.ErrLiquidationExceedsMaxInsuranceLost,
		},
		"success - insurance fund delta at limit": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			liquidationsConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm: 10_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           10_000_000,
					SpreadToMaintenanceMarginRatioPpm: 10_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000),
			},
			insuranceFundDelta: big.NewInt(-1_000_000),
			perpetualId:        0,
			expectedError:      nil,
		},
		"success - positive insurance fund delta": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			liquidationsConfig: types.LiquidationsConfig{
				InsuranceFundFeePpm: 10_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           10_000_000,
					SpreadToMaintenanceMarginRatioPpm: 10_000,
				},
				MaxCumulativeInsuranceFundDelta: uint64(1_000_000),
			},
			insuranceFundDelta: big.NewInt(2_000_000),
			perpetualId:        0,
			expectedError:      nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())

			ctx := ks.Ctx.WithIsCheckTx(true)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			// Set up TDai asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			// Set the liquidations config.
			err = ks.ClobKeeper.InitializeLiquidationsConfig(ctx, tc.liquidationsConfig)
			require.NoError(t, err)

			// Run the test.
			if tc.expectPanic {
				require.Panics(t, func() {
					_ = ks.ClobKeeper.CheckInsuranceFundLimits(
						ctx,
						tc.perpetualId,
						tc.insuranceFundDelta,
					)
				})
			} else {
				err := ks.ClobKeeper.CheckInsuranceFundLimits(
					ctx,
					tc.perpetualId,
					tc.insuranceFundDelta,
				)
				if tc.expectedError != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestIsIsolatedPerpetualError_InLiquidateSubaccountsAgainstOrderbookInternal(t *testing.T) {
	tests := map[string]struct {
		perpetuals                             []perptypes.Perpetual
		subaccounts                            []satypes.Subaccount
		subaccountIds                          *heap.LiquidationPriorityHeap
		isolatedPositionsPriorityHeap          *heap.LiquidationPriorityHeap
		MaxLiquidationAttemptsPerBlock         uint32
		MaxIsolatedLiquidationAttemptsPerBlock uint32
		expectedError                          error
	}{
		"Perpetual does not exist and returns err": {
			perpetuals: []perptypes.Perpetual{},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Carl_Num0,
					Priority:     big.NewFloat(0),
				},
			},
			isolatedPositionsPriorityHeap:          heap.NewLiquidationPriorityHeap(),
			MaxLiquidationAttemptsPerBlock:         1,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,
			expectedError:                          errors.New("0: Perpetual does not exist"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(1, 1))

			ctx := ks.Ctx.WithIsCheckTx(true)
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Set up TDai asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					"0/1",
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Initialize the liquidations config.
			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default),
			)

			ks.ClobKeeper.Flags.MaxLiquidationAttemptsPerBlock = tc.MaxLiquidationAttemptsPerBlock
			ks.ClobKeeper.Flags.MaxIsolatedLiquidationAttemptsPerBlock = tc.MaxIsolatedLiquidationAttemptsPerBlock

			_, err = ks.ClobKeeper.LiquidateSubaccountsAgainstOrderbookInternal(ctx, tc.subaccountIds, tc.isolatedPositionsPriorityHeap)
			require.Error(t, err)
			if tc.expectedError != nil {
				require.Contains(t, err.Error(), tc.expectedError.Error())
			}
		})
	}

}

func TestPlacePerpetualLiquidation_InLiquidateSubaccountsAgainstOrderbookInternal(t *testing.T) {
	tests := map[string]struct {
		perpetuals                             []perptypes.Perpetual
		subaccounts                            []satypes.Subaccount
		feeParams                              feetypes.PerpetualFeeParams
		subaccountIds                          *heap.LiquidationPriorityHeap
		isolatedPositionsPriorityHeap          *heap.LiquidationPriorityHeap
		MaxLiquidationAttemptsPerBlock         uint32
		MaxIsolatedLiquidationAttemptsPerBlock uint32
	}{
		"clob does not exists throws an error": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			feeParams: constants.PerpetualFeeParams,
			subaccountIds: &heap.LiquidationPriorityHeap{
				{
					SubaccountId: constants.Carl_Num0,
					Priority:     big.NewFloat(0),
				},
			},
			isolatedPositionsPriorityHeap:          heap.NewLiquidationPriorityHeap(),
			MaxLiquidationAttemptsPerBlock:         1,
			MaxIsolatedLiquidationAttemptsPerBlock: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())

			ctx := ks.Ctx.WithIsCheckTx(true)
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, tc.feeParams))

			// Set up TDai asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					"0/1",
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			// Create all subaccounts.
			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			// Initialize the liquidations config.
			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default),
			)

			ks.ClobKeeper.Flags.MaxLiquidationAttemptsPerBlock = tc.MaxLiquidationAttemptsPerBlock
			ks.ClobKeeper.Flags.MaxIsolatedLiquidationAttemptsPerBlock = tc.MaxIsolatedLiquidationAttemptsPerBlock

			require.Panics(t, func() {
				_, err = ks.ClobKeeper.LiquidateSubaccountsAgainstOrderbookInternal(ctx, tc.subaccountIds, tc.isolatedPositionsPriorityHeap)
			}, "Expected panic did not occur")

			// Check the panic message.
			defer func() {
				if r := recover(); r != nil {
					require.Contains(t, r.(string), "Perpetual ID 0 has no associated CLOB pairs")
				}
			}()
		})
	}

}

func TestGetValidatorAndLiquidityFee(t *testing.T) {
	tests := map[string]struct {
		remainingQuoteQuantums            *big.Int
		expectedValidatorFeeQuoteQuantums *big.Int
		expectedLiquidityFeeQuoteQuantums *big.Int
		expectedError                     error
	}{
		"remaining quote quantums is negative - throws error": {
			remainingQuoteQuantums:            big.NewInt(-1),
			expectedValidatorFeeQuoteQuantums: big.NewInt(0),
			expectedLiquidityFeeQuoteQuantums: big.NewInt(0),
			expectedError:                     errors.New("Remaining quote quantums -1 is negative"),
		},
		"remaining quote quantums is zero - returns zero fees": {
			remainingQuoteQuantums:            big.NewInt(0),
			expectedValidatorFeeQuoteQuantums: big.NewInt(0),
			expectedLiquidityFeeQuoteQuantums: big.NewInt(0),
			expectedError:                     nil,
		},
		"remaining quote quantums is positive - returns fees": {
			remainingQuoteQuantums:            big.NewInt(100),
			expectedValidatorFeeQuoteQuantums: big.NewInt(20),
			expectedLiquidityFeeQuoteQuantums: big.NewInt(80),
			expectedError:                     nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)

			mockBankKeeper := &mocks.BankKeeper{}

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())

			ctx := ks.Ctx.WithIsCheckTx(true)

			require.NoError(
				t,
				ks.ClobKeeper.InitializeLiquidationsConfig(ctx, types.LiquidationsConfig_Default),
			)

			validatorFeeQuoteQuantums, liquidityFeeQuoteQuantums, err := ks.ClobKeeper.GetValidatorAndLiquidityFee(
				ctx,
				tc.remainingQuoteQuantums,
			)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedValidatorFeeQuoteQuantums, validatorFeeQuoteQuantums)
				require.Equal(t, tc.expectedLiquidityFeeQuoteQuantums, liquidityFeeQuoteQuantums)
			}
		})
	}
}

func TestGetInsuranceFundDeltaBlockLimit(t *testing.T) {
	tests := map[string]struct {
		perpetuals                           []perptypes.Perpetual
		feeParams                            feetypes.PerpetualFeeParams
		perpetualId                          uint32
		expectedInsuranceFundDeltaBlockLimit *big.Int
		expectedError                        error
	}{
		"perpetual does not exist - returns error": {
			perpetuals:                           []perptypes.Perpetual{},
			feeParams:                            constants.PerpetualFeeParams,
			perpetualId:                          0,
			expectedInsuranceFundDeltaBlockLimit: big.NewInt(0),
			expectedError:                        errors.New("Perpetual does not exist"),
		},
		"isolated perpetual returns isolated market max cummalitive insurance fund delta per block": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_SmallMarginRequirement_Isolated,
			},
			feeParams:                            constants.PerpetualFeeParams,
			perpetualId:                          0,
			expectedInsuranceFundDeltaBlockLimit: big.NewInt(1_000_000),
			expectedError:                        nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockBankKeeper := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())

			ctx := ks.Ctx.WithIsCheckTx(true)
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, tc.feeParams))

			// Set up TDai asset in assets module.
			err := keepertest.CreateTDaiAsset(ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					"0/1",
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				tc.perpetuals,
			)

			insuranceFundDeltaBlockLimit, err := ks.ClobKeeper.GetInsuranceFundDeltaBlockLimit(
				ctx,
				tc.perpetualId,
			)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedInsuranceFundDeltaBlockLimit, insuranceFundDeltaBlockLimit)
			}

		})
	}
}
