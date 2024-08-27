package aggregator

import (
	"fmt"
	"math/big"

	"cosmossdk.io/log"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	AggregateDaemonVEIntoFinalPrices(ctx sdk.Context, votes []Vote) (map[string]veaggregator.AggregatorPricePair, error)

	// GetPriceForValidator gets the prices reported by a given validator. This method depends
	// on the prices from the latest set of aggregated votes.
	GetPriceForValidator(validator sdk.ConsAddress) map[string]veaggregator.AggregatorPricePair
}

type MedianAggregator struct {
	logger log.Logger

	// keeper is used to fetch the marketParam object
	pricesKeeper pk.Keeper

	// prices is a map of validator address to a map of currency pair to price
	perValidatorPrices map[string]map[string]veaggregator.AggregatorPricePair

	aggregateFn veaggregator.AggregateFn
}

func NewVeAggregator(
	logger log.Logger,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	pricekeeper pk.Keeper,
	aggregateFn veaggregator.AggregateFn,
) VoteAggregator {
	return &MedianAggregator{
		logger:             logger,
		perValidatorPrices: make(map[string]map[string]veaggregator.AggregatorPricePair),
		aggregateFn:        aggregateFn,
		pricesKeeper:       pricekeeper,
	}
}

func (ma *MedianAggregator) AggregateDaemonVEIntoFinalPrices(
	ctx sdk.Context,
	votes []Vote,
) (map[string]veaggregator.AggregatorPricePair, error) {
	// wipe the previous prices
	ma.perValidatorPrices = make(map[string]map[string]veaggregator.AggregatorPricePair)

	for _, vote := range votes {
		consAddr := vote.ConsAddress.String()
		voteExtension := vote.DaemonVoteExtension
		ma.addVoteToAggregator(ctx, consAddr, voteExtension)
	}

	prices, err := ma.aggregateFn(ctx, ma.perValidatorPrices)
	if err != nil {
		ma.logger.Error(
			"failed to aggregate prices",
			"err", err,
		)
		return nil, err
	}

	return prices, nil
}

func (ma *MedianAggregator) addVoteToAggregator(
	ctx sdk.Context,
	address string,
	ve vetypes.DaemonVoteExtension,
) {
	if len(ve.Prices) == 0 {
		return
	}

	prices := make(map[string]veaggregator.AggregatorPricePair, len(ve.Prices))
	for _, pricePair := range ve.Prices {
		var spotPrice, pnlPrice *big.Int

		market, exists := ma.pricesKeeper.GetMarketParam(ctx, pricePair.MarketId)
		if !exists {
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

		if spotPrice == nil {
			continue
		}

		if pnlPrice == nil {
			pnlPrice = spotPrice
		}

		prices[market.Pair] = veaggregator.AggregatorPricePair{
			SpotPrice: spotPrice,
			PnlPrice:  pnlPrice,
		}
	}
	ma.perValidatorPrices[address] = prices
}

func (ma *MedianAggregator) GetPriceForValidator(validator sdk.ConsAddress) map[string]veaggregator.AggregatorPricePair {
	return ma.perValidatorPrices[validator.String()]
}

func GetDaemonVotesFromBlock(
	proposal [][]byte,
	veCodec codec.VoteExtensionCodec,
	extCommitCodec codec.ExtendedCommitCodec,
) ([]Vote, error) {
	if len(proposal) < constants.InjectedNonTxCount {
		return nil, fmt.Errorf("proposal does not contain enough set messages (VE's, proposed operations, or premium votes): %d", len(proposal))
	}

	extCommitInfoBytes := proposal[constants.DaemonInfoIndex]

	extCommitInfo, err := extCommitCodec.Decode(extCommitInfoBytes)
	if err != nil {
		return nil, fmt.Errorf("error decoding extended-commit-info: %w", err)
	}

	votes := make([]Vote, len(extCommitInfo.Votes))
	for i, voteInfo := range extCommitInfo.Votes {
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
