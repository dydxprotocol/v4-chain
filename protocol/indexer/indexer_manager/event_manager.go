package indexer_manager

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ante_types "github.com/dydxprotocol/v4-chain/protocol/app/ante/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
)

type IndexerEventManager interface {
	Enabled() bool
	AddTxnEvent(ctx sdk.Context, subType string, version uint32, dataByes []byte)
	AddOnchainStreamEvent(
		ctx sdk.Context,
		dataBytes []byte,
	)
	SendOffchainData(message msgsender.Message)
	SendOnchainData(block *IndexerTendermintBlock)
	ProduceBlock(ctx sdk.Context) *IndexerTendermintBlock
	AddBlockEvent(
		ctx sdk.Context,
		subType string,
		blockEvent IndexerTendermintEvent_BlockEvent,
		version uint32,
		dataBytes []byte,
	)
	GetOnchainStreamEvents(ctx sdk.Context) []*IndexerTendermintEventWrapper
	ClearEvents(ctx sdk.Context)
}

// Ensure the `IndexerEventManager` interface is implemented at compile time.
var _ IndexerEventManager = (*indexerEventManagerImpl)(nil)

type indexerEventManagerImpl struct {
	indexerMessageSender           msgsender.IndexerMessageSender
	indexerEventsTransientStoreKey storetypes.StoreKey
	// stores the onchain events for full node streaming.
	onchainStreamTransientStoreKey storetypes.StoreKey
	sendOffchainData               bool
}

func NewIndexerEventManager(
	indexerMessageSender msgsender.IndexerMessageSender,
	indexerEventsTransientStoreKey storetypes.StoreKey,
	onchainStreamTransientStoreKey storetypes.StoreKey,
	sendOffchainData bool,
) IndexerEventManager {
	return &indexerEventManagerImpl{
		indexerMessageSender:           indexerMessageSender,
		indexerEventsTransientStoreKey: indexerEventsTransientStoreKey,
		onchainStreamTransientStoreKey: onchainStreamTransientStoreKey,
		sendOffchainData:               sendOffchainData,
	}
}

func (i *indexerEventManagerImpl) Enabled() bool {
	return i.indexerMessageSender.Enabled()
}

func (i *indexerEventManagerImpl) GetIndexerEventsTransientStoreKey() storetypes.StoreKey {
	return i.indexerEventsTransientStoreKey
}

func (i *indexerEventManagerImpl) SendOffchainData(message msgsender.Message) {
	if i.Enabled() && i.sendOffchainData {
		i.indexerMessageSender.SendOffchainData(message)
	}
}

func (i *indexerEventManagerImpl) SendOnchainData(block *IndexerTendermintBlock) {
	if i.Enabled() {
		message := CreateIndexerBlockEventMessage(block)
		i.indexerMessageSender.SendOnchainData(message)
	}
}

// AddOnchainStreamEvent adds a onchain stream event to the context's transient store.
// TODO(CT-939): Consolidate FNS with indexer events; share storage and event emission logic.
func (i *indexerEventManagerImpl) AddOnchainStreamEvent(
	ctx sdk.Context,
	dataBytes []byte,
) {
	event := IndexerTendermintEventWrapper{
		Event: &IndexerTendermintEvent{
			DataBytes: dataBytes,
		},
	}
	addEvent(ctx, event, i.onchainStreamTransientStoreKey)
}

// AddTxnEvent adds a transaction event to the context's transient store of indexer events.
func (i *indexerEventManagerImpl) AddTxnEvent(
	ctx sdk.Context,
	subType string,
	version uint32,
	dataBytes []byte,
) {
	if i.Enabled() {
		addTxnEvent(ctx, subType, version, i.indexerEventsTransientStoreKey, dataBytes)
	}
}

// ClearEvents clears all events in the context's transient store of indexer events.
func (i *indexerEventManagerImpl) ClearEvents(
	ctx sdk.Context,
) {
	if i.Enabled() {
		clearEvents(ctx, i.indexerEventsTransientStoreKey)
	}
}

// AddBlockEvent adds a block event to the context's transient store of indexer events.
func (i *indexerEventManagerImpl) AddBlockEvent(
	ctx sdk.Context,
	subType string,
	blockEvent IndexerTendermintEvent_BlockEvent,
	version uint32,
	dataBytes []byte,
) {
	if i.Enabled() {
		addBlockEvent(ctx, subType, i.indexerEventsTransientStoreKey, blockEvent, version, dataBytes)
	}
}

// ProduceBlock returns an `IndexerTendermintBlock` containing all the indexer events in the block.
// It should only be called in EndBlocker when the transient store contains all onchain events from
// a ready-to-be-committed block.
func (i *indexerEventManagerImpl) ProduceBlock(
	ctx sdk.Context,
) *IndexerTendermintBlock {
	if i.Enabled() {
		return produceBlock(ctx, i.indexerEventsTransientStoreKey)
	}
	return nil
}

// GetOnchainStreamEvents returns all onchain stream events (fills) stored in the transient store.
func (i *indexerEventManagerImpl) GetOnchainStreamEvents(ctx sdk.Context) []*IndexerTendermintEventWrapper {
	noGasCtx := ctx.WithGasMeter(ante_types.NewFreeInfiniteGasMeter())
	return getIndexerEvents(noGasCtx, i.onchainStreamTransientStoreKey)
}
