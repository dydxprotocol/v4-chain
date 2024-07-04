package aggregator

import (
	"cosmossdk.io/log"

	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	ptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceWriter is an interface that defines the methods required to aggregate and apply prices from VE's
type PriceWriter struct {
	// va is a VoteAggregator that is used to aggregate votes into prices.
	va VoteAggregator

	// pk is the prices keeper that is used to write prices to state.
	pk pk.Keeper

	// logger
	logger log.Logger

	// codecs
	voteExtensionCodec  codec.VoteExtensionCodec
	extendedCommitCodec codec.ExtendedCommitCodec
}

type PriceApplier interface {
	// ApplyPricesFromVoteExtensions derives the aggregate prices per asset in accordance with the given
	// vote extensions + VoteAggregator. If a price exists for an asset, it is written to state. The
	// prices aggregated from vote-extensions are returned if no errors are encountered in execution,
	// otherwise an error is returned + nil prices.
	ApplyPricesFromVoteExtensions(ctx sdk.Context, req *abci.RequestFinalizeBlock) (map[string]*big.Int, error)

	// GetPriceForValidator gets the prices reported by a given validator. This method depends
	// on the prices from the latest set of aggregated votes.
	GetPricesForValidator(validator sdk.ConsAddress) map[string]*big.Int
}

func NewPriceWriter(
	va VoteAggregator,
	pk pk.Keeper,
	voteExtensionCodec codec.VoteExtensionCodec,
	extendedCommitCodec codec.ExtendedCommitCodec,
	logger log.Logger,
) PriceApplier {
	return &PriceWriter{
		va:                  va,
		pk:                  pk,
		logger:              logger,
		voteExtensionCodec:  voteExtensionCodec,
		extendedCommitCodec: extendedCommitCodec,
	}
}

func (pw *PriceWriter) ApplyPricesFromVoteExtensions(ctx sdk.Context, req *abci.RequestFinalizeBlock) (map[string]*big.Int, error) {
	votes, err := GetDaemonVotes(req.Txs, pw.voteExtensionCodec, pw.extendedCommitCodec)
	if err != nil {
		pw.logger.Error(
			"failed to get extended commit info from proposal",
			"height", req.Height,
			"num_txs", len(req.Txs),
			"err", err,
		)

		return nil, err
	}

	pw.logger.Debug(
		"got oracle vote extensions",
		"height", req.Height,
		"num_votes", len(votes),
	)

	prices, err := pw.va.AggregateDaemonVE(ctx, votes)
	if err != nil {
		pw.logger.Error(
			"failed to aggregate prices",
			"height", req.Height,
			"err", err,
		)

		return nil, err
	}

	marketParams := pw.pk.GetAllMarketParams(ctx)

	for _, market := range marketParams {
		pair := market.Pair
		price, ok := prices[pair]
		if !ok || price == nil {
			pw.logger.Debug(
				"no price for currency pair",
				"currency_pair", pair,
			)

			continue
		}

		if price.Sign() == -1 {
			pw.logger.Error(
				"price is negative",
				"currency_pair", pair,
				"price", price.String(),
			)

			continue
		}

		newPrice := ptypes.MarketPriceUpdates_MarketPriceUpdate{
			MarketId: market.Id,
			Price:    price.Uint64(),
		}

		if err := pw.pk.UpdateMarketPrice(ctx, &newPrice); err != nil {
			pw.logger.Error(
				"failed to set price for currency pair",
				"currency_pair", pair,
				"err", err,
			)

			return nil, err
		}
		pw.logger.Debug(
			"set price for currency pair",
			"currency_pair", pair,
			"quote_price", newPrice.Price,
		)
	}
	return prices, nil

}

func (pw *PriceWriter) GetPricesForValidator(validator sdk.ConsAddress) map[string]*big.Int {
	return pw.va.GetPriceForValidator(validator)
}
