package indexer_manager

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/common"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
)

const (
	createErrMsg          = "Cannot create message."
	onChainEventsKafkaKey = "on_chain_events"
)

func marshalIndexerTendermintBlock(
	indexerTendermintBlock *IndexerTendermintBlock,
	marshaler common.Marshaler,
) ([]byte, error) {
	bytes, err := marshaler.Marshal(indexerTendermintBlock)
	return bytes, err
}

// CreateIndexerBlockEventMessage creates an on-chain update message for all the Indexer events in a block.
func CreateIndexerBlockEventMessage(
	block *IndexerTendermintBlock,
) msgsender.Message {
	errMessage := "Error creating on-chain Indexer block event message."
	errDetails := fmt.Sprintf("Block: %+v", *block)

	update, err := marshalIndexerTendermintBlock(block, &common.MarshalerImpl{})
	if err != nil {
		panic(fmt.Sprintf("%s %s Err: %+v %s\n", errMessage, createErrMsg, err, errDetails))
	}

	return msgsender.Message{Value: update}
}
