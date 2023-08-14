package okx_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/okx"
	"github.com/stretchr/testify/require"
)

func TestOkxUrl(t *testing.T) {
	require.Equal(t, "https://www.okx.com/api/v5/market/tickers?instType=SPOT", okx.OkxDetails.Url)
}

func TestOkxIsMultiMarket(t *testing.T) {
	require.True(t, okx.OkxDetails.IsMultiMarket)
}
