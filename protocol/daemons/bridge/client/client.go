package client

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	daemontypes "github.com/dydxprotocol/v4-chain/protocol/daemons/types"
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
	appFlags appflags.Flags,
	logger log.Logger,
	grpcClient daemontypes.GrpcClient,
) error {
	// Make a connection to the Cosmos gRPC query services.
	queryConn, err := grpcClient.NewTcpConnection(ctx, appFlags.GrpcAddress)
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

	// Run the main task loop at an interval.
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
// 2) Fetches Ethereum events from a configured Ethereum client.
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

	// Fetch parameters from x/bridge module. Relevant ones to bridge daemon are:
	// - EventParams
	//   - ChainId: Ethereum chain ID that bridge contract resides on.
	//   - EthAddress: Address of the bridge contract to query events from.
	// - ProposeParams
	//   - MaxBridgesPerBlock: Number of bridge events to query for.
	// - RecognizedEventInfo
	//   - EthBlockHeight: Ethereum block height from which to start querying events.
	//   - NextId: Next bridge event ID to query for.
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
	newBridgeEvents := make([]bridgetypes.BridgeEvent, 0, len(logs))
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

	// Success.
	return nil
}

// getFilterQuery returns a query to fetch logs of bridge events with following filters:
// - logs are emitted by contract at address `contractAddressHex`.
// - block height is between `fromBlock` and current finalized block height (both inclusive).
// - event IDs are sequential integers between `firstId` and `firstId + numIds - 1` (both inclusive).
func getFilterQuery(
	contractAddressHex string,
	fromBlock uint64,
	firstId uint32,
	numIds uint32,
) eth.FilterQuery {
	// Generate `ethcommon.Hash`s of the next `numIds` event IDs.
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
