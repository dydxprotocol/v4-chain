package msgsender

// Ensure the `IndexerMessageSender` interface is implemented at compile time.
var _ IndexerMessageSender = (*IndexerMessageSenderNoop)(nil)

// No-op implementation of the IndexerMessageSender interface.
// Will be used in tests or when the V4 application is not connected to an Indexer.
type IndexerMessageSenderNoop struct{}

func NewIndexerMessageSenderNoop() *IndexerMessageSenderNoop {
	return &IndexerMessageSenderNoop{}
}

func (msgSender *IndexerMessageSenderNoop) Enabled() bool {
	return false
}

func (msgSender *IndexerMessageSenderNoop) SendOnchainData(message Message) {}

func (msgSender *IndexerMessageSenderNoop) SendOffchainData(message Message) {}

func (msgSender *IndexerMessageSenderNoop) Close() error {
	return nil
}
