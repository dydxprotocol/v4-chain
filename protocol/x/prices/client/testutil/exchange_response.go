package testutil

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	client_types "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

// Response represents an API response from an exchange.
type Response struct {
	ResponseCode int
	Tickers      []JsonResponse
}

// JsonResponse is an interface that encapsulates a JSON API response.
type JsonResponse interface {
	toJson() interface{}
}

// ExchangeIdAndName represents an exchange with an ID and a name.
type ExchangeIdAndName struct {
	exchangeId   client_types.ExchangeId
	exchangeName string
}

// SetupExchangeResponses validates and sets up responses returned by exchange APIs using `gock`.
func SetupExchangeResponses(
	t testing.TB,
	responses map[ExchangeIdAndName]Response,
) {
	// Setup `gock` responses.
	for exchange, response := range responses {
		setupGockResponse(t, exchange, response)
	}
}

// setupGockResponse sets up the mock API responses returned by exchanges using `gock`.
func setupGockResponse(
	t testing.TB,
	exchange ExchangeIdAndName,
	response Response,
) {
	exchangeName := exchange.exchangeName

	var gockResponse *gock.Response
	switch exchangeName {
	case exchange_common.EXCHANGE_ID_BINANCE:
		gockResponse = NewGockBinanceResponse(
			response.ResponseCode,
			response.Tickers,
		)
	case exchange_common.EXCHANGE_ID_BINANCE_US:
		gockResponse = NewGockBinanceUSResponse(
			response.ResponseCode,
			response.Tickers,
		)
	case exchange_common.EXCHANGE_ID_BITFINEX:
		gockResponse = NewGockBitfinexResponse(
			response.ResponseCode,
			response.Tickers,
		)
	default:
		panic(fmt.Errorf("unsupported exchange: %s", exchangeName))
	}

	require.NotNil(t, gockResponse)
}
