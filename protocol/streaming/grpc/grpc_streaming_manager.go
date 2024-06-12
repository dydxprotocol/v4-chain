package grpc

import (
	"fmt"
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

const (
	// TODO CHANGE THIS CONFIGURABLE
	MAX_BUFF_LEN = 1000
)

// GrpcStreamingManagerImpl is an implementation for managing gRPC streaming subscriptions.
type GrpcStreamingManagerImpl struct {
	sync.Mutex

	logger log.Logger

	// orderbookSubscriptions maps subscription IDs to their respective orderbook subscriptions.
	orderbookSubscriptions map[uint32]*OrderbookSubscription
	nextSubscriptionId     uint32

	// grpc stream will batch and flush out messages every 10 ms.
	ticker *time.Ticker
	done   chan bool
	// map of clob pair id to stream updates.
	streamUpdateCache map[uint32][]clobtypes.StreamUpdate
	numUpdatesInCache uint32

	maxUpdatesInCache uint32
}

// OrderbookSubscription represents a active subscription to the orderbook updates stream.
type OrderbookSubscription struct {
	subscriptionId uint32

	// Initialize the subscription with orderbook snapshots.
	initialize sync.Once

	// Clob pair ids to subscribe to.
	clobPairIds []uint32

	// Stream
	srv clobtypes.Query_StreamOrderbookUpdatesServer

	// Channel to buffer writes before the stream
	updatesChannel chan []clobtypes.StreamUpdate
}

func NewGrpcStreamingManager(
	logger log.Logger,
	flushIntervalMs uint32,
	maxUpdatesInCache uint32,
) *GrpcStreamingManagerImpl {
	logger = logger.With(log.ModuleKey, "grpc-streaming")
	grpcStreamingManager := &GrpcStreamingManagerImpl{
		logger:                 logger,
		orderbookSubscriptions: make(map[uint32]*OrderbookSubscription),
		nextSubscriptionId:     0,

		ticker:            time.NewTicker(time.Duration(flushIntervalMs) * time.Millisecond),
		done:              make(chan bool),
		streamUpdateCache: make(map[uint32][]clobtypes.StreamUpdate),
		numUpdatesInCache: 0,

		maxUpdatesInCache: maxUpdatesInCache,
	}

	// Start the goroutine for pushing order updates through.
	// Sender goroutine for the subscription channels.
	go func() {
		for {
			select {
			case <-grpcStreamingManager.ticker.C:
				grpcStreamingManager.FlushStreamUpdates()
			case <-grpcStreamingManager.done:
				grpcStreamingManager.logger.Info(
					"GRPC Stream poller goroutine shutting down",
				)
				return
			}
		}
	}()

	return grpcStreamingManager
}

func (sm *GrpcStreamingManagerImpl) Enabled() bool {
	return true
}

func (sm *GrpcStreamingManagerImpl) EmitMetrics() {
	metrics.SetGauge(
		metrics.GrpcStreamNumUpdatesBuffered,
		float32(sm.numUpdatesInCache),
	)
	metrics.SetGauge(
		metrics.GrpcStreamSubscriberCount,
		float32(len(sm.orderbookSubscriptions)),
	)
	for _, subscription := range sm.orderbookSubscriptions {
		metrics.AddSample(
			metrics.GrpcSubscriptionChannelLength,
			float32(len(subscription.updatesChannel)),
		)
	}
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

	sm.Lock()
	subscription := &OrderbookSubscription{
		subscriptionId: sm.nextSubscriptionId,
		clobPairIds:    clobPairIds,
		srv:            srv,
		updatesChannel: make(chan []clobtypes.StreamUpdate, MAX_BUFF_LEN),
	}

	sm.logger.Info(
		fmt.Sprintf(
			"New subscription id %+v for clob pair ids: %+v",
			subscription.subscriptionId,
			clobPairIds,
		),
	)
	sm.orderbookSubscriptions[subscription.subscriptionId] = subscription
	sm.nextSubscriptionId++
	sm.EmitMetrics()
	sm.Unlock()

	// Use current goroutine to consistently poll subscription channel for updates
	// to send through stream.
	for updates := range subscription.updatesChannel {
		metrics.IncrCounter(
			metrics.GrpcSendResponseToSubscriberCount,
			1,
		)
		err = subscription.srv.Send(
			&clobtypes.StreamOrderbookUpdatesResponse{
				Updates: updates,
			},
		)
		if err != nil {
			// On error, remove the subscription from the streaming manager
			sm.logger.Error(
				fmt.Sprintf(
					"Error sending out update for grpc streaming subscription %+v. Dropping subsciption connection.",
					subscription.subscriptionId,
				),
				"err", err,
			)
			delete(sm.orderbookSubscriptions, subscription.subscriptionId)
			break
		}
	}

	sm.logger.Info(
		fmt.Sprintf(
			"Terminating poller for subscription id %+v",
			subscription.subscriptionId,
		),
	)
	return err
}

// removeSubscription removes a subscription from the grpc streaming manager.
// The streaming manager's lock should already be acquired before calling this.
func (sm *GrpcStreamingManagerImpl) removeSubscription(
	subscriptionIdToRemove uint32,
) {
	subscription := sm.orderbookSubscriptions[subscriptionIdToRemove]
	if subscription == nil {
		return
	}
	close(subscription.updatesChannel)
	delete(sm.orderbookSubscriptions, subscriptionIdToRemove)
	sm.logger.Info(
		fmt.Sprintf("Removed grpc streaming subscription id %+v", subscriptionIdToRemove),
	)
}

func (sm *GrpcStreamingManagerImpl) Stop() {
	sm.done <- true
}

// SendSnapshot groups updates by their clob pair ids and
// sends messages to the subscribers. It groups out updates differently
// and bypasses the buffer.
func (sm *GrpcStreamingManagerImpl) SendSnapshot(
	offchainUpdates *clobtypes.OffchainUpdates,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendOrderbookSnapshotLatency,
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
	updatesByClobPairId := make(map[uint32][]ocutypes.OffChainUpdateV1)
	for clobPairId, update := range updates {
		v1updates, err := GetOffchainUpdatesV1(update)
		if err != nil {
			panic(err)
		}
		updatesByClobPairId[clobPairId] = v1updates
	}

	sm.Lock()
	defer sm.Unlock()

	idsToRemove := make([]uint32, 0)
	for id, subscription := range sm.orderbookSubscriptions {
		// Consolidate orderbook updates into a single `StreamUpdate`.
		v1updates := make([]ocutypes.OffChainUpdateV1, 0)
		for _, clobPairId := range subscription.clobPairIds {
			if update, ok := updatesByClobPairId[clobPairId]; ok {
				v1updates = append(v1updates, update...)
			}
		}

		if len(v1updates) > 0 {
			streamUpdates := []clobtypes.StreamUpdate{
				{
					UpdateMessage: &clobtypes.StreamUpdate_OrderbookUpdate{
						OrderbookUpdate: &clobtypes.StreamOrderbookUpdate{
							Updates:  v1updates,
							Snapshot: true,
						},
					},
					BlockHeight: blockHeight,
					ExecMode:    uint32(execMode),
				},
			}
			metrics.IncrCounter(
				metrics.GrpcAddToSubscriptionChannelCount,
				1,
			)
			select {
			case subscription.updatesChannel <- streamUpdates:
			default:
				sm.logger.Error(
					fmt.Sprintf(
						"GRPC Streaming subscription id %+v channel full capacity. Dropping subscription connection.",
						id,
					),
				)
				idsToRemove = append(idsToRemove, subscription.subscriptionId)
			}
		}
	}

	// Clean up subscriptions that have been closed.
	// If a Send update has failed for any clob pair id, the whole subscription will be removed.
	for _, id := range idsToRemove {
		sm.removeSubscription(id)
	}
}

// SendOrderbookUpdates groups updates by their clob pair ids and
// sends messages to the subscribers.
func (sm *GrpcStreamingManagerImpl) SendOrderbookUpdates(
	offchainUpdates *clobtypes.OffchainUpdates,
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
						Snapshot: false,
					},
				},
				BlockHeight: blockHeight,
				ExecMode:    uint32(execMode),
			},
		}
	}

	sm.AddUpdatesToCache(updatesByClobPairId, uint32(len(updates)))
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

	sm.AddUpdatesToCache(updatesByClobPairId, uint32(len(orderbookFills)))
}

func (sm *GrpcStreamingManagerImpl) AddUpdatesToCache(
	updatesByClobPairId map[uint32][]clobtypes.StreamUpdate,
	numUpdatesToAdd uint32,
) {
	sm.Lock()
	defer sm.Unlock()

	metrics.IncrCounter(
		metrics.GrpcAddUpdateToBufferCount,
		1,
	)

	for clobPairId, streamUpdates := range updatesByClobPairId {
		sm.streamUpdateCache[clobPairId] = append(sm.streamUpdateCache[clobPairId], streamUpdates...)
	}
	sm.numUpdatesInCache += numUpdatesToAdd

	// Remove all subscriptions and wipe the buffer if buffer overflows.
	if sm.numUpdatesInCache > sm.maxUpdatesInCache {
		sm.logger.Error("GRPC Streaming buffer full capacity. Dropping messages and all subscriptions. " +
			"Disconnect all clients and increase buffer size via the grpc-stream-buffer-size flag.")
		for id := range sm.orderbookSubscriptions {
			sm.removeSubscription(id)
		}
		clear(sm.streamUpdateCache)
		sm.numUpdatesInCache = 0
	}
	sm.EmitMetrics()
}

// FlushStreamUpdates takes in a map of clob pair id to stream updates and emits them to subscribers.
func (sm *GrpcStreamingManagerImpl) FlushStreamUpdates() {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcFlushUpdatesLatency,
		time.Now(),
	)

	sm.Lock()
	defer sm.Unlock()

	// Non-blocking send updates through subscriber's buffered channel.
	// If the buffer is full, drop the subscription.
	idsToRemove := make([]uint32, 0)
	for id, subscription := range sm.orderbookSubscriptions {
		streamUpdatesForSubscription := make([]clobtypes.StreamUpdate, 0)
		for _, clobPairId := range subscription.clobPairIds {
			if update, ok := sm.streamUpdateCache[clobPairId]; ok {
				streamUpdatesForSubscription = append(streamUpdatesForSubscription, update...)
			}
		}

		if len(streamUpdatesForSubscription) > 0 {
			metrics.IncrCounter(
				metrics.GrpcAddToSubscriptionChannelCount,
				1,
			)
			select {
			case subscription.updatesChannel <- streamUpdatesForSubscription:
			default:
				idsToRemove = append(idsToRemove, id)
			}
		}
	}

	clear(sm.streamUpdateCache)
	sm.numUpdatesInCache = 0

	for _, id := range idsToRemove {
		sm.logger.Error(
			fmt.Sprintf(
				"GRPC Streaming subscription id %+v channel full capacity. Dropping subscription connection.",
				id,
			),
		)
		sm.removeSubscription(id)
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
