package client

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/log"
	appflags "github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	daemonflags "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/types"
	daemontypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/types"
	libtime "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/time"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client encapsulates the logic and interface for the sDAI daemon. The sDAI daemon periodically queries the
// Ethereum blockchain for new sDAI conversion rate and relays them to the Cosmos gRPC server.
type Client struct {
	daemontypes.HealthCheckable

	// logger is the logger used by the sDAI daemon.
	logger log.Logger
}

func NewClient(logger log.Logger) *Client {
	logger = logger.With(log.ModuleKey, types.SDAIOracleDaemonModuleName)
	return &Client{
		HealthCheckable: daemontypes.NewTimeBoundedHealthCheckable(
			types.SDAIOracleDaemonModuleName,
			&libtime.TimeProviderImpl{},
			logger,
		),
		logger: logger,
	}
}

// Start begins a job that periodically runs the RunSDAIDaemonTaskLoop function.
func (c *Client) Start(
	ctx context.Context,
	flags daemonflags.DaemonFlags,
	appFlags appflags.Flags,
	grpcClient daemontypes.GrpcClient,
) error {
	// Log the daemon flags.
	c.logger.Info(
		"Starting sDAI daemon with flags",
		"SDAIFlags", flags.SDAI,
	)

	// Panic if EthRpcEndpoint is empty.
	if flags.SDAI.EthRpcEndpoint == "" {
		return fmt.Errorf("flag %s is not set", daemonflags.FlagSDAIDaemonEthRpcEndpoint)
	}

	// Make a connection to the private daemon gRPC server.
	daemonConn, err := grpcClient.NewGrpcConnection(ctx, flags.Shared.SocketAddress)
	if err != nil {
		c.logger.Error("Failed to establish gRPC connection to socket address", "error", err)
		return err
	}
	defer func() {
		if connErr := grpcClient.CloseConnection(daemonConn); connErr != nil {
			c.logger.Error("Failed to close gRPC connection to Cosmos gRPC query services", "error", connErr)
		}
	}()

	// Initialize gRPC clients from query connection and daemon server connection.
	serviceClient := api.NewSDAIServiceClient(daemonConn)

	// Initialize an Ethereum client from an RPC endpoint.
	ethClient, err := ethclient.Dial(flags.SDAI.EthRpcEndpoint)
	if err != nil {
		c.logger.Error("Failed to establish connection to Ethereum node", "error", err)
		return err
	}
	defer func() { ethClient.Close() }()

	ticker := time.NewTicker(time.Duration(flags.SDAI.LoopDelayMs) * time.Millisecond)
	stop := make(chan bool, 1)
	// Run the main task loop at an interval.
	StartsDAIDaemonTaskLoop(
		ctx,
		c,
		ticker,
		stop,
		&SubTaskRunnerImpl{},
		ethClient,
		serviceClient,
	)

	return nil
}

// StartsDAIDaemonTaskLoop operates the continuous loop that runs the sDAI daemon. It receives as arguments
// a ticker and a stop channel that are used to control and halt the loop.
func StartsDAIDaemonTaskLoop(
	ctx context.Context,
	c *Client,
	ticker *time.Ticker,
	stop <-chan bool,
	s SubTaskRunner,
	ethClient *ethclient.Client,
	serviceClient api.SDAIServiceClient,
) {
	// Run the main task loop at an interval.
	for {
		select {
		case <-ticker.C:
			if err := s.RunsDAIDaemonTaskLoop(
				ctx,
				c,
				c.logger,
				ethClient,
				serviceClient,
			); err == nil {
				c.ReportSuccess()
			} else {
				// TODO(DEC-947): Move daemon shutdown to application.
				c.logger.Error("SDAI daemon returned error", "error", err)
				c.ReportFailure(err)
			}
		case <-stop:
			return
		}
	}
}
