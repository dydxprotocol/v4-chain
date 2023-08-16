package constants_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestStaticMarketNames(t *testing.T) {
	tests := map[string]struct {
		// parameters
		marketId types.MarketId

		// expectations
		expectedValue string
		expectedFound bool
	}{
		"Get BTC-USD name": {
			marketId:      exchange_common.MARKET_BTC_USD,
			expectedValue: "BTC-USD",
			expectedFound: true,
		},
		"Get ETH-USD name": {
			marketId:      exchange_common.MARKET_ETH_USD,
			expectedValue: "ETH-USD",
			expectedFound: true,
		},
		"Get LINK-USD name": {
			marketId:      exchange_common.MARKET_LINK_USD,
			expectedValue: "LINK-USD",
			expectedFound: true,
		},
		"Get unknown name": {
			marketId:      99999999,
			expectedValue: "",
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := constants.StaticMarketNames[tc.marketId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}

func TestStaticMarketNamesLength(t *testing.T) {
	require.Len(t, constants.StaticMarketNames, 34)
}
