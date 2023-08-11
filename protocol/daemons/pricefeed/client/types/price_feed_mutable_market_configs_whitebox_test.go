package types

import (
	"errors"
	"fmt"
	prices_types "github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
	"testing"
)

// Define a simple mock class to avoid import loops.
type MockPriceFetcher struct {
	Called            bool
	ExchangeId        ExchangeId
	NewExchangeConfig *MutableExchangeMarketConfig
	NewMarketConfigs  []*MutableMarketConfig
}

func (m *MockPriceFetcher) GetExchangeId() ExchangeId { return m.ExchangeId }
func (m *MockPriceFetcher) UpdateMutableExchangeConfig(
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
		MarketToTicker: map[uint32]string{
			5: "BTC-USD",
			6: "ETH-USD",
		},
	}
	updatedBinanceMutableExchangeConfig = &MutableExchangeMarketConfig{
		Id: exchangeIdBinance,
		MarketToTicker: map[uint32]string{
			5: "BTCUSDT",
			6: "ETHUSDT",
			7: "SOLUSDT",
		},
	}
	binanceMutableExchangeConfig = &MutableExchangeMarketConfig{
		Id: exchangeIdBinance,
		MarketToTicker: map[uint32]string{
			5: "BTCUSDT",
			6: "ETHUSDT",
		},
	}

	testMutableExchangeMarketConfigs = map[string]*MutableExchangeMarketConfig{
		exchangeIdCoinbase: coinbaseMutableExchangeConfig,
		exchangeIdBinance:  binanceMutableExchangeConfig,
	}
	updatedMutableExchangeMarketConfigs = map[string]*MutableExchangeMarketConfig{
		exchangeIdCoinbase: coinbaseMutableExchangeConfig,
		exchangeIdBinance:  updatedBinanceMutableExchangeConfig,
	}

	market5MutableMarketConfig = &MutableMarketConfig{
		Id:       5,
		Pair:     "BTC-USD",
		Exponent: -5,
	}
	market6MutableMarketConfig = &MutableMarketConfig{
		Id:       6,
		Pair:     "ETH-USD",
		Exponent: -6,
	}
	market7MutableMarketConfig = &MutableMarketConfig{
		Id:       7,
		Pair:     "SOL-USD",
		Exponent: -7,
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

	market5BinanceConfig  = `{"exchangeName":"binance","ticker":"BTCUSDT"}`
	market5CoinbaseConfig = `{"exchangeName":"coinbase","ticker":"BTC-USD"}`
	market6BinanceConfig  = `{"exchangeName":"binance","ticker":"ETHUSDT"}`
	market6CoinbaseConfig = `{"exchangeName":"coinbase","ticker":"ETH-USD"}`
	market7BinanceConfig  = `{"exchangeName":"binance","ticker":"SOLUSDT"}`
)

func TestAddExchangeConfigUpdater(t *testing.T) {
	pfmc := NewPriceFeedMutableMarketConfigs(nil)

	mockPriceFetcher := MockPriceFetcher{ExchangeId: exchangeIdCoinbase}
	pfmc.AddExchangeConfigUpdater(&mockPriceFetcher)

	exchangeConfigUpdater, ok := pfmc.mutableExchangeConfigUpdaters[exchangeIdCoinbase]
	require.True(t, ok)
	require.Equal(t, &mockPriceFetcher, exchangeConfigUpdater)
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
		expectedError error
	}{
		"Error: invalid market params": {
			marketParams: []prices_types.MarketParam{
				{
					// Empty pair is invalid
					Pair: "",
				},
			},
			expectedError: errors.New(
				"UpdateMarkets market param validation failed: invalid market param 0: pair cannot be empty",
			),
		},
		"Success: No updates": {
			marketParams: []prices_types.MarketParam{
				{
					Id:       5,
					Exponent: -5,
					Pair:     "BTC-USD",
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						market5CoinbaseConfig,
						market5BinanceConfig,
					),
				},
				{
					Id:       6,
					Exponent: -6,
					Pair:     "ETH-USD",
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						market6CoinbaseConfig,
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
					Id:       5,
					Exponent: -5,
					Pair:     "BTC-USD",
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						market5CoinbaseConfig,
						market5BinanceConfig,
					),
				},
				{
					Id:       6,
					Exponent: -6,
					Pair:     "ETH-USD",
					ExchangeConfigJson: fmt.Sprintf(
						`{"exchanges":[%s,%s]}`,
						market6CoinbaseConfig,
						market6BinanceConfig,
					),
				},
				// Add market 7 to binance
				{
					Id:                 7,
					Exponent:           -7,
					Pair:               "SOL-USD",
					ExchangeConfigJson: fmt.Sprintf(`{"exchanges":[%s]}`, market7BinanceConfig),
				},
			},
			updatedExchangeConfigs: updatedMutableExchangeMarketConfigs,
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
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			// Initialize the price feed mutable market config.
			pfmmc := PricefeedMutableMarketConfigsImpl{
				mutableExchangeToConfigs:      testMutableExchangeMarketConfigs,
				mutableMarketToConfigs:        testMutableMarketConfigs,
				mutableExchangeConfigUpdaters: map[ExchangeId]ExchangeConfigUpdater{},
			}

			// Create mock fetchers for each exchange and add them to the price feed mutable market config.
			// These fetchers should receive updates when an exchange config changes.
			priceFetchers := make(map[ExchangeId]*MockPriceFetcher)
			for exchangeId := range testMutableExchangeMarketConfigs {
				priceFetcher := MockPriceFetcher{ExchangeId: exchangeId}
				priceFetchers[exchangeId] = &priceFetcher
				pfmmc.AddExchangeConfigUpdater(&priceFetcher)
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
			err := pfmmc.UpdateMarkets(tc.marketParams)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())

				// Under normal circumstances, we expect that errors will come from validation and price fetchers
				// will not be updated.
				for _, priceFetcher := range priceFetchers {
					require.False(t, priceFetcher.Called)
				}

				// PricefeedMutableMarketConfigsImpl should not have updated its internal state.
				exchangeConfigsEqual(t, oldExchangeConfigsSnapshot, pfmmc.mutableExchangeToConfigs)
				marketConfigsEqual(t, oldMarketConfigsSnapshot, pfmmc.mutableMarketToConfigs)
			} else {
				require.Nil(t, err)

				// If the exchange config was updated, expect that the appropriate price fetchers were updated.
				// Otherwise, the price fetcher should be untouched.
				for exchangeId, priceFetcher := range priceFetchers {
					if parameters, ok := tc.expectedUpdates[exchangeId]; ok {
						require.True(t, priceFetcher.Called)

						// Expect that a copy of the new state was sent to the price fetcher.
						require.NotSame(t, priceFetcher.NewExchangeConfig, pfmmc.mutableExchangeToConfigs[exchangeId])
						require.EqualValues(t, priceFetcher.NewExchangeConfig, pfmmc.mutableExchangeToConfigs[exchangeId])

						// Expect that price fetcher update parameters match the expected update values.
						require.EqualValues(t, parameters.updatedMarketConfigs, priceFetcher.NewMarketConfigs)
						require.EqualValues(t, parameters.updatedExchangeConfig, priceFetcher.NewExchangeConfig)
					} else {
						// Expect price fetcher was not updated when the exchange config does not change.
						require.False(t, priceFetcher.Called)
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
