package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/metrics"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateMarket creates a new market param in the store along with a new market price
// for that market param.
// This is the only path to creating new MarketPrices, so if we have a param
// defined for a market, we should expect to see a price defined, and vice versa.
func (k Keeper) CreateMarket(
	ctx sdk.Context,
	marketParam types.MarketParam,
	marketPrice types.MarketPrice,
) (types.MarketParam, error) {
	if _, exists := k.GetMarketParam(ctx, marketParam.Id); exists {
		return types.MarketParam{}, errorsmod.Wrapf(
			types.ErrMarketParamAlreadyExists,
			"market param with id %d already exists",
			marketParam.Id,
		)
	}

	// Validate input.
	if err := marketParam.Validate(); err != nil {
		return types.MarketParam{}, err
	}
	if err := marketPrice.ValidateFromParam(marketParam); err != nil {
		return types.MarketParam{}, err
	}

	if marketPrice.SpotPrice != marketPrice.PnlPrice {
		return types.MarketParam{}, errorsmod.Wrapf(
			types.ErrInvalidInput,
			"spot price and pnl price must be the same",
		)
	}

	paramBytes := k.cdc.MustMarshal(&marketParam)
	priceBytes := k.cdc.MustMarshal(&marketPrice)

	marketParamStore := k.getMarketParamStore(ctx)
	marketParamStore.Set(lib.Uint32ToKey(marketParam.Id), paramBytes)

	marketPriceStore := k.getMarketPriceStore(ctx)
	marketPriceStore.Set(lib.Uint32ToKey(marketPrice.Id), priceBytes)

	// Generate indexer event.
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexerevents.MarketEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewMarketCreateEvent(
				marketParam.Id,
				marketParam.Pair,
				marketParam.MinPriceChangePpm,
				marketParam.Exponent,
			),
		),
	)

	k.marketToCreatedAt[marketParam.Id] = k.timeProvider.Now()
	metrics.SetMarketPairForTelemetry(marketParam.Id, marketParam.Pair)

	return marketParam, nil
}

// IsRecentlyAvailable returns true if the market was recently made available to the pricefeed daemon. A market is
// considered recently available either if it was recently created, or if the pricefeed daemon was recently started. If
// an daemon price does not exist for a recently available market, the protocol does not consider this an error
// condition, as it is expected that the pricefeed daemon will eventually provide a price for the market within a
// few seconds.
func (k Keeper) IsRecentlyAvailable(ctx sdk.Context, marketId uint32) bool {
	createdAt, ok := k.marketToCreatedAt[marketId]

	if !ok {
		return false
	}

	// The comparison condition considers both market age and price daemon warmup time because a market can be
	// created before or after the daemon starts. We use block height as a proxy for daemon warmup time because
	// the price daemon is started when the gRPC service comes up, which typically occurs just before the first
	// block is processed.
	return k.timeProvider.Now().Sub(createdAt) < types.MarketIsRecentDuration ||
		ctx.BlockHeight() < types.PriceDaemonInitializationBlocks
}

// GetAllMarketParamPrices returns a slice of MarketParam, MarketPrice tuples for all markets.
func (k Keeper) GetAllMarketParamPrices(ctx sdk.Context) ([]types.MarketParamPrice, error) {
	marketParams := k.GetAllMarketParams(ctx)
	marketPrices := k.GetAllMarketPrices(ctx)

	if len(marketParams) != len(marketPrices) {
		return nil, errorsmod.Wrap(types.ErrMarketPricesAndParamsDontMatch, "market param and price lengths do not match")
	}

	marketParamPrices := make([]types.MarketParamPrice, len(marketParams))
	for i, param := range marketParams {
		marketParamPrices[i].Param = param
		price := marketPrices[i]
		if param.Id != price.Id {
			return nil, errorsmod.Wrap(types.ErrMarketPricesAndParamsDontMatch,
				fmt.Sprintf("market param and price ids do not match: %d != %d", param.Id, price.Id))
		}
		marketParamPrices[i].Price = price
	}
	return marketParamPrices, nil
}
