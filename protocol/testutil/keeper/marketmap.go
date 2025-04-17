package keeper

import (
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/slinky/x/marketmap/types"

	storetypes "cosmossdk.io/store/types"
	keeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
)

func createMarketMapKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	ss := runtime.NewKVStoreService(storeKey)
	k := keeper.NewKeeper(
		ss, cdc, sdk.AccAddress("authority"),
	)

	return k, storeKey
}
