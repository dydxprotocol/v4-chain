package keeper

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	dbm "github.com/cosmos/cosmos-db"

	storetypes "cosmossdk.io/store/types"
	delaymsgtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	statskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
)

func createFeeTiersKeeper(
	stateStore storetypes.CommitMultiStore,
	statsKeeper *statskeeper.Keeper,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
) (*keeper.Keeper, storetypes.StoreKey) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)

	authorities := []string{
		delaymsgtypes.ModuleAddress.String(),
		lib.GovModuleAddress.String(),
	}
	k := keeper.NewKeeper(
		cdc,
		statsKeeper,
		storeKey,
		authorities,
	)

	return k, storeKey
}
