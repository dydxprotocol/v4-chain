package grpc

import (
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/grpc/client"
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

	// grpc stream will batch and flush out messages every 10 ms.
	ticker *time.Ticker
	done   chan bool
	// map of clob pair id to stream updates.
	streamUpdateCache map[uint32][]clobtypes.StreamUpdate
	numUpdatesInCache uint32

	maxUpdatesInCache          uint32
	maxSubscriptionChannelSize uint32
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

	// Testing
	client *client.GrpcClient
}

func NewGrpcStreamingManager(
	logger log.Logger,
	flushIntervalMs uint32,
	maxUpdatesInCache uint32,
	maxSubscriptionChannelSize uint32,
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

		maxUpdatesInCache:          maxUpdatesInCache,
		maxSubscriptionChannelSize: maxSubscriptionChannelSize,
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
		updatesChannel: make(chan []clobtypes.StreamUpdate, sm.maxSubscriptionChannelSize),
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
			// Break out of the loop, stopping this goroutine.
			// The channel will fill up and the main thread will prune the subscription.
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

func (sm *GrpcStreamingManagerImpl) SubscribeTestClient(client *client.GrpcClient) {
	subscription := &OrderbookSubscription{
		clobPairIds: []uint32{0, 1},
		client:      client,
	}

	sm.Lock()
	defer sm.Unlock()

	sm.orderbookSubscriptions[sm.nextSubscriptionId] = subscription
	sm.nextSubscriptionId++
}

// removeSubscription removes a subscription from the streaming manager.
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

// SendSnapshot sends messages to a particular subscriber without buffering.
// Note this method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *GrpcStreamingManagerImpl) SendSnapshot(
	offchainUpdates *clobtypes.OffchainUpdates,
	subscriptionId uint32,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendOrderbookSnapshotLatency,
		time.Now(),
	)

	v1updates, err := GetOffchainUpdatesV1(offchainUpdates)
	if err != nil {
		panic(err)
	}

	removeSubscription := false
	if len(v1updates) > 0 {
		subscription, ok := sm.orderbookSubscriptions[subscriptionId]
		if !ok {
			sm.logger.Error(
				fmt.Sprintf(
					"GRPC Streaming subscription id %+v not found. This should not happen.",
					subscriptionId,
				),
			)
			return
		}
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
		subscription.client.Update(
			&clobtypes.StreamOrderbookUpdatesResponse{
				Updates: streamUpdates,
			},
		)
	}

	// Clean up subscriptions that have been closed.
	// If a Send update has failed for any clob pair id, the whole subscription will be removed.
	if removeSubscription {
		sm.removeSubscription(subscriptionId)
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

	sm.Lock()
	defer sm.Unlock()

	for _, subscription := range sm.orderbookSubscriptions {
		updatesToSend := make([]clobtypes.StreamUpdate, 0)
		for _, clobPairId := range subscription.clobPairIds {
			if updates, ok := updatesByClobPairId[clobPairId]; ok {
				updatesToSend = append(updatesToSend, updates...)
			}
		}

		if len(updatesToSend) > 0 {
			subscription.client.Update(
				&clobtypes.StreamOrderbookUpdatesResponse{
					Updates: updatesToSend,
				},
			)
		}
	}
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

	sm.Lock()
	defer sm.Unlock()

	for _, subscription := range sm.orderbookSubscriptions {
		updatesToSend := make([]clobtypes.StreamUpdate, 0)
		for _, clobPairId := range subscription.clobPairIds {
			if updates, ok := updatesByClobPairId[clobPairId]; ok {
				updatesToSend = append(updatesToSend, updates...)
			}
		}

		if len(updatesToSend) > 0 {
			subscription.client.Update(
				&clobtypes.StreamOrderbookUpdatesResponse{
					Updates: updatesToSend,
				},
			)
		}
	}
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

func (sm *GrpcStreamingManagerImpl) FlushStreamUpdates() {
	sm.Lock()
	defer sm.Unlock()
	sm.FlushStreamUpdatesWithLock()
}

// FlushStreamUpdatesWithLock takes in a map of clob pair id to stream updates and emits them to subscribers.
// Note this method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *GrpcStreamingManagerImpl) FlushStreamUpdatesWithLock() {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcFlushUpdatesLatency,
		time.Now(),
	)

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

func (sm *GrpcStreamingManagerImpl) InitializeNewGrpcStreams(
	getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	sm.Lock()
	defer sm.Unlock()

	// Flush any pending updates before sending the snapshot to avoid
	// race conditions with the snapshot.
	sm.FlushStreamUpdatesWithLock()

	updatesByClobPairId := make(map[uint32]*clobtypes.OffchainUpdates)
	for subscriptionId, subscription := range sm.orderbookSubscriptions {
		subscription.initialize.Do(
			func() {
				allUpdates := clobtypes.NewOffchainUpdates()
				for _, clobPairId := range subscription.clobPairIds {
					if _, ok := updatesByClobPairId[clobPairId]; !ok {
						updatesByClobPairId[clobPairId] = getOrderbookSnapshot(clobtypes.ClobPairId(clobPairId))
					}
					allUpdates.Append(updatesByClobPairId[clobPairId])
				}

				sm.SendSnapshot(allUpdates, subscriptionId, blockHeight, execMode)
			},
		)
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
