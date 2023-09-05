package client

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	libeth "github.com/dydxprotocol/v4-chain/protocol/lib/eth"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
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
// 1) Fetches configuration information by querying the gRPC server.
// 2) Fetches Ethereum events from a configured node.
// 3) Sends newly-recognized bridge events to the gRPC server.
func RunBridgeDaemonTaskLoop(
	ctx context.Context,
	logger log.Logger,
	ethClient types.EthClient,
	queryClient bridgetypes.QueryClient,
	serviceClient api.BridgeServiceClient,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.BridgeDaemon,
		time.Now(),
		metrics.MainTaskLoop,
		metrics.Latency,
	)

	// Fetch parameters from x/bridge module.
	eventParams, err := queryClient.EventParams(ctx, &bridgetypes.QueryEventParamsRequest{})
	if err != nil {
		return err
	}
	proposeParams, err := queryClient.ProposeParams(ctx, &bridgetypes.QueryProposeParamsRequest{})
	if err != nil {
		return err
	}
	recognizedEventInfo, err := queryClient.RecognizedEventInfo(ctx, &bridgetypes.QueryRecognizedEventInfoRequest{})
	if err != nil {
		return err
	}

	// Verify Chain ID.
	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		return err
	}
	if chainId.Uint64() != eventParams.Params.EthChainId {
		return fmt.Errorf(
			"Expected chain ID %d but node has chain ID %d",
			eventParams.Params.EthChainId,
			chainId,
		)
	}

	// Fetch logs from Ethereum Node.
	filterQuery := getFilterQuery(
		eventParams.Params.EthAddress,
		recognizedEventInfo.Info.EthBlockHeight,
		recognizedEventInfo.Info.NextId,
		proposeParams.Params.MaxBridgesPerBlock,
	)
	logs, err := ethClient.FilterLogs(ctx, filterQuery)
	if err != nil {
		return err
	}
	telemetry.IncrCounter(
		float32(len(logs)),
		metrics.BridgeDaemon,
		metrics.NewEthLogs,
		metrics.Count,
	)

	// Parse logs into bridge events.
	newBridgeEvents := []bridgetypes.BridgeEvent{}
	for _, log := range logs {
		newBridgeEvents = append(
			newBridgeEvents,
			libeth.BridgeLogToEvent(log, eventParams.Params.Denom),
		)
	}

	// Send bridge events to bridge server.
	if _, err = serviceClient.AddBridgeEvents(ctx, &api.AddBridgeEventsRequest{
		BridgeEvents: newBridgeEvents,
	}); err != nil {
		return err
	}

	// Success
	return nil
}

// getFilterQuery returns a FilterQuery for fetching logs for the next `numIds`
// bridge events after block height `fromBlock` and before current finalized
// block height.
func getFilterQuery(
	contractAddressHex string,
	fromBlock uint64,
	firstId uint32,
	numIds uint32,
) eth.FilterQuery {
	// Generate bytes32 of the next x ids.
	eventIdHashes := make([]ethcommon.Hash, numIds)
	for i := uint32(0); i < numIds; i++ {
		h := ethcommon.BigToHash(big.NewInt(int64(firstId + i)))
		eventIdHashes = append(eventIdHashes, h)
	}

	return eth.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   big.NewInt(ethrpc.FinalizedBlockNumber.Int64()),
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(contractAddressHex)},
		Topics: [][]ethcommon.Hash{
			{ethcommon.HexToHash(constants.BridgeEventSignature)},
			eventIdHashes,
		},
	}
}
