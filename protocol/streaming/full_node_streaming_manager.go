package streaming

import (
	"fmt"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"sync"
	"time"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	streaming_util "github.com/dydxprotocol/v4-chain/protocol/streaming/util"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.FullNodeStreamingManager = (*FullNodeStreamingManagerImpl)(nil)

// FullNodeStreamingManagerImpl is an implementation for managing streaming subscriptions.
type FullNodeStreamingManagerImpl struct {
	sync.Mutex

	logger log.Logger

	// orderbookSubscriptions maps subscription IDs to their respective orderbook subscriptions.
	orderbookSubscriptions      map[uint32]*OrderbookSubscription
	nextOrderbookSubscriptionId uint32

	// subaccountSubscriptions maps subscription IDs to their respective subaccount subscriptions.
	subaccountSubscriptions      map[uint32]*SubaccountSubscription
	nextSubaccountSubscriptionId uint32

	// stream will batch and flush out messages every 10 ms.
	ticker *time.Ticker
	done   chan bool

	// map of clob pair id to stream updates.
	orderbookStreamUpdateCache map[uint32][]clobtypes.StreamUpdate
	numOrderbookUpdatesInCache uint32

	// map of subaccount id to stream updates.
	subaccountStreamUpdateCache map[satypes.SubaccountId][]*satypes.StreamSubaccountUpdate
	numSubaccountUpdatesInCache uint32

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
	messageSender types.OutgoingOrderbookMessageSender

	// Channel to buffer writes before the stream
	updatesChannel chan []clobtypes.StreamUpdate
}

// SubaccountSubscription represents a active subscription to the subaccount updates stream.
type SubaccountSubscription struct {
	subscriptionId uint32

	// Initialize the subscription with subaccount snapshots.
	initialize sync.Once

	// Subaccount ids to subscribe to.
	subaccountIds []*satypes.SubaccountId

	// Stream
	messageSender types.OutgoingSubaccountMessageSender

	// Channel to buffer writes before the stream
	updatesChannel chan []*satypes.StreamSubaccountUpdate
}

func NewFullNodeStreamingManager(
	logger log.Logger,
	flushIntervalMs uint32,
	maxUpdatesInCache uint32,
	maxSubscriptionChannelSize uint32,
) *FullNodeStreamingManagerImpl {
	logger = logger.With(log.ModuleKey, "full-node-streaming")
	fullNodeStreamingManager := &FullNodeStreamingManagerImpl{
		logger:                      logger,
		orderbookSubscriptions:      make(map[uint32]*OrderbookSubscription),
		nextOrderbookSubscriptionId: 0,

		subaccountSubscriptions:      make(map[uint32]*SubaccountSubscription),
		nextSubaccountSubscriptionId: 0,

		ticker:                     time.NewTicker(time.Duration(flushIntervalMs) * time.Millisecond),
		done:                       make(chan bool),
		orderbookStreamUpdateCache: make(map[uint32][]clobtypes.StreamUpdate),
		numOrderbookUpdatesInCache: 0,

		subaccountStreamUpdateCache: make(map[satypes.SubaccountId][]*satypes.StreamSubaccountUpdate),
		numSubaccountUpdatesInCache: 0,

		maxUpdatesInCache:          maxUpdatesInCache,
		maxSubscriptionChannelSize: maxSubscriptionChannelSize,
	}

	// Start the goroutine for pushing order updates through.
	// Sender goroutine for the subscription channels.
	go func() {
		for {
			select {
			case <-fullNodeStreamingManager.ticker.C:
				fullNodeStreamingManager.FlushStreamUpdates()
			case <-fullNodeStreamingManager.done:
				fullNodeStreamingManager.logger.Info(
					"Stream poller goroutine shutting down",
				)
				return
			}
		}
	}()

	return fullNodeStreamingManager
}

func (sm *FullNodeStreamingManagerImpl) Enabled() bool {
	return true
}

func (sm *FullNodeStreamingManagerImpl) EmitMetrics() {
	metrics.SetGauge(
		metrics.GrpcStreamNumUpdatesBuffered,
		float32(sm.numOrderbookUpdatesInCache),
	)
	metrics.SetGauge(
		metrics.GrpcStreamSubscriberCount,
		float32(len(sm.orderbookSubscriptions)),
	)
	for _, subscription := range sm.orderbookSubscriptions {
		metrics.AddSample(
			metrics.GrpcOrderbookSubscriptionChannelLength,
			float32(len(subscription.updatesChannel)),
		)
	}
	for _, subscription := range sm.subaccountSubscriptions {
		metrics.AddSample(
			metrics.GrpcSubaccountSubscriptionChannelLength,
			float32(len(subscription.updatesChannel)),
		)
	}
}

// Subscribe subscribes to the orderbook updates stream.
func (sm *FullNodeStreamingManagerImpl) SubscribeToOrderbookStream(
	clobPairIds []uint32,
	messageSender types.OutgoingOrderbookMessageSender,
) (
	err error,
) {
	// Perform some basic validation on the request.
	if len(clobPairIds) == 0 {
		return types.ErrInvalidStreamingRequest
	}

	sm.Lock()
	subscription := &OrderbookSubscription{
		subscriptionId: sm.nextOrderbookSubscriptionId,
		clobPairIds:    clobPairIds,
		messageSender:  messageSender,
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
	sm.nextOrderbookSubscriptionId++
	sm.EmitMetrics()
	sm.Unlock()

	// Use current goroutine to consistently poll subscription channel for updates
	// to send through stream.
	for updates := range subscription.updatesChannel {
		metrics.IncrCounter(
			metrics.GrpcSendResponseToSubscriberCount,
			1,
		)
		err = subscription.messageSender.Send(
			&clobtypes.StreamOrderbookUpdatesResponse{
				Updates: updates,
			},
		)
		if err != nil {
			// On error, remove the subscription from the streaming manager
			sm.logger.Error(
				fmt.Sprintf(
					"Error sending out update for streaming subscription %+v. Dropping subsciption connection.",
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

// Subscribe subscribes to the subaccount updates stream.
func (sm *FullNodeStreamingManagerImpl) SubscribeToSubaccountStream(
	subaccountIds []*satypes.SubaccountId,
	messageSender types.OutgoingSubaccountMessageSender,
) (
	err error,
) {
	return types.ErrNotImplemented
}

// removeOrderbookSubscription removes a subscription from the streaming manager.
// The streaming manager's lock should already be acquired before calling this.
func (sm *FullNodeStreamingManagerImpl) removeOrderbookSubscription(
	subscriptionIdToRemove uint32,
) {
	subscription := sm.orderbookSubscriptions[subscriptionIdToRemove]
	if subscription == nil {
		return
	}
	close(subscription.updatesChannel)
	delete(sm.orderbookSubscriptions, subscriptionIdToRemove)
	sm.logger.Info(
		fmt.Sprintf("Removed streaming subscription id %+v", subscriptionIdToRemove),
	)
}

// removeSubaccountSubscription removes a subaccount subscription from the streaming manager.
// The streaming manager's lock should already be acquired before calling this.
func (sm *FullNodeStreamingManagerImpl) removeSubaccountSubscription(subscriptionIdToRemove uint32) {
	subscription := sm.subaccountSubscriptions[subscriptionIdToRemove]
	if subscription == nil {
		return
	}
	close(subscription.updatesChannel)
	delete(sm.subaccountSubscriptions, subscriptionIdToRemove)
	sm.logger.Info(
		fmt.Sprintf("Removed streaming subaccount subscription id %+v", subscriptionIdToRemove),
	)
}

func (sm *FullNodeStreamingManagerImpl) Stop() {
	sm.done <- true
}

// SendSnapshot sends messages to a particular subscriber without buffering.
// Note this method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) SendSnapshot(
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

	v1updates, err := streaming_util.GetOffchainUpdatesV1(offchainUpdates)
	if err != nil {
		panic(err)
	}

	removeSubscription := false
	if len(v1updates) > 0 {
		subscription, ok := sm.orderbookSubscriptions[subscriptionId]
		if !ok {
			sm.logger.Error(
				fmt.Sprintf(
					"Streaming subscription id %+v not found. This should not happen.",
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
		select {
		case subscription.updatesChannel <- streamUpdates:
		default:
			sm.logger.Error(
				fmt.Sprintf(
					"Streaming subscription id %+v channel full capacity. Dropping subscription connection.",
					subscriptionId,
				),
			)
			removeSubscription = true
		}
	}

	// Clean up subscriptions that have been closed.
	// If a Send update has failed for any clob pair id, the whole subscription will be removed.
	if removeSubscription {
		sm.removeOrderbookSubscription(subscriptionId)
	}
}

// SendOrderbookUpdates groups updates by their clob pair ids and
// sends messages to the subscribers.
func (sm *FullNodeStreamingManagerImpl) SendOrderbookUpdates(
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
		v1updates, err := streaming_util.GetOffchainUpdatesV1(update)
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

	sm.AddOrderbookUpdatesToCache(updatesByClobPairId, uint32(len(updates)))
}

// SendOrderbookFillUpdates groups fills by their clob pair ids and
// sends messages to the subscribers.
func (sm *FullNodeStreamingManagerImpl) SendOrderbookFillUpdates(
	orderbookFills []clobtypes.StreamOrderbookFill,
	blockHeight uint32,
	execMode sdk.ExecMode,
	perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendOrderbookFillsLatency,
		time.Now(),
	)

	// Group fills by clob pair id.
	updatesByClobPairId := make(map[uint32][]clobtypes.StreamUpdate)
	for _, orderbookFill := range orderbookFills {
		// If this is a deleveraging fill, fetch the clob pair id from the deleveraged
		// perpetual id.
		// Otherwise, fetch the clob pair id from the first order in `OrderBookMatchFill`.
		// We can assume there must be an order, and that all orders share the same
		// clob pair id.
		clobPairId := uint32(0)
		if match := orderbookFill.GetClobMatch().GetMatchPerpetualDeleveraging(); match != nil {
			clobPairId = uint32(perpetualIdToClobPairId[match.PerpetualId][0])
		} else {
			clobPairId = orderbookFill.Orders[0].OrderId.ClobPairId
		}
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

	sm.AddOrderbookUpdatesToCache(updatesByClobPairId, uint32(len(orderbookFills)))
}

func (sm *FullNodeStreamingManagerImpl) AddOrderbookUpdatesToCache(
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
		sm.orderbookStreamUpdateCache[clobPairId] = append(sm.orderbookStreamUpdateCache[clobPairId], streamUpdates...)
	}
	sm.numOrderbookUpdatesInCache += numUpdatesToAdd

	// Remove all subscriptions and wipe the buffer if buffer overflows.
	if sm.numOrderbookUpdatesInCache > sm.maxUpdatesInCache {
		sm.logger.Error("Streaming buffer full capacity. Dropping messages and all subscriptions. " +
			"Disconnect all clients and increase buffer size via the grpc-stream-buffer-size flag.")
		for id := range sm.orderbookSubscriptions {
			sm.removeOrderbookSubscription(id)
		}
		clear(sm.orderbookStreamUpdateCache)
		sm.numOrderbookUpdatesInCache = 0
	}
	sm.EmitMetrics()
}

// AddSubaccountUpdatesToCache adds subaccount updates to the cache.
// If the cache exceeds its maximum size, all subscriptions are removed and the cache is cleared.
func (sm *FullNodeStreamingManagerImpl) AddSubaccountUpdatesToCache(
	updatesBySubaccountId map[satypes.SubaccountId][]*satypes.StreamSubaccountUpdate,
	numUpdatesToAdd uint32,
) {
	sm.Lock()
	defer sm.Unlock()

	metrics.IncrCounter(
		metrics.GrpcAddUpdateToBufferCount,
		1,
	)

	for subaccountId, streamUpdates := range updatesBySubaccountId {
		sm.subaccountStreamUpdateCache[subaccountId] = append(sm.subaccountStreamUpdateCache[subaccountId], streamUpdates...)
	}
	sm.numSubaccountUpdatesInCache += numUpdatesToAdd

	// Remove all subscriptions and wipe the buffer if buffer overflows.
	if sm.numSubaccountUpdatesInCache > sm.maxUpdatesInCache {
		sm.logger.Error("Streaming subaccount buffer full capacity. Dropping messages and all subscriptions. " +
			"Disconnect all clients and increase buffer size via the grpc-stream-buffer-size flag.")
		for id := range sm.subaccountSubscriptions {
			sm.removeSubaccountSubscription(id)
		}
		clear(sm.subaccountStreamUpdateCache)
		sm.numSubaccountUpdatesInCache = 0
	}
	sm.EmitMetrics()
}

func (sm *FullNodeStreamingManagerImpl) FlushStreamUpdates() {
	sm.Lock()
	defer sm.Unlock()
	sm.FlushStreamUpdatesWithLock()
}

// FlushStreamUpdatesWithLock takes in a map of clob pair id to stream updates and emits them to subscribers.
// Note this method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) FlushStreamUpdatesWithLock() {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcFlushUpdatesLatency,
		time.Now(),
	)

	// Non-blocking send updates through subscriber's buffered channel.
	// If the buffer is full, drop the subscription.
	orderbookIdsToRemove := make([]uint32, 0)
	for id, subscription := range sm.orderbookSubscriptions {
		streamUpdatesForSubscription := make([]clobtypes.StreamUpdate, 0)
		for _, clobPairId := range subscription.clobPairIds {
			if update, ok := sm.orderbookStreamUpdateCache[clobPairId]; ok {
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
				orderbookIdsToRemove = append(orderbookIdsToRemove, id)
			}
		}
	}

	// Clear the orderbook stream update cache and reset the count
	clear(sm.orderbookStreamUpdateCache)
	sm.numOrderbookUpdatesInCache = 0

	for _, id := range orderbookIdsToRemove {
		sm.logger.Error(
			fmt.Sprintf(
				"Streaming subscription id %+v channel full capacity. Dropping subscription connection.",
				id,
			),
		)
		sm.removeOrderbookSubscription(id)
	}

	// Handling subaccount subscriptions
	subaccountIdsToRemove := make([]uint32, 0)
	for id, subscription := range sm.subaccountSubscriptions {
		streamUpdatesForSubscription := make([]*satypes.StreamSubaccountUpdate, 0)
		for _, subaccountId := range subscription.subaccountIds {
			if update, ok := sm.subaccountStreamUpdateCache[*subaccountId]; ok {
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
				subaccountIdsToRemove = append(subaccountIdsToRemove, id)
			}
		}
	}

	// Clear the subaccount stream update cache and reset the count
	clear(sm.subaccountStreamUpdateCache)
	sm.numSubaccountUpdatesInCache = 0

	for _, id := range subaccountIdsToRemove {
		sm.logger.Error(
			fmt.Sprintf(
				"Streaming subaccount subscription id %+v channel full capacity. Dropping subscription connection.",
				id,
			),
		)
		sm.removeSubaccountSubscription(id)
	}

	sm.EmitMetrics()
}

func (sm *FullNodeStreamingManagerImpl) InitializeNewStreams(
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
