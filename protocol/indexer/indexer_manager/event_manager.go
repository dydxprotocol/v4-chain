package indexer_manager

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
)

type IndexerEventManager interface {
	Enabled() bool
	AddTxnEvent(ctx sdk.Context, subType string, version uint32, dataByes []byte)
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
	ClearEvents(ctx sdk.Context)
}

// Ensure the `IndexerEventManager` interface is implemented at compile time.
var _ IndexerEventManager = (*indexerEventManagerImpl)(nil)

type indexerEventManagerImpl struct {
	indexerMessageSender           msgsender.IndexerMessageSender
	indexerEventsTransientStoreKey storetypes.StoreKey
	sendOffchainData               bool
}

func NewIndexerEventManager(
	indexerMessageSender msgsender.IndexerMessageSender,
	indexerEventsTransientStoreKey storetypes.StoreKey,
	sendOffchainData bool,
) IndexerEventManager {
	return &indexerEventManagerImpl{
		indexerMessageSender:           indexerMessageSender,
		indexerEventsTransientStoreKey: indexerEventsTransientStoreKey,
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
