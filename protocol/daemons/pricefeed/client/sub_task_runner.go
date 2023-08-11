package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"syscall"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"gopkg.in/typ.v4/lists"
)

const (
	// https://stackoverflow.com/questions/37774624/go-http-get-concurrency-and-connection-reset-by-peer.
	// This is a good number to start with based on the above link. Adjustments can/will be made accordingly.
	MaxConnectionsPerExchange = 50
)

var (
	HttpClient = http.Client{
		Transport: &http.Transport{MaxConnsPerHost: MaxConnectionsPerExchange},
	}
)

// SubTaskRunnerImpl is the struct that implements the `SubTaskRunner` interface.
type SubTaskRunnerImpl struct{}

// Ensure the `SubTaskRunnerImpl` struct is implemented at compile time.
var _ SubTaskRunner = (*SubTaskRunnerImpl)(nil)

// SubTaskRunner is the interface for running pricefeed client task functions.
type SubTaskRunner interface {
	StartPriceUpdater(
		ctx context.Context,
		exchangeToMarketPrices *types.ExchangeToMarketPrices,
		priceFeedServiceClient api.PriceFeedServiceClient,
		loopDelayMs uint32,
		logger log.Logger,
	)
	StartPriceEncoder(
		exchangeFeedId types.ExchangeFeedId,
		exchangeToMarketPrices *types.ExchangeToMarketPrices,
		logger log.Logger,
		bCh chan *PriceFetcherSubtaskResponse,
	)
	StartPriceFetcher(
		exchangeConfig types.ExchangeConfig,
		queryHandler handler.ExchangeQueryHandler,
		logger log.Logger,
		bCh chan *PriceFetcherSubtaskResponse,
	)
}

// StartPriceUpdater periodically runs a task loop to send price updates to the pricefeed server
// via:
// 1) Get `MarketPriceTimestamps` for all exchanges in an `ExchangeToMarketPrices` struct.
// 2) Transform `MarketPriceTimestamps` and exchange ids into an `UpdateMarketPricesRequest` struct.
func (s *SubTaskRunnerImpl) StartPriceUpdater(
	ctx context.Context,
	exchangeToMarketPrices *types.ExchangeToMarketPrices,
	priceFeedServiceClient api.PriceFeedServiceClient,
	loopDelayMs uint32,
	logger log.Logger,
) {
	// Start a `ticker` to run `RunPriceUpdaterTaskLoop` immediately and then every `loopDelayMs`
	// milliseconds.
	ticker := time.NewTicker(time.Duration(loopDelayMs) * time.Millisecond)
	for ; true; <-ticker.C {
		err := RunPriceUpdaterTaskLoop(ctx, exchangeToMarketPrices, priceFeedServiceClient, logger)
		if err != nil {
			panic(err)
		}
	}
}

// StartPriceEncoder continuously reads from a buffered channel, reading encoded API responses for exchange
// requests and inserting them into an `ExchangeToMarketPrices` cache.
func (s *SubTaskRunnerImpl) StartPriceEncoder(
	exchangeFeedId types.ExchangeFeedId,
	exchangeToMarketPrices *types.ExchangeToMarketPrices,
	logger log.Logger,
	bCh chan *PriceFetcherSubtaskResponse,
) {
	for response := range bCh {
		if response.Err == nil {
			exchangeToMarketPrices.UpdatePrice(exchangeFeedId, response.Price)
		} else {
			if errors.Is(response.Err, context.DeadlineExceeded) {
				// Log info if there are timeout errors in the ingested buffered channel prices.
				// This is only an info so that there aren't noisy errors when undesirable but
				// expected behavior occurs.
				logger.Info(
					"Failed to update exchange price in price daemon priceEncoder due to timeout",
					"error",
					response.Err,
					"exchangeFeedId",
					exchangeFeedId,
				)

				// Measure timeout failures.
				telemetry.IncrCounterWithLabels(
					[]string{
						metrics.PricefeedDaemon,
						metrics.PriceEncoderUpdatePrice,
						metrics.HttpGetTimeout,
						metrics.Error,
					},
					1,
					[]gometrics.Label{
						pricefeedmetrics.GetLabelForExchangeFeedId(exchangeFeedId),
					},
				)
			} else if strings.Contains(response.Err.Error(), fmt.Sprintf("%s 5", constants.UnexpectedResponseStatusMessage)) {
				// Log info if there are 5xx errors in the ingested buffered channel prices.
				// This is only an info so that there aren't noisy errors when undesirable but
				// expected behavior occurs.
				logger.Info(
					"Failed to update exchange price in price daemon priceEncoder due to exchange-side error",
					"error",
					response.Err,
					"exchangeFeedId",
					exchangeFeedId,
				)

				// Measure 5xx failures.
				telemetry.IncrCounterWithLabels(
					[]string{
						metrics.PricefeedDaemon,
						metrics.PriceEncoderUpdatePrice,
						metrics.HttpGet5xxx,
						metrics.Error,
					},
					1,
					[]gometrics.Label{
						pricefeedmetrics.GetLabelForExchangeFeedId(exchangeFeedId),
					},
				)
			} else if errors.Is(response.Err, syscall.ECONNRESET) {
				// Log info if there are connections reset by the exchange.
				// This is only an info so that there aren't noisy errors when undesirable but
				// expected behavior occurs.
				logger.Info(
					"Failed to update exchange price in price daemon priceEncoder due to exchange-side hang-up",
					"error",
					response.Err,
					"exchangeFeedId",
					exchangeFeedId,
				)

				// Measure HTTP GET hangups.
				telemetry.IncrCounterWithLabels(
					[]string{
						metrics.PricefeedDaemon,
						metrics.PriceEncoderUpdatePrice,
						metrics.HttpGetHangup,
						metrics.Error,
					},
					1,
					[]gometrics.Label{
						pricefeedmetrics.GetLabelForExchangeFeedId(exchangeFeedId),
					},
				)
			} else {
				// Log error if there are errors in the ingested buffered channel prices.
				logger.Error(
					"Failed to update exchange price in price daemon priceEncoder",
					"error",
					response.Err,
					"exchangeFeedId",
					exchangeFeedId,
				)

				// Measure all failures in querying other than timeout.
				telemetry.IncrCounterWithLabels(
					[]string{
						metrics.PricefeedDaemon,
						metrics.PriceEncoderUpdatePrice,
						metrics.Error,
					},
					1,
					[]gometrics.Label{
						pricefeedmetrics.GetLabelForExchangeFeedId(exchangeFeedId),
					},
				)
			}
		}
	}
}

// StartPriceFetcher periodically starts goroutines to "fetch" market prices from a specific exchange. Each
// goroutine does the following:
// 1) query a single market price from a specific exchange
// 2) transform response to `MarketPriceTimestamp`
// 3) send transformed response to a buffered channel that's shared across multiple goroutines
// NOTE: the shared channel has a buffer size and goroutines will block if the buffer is full.
func (s *SubTaskRunnerImpl) StartPriceFetcher(
	exchangeConfig types.ExchangeConfig,
	queryHandler handler.ExchangeQueryHandler,
	logger log.Logger,
	bCh chan *PriceFetcherSubtaskResponse,
) {
	marketIds := exchangeConfig.Markets

	// Create ring that holds all markets for an exchange.
	marketIdsRing := lists.NewRing[types.MarketId](len(marketIds))
	for _, marketId := range marketIds {
		marketIdsRing.Value = marketId
		marketIdsRing = marketIdsRing.Next()
	}

	// Create PriceFetcher to begin querying with.
	priceFetcher := NewPriceFetcher(
		exchangeConfig,
		queryHandler,
		marketIdsRing,
		logger,
		bCh,
	)

	requestHandler := lib.NewRequestHandlerImpl(
		&HttpClient,
	)
	// Begin loop to periodically start goroutines to query market prices.
	for {
		// Start goroutines to query exchange markets.
		priceFetcher.RunTaskLoop(requestHandler)
		// Wait exchange specific time until next loop begins.
		time.Sleep(time.Duration(exchangeConfig.ExchangeStartupConfig.IntervalMs) * time.Millisecond)
	}
}

// -------------------- Task Loops -------------------- //

// RunPriceUpdaterTaskLoop copies the map of current `exchangeFeedId -> MarketPriceTimestamp`,
// transforms the map values into a market price update request and sends the request to the socket
// where the pricefeed server is listening.
func RunPriceUpdaterTaskLoop(
	ctx context.Context,
	exchangeToMarketPrices *types.ExchangeToMarketPrices,
	priceFeedServiceClient api.PriceFeedServiceClient,
	logger log.Logger,
) error {
	priceUpdates := exchangeToMarketPrices.GetAllPrices()
	request := transformPriceUpdates(priceUpdates)

	// Measure latency to send prices over gRPC.
	// Note: intentionally skipping latency for `GetAllPrices`.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedDaemon,
		time.Now(),
		metrics.PriceUpdaterSendPrices,
		metrics.Latency,
	)

	// On startup the length of request will likely be 0. However, sending a request of length 0
	// is a fatal error.
	// panic: rpc error: code = Unknown desc = Market price update has length of 0.
	if len(request.MarketPriceUpdates) > 0 {
		_, err := priceFeedServiceClient.UpdateMarketPrices(ctx, request)
		if err != nil {
			// Log error if an error will be returned from the task loop and measure failure.
			logger.Error("Failed to run price updater task loop for price daemon", "error", err)
			telemetry.IncrCounter(
				1,
				metrics.PricefeedDaemon,
				metrics.PriceUpdaterTaskLoop,
				metrics.Error,
			)
			return err
		}
	} else {
		// This is expected to happen on startup until prices have been encoded into the in-memory
		// `exchangeToMarketPrices` map. After that point, there should be no price updates of length 0.
		logger.Info(
			"Price update had length of 0",
		)
		telemetry.IncrCounter(
			1,
			metrics.PricefeedDaemon,
			metrics.PriceUpdaterZeroPrices,
			metrics.Count,
		)
	}

	return nil
}

// -------------------- Task Loop Helpers -------------------- //

// transformPriceUpdates transforms a map (key: exchangeFeedId, value: list of market prices) into a
// market price update request.
func transformPriceUpdates(
	updates map[types.ExchangeFeedId][]types.MarketPriceTimestamp,
) *api.UpdateMarketPricesRequest {
	// Measure latency to transform prices being sent over gRPC.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedDaemon,
		time.Now(),
		metrics.PriceUpdaterTransformPrices,
		metrics.Latency,
	)

	marketPriceUpdateMap := make(map[types.ExchangeFeedId]*api.MarketPriceUpdate)

	// Invert to marketId -> `api.MarketPriceUpdate`.
	for exchangeFeedId, marketPriceTimestamps := range updates {
		for _, marketPriceTimestamp := range marketPriceTimestamps {
			telemetry.IncrCounterWithLabels(
				[]string{
					metrics.PricefeedDaemon,
					metrics.PriceUpdateCount,
					metrics.Count,
				},
				1,
				[]gometrics.Label{
					pricefeedmetrics.GetLabelForExchangeFeedId(exchangeFeedId),
					pricefeedmetrics.GetLabelForMarketId(marketPriceTimestamp.MarketId),
				},
			)

			marketPriceUpdate, exists := marketPriceUpdateMap[marketPriceTimestamp.MarketId]
			// Add key with empty `api.MarketPriceUpdate` if entry does not exist.
			if !exists {
				marketPriceUpdate = &api.MarketPriceUpdate{
					MarketId:       marketPriceTimestamp.MarketId,
					ExchangePrices: []*api.ExchangePrice{},
				}
				marketPriceUpdateMap[marketPriceTimestamp.MarketId] = marketPriceUpdate
			}

			// Add `api.ExchangePrice`.
			priceUpdateTime := marketPriceTimestamp.LastUpdatedAt
			exchangePrice := &api.ExchangePrice{
				ExchangeFeedId: exchangeFeedId,
				Price:          marketPriceTimestamp.Price,
				LastUpdateTime: &priceUpdateTime,
			}
			marketPriceUpdate.ExchangePrices = append(marketPriceUpdate.ExchangePrices, exchangePrice)
		}
	}

	// Add all `api.MarketPriceUpdate` to request to be sent by `client.UpdateMarketPrices`.
	request := &api.UpdateMarketPricesRequest{
		MarketPriceUpdates: make([]*api.MarketPriceUpdate, 0, len(marketPriceUpdateMap)),
	}
	for _, update := range marketPriceUpdateMap {
		request.MarketPriceUpdates = append(
			request.MarketPriceUpdates,
			update,
		)
	}
	return request
}
