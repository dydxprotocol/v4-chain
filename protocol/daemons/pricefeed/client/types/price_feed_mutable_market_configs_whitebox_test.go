package types

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed"
	prices_types "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
	"testing"
)

// Define a simple mock class to avoid import loops caused by importing the mock class.
type MockUpdater struct {
	Called            bool
	ExchangeId        ExchangeId
	NewExchangeConfig *MutableExchangeMarketConfig
	NewMarketConfigs  []*MutableMarketConfig
}

func (m *MockUpdater) GetExchangeId() ExchangeId { return m.ExchangeId }
func (m *MockUpdater) UpdateMutableExchangeConfig(
	newConfig *MutableExchangeMarketConfig,
	newMarketConfigs []*MutableMarketConfig,
) error {
	if m.Called {
		panic("UpdateMutableExchangeConfig called twice")
	}
	m.Called = true
	m.NewMarketConfigs = newMarketConfigs
	m.NewExchangeConfig = newConfig
	return nil
}

const (
	exchangeIdCoinbase = "coinbase"
	exchangeIdBinance  = "binance"
)

var (
	// Some pricefeed test constants need to be reproduced here in order to avoid import loops with the types package
	// since this is a whitebox test in the types package.
	coinbaseMutableExchangeConfig = &MutableExchangeMarketConfig{
		Id: exchangeIdCoinbase,
		MarketToMarketConfig: map[MarketId]MarketConfig{
			5: {
				Ticker: "BTC-USD",
			},
		},
	}
	updatedCoinbaseMutableExchangeConfig = &MutableExchangeMarketConfig{
		Id: exchangeIdCoinbase,
		MarketToMarketConfig: map[MarketId]MarketConfig{
			5: {
				Ticker: "BTC-USD",
			},
			7: {
				Ticker:         "SOL-USD",
				AdjustByMarket: newMarketIdWithValue(6),
			},
		},
	}
	binanceMutableExchangeConfig = &MutableExchangeMarketConfig{
		Id: exchangeIdBinance,
		MarketToMarketConfig: map[MarketId]MarketConfig{
			5: {
				Ticker: "BTCUSDT",
			},
			6: {
				Ticker: "ETHUSDT",
			},
		},
	}
	updatedBinanceMutableExchangeConfig = &MutableExchangeMarketConfig{
		Id: exchangeIdBinance,
		MarketToMarketConfig: map[MarketId]MarketConfig{
			5: {
				Ticker: "BTCUSDT",
			},
			6: {
				Ticker: "ETHUSDT",
			},
			7: {
				Ticker: "SOLUSDT",
			},
		},
	}

	testMutableExchangeMarketConfigs = map[string]*MutableExchangeMarketConfig{
		exchangeIdCoinbase: coinbaseMutableExchangeConfig,
		exchangeIdBinance:  binanceMutableExchangeConfig,
	}
	updatedBinanceExchangeMarketConfigs = map[string]*MutableExchangeMarketConfig{
		exchangeIdCoinbase: coinbaseMutableExchangeConfig,
		exchangeIdBinance:  updatedBinanceMutableExchangeConfig,
	}
	updatedCoinbaseExchangeMarketConfigs = map[string]*MutableExchangeMarketConfig{
		exchangeIdCoinbase: updatedCoinbaseMutableExchangeConfig,
		exchangeIdBinance:  binanceMutableExchangeConfig,
	}

	market5MutableMarketConfig = &MutableMarketConfig{
		Id:           5,
		Pair:         "BTC-USD",
		Exponent:     -5,
		MinExchanges: 1,
	}
	market6MutableMarketConfig = &MutableMarketConfig{
		Id:           6,
		Pair:         "ETH-USD",
		Exponent:     -6,
		MinExchanges: 2,
	}
	market7MutableMarketConfig = &MutableMarketConfig{
		Id:           7,
		Pair:         "SOL-USD",
		Exponent:     -7,
		MinExchanges: 1,
	}

	testMutableMarketConfigs = map[MarketId]*MutableMarketConfig{
		5: market5MutableMarketConfig,
		6: market6MutableMarketConfig,
	}
	updatedMutableMarketConfigs = map[MarketId]*MutableMarketConfig{
		5: market5MutableMarketConfig,
		6: market6MutableMarketConfig,
		7: market7MutableMarketConfig,
	}

	market5BinanceConfig                  = binanceConfig("BTCUSDT")
	market5CoinbaseConfig                 = `{"exchangeName":"coinbase","ticker":"BTC-USD"}`
	market6BinanceConfig                  = binanceConfig("ETHUSDT")
	market7BinanceConfig                  = binanceConfig("SOLUSDT")
	market7CoinbaseConfig_AdjustByMarket6 = `{"exchangeName":"coinbase","ticker":"SOL-USD","adjustByMarket":"ETH-USD"}`
	market8BinanceConfig                  = binanceConfig("DOTUSDT")
)

// binanceConfig returns a market exchange config for binance using the given ticker.
func binanceConfig(ticker string) string {
	return fmt.Sprintf(`{"exchangeName":"binance","ticker":"%s"}`, ticker)
}

func newMarketIdWithValue(id MarketId) *MarketId {
	ptr := new(MarketId)
	*ptr = id
	return ptr
}

func TestAddExchangeConfigUpdater(t *testing.T) {
	pfmmc := NewPriceFeedMutableMarketConfigs([]ExchangeId{exchangeIdCoinbase})

	mockPriceFetcher := MockUpdater{ExchangeId: exchangeIdCoinbase}
	pfmmc.AddPriceFetcher(&mockPriceFetcher)

	mockPriceEncoder := MockUpdater{ExchangeId: exchangeIdCoinbase}
	pfmmc.AddPriceEncoder(&mockPriceEncoder)

	exchangeConfigUpdaters, ok := pfmmc.mutableExchangeConfigUpdaters[exchangeIdCoinbase]
	require.True(t, ok)
	require.Equal(t, &mockPriceFetcher, exchangeConfigUpdaters.PriceFetcher)
	require.Equal(t, &mockPriceEncoder, exchangeConfigUpdaters.PriceEncoder)
}

func TestUpdateMarkets_Mixed(t *testing.T) {
	tests := map[string]struct {
		marketParams           []prices_types.MarketParam
		updatedExchangeConfigs map[ExchangeId]*MutableExchangeMarketConfig
		updatedMarketConfigs   map[MarketId]*MutableMarketConfig
		expectedUpdates        map[ExchangeId]struct {
			updatedExchangeConfig *MutableExchangeMarketConfig
			updatedMarketConfigs  []*MutableMarketConfig
		}
		expectedError             error
		expectedMarketParamErrors map[MarketId]error
	}{
		"Error: market params nil": {
			expectedError: errors.New(
				"UpdateMarkets: marketParams cannot be nil",
			),
		},
		"Success: No updates": {
			marketParams: []prices_types.MarketParam{
				{
					Id:                5,
					Exponent:          -5,
					Pair:              "BTC-USD",
					MinExchanges:      1,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						market5CoinbaseConfig,
						market5BinanceConfig,
					),
				},
				{
					Id:                6,
					Exponent:          -6,
					Pair:              "ETH-USD",
					MinExchanges:      2,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s]}`,
						market6BinanceConfig,
					),
				},
			},
			updatedExchangeConfigs: testMutableExchangeMarketConfigs,
			updatedMarketConfigs:   testMutableMarketConfigs,
		},
		"Success: Added market to 1 exchange": {
			marketParams: []prices_types.MarketParam{
				{
					Id:                5,
					Exponent:          -5,
					Pair:              "BTC-USD",
					MinExchanges:      1,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						market5CoinbaseConfig,
						market5BinanceConfig,
					),
				},
				{
					Id:                6,
					Exponent:          -6,
					Pair:              "ETH-USD",
					MinExchanges:      2,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s]}`,
						market6BinanceConfig,
					),
				},
				// Add market 7 to binance
				{
					Id:                 7,
					Exponent:           -7,
					Pair:               "SOL-USD",
					MinExchanges:       1,
					MinPriceChangePpm:  1,
					ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, market7BinanceConfig),
				},
			},
			updatedExchangeConfigs: updatedBinanceExchangeMarketConfigs,
			updatedMarketConfigs:   updatedMutableMarketConfigs,
			expectedUpdates: map[ExchangeId]struct {
				updatedExchangeConfig *MutableExchangeMarketConfig
				updatedMarketConfigs  []*MutableMarketConfig
			}{
				exchangeIdBinance: {
					updatedExchangeConfig: updatedBinanceMutableExchangeConfig,
					updatedMarketConfigs: []*MutableMarketConfig{
						market5MutableMarketConfig,
						market6MutableMarketConfig,
						market7MutableMarketConfig,
					},
				},
			},
		},
		"Partial update - 1 market add succeeds, 1 fails, existing markets retained": {
			marketParams: []prices_types.MarketParam{
				{
					Id:                5,
					Exponent:          -5,
					Pair:              "BTC-USD",
					MinExchanges:      1,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						market5CoinbaseConfig,
						market5BinanceConfig,
					),
				},
				{
					Id:                6,
					Exponent:          -6,
					Pair:              "ETH-USD",
					MinExchanges:      2,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s]}`,
						market6BinanceConfig,
					),
				},
				// Add market 7 to binance
				{
					Id:                 7,
					Exponent:           -7,
					Pair:               "SOL-USD",
					MinExchanges:       1,
					MinPriceChangePpm:  1,
					ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, market7BinanceConfig),
				},
				// Market 8 will fail to add because it has an invalid MinPriceChangePpm
				{
					Id:                 8,
					Exponent:           -8,
					Pair:               "DOT-USD",
					MinExchanges:       1,
					MinPriceChangePpm:  0, // Invalid
					ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, market8BinanceConfig),
				},
			},
			expectedMarketParamErrors: map[MarketId]error{
				8: errors.New(
					"invalid market param 8: Min price change in parts-per-million must be greater than 0 " +
						"and less than 10000",
				),
			},
			updatedExchangeConfigs: updatedBinanceExchangeMarketConfigs,
			updatedMarketConfigs:   updatedMutableMarketConfigs,
			expectedUpdates: map[ExchangeId]struct {
				updatedExchangeConfig *MutableExchangeMarketConfig
				updatedMarketConfigs  []*MutableMarketConfig
			}{
				exchangeIdBinance: {
					updatedExchangeConfig: updatedBinanceMutableExchangeConfig,
					updatedMarketConfigs: []*MutableMarketConfig{
						market5MutableMarketConfig,
						market6MutableMarketConfig,
						market7MutableMarketConfig,
					},
				},
			},
		},
		"Success - Added market with un-supported adjustment market to 1 exchange": {
			marketParams: []prices_types.MarketParam{
				{
					Id:                5,
					Exponent:          -5,
					Pair:              "BTC-USD",
					MinExchanges:      1,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						market5CoinbaseConfig,
						market5BinanceConfig,
					),
				},
				{
					Id:                6,
					Exponent:          -6,
					Pair:              "ETH-USD",
					MinExchanges:      2,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s]}`,
						market6BinanceConfig,
					),
				},
				// Add market 7 to coinbase with an adjustment market of 6.
				{
					Id:                7,
					Exponent:          -7,
					Pair:              "SOL-USD",
					MinExchanges:      1,
					MinPriceChangePpm: 1,
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s]}`,
						market7CoinbaseConfig_AdjustByMarket6,
					),
				},
			},
			updatedExchangeConfigs: updatedCoinbaseExchangeMarketConfigs,
			updatedMarketConfigs:   updatedMutableMarketConfigs,
			expectedUpdates: map[ExchangeId]struct {
				updatedExchangeConfig *MutableExchangeMarketConfig
				updatedMarketConfigs  []*MutableMarketConfig
			}{
				exchangeIdCoinbase: {
					updatedExchangeConfig: updatedCoinbaseMutableExchangeConfig,
					updatedMarketConfigs: []*MutableMarketConfig{
						market5MutableMarketConfig,
						market6MutableMarketConfig,
						market7MutableMarketConfig,
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			// Initialize the price feed mutable market config manually.
			pfmmc := PricefeedMutableMarketConfigsImpl{
				mutableExchangeToConfigs:      testMutableExchangeMarketConfigs,
				mutableMarketToConfigs:        testMutableMarketConfigs,
				mutableExchangeConfigUpdaters: map[ExchangeId]UpdatersForExchange{},
			}

			// Create mock updaters for each exchange and add them to the price feed mutable market config.
			// These should receive updates when an exchange config changes.
			exchangeToUpdaters := make(map[ExchangeId]struct {
				PriceEncoder *MockUpdater
				PriceFetcher *MockUpdater
			})
			for exchangeId := range testMutableExchangeMarketConfigs {
				// Increment the wait group to ensure that registering the updaters does not cause the waitgroup
				// to become negative.
				pfmmc.updatersInitialized.Add(2)

				updatersForExchange := struct {
					PriceEncoder *MockUpdater
					PriceFetcher *MockUpdater
				}{
					PriceEncoder: &MockUpdater{ExchangeId: exchangeId},
					PriceFetcher: &MockUpdater{ExchangeId: exchangeId},
				}
				exchangeToUpdaters[exchangeId] = updatersForExchange
				pfmmc.AddPriceFetcher(updatersForExchange.PriceFetcher)
				pfmmc.AddPriceEncoder(updatersForExchange.PriceEncoder)
			}

			// Take a snapshot of the old mutable state to validate updates occurred.
			oldExchangeConfigsSnapshot := make(map[string]*MutableExchangeMarketConfig, len(pfmmc.mutableExchangeToConfigs))
			for exchangeId, exchangeConfig := range pfmmc.mutableExchangeToConfigs {
				oldExchangeConfigsSnapshot[exchangeId] = exchangeConfig.Copy()
			}
			oldMarketConfigsSnapshot := make(map[MarketId]*MutableMarketConfig, len(pfmmc.mutableMarketToConfigs))
			for marketId, marketConfig := range pfmmc.mutableMarketToConfigs {
				oldMarketConfigsSnapshot[marketId] = marketConfig.Copy()
			}

			// Execute the update.
			marketParamErrors, err := pfmmc.UpdateMarkets(tc.marketParams)

			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
				require.Empty(t, marketParamErrors)

				// Under normal circumstances, we expect that errors will come from validation and price fetchers and
				// encoders will not be updated.
				for _, exchangeFeedUpdaters := range exchangeToUpdaters {
					require.False(t, exchangeFeedUpdaters.PriceEncoder.Called)
					require.False(t, exchangeFeedUpdaters.PriceFetcher.Called)
				}

				// PricefeedMutableMarketConfigsImpl should not have updated its internal state.
				exchangeConfigsEqual(t, oldExchangeConfigsSnapshot, pfmmc.mutableExchangeToConfigs)
				marketConfigsEqual(t, oldMarketConfigsSnapshot, pfmmc.mutableMarketToConfigs)
			} else {
				require.Nil(t, err)
				pricefeed.MarketParamErrorsEqual(t, tc.expectedMarketParamErrors, marketParamErrors)

				// If the exchange config was updated, expect that the appropriate updaters were updated.
				// Otherwise, each updater should be untouched.
				for exchangeId, exchangeUpdaters := range exchangeToUpdaters {
					updaters := []*MockUpdater{exchangeUpdaters.PriceEncoder, exchangeUpdaters.PriceFetcher}

					// Keep track of each mutable exchange config to ensure that the config owned by the pfmmc and the
					// config used to update each updater is unique. That way, we can ensure that only copies of pfmmc
					// state are used to update the updaters.
					uniqueMutableExchangeConfigs := map[*MutableExchangeMarketConfig]struct{}{}
					uniqueMutableExchangeConfigs[pfmmc.mutableExchangeToConfigs[exchangeId]] = struct{}{}

					// Likewise, keep track of each mutable market config in order to ensure that the configs owned
					// by the pfmmc and the configs sent to each updater are also unique, and that only copies of pfmmc
					// mutable market state are used to update the updaters.
					uniqueMutableMarketConfigs := map[*MutableMarketConfig]struct{}{}
					for _, marketConfig := range pfmmc.mutableMarketToConfigs {
						uniqueMutableMarketConfigs[marketConfig] = struct{}{}
					}

					if parameters, ok := tc.expectedUpdates[exchangeId]; ok {
						for _, updater := range updaters {
							require.True(t, updater.Called)

							// Expect that the parameters of the update match the expected update values.
							marketConfigsSliceEqual(t, parameters.updatedMarketConfigs, updater.NewMarketConfigs)
							require.True(t, parameters.updatedExchangeConfig.Equal(updater.NewExchangeConfig))

							uniqueMutableExchangeConfigs[updater.NewExchangeConfig] = struct{}{}
							for _, marketConfig := range updater.NewMarketConfigs {
								uniqueMutableMarketConfigs[marketConfig] = struct{}{}
							}
						}

						require.Len(
							t,
							uniqueMutableExchangeConfigs,
							len(pfmmc.mutableExchangeToConfigs)+1,
							"Expected a new copy of the exchange config to be created for each updater",
						)

						require.Len(
							t,
							uniqueMutableMarketConfigs,
							len(parameters.updatedMarketConfigs)*3,
							"Expected a new copy of each market config to be created for each updater",
						)
					} else {
						for _, updater := range updaters {
							// Expect updaters are not updated when the exchange config does not change.
							require.False(t, updater.Called)
						}
					}
				}

				// PricefeedMutableMarketConfigsImpl should have updated its internal state.
				exchangeConfigsEqual(t, tc.updatedExchangeConfigs, pfmmc.mutableExchangeToConfigs)
				marketConfigsEqual(t, tc.updatedMarketConfigs, pfmmc.mutableMarketToConfigs)
			}
		})
	}
}

// exchangeConfigsEqual compares two maps of exchange configs for equality. This was placed in this file to avoid
// import cycles.
func exchangeConfigsEqual(
	t *testing.T,
	expected map[ExchangeId]*MutableExchangeMarketConfig,
	actual map[ExchangeId]*MutableExchangeMarketConfig,
) {
	require.ElementsMatch(
		t,
		maps.Keys(expected),
		maps.Keys(actual),
	)
	for exchangeId, expectedExchangeConfig := range expected {
		actualExchangeConfig, ok := actual[exchangeId]
		require.True(t, ok)
		require.True(t, expectedExchangeConfig.Equal(actualExchangeConfig))
	}
}

// marketConfigsEqual compares two maps of market configs for equality. This was placed in this file to avoid
// import cycles.
func marketConfigsEqual(
	t *testing.T,
	expected map[MarketId]*MutableMarketConfig,
	actual map[MarketId]*MutableMarketConfig,
) {
	require.ElementsMatch(
		t,
		maps.Keys(expected),
		maps.Keys(actual),
	)
	for marketId, expectedMarketConfig := range expected {
		actualMarketConfig, ok := actual[marketId]
		require.True(t, ok)
		require.Equal(t, *expectedMarketConfig, *actualMarketConfig)
	}
}

// marketConfigsSliceEqual compares two slices of market configs for equality.
func marketConfigsSliceEqual(t *testing.T,
	expected []*MutableMarketConfig,
	actual []*MutableMarketConfig,
) {
	require.Equal(t, len(expected), len(actual), "market config slice lengths do not match")

	for i, expectedConfig := range expected {
		require.Equal(t, *expectedConfig, *actual[i])
	}
}
