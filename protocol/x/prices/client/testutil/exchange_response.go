package testutil

import (
	"fmt"
	"sort"
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	client_types "github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

// MarketToExchangeResponse is a map of `Market` to API responses from exchanges.
type MarketToExchangeResponse struct {
	marketToExchangesResponse map[MarketIdAndName]ExchangeResponse
}

// ExchangeResponse is a map of `Exchange` to API response from an exchange.
type ExchangeResponse struct {
	exchangeToResponse map[ExchangeIdAndName]Response
}

// Response represents an API response from an exchange.
type Response struct {
	ResponseCode int
	Response     JsonResponse
}

// JsonResponse is an interface that encapsulates a JSON API response.
type JsonResponse interface {
	toJson() interface{}
}

// MarketIdAndName represents a market with an ID and a name.
type MarketIdAndName struct {
	marketId   client_types.MarketId
	marketName string
}

// ExchangeIdAndName represents an exchange with an ID and a name.
type ExchangeIdAndName struct {
	exchangeId   client_types.ExchangeFeedId
	exchangeName string
}

// SetupExchangeResponses validates and sets up responses returned by exchange APIs using `gock`.
func SetupExchangeResponses(
	t *testing.T,
	responses MarketToExchangeResponse,
	genesisState types.GenesisState,
) {
	validateExchangeResponses(t, responses, genesisState)

	// Setup `gock` responses.
	for market, exchangeResponse := range responses.marketToExchangesResponse {
		for exchange, response := range exchangeResponse.exchangeToResponse {
			setupGockResponse(t, market, exchange, response)
		}
	}
}

// validateExchangeResponses validates that the genesis state and the exchange API responses match.
func validateExchangeResponses(
	t *testing.T,
	responses MarketToExchangeResponse,
	genesisState types.GenesisState,
) {
	// Gather markets and exchanges from genesis state.
	genesisExchanges := make([]ExchangeIdAndName, len(genesisState.ExchangeFeeds))
	for i, exchange := range genesisState.ExchangeFeeds {
		genesisExchanges[i] = ExchangeIdAndName{
			exchangeId:   exchange.Id,
			exchangeName: exchange.Name,
		}
	}
	genesisMarkets := make([]MarketIdAndName, len(genesisState.Markets))
	for i, market := range genesisState.Markets {
		genesisMarkets[i] = MarketIdAndName{
			marketId:   market.Id,
			marketName: market.Pair,
		}
	}

	// Validate that `responses` have the same set of markets and exchanges as genesis state.
	respMarkets := make([]MarketIdAndName, 0)
	for market, exchangeResponse := range responses.marketToExchangesResponse {
		respMarkets = append(respMarkets, market)

		respExchanges := make([]ExchangeIdAndName, 0)
		for exchange := range exchangeResponse.exchangeToResponse {
			respExchanges = append(respExchanges, exchange)
		}
		sort.SliceStable(respExchanges, func(i, j int) bool {
			return respExchanges[i].exchangeId < respExchanges[j].exchangeId
		})
		require.ElementsMatch(t, genesisExchanges, respExchanges)
	}
	sort.SliceStable(respMarkets, func(i, j int) bool {
		return respMarkets[i].marketId < respMarkets[j].marketId
	})
	require.ElementsMatch(t, genesisMarkets, respMarkets)
}

// setupGockResponse sets up the mock API responses returned by exchanges using `gock`.
func setupGockResponse(
	t *testing.T,
	market MarketIdAndName,
	exchange ExchangeIdAndName,
	response Response,
) {
	marketId := market.marketId
	exchangeName := exchange.exchangeName

	var gockResponse *gock.Response
	switch exchangeName {
	case exchange_common.EXCHANGE_NAME_BINANCE:
		gockResponse = NewGockBinanceResponse(
			marketId,
			response.ResponseCode,
			response.Response.(BinanceResponse),
		)
	case exchange_common.EXCHANGE_NAME_BINANCEUS:
		gockResponse = NewGockBinanceUSResponse(
			marketId,
			response.ResponseCode,
			response.Response.(BinanceResponse),
		)
	case exchange_common.EXCHANGE_NAME_BITFINEX:
		gockResponse = NewGockBitfinexResponse(
			marketId,
			response.ResponseCode,
			response.Response.(BitfinexResponse),
		)
	default:
		panic(fmt.Errorf("unsupported exchange: %s", exchangeName))
	}

	require.NotNil(t, gockResponse)
}
