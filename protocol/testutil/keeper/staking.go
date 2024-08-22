package keeper

import (
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

func createStakingKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	accountKeeper *authkeeper.AccountKeeper,
	bankKeeper types.BankKeeper,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	ss := runtime.NewKVStoreService(storeKey)

	k := keeper.NewKeeper(
		cdc,
		ss,
		accountKeeper,
		bankKeeper,
		lib.GovModuleAddress.String(),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
	)

	return k, storeKey
}
