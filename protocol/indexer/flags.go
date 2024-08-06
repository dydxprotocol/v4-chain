package indexer

import (
	"strings"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

type IndexerFlags struct {
	KafkaAddrs       []string
	MaxRetries       int
	SendOffchainData bool
}

// List of default values
const (
	DefaultMaxRetries = 20
)

// List of CLI flags
const (
	FlagKafkaConnStr = "indexer-kafka-conn-str"
	// max retry should be set so that max retry * retry backoff > Zookeeper session.timeout + some buffer
	FlagKafkaMaxRetry        = "indexer-kafka-max-retry"
	FlagSendOffchainData     = "indexer-send-offchain-data"
	MsgSenderInstanceForTest = "msgsender-instance-for-test"
)

// AddIndexerFlagsToCmd adds the required flags to instantiate a connection to Kafka during App
// initialization to a command. These flags should be applied to the `start` command of the
// V4 Cosmos application.
// E.g. `dydxprotocold start --indexer-kafka-conn-str kafka:9092
func AddIndexerFlagsToCmd(cmd *cobra.Command) {
	cmd.
		Flags().
		String(
			FlagKafkaConnStr,
			"",
			"Comma delimited list of Kafka Broker addresses in the form of <url>:<port> to connect to "+
				"for sending data to the Indexer, no connections to Kafka Brokers will be made if the "+
				"value is an empty string. E.g. \"localhost:9092,localhost:29092\"",
		)
	cmd.
		Flags().
		Int(
			FlagKafkaMaxRetry,
			DefaultMaxRetries,
			"Maximum number of times to retry sending a message to the Indexer",
		)
	cmd.
		Flags().
		Bool(
			FlagSendOffchainData,
			true,
			"Whether to send offchain data to the Indexer. This should be set to false when the Indexer "+
				"full node is being restarted from a snapshot and is behind the Indexer's view of the "+
				"chain during the fast sync process.",
		)
}

// GetIndexerFlagValuesFromOptions gets values for connecting to Kafka from the `AppOptions`
// struct which contains values from the Indexer command-line flags.
func GetIndexerFlagValuesFromOptions(
	appOpts servertypes.AppOptions,
) IndexerFlags {
	option := appOpts.Get(FlagKafkaConnStr)
	kafkaConnStr, err := cast.ToStringE(option)
	if option == nil || err != nil {
		return IndexerFlags{
			KafkaAddrs:       []string{},
			MaxRetries:       DefaultMaxRetries,
			SendOffchainData: false,
		}
	}

	maxRetries := cast.ToInt(appOpts.Get(FlagKafkaMaxRetry))
	sendOffchainData := cast.ToBool(appOpts.Get(FlagSendOffchainData))

	var kafkaAddrs []string
	if kafkaConnStr == "" {
		kafkaAddrs = []string{}
	} else {
		kafkaAddrs = strings.Split(kafkaConnStr, ",")
	}

	return IndexerFlags{
		KafkaAddrs:       kafkaAddrs,
		MaxRetries:       maxRetries,
		SendOffchainData: sendOffchainData,
	}
}
