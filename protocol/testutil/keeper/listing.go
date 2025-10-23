package keeper

import (
	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	marketmapkeeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assetskeeper "github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	perpetualskeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	subaccountskeeper "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	vaultkeeper "github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	"github.com/stretchr/testify/mock"

	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func ListingKeepers(
	t testing.TB,
	bankKeeper bankkeeper.Keeper,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	mockTimeProvider *mocks.TimeProvider,
	pricesKeeper *priceskeeper.Keeper,
	perpetualsKeeper *perpetualskeeper.Keeper,
	clobKeeper *clobkeeper.Keeper,
	marketMapKeeper *marketmapkeeper.Keeper,
	assetsKeeper *assetskeeper.Keeper,
	bankKeeper_out *bankkeeper.BaseKeeper,
	subaccountsKeeper *subaccountskeeper.Keeper,
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
				registry,
			)
			accountPlusKeeper, _, _ := createAccountPlusKeeper(
				stateStore,
				db,
				cdc,
			)
			bankKeeper_out, _ = createBankKeeper(stateStore, db, cdc, accountsKeeper)
			stakingKeeper, _ := createStakingKeeper(
				stateStore,
				db,
				cdc,
				accountsKeeper,
				bankKeeper_out,
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
			feeTiersKeeper, _ := createFeeTiersKeeper(stateStore, statsKeeper, vaultKeeper, affiliatesKeeper, db, cdc)
			revShareKeeper, _, _ := createRevShareKeeper(stateStore, db, cdc, affiliatesKeeper, feeTiersKeeper, statsKeeper)
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
			assetsKeeper, _ = createAssetsKeeper(
				stateStore,
				db,
				cdc,
				pricesKeeper,
				transientStoreKey,
				true,
			)

			mockMsgSender := &mocks.IndexerMessageSender{}
			mockMsgSender.On("Enabled").Return(true)
			mockIndexerEventManager := indexer_manager.NewIndexerEventManager(mockMsgSender, transientStoreKey, true)

			blockTimeKeeper, _ := createBlockTimeKeeper(stateStore, db, cdc)
			rewardsKeeper, _ := createRewardsKeeper(
				stateStore,
				assetsKeeper,
				bankKeeper_out,
				feeTiersKeeper,
				pricesKeeper,
				mockIndexerEventManager,
				db,
				cdc,
			)
			// Create subaccounts keeper first with nil leverageKeeper
			subaccountsKeeper, _ = createSubaccountsKeeper(
				stateStore,
				db,
				cdc,
				assetsKeeper,
				bankKeeper_out,
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
				bankKeeper_out,
				feeTiersKeeper,
				perpetualsKeeper,
				pricesKeeper,
				statsKeeper,
				rewardsKeeper,
				affiliatesKeeper,
				subaccountsKeeper,
				revShareKeeper,
				accountPlusKeeper,
				mockIndexerEventManager,
				transientStoreKey,
			)

			// Create the listing keeper
			keeper, storeKey, _ = createListingKeeper(
				stateStore,
				db,
				cdc,
				mockIndexerEventManager,
				pricesKeeper,
				perpetualsKeeper,
				clobKeeper,
				marketMapKeeper,
				subaccountsKeeper,
				vaultKeeper,
			)

			return []GenesisInitializer{keeper}
		},
	)

	return ctx, keeper, storeKey, mockTimeProvider, pricesKeeper, perpetualsKeeper,
		clobKeeper, marketMapKeeper, assetsKeeper, bankKeeper_out, subaccountsKeeper
}

func createListingKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	indexerEventManager indexer_manager.IndexerEventManager,
	pricesKeeper *priceskeeper.Keeper,
	perpetualsKeeper *perpetualskeeper.Keeper,
	clobKeeper *clobkeeper.Keeper,
	marketMapKeeper *marketmapkeeper.Keeper,
	subaccountsKeeper *subaccountskeeper.Keeper,
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
		indexerEventManager,
		pricesKeeper,
		clobKeeper,
		marketMapKeeper,
		perpetualsKeeper,
		subaccountsKeeper,
		vaultkeeper,
	)

	return k, storeKey, mockTimeProvider
}
