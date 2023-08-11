package constants_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
)

func TestMarketPriceExponentCache(t *testing.T) {
	tests := map[string]struct {
		// parameters
		marketId types.MarketId

		// expectations
		expectedValue int32
		expectedFound bool
	}{
		"Get BTCUSD exponent": {
			marketId:      exchange_common.MARKET_BTC_USD,
			expectedValue: -5,
			expectedFound: true,
		},
		"Get WETHUSD exponent": {
			marketId:      exchange_common.MARKET_ETH_USD,
			expectedValue: -6,
			expectedFound: true,
		},
		"Get LINKUSD exponent": {
			marketId:      exchange_common.MARKET_LINK_USD,
			expectedValue: -8,
			expectedFound: true,
		},
		"Get unknown exponent": {
			marketId:      99999999,
			expectedValue: 0,
			expectedFound: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := constants.StaticMarketPriceExponent[tc.marketId]
			require.Equal(t, tc.expectedValue, value)
			require.Equal(t, tc.expectedFound, ok)
		})
	}
}
