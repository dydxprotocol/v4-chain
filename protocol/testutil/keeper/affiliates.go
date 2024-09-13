package keeper

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	affiliateskeeper "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	statskeeper "github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"
)

func AffiliatesKeepers(t testing.TB, msgSenderEnabled bool) (
	ctx sdk.Context,
	keeper *affiliateskeeper.Keeper,
	statsKeeper *statskeeper.Keeper,
	storeKey storetypes.StoreKey,
) {
	ctx = initKeepers(t, func(
		db *dbm.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)

		accountKeeper, _ := createAccountKeeper(stateStore, db, cdc, registry)
		bankKeeper, _ := createBankKeeper(stateStore, db, cdc, accountKeeper)
		stakingKeeper, _ := createStakingKeeper(stateStore, db, cdc, accountKeeper, bankKeeper)
		statsKeeper, _ := createStatsKeeper(stateStore, epochsKeeper, db, cdc, stakingKeeper)

		keeper, storeKey = createAffiliatesKeeper(
			stateStore,
			db,
			cdc,
			statsKeeper,
			transientStoreKey,
			msgSenderEnabled,
		)

		return []GenesisInitializer{statsKeeper, keeper}
	})

	return ctx,
		keeper,
		statsKeeper,
		storeKey
}

func createAffiliatesKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	statsKeeper *statskeeper.Keeper,
	transientStoreKey storetypes.StoreKey,
	msgSenderEnabled bool,
) (*affiliateskeeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(msgSenderEnabled)
	mockIndexerEventsManager := indexer_manager.NewIndexerEventManager(mockMsgSender, transientStoreKey, true)

	k := affiliateskeeper.NewKeeper(
		cdc,
		storeKey,
		[]string{},
		statsKeeper,
		mockIndexerEventsManager,
	)
	return k, storeKey
}
