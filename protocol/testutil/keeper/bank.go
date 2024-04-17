package keeper

import (
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func createBankKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	accountKeeper *authkeeper.AccountKeeper,
) (*keeper.BaseKeeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	k := keeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		accountKeeper,
		map[string]bool{},
		lib.GovModuleAddress.String(),
		log.NewNopLogger(),
	)

	return &k, storeKey
}
