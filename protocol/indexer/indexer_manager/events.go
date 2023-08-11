package indexer_manager

import (
	"encoding/base64"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	ante_types "github.com/dydxprotocol/v4/app/ante/types"
	"github.com/dydxprotocol/v4/indexer/common"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
)

const (
	// TransientStoreKey is the transient store key for indexer events.
	TransientStoreKey = "transient_indexer_events"

	// IndexerEventsKey is the key to retrieve the indexer events
	// within the last block.
	IndexerEventsKey = "IndexerEvents"

	ModuleName = "indexer_events"
)

func getIndexerEvents(ctx sdk.Context, storeKey storetypes.StoreKey) []*IndexerTendermintEventWrapper {
	// This is necessary to prevent GasConsumed from being incremented when indexer events are recorded.
	// Without this, consensus failure would occur due to lastResultsHash mismatch from different gas costs.
	noGasCtx := ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())
	store := noGasCtx.TransientStore(storeKey)
	indexerEventsSliceBytes := store.Get([]byte(IndexerEventsKey))
	if indexerEventsSliceBytes == nil {
		return []*IndexerTendermintEventWrapper{}
	}
	var events IndexerEventsStoreValue
	unmarshaler := &common.UnmarshalerImpl{}
	err := unmarshaler.Unmarshal(indexerEventsSliceBytes, &events)
	if err != nil {
		panic(err)
	}
	return events.Events
}

// GetB64EncodedEventMessage returns the base64 encoded event message.
// TODO(DEC-1720): Deprecate this function once we change the underlying proto to use bytes.
func GetB64EncodedEventMessage(
	eventMessage proto.Message,
) string {
	marshaler := &common.MarshalerImpl{}
	eventMessageBytes, err := marshaler.Marshal(eventMessage)
	if err != nil {
		panic(err)
	}

	b64encodedEventMessage := base64.StdEncoding.EncodeToString(eventMessageBytes)
	return b64encodedEventMessage
}

// GetBytesFromEventData returns the decoded bytes of the base64 event data string.
func GetBytesFromEventData(
	event string,
) []byte {
	bytes, err := base64.StdEncoding.DecodeString(event)
	if err != nil {
		panic(err)
	}
	return bytes
}

// addTxnEvent adds a transaction event to the context's transient store of indexer events.
func addTxnEvent(
	ctx sdk.Context,
	subType string,
	data string,
	storeKey storetypes.StoreKey,
) {
	event := IndexerTendermintEventWrapper{
		Event: &IndexerTendermintEvent{
			Subtype:             subType,
			Data:                data,
			OrderingWithinBlock: &IndexerTendermintEvent_TransactionIndex{},
		},
		TxnHash: string(lib.GetTxHash(ctx.TxBytes())),
	}
	addEvent(ctx, event, storeKey)
}

// addBlockEvent adds a block event to the context's transient store of indexer events.
func addBlockEvent(
	ctx sdk.Context,
	subType string,
	data string,
	storeKey storetypes.StoreKey,
	blockEvent IndexerTendermintEvent_BlockEvent,
) {
	event := IndexerTendermintEventWrapper{
		Event: &IndexerTendermintEvent{
			Subtype: subType,
			Data:    data,
			OrderingWithinBlock: &IndexerTendermintEvent_BlockEvent_{
				BlockEvent: blockEvent,
			},
		},
	}
	addEvent(ctx, event, storeKey)
}

// addEvent adds an event to the context's transient store of indexer events.
func addEvent(
	ctx sdk.Context,
	event IndexerTendermintEventWrapper,
	storeKey storetypes.StoreKey,
) {
	noGasCtx := ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())
	indexerEvents := getIndexerEvents(noGasCtx, storeKey)
	indexerEvents = append(indexerEvents, &event)
	newEventsValue := IndexerEventsStoreValue{
		Events: indexerEvents,
	}
	marshaler := &common.MarshalerImpl{}
	newEventsValueBytes, err := marshaler.Marshal(&newEventsValue)
	if err != nil {
		panic(err)
	}
	store := noGasCtx.TransientStore(storeKey)
	store.Set([]byte(IndexerEventsKey), newEventsValueBytes)
}

// produceBlock returns the block. It should only be called in EndBlocker when the
// transient store contains all onchain events from a ready-to-be-committed block.
func produceBlock(ctx sdk.Context, storeKey storetypes.StoreKey) *IndexerTendermintBlock {
	noGasCtx := ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())
	txHashes := []string{}
	txEventsMap := make(map[string][]*IndexerTendermintEvent)
	blockEvents := []*IndexerTendermintEvent{}
	blockHeight := lib.MustConvertIntegerToUint32(noGasCtx.BlockHeight())
	blockTime := noGasCtx.BlockTime()
	events := getIndexerEvents(noGasCtx, storeKey)

	for _, event := range events {
		switch event.Event.OrderingWithinBlock.(type) {
		case *IndexerTendermintEvent_BlockEvent_:
			blockEvents = append(blockEvents, event.Event)
		case *IndexerTendermintEvent_TransactionIndex:
			txHash := event.TxnHash
			if txEvents, ok := txEventsMap[txHash]; ok {
				txEventsMap[txHash] = append(txEvents, event.Event)
			} else {
				txHashes = append(txHashes, txHash)
				txEventsMap[txHash] = []*IndexerTendermintEvent{event.Event}
			}
		}
	}
	// create map from txHash to index
	txHashesMap := make(map[string]int)
	for i, txHash := range txHashes {
		txHashesMap[txHash] = i
	}
	// iterate through txEventsMap and add transaction/event indices to each event
	numTxnEvents := 0
	for txHash, events := range txEventsMap {
		for i, event := range events {
			event.OrderingWithinBlock = &IndexerTendermintEvent_TransactionIndex{
				TransactionIndex: uint32(txHashesMap[txHash]),
			}
			event.EventIndex = uint32(i)
			events[i] = event
			numTxnEvents++
		}
		txEventsMap[txHash] = events
	}
	// build list of all events
	allEvents := make([]*IndexerTendermintEvent, 0, numTxnEvents+len(blockEvents))
	for _, txHash := range txHashes {
		allEvents = append(allEvents, txEventsMap[txHash]...)
	}
	// append block events
	allEvents = append(allEvents, blockEvents...)
	recordMetrics(numTxnEvents, len(blockEvents))

	return &IndexerTendermintBlock{
		Height:   blockHeight,
		Time:     blockTime,
		Events:   allEvents,
		TxHashes: txHashes,
	}
}

func recordMetrics(
	totalNumTxnEvents int,
	totalNumBlockEvents int,
) {
	telemetry.SetGauge(
		float32(totalNumTxnEvents),
		ModuleName,
		metrics.TotalNumIndexerTxnEvents,
	)
	telemetry.SetGauge(
		float32(totalNumBlockEvents),
		ModuleName,
		metrics.TotalNumIndexerBlockEvents,
	)
}
