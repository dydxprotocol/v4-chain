//go:build all || integration_test

package msgsender

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	sdklog "cosmossdk.io/log"
	"github.com/Shopify/sarama"
	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"
)

const (
	port                 = "9092"
	ip                   = "localhost"
	protocol             = "tcp"
	messageKeyOnchain    = "keyOnchain"
	messageValueOnchain  = "valueOnchain"
	messageKeyOffchain   = "keyOffChain"
	messageValueOffchain = "valueOffchain"
)

// TestMain sets up an instance of Kafka running in a docker container and runs the integration
// tests in this file against the Kafka instance.
func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Creating docker pool: %v", err)
	}

	if err = pool.Client.Ping(); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "blacktop/kafka",
		Tag:        "2.6",
		Env: []string{
			"KAFKA_ADVERTISED_HOST_NAME=localhost",
			fmt.Sprintf("KAFKA_CREATE_TOPICS=%s:1:1,%s:1:1", ON_CHAIN_KAFKA_TOPIC, OFF_CHAIN_KAFKA_TOPIC),
		},
		ExposedPorts: []string{fmt.Sprintf("%s/%s", port, protocol)},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9092/tcp": {{HostIP: ip, HostPort: fmt.Sprintf("%s/%s", port, protocol)}},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatalf("could not start kafka: %s", err)
	}

	// Tell docker to hard kill the container in 10 minutes.
	if err := resource.Expire(600); err != nil {
		log.Fatalf("could not set kafka to expire: %s", err)
	}

	// Function to check if Kafka has been stood up fully, returns an error until all conditions for
	// Kafka to be fully stood up are met.
	// These conditions are:
	// - a consumer is able to connect to Kafka
	// - both on-chain and off-chain topcis are created
	retryFn := func() error {
		consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%s", ip, port)}, nil)
		if err != nil {
			return err
		}
		defer consumer.Close()

		topics, err := consumer.Topics()
		if err != nil {
			return err
		}

		if !reflect.DeepEqual(topics, []string{ON_CHAIN_KAFKA_TOPIC, OFF_CHAIN_KAFKA_TOPIC}) {
			return errors.New("waiting for topics to be created")
		}

		return nil
	}

	// Wait up to 5 minutes for the Kafka docker container to be stood up.
	pool.MaxWait = 300 * time.Second

	if err = pool.Retry(retryFn); err != nil {
		log.Fatalf("could not connect to kafka: %s", err)
	}

	code := m.Run()

	// Remove the docker container once the test is complete.
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestIndexerMessageSenderKafka_VerifySend(t *testing.T) {
	messageSender, err := NewIndexerMessageSenderKafka(
		indexer.IndexerFlags{
			KafkaAddrs: []string{fmt.Sprintf("%s:%s", ip, port)},
			MaxRetries: indexer.DefaultMaxRetries,
		},
		nil,
		sdklog.NewNopLogger(),
	)
	require.NoError(t, err)

	// Test sending data to on-chain topic.
	messageSender.SendOnchainData(Message{
		Key:   []byte(messageKeyOnchain + "VerifySend"),
		Value: []byte(messageValueOnchain + "VerifySend"),
	})
	verifyMessage(
		t,
		ON_CHAIN_KAFKA_TOPIC,
		[]byte(messageKeyOnchain+"VerifySend"),
		[]byte(messageValueOnchain+"VerifySend"),
	)

	// Test sending data to off-chain topic.
	messageSender.SendOffchainData(Message{
		Key:   []byte(messageKeyOffchain + "VerifySend"),
		Value: []byte(messageValueOffchain + "VerifySend"),
	})
	verifyMessage(
		t,
		OFF_CHAIN_KAFKA_TOPIC,
		[]byte(messageKeyOffchain+"VerifySend"),
		[]byte(messageValueOffchain+"VerifySend"),
	)

	messageSender.Close()
}

func TestIndexerMessageSenderKafka_SendAfterClosed(t *testing.T) {
	messageSender, err := NewIndexerMessageSenderKafka(
		indexer.IndexerFlags{
			KafkaAddrs: []string{fmt.Sprintf("%s:%s", ip, port)},
			MaxRetries: indexer.DefaultMaxRetries,
		},
		nil,
		sdklog.NewNopLogger(),
	)
	require.NoError(t, err)
	messageSender.Close()

	// Test send when closed is a no-op and doesn't panic
	messageSender.SendOnchainData(Message{
		Key:   []byte(messageKeyOnchain + "SendAfterClosed"),
		Value: []byte(messageValueOnchain + "SendAfterClosed"),
	})
	messageSender.SendOffchainData(Message{
		Key:   []byte(messageKeyOffchain + "SendAfterClosed"),
		Value: []byte(messageValueOffchain + "SendAfterClosed"),
	})
}

func TestIndexerMessageSenderKafka_ConcurrentSendAndClosed(t *testing.T) {
	max := 100
	waitForEnd := sync.WaitGroup{}
	waitForEnd.Add(3 * max)
	for i := 0; i < max; i++ {
		i := i
		messageSender, err := NewIndexerMessageSenderKafka(
			indexer.IndexerFlags{
				KafkaAddrs: []string{fmt.Sprintf("%s:%s", ip, port)},
				MaxRetries: indexer.DefaultMaxRetries,
			},
			nil,
			sdklog.NewNopLogger(),
		)
		require.NoError(t, err)

		waitTillReady := sync.WaitGroup{}
		waitTillReady.Add(1)

		go func() {
			defer waitForEnd.Done()
			waitTillReady.Wait()
			messageSender.SendOnchainData(Message{
				Key:   []byte(messageKeyOnchain + "ConcurrentSendAndClosed" + strconv.Itoa(i)),
				Value: []byte(messageValueOnchain + "ConcurrentSendAndClosed" + strconv.Itoa(i)),
			})
		}()
		go func() {
			defer waitForEnd.Done()
			waitTillReady.Wait()
			messageSender.SendOffchainData(Message{
				Key:   []byte(messageKeyOffchain + "ConcurrentSendAndClosed" + strconv.Itoa(i)),
				Value: []byte(messageValueOffchain + "ConcurrentSendAndClosed" + strconv.Itoa(i)),
			})
		}()
		go func() {
			defer waitForEnd.Done()
			waitTillReady.Wait()
			messageSender.Close()
		}()
		waitTillReady.Done()
	}
	waitForEnd.Wait()
}

// verifyMessage checks that the first message in a topic has the given key/value.
func verifyMessage(t *testing.T, topic string, key []byte, value []byte) {
	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%s", ip, port)}, nil)
	require.NoError(t, err)

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, 0)
	require.NoError(t, err)

	msg := <-partitionConsumer.Messages()
	partitionConsumer.Close()

	require.Equal(t, key, msg.Key)
	require.Equal(t, value, msg.Value)
}
