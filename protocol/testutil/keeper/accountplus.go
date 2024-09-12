package keeper

import (
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	keeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
)

func TimestampNonceKeepers(t testing.TB) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	mockTimeProvider *mocks.TimeProvider,
) {
	ctx = initKeepers(
		t, func(
			db *dbm.MemDB,
			registry codectypes.InterfaceRegistry,
			cdc *codec.ProtoCodec,
			stateStore storetypes.CommitMultiStore,
			transientStoreKey storetypes.StoreKey,
		) []GenesisInitializer {
			// Define necessary keepers here for unit tests
			keeper, storeKey, mockTimeProvider =
				createTimestampNonceKeeper(stateStore, db, cdc)

			return []GenesisInitializer{keeper}
		},
	)

	return ctx, keeper, storeKey, mockTimeProvider
}

func createTimestampNonceKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
	*mocks.TimeProvider,
) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	mockTimeProvider := &mocks.TimeProvider{}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		authenticator.NewAuthenticatorManager(),
		[]string{
			lib.GovModuleAddress.String(),
		},
	)

	return k, storeKey, mockTimeProvider
}
