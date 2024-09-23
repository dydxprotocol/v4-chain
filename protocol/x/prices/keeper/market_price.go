package keeper

import (
	"math/big"
	"sort"
	"time"

	errorsmod "cosmossdk.io/errors"

	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"

	gometrics "github.com/hashicorp/go-metrics"

	"cosmossdk.io/store/prefix"
	pricefeedmetrics "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/metrics"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// getMarketPriceStore returns a prefix store for MarketPrices.
func (k Keeper) getMarketPriceStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketPriceKeyPrefix))
}

func (k Keeper) UpdatePnlPrice(
	ctx sdk.Context,
	update *types.MarketPnlPriceUpdate,
) error {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.UpdateMarketPrices,
		metrics.Latency,
	)

	// Get necessary store.
	marketPriceStore := k.getMarketPriceStore(ctx)

	// Get market price.
	marketPrice, err := k.GetMarketPrice(ctx, update.MarketId)
	if err != nil {
		return err
	}

	// Report rate of price change for each market.
	if diffRate, err := lib.ChangeRateUint64(marketPrice.PnlPrice, update.PnlPrice); err == nil {
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.PriceChangeRate},
			diffRate,
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
			},
		)
	}

	// Update market price.
	marketPrice.PnlPrice = update.PnlPrice

	// Report the oracle price.
	updatedPrice, _ := lib.BigMulPow10(
		new(big.Int).SetUint64(update.PnlPrice),
		marketPrice.Exponent,
	).Float32()
	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, metrics.CurrentMarketPrices},
		updatedPrice,
		[]gometrics.Label{ // To track per market, include the id as a label.
			pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
		},
	)
	// Writes to the store are delayed so that the updates are atomically applied to state.
	// Store the modified market price.
	b := k.cdc.MustMarshal(&marketPrice)
	marketPriceStore.Set(lib.Uint32ToKey(marketPrice.Id), b)

	// Monitor the last block a market price is updated.
	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, metrics.LastPriceUpdateForMarketBlock},
		float32(ctx.BlockHeight()),
		[]gometrics.Label{ // To track per market, include the id as a label.
			pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
		},
	)

	// Generate indexer events.
	priceUpdateIndexerEvents := GenerateMarketPriceUpdateIndexerEvent(marketPrice)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexerevents.MarketEventVersion,
		indexer_manager.GetBytes(
			priceUpdateIndexerEvents,
		),
	)

	return nil
}

func (k Keeper) UpdateSpotPrice(
	ctx sdk.Context,
	update *types.MarketSpotPriceUpdate,
) error {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.UpdateMarketPrices,
		metrics.Latency,
	)

	// Get necessary store.
	marketPriceStore := k.getMarketPriceStore(ctx)

	// Get market price.
	marketPrice, err := k.GetMarketPrice(ctx, update.MarketId)
	if err != nil {
		return err
	}

	// Report rate of price change for each market.
	if diffRate, err := lib.ChangeRateUint64(marketPrice.SpotPrice, update.SpotPrice); err == nil {
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.PriceChangeRate},
			diffRate,
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
			},
		)
	}

	// Update market price.
	marketPrice.SpotPrice = update.SpotPrice

	// Report the oracle price.
	updatedPrice, _ := lib.BigMulPow10(
		new(big.Int).SetUint64(update.SpotPrice),
		marketPrice.Exponent,
	).Float32()
	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, metrics.CurrentMarketPrices},
		updatedPrice,
		[]gometrics.Label{ // To track per market, include the id as a label.
			pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
		},
	)
	// Writes to the store are delayed so that the updates are atomically applied to state.
	// Store the modified market price.
	b := k.cdc.MustMarshal(&marketPrice)
	marketPriceStore.Set(lib.Uint32ToKey(marketPrice.Id), b)

	// Monitor the last block a market price is updated.
	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, metrics.LastPriceUpdateForMarketBlock},
		float32(ctx.BlockHeight()),
		[]gometrics.Label{ // To track per market, include the id as a label.
			pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
		},
	)

	// Generate indexer events.
	priceUpdateIndexerEvents := GenerateMarketPriceUpdateIndexerEvent(marketPrice)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexerevents.MarketEventVersion,
		indexer_manager.GetBytes(
			priceUpdateIndexerEvents,
		),
	)

	return nil
}

// UpdateMarketPrices updates the prices for markets.
func (k Keeper) UpdateSpotAndPnlMarketPrices(
	ctx sdk.Context,
	update *types.MarketPriceUpdate,
) error {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.UpdateMarketPrices,
		metrics.Latency,
	)

	// Get necessary store.
	marketPriceStore := k.getMarketPriceStore(ctx)

	// Get market price.
	marketPrice, err := k.GetMarketPrice(ctx, update.MarketId)
	if err != nil {
		return err
	}

	// Report rate of price change for each market.
	if diffSpotRate, err := lib.ChangeRateUint64(marketPrice.SpotPrice, update.SpotPrice); err == nil {
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.PriceChangeRate},
			diffSpotRate,
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
			},
		)
	}

	if diffPnlRate, err := lib.ChangeRateUint64(marketPrice.PnlPrice, update.PnlPrice); err == nil {
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.PriceChangeRate},
			diffPnlRate,
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
			},
		)
	}

	// Update market price.
	marketPrice.SpotPrice = update.SpotPrice
	marketPrice.PnlPrice = update.PnlPrice

	// Report the oracle price.
	updatedSpotPrice, _ := lib.BigMulPow10(
		new(big.Int).SetUint64(update.SpotPrice),
		marketPrice.Exponent,
	).Float32()

	updatedPnlPrice, _ := lib.BigMulPow10(
		new(big.Int).SetUint64(update.PnlPrice),
		marketPrice.Exponent,
	).Float32()

	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, metrics.CurrentMarketPrices},
		updatedSpotPrice,
		[]gometrics.Label{ // To track per market, include the id as a label.
			pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
		},
	)

	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, metrics.CurrentMarketPrices},
		updatedPnlPrice,
		[]gometrics.Label{ // To track per market, include the id as a label.
			pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
		},
	)

	// Writes to the store are delayed so that the updates are atomically applied to state.
	// Store the modified market price.
	b := k.cdc.MustMarshal(&marketPrice)
	marketPriceStore.Set(lib.Uint32ToKey(marketPrice.Id), b)

	// Monitor the last block a market price is updated.
	telemetry.SetGaugeWithLabels(
		[]string{types.ModuleName, metrics.LastPriceUpdateForMarketBlock},
		float32(ctx.BlockHeight()),
		[]gometrics.Label{ // To track per market, include the id as a label.
			pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
		},
	)

	// Generate indexer events.
	priceUpdateIndexerEvents := GenerateMarketPriceUpdateIndexerEvent(marketPrice)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexerevents.MarketEventVersion,
		indexer_manager.GetBytes(
			priceUpdateIndexerEvents,
		),
	)

	return nil
}

// GetMarketPrice returns a market price from its id.
func (k Keeper) GetMarketPrice(
	ctx sdk.Context,
	id uint32,
) (types.MarketPrice, error) {
	store := k.getMarketPriceStore(ctx)
	b := store.Get(lib.Uint32ToKey(id))
	if b == nil {
		return types.MarketPrice{}, errorsmod.Wrap(types.ErrMarketPriceDoesNotExist, lib.UintToString(id))
	}

	var marketPrice = types.MarketPrice{}
	k.cdc.MustUnmarshal(b, &marketPrice)
	return marketPrice, nil
}

// GetAllMarketPrices returns all market prices.
func (k Keeper) GetAllMarketPrices(ctx sdk.Context) []types.MarketPrice {
	marketPriceStore := k.getMarketPriceStore(ctx)

	marketPrices := make([]types.MarketPrice, 0)

	iterator := marketPriceStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		marketPrice := types.MarketPrice{}
		k.cdc.MustUnmarshal(iterator.Value(), &marketPrice)
		marketPrices = append(marketPrices, marketPrice)
	}

	// Sort the market prices to return them in ascending order based on Id.
	sort.Slice(marketPrices, func(i, j int) bool {
		return marketPrices[i].Id < marketPrices[j].Id
	})

	return marketPrices
}

// GetMarketIdToValidDaemonPrice returns a map of market id to valid daemon price.
// An daemon price is valid iff:
// 1) the last update time is within a predefined threshold away from the given
// read time.
// 2) the number of prices that meet 1) are greater than the minimum number of
// exchanges specified in the given input.
// If a market does not have a valid daemon price, its `marketId` is not included
// in returned map.
func (k Keeper) GetMarketIdToValidDaemonPrice(
	ctx sdk.Context,
) map[uint32]types.MarketSpotPrice {
	allMarketParams := k.GetAllMarketParams(ctx)
	marketIdToValidDaemonPrice := k.daemonPriceCache.GetValidMedianPrices(
		allMarketParams,
		k.timeProvider.Now(),
	)

	ret := make(map[uint32]types.MarketSpotPrice)
	for _, marketParam := range allMarketParams {
		if daemonPrice, exists := marketIdToValidDaemonPrice[marketParam.Id]; exists {
			ret[marketParam.Id] = types.MarketSpotPrice{
				Id:        marketParam.Id,
				SpotPrice: daemonPrice,
				Exponent:  marketParam.Exponent,
			}
		}
	}
	return ret
}
