package binance_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/binance"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestBinanceSymbols(t *testing.T) {
	tests := map[string]struct {
		// parameters
		marketId types.MarketId

		// expectations
		expectedValue string
		expectedFound bool
	}{
		"Get btcUsd information": {
			marketId:      exchange_common.MARKET_BTC_USD,
			expectedValue: "BTCUSDT",
			expectedFound: true,
		},
		"Get ethUsd information": {
			marketId:      exchange_common.MARKET_ETH_USD,
			expectedValue: "ETHUSDT",
			expectedFound: true,
		},
		"Get linkUsd information": {
			marketId:      exchange_common.MARKET_LINK_USD,
			expectedValue: "LINKUSDT",
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
			value, ok := binance.BinanceDetails.MarketSymbols[tc.marketId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestBinanceUSSymbols(t *testing.T) {
	tests := map[string]struct {
		// parameters
		marketId types.MarketId

		// expectations
		expectedValue string
		expectedFound bool
	}{
		"Get btcUsd information": {
			marketId:      exchange_common.MARKET_BTC_USD,
			expectedValue: "BTCUSD",
			expectedFound: true,
		},
		"Get ethUsd information": {
			marketId:      exchange_common.MARKET_ETH_USD,
			expectedValue: "ETHUSD",
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
			value, ok := binance.BinanceUSDetails.MarketSymbols[tc.marketId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestBinanceShareSymbolsLength(t *testing.T) {
	require.Len(t, binance.BinanceDetails.MarketSymbols, 3)
}

func TestBinanceUSShareSymbolsLength(t *testing.T) {
	require.Len(t, binance.BinanceUSDetails.MarketSymbols, 3)
}

func TestBinanceUrl(t *testing.T) {
	require.Equal(t, "https://data.binance.com/api/v3/ticker/24hr?symbol=$", binance.BinanceDetails.Url)
}

func TestBinanceUsUrl(t *testing.T) {
	require.Equal(t, "https://api.binance.us/api/v3/ticker/24hr?symbol=$", binance.BinanceUSDetails.Url)
}

func TestBinanceIsMultiMarket(t *testing.T) {
	require.False(t, binance.BinanceDetails.IsMultiMarket)
}

func TestBinanceUSIsMultiMarket(t *testing.T) {
	require.False(t, binance.BinanceUSDetails.IsMultiMarket)
}
