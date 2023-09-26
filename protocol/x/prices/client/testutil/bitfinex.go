package testutil

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/bitfinex"
	"github.com/h2non/gock"
)

// BitfinexTicker represents ticker in Bitfinex response for testing purposes.
type BitfinexTicker struct {
	Symbol    string  // index 0
	BidPrice  float64 // index 1
	AskPrice  float64 // index 3
	LastPrice float64 // index 7
}

// NewBitfinexTicker returns a new BitfinexTicker.
func NewBitfinexTicker(symbol string, askPrice, bidPrice, lastPrice float64) BitfinexTicker {
	return BitfinexTicker{
		Symbol:    symbol,
		AskPrice:  askPrice,
		BidPrice:  bidPrice,
		LastPrice: lastPrice,
	}
}

// toJson returns a JSON representation of a valid ticker in Bitfinex response.
func (t BitfinexTicker) toJson() interface{} {
	return []interface{}{
		t.Symbol,    // idx = 0
		t.BidPrice,  // idx = 1
		0.0,         // idx = 2
		t.AskPrice,  // idx = 3
		0.0,         // idx = 4
		0.0,         // idx = 5
		0.0,         // idx = 6
		t.LastPrice, // idx = 7
		0.0,         // idx = 8
		0.0,         // idx = 9
		0.0,         // idx = 10
	}
}

// NewGockBitfinexResponse creates and registers a new HTTP mock using `gock` for Bitfinex.
func NewGockBitfinexResponse(
	responseCode int,
	tickers []JsonResponse,
) *gock.Response {
	// Construct Bitfinex request URL.
	sortedSymbols := GetTickersSortedByMarketId(BitfinexExchangeConfig)
	url := handler.CreateRequestUrl(bitfinex.BitfinexDetails.Url, sortedSymbols)

	// Construct Bitfinex response as a list of tickers.
	jsonResponse := []interface{}{}
	for _, ticker := range tickers {
		jsonResponse = append(jsonResponse, ticker.(BitfinexTicker).toJson())
	}

	return gock.New(url).Persist().Get("/").Reply(responseCode).JSON(jsonResponse)
}
