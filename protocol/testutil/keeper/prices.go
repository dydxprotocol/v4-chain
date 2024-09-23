package keeper

import (
	"fmt"
	"testing"

	storetypes "cosmossdk.io/store/types"
	pricefeed_types "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/types"
	pricefeedserver_types "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	delaymsgmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func PricesKeepers(
	t testing.TB,
) (
	ctx sdk.Context,
	keeper *keeper.Keeper,
	storeKey storetypes.StoreKey,
	daemonPriceCache *pricefeedserver_types.MarketToExchangePrices,
	marketToSmoothedPrices types.MarketToSmoothedPrices,
	mockTimeProvider *mocks.TimeProvider,
) {
	ctx = initKeepers(t, func(
		db *dbm.MemDB,
		registry codectypes.InterfaceRegistry,
		cdc *codec.ProtoCodec,
		stateStore storetypes.CommitMultiStore,
		transientStoreKey storetypes.StoreKey,
	) []GenesisInitializer {
		// Define necessary keepers here for unit tests
		keeper, storeKey, daemonPriceCache, marketToSmoothedPrices, mockTimeProvider =
			createPricesKeeper(stateStore, db, cdc, transientStoreKey)

		return []GenesisInitializer{keeper}
	})

	return ctx, keeper, storeKey, daemonPriceCache, marketToSmoothedPrices, mockTimeProvider
}

func createPricesKeeper(
	stateStore storetypes.CommitMultiStore,
	db *dbm.MemDB,
	cdc *codec.ProtoCodec,
	transientStoreKey storetypes.StoreKey,
) (
	*keeper.Keeper,
	storetypes.StoreKey,
	*pricefeedserver_types.MarketToExchangePrices,
	types.MarketToSmoothedPrices,
	*mocks.TimeProvider,
) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	daemonPriceCache := pricefeedserver_types.NewMarketToExchangePrices(pricefeed_types.MaxPriceAge)

	marketToSmoothedPrices := types.NewMarketToSmoothedSpotPrices(types.SmoothedPriceTrackingBlockHistoryLength)

	mockTimeProvider := &mocks.TimeProvider{}

	mockMsgSender := &mocks.IndexerMessageSender{}
	mockMsgSender.On("Enabled").Return(true)
	mockMsgSender.On("SendOnchainData", mock.Anything).Return()
	mockMsgSender.On("SendOffchainData", mock.Anything).Return()

	mockIndexerEventsManager := indexer_manager.NewIndexerEventManager(mockMsgSender, transientStoreKey, true)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		daemonPriceCache,
		marketToSmoothedPrices,
		mockTimeProvider,
		mockIndexerEventsManager,
		[]string{
			delaymsgmoduletypes.ModuleAddress.String(),
			lib.GovModuleAddress.String(),
		},
	)

	return k, storeKey, daemonPriceCache, marketToSmoothedPrices, mockTimeProvider
}

// CreateTestMarkets creates a standard set of test markets for testing.
// This function assumes no markets exist and will create markets as id `0`, `1`, and `2`, ... using markets
// defined in constants.TestMarkets.
func CreateTestMarkets(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
	for i, marketParam := range constants.TestMarketParams {
		_, err := k.CreateMarket(
			ctx,
			marketParam,
			constants.TestMarketPrices[i],
		)
		require.NoError(t, err)
		err = k.UpdateSpotAndPnlMarketPrices(ctx, &types.MarketPriceUpdate{
			MarketId:  uint32(i),
			SpotPrice: constants.TestMarketPrices[i].SpotPrice,
			PnlPrice:  constants.TestMarketPrices[i].PnlPrice,
		})
		require.NoError(t, err)
	}
}

// CreateNMarkets creates N MarketParam, MarketPrice pairs for testing.
func CreateNMarkets(t *testing.T, ctx sdk.Context, keeper *keeper.Keeper, n int) []types.MarketParamPrice {
	items := make([]types.MarketParamPrice, n)
	numExistingMarkets := GetNumMarkets(t, ctx, keeper)
	for i := range items {
		items[i].Param.Id = uint32(i) + numExistingMarkets
		items[i].Param.Pair = fmt.Sprintf("%v", i)
		items[i].Param.Exponent = int32(i)
		items[i].Param.ExchangeConfigJson = ""
		items[i].Param.MinExchanges = uint32(1)
		items[i].Param.MinPriceChangePpm = uint32(i + 1)
		items[i].Price.Id = uint32(i) + numExistingMarkets
		items[i].Price.Exponent = int32(i)
		items[i].Price.SpotPrice = uint64(1_000 + i)
		items[i].Price.PnlPrice = uint64(1_000 + i)
		items[i].Param.ExchangeConfigJson = "{}" // Use empty, valid JSON for testing.

		_, err := keeper.CreateMarket(
			ctx,
			items[i].Param,
			items[i].Price,
		)
		require.NoError(t, err)
		items[i].Price, err = keeper.GetMarketPrice(ctx, items[i].Param.Id)
		require.NoError(t, err)
	}

	return items
}

// AssertPriceUpdateEventsInIndexerBlock verifies that the market update has a corresponding price update
// event included in the Indexer block message.
func AssertPriceUpdateEventsInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	updatedMarketPrices []types.MarketPrice,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	var expectedEvents []*indexerevents.MarketEventV1
	for _, updatedMarketPrice := range updatedMarketPrices {
		expectedEvent := keeper.GenerateMarketPriceUpdateIndexerEvent(updatedMarketPrice)
		expectedEvents = append(expectedEvents, expectedEvent)
	}
	for _, expectedEvent := range expectedEvents {
		require.Contains(t, marketEvents, expectedEvent)
	}
}

// AssertMarketEventsNotInIndexerBlock verifies that no market events were included in the Indexer block message.
func AssertMarketEventsNotInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
) {
	indexerMarketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	require.Equal(t, 0, len(indexerMarketEvents))
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
	t *testing.T,
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
	t *testing.T,
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
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	updatedMarketPrice types.MarketPrice,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	expectedEvent := indexerevents.NewMarketPriceUpdateEvent(
		updatedMarketPrice.Id,
		updatedMarketPrice.SpotPrice,
		updatedMarketPrice.PnlPrice,
	)
	require.Contains(t, marketEvents, expectedEvent)
}

// CreateTestPriceMarkets is a test utility function that creates list of given
// price markets in state.
func CreateTestPriceMarkets(
	t *testing.T,
	ctx sdk.Context,
	pricesKeeper *keeper.Keeper,
	markets []types.MarketParamPrice,
) {
	// Create a new market param and price.
	marketId := uint32(0)
	for _, m := range markets {
		_, err := pricesKeeper.CreateMarket(
			ctx,
			m.Param,
			m.Price,
		)
		require.NoError(t, err)
		marketId++
	}
}

func GetNumMarkets(t *testing.T, ctx sdk.Context, keeper *keeper.Keeper) uint32 {
	allMarkets, err := keeper.GetAllMarketParamPrices(ctx)
	require.NoError(t, err)
	return lib.MustConvertIntegerToUint32(len(allMarkets))
}
