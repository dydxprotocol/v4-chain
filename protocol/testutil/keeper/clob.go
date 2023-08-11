package keeper

import (
	"testing"

	"github.com/dydxprotocol/v4/indexer/indexer_manager"

	db "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	asskeeper "github.com/dydxprotocol/v4/x/assets/keeper"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/types"
	perpkeeper "github.com/dydxprotocol/v4/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4/x/prices/keeper"
	subkeeper "github.com/dydxprotocol/v4/x/subaccounts/keeper"
)

func ClobKeepers(
	t testing.TB,
	memClob types.MemClob,
	bankKeeper bankkeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	assetsKeeper *asskeeper.Keeper,
	perpetualsKeeper *perpkeeper.Keeper,
	subaccountsKeeper *subkeeper.Keeper,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
) {
	ctx,
		keeper,
		pricesKeeper,
		assetsKeeper,
		perpetualsKeeper,
		subaccountsKeeper,
		storeKey,
		memKey = ClobKeepersWithUninitializedMemStore(t, memClob, bankKeeper, indexerEventManager)

	// Initialize the memstore.
	keeper.InitMemStore(ctx)

	return ctx, keeper, pricesKeeper, assetsKeeper, perpetualsKeeper, subaccountsKeeper, storeKey, memKey
}

func ClobKeepersWithUninitializedMemStore(
	t testing.TB,
	memClob types.MemClob,
	bankKeeper bankkeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	assetsKeeper *asskeeper.Keeper,
	perpetualsKeeper *perpkeeper.Keeper,
	subaccountsKeeper *subkeeper.Keeper,
	storeKey storetypes.StoreKey,
	memKey storetypes.StoreKey,
) {
	ctx = initKeepers(t, func(
		db *db.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		indexerEventsTransientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		pricesKeeper, _, _, _ = createPricesKeeper(stateStore, db, cdc, indexerEventsTransientStoreKey)
		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)
		perpetualsKeeper, _ = createPerpetualsKeeper(
			stateStore,
			db,
			cdc,
			pricesKeeper,
			epochsKeeper,
			indexerEventsTransientStoreKey,
		)
		assetsKeeper, _ = createAssetsKeeper(stateStore, db, cdc, pricesKeeper)
		subaccountsKeeper, _ = createSubaccountsKeeper(
			stateStore,
			db,
			cdc,
			assetsKeeper,
			bankKeeper,
			perpetualsKeeper,
			indexerEventsTransientStoreKey,
			true,
		)
		keeper, storeKey, memKey = createClobKeeper(
			stateStore,
			db,
			cdc,
			memClob,
			assetsKeeper,
			bankKeeper,
			perpetualsKeeper,
			subaccountsKeeper,
			indexerEventManager,
			indexerEventsTransientStoreKey,
		)

		return []GenesisInitializer{pricesKeeper, perpetualsKeeper, assetsKeeper, subaccountsKeeper, keeper}
	})

	return ctx, keeper, pricesKeeper, assetsKeeper, perpetualsKeeper, subaccountsKeeper, storeKey, memKey
}

func createClobKeeper(
	stateStore storetypes.CommitMultiStore,
	db *db.MemDB,
	cdc *codec.ProtoCodec,
	memClob types.MemClob,
	aKeeper *asskeeper.Keeper,
	bankKeeper types.BankKeeper,
	perpKeeper *perpkeeper.Keeper,
	saKeeper *subkeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
	indexerEventsTransientStoreKey storetypes.StoreKey,
) (*keeper.Keeper, storetypes.StoreKey, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	transientStoreKey := sdk.NewTransientStoreKey(types.TransientStoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		transientStoreKey,
		memClob,
		saKeeper,
		aKeeper,
		bankKeeper,
		perpKeeper,
		indexerEventManager,
	)

	return k, storeKey, memKey
}
