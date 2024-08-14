package streaming

import (
	"fmt"
	"sync"
	"time"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

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
	orderbookSubscriptions map[uint32]*OrderbookSubscription
	nextSubscriptionId     uint32

	// stream will batch and flush out messages every 10 ms.
	ticker *time.Ticker
	done   chan bool

	// TODO: Consolidate the streamUpdateCache and streamUpdateSubscriptionCache into a single
	// struct to avoid the need to maintain two separate slices for the same data.

	// list of stream updates.
	streamUpdateCache []clobtypes.StreamUpdate
	// list of subscription ids for each stream update.
	streamUpdateSubscriptionCache [][]uint32
	// map from clob pair id to subscription ids.
	clobPairIdToSubscriptionIdMapping map[uint32][]uint32
	// map from subaccount id to subscription ids.
	subaccountIdToSubscriptionIdMapping map[satypes.SubaccountId][]uint32

	maxUpdatesInCache          uint32
	maxSubscriptionChannelSize uint32

	// Block interval in which snapshot info should be sent out in.
	// Defaults to 0, which means only one snapshot will be sent out.
	snapshotBlockInterval uint32
}

// OrderbookSubscription represents a active subscription to the orderbook updates stream.
type OrderbookSubscription struct {
	subscriptionId uint32

	// Initialize the subscription with orderbook snapshots.
	initialize *sync.Once

	// Clob pair ids to subscribe to.
	clobPairIds []uint32

	// Subaccount ids to subscribe to.
	subaccountIds []satypes.SubaccountId

	// Stream
	messageSender types.OutgoingMessageSender

	// Channel to buffer writes before the stream
	updatesChannel chan []clobtypes.StreamUpdate

	// If interval snapshots are turned on, the next block height at which
	// a snapshot should be sent out.
	nextSnapshotBlock uint32
}

func NewFullNodeStreamingManager(
	logger log.Logger,
	flushIntervalMs uint32,
	maxUpdatesInCache uint32,
	maxSubscriptionChannelSize uint32,
	snapshotBlockInterval uint32,
) *FullNodeStreamingManagerImpl {
	logger = logger.With(log.ModuleKey, "full-node-streaming")
	fullNodeStreamingManager := &FullNodeStreamingManagerImpl{
		logger:                 logger,
		orderbookSubscriptions: make(map[uint32]*OrderbookSubscription),
		nextSubscriptionId:     0,

		ticker:                              time.NewTicker(time.Duration(flushIntervalMs) * time.Millisecond),
		done:                                make(chan bool),
		streamUpdateCache:                   make([]clobtypes.StreamUpdate, 0),
		streamUpdateSubscriptionCache:       make([][]uint32, 0),
		clobPairIdToSubscriptionIdMapping:   make(map[uint32][]uint32),
		subaccountIdToSubscriptionIdMapping: make(map[satypes.SubaccountId][]uint32),

		maxUpdatesInCache:          maxUpdatesInCache,
		maxSubscriptionChannelSize: maxSubscriptionChannelSize,
		snapshotBlockInterval:      snapshotBlockInterval,
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
		float32(len(sm.streamUpdateCache)),
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
func (sm *FullNodeStreamingManagerImpl) Subscribe(
	clobPairIds []uint32,
	subaccountIds []*satypes.SubaccountId,
	messageSender types.OutgoingMessageSender,
) (
	err error,
) {
	// Perform some basic validation on the request.
	if len(clobPairIds) == 0 {
		return types.ErrInvalidStreamingRequest
	}

	sm.Lock()
	sIds := make([]satypes.SubaccountId, len(subaccountIds))
	for i, subaccountId := range subaccountIds {
		sIds[i] = *subaccountId
	}
	subscription := &OrderbookSubscription{
		subscriptionId: sm.nextSubscriptionId,
		initialize:     &sync.Once{},
		clobPairIds:    clobPairIds,
		subaccountIds:  sIds,
		messageSender:  messageSender,
		updatesChannel: make(chan []clobtypes.StreamUpdate, sm.maxSubscriptionChannelSize),
	}
	for _, clobPairId := range clobPairIds {
		// if clobPairId exists in the map, append the subscription id to the slice
		// otherwise, create a new slice with the subscription id
		if _, ok := sm.clobPairIdToSubscriptionIdMapping[clobPairId]; !ok {
			sm.clobPairIdToSubscriptionIdMapping[clobPairId] = []uint32{}
		}
		sm.clobPairIdToSubscriptionIdMapping[clobPairId] = append(
			sm.clobPairIdToSubscriptionIdMapping[clobPairId],
			sm.nextSubscriptionId,
		)
	}
	for _, subaccountId := range sIds {
		// if subaccountId exists in the map, append the subscription id to the slice
		// otherwise, create a new slice with the subscription id
		if _, ok := sm.subaccountIdToSubscriptionIdMapping[subaccountId]; !ok {
			sm.subaccountIdToSubscriptionIdMapping[subaccountId] = []uint32{}
		}
		sm.subaccountIdToSubscriptionIdMapping[subaccountId] = append(
			sm.subaccountIdToSubscriptionIdMapping[subaccountId],
			sm.nextSubscriptionId,
		)
	}

	sm.logger.Info(
		fmt.Sprintf(
			"New subscription id %+v for clob pair ids: %+v and subaccount ids: %+v",
			subscription.subscriptionId,
			clobPairIds,
			subaccountIds,
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

// removeSubscription removes a subscription from the streaming manager.
// The streaming manager's lock should already be acquired before calling this.
func (sm *FullNodeStreamingManagerImpl) removeSubscription(
	subscriptionIdToRemove uint32,
) {
	subscription := sm.orderbookSubscriptions[subscriptionIdToRemove]
	if subscription == nil {
		return
	}
	close(subscription.updatesChannel)
	delete(sm.orderbookSubscriptions, subscriptionIdToRemove)

	// Iterate over the clobPairIdToSubscriptionIdMapping to remove the subscriptionIdToRemove
	for pairId, subscriptionIds := range sm.clobPairIdToSubscriptionIdMapping {
		for i, id := range subscriptionIds {
			if id == subscriptionIdToRemove {
				// Remove the subscription ID from the slice
				sm.clobPairIdToSubscriptionIdMapping[pairId] = append(subscriptionIds[:i], subscriptionIds[i+1:]...)
				break
			}
		}
		// If the list is empty after removal, delete the key from the map
		if len(sm.clobPairIdToSubscriptionIdMapping[pairId]) == 0 {
			delete(sm.clobPairIdToSubscriptionIdMapping, pairId)
		}
	}

	// Iterate over the subaccountIdToSubscriptionIdMapping to remove the subscriptionIdToRemove
	for subaccountId, subscriptionIds := range sm.subaccountIdToSubscriptionIdMapping {
		for i, id := range subscriptionIds {
			if id == subscriptionIdToRemove {
				// Remove the subscription ID from the slice
				sm.subaccountIdToSubscriptionIdMapping[subaccountId] = append(subscriptionIds[:i], subscriptionIds[i+1:]...)
				break
			}
		}
		// If the list is empty after removal, delete the key from the map
		if len(sm.subaccountIdToSubscriptionIdMapping[subaccountId]) == 0 {
			delete(sm.subaccountIdToSubscriptionIdMapping, subaccountId)
		}
	}

	sm.logger.Info(
		fmt.Sprintf("Removed streaming subscription id %+v", subscriptionIdToRemove),
	)
}

func (sm *FullNodeStreamingManagerImpl) Stop() {
	sm.done <- true
}

func toOrderbookStreamUpdate(
	offchainUpdates *clobtypes.OffchainUpdates,
	blockHeight uint32,
	execMode sdk.ExecMode,
) []clobtypes.StreamUpdate {
	v1updates, err := streaming_util.GetOffchainUpdatesV1(offchainUpdates)
	if err != nil {
		panic(err)
	}
	return []clobtypes.StreamUpdate{
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
}

func toSubaccountStreamUpdates(
	saUpdates []*satypes.StreamSubaccountUpdate,
	blockHeight uint32,
	execMode sdk.ExecMode,
) []clobtypes.StreamUpdate {
	streamUpdates := make([]clobtypes.StreamUpdate, 0)
	for _, saUpdate := range saUpdates {
		streamUpdates = append(streamUpdates, clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_SubaccountUpdate{
				SubaccountUpdate: saUpdate,
			},
			BlockHeight: blockHeight,
			ExecMode:    uint32(execMode),
		})
	}
	return streamUpdates
}

func (sm *FullNodeStreamingManagerImpl) sendStreamUpdates(
	subscriptionId uint32,
	streamUpdates []clobtypes.StreamUpdate,
) {
	removeSubscription := false
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

	if removeSubscription {
		sm.removeSubscription(subscriptionId)
	}
}

// SendCombinedSnapshot sends messages to a particular subscriber without buffering.
// Note this method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) SendCombinedSnapshot(
	offchainUpdates *clobtypes.OffchainUpdates,
	saUpdates []*satypes.StreamSubaccountUpdate,
	subscriptionId uint32,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendOrderbookSnapshotLatency,
		time.Now(),
	)

	var streamUpdates []clobtypes.StreamUpdate
	streamUpdates = append(streamUpdates, toOrderbookStreamUpdate(offchainUpdates, blockHeight, execMode)...)
	streamUpdates = append(streamUpdates, toSubaccountStreamUpdates(saUpdates, blockHeight, execMode)...)
	sm.sendStreamUpdates(subscriptionId, streamUpdates)
}

// TracksSubaccountId checks if a subaccount id is being tracked by the streaming manager.
func (sm *FullNodeStreamingManagerImpl) TracksSubaccountId(subaccountId satypes.SubaccountId) bool {
	sm.Lock()
	defer sm.Unlock()
	_, exists := sm.subaccountIdToSubscriptionIdMapping[subaccountId]
	return exists
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
	streamUpdates := make([]clobtypes.StreamUpdate, 0)
	clobPairIds := make([]uint32, 0)
	for clobPairId, update := range updates {
		v1updates, err := streaming_util.GetOffchainUpdatesV1(update)
		if err != nil {
			panic(err)
		}
		streamUpdate := clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_OrderbookUpdate{
				OrderbookUpdate: &clobtypes.StreamOrderbookUpdate{
					Updates:  v1updates,
					Snapshot: false,
				},
			},
			BlockHeight: blockHeight,
			ExecMode:    uint32(execMode),
		}
		streamUpdates = append(streamUpdates, streamUpdate)
		clobPairIds = append(clobPairIds, clobPairId)
	}

	sm.AddOrderUpdatesToCache(streamUpdates, clobPairIds)
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
	streamUpdates := make([]clobtypes.StreamUpdate, 0)
	clobPairIds := make([]uint32, 0)
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
		streamUpdate := clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_OrderFill{
				OrderFill: &orderbookFill,
			},
			BlockHeight: blockHeight,
			ExecMode:    uint32(execMode),
		}
		streamUpdates = append(streamUpdates, streamUpdate)
		clobPairIds = append(clobPairIds, clobPairId)
	}

	sm.AddOrderUpdatesToCache(streamUpdates, clobPairIds)
}

// SendTakerOrderStatus sends out a taker order and its status to the full node streaming service.
func (sm *FullNodeStreamingManagerImpl) SendTakerOrderStatus(
	streamTakerOrder clobtypes.StreamTakerOrder,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	clobPairId := uint32(0)
	if liqOrder := streamTakerOrder.GetLiquidationOrder(); liqOrder != nil {
		clobPairId = liqOrder.ClobPairId
	}
	if takerOrder := streamTakerOrder.GetOrder(); takerOrder != nil {
		clobPairId = takerOrder.OrderId.ClobPairId
	}

	sm.AddOrderUpdatesToCache(
		[]clobtypes.StreamUpdate{
			{
				UpdateMessage: &clobtypes.StreamUpdate_TakerOrder{
					TakerOrder: &streamTakerOrder,
				},
				BlockHeight: blockHeight,
				ExecMode:    uint32(execMode),
			},
		},
		[]uint32{clobPairId},
	)
}

// SendSubaccountUpdates groups subaccount updates by their subaccount ids and
// sends messages to the subscribers.
func (sm *FullNodeStreamingManagerImpl) SendSubaccountUpdates(
	subaccountUpdates []satypes.StreamSubaccountUpdate,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendSubaccountUpdatesLatency,
		time.Now(),
	)

	// Group subaccount updates by subaccount id.
	streamUpdates := make([]clobtypes.StreamUpdate, 0)
	subaccountIds := make([]*satypes.SubaccountId, 0)
	for _, subaccountUpdate := range subaccountUpdates {
		streamUpdate := clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_SubaccountUpdate{
				SubaccountUpdate: &subaccountUpdate,
			},
			BlockHeight: blockHeight,
			ExecMode:    uint32(execMode),
		}
		streamUpdates = append(streamUpdates, streamUpdate)
		subaccountIds = append(subaccountIds, subaccountUpdate.SubaccountId)
	}

	sm.AddSubaccountUpdatesToCache(streamUpdates, subaccountIds)
}

// AddOrderUpdatesToCache adds a series of updates to the full node streaming cache.
// Clob pair ids are the clob pair id each update is relevant to.
func (sm *FullNodeStreamingManagerImpl) AddOrderUpdatesToCache(
	updates []clobtypes.StreamUpdate,
	clobPairIds []uint32,
) {
	sm.Lock()
	defer sm.Unlock()

	metrics.IncrCounter(
		metrics.GrpcAddUpdateToBufferCount,
		float32(len(updates)),
	)

	sm.streamUpdateCache = append(sm.streamUpdateCache, updates...)
	for _, clobPairId := range clobPairIds {
		sm.streamUpdateSubscriptionCache = append(
			sm.streamUpdateSubscriptionCache,
			sm.clobPairIdToSubscriptionIdMapping[clobPairId],
		)
	}

	// Remove all subscriptions and wipe the buffer if buffer overflows.
	sm.RemoveSubscriptionsAndClearBufferIfFull()
	sm.EmitMetrics()
}

// AddSubaccountUpdatesToCache adds a series of updates to the full node streaming cache.
// Subaccount ids are the subaccount id each update is relevant to.
func (sm *FullNodeStreamingManagerImpl) AddSubaccountUpdatesToCache(
	updates []clobtypes.StreamUpdate,
	subaccountIds []*satypes.SubaccountId,
) {
	sm.Lock()
	defer sm.Unlock()

	metrics.IncrCounter(
		metrics.GrpcAddUpdateToBufferCount,
		float32(len(updates)),
	)

	sm.streamUpdateCache = append(sm.streamUpdateCache, updates...)
	for _, subaccountId := range subaccountIds {
		sm.streamUpdateSubscriptionCache = append(
			sm.streamUpdateSubscriptionCache,
			sm.subaccountIdToSubscriptionIdMapping[*subaccountId],
		)
	}
	sm.RemoveSubscriptionsAndClearBufferIfFull()
	sm.EmitMetrics()
}

// RemoveSubscriptionsAndClearBufferIfFull removes all subscriptions and wipes the buffer if buffer overflows.
// Note this method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) RemoveSubscriptionsAndClearBufferIfFull() {
	// Remove all subscriptions and wipe the buffer if buffer overflows.
	if len(sm.streamUpdateCache) > int(sm.maxUpdatesInCache) {
		sm.logger.Error("Streaming buffer full capacity. Dropping messages and all subscriptions. " +
			"Disconnect all clients and increase buffer size via the grpc-stream-buffer-size flag.")
		for id := range sm.orderbookSubscriptions {
			sm.removeSubscription(id)
		}
		sm.streamUpdateCache = nil
		sm.streamUpdateSubscriptionCache = nil
	}
}

func (sm *FullNodeStreamingManagerImpl) FlushStreamUpdates() {
	sm.Lock()
	defer sm.Unlock()
	sm.FlushStreamUpdatesWithLock()
}

// FlushStreamUpdatesWithLock takes in a list of stream updates and their corresponding subscription IDs,
// and emits them to subscribers. Note this method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) FlushStreamUpdatesWithLock() {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcFlushUpdatesLatency,
		time.Now(),
	)

	// Map to collect updates for each subscription.
	subscriptionUpdates := make(map[uint32][]clobtypes.StreamUpdate)
	idsToRemove := make([]uint32, 0)

	// Collect updates for each subscription.
	for i, update := range sm.streamUpdateCache {
		subscriptionIds := sm.streamUpdateSubscriptionCache[i]
		for _, id := range subscriptionIds {
			subscriptionUpdates[id] = append(subscriptionUpdates[id], update)
		}
	}

	// Non-blocking send updates through subscriber's buffered channel.
	// If the buffer is full, drop the subscription.
	for id, updates := range subscriptionUpdates {
		if subscription, ok := sm.orderbookSubscriptions[id]; ok {
			metrics.IncrCounter(
				metrics.GrpcAddToSubscriptionChannelCount,
				1,
			)
			select {
			case subscription.updatesChannel <- updates:
			default:
				idsToRemove = append(idsToRemove, id)
			}
		}
	}

	sm.streamUpdateCache = nil
	sm.streamUpdateSubscriptionCache = nil

	for _, id := range idsToRemove {
		sm.logger.Error(
			fmt.Sprintf(
				"Streaming subscription id %+v channel full capacity. Dropping subscription connection.",
				id,
			),
		)
		sm.removeSubscription(id)
	}

	sm.EmitMetrics()
}

func (sm *FullNodeStreamingManagerImpl) InitializeNewStreams(
	getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
	getSubaccountSnapshot func(subaccountId satypes.SubaccountId) *satypes.StreamSubaccountUpdate,
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
		// If the snapshot block interval is enabled, reset the sync.Once in order to
		// re-send snapshots out.
		if sm.snapshotBlockInterval > 0 &&
			blockHeight == subscription.nextSnapshotBlock {
			subscription.initialize = &sync.Once{}
		}

		subscription.initialize.Do(
			func() {
				allUpdates := clobtypes.NewOffchainUpdates()
				for _, clobPairId := range subscription.clobPairIds {
					if _, ok := updatesByClobPairId[clobPairId]; !ok {
						updatesByClobPairId[clobPairId] = getOrderbookSnapshot(clobtypes.ClobPairId(clobPairId))
					}
					allUpdates.Append(updatesByClobPairId[clobPairId])
				}
				saUpdates := []*satypes.StreamSubaccountUpdate{}
				for _, subaccountId := range subscription.subaccountIds {
					saUpdates = append(saUpdates, getSubaccountSnapshot(subaccountId))
				}
				sm.SendCombinedSnapshot(allUpdates, saUpdates, subscriptionId, blockHeight, execMode)
				if sm.snapshotBlockInterval != 0 {
					subscription.nextSnapshotBlock = blockHeight + sm.snapshotBlockInterval
				}
			},
		)
	}
}
