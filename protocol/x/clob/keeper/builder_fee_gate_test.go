package keeper_test

import (
	"math/big"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	affiliatetypes "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

const (
	builderFeeGateBlockHeight = uint32(100)
	// 1% of the $50,000 fill.
	builderFeeGateFee = 500_000_000
)

type gateSetup = func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext)

func builderFeeGateTakerOrder() types.Order {
	return types.Order{
		OrderId:      types.OrderId{SubaccountId: constants.Bob_Num0, ClientId: 14, ClobPairId: 0},
		Side:         types.Order_SIDE_SELL,
		Quantums:     100_000_000,    // 1 BTC
		Subticks:     50_000_000_000, // $50,000
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: builderFeeGateBlockHeight + 10},
		BuilderCodeParameters: &types.BuilderCodeParameters{
			BuilderAddress: constants.Carl_Num0.Owner,
			FeePpm:         10_000,
		},
	}
}

func builderFeeGateMakerOrder(withBuilderCode bool) types.Order {
	order := constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy1BTC_Price50000_GTBT15
	if withBuilderCode {
		order.BuilderCodeParameters = builderFeeGateTakerOrder().BuilderCodeParameters
	}
	return order
}

// builderFeeGateTestCase fills Bob's 1 BTC sell against Alice's resting buy at block 100. Bob, and
// Alice too if makerBuilderCode, owe Carl a 1% builder fee, which the case expects charged or waived.
func builderFeeGateTestCase(
	setupGate gateSetup,
	makerBuilderCode bool,
	charged bool,
	extraPerpetuals ...perptypes.Perpetual,
) processProposerOperationsTestCase {
	taker := builderFeeGateTakerOrder()
	maker := builderFeeGateMakerOrder(makerBuilderCode)
	match := types.MatchWithOrders{
		TakerOrder: &taker,
		MakerOrder: &maker,
		FillAmount: 100_000_000,
		MakerFee:   10_000_000,
		TakerFee:   25_000_000,
	}
	aliceQuote := int64(100_000_000_000 - 50_000_000_000 - 10_000_000)
	bobQuote := int64(100_000_000_000 + 50_000_000_000 - 25_000_000)
	builderPayments := 0
	if charged {
		match.TakerBuilderFee = builderFeeGateFee
		bobQuote -= builderFeeGateFee
		builderPayments++
		if makerBuilderCode {
			match.MakerBuilderFee = builderFeeGateFee
			aliceQuote -= builderFeeGateFee
			builderPayments++
		}
	}

	subaccountWith10Btc := func(id satypes.SubaccountId) satypes.Subaccount {
		return satypes.Subaccount{
			Id:             &id,
			AssetPositions: []*satypes.AssetPosition{&constants.Usdc_Asset_100_000},
			PerpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(0, big.NewInt(1_000_000_000), big.NewInt(0), big.NewInt(0)),
			},
		}
	}
	return processProposerOperationsTestCase{
		perpetuals:         append([]perptypes.Perpetual{constants.BtcUsd_100PercentMarginRequirement}, extraPerpetuals...),
		perpetualFeeParams: &constants.PerpetualFeeParams,
		clobPairs:          []types.ClobPair{constants.ClobPair_Btc},
		subaccounts: []satypes.Subaccount{
			subaccountWith10Btc(constants.Alice_Num0),
			subaccountWith10Btc(constants.Bob_Num0),
			{Id: &constants.Carl_Num0, AssetPositions: []*satypes.AssetPosition{&constants.Usdc_Asset_0}},
		},
		preExistingStatefulOrders: []types.Order{maker},
		setupState:                setupGate,
		setupMockBankKeeper: func(bk *mocks.BankKeeper) {
			bk.On("SendCoinsFromModuleToModule", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			bk.On("BlockedAddr", mock.Anything).Return(false)
			// With no expectation for Carl, an unexpected builder payment fails the test.
			if builderPayments > 0 {
				bk.On(
					"SendCoins",
					mock.Anything,
					mock.Anything,
					constants.CarlAccAddress,
					mock.MatchedBy(func(coins sdk.Coins) bool {
						return coins.AmountOf(assettypes.UusdcDenom).Equal(sdkmath.NewInt(builderFeeGateFee))
					}),
				).Return(nil).Times(builderPayments)
			}
			bk.On(
				"SendCoins",
				mock.Anything,
				mock.Anything,
				authtypes.NewModuleAddress(authtypes.FeeCollectorName),
				mock.Anything,
			).Return(nil).Once()
			bk.On("GetBalance", mock.Anything, mock.Anything, mock.Anything).
				Return(sdk.NewCoin("USDC", sdkmath.NewIntFromUint64(100_000_000_000)))
		},
		rawOperations: []types.OperationRaw{
			clobtest.NewShortTermOrderPlacementOperationRaw(taker),
			clobtest.NewMatchOperationRaw(&taker, []types.MakerFill{
				{FillAmount: 100_000_000, MakerOrderId: maker.OrderId},
			}),
		},
		expectedMatches: []*MatchWithOrdersForTesting{
			{MatchWithOrders: match, TotalFilledMaker: 100_000_000, TotalFilledTaker: 100_000_000},
		},
		expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
			OrderIdsFilledInLastBlock: []types.OrderId{taker.OrderId, maker.OrderId},
			RemovedStatefulOrderIds:   []types.OrderId{maker.OrderId},
			BlockHeight:               builderFeeGateBlockHeight,
		},
		expectedQuoteBalances: map[satypes.SubaccountId]int64{
			constants.Alice_Num0: aliceQuote,
			constants.Bob_Num0:   bobQuote,
		},
		expectedPerpetualPositions: map[satypes.SubaccountId][]*satypes.PerpetualPosition{
			constants.Alice_Num0: {
				testutil.CreateSinglePerpetualPosition(0, big.NewInt(1_100_000_000), big.NewInt(0), big.NewInt(0)),
			},
			constants.Bob_Num0: {
				testutil.CreateSinglePerpetualPosition(0, big.NewInt(900_000_000), big.NewInt(0), big.NewInt(0)),
			},
			constants.Carl_Num0: {},
		},
	}
}

func negativeTncSeenAt(t *testing.T, perpetualId uint32, height uint32) gateSetup {
	return func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
		require.NoError(t, ks.SubaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, height))
	}
}

func chainOutageEndedAt(height uint32) gateSetup {
	return func(ctx sdk.Context, ks keepertest.ClobKeepersTestContext) {
		ks.BlockTimeKeeper.SetAllDowntimeInfo(ctx, &blocktimetypes.AllDowntimeInfo{
			Infos: []*blocktimetypes.AllDowntimeInfo_DowntimeInfo{{
				Duration: satypes.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_CHAIN_OUTAGE_DURATION,
				BlockInfo: blocktimetypes.BlockInfo{
					Height:    height,
					Timestamp: time.Unix(4, 0),
				},
			}},
		})
	}
}

func TestProcessProposerOperations_BuilderFeeGate(t *testing.T) {
	const window = satypes.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS
	btc := constants.ClobPair_Btc.MustGetPerpetualId()
	iso := constants.IsoUsd_IsolatedMarket.Params.Id

	tests := map[string]processProposerOperationsTestCase{
		"waived in the gate's last block": builderFeeGateTestCase(
			negativeTncSeenAt(t, btc, builderFeeGateBlockHeight-(window-1)), false, false,
		),
		"charged in the first block after the gate's window": builderFeeGateTestCase(
			negativeTncSeenAt(t, btc, builderFeeGateBlockHeight-window), false, true,
		),
		"maker and taker builder fees are both waived": builderFeeGateTestCase(
			negativeTncSeenAt(t, btc, builderFeeGateBlockHeight), true, false,
		),
		"maker and taker builder fees are both charged without the gate": builderFeeGateTestCase(
			nil, true, true,
		),
		"waived after a chain outage": builderFeeGateTestCase(
			chainOutageEndedAt(builderFeeGateBlockHeight-1), false, false,
		),
		"charged in the first block after a chain outage's window": builderFeeGateTestCase(
			chainOutageEndedAt(builderFeeGateBlockHeight-window), false, true,
		),
		"charged when only another collateral pool is gated": builderFeeGateTestCase(
			negativeTncSeenAt(t, iso, builderFeeGateBlockHeight), false, true, constants.IsoUsd_IsolatedMarket,
		),
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runProcessProposerOperationsTestCase(t, tc)
		})
	}
}

// CheckTx state sits at the last committed height, and the matches it validates are persisted in
// the next block. In the gate's last block, CheckTx must therefore already charge the builder fee
// that DeliverTx in the next block will charge.
func TestProcessSingleMatch_BuilderFeeGateLastBlock(t *testing.T) {
	const window = satypes.WITHDRAWAL_AND_TRANSFERS_BLOCKED_AFTER_NEGATIVE_TNC_SUBACCOUNT_SEEN_BLOCKS
	btc := constants.ClobPair_Btc.MustGetPerpetualId()

	for name, isCheckTx := range map[string]bool{
		"DeliverTx waives the fee in the gate's last block": false,
		"CheckTx charges the fee in the gate's last block":  true,
	} {
		t.Run(name, func(t *testing.T) {
			tc := builderFeeGateTestCase(
				negativeTncSeenAt(t, btc, builderFeeGateBlockHeight-(window-1)), false, isCheckTx,
			)
			// Only the match is processed here, so no fill events are expected.
			tc.expectedMatches = nil
			ctx, ks, _ := setupProcessProposerOperationsTestCase(t, tc)

			taker := builderFeeGateTakerOrder()
			maker := builderFeeGateMakerOrder(false)
			match := &types.MatchWithOrders{TakerOrder: &taker, MakerOrder: &maker, FillAmount: 100_000_000}
			_, _, _, _, err := ks.ClobKeeper.ProcessSingleMatch(
				ctx.WithIsCheckTx(isCheckTx), match, map[string]bool{}, affiliatetypes.AffiliateParameters{},
			)
			require.NoError(t, err)

			expectedFee := uint64(0)
			if isCheckTx {
				expectedFee = builderFeeGateFee
			}
			require.Equal(t, expectedFee, match.TakerBuilderFee)
		})
	}
}
