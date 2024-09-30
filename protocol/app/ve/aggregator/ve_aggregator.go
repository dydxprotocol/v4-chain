package aggregator

import (
	"fmt"
	"math/big"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/log"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
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
	AggregateDaemonVEIntoFinalPricesAndConversionRate(ctx sdk.Context, votes []Vote) (map[string]veaggregator.AggregatorPricePair, *big.Int, error)

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

	perValidatorSDaiConversionRate map[string]*big.Int

	pricesAggregateFn veaggregator.PricesAggregateFn

	conversionRateAggregateFn veaggregator.ConversionRateAggregateFn
}

func NewVeAggregator(
	logger log.Logger,
	pricekeeper pk.Keeper,
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
	}
}

func (ma *MedianAggregator) AggregateDaemonVEIntoFinalPricesAndConversionRate(
	ctx sdk.Context,
	votes []Vote,
) (map[string]veaggregator.AggregatorPricePair, *big.Int, error) {
	// wipe the previous prices
	ma.perValidatorPrices = make(map[string]map[string]veaggregator.AggregatorPricePair)

	for _, vote := range votes {
		consAddr := vote.ConsAddress.String()
		voteExtension := vote.DaemonVoteExtension
		ma.addVoteToAggregator(ctx, consAddr, voteExtension)
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

func (ma *MedianAggregator) addVoteToAggregator(
	ctx sdk.Context,
	address string,
	ve vetypes.DaemonVoteExtension,
) {
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

	var sDaiConversionRate *big.Int
	if ve.SDaiConversionRate != "" {
		newConversionRate, ok := new(big.Int).SetString(ve.SDaiConversionRate, 10)

		if ok {
			sDaiConversionRate = newConversionRate
		}
	}

	ma.perValidatorSDaiConversionRate[address] = sDaiConversionRate
}

func (ma *MedianAggregator) GetPriceForValidator(validator sdk.ConsAddress) map[string]veaggregator.AggregatorPricePair {
	return ma.perValidatorPrices[validator.String()]
}

func GetDaemonVotesFromBlock(
	proposal [][]byte,
	veCodec codec.VoteExtensionCodec,
	extCommitCodec codec.ExtendedCommitCodec,
) ([]Vote, error) {

	extCommitInfo, err := FetchExtCommitInfoFromProposal(proposal, extCommitCodec)
	if err != nil {
		return nil, fmt.Errorf("error fetching extended-commit-info: %w", err)
	}

	votes, err := FetchVotesFromExtCommitInfo(extCommitInfo, veCodec)
	if err != nil {
		return nil, fmt.Errorf("error fetching votes: %w", err)
	}

	return votes, nil
}

func FetchExtCommitInfoFromProposal(
	proposal [][]byte,
	extCommitCodec codec.ExtendedCommitCodec,
) (abci.ExtendedCommitInfo, error) {

	if len(proposal) <= constants.DaemonInfoIndex {
		return abci.ExtendedCommitInfo{}, fmt.Errorf("proposal slice is too short, expected at least %d elements but got %d", constants.DaemonInfoIndex+1, len(proposal))
	}

	extCommitInfoBytes := proposal[constants.DaemonInfoIndex]

	extCommitInfo, err := extCommitCodec.Decode(extCommitInfoBytes)
	if err != nil {
		return abci.ExtendedCommitInfo{}, fmt.Errorf("error decoding extended-commit-info: %w", err)
	}

	return extCommitInfo, nil
}

func FetchVotesFromExtCommitInfo(
	extCommitInfo abci.ExtendedCommitInfo,
	veCodec codec.VoteExtensionCodec,
) ([]Vote, error) {
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
