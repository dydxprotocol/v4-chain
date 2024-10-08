package ve_writer

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"

	aggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	bigintcache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/bigintcache"
	pricecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/pricecache"
	vecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/vecache"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VEApplier struct {
	// used to aggregate votes into final prices
	voteAggregator aggregator.VoteAggregator

	// prices keeper that is used to write prices to state.
	pricesKeeper VEApplierPricesKeeper

	// ratelimit keeper that is used to write sDAI conversion rate to state.
	ratelimitKeeper VEApplierRatelimitKeeper

	// spotPriceUpdateCache is the cache that stores the final spot prices
	spotPriceUpdateCache pricecache.PriceUpdatesCache

	// pnlPriceUpdateCache is the cache that stores the final pnl prices
	pnlPriceUpdateCache pricecache.PriceUpdatesCache

	// sDaiConversionRateCache is the cache that stores the sDAI conversion rate
	sDaiConversionRateCache bigintcache.BigIntCache

	// veCache is the cache that is used to store the seen votes
	veCache *vecache.VeCache

	// logger
	logger log.Logger

	// codecs
	voteExtensionCodec  codec.VoteExtensionCodec
	extendedCommitCodec codec.ExtendedCommitCodec
}

func NewVEApplier(
	logger log.Logger,
	voteAggregator aggregator.VoteAggregator,
	pricesKeeper VEApplierPricesKeeper,
	ratelimitKeeper VEApplierRatelimitKeeper,
	voteExtensionCodec codec.VoteExtensionCodec,
	extendedCommitCodec codec.ExtendedCommitCodec,
	spotPriceUpdateCache pricecache.PriceUpdatesCache,
	pnlPriceUpdateCache pricecache.PriceUpdatesCache,
	sDaiConversionRateCache bigintcache.BigIntCache,
	vecache *vecache.VeCache,
) *VEApplier {
	return &VEApplier{
		voteAggregator:          voteAggregator,
		pricesKeeper:            pricesKeeper,
		ratelimitKeeper:         ratelimitKeeper,
		spotPriceUpdateCache:    spotPriceUpdateCache,
		pnlPriceUpdateCache:     pnlPriceUpdateCache,
		sDaiConversionRateCache: sDaiConversionRateCache,
		logger:                  logger,
		veCache:                 vecache,
		voteExtensionCodec:      voteExtensionCodec,
		extendedCommitCodec:     extendedCommitCodec,
	}
}

func (vea *VEApplier) VoteAggregator() aggregator.VoteAggregator {
	return vea.voteAggregator
}

func (vea *VEApplier) GetVECache() *vecache.VeCache {
	return vea.veCache
}

func (vea *VEApplier) GetSpotPriceUpdateCache() pricecache.PriceUpdatesCache {
	return vea.spotPriceUpdateCache
}

func (vea *VEApplier) GetPnlPriceUpdateCache() pricecache.PriceUpdatesCache {
	return vea.pnlPriceUpdateCache
}

func (vea *VEApplier) GetSDaiConversionRateCache() bigintcache.BigIntCache {
	return vea.sDaiConversionRateCache
}
func (vea *VEApplier) CheckCacheHasValidValues(
	ctx sdk.Context,
	request *abci.RequestFinalizeBlock,
) bool {
	return vea.spotPriceUpdateCache.HasValidValues(ctx.BlockHeight(), request.DecidedLastCommit.Round) &&
		vea.pnlPriceUpdateCache.HasValidValues(ctx.BlockHeight(), request.DecidedLastCommit.Round) &&
		vea.sDaiConversionRateCache.HasValidValue(ctx.BlockHeight(), request.DecidedLastCommit.Round)
}

func (vea *VEApplier) ApplyVE(
	ctx sdk.Context,
	request *abci.RequestFinalizeBlock,
	writeToCache bool,
) error {
	if vea.CheckCacheHasValidValues(ctx, request) {
		err := vea.writeSpotPricesToStoreFromCache(ctx)
		if err != nil {
			return err
		}
		err = vea.writePnlPricesToStoreFromCache(ctx)
		if err != nil {
			return err
		}
		err = vea.writeConversionRateToStoreFromCache(ctx)
		if err != nil {
			return err
		}
		return nil
	} else {
		prices, conversionRate, err := vea.getAggregatePricesAndConversionRateFromVE(ctx, request)
		if err != nil {
			return err
		}

		err = vea.WritePricesToStoreAndMaybeCache(ctx, prices, request.DecidedLastCommit.Round, writeToCache)
		if err != nil {
			return err
		}

		err = vea.WriteSDaiConversionRateToStoreAndMaybeCache(ctx, conversionRate, request.DecidedLastCommit.Round, writeToCache)
		if err != nil {
			return err
		}
		return nil
	}
}

func (vea *VEApplier) getAggregatePricesAndConversionRateFromVE(
	ctx sdk.Context,
	request *abci.RequestFinalizeBlock,
) (map[string]voteweighted.AggregatorPricePair, *big.Int, error) {
	votes, err := aggregator.GetDaemonVotesFromBlock(
		request.Txs,
		vea.voteExtensionCodec,
		vea.extendedCommitCodec,
	)
	if err != nil {
		vea.logger.Error(
			"failed to get extended commit info from proposal",
			"height", request.Height,
			"num_txs", len(request.Txs),
			"err", err,
		)

		return nil, nil, err
	}
	prices, conversionRate, err := vea.voteAggregator.AggregateDaemonVEIntoFinalPricesAndConversionRate(ctx, votes)
	if err != nil {
		vea.logger.Error(
			"failed to aggregate prices",
			"height", request.Height,
			"err", err,
		)

		return nil, nil, err
	}

	return prices, conversionRate, nil
}

func (vea *VEApplier) writeSpotPricesToStoreFromCache(ctx sdk.Context) error {
	pricesFromCache := vea.spotPriceUpdateCache.GetPriceUpdates()
	for _, price := range pricesFromCache {
		if price.Price == nil {
			return fmt.Errorf("cache spot price is nil. spot price is %v", price.Price)
		}

		spotPriceUpdate := &pricestypes.MarketSpotPriceUpdate{
			MarketId:  price.MarketId,
			SpotPrice: price.Price.Uint64(),
		}

		if err := vea.pricesKeeper.UpdateSpotPrice(ctx, spotPriceUpdate); err != nil {
			vea.logger.Error(
				"failed to set spot price for currency pair",
				"market_id", price.MarketId,
				"err", err,
			)
			return err
		}

		vea.logger.Info(
			"set spot price for currency pair",
			"market_id", price.MarketId,
			"spot_price", price.Price.Uint64(),
		)
	}
	return nil
}

func (vea *VEApplier) writePnlPricesToStoreFromCache(ctx sdk.Context) error {
	pricesFromCache := vea.pnlPriceUpdateCache.GetPriceUpdates()
	for _, price := range pricesFromCache {
		if price.Price == nil {
			return fmt.Errorf("cache pnl price is nil. pnl price is %v", price.Price)
		}

		pnlPriceUpdate := &pricestypes.MarketPnlPriceUpdate{
			MarketId: price.MarketId,
			PnlPrice: price.Price.Uint64(),
		}

		if err := vea.pricesKeeper.UpdatePnlPrice(ctx, pnlPriceUpdate); err != nil {
			vea.logger.Error(
				"failed to set pnl price for currency pair",
				"market_id", price.MarketId,
				"err", err,
			)
			return err
		}

		vea.logger.Info(
			"set spot price for currency pair",
			"market_id", price.MarketId,
			"pnl_price", price.Price.Uint64(),
		)
	}
	return nil
}

func (vea *VEApplier) writeConversionRateToStoreFromCache(ctx sdk.Context) error {
	sDaiConversionRate := vea.sDaiConversionRateCache.GetValue()
	if sDaiConversionRate == nil {
		return nil
	}

	return vea.ratelimitKeeper.ProcessNewSDaiConversionRateUpdate(ctx, sDaiConversionRate, big.NewInt(ctx.BlockHeight()))
}

func (vea *VEApplier) WritePricesToStoreAndMaybeCache(
	ctx sdk.Context,
	prices map[string]voteweighted.AggregatorPricePair,
	round int32,
	writeToCache bool,
) error {
	marketParams := vea.pricesKeeper.GetAllMarketParams(ctx)
	var spotPriceUpdates pricecache.PriceUpdates
	var pnlPriceUpdates pricecache.PriceUpdates
	for _, market := range marketParams {
		pair := market.Pair
		pricePair, ok := prices[pair]
		if !ok {
			continue
		}

		shouldWriteSpotPrice, shouldWritePnlPrice := vea.shouldWritePriceToStore(ctx, pricePair, market.Id)

		if shouldWriteSpotPrice {
			spotPriceUpdate, err := vea.writeSpotPriceToStore(ctx, pricePair, market.Id)
			if err != nil {
				return err
			}
			spotPriceUpdates = append(spotPriceUpdates, spotPriceUpdate)
		}

		if shouldWritePnlPrice {
			pnlPriceUpdate, err := vea.writePnlPriceToStore(ctx, pricePair, market.Id)
			if err != nil {
				return err
			}
			pnlPriceUpdates = append(pnlPriceUpdates, pnlPriceUpdate)
		}
	}

	if writeToCache {
		vea.spotPriceUpdateCache.SetPriceUpdates(ctx, spotPriceUpdates, round)
		vea.pnlPriceUpdateCache.SetPriceUpdates(ctx, pnlPriceUpdates, round)
	}

	return nil
}

func (vea *VEApplier) WriteSDaiConversionRateToStoreAndMaybeCache(
	ctx sdk.Context,
	sDaiConversionRate *big.Int,
	round int32,
	writeToCache bool,
) error {
	if sDaiConversionRate != nil {
		tenScaledBySDaiDecimals := new(big.Int).Exp(
			big.NewInt(ratelimittypes.BASE_10),
			big.NewInt(ratelimittypes.SDAI_DECIMALS),
			nil,
		)
		if sDaiConversionRate.Cmp(tenScaledBySDaiDecimals) < 0 {
			return fmt.Errorf("invalid sDAI conversion rate: %s", sDaiConversionRate.String())
		}

		err := vea.ratelimitKeeper.ProcessNewSDaiConversionRateUpdate(ctx, sDaiConversionRate, big.NewInt(ctx.BlockHeight()))
		if err != nil {
			return err
		}
	}

	if writeToCache {
		vea.sDaiConversionRateCache.SetValue(ctx, sDaiConversionRate, round)
	}

	return nil
}

func (vea *VEApplier) writePnlPriceToStore(
	ctx sdk.Context,
	pricePair voteweighted.AggregatorPricePair,
	marketId uint32,
) (pricecache.PriceUpdate, error) {
	pnlPriceUpdate := &pricestypes.MarketPnlPriceUpdate{
		MarketId: marketId,
		PnlPrice: pricePair.PnlPrice.Uint64(),
	}

	if err := vea.pricesKeeper.UpdatePnlPrice(ctx, pnlPriceUpdate); err != nil {
		return pricecache.PriceUpdate{}, err
	}

	vea.logger.Info(
		"set price for currency pair",
		"market_id", marketId,
		"pnl_price", pricePair.PnlPrice.Uint64(),
	)

	return pricecache.PriceUpdate{
		MarketId: marketId,
		Price:    pricePair.PnlPrice,
	}, nil
}

func (vea *VEApplier) writeSpotPriceToStore(
	ctx sdk.Context,
	pricePair voteweighted.AggregatorPricePair,
	marketId uint32,
) (pricecache.PriceUpdate, error) {
	spotPriceUpdate := &pricestypes.MarketSpotPriceUpdate{
		MarketId:  marketId,
		SpotPrice: pricePair.SpotPrice.Uint64(),
	}

	if err := vea.pricesKeeper.UpdateSpotPrice(ctx, spotPriceUpdate); err != nil {
		return pricecache.PriceUpdate{}, err
	}

	vea.logger.Info(
		"set price for currency pair",
		"market_id", marketId,
		"spot_price", pricePair.SpotPrice.Uint64(),
	)

	return pricecache.PriceUpdate{
		MarketId: marketId,
		Price:    pricePair.SpotPrice,
	}, nil
}

func (vea *VEApplier) shouldWritePriceToStore(
	ctx sdk.Context,
	prices voteweighted.AggregatorPricePair,
	marketId uint32,
) (
	shouldWriteSpot bool,
	shouldWritePnl bool,
) {
	if prices.SpotPrice.Sign() == -1 || prices.PnlPrice.Sign() == -1 {
		vea.logger.Error(
			"price is negative",
			"market_id", marketId,
			"spot_price", prices.SpotPrice.String(),
			"pnl_price", prices.PnlPrice.String(),
		)

		return false, false
	}

	priceUpdate := &pricestypes.MarketPriceUpdate{
		MarketId:  marketId,
		SpotPrice: prices.SpotPrice.Uint64(),
		PnlPrice:  prices.PnlPrice.Uint64(),
	}

	return vea.pricesKeeper.PerformStatefulPriceUpdateValidation(ctx, priceUpdate)
}

func (vea *VEApplier) GetCachedSpotPrices() pricecache.PriceUpdates {
	return vea.spotPriceUpdateCache.GetPriceUpdates()
}

func (vea *VEApplier) GetCachedPnlPrices() pricecache.PriceUpdates {
	return vea.pnlPriceUpdateCache.GetPriceUpdates()
}

func (vea *VEApplier) GetCachedSDaiConversionRate() *big.Int {
	return vea.sDaiConversionRateCache.GetValue()
}

func (vea *VEApplier) CacheSeenExtendedVotes(
	ctx sdk.Context,
	req *abci.RequestCommit,
) error {
	if req.ExtendedCommitInfo == nil {
		return nil
	}

	votes, err := aggregator.FetchVotesFromExtCommitInfo(*req.ExtendedCommitInfo, vea.voteExtensionCodec)
	if err != nil {
		return err
	}

	seenValidators := make(map[string]struct{})
	for _, vote := range votes {
		seenValidators[vote.ConsAddress.String()] = struct{}{}
	}

	vea.veCache.SetSeenVotesInCache(ctx, seenValidators)

	return nil
}
