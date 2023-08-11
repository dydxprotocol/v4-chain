package indexer_manager

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/indexer/msgsender"
)

type IndexerEventManager interface {
	Enabled() bool
	AddTxnEvent(ctx sdk.Context, subType string, data string)
	SendOffchainData(message msgsender.Message)
	ProduceBlock(ctx sdk.Context) *IndexerTendermintBlock
	AddBlockEvent(ctx sdk.Context, subType string, data string, blockEvent IndexerTendermintEvent_BlockEvent)
}

type IndexerEventManagerImpl struct {
	indexerMessageSender           msgsender.IndexerMessageSender
	indexerEventsTransientStoreKey storetypes.StoreKey
}

func NewIndexerEventManager(
	indexerMessageSender msgsender.IndexerMessageSender,
	indexerEventsTransientStoreKey storetypes.StoreKey,
) *IndexerEventManagerImpl {
	return &IndexerEventManagerImpl{
		indexerMessageSender:           indexerMessageSender,
		indexerEventsTransientStoreKey: indexerEventsTransientStoreKey,
	}
}

func (i *IndexerEventManagerImpl) Enabled() bool {
	return i.indexerMessageSender.Enabled()
}

func (i *IndexerEventManagerImpl) GetIndexerEventsTransientStoreKey() storetypes.StoreKey {
	return i.indexerEventsTransientStoreKey
}

func (i *IndexerEventManagerImpl) SendOffchainData(message msgsender.Message) {
	if i.indexerMessageSender.Enabled() {
		i.indexerMessageSender.SendOffchainData(message)
	}
}

func (i *IndexerEventManagerImpl) SendOnchainData(block *IndexerTendermintBlock) {
	if i.indexerMessageSender.Enabled() {
		message := CreateIndexerBlockEventMessage(block)
		i.indexerMessageSender.SendOnchainData(message)
	}
}

// AddTxnEvent adds a transaction event to the context's transient store of indexer events.
func (i *IndexerEventManagerImpl) AddTxnEvent(
	ctx sdk.Context,
	subType string,
	data string,
) {
	if i.indexerMessageSender.Enabled() {
		addTxnEvent(ctx, subType, data, i.indexerEventsTransientStoreKey)
	}
}

// AddBlockEvent adds a block event to the context's transient store of indexer events.
func (i *IndexerEventManagerImpl) AddBlockEvent(
	ctx sdk.Context,
	subType string,
	data string,
	blockEvent IndexerTendermintEvent_BlockEvent,
) {
	if i.indexerMessageSender.Enabled() {
		addBlockEvent(ctx, subType, data, i.indexerEventsTransientStoreKey, blockEvent)
	}
}

// ProduceBlock returns an `IndexerTendermintBlock` containing all the indexer events in the block.
// It should only be called in EndBlocker when the transient store contains all onchain events from
// a ready-to-be-committed block.
func (i *IndexerEventManagerImpl) ProduceBlock(
	ctx sdk.Context,
) *IndexerTendermintBlock {
	if i.indexerMessageSender.Enabled() {
		return produceBlock(ctx, i.indexerEventsTransientStoreKey)
	}
	return nil
}
