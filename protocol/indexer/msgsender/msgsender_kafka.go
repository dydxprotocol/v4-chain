package msgsender

import (
	"sync"
	"time"

	"cosmossdk.io/log"
	"github.com/Shopify/sarama"
	"github.com/burdiyan/kafkautil"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

// Ensure the `IndexerMessageSender` interface is implemented at compile time.
var _ IndexerMessageSender = (*IndexerMessageSenderKafka)(nil)

// Implementation of the IndexerMessageSender interface that sends data to Kafka.
// Will be used when the V4 application is connected to an Indexer.
// NOTE: This struct is go-routine safe. Messages are sent by writing to a single-channel, and a
// mutex and boolean variable is used to ensure `Close` only closes the underlying Kafka producer
// once.
type IndexerMessageSenderKafka struct {
	mutex      sync.Mutex
	closed     bool
	inputsDone sync.WaitGroup
	producer   sarama.AsyncProducer
	logger     log.Logger
	successes  int
	errors     int
}

func NewIndexerMessageSenderKafka(
	indexerFlags indexer.IndexerFlags,
	config *sarama.Config,
	logger log.Logger,
) (*IndexerMessageSenderKafka, error) {
	if config == nil {
		config = sarama.NewConfig()
	}

	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = indexerFlags.MaxRetries
	// max retry should be set so that max retry * retry backoff > Zookeeper session.timeout + some buffer
	config.Producer.Retry.Backoff = 1000 * time.Millisecond
	config.Producer.MaxMessageBytes = 4194304 // 4MB
	config.Producer.RequiredAcks = sarama.WaitForAll
	// Use the JVM compatible parititoner to match `kafkajs` which is used in the indexer services.
	config.Producer.Partitioner = kafkautil.NewJVMCompatiblePartitioner
	producer, err := sarama.NewAsyncProducer(indexerFlags.KafkaAddrs, config)

	if err != nil {
		return nil, err
	}

	sender := NewIndexerMessageSenderKafkaWithProducer(producer, logger)

	return sender, nil
}

func NewIndexerMessageSenderKafkaWithProducer(
	producer sarama.AsyncProducer,
	logger log.Logger,
) *IndexerMessageSenderKafka {
	sender := &IndexerMessageSenderKafka{
		inputsDone: sync.WaitGroup{},
		closed:     false,
		producer:   producer,
		logger:     logger,
		successes:  0,
		errors:     0,
	}
	// The wait group waits for successes and errors which is why it is 2.
	sender.inputsDone.Add(2)
	go sender.handleSuccesses()
	go sender.handleErrors()

	return sender
}

func (msgSender *IndexerMessageSenderKafka) Enabled() bool {
	return true
}

// SendOnchainData sends a key/value pair of byte slices to the on-chain data kafka topic.
// This method is go-routine safe.
func (msgSender *IndexerMessageSenderKafka) SendOnchainData(message Message) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.SendOnchainData,
		metrics.Latency,
	)

	value := sarama.ByteEncoder(message.Value)
	telemetry.SetGauge(float32(value.Length()), types.ModuleName, metrics.OnchainMessageLength)
	msgSender.send(&sarama.ProducerMessage{
		Topic:   ON_CHAIN_KAFKA_TOPIC,
		Key:     sarama.ByteEncoder(message.Key),
		Value:   value,
		Headers: message.Headers,
	})
}

// SendOffchainData sends a key/value pair of byte slices to the off-chain data kafka topic.
// This method is go-routine safe.
func (msgSender *IndexerMessageSenderKafka) SendOffchainData(message Message) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.SendOffchainData,
		metrics.Latency,
	)

	value := sarama.ByteEncoder(message.Value)
	telemetry.SetGauge(float32(value.Length()), types.ModuleName, metrics.OffchainMessageLength)
	msgSender.send(&sarama.ProducerMessage{
		Topic:   OFF_CHAIN_KAFKA_TOPIC,
		Key:     sarama.ByteEncoder(message.Key),
		Value:   value,
		Headers: message.Headers,
	})
}

// send sends a message to Kafka. This method is go-routine safe.
func (msgSender *IndexerMessageSenderKafka) send(message *sarama.ProducerMessage) {
	msgSender.mutex.Lock()
	defer msgSender.mutex.Unlock()
	if msgSender.closed {
		msgSender.logger.Error("Cannot send to a closed IndexerMessageSenderKafka.")
		return
	}

	msgSender.producer.Input() <- message
}

// Close closes the underlying `AsyncProducer` and waits for all errors/success messages to be
// processed before returning.
func (msgSender *IndexerMessageSenderKafka) Close() error {
	// Lock to make this function go-routine safe.
	msgSender.mutex.Lock()
	defer msgSender.mutex.Unlock()

	// Ensure that the producer is only closed once.
	if msgSender.closed {
		return ErrKafkaAlreadyClosed
	}

	err := msgSender.producer.Close()
	if err != nil {
		return err
	}

	// Wait for success and error messages from the `AsyncProducer` to finish processing before
	// returning. Each goroutine will signal the channel.
	msgSender.inputsDone.Wait()
	msgSender.closed = true

	return nil
}

// handleSuccesses reads messages from the success channel of the `AsyncProducer`
// This is required so that the producer will not deadlock due to the channel becoming full.
func (msgSender *IndexerMessageSenderKafka) handleSuccesses() {
	c := msgSender.producer.Successes()
	for {
		_, ok := <-c
		if !ok {
			msgSender.inputsDone.Done()
			return
		}
		msgSender.successes = msgSender.successes + 1
		telemetry.IncrCounter(1, types.ModuleName, metrics.MessageSendSuccess)
	}
}

// handleErrors reads messages from the error channel of the `AsyncProducer`
// This is required so that the producer will not deadlock due to the channel becoming full.
func (msgSender *IndexerMessageSenderKafka) handleErrors() {
	c := msgSender.producer.Errors()
	for {
		err, ok := <-c
		if !ok {
			msgSender.inputsDone.Done()
			return
		}
		msgSender.logger.Error(
			"Failed to deliver message to Indexer",
			"message",
			err.Msg,
			"error",
			err.Err,
		)
		msgSender.errors = msgSender.errors + 1
		telemetry.IncrCounter(1, types.ModuleName, metrics.MessageSendError)
	}
}
