package bitfinex_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/stretchr/testify/require"
)

func TestBitfinexUrl(t *testing.T) {
	require.Equal(t, "https://api-pub.bitfinex.com/v2/tickers?symbols=ALL", bitfinex.BitfinexDetails.Url)
}

func TestBitfinexIsMultiMarket(t *testing.T) {
	require.True(t, bitfinex.BitfinexDetails.IsMultiMarket)
}
