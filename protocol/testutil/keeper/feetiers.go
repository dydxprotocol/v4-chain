package keeper

import (
	"github.com/dydxprotocol/v4-chain/protocol/mocks"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	statskeeper "github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"
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

	authorities := []string{
		authtypes.NewModuleAddress(delaymsgtypes.ModuleName).String(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	}
	k := keeper.NewKeeper(
		cdc,
		statsKeeper,
		storeKey,
		authorities,
	)

	return k, storeKey
}
