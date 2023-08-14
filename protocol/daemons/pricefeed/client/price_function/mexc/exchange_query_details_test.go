package mexc_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/mexc"
	"github.com/stretchr/testify/require"
)

func TestMexcUrl(t *testing.T) {
	require.Equal(t, "https://www.mexc.com/open/api/v2/market/ticker", mexc.MexcDetails.Url)
}

func TestMexcIsMultiMarket(t *testing.T) {
	require.True(t, mexc.MexcDetails.IsMultiMarket)
}
