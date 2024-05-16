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
	logger log.Logger
	sync.Mutex

	// orderbookSubscriptions maps subscription IDs to their respective orderbook subscriptions.
	orderbookSubscriptions map[uint32]*OrderbookSubscription
	nextSubscriptionId     uint32

	// Readonly buffer to enqueue orderbook updates before pushing them through grpc streams.
	// Decouples the execution of abci logic with full node streaming.
	updateBuffer chan bufferInternalResponse
}

// bufferInternalResponse is enqueued into the readonly buffer.
// It contains an update respnose and the clob pair id to send this information to.
type bufferInternalResponse struct {
	response clobtypes.StreamOrderbookUpdatesResponse

	// Information relevant to which Orderbook Subscription to send out to
	clobPairId uint32
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
	bufferWindow uint32,
) *GrpcStreamingManagerImpl {
	grpcStreamingManager := &GrpcStreamingManagerImpl{
		logger:                 logger.With("module", "grpc-streaming"),
		orderbookSubscriptions: make(map[uint32]*OrderbookSubscription),
		nextSubscriptionId:     0,

		updateBuffer: make(chan bufferInternalResponse, bufferWindow),
	}

	// Worker goroutine to consistently read from channel and send out updates
	go func() {
		for internalResponse := range grpcStreamingManager.updateBuffer {
			grpcStreamingManager.sendUpdateResponse(internalResponse)
		}
	}()

	return grpcStreamingManager
}

func (sm *GrpcStreamingManagerImpl) Enabled() bool {
	return true
}

func (sm *GrpcStreamingManagerImpl) Stop() {
	close(sm.updateBuffer)
}

func (sm *GrpcStreamingManagerImpl) sendUpdateResponse(
	internalResponse bufferInternalResponse,
) {
	// Send update to subscribers.
	subscriptionIdsToRemove := make([]uint32, 0)

	for id, subscription := range sm.orderbookSubscriptions {
		for _, clobPairId := range subscription.clobPairIds {
			if clobPairId == internalResponse.clobPairId {
				if err := subscription.srv.Send(
					&internalResponse.response,
				); err != nil {
					subscriptionIdsToRemove = append(subscriptionIdsToRemove, id)
				}
			}
		}
	}
	// Clean up subscriptions that have been closed.
	// If a Send update has failed for any clob pair id, the whole subscription will be removed.
	for _, id := range subscriptionIdsToRemove {
		delete(sm.orderbookSubscriptions, id)
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

// SendOrderbookFillUpdates groups fills by their clob pair ids and
// enqueues messages to be sent to the subscribers.
func (sm *GrpcStreamingManagerImpl) SendOrderbookFillUpdates(
	orderbookFills []clobtypes.StreamOrderbookFill,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
	defer metrics.ModuleMeasureSince(
		metrics.FullNodeGrpc,
		metrics.GrpcSendOrderbookFillsLatency,
		time.Now(),
	)
	sm.Lock()
	defer sm.Unlock()

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
		}
		updatesByClobPairId[clobPairId] = append(updatesByClobPairId[clobPairId], streamUpdate)
	}

	// Send response updates into the stream buffer
	for clobPairId, streamUpdates := range updatesByClobPairId {
		streamResponse := clobtypes.StreamOrderbookUpdatesResponse{
			Updates:     streamUpdates,
			BlockHeight: blockHeight,
			ExecMode:    uint32(execMode),
		}

		sm.mustEnqueueOrderbookUpdate(bufferInternalResponse{
			response:   streamResponse,
			clobPairId: clobPairId,
		})
	}
}

// SendOrderbookUpdates groups updates by their clob pair ids and
// enqueues messages to be sent to the subscribers.
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
	sm.Lock()
	defer sm.Unlock()

	// Group updates by clob pair id.
	updatesByClobPairId := make(map[uint32]*clobtypes.OffchainUpdates)
	for _, message := range offchainUpdates.Messages {
		clobPairId := message.OrderId.ClobPairId
		if _, ok := updatesByClobPairId[clobPairId]; !ok {
			updatesByClobPairId[clobPairId] = clobtypes.NewOffchainUpdates()
		}
		updatesByClobPairId[clobPairId].Messages = append(updatesByClobPairId[clobPairId].Messages, message)
	}

	// Unmarshal messages to v1 updates and enqueue in buffer to be sent.
	for clobPairId, update := range updatesByClobPairId {
		v1updates, err := GetOffchainUpdatesV1(update)
		if err != nil {
			panic(err)
		}
		streamUpdate := clobtypes.StreamUpdate{
			UpdateMessage: &clobtypes.StreamUpdate_OrderbookUpdate{
				OrderbookUpdate: &clobtypes.StreamOrderbookUpdate{
					Updates:  v1updates,
					Snapshot: snapshot,
				},
			},
		}
		sm.mustEnqueueOrderbookUpdate(bufferInternalResponse{
			response: clobtypes.StreamOrderbookUpdatesResponse{
				Updates:     []clobtypes.StreamUpdate{streamUpdate},
				BlockHeight: blockHeight,
				ExecMode:    uint32(execMode),
			},
			clobPairId: clobPairId,
		})
	}
}

// mustEnqueueOrderbookUpdate tries to enqueue an orderbook update to the buffer via non-blocking send.
// If the buffer is full, *all* streaming subscriptions will be shut down.
func (sm *GrpcStreamingManagerImpl) mustEnqueueOrderbookUpdate(internalResponse bufferInternalResponse) {
	select {
	case sm.updateBuffer <- internalResponse:
		return
	default:
		sm.logger.Info("GRPC Streaming buffer full. Clearing all subscriptions")
		for k := range sm.orderbookSubscriptions {
			delete(sm.orderbookSubscriptions, k)
		}
	}
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
