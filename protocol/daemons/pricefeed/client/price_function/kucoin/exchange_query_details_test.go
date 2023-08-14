package kucoin_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/kucoin"
	"github.com/stretchr/testify/require"
)

func TestKucoinUrl(t *testing.T) {
	require.Equal(t, "https://api.kucoin.com/api/v1/market/allTickers", kucoin.KucoinDetails.Url)
}

func TestKucoinIsMultiMarket(t *testing.T) {
	require.True(t, kucoin.KucoinDetails.IsMultiMarket)
}
