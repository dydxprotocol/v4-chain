package msgsender

import "sync"

// Ensure the `IndexerMessageSender` interface is implemented at compile time.
var _ IndexerMessageSender = (*IndexerMessageSenderInMemoryCollector)(nil)

func NewIndexerMessageSenderInMemoryCollector() *IndexerMessageSenderInMemoryCollector {
	return &IndexerMessageSenderInMemoryCollector{}
}

type IndexerMessageSenderInMemoryCollector struct {
	mutex            sync.Mutex
	onchainMessages  []Message
	offchainMessages []Message
	closed           bool
}

func (i *IndexerMessageSenderInMemoryCollector) Enabled() bool {
	return true
}

func (i *IndexerMessageSenderInMemoryCollector) SendOnchainData(message Message) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.closed {
		return
	}
	i.onchainMessages = append(i.onchainMessages, message)
}

func (i *IndexerMessageSenderInMemoryCollector) SendOffchainData(message Message) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.closed {
		return
	}
	i.offchainMessages = append(i.offchainMessages, message)
}

func (i *IndexerMessageSenderInMemoryCollector) Close() error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.closed = true
	return nil
}

func (i *IndexerMessageSenderInMemoryCollector) Clear() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.offchainMessages = i.offchainMessages[:0]
	i.onchainMessages = i.onchainMessages[:0]
}

func (i *IndexerMessageSenderInMemoryCollector) GetOffchainMessages() []Message {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.offchainMessages
}

func (i *IndexerMessageSenderInMemoryCollector) GetOnchainMessages() []Message {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.onchainMessages
}
