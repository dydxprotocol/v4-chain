package keeper

import (
	dbm "github.com/cosmos/cosmos-db"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bridgeserver_types "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/bridge"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
)

func BridgeKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	mockTimeProvider *mocks.TimeProvider,
	bridgeEventManager *bridgeserver_types.BridgeEventManager,
	bankKeeper *bankkeeper.BaseKeeper,
	mockDelayMsgKeeper *mocks.DelayMsgKeeper,
) {
	ctx = initKeepers(t, func(
		db *dbm.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		accountKeeper, _ := createAccountKeeper(stateStore, db, cdc, registry)
		bankKeeper, _ = createBankKeeper(stateStore, db, cdc, accountKeeper)
		keeper, storeKey, mockTimeProvider, bridgeEventManager, mockDelayMsgKeeper =
			createBridgeKeeper(stateStore, db, cdc, transientStoreKey, bankKeeper)
		return []GenesisInitializer{keeper}
	})

	return ctx, keeper, storeKey, mockTimeProvider, bridgeEventManager, bankKeeper, mockDelayMsgKeeper
}

func createBridgeKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	transientStoreKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
	*mocks.TimeProvider,
	*bridgeserver_types.BridgeEventManager,
	*mocks.DelayMsgKeeper,
) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	mockTimeProvider := &mocks.TimeProvider{}
	bridgeEventManager := bridgeserver_types.NewBridgeEventManager(mockTimeProvider)

	mockDelayMsgKeeper := &mocks.DelayMsgKeeper{}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		bridgeEventManager,
		bankKeeper,
		mockDelayMsgKeeper,
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgtypes.ModuleAddress.String(),
		},
	)

	return k, storeKey, mockTimeProvider, bridgeEventManager, mockDelayMsgKeeper
}
