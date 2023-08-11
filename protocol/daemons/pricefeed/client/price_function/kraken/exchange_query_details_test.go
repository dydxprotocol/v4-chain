package kraken_test

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/kraken"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKrakenSymbols(t *testing.T) {
	tests := map[string]struct {
		// parameters
		marketId types.MarketId

		// expectations
		expectedValue string
		expectedFound bool
	}{
		"Get btcUsd information": {
			marketId:      exchange_common.MARKET_BTC_USD,
			expectedValue: "XXBTZUSD",
			expectedFound: true,
		},
		"Get ethUsd information": {
			marketId:      exchange_common.MARKET_ETH_USD,
			expectedValue: "XETHZUSD",
			expectedFound: true,
		},
		"Get linkUsd information": {
			marketId:      exchange_common.MARKET_LINK_USD,
			expectedValue: "LINKUSD",
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
			value, ok := kraken.KrakenDetails.MarketSymbols[tc.marketId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestKrakenShareSymbolsLength(t *testing.T) {
	require.Len(t, kraken.KrakenDetails.MarketSymbols, 3)
}

func TestKrakenUrl(t *testing.T) {
	require.Equal(t, "https://api.kraken.com/0/public/Ticker?pair=$", kraken.KrakenDetails.Url)
}

func TestKrakenIsMultiMarket(t *testing.T) {
	require.True(t, kraken.KrakenDetails.IsMultiMarket)
}
