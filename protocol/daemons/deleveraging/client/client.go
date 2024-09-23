package client

import (
	"context"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types"
	timelib "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/time"

	"cosmossdk.io/log"
	appflags "github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/deleveraging/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	daemontypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/types"
	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

// Client implements a daemon service client that periodically gets a list of all subaccounts with open positions on perpetuals
type Client struct {
	// Query clients
	BlocktimeQueryClient      blocktimetypes.QueryClient
	SubaccountQueryClient     satypes.QueryClient
	DeleveragingServiceClient api.DeleveragingServiceClient

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
			types.DeleveragingDaemonServiceName,
			&timelib.TimeProviderImpl{},
			logger,
		),
		logger: logger,
	}
}

// Start begins a job that periodically:
// 1) Queries a gRPC server for all subaccounts including their open positions.
// 2) Sends a list of subaccount ids with open positions for each perpetual.
func (c *Client) Start(
	ctx context.Context,
	flags flags.DaemonFlags,
	appFlags appflags.Flags,
	grpcClient daemontypes.GrpcClient,
) error {
	// Log the daemon flags.
	c.logger.Info(
		"Starting deleveraging daemon with flags",
		"DeleveragingFlags", flags.Deleveraging,
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
	c.DeleveragingServiceClient = api.NewDeleveragingServiceClient(daemonConn)

	ticker := time.NewTicker(time.Duration(flags.Deleveraging.LoopDelayMs) * time.Millisecond)
	stop := make(chan bool)

	s := &SubTaskRunnerImpl{}
	StartDeleveragingDaemonTaskLoop(
		c,
		ctx,
		s,
		flags,
		ticker,
		stop,
	)

	return nil
}

// StartDeleveragingDaemonTaskLoop contains the logic to periodically run the deleveraging daemon task.
func StartDeleveragingDaemonTaskLoop(
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
			if err := s.RunDeleveragingDaemonTaskLoop(
				ctx,
				client,
				flags.Deleveraging,
			); err != nil {
				// TODO(DEC-947): Move daemon shutdown to application.
				client.logger.Error("Deleveraging daemon returned error", "error", err)
				client.ReportFailure(err)
			} else {
				client.ReportSuccess()
			}
		case <-stop:
			return
		}
	}
}
