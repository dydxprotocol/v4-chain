package indexer_manager_test

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/stretchr/testify/require"
)

func TestCreateIndexerBlockEventMessage(t *testing.T) {
	blockEvent := &indexer_manager.IndexerTendermintBlock{
		Height: uint32(BlockHeight),
		Time:   BlockTime,
		Events: []*indexer_manager.IndexerTendermintEvent{
			&OrderFillTendermintEvent,
			&TransferTendermintEvent,
			&SubaccountTendermintEvent,
		},
		TxHashes: []string{TxHash, TxHash1},
	}
	actualMessage := indexer_manager.CreateIndexerBlockEventMessage(blockEvent)
	blockEventBytes, err := proto.Marshal(blockEvent)
	require.NoError(t, err)
	expectedMessage := msgsender.Message{
		Value: blockEventBytes,
	}
	require.Equal(t, expectedMessage, actualMessage)
}
