package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"

	sdk "github.com/cosmos/cosmos-sdk/types"
	marketmapkeeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestCreateOracleMarket(t *testing.T) {
	testMarket1 := *pricestest.GenerateMarketParamPrice(
		pricestest.WithId(1),
		pricestest.WithPair(constants.BtcUsdPair),
		pricestest.WithExponent(-8), // for both Param and Price
		pricestest.WithPriceValue(0),
	)
	testCases := map[string]struct {
		setup           func(*testing.T, sdk.Context, *keeper.Keeper)
		msg             *pricestypes.MsgCreateOracleMarket
		expectedMarkets []pricestypes.MarketParamPrice
		expectedErr     string
	}{
		"Succeeds: create new oracle market (id = 1)": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {
				keepertest.CreateMarketsInMarketMapFromParams(
					t,
					ctx,
					pricesKeeper.MarketMapKeeper.(*marketmapkeeper.Keeper),
					[]pricestypes.MarketParam{testMarket1.Param},
				)
			},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: lib.GovModuleAddress.String(),
				Params:    testMarket1.Param,
			},
			expectedMarkets: []pricestypes.MarketParamPrice{testMarket1},
		},
		"Failure: empty pair": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: lib.GovModuleAddress.String(),
				Params: pricestest.GenerateMarketParamPrice(
					pricestest.WithPair(""),
					pricestest.WithExponent(-8), // for both Param and Price
				).Param,
			},
			expectedMarkets: []pricestypes.MarketParamPrice{},
			expectedErr:     "incorrectly formatted CurrencyPair",
		},
		"Failure: market not found in MarketMap": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: lib.GovModuleAddress.String(),
				Params:    testMarket1.Param,
			},
			expectedErr:     "Ticker not found in market map",
			expectedMarkets: []pricestypes.MarketParamPrice{},
		},
		"Failure: oracle market id already exists": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {
				keepertest.CreateTestPriceMarkets(
					t,
					ctx,
					pricesKeeper,
					[]pricestypes.MarketParamPrice{testMarket1},
				)
				keepertest.CreateMarketsInMarketMapFromParams(
					t,
					ctx,
					pricesKeeper.MarketMapKeeper.(*marketmapkeeper.Keeper),
					[]pricestypes.MarketParam{pricestest.GenerateMarketParamPrice(
						pricestest.WithId(1), // same id as testMarket1
					).Param},
				)
			},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: lib.GovModuleAddress.String(),
				Params: pricestest.GenerateMarketParamPrice(
					pricestest.WithId(1), // same id as testMarket1
				).Param,
			},
			expectedErr:     "Market params already exists",
			expectedMarkets: []pricestypes.MarketParamPrice{testMarket1},
		},
		"Failure: invalid authority": {
			setup: func(t *testing.T, ctx sdk.Context, pricesKeeper *keeper.Keeper) {},
			msg: &pricestypes.MsgCreateOracleMarket{
				Authority: "invalid",
				Params:    testMarket1.Param,
			},
			expectedMarkets: []pricestypes.MarketParamPrice{},
			expectedErr:     "invalid authority",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx, pricesKeeper, _, _, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			msgServer := keeper.NewMsgServerImpl(pricesKeeper)
			tc.setup(t, ctx, pricesKeeper)

			// Check that market is disabled in MarketMap before creating Oracle market for it
			if len(tc.expectedMarkets) > 0 && tc.expectedErr == "" {
				for i := range tc.expectedMarkets {
					market := tc.expectedMarkets[i]

					currencyPair, _ := slinky.MarketPairToCurrencyPair(market.Param.Pair)
					mmMarket, _ := pricesKeeper.MarketMapKeeper.GetMarket(ctx, currencyPair.String())
					require.False(t, mmMarket.Ticker.Enabled)
				}
			}

			_, err := msgServer.CreateOracleMarket(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			gotAllMarketParamPrices, err := pricesKeeper.GetAllMarketParamPrices(ctx)
			require.NoError(t, err)
			require.Equal(t, tc.expectedMarkets, gotAllMarketParamPrices)

			// Check if the market is enabled in MarketMap
			if len(tc.expectedMarkets) > 0 {
				for i := range tc.expectedMarkets {
					market := tc.expectedMarkets[i]

					currencyPair, _ := slinky.MarketPairToCurrencyPair(market.Param.Pair)
					mmMarket, _ := pricesKeeper.MarketMapKeeper.GetMarket(ctx, currencyPair.String())
					require.True(t, mmMarket.Ticker.Enabled)
				}
			}
		})
	}
}

func TestMarketPriceExponentIsFromMarketmap(t *testing.T) {
	ctx, pricesKeeper, _, _, _, _, marketMapKeeper := keepertest.PricesKeepers(t)
	msgServer := keeper.NewMsgServerImpl(pricesKeeper)

	// Create test market in marketmap
	currencyPair, err := slinky.MarketPairToCurrencyPair(constants.BtcUsdPair)
	require.NoError(t, err)

	marketMapDetails := marketmaptypes.Market{
		Ticker: marketmaptypes.Ticker{
			CurrencyPair:     currencyPair,
			Decimals:         uint64(8),
			MinProviderCount: 1,
		},
		ProviderConfigs: []marketmaptypes.ProviderConfig{},
	}
	err = marketMapKeeper.CreateMarket(ctx, marketMapDetails)
	require.NoError(t, err)

	// Send message to create oracle market without setting exponent
	// because the market price exponent is calculated from the marketmap Decimals
	testMarket := pricestest.GenerateMarketParamPrice(
		pricestest.WithId(1),
		pricestest.WithPair(constants.BtcUsdPair),
		pricestest.WithPriceValue(0),
		// Do not set exponent
	)

	msg := &pricestypes.MsgCreateOracleMarket{
		Authority: lib.GovModuleAddress.String(),
		Params:    testMarket.Param,
	}
	_, err = msgServer.CreateOracleMarket(ctx, msg)
	require.NoError(t, err)

	// Verify that the market price exponent matches negation of the marketmap Decimals
	marketPrice, err := pricesKeeper.GetMarketPrice(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, int32(-8), marketPrice.Exponent)
}
