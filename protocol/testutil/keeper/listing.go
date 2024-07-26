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
	clobkeeper "github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	perpetualskeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	marketmapkeeper "github.com/skip-mev/slinky/x/marketmap/keeper"
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
			ck := NewClobKeepersTestContext(t, nil, bankKeeper, indexerEventManager)

			keeper, storeKey, mockTimeProvider =
				createListingKeeper(
					stateStore,
					db,
					cdc,
					ck.PricesKeeper,
					ck.PerpetualsKeeper,
					ck.ClobKeeper,
					ck.MarketMapKeeper,
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
	)

	return k, storeKey, mockTimeProvider
}
