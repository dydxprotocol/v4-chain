package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"math/big"
	"sort"
	"time"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	gometrics "github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// newMarketPriceStore creates a new prefix store for MarketPrices.
func (k Keeper) newMarketPriceStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MarketPriceKeyPrefix))
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
	marketPriceStore := k.newMarketPriceStore(ctx)
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
		updatedPrice, _ := lib.BigMulPow10(
			new(big.Int).SetUint64(update.Price),
			marketPrice.Exponent,
		).Float32()
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.CurrentMarketPrices},
			updatedPrice,
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
			},
		)
	}

	// Writes to the store are delayed so that the updates are atomically applied to state.
	for _, marketPrice := range updatedMarketPrices {
		// Store the modified market price.
		b := k.cdc.MustMarshal(&marketPrice)
		marketPriceStore.Set(types.MarketKey(marketPrice.Id), b)

		// Monitor the last block a market price is updated.
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.LastPriceUpdateForMarketBlock},
			float32(ctx.BlockHeight()),
			[]gometrics.Label{ // To track per market, include the id as a label.
				pricefeedmetrics.GetLabelForMarketId(marketPrice.Id),
			},
		)
	}

	marketPriceUpdates := GenerateMarketPriceUpdateEvents(updatedMarketPrices)
	for _, update := range marketPriceUpdates {
		k.GetIndexerEventManager().AddTxnEvent(
			ctx,
			indexerevents.SubtypeMarket,
			indexer_manager.GetB64EncodedEventMessage(
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
	store := k.newMarketPriceStore(ctx)
	b := store.Get(types.MarketKey(id))
	if b == nil {
		return types.MarketPrice{}, errorsmod.Wrap(types.ErrMarketPriceDoesNotExist, lib.Uint32ToString(id))
	}

	var marketPrice = types.MarketPrice{}
	k.cdc.MustUnmarshal(b, &marketPrice)
	return marketPrice, nil
}

// GetAllMarketPrices returns all market prices.
func (k Keeper) GetAllMarketPrices(ctx sdk.Context) []types.MarketPrice {
	marketPriceStore := k.newMarketPriceStore(ctx)

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
