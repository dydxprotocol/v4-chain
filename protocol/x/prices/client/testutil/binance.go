package testutil

import (
	"fmt"
	"strings"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/binance"
	"github.com/h2non/gock"
)

// BinanceTicker represents ticker returned by Binance for testing purposes.
type BinanceTicker struct {
	Symbol    string
	AskPrice  string
	BidPrice  string
	LastPrice string
}

// NewBinanceTicker returns a new BinanceTicker.
func NewBinanceTicker(symbol, askPrice, bidPrice, lastPrice string) BinanceTicker {
	return BinanceTicker{
		Symbol:    symbol,
		AskPrice:  askPrice,
		BidPrice:  bidPrice,
		LastPrice: lastPrice,
	}
}

// toJson returns a JSON representation of a valid ticker in Binance response.
func (t BinanceTicker) toJson() interface{} {
	return map[string]string{
		"symbol":    t.Symbol,
		"askPrice":  t.AskPrice,
		"bidPrice":  t.BidPrice,
		"lastPrice": t.LastPrice,
	}
}

// NewGockBinanceResponse creates and registers a new HTTP mock using `gock` for Binance.
func NewGockBinanceResponse(
	responseCode int,
	tickers []JsonResponse,
) *gock.Response {
	rootUrl := binance.BinanceDetails.Url
	rootUrl = rootUrl[:strings.Index(rootUrl, "[")]

	// Construct `symbols` parameter in Binance API request.
	sortedTickers := GetTickersSortedByMarketId(
		constants.StaticExchangeMarketConfig[exchange_common.EXCHANGE_ID_BINANCE].MarketToMarketConfig,
	)
	symbolsParam := fmt.Sprintf(
		"[%s]",
		strings.Join(sortedTickers, ","),
	)

	// Construct Binance API response as a list of tickers.
	jsonResponse := []interface{}{}
	for _, ticker := range tickers {
		jsonResponse = append(jsonResponse, ticker.(BinanceTicker).toJson())
	}

	return gock.New(rootUrl).
		Persist().
		MatchParam("symbols", symbolsParam).
		Reply(responseCode).
		JSON(jsonResponse)
}

// NewGockBinanceResponse creates and registers a new HTTP mock using `gock` for BinanceUS.
func NewGockBinanceUSResponse(
	responseCode int,
	tickers []JsonResponse,
) *gock.Response {
	rootUrl := binance.BinanceUSDetails.Url
	rootUrl = rootUrl[:strings.Index(rootUrl, "[")]

	// Construct `symbols` parameter in BinanceUS API request.
	sortedTickers := GetTickersSortedByMarketId(
		constants.StaticExchangeMarketConfig[exchange_common.EXCHANGE_ID_BINANCE_US].MarketToMarketConfig,
	)
	symbolsParam := fmt.Sprintf(
		"[%s]",
		strings.Join(sortedTickers, ","),
	)

	// Construct BinanceUS API response as a list of tickers.
	jsonResponse := []interface{}{}
	for _, ticker := range tickers {
		jsonResponse = append(jsonResponse, ticker.(BinanceTicker).toJson())
	}

	return gock.New(rootUrl).
		Persist().
		MatchParam("symbols", symbolsParam).
		Reply(responseCode).
		JSON(jsonResponse)
}
