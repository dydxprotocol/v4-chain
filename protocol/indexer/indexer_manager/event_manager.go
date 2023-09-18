package indexer_manager

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
)

type IndexerEventManager interface {
	Enabled() bool
	AddTxnEvent(ctx sdk.Context, subType string, data string, version uint32)
	SendOffchainData(message msgsender.Message)
	SendOnchainData(block *IndexerTendermintBlock)
	ProduceBlock(ctx sdk.Context) *IndexerTendermintBlock
	AddBlockEvent(
		ctx sdk.Context,
		subType string,
		data string,
		blockEvent IndexerTendermintEvent_BlockEvent,
		version uint32,
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
	if i.indexerMessageSender.Enabled() && i.sendOffchainData {
		i.indexerMessageSender.SendOffchainData(message)
	}
}

func (i *indexerEventManagerImpl) SendOnchainData(block *IndexerTendermintBlock) {
	if i.indexerMessageSender.Enabled() {
		message := CreateIndexerBlockEventMessage(block)
		i.indexerMessageSender.SendOnchainData(message)
	}
}

// AddTxnEvent adds a transaction event to the context's transient store of indexer events.
func (i *indexerEventManagerImpl) AddTxnEvent(
	ctx sdk.Context,
	subType string,
	data string,
	version uint32,
) {
	if i.indexerMessageSender.Enabled() {
		addTxnEvent(ctx, subType, data, version, i.indexerEventsTransientStoreKey)
	}
}

// ClearEvents clears all events in the context's transient store of indexer events.
func (i *indexerEventManagerImpl) ClearEvents(
	ctx sdk.Context,
) {
	if i.indexerMessageSender.Enabled() {
		clearEvents(ctx, i.indexerEventsTransientStoreKey)
	}
}

// AddBlockEvent adds a block event to the context's transient store of indexer events.
func (i *indexerEventManagerImpl) AddBlockEvent(
	ctx sdk.Context,
	subType string,
	data string,
	blockEvent IndexerTendermintEvent_BlockEvent,
	version uint32,
) {
	if i.indexerMessageSender.Enabled() {
		addBlockEvent(ctx, subType, data, i.indexerEventsTransientStoreKey, blockEvent, version)
	}
}

// ProduceBlock returns an `IndexerTendermintBlock` containing all the indexer events in the block.
// It should only be called in EndBlocker when the transient store contains all onchain events from
// a ready-to-be-committed block.
func (i *indexerEventManagerImpl) ProduceBlock(
	ctx sdk.Context,
) *IndexerTendermintBlock {
	if i.indexerMessageSender.Enabled() {
		return produceBlock(ctx, i.indexerEventsTransientStoreKey)
	}
	return nil
}
