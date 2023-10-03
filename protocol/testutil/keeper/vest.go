package keeper

import (
	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	blocktimekeeper "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/keeper"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"

	"testing"
)

func VestKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	vestKeeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	bankKeeper *bankkeeper.BaseKeeper,
	blocktimeKeeper *blocktimekeeper.Keeper,
	authorities []string,
) {
	ctx = initKeepers(t, func(
		db *tmdb.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		authorities = []string{
			authtypes.NewModuleAddress(bridgetypes.ModuleName).String(),
			authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		}
		accountKeeper, _ := createAccountKeeper(stateStore, db, cdc, registry)
		bankKeeper, _ = createBankKeeper(stateStore, db, cdc, accountKeeper)
		blocktimeKeeper, _ = createBlockTimeKeeper(stateStore, db, cdc)
		vestKeeper, storeKey = createVestKeeper(
			stateStore,
			db,
			cdc,
			bankKeeper,
			blocktimeKeeper,
			authorities,
		)
		return []GenesisInitializer{blocktimeKeeper}
	})
	return ctx, vestKeeper, storeKey, bankKeeper, blocktimeKeeper, authorities
}

func createVestKeeper(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc codec.BinaryCodec,
	bankKeeper *bankkeeper.BaseKeeper,
	blocktimeKeeper *blocktimekeeper.Keeper,
	authorities []string,
) (*keeper.Keeper, *storetypes.KVStoreKey) {
	vestStoreKey := storetypes.NewKVStoreKey(types.StoreKey)
	vestKeeper := keeper.NewKeeper(cdc, vestStoreKey, bankKeeper, blocktimeKeeper, authorities)
	stateStore.MountStoreWithDB(vestStoreKey, storetypes.StoreTypeIAVL, db)
	return vestKeeper, vestStoreKey
}
