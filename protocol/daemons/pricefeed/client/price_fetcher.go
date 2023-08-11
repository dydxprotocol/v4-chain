package client

import (
	"context"
	"fmt"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"gopkg.in/typ.v4/lists"
)

// PriceFetcherSubtaskResponse represents a transformed exchange API response that contains price
// info or an error.
type PriceFetcherSubtaskResponse struct {
	Price *types.MarketPriceTimestamp
	Err   error
}

// PriceFetcher fetches prices from an exchange by making a query based on the
// `exchangeConfig` specifications and then encodes the price or any associated error.
type PriceFetcher struct {
	exchangeConfig types.ExchangeConfig
	queryHandler   handler.ExchangeQueryHandler
	marketIdsRing  *lists.Ring[types.MarketId]
	logger         log.Logger
	bCh            chan *PriceFetcherSubtaskResponse
}

// NewPriceFetcher creates a new PriceFetcher struct. It manages querying markets via goroutine
// queries to an exchange and encodes the responses or related erros into the shared buffered
// channel `bCh`.
func NewPriceFetcher(
	exchangeConfig types.ExchangeConfig,
	queryHandler handler.ExchangeQueryHandler,
	marketIdsRing *lists.Ring[types.MarketId],
	logger log.Logger,
	bCh chan *PriceFetcherSubtaskResponse,
) *PriceFetcher {
	return &PriceFetcher{
		exchangeConfig,
		queryHandler,
		marketIdsRing,
		logger,
		bCh,
	}
}

// RunTaskLoop queries the exchange for market prices.
// Each goroutine makes a single exchange query for a specific market.
func (pf *PriceFetcher) RunTaskLoop(requestHandler lib.RequestHandler) {
	if pf.exchangeConfig.IsMultiMarket {
		go pf.runSubTask(requestHandler, pf.exchangeConfig.Markets)
	} else {
		maxQueries := lib.Min(len(pf.exchangeConfig.Markets), int(pf.exchangeConfig.ExchangeStartupConfig.MaxQueries))
		for i := 0; i < maxQueries; i++ {
			go pf.runSubTask(requestHandler, []types.MarketId{pf.marketIdsRing.Value})
			pf.marketIdsRing = pf.marketIdsRing.Next()
		}
	}
}

// runSubTask make a single query to an exchange for market prices.
// The price or the error generated during querying is written to the `PriceFetcher`
// buffered channel.
func (pf *PriceFetcher) runSubTask(requestHandler lib.RequestHandler, marketIds []types.MarketId) {
	exchangeFeedId := pf.exchangeConfig.ExchangeStartupConfig.ExchangeFeedId

	// Measure total latency for subtask to run for one API call and creating a context with timeout.
	defer metrics.ModuleMeasureSinceWithLabels(
		metrics.PricefeedDaemon,
		[]string{
			metrics.PricefeedDaemon,
			metrics.PriceFetcherSubtaskLoopAndSetCtxTimeout,
			metrics.Latency,
		},
		time.Now(),
		[]gometrics.Label{pricefeedmetrics.GetLabelForExchangeFeedId(exchangeFeedId)},
	)

	ctxWithTimeout, cancelFunc := context.WithTimeout(
		context.Background(),
		time.Duration(pf.exchangeConfig.ExchangeStartupConfig.TimeoutMs)*time.Millisecond,
	)

	defer cancelFunc()

	// Measure total latency for subtask to run for one API call.
	defer metrics.ModuleMeasureSinceWithLabels(
		metrics.PricefeedDaemon,
		[]string{
			metrics.PricefeedDaemon,
			metrics.PriceFetcherSubtaskLoop,
			metrics.Latency,
		},
		time.Now(),
		[]gometrics.Label{pricefeedmetrics.GetLabelForExchangeFeedId(exchangeFeedId)},
	)

	exchangeDetails := constants.StaticExchangeDetails[exchangeFeedId]
	prices, unavailableMarkets, err := pf.queryHandler.Query(
		ctxWithTimeout,
		&exchangeDetails,
		marketIds,
		requestHandler,
		constants.StaticMarketPriceExponent,
	)

	if err != nil {
		pf.writeToBufferedChannel(exchangeFeedId, nil, err)
		return
	}

	for _, price := range prices {
		// No price should validly be zero. A price of zero points to an error in the API queried.
		if price.Price == uint64(0) {
			pf.writeToBufferedChannel(
				exchangeFeedId,
				nil,
				fmt.Errorf(
					"Invalid price of 0 for exchange: %v and market: %v",
					exchangeFeedId,
					price.MarketId,
				),
			)

			continue
		}

		// Log each new price (per-market per-exchange).
		pf.logger.Debug(
			fmt.Sprintf(
				"Adding new price for market. Price: %d. MarketId: %d. LastUpdatedAt: %v. ExchangeFeedId: %d",
				price.Price,
				price.MarketId,
				price.LastUpdatedAt,
				exchangeFeedId,
			),
		)

		pf.writeToBufferedChannel(exchangeFeedId, price, err)
	}
	for market, error := range unavailableMarkets {
		pf.writeToBufferedChannel(
			exchangeFeedId,
			nil,
			fmt.Errorf("Market %d unavailable on exchange %d (%w)", market, exchangeFeedId, error),
		)
	}
}

func (pf *PriceFetcher) writeToBufferedChannel(
	exchangeFeedId types.ExchangeFeedId,
	price *types.MarketPriceTimestamp,
	err error,
) {
	// Sanity check that the channel is not full already.
	if len(pf.bCh) == FixedBufferSize {
		// Log if writing to buffered channel failed.
		pf.logger.Error(
			fmt.Sprintf(
				"Pricefeed daemon's shared buffer is full. Exchange %d.",
				exchangeFeedId,
			),
		)
	}

	pf.bCh <- &PriceFetcherSubtaskResponse{
		Err:   err,
		Price: price,
	}
}
