package keeper

import (
	db "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	db *db.MemDB,
	cdc *codec.ProtoCodec,
	registry codectypes.InterfaceRegistry,
) (*keeper.AccountKeeper, storetypes.StoreKey) {
	types.RegisterInterfaces(registry)

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	paramKey := sdk.NewKVStoreKey(paramtypes.StoreKey)
	stateStore.MountStoreWithDB(paramKey, storetypes.StoreTypeIAVL, db)
	paramTKey := sdk.NewTransientStoreKey(paramtypes.TStoreKey)
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
		storeKey,
		types.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32MainPrefix,
		types.NewModuleAddress(govtypes.ModuleName).String(),
	)

	return &k, storeKey
}
