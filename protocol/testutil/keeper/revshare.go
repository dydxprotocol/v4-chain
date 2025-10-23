package keeper

import (
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	affiliateskeeper "github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	feetierskeeper "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	statskeeper "github.com/dydxprotocol/v4-chain/protocol/x/stats/keeper"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	keeper "github.com/dydxprotocol/v4-chain/protocol/x/revshare/keeper"
)

func RevShareKeepers(t testing.TB) (
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
			epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)

			accountsKeeper, _ := createAccountKeeper(
				stateStore,
				db,
				cdc,
				registry)
			bankKeeper, _ := createBankKeeper(stateStore, db, cdc, accountsKeeper)
			stakingKeeper, _ := createStakingKeeper(
				stateStore,
				db,
				cdc,
				accountsKeeper,
				bankKeeper,
			)
			statsKeeper, _ := createStatsKeeper(
				stateStore,
				epochsKeeper,
				db,
				cdc,
				stakingKeeper,
			)
			affiliatesKeeper, _ := createAffiliatesKeeper(stateStore, db, cdc, statsKeeper, transientStoreKey, true)
			vaultKeeper, _ := createVaultKeeper(stateStore, db, cdc, transientStoreKey)
			feetiersKeeper, _ := createFeeTiersKeeper(stateStore, statsKeeper, vaultKeeper, affiliatesKeeper, db, cdc)
			keeper, storeKey, mockTimeProvider =
				createRevShareKeeper(stateStore, db, cdc, affiliatesKeeper, feetiersKeeper, statsKeeper)

			return []GenesisInitializer{keeper}
		},
	)

	return ctx, keeper, storeKey, mockTimeProvider
}

func createRevShareKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	affiliatesKeeper *affiliateskeeper.Keeper,
	feetiersKeeper *feetierskeeper.Keeper,
	statsKeeper *statskeeper.Keeper,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
	*mocks.TimeProvider,
) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	mockTimeProvider := &mocks.TimeProvider{}

	k := keeper.NewKeeper(
		cdc, storeKey, []string{
			lib.GovModuleAddress.String(),
		},
		*affiliatesKeeper,
		*feetiersKeeper,
		*statsKeeper,
	)

	return k, storeKey, mockTimeProvider
}
