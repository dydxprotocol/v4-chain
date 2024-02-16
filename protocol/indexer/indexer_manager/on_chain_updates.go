package indexer_manager

import (
	"fmt"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
)

const (
	createErrMsg          = "Cannot create message."
	onChainEventsKafkaKey = "on_chain_events"
)

// CreateIndexerBlockEventMessage creates an on-chain update message for all the Indexer events in a block.
func CreateIndexerBlockEventMessage(
	block *IndexerTendermintBlock,
) msgsender.Message {
	errMessage := "Error creating on-chain Indexer block event message."
	errDetails := fmt.Sprintf("Block: %+v", *block)

	update, err := proto.Marshal(block)
	if err != nil {
		panic(fmt.Sprintf("%s %s Err: %+v %s\n", errMessage, createErrMsg, err, errDetails))
	}

	return msgsender.Message{Value: update}
}
