package client

import (
	"context"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	timelib "github.com/dydxprotocol/v4-chain/protocol/lib/time"

	"cosmossdk.io/log"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// Client implements a daemon service client that periodically calculates and reports liquidatable subaccounts
// to the protocol.
type Client struct {
	// Query clients
	BlocktimeQueryClient     blocktimetypes.QueryClient
	SubaccountQueryClient    satypes.QueryClient
	PerpetualsQueryClient    perptypes.QueryClient
	PricesQueryClient        pricestypes.QueryClient
	ClobQueryClient          clobtypes.QueryClient
	LiquidationServiceClient api.LiquidationServiceClient

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

	// Initialize the query clients. These are used to query the Cosmos gRPC query services.
	c.BlocktimeQueryClient = blocktimetypes.NewQueryClient(queryConn)
	c.SubaccountQueryClient = satypes.NewQueryClient(queryConn)
	c.PerpetualsQueryClient = perptypes.NewQueryClient(queryConn)
	c.PricesQueryClient = pricestypes.NewQueryClient(queryConn)
	c.ClobQueryClient = clobtypes.NewQueryClient(queryConn)
	c.LiquidationServiceClient = api.NewLiquidationServiceClient(daemonConn)

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
) {
	for {
		select {
		case <-ticker.C:
			if err := s.RunLiquidationDaemonTaskLoop(
				ctx,
				client,
				flags.Liquidation,
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
