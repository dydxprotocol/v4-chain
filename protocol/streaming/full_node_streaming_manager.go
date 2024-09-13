package streaming

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ante_types "github.com/dydxprotocol/v4-chain/protocol/app/ante/types"
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

	// stores the staged FinalizeBlock events for full node streaming.
	streamingManagerTransientStoreKey storetypes.StoreKey
}

// OrderbookSubscription represents a active subscription to the orderbook updates stream.
type OrderbookSubscription struct {
	subscriptionId uint32

	// Whether the subscription is initialized with snapshot.
	initialized *atomic.Bool

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

func (sub *OrderbookSubscription) IsInitialized() bool {
	return sub.initialized.Load()
}

func NewFullNodeStreamingManager(
	logger log.Logger,
	flushIntervalMs uint32,
	maxUpdatesInCache uint32,
	maxSubscriptionChannelSize uint32,
	snapshotBlockInterval uint32,
	streamingManagerTransientStoreKey storetypes.StoreKey,
) *FullNodeStreamingManagerImpl {
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

		streamingManagerTransientStoreKey: streamingManagerTransientStoreKey,
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
	if len(clobPairIds) == 0 && len(subaccountIds) == 0 {
		return types.ErrInvalidStreamingRequest
	}

	sm.Lock()
	sIds := make([]satypes.SubaccountId, len(subaccountIds))
	for i, subaccountId := range subaccountIds {
		sIds[i] = *subaccountId
	}
	subscription := &OrderbookSubscription{
		subscriptionId: sm.nextSubscriptionId,
		initialized:    &atomic.Bool{}, // False by default.
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

func getStagedEventsCount(store storetypes.KVStore) uint32 {
	countsBytes := store.Get([]byte(StagedEventsCountKey))
	if countsBytes == nil {
		return 0
	}
	return binary.BigEndian.Uint32(countsBytes)
}

// Stage a subaccount update event in transient store, during `FinalizeBlock`.
func (sm *FullNodeStreamingManagerImpl) StageFinalizeBlockSubaccountUpdate(
	ctx sdk.Context,
	subaccountUpdate satypes.StreamSubaccountUpdate,
) {
	stagedEvent := clobtypes.StagedFinalizeBlockEvent{
		Event: &clobtypes.StagedFinalizeBlockEvent_SubaccountUpdate{
			SubaccountUpdate: &subaccountUpdate,
		},
	}
	sm.stageFinalizeBlockEvent(
		ctx,
		clobtypes.Amino.MustMarshal(stagedEvent),
	)
}

// Stage a fill event in transient store, during `FinalizeBlock`.
func (sm *FullNodeStreamingManagerImpl) StageFinalizeBlockFill(
	ctx sdk.Context,
	fill clobtypes.StreamOrderbookFill,
) {
	stagedEvent := clobtypes.StagedFinalizeBlockEvent{
		Event: &clobtypes.StagedFinalizeBlockEvent_OrderFill{
			OrderFill: &fill,
		},
	}
	sm.stageFinalizeBlockEvent(
		ctx,
		clobtypes.Amino.MustMarshal(stagedEvent),
	)
}

func getStagedFinalizeBlockEvents(store storetypes.KVStore) []clobtypes.StagedFinalizeBlockEvent {
	count := getStagedEventsCount(store)
	events := make([]clobtypes.StagedFinalizeBlockEvent, count)
	store = prefix.NewStore(store, []byte(StagedEventsKeyPrefix))
	for i := uint32(0); i < count; i++ {
		var event clobtypes.StagedFinalizeBlockEvent
		bytes := store.Get(lib.Uint32ToKey(i))
		clobtypes.Amino.MustUnmarshal(bytes, &event)
		events[i] = event
	}
	return events
}

// Retrieve all events staged during `FinalizeBlock`.
func (sm *FullNodeStreamingManagerImpl) GetStagedFinalizeBlockEvents(
	ctx sdk.Context,
) []clobtypes.StagedFinalizeBlockEvent {
	noGasCtx := ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())
	store := noGasCtx.TransientStore(sm.streamingManagerTransientStoreKey)
	return getStagedFinalizeBlockEvents(store)
}

func (sm *FullNodeStreamingManagerImpl) stageFinalizeBlockEvent(
	ctx sdk.Context,
	eventBytes []byte,
) {
	noGasCtx := ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())
	store := noGasCtx.TransientStore(sm.streamingManagerTransientStoreKey)

	// Increment events count.
	count := getStagedEventsCount(store)
	store.Set([]byte(StagedEventsCountKey), lib.Uint32ToKey(count+1))

	// Store events keyed by index.
	store = prefix.NewStore(store, []byte(StagedEventsKeyPrefix))
	store.Set(lib.Uint32ToKey(count), eventBytes)
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

// SendFinalizedSubaccountUpdates groups subaccount updates by their subaccount ids and
// sends messages to the subscribers.
func (sm *FullNodeStreamingManagerImpl) SendFinalizedSubaccountUpdates(
	subaccountUpdates []satypes.StreamSubaccountUpdate,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendFinalizedSubaccountUpdatesLatency,
		time.Now(),
	)

	if execMode != sdk.ExecModeFinalize {
		panic("SendFinalizedSubaccountUpdates should only be called in ExecModeFinalize")
	}

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

func (sm *FullNodeStreamingManagerImpl) GetSubaccountSnapshotsForInitStreams(
	getSubaccountSnapshot func(subaccountId satypes.SubaccountId) *satypes.StreamSubaccountUpdate,
) map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate {
	sm.Lock()
	defer sm.Unlock()

	ret := make(map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate)
	for _, subscription := range sm.orderbookSubscriptions {
		// If the subscription has been initialized, no need to grab the subaccount snapshot.
		if alreadyInitialized := subscription.initialized.Load(); alreadyInitialized {
			continue
		}

		for _, subaccountId := range subscription.subaccountIds {
			if _, exists := ret[subaccountId]; exists {
				continue
			}

			ret[subaccountId] = getSubaccountSnapshot(subaccountId)
		}
	}
	return ret
}

// Grpc Streaming logic after consensus agrees on a block.
// - Stream all events staged during `FinalizeBlock`.
// - Stream orderbook updates to sync fills in local ops queue.
func (sm *FullNodeStreamingManagerImpl) StreamBatchUpdatesAfterFinalizeBlock(
	ctx sdk.Context,
	orderBookUpdatesToSyncLocalOpsQueue *clobtypes.OffchainUpdates,
	perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
) {
	// Flush all pending updates, since we want the onchain updates to arrive in a batch.
	sm.FlushStreamUpdates()

	finalizedFills, finalizedSubaccountUpdates := sm.getStagedEventsFromFinalizeBlock(ctx)

	// TODO(CT-1190): Stream below in a single batch.
	// Send orderbook updates to sync optimistic orderbook onchain state after FinalizeBlock.
	sm.SendOrderbookUpdates(
		orderBookUpdatesToSyncLocalOpsQueue,
		uint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)

	// Send finalized fills from FinalizeBlock.
	sm.SendOrderbookFillUpdates(
		finalizedFills,
		uint32(ctx.BlockHeight()),
		ctx.ExecMode(),
		perpetualIdToClobPairId,
	)

	// Send finalized subaccount updates from FinalizeBlock.
	sm.SendFinalizedSubaccountUpdates(
		finalizedSubaccountUpdates,
		uint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)
}

// getStagedEventsFromFinalizeBlock returns staged events from `FinalizeBlock`.
// It should be called after the consensus agrees on a block (e.g. Precommitter).
func (sm *FullNodeStreamingManagerImpl) getStagedEventsFromFinalizeBlock(
	ctx sdk.Context,
) (
	finalizedFills []clobtypes.StreamOrderbookFill,
	finalizedSubaccountUpdates []satypes.StreamSubaccountUpdate,
) {
	// Get onchain stream events stored in transient store.
	stagedEvents := sm.GetStagedFinalizeBlockEvents(ctx)

	telemetry.SetGauge(
		float32(len(stagedEvents)),
		types.ModuleName,
		metrics.GrpcStagedAllFinalizeBlockUpdates,
		metrics.Count,
	)

	for _, stagedEvent := range stagedEvents {
		switch event := stagedEvent.Event.(type) {
		case *clobtypes.StagedFinalizeBlockEvent_OrderFill:
			finalizedFills = append(finalizedFills, *event.OrderFill)
		case *clobtypes.StagedFinalizeBlockEvent_SubaccountUpdate:
			finalizedSubaccountUpdates = append(finalizedSubaccountUpdates, *event.SubaccountUpdate)
		}
	}

	telemetry.SetGauge(
		float32(len(finalizedSubaccountUpdates)),
		types.ModuleName,
		metrics.GrpcStagedSubaccountFinalizeBlockUpdates,
		metrics.Count,
	)
	telemetry.SetGauge(
		float32(len(finalizedFills)),
		types.ModuleName,
		metrics.GrpcStagedFillFinalizeBlockUpdates,
		metrics.Count,
	)

	return finalizedFills, finalizedSubaccountUpdates
}

func (sm *FullNodeStreamingManagerImpl) InitializeNewStreams(
	getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
	subaccountSnapshots map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate,
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
		if alreadyInitialized := subscription.initialized.Swap(true); !alreadyInitialized {
			allUpdates := clobtypes.NewOffchainUpdates()
			for _, clobPairId := range subscription.clobPairIds {
				if _, ok := updatesByClobPairId[clobPairId]; !ok {
					updatesByClobPairId[clobPairId] = getOrderbookSnapshot(clobtypes.ClobPairId(clobPairId))
				}
				allUpdates.Append(updatesByClobPairId[clobPairId])
			}

			saUpdates := []*satypes.StreamSubaccountUpdate{}
			for _, subaccountId := range subscription.subaccountIds {
				// The subaccount snapshot may not exist due to the following race condition
				// 1. At beginning of PrepareCheckState we get snapshot for all subscribed subaccounts.
				// 2. A new subaccount is subscribed to by a new subscription.
				// 3. InitializeNewStreams is called.
				// Then the new subaccount would not be included in the snapshot.
				// We are okay with this behavior.
				if saUpdate, ok := subaccountSnapshots[subaccountId]; ok {
					saUpdates = append(saUpdates, saUpdate)
				}
			}

			sm.SendCombinedSnapshot(allUpdates, saUpdates, subscriptionId, blockHeight, execMode)

			if sm.snapshotBlockInterval != 0 {
				subscription.nextSnapshotBlock = blockHeight + sm.snapshotBlockInterval
			}
		}

		// If the snapshot block interval is enabled and the next block is a snapshot block,
		// reset the `atomic.Bool` so snapshots are sent for the next block.
		if sm.snapshotBlockInterval > 0 &&
			blockHeight+1 == subscription.nextSnapshotBlock {
			subscription.initialized = &atomic.Bool{} // False by default.
		}
	}
}
