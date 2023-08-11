package keeper

import (
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	indexer_manager "github.com/dydxprotocol/v4/indexer/indexer_manager"
	"math/big"
	"sort"
	"time"

	gometrics "github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// CreateMarket creates a new market in the store.
// Note: this func sorts `exchanges` in ascending order, so the slice may get updated.
func (k Keeper) CreateMarket(
	ctx sdk.Context,
	pair string,
	exponent int32,
	exchanges []uint32,
	minExchanges uint32,
	minPriceChangePpm uint32,
) (types.Market, error) {
	sort.Slice(exchanges, func(i, j int) bool { return exchanges[i] < exchanges[j] })

	// Validate input.
	if err := k.validateMarketFields(
		ctx,
		pair,
		exchanges,
		minExchanges,
		minPriceChangePpm,
	); err != nil {
		return types.Market{}, err
	}

	// Get the `nextId`.
	nextId := k.GetNumMarkets(ctx)

	// Create the market.
	market := types.Market{
		Id:                nextId,
		Pair:              pair,
		Exponent:          exponent,
		Exchanges:         exchanges,
		MinExchanges:      minExchanges,
		MinPriceChangePpm: minPriceChangePpm,
	}

	// Store the new market.
	b := k.cdc.MustMarshal(&market)
	marketStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MarketKeyPrefix))
	marketStore.Set(types.MarketKey(market.Id), b)

	// Store the new `numMarkets`.
	k.setNumMarkets(ctx, nextId+1)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewMarketCreateEvent(
				market.Id,
				market.Pair,
				market.MinPriceChangePpm,
				market.Exponent,
			),
		),
	)
	return market, nil
}

// ModifyMarket modifies an existing market in the store.
// Note: this func sorts `exchanges` in ascending order, so the slice may get updated.
func (k Keeper) ModifyMarket(
	ctx sdk.Context,
	id uint32,
	pair string,
	exchanges []uint32,
	minExchanges uint32,
	minPriceChangePpm uint32,
) (types.Market, error) {
	sort.Slice(exchanges, func(i, j int) bool { return exchanges[i] < exchanges[j] })

	// Validate input.
	if err := k.validateMarketFields(
		ctx,
		pair,
		exchanges,
		minExchanges,
		minPriceChangePpm,
	); err != nil {
		return types.Market{}, err
	}

	// Get necessary store.
	marketStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MarketKeyPrefix))

	// Get market.
	market, err := k.GetMarket(ctx, id)
	if err != nil {
		return types.Market{}, err
	}

	// Modify market.
	market.Pair = pair
	market.Exchanges = exchanges
	market.MinExchanges = minExchanges
	market.MinPriceChangePpm = minPriceChangePpm

	// Store the modified market.
	b := k.cdc.MustMarshal(&market)
	marketStore.Set(types.MarketKey(market.Id), b)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewMarketModifyEvent(
				market.Id,
				market.Pair,
				market.MinPriceChangePpm,
			),
		),
	)

	return market, nil
}

// validateMarketFields validates field inputs.
func (k Keeper) validateMarketFields(
	ctx sdk.Context,
	pair string,
	exchanges []uint32,
	minExchanges uint32,
	minPriceChangePpm uint32,
) error {
	// Validate pair.
	if pair == "" {
		return sdkerrors.Wrap(types.ErrInvalidInput, "Pair cannot be empty")
	}

	if minExchanges == 0 {
		return types.ErrZeroMinExchanges
	}

	// Verify we have enough exchanges.
	if uint32(len(exchanges)) < minExchanges {
		return types.ErrTooFewExchanges
	}

	// Validate min price change.
	if minPriceChangePpm == 0 || minPriceChangePpm >= lib.MaxPriceChangePpm {
		return sdkerrors.Wrapf(
			types.ErrInvalidInput,
			"Min price change in parts-per-million must be greater than 0 and less than %d",
			lib.MaxPriceChangePpm)
	}

	// Validate exchanges.
	for i, exchange := range exchanges {
		// Don't allow duplicates.
		if i > 0 && exchanges[i-1] >= exchange {
			return sdkerrors.Wrap(types.ErrDuplicateExchanges, lib.Uint32ToString(exchange))
		}

		// Verify exchanges exist.
		if _, err := k.GetExchangeFeed(ctx, exchange); err != nil {
			return err
		}
	}

	return nil
}

// UpdateMarketPrices updates the prices for markets.
func (k Keeper) UpdateMarketPrices(
	ctx sdk.Context,
	updates []*types.MsgUpdateMarketPrices_MarketPrice,
	sendIndexerPriceUpdates bool,
) error {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.UpdateMarketPrices,
		metrics.Latency,
	)

	// Get necessary store.
	marketStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MarketKeyPrefix))
	updatedMarkets := make([]types.Market, 0, len(updates))

	for _, update := range updates {
		// Get market.
		market, err := k.GetMarket(ctx, update.MarketId)
		if err != nil {
			return err
		}

		// Report rate of price change for each market.
		if diffRate, err := lib.ChangeRateUint64(market.Price, update.Price); err == nil {
			telemetry.SetGaugeWithLabels(
				[]string{types.ModuleName, metrics.PriceChangeRate},
				diffRate,
				[]gometrics.Label{ // To track per market, include the id as a label.
					pricefeedmetrics.GetLabelForMarketId(market.Id),
				},
			)
		}

		// Update price.
		market.Price = update.Price
		updatedMarkets = append(updatedMarkets, market)

		// Report the oracle price.
		updatedPrice, _ := lib.BigMulPow10(
			new(big.Int).SetUint64(update.Price),
			market.Exponent,
		).Float32()
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.CurrentMarketPrices},
			updatedPrice,
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(market.Id),
			},
		)
	}

	// Writes to the store are delayed so that the updates are atomically applied to state.
	for _, market := range updatedMarkets {
		// Store the modified market.
		b := k.cdc.MustMarshal(&market)
		marketStore.Set(types.MarketKey(market.Id), b)

		// Monitor the last block a market is updated.
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.LastPriceUpdateForMarketBlock},
			float32(ctx.BlockHeight()),
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(market.Id),
			},
		)
	}

	if sendIndexerPriceUpdates {
		marketPriceUpdates := GenerateMarketPriceUpdateEvents(updatedMarkets)
		for _, update := range marketPriceUpdates {
			k.GetIndexerEventManager().AddTxnEvent(
				ctx,
				indexerevents.SubtypeMarket,
				indexer_manager.GetB64EncodedEventMessage(
					update,
				),
			)
		}
	}

	return nil
}

// GetMarket returns a market from its id.
func (k Keeper) GetMarket(
	ctx sdk.Context,
	id uint32,
) (types.Market, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MarketKeyPrefix))

	b := store.Get(types.MarketKey(id))
	if b == nil {
		return types.Market{}, sdkerrors.Wrap(types.ErrMarketDoesNotExist, lib.Uint32ToString(id))
	}

	var market = types.Market{}
	k.cdc.MustUnmarshal(b, &market)
	return market, nil
}

// GetAllMarkets returns all markets.
func (k Keeper) GetAllMarkets(ctx sdk.Context) []types.Market {
	num := k.GetNumMarkets(ctx)
	markets := make([]types.Market, num)

	for i := uint32(0); i < num; i++ {
		market, err := k.GetMarket(ctx, i)
		if err != nil {
			panic(err)
		}

		markets[i] = market
	}

	return markets
}

// GetNumMarkets returns the total number of markets.
func (k Keeper) GetNumMarkets(
	ctx sdk.Context,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	var numMarketsBytes []byte = store.Get(types.KeyPrefix(types.NumMarketsKey))
	return lib.BytesToUint32(numMarketsBytes)
}

// setNumberMarkets sets the number of markets with a new value.
func (k Keeper) setNumMarkets(
	ctx sdk.Context,
	newValue uint32,
) {
	// Get necessary stores.
	store := ctx.KVStore(k.storeKey)

	// Set `numMarkets`.
	store.Set(types.KeyPrefix(types.NumMarketsKey), lib.Uint32ToBytes(newValue))
}
