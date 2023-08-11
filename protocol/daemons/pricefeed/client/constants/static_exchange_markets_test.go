package constants_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestStaticExchangeMarketsCache(t *testing.T) {
	tests := map[string]struct {
		// parameters
		exchangeId types.ExchangeFeedId

		// expectations
		expectedValue []types.MarketId
		expectedFound bool
	}{
		"Get BINANCE exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_FEED_BINANCE,
			expectedValue: []types.MarketId{
				exchange_common.MARKET_BTC_USD,
				exchange_common.MARKET_ETH_USD,
				exchange_common.MARKET_LINK_USD,
			},
			expectedFound: true,
		},
		"Get BINANCEUS exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_FEED_BINANCE_US,
			expectedValue: []types.MarketId{
				exchange_common.MARKET_BTC_USD,
				exchange_common.MARKET_ETH_USD,
				exchange_common.MARKET_LINK_USD,
			},
			expectedFound: true,
		},
		"Get Bitfinex exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_FEED_BITFINEX,
			expectedValue: []types.MarketId{
				exchange_common.MARKET_BTC_USD,
				exchange_common.MARKET_ETH_USD,
			},
			expectedFound: true,
		},
		"Get Kraken exchangeDetails": {
			exchangeId: exchange_common.EXCHANGE_FEED_KRAKEN,
			expectedValue: []types.MarketId{
				exchange_common.MARKET_BTC_USD,
				exchange_common.MARKET_ETH_USD,
				exchange_common.MARKET_LINK_USD,
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
			value, ok := constants.StaticExchangeMarkets[tc.exchangeId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestStaticExchangeMarketsCacheLength(t *testing.T) {
	require.Len(t, constants.StaticExchangeMarkets, 4)
}
