package msgsender

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/stretchr/testify/require"
)

const (
	msgKey   = "messageKey"
	msgValue = "messageValue"
)

var (
	produceError = fmt.Errorf("ProduceError")
)

func getMockProducer(t *testing.T, topic string, successes int, errors int) sarama.AsyncProducer {
	mockConfig := mocks.NewTestConfig()
	mockConfig.Producer.Return.Successes = true
	mockProducer := mocks.NewAsyncProducer(t, mockConfig)
	for i := 0; i < successes; i++ {
		mockProducer.ExpectInputWithMessageCheckerFunctionAndSucceed(
			func(msg *sarama.ProducerMessage) error {
				require.Equal(t, topic, msg.Topic)
				require.Equal(t, sarama.ByteEncoder(msgKey), msg.Key)
				require.Equal(t, sarama.ByteEncoder(msgValue), msg.Value)

				return nil
			},
		)
	}

	for i := 0; i < errors; i++ {
		mockProducer.ExpectInputWithMessageCheckerFunctionAndFail(
			func(msg *sarama.ProducerMessage) error {
				require.Equal(t, topic, msg.Topic)
				require.Equal(t, sarama.ByteEncoder(msgKey), msg.Key)
				require.Equal(t, sarama.ByteEncoder(msgValue), msg.Value)

				return nil
			},
			produceError,
		)
	}

	return mockProducer
}

// getMockBrokers creates mock seed/leader Kafka brokers for testing, and instantiates an
// IndexerMessageSenderKafka struct connected to the mock brokers.
func getMockBrokersAndSender(
	t *testing.T,
	topic string, // topic that messages will be sent to
	numMessages int, // number of messages that will be sent
) (
	seed *sarama.MockBroker,
	leader *sarama.MockBroker,
	sender *IndexerMessageSenderKafka,
) {
	seed = sarama.NewMockBroker(t, 1)
	leader = sarama.NewMockBroker(t, 2)

	// Metadata response from the seed broker should include the topic and the leader broker.
	metadataResponse := new(sarama.MetadataResponse)
	metadataResponse.AddTopicPartition(topic, 0, leader.BrokerID(), nil, nil, nil, sarama.ErrNoError)
	metadataResponse.AddBroker(leader.Addr(), leader.BrokerID())
	seed.Returns(metadataResponse)

	// Response the leader broker sends when messages are produced to it.
	produceResponse := new(sarama.ProduceResponse)
	produceResponse.AddTopicPartition(topic, 0, sarama.ErrNoError)
	leader.Returns(produceResponse)

	config := mocks.NewTestConfig()
	config.Producer.Flush.Messages = numMessages
	sender, err := NewIndexerMessageSenderKafka(
		indexer.IndexerFlags{
			KafkaAddrs: []string{seed.Addr()},
			MaxRetries: indexer.DefaultMaxRetries,
		},
		config,
		log.NewNopLogger(),
	)
	require.NoError(t, err)

	return seed, leader, sender
}

func TestIndexerMessageSenderKafka_SendOnchainData_WithMockProducer(t *testing.T) {
	tests := map[string]struct {
		// Expectations
		numSuccesses int
		numErrors    int
	}{
		"Only successes": {
			numSuccesses: 10,
			numErrors:    0,
		},
		"Only errors": {
			numSuccesses: 0,
			numErrors:    5,
		},
		"Mixed successes and errors": {
			numSuccesses: 10,
			numErrors:    5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockProducer := getMockProducer(t, ON_CHAIN_KAFKA_TOPIC, tc.numSuccesses, tc.numErrors)

			sender := NewIndexerMessageSenderKafkaWithProducer(mockProducer, log.NewNopLogger())
			for i := 0; i < tc.numSuccesses+tc.numErrors; i++ {
				sender.SendOnchainData(Message{Key: []byte(msgKey), Value: []byte(msgValue)})
			}

			err := sender.Close()
			require.NoError(t, err)

			require.Equal(t, tc.numSuccesses, sender.successes)
			require.Equal(t, tc.numErrors, sender.errors)
		})
	}
}

// Tests connecting to and sending data to a Kafka broker, but does not test that the correct data
// has been sent.
func TestIndexerMessageSenderKafka_SendOnchainData_WithMockBroker(t *testing.T) {
	toSend := 10
	seed, leader, sender := getMockBrokersAndSender(t, ON_CHAIN_KAFKA_TOPIC, toSend)

	for i := 0; i < toSend; i++ {
		sender.SendOnchainData(Message{Key: []byte(msgKey), Value: []byte(msgValue)})
	}

	// Wait for communication between brokers and producer.
	time.Sleep(time.Second)

	err := sender.Close()
	require.NoError(t, err)
	require.Equal(t, toSend, sender.successes)
	require.Equal(t, 0, sender.errors)

	seed.Close()
	leader.Close()
}

func TestIndexerMessageSenderKafka_SendOffchainData_WithMockProducer(t *testing.T) {
	tests := map[string]struct {
		// Expectations
		numSuccesses int
		numErrors    int
	}{
		"Only successes": {
			numSuccesses: 10,
			numErrors:    0,
		},
		"Only errors": {
			numSuccesses: 0,
			numErrors:    5,
		},
		"Mixed successes and errors": {
			numSuccesses: 10,
			numErrors:    5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockProducer := getMockProducer(t, OFF_CHAIN_KAFKA_TOPIC, tc.numSuccesses, tc.numErrors)

			sender := NewIndexerMessageSenderKafkaWithProducer(mockProducer, log.NewNopLogger())
			for i := 0; i < tc.numSuccesses+tc.numErrors; i++ {
				sender.SendOffchainData(Message{Key: []byte(msgKey), Value: []byte(msgValue)})
			}

			err := sender.Close()
			require.NoError(t, err)

			require.Equal(t, tc.numSuccesses, sender.successes)
			require.Equal(t, tc.numErrors, sender.errors)
		})
	}
}

// Tests connecting to and sending data to a Kafka broker, but does not test that the correct data
// has been sent.
func TestIndexerMessageSenderKafka_SendOffchainData_WithMockBroker(t *testing.T) {
	toSend := 10
	seed, leader, sender := getMockBrokersAndSender(t, OFF_CHAIN_KAFKA_TOPIC, toSend)

	for i := 0; i < toSend; i++ {
		sender.SendOffchainData(Message{Key: []byte(msgKey), Value: []byte(msgValue)})
	}

	// Wait for communication between brokers and producer.
	time.Sleep(time.Second)

	err := sender.Close()
	require.NoError(t, err)
	require.Equal(t, toSend, sender.successes)
	require.Equal(t, 0, sender.errors)

	seed.Close()
	leader.Close()
}

func TestIndexerMessageSenderKafka_Close(t *testing.T) {
	mockProducer := getMockProducer(t, ON_CHAIN_KAFKA_TOPIC, 0, 0)
	sender := NewIndexerMessageSenderKafkaWithProducer(mockProducer, log.NewNopLogger())

	err := sender.Close()
	require.NoError(t, err)

	// Closing again should return an error.
	err = sender.Close()
	require.EqualError(t, err, ErrKafkaAlreadyClosed.Error())
}
