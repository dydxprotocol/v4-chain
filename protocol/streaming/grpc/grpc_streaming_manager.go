package grpc

import (
	"sync"

	"github.com/cosmos/gogoproto/proto"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.GrpcStreamingManager = (*GrpcStreamingManagerImpl)(nil)

// GrpcStreamingManagerImpl is an implementation for managing gRPC streaming subscriptions.
type GrpcStreamingManagerImpl struct {
	sync.Mutex

	// orderbookSubscriptions maps subscription IDs to their respective orderbook subscriptions.
	orderbookSubscriptions map[uint32]*OrderbookSubscription
	nextId                 uint32
}

// OrderbookSubscription represents a active subscription to the orderbook updates stream.
type OrderbookSubscription struct {
	clobPairIds []uint32
	srv         clobtypes.Query_StreamOrderbookUpdatesServer
	finished    chan bool
}

func NewGrpcStreamingManager() *GrpcStreamingManagerImpl {
	return &GrpcStreamingManagerImpl{
		orderbookSubscriptions: make(map[uint32]*OrderbookSubscription),
	}
}

func (sm *GrpcStreamingManagerImpl) Enabled() bool {
	return true
}

// Subscribe subscribes to the orderbook updates stream.
// This function returns a channel that is used to signal termination when an error occurs.
func (sm *GrpcStreamingManagerImpl) Subscribe(
	req clobtypes.StreamOrderbookUpdatesRequest,
	srv clobtypes.Query_StreamOrderbookUpdatesServer,
) (
	finished chan bool,
	err error,
) {
	finished = make(chan bool)
	subscription := &OrderbookSubscription{
		clobPairIds: req.GetClobPairId(),
		srv:         srv,
		finished:    finished,
	}

	sm.Lock()
	defer sm.Unlock()

	sm.orderbookSubscriptions[sm.nextId] = subscription
	sm.nextId++

	return finished, nil
}

// SendOrderbookUpdates groups updates by their clob pair ids and
// sends messages to the subscribers.
func (sm *GrpcStreamingManagerImpl) SendOrderbookUpdates(
	offchainUpdates *clobtypes.OffchainUpdates,
) {
	// Group updates by clob pair id.
	updates := make(map[uint32]*clobtypes.OffchainUpdates)
	for _, message := range offchainUpdates.Messages {
		clobPairId := message.OrderId.ClobPairId
		if _, ok := updates[clobPairId]; !ok {
			updates[clobPairId] = clobtypes.NewOffchainUpdates()
		}
		updates[clobPairId].Messages = append(updates[clobPairId].Messages, message)
	}

	// Unmarshal messages to v1 updates.
	v1updates := make(map[uint32][]ocutypes.OffChainUpdateV1)
	for clobPairId, update := range updates {
		v1updates[clobPairId] = GetOffchainUpdatesV1(update)
	}

	sm.Lock()
	defer sm.Unlock()

	// Send updates to subscribers.
	idsToRemove := make([]uint32, 0)
	for id, subscription := range sm.orderbookSubscriptions {
		for _, clobPairId := range subscription.clobPairIds {
			if updates, ok := v1updates[clobPairId]; ok {
				if err := subscription.srv.Send(
					&clobtypes.StreamOrderbookUpdatesResponse{
						Updates:  updates,
						Snapshot: false,
					},
				); err != nil {
					idsToRemove = append(idsToRemove, id)
					break
				}
			}
		}
	}

	// Clean up subscriptions that have been closed.
	for _, id := range idsToRemove {
		sm.orderbookSubscriptions[id].finished <- true
		delete(sm.orderbookSubscriptions, id)
	}
}

// GetOffchainUpdatesV1 unmarshals messages in offchain updates to OffchainUpdateV1.
func GetOffchainUpdatesV1(offchainUpdates *clobtypes.OffchainUpdates) []ocutypes.OffChainUpdateV1 {
	v1updates := make([]ocutypes.OffChainUpdateV1, 0)
	for _, message := range offchainUpdates.Messages {
		var update ocutypes.OffChainUpdateV1
		proto.Unmarshal(message.Message.Value, &update)
		v1updates = append(v1updates, update)
	}
	return v1updates
}
