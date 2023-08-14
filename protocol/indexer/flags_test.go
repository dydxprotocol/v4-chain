package indexer_test

import (
	"fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddIndexerFlagsToCommand(t *testing.T) {
	cmd := cobra.Command{}

	indexer.AddIndexerFlagsToCmd(&cmd)
	tests := map[string]struct {
		flagName string
	}{
		fmt.Sprintf("Has %s flag", indexer.FlagKafkaConnStr): {
			flagName: indexer.FlagKafkaConnStr,
		},
		fmt.Sprintf("Has %s flag", indexer.FlagKafkaMaxRetry): {
			flagName: indexer.FlagKafkaMaxRetry,
		},
		fmt.Sprintf("Has %s flag", indexer.FlagSendOffchainData): {
			flagName: indexer.FlagSendOffchainData,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Contains(t, cmd.Flags().FlagUsages(), tc.flagName)
		})
	}
}

func TestGetIndexerFlagValuesFromOptions(t *testing.T) {
	tests := map[string]struct {
		// Parameters.
		kafkaConnStr     string
		maxRetries       int
		nilConnStr       bool
		sendOffchainData bool

		// Expectations.
		expectedIndexerFlags indexer.IndexerFlags
	}{
		"Sets KafkaAddrs to empty slice if kafkaConnStr is empty string": {
			kafkaConnStr:     "",
			maxRetries:       0,
			nilConnStr:       false,
			sendOffchainData: false,
			expectedIndexerFlags: indexer.IndexerFlags{
				KafkaAddrs:       []string{},
				MaxRetries:       0,
				SendOffchainData: false,
			},
		},
		"Sets KafkaAddrs to slice of 1 string if no commas in kafkaConnStr": {
			kafkaConnStr:     "kafka:9092",
			maxRetries:       0,
			nilConnStr:       false,
			sendOffchainData: true,
			expectedIndexerFlags: indexer.IndexerFlags{
				KafkaAddrs:       []string{"kafka:9092"},
				MaxRetries:       0,
				SendOffchainData: true,
			},
		},
		"Sets KafkaAddrs to slice of multiple strings if commas in kafkaConnStr": {
			kafkaConnStr:     "kafka:9092,kafka:9093,kafka:9094",
			maxRetries:       0,
			nilConnStr:       false,
			sendOffchainData: true,
			expectedIndexerFlags: indexer.IndexerFlags{
				KafkaAddrs:       []string{"kafka:9092", "kafka:9093", "kafka:9094"},
				MaxRetries:       0,
				SendOffchainData: true,
			},
		},
		"Sets MaxRetries": {
			kafkaConnStr:     "",
			maxRetries:       5,
			nilConnStr:       false,
			sendOffchainData: false,
			expectedIndexerFlags: indexer.IndexerFlags{
				KafkaAddrs:       []string{},
				MaxRetries:       5,
				SendOffchainData: false,
			},
		},
		"Sets KafkaAddrs to empty slice and MaxRetries to default if kafkaConnStr is nil": {
			kafkaConnStr:     "kafka:9092",
			maxRetries:       5,
			nilConnStr:       true,
			sendOffchainData: false,
			expectedIndexerFlags: indexer.IndexerFlags{
				KafkaAddrs:       []string{},
				MaxRetries:       indexer.DefaultMaxRetries,
				SendOffchainData: false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			optsMap := make(map[string]interface{})
			if tc.nilConnStr {
				optsMap[indexer.FlagKafkaConnStr] = nil
			} else {
				optsMap[indexer.FlagKafkaConnStr] = tc.kafkaConnStr
			}
			optsMap[indexer.FlagKafkaMaxRetry] = tc.maxRetries
			optsMap[indexer.FlagSendOffchainData] = tc.sendOffchainData
			mockOpts := mocks.AppOptions{}
			mockOpts.On("Get", mock.AnythingOfType("string")).
				Return(func(key string) interface{} {
					return optsMap[key]
				})

			indexerFlags := indexer.GetIndexerFlagValuesFromOptions(&mockOpts)
			require.Equal(t, tc.expectedIndexerFlags, indexerFlags)
		})
	}
}
