package binance_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/binance"
	"github.com/stretchr/testify/require"
)

func TestBinanceUrl(t *testing.T) {
	require.Equal(t, "https://data-api.binance.vision/api/v3/ticker/24hr", binance.BinanceDetails.Url)
}

func TestBinanceUsUrl(t *testing.T) {
	require.Equal(t, "https://api.binance.us/api/v3/ticker/24hr", binance.BinanceUSDetails.Url)
}

func TestBinanceIsMultiMarket(t *testing.T) {
	require.True(t, binance.BinanceDetails.IsMultiMarket)
}

func TestBinanceUSIsMultiMarket(t *testing.T) {
	require.True(t, binance.BinanceUSDetails.IsMultiMarket)
}
