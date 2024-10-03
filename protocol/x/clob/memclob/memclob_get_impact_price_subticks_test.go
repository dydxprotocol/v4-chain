package memclob

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestMemClobPriceTimePriority_getImpactPriceSubticks(t *testing.T) {
	// the perpetual's atomic resolution must match that of makerOrders
	perpetual := constants.BtcUsd_0DefaultFunding_10AtomicResolution
	// perpetual := constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM
	clobPair := constants.ClobPair_Btc
	oraclePrice := types.SubticksToPrice(
		types.Subticks(500_000_000), // 50k
		constants.BtcUsdExponent,
		clobPair,
		perpetual.Params.AtomicResolution,
		lib.QuoteCurrencyAtomicResolution,
	)
	feeParams := constants.PerpetualFeeParamsNoFee

	// All subaccounts have high collat
	subaccountsForSufficientCollatTests := []satypes.Subaccount{
		constants.Bob_Num0_100_000USD,
		constants.Alice_Num0_100_000USD,
	}

	// One subaccount has low collat so its orders are skipped during collat check
	subaccountsForInsufficientCollatTests := []satypes.Subaccount{
		constants.Bob_Num0_1USD,
		constants.Alice_Num0_100_000USD,
	}

	subaccountsForPreciseCollatTests := []satypes.Subaccount{
		constants.Bob_Num0_50_000USD,
	}

	// makerBids/Asks should be set close to the oralce price so as to not affect
	// the net collateral of the account (which is updated based on trade cashflows)
	// This way, we can isolate IMR effects on collat check.
	makerAsks := []types.Order{
		{
			OrderId: types.OrderId{
				SubaccountId: constants.Bob_Num0,
				ClientId:     0,
				ClobPairId:   clobPair.Id,
			},
			Side:         types.Order_SIDE_SELL,
			Quantums:     10_000_000_000, // 1 BTC
			Subticks:     500_000_000,    // $50_000
			GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
		},
		{
			OrderId: types.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0,
				ClobPairId:   clobPair.Id,
			},
			Side:         types.Order_SIDE_SELL,
			Quantums:     10_000_000_000,
			Subticks:     500_010_000,
			GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
		},
	}

	makerBids := []types.Order{
		{
			OrderId: types.OrderId{
				SubaccountId: constants.Bob_Num0,
				ClientId:     1,
				ClobPairId:   clobPair.Id,
			},
			Side:         types.Order_SIDE_BUY,
			Quantums:     10_000_000_000,
			Subticks:     500_000_000,
			GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
		},
		{
			OrderId: types.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     1,
				ClobPairId:   clobPair.Id,
			},
			Side:         types.Order_SIDE_BUY,
			Quantums:     10_000_000_000,
			Subticks:     499_990_000,
			GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 1},
		},
	}

	tests := map[string]struct {
		// Params
		clobPair                    types.ClobPair
		subaccounts                 []satypes.Subaccount
		makerOrders                 []types.Order
		isBid                       bool
		impactNotionalQuoteQuantums *big.Int

		// Expectations
		expectedImpactPriceSubticks *big.Rat
		expectedHasEnoughLiquidity  bool
	}{
		// Collateralized tests
		`Buy crosses single level (orders collateralized) with sufficient liquidity`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForSufficientCollatTests,
			makerOrders:                 makerAsks,
			isBid:                       false,
			impactNotionalQuoteQuantums: big.NewInt(10_000_000_000), // $10000
			expectedImpactPriceSubticks: big.NewRat(500_000_000, 1),
			expectedHasEnoughLiquidity:  true,
		},
		`Buy crosses multiple levels (orders collateralized) and clears book exactly`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForSufficientCollatTests,
			makerOrders:                 makerAsks,
			isBid:                       false,
			impactNotionalQuoteQuantums: big.NewInt(100_001_000_000),
			expectedImpactPriceSubticks: new(big.Rat).Add(
				new(big.Rat).Mul(
					new(big.Rat).SetFrac(big.NewInt(1), big.NewInt(2)), big.NewRat(500_000_000, 1),
				),
				new(big.Rat).Mul(
					new(big.Rat).SetFrac(big.NewInt(1), big.NewInt(2)), big.NewRat(500_010_000, 1),
				),
			),
			expectedHasEnoughLiquidity: true,
		},
		`Buy notional too large for book (orders collateralized)`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForSufficientCollatTests,
			makerOrders:                 makerAsks,
			isBid:                       false,
			impactNotionalQuoteQuantums: big.NewInt(300_000_000_000),
			expectedImpactPriceSubticks: nil,
			expectedHasEnoughLiquidity:  false,
		},
		`Sell crosses single level (orders collateralized) with sufficient liquidity`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForSufficientCollatTests,
			makerOrders:                 makerBids,
			isBid:                       true,
			impactNotionalQuoteQuantums: big.NewInt(5_000_000_000),
			expectedImpactPriceSubticks: big.NewRat(500_000_000, 1),
			expectedHasEnoughLiquidity:  true,
		},
		`Sell crosses multiple levels (orders collateralized) and clears book exactly`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForSufficientCollatTests,
			makerOrders:                 makerBids,
			isBid:                       true,
			impactNotionalQuoteQuantums: big.NewInt(99_999_000_000),
			expectedImpactPriceSubticks: new(big.Rat).Add(
				new(big.Rat).Mul(
					new(big.Rat).SetFrac(big.NewInt(1), big.NewInt(2)), big.NewRat(500_000_000, 1),
				),
				new(big.Rat).Mul(
					new(big.Rat).SetFrac(big.NewInt(1), big.NewInt(2)), big.NewRat(499_990_000, 1),
				),
			),
			expectedHasEnoughLiquidity: true,
		},
		`Sell notional too large for book (orders collateralized)`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForSufficientCollatTests,
			makerOrders:                 makerBids,
			isBid:                       true,
			impactNotionalQuoteQuantums: big.NewInt(300_000_000_000),
			expectedImpactPriceSubticks: nil,
			expectedHasEnoughLiquidity:  false,
		},

		// Uncollateralized tests
		`Buy matches with ask that has insufficient collateral causing higher impact price`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForInsufficientCollatTests,
			makerOrders:                 makerAsks,
			isBid:                       false,
			impactNotionalQuoteQuantums: big.NewInt(5_000_000_000),
			expectedImpactPriceSubticks: big.NewRat(500_010_000, 1),
			expectedHasEnoughLiquidity:  true,
		},
		`Buy notional too large due to order with insufficient collateral`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForInsufficientCollatTests,
			makerOrders:                 makerAsks,
			isBid:                       false,
			impactNotionalQuoteQuantums: big.NewInt(70_000_000_000),
			expectedImpactPriceSubticks: nil,
			expectedHasEnoughLiquidity:  false,
		},
		`Sell matches with bid that has insufficient collateral causing lower impact price`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForInsufficientCollatTests,
			makerOrders:                 makerBids,
			isBid:                       true,
			impactNotionalQuoteQuantums: big.NewInt(5_000_000_000),
			expectedImpactPriceSubticks: big.NewRat(499_990_000, 1),
			expectedHasEnoughLiquidity:  true,
		},
		`Sell notional too large due to order with insufficient collateral`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForInsufficientCollatTests,
			makerOrders:                 makerBids,
			isBid:                       true,
			impactNotionalQuoteQuantums: big.NewInt(70_000_000_000),
			expectedImpactPriceSubticks: nil,
			expectedHasEnoughLiquidity:  false,
		},

		// Precise collat tests
		// One test should set impactNotionalQuoteQuantums = subaccount collat
		// The other test should set impactNotionalQuoteQuantums = subaccount collat + 1
		`Maker ask has exact collat`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForPreciseCollatTests,
			makerOrders:                 makerAsks[:1], // only use single order
			isBid:                       false,
			impactNotionalQuoteQuantums: big.NewInt(50_000_000_000),
			expectedImpactPriceSubticks: big.NewRat(500_000_000, 1),
			expectedHasEnoughLiquidity:  true,
		},
		`Maker ask has just insufficient collat`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForPreciseCollatTests,
			makerOrders:                 makerAsks[:1], // only use single order
			isBid:                       false,
			impactNotionalQuoteQuantums: big.NewInt(50_001_000_000),
			expectedImpactPriceSubticks: nil,
			expectedHasEnoughLiquidity:  false,
		},
		`Maker bid has exact collat`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForPreciseCollatTests,
			makerOrders:                 makerBids[:1], // only use single order
			isBid:                       true,
			impactNotionalQuoteQuantums: big.NewInt(50_000_000_000),
			expectedImpactPriceSubticks: big.NewRat(500_000_000, 1),
			expectedHasEnoughLiquidity:  true,
		},
		`Maker bid has just insufficient collat`: {
			clobPair:                    constants.ClobPair_Btc,
			subaccounts:                 subaccountsForPreciseCollatTests,
			makerOrders:                 makerBids[:1], // only use single order
			isBid:                       true,
			impactNotionalQuoteQuantums: big.NewInt(50_001_000_000),
			expectedImpactPriceSubticks: nil,
			expectedHasEnoughLiquidity:  false,
		},
	}

	for name, tc := range tests {
		memclob, ctx := initializeMemclobForTest(
			t, tc.clobPair, perpetual, oraclePrice, feeParams, tc.subaccounts, tc.makerOrders,
		)
		orderbook := memclob.orderbooks[tc.clobPair.GetClobPairId()]

		t.Run(name, func(t *testing.T) {
			impactPriceSubticks, hasEnoughLiquidity := memclob.getImpactPriceSubticks(
				ctx,
				tc.clobPair,
				orderbook,
				tc.isBid,
				tc.impactNotionalQuoteQuantums,
			)

			requireLargeValuesAlmostEqual(t, tc.expectedImpactPriceSubticks, impactPriceSubticks)
			require.Equal(t, tc.expectedHasEnoughLiquidity, hasEnoughLiquidity)
		})
	}
}

func initializeMemclobForTest(
	t *testing.T,
	clobPair types.ClobPair,
	perpetual perptypes.Perpetual, oralcePrice uint64,
	feeParams feetypes.PerpetualFeeParams,
	subaccounts []satypes.Subaccount,
	makerOrders []types.Order,
) (*MemClobPriceTimePriority, sdk.Context) {
	memClob := NewMemClobPriceTimePriority(false)
	mockBankKeeper := &mocks.BankKeeper{}
	ks := keepertest.NewClobKeepersTestContext(t, memClob, mockBankKeeper, indexer_manager.NewIndexerEventManagerNoop())

	ctx := ks.Ctx.WithIsCheckTx(true)

	// Create the default markets
	keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

	// Create liquidity tiers
	keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

	require.NoError(t, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, feeParams))

	// Set up USDC asset in assets module
	err := keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper)
	require.NoError(t, err)

	// Create perpetual
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

	perptest.SetUpDefaultPerpOIsForTest(
		t,
		ks.Ctx,
		ks.PerpetualsKeeper,
		[]perptypes.Perpetual{perpetual},
	)

	// Create all subaccounts
	for _, subaccount := range subaccounts {
		ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)
	}

	// Create CLOB
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

	err = ks.ClobKeeper.InitializeEquityTierLimit(
		ctx,
		types.EquityTierLimitConfiguration{
			ShortTermOrderEquityTiers: []types.EquityTierLimit{
				{
					UsdTncRequired: dtypes.NewInt(20_000_000),
					Limit:          5,
				},
			},
			StatefulOrderEquityTiers: []types.EquityTierLimit{
				{
					UsdTncRequired: dtypes.NewInt(20_000_000),
					Limit:          5,
				},
			},
		},
	)
	require.NoError(t, err)

	for _, order := range makerOrders {
		_, _, err := ks.ClobKeeper.PlaceShortTermOrder(ctx, &types.MsgPlaceOrder{Order: order})
		require.NoError(t, err)
	}

	// Set oracle price
	err = ks.PricesKeeper.UpdateMarketPrices(ctx, []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
		{
			MarketId: clobPair.Id,
			Price:    oralcePrice,
		},
	})
	require.NoError(t, err)

	return memClob, ctx
}

func requireLargeValuesAlmostEqual(t *testing.T, a, b *big.Rat) {
	if a == nil && b == nil {
		require.True(t, true)
	} else if a != nil && b != nil {
		// epsilon value set large because this function compares very large values
		epsilon := new(big.Rat).SetFloat64(1e0)
		diff := new(big.Rat).Sub(a, b)
		diff.Abs(diff)
		require.True(t, diff.Cmp(epsilon) < 0, "a: %v, b: %v, diff: %v", a, b, diff)
	} else {
		require.Fail(t, "One of the values is nil", "a: %v, b: %v", a, b)
	}
}
