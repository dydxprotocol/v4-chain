package client

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/log"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_fetcher"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// Client encapsulates the logic for executing and cleanly stopping all subtasks associated with the
// pricefeed client daemon. Access to the client's internal state is synchronized.
// The pricefeed daemon is a job that periodically queries external exchanges and transmits
// price data to the pricefeed service, which is then used by the application to compute index
// prices for proposing and validating oracle price updates on the blockchain.
// Note: price fetchers manage their own subtasks by blocking on their completion on every subtask run.
// When the price fetcher is stopped, it will wait for all of its own subtasks to complete before returning.
type Client struct {
	// include HealthCheckable to track the health of the daemon.
	daemontypes.HealthCheckable

	// daemonStartup tracks whether the daemon has finished startup. The daemon
	// cannot be stopped until all persistent daemon subtasks have been launched within `Start`.
	daemonStartup sync.WaitGroup

	// runningSubtasksWaitGroup tracks the number of running subtasks on the daemon.
	// This is used to block the daemon from stopping until all running processes have completed.
	runningSubtasksWaitGroup sync.WaitGroup

	// tickers tracks the list of tickers that are used to execute subtasks that repeat periodically on the daemon.
	// Access to tickers is implicitly synchronized by the daemonStartup WaitGroup.
	tickers []*time.Ticker

	// stops tracks the list of channels that are used to send a stop signal to subtasks on the daemon.
	// Access to stops is implicitly synchronized by the daemonStartup WaitGroup.
	stops []chan bool

	// Ensure stop only executes one time.
	stopDaemon sync.Once

	// logger is the logger for the daemon.
	logger log.Logger
}

// Ensure Client implements the HealthCheckable interface.
var _ daemontypes.HealthCheckable = (*Client)(nil)

func newClient(logger log.Logger) *Client {
	logger = logger.With(log.ModuleKey, constants.PricefeedDaemonModuleName)
	client := &Client{
		tickers: []*time.Ticker{},
		stops:   []chan bool{},
		HealthCheckable: daemontypes.NewTimeBoundedHealthCheckable(
			constants.PricefeedDaemonModuleName,
			&libtime.TimeProviderImpl{},
			logger,
		),
		logger: logger,
	}

	// Set the client's daemonStartup state to indicate that the daemon has not finished starting up.
	client.daemonStartup.Add(1)
	return client
}

// newTickerWithStop creates a new ticker and a channel for iteratively looping through a subtask with a stop signal
// for any subtask kicked off by the client. The ticker and channel are tracked in order to properly clean up and send
// all needed stop signals when the daemon is stopped.
// Note: this method is not synchronized. It is expected to be called from the client's `StartNewClient` method before
// the daemonStartup waitgroup signals.
func (c *Client) newTickerWithStop(intervalMs int) (*time.Ticker, <-chan bool) {
	ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
	c.tickers = append(c.tickers, ticker)

	stop := make(chan bool)
	c.stops = append(c.stops, stop)

	return ticker, stop
}

// Stop stops the daemon and all running subtasks. This method is synchronized by the daemonStartup WaitGroup.
func (c *Client) Stop() {
	c.stopDaemon.Do(func() {
		c.daemonStartup.Wait()

		// Send a signal to all tickers and stop channels to stop all running subtasks managed by the client.
		for _, stop := range c.stops {
			close(stop)
		}
		for _, ticker := range c.tickers {
			ticker.Stop()
		}

		c.runningSubtasksWaitGroup.Wait()
	})
}

// start begins a job that:
// A) periodically queries prices from external data sources and saves the retrieved prices in an
// in-memory datastore
// B) periodically sends the most recent prices to a gRPC server
// C) periodically queries the prices module for the latest market/exchange configuration and then updates
// the shared, in-memory datastore with the latest configuration.
// The exchangeIdToQueryConfig map dictates which exchanges the pricefeed client queries against.
// For all exchanges included in this map, the pricefeed client expects an exchangeQueryDetails and an
// initialExchangeMarketConfig object to be defined in the parameter maps. To initialize an exchange with
// zero markets, pass in an initialExchangeMarketConfig object with an empty map of market tickers for that
// exchange.
// Implementation:
//  1. Establish connections to gRPC servers.
//  2. Validate daemon configuration.
//  3. Initialize synchronized, in-memory shared daemon configuration.
//  4. Start PriceEncoder and PriceFetcher per exchange. Each price fetcher adds itself to the shared
//     daemon config.
//  5. Start MarketUpdater subtask to periodically update the market configs.
//  6. Start PriceUpdater to begin broadcasting prices.
func (c *Client) start(ctx context.Context,
	daemonFlags flags.DaemonFlags,
	appFlags appflags.Flags,
	grpcClient daemontypes.GrpcClient,
	exchangeIdToQueryConfig map[types.ExchangeId]*types.ExchangeQueryConfig,
	exchangeIdToExchangeDetails map[types.ExchangeId]types.ExchangeQueryDetails,
	subTaskRunner SubTaskRunner,
) (err error) {
	// 1. Establish connections to gRPC servers.
	queryConn, err := grpcClient.NewTcpConnection(ctx, appFlags.GrpcAddress)
	if err != nil {
		c.logger.Error("Failed to establish gRPC connection to Cosmos gRPC query services", "error", err)
		return err
	}
	// Defer closing gRPC connection until job completes.
	defer func() {
		if connErr := grpcClient.CloseConnection(queryConn); connErr != nil {
			err = connErr
		}
	}()

	daemonConn, err := grpcClient.NewGrpcConnection(ctx, daemonFlags.Shared.SocketAddress)
	if err != nil {
		c.logger.Error("Failed to establish gRPC connection to socket address", "error", err)
		return err
	}
	// Defer closing gRPC connection until job completes.
	defer func() {
		if connErr := grpcClient.CloseConnection(daemonConn); connErr != nil {
			err = connErr
		}
	}()

	pricesQueryClient := pricestypes.NewQueryClient(queryConn)

	// 2. Validate daemon configuration.
	if err := validateDaemonConfiguration(
		exchangeIdToQueryConfig,
		exchangeIdToExchangeDetails,
	); err != nil {
		return err
	}

	// Let the canonical list of exchange feeds be the keys of the map of exchange feed ids to startup configs.
	canonicalExchangeIds := make([]types.ExchangeId, 0, len(exchangeIdToQueryConfig))
	for exchangeId := range exchangeIdToQueryConfig {
		canonicalExchangeIds = append(canonicalExchangeIds, exchangeId)
	}

	// 3. Initialize synchronized, in-memory shared daemon configuration.
	priceFeedMutableMarketConfigs := types.NewPriceFeedMutableMarketConfigs(
		canonicalExchangeIds,
	)

	exchangeToMarketPrices, err := types.NewExchangeToMarketPrices(canonicalExchangeIds)
	if err != nil {
		return err
	}

	// 4. Start PriceEncoder and PriceFetcher per exchange.
	timeProvider := &libtime.TimeProviderImpl{}
	for _exchangeId := range exchangeIdToQueryConfig {
		// Assign these within the loop to avoid unexpected values being passed to the goroutines.
		exchangeId := _exchangeId
		exchangeConfig := exchangeIdToQueryConfig[exchangeId]

		// Expect an ExchangeQueryDetails to exist for each supported exchange feed id.
		exchangeDetails, exists := exchangeIdToExchangeDetails[exchangeId]
		if !exists {
			return fmt.Errorf("no exchange details exists for exchangeId: %v", exchangeId)
		}

		// Instantiate shared buffered channel to be written to by the price fetcher and read from
		// by the price encoder.
		bCh := make(chan *price_fetcher.PriceFetcherSubtaskResponse, constants.FixedBufferSize)

		c.runningSubtasksWaitGroup.Add(1)
		go func() {
			defer c.runningSubtasksWaitGroup.Done()
			subTaskRunner.StartPriceEncoder(
				exchangeId,
				priceFeedMutableMarketConfigs,
				exchangeToMarketPrices,
				c.logger,
				bCh,
			)
		}()

		ticker, stop := c.newTickerWithStop(int(exchangeConfig.IntervalMs))
		c.runningSubtasksWaitGroup.Add(1)
		go func() {
			defer c.runningSubtasksWaitGroup.Done()
			subTaskRunner.StartPriceFetcher(
				ticker,
				stop,
				priceFeedMutableMarketConfigs,
				*exchangeConfig,
				exchangeDetails,
				&handler.ExchangeQueryHandlerImpl{TimeProvider: timeProvider},
				c.logger,
				bCh,
			)
		}()
	}

	// 5. Start MarketUpdater go routine to periodically update the market configs.
	marketParamUpdaterTicker, marketParamUpdaterStop := c.newTickerWithStop(constants.MarketUpdateIntervalMs)
	c.runningSubtasksWaitGroup.Add(1)
	go func() {
		defer c.runningSubtasksWaitGroup.Done()
		subTaskRunner.StartMarketParamUpdater(
			ctx,
			marketParamUpdaterTicker,
			marketParamUpdaterStop,
			priceFeedMutableMarketConfigs,
			pricesQueryClient,
			c.logger,
		)
	}()

	// 6. Start PriceUpdater to begin broadcasting prices.
	// `StartPriceUpdater` does not run in a go-routine since it is used to block indefinitely
	// until the pricefeed daemon ends.
	// The price updater will read from an in-memory cache and send updates over gRPC for the
	// server to read.

	priceUpdaterTicker, priceUpdaterStop := c.newTickerWithStop(int(daemonFlags.Price.LoopDelayMs))

	// Now that all persistent subtasks have been started and all tickers and stop channels are created,
	// signal that the startup process is complete. This needs to be called before entering the
	// price updater loop, which loops indefinitely until the daemon is stopped.
	c.daemonStartup.Done()

	pricefeedClient := api.NewPriceFeedServiceClient(daemonConn)
	subTaskRunner.StartPriceUpdater(
		c,
		ctx,
		priceUpdaterTicker,
		priceUpdaterStop,
		exchangeToMarketPrices,
		pricefeedClient,
		c.logger,
	)
	return nil
}

// StartNewClient initializes and starts a new pricefeed daemon as a subtask of the calling process.
// The pricefeed daemon is a job that periodically queries external exchanges and transmits
// price data to the pricefeed service, which is then used by the application to compute index
// prices for proposing and validating oracle price updates on the blockchain.
// Note: the daemon will panic if it fails to start up.
func StartNewClient(
	ctx context.Context,
	daemonFlags flags.DaemonFlags,
	appFlags appflags.Flags,
	logger log.Logger,
	grpcClient daemontypes.GrpcClient,
	exchangeIdToQueryConfig map[types.ExchangeId]*types.ExchangeQueryConfig,
	exchangeIdToExchangeDetails map[types.ExchangeId]types.ExchangeQueryDetails,
	subTaskRunner SubTaskRunner,
) (client *Client) {
	// Log the daemon flags.
	logger.Info(
		"Starting pricefeed daemon with flags",
		"PriceFlags", daemonFlags.Price,
	)

	client = newClient(logger)
	client.runningSubtasksWaitGroup.Add(1)
	go func() {
		defer client.runningSubtasksWaitGroup.Done()
		err := client.start(
			ctx,
			daemonFlags,
			appFlags,
			grpcClient,
			exchangeIdToQueryConfig,
			exchangeIdToExchangeDetails,
			subTaskRunner,
		)
		if err != nil {
			logger.Error("Error initializing pricefeed daemon: %w", err.Error())
			panic(err)
		}
	}()
	return client
}

// validateDaemonConfiguration validates the daemon configuration.
// The list of exchanges used as keys for the exchangeIdToQueryConfig defines the exchanges used
// by the daemon.
// The daemon configuration is valid iff:
// 1) The exchangeIdToExchangeDetails map has an entry for each exchange.
// 2) The static exchange names map has an entry for each exchange, and each name is unique.
func validateDaemonConfiguration(
	exchangeIdToQueryConfig map[types.ExchangeId]*types.ExchangeQueryConfig,
	exchangeIdToExchangeDetails map[types.ExchangeId]types.ExchangeQueryDetails,
) (
	err error,
) {
	// Loop through all exchanges, which are defined by the exchangeIdToQueryConfig map,
	// and validate all ids are unique and have a corresponding ExchangeQueryDetails.
	exchangeIds := make(map[string]struct{}, len(exchangeIdToQueryConfig))
	for exchangeId := range exchangeIdToQueryConfig {
		if _, exists := exchangeIds[exchangeId]; exists {
			return fmt.Errorf("duplicate exchange id '%v' found for exchangeIds", exchangeId)
		}
		exchangeIds[exchangeId] = struct{}{}

		// Expect an ExchangeQueryDetails to exist for each supported exchange feed id.
		if _, exists := exchangeIdToExchangeDetails[exchangeId]; !exists {
			return fmt.Errorf("no exchange details exists for exchangeId: %v", exchangeId)
		}
	}

	// Validate that there is at least 1 exchange.
	if len(exchangeIds) == 0 {
		return errors.New("exchangeIds must not be empty")
	}

	return nil
}
