package msgsender

import (
	"github.com/Shopify/sarama"
)

var TransactionHashHeaderKey = []byte("TransactionHash")

// Message is a key/value pair of byte slices that can be sent via the send functions in the
// IndexerMessageSender.
type Message struct {
	Key     []byte
	Value   []byte
	Headers []sarama.RecordHeader
}

// Message header is a key/value pair of byte slices that can be added to a Message to be sent
// along with the main key/value pair of the Message. This is converted to a `sarama.RecordHeader`
// in the `Message`.
type MessageHeader struct {
	Key   []byte
	Value []byte
}

// IndexerMessageSender is an interface that exposes methods to send messages to the
// on-chain/off-chain data archival services in the Indexer.
// The `Enabled` function is used to determine if any additional computations needed to generate
// Indexer-specific data should be run in various modules.
type IndexerMessageSender interface {
	Enabled() bool // whether the IndexerMessageSender will send messages to the Indexer
	SendOnchainData(message Message)
	SendOffchainData(message Message)
	Close() error
}

// AddHeader adds a `RecordHeader` to a `Message`. If there are already existing headers in the
// `Message`, the new header will be appended to the slice of existing headers.
func (msg Message) AddHeader(header MessageHeader) Message {
	return Message{
		Key:   msg.Key,
		Value: msg.Value,
		Headers: append(
			msg.Headers,
			sarama.RecordHeader{
				Key:   header.Key,
				Value: header.Value,
			},
		),
	}
}
