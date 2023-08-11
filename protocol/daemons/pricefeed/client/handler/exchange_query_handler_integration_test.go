//go:build exchange_tests

package handler

import (
	"context"
	"fmt"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/kraken"
	"net/http"
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/binance"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/dydxprotocol/v4/lib"
	"github.com/stretchr/testify/require"
)

func TestQueryingActualExchanges(t *testing.T) {
	tests := map[string]struct {
		// parameters
		url string
	}{
		"Binance": {
			url: CreateRequestUrl(binance.BinanceDetails.Url, []string{"ETHUSDC"}),
		},
		"BinanceUS": {
			url: CreateRequestUrl(binance.BinanceUSDetails.Url, []string{"BTCUSDC"}),
		},
		"Bitfinex": {
			url: CreateRequestUrl(bitfinex.BitfinexDetails.Url, []string{"tBTCUSD"}),
		},
		"Kraken": {
			url: CreateRequestUrl(kraken.KrakenDetails.Url, []string{"XXBTZUSD", "XETHZUSD", "LINKUSD"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			requestHandler := lib.NewRequestHandlerImpl(http.DefaultClient)

			response, err := requestHandler.Get(context.Background(), tc.url)

			if response.StatusCode != 200 {
				fmt.Println(response)
			}

			require.NoError(t, err)
			require.Equal(t, 200, response.StatusCode)
		})
	}
}
