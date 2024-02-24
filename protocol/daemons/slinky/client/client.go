package client

import (
	"context"
	"cosmossdk.io/errors"
	"sync"
	"time"

	"cosmossdk.io/log"

	oracleclient "github.com/skip-mev/slinky/service/clients/oracle"

	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
)

// Ensure Client is HealthCheckable
var _ daemontypes.HealthCheckable = (*Client)(nil)

// Client is the daemon implementation for pulling price data from the slinky sidecar.
type Client struct {
	daemontypes.HealthCheckable
	ctx               context.Context
	cf                context.CancelFunc
	marketPairFetcher *MarketPairFetcher
	priceFetcher      *PriceFetcher
	wg                sync.WaitGroup
	logger            log.Logger
}

func newClient(ctx context.Context, logger log.Logger) *Client {
	logger = logger.With(log.ModuleKey, SlinkyClientDaemonModuleName)
	client := &Client{
		HealthCheckable: daemontypes.NewTimeBoundedHealthCheckable(
			SlinkyClientDaemonModuleName,
			&libtime.TimeProviderImpl{},
			logger,
		),
		logger: logger,
	}
	client.ctx, client.cf = context.WithCancel(ctx)
	return client
}

// start creates the main goroutines of the Client.
func (c *Client) start(
	slinky oracleclient.OracleClient,
	grpcClient daemontypes.GrpcClient,
	daemonFlags flags.DaemonFlags,
	appFlags appflags.Flags,
) error {
	// 1. Start the MarketPairFetcher
	c.marketPairFetcher = NewMarketPairFetcher(c.logger)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.RunMarketPairFetcher(c.ctx, appFlags, grpcClient)
	}()
	// 2. Start the PriceFetcher
	c.priceFetcher = NewPriceFetcher(
		c.marketPairFetcher,
		grpcClient,
		daemonFlags.Shared.SocketAddress,
		slinky,
		c.logger,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.RunPriceFetcher(c.ctx)
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
				c.ReportFailure(errors.Wrap(err, "failed to run PriceFetcher for slinky daemon"))
			} else {
				c.ReportSuccess()
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
			err := c.marketPairFetcher.FetchIdMappings(ctx)
			if err != nil {
				c.logger.Error("Failed to run fetch id mappings for slinky daemon", "error", err)
				c.ReportFailure(errors.Wrap(err, "failed to run FetchIdMappings for slinky daemon"))
			}
			c.ReportSuccess()
		case <-ctx.Done():
			return
		}
	}
}

// StartNewClient creates and runs a Client.
func StartNewClient(
	ctx context.Context,
	slinky oracleclient.OracleClient,
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
	err := client.start(slinky, grpcClient, daemonFlags, appFlags)
	if err != nil {
		logger.Error("Error initializing slinky daemon: %w", err)
		panic(err)
	}
	return client
}
