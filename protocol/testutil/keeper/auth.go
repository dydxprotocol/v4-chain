package keeper

import (
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func createAccountKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	registry codectypes.InterfaceRegistry,
) (*keeper.AccountKeeper, storetypes.StoreKey) {
	types.RegisterInterfaces(registry)

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	paramKey := storetypes.NewKVStoreKey(paramtypes.StoreKey)
	stateStore.MountStoreWithDB(paramKey, storetypes.StoreTypeIAVL, db)
	paramTKey := storetypes.NewTransientStoreKey(paramtypes.TStoreKey)
	stateStore.MountStoreWithDB(paramTKey, storetypes.StoreTypeTransient, db)

	// Create default module account permissions for test.
	maccPerms := map[string][]string{
		minttypes.ModuleName:        {types.Minter},
		bridgetypes.ModuleName:      {types.Minter},
		types.FeeCollectorName:      nil,
		satypes.ModuleName:          nil,
		clobtypes.InsuranceFundName: nil,
	}

	k := keeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		types.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.Bech32MainPrefix),
		sdk.Bech32MainPrefix,
		types.NewModuleAddress(govtypes.ModuleName).String(),
	)

	return &k, storeKey
}
