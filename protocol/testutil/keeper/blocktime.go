package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	delaymsgtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
)

func createBlockTimeKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	authorities := []string{
		delaymsgtypes.ModuleAddress.String(),
		lib.GovModuleAddress.String(),
	}
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		authorities,
	)

	return k, storeKey
}
