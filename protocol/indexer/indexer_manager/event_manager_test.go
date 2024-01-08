package indexer_manager_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/msgsender"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sdk"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var ExpectedEvent0 = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeOrderFill,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{
		TransactionIndex: 0,
	},
	EventIndex: 0,
	Version:    indexerevents.OrderFillEventVersion,
	DataBytes: indexer_manager.GetBytes(
		&OrderFillEvent,
	),
}

var ExpectedEvent1 = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeSubaccountUpdate,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{
		TransactionIndex: 0,
	},
	EventIndex: 1,
	Version:    indexerevents.SubaccountUpdateEventVersion,
	DataBytes: indexer_manager.GetBytes(
		&SubaccountEvent,
	),
}

var ExpectedEvent2 = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeTransfer,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_TransactionIndex{
		TransactionIndex: 1,
	},
	EventIndex: 0,
	Version:    indexerevents.TransferEventVersion,
	DataBytes: indexer_manager.GetBytes(
		&TransferEvent,
	),
}

var ExpectedEvent3 = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeFundingValues,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_BlockEvent_{
		BlockEvent: indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK,
	},
	EventIndex: 0,
	Version:    indexerevents.FundingValuesEventVersion,
	DataBytes: indexer_manager.GetBytes(
		&FundingRateAndIndexEvent,
	),
}

var ExpectedEvent4 = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeFundingValues,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_BlockEvent_{
		BlockEvent: indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK,
	},
	EventIndex: 1,
	Version:    indexerevents.FundingValuesEventVersion,
	DataBytes: indexer_manager.GetBytes(
		&FundingPremiumSampleEvent,
	),
}

var ExpectedEvent5 = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeFundingValues,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_BlockEvent_{
		BlockEvent: indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_BEGIN_BLOCK,
	},
	EventIndex: 0,
	Version:    indexerevents.FundingValuesEventVersion,
	DataBytes: indexer_manager.GetBytes(
		&FundingPremiumSampleEvent,
	),
}

var ExpectedEvent6 = indexer_manager.IndexerTendermintEvent{
	Subtype: indexerevents.SubtypeFundingValues,
	OrderingWithinBlock: &indexer_manager.IndexerTendermintEvent_BlockEvent_{
		BlockEvent: indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_BEGIN_BLOCK,
	},
	EventIndex: 1,
	Version:    indexerevents.FundingValuesEventVersion,
	DataBytes: indexer_manager.GetBytes(
		&FundingRateAndIndexEvent,
	),
}

var EventVersion uint32 = 1

func assertIsEnabled(t *testing.T, isEnabled bool) {
	storeKey := storetypes.NewTransientStoreKey(indexer_manager.TransientStoreKey)
	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(isEnabled)
	indexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, storeKey, isEnabled)
	require.Equal(t, isEnabled, indexerEventManager.Enabled())
}

func TestIsEnabled(t *testing.T) {
	assertIsEnabled(t, true)
	assertIsEnabled(t, false)
}

func TestSendOffchainData(t *testing.T) {
	storeKey := storetypes.NewTransientStoreKey(indexer_manager.TransientStoreKey)
	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	mockMsgSender.On("SendOffchainData", mock.Anything).Return(nil)
	indexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, storeKey, true)
	var message msgsender.Message
	indexerEventManager.SendOffchainData(message)
	mockMsgSender.AssertExpectations(t)
}

func TestSendOnchainData(t *testing.T) {
	storeKey := storetypes.NewTransientStoreKey(indexer_manager.TransientStoreKey)
	indexerTendermintBlock := &indexer_manager.IndexerTendermintBlock{}
	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	mockMsgSender.On("SendOnchainData", mock.Anything).Return(nil)
	indexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, storeKey, true)
	indexerEventManager.SendOnchainData(indexerTendermintBlock)
	mockMsgSender.AssertExpectations(t)
}

func TestProduceBlockBasicTxnEvent(t *testing.T) {
	ctx, stateStore, db := sdk.NewSdkContextWithMultistore()
	storeKey := storetypes.NewTransientStoreKey(indexer_manager.TransientStoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeTransient, db)
	ctx = ctx.WithBlockTime(BlockTime).WithBlockHeight(BlockHeight).WithTxBytes(constants.TestTxBytes)
	ctx.GasMeter().ConsumeGas(ConsumedGas, "beforeWrite")
	require.NoError(t, stateStore.LoadLatestVersion())
	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	indexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, storeKey, true)
	indexerEventManager.AddTxnEvent(
		ctx,
		indexerevents.SubtypeOrderFill,
		EventVersion,
		indexer_manager.GetBytes(
			&OrderFillEvent,
		),
	)

	block := indexerEventManager.ProduceBlock(ctx)
	require.Len(t, block.Events, 1)
	require.Equal(t, ExpectedEvent0, *block.Events[0])
	require.Equal(t, []string{string(constants.TestTxHashString)}, block.TxHashes)
	require.Equal(t, uint32(BlockHeight), block.Height)
	require.Equal(t, BlockTime, block.Time)
	require.Equal(t, ConsumedGas, ctx.GasMeter().GasConsumed())
}

func TestProduceBlockBasicBlockEvent(t *testing.T) {
	ctx, stateStore, db := sdk.NewSdkContextWithMultistore()
	storeKey := storetypes.NewTransientStoreKey(indexer_manager.TransientStoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeTransient, db)
	ctx = ctx.WithBlockTime(BlockTime).WithBlockHeight(BlockHeight)
	ctx.GasMeter().ConsumeGas(ConsumedGas, "beforeWrite")
	require.NoError(t, stateStore.LoadLatestVersion())
	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	indexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, storeKey, true)
	indexerEventManager.AddBlockEvent(
		ctx,
		indexerevents.SubtypeFundingValues,
		indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK,
		EventVersion,
		indexer_manager.GetBytes(
			&FundingRateAndIndexEvent,
		),
	)

	block := indexerEventManager.ProduceBlock(ctx)
	require.Len(t, block.Events, 1)
	require.Equal(t, ExpectedEvent3, *block.Events[0])
	require.Empty(t, block.TxHashes)
	require.Equal(t, uint32(BlockHeight), block.Height)
	require.Equal(t, BlockTime, block.Time)
	require.Equal(t, ConsumedGas, ctx.GasMeter().GasConsumed())
}

func TestProduceBlockMultipleTxnEvents(t *testing.T) {
	ctx, stateStore, db := sdk.NewSdkContextWithMultistore()
	storeKey := storetypes.NewTransientStoreKey(indexer_manager.TransientStoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeTransient, db)
	ctx = ctx.WithBlockTime(BlockTime).WithBlockHeight(BlockHeight).WithTxBytes(constants.TestTxBytes)
	ctx.GasMeter().ConsumeGas(ConsumedGas, "beforeWrite")
	require.NoError(t, stateStore.LoadLatestVersion())
	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	indexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, storeKey, true)
	indexerEventManager.AddTxnEvent(
		ctx,
		indexerevents.SubtypeOrderFill,
		EventVersion,
		indexer_manager.GetBytes(
			&OrderFillEvent,
		),
	)
	indexerEventManager.AddTxnEvent(
		ctx,
		indexerevents.SubtypeSubaccountUpdate,
		EventVersion,
		indexer_manager.GetBytes(
			&SubaccountEvent,
		),
	)
	ctx = ctx.WithTxBytes(constants.TestTxBytes1)
	indexerEventManager.AddTxnEvent(
		ctx,
		indexerevents.SubtypeTransfer,
		EventVersion,
		indexer_manager.GetBytes(
			&TransferEvent,
		),
	)

	block := indexerEventManager.ProduceBlock(ctx)
	require.Len(t, block.Events, 3)
	require.Equal(t, ExpectedEvent0, *block.Events[0])
	require.Equal(t, ExpectedEvent1, *block.Events[1])
	require.Equal(t, ExpectedEvent2, *block.Events[2])
	require.Equal(t, []string{
		string(constants.TestTxHashString),
		string(constants.TestTxHashString1),
	}, block.TxHashes)
	require.Equal(t, uint32(BlockHeight), block.Height)
	require.Equal(t, BlockTime, block.Time)
	require.Equal(t, ConsumedGas, ctx.GasMeter().GasConsumed())
}

func TestProduceBlockMultipleTxnAndBlockEvents(t *testing.T) {
	ctx, stateStore, db := sdk.NewSdkContextWithMultistore()
	storeKey := storetypes.NewTransientStoreKey(indexer_manager.TransientStoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeTransient, db)
	ctx = ctx.WithBlockTime(BlockTime).WithBlockHeight(BlockHeight).WithTxBytes(constants.TestTxBytes)
	ctx.GasMeter().ConsumeGas(ConsumedGas, "beforeWrite")
	require.NoError(t, stateStore.LoadLatestVersion())
	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	indexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, storeKey, true)
	indexerEventManager.AddTxnEvent(
		ctx,
		indexerevents.SubtypeOrderFill,
		EventVersion,
		indexer_manager.GetBytes(
			&OrderFillEvent,
		),
	)
	indexerEventManager.AddTxnEvent(
		ctx,
		indexerevents.SubtypeSubaccountUpdate,
		EventVersion,
		indexer_manager.GetBytes(
			&SubaccountEvent,
		),
	)
	ctx = ctx.WithTxBytes(constants.TestTxBytes1)
	indexerEventManager.AddTxnEvent(
		ctx,
		indexerevents.SubtypeTransfer,
		EventVersion,
		indexer_manager.GetBytes(
			&TransferEvent,
		),
	)
	indexerEventManager.AddBlockEvent(
		ctx,
		indexerevents.SubtypeFundingValues,
		indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK,
		EventVersion,
		indexer_manager.GetBytes(
			&FundingRateAndIndexEvent,
		),
	)
	indexerEventManager.AddBlockEvent(
		ctx,
		indexerevents.SubtypeFundingValues,
		indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK,
		EventVersion,
		indexer_manager.GetBytes(
			&FundingPremiumSampleEvent,
		),
	)
	indexerEventManager.AddBlockEvent(
		ctx,
		indexerevents.SubtypeFundingValues,
		indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_BEGIN_BLOCK,
		EventVersion,
		indexer_manager.GetBytes(
			&FundingPremiumSampleEvent,
		),
	)
	indexerEventManager.AddBlockEvent(
		ctx,
		indexerevents.SubtypeFundingValues,
		indexer_manager.IndexerTendermintEvent_BLOCK_EVENT_BEGIN_BLOCK,
		EventVersion,
		indexer_manager.GetBytes(
			&FundingRateAndIndexEvent,
		),
	)

	block := indexerEventManager.ProduceBlock(ctx)
	require.Len(t, block.Events, 7)
	require.Equal(t, ExpectedEvent0, *block.Events[0])
	require.Equal(t, ExpectedEvent1, *block.Events[1])
	require.Equal(t, ExpectedEvent2, *block.Events[2])
	require.Equal(t, ExpectedEvent3, *block.Events[3])
	require.Equal(t, ExpectedEvent4, *block.Events[4])
	require.Equal(t, ExpectedEvent5, *block.Events[5])
	require.Equal(t, ExpectedEvent6, *block.Events[6])
	require.Equal(t, []string{
		string(constants.TestTxHashString),
		string(constants.TestTxHashString1),
	}, block.TxHashes)
	require.Equal(t, uint32(BlockHeight), block.Height)
	require.Equal(t, BlockTime, block.Time)
	require.Equal(t, ConsumedGas, ctx.GasMeter().GasConsumed())
}

func TestClearEvents(t *testing.T) {
	ctx, stateStore, db := sdk.NewSdkContextWithMultistore()
	storeKey := storetypes.NewTransientStoreKey(indexer_manager.TransientStoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeTransient, db)
	ctx = ctx.WithBlockTime(BlockTime).WithBlockHeight(BlockHeight).WithTxBytes(constants.TestTxBytes)
	ctx.GasMeter().ConsumeGas(ConsumedGas, "beforeWrite")
	require.NoError(t, stateStore.LoadLatestVersion())
	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	indexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, storeKey, true)
	indexerEventManager.AddTxnEvent(
		ctx,
		indexerevents.SubtypeOrderFill,
		EventVersion,
		indexer_manager.GetBytes(
			&OrderFillEvent,
		),
	)

	block := indexerEventManager.ProduceBlock(ctx)
	require.Len(t, block.Events, 1)
	indexerEventManager.ClearEvents(ctx)
	block = indexerEventManager.ProduceBlock(ctx)
	require.Len(t, block.Events, 0)
	require.Equal(t, ConsumedGas, ctx.GasMeter().GasConsumed())
}
