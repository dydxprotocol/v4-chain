package grpc

import (
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.GrpcStreamingManager = (*NoopGrpcStreamingManager)(nil)

// GrpcStreamingManager is an implementation for managing gRPC streaming subscriptions.
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
	finished chan<- bool,
	err error,
) {
	return nil, nil
}

// SendMessages groups messages by their clob pair ids and
// sends messages to the subscribers.
func (sm *NoopGrpcStreamingManager) SendMessages(
	msg *clobtypes.OffchainUpdateMessage,
) {
}
