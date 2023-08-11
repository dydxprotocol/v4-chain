package constants_test

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/kraken"
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/binance"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestStaticExchangeDetailsCache(t *testing.T) {
	tests := map[string]struct {
		// parameters
		exchangeId types.ExchangeFeedId

		// expectations
		expectedValue types.ExchangeQueryDetails
		expectedFound bool
	}{
		"Get BINANCE exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_FEED_BINANCE,
			expectedValue: binance.BinanceDetails,
			expectedFound: true,
		},
		"Get BINANCEUS exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_FEED_BINANCE_US,
			expectedValue: binance.BinanceUSDetails,
			expectedFound: true,
		},
		"Get Bitfinex exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_FEED_BITFINEX,
			expectedValue: bitfinex.BitfinexDetails,
			expectedFound: true,
		},
		"Get Kraken exchangeDetails": {
			exchangeId:    exchange_common.EXCHANGE_FEED_KRAKEN,
			expectedValue: kraken.KrakenDetails,
			expectedFound: true,
		},
		"Get unknown exchangeDetails": {
			exchangeId:    99999999,
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := constants.StaticExchangeDetails[tc.exchangeId]
			require.Equal(t, tc.expectedValue.Exchange, value.Exchange)
			require.Equal(t, tc.expectedValue.Url, value.Url)
			require.Equal(t, tc.expectedValue.MarketSymbols, value.MarketSymbols)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestStaticExchangeDetailsCacheLength(t *testing.T) {
	require.Len(t, constants.StaticExchangeDetails, 4)
}
