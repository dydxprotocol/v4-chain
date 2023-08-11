package indexer_manager

import "github.com/dydxprotocol/v4/indexer/msgsender"

func NewIndexerEventManagerNoop() *IndexerEventManagerImpl {
	return NewIndexerEventManager(
		msgsender.NewIndexerMessageSenderNoop(),
		nil,
	)
}
