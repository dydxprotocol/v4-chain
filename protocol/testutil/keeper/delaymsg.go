package keeper

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ante"
	dbm "github.com/cosmos/cosmos-db"

	storetypes "cosmossdk.io/store/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	perpetualskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/keeper"
	perpetualstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	priceskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DelayMsgKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	delayMsgKeeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	perpsKeeper *perpetualskeeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	authorities []string,
) {
	ctx = initKeepers(t, func(
		db *dbm.MemDB,
		_ codectypes.InterfaceRegistry,
		_ *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		encCfg := ante.MakeTestEncodingConfig()
		cdc := encCfg.Codec.(*codec.ProtoCodec)
		registry := encCfg.InterfaceRegistry

		router := baseapp.NewMsgServiceRouter()
		router.SetInterfaceRegistry(registry)

		// Register perpetuals messages for encoding / decoding.
		perpetualstypes.RegisterInterfaces(registry)

		epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)
		pricesKeeper, _, _, _, _ = createPricesKeeper(stateStore, db, cdc, transientStoreKey)
		perpsKeeper, _ = createPerpetualsKeeper(stateStore, db, cdc, pricesKeeper, epochsKeeper, transientStoreKey)

		// Register perps keeper msg server for msg routing.
		perpetualstypes.RegisterMsgServer(router, perpetualskeeper.NewMsgServerImpl(perpsKeeper))

		authorities = []string{
			lib.GovModuleAddress.String(),
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
	return ctx, delayMsgKeeper, storeKey, perpsKeeper, pricesKeeper, authorities
}

func createDelayMsgKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	router *baseapp.MsgServiceRouter,
	authorities []string,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		router,
		authorities,
	)
	return k, storeKey
}
