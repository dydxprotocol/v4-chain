package client

import (
	"context"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	timelib "github.com/dydxprotocol/v4-chain/protocol/lib/time"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// Client implements a daemon service client that periodically calculates and reports liquidatable subaccounts
// to the protocol.
type Client struct {
	// include HealthCheckable to track the health of the daemon.
	daemontypes.HealthCheckable

	// logger is the logger for the daemon.
	logger log.Logger
}

// Ensure Client implements the HealthCheckable interface.
var _ daemontypes.HealthCheckable = (*Client)(nil)

func NewClient(logger log.Logger) *Client {
	return &Client{
		HealthCheckable: daemontypes.NewTimeBoundedHealthCheckable(
			types.LiquidationsDaemonServiceName,
			&timelib.TimeProviderImpl{},
			logger,
		),
		logger: logger,
	}
}

// Start begins a job that periodically:
// 1) Queries a gRPC server for all subaccounts including their open positions.
// 2) Checks collateralization statuses of subaccounts with at least one open position.
// 3) Sends a list of subaccount ids that potentially need to be liquidated to the application.
func (c *Client) Start(
	ctx context.Context,
	flags flags.DaemonFlags,
	appFlags appflags.Flags,
	grpcClient daemontypes.GrpcClient,
) error {
	// Log the daemon flags.
	c.logger.Info(
		"Starting liquidations daemon with flags",
		"LiquidationFlags", flags.Liquidation,
	)

	// Make a connection to the Cosmos gRPC query services.
	queryConn, err := grpcClient.NewTcpConnection(ctx, appFlags.GrpcAddress)
	if err != nil {
		c.logger.Error("Failed to establish gRPC connection to Cosmos gRPC query services", "error", err)
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
		c.logger.Error("Failed to establish gRPC connection to socket address", "error", err)
		return err
	}
	defer func() {
		if connErr := grpcClient.CloseConnection(daemonConn); connErr != nil {
			err = connErr
		}
	}()

	subaccountQueryClient := satypes.NewQueryClient(queryConn)
	clobQueryClient := clobtypes.NewQueryClient(queryConn)
	liquidationServiceClient := api.NewLiquidationServiceClient(daemonConn)

	ticker := time.NewTicker(time.Duration(flags.Liquidation.LoopDelayMs) * time.Millisecond)
	stop := make(chan bool)

	s := &SubTaskRunnerImpl{}
	StartLiquidationsDaemonTaskLoop(
		c,
		ctx,
		s,
		flags,
		ticker,
		stop,
		subaccountQueryClient,
		clobQueryClient,
		liquidationServiceClient,
	)

	return nil
}

// StartLiquidationsDaemonTaskLoop contains the logic to periodically run the liquidations daemon task.
func StartLiquidationsDaemonTaskLoop(
	client *Client,
	ctx context.Context,
	s SubTaskRunner,
	flags flags.DaemonFlags,
	ticker *time.Ticker,
	stop <-chan bool,
	subaccountQueryClient satypes.QueryClient,
	clobQueryClient clobtypes.QueryClient,
	liquidationServiceClient api.LiquidationServiceClient,
) {
	for {
		select {
		case <-ticker.C:
			if err := s.RunLiquidationDaemonTaskLoop(
				client,
				ctx,
				flags.Liquidation,
				subaccountQueryClient,
				clobQueryClient,
				liquidationServiceClient,
			); err != nil {
				// TODO(DEC-947): Move daemon shutdown to application.
				client.logger.Error("Liquidations daemon returned error", "error", err)
				client.ReportFailure(err)
			} else {
				client.ReportSuccess()
			}
		case <-stop:
			return
		}
	}
}
