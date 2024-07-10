package keeper_test

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"

	storetypes "cosmossdk.io/store/types"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testutil_bank "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Process two matches where the first fill succeeds and the second fails due to
// undercollateralization. Run this 100 times, and verify that the gasConsumed for
// each run is equal.
func TestProcessProposerMatches_Liquidation_Undercollateralized_Determinism(t *testing.T) {
	// TODO(DEC-908): Set up correct `bankKeeper` mock to verify fee transfer.
	tc := processProposerOperationsTestCase{
		perpetuals: []perptypes.Perpetual{
			constants.BtcUsd_100PercentMarginRequirement,
			constants.EthUsd_20PercentInitial_10PercentMaintenance,
		},
		subaccounts: []satypes.Subaccount{
			constants.Carl_Num0_1BTC_Short,
			{
				Id: &constants.Dave_Num0,
				AssetPositions: []*satypes.AssetPosition{
					testutil.CreateSingleAssetPosition(
						0,
						big.NewInt(-45_001_000_000), // -$45,001
					),
				},
				PerpetualPositions: []*satypes.PerpetualPosition{
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(100_000_000), // 1 BTC
						big.NewInt(0),
						big.NewInt(0),
					),
					testutil.CreateSinglePerpetualPosition(
						1,
						big.NewInt(1000),
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
		},
		perpetualFeeParams: &constants.PerpetualFeeParams,
		rawOperations: []types.OperationRaw{
			clobtest.NewShortTermOrderPlacementOperationRaw(
				types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Carl_Num0,
						ClientId:     0,
						ClobPairId:   0,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     10,
					Subticks:     90_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
				},
			),
			clobtest.NewShortTermOrderPlacementOperationRaw(
				types.Order{
					OrderId: types.OrderId{
						SubaccountId: constants.Carl_Num0,
						ClientId:     1,
						ClobPairId:   1,
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     1000,
					Subticks:     200_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
				},
			),
			clobtest.NewMatchOperationRawFromPerpetualLiquidation(
				types.MatchPerpetualLiquidation{
					Liquidated:  constants.Dave_Num0,
					ClobPairId:  0,
					PerpetualId: 0,
					TotalSize:   100_000_000,
					IsBuy:       false,
					Fills: []types.MakerFill{
						// Fill would be processed successfully.
						{
							MakerOrderId: types.OrderId{
								SubaccountId: constants.Carl_Num0,
								ClientId:     0,
								ClobPairId:   0,
							},
							FillAmount: 10,
						},
					},
				},
			),
			clobtest.NewMatchOperationRawFromPerpetualLiquidation(
				types.MatchPerpetualLiquidation{
					Liquidated:  constants.Dave_Num0,
					ClobPairId:  1,
					PerpetualId: 1,
					TotalSize:   1000,
					IsBuy:       false,
					Fills: []types.MakerFill{
						// Fill would lead to undercollateralization.
						{
							MakerOrderId: types.OrderId{
								SubaccountId: constants.Carl_Num0,
								ClientId:     1,
								ClobPairId:   1,
							},
							FillAmount: 1000,
						},
					},
				},
			),
		},
		clobPairs: []types.ClobPair{
			constants.ClobPair_Btc,
			constants.ClobPair_Eth,
		},

		expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
			BlockHeight: 5,
		},
		expectedError: fmt.Errorf(
			"Subaccount with id %v failed with UpdateResult: NewlyUndercollateralized:",
			constants.Carl_Num0,
		),
	}

	// Should be the same among all runs.
	var gasConsumed storetypes.Gas

	for i := 0; i < 100; i++ {
		ctx, _ := runProcessProposerOperationsTestCase(t, tc)

		if i == 0 {
			gasConsumed = ctx.GasMeter().GasConsumed()
		} else {
			require.NotEqual(t,
				0,
				gasConsumed,
			)
			// Assert that gas consumed is the same across all runs.
			require.Equal(t,
				gasConsumed,
				ctx.GasMeter().GasConsumed(),
			)
		}
	}
}

func TestProcessProposerMatches_Liquidation_Success(t *testing.T) {
	blockHeight := uint32(5)
	tests := map[string]processProposerOperationsTestCase{
		"Liquidation succeeds no fills": {
			perpetuals:                 []perptypes.Perpetual{constants.BtcUsd_100PercentMarginRequirement},
			subaccounts:                []satypes.Subaccount{},
			perpetualFeeParams:         &constants.PerpetualFeeParams,
			setupMockBankKeeper:        func(bk *mocks.BankKeeper) {},
			rawOperations:              []types.OperationRaw{},
			expectedFillAmounts:        map[types.OrderId]satypes.BaseQuantums{},
			expectedQuoteBalances:      map[satypes.SubaccountId]int64{},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{},
		},
		"Liquidation succeeds when order is completely filled": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(10_000_000)),
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					perptypes.InsuranceFundModuleAddress,
					// Subaccount pays $250 to insurance fund for liquidating 1 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(250_000_000)),
				).Return(nil).Once()
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
								FillAmount:   100_000_000, // 1 BTC
							},
						},
					},
				),
			},

			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 100_000_000,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $4,749, no taker fees, pays $250 insurance fee
				constants.Carl_Num0: 4_999_000_000 - 250_000_000,
				// $99,990
				constants.Dave_Num0: 100_000_000_000 - 10_000_000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Dave_Num0: {},
				constants.Carl_Num0: {},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_000_000_000, // Liquidated 1BTC at $50,000.
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		"Liquidation succeeds with negative insurance fund delta when order is completely filled": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(10_100_000)),
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(math.MaxUint64)))
				bk.On(
					"SendCoins",
					mock.Anything,
					perptypes.InsuranceFundModuleAddress,
					satypes.ModuleAddress,
					// Insurance fund covers $1 loss for liquidating 1 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(1_000_000)),
				).Return(nil).Once()
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId,
								FillAmount:   100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId: 100_000_000,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// The subaccount had $50,499 initially, bought 1BTC at $50,500
				// to cover the short position, and received $1 from insurance fund.
				constants.Carl_Num0: 0,
				// $100,489.9
				constants.Dave_Num0: 50_000_000_000 + 50_500_000_000 - 10_100_000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Dave_Num0: {},
				constants.Carl_Num0: {},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_000_000_000, // Liquidated 1BTC at $50,000
					QuantumsInsuranceLost: 1_000_000,
				},
				constants.Dave_Num0: {},
			},
		},
		"Liquidation succeeds with multiple partial fills": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(2_500_000)),
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					perptypes.InsuranceFundModuleAddress,
					// Subaccount pays $62.5 to insurance fund for liquidating 0.25 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(62_500_000)),
				).Return(nil).Twice()
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
							{
								MakerOrderId: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId: 25_000_000,
				constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12.OrderId: 25_000_000,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $29874, no taker fees, pays $125 insurance fee
				constants.Carl_Num0: 29_999_000_000 - 125_000_000,
				// $74,995
				constants.Dave_Num0: 75_000_000_000 - 5_000_000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(-50_000_000), // .5 BTC
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Dave_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(50_000_000), // .5 BTC
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId,
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    25_000_000_000, // Liquidated 0.5 BTC at $50,000
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		"Liquidation succeeds with multiple partial fills - negative insurance fund delta": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(2_525_000)),
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(math.MaxUint64)))
				bk.On(
					"SendCoins",
					mock.Anything,
					perptypes.InsuranceFundModuleAddress,
					satypes.ModuleAddress,
					// Insurance fund covers $0.25 loss for liquidating 0.25 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(250_000)),
				).Return(nil).Twice()
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
							{
								MakerOrderId: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11.OrderId: 25_000_000,
				constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId: 25_000_000,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// The subaccount had $50,499 initially, bought 0.5BTC at $50,500
				// to cover the short position, and received $0.5 from insurance fund.
				constants.Carl_Num0: 25_249_500_000,
				// $75,244.5
				constants.Dave_Num0: 50_000_000_000 + 25_250_000_000 - 5_050_000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(-50_000_000), // .5 BTC
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Dave_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(50_000_000), // .5 BTC
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11.OrderId,
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    25_000_000_000, // Liquidated 0.5 BTC at $50,000
					QuantumsInsuranceLost: 500_000,
				},
				constants.Dave_Num0: {},
			},
		},
		"Liquidation succeeds with both positive and negative insurance fund delta": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.Anything,
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(math.MaxUint64)))
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					perptypes.InsuranceFundModuleAddress,
					// Pays insurance fund $0.75 for liquidating 0.75 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(750_000)),
				).Return(nil).Once()
				bk.On(
					"SendCoins",
					mock.Anything,
					perptypes.InsuranceFundModuleAddress,
					satypes.ModuleAddress,
					// Insurance fund covers $0.25 loss for liquidating 0.25 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(250_000)),
				).Return(nil).Once()
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Above bankruptcy price.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10,
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Below bankruptcy price.
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10.OrderId,
								FillAmount:   75_000_000, // .75 BTC
							},
							{
								MakerOrderId: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10.OrderId:   75_000_000,
				constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId: 25_000_000,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// The subaccount had $50,499 initially, bought 0.75BTC at $50,498
				// and 0.25BTC at $50,500.
				// The subaccount pays $0.5 total to insurance fund.
				constants.Carl_Num0: 0,
				// // $50,000 + (50498 * 0.75 + 50500 * 0.25) * (1 - 0.02%)
				constants.Dave_Num0: 100_488_400_300,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10.OrderId,
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_000_000_000,
					QuantumsInsuranceLost: 250_000, // Insurance fund covered $0.25.
				},
				constants.Dave_Num0: {},
			},
		},
		"Insurance fund delta calculation accounts for state changes from previous fills": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.Anything,
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(math.MaxUint64)))
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					perptypes.InsuranceFundModuleAddress,
					// Pays insurance fund $0.378735 (capped by MaxLiquidationFeePpm)
					// for liquidating 0.75 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(378_735)),
				).Return(nil).Once()
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					perptypes.InsuranceFundModuleAddress,
					// Pays insurance fund $0.121265.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(121_265)),
				).Return(nil).Once()
			},
			liquidationConfig: &types.LiquidationsConfig{
				// Cap the max liquidation fee ppm so that the bankruptcy price changes
				// in the insurance fund delta calculation.
				MaxLiquidationFeePpm:  10,
				FillablePriceConfig:   constants.FillablePriceConfig_Default,
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Above bankruptcy price.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10,
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Below bankruptcy price.
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10.OrderId,
								FillAmount:   75_000_000, // .75 BTC
							},
							{
								MakerOrderId: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10.OrderId:   75_000_000,
				constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId: 25_000_000,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// The subaccount had $50,499 initially, bought 0.75BTC at $50,498
				// and 0.25BTC at $50,500.
				// The subaccount pays $0.5 total to insurance fund.
				constants.Carl_Num0: 0,
				// // $50,000 + (50498 * 0.75 + 50500 * 0.25) * (1 - 0.02%)
				constants.Dave_Num0: 100_488_400_300,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10.OrderId,
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_000_000_000,
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		"Liquidation succeeds if matches does not exceed the order quantums when considering state fill amounts": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(5_000_000)),
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					perptypes.InsuranceFundModuleAddress,
					// Subaccount pays $125 to insurance fund for liquidating 0.5 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(125_000_000)),
				).Return(nil).Once()
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					types.OrderId{
						SubaccountId: constants.Dave_Num0, ClientId: 0,
					},
					satypes.BaseQuantums(50_000_000),
					50,
				)
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
								FillAmount:   50_000_000, // .50 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 100_000_000,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $29874, no taker fees, pays $125 insurance fee
				constants.Carl_Num0: 29_999_000_000 - 125_000_000,
				// $74,995
				constants.Dave_Num0: 75_000_000_000 - 5_000_000,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(-50_000_000), // .5 BTC
						big.NewInt(0),
						big.NewInt(0),
					),
				},
				constants.Dave_Num0: {
					testutil.CreateSinglePerpetualPosition(
						0,
						big.NewInt(50_000_000), // .5 BTC
						big.NewInt(0),
						big.NewInt(0),
					),
				},
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    25_000_000_000,
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		"Liquidation succeeds with position size smaller than clobPair.StepBaseQuantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc3,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						testutil.CreateSingleAssetPosition(
							0,
							big.NewInt(5_499),
						),
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(
							0,
							big.NewInt(-10), // Liquidatable position is smaller than StepBaseQuantums
							big.NewInt(0),
							big.NewInt(0),
						),
					},
				},
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(1)),
				).Return(nil)
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					perptypes.InsuranceFundModuleAddress,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(25)),
				).Return(nil)
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId: types.OrderId{
							SubaccountId: constants.Dave_Num0,
							ClientId:     1,
							ClobPairId:   0,
						},
						Side:         types.Order_SIDE_SELL,
						Quantums:     25_000_000,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 11},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   10,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Dave_Num0,
									ClientId:     1,
									ClobPairId:   0,
								},
								FillAmount: 10,
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				{
					SubaccountId: constants.Dave_Num0,
					ClientId:     1,
					ClobPairId:   0,
				}: satypes.BaseQuantums(10),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					{
						SubaccountId: constants.Dave_Num0,
						ClientId:     1,
						ClobPairId:   0,
					},
				},
				BlockHeight: blockHeight,
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    5_000,
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		// TODO(CLOB-824): Re-enable reduce-only tests.
		// "Liquidation succeeds if maker order is reduce-only": {
		// 	perpetuals: []*perptypes.Perpetual{
		// 		&constants.BtcUsd_100PercentMarginRequirement,
		// 	},
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_1BTC_Short_54999USD,
		// 		constants.Dave_Num0_1BTC_Long_50000USD,
		// 	},
		// 	perpetualFeeParams: &constants.PerpetualFeeParams,
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		// 	setupMockBankKeeper: func(bk *mocks.BankKeeper) {
		// 		bk.On(
		// 			"SendCoinsFromModuleToModule",
		// 			mock.Anything,
		// 			satypes.ModuleName,
		// 			authtypes.FeeCollectorName,
		// 			mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(5_000_000)),
		// 		).Return(nil)
		// 		bk.On(
		// 			"SendCoinsFromModuleToModule",
		// 			mock.Anything,
		// 			satypes.ModuleName,
		// 			perptypes.InsuranceFundName,
		// 			// Subaccount pays $125 to insurance fund for liquidating 0.5 BTC.
		// 			mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(125_000_000)),
		// 		).Return(nil).Once()
		// 	},
		// 	rawOperations: []types.OperationRaw{
		// 		clobtest.NewShortTermOrderPlacementOperationRaw(
		// 			constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO,
		// 		),
		// 		clobtest.NewMatchOperationRawFromPerpetualLiquidation(
		// 			types.MatchPerpetualLiquidation{
		// 				Liquidated:  constants.Carl_Num0,
		// 				ClobPairId:  0,
		// 				PerpetualId: 0,
		// 				TotalSize:   100_000_000, // 1 BTC
		// 				IsBuy:       true,
		// 				Fills: []types.MakerFill{
		// 					{
		// 						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO.OrderId,
		// 						FillAmount:   50_000_000, // .50 BTC
		// 					},
		// 				},
		// 			},
		// 		),
		// 	},
		// expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
		// 	constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO.OrderId: 50_000_000,
		// },
		// 	expectedQuoteBalances: map[satypes.SubaccountId]int64{
		// 		// $29874, no taker fees, pays $125 insurance fee
		// 		constants.Carl_Num0: 29_999_000_000 - 125_000_000,
		// 		// $74,995
		// 		constants.Dave_Num0: 75_000_000_000 - 5_000_000,
		// 	},
		// 	expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
		// 		constants.Carl_Num0: {
		// 			{
		// 				PerpetualId:  0,
		// 				Quantums:     dtypes.NewInt(-50_000_000), // .5 BTC
		// 				FundingIndex: dtypes.ZeroInt(),
		// 			},
		// 		},
		// 		constants.Dave_Num0: {
		// 			{
		// 				PerpetualId:  0,
		// 				Quantums:     dtypes.NewInt(50_000_000), // .5 BTC
		// 				FundingIndex: dtypes.ZeroInt(),
		// 			},
		// 		},
		// 	},
		// 	expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
		// 		OrderIdsFilledInLastBlock: []types.OrderId{
		// 			constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO.OrderId,
		// 		},
		// 		BlockHeight: blockHeight,
		// 	},
		// 	expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
		// 		constants.Carl_Num0: {
		// 			PerpetualsLiquidated:  []uint32{0},
		// 			NotionalLiquidated:    25_000_000_000,
		// 			QuantumsInsuranceLost: 0,
		// 		},
		// 		constants.Dave_Num0: {},
		// 	},
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}

func TestProcessProposerMatches_Liquidation_Failure(t *testing.T) {
	tests := map[string]processProposerOperationsTestCase{
		"Liquidation returns error if order quantums is not divisible by StepBaseQuantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
						Side:         types.Order_SIDE_SELL,
						Quantums:     9, // StepBaseQuantums is 5
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
								FillAmount:   5,
							},
						},
					},
				),
			},
			expectedError: errors.New("Order Quantums 9 must be a multiple of the ClobPair's StepBaseQuantums"),
		},
		"Liquidation returns error if fillAmount is not divisible by StepBaseQuantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
								FillAmount:   9, // StepBaseQuantums is 5
							},
						},
					},
				),
			},
			expectedError: types.ErrFillAmountNotDivisibleByStepSize,
		},
		"Liquidation returns error if collateralization check fails with non-success": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_45001USD_Short,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     10,
						Subticks:     90_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 1, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     10,
						Subticks:     200_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Dave_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       false,
						Fills: []types.MakerFill{
							// Fill would be processed successfully.
							{
								MakerOrderId: types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 0, ClobPairId: 0},
								FillAmount:   10,
							},
							// Fill would lead to undercollateralization.
							{
								MakerOrderId: types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 1, ClobPairId: 0},
								FillAmount:   10,
							},
						},
					},
				),
			},
			expectedError: fmt.Errorf(
				"Subaccount with id %v failed with UpdateResult: NewlyUndercollateralized",
				constants.Carl_Num0,
			),
		},
		"Liquidation fails if matches exceed the order quantums when considering state fill amounts": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 0},
					satypes.BaseQuantums(50_000_001),
					50,
				)
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
								FillAmount:   50_000_000,
							},
						},
					},
				),
			},
			expectedError: fmt.Errorf(
				"Match with Quantums 50000000 would exceed total Quantums 100000000 of "+
					"OrderId %v. New total filled quantums would be 100000001",
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
			),
		},
		"Returns error when order filled, subaccounts updated, but transfer to fee module acc failed": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.Anything,
				).Return(fmt.Errorf("transfer failed"))
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					perptypes.InsuranceFundModuleAddress,
					mock.Anything,
				).Return(nil)
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			expectedError: fmt.Errorf(
				"persistMatchedOrders: subaccounts (%v, %v)"+
					" updated, but fee transfer (bigFeeQuoteQuantums: %v) to fee-collector failed. Err: transfer failed:"+
					" Subaccounts updated for a matched order, but fee transfer to fee-collector failed",
				constants.Dave_Num0,
				constants.Carl_Num0,
				10_000_000,
			),
		},
		// TODO(CLOB-824): Re-enable reduce-only tests.
		// "Returns error when maker order is reduce-only and would increase position size": {
		// 	perpetuals: []*perptypes.Perpetual{
		// 		&constants.BtcUsd_100PercentMarginRequirement,
		// 	},
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_1BTC_Short_54999USD,
		// 		{
		// 			Id: &constants.Dave_Num0,
		// 			AssetPositions: []*satypes.AssetPosition{
		// 				&constants.Usdc_Asset_50_000,
		// 			},
		// 			PerpetualPositions: []*satypes.PerpetualPosition{
		// 				{
		// 					PerpetualId: 0,
		// 					Quantums:    dtypes.NewInt(-100_000_000), // 1 BTC
		// 				},
		// 			},
		// 		},
		// 	},
		// 	perpetualFeeParams: &constants.PerpetualFeeParams,
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		// 	rawOperations: []types.OperationRaw{
		// 		clobtest.NewShortTermOrderPlacementOperationRaw(
		// 			constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO,
		// 		),
		// 		clobtest.NewMatchOperationRawFromPerpetualLiquidation(
		// 			types.MatchPerpetualLiquidation{
		// 				Liquidated:  constants.Carl_Num0,
		// 				ClobPairId:  0,
		// 				PerpetualId: 0,
		// 				TotalSize:   100_000_000,
		// 				IsBuy:       true,
		// 				Fills: []types.MakerFill{
		// 					{
		// 						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO.OrderId,
		// 						FillAmount:   100_000_000,
		// 					},
		// 				},
		// 			},
		// 		),
		// 	},
		// 	expectedError: types.ErrReduceOnlyWouldIncreasePositionSize,
		// },
		// "Returns error when maker order is reduce-only and would change position side": {
		// 	perpetuals: []*perptypes.Perpetual{
		// 		&constants.BtcUsd_100PercentMarginRequirement,
		// 	},
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_1BTC_Short_54999USD,
		// 		{
		// 			Id: &constants.Dave_Num0,
		// 			AssetPositions: []*satypes.AssetPosition{
		// 				&constants.Usdc_Asset_50_000,
		// 			},
		// 			PerpetualPositions: []*satypes.PerpetualPosition{
		// 				{
		// 					PerpetualId: 0,
		// 					Quantums:    dtypes.NewInt(99_000_000), // 0.99 BTC
		// 				},
		// 			},
		// 		},
		// 	},
		// 	perpetualFeeParams: &constants.PerpetualFeeParams,
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		// 	rawOperations: []types.OperationRaw{
		// 		clobtest.NewShortTermOrderPlacementOperationRaw(
		// 			constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO,
		// 		),
		// 		clobtest.NewMatchOperationRawFromPerpetualLiquidation(
		// 			types.MatchPerpetualLiquidation{
		// 				Liquidated:  constants.Carl_Num0,
		// 				ClobPairId:  0,
		// 				PerpetualId: 0,
		// 				TotalSize:   100_000_000,
		// 				IsBuy:       true,
		// 				Fills: []types.MakerFill{
		// 					{
		// 						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO.OrderId,
		// 						FillAmount:   100_000_000,
		// 					},
		// 				},
		// 			},
		// 		),
		// 	},
		// 	expectedError: types.ErrReduceOnlyWouldChangePositionSide,
		// },
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}

func TestProcessProposerMatches_Liquidation_Validation_Failure(t *testing.T) {
	tests := map[string]processProposerOperationsTestCase{
		"Stateful order validation: subaccount is not liquidatable": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrSubaccountNotLiquidatable,
		},
		"Stateful order validation: invalid clob": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  999,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrInvalidClob,
		},
		"Stateful order validation: subaccount has no open position for perpetual id": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  1,
						PerpetualId: 1,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
								FillAmount:   25_000_000,
							},
						},
					},
				),
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.ClobKeeper.MustUpdateSubaccountPerpetualLiquidated(ctx, constants.Carl_Num0, 0)
			},
			expectedError: types.ErrNoPerpetualPositionsToLiquidate,
		},
		"Stateful order validation: size of liquidation order exceeds position size": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   150_000_000, // 1.5 BTC exceeding position size of 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
								FillAmount:   25_000_000,
							},
						},
					},
				),
			},
			expectedError: types.ErrInvalidLiquidationOrderTotalSize,
		},
		"Stateful order validation: liquidation order is on the wrong side": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       false,       // wrong side
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
								FillAmount:   25_000_000,
							},
						},
					},
				),
			},
			expectedError: types.ErrInvalidLiquidationOrderSide,
		},
		"Stateful match validation: clob pair and perpetual ids do not match": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     1000,
						Subticks:     1000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  1,
						PerpetualId: 0,           // does not match clob pair id 1
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Alice_Num0,
									ClientId:     0,
									ClobPairId:   0,
								},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrClobPairAndPerpetualDoNotMatch,
		},
		"Stateful match validation: fails if collateralization check does not succeed": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_45001USD_Short,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId: types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 0, ClobPairId: 0},
						Side:    types.Order_SIDE_BUY,
						// Note: This perpetual has a `MaintenanceMargin` of 100%.
						// This account currently has a 1 BTC short, with $100,000 in `QuoteBalance`.
						// If the account loses a single unit of `QuoteBalance`, it will be
						// considered undercollateralized.
						//
						// Making this trade shrinks the account's BTC position by 10 base quantums
						// (5,000 quote quantums worth of BTC) which lowers their maintenance margin
						// from 50,000,000,000 to 49,999,995,000, which means they need at least
						// 99,999,990,000 quote quantums of `QuoteBalance` to remain collateralized.
						//
						// For this reason, we need for this account to spend at least
						// 10,000 quote quantums in order to lower their `QuoteBalance` and bring
						// them under their margin requirement by a single Quote Quantum.
						Quantums:     10,              // 5,000 quote quantums worth of BTC
						Subticks:     100_010_000_000, // Spending 10,001 quote quantums
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Dave_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       false,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Carl_Num0,
									ClientId:     0,
									ClobPairId:   0,
								},
								FillAmount: 10,
							},
						},
					},
				),
			},
			expectedError: fmt.Errorf(
				"Subaccount with id %v failed with UpdateResult: NewlyUndercollateralized",
				constants.Carl_Num0,
			),
		},
		"Stateless match validation: self trade": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId,
								FillAmount:   25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedError: errors.New("Match constitutes a self-trade"),
		},
		"Stateless match validation: fillAmount must be greater than 0": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
								FillAmount:   0,
							},
							{
								MakerOrderId: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12.OrderId,
								FillAmount:   100,
							},
						},
					},
				),
			},
			expectedError: types.ErrFillAmountIsZero,
		},
		"Stateless match validation: clobPairIds do not match": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 1},
						Side:         types.Order_SIDE_BUY,
						Quantums:     1000,
						Subticks:     1000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,           // Corresponds to ClobPairId 0.
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Alice_Num0,
									ClientId:     0,
									ClobPairId:   1,
								},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedError: errors.New("ClobPairIds do not match"),
		},
		"Stateless match validation: maker and taker on the same side": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD, // Buy to cover short position.
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     100_000_000,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Alice_Num0,
									ClientId:     0,
									ClobPairId:   0,
								},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedError: errors.New("Orders are not on opposing sides of the book in match"),
		},
		"Stateless match validation: liquidation buy order doesn't cross with maker sell order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000,
						Subticks:     1_000_000_000_000, // Maker order selling at $1,000,000, higher than fillable price
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Dave_Num0,
									ClientId:     0,
									ClobPairId:   0,
								},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			insuranceFundBalance: math.MaxUint64,
			expectedError:        errors.New("Orders do not cross in match"),
		},
		"Stateless match validation: liquidation sell order doesn't cross with maker buy order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_45001USD_Short,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     100_000_000,
						Subticks:     500_000_000, // Maker order buying at $500, lower than fillable price
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Dave_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       false,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Carl_Num0,
									ClientId:     0,
									ClobPairId:   0,
								},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			insuranceFundBalance: math.MaxUint64,
			expectedError:        errors.New("Orders do not cross in match"),
		},
		"Stateless match validation: minimum initial order quantums exceeds fill amount": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId,
								FillAmount:   50_000_000, // 0.5 BTC. Too big!
							},
						},
					},
				),
			},
			expectedError: errors.New("Minimum initial order quantums exceeds fill amount"),
		},
		// TODO(CLOB-816): validate liquidation order size against liquidation config
		// position limit.
		// "Position limit: fails when liquidation order size is greater than" +
		// 	" max portion of the position that can be liquidated": {
		// 	perpetuals: []*perptypes.Perpetual{
		// 		&constants.BtcUsd_100PercentMarginRequirement,
		// 	},
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_1BTC_Short_54999USD,
		// 		constants.Dave_Num0_1BTC_Long_50000USD,
		// 	},
		// 	perpetualFeeParams: &constants.PerpetualFeeParams,
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		// 	rawOperations: []types.OperationRaw{
		// 		clobtest.NewShortTermOrderPlacementOperationRaw(
		// 			constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
		// 		),
		// 		clobtest.NewMatchOperationRawFromPerpetualLiquidation(
		// 			types.MatchPerpetualLiquidation{
		// 				Liquidated:  constants.Carl_Num0,
		// 				ClobPairId:  0,
		// 				PerpetualId: 0,
		// 				TotalSize:   100_000_000, // 1 BTC, liquidating entire position
		// 				IsBuy:       true,
		// 				Fills: []types.MakerFill{
		// 					{
		// 						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
		// 						FillAmount:   51_000_000, // 0.51 BTC
		// 					},
		// 				},
		// 			},
		// 		),
		// 	},
		// 	// Can only liquidate 50% of any position at most.
		// 	liquidationConfig: &constants.LiquidationsConfig_Position_Min10m_Max05mPpm,
		// 	expectedError:     types.ErrLiquidationOrderSizeGreaterThanMax,
		// },
		// "Position limit: fails when liquidation order size is smaller than min notional liquidated": {
		// 	perpetuals: []*perptypes.Perpetual{
		// 		&constants.BtcUsd_100PercentMarginRequirement,
		// 	},
		// 	subaccounts: []satypes.Subaccount{
		// 		constants.Carl_Num0_1BTC_Short_54999USD,
		// 		constants.Dave_Num0_1BTC_Long_50000USD,
		// 	},
		// 	perpetualFeeParams: &constants.PerpetualFeeParams,
		// 	clobPairs: []types.ClobPair{
		// 		constants.ClobPair_Btc,
		// 	},
		// 	rawOperations: []types.OperationRaw{
		// 		clobtest.NewShortTermOrderPlacementOperationRaw(
		// 			constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
		// 		),
		// 		clobtest.NewMatchOperationRawFromPerpetualLiquidation(
		// 			types.MatchPerpetualLiquidation{
		// 				Liquidated:  constants.Carl_Num0,
		// 				ClobPairId:  0,
		// 				PerpetualId: 0,
		// 				TotalSize:   10_000, // $5 notional
		// 				IsBuy:       true,
		// 				Fills: []types.MakerFill{
		// 					{
		// 						MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
		// 						FillAmount:   10_000, // $5 notional
		// 					},
		// 				},
		// 			},
		// 		),
		// 	},
		// 	liquidationConfig: &constants.LiquidationsConfig_Position_Min10m_Max05mPpm,
		// 	expectedError:     types.ErrLiquidationOrderSizeSmallerThanMin,
		// },
		"Subaccount block limit: fails when trying to liquidate the same perpetual id": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
								FillAmount:   50_000_000, // 0.50 BTC
							},
						},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   50_000_000, // 0.5 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
								FillAmount:   50_000_000, // 0.50 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrSubaccountHasLiquidatedPerpetual,
		},
		"Subaccount block limit: fails when liquidation exceeds subaccount notional amount limit": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000.OrderId,
								FillAmount:   50_000_000, // 0.50 BTC, $25,000 notional
							},
						},
					},
				),
			},
			liquidationConfig: &constants.LiquidationsConfig_Subaccount_Max10bNotionalLiquidated_Max10bInsuranceLost,
			expectedError:     types.ErrInvalidLiquidationOrderTotalSize,
		},
		"Subaccount block limit: fails when a single liquidation fill exceeds max insurance lost block limit": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId,
								FillAmount:   100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    math.MaxUint64,
					MaxQuantumsInsuranceLost: 999_999, // $0.999999
				},
			},
			insuranceFundBalance: math.MaxUint64,
			expectedError:        types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
		},
		"Subaccount block limit: fails when insurance lost from multiple liquidation fills exceed block limit": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Insurance fund delta is -$0.25.
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Insurance fund delta is -$0.25.
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11.OrderId,
								FillAmount:   25_000_000, // 0.25 BTC
							},
							{
								MakerOrderId: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId,
								FillAmount:   25_000_000, // 0.25 BTC
							},
						},
					},
				),
			},
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated:    math.MaxUint64,
					MaxQuantumsInsuranceLost: 499_999, // $0.499999
				},
			},
			insuranceFundBalance: math.MaxUint64,
			expectedError:        types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
		},
		"Liquidation checks insurance fund delta for individual fills and not the entire liquidation order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Above bankruptcy price.
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     25_000_000,
						Subticks:     50_498_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Below bankruptcy price.
					types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 2, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     75_000_000,
						Subticks:     50_500_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 12},
					},
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Dave_Num0,
									ClientId:     0,
									ClobPairId:   0,
								},
								FillAmount: 25_000_000, // .25 BTC, insurance fund delta is $0.25.
							},
							{
								MakerOrderId: types.OrderId{
									SubaccountId: constants.Dave_Num0,
									ClientId:     2,
									ClobPairId:   0,
								},
								FillAmount: 75_000_000, // .75 BTC, insurance fund delta is -$0.75
							},
						},
					},
				),
			},
			insuranceFundBalance: 10_000_000,
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated: math.MaxUint64,
					// Max insurance lost that a subaccount can have is $0.5.
					// For this liquidation, overall insurance fund delta is -$0.5, which is within the limit.
					// but the delta for the second fill is -$0.75, therefore, still considered to be exceeding the limit.
					MaxQuantumsInsuranceLost: 500_000,
				},
			},
			expectedError: types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}

func TestValidateProposerMatches_InsuranceFund(t *testing.T) {
	tests := map[string]processProposerOperationsTestCase{
		"Fails when insurance fund is empty": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId,
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			insuranceFundBalance: 0, // Insurance fund is empty
			expectedError:        types.ErrInsuranceFundHasInsufficientFunds,
		},
		"Fails when insurance fund is non empty but does not have enough to cover liquidation": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId,
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			insuranceFundBalance: 999_999, // Insurance fund only has $0.999999
			expectedError:        types.ErrInsuranceFundHasInsufficientFunds,
		},
		"Succeeds when insurance fund has enough balance": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				),
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId,
								FillAmount:   100_000_000,
							},
						},
					},
				),
			},
			insuranceFundBalance: 2_000_000, // Insurance fund has $2
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm:  5_000,
				FillablePriceConfig:   constants.FillablePriceConfig_Default,
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},
			expectedError: nil,
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: 5,
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}
