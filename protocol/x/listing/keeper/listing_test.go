package keeper_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

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
			ticker:      "TEST-USD",
			expectedErr: nil,
		},
		"failure - market not found": {
			ticker:      "INVALID-USD",
			expectedErr: types.ErrMarketNotFound,
		},
		"failure - duplicate market": {
			ticker:      "BTC-USD",
			expectedErr: pricestypes.ErrMarketParamPairAlreadyExists,
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

				keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

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
					require.ErrorContains(t, err, tc.expectedErr.Error())
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
			referencePrice: 1000000000, // $1000
			expectedErr:    nil,
		},
		"failure - reference price 0": {
			ticker:         "TEST-USD",
			referencePrice: 0,
			expectedErr:    types.ErrReferencePriceZero,
		},
		"failure - market not found": {
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
					require.Equal(t, uint32(types.LiquidityTier_Isolated), perpetual.Params.LiquidityTier)
					require.Equal(
						t, perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
						perpetual.Params.MarketType,
					)
				}
			},
		)
	}
}

func TestCreateClobPair(t *testing.T) {
	tests := map[string]struct {
		ticker      string
		isDeliverTx bool
	}{
		"deliverTx - true": {
			ticker:      "TEST-USD",
			isDeliverTx: true,
		},
		"deliverTx - false": {
			ticker:      "TEST-USD",
			isDeliverTx: false,
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				mockIndexerEventManager := &mocks.IndexerEventManager{}
				ctx, keeper, _, _, pricesKeeper, perpetualsKeeper, clobKeeper, marketMapKeeper := keepertest.ListingKeepers(
					t,
					&mocks.BankKeeper{},
					mockIndexerEventManager,
				)
				mockIndexerEventManager.On(
					"AddTxnEvent",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return()
				keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, perpetualsKeeper, pricesKeeper, 10)

				// Set deliverTx mode
				if tc.isDeliverTx {
					ctx = ctx.WithIsCheckTx(false).WithIsReCheckTx(false)
					lib.AssertDeliverTxMode(ctx)
				} else {
					ctx = ctx.WithIsCheckTx(true)
					lib.AssertCheckTxMode(ctx)
				}

				// Create a marketmap with a single market
				dydxMetadata, err := tickermetadata.MarshalDyDx(
					tickermetadata.DyDx{
						ReferencePrice: 1000000000,
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
				require.NoError(t, err)

				perpetualId, err := keeper.CreatePerpetual(ctx, marketId, tc.ticker)
				require.NoError(t, err)

				clobPairId, err := keeper.CreateClobPair(ctx, perpetualId)
				require.NoError(t, err)

				// Check if the clob pair was created only if we are in deliverTx mode
				if tc.isDeliverTx {
					clobPair, found := clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(clobPairId))
					require.True(t, found)
					require.Equal(t, clobtypes.ClobPair_STATUS_ACTIVE, clobPair.Status)
					require.Equal(
						t,
						clobtypes.SubticksPerTick(types.SubticksPerTick_LongTail),
						clobPair.GetClobPairSubticksPerTick(),
					)
					require.Equal(
						t,
						types.DefaultStepBaseQuantums,
						clobPair.GetClobPairMinOrderBaseQuantums().ToUint64(),
					)
					require.Equal(t, perpetualId, clobPair.MustGetPerpetualId())
				} else {
					_, found := clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(clobPairId))
					require.False(t, found)
				}
			},
		)
	}
}
