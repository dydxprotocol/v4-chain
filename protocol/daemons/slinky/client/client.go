package client

import (
	"context"
	"sync"
	"time"

	"cosmossdk.io/errors"

	"cosmossdk.io/log"

	oracleclient "github.com/dydxprotocol/slinky/service/clients/oracle"

	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	pricefeedtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
)

// Client is the daemon implementation for pulling price data from the slinky sidecar.
type Client struct {
	ctx                   context.Context
	cf                    context.CancelFunc
	marketPairFetcher     MarketPairFetcher
	marketPairHC          daemontypes.HealthCheckable
	priceFetcher          PriceFetcher
	priceHC               daemontypes.HealthCheckable
	sidecarVersionChecker SidecarVersionChecker
	sidecarVersionHC      daemontypes.HealthCheckable
	wg                    sync.WaitGroup
	logger                log.Logger
}

func newClient(ctx context.Context, logger log.Logger) *Client {
	logger = logger.With(log.ModuleKey, SlinkyClientDaemonModuleName)
	client := &Client{
		marketPairHC: daemontypes.NewTimeBoundedHealthCheckable(
			SlinkyClientMarketPairFetcherDaemonModuleName,
			&libtime.TimeProviderImpl{},
			logger,
		),
		priceHC: daemontypes.NewTimeBoundedHealthCheckable(
			SlinkyClientPriceFetcherDaemonModuleName,
			&libtime.TimeProviderImpl{},
			logger,
		),
		sidecarVersionHC: daemontypes.NewTimeBoundedHealthCheckable(
			SlinkyClientSidecarVersionFetcherDaemonModuleName,
			&libtime.TimeProviderImpl{},
			logger,
		),
		logger: logger,
	}
	client.ctx, client.cf = context.WithCancel(ctx)
	return client
}

func (c *Client) GetMarketPairHC() daemontypes.HealthCheckable {
	return c.marketPairHC
}

func (c *Client) GetPriceHC() daemontypes.HealthCheckable {
	return c.priceHC
}

func (c *Client) GetSidecarVersionHC() daemontypes.HealthCheckable {
	return c.sidecarVersionHC
}

// start creates the main goroutines of the Client.
func (c *Client) start(
	slinky oracleclient.OracleClient,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	grpcClient daemontypes.GrpcClient,
	appFlags appflags.Flags,
) error {
	// 1. Start the MarketPairFetcher
	c.marketPairFetcher = NewMarketPairFetcher(c.logger)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.RunMarketPairFetcher(c.ctx, appFlags, grpcClient)
	}()

	// 2. Start the PriceFetcher
	c.priceFetcher = NewPriceFetcher(
		c.marketPairFetcher,
		indexPriceCache,
		slinky,
		c.logger,
	)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.RunPriceFetcher(c.ctx)
	}()

	// 3. Start the SidecarVersionChecker
	c.sidecarVersionChecker = NewSidecarVersionChecker(
		slinky,
		c.logger,
	)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.RunSidecarVersionChecker(c.ctx)
	}()
	return nil
}

// RunPriceFetcher periodically calls the priceFetcher to grab prices from the slinky sidecar and
// push them to the pricefeed server.
func (c *Client) RunPriceFetcher(ctx context.Context) {
	err := c.priceFetcher.Start(ctx)
	if err != nil {
		c.logger.Error("Error initializing PriceFetcher in slinky daemon: %w", err)
		panic(err)
	}
	ticker := time.NewTicker(SlinkyPriceFetchDelay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := c.priceFetcher.FetchPrices(ctx)
			if err != nil {
				c.logger.Error("Failed to run fetch prices for slinky daemon", "error", err)
				c.priceHC.ReportFailure(errors.Wrap(err, "failed to run PriceFetcher for slinky daemon"))
			} else {
				c.priceHC.ReportSuccess()
			}
		case <-ctx.Done():
			return
		}
	}
}

// Stop closes all connections and waits for goroutines to exit.
func (c *Client) Stop() {
	c.cf()
	c.priceFetcher.Stop()
	c.marketPairFetcher.Stop()
	c.wg.Wait()
}

// RunMarketPairFetcher periodically calls the marketPairFetcher to cache mappings between
// currency pair and market param ID.
func (c *Client) RunMarketPairFetcher(ctx context.Context, appFlags appflags.Flags, grpcClient daemontypes.GrpcClient) {
	err := c.marketPairFetcher.Start(ctx, appFlags, grpcClient)
	if err != nil {
		c.logger.Error("Error initializing MarketPairFetcher in slinky daemon: %w", err)
		panic(err)
	}
	ticker := time.NewTicker(SlinkyMarketParamFetchDelay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err = c.marketPairFetcher.FetchIdMappings(ctx)
			if err != nil {
				c.logger.Error("Failed to run fetch id mappings for slinky daemon", "error", err)
				c.marketPairHC.ReportFailure(errors.Wrap(err, "failed to run FetchIdMappings for slinky daemon"))
			} else {
				c.marketPairHC.ReportSuccess()
			}
		case <-ctx.Done():
			return
		}
	}
}

// RunSidecarVersionChecker periodically calls the sidecarVersionChecker to check if the running sidecar version
// is at least a minimum acceptable version.
func (c *Client) RunSidecarVersionChecker(ctx context.Context) {
	err := c.sidecarVersionChecker.Start(ctx)
	if err != nil {
		c.logger.Error("Error initializing sidecarVersionChecker in slinky daemon", "error", err)
		panic(err)
	}
	ticker := time.NewTicker(SlinkySidecarCheckDelay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err = c.sidecarVersionChecker.CheckSidecarVersion(ctx)
			if err != nil {
				c.logger.Error("Sidecar version check failed", "error", err)
				c.sidecarVersionHC.ReportFailure(errors.Wrap(err, "Sidecar version check failed for slinky daemon"))
			} else {
				c.sidecarVersionHC.ReportSuccess()
			}
		case <-ctx.Done():
			return
		}
	}
}

// StartNewClient creates and runs a Client.
// The client creates the MarketPairFetcher, PriceFetcher, and SidecarVersionChecker,
// connects to the required grpc services, and launches them in goroutines.
// It is non-blocking and returns on successful startup.
// If it hits a critical error in startup it panics.
func StartNewClient(
	ctx context.Context,
	slinky oracleclient.OracleClient,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	grpcClient daemontypes.GrpcClient,
	daemonFlags flags.DaemonFlags,
	appFlags appflags.Flags,
	logger log.Logger,
) *Client {
	logger.Info(
		"Starting slinky daemon with flags",
		"SlinkyFlags", daemonFlags.Slinky,
	)

	client := newClient(ctx, logger)
	err := client.start(slinky, indexPriceCache, grpcClient, appFlags)
	if err != nil {
		logger.Error("Error initializing slinky daemon: %w", err)
		panic(err)
	}
	return client
}
