package grpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.FullNodeStreamingManager = (*NoopGrpcStreamingManager)(nil)

type NoopGrpcStreamingManager struct{}

func NewNoopGrpcStreamingManager() *NoopGrpcStreamingManager {
	return &NoopGrpcStreamingManager{}
}

func (sm *NoopGrpcStreamingManager) Enabled() bool {
	return false
}

func (sm *NoopGrpcStreamingManager) Subscribe(
	req clobtypes.StreamOrderbookUpdatesRequest,
	srv clobtypes.Query_StreamOrderbookUpdatesServer,
) (
	err error,
) {
	return types.ErrNotImplemented
}

func (sm *NoopGrpcStreamingManager) SendOrderbookUpdates(
	updates *clobtypes.OffchainUpdates,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
}

func (sm *NoopGrpcStreamingManager) SendOrderbookFillUpdates(
	orderbookFills []clobtypes.StreamOrderbookFill,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
}

func (sm *NoopGrpcStreamingManager) InitializeNewStreams(
	getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
}

func (sm *NoopGrpcStreamingManager) Stop() {
}
