package keeper_test

import (
	"fmt"
	"github.com/dydxprotocol/v4/indexer/common"
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/indexer/indexer_manager"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/testutil/nullify"
	"github.com/dydxprotocol/v4/x/prices/keeper"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
)

const (
	ErrorMsgInvalidMinPriceChange = "Min price change in parts-per-million must be greater than 0 and less than 10000"
)

func createNMarkets(t *testing.T, keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Market {
	// This should create exchanges with id `0`, `1` and `2`
	keepertest.CreateTestExchangeFeeds(t, ctx, keeper)

	items := make([]types.Market, n)
	for i := range items {
		items[i].Id = uint32(i)
		items[i].Pair = fmt.Sprintf("%v", i)
		items[i].Exponent = int32(i)
		items[i].Exchanges = []uint32{0, 1}
		items[i].MinExchanges = uint32(1)
		items[i].MinPriceChangePpm = uint32(i + 1)

		_, err := keeper.CreateMarket(
			ctx,
			items[i].Pair,
			items[i].Exponent,
			items[i].Exchanges,
			items[i].MinExchanges,
			items[i].MinPriceChangePpm,
		)

		require.NoError(t, err)
	}

	return items
}

func createNMarketPriceUpdates(
	keeper *keeper.Keeper,
	ctx sdk.Context,
	n int,
) []*types.MsgUpdateMarketPrices_MarketPrice {
	items := make([]*types.MsgUpdateMarketPrices_MarketPrice, n)
	for i := range items {
		items[i] = &types.MsgUpdateMarketPrices_MarketPrice{
			MarketId: uint32(i),
			Price:    uint64(i),
		}
	}

	return items
}

// assertPriceUpdateEventsInIndexerBlock verifies that the market update has a corresponding price update
// event included in the Indexer block message.
func assertPriceUpdateEventsInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	updatedMarkets []types.Market,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	expectedEvents := keeper.GenerateMarketPriceUpdateEvents(updatedMarkets)
	for _, expectedEvent := range expectedEvents {
		require.Contains(t, marketEvents, expectedEvent)
	}
}

// getMarketEventsFromIndexerBlock returns the market events from the Indexer Block event Kafka message.
func getMarketEventsFromIndexerBlock(
	ctx sdk.Context,
	k *keeper.Keeper,
) []*indexerevents.MarketEvent {
	block := k.GetIndexerEventManager().ProduceBlock(ctx)
	var marketEvents []*indexerevents.MarketEvent
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeMarket {
			continue
		}
		bytes := indexer_manager.GetBytesFromEventData(event.Data)
		unmarshaler := common.UnmarshalerImpl{}
		var marketEvent indexerevents.MarketEvent
		err := unmarshaler.Unmarshal(bytes, &marketEvent)
		if err != nil {
			panic(err)
		}
		marketEvents = append(marketEvents, &marketEvent)
	}
	return marketEvents
}

// assertNoMarketEventsFromIndexerBlock returns true if there are no price update events in the
// Indexer block message.
func assertNoMarketEventsFromIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	exists := false
	for _, event := range marketEvents {
		if _, ok := event.Event.(*indexerevents.MarketEvent_PriceUpdate); ok {
			exists = true
		}
	}
	require.False(t, exists)
}

// assertMarketModifyEventInIndexerBlock verifies that the market update has a corresponding market modify
// event included in the Indexer block message.
func assertMarketModifyEventInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	updatedMarket types.Market,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	expectedEvent := indexerevents.NewMarketModifyEvent(
		updatedMarket.Id,
		updatedMarket.Pair,
		updatedMarket.MinPriceChangePpm,
	)
	require.Contains(t, marketEvents, expectedEvent)
}

// assertMarketCreateEventInIndexerBlock verifies that the market create has a corresponding market create
// event included in the Indexer block message.
func assertMarketCreateEventInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	createdMarket types.Market,
) {
	marketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	expectedEvent := indexerevents.NewMarketCreateEvent(
		createdMarket.Id,
		createdMarket.Pair,
		createdMarket.MinPriceChangePpm,
		createdMarket.Exponent,
	)
	require.Contains(t, marketEvents, expectedEvent)
}

// assertMarketEventsNotInIndexerBlock verifies that no market events were included in the Indexer block message.
func assertMarketEventsNotInIndexerBlock(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
) {
	indexerMarketEvents := getMarketEventsFromIndexerBlock(ctx, k)
	require.Equal(t, 0, len(indexerMarketEvents))
}

func TestCreateMarket(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	createNExchangeFeeds(t, keeper, ctx, 2)

	market, err := keeper.CreateMarket(
		ctx,
		constants.BtcUsdPair,
		int32(-6),
		[]uint32{0, 1},
		uint32(2),
		uint32(9999),
	)

	require.NoError(t, err)
	require.Equal(t, uint32(0), market.Id)
	require.Equal(t, constants.BtcUsdPair, market.Pair)
	require.Equal(t, int32(-6), market.Exponent)
	require.Equal(t, []uint32{0, 1}, market.Exchanges)
	require.Equal(t, uint32(2), market.MinExchanges)
	require.Equal(t, uint32(9999), market.MinPriceChangePpm)
	assertMarketCreateEventInIndexerBlock(t, keeper, ctx, market)
}

func TestCreateMarket_Errors(t *testing.T) {
	tests := map[string]struct {
		// Setup
		pair                 string
		exchanges            []uint32
		minExchanges         uint32
		minPriceChangePpm    uint32
		numExchangesToCreate int

		// Expected
		expectedErr string
	}{
		"Empty pair": {
			pair:                 "", // pair cannot be empty
			exchanges:            []uint32{0, 1},
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrInvalidInput, "Pair cannot be empty").Error(),
		},
		"Invalid min price change: zero": {
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 1},
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(0), // must be > 0
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrInvalidInput, ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Invalid min price change: ten thousand": {
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 1},
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(10_000), // must be < 10,000
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrInvalidInput, ErrorMsgInvalidMinPriceChange).Error(),
		},
		"ExchangeFeed does not exist": {
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 1, 6}, // exchange with id `6` does not exist
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrExchangeFeedDoesNotExist, "6").Error(),
		},
		"Too few exchanges": {
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0}, // does not meet min exchanges
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numExchangesToCreate: 1,
			expectedErr:          types.ErrTooFewExchanges.Error(),
		},
		"Min exchanges cannot be zero": {
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{},
			minExchanges:         uint32(0), // cannot be zero
			minPriceChangePpm:    uint32(50),
			numExchangesToCreate: 1,
			expectedErr:          types.ErrZeroMinExchanges.Error(),
		},
		"Duplicate exchanges": {
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 0}, // duplicates
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrDuplicateExchanges, "0").Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			createNExchangeFeeds(t, keeper, ctx, tc.numExchangesToCreate)
			_, err := keeper.CreateMarket(
				ctx,
				tc.pair,
				int32(-6),
				tc.exchanges,
				tc.minExchanges,
				tc.minPriceChangePpm,
			)
			require.EqualError(t, err, tc.expectedErr)
			assertMarketEventsNotInIndexerBlock(t, keeper, ctx)
		})
	}
}

func TestUpdateMarketPrices(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	items := createNMarkets(t, keeper, ctx, 10)
	require.Equal(t, uint32(10), keeper.GetNumMarkets(ctx))

	// Create firstPriceUpdates which should be overwritten by secondPriceUpdates
	firstPriceUpdates := createNMarketPriceUpdates(keeper, ctx, 10)
	secondPriceUpdates := createNMarketPriceUpdates(keeper, ctx, 10)
	for _, pu := range secondPriceUpdates {
		pu.Price = 10 + (pu.Price * 10)
	}

	priceUpdates := append(firstPriceUpdates, secondPriceUpdates...)
	err := keeper.UpdateMarketPrices(
		ctx,
		priceUpdates,
		true,
	)
	require.NoError(t, err)

	markets := make([]types.Market, 10)
	for i, item := range items {
		market, err := keeper.GetMarket(ctx, item.Id)
		require.NoError(t, err)
		require.Equal(t,
			secondPriceUpdates[i].Price,
			market.Price,
		)
		markets = append(markets, market)
	}

	assertPriceUpdateEventsInIndexerBlock(t, keeper, ctx, markets)
}

func TestUpdateMarketPricesGenesis(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	createNMarkets(t, keeper, ctx, 10)
	require.Equal(t, uint32(10), keeper.GetNumMarkets(ctx))

	// Create firstPriceUpdates which should be overwritten by secondPriceUpdates
	firstPriceUpdates := createNMarketPriceUpdates(keeper, ctx, 10)
	secondPriceUpdates := createNMarketPriceUpdates(keeper, ctx, 10)
	for _, pu := range secondPriceUpdates {
		pu.Price = 10 + (pu.Price * 10)
	}

	priceUpdates := append(firstPriceUpdates, secondPriceUpdates...)
	err := keeper.UpdateMarketPrices(
		ctx,
		priceUpdates,
		// set to false for genesis
		false,
	)
	require.NoError(t, err)
	assertNoMarketEventsFromIndexerBlock(t, keeper, ctx)
}

func TestUpdateMarketPrices_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	priceUpdates := createNMarketPriceUpdates(keeper, ctx, 10)
	err := keeper.UpdateMarketPrices(
		ctx,
		priceUpdates,
		true,
	)
	require.EqualError(t, err, "0: Market does not exist")
	assertMarketEventsNotInIndexerBlock(t, keeper, ctx)

	items := createNMarkets(t, keeper, ctx, 10)
	err = keeper.UpdateMarketPrices(
		ctx,
		priceUpdates,
		true,
	)
	require.NoError(t, err)

	for i, item := range items {
		price, err := keeper.GetMarket(ctx, item.Id)
		require.NoError(t, err)
		require.Equal(t,
			priceUpdates[i].Price,
			price.Price,
		)
	}
}

func TestModifyMarket(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	ctx = ctx.WithTxBytes(constants.TestTxBytes)
	items := createNMarkets(t, keeper, ctx, 10)
	createNExchangeFeeds(t, keeper, ctx, 2)
	for i, item := range items {
		// Modify each field arbitrarily and
		// verify the fields were modified in state
		newItem, err := keeper.ModifyMarket(
			ctx,
			item.Id,
			fmt.Sprintf("foo_%v", i),
			[]uint32{0, 1},
			uint32(2),
			uint32(9999-i),
		)
		require.NoError(t, err)
		require.Equal(t, uint32(i), newItem.Id)
		require.Equal(t, fmt.Sprintf("foo_%v", i), newItem.Pair)
		require.Equal(t, []uint32{0, 1}, newItem.Exchanges)
		require.Equal(t, uint32(2), newItem.MinExchanges)
		require.Equal(t, uint32(9999-i), newItem.MinPriceChangePpm)
		assertMarketModifyEventInIndexerBlock(t, keeper, ctx, newItem)
	}
}

func TestModifyMarket_Errors(t *testing.T) {
	tests := map[string]struct {
		// Setup
		targetId             uint32
		pair                 string
		exchanges            []uint32
		minExchanges         uint32
		minPriceChangePpm    uint32
		numMarketsToCreate   int
		numExchangesToCreate int

		// Expected
		expectedErr string
	}{
		"Market not found": {
			targetId:             99, // this market id does not exist
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 1},
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numMarketsToCreate:   2,
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrMarketDoesNotExist, "99").Error(),
		},
		"Empty pair": {
			targetId:             0,
			pair:                 "", // pair cannot be empty
			exchanges:            []uint32{0, 1},
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numMarketsToCreate:   1,
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrInvalidInput, "Pair cannot be empty").Error(),
		},
		"Invalid min price change: zero": {
			targetId:             0,
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 1},
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(0), // must be > 0
			numMarketsToCreate:   1,
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrInvalidInput, ErrorMsgInvalidMinPriceChange).Error(),
		},
		"Invalid min price change: ten thousand": {
			targetId:             0,
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 1},
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(10_000), // must be < 10,000
			numMarketsToCreate:   1,
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrInvalidInput, ErrorMsgInvalidMinPriceChange).Error(),
		},
		"ExchangeFeed does not exist": {
			targetId:             0,
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 1, 6}, // exchange with id `6` does not exist
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numMarketsToCreate:   1,
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrExchangeFeedDoesNotExist, "6").Error(),
		},
		"Too few exchanges": {
			targetId:             0,
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0}, // this does not match minExchanges
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numMarketsToCreate:   1,
			numExchangesToCreate: 1,
			expectedErr:          types.ErrTooFewExchanges.Error(),
		},
		"Min exchanges cannot be zero": {
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{},
			minExchanges:         uint32(0), // cannot be zero
			minPriceChangePpm:    uint32(50),
			numMarketsToCreate:   1,
			numExchangesToCreate: 1,
			expectedErr:          types.ErrZeroMinExchanges.Error(),
		},
		"Duplicate exchanges": {
			targetId:             0,
			pair:                 constants.BtcUsdPair,
			exchanges:            []uint32{0, 0}, // there are duplicates
			minExchanges:         uint32(2),
			minPriceChangePpm:    uint32(50),
			numMarketsToCreate:   1,
			numExchangesToCreate: 2,
			expectedErr:          sdkerrors.Wrap(types.ErrDuplicateExchanges, "0").Error(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			createNMarkets(t, keeper, ctx, tc.numMarketsToCreate)
			createNExchangeFeeds(t, keeper, ctx, tc.numExchangesToCreate)
			_, err := keeper.ModifyMarket(
				ctx,
				tc.targetId,
				tc.pair,
				tc.exchanges,
				tc.minExchanges,
				tc.minPriceChangePpm)
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestGetMarket(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	items := createNMarkets(t, keeper, ctx, 10)
	for _, item := range items {
		rst, err := keeper.GetMarket(ctx, item.Id)
		require.NoError(t, err)
		require.Equal(
			t,
			nullify.Fill(&item), //nolint:staticcheck
			nullify.Fill(&rst),  //nolint:staticcheck
		)
	}
}

func TestGetMarket_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	_, err := keeper.GetMarket(ctx, uint32(0))
	require.EqualError(t, err, "0: Market does not exist")
}

func TestGetAllMarkets(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	items := createNMarkets(t, keeper, ctx, 10)
	require.ElementsMatch(
		t,
		nullify.Fill(items),                     //nolint:staticcheck
		nullify.Fill(keeper.GetAllMarkets(ctx)), //nolint:staticcheck
	)
}

func TestGetAllMarkets_MissingMarket(t *testing.T) {
	ctx, keeper, storeKey, _, _, _ := keepertest.PricesKeepers(t)

	// Write some bad data to the store
	store := ctx.KVStore(storeKey)
	store.Set(types.KeyPrefix(types.NumMarketsKey), lib.Uint32ToBytes(20))

	// Expect a panic
	require.PanicsWithError(
		t,
		"0: Market does not exist",
		func() { keeper.GetAllMarkets(ctx) },
	)
}

func TestGetNumMarkets(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.PricesKeepers(t)
	require.Equal(t, uint32(0), keeper.GetNumMarkets(ctx))

	createNMarkets(t, keeper, ctx, 10)
	require.Equal(t, uint32(10), keeper.GetNumMarkets(ctx))
}
