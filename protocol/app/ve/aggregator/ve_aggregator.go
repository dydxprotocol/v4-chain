package aggregator

import (
	"math/big"

	"cosmossdk.io/log"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
)

// Vote encapsulates the validator and oracle data contained within a vote extension.
type Vote struct {
	// ConsAddress is the validator that submitted the vote extension.
	ConsAddress sdk.ConsAddress
	// OracleVoteExtension
	DaemonVoteExtension vetypes.DaemonVoteExtension
}

// VoteAggregator is an interface that defines the methods for aggregating oracle votes into a set of prices.
// This object holds both the aggregated price resulting from a given set of votes, and the prices
// reported by each validator.
//
//go:generate mockery --name VoteAggregator --filename mock_vote_aggregator.go
type VoteAggregator interface {

	// In order for a currency pair to be included in the final oracle price, the currency
	// pair must be provided by a super-majority (2/3+) of validators. This is enforced by the
	// price aggregator but can be replaced by the application.
	AggregateDaemonVEIntoFinalPricesAndConversionRate(ctx sdk.Context, votes []Vote) (map[string]veaggregator.AggregatorPricePair, *big.Int, error)
}

type MedianAggregator struct {
	logger log.Logger

	// keeper is used to fetch the marketParam object
	pricesKeeper pk.Keeper

	ratelimitKeeper ratelimitkeeper.Keeper

	// prices is a map of validator address to a map of currency pair to price
	perValidatorPrices map[string]map[string]veaggregator.AggregatorPricePair

	perValidatorSDaiConversionRate map[string]*big.Int

	pricesAggregateFn veaggregator.PricesAggregateFn

	conversionRateAggregateFn veaggregator.ConversionRateAggregateFn
}

func NewVeAggregator(
	logger log.Logger,
	pricekeeper pk.Keeper,
	ratelimitKeeper ratelimitkeeper.Keeper,
	pricesAggregateFn veaggregator.PricesAggregateFn,
	conversionRateAggregateFn veaggregator.ConversionRateAggregateFn,
) VoteAggregator {
	return &MedianAggregator{
		logger:                         logger,
		perValidatorPrices:             make(map[string]map[string]veaggregator.AggregatorPricePair),
		perValidatorSDaiConversionRate: make(map[string]*big.Int),
		pricesAggregateFn:              pricesAggregateFn,
		conversionRateAggregateFn:      conversionRateAggregateFn,
		pricesKeeper:                   pricekeeper,
		ratelimitKeeper:                ratelimitKeeper,
	}
}

func (ma *MedianAggregator) AggregateDaemonVEIntoFinalPricesAndConversionRate(
	ctx sdk.Context,
	votes []Vote,
) (map[string]veaggregator.AggregatorPricePair, *big.Int, error) {
	// wipe the previous prices
	ma.perValidatorPrices = make(map[string]map[string]veaggregator.AggregatorPricePair)
	lastCommittedPrices := ma.getLastCommittedPrices(ctx)
	lastCommittedSDAIPrice, found := ma.ratelimitKeeper.GetSDAIPrice(ctx)
	if !found {
		ma.logger.Error("failed to get last committed sDai conversion rate")
	}

	for _, vote := range votes {
		consAddr := vote.ConsAddress.String()
		voteExtension := vote.DaemonVoteExtension
		deepCopyPrices := ma.deepCopyLastCommittedPrices(lastCommittedPrices)
		ma.addVoteToAggregator(ctx, consAddr, voteExtension, deepCopyPrices, lastCommittedSDAIPrice)
	}

	prices, err := ma.pricesAggregateFn(ctx, ma.perValidatorPrices)
	if err != nil {
		ma.logger.Error(
			"failed to aggregate prices",
			"err", err,
		)
		return nil, nil, err
	}

	sDaiConversionRate, err := ma.conversionRateAggregateFn(ctx, ma.perValidatorSDaiConversionRate)
	if err != nil {
		ma.logger.Error(
			"failed to aggregate sDai conversion rate",
			"err", err,
		)
		return nil, nil, err
	}

	return prices, sDaiConversionRate, nil
}

func (ma *MedianAggregator) getLastCommittedPrices(ctx sdk.Context) map[string]veaggregator.AggregatorPricePair {
	lastCommittedPrices := make(map[string]veaggregator.AggregatorPricePair)

	markets := ma.pricesKeeper.GetAllMarketParams(ctx)
	marketPrices := ma.pricesKeeper.GetAllMarketPrices(ctx)

	// Create a map of market prices by ID for easier lookup
	pricesById := make(map[uint32]pricetypes.MarketPrice)
	for _, marketPrice := range marketPrices {
		pricesById[marketPrice.Id] = marketPrice
	}

	for _, market := range markets {
		marketPrice, exists := pricesById[market.Id]
		if !exists {
			ma.logger.Error("failed to find matching price for market param", "market", market.Pair, "id", market.Id)
			continue
		}

		lastCommittedPrices[market.Pair] = veaggregator.AggregatorPricePair{
			SpotPrice: new(big.Int).SetUint64(marketPrice.SpotPrice),
			PnlPrice:  new(big.Int).SetUint64(marketPrice.PnlPrice),
		}
	}

	return lastCommittedPrices
}

func (ma *MedianAggregator) deepCopyLastCommittedPrices(lastCommittedPrices map[string]veaggregator.AggregatorPricePair) map[string]veaggregator.AggregatorPricePair {
	copy := make(map[string]veaggregator.AggregatorPricePair)
	for k, v := range lastCommittedPrices {
		copy[k] = veaggregator.AggregatorPricePair{
			SpotPrice: new(big.Int).Set(v.SpotPrice),
			PnlPrice:  new(big.Int).Set(v.PnlPrice),
		}
	}
	return copy
}

func (ma *MedianAggregator) addVoteToAggregator(
	ctx sdk.Context,
	address string,
	ve vetypes.DaemonVoteExtension,
	lastCommittedPrices map[string]veaggregator.AggregatorPricePair,
	lastCommittedSDAIPrice *big.Int,
) {
	for _, pricePair := range ve.Prices {
		var spotPrice, pnlPrice *big.Int

		market, exists := ma.pricesKeeper.GetMarketParam(ctx, pricePair.MarketId)
		if !exists {
			continue
		}

		if _, ok := lastCommittedPrices[market.Pair]; !ok {
			ma.logger.Error("market pair not found in last committed prices", "pair", market.Pair)
			continue
		}

		if len(pricePair.SpotPrice) <= constants.MaximumPriceSizeInBytes {
			price, err := veutils.GetPriceFromBytes(pricePair.MarketId, pricePair.SpotPrice)
			if err == nil {
				spotPrice = price
			}
		}

		if len(pricePair.PnlPrice) <= constants.MaximumPriceSizeInBytes {
			price, err := veutils.GetPriceFromBytes(pricePair.MarketId, pricePair.PnlPrice)
			if err == nil {
				pnlPrice = price
			}
		}

		if spotPrice == nil || spotPrice.Sign() <= 0 {
			continue
		}

		if pnlPrice == nil || pnlPrice.Sign() <= 0 {
			pnlPrice = spotPrice
		}

		lastCommittedPrices[market.Pair] = veaggregator.AggregatorPricePair{
			SpotPrice: spotPrice,
			PnlPrice:  pnlPrice,
		}
	}
	ma.perValidatorPrices[address] = lastCommittedPrices

	sDaiConversionRate := lastCommittedSDAIPrice
	if ve.SDaiConversionRate != "" {
		newConversionRate, ok := new(big.Int).SetString(ve.SDaiConversionRate, 10)

		if ok {
			sDaiConversionRate = newConversionRate
		}
	}

	ma.perValidatorSDaiConversionRate[address] = sDaiConversionRate
}
