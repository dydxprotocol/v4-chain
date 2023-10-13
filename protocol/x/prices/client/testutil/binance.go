package testutil

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	"strings"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/binance"
	"github.com/h2non/gock"
)

var (
	binanceUsMarketConfig = map[types.MarketId]types.MarketConfig{
		exchange_config.MARKET_BTC_USD: {
			Ticker:         "BTCUSDT",
			AdjustByMarket: newMarketIdWithValue(exchange_config.MARKET_USDT_USD),
		},
		exchange_config.MARKET_ETH_USD: {
			Ticker:         "ETHUSDT",
			AdjustByMarket: newMarketIdWithValue(exchange_config.MARKET_USDT_USD),
		},
		exchange_config.MARKET_USDT_USD: {
			Ticker: "USDTUSD",
		},
	}
)

// newMarketIdWithValue returns a pointer to a new MarketId with the given value.
func newMarketIdWithValue(id types.MarketId) *types.MarketId {
	val := new(types.MarketId)
	*val = id
	return val
}

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

	// Construct `symbols` parameter in Binance API request.
	sortedTickers := GetTickersSortedByMarketId(
		binanceUsMarketConfig,
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

	// Construct `symbols` parameter in BinanceUS API request.
	sortedTickers := GetTickersSortedByMarketId(
		binanceUsMarketConfig,
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
