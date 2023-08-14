package metrics_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/stretchr/testify/require"
)

func TestStaticMarketPairs(t *testing.T) {
	tests := map[string]struct {
		// parameters
		marketId types.MarketId

		// expectations
		expectedValue string
		expectedFound bool
	}{
		"Get BTCUSD pair": {
			marketId:      exchange_common.MARKET_BTC_USD,
			expectedValue: "BTCUSD",
			expectedFound: true,
		},
		"Get ETHUSD pair": {
			marketId:      exchange_common.MARKET_ETH_USD,
			expectedValue: "ETHUSD",
			expectedFound: true,
		},
		"Get LINKUSD pair": {
			marketId:      exchange_common.MARKET_LINK_USD,
			expectedValue: "LINKUSD",
			expectedFound: true,
		},
		"Get unknown pair": {
			marketId:      99999999,
			expectedValue: "",
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := metrics.StaticMarketPairs[tc.marketId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestStaticMarketPairsLength(t *testing.T) {
	require.Len(t, metrics.StaticMarketPairs, 34)
}
