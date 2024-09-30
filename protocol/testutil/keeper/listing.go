package keeper

import (
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	perpetualskeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	vaultkeeper "github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	"github.com/stretchr/testify/mock"

	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func ListingKeepers(
	t testing.TB,
	bankKeeper bankkeeper.Keeper,
	indexerEventManager indexer_manager.IndexerEventManager,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	mockTimeProvider *mocks.TimeProvider,
	pricesKeeper *priceskeeper.Keeper,
	perpetualsKeeper *perpetualskeeper.Keeper,
	clobKeeper *clobkeeper.Keeper,
	marketMapKeeper *marketmapkeeper.Keeper,
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
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			memClob.On("CreateOrderbook", mock.Anything, mock.Anything).Return(nil)
			epochsKeeper, _ := createEpochsKeeper(stateStore, db, cdc)

			accountsKeeper, _ := createAccountKeeper(
				stateStore,
				db,
				cdc,
				registry)
			bankKeeper, _ = createBankKeeper(stateStore, db, cdc, accountsKeeper)
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
			vaultKeeper, _ := createVaultKeeper(
				stateStore,
				db,
				cdc,
				transientStoreKey,
			)
			feeTiersKeeper, _ := createFeeTiersKeeper(
				stateStore,
				statsKeeper,
				vaultKeeper,
				affiliatesKeeper,
				db,
				cdc,
			)
			revShareKeeper, _, _ := createRevShareKeeper(stateStore, db, cdc, affiliatesKeeper, feeTiersKeeper)
			marketMapKeeper, _ = createMarketMapKeeper(stateStore, db, cdc)
			pricesKeeper, _, _, mockTimeProvider = createPricesKeeper(
				stateStore,
				db,
				cdc,
				transientStoreKey,
				revShareKeeper,
				marketMapKeeper,
			)
			// Mock time provider response for market creation.
			mockTimeProvider.On("Now").Return(constants.TimeT)
			perpetualsKeeper, _ = createPerpetualsKeeper(
				stateStore,
				db,
				cdc,
				pricesKeeper,
				epochsKeeper,
				transientStoreKey,
			)
			assetsKeeper, _ := createAssetsKeeper(
				stateStore,
				db,
				cdc,
				pricesKeeper,
				transientStoreKey,
				true,
			)

			blockTimeKeeper, _ := createBlockTimeKeeper(stateStore, db, cdc)
			rewardsKeeper, _ := createRewardsKeeper(
				stateStore,
				assetsKeeper,
				bankKeeper,
				feeTiersKeeper,
				pricesKeeper,
				indexerEventManager,
				db,
				cdc,
			)
			subaccountsKeeper, _ := createSubaccountsKeeper(
				stateStore,
				db,
				cdc,
				assetsKeeper,
				bankKeeper,
				perpetualsKeeper,
				blockTimeKeeper,
				transientStoreKey,
				true,
			)
			clobKeeper, _, _ = createClobKeeper(
				stateStore,
				db,
				cdc,
				memClob,
				assetsKeeper,
				blockTimeKeeper,
				bankKeeper,
				feeTiersKeeper,
				perpetualsKeeper,
				pricesKeeper,
				statsKeeper,
				rewardsKeeper,
				affiliatesKeeper,
				subaccountsKeeper,
				revShareKeeper,
				indexerEventManager,
				transientStoreKey,
			)
			// Create the listing keeper
			keeper, storeKey, _ = createListingKeeper(
				stateStore,
				db,
				cdc,
				pricesKeeper,
				perpetualsKeeper,
				clobKeeper,
				marketMapKeeper,
				vaultKeeper,
			)

			return []GenesisInitializer{keeper}
		},
	)

	return ctx, keeper, storeKey, mockTimeProvider, pricesKeeper, perpetualsKeeper, clobKeeper, marketMapKeeper
}

func createListingKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	pricesKeeper *priceskeeper.Keeper,
	perpetualsKeeper *perpetualskeeper.Keeper,
	clobKeeper *clobkeeper.Keeper,
	marketMapKeeper *marketmapkeeper.Keeper,
	vaultkeeper *vaultkeeper.Keeper,
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
		[]string{
			lib.GovModuleAddress.String(),
		},
		pricesKeeper,
		clobKeeper,
		marketMapKeeper,
		perpetualsKeeper,
		vaultkeeper,
	)

	return k, storeKey, mockTimeProvider
}
