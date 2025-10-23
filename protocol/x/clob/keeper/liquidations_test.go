package keeper_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
			// Give the insurance fund a 1M USDC balance.
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				perptypes.InsuranceFundModuleAddress,
				constants.Usdc.Denom,
			).Return(
				sdk.NewCoin(
					constants.Usdc.Denom,
					sdkmath.NewIntFromBigInt(big.NewInt(1_000_000_000_000)),
				),
			)

			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())

			ctx := ks.Ctx.WithIsCheckTx(true).WithBlockTime(time.Unix(5, 0))
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, tc.feeParams))

			// Set up USDC asset in assets module.
			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
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
				_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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
				msg := &types.MsgPlaceOrder{Order: order}

				txBuilder := constants.TestEncodingCfg.TxConfig.NewTxBuilder()
				err := txBuilder.SetMsgs(msg)
				require.NoError(t, err)
				bytes, err := constants.TestEncodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
				require.NoError(t, err)
				ctx = ctx.WithTxBytes(bytes)

				_, _, err = ks.ClobKeeper.PlaceShortTermOrder(ctx, msg)
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
			ctx := ks.Ctx.WithIsCheckTx(true).WithBlockTime(time.Unix(5, 0))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper)
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
				)
				require.NoError(t, err)
			}

			clobPair := constants.ClobPair_Btc
			_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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
	}{
		`PlacePerpetualLiquidation succeeds with pre-existing liquidations in the block`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(54_999_000_000), // $54,999
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-100_000_000), // -1 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
						testutil.CreateSinglePerpetualPosition(
							1,
							big.NewInt(-1_000_000_000), // -1 ETH
							big.NewInt(0),
							big.NewInt(0),
						),
					},
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
					PerpetualsLiquidated:  []uint32{1, 0},
					NotionalLiquidated:    53_000_000_000, // $53,000
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		`PlacePerpetualLiquidation considers pre-existing liquidations and stops before exceeding
		max notional liquidated per block`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(54_999_000_000), // $54,999
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-100_000_000), // -1 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
						testutil.CreateSinglePerpetualPosition(
							1,
							big.NewInt(-1_000_000_000), // -1 ETH
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    10_000_000_000, // $10,000
					MaxQuantumsInsuranceLost: math.MaxUint64,
				},
			},
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000,
				&constants.LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000,
				&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,

			// Only matches one order since matching both orders would exceed `MaxNotionalLiquidated`.
			expectedOrderStatus: types.LiquidationExceededSubaccountMaxNotionalLiquidated,
			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id3_Clob1_Sell1ETH_Price3000,
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
					PerpetualsLiquidated:  []uint32{1, 0},
					NotionalLiquidated:    3_000_000_000, // $3,000
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		`PlacePerpetualLiquidation matches some order and stops before exceeding max notional liquidated per block`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    20_000_000_000, // $20,000
					MaxQuantumsInsuranceLost: math.MaxUint64,
				},
			},
			placedMatchableOrders: []types.MatchableOrder{
				&constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				&constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12,
			},
			order: constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,

			// Only matches one order since matching both orders would exceed `MaxNotionalLiquidated`.
			expectedOrderStatus: types.LiquidationExceededSubaccountMaxNotionalLiquidated,
			expectedPlacedOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    12_500_000_000, // $12,500
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		`PlacePerpetualLiquidation considers pre-existing liquidations and stops before exceeding
		max insurance fund lost per block`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(53_000_000_000), // $53,000
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-100_000_000), // -1 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
						testutil.CreateSinglePerpetualPosition(
							1,
							big.NewInt(-1_000_000_000), // -1 ETH
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    math.MaxUint64,
					MaxQuantumsInsuranceLost: 50_000_000, // $50
				},
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
					PerpetualsLiquidated:  []uint32{1, 0},
					NotionalLiquidated:    3_000_000_000, // $3,000
					QuantumsInsuranceLost: 30_000_000,    // $30
				},
				constants.Dave_Num0: {},
			},
		},
		`PlacePerpetualLiquidation matches some order and stops before exceeding max insurance lost per block`: {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},

			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    math.MaxUint64,
					MaxQuantumsInsuranceLost: 500_000, // $0.5
				},
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    12_500_000_000, // $12,500
					QuantumsInsuranceLost: 250_000,
				},
				constants.Dave_Num0: {},
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
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(0))) // Insurance fund is empty.
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    0,
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
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
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(0))) // Insurance fund is empty.
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    0,
					QuantumsInsuranceLost: 0,
				},
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
					mock.Anything,
					mock.Anything,
				).Return(
					// Insurance fund has $0.99 initially.
					sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(990_000)),
				).Once()
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(
					// Insurance fund has $0.74 after covering the loss of the first match.
					sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(740_000)),
				).Twice()
			},

			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm:  5_000,
				FillablePriceConfig:   constants.FillablePriceConfig_Default,
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    12_500_000_000, // $12,500
					QuantumsInsuranceLost: 250_000,
				},
				constants.Dave_Num0: {},
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
					mock.Anything,
					mock.Anything,
				).Return(
					// Insurance fund has $0.99 initially.
					sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(990_000)),
				).Once()
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(
					// Insurance fund has $0.74 after covering the loss of the first match.
					sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(740_000)),
				).Once()
			},

			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm:  5_000,
				FillablePriceConfig:   constants.FillablePriceConfig_Default,
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    12_500_000_000, // $12,500
					QuantumsInsuranceLost: 250_000,
				},
			},
		},
		`PlacePerpetualLiquidation panics when trying to liquidate the same perpetual in a block`: {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(54_999_000_000), // $54,999
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-100_000_000), // -1 BTC
							big.NewInt(0),
							big.NewInt(0),
						),
						testutil.CreateSinglePerpetualPosition(
							1,
							big.NewInt(-2_000_000_000), // -2 ETH
							big.NewInt(0),
							big.NewInt(0),
						),
					},
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
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(math.MaxUint64)))
			}

			mockIndexerEventManager := &mocks.IndexerEventManager{}
			mockIndexerEventManager.On("Enabled").Return(false)
			ks := keepertest.NewClobKeepersTestContext(t, memclob, bankKeeper, mockIndexerEventManager)

			ctx := ks.Ctx.WithIsCheckTx(true).WithBlockTime(time.Unix(5, 0))

			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParams))

			// Set up USDC asset in assets module.
			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
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
						constants.BtcUsd_100PercentMarginRequirement.Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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
						constants.EthUsd_100PercentMarginRequirement.Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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
					msg := &types.MsgPlaceOrder{Order: order}

					txBuilder := constants.TestEncodingCfg.TxConfig.NewTxBuilder()
					err := txBuilder.SetMsgs(msg)
					require.NoError(t, err)
					bytes, err := constants.TestEncodingCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
					require.NoError(t, err)
					ctx = ctx.WithTxBytes(bytes)

					_, _, err = ks.ClobKeeper.PlaceShortTermOrder(ctx, msg)
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

				// Verify test expectations.
				// TODO(DEC-1979): Refactor these tests to support the operations queue refactor.
				// placedOrders, matchedOrders := memclob.GetPendingFills(ctx)

				// require.Equal(t, tc.expectedPlacedOrders, placedOrders, "Placed orders lists are not equal")
				// require.Equal(t, tc.expectedMatchedOrders, matchedOrders, "Matched orders lists are not equal")
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_000_000_000, // $50,000
					QuantumsInsuranceLost: 0,
				},
			},
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    12_500_000_000, // $12,500
					QuantumsInsuranceLost: 0,
				},
			},
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    0, // $0
					QuantumsInsuranceLost: 0,
				},
			},
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_500_000_000 / 4,
					QuantumsInsuranceLost: 0,
				},
			},
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    0,
					QuantumsInsuranceLost: 0,
				},
			},
			expectedSubaccounts: []satypes.Subaccount{
				// Deleveraging fails.
				// Dave's bankruptcy price to close 1 BTC long is $50,000, and deleveraging can not be
				// performed due to non overlapping bankruptcy prices.
				constants.Carl_Num0_1BTC_Short_49999USD,
				constants.Dave_Num0_1BTC_Long_50000USD_Short,
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    12_500_000_000,
					QuantumsInsuranceLost: 0,
				},
			},
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
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+12_499_750_000),
						),
					},
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    0,
					QuantumsInsuranceLost: 0,
				},
			},
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
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+24_999_500_000),
						),
					},
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_500_000_000 / 4,
					QuantumsInsuranceLost: 0,
				},
			},
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_500_000_000,
					QuantumsInsuranceLost: 750_000,
				},
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(50_000_000_000+50_499_500_000),
						),
					},
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
				// Current maxLiquidationFeePpm = 5000.
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
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_000_000_000 / 4,
					QuantumsInsuranceLost: 0,
				},
			},
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(
								54_999_000_000-50_000_000_000/4-
									lib.BigIntMulPpm(
										big.NewInt(50_000_000_000/4),
										constants.LiquidationsConfig_No_Limit.MaxLiquidationFeePpm,
									).Int64(),
							),
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
							big.NewInt(50_000_000_000+50_000_000_000/4),
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
				mock.Anything,
				mock.Anything,
			).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(tc.insuranceFundBalance))).Twice()

			mockIndexerEventManager := &mocks.IndexerEventManager{}
			mockIndexerEventManager.On("Enabled").Return(false)
			ks := keepertest.NewClobKeepersTestContext(t, memclob, bankKeeper, mockIndexerEventManager)

			ctx := ks.Ctx.WithIsCheckTx(true).WithBlockTime(time.Unix(5, 0))

			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParamsNoFee))

			// Set up USDC asset in assets module.
			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
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

			ks.ClobKeeper.DaemonLiquidationInfo.UpdateSubaccountsWithPositions(
				clobtest.GetOpenPositionsFromSubaccounts(tc.subaccounts),
				uint32(ctx.BlockHeight()),
			)

			for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {
				err := ks.PricesKeeper.UpdateMarketPrices(
					ctx,
					[]*pricestypes.MsgUpdateMarketPrices_MarketPrice{
						{
							MarketId: marketId,
							Price:    oraclePrice,
						},
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
							perpetuals[i].Params.DefaultFundingPpm,
						),
					),
				).Once().Return()
				_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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

	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, indexerEventManager)
	ctx := ks.Ctx.WithTxBytes(constants.TestTxBytes)
	// CheckTx mode set correctly
	ctx = ctx.WithIsCheckTx(true).WithBlockTime(time.Unix(5, 0))

	ks.MarketMapKeeper.InitGenesis(ks.Ctx, constants.MarketMap_DefaultGenesisState)
	prices.InitGenesis(ctx, *ks.PricesKeeper, constants.Prices_DefaultGenesisState)
	perpetuals.InitGenesis(ctx, *ks.PerpetualsKeeper, constants.Perpetuals_DefaultGenesisState)

	memClob.On("CreateOrderbook", constants.ClobPair_Btc).Return()
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
				constants.Perpetuals_DefaultGenesisState.Perpetuals[0].Params.DefaultFundingPpm,
			),
		),
	).Once().Return()
	_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * 1),
			),
		},
		"Subaccount with no open positions but negative net collateral is not liquidatable": {
			expectedIsLiquidatable: false,
			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -1),
			),
		},
		"Subaccount at initial margin requirements is not liquidatable": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			perpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(
					uint32(0),
					big.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
					big.NewInt(0),
					big.NewInt(0),
				),
			},
			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_000),
			),
			expectedIsLiquidatable: false,
		},
		"Subaccount below initial but at maintenance margin requirements is not liquidatable": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			perpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(
					uint32(0),
					big.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
					big.NewInt(0),
					big.NewInt(0),
				),
			},
			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_500),
			),
			expectedIsLiquidatable: false,
		},
		"Subaccount below maintenance margin requirements is liquidatable": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			perpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(
					uint32(0),
					big.NewInt(10_000_000), // 0.1 BTC, $5,000 notional.
					big.NewInt(0),
					big.NewInt(0),
				),
			},
			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			expectedIsLiquidatable: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			require.NoError(t, keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper))

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
							AssetUpdates: testutil.CreateUsdcAssetUpdates(bankruptcyPriceInQuoteQuantums),
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

func TestGetFillablePrice(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		perpetualId   uint32
		deltaQuantums int64

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
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           2_000_000,
					SpreadToMaintenanceMarginRatioPpm: 100_000,
				},
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},
			// $49,998 = (49,998 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,999.8 notional sell order.
			expectedFillablePrice: big.NewRat(49_998, 100),
		},
		`Can calculate fillable price for a subaccount with one long position when 
		spreadToMaintenanceMarginRatioPpm is 200_000`: {
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           lib.OneMillion,
					SpreadToMaintenanceMarginRatioPpm: 200_000,
				},
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},
			// $49,998 = (49,998 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,999.8 notional sell order.
			expectedFillablePrice: big.NewRat(49_998, 100),
		},
		`Can calculate fillable price for a subaccount with one short position that is slightly
		below maintenance margin requirements`: {
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_499),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           2_000_000,
					SpreadToMaintenanceMarginRatioPpm: 100_000,
				},
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			// $50,002 = (50,002 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC short with a $5,000.2 notional buy order.
			expectedFillablePrice: big.NewRat(50_002, 100),
		},
		`Can calculate fillable price for a subaccount with one short position when 
		SpreadToMaintenanceMarginRatioPpm is 200_000`: {
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * 5_499),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           lib.OneMillion,
					SpreadToMaintenanceMarginRatioPpm: 200_000,
				},
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			// $50,002 = (50,002 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC short with a $5,000.2 notional buy order.
			expectedFillablePrice: big.NewRat(50_002, 100),
		},
		"Can calculate fillable price for a subaccount with one long position at the bankruptcy price": {
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			perpetualId:   0,
			deltaQuantums: -5_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			perpetualId:   0,
			deltaQuantums: 10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			perpetualId:   1,
			deltaQuantums: -100_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           math.MaxUint32,
					SpreadToMaintenanceMarginRatioPpm: 100_000,
				},
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			// $49,500 = (495 / 1) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,950 notional sell order.
			expectedFillablePrice: big.NewRat(495, 1),
		},
		`Can calculate fillable price when SpreadTomaintenanceMarginRatioPpm is 1`: {
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           lib.OneMillion,
					SpreadToMaintenanceMarginRatioPpm: 1,
				},
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			expectedFillablePrice: big.NewRat(4_999_999_999, 10_000_000),
		},
		`Can calculate fillable price when SpreadTomaintenanceMarginRatioPpm is one million`: {
			perpetualId:   0,
			deltaQuantums: -10_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -4_501),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm:           lib.OneMillion,
					SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
				},
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			// $49,990 = (49990 / 100) subticks * 10^(QuoteCurrencyAtomicResolution - BaseCurrencyAtomicResolution).
			// This means we should close the 0.1 BTC long with a $4,999 notional sell order.
			expectedFillablePrice: big.NewRat(49_990, 100),
		},
		`Returns error when deltaQuantums is zero`: {
			perpetualId:   0,
			deltaQuantums: 0,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			assetPositions: testutil.CreateUsdcAssetPositions(
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
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

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
				big.NewInt(tc.deltaQuantums),
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

func TestGetLiquidationInsuranceFundDelta(t *testing.T) {
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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
		},
		`Fully closing one long position above the bankruptcy price pays max liquidation fee 
		when MaxLiquidationFeePpm is 25_000`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    56_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm:  25_000,
				FillablePriceConfig:   constants.FillablePriceConfig_Default,
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, 5,610,000,000 quote quantums.
			// abs(5,610,000,000) * 2.5% max liquidation fee < 5,610,000,000 - 5,100,000,000, so the max
			// liquidation fee is returned.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(140_250_000),
		},
		`Fully closing one long position above the bankruptcy price pays less than max liquidation fee 
		when MaxLiquidationFeePpm is one million`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    56_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm:  1_000_000,
				FillablePriceConfig:   constants.FillablePriceConfig_Default,
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, 5,610,000,000 quote quantums.
			// abs(5,610,000,000) * 100% max liquidation fee > 5,610,000,000 - 5,100,000,000, so all
			// of the leftover collateral is transferred to the insurance fund.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(510_000_000),
		},
		`Fully closing one short position above the bankruptcy price and pays max liquidation fee`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    44_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
		},
		`Fully closing one short position above the bankruptcy price and pays max liquidation fee
		when MaxLiquidationFeePpm is 25_000`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    44_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm:  25_000,
				FillablePriceConfig:   constants.FillablePriceConfig_Default,
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, -4,410,000,000 quote quantums.
			// abs(-4,410,000,000) * 2.5% max liquidation fee < -4,900,000,000 - -4,410,000,000, so
			// the max liquidation fee is returned.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(110_250_000),
		},
		`Fully closing one short position above the bankruptcy price and pays less than max liquidation fee
		when MaxLiquidationFeePpm is one million`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    44_100_000_000, // 10% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm:  1_000_000,
				FillablePriceConfig:   constants.FillablePriceConfig_Default,
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 10% above bankruptcy price, -4,410,000,000 quote quantums.
			// abs(-4,410,000,000) * 100% max liquidation fee > -4,900,000,000 - -4,410,000,000, so all
			// of the leftover collateral is transferred to the insurance fund.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(490_000_000),
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
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
		},
		`Fully closing one long position at the bankruptcy price and the delta is 0`: {
			perpetualId: 0,
			isBuy:       false,
			fillAmount:  10_000_000,     // -0.1 BTC delta.
			subticks:    51_000_000_000, // 0% above bankruptcy price (equal).

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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
		},
		`Fully closing one short position above the bankruptcy price and the delta is 0`: {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  10_000_000,     // 0.1 BTC delta.
			subticks:    49_000_000_000, // 0% above bankruptcy price.

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * -5_100),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},

			// Bankruptcy price in quote quantums is 5,100,000,000 quote quantums.
			// Liquidation price is 1% below the bankruptcy price, 5,049,000,000 quote quantums.
			// 5,049,000,000 - 5,100,000,000 < 0, so the insurance fund must cover the losses.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(-51_000_000),
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

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Liquidation price is 1% below the bankruptcy price, -4,949,000,000 quote quantums.
			// -4,949,000,000 - -4,900,000,000 < 0, so the insurance fund msut cover the losses.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(-49_000_000),
		},
		"Returns error when delta quantums is zero": {
			perpetualId: 0,
			isBuy:       true,
			fillAmount:  0,
			subticks:    50_000_000_000,

			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},

			assetPositions: testutil.CreateUsdcAssetPositions(
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

			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(constants.QuoteBalance_OneDollar * 4_900),
			),
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},

			// Bankruptcy price in quote quantums is -4,900,000,000 quote quantums.
			// Insurance fund delta before applying position limit is 0 - -4,900,000,000 = 4,900,000,000.
			// abs(0) * 0.5% max liquidation fee < 4,900,000,000, so overall delta is zero.
			expectedLiquidationInsuranceFundDeltaBig: big.NewInt(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

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
						tc.perpetuals[0].Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
			_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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
			liquidationInsuranceFundDeltaBig, err := ks.ClobKeeper.GetLiquidationInsuranceFundDelta(
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
			}
		})
	}
}

func TestConvertFillablePriceToSubticks(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		fillablePrice     *big.Rat
		isLiquidatingLong bool
		clobPair          types.ClobPair

		// Expectations.
		expectedSubticks types.Subticks
	}{
		`Converts fillable price to subticks for liquidating a BTC long position`: {
			fillablePrice: big.NewRat(
				int64(constants.FiveBillion),
				1,
			),
			isLiquidatingLong: true,
			clobPair:          constants.ClobPair_Btc,

			expectedSubticks: 500_000_000_000_000_000,
		},
		`Converts fillable price to subticks for liquidating a BTC short position`: {
			fillablePrice: big.NewRat(
				int64(constants.FiveBillion),
				1,
			),
			isLiquidatingLong: false,
			clobPair:          constants.ClobPair_Btc,

			expectedSubticks: 500_000_000_000_000_000,
		},
		`Converts fillable price to subticks for liquidating a long position and rounds up`: {
			fillablePrice: big.NewRat(
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
		`Converts fillable price to subticks for liquidating a short position and rounds down`: {
			fillablePrice: big.NewRat(
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
		`Converts fillable price to subticks for liquidating a short position and rounds down, but
		the result is lower bounded at SubticksPerTick`: {
			fillablePrice: big.NewRat(
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
		`Converts zero fillable price to subticks for liquidating a short position and rounds down,
		but the result is lower bounded at SubticksPerTick`: {
			fillablePrice: big.NewRat(
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
		`Converts fillable price to subticks for liquidating a long position and rounds up, but
		the result is upper bounded at the max Uint64 that is most aligned with SubticksPerTick`: {
			fillablePrice: big_testutil.MustFirst(
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
			subticks := ks.ClobKeeper.ConvertFillablePriceToSubticks(
				ks.Ctx,
				tc.fillablePrice,
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

func TestConvertFillablePriceToSubticks_PanicsOnNegativeFillablePrice(t *testing.T) {
	// Setup keeper state.
	memClob := memclob.NewMemClobPriceTimePriority(false)
	ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

	// Run the test.
	require.Panics(t, func() {
		ks.ClobKeeper.ConvertFillablePriceToSubticks(
			ks.Ctx,
			big.NewRat(-1, 1),
			false,
			constants.ClobPair_Btc,
		)
	})
}

func TestGetPerpetualPositionToLiquidate(t *testing.T) {
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
		`Full position size is returned when MinPositionNotionalLiquidated is greater than position size`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   10_000_000_000,
					MaxPositionPortionLiquidatedPpm: lib.OneMillion,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(-10_000_000),
		},
		`Half position size is returned when MaxPositionPortionLiquidatedPpm is 500,000`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   1_000,
					MaxPositionPortionLiquidatedPpm: 500_000,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(-5_000_000),
		},
		`full position is returned when position size is smaller than StepBaseQuantums`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(
					0,
					big.NewInt(5),
					big.NewInt(0),
					big.NewInt(0),
				),
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   1,
					MaxPositionPortionLiquidatedPpm: 100_000,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			clobPairs: []types.ClobPair{
				// StepBaseQuantums is 10.
				constants.ClobPair_Btc3,
			},

			expectedClobPair: constants.ClobPair_Btc3,
			expectedQuantums: new(big.Int).SetInt64(-5),
		},
		`returned position size is rounded down to the nearest clob.stepBaseQuantums`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(
					0,
					big.NewInt(140),
					big.NewInt(0),
					big.NewInt(0),
				),
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   1,
					MaxPositionPortionLiquidatedPpm: 100_000,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			clobPairs: []types.ClobPair{
				// StepBaseQuantums is 5.
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			// 140 * 10% = 14, which is rounded down to 10.
			expectedQuantums: new(big.Int).SetInt64(-10),
		},
		`returned position size is at least clob.stepBaseQuantums`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(
					0,
					big.NewInt(20),
					big.NewInt(0),
					big.NewInt(0),
				),
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   1,
					MaxPositionPortionLiquidatedPpm: 100_000,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			clobPairs: []types.ClobPair{
				// StepBaseQuantums is 5.
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			// 20 * 10% = 2, however, clobPair.StepBaseQuantum is 5,
			// so the returned position size is 5.
			expectedQuantums: new(big.Int).SetInt64(-5),
		},
		`Full position is returned when position smaller than subaccount limit`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong, // 0.1 BTC, $5,000 notional
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    10_000_000_000, // $10,000
					MaxQuantumsInsuranceLost: math.MaxUint64,
				},
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(-10_000_000), // -0.1 BTC
		},
		`Max subaccount limit is returned when position larger than subaccount limit`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong, // 0.1 BTC, $5,000 notional
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    2_500_000_000, // $2,500
					MaxQuantumsInsuranceLost: math.MaxUint64,
				},
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(-5_000_000), // -0.05 BTC
		},
		`position size is capped by subaccount block limit when subaccount limit is lower than 
		position block limit`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   1_000,
					MaxPositionPortionLiquidatedPpm: 500_000,
				},
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    2_000_000_000, // $2,000
					MaxQuantumsInsuranceLost: math.MaxUint64,
				},
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(-4_000_000), // capped by subaccount block limit
		},
		`position size is capped by position block limit when position limit is lower than 
		subaccount block limit`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCLong,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   1_000,
					MaxPositionPortionLiquidatedPpm: 400_000, // 40%
				},
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    2_500_000_000, // $2,500
					MaxQuantumsInsuranceLost: math.MaxUint64,
				},
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(-4_000_000), // capped by position block limit
		},
		`Result is rounded to nearest step size`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(
					0,
					big.NewInt(21),
					big.NewInt(0),
					big.NewInt(0),
				),
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   1_000,
					MaxPositionPortionLiquidatedPpm: 500_000,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			clobPairs: []types.ClobPair{
				{
					Metadata: &types.ClobPair_PerpetualClobMetadata{
						PerpetualClobMetadata: &types.PerpetualClobMetadata{
							PerpetualId: 0,
						},
					},
					Status:                    types.ClobPair_STATUS_ACTIVE,
					StepBaseQuantums:          3, // step size is 3
					SubticksPerTick:           100,
					QuantumConversionExponent: -8,
				},
			},

			expectedClobPair: types.ClobPair{
				Id: 0,
				Metadata: &types.ClobPair_PerpetualClobMetadata{
					PerpetualClobMetadata: &types.PerpetualClobMetadata{
						PerpetualId: 0,
					},
				},
				Status:                    types.ClobPair_STATUS_ACTIVE,
				StepBaseQuantums:          3, // step size is 3
				SubticksPerTick:           100,
				QuantumConversionExponent: -8,
			},
			expectedQuantums: new(big.Int).SetInt64(-9), // result is rounded down
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
		`Full position size (short) is returned when MinPositionNotionalLiquidated is 
		greater than position size`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   10_000_000_000,
					MaxPositionPortionLiquidatedPpm: lib.OneMillion,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(10_000_000),
		},
		`Half position size (short) is returned when MaxPositionPortionLiquidatedPpm is 500,000`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   1_000,
					MaxPositionPortionLiquidatedPpm: 500_000,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(5_000_000),
		},
		`Full position (short) is returned when position smaller than subaccount limit`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort, // 0.1 BTC, $5,000 notional
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    10_000_000_000, // $10,000
					MaxQuantumsInsuranceLost: math.MaxUint64,
				},
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(10_000_000), // 0.1 BTC
		},
		`Max subaccount limit is returned when short position larger than subaccount limit`: {
			perpetualPositions: []*satypes.PerpetualPosition{
				&constants.PerpetualPosition_OneTenthBTCShort, // 0.1 BTC, $5,000 notional
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    2_500_000_000, // $2,500
					MaxQuantumsInsuranceLost: math.MaxUint64,
				},
			},

			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedClobPair: constants.ClobPair_Btc,
			expectedQuantums: new(big.Int).SetInt64(5_000_000), // 0.05 BTC
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
				new(big.Int).SetUint64(6148914691236517000),
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
					new(big.Int).SetString("-6148914691236517000", 10),
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
							tc.perpetuals[perpetualId].Params.DefaultFundingPpm,
						),
					),
				).Once().Return()
				_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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

			perpetualId, err := ks.ClobKeeper.GetPerpetualPositionToLiquidate(
				ks.Ctx,
				*subaccount.Id,
			)
			require.NoError(t, err)

			deltaQuantums, err := ks.ClobKeeper.GetLiquidatablePositionSizeDelta(
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
		`Does not place a liquidation order for a non-liquidatable subaccount`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
			},
			clobs: []types.ClobPair{constants.ClobPair_Btc},
			existingOrders: []types.Order{
				constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
			},

			liquidatableSubaccount: constants.Carl_Num0,

			expectedErr:           types.ErrSubaccountNotLiquidatable,
			expectedPlacedOrders:  []*types.MsgPlaceOrder{},
			expectedMatchedOrders: []*types.ClobMatch{},
		},
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
				"SendCoinsFromModuleToModule",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(nil)
			// Give the insurance fund a 1M USDC balance.
			mockBankKeeper.On(
				"GetBalance",
				mock.Anything,
				perptypes.InsuranceFundModuleAddress,
				constants.Usdc.Denom,
			).Return(
				sdk.NewCoin(
					constants.Usdc.Denom,
					sdkmath.NewIntFromBigInt(big.NewInt(1_000_000_000_000)),
				),
			)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())
			ctx := ks.Ctx.WithIsCheckTx(true).WithBlockTime(time.Unix(5, 0))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParams))

			err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
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
				_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
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

func TestGetMaxLiquidatableNotionalAndInsuranceLost(t *testing.T) {
	tests := map[string]struct {
		// Setup
		liquidationConfig               types.LiquidationsConfig
		previouslyLiquidatedPerpetualId uint32
		previousNotionalLiquidated      *big.Int
		previousInsuranceFundLost       *big.Int

		// Expectations.
		expectedMaxNotionalLiquidatablePanic bool
		expectedMaxNotionalLiquidatableErr   error
		expectedMaxInsuranceLostPanic        bool
		expectedMaxInsuranceLostErr          error
		expectedMaxNotionalLiquidatable      *big.Int
		expectedMaxInsuranceLost             *big.Int
	}{
		"Can get max notional liquidatable and insurance lost": {
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    150,
					MaxQuantumsInsuranceLost: 150,
				},
			},
			previouslyLiquidatedPerpetualId: uint32(1),
			previousNotionalLiquidated:      big.NewInt(100),
			previousInsuranceFundLost:       big.NewInt(-100),

			expectedMaxNotionalLiquidatable: big.NewInt(50),
			expectedMaxInsuranceLost:        big.NewInt(50),
		},
		"Same perpetual id": {
			liquidationConfig:          constants.LiquidationsConfig_No_Limit,
			previousNotionalLiquidated: big.NewInt(100),
			previousInsuranceFundLost:  big.NewInt(-100),

			expectedMaxNotionalLiquidatableErr: types.ErrSubaccountHasLiquidatedPerpetual,
			expectedMaxInsuranceLostErr:        types.ErrSubaccountHasLiquidatedPerpetual,
		},
		"invalid notional liquidated": {
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    50,
					MaxQuantumsInsuranceLost: 150,
				},
			},
			previouslyLiquidatedPerpetualId: uint32(1),
			previousNotionalLiquidated:      big.NewInt(100),
			previousInsuranceFundLost:       big.NewInt(-100),

			expectedMaxInsuranceLost:             big.NewInt(50),
			expectedMaxNotionalLiquidatablePanic: true,
			expectedMaxNotionalLiquidatableErr: errorsmod.Wrapf(
				types.ErrLiquidationExceedsSubaccountMaxNotionalLiquidated,
				"Subaccount %+v notional liquidated exceeds block limit. Current notional liquidated: %v, block limit: %v",
				constants.Alice_Num0,
				100,
				50,
			),
		},
		"invalid insurance lost": {
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    150,
					MaxQuantumsInsuranceLost: 50,
				},
			},
			previouslyLiquidatedPerpetualId: uint32(1),
			previousNotionalLiquidated:      big.NewInt(100),
			previousInsuranceFundLost:       big.NewInt(-100),

			expectedMaxNotionalLiquidatable: big.NewInt(50),
			expectedMaxInsuranceLostPanic:   true,
			expectedMaxInsuranceLostErr: errorsmod.Wrapf(
				types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
				"Subaccount %+v insurance lost exceeds block limit. Current insurance lost: %v, block limit: %v",
				constants.Alice_Num0,
				100,
				50,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

			err := ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, tc.liquidationConfig)
			require.NoError(t, err)

			subaccountId := constants.Alice_Num0
			perpetualId := uint32(0)
			ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(
				ks.Ctx,
				subaccountId,
				tc.previouslyLiquidatedPerpetualId,
			)
			ks.ClobKeeper.UpdateSubaccountLiquidationInfo(
				ks.Ctx,
				subaccountId,
				tc.previousNotionalLiquidated,
				tc.previousInsuranceFundLost,
			)

			if tc.expectedMaxNotionalLiquidatablePanic {
				require.PanicsWithError(
					t,
					tc.expectedMaxNotionalLiquidatableErr.Error(),
					func() {
						//nolint: errcheck
						ks.ClobKeeper.GetSubaccountMaxNotionalLiquidatable(
							ks.Ctx,
							subaccountId,
							perpetualId,
						)
					},
				)
			} else {
				actualMaxNotionalLiquidatable, err := ks.ClobKeeper.GetSubaccountMaxNotionalLiquidatable(
					ks.Ctx,
					subaccountId,
					perpetualId,
				)
				if tc.expectedMaxNotionalLiquidatableErr != nil {
					require.ErrorContains(t, err, tc.expectedMaxNotionalLiquidatableErr.Error())
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedMaxNotionalLiquidatable, actualMaxNotionalLiquidatable)
				}
			}

			if tc.expectedMaxInsuranceLostPanic {
				require.PanicsWithError(
					t,
					tc.expectedMaxInsuranceLostErr.Error(),
					func() {
						//nolint: errcheck
						ks.ClobKeeper.GetSubaccountMaxInsuranceLost(
							ks.Ctx,
							subaccountId,
							perpetualId,
						)
					},
				)
			} else {
				actualMaxInsuranceLost, err := ks.ClobKeeper.GetSubaccountMaxInsuranceLost(
					ks.Ctx,
					subaccountId,
					perpetualId,
				)
				if tc.expectedMaxInsuranceLostErr != nil {
					require.ErrorContains(t, err, tc.expectedMaxInsuranceLostErr.Error())
				} else {
					require.NoError(t, err)
					require.Equal(t, tc.expectedMaxInsuranceLost, actualMaxInsuranceLost)
				}
			}
		})
	}
}

func TestGetMaxAndMinPositionNotionalLiquidatable(t *testing.T) {
	tests := map[string]struct {
		// Setup
		liquidationConfig   types.LiquidationsConfig
		positionToLiquidate *satypes.PerpetualPosition

		// Expectations.
		expectedErr                        error
		expectedMinPosNotionalLiquidatable *big.Int
		expectedMaxPosNotionalLiquidatable *big.Int
	}{
		"Can get min notional liquidatable": {
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   100,
					MaxPositionPortionLiquidatedPpm: lib.OneMillion,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},
			positionToLiquidate: testutil.CreateSinglePerpetualPosition(
				uint32(0),
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
			expectedMinPosNotionalLiquidatable: big.NewInt(100),
			expectedMaxPosNotionalLiquidatable: big.NewInt(50_000_000_000), // $50,000
		},
		"Can get max notional liquidatable": {
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   100,
					MaxPositionPortionLiquidatedPpm: 500_000,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},
			positionToLiquidate: testutil.CreateSinglePerpetualPosition(
				uint32(0),
				big.NewInt(100_000_000), // 1 BTC
				big.NewInt(0),
				big.NewInt(0),
			),
			expectedMinPosNotionalLiquidatable: big.NewInt(100),
			expectedMaxPosNotionalLiquidatable: big.NewInt(25_000_000_000), // $25,000
		},
		"min and max notional liquidatable can be overridden": {
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits: types.PositionBlockLimits{
					MinPositionNotionalLiquidated:   10_000_000, // $10
					MaxPositionPortionLiquidatedPpm: lib.OneMillion,
				},
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},
			positionToLiquidate: testutil.CreateSinglePerpetualPosition(
				uint32(0),
				big.NewInt(10_000), // $5 notional
				big.NewInt(0),
				big.NewInt(0),
			),
			expectedMinPosNotionalLiquidatable: big.NewInt(5_000_000), // $5
			expectedMaxPosNotionalLiquidatable: big.NewInt(5_000_000), // $5
		},
		"errors are propagated": {
			liquidationConfig: constants.LiquidationsConfig_No_Limit,
			positionToLiquidate: testutil.CreateSinglePerpetualPosition(
				uint32(999), // non-existent
				big.NewInt(1),
				big.NewInt(0),
				big.NewInt(0),
			),
			expectedErr: perptypes.ErrPerpetualDoesNotExist,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Create perpetual.
			_, err := ks.PerpetualsKeeper.CreatePerpetual(
				ks.Ctx,
				constants.BtcUsd_100PercentMarginRequirement.Params.Id,
				constants.BtcUsd_100PercentMarginRequirement.Params.Ticker,
				constants.BtcUsd_100PercentMarginRequirement.Params.MarketId,
				constants.BtcUsd_100PercentMarginRequirement.Params.AtomicResolution,
				constants.BtcUsd_100PercentMarginRequirement.Params.DefaultFundingPpm,
				constants.BtcUsd_100PercentMarginRequirement.Params.LiquidityTier,
				constants.BtcUsd_100PercentMarginRequirement.Params.MarketType,
			)
			require.NoError(t, err)

			// Create all CLOBs.
			mockIndexerEventManager.On("AddTxnEvent",
				ks.Ctx,
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
						constants.BtcUsd_100PercentMarginRequirement.Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ks.Ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)

			err = ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, tc.liquidationConfig)
			require.NoError(t, err)

			actualMinPosNotionalLiquidatable,
				actualMaxPosNotionalLiquidatable,
				err := ks.ClobKeeper.GetMaxAndMinPositionNotionalLiquidatable(
				ks.Ctx,
				tc.positionToLiquidate,
			)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedMinPosNotionalLiquidatable, actualMinPosNotionalLiquidatable)
				require.Equal(t, tc.expectedMaxPosNotionalLiquidatable, actualMaxPosNotionalLiquidatable)
			}
		})
	}
}

func TestSortLiquidationOrders(t *testing.T) {
	tests := map[string]struct {
		orders   []types.LiquidationOrder
		expected []types.LiquidationOrder
	}{
		"Sorts liquidations by abs percentage difference from oracle price (long vs long)": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50501_01,
			},
			expected: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50501_01,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500,
			},
		},
		"Sorts liquidations by abs percentage difference from oracle price (short vs short)": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price49500,
			},
			expected: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price49500,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
			},
		},
		"Sorts liquidations by abs percentage difference from oracle price (long vs short)": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price49500,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50501_01,
			},
			expected: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50501_01,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price49500,
			},
		},
		"Sorts liquidations by order size in quote quantums (long vs long)": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy01BTC_Price50000,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,
			},
			expected: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy01BTC_Price50000,
			},
		},
		"Sorts liquidations by order size in quote quantums (short vs short)": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
			},
			expected: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
				constants.LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000,
			},
		},
		"Sorts liquidations by order size in quote quantums (long vs short)": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,
			},
			expected: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,
				constants.LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000,
			},
		},
		"Sorts liquidations by order hash": {
			orders: []types.LiquidationOrder{
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,
			},
			expected: []types.LiquidationOrder{
				constants.LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000,
				constants.LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Create perpetual.
			_, err := ks.PerpetualsKeeper.CreatePerpetual(
				ks.Ctx,
				constants.BtcUsd_100PercentMarginRequirement.Params.Id,
				constants.BtcUsd_100PercentMarginRequirement.Params.Ticker,
				constants.BtcUsd_100PercentMarginRequirement.Params.MarketId,
				constants.BtcUsd_100PercentMarginRequirement.Params.AtomicResolution,
				constants.BtcUsd_100PercentMarginRequirement.Params.DefaultFundingPpm,
				constants.BtcUsd_100PercentMarginRequirement.Params.LiquidityTier,
				constants.BtcUsd_100PercentMarginRequirement.Params.MarketType,
			)
			require.NoError(t, err)

			// Create all CLOBs.
			mockIndexerEventManager.On("AddTxnEvent",
				ks.Ctx,
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
						constants.BtcUsd_100PercentMarginRequirement.Params.DefaultFundingPpm,
					),
				),
			).Once().Return()
			_, err = ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
				ks.Ctx,
				constants.ClobPair_Btc.Id,
				clobtest.MustPerpetualId(constants.ClobPair_Btc),
				satypes.BaseQuantums(constants.ClobPair_Btc.StepBaseQuantums),
				constants.ClobPair_Btc.QuantumConversionExponent,
				constants.ClobPair_Btc.SubticksPerTick,
				constants.ClobPair_Btc.Status,
			)
			require.NoError(t, err)

			err = ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, types.LiquidationsConfig_Default)
			require.NoError(t, err)

			ks.ClobKeeper.SortLiquidationOrders(
				ks.Ctx,
				tc.orders,
			)
			require.Equal(t, tc.expected, tc.orders)
		})
	}
}
