package client

import (
	"context"
	"cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_encoder"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_fetcher"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"net/http"
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	daemonlib "github.com/dydxprotocol/v4-chain/protocol/daemons/shared"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	gometrics "github.com/hashicorp/go-metrics"
)

var (
	HttpClient = http.Client{
		Transport: &http.Transport{MaxConnsPerHost: constants.MaxConnectionsPerExchange},
	}
)

// SubTaskRunnerImpl is the struct that implements the `SubTaskRunner` interface.
type SubTaskRunnerImpl struct{}

// Ensure the `SubTaskRunnerImpl` struct is implemented at compile time.
var _ SubTaskRunner = (*SubTaskRunnerImpl)(nil)

// SubTaskRunner is the interface for running pricefeed client task functions.
type SubTaskRunner interface {
	StartPriceUpdater(
		c *Client,
		ctx context.Context,
		ticker *time.Ticker,
		stop <-chan bool,
		exchangeToMarketPrices types.ExchangeToMarketPrices,
		priceFeedServiceClient api.PriceFeedServiceClient,
		logger log.Logger,
	)
	StartPriceEncoder(
		exchangeId types.ExchangeId,
		configs types.PricefeedMutableMarketConfigs,
		exchangeToMarketPrices types.ExchangeToMarketPrices,
		logger log.Logger,
		bCh <-chan *price_fetcher.PriceFetcherSubtaskResponse,
	)
	StartPriceFetcher(
		ticker *time.Ticker,
		stop <-chan bool,
		configs types.PricefeedMutableMarketConfigs,
		exchangeQueryConfig types.ExchangeQueryConfig,
		exchangeDetails types.ExchangeQueryDetails,
		queryHandler handler.ExchangeQueryHandler,
		logger log.Logger,
		bCh chan<- *price_fetcher.PriceFetcherSubtaskResponse,
	)
	StartMarketParamUpdater(
		ctx context.Context,
		ticker *time.Ticker,
		stop <-chan bool,
		configs types.PricefeedMutableMarketConfigs,
		pricesQueryClient pricetypes.QueryClient,
		logger log.Logger,
	)
}

// StartPriceUpdater periodically runs a task loop to send price updates to the pricefeed server
// via:
// 1) Get `MarketPriceTimestamps` for all exchanges in an `ExchangeToMarketPrices` struct.
// 2) Transform `MarketPriceTimestamps` and exchange ids into an `UpdateMarketPricesRequest` struct.
// StartPriceUpdater runs in the daemon's main goroutine and does not need access to the daemon's wait group
// to signal task completion.
func (s *SubTaskRunnerImpl) StartPriceUpdater(
	c *Client,
	ctx context.Context,
	ticker *time.Ticker,
	stop <-chan bool,
	exchangeToMarketPrices types.ExchangeToMarketPrices,
	priceFeedServiceClient api.PriceFeedServiceClient,
	logger log.Logger,
) {
	for {
		select {
		case <-ticker.C:
			err := RunPriceUpdaterTaskLoop(ctx, exchangeToMarketPrices, priceFeedServiceClient, logger)

			if err == nil {
				// Record update success for the daemon health check.
				c.ReportSuccess()
			} else {
				logger.Error("Failed to run price updater task loop for price daemon", constants.ErrorLogKey, err)
				// Record update failure for the daemon health check.
				c.ReportFailure(errors.Wrap(err, "failed to run price updater task loop for price daemon"))
			}

		case <-stop:
			return
		}
	}
}

// StartPriceEncoder continuously reads from a buffered channel, reading encoded API responses for exchange
// requests and inserting them into an `ExchangeToMarketPrices` cache, performing currency conversions based
// on the index price of other markets as necessary.
// StartPriceEncoder reads price fetcher responses from a shared channel, and does not need a ticker or stop
// signal from the daemon to exit. It marks itself as done in the daemon's wait group when the price fetcher
// closes the shared channel.
func (s *SubTaskRunnerImpl) StartPriceEncoder(
	exchangeId types.ExchangeId,
	configs types.PricefeedMutableMarketConfigs,
	exchangeToMarketPrices types.ExchangeToMarketPrices,
	logger log.Logger,
	bCh <-chan *price_fetcher.PriceFetcherSubtaskResponse,
) {
	exchangeMarketConfig, err := configs.GetExchangeMarketConfigCopy(exchangeId)
	if err != nil {
		panic(err)
	}

	marketConfigs, err := configs.GetMarketConfigCopies(exchangeMarketConfig.GetMarketIds())
	if err != nil {
		panic(err)
	}

	priceEncoder, err := price_encoder.NewPriceEncoder(
		exchangeMarketConfig,
		marketConfigs,
		exchangeToMarketPrices,
		logger,
		bCh,
	)

	if err != nil {
		panic(err)
	}

	configs.AddPriceEncoder(priceEncoder)

	// Listen for prices from the buffered channel and update the exchangeToMarketPrices cache.
	// Also log any errors that occur.
	for response := range bCh {
		priceEncoder.ProcessPriceFetcherResponse(response)
	}
}

// StartPriceFetcher periodically starts goroutines to "fetch" market prices from a specific exchange. Each
// goroutine does the following:
// 1) query a single market price from a specific exchange
// 2) transform response to `MarketPriceTimestamp`
// 3) send transformed response to a buffered channel that's shared across multiple goroutines
// NOTE: the subtask response shared channel has a buffer size and goroutines will block if the buffer is full.
// NOTE: the price fetcher kicks off 1 to n go routines every time the subtask loop runs, but the subtask
// loop blocks until all go routines are done. This means that these go routines are not tracked by the wait group.
func (s *SubTaskRunnerImpl) StartPriceFetcher(
	ticker *time.Ticker,
	stop <-chan bool,
	configs types.PricefeedMutableMarketConfigs,
	exchangeQueryConfig types.ExchangeQueryConfig,
	exchangeDetails types.ExchangeQueryDetails,
	queryHandler handler.ExchangeQueryHandler,
	logger log.Logger,
	bCh chan<- *price_fetcher.PriceFetcherSubtaskResponse,
) {
	exchangeMarketConfig, err := configs.GetExchangeMarketConfigCopy(exchangeQueryConfig.ExchangeId)
	if err != nil {
		panic(err)
	}

	marketConfigs, err := configs.GetMarketConfigCopies(exchangeMarketConfig.GetMarketIds())
	if err != nil {
		panic(err)
	}

	// Create PriceFetcher to begin querying with.
	priceFetcher, err := price_fetcher.NewPriceFetcher(
		exchangeQueryConfig,
		exchangeDetails,
		exchangeMarketConfig,
		marketConfigs,
		queryHandler,
		logger,
		bCh,
	)
	if err != nil {
		panic(err)
	}

	// The PricefeedMutableMarketConfigs object that stores the configs for each exchange
	// is not initialized with the price fetcher, because both objects have references to
	// each other defined during normal daemon operation. Instead, the price fetcher is
	// initialized with the configs object after the price fetcher is created, and then adds
	// itself to the config's list of exchange config updaters here.
	configs.AddPriceFetcher(priceFetcher)

	requestHandler := daemontypes.NewRequestHandlerImpl(
		&HttpClient,
	)
	// Begin loop to periodically start goroutines to query market prices.
	for {
		select {
		case <-ticker.C:
			// Start goroutines to query exchange markets. The goroutines started by the price
			// fetcher are not tracked by the global wait group, because RunTaskLoop will
			// block until all goroutines are done.
			priceFetcher.RunTaskLoop(requestHandler)

		case <-stop:
			// Signal to the encoder that the price fetcher is done.
			close(bCh)
			return
		}
	}
}

// StartMarketParamUpdater periodically starts a goroutine to update the market parameters that control which
// markets the daemon queries and how they are queried and computed from each exchange.
func (s *SubTaskRunnerImpl) StartMarketParamUpdater(
	ctx context.Context,
	ticker *time.Ticker,
	stop <-chan bool,
	configs types.PricefeedMutableMarketConfigs,
	pricesQueryClient pricetypes.QueryClient,
	logger log.Logger,
) {
	// Delay reporting certain errors for a grace period to allow the daemon to start up. There is a bit of a race
	// condition here with reading/writing the variable, but it's not a big deal if there is some jitter in the
	// timing of the grace period ending.
	isPastGracePeriod := false
	go func() {
		time.Sleep(constants.PriceDaemonStartupErrorGracePeriod)
		isPastGracePeriod = true
	}()

	// Periodically update market parameters.
	for {
		select {
		case <-ticker.C:
			RunMarketParamUpdaterTaskLoop(ctx, configs, pricesQueryClient, logger, isPastGracePeriod)

		case <-stop:
			return
		}
	}
}

// -------------------- Task Loops -------------------- //

// RunPriceUpdaterTaskLoop copies the map of current `exchangeId -> MarketPriceTimestamp`,
// transforms the map values into a market price update request and sends the request to the socket
// where the pricefeed server is listening.
func RunPriceUpdaterTaskLoop(
	ctx context.Context,
	exchangeToMarketPrices types.ExchangeToMarketPrices,
	priceFeedServiceClient api.PriceFeedServiceClient,
	logger log.Logger,
) error {
	logger = logger.With(constants.SubmoduleLogKey, constants.PriceUpdaterSubmoduleName)
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

	// On startup the length of request will likely be 0. Even so, we return an error here because this
	// is unexpected behavior once the daemon reaches a steady state. The daemon health check process should
	// be robust enough to ignore temporarily unhealthy daemons.
	// Sending a request of length 0, however, causes a panic.
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
		logger.Info("Price update had length of 0")
		telemetry.IncrCounter(
			1,
			metrics.PricefeedDaemon,
			metrics.PriceUpdaterZeroPrices,
			metrics.Count,
		)
		return types.ErrEmptyMarketPriceUpdate
	}

	return nil
}

// RunMarketParamUpdaterTaskLoop queries all market params from the query client, and then updates the
// shared, in-memory `PricefeedMutableMarketConfigs` object with the latest market params.
func RunMarketParamUpdaterTaskLoop(
	ctx context.Context,
	configs types.PricefeedMutableMarketConfigs,
	pricesQueryClient pricetypes.QueryClient,
	logger log.Logger,
	isPastGracePeriod bool,
) {
	// Measure latency to fetch and parse the market params, and propagate all updates.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedDaemon,
		time.Now(),
		metrics.MarketUpdaterUpdateMarkets,
		metrics.Latency,
	)

	logger = logger.With(constants.SubmoduleLogKey, constants.MarketParamUpdaterSubmoduleName)

	// Query all market params from the query client.
	marketParams, err := daemonlib.AllPaginatedMarketParams(ctx, pricesQueryClient)
	if err != nil {
		var logMethod = logger.Info
		if isPastGracePeriod {
			// When the daemon starts, there is normally a delay between when the prices daemon starts and the prices
			// query service becomes available. This is not a true error condition, so we log it as info instead of
			// error in order to avoid spurious error logs and alerts.
			logMethod = logger.Error
		}
		logMethod("Failed to get all market params",
			"error",
			err,
		)
		// Measure all failures to retrieve market params from the query client.
		telemetry.IncrCounter(
			1,
			metrics.PricefeedDaemon,
			metrics.MarketUpdaterGetAllMarketParams,
			metrics.Error,
		)
		return
	}

	// Update shared, in-memory config with the latest market params. Report update success/failure via logging/metrics.
	marketParamErrors, err := configs.UpdateMarkets(marketParams)

	for _, marketParam := range marketParams {
		// Update the market id -> pair for telemetry.
		pricefeedmetrics.SetMarketPairForTelemetry(marketParam.Id, marketParam.Pair)

		outcome := metrics.Success

		// Mark this update as an error either if this market failed to update, or if all markets failed.
		if _, ok := marketParamErrors[marketParam.Id]; ok || err != nil {
			outcome = metrics.Error
		}

		telemetry.IncrCounterWithLabels(
			[]string{metrics.PricefeedDaemon, metrics.MarketUpdaterApplyMarketUpdates, outcome},
			1,
			[]gometrics.Label{
				pricefeedmetrics.GetLabelForMarketId(marketParam.Id),
			},
		)
	}
	if err != nil {
		logger.Error(
			"Failed to apply all market updates",
			"error",
			err,
			"marketParamErrors",
			marketParamErrors,
		)
	} else if len(marketParamErrors) > 0 {
		logger.Error(
			"Failed to apply some market updates",
			"marketParamErrors",
			marketParamErrors,
		)
	}
}

// -------------------- Task Loop Helpers -------------------- //

// transformPriceUpdates transforms a map (key: exchangeId, value: list of market prices) into a
// market price update request.
func transformPriceUpdates(
	updates map[types.ExchangeId][]types.MarketPriceTimestamp,
) *api.UpdateMarketPricesRequest {
	// Measure latency to transform prices being sent over gRPC.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedDaemon,
		time.Now(),
		metrics.PriceUpdaterTransformPrices,
		metrics.Latency,
	)

	marketPriceUpdateMap := make(map[types.MarketId]*api.MarketPriceUpdate)

	// Invert to marketId -> `api.MarketPriceUpdate`.
	for exchangeId, marketPriceTimestamps := range updates {
		for _, marketPriceTimestamp := range marketPriceTimestamps {
			telemetry.IncrCounterWithLabels(
				[]string{
					metrics.PricefeedDaemon,
					metrics.PriceUpdateCount,
					metrics.Count,
				},
				1,
				[]gometrics.Label{
					pricefeedmetrics.GetLabelForExchangeId(exchangeId),
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
				ExchangeId:     exchangeId,
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
