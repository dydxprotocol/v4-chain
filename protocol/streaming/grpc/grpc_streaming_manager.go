package grpc

import (
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.GrpcStreamingManager = (*GrpcStreamingManagerImpl)(nil)

// GrpcStreamingManager is an implementation for managing gRPC streaming subscriptions.
type GrpcStreamingManagerImpl struct {
}

func NewGrpcStreamingManager() *GrpcStreamingManagerImpl {
	return &GrpcStreamingManagerImpl{}
}

func (sm *GrpcStreamingManagerImpl) Enabled() bool {
	return true
}

// Subscribe subscribes to the orderbook updates stream.
func (sm *GrpcStreamingManagerImpl) Subscribe(
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
func (sm *GrpcStreamingManagerImpl) SendMessages(
	msg *clobtypes.OffchainUpdateMessage,
) {
}
