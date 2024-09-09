package keeper

import (
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	affiliateskeeper "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	statskeeper "github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"
)

func createAffiliatesKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	statsKeeper *statskeeper.Keeper,
) (*affiliateskeeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	k := affiliateskeeper.NewKeeper(
		cdc,
		storeKey,
		[]string{},
		statsKeeper,
	)
	return k, storeKey
}
