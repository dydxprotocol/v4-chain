package keeper_test

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	testutil_bank "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
)

// func TestStatefulValidateProposedOperations(t *testing.T) {
// 	// TODO(DEC-908): Set up correct `bankKeeper` mock to verify fee transfer.
// 	tests := map[string]struct {
// 		// Setup.
// 		internalOperations []types.InternalOperation
// 		// orders
// 		statefulOrders []types.Order

// 		expectedErr error
// 	}{
// 		"test1": {},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			memclob := memclob.NewMemClobPriceTimePriority(true)
// 			ks := keepertest.NewClobKeepersTestContext(t, memclob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})
// 			ctx := ks.Ctx.WithBlockHeight(int64(5))

// 			for _, order := tc.statefulOrders {

// 			}
// 			require.Equal(t, tc.expectedProcessProposerMatchesEvents, processProposerMatchesEvents)
// 		})
// 	}
// }

func TestProcessProposerMatches_LongTerm_StatefulValidation_Failure(t *testing.T) {
	tests := map[string]processProposerOperationsTestCase{
		`Stateful order validation: referenced maker order does not exist in state`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
			expectedError: sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful long term order id %+v does not exist in state.",
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced taker order does not exist in state`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful long term order id %+v does not exist in state.",
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced maker order in liquidation match does not exist in state`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								// Maker order is a long-term order.
								MakerOrderId: constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
								FillAmount:   100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful long term order id %+v does not exist in state.",
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced long-term order is on the wrong side`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
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
							MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrInvalidMatchOrder,
				"Taker Order %+v and Maker order %+v are not on opposing sides of the book",
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderTextString(),
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.GetOrderTextString(),
			),
		},
		`Stateful order validation: referenced long-term order is for the wrong clob pair`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
				&constants.EthUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			},
			preExistingStatefulOrders: []types.Order{
				{
					OrderId: types.OrderId{
						SubaccountId: constants.Carl_Num0,
						ClientId:     0,
						OrderFlags:   types.OrderIdFlags_LongTerm,
						ClobPairId:   1, // ETH.
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
				},
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// This is a BTC order.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							MakerOrderId: types.OrderId{
								SubaccountId: constants.Carl_Num0,
								ClientId:     0,
								OrderFlags:   types.OrderIdFlags_LongTerm,
								ClobPairId:   1, // ETH.
							},
							FillAmount: 100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: errors.New("ClobPairIds do not match in match"),
		},
		"Fails with Long-Term order when considering state fill amount": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
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
					50_000_001,
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
			expectedError: fmt.Errorf(
				"Match with Quantums 50000000 would exceed total Quantums 100000000 of "+
					"OrderId %v. New total filled quantums would be 100000001",
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}

func TestProcessProposerMatches_Conditional_Validation_Failure(t *testing.T) {
	tests := map[string]processProposerOperationsTestCase{
		`Stateful order validation: referenced maker order does not exist in state`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced taker order does not exist in state`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
					constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced maker order in liquidation match does not exist in state`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
				clobtest.NewMatchOperationRawFromPerpetualLiquidation(
					types.MatchPerpetualLiquidation{
						Liquidated:  constants.Carl_Num0,
						ClobPairId:  0,
						PerpetualId: 0,
						TotalSize:   100_000_000, // 1 BTC
						IsBuy:       true,
						Fills: []types.MakerFill{
							{
								// Maker order is a conditional order.
								MakerOrderId: constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
								FillAmount:   100_000_000, // 1 BTC
							},
						},
					},
				),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced maker order exist in state but is untriggered`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
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
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced conditional order is on the wrong side`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			triggeredConditionalOrders: []types.Order{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
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
							MakerOrderId: constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrInvalidMatchOrder,
				"Taker Order %+v and Maker order %+v are not on opposing sides of the book",
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.GetOrderTextString(),
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.GetOrderTextString(),
			),
		},
		`Stateful order validation: referenced conditional order is for the wrong clob pair`: {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			triggeredConditionalOrders: []types.Order{
				{
					OrderId: types.OrderId{
						SubaccountId: constants.Carl_Num0,
						ClientId:     0,
						OrderFlags:   types.OrderIdFlags_Conditional,
						ClobPairId:   1, // ETH.
					},
					Side:         types.Order_SIDE_BUY,
					Quantums:     100_000_000,
					Subticks:     50_000_000_000,
					GoodTilOneof: &types.Order_GoodTilBlockTime{GoodTilBlockTime: 10},
				},
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewShortTermOrderPlacementOperationRaw(
					// This is a BTC order.
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					[]types.MakerFill{
						{
							MakerOrderId: types.OrderId{
								SubaccountId: constants.Carl_Num0,
								ClientId:     0,
								OrderFlags:   types.OrderIdFlags_Conditional,
								ClobPairId:   1, // ETH.
							},
							FillAmount: 100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: errors.New("ClobPairIds do not match in match"),
		},
		"Fails with conditional order when considering state fill amount": {
			perpetuals: []*perptypes.Perpetual{
				&constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			triggeredConditionalOrders: []types.Order{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			setupMockBankKeeper: func(bk *mocks.BankKeeper) {
				bk.On(
					"SendCoinsFromModuleToModule",
					mock.Anything,
					satypes.ModuleName,
					authtypes.FeeCollectorName,
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
					constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					50_000_001,
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
							MakerOrderId: constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   50_000_000,
						},
					},
				),
			},
			expectedError: fmt.Errorf(
				"Match with Quantums 50000000 would exceed total Quantums 100000000 of "+
					"OrderId %v. New total filled quantums would be 100000001",
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}
