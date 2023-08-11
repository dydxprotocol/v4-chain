package bitfinex_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestBitfinexSymbols(t *testing.T) {
	tests := map[string]struct {
		// parameters
		marketId types.MarketId

		// expectations
		expectedValue string
		expectedFound bool
	}{
		"Get btcUsd information": {
			marketId:      exchange_common.MARKET_BTC_USD,
			expectedValue: "tBTCUSD",
			expectedFound: true,
		},
		"Get ethUsd information": {
			marketId:      exchange_common.MARKET_ETH_USD,
			expectedValue: "tETHUSD",
			expectedFound: true,
		},
		"Get unknown information": {
			marketId:      99999999,
			expectedValue: "",
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := bitfinex.BitfinexDetails.MarketSymbols[tc.marketId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestBitfinexSymbolsLength(t *testing.T) {
	require.Len(t, bitfinex.BitfinexDetails.MarketSymbols, 2)
}

func TestBitfinexUrl(t *testing.T) {
	require.Equal(t, "https://api.bitfinex.com/v2/ticker/$", bitfinex.BitfinexDetails.Url)
}

func TestBitfinexIsMultiMarket(t *testing.T) {
	require.False(t, bitfinex.BitfinexDetails.IsMultiMarket)
}
