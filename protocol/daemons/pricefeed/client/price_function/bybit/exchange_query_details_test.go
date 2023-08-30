package bybit_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bybit"
	"github.com/stretchr/testify/require"
)

func TestBybitUrl(t *testing.T) {
	require.Equal(t, "https://api.bybit.com/v5/market/tickers?category=spot", bybit.BybitDetails.Url)
}

func TestBybitIsMultiMarket(t *testing.T) {
	require.True(t, bybit.BybitDetails.IsMultiMarket)
}
