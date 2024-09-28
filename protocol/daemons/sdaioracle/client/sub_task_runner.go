package client

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	ethqueryclienttypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/eth_query_client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SubTaskRunner interface {
	RunsDAIDaemonTaskLoop(
		ctx context.Context,
		logger log.Logger,
		ethClient *ethclient.Client,
		queryClient ethqueryclienttypes.EthQueryClient,
		serviceClient api.SDAIServiceClient,
	) error
}

type SubTaskRunnerImpl struct{}

var _ SubTaskRunner = (*SubTaskRunnerImpl)(nil)

// RunsDAIDaemonTaskLoop does the following:
// 1) Fetches sDAI conversion rate from a configured Ethereum client.
// 2) Sends sDAI conversion rate to the gRPC server.
func (s *SubTaskRunnerImpl) RunsDAIDaemonTaskLoop(
	ctx context.Context,
	logger log.Logger,
	ethClient *ethclient.Client,
	queryClient ethqueryclienttypes.EthQueryClient,
	serviceClient api.SDAIServiceClient,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.SDAIDaemon,
		time.Now(),
		metrics.MainTaskLoop,
		metrics.Latency,
	)

	// Verify Chain ID.
	chainId, err := queryClient.ChainID(ctx, ethClient)
	if err != nil {
		return fmt.Errorf("failed to fetch chain ID: %w", err)
	}
	if chainId.Uint64() != types.EthChainID {
		return fmt.Errorf(
			"expected chain ID %d but node has chain ID %d",
			types.EthChainID,
			chainId,
		)
	}

	// Call the QueryDaiConversionRate function
	sDAIConversionRate, err := queryClient.QueryDaiConversionRate(ethClient)
	if err != nil {
		return fmt.Errorf("failed to query DAI conversion rate: %w", err)
	}

	telemetry.IncrCounter(
		1,
		metrics.SDAIDaemon,
		metrics.Count,
	)

	// Send sDAI events to sDAI server.
	if _, err = serviceClient.AddsDAIEvent(ctx, &api.AddsDAIEventsRequest{
		ConversionRate: sDAIConversionRate,
	}); err != nil {
		return fmt.Errorf("failed to add sDAI events: %w", err)
	}

	// Success.
	return nil
}
