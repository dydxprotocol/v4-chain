package testutil

import (
	"fmt"
	"strings"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/price_function/binance"
	clienttypes "github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/h2non/gock"
)

// BinanceResponse represents response returned by Binance for testing purposes.
type BinanceResponse struct {
	AskPrice  string
	BidPrice  string
	LastPrice string
}

// NewBinanceResponse returns a new BinanceResponse.
func NewBinanceResponse(askPrice, bidPrice, lastPrice string) BinanceResponse {
	return BinanceResponse{
		AskPrice:  askPrice,
		BidPrice:  bidPrice,
		LastPrice: lastPrice,
	}
}

// toJson returns a JSON representation of a valid Binance response.
func (r BinanceResponse) toJson() interface{} {
	return map[string]string{
		"askPrice":  r.AskPrice,
		"bidPrice":  r.BidPrice,
		"lastPrice": r.LastPrice,
	}
}

// NewGockBinanceResponse creates and registers a new HTTP mock using `gock` for Binance.
func NewGockBinanceResponse(
	marketId clienttypes.MarketId,
	responseCode int,
	response BinanceResponse,
) *gock.Response {
	rootUrl := binance.BinanceDetails.Url
	rootUrl = rootUrl[:strings.Index(rootUrl, "?")]
	symbol, exists := binance.BinanceDetails.MarketSymbols[marketId]
	if !exists {
		panic(fmt.Sprintf("Binance: market (%d) does not exist!", marketId))
	}

	return gock.New(rootUrl).
		Persist().
		MatchParam("symbol", symbol).
		Reply(responseCode).
		JSON(response.toJson())
}

// NewGockBinanceResponse creates and registers a new HTTP mock using `gock` for BinanceUS.
func NewGockBinanceUSResponse(
	marketId clienttypes.MarketId,
	responseCode int,
	response BinanceResponse,
) *gock.Response {
	rootUrl := binance.BinanceUSDetails.Url
	rootUrl = rootUrl[:strings.Index(rootUrl, "?")]
	symbol, exists := binance.BinanceUSDetails.MarketSymbols[marketId]
	if !exists {
		panic(fmt.Sprintf("BinanceUS:arket (%d) does not exist!", marketId))
	}

	return gock.New(rootUrl).
		Persist().
		MatchParam("symbol", symbol).
		Reply(responseCode).
		JSON(response.toJson())
}
