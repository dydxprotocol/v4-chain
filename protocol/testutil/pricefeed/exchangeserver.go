package pricefeed

import (
	"context"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/testexchange"
	pricefeed "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/client/testutil"
	"io"
	"log"
	"net/http"
	"strings"
)

// This file implements an HTTP server that is used to fake price data from exchanges.
// It's accessible by mapping the testexchange.TestExchangeHost to the host running the test.

const (
	defaultPrice        float64 = 1000
	bitfinexPriceOffset float64 = 100
)

var (
	testExchangeSymbolToMarketId     = map[string]pricefeed.MarketId{}
	bitfinexExchangeSymbolToMarketId = map[string]pricefeed.MarketId{}
)

func init() {
	testExchangeConfig := exchange_config.TestnetExchangeMarketConfig[exchange_common.EXCHANGE_ID_TEST_EXCHANGE]
	for marketId, config := range testExchangeConfig.MarketToMarketConfig {
		testExchangeSymbolToMarketId[config.Ticker] = marketId
	}

	for marketId, config := range testutil.BitfinexExchangeConfig {
		bitfinexExchangeSymbolToMarketId[config.Ticker] = marketId
	}
}

type ExchangeServer struct {
	fakeServer *http.Server
	priceMap   map[pricefeed.MarketId]float64
}

// NewExchangeServer creates a new ExchangeServer that can be used to fake price data from exchanges.
// This server responds queries for the test exchange built into the pricefeed daemon and for bitfinex.
// It is appropriate to run locally or on containers.
func NewExchangeServer() *ExchangeServer {
	ret := &ExchangeServer{
		priceMap: map[pricefeed.MarketId]float64{},
	}
	ret.startFakeServer()
	return ret
}

func (p *ExchangeServer) GetPrice(marketId pricefeed.MarketId) float64 {
	val, ok := p.priceMap[marketId]
	if ok {
		return val
	}
	return defaultPrice
}

func (p *ExchangeServer) SetPrice(marketId pricefeed.MarketId, price float64) {
	p.priceMap[marketId] = price
	print(p.priceMap)
}

// addTestExchangeHandler updates the mux to respond to requests for the test exchange with the
// standard coinbase ticker response. This is the default configuration for the test exchange server.
func addTestExchangeHandler(mux *http.ServeMux, es *ExchangeServer) {
	mux.HandleFunc("/ticker", func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
		currentPrice := es.GetPrice(testExchangeSymbolToMarketId[symbol])
		_, _ = io.WriteString(
			w,
			fmt.Sprintf(
				`{"ask":"%g","bid":"%g","price":"%g"}`,
				currentPrice,
				currentPrice,
				currentPrice,
			),
		)
	})
}

// addTestBitfinexExchangeHandler updates the mux to respond to requests for the test exchange with bitfinex
// symbols to return a bitfinex response, using default prices plus a constant for all symbols.
func addTestBitfinexExchangeHandler(mux *http.ServeMux, es *ExchangeServer) {
	mux.HandleFunc("/bitfinex-ticker", func(w http.ResponseWriter, r *http.Request) {
		symbols := strings.Split(r.URL.Query().Get("symbols"), ",")
		tickers := make([]string, 0, len(symbols))
		for _, symbol := range symbols {
			currentPrice := es.GetPrice(bitfinexExchangeSymbolToMarketId[symbol]) + bitfinexPriceOffset
			tickers = append(tickers, fmt.Sprintf(
				`["%s",%g,"",%g,"","","",%g,"","",""]`,
				symbol,
				currentPrice,
				currentPrice,
				currentPrice,
			))
		}

		_, _ = io.WriteString(
			w,
			fmt.Sprintf(`[%s]`, strings.Join(tickers, ",")),
		)
	})
}

// startFakeServer starts up the server with endpoint handling for both the test exchange and bitfinex.
func (p *ExchangeServer) startFakeServer() {
	mux := http.NewServeMux()
	addTestExchangeHandler(mux, p)
	addTestBitfinexExchangeHandler(mux, p)

	p.fakeServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", testexchange.TestExchangePort),
		Handler: mux,
	}

	go func() {
		if err := p.fakeServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Exchange ListenAndServe() failed: %v", err)
		}
	}()
}

// CleanUp shuts down the server and validates that it shut down correctly.
func (p *ExchangeServer) CleanUp() error {
	if err := p.fakeServer.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}
