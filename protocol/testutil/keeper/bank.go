package keeper

import (
	db "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

func createBankKeeper(
	stateStore storetypes.CommitMultiStore,
	db *db.MemDB,
	cdc *codec.ProtoCodec,
	accountKeeper *authkeeper.AccountKeeper,
) (*keeper.BaseKeeper, storetypes.StoreKey) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	k := keeper.NewBaseKeeper(
		cdc,
		storeKey,
		accountKeeper,
		map[string]bool{},
		lib.GovModuleAddress.String(),
	)

	return &k, storeKey
}
