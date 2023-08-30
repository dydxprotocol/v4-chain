package crypto_com_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/crypto_com"
	"github.com/stretchr/testify/require"
)

func TestCryptoComUrl(t *testing.T) {
	require.Equal(t, "https://api.crypto.com/v2/public/get-ticker", crypto_com.CryptoComDetails.Url)
}

func TestCryptoComIsMultiMarket(t *testing.T) {
	require.True(t, crypto_com.CryptoComDetails.IsMultiMarket)
}
