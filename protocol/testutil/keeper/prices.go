package keeper

import (
	"fmt"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedtypes "github.com/dydxprotocol/v4/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/prices/keeper"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
)

func PricesKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	mockTimeProvider *mocks.TimeProvider,
) {
	ctx = initKeepers(t, func(
		db *tmdb.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		keeper, storeKey, indexPriceCache, mockTimeProvider = createPricesKeeper(stateStore, db, cdc, transientStoreKey)

		return []GenesisInitializer{keeper}
	})

	return ctx, keeper, storeKey, indexPriceCache, mockTimeProvider
}

func createPricesKeeper(
	stateStore storetypes.CommitMultiStore,
	db *tmdb.MemDB,
	cdc *codec.ProtoCodec,
	transientStoreKey storetypes.StoreKey,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
	*pricefeedtypes.MarketToExchangePrices,
	*mocks.TimeProvider,
) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	indexPriceCache := pricefeedtypes.NewMarketToExchangePrices()

	mockTimeProvider := &mocks.TimeProvider{}

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	mockIndexerEventsManager := indexer_manager.NewIndexerEventManager(mockMsgSender, transientStoreKey)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		indexPriceCache,
		mockTimeProvider,
		mockIndexerEventsManager,
	)

	return k, storeKey, indexPriceCache, mockTimeProvider
}

// CreateTestExchangeFeeds creates exchanges for testing.
// This function assumes no exchanges exist and will create exchanges as id `0`, `1`, and `2`, ... using markets
// defined in constants.TestExchangeFeeds.
func CreateTestExchangeFeeds(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
	for _, exchange := range constants.TestExchangeFeeds {
		_, err := k.CreateExchangeFeed(ctx, exchange.Name, exchange.Memo)
		require.NoError(t, err)
	}
}

// CreateTestMarketsAndExchangeFeeds creates markets and exchanges for testing.
// This function assumes no markets exist and will create markets as id `0`, `1`, and `2`, ... using markets
// defined in constants.TestMarkets.
// This function assumes no exchanges exist and will create exchanges as id `0`, `1`, and `2`, ... using markets
// defined in constants.TestExchangeFeeds.
func CreateTestMarketsAndExchangeFeeds(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
	CreateTestExchangeFeeds(t, ctx, k)

	for i, market := range constants.TestMarkets {
		_, err := k.CreateMarket(
			ctx,
			market.Pair,
			market.Exponent,
			market.Exchanges,
			market.MinExchanges,
			market.MinPriceChangePpm,
		)
		require.NoError(t, err)
		err = k.UpdateMarketPrices(ctx, []*types.MsgUpdateMarketPrices_MarketPrice{
			{
				MarketId: uint32(i),
				Price:    market.Price,
			},
		},
			true)
		require.NoError(t, err)
	}
}

// CreateNMarketsWithExchangeFeeds creates specified number of markets along with a few exchanges for testing.
func CreateNMarketsWithExchangeFeeds(
	t *testing.T,
	ctx sdk.Context,
	k *keeper.Keeper,
	numMarkets int,
) error {
	// This should create exchange feeds with id `0`, `1` and `2`
	CreateTestExchangeFeeds(t, ctx, k)

	for i := 0; i < numMarkets; i++ {
		if _, err := k.CreateMarket(
			ctx,
			fmt.Sprintf("Market-%v", i),
			int32(0),
			[]uint32{0, 1},
			uint32(1),
			uint32(50),
		); err != nil {
			return err
		}
	}

	return nil
}
