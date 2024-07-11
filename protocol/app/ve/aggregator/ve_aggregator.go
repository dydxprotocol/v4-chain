package aggregator

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"

	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Vote encapsulates the validator and oracle data contained within a vote extension.
type Vote struct {
	// ConsAddress is the validator that submitted the vote extension.
	ConsAddress sdk.ConsAddress
	// OracleVoteExtension
	DaemonVoteExtension vetypes.DaemonVoteExtension
}

func GetDaemonVotes(
	proposal [][]byte,
	veCodec codec.VoteExtensionCodec,
	extCommitCodec codec.ExtendedCommitCodec,
) ([]Vote, error) {
	if len(proposal) < constants.InjectedNonTxCount {
		return nil, fmt.Errorf("proposal does not contain enough set messages (VE's, proposed operations, or premium votes): %d", len(proposal))
	}

	extendedCommitInfo, err := extCommitCodec.Decode(proposal[constants.DaemonInfoIndex])
	if err != nil {
		return nil, fmt.Errorf("error decoding extended-commit-info: %w", err)
	}

	votes := make([]Vote, len(extendedCommitInfo.Votes))
	for i, voteInfo := range extendedCommitInfo.Votes {
		voteExtension, err := veCodec.Decode(voteInfo.VoteExtension)
		if err != nil {
			return nil, fmt.Errorf("error decoding vote-extension: %w", err)
		}

		votes[i] = Vote{
			ConsAddress:         voteInfo.Validator.Address,
			DaemonVoteExtension: voteExtension,
		}
	}

	return votes, nil
}

// VoteAggregator is an interface that defines the methods for aggregating oracle votes into a set of prices.
// This object holds both the aggregated price resulting from a given set of votes, and the prices
// reported by each validator.
//
//go:generate mockery --name VoteAggregator --filename mock_vote_aggregator.go
type VoteAggregator interface {
	// AggregateDaemonVotes ingresses vote information which contains all
	// vote extensions each validator extended in the previous block. it is important
	// to note that
	//  1. The vote extension may be nil, in which case the validator is not providing
	//     any daemon price data for the current block. This could have occurred because the
	//     validator was offline, or its local price daemon service was down.
	//  2. The vote extension may contain prices updates for only a subset of currency pairs.
	//     This could have occurred because the price providers for the validator were
	//     offline, or the price providers did not provide a price update for a given
	//     currency pair.
	//
	// In order for a currency pair to be included in the final oracle price, the currency
	// pair must be provided by a super-majority (2/3+) of validators. This is enforced by the
	// price aggregator but can be replaced by the application.
	//
	// Notice: This method overwrites the VoteAggregator's local view of prices.
	AggregateDaemonVE(ctx sdk.Context, votes []Vote) (map[string]*big.Int, error)

	// GetPriceForValidator gets the prices reported by a given validator. This method depends
	// on the prices from the latest set of aggregated votes.
	GetPriceForValidator(validator sdk.ConsAddress) map[string]*big.Int
}

type MedianAggregator struct {
	logger log.Logger

	// used to decode prices from the vote extension
	indexPriceCache *pricefeedtypes.MarketToExchangePrices

	// keeper is used to fetch the marketParam object
	pk pk.Keeper

	// prices is a map of validator address to a map of currency pair to price
	prices map[string]map[string]*big.Int

	aggregateFn func(ctx sdk.Context, vePrices map[string]map[string]*big.Int) (map[string]*big.Int, error)
}

func NewVeAggregator(
	logger log.Logger,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	pricekeeper pk.Keeper,
	aggregateFn func(ctx sdk.Context, vePrices map[string]map[string]*big.Int) (map[string]*big.Int, error),
) VoteAggregator {
	return &MedianAggregator{
		logger:          logger,
		indexPriceCache: indexPriceCache,
		prices:          make(map[string]map[string]*big.Int),
		aggregateFn:     aggregateFn,
		pk:              pricekeeper,
	}
}
func (ma *MedianAggregator) AggregateDaemonVE(ctx sdk.Context, votes []Vote) (map[string]*big.Int, error) {

	for _, vote := range votes {
		consAddr := vote.ConsAddress.String()
		if err := ma.addVoteToAggregator(ctx, vote.ConsAddress.String(), vote.DaemonVoteExtension); err != nil {
			ma.logger.Error(
				"failed to add vote to aggregator",
				"validator_address", consAddr,
				"err", err,
			)
			return nil, err
		}
	}

	prices, err := ma.aggregateFn(ctx, ma.prices)
	if err != nil {
		ma.logger.Error(
			"failed to aggregate prices",
			"err", err,
		)
		return nil, err
	}

	ma.logger.Debug(
		"aggregated daemon price data",
		"num_prices", len(prices),
	)

	return prices, nil
}

func (ma *MedianAggregator) addVoteToAggregator(ctx sdk.Context, address string, ve vetypes.DaemonVoteExtension) error {
	if len(ve.Prices) == 0 {
		return nil
	}
	var priceupdates pricestypes.MarketPriceUpdates

	prices := make(map[string]*big.Int, len(ve.Prices))
	for marketId, priceBz := range ve.Prices {
		if len(priceBz) > constants.MaximumPriceSize {
			ma.logger.Debug(
				"failed to store price, bytes are too long",
				"market_id", marketId,
				"num_bytes", len(priceBz),
			)

			continue
		}

		market, exists := ma.pk.GetMarketParam(ctx, marketId)
		if !exists {
			ma.logger.Debug("market id not found", "market_id", marketId)
			continue
		}

		pu, err := ma.pk.GetMarketPriceUpdateFromBytes(marketId, priceBz)
		if err != nil {
			ma.logger.Debug(
				"failed to decode price",
				"marketId", marketId,
				"err", err,
			)
			continue
		}
		priceupdates.MarketPriceUpdates = append(priceupdates.MarketPriceUpdates, pu)

		prices[market.Pair] = new(big.Int).SetUint64(pu.Price)
	}

	if ma.pk.PerformStatefulPriceUpdateValidation(ctx, &priceupdates, false) != nil {
		ma.logger.Debug(
			"failed to validate price updates",
			"num_price_updates", len(priceupdates.MarketPriceUpdates),
		)

		ma.prices[address] = nil

	} else {
		ma.logger.Debug(
			"adding prices to aggregator",
			"validator_address", address,
			"num_prices", len(prices),
		)

		ma.prices[address] = prices

	}
	return nil

}

func (ma *MedianAggregator) GetPriceForValidator(validator sdk.ConsAddress) map[string]*big.Int {
	return ma.prices[validator.String()]
}
