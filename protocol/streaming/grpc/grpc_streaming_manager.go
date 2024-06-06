package grpc

import (
	"sync"
	"time"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.GrpcStreamingManager = (*GrpcStreamingManagerImpl)(nil)

// GrpcStreamingManagerImpl is an implementation for managing gRPC streaming subscriptions.
type GrpcStreamingManagerImpl struct {
	sync.Mutex

	logger log.Logger

	// orderbookSubscriptions maps subscription IDs to their respective orderbook subscriptions.
	orderbookSubscriptions map[uint32]*OrderbookSubscription
	nextSubscriptionId     uint32
}

// OrderbookSubscription represents a active subscription to the orderbook updates stream.
type OrderbookSubscription struct {
	// Initialize the subscription with orderbook snapshots.
	initialize sync.Once

	// Clob pair ids to subscribe to.
	clobPairIds []uint32

	// Stream
	srv clobtypes.Query_StreamOrderbookUpdatesServer
}

func NewGrpcStreamingManager(
	logger log.Logger,
) *GrpcStreamingManagerImpl {
	logger = logger.With(log.ModuleKey, "grpc-streaming")
	return &GrpcStreamingManagerImpl{
		logger:                 logger,
		orderbookSubscriptions: make(map[uint32]*OrderbookSubscription),
	}
}

func (sm *GrpcStreamingManagerImpl) Enabled() bool {
	return true
}

func (sm *GrpcStreamingManagerImpl) EmitMetrics() {
	metrics.SetGauge(
		metrics.GrpcStreamSubscriberCount,
		float32(len(sm.orderbookSubscriptions)),
	)
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
	sm.EmitMetrics()
	return nil
}

// SendOrderbookUpdates groups updates by their clob pair ids and
// sends messages to the subscribers.
func (sm *GrpcStreamingManagerImpl) SendOrderbookUpdates(
	offchainUpdates *clobtypes.OffchainUpdates,
	snapshot bool,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendOrderbookUpdatesLatency,
		time.Now(),
	)

	// Group updates by clob pair id.
	updates := make(map[uint32]*clobtypes.OffchainUpdates)
	for _, message := range offchainUpdates.Messages {
		clobPairId := message.OrderId.ClobPairId
		if _, ok := updates[clobPairId]; !ok {
			updates[clobPairId] = clobtypes.NewOffchainUpdates()
		}
		updates[clobPairId].Messages = append(updates[clobPairId].Messages, message)
	}

	// Unmarshal each per-clob pair message to v1 updates.
	updatesByClobPairId := make(map[uint32][]clobtypes.StreamUpdate)
	for clobPairId, update := range updates {
		v1updates, err := GetOffchainUpdatesV1(update)
		if err != nil {
			panic(err)
		}
		updatesByClobPairId[clobPairId] = []clobtypes.StreamUpdate{
			{
				UpdateMessage: &clobtypes.StreamUpdate_OrderbookUpdate{
					OrderbookUpdate: &clobtypes.StreamOrderbookUpdate{
						Updates:  v1updates,
						Snapshot: snapshot,
					},
				},
				BlockHeight: blockHeight,
				ExecMode:    uint32(execMode),
			},
		}
	}

	sm.sendStreamUpdate(
		updatesByClobPairId,
	)
}

// SendOrderbookFillUpdates groups fills by their clob pair ids and
// sends messages to the subscribers.
func (sm *GrpcStreamingManagerImpl) SendOrderbookFillUpdates(
	ctx sdk.Context,
	orderbookFills []clobtypes.StreamOrderbookFill,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendOrderbookFillsLatency,
		time.Now(),
	)

	// Group fills by clob pair id.
	updatesByClobPairId := make(map[uint32][]clobtypes.StreamUpdate)
	for _, orderbookFill := range orderbookFills {
		// Fetch the clob pair id from the first order in `OrderBookMatchFill`.
		// We can assume there must be an order, and that all orders share the same
		// clob pair id.
		clobPairId := orderbookFill.Orders[0].OrderId.ClobPairId
		if _, ok := updatesByClobPairId[clobPairId]; !ok {
			updatesByClobPairId[clobPairId] = []clobtypes.StreamUpdate{}
		}
		streamUpdate := clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_OrderFill{
				OrderFill: &orderbookFill,
			},
			BlockHeight: blockHeight,
			ExecMode:    uint32(execMode),
		}
		updatesByClobPairId[clobPairId] = append(updatesByClobPairId[clobPairId], streamUpdate)
	}

	sm.sendStreamUpdate(
		updatesByClobPairId,
	)
}

// sendStreamUpdate takes in a map of clob pair id to stream updates and emits them to subscribers.
func (sm *GrpcStreamingManagerImpl) sendStreamUpdate(
	updatesByClobPairId map[uint32][]clobtypes.StreamUpdate,
) {
	metrics.IncrCounter(
		metrics.GrpcEmitProtocolUpdateCount,
		1,
	)

	sm.Lock()
	defer sm.Unlock()

	// Send updates to subscribers.
	idsToRemove := make([]uint32, 0)
	for id, subscription := range sm.orderbookSubscriptions {
		streamUpdatesForSubscription := make([]clobtypes.StreamUpdate, 0)
		for _, clobPairId := range subscription.clobPairIds {
			if update, ok := updatesByClobPairId[clobPairId]; ok {
				streamUpdatesForSubscription = append(streamUpdatesForSubscription, update...)
			}
		}

		if len(streamUpdatesForSubscription) > 0 {
			metrics.IncrCounter(
				metrics.GrpcSendResponseToSubscriberCount,
				1,
			)
			if err := subscription.srv.Send(
				&clobtypes.StreamOrderbookUpdatesResponse{
					Updates: streamUpdatesForSubscription,
				},
			); err != nil {
				idsToRemove = append(idsToRemove, id)
			}
		}
	}

	// Clean up subscriptions that have been closed.
	// If a Send update has failed for any clob pair id, the whole subscription will be removed.
	for _, id := range idsToRemove {
		delete(sm.orderbookSubscriptions, id)
	}
	sm.EmitMetrics()
}

// GetUninitializedClobPairIds returns the clob pair ids that have not been initialized.
func (sm *GrpcStreamingManagerImpl) GetUninitializedClobPairIds() []uint32 {
	sm.Lock()
	defer sm.Unlock()

	clobPairIds := make(map[uint32]bool)
	for _, subscription := range sm.orderbookSubscriptions {
		subscription.initialize.Do(
			func() {
				for _, clobPairId := range subscription.clobPairIds {
					clobPairIds[clobPairId] = true
				}
			},
		)
	}

	return lib.GetSortedKeys[lib.Sortable[uint32]](clobPairIds)
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
