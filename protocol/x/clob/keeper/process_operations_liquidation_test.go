package keeper_test

/*

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testutil_bank "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
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
	tc := processProposerMatchesTestCase{
		perpetuals: []perptypes.Perpetual{
			constants.BtcUsd_100PercentMarginRequirement,
			constants.EthUsd_20PercentInitial_10PercentMaintenance,
		},
		subaccounts: []satypes.Subaccount{
			constants.Carl_Num0_1BTC_Short,
			{
				Id: &constants.Dave_Num0,
				AssetPositions: []*satypes.AssetPosition{
					{
						AssetId:  0,
						Quantums: dtypes.NewInt(-45_001_000_000), // -$45,001
					},
				},
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId: 0,
						Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
					},
					{
						PerpetualId: 1,
						Quantums:    dtypes.NewInt(1000),
					},
				},
			},
		},
		placeOrders: []*types.MsgPlaceOrder{
			{
				Order: types.Order{
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
			},
			{
				Order: types.Order{
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
			},
		},
		clobMatches: []*types.ClobMatch{
			types.NewClobMatchFromMatchPerpetualLiquidation(
				&types.MatchPerpetualLiquidation{
					Liquidated:  constants.Dave_Num0,
					ClobPairId:  0,
					PerpetualId: 0,
					TotalSize:   10,
					IsBuy:       false,
					Fills: []types.MatchPerpetualLiquidation_Fill{
						// Fill would be processed successfully.
						{
							MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
							FillAmount: 10,
						},
					},
				},
			),
			types.NewClobMatchFromMatchPerpetualLiquidation(
				&types.MatchPerpetualLiquidation{
					Liquidated:  constants.Dave_Num0,
					ClobPairId:  1,
					PerpetualId: 1,
					TotalSize:   1000,
					IsBuy:       false,
					Fills: []types.MatchPerpetualLiquidation_Fill{
						// Fill would lead to undercollateralization.
						{
							MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
							FillAmount: 1000,
						},
					},
				},
			),
		},
		clobPairs: []*types.ClobPair{
			&constants.ClobPair_Btc,
			&constants.ClobPair_Eth,
		},
	}

	// Should be the same among all runs.
	var gasConsumed sdk.Gas

	for i := 0; i < 100; i++ {
		ctx, keeper, _, _, mockMsgSender, _, _ := setupProcessProposerMatchesTestCase(t, tc)
		mockMsgSender.On("Enabled").Return(false)

		err := keeper.ProcessProposerMatches(ctx, tc.placeOrders, tc.clobMatches)
		require.ErrorContains(t, err,
			fmt.Sprintf(
				"Subaccount with id {%s 0} failed with UpdateResult: NewlyUndercollateralized:",
				constants.Carl_Num0.Owner,
			),
		)

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
	tests := map[string]processProposerMatchesTestCase{
		"Liquidation succeeds no fills": {
			perpetuals:                        []perptypes.Perpetual{constants.BtcUsd_100PercentMarginRequirement},
			subaccounts:                       []satypes.Subaccount{},
			setupMockBankKeeper:               func(bk *mocks.BankKeeper) {},
			placeOrders:                       []*types.MsgPlaceOrder{},
			clobMatches:                       []*types.ClobMatch{},
			expectedFillAmounts:               map[types.OrderId]satypes.BaseQuantums{},
			expectedPruneableBlockHeights:     map[uint32][]types.OrderId{},
			expectedQuoteBalances:             map[satypes.SubaccountId]*big.Int{},
			expectedPerpetualPositions:        map[satypes.SubaccountId][]*satypes.PerpetualPosition{},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{},
		},
		"Liquidation succeeds when order is completely filled": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(10_000_000)),
				).Return(nil)
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Subaccount pays $250 to insurance fund for liquidating 1 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(250_000_000)),
				).Return(nil).Once()
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $4,749, no taker fees, pays $250 insurance fee
				constants.Carl_Num0: big.NewInt(4_999_000_000 - 250_000_000),
				// $99,990
				constants.Dave_Num0: big.NewInt(100_000_000_000 - 10_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Dave_Num0: {},
				constants.Carl_Num0: {},
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
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(10_100_000)),
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdk.NewIntFromUint64(math.MaxUint64)))
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					types.InsuranceFundName,
					satypes.ModuleName,
					// Insurance fund covers $1 loss for liquidating 1 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(1_000_000)),
				).Return(nil).Once()
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// The subaccount had $50,499 initially, bought 1BTC at $50,500
				// to cover the short position, and received $1 from insurance fund.
				constants.Carl_Num0: big.NewInt(0),
				// $100,489.9
				constants.Dave_Num0: big.NewInt(50_000_000_000 + 50_500_000_000 - 10_100_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Dave_Num0: {},
				constants.Carl_Num0: {},
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_500_000_000, // Liquidated 1BTC at $50,500
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
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(2_500_000)),
				).Return(nil)
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Subaccount pays $62.5 to insurance fund for liquidating 0.25 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(62_500_000)),
				).Return(nil).Twice()
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				},
				{
					Order: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000, // .25 BTC
							},
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId: 25_000_000,
				constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12.OrderId: 25_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {},
				11 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId,
				},
				12 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $29874, no taker fees, pays $125 insurance fee
				constants.Carl_Num0: big.NewInt(29_999_000_000 - 125_000_000),
				// $74,995
				constants.Dave_Num0: big.NewInt(75_000_000_000 - 5_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Dave_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
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
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(2_525_000)),
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdk.NewIntFromUint64(math.MaxUint64)))
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					types.InsuranceFundName,
					satypes.ModuleName,
					// Insurance fund covers $0.25 loss for liquidating 0.25 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(250_000)),
				).Return(nil).Twice()
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
				},
				{
					Order: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000, // .25 BTC
							},
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11.OrderId: 25_000_000,
				constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId: 25_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {},
				11 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11.OrderId,
				},
				12 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// The subaccount had $50,499 initially, bought 0.5BTC at $50,500
				// to cover the short position, and received $0.5 from insurance fund.
				constants.Carl_Num0: big.NewInt(25_249_500_000),
				// $75,244.5
				constants.Dave_Num0: big.NewInt(50_000_000_000 + 25_250_000_000 - 5_050_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Dave_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    25_250_000_000, // Liquidated 0.5 BTC at $50,500
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
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.Anything,
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdk.NewIntFromUint64(math.MaxUint64)))
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Pays insurance fund $0.75 for liquidating 0.75 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(750_000)),
				).Return(nil).Once()
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					types.InsuranceFundName,
					satypes.ModuleName,
					// Insurance fund covers $0.25 loss for liquidating 0.25 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(250_000)),
				).Return(nil).Once()
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Above bankruptcy price.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10,
				},
				{
					// Below bankruptcy price.
					Order: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 75_000_000, // .75 BTC
							},
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10.OrderId:   75_000_000,
				constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId: 25_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				12 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// The subaccount had $50,499 initially, bought 0.75BTC at $50,498
				// and 0.25BTC at $50,500.
				// The subaccount pays $0.5 total to insurance fund.
				constants.Carl_Num0: big.NewInt(0),
				// // $50,000 + (50498 * 0.75 + 50500 * 0.25) * (1 - 0.02%)
				constants.Dave_Num0: big.NewInt(100_488_400_300),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_498_500_000,
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
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.Anything,
				).Return(nil)
				bk.On(
					"GetBalance",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(sdk.NewCoin("USDC", sdk.NewIntFromUint64(math.MaxUint64)))
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Pays insurance fund $0.378735 (capped by MaxLiquidationFeePpm)
					// for liquidating 0.75 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(378_735)),
				).Return(nil).Once()
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
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
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Above bankruptcy price.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10,
				},
				{
					// Below bankruptcy price.
					Order: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 75_000_000, // .75 BTC
							},
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 25_000_000, // .25 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50498_GTB10.OrderId:   75_000_000,
				constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12.OrderId: 25_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				12 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// The subaccount had $50,499 initially, bought 0.75BTC at $50,498
				// and 0.25BTC at $50,500.
				// The subaccount pays $0.5 total to insurance fund.
				constants.Carl_Num0: big.NewInt(0),
				// // $50,000 + (50498 * 0.75 + 50500 * 0.25) * (1 - 0.02%)
				constants.Dave_Num0: big.NewInt(100_488_400_300),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    50_498_500_000,
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
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(5_000_000)),
				).Return(nil)
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Subaccount pays $125 to insurance fund for liquidating 0.5 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(125_000_000)),
				).Return(nil).Once()
			},
			stateFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				{SubaccountId: constants.Dave_Num0, ClientId: 0}: satypes.BaseQuantums(50_000_000),
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 50_000_000, // .50 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 100_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				1000: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $29874, no taker fees, pays $125 insurance fee
				constants.Carl_Num0: big.NewInt(29_999_000_000 - 125_000_000),
				// $74,995
				constants.Dave_Num0: big.NewInt(75_000_000_000 - 5_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Dave_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
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
		"Liquidation succeeds with multiple fills if one order is a replacement and " +
			"the fills are not ordered ascending by GTB": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(1_000_000)),
				).Return(nil)
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Subaccount pays $25 to insurance fund for liquidating 0.1 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(25_000_000)),
				).Return(nil).Twice()
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB12,
				},
				{
					Order: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 10_000_000, // .10 BTC
							},
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 10_000_000, // .10 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11.OrderId: 20_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {},
				12 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB12.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $44949, no taker fees, pays $50 insurance fee
				constants.Carl_Num0: big.NewInt(44_999_000_000 - 50_000_000),
				// $59,998
				constants.Dave_Num0: big.NewInt(60_000_000_000 - 2_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-80_000_000), // .8 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Dave_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(80_000_000), // .8 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
			expectedSubaccountLiquidationInfo: map[satypes.SubaccountId]types.SubaccountLiquidationInfo{
				constants.Carl_Num0: {
					PerpetualsLiquidated:  []uint32{0},
					NotionalLiquidated:    10_000_000_000,
					QuantumsInsuranceLost: 0,
				},
				constants.Dave_Num0: {},
			},
		},
		"Liquidation succeeds with position size smaller than clobPair.MinOrderBaseQuantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			clobPairs: []*types.ClobPair{
				&constants.ClobPair_Btc3,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					AssetPositions: []*satypes.AssetPosition{
						{
							AssetId:  0,
							Quantums: dtypes.NewInt(5_499),
						},
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(-10), // Liquidatable position is smaller than MinOrderBaseQuantums
						},
					},
				},
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
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
				},
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(1)),
				).Return(nil)
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(25)),
				).Return(nil)
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   10,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 10,
							},
						},
					},
				),
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
		"Liquidation succeeds if maker order is reduce-only": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(5_000_000)),
				).Return(nil)
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					types.InsuranceFundName,
					// Subaccount pays $125 to insurance fund for liquidating 0.5 BTC.
					mock.MatchedBy(testutil_bank.MatchUsdcOfAmount(125_000_000)),
				).Return(nil).Once()
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 50_000_000, // .50 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO.OrderId: 50_000_000,
			},
			expectedPruneableBlockHeights: map[uint32][]types.OrderId{
				10 + types.ShortBlockWindow: {
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO.OrderId,
				},
			},
			expectedQuoteBalances: map[satypes.SubaccountId]*big.Int{
				// $29874, no taker fees, pays $125 insurance fee
				constants.Carl_Num0: big.NewInt(29_999_000_000 - 125_000_000),
				// $74,995
				constants.Dave_Num0: big.NewInt(75_000_000_000 - 5_000_000),
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
				constants.Dave_Num0: {
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(50_000_000), // .5 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerMatchSuccessTest(t, tc)
		})
	}
}

func TestProcessProposerMatches_Liquidation_Failure(t *testing.T) {
	tests := map[string]processProposerMatchesTestCase{
		"Liquidation returns error if order quantums is not divisible by StepBaseQuantums": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId:      constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
						Side:         types.Order_SIDE_SELL,
						Quantums:     9, // StepBaseQuantums is 5
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 5,
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId:      constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
						Side:         types.Order_SIDE_SELL,
						Quantums:     10,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 9, // StepBaseQuantums is 5
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
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     10,
						Subticks:     90_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				},
				{
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 1, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     10,
						Subticks:     200_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 20},
					},
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Dave_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   20,
						IsBuy:       false,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							// Fill would be processed successfully.
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 10,
							},
							// Fill would lead to undercollateralization.
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 10,
							},
						},
					},
				),
			},
			expectedError: fmt.Errorf(
				"Subaccount with id {%s 0} failed with UpdateResult: NewlyUndercollateralized",
				constants.Carl_Num0.Owner,
			),
		},
		"Liquidation fails if matches exceed the order quantums when considering state fill amounts": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			stateFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				{SubaccountId: constants.Dave_Num0, ClientId: 0}: satypes.BaseQuantums(50_000_001),
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   50_000_000,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 50_000_000,
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
				constants.Dave_Num0_1BTC_Long,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					mock.Anything,
					authtypes.FeeCollectorName,
					mock.Anything,
				).Return(fmt.Errorf("transfer failed"))
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					mock.Anything,
					types.InsuranceFundName,
					mock.Anything,
				).Return(nil)
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000,
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
		"Returns error when maker order is reduce-only and would increase position size": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_50_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(-100_000_000), // 1 BTC
						},
					},
				},
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000,
							},
						},
					},
				),
			},
			expectedError: types.ErrReduceOnlyWouldIncreasePositionSize,
		},
		"Returns error when maker order is reduce-only and would change position side": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_50_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(99_000_000), // 0.99 BTC
						},
					},
				},
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10_RO,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000,
							},
						},
					},
				),
			},
			expectedError: types.ErrReduceOnlyWouldChangePositionSide,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerMatchFailureTest(t, tc)
		})
	}
}

func TestProcessProposerMatches_Liquidation_Validation_Failure(t *testing.T) {
	tests := map[string]processProposerMatchesTestCase{
		"Stateful order validation: subaccount is not liquidatable": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000, // .25 BTC
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  999,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000, // .25 BTC
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  1,
						PerpetualId: 1,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000,
							},
						},
					},
				),
			},
			expectedError: types.ErrNoOpenPositionForPerpetual,
		},
		"Stateful order validation: size of liquidation order exceeds position size": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   150_000_000, // 1.5 BTC exceeding position size of 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000,
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       false,       // wrong side
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000,
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     1000,
						Subticks:     1000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				},
			},
			clobPairs: []*types.ClobPair{
				&constants.ClobPair_Btc,
				&constants.ClobPair_Eth,
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  1,
						PerpetualId: 0,           // does not match clob pair id 1
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
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
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
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
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Dave_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       false,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000, // .25 BTC
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				},
				{
					Order: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50000_GTB12,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 0,
							},
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 100,
							},
						},
					},
				),
			},
			expectedError: errors.New("fillAmount must be greater than 0"),
		},
		"Stateless match validation: clobPairIds do not match": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 1},
						Side:         types.Order_SIDE_BUY,
						Quantums:     1000,
						Subticks:     1000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				},
			},
			clobPairs: []*types.ClobPair{
				&constants.ClobPair_Btc,
				&constants.ClobPair_Eth,
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,           // Corresponds to ClobPairId 0.
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Alice_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     100_000_000,
						Subticks:     50_000_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     100_000_000,
						Subticks:     1_000_000_000_000, // Maker order selling at $1,000,000, higher than fillable price
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
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
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Carl_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_BUY,
						Quantums:     100_000_000,
						Subticks:     500_000_000, // Maker order buying at $500, lower than fillable price
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Dave_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       false,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50000_GTB11,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 50_000_000, // 0.5 BTC. Too big!
							},
						},
					},
				),
			},
			expectedError: errors.New("Minimum initial order quantums exceeds fill amount"),
		},
		"Position limit: fails when liquidation order size is greater than" +
			" max portion of the position that can be liquidated": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			// Can only liquidate 50% of any position at most.
			liquidationConfig: &constants.LiquidationsConfig_Position_Min10m_Max05mPpm,
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC, liquidating entire position
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 51_000_000, // 0.51 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrLiquidationOrderSizeGreaterThanMax,
		},
		"Position limit: fails when liquidation order size is smaller than min notional liquidated": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			liquidationConfig: &constants.LiquidationsConfig_Position_Min10m_Max05mPpm,
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   10_000, // $5 notional
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 10_000, // $5 notional
							},
						},
					},
				),
			},
			expectedError: types.ErrLiquidationOrderSizeSmallerThanMin,
		},
		"Subaccount block limit: fails when trying to liquidate the same perpetual id": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 50_000_000, // 0.50 BTC
							},
						},
					},
				),
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   50_000_000, // 0.5 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 50_000_000, // 0.50 BTC
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
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000,
				},
			},
			liquidationConfig: &constants.LiquidationsConfig_Subaccount_Max10bNotionalLiquidated_Max10bInsuranceLost,
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 50_000_000, // 0.50 BTC, $25,000 notional
							},
						},
					},
				),
			},
			expectedError: types.ErrLiquidationExceedsSubaccountMaxNotionalLiquidated,
		},
		"Subaccount block limit: fails when a single liquidation fill exceeds max insurance lost block limit": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				},
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
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
		},
		"Subaccount block limit: fails when insurance lost from multiple liquidation fills exceed block limit": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Insurance fund delta is -$0.25.
					Order: constants.Order_Dave_Num0_Id1_Clob0_Sell025BTC_Price50500_GTB11,
				},
				{
					// Insurance fund delta is -$0.25.
					Order: constants.Order_Dave_Num0_Id2_Clob0_Sell025BTC_Price50500_GTB12,
				},
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
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000, // 0.25 BTC
							},
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 25_000_000, // 0.25 BTC
							},
						},
					},
				),
			},
			expectedError: types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
		},
		"Liquidation checks insurance fund delta for individual fills and not the entire liquidation order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Above bankruptcy price.
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 0, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     25_000_000,
						Subticks:     50_498_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 10},
					},
				},
				{
					// Below bankruptcy price.
					Order: types.Order{
						OrderId:      types.OrderId{SubaccountId: constants.Dave_Num0, ClientId: 2, ClobPairId: 0},
						Side:         types.Order_SIDE_SELL,
						Quantums:     75_000_000,
						Subticks:     50_500_000_000,
						GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 12},
					},
				},
			},
			insuranceFundBalance: 10_000_000,
			liquidationConfig: &types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig:  constants.FillablePriceConfig_Default,
				PositionBlockLimits:  constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: types.SubaccountBlockLimits{
					MaxNotionalLiquidated: math.MaxUint64,
					// Max insuracen lost that a subaccount can have is $0.5.
					// For this liquidation, overall insurance fund delta is -$0.5, which is within the limit.
					// but the delta for the second fill is -$0.75, therefore, still considered to be exceeding the limit.
					MaxQuantumsInsuranceLost: 500_000,
				},
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 25_000_000, // .25 BTC, insurance fund delta is $0.25.
							},
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 1},
								FillAmount: 75_000_000, // .75 BTC, insurance fund delta is -$0.75
							},
						},
					},
				),
			},
			expectedError: types.ErrLiquidationExceedsSubaccountMaxInsuranceLost,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerMatchFailureTest(t, tc)
		})
	}
}

func TestValidateProposerMatches_InsuranceFund(t *testing.T) {
	tests := map[string]processProposerMatchesTestCase{
		"Fails when insurance fund is empty": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				},
			},
			insuranceFundBalance: 0, // Insurance fund is empty
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000,
							},
						},
					},
				),
			},
			expectedError: types.ErrInsuranceFundHasInsufficientFunds,
		},
		"Fails when insurance fund is non empty but does not have enough to cover liquidation": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				},
			},
			insuranceFundBalance: 999_999, // Insurance fund only has $0.999999
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000,
							},
						},
					},
				),
			},
			expectedError: types.ErrInsuranceFundHasInsufficientFunds,
		},
		"Fails when insurance fund has enough balance but is less than MaxInsuranceFundQuantumsForDeleveraging": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50499USD,
				constants.Dave_Num0_1BTC_Long,
			},
			placeOrders: []*types.MsgPlaceOrder{
				{
					// Bankruptcy price in quote quantums is $50499 for 1 BTC.
					// When subticks is $50,500, the insurance fund delta is -$1.
					Order: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50500_GTB10,
				},
			},
			insuranceFundBalance: 2_000_000, // Insurance fund has $2
			liquidationConfig: &types.LiquidationsConfig{
				MaxInsuranceFundQuantumsForDeleveraging: 5_000_000, // $5
				MaxLiquidationFeePpm:                    5_000,
				FillablePriceConfig:                     constants.FillablePriceConfig_Default,
				PositionBlockLimits:                     constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits:                   constants.SubaccountBlockLimits_No_Limit,
			},
			clobMatches: []*types.ClobMatch{
				types.NewClobMatchFromMatchPerpetualLiquidation(
					&types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000,
						IsBuy:       true,
						Fills: []types.MatchPerpetualLiquidation_Fill{
							{
								MakerOneof: &types.MatchPerpetualLiquidation_Fill_MakerOrderIndex{MakerOrderIndex: 0},
								FillAmount: 100_000_000,
							},
						},
					},
				),
			},
			expectedError: types.ErrInsuranceFundHasInsufficientFunds,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerMatchFailureTest(t, tc)
		})
	}
}

*/
