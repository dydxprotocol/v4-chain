package metrics_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/stretchr/testify/require"
)

func TestStaticMarketSymbols(t *testing.T) {
	tests := map[string]struct {
		// parameters
		marketId types.MarketId

		// expectations
		expectedValue string
		expectedFound bool
	}{
		"Get BTCUSD symbol": {
			marketId:      exchange_common.MARKET_BTC_USD,
			expectedValue: "BTCUSD",
			expectedFound: true,
		},
		"Get ETHUSD symbol": {
			marketId:      exchange_common.MARKET_ETH_USD,
			expectedValue: "ETHUSD",
			expectedFound: true,
		},
		"Get LINKUSD symbol": {
			marketId:      exchange_common.MARKET_LINK_USD,
			expectedValue: "LINKUSD",
			expectedFound: true,
		},
		"Get unknown symbol": {
			marketId:      99999999,
			expectedValue: "",
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := metrics.StaticMarketSymbols[tc.marketId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestStaticMarketSymbolsLength(t *testing.T) {
	require.Len(t, metrics.StaticMarketSymbols, 3)
}
