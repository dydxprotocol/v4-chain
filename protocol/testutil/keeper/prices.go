package keeper

import (
	"fmt"
	"testing"

	streaming "github.com/dydxprotocol/v4-chain/protocol/streaming"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"

	"github.com/cosmos/gogoproto/proto"

	storetypes "cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	marketmapkeeper "github.com/dydxprotocol/slinky/x/marketmap/keeper"
	pricefeed_types "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
	pricefeedserver_types "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/pricefeed"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/marketmap"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	delaymsgmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	revsharekeeper "github.com/dydxprotocol/v4-chain/protocol/x/revshare/keeper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func PricesKeepers(t testing.TB) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	indexPriceCache *pricefeedserver_types.MarketToExchangePrices,
	mockTimeProvider *mocks.TimeProvider,
	revShareKeeper *revsharekeeper.Keeper,
	marketMapKeeper *marketmapkeeper.Keeper,
) {
	ctx = initKeepers(t, func(
		db *dbm.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Necessary keeper for testing
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
		revShareKeeper, _, _ = createRevShareKeeper(stateStore, db, cdc, affiliatesKeeper, feetiersKeeper, statsKeeper)
		marketMapKeeper, _ = createMarketMapKeeper(stateStore, db, cdc)
		// Define necessary keepers here for unit tests
		keeper, storeKey, indexPriceCache, mockTimeProvider =
			createPricesKeeper(stateStore, db, cdc, transientStoreKey, revShareKeeper, marketMapKeeper)

		return []GenesisInitializer{keeper}
	})

	return ctx, keeper, storeKey, indexPriceCache, mockTimeProvider, revShareKeeper, marketMapKeeper
}

func createPricesKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	transientStoreKey storetypes.StoreKey,
	revShareKeeper *revsharekeeper.Keeper,
	marketMapKeeper *marketmapkeeper.Keeper,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
	*pricefeedserver_types.MarketToExchangePrices,
	*mocks.TimeProvider,
) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	indexPriceCache := pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)

	mockTimeProvider := &mocks.TimeProvider{}

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	mockMsgSender.On("SendOnchainData", mock.Anything).Return()
	mockMsgSender.On("SendOffchainData", mock.Anything).Return()

	mockIndexerEventsManager := indexer_manager.NewIndexerEventManager(mockMsgSender, transientStoreKey, true)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		indexPriceCache,
		mockTimeProvider,
		mockIndexerEventsManager,
		[]string{
			delaymsgmoduletypes.ModuleAddress.String(),
			lib.GovModuleAddress.String(),
		},
		revShareKeeper,
		marketMapKeeper,
		streaming.NewNoopGrpcStreamingManager(),
	)

	return k, storeKey, indexPriceCache, mockTimeProvider
}

// Convert MarketParams into MarketMap markets and create them in the MarketMap keeper.
func CreateMarketsInMarketMapFromParams(
	t testing.TB,
	ctx sdk.Context,
	mmk *marketmapkeeper.Keeper,
	allMarketParams []types.MarketParam,
) {
	marketMap, err := marketmap.ConstructMarketMapFromParams(allMarketParams)
	require.NoError(t, err)
	for _, market := range marketMap.Markets {
		market.Ticker.Enabled = false
		require.NoError(t, mmk.CreateMarket(ctx, market))
	}
}

// CreateTestMarket creates a market with the given MarketParam and MarketPrice. It creates this market
// in the market map first.
func CreateTestMarket(
	t testing.TB,
	ctx sdk.Context,
	k *keeper.Keeper,
	marketParam types.MarketParam,
	marketPrice types.MarketPrice,
) (types.MarketParam, error) {
	CreateMarketsInMarketMapFromParams(
		t,
		ctx,
		k.MarketMapKeeper.(*marketmapkeeper.Keeper),
		[]types.MarketParam{marketParam},
	)

	return k.CreateMarket(
		ctx,
		marketParam,
		marketPrice,
	)
}

// CreateTestMarkets creates a standard set of test markets for testing.
// This function assumes no markets exist and will create markets as id `0`, `1`, and `2`, ... using markets
// defined in constants.TestMarkets.
func CreateTestMarkets(t testing.TB, ctx sdk.Context, k *keeper.Keeper) {
	for i, marketParam := range constants.TestMarketParams {
		_, err := CreateTestMarket(
			t,
			ctx,
			k,
			marketParam,
			constants.TestMarketPrices[i],
		)
		require.NoError(t, err)
		err = k.UpdateMarketPrices(ctx, []*types.MsgUpdateMarketPrices_MarketPrice{
			{
				MarketId: uint32(i),
				Price:    constants.TestMarketPrices[i].Price,
			},
		})
		require.NoError(t, err)

		// update all markets to not have revenue share
		k.RevShareKeeper.SetMarketMapperRevShareDetails(
			ctx,
			uint32(i),
			revsharetypes.MarketMapperRevShareDetails{ExpirationTs: 0},
		)
	}
}

// CreateNMarkets creates N MarketParam, MarketPrice pairs for testing.
func CreateNMarkets(t testing.TB, ctx sdk.Context, keeper *keeper.Keeper, n int) []types.MarketParamPrice {
	items := make([]types.MarketParamPrice, n)
	numExistingMarkets := GetNumMarkets(t, ctx, keeper)
	for i := range items {
		items[i].Param.Id = uint32(i) + numExistingMarkets
		items[i].Param.Pair = fmt.Sprintf("%v-%v", i, i)
		items[i].Param.Exponent = -int32(i%36 + 1) // must be between -1 and -36
		items[i].Param.ExchangeConfigJson = ""
		items[i].Param.MinExchanges = uint32(1)
		items[i].Param.MinPriceChangePpm = uint32(i + 1)
		items[i].Price.Id = uint32(i) + numExistingMarkets
		items[i].Price.Exponent = -int32(i%36 + 1) // must be between -1 and -36
		items[i].Price.Price = uint64(1_000 + i)
		items[i].Param.ExchangeConfigJson = "{}" // Use empty, valid JSON for testing.

		_, err := CreateTestMarket(t, ctx, keeper, items[i].Param, items[i].Price)
		require.NoError(t, err)
		items[i].Price, err = keeper.GetMarketPrice(ctx, items[i].Param.Id)
		require.NoError(t, err)
	}

	return items
}

// AssertPriceUpdateEventsInIndexerBlock verifies that the market update has a corresponding price update
// event included in the Indexer block message.
func AssertPriceUpdateEventsInIndexerBlock(
	t testing.TB,
	k *keeper.Keeper,
	ctx sdk.Context,
	updatedMarketPrices []types.MarketPrice,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	expectedEvents := keeper.GenerateMarketPriceUpdateIndexerEvents(updatedMarketPrices)
	for _, expectedEvent := range expectedEvents {
		require.Contains(t, marketEvents, expectedEvent)
	}
}

// AssertMarketEventsNotInIndexerBlock verifies that no market events were included in the Indexer block message.
func AssertMarketEventsNotInIndexerBlock(
	t testing.TB,
	k *keeper.Keeper,
	ctx sdk.Context,
) {
	indexerMarketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	require.Equal(t, 0, len(indexerMarketEvents))
}

// AssertNMarketEventsNotInIndexerBlock verifies that N market events were included in the Indexer block message.
func AssertNMarketEventsNotInIndexerBlock(
	t testing.TB,
	k *keeper.Keeper,
	ctx sdk.Context,
	n int,
) {
	indexerMarketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	require.Equal(t, n, len(indexerMarketEvents))
}

// getMarketEventsFromIndexerBlock returns the market events from the Indexer Block event Kafka message.
func getMarketEventsFromIndexerBlock(
	ctx sdk.Context,
	k *keeper.Keeper,
) []*indexerevents.MarketEventV1 {
	block := k.GetIndexerEventManager().ProduceBlock(ctx)
	var marketEvents []*indexerevents.MarketEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeMarket {
			continue
		}
		var marketEvent indexerevents.MarketEventV1
		err := proto.Unmarshal(event.DataBytes, &marketEvent)
		if err != nil {
			panic(err)
		}
		marketEvents = append(marketEvents, &marketEvent)
	}
	return marketEvents
}

// AssertMarketModifyEventInIndexerBlock verifies that the market update has a corresponding market modify
// event included in the Indexer block message.
func AssertMarketModifyEventInIndexerBlock(
	t testing.TB,
	k *keeper.Keeper,
	ctx sdk.Context,
	updatedMarketParam types.MarketParam,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	expectedEvent := indexerevents.NewMarketModifyEvent(
		updatedMarketParam.Id,
		updatedMarketParam.Pair,
		updatedMarketParam.MinPriceChangePpm,
	)
	require.Contains(t, marketEvents, expectedEvent)
}

// AssertMarketCreateEventInIndexerBlock verifies that the market create has a corresponding market create
// event included in the Indexer block message.
func AssertMarketCreateEventInIndexerBlock(
	t testing.TB,
	k *keeper.Keeper,
	ctx sdk.Context,
	createdMarketParam types.MarketParam,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	expectedEvent := indexerevents.NewMarketCreateEvent(
		createdMarketParam.Id,
		createdMarketParam.Pair,
		createdMarketParam.MinPriceChangePpm,
		createdMarketParam.Exponent,
	)
	require.Contains(t, marketEvents, expectedEvent)
}

func AssertMarketPriceUpdateEventInIndexerBlock(
	t testing.TB,
	k *keeper.Keeper,
	ctx sdk.Context,
	updatedMarketPrice types.MarketPrice,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	expectedEvent := indexerevents.NewMarketPriceUpdateEvent(
		updatedMarketPrice.Id,
		updatedMarketPrice.Price,
	)
	require.Contains(t, marketEvents, expectedEvent)
}

// CreateTestPriceMarkets is a test utility function that creates list of given
// price markets in state.
func CreateTestPriceMarkets(
	t testing.TB,
	ctx sdk.Context,
	pricesKeeper *keeper.Keeper,
	markets []types.MarketParamPrice,
) {
	// Create a new market param and price.
	marketId := uint32(0)
	for _, m := range markets {
		_, err := CreateTestMarket(
			t,
			ctx,
			pricesKeeper,
			m.Param,
			m.Price,
		)
		require.NoError(t, err)
		marketId++
	}
}

func GetNumMarkets(t testing.TB, ctx sdk.Context, keeper *keeper.Keeper) uint32 {
	allMarkets, err := keeper.GetAllMarketParamPrices(ctx)
	require.NoError(t, err)
	return lib.MustConvertIntegerToUint32(len(allMarkets))
}
