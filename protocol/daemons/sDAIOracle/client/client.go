package client

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/log"
	appflags "github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	daemonflags "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	ethqueryclienttypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/eth_query_client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/types"
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

	stop chan bool
}

func NewClient(logger log.Logger) *Client {
	logger = logger.With(log.ModuleKey, types.SDAIOracleDaemonModuleName)
	stop := make(chan bool, 1)
	return &Client{
		HealthCheckable: daemontypes.NewTimeBoundedHealthCheckable(
			types.SDAIOracleDaemonModuleName,
			&libtime.TimeProviderImpl{},
			logger,
		),
		logger: logger,
		stop:   stop,
	}
}

// Start begins a job that periodically runs the RunSDAIDaemonTaskLoop function.
func (c *Client) Start(
	ctx context.Context,
	flags daemonflags.DaemonFlags,
	appFlags appflags.Flags,
	grpcClient daemontypes.GrpcClient,
) error {

	c.logger.Info(
		"Starting sDAI daemon with flags",
		"SDAIFlags", flags.SDAI,
	)

	if flags.SDAI.EthRpcEndpoint == "" {
		return fmt.Errorf("flag %s is not set", daemonflags.FlagSDAIDaemonEthRpcEndpoint)
	}

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

	serviceClient := api.NewSDAIServiceClient(daemonConn)

	ethClient, err := ethclient.Dial(flags.SDAI.EthRpcEndpoint)
	if err != nil {
		c.logger.Error("Failed to establish connection to Ethereum node", "error", err)
		return err
	}
	defer func() { ethClient.Close() }()

	ticker := time.NewTicker(time.Duration(flags.SDAI.LoopDelayMs) * time.Millisecond)

	queryClient := &ethqueryclienttypes.EthQueryClientImpl{}
	StartsDAIDaemonTaskLoop(
		ctx,
		c,
		ticker,
		c.stop,
		&SubTaskRunnerImpl{},
		ethClient,
		queryClient,
		serviceClient,
	)

	return nil
}

// Stop signals the daemon to stop.
func (c *Client) Stop() {
	c.stop <- true
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
	queryClient ethqueryclienttypes.EthQueryClient,
	serviceClient api.SDAIServiceClient,
) {
	// Run the main task loop at an interval.
	for {
		select {
		case <-ticker.C:
			if err := s.RunsDAIDaemonTaskLoop(
				ctx,
				c.logger,
				ethClient,
				queryClient,
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
