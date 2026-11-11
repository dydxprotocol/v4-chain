package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtest "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mainnet short-term ladder.
var benchShortTermEquityTiers = []types.EquityTierLimit{
	{UsdTncRequired: dtypes.NewInt(0), Limit: 0},
	{UsdTncRequired: dtypes.NewInt(20_000_000), Limit: 1},
	{UsdTncRequired: dtypes.NewInt(100_000_000), Limit: 5},
	{UsdTncRequired: dtypes.NewInt(1_000_000_000), Limit: 10},
	{UsdTncRequired: dtypes.NewInt(10_000_000_000), Limit: 100},
	{UsdTncRequired: dtypes.NewInt(100_000_000_000), Limit: 1000},
}

const benchNumPerpetuals = 20

func benchSubaccountWithPositions(n int) satypes.Subaccount {
	positions := make([]*satypes.PerpetualPosition, n)
	for i := range positions {
		positions[i] = testutil.CreateSinglePerpetualPosition(
			uint32(i),
			big.NewInt(1_000_000), // 0.01 BTC
			big.NewInt(0),
			big.NewInt(0),
		)
	}
	return satypes.Subaccount{
		Id:                 &constants.Carl_Num0,
		AssetPositions:     testutil.CreateUsdcAssetPositions(big.NewInt(100_000_000_000)), // $100,000
		PerpetualPositions: positions,
	}
}

var benchSubaccounts = []struct {
	name       string
	subaccount satypes.Subaccount
}{
	{"no_positions", constants.Alice_Num0_10_000USD},
	{"two_positions", constants.Dave_Num0_1BTC_Long_1ETH_Long_46000USD_Short},
	{"twenty_positions", benchSubaccountWithPositions(benchNumPerpetuals)},
}

func setupShortTermOrderKeeper(
	b testing.TB,
	subaccount satypes.Subaccount,
	shortTermTiers []types.EquityTierLimit,
) (keepertest.ClobKeepersTestContext, sdk.Context) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	bankKeeper := &mocks.BankKeeper{}
	bankKeeper.On("SendCoins", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	ks := keepertest.NewClobKeepersTestContext(b, memClob, bankKeeper, indexer_manager.NewIndexerEventManagerNoop())
	ctx := ks.Ctx.WithIsCheckTx(true)

	keepertest.CreateTestMarkets(b, ctx, ks.PricesKeeper)
	keepertest.CreateTestLiquidityTiers(b, ctx, ks.PerpetualsKeeper)
	require.NoError(b, ks.FeeTiersKeeper.SetPerpetualFeeParams(ctx, constants.PerpetualFeeParamsNoFee))
	require.NoError(b, keepertest.CreateUsdcAsset(ctx, ks.AssetsKeeper))

	perpetuals := []perptypes.Perpetual{
		constants.BtcUsd_20PercentInitial_10PercentMaintenance,
		constants.EthUsd_20PercentInitial_10PercentMaintenance,
	}
	for id := uint32(len(perpetuals)); id < benchNumPerpetuals; id++ {
		p := constants.BtcUsd_20PercentInitial_10PercentMaintenance
		p.Params.Id = id
		p.Params.Ticker = fmt.Sprintf("BENCH-%d", id)
		perpetuals = append(perpetuals, p)
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
		require.NoError(b, err)
	}
	perptest.SetUpDefaultPerpOIsForTest(b, ks.Ctx, ks.PerpetualsKeeper, perpetuals)
	ks.SubaccountsKeeper.SetSubaccount(ctx, subaccount)

	for _, clobPair := range []types.ClobPair{constants.ClobPair_Btc, constants.ClobPair_Eth} {
		_, err := ks.ClobKeeper.CreatePerpetualClobPairAndMemStructs(
			ctx,
			clobPair.Id,
			clobtest.MustPerpetualId(clobPair),
			satypes.BaseQuantums(clobPair.StepBaseQuantums),
			clobPair.QuantumConversionExponent,
			clobPair.SubticksPerTick,
			clobPair.Status,
		)
		require.NoError(b, err)
	}
	require.NoError(b, ks.ClobKeeper.InitializeEquityTierLimit(ctx, types.EquityTierLimitConfiguration{
		ShortTermOrderEquityTiers: shortTermTiers,
	}))
	return ks, ctx
}

func restingBuyOrder(subaccountId satypes.SubaccountId, clientId uint32) types.Order {
	return types.Order{
		OrderId:      types.OrderId{SubaccountId: subaccountId, ClientId: clientId, ClobPairId: 0},
		Side:         types.Order_SIDE_BUY,
		Quantums:     5,
		Subticks:     10,
		GoodTilOneof: &types.Order_GoodTilBlock{GoodTilBlock: 15},
	}
}

func BenchmarkPlaceShortTermOrder(b *testing.B) {
	for _, sa := range benchSubaccounts {
		for _, gate := range []struct {
			name  string
			tiers []types.EquityTierLimit
		}{{"gate_off", nil}, {"gate_on", benchShortTermEquityTiers}} {
			b.Run(sa.name+"/"+gate.name, func(b *testing.B) {
				ks, ctx := setupShortTermOrderKeeper(b, sa.subaccount, gate.tiers)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					order := restingBuyOrder(*sa.subaccount.Id, uint32(i))
					_, _, err := ks.ClobKeeper.PlaceShortTermOrder(
						ctx.WithTxBytes([]byte(fmt.Sprint(i))),
						&types.MsgPlaceOrder{Order: order},
					)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

func BenchmarkValidateSubaccountEquityTierLimitForShortTermOrder(b *testing.B) {
	for _, sa := range benchSubaccounts {
		for _, cold := range []bool{true, false} {
			name := sa.name + "/warm"
			if cold {
				name = sa.name + "/cold"
			}
			b.Run(name, func(b *testing.B) {
				ks, ctx := setupShortTermOrderKeeper(b, sa.subaccount, benchShortTermEquityTiers)
				order := restingBuyOrder(*sa.subaccount.Id, 0)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					if cold {
						ks.ClobKeeper.ClearShortTermEquityTierResults()
					}
					if err := ks.ClobKeeper.ValidateSubaccountEquityTierLimitForShortTermOrder(ctx, order); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

func BenchmarkGetEquityTierLimitConfiguration(b *testing.B) {
	ks, ctx := setupShortTermOrderKeeper(b, constants.Alice_Num0_10_000USD, benchShortTermEquityTiers)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ks.ClobKeeper.GetEquityTierLimitConfiguration(ctx)
	}
}
