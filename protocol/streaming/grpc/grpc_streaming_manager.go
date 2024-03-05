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
	nextSubscriptionId     uint32
}

// OrderbookSubscription represents a active subscription to the orderbook updates stream.
type OrderbookSubscription struct {
	clobPairIds []uint32
	srv         clobtypes.Query_StreamOrderbookUpdatesServer
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
func (sm *GrpcStreamingManagerImpl) Subscribe(
	req clobtypes.StreamOrderbookUpdatesRequest,
	srv clobtypes.Query_StreamOrderbookUpdatesServer,
) (
	err error,
) {
	clobPairIds := req.GetClobPairId()

	// Perform some basic validation on the request.
	if len(clobPairIds) == 0 {
		return clobtypes.ErrInvalidGrpcStreamingRequest
	}

	subscription := &OrderbookSubscription{
		clobPairIds: clobPairIds,
		srv:         srv,
	}

	sm.Lock()
	defer sm.Unlock()

	sm.orderbookSubscriptions[sm.nextSubscriptionId] = subscription
	sm.nextSubscriptionId++

	return nil
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
		v1update, err := GetOffchainUpdatesV1(update)
		if err != nil {
			panic(err)
		}
		v1updates[clobPairId] = v1update
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
	// If a Send update has failed for any clob pair id, the whole subscription will be removed.
	for _, id := range idsToRemove {
		delete(sm.orderbookSubscriptions, id)
	}
}

// GetOffchainUpdatesV1 unmarshals messages in offchain updates to OffchainUpdateV1.
func GetOffchainUpdatesV1(offchainUpdates *clobtypes.OffchainUpdates) ([]ocutypes.OffChainUpdateV1, error) {
	v1updates := make([]ocutypes.OffChainUpdateV1, 0)
	for _, message := range offchainUpdates.Messages {
		var update ocutypes.OffChainUpdateV1
		err := proto.Unmarshal(message.Message.Value, &update)
		if err != nil {
			return nil, err
		}
		v1updates = append(v1updates, update)
	}
	return v1updates, nil
}
