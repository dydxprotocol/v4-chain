package keeper

import (
	"github.com/dydxprotocol/v4/mocks"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/x/feetiers/keeper"
	"github.com/dydxprotocol/v4/x/feetiers/types"
	statskeeper "github.com/dydxprotocol/v4/x/stats/keeper"
)

func createFeeTiersKeeper(
	stateStore storetypes.CommitMultiStore,
	statsKeeper *statskeeper.Keeper,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)

	k := keeper.NewKeeper(
		cdc,
		statsKeeper,
		storeKey,
	)

	return k, storeKey
}
