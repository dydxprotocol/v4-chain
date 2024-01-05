package client

import (
	"context"
	"cosmossdk.io/log"
	"fmt"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	libeth "github.com/dydxprotocol/v4-chain/protocol/lib/eth"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"
)

type SubTaskRunner interface {
	RunBridgeDaemonTaskLoop(
		ctx context.Context,
		logger log.Logger,
		ethClient types.EthClient,
		queryClient bridgetypes.QueryClient,
		serviceClient api.BridgeServiceClient,
	) error
}

type SubTaskRunnerImpl struct{}

var _ SubTaskRunner = (*SubTaskRunnerImpl)(nil)

// RunBridgeDaemonTaskLoop does the following:
// 1) Fetches configuration information by querying the gRPC server.
// 2) Fetches Ethereum events from a configured Ethereum client.
// 3) Sends newly-recognized bridge events to the gRPC server.
func (s *SubTaskRunnerImpl) RunBridgeDaemonTaskLoop(
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
		return fmt.Errorf("failed to fetch event params: %w", err)
	}
	proposeParams, err := queryClient.ProposeParams(ctx, &bridgetypes.QueryProposeParamsRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch propose params: %w", err)
	}
	recognizedEventInfo, err := queryClient.RecognizedEventInfo(ctx, &bridgetypes.QueryRecognizedEventInfoRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch recognized event info: %w", err)
	}

	// Verify Chain ID.
	chainId, err := ethClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch chain ID: %w", err)
	}
	if chainId.Uint64() != eventParams.Params.EthChainId {
		return fmt.Errorf(
			"expected chain ID %d but node has chain ID %d",
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
		return fmt.Errorf("failed to fetch logs: %w", err)
	}
	telemetry.IncrCounter(
		float32(len(logs)),
		metrics.BridgeDaemon,
		metrics.NewEthLogs,
		metrics.Count,
	)

	// Parse logs into bridge events.
	newBridgeEvents := make([]bridgetypes.BridgeEvent, len(logs))
	for i, log := range logs {
		newBridgeEvents[i] = libeth.BridgeLogToEvent(log, eventParams.Params.Denom)
	}

	// Send bridge events to bridge server.
	if _, err = serviceClient.AddBridgeEvents(ctx, &api.AddBridgeEventsRequest{
		BridgeEvents: newBridgeEvents,
	}); err != nil {
		return fmt.Errorf("failed to add bridge events: %w", err)
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
		eventIdHashes[i] = ethcommon.BigToHash(new(big.Int).SetUint64(uint64(firstId + i)))
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
