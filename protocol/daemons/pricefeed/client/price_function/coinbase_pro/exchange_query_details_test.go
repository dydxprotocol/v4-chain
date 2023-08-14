package coinbase_pro_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/coinbase_pro"
	"github.com/stretchr/testify/require"
)

func TestCoinbaseProUrl(t *testing.T) {
	require.Equal(t, "https://api.pro.coinbase.com/products/$/ticker", coinbase_pro.CoinbaseProDetails.Url)
}

func TestCoinbaseProIsMultiMarket(t *testing.T) {
	require.False(t, coinbase_pro.CoinbaseProDetails.IsMultiMarket)
}
