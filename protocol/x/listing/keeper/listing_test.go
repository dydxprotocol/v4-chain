package keeper_test

import (
	"errors"
	"testing"

	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	oracletypes "github.com/skip-mev/slinky/pkg/types"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/skip-mev/slinky/x/marketmap/types/tickermetadata"

	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestCreateMarket(t *testing.T) {
	tests := map[string]struct {
		ticker          string
		duplicateMarket bool

		expectedErr error
	}{
		"success": {
			ticker:          "TEST-USD",
			duplicateMarket: false,
			expectedErr:     nil,
		},
		"failure - invalid market": {
			ticker:          "INVALID-USD",
			duplicateMarket: false,
			expectedErr:     types.ErrMarketNotFound,
		},
		"failure - duplicate market": {
			ticker:          "TEST-USD",
			duplicateMarket: true,
			expectedErr:     nil,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				mockIndexerEventManager := &mocks.IndexerEventManager{}
				ctx, keeper, _, _, pricesKeeper, _, _, marketMapKeeper := keepertest.ListingKeepers(
					t,
					&mocks.BankKeeper{},
					mockIndexerEventManager,
				)

				testMarketParams := pricestypes.MarketParam{
					Pair:               "TEST-USD",
					Exponent:           int32(-6),
					ExchangeConfigJson: `{"test_config_placeholder":{}}`,
					MinExchanges:       2,
					MinPriceChangePpm:  uint32(800),
				}

				keepertest.CreateMarketsInMarketMapFromParams(
					t,
					ctx,
					marketMapKeeper,
					[]pricestypes.MarketParam{
						testMarketParams,
					},
				)

				marketId, err := keeper.CreateMarket(ctx, tc.ticker)
				if tc.expectedErr != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					// Check if the market was created
					market, exists := pricesKeeper.GetMarketParam(ctx, marketId)
					require.True(t, exists)
					require.Equal(t, testMarketParams.Pair, market.Pair)
					require.Equal(t, testMarketParams.Exponent, market.Exponent)
					require.Equal(t, testMarketParams.MinExchanges, market.MinExchanges)
					require.Equal(t, testMarketParams.MinPriceChangePpm, types.MinPriceChangePpm_LongTail)
				}

				if tc.duplicateMarket {
					_, err = keeper.CreateMarket(ctx, tc.ticker)
					require.ErrorContains(t, err, pricestypes.ErrMarketParamPairAlreadyExists.Error())
				}
			},
		)
	}
}

func TestCreatePerpetual(t *testing.T) {
	tests := map[string]struct {
		ticker         string
		referencePrice uint64

		expectedErr error
	}{
		"success": {
			ticker:         "TEST-USD",
			referencePrice: 1000000000,
			expectedErr:    nil,
		},
		"failure - reference price 0": {
			ticker:         "TEST-USD",
			referencePrice: 0,
			expectedErr:    types.ErrReferencePriceZero,
		},
		"failure - invalid market": {
			ticker:      "INVALID-USD",
			expectedErr: types.ErrMarketNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				mockIndexerEventManager := &mocks.IndexerEventManager{}
				ctx, keeper, _, _, pricesKeeper, perpetualsKeeper, _, marketMapKeeper := keepertest.ListingKeepers(
					t,
					&mocks.BankKeeper{},
					mockIndexerEventManager,
				)
				keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, perpetualsKeeper, pricesKeeper, 10)

				// Create a marketmap with a single market
				dydxMetadata, err := tickermetadata.MarshalDyDx(
					tickermetadata.DyDx{
						ReferencePrice: tc.referencePrice,
						Liquidity:      0,
						AggregateIDs:   nil,
					},
				)
				require.NoError(t, err)

				market := marketmaptypes.Market{
					Ticker: marketmaptypes.Ticker{
						CurrencyPair:     oracletypes.CurrencyPair{Base: "TEST", Quote: "USD"},
						Decimals:         6,
						MinProviderCount: 2,
						Enabled:          false,
						Metadata_JSON:    string(dydxMetadata),
					},
					ProviderConfigs: []marketmaptypes.ProviderConfig{
						{
							Name:           "binance_ws",
							OffChainTicker: "TESTUSDT",
						},
					},
				}
				err = marketMapKeeper.CreateMarket(ctx, market)
				require.NoError(t, err)

				marketId, err := keeper.CreateMarket(ctx, tc.ticker)
				if errors.Is(tc.expectedErr, types.ErrMarketNotFound) {
					require.ErrorContains(t, err, tc.expectedErr.Error())
					return
				}

				perpetualId, err := keeper.CreatePerpetual(ctx, marketId, tc.ticker)
				if tc.expectedErr != nil {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					// Check if the perpetual was created
					perpetual, err := perpetualsKeeper.GetPerpetual(ctx, perpetualId)
					require.NoError(t, err)
					require.Equal(t, uint32(10), perpetual.GetId())
					require.Equal(t, marketId, perpetual.Params.MarketId)
					require.Equal(t, tc.ticker, perpetual.Params.Ticker)
					// Expected resolution = -6 - Floor(log10(1000000000)) = -15
					require.Equal(t, int32(-15), perpetual.Params.AtomicResolution)
					require.Equal(t, int32(types.DefaultFundingPpm), perpetual.Params.DefaultFundingPpm)
					require.Equal(t, uint32(types.LiquidityTier_LongTail), perpetual.Params.LiquidityTier)
					require.Equal(
						t, perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
						perpetual.Params.MarketType,
					)
				}
			},
		)
	}
}
