package keeper

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ante"
	dbm "github.com/cosmos/cosmos-db"

	storetypes "cosmossdk.io/store/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

func DelayMsgKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	delayMsgKeeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	bankKeeper bankkeeper.Keeper,
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

		accountKeeper, _ := createAccountKeeper(stateStore, db, cdc, registry)
		bankKeeper, _ = createBankKeeper(stateStore, db, cdc, accountKeeper)

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
	return ctx, delayMsgKeeper, storeKey, bankKeeper, authorities
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
