package constants_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestStaticExchangeStartupConfigCache(t *testing.T) {
	tests := map[string]struct {
		// parameters
		exchangeId types.ExchangeFeedId

		// expectations
		expectedValue *types.ExchangeStartupConfig
		expectedFound bool
	}{
		"Get BINANCE exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_FEED_BINANCE,
			expectedValue: &types.ExchangeStartupConfig{
				ExchangeFeedId: exchange_common.EXCHANGE_FEED_BINANCE,
				IntervalMs:     2_000,
				TimeoutMs:      3_000,
				MaxQueries:     3,
			},
			expectedFound: true,
		},
		"Get BINANCEUS exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_FEED_BINANCE_US,
			expectedValue: &types.ExchangeStartupConfig{
				ExchangeFeedId: exchange_common.EXCHANGE_FEED_BINANCE_US,
				IntervalMs:     2_000,
				TimeoutMs:      3_000,
				MaxQueries:     3,
			},
			expectedFound: true,
		},
		"Get BITFINEX exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_FEED_BITFINEX,
			expectedValue: &types.ExchangeStartupConfig{
				ExchangeFeedId: exchange_common.EXCHANGE_FEED_BITFINEX,
				IntervalMs:     2_000,
				TimeoutMs:      3_000,
				MaxQueries:     2,
			},
			expectedFound: true,
		},
		"Get Kraken exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_FEED_KRAKEN,
			expectedValue: &types.ExchangeStartupConfig{
				ExchangeFeedId: exchange_common.EXCHANGE_FEED_KRAKEN,
				IntervalMs:     2_000,
				TimeoutMs:      3_000,
				MaxQueries:     1,
			},
			expectedFound: true,
		},
		"Get unknown exchangeDetails": {
			exchangeId:    99999999,
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := constants.StaticExchangeStartupConfig[tc.exchangeId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, ok, tc.expectedFound)
		})
	}
}

func TestStaticExchangeStartupConfigCacheLength(t *testing.T) {
	require.Len(t, constants.StaticExchangeStartupConfig, 4)
}
