package client

import (
	"context"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/client/types/constants"
	libtime "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/api"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client encapsulates the logic and interface for the bridge daemon. The bridge daemon periodically queries the
// Ethereum blockchain for new bridge events and relays them to the Cosmos gRPC server.
type Client struct {
	daemontypes.HealthCheckable
}

func NewClient() *Client {
	return &Client{
		HealthCheckable: daemontypes.NewTimeBoundedHealthCheckable(
			constants.BridgeDaemonModuleName,
			&libtime.TimeProviderImpl{},
		),
	}
}

// Start begins a job that periodically runs the RunBridgeDaemonTaskLoop function.
func (c *Client) Start(
	ctx context.Context,
	flags daemonflags.DaemonFlags,
	appFlags appflags.Flags,
	logger log.Logger,
	grpcClient daemontypes.GrpcClient,
) error {
	// Log the daemon flags.
	logger.Info(
		"Starting bridge daemon with flags",
		"BridgeFlags", flags.Bridge,
	)

	// Panic if EthRpcEndpoint is empty.
	if flags.Bridge.EthRpcEndpoint == "" {
		return fmt.Errorf("flag %s is not set", daemonflags.FlagBridgeDaemonEthRpcEndpoint)
	}

	// Make a connection to the Cosmos gRPC query services.
	queryConn, err := grpcClient.NewTcpConnection(ctx, appFlags.GrpcAddress)
	if err != nil {
		logger.Error("Failed to establish gRPC connection to Cosmos gRPC query services", "error", err)
		return err
	}
	defer func() {
		if connErr := grpcClient.CloseConnection(queryConn); connErr != nil {
			logger.Error("Failed to close gRPC connection to Cosmos gRPC query services", "error", connErr)
		}
	}()

	// Make a connection to the private daemon gRPC server.
	daemonConn, err := grpcClient.NewGrpcConnection(ctx, flags.Shared.SocketAddress)
	if err != nil {
		logger.Error("Failed to establish gRPC connection to socket address", "error", err)
		return err
	}
	defer func() {
		if connErr := grpcClient.CloseConnection(daemonConn); connErr != nil {
			logger.Error("Failed to close gRPC connection to Cosmos gRPC query services", "error", connErr)
		}
	}()

	// Initialize gRPC clients from query connection and daemon server connection.
	queryClient := bridgetypes.NewQueryClient(queryConn)
	serviceClient := api.NewBridgeServiceClient(daemonConn)

	// Initialize an Ethereum client from an RPC endpoint.
	ethClient, err := ethclient.Dial(flags.Bridge.EthRpcEndpoint)
	if err != nil {
		logger.Error("Failed to establish connection to Ethereum node", "error", err)
		return err
	}
	defer func() { ethClient.Close() }()

	ticker := time.NewTicker(time.Duration(flags.Bridge.LoopDelayMs) * time.Millisecond)
	stop := make(chan bool, 1)
	// Run the main task loop at an interval.
	StartBridgeDaemonTaskLoop(
		ctx,
		c,
		ticker,
		stop,
		&SubTaskRunnerImpl{},
		logger,
		ethClient,
		queryClient,
		serviceClient,
	)

	return nil
}

// StartBridgeDaemonTaskLoop operates the continuous loop that runs the bridge daemon. It receives as arguments
// a ticker and a stop channel that are used to control and halt the loop.
func StartBridgeDaemonTaskLoop(
	ctx context.Context,
	c *Client,
	ticker *time.Ticker,
	stop <-chan bool,
	s SubTaskRunner,
	logger log.Logger,
	ethClient types.EthClient,
	queryClient bridgetypes.QueryClient,
	serviceClient api.BridgeServiceClient,
) {
	// Run the main task loop at an interval.
	for {
		select {
		case <-ticker.C:
			if err := s.RunBridgeDaemonTaskLoop(
				ctx,
				logger,
				ethClient,
				queryClient,
				serviceClient,
			); err == nil {
				c.ReportSuccess()
			} else {
				// TODO(DEC-947): Move daemon shutdown to application.
				logger.Error("Bridge daemon returned error", "error", err)
				c.ReportFailure(err)
			}
		case <-stop:
			return
		}
	}
}
