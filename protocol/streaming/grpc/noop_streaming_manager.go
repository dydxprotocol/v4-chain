package grpc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/streaming/grpc/types"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
)

var _ types.GrpcStreamingManager = (*NoopGrpcStreamingManager)(nil)

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
	return clobtypes.ErrGrpcStreamingManagerNotEnabled
}

func (sm *NoopGrpcStreamingManager) SendOrderbookUpdates(
	updates *clobtypes.OffchainUpdates,
	snapshot bool,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
}

func (sm *NoopGrpcStreamingManager) GetUninitializedClobPairIds() []uint32 {
	return []uint32{}
}
