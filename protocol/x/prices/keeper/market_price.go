package keeper

import (
	"math/big"
	"sort"
	"time"

	errorsmod "cosmossdk.io/errors"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	gometrics "github.com/hashicorp/go-metrics"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// getMarketPriceStore returns a prefix store for MarketPrices.
func (k Keeper) getMarketPriceStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.MarketPriceKeyPrefix))
}

// UpdateMarketPrices updates the prices for markets.
func (k Keeper) UpdateMarketPrices(
	ctx sdk.Context,
	updates []*types.MsgUpdateMarketPrices_MarketPrice,
) error {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.UpdateMarketPrices,
		metrics.Latency,
	)

	// Get necessary store.
	marketPriceStore := k.getMarketPriceStore(ctx)
	updatedMarketPrices := make([]types.MarketPrice, 0, len(updates))

	for _, update := range updates {
		// Get market price.
		marketPrice, err := k.GetMarketPrice(ctx, update.MarketId)
		if err != nil {
			return err
		}

		// Report rate of price change for each market.
		if diffRate, err := lib.ChangeRateUint64(marketPrice.Price, update.Price); err == nil {
			telemetry.SetGaugeWithLabels(
				[]string{types.ModuleName, metrics.PriceChangeRate},
				diffRate,
				[]gometrics.Label{ // To track per market, include the id as a label.
					pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
				},
			)
		}

		// Update market price.
		marketPrice.Price = update.Price
		updatedMarketPrices = append(updatedMarketPrices, marketPrice)

		// Report the oracle price.
		p10, inverse := lib.BigPow10(marketPrice.Exponent)
		updatePrice := new(big.Int).SetUint64(update.Price)
		var result float32
		if inverse {
			result, _ = new(big.Rat).SetFrac(updatePrice, p10).Float32()
		} else {
			result, _ = new(big.Rat).SetInt(updatePrice.Mul(updatePrice, p10)).Float32()
		}
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.CurrentMarketPrices},
			result,
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
			},
		)
	}

	// Writes to the store are delayed so that the updates are atomically applied to state.
	for _, marketPrice := range updatedMarketPrices {
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

		// If GRPC streaming is on, emit a price update to stream.
		if k.GetFullNodeStreamingManager().Enabled() {
			if k.GetFullNodeStreamingManager().TracksMarketId(marketPrice.Id) {
				k.GetFullNodeStreamingManager().SendPriceUpdate(
					ctx,
					types.StreamPriceUpdate{
						MarketId: marketPrice.Id,
						Price:    marketPrice,
						Snapshot: false,
					},
				)
			}
		}
	}

	// Generate indexer events.
	priceUpdateIndexerEvents := GenerateMarketPriceUpdateIndexerEvents(updatedMarketPrices)
	for _, update := range priceUpdateIndexerEvents {
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeMarket,
			indexerevents.MarketEventVersion,
			indexer_manager.GetBytes(
				update,
			),
		)
	}

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

// GetMarketIdToValidIndexPrice returns a map of market id to valid index price.
// An index price is valid iff:
// 1) the last update time is within a predefined threshold away from the given
// read time.
// 2) the number of prices that meet 1) are greater than the minimum number of
// exchanges specified in the given input.
// If a market does not have a valid index price, its `marketId` is not included
// in returned map.
func (k Keeper) GetMarketIdToValidIndexPrice(
	ctx sdk.Context,
) map[uint32]types.MarketPrice {
	allMarketParams := k.GetAllMarketParams(ctx)
	marketIdToValidIndexPrice := k.indexPriceCache.GetValidMedianPrices(
		k.Logger(ctx),
		allMarketParams,
		k.timeProvider.Now(),
	)

	ret := make(map[uint32]types.MarketPrice)
	for _, marketParam := range allMarketParams {
		if indexPrice, exists := marketIdToValidIndexPrice[marketParam.Id]; exists {
			exponent, err := k.GetExponent(ctx, marketParam.Pair)
			if err != nil {
				k.Logger(ctx).Error(
					"failed to get exponent for market",
					"market id", marketParam.Id,
					"market pair", marketParam.Pair,
					"error", err,
				)
				continue
			}
			ret[marketParam.Id] = types.MarketPrice{
				Id:       marketParam.Id,
				Price:    indexPrice,
				Exponent: exponent,
			}
		}
	}
	return ret
}

// GetStreamPriceUpdate returns a stream price update from its id.
func (k Keeper) GetStreamPriceUpdate(
	ctx sdk.Context,
	id uint32,
	snapshot bool,
) (
	val types.StreamPriceUpdate,
) {
	price, err := k.GetMarketPrice(ctx, id)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to get market price for streaming",
			"market id", id,
			"error", err,
		)
	}
	return types.StreamPriceUpdate{
		MarketId: id,
		Price:    price,
		Snapshot: snapshot,
	}
}
