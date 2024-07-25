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

// VoteAggregator is an interface that defines the methods for aggregating oracle votes into a set of prices.
// This object holds both the aggregated price resulting from a given set of votes, and the prices
// reported by each validator.
//
//go:generate mockery --name VoteAggregator --filename mock_vote_aggregator.go
type VoteAggregator interface {

	// In order for a currency pair to be included in the final oracle price, the currency
	// pair must be provided by a super-majority (2/3+) of validators. This is enforced by the
	// price aggregator but can be replaced by the application.
	AggregateDaemonVEIntoFinalPrices(ctx sdk.Context, votes []Vote) (map[string]*big.Int, error)

	// GetPriceForValidator gets the prices reported by a given validator. This method depends
	// on the prices from the latest set of aggregated votes.
	GetPriceForValidator(validator sdk.ConsAddress) map[string]*big.Int
}

type MedianAggregator struct {
	logger log.Logger

	// keeper is used to fetch the marketParam object
	pricesKeeper pk.Keeper

	// prices is a map of validator address to a map of currency pair to price
	perValidatorPrices map[string]map[string]*big.Int

	aggregateFn func(ctx sdk.Context, perValidatorPrices map[string]map[string]*big.Int) (map[string]*big.Int, error)
}

func NewVeAggregator(
	logger log.Logger,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	pricekeeper pk.Keeper,
	aggregateFn veaggregator.AggregateFn,
) VoteAggregator {
	return &MedianAggregator{
		logger:             logger,
		perValidatorPrices: make(map[string]map[string]*big.Int),
		aggregateFn:        aggregateFn,
		pricesKeeper:       pricekeeper,
	}
}

func (ma *MedianAggregator) AggregateDaemonVEIntoFinalPrices(ctx sdk.Context, votes []Vote) (map[string]*big.Int, error) {
	// wipe the previous prices
	ma.perValidatorPrices = make(map[string]map[string]*big.Int)

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

	var priceupdates pricestypes.MarketPriceUpdates

	prices := make(map[string]*big.Int, len(ve.Prices))

	for marketId, priceBytes := range ve.Prices {
		if len(priceBytes) > constants.MaximumPriceSizeInBytes {
			ma.logger.Debug(
				"failed to store price, bytes are too long",
				"market_id", marketId,
				"num_bytes", len(priceBytes),
			)
			continue
		}

		market, exists := ma.pricesKeeper.GetMarketParam(ctx, marketId)
		if !exists {
			ma.logger.Debug("market id not found", "market_id", marketId)
			continue
		}

		pu, err := veutils.GetMarketPriceUpdateFromBytes(marketId, priceBytes)
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

	ma.perValidatorPrices[address] = prices

}

func (ma *MedianAggregator) GetPriceForValidator(validator sdk.ConsAddress) map[string]*big.Int {
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
