package testutil

import (
	"fmt"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/bitfinex"
	clienttypes "github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/h2non/gock"
)

// BitfinexResponse represents response returned by Bitfinex for testing purposes.
type BitfinexResponse struct {
	BidPrice  float64 // index 0
	AskPrice  float64 // index 2
	LastPrice float64 // index 6
}

// NewBitfinexResponse returns a new BitfinexResponse.
func NewBitfinexResponse(askPrice, bidPrice, lastPrice float64) BitfinexResponse {
	return BitfinexResponse{
		AskPrice:  askPrice,
		BidPrice:  bidPrice,
		LastPrice: lastPrice,
	}
}

// toJson returns a JSON representation of a valid Bitfinex response.
func (r BitfinexResponse) toJson() interface{} {
	return []float64{
		r.BidPrice,  // idx = 0
		0.0,         // idx = 1
		r.AskPrice,  // idx = 2
		0.0,         // idx = 3
		0.0,         // idx = 4
		0.0,         // idx = 5
		r.LastPrice, // idx = 6
		0.0,         // idx = 7
		0.0,         // idx = 8
		0.0,         // idx = 9
	}
}

// NewGockBitfinexResponse creates and registers a new HTTP mock using `gock` for Bitfinex.
func NewGockBitfinexResponse(
	marketId clienttypes.MarketId,
	responseCode int,
	response BitfinexResponse,
) *gock.Response {
	symbol, exists := bitfinex.BitfinexDetails.MarketSymbols[marketId]
	if !exists {
		panic(fmt.Sprintf("Bitfinex: market (%d) does not exist!", marketId))
	}

	url := handler.CreateRequestUrl(bitfinex.BitfinexDetails.Url, []string{symbol})
	return gock.New(url).Persist().Get("/").Reply(responseCode).JSON(response.toJson())
}
