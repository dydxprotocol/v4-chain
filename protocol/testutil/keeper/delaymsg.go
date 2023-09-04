package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	bridgekeeper "github.com/dydxprotocol/v4-chain/protocol/x/bridge/keeper"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
)

func DelayMsgKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	delayMsgKeeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	bridgeKeeper *bridgekeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	authorities []string,
) {
	ctx = initKeepers(t, func(
		db *tmdb.MemDB,
		_ codectypes.InterfaceRegistry,
		_ *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		encCfg := testutil.MakeTestEncodingConfig()
		cdc := encCfg.Codec.(*codec.ProtoCodec)
		registry := encCfg.InterfaceRegistry

		router := baseapp.NewMsgServiceRouter()
		router.SetInterfaceRegistry(registry)

		// Register bridge messages for encoding / decoding.
		bridgetypes.RegisterInterfaces(registry)

		accountKeeper, _ := createAccountKeeper(stateStore, db, cdc, registry)
		bankKeeper, _ = createBankKeeper(stateStore, db, cdc, accountKeeper)
		bridgeKeeper, _, _, _, _ =
			createBridgeKeeper(stateStore, db, cdc, transientStoreKey, bankKeeper)

		// Register bridge keeper msg server for msg routing.
		bridgetypes.RegisterMsgServer(router, bridgekeeper.NewMsgServerImpl(bridgeKeeper))

		authorities = []string{
			authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		}
		delayMsgKeeper, storeKey = createDelayMsgKeeper(
			stateStore,
			db,
			cdc,
			router,
			authorities,
		)

		return []GenesisInitializer{
			delayMsgKeeper,
		}
	})
	return ctx, delayMsgKeeper, storeKey, bridgeKeeper, bankKeeper, authorities
}

func DelayMsgKeeperWithMockBridgeKeeper(
	t testing.TB,
) (
	ctx sdk.Context,
	delayMsgKeeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	bridgeKeeper *mocks.BridgeKeeper,
	authorities []string,
) {
	ctx = initKeepers(t, func(
		db *tmdb.MemDB,
		_ codectypes.InterfaceRegistry,
		_ *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		encCfg := testutil.MakeTestEncodingConfig()
		cdc := encCfg.Codec.(*codec.ProtoCodec)
		registry := encCfg.InterfaceRegistry

		router := baseapp.NewMsgServiceRouter()
		router.SetInterfaceRegistry(registry)

		// Register bridge messages for encoding / decoding.
		bridgetypes.RegisterInterfaces(registry)

		bridgeKeeper = &mocks.BridgeKeeper{}

		// Register bridge keeper msg server for msg routing.
		bridgetypes.RegisterMsgServer(router, bridgekeeper.NewMsgServerImpl(bridgeKeeper))

		authorities = []string{
			authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		}
		delayMsgKeeper, storeKey = createDelayMsgKeeper(
			stateStore,
			db,
			cdc,
			router,
			authorities,
		)

		return []GenesisInitializer{delayMsgKeeper}
	})
	return ctx, delayMsgKeeper, storeKey, bridgeKeeper, authorities
}

func createDelayMsgKeeper(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	router *baseapp.MsgServiceRouter,
	authorities []string,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		router,
		authorities,
	)
	return k, storeKey
}
