package metrics_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/stretchr/testify/require"
)

func TestStaticExchangeNames(t *testing.T) {
	tests := map[string]struct {
		// parameters
		exchangeFeedId types.ExchangeFeedId

		// expectations
		expectedValue string
		expectedFound bool
	}{
		"Get Binance name": {
			exchangeFeedId: exchange_common.EXCHANGE_FEED_BINANCE,
			expectedValue:  exchange_common.EXCHANGE_NAME_BINANCE,
			expectedFound:  true,
		},
		"Get BinanceUS name": {
			exchangeFeedId: exchange_common.EXCHANGE_FEED_BINANCE_US,
			expectedValue:  exchange_common.EXCHANGE_NAME_BINANCEUS,
			expectedFound:  true,
		},
		"Get Bitfinex name": {
			exchangeFeedId: exchange_common.EXCHANGE_FEED_BITFINEX,
			expectedValue:  exchange_common.EXCHANGE_NAME_BITFINEX,
			expectedFound:  true,
		},
		"Get Kraken name": {
			exchangeFeedId: exchange_common.EXCHANGE_FEED_KRAKEN,
			expectedValue:  exchange_common.EXCHANGE_NAME_KRAKEN,
			expectedFound:  true,
		},
		"Get unknown symbol": {
			exchangeFeedId: 99999999,
			expectedValue:  "",
			expectedFound:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := metrics.StaticExchangeNames[tc.exchangeFeedId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestStaticExchangeNamesLength(t *testing.T) {
	require.Len(t, metrics.StaticExchangeNames, 4)
}
