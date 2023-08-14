package bitstamp_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bitstamp"
	"github.com/stretchr/testify/require"
)

func TestBitstampUrl(t *testing.T) {
	require.Equal(t, "https://www.bitstamp.net/api/v2/ticker/", bitstamp.BitstampDetails.Url)
}

func TestBitstampIsMultiMarket(t *testing.T) {
	require.True(t, bitstamp.BitstampDetails.IsMultiMarket)
}
