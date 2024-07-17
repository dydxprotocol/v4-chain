package price_writer

import (
	"cosmossdk.io/log"

	"math/big"

	aggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	ptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceWriter is an interface that defines the methods required to aggregate and apply prices from VE's
type PriceApplier struct {
	// va is a VoteAggregator that is used to aggregate votes into prices.
	va aggregator.VoteAggregator

	// pk is the prices keeper that is used to write prices to state.
	pk PriceApplierPricesKeeper

	// logger
	logger log.Logger

	// codecs
	voteExtensionCodec  codec.VoteExtensionCodec
	extendedCommitCodec codec.ExtendedCommitCodec
}

func NewPriceApplier(
	va aggregator.VoteAggregator,
	pk PriceApplierPricesKeeper,
	voteExtensionCodec codec.VoteExtensionCodec,
	extendedCommitCodec codec.ExtendedCommitCodec,
	logger log.Logger,
) *PriceApplier {
	return &PriceApplier{
		va:                  va,
		pk:                  pk,
		logger:              logger,
		voteExtensionCodec:  voteExtensionCodec,
		extendedCommitCodec: extendedCommitCodec,
	}
}

func (pw *PriceApplier) ApplyPricesFromVoteExtensions(ctx sdk.Context, req *abci.RequestFinalizeBlock) (map[string]*big.Int, error) {
	votes, err := aggregator.GetDaemonVotes(req.Txs, pw.voteExtensionCodec, pw.extendedCommitCodec)
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
		pw.logger.Info(
			"set price for currency pair",
			"currency_pair", pair,
			"quote_price", newPrice.Price,
		)
	}
	return prices, nil

}
