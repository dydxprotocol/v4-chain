package keeper

import (
	dbm "github.com/cosmos/cosmos-db"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	delaymsgtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	epochskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/types"
)

func createStatsKeeper(
	stateStore storetypes.CommitMultiStore,
	epochsKeeper *epochskeeper.Keeper,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	transientStoreKey := storetypes.NewTransientStoreKey(types.TransientStoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)

	authorities := []string{
		delaymsgtypes.ModuleAddress.String(),
		lib.GovModuleAddress.String(),
	}
	k := keeper.NewKeeper(
		cdc,
		epochsKeeper,
		storeKey,
		transientStoreKey,
		authorities,
	)

	return k, storeKey
}
