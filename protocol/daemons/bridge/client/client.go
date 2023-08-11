package client

import (
	"context"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4/daemons/bridge/api"
	"github.com/dydxprotocol/v4/daemons/flags"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	bridgetypes "github.com/dydxprotocol/v4/x/bridge/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Start begins a job that periodically runs the RunBridgeDaemonTaskLoop function.
func Start(
	ctx context.Context,
	flags flags.DaemonFlags,
	logger log.Logger,
	grpcClient lib.GrpcClient,
) error {
	// Make a connection to the Cosmos gRPC query services.
	queryConn, err := grpcClient.NewTcpConnection(ctx, flags.Shared.GrpcServerAddress)
	if err != nil {
		logger.Error("Failed to establish gRPC connection to Cosmos gRPC query services", "error", err)
		return err
	}
	defer func() {
		if connErr := grpcClient.CloseConnection(queryConn); connErr != nil {
			err = connErr
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
			err = connErr
		}
	}()

	queryClient := bridgetypes.NewQueryClient(queryConn)
	serviceClient := api.NewBridgeServiceClient(daemonConn)

	ethClient, err := ethclient.Dial(flags.Bridge.EthRpcEndpoint)
	if err != nil {
		logger.Error("Failed to establish connection to Ethereum Node", "error", err)
		return err
	}
	defer func() { ethClient.Close() }()

	ticker := time.NewTicker(time.Duration(flags.Bridge.LoopDelayMs) * time.Millisecond)
	for ; true; <-ticker.C {
		if err := RunBridgeDaemonTaskLoop(
			ctx,
			logger,
			ethClient,
			queryClient,
			serviceClient,
		); err != nil {
			// TODO(DEC-947): Move daemon shutdown to application.
			logger.Error("Bridge daemon returned error", "error", err)
		}
	}

	return nil
}

// RunBridgeDaemonTaskLoop does the following:
// 1) Fetches configuration information by querying the gRPC server. (TODO: CORE-318)
// 2) Fetches Ethereum events from a configured node. (TODO: CORE-319)
// 3) Sends newly-recognized bridge events to the gRPC server. (TODO: CORE-320)
func RunBridgeDaemonTaskLoop(
	ctx context.Context,
	logger log.Logger,
	ethClient *ethclient.Client,
	queryClient bridgetypes.QueryClient,
	serviceClient api.BridgeServiceClient,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.BridgeDaemon,
		time.Now(),
		metrics.MainTaskLoop,
		metrics.Latency,
	)

	// Success
	return nil
}
