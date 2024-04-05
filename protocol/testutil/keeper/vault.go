package keeper

import (
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	delaymsgtypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func VaultKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
) {
	ctx = initKeepers(t, func(
		db *dbm.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		keeper, storeKey = createVaultKeeper(stateStore, db, cdc, transientStoreKey)
		return []GenesisInitializer{keeper}
	})

	return ctx, keeper, storeKey
}

func createVaultKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	transientStoreKey storetypes.StoreKey,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		&mocks.ClobKeeper{},
		&mocks.PerpetualsKeeper{},
		&mocks.PricesKeeper{},
		&mocks.SendingKeeper{},
		&mocks.SubaccountsKeeper{},
		[]string{
			lib.GovModuleAddress.String(),
			delaymsgtypes.ModuleAddress.String(),
		},
	)

	return k, storeKey
}
