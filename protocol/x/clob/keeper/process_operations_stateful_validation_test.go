package keeper_test

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"

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

func TestProcessProposerMatches_LongTerm_StatefulValidation_Failure(t *testing.T) {
	tests := map[string]processProposerOperationsTestCase{
		`Stateful order validation: referenced maker order does not exist in state`: {
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
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful long term order id %+v does not exist in state.",
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced taker order does not exist in state`: {
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
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful long term order id %+v does not exist in state.",
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced maker order in liquidation match does not exist in state`: {
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
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful long term order id %+v does not exist in state.",
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced long-term order is on the wrong side`: {
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
			expectedError: errors.New("Orders are not on opposing sides of the book in match"),
		},
		`Stateful match validation: taker order cannot be post only`: {
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
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_PO,
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_PO,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: errorsmod.Wrapf(
				types.ErrInvalidMatchOrder,
				"Taker order %+v cannot be post only.",
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10_PO.GetOrderTextString(),
			),
		},
		`Stateful match validation: maker order cannot be IOC`: {
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
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_IOC,
			},
			setupState: func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
				ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
					Timestamp: time.Unix(5, 0),
				})
			},
			rawOperations: []types.OperationRaw{
				clobtest.NewMatchOperationRaw(
					&constants.LongTermOrder_Dave_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
					[]types.MakerFill{
						{
							MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10_IOC.OrderId,
							FillAmount:   100_000_000, // 1 BTC
						},
					},
				),
			},
			expectedError: errors.New("IOC order cannot be matched as a maker order"),
		},
		`Stateful order validation: referenced long-term order is for the wrong clob pair`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_100PercentMarginRequirement,
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

// TestProcessProposerMatches_ExpiredStatefulOrderSkipped verifies that a proposed match that
// references a stateful order which has expired (its GoodTilBlockTime <= block time) is SKIPPED
// rather than aborting the whole MsgProposedOperations transaction — so healthy matches processed
// alongside it are not rolled back. Expiry cleanup is budgeted in the EndBlocker, so an order can
// be valid when the proposer assembles the block and expired by FinalizeBlock.
func TestProcessProposerMatches_ExpiredStatefulOrderSkipped(t *testing.T) {
	tests := map[string]processProposerOperationsTestCase{
		`Expired maker order in a proposed match is skipped, not fatal`: {
			blockTime: time.Unix(10, 0),
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
				// GoodTilBlockTime 10 == block time 10 => expired.
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
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
							FillAmount:   100_000_000,
						},
					},
				),
			},
			// No error; the match against the expired maker is skipped, so nothing fills.
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: 5,
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         0,
			},
		},
		`Expired taker order in a proposed match is skipped, not fatal`: {
			blockTime: time.Unix(10, 0),
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
				// GoodTilBlockTime 10 == block time 10 => the taker is expired.
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10,
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
							FillAmount:   100_000_000,
						},
					},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				BlockHeight: 5,
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				constants.LongTermOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId: 0,
				constants.Order_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTB10.OrderId:           0,
			},
		},
		`Healthy matches before and after an expired match all succeed`: {
			blockTime: time.Unix(10, 0),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id:             &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{&constants.Usdc_Asset_100_000},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(0, big.NewInt(1_000_000_000), big.NewInt(0), big.NewInt(0)),
					},
				},
				{
					Id:             &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{&constants.Usdc_Asset_100_000},
					PerpetualPositions: []*satypes.PerpetualPosition{
						testutil.CreateSinglePerpetualPosition(0, big.NewInt(1_000_000_000), big.NewInt(0), big.NewInt(0)),
					},
				},
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			perpetualFeeParams: &constants.PerpetualFeeParams,
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},
			preExistingStatefulOrders: []types.Order{
				// Expired maker (GoodTilBlockTime 10 == block time 10).
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
			},
			rawOperations: []types.OperationRaw{
				// Healthy match BEFORE the expired one: Bob (taker) sells to Alice (maker).
				clobtest.NewShortTermOrderPlacementOperationRaw(healthyBuy(constants.Alice_Num0, 14)),
				clobtest.NewShortTermOrderPlacementOperationRaw(healthySell(constants.Bob_Num0, 14)),
				clobtest.NewMatchOperationRaw(
					ptrOrder(healthySell(constants.Bob_Num0, 14)),
					[]types.MakerFill{{FillAmount: 100_000_000, MakerOrderId: healthyBuy(constants.Alice_Num0, 14).OrderId}},
				),
				// Expired match: Dave (taker) vs the expired Carl maker — skipped, not fatal.
				clobtest.NewShortTermOrderPlacementOperationRaw(
					constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
				),
				clobtest.NewMatchOperationRaw(
					&constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10,
					[]types.MakerFill{{
						FillAmount:   100_000_000,
						MakerOrderId: constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
					}},
				),
				// Healthy match AFTER the expired one: Bob (taker) sells to Alice (maker) again.
				clobtest.NewShortTermOrderPlacementOperationRaw(healthyBuy(constants.Alice_Num0, 15)),
				clobtest.NewShortTermOrderPlacementOperationRaw(healthySell(constants.Bob_Num0, 15)),
				clobtest.NewMatchOperationRaw(
					ptrOrder(healthySell(constants.Bob_Num0, 15)),
					[]types.MakerFill{{FillAmount: 100_000_000, MakerOrderId: healthyBuy(constants.Alice_Num0, 15).OrderId}},
				),
			},
			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				// Sorted by SortedOrders (Bob's subaccount sorts before Alice's).
				OrderIdsFilledInLastBlock: []types.OrderId{
					healthySell(constants.Bob_Num0, 14).OrderId,
					healthySell(constants.Bob_Num0, 15).OrderId,
					healthyBuy(constants.Alice_Num0, 14).OrderId,
					healthyBuy(constants.Alice_Num0, 15).OrderId,
				},
				BlockHeight: 5,
			},
			expectedFillAmounts: map[types.OrderId]satypes.BaseQuantums{
				// Both healthy matches filled despite the expired one being skipped between them.
				healthyBuy(constants.Alice_Num0, 14).OrderId:                                  100_000_000,
				healthySell(constants.Bob_Num0, 14).OrderId:                                   100_000_000,
				healthyBuy(constants.Alice_Num0, 15).OrderId:                                  100_000_000,
				healthySell(constants.Bob_Num0, 15).OrderId:                                   100_000_000,
				constants.LongTermOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId: 0,
				constants.Order_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTB10.OrderId:         0,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}

// healthyBuy / healthySell build simple crossing short-term BTC orders (never expire by block time)
// used to assert that healthy matches survive alongside a skipped expired match.
func healthyBuy(subaccountId satypes.SubaccountId, clientId uint32) types.Order {
	return types.Order{
		OrderId:      types.OrderId{SubaccountId: subaccountId, ClientId: clientId, ClobPairId: 0},
		Side:         types.Order_SIDE_BUY,
		Quantums:     100_000_000,
		Subticks:     50_000_000,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
	}
}

func healthySell(subaccountId satypes.SubaccountId, clientId uint32) types.Order {
	return types.Order{
		OrderId:      types.OrderId{SubaccountId: subaccountId, ClientId: clientId, ClobPairId: 0},
		Side:         types.Order_SIDE_SELL,
		Quantums:     100_000_000,
		Subticks:     50_000_000,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 25},
	}
}

func ptrOrder(o types.Order) *types.Order {
	return &o
}

func TestProcessProposerMatches_Conditional_Validation_Failure(t *testing.T) {
	tests := map[string]processProposerOperationsTestCase{
		`Stateful order validation: referenced maker order does not exist in state`: {
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
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced taker order does not exist in state`: {
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
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced maker order in liquidation match does not exist in state`: {
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
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Dave_Num0_Id0_Clob0_Sell1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced maker order exist in state but is untriggered`: {
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
			expectedError: errorsmod.Wrapf(
				types.ErrStatefulOrderDoesNotExist,
				"stateful conditional order id %+v does not exist in triggered conditional state.",
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10.OrderId,
			),
		},
		`Stateful order validation: referenced conditional order is on the wrong side`: {
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
			expectedError: errors.New("Orders are not on opposing sides of the book in match"),
		},
		`Stateful order validation: referenced conditional order is for the wrong clob pair`: {
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_100PercentMarginRequirement,
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
			triggeredConditionalOrders: []types.Order{
				constants.ConditionalOrder_Carl_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT10,
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
