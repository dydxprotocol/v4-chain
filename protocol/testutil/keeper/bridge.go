package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	bridgeserver_types "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/bridge"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
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
		db *tmdb.MemDB,
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
	db *tmdb.MemDB,
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
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
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
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String(),
		},
	)

	return k, storeKey, mockTimeProvider, bridgeEventManager, mockDelayMsgKeeper
}
