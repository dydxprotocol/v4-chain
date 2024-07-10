package keeper_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testutil_bank "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
)

func TestProcessProposerMatches_LongTerm_Success(t *testing.T) {
	blockHeight := uint32(5)
	tests := map[string]processProposerOperationsTestCase{
		"Succeeds with new maker Long-Term order": {
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
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 100_000_000,
				// Fully filled orders are removed.
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $99,975
				constants.Dave_Num0: 100_000_000_000 - 25_000_000,
				// $49,990
				constants.Carl_Num0: 50_000_000_000 - 10_000_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with new taker Long-Term order": {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 100_000_000,
				// Fully filled orders are removed.
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $99,990
				constants.Dave_Num0: 100_000_000_000 - 10_000_000,
				// $49,975
				constants.Carl_Num0: 50_000_000_000 - 25_000_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with existing maker Long-Term order": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 100_000_000,
				// Fully filled orders are removed.
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $99,975
				constants.Dave_Num0: 100_000_000_000 - 25_000_000,
				// $49,990
				constants.Carl_Num0: 50_000_000_000 - 10_000_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with existing taker Long-Term order": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 100_000_000,
				// Fully filled orders are removed.
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $99,990
				constants.Dave_Num0: 100_000_000_000 - 10_000_000,
				// $49,975
				constants.Carl_Num0: 50_000_000_000 - 25_000_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with new maker and taker Long-Term orders completely filled": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							25_000_000+10_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   100_000_000, // 1 BTC,
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				// Fully filled orders are removed.
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 0,
				// Fully filled orders are removed.
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $99,975
				constants.Dave_Num0: 100_000_000_000 - 25_000_000,
				// $49,990
				constants.Carl_Num0: 50_000_000_000 - 10_000_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with new maker and taker Long-Term orders partially filled": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   50_000_000, // 0.5 BTC,
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 50_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId:  50_000_000,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $74,987.5
				constants.Dave_Num0: 75_000_000_000 - 12_500_000,
				// $74,995
				constants.Carl_Num0: 75_000_000_000 - 5_000_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
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
		},
		"Succeeds with Long-Term order and multiple fills": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Twice()
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   50_000_000, // 0.5 BTC,
						},
					},
				),
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   50_000_000, // 0.5 BTC,
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC.OrderId: 50_000_000,
				constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000.OrderId:           50_000_000,
				// Fully filled orders are removed.
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 0,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $99,990
				constants.Dave_Num0: 100_000_000_000 - 10_000_000,
				// $49,975
				constants.Carl_Num0: 50_000_000_000 - 25_000_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.Order_Carl_Num0_Id0_Clob0_Buy05BTC_Price50000_GTB10_IOC.OrderId,
					constants.Order_Carl_Num0_Id2_Clob0_Buy05BTC_Price50000.OrderId,
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
			},
			expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
				constants.Carl_Num0: {},
				constants.Dave_Num0: {},
			},
		},
		"Succeeds with new maker Long-Term order in liquidation match": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
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
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(
					ks.Ctx, &blocktimetypes.BlockInfo{
						Timestamp: time.Unix(5, 0),
					},
				)
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
								FillAmount:   100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				// Fully filled orders are removed.
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 0,
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
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
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
		"Succeeds with existing maker Long-Term order in liquidation match": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
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
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(
					ks.Ctx, &blocktimetypes.BlockInfo{
						Timestamp: time.Unix(5, 0),
					},
				)
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								MakerOrderId: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
								FillAmount:   100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				// Fully filled orders are removed.
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 0,
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
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
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
		"Succeeds with maker Long-Term order when considering state fill amount": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					50_000_000,
					math.MaxUint32,
				)
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   50_000_000,
						},
					},
				),
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 50_000_000,
				// Fully filled orders are removed.
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $74,975
				constants.Dave_Num0: 75_000_000_000 - 12_500_000,
				// $74,990
				constants.Carl_Num0: 75_000_000_000 - 5_000_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
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
		},
		"Succeeds with taker Long-Term order when considering state fill amount": {
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
			preExistingStatefulOrders: []types.Order{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoins",
					mock.Anything,
					satypes.ModuleAddress,
					authtypes.NewModuleAddress(authtypes.FeeCollectorName),
					mock.MatchedBy(
						testutil_bank.MatchUsdcOfAmount(
							12_500_000+5_000_000,
						),
					),
				).Return(nil).Once()
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
				ks.ClobKeeper.SetOrderFillAmount(
					ctx,
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					50_000_000,
					math.MaxUint32,
				)
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
							FillAmount:   50_000_000,
						},
					},
				),
			},

			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId: 50_000_000,
				// Fully filled orders are removed.
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
			},
			expectedQuoteBalances: map[satypes.SubaccountId]int64{
				// $74,990
				constants.Dave_Num0: 75_000_000_000 - 5_000_000,
				// $74,975
				constants.Carl_Num0: 75_000_000_000 - 12_500_000,
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				OrderIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId,
				},
				RemovedStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
				},
				BlockHeight: blockHeight,
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
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}
