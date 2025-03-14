package streaming

import (
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ante_types "github.com/dydxprotocol/v4-chain/protocol/app/ante/types"
	"github.com/dydxprotocol/v4-chain/protocol/finalizeblock"
	ocutypes "github.com/dydxprotocol/v4-chain/protocol/indexer/off_chain_updates/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	streaming_util "github.com/dydxprotocol/v4-chain/protocol/streaming/util"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var _ types.FullNodeStreamingManager = (*FullNodeStreamingManagerImpl)(nil)

// FullNodeStreamingManagerImpl is an implementation for managing streaming subscriptions.
type FullNodeStreamingManagerImpl struct {
	sync.Mutex

	cdc    codec.BinaryCodec
	logger log.Logger

	// orderbookSubscriptions maps subscription IDs to their respective orderbook subscriptions.
	orderbookSubscriptions map[uint32]*OrderbookSubscription
	activeSubscriptionIds  map[uint32]bool

	// stream will batch and flush out messages every 10 ms.
	ticker *time.Ticker
	done   chan bool

	// TODO: Consolidate the streamUpdateCache and streamUpdateSubscriptionCache into a single
	// struct to avoid the need to maintain two separate slices for the same data.

	streamUpdateCache                   []clobtypes.StreamUpdate
	streamUpdateSubscriptionCache       [][]uint32
	clobPairIdToSubscriptionIdMapping   map[uint32][]uint32
	subaccountIdToSubscriptionIdMapping map[satypes.SubaccountId][]uint32
	marketIdToSubscriptionIdMapping     map[uint32][]uint32
	allClobPairSubscriptionIdMapping    map[uint32]struct{}

	maxUpdatesInCache          uint32
	maxSubscriptionChannelSize uint32

	// Block interval in which snapshot info should be sent out in.
	// Defaults to 0, which means only one snapshot will be sent out.
	snapshotBlockInterval uint32

	// stores the staged FinalizeBlock events for full node streaming.
	streamingManagerTransientStoreKey storetypes.StoreKey

	finalizeBlockStager finalizeblock.EventStager[*clobtypes.StagedFinalizeBlockEvent]
}

// OrderbookSubscription represents a active subscription to the orderbook updates stream.
type OrderbookSubscription struct {
	subscriptionId uint32

	initializedWithSnapshot *atomic.Bool
	clobPairIds             []uint32
	subaccountIds           []satypes.SubaccountId
	marketIds               []uint32

	messageSender        types.OutgoingMessageSender
	streamUpdatesChannel chan []clobtypes.StreamUpdate

	// If interval snapshots are turned on, the next block height at which
	// a snapshot should be sent out.
	nextSnapshotBlock uint32
}

func (sm *FullNodeStreamingManagerImpl) AllClobPairSubscriptionIds() []uint32 {
	allClobPairSubscriptionIds := []uint32{}
	for subscriptionId := range sm.allClobPairSubscriptionIdMapping {
		allClobPairSubscriptionIds = append(allClobPairSubscriptionIds, subscriptionId)
	}
	return allClobPairSubscriptionIds
}

func (sm *FullNodeStreamingManagerImpl) NewOrderbookSubscription(
	clobPairIds []uint32,
	subaccountIds []satypes.SubaccountId,
	marketIds []uint32,
	messageSender types.OutgoingMessageSender,
) *OrderbookSubscription {
	return &OrderbookSubscription{
		subscriptionId:          sm.getNextAvailableSubscriptionId(),
		initializedWithSnapshot: &atomic.Bool{}, // False by default.
		clobPairIds:             clobPairIds,
		subaccountIds:           subaccountIds,
		marketIds:               marketIds,
		messageSender:           messageSender,
		streamUpdatesChannel:    make(chan []clobtypes.StreamUpdate, sm.maxSubscriptionChannelSize),
	}
}

func (sub *OrderbookSubscription) IsInitialized() bool {
	return sub.initializedWithSnapshot.Load()
}

func NewFullNodeStreamingManager(
	logger log.Logger,
	flushIntervalMs uint32,
	maxUpdatesInCache uint32,
	maxSubscriptionChannelSize uint32,
	snapshotBlockInterval uint32,
	streamingManagerTransientStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) *FullNodeStreamingManagerImpl {
	fullNodeStreamingManager := &FullNodeStreamingManagerImpl{
		logger:                 logger,
		orderbookSubscriptions: make(map[uint32]*OrderbookSubscription),
		activeSubscriptionIds:  make(map[uint32]bool),

		ticker:                              time.NewTicker(time.Duration(flushIntervalMs) * time.Millisecond),
		done:                                make(chan bool),
		streamUpdateCache:                   make([]clobtypes.StreamUpdate, 0),
		streamUpdateSubscriptionCache:       make([][]uint32, 0),
		clobPairIdToSubscriptionIdMapping:   make(map[uint32][]uint32),
		subaccountIdToSubscriptionIdMapping: make(map[satypes.SubaccountId][]uint32),
		marketIdToSubscriptionIdMapping:     make(map[uint32][]uint32),
		allClobPairSubscriptionIdMapping:    make(map[uint32]struct{}),

		maxUpdatesInCache:          maxUpdatesInCache,
		maxSubscriptionChannelSize: maxSubscriptionChannelSize,
		snapshotBlockInterval:      snapshotBlockInterval,

		streamingManagerTransientStoreKey: streamingManagerTransientStoreKey,
		cdc:                               cdc,
		finalizeBlockStager: finalizeblock.NewEventStager[*clobtypes.StagedFinalizeBlockEvent](
			streamingManagerTransientStoreKey,
			cdc,
			StagedEventsCountKey,
			StagedEventsKeyPrefix,
		),
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
	metrics.AddSample(
		metrics.GrpcStreamNumUpdatesBuffered,
		float32(len(sm.streamUpdateCache)),
	)
	metrics.SetGauge(
		metrics.GrpcStreamSubscriberCount,
		float32(len(sm.orderbookSubscriptions)),
	)
	for _, subscription := range sm.orderbookSubscriptions {
		metrics.AddSampleWithLabels(
			metrics.GrpcSubscriptionChannelLength,
			float32(len(subscription.streamUpdatesChannel)),
			metrics.GetLabelForIntValue(metrics.SubscriptionId, int(subscription.subscriptionId)),
		)
	}
}

// getNextAvailableSubscriptionId returns next available subscription id. Assumes the
// lock has been acquired.
func (sm *FullNodeStreamingManagerImpl) getNextAvailableSubscriptionId() uint32 {
	id := uint32(0)
	for _, inUse := sm.activeSubscriptionIds[id]; inUse; _, inUse = sm.activeSubscriptionIds[id] {
		id = id + uint32(1)
	}
	return id
}

func doFilterOrderbookUpdateBySubaccount(
	orderBookUpdate *clobtypes.StreamUpdate_OrderbookUpdate,
	subaccountIds []satypes.SubaccountId,
) (bool, error) {
	if orderBookUpdate.OrderbookUpdate.Snapshot {
		return true, nil
	}
	for _, orderBookUpdate := range orderBookUpdate.OrderbookUpdate.Updates {
		orderBookUpdateSubaccountId, err := streaming_util.GetOffChainUpdateV1SubaccountId(orderBookUpdate)
		if err != nil {
			return false, err
		}
		if slices.Contains(subaccountIds, orderBookUpdateSubaccountId) {
			return true, nil
		}
	}
	return false, nil
}

func doFilterTakerOrderBySubaccount(
	takerOrder *clobtypes.StreamUpdate_TakerOrder,
	subaccountIds []satypes.SubaccountId,
) bool {
	return slices.Contains(subaccountIds, takerOrder.TakerOrder.GetOrder().OrderId.SubaccountId)
}

func doFilterOrderFillBySubaccount(
	orderFill *clobtypes.StreamUpdate_OrderFill,
	subaccountIds []satypes.SubaccountId,
) bool {
	switch match := orderFill.OrderFill.GetClobMatch().GetMatch().(type) {
	case *clobtypes.ClobMatch_MatchOrders:
		if slices.Contains(subaccountIds, match.MatchOrders.TakerOrderId.SubaccountId) {
			return true
		}
		for _, fill := range match.MatchOrders.Fills {
			if slices.Contains(subaccountIds, fill.MakerOrderId.SubaccountId) {
				return true
			}
		}
		return false
	case *clobtypes.ClobMatch_MatchPerpetualLiquidation:
		if slices.Contains(subaccountIds, match.MatchPerpetualLiquidation.Liquidated) {
			return true
		}
		for _, fill := range match.MatchPerpetualLiquidation.Fills {
			if slices.Contains(subaccountIds, fill.MakerOrderId.SubaccountId) {
				return true
			}
		}
		return false
	case *clobtypes.ClobMatch_MatchPerpetualDeleveraging:
		if slices.Contains(subaccountIds, match.MatchPerpetualDeleveraging.Liquidated) {
			return true
		}
		for _, fill := range match.MatchPerpetualDeleveraging.Fills {
			if slices.Contains(subaccountIds, fill.OffsettingSubaccountId) {
				return true
			}
		}
		return false
	case nil:
		return false
	}
	return true
}

func doFilterStreamUpdateBySubaccount(
	update *clobtypes.StreamUpdate,
	subaccountIds []satypes.SubaccountId,
) (bool, error) {
	// If reflection becomes too expensive, split updatesChannel by message type
	switch updateMessage := update.UpdateMessage.(type) {
	case *clobtypes.StreamUpdate_OrderbookUpdate:
		return doFilterOrderbookUpdateBySubaccount(updateMessage, subaccountIds)
	case *clobtypes.StreamUpdate_TakerOrder:
		return doFilterTakerOrderBySubaccount(updateMessage, subaccountIds), nil
	case *clobtypes.StreamUpdate_OrderFill:
		return doFilterOrderFillBySubaccount(updateMessage, subaccountIds), nil
	}
	return true, nil
}

// If UpdateMessage is not a StreamUpdate_OrderUpdate, filter it
// If a StreamUpdate_OrderUpdate contains updates for subscribed subaccounts, filter it
// If a StreamUpdate_OrderUpdate contains no updates for subscribed subaccounts, drop it
// If checking subaccount ids in a StreamUpdate_OrderUpdate results in an error, log error and drop it
func FilterStreamUpdateBySubaccount(
	updates []clobtypes.StreamUpdate,
	subaccountIds []satypes.SubaccountId,
	logger log.Logger,
) []clobtypes.StreamUpdate {
	filteredUpdates := []clobtypes.StreamUpdate{}
	for _, update := range updates {
		doFilter, err := doFilterStreamUpdateBySubaccount(&update, subaccountIds)
		if err != nil {
			logger.Error(err.Error())
		}
		if doFilter {
			filteredUpdates = append(filteredUpdates, update)
		}
	}
	return filteredUpdates
}

func (sm *FullNodeStreamingManagerImpl) Subscribe(
	clobPairIds []uint32,
	subaccountIds []*satypes.SubaccountId,
	marketIds []uint32,
	filterOrdersBySubAccountId bool,
	messageSender types.OutgoingMessageSender,
) (
	err error,
) {
	// Perform some basic validation on the request.
	if len(clobPairIds) == 0 && len(subaccountIds) == 0 && len(marketIds) == 0 {
		return types.ErrInvalidStreamingRequest
	}
	if filterOrdersBySubAccountId && (len(subaccountIds) == 0) {
		sm.logger.Error("filterOrdersBySubaccountId with no subaccountIds")
		return types.ErrInvalidSubaccountFilteringRequest
	}

	sm.Lock()
	sIds := make([]satypes.SubaccountId, len(subaccountIds))
	for i, subaccountId := range subaccountIds {
		sIds[i] = *subaccountId
	}

	subscription := sm.NewOrderbookSubscription(clobPairIds, sIds, marketIds, messageSender)
	if len(clobPairIds) == 0 {
		sm.allClobPairSubscriptionIdMapping[subscription.subscriptionId] = struct{}{}
	}
	for _, clobPairId := range clobPairIds {
		// if clobPairId exists in the map, append the subscription id to the slice
		// otherwise, create a new slice with the subscription id
		if _, ok := sm.clobPairIdToSubscriptionIdMapping[clobPairId]; !ok {
			sm.clobPairIdToSubscriptionIdMapping[clobPairId] = []uint32{}
		}
		sm.clobPairIdToSubscriptionIdMapping[clobPairId] = append(
			sm.clobPairIdToSubscriptionIdMapping[clobPairId],
			subscription.subscriptionId,
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
			subscription.subscriptionId,
		)
	}
	for _, marketId := range marketIds {
		// if subaccountId exists in the map, append the subscription id to the slice
		// otherwise, create a new slice with the subscription id
		if _, ok := sm.marketIdToSubscriptionIdMapping[marketId]; !ok {
			sm.marketIdToSubscriptionIdMapping[marketId] = []uint32{}
		}
		sm.marketIdToSubscriptionIdMapping[marketId] = append(
			sm.marketIdToSubscriptionIdMapping[marketId],
			subscription.subscriptionId,
		)
	}

	var clobPairIdString string
	if _, ok := sm.allClobPairSubscriptionIdMapping[subscription.subscriptionId]; ok {
		clobPairIdString = "*"
	} else {
		clobPairIdString = fmt.Sprintf("%+v", clobPairIds)
	}
	sm.logger.Info(
		fmt.Sprintf(
			"New subscription id %+v for clob pair ids: %v and subaccount ids: %+v. filter orders by subaccount ids: %+v",
			subscription.subscriptionId,
			clobPairIdString,
			subaccountIds,
			filterOrdersBySubAccountId,
		),
	)
	sm.orderbookSubscriptions[subscription.subscriptionId] = subscription
	sm.activeSubscriptionIds[subscription.subscriptionId] = true
	sm.EmitMetrics()
	sm.Unlock()

	// Use current goroutine to consistently poll subscription channel for updates
	// to send through stream.
	for updates := range subscription.streamUpdatesChannel {
		if filterOrdersBySubAccountId {
			updates = FilterStreamUpdateBySubaccount(updates, sIds, sm.logger)
		}
		if len(updates) == 0 {
			continue
		}
		metrics.IncrCounterWithLabels(
			metrics.GrpcSendResponseToSubscriberCount,
			1,
			metrics.GetLabelForIntValue(metrics.SubscriptionId, int(subscription.subscriptionId)),
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
	close(subscription.streamUpdatesChannel)
	delete(sm.orderbookSubscriptions, subscriptionIdToRemove)
	delete(sm.activeSubscriptionIds, subscriptionIdToRemove)
	delete(sm.allClobPairSubscriptionIdMapping, subscriptionIdToRemove)

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

	// Iterate over the marketIdToSubscriptionIdMapping to remove the subscriptionIdToRemove
	for marketId, subscriptionIds := range sm.marketIdToSubscriptionIdMapping {
		for i, id := range subscriptionIds {
			if id == subscriptionIdToRemove {
				// Remove the subscription ID from the slice
				sm.marketIdToSubscriptionIdMapping[marketId] = append(subscriptionIds[:i], subscriptionIds[i+1:]...)
				break
			}
		}
		// If the list is empty after removal, delete the key from the map
		if len(sm.marketIdToSubscriptionIdMapping[marketId]) == 0 {
			delete(sm.marketIdToSubscriptionIdMapping, marketId)
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
	v1updates := streaming_util.GetOffchainUpdatesV1(offchainUpdates)
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

func toPriceStreamUpdates(
	priceUpdates []*pricestypes.StreamPriceUpdate,
	blockHeight uint32,
	execMode sdk.ExecMode,
) []clobtypes.StreamUpdate {
	streamUpdates := make([]clobtypes.StreamUpdate, 0)
	for _, update := range priceUpdates {
		streamUpdates = append(streamUpdates, clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_PriceUpdate{
				PriceUpdate: update,
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

	metrics.IncrCounterWithLabels(
		metrics.GrpcAddToSubscriptionChannelCount,
		1,
		metrics.GetLabelForIntValue(metrics.SubscriptionId, int(subscriptionId)),
	)

	select {
	case subscription.streamUpdatesChannel <- streamUpdates:
	default:
		// Buffer is full. Emit metric and drop subscription.
		sm.EmitMetrics()
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

// Send a subaccount update event.
func (sm *FullNodeStreamingManagerImpl) SendSubaccountUpdate(
	ctx sdk.Context,
	subaccountUpdate satypes.StreamSubaccountUpdate,
) {
	// If not `DeliverTx`, return since we don't stream optimistic subaccount updates.
	if !lib.IsDeliverTxMode(ctx) {
		return
	}

	metrics.IncrCounter(
		metrics.GrpcSendSubaccountUpdateCount,
		1,
	)

	// If `DeliverTx`, updates should be staged to be streamed after consensus finalizes on a block.
	stagedEvent := clobtypes.StagedFinalizeBlockEvent{
		Event: &clobtypes.StagedFinalizeBlockEvent_SubaccountUpdate{
			SubaccountUpdate: &subaccountUpdate,
		},
	}
	sm.finalizeBlockStager.StageFinalizeBlockEvent(
		ctx,
		&stagedEvent,
	)
}

// SendPriceUpdates sends price updates to the subscribers.
func (sm *FullNodeStreamingManagerImpl) SendPriceUpdate(
	ctx sdk.Context,
	priceUpdate pricestypes.StreamPriceUpdate,
) {
	if !lib.IsDeliverTxMode(ctx) {
		// If not `DeliverTx`, return since there is no optimistic price updates.
		return
	}

	metrics.IncrCounter(
		metrics.GrpcSendPriceUpdateCount,
		1,
	)

	// If `DeliverTx`, updates should be staged to be streamed after consensus finalizes on a block.
	stagedEvent := clobtypes.StagedFinalizeBlockEvent{
		Event: &clobtypes.StagedFinalizeBlockEvent_PriceUpdate{
			PriceUpdate: &priceUpdate,
		},
	}
	sm.finalizeBlockStager.StageFinalizeBlockEvent(
		ctx,
		&stagedEvent,
	)
}

// Retrieve all events staged during `FinalizeBlock`.
func (sm *FullNodeStreamingManagerImpl) GetStagedFinalizeBlockEvents(
	ctx sdk.Context,
) []clobtypes.StagedFinalizeBlockEvent {
	events := sm.finalizeBlockStager.GetStagedFinalizeBlockEvents(
		ctx,
		func() *clobtypes.StagedFinalizeBlockEvent {
			return &clobtypes.StagedFinalizeBlockEvent{}
		},
	)
	results := make([]clobtypes.StagedFinalizeBlockEvent, len(events))
	for i, event := range events {
		if event == nil {
			panic("Got nil event from finalizeBlockStager")
		}
		results[i] = *event
	}
	return results
}

// SendCombinedSnapshot sends messages to a particular subscriber without buffering.
// Note this method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) SendCombinedSnapshot(
	offchainUpdates *clobtypes.OffchainUpdates,
	saUpdates []*satypes.StreamSubaccountUpdate,
	priceUpdates []*pricestypes.StreamPriceUpdate,
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
	streamUpdates = append(streamUpdates, toPriceStreamUpdates(priceUpdates, blockHeight, execMode)...)
	sm.sendStreamUpdates(subscriptionId, streamUpdates)
}

// TracksSubaccountId checks if a subaccount id is being tracked by the streaming manager.
func (sm *FullNodeStreamingManagerImpl) TracksSubaccountId(subaccountId satypes.SubaccountId) bool {
	sm.Lock()
	defer sm.Unlock()
	_, exists := sm.subaccountIdToSubscriptionIdMapping[subaccountId]
	return exists
}

// TracksMarketId checks if a market id is being tracked by the streaming manager.
func (sm *FullNodeStreamingManagerImpl) TracksMarketId(marketId uint32) bool {
	sm.Lock()
	defer sm.Unlock()
	_, exists := sm.marketIdToSubscriptionIdMapping[marketId]
	return exists
}

func getStreamUpdatesFromOffchainUpdates(
	v1updates []ocutypes.OffChainUpdateV1,
	blockHeight uint32,
	execMode sdk.ExecMode,
) (streamUpdates []clobtypes.StreamUpdate, clobPairIds []uint32) {
	// Group updates by clob pair id.
	clobPairIdToV1Updates := make(map[uint32][]ocutypes.OffChainUpdateV1)
	// unique list of clob pair Ids to send updates for.
	clobPairIds = make([]uint32, 0)
	for _, v1update := range v1updates {
		var clobPairId uint32
		switch u := v1update.UpdateMessage.(type) {
		case *ocutypes.OffChainUpdateV1_OrderPlace:
			clobPairId = u.OrderPlace.Order.OrderId.ClobPairId
		case *ocutypes.OffChainUpdateV1_OrderReplace:
			clobPairId = u.OrderReplace.OldOrderId.ClobPairId
		case *ocutypes.OffChainUpdateV1_OrderRemove:
			clobPairId = u.OrderRemove.RemovedOrderId.ClobPairId
		case *ocutypes.OffChainUpdateV1_OrderUpdate:
			clobPairId = u.OrderUpdate.OrderId.ClobPairId
		default:
			panic(fmt.Sprintf("Unhandled UpdateMessage type: %v", u))
		}

		if _, ok := clobPairIdToV1Updates[clobPairId]; !ok {
			clobPairIdToV1Updates[clobPairId] = []ocutypes.OffChainUpdateV1{}
			clobPairIds = append(clobPairIds, clobPairId)
		}
		clobPairIdToV1Updates[clobPairId] = append(clobPairIdToV1Updates[clobPairId], v1update)
	}

	// Unmarshal each per-clob pair message to v1 updates.
	streamUpdates = make([]clobtypes.StreamUpdate, len(clobPairIds))

	for i, clobPairId := range clobPairIds {
		v1updates, exists := clobPairIdToV1Updates[clobPairId]
		if !exists {
			panic(fmt.Sprintf(
				"clob pair id %v not found in clobPairIdToV1Updates: %v",
				clobPairId,
				clobPairIdToV1Updates,
			))
		}
		streamUpdates[i] = clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_OrderbookUpdate{
				OrderbookUpdate: &clobtypes.StreamOrderbookUpdate{
					Updates:  v1updates,
					Snapshot: false,
				},
			},
			BlockHeight: blockHeight,
			ExecMode:    uint32(execMode),
		}
	}

	return streamUpdates, clobPairIds
}

// SendOrderbookUpdates groups updates by their clob pair ids and
// sends messages to the subscribers.
func (sm *FullNodeStreamingManagerImpl) SendOrderbookUpdates(
	offchainUpdates *clobtypes.OffchainUpdates,
	ctx sdk.Context,
) {
	v1updates := streaming_util.GetOffchainUpdatesV1(offchainUpdates)

	// If not `DeliverTx`, then updates are optimistic. Stream them directly.
	if !lib.IsDeliverTxMode(ctx) {
		defer metrics.ModuleMeasureSince(
			metrics.FullNodeGrpc,
			metrics.GrpcSendOrderbookUpdatesLatency,
			time.Now(),
		)

		streamUpdates, clobPairIds := getStreamUpdatesFromOffchainUpdates(
			v1updates,
			lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
			ctx.ExecMode(),
		)
		sm.AddOrderUpdatesToCache(streamUpdates, clobPairIds)
		return
	}

	// If `DeliverTx`, updates should be staged to be streamed after consensus finalizes on a block.
	stagedEvent := clobtypes.StagedFinalizeBlockEvent{
		Event: &clobtypes.StagedFinalizeBlockEvent_OrderbookUpdate{
			OrderbookUpdate: &clobtypes.StreamOrderbookUpdate{
				Updates:  v1updates,
				Snapshot: false,
			},
		},
	}
	sm.finalizeBlockStager.StageFinalizeBlockEvent(
		ctx,
		&stagedEvent,
	)
}

func (sm *FullNodeStreamingManagerImpl) getStreamUpdatesForOrderbookFills(
	orderbookFills []clobtypes.StreamOrderbookFill,
	blockHeight uint32,
	execMode sdk.ExecMode,
	perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
) (
	streamUpdates []clobtypes.StreamUpdate,
	clobPairIds []uint32,
) {
	// Group fills by clob pair id.
	streamUpdates = make([]clobtypes.StreamUpdate, 0)
	clobPairIds = make([]uint32, 0)
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
	return streamUpdates, clobPairIds
}

// SendOrderbookFillUpdate groups fills by their clob pair ids and
// sends messages to the subscribers.
func (sm *FullNodeStreamingManagerImpl) SendOrderbookFillUpdate(
	orderbookFill clobtypes.StreamOrderbookFill,
	ctx sdk.Context,
	perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
) {
	// If not `DeliverTx`, then updates are optimistic. Stream them directly.
	if !lib.IsDeliverTxMode(ctx) {
		defer metrics.ModuleMeasureSince(
			metrics.FullNodeGrpc,
			metrics.GrpcSendOrderbookFillsLatency,
			time.Now(),
		)

		streamUpdates, clobPairIds := sm.getStreamUpdatesForOrderbookFills(
			[]clobtypes.StreamOrderbookFill{orderbookFill},
			lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
			ctx.ExecMode(),
			perpetualIdToClobPairId,
		)
		sm.AddOrderUpdatesToCache(streamUpdates, clobPairIds)
		return
	}

	// If `DeliverTx`, updates should be staged to be streamed after consensus finalizes on a block.
	stagedEvent := clobtypes.StagedFinalizeBlockEvent{
		Event: &clobtypes.StagedFinalizeBlockEvent_OrderFill{
			OrderFill: &orderbookFill,
		},
	}

	sm.finalizeBlockStager.StageFinalizeBlockEvent(
		ctx,
		&stagedEvent,
	)
}

// SendTakerOrderStatus sends out a taker order and its status to the full node streaming service.
func (sm *FullNodeStreamingManagerImpl) SendTakerOrderStatus(
	streamTakerOrder clobtypes.StreamTakerOrder,
	ctx sdk.Context,
) {
	// In current design, we never send this during `DeliverTx` (`FinalizeBlock`).
	lib.AssertCheckTxMode(ctx)

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
				BlockHeight: lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
				ExecMode:    uint32(ctx.ExecMode()),
			},
		},
		[]uint32{clobPairId},
	)
}

func getStreamUpdatesForSubaccountUpdates(
	subaccountUpdates []satypes.StreamSubaccountUpdate,
	blockHeight uint32,
	execMode sdk.ExecMode,
) (
	streamUpdates []clobtypes.StreamUpdate,
	subaccountIds []*satypes.SubaccountId,
) {
	// Group subaccount updates by subaccount id.
	streamUpdates = make([]clobtypes.StreamUpdate, 0)
	subaccountIds = make([]*satypes.SubaccountId, 0)
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
	return streamUpdates, subaccountIds
}

func getStreamUpdatesForPriceUpdates(
	priceUpdates []pricestypes.StreamPriceUpdate,
	blockHeight uint32,
	execMode sdk.ExecMode,
) (
	streamUpdates []clobtypes.StreamUpdate,
	marketIds []uint32,
) {
	// Group subaccount updates by subaccount id.
	streamUpdates = make([]clobtypes.StreamUpdate, 0)
	marketIds = make([]uint32, 0)
	for _, priceUpdate := range priceUpdates {
		streamUpdate := clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_PriceUpdate{
				PriceUpdate: &priceUpdate,
			},
			BlockHeight: blockHeight,
			ExecMode:    uint32(execMode),
		}
		streamUpdates = append(streamUpdates, streamUpdate)
		marketIds = append(marketIds, priceUpdate.MarketId)
	}
	return streamUpdates, marketIds
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

	sm.cacheStreamUpdatesByClobPairWithLock(updates, clobPairIds)

	sm.EmitMetrics()
	// Remove all subscriptions and wipe the buffer if buffer overflows.
	sm.RemoveSubscriptionsAndClearBufferIfFull()
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

	sm.cacheStreamUpdatesBySubaccountWithLock(updates, subaccountIds)

	sm.EmitMetrics()
	sm.RemoveSubscriptionsAndClearBufferIfFull()
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
		sm.EmitMetrics()
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
			metrics.IncrCounterWithLabels(
				metrics.GrpcAddToSubscriptionChannelCount,
				1,
				metrics.GetLabelForIntValue(metrics.SubscriptionId, int(id)),
			)
			select {
			case subscription.streamUpdatesChannel <- updates:
			default:
				// Buffer is full. Emit metric and drop subscription.
				sm.EmitMetrics()
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
		if alreadyInitialized := subscription.initializedWithSnapshot.Load(); alreadyInitialized {
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

func (sm *FullNodeStreamingManagerImpl) GetPriceSnapshotsForInitStreams(
	getPriceSnapshot func(marketId uint32) *pricestypes.StreamPriceUpdate,
) map[uint32]*pricestypes.StreamPriceUpdate {
	sm.Lock()
	defer sm.Unlock()

	ret := make(map[uint32]*pricestypes.StreamPriceUpdate)
	for _, subscription := range sm.orderbookSubscriptions {
		// If the subscription has been initialized, no need to grab the price snapshot.
		if alreadyInitialized := subscription.initializedWithSnapshot.Load(); alreadyInitialized {
			continue
		}

		for _, marketId := range subscription.marketIds {
			if _, exists := ret[marketId]; exists {
				continue
			}

			ret[marketId] = getPriceSnapshot(marketId)
		}
	}
	return ret
}

// cacheStreamUpdatesByClobPairWithLock adds stream updates to cache,
// and store corresponding clob pair Ids.
// This method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) cacheStreamUpdatesByClobPairWithLock(
	streamUpdates []clobtypes.StreamUpdate,
	clobPairIds []uint32,
) {
	sm.streamUpdateCache = append(sm.streamUpdateCache, streamUpdates...)
	allClobPairSubscriptionIds := sm.AllClobPairSubscriptionIds()
	for _, clobPairId := range clobPairIds {
		subscriptionIds := append(sm.clobPairIdToSubscriptionIdMapping[clobPairId], allClobPairSubscriptionIds...)
		sm.streamUpdateSubscriptionCache = append(
			sm.streamUpdateSubscriptionCache,
			subscriptionIds,
		)
	}
}

// cacheStreamUpdatesBySubaccountWithLock adds subaccount stream updates to cache,
// and store corresponding subaccount Ids.
// This method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) cacheStreamUpdatesBySubaccountWithLock(
	subaccountStreamUpdates []clobtypes.StreamUpdate,
	subaccountIds []*satypes.SubaccountId,
) {
	sm.streamUpdateCache = append(sm.streamUpdateCache, subaccountStreamUpdates...)
	for _, subaccountId := range subaccountIds {
		sm.streamUpdateSubscriptionCache = append(
			sm.streamUpdateSubscriptionCache,
			sm.subaccountIdToSubscriptionIdMapping[*subaccountId],
		)
	}
}

// cacheStreamUpdatesByMarketIdWithLock adds stream updates to cache,
// and store corresponding market ids.
// This method requires the lock and assumes that the lock has already been
// acquired by the caller.
func (sm *FullNodeStreamingManagerImpl) cacheStreamUpdatesByMarketIdWithLock(
	streamUpdates []clobtypes.StreamUpdate,
	marketIds []uint32,
) {
	if len(streamUpdates) != len(marketIds) {
		sm.logger.Error("Mismatch between stream updates and market IDs lengths")
		return
	}
	sm.streamUpdateCache = append(sm.streamUpdateCache, streamUpdates...)
	for _, marketId := range marketIds {
		sm.streamUpdateSubscriptionCache = append(
			sm.streamUpdateSubscriptionCache,
			sm.marketIdToSubscriptionIdMapping[marketId],
		)
	}
}

// Grpc Streaming logic after consensus agrees on a block.
// - Stream all events staged during `FinalizeBlock`.
// - Stream orderbook updates to sync fills in local ops queue.
func (sm *FullNodeStreamingManagerImpl) StreamBatchUpdatesAfterFinalizeBlock(
	ctx sdk.Context,
	orderBookUpdatesToSyncLocalOpsQueue *clobtypes.OffchainUpdates,
	perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
) {
	// Prevent gas metering from state read.
	ctx = ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())

	finalizedFills,
		finalizedSubaccountUpdates,
		finalizedOrderbookUpdates,
		finalizedPriceUpdates := sm.getStagedEventsFromFinalizeBlock(ctx)

	sm.Lock()
	defer sm.Unlock()

	// Flush all pending updates, since we want the onchain updates to arrive in a batch.
	sm.FlushStreamUpdatesWithLock()

	// Cache updates to sync local ops queue
	syncLocalUpdates, syncLocalClobPairIds := getStreamUpdatesFromOffchainUpdates(
		streaming_util.GetOffchainUpdatesV1(orderBookUpdatesToSyncLocalOpsQueue),
		lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)
	sm.cacheStreamUpdatesByClobPairWithLock(syncLocalUpdates, syncLocalClobPairIds)

	// Cache updates for finalized fills.
	fillStreamUpdates, fillClobPairIds := sm.getStreamUpdatesForOrderbookFills(
		finalizedFills,
		lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
		ctx.ExecMode(),
		perpetualIdToClobPairId,
	)
	sm.cacheStreamUpdatesByClobPairWithLock(fillStreamUpdates, fillClobPairIds)

	// Cache updates for finalized orderbook updates (e.g. RemoveOrderFillAmount in `EndBlocker`).
	for _, finalizedUpdate := range finalizedOrderbookUpdates {
		streamUpdates, clobPairIds := getStreamUpdatesFromOffchainUpdates(
			finalizedUpdate.Updates,
			lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
			ctx.ExecMode(),
		)
		sm.cacheStreamUpdatesByClobPairWithLock(streamUpdates, clobPairIds)
	}

	// Finally, cache updates for finalized subaccount updates
	subaccountStreamUpdates, subaccountIds := getStreamUpdatesForSubaccountUpdates(
		finalizedSubaccountUpdates,
		lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)
	sm.cacheStreamUpdatesBySubaccountWithLock(subaccountStreamUpdates, subaccountIds)

	// Finally, cache updates for finalized subaccount updates
	priceStreamUpdates, marketIds := getStreamUpdatesForPriceUpdates(
		finalizedPriceUpdates,
		lib.MustConvertIntegerToUint32(ctx.BlockHeight()),
		ctx.ExecMode(),
	)
	sm.cacheStreamUpdatesByMarketIdWithLock(priceStreamUpdates, marketIds)

	// Emit all stream updates in a single batch.
	// Note we still have the lock, which is released right before function returns.
	sm.FlushStreamUpdatesWithLock()
}

// getStagedEventsFromFinalizeBlock returns staged events from `FinalizeBlock`.
// It should be called after the consensus agrees on a block (e.g. Precommitter).
func (sm *FullNodeStreamingManagerImpl) getStagedEventsFromFinalizeBlock(
	ctx sdk.Context,
) (
	finalizedFills []clobtypes.StreamOrderbookFill,
	finalizedSubaccountUpdates []satypes.StreamSubaccountUpdate,
	finalizedOrderbookUpdates []clobtypes.StreamOrderbookUpdate,
	finalizedPriceUpdates []pricestypes.StreamPriceUpdate,
) {
	// Get onchain stream events stored in transient store.
	stagedEvents := sm.GetStagedFinalizeBlockEvents(ctx)

	metrics.SetGauge(
		metrics.GrpcStagedAllFinalizeBlockUpdatesCount,
		float32(len(stagedEvents)),
	)

	for _, stagedEvent := range stagedEvents {
		switch event := stagedEvent.Event.(type) {
		case *clobtypes.StagedFinalizeBlockEvent_OrderFill:
			finalizedFills = append(finalizedFills, *event.OrderFill)
		case *clobtypes.StagedFinalizeBlockEvent_SubaccountUpdate:
			finalizedSubaccountUpdates = append(finalizedSubaccountUpdates, *event.SubaccountUpdate)
		case *clobtypes.StagedFinalizeBlockEvent_OrderbookUpdate:
			finalizedOrderbookUpdates = append(finalizedOrderbookUpdates, *event.OrderbookUpdate)
		case *clobtypes.StagedFinalizeBlockEvent_PriceUpdate:
			finalizedPriceUpdates = append(finalizedPriceUpdates, *event.PriceUpdate)
		default:
			panic(fmt.Sprintf("Unhandled staged event type: %v\n", stagedEvent.Event))
		}
	}

	metrics.SetGauge(
		metrics.GrpcStagedSubaccountFinalizeBlockUpdatesCount,
		float32(len(finalizedSubaccountUpdates)),
	)
	metrics.SetGauge(
		metrics.GrpcStagedFillFinalizeBlockUpdatesCount,
		float32(len(finalizedFills)),
	)

	return finalizedFills, finalizedSubaccountUpdates, finalizedOrderbookUpdates, finalizedPriceUpdates
}

func (sm *FullNodeStreamingManagerImpl) InitializeNewStreams(
	getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
	subaccountSnapshots map[satypes.SubaccountId]*satypes.StreamSubaccountUpdate,
	pricesSnapshots map[uint32]*pricestypes.StreamPriceUpdate,
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
		if alreadyInitialized := subscription.initializedWithSnapshot.Swap(true); !alreadyInitialized {
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

			priceUpdates := []*pricestypes.StreamPriceUpdate{}
			for _, marketId := range subscription.marketIds {
				if priceUpdate, ok := pricesSnapshots[marketId]; ok {
					priceUpdates = append(priceUpdates, priceUpdate)
				} else {
					sm.logger.Error(
						fmt.Sprintf(
							"Price update not found for market id %v. This should not happen.",
							marketId,
						),
					)
				}
			}

			sm.SendCombinedSnapshot(
				allUpdates,
				saUpdates,
				priceUpdates,
				subscriptionId,
				blockHeight,
				execMode,
			)

			if sm.snapshotBlockInterval != 0 {
				subscription.nextSnapshotBlock = blockHeight + sm.snapshotBlockInterval
			}
		}

		// If the snapshot block interval is enabled and the next block is a snapshot block,
		// reset the `atomic.Bool` so snapshots are sent for the next block.
		if sm.snapshotBlockInterval > 0 &&
			blockHeight+1 == subscription.nextSnapshotBlock {
			subscription.initializedWithSnapshot = &atomic.Bool{} // False by default.
		}
	}
}
