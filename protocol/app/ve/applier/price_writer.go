package price_writer

import (
	"cosmossdk.io/log"

	"math/big"

	aggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PriceWriter is an interface that defines the methods required to aggregate and apply prices from VE's
type PriceApplier struct {
	// used to aggregate votes into final prices
	voteAggregator aggregator.VoteAggregator

	// pk is the prices keeper that is used to write prices to state.
	pricesKeeper PriceApplierPricesKeeper

	// logger
	logger log.Logger

	// codecs
	voteExtensionCodec  codec.VoteExtensionCodec
	extendedCommitCodec codec.ExtendedCommitCodec
}

func NewPriceApplier(
	logger log.Logger,
	voteAggregator aggregator.VoteAggregator,
	pricesKeeper PriceApplierPricesKeeper,
	voteExtensionCodec codec.VoteExtensionCodec,
	extendedCommitCodec codec.ExtendedCommitCodec,
) *PriceApplier {
	return &PriceApplier{
		voteAggregator:      voteAggregator,
		pricesKeeper:        pricesKeeper,
		logger:              logger,
		voteExtensionCodec:  voteExtensionCodec,
		extendedCommitCodec: extendedCommitCodec,
	}
}

func (pa *PriceApplier) ApplyPricesFromVE(
	ctx sdk.Context,
	request *abci.RequestFinalizeBlock,
) (map[string]*big.Int, error) {
	votes, err := aggregator.GetDaemonVotesFromBlock(request.Txs, pa.voteExtensionCodec, pa.extendedCommitCodec)
	if err != nil {
		pa.logger.Error(
			"failed to get extended commit info from proposal",
			"height", request.Height,
			"num_txs", len(request.Txs),
			"err", err,
		)

		return nil, err
	}

	prices, err := pa.voteAggregator.AggregateDaemonVEIntoFinalPrices(ctx, votes)
	if err != nil {
		pa.logger.Error(
			"failed to aggregate prices",
			"height", request.Height,
			"err", err,
		)

		return nil, err
	}

	marketParams := pa.pricesKeeper.GetAllMarketParams(ctx)

	pricesWrittenToState := make(map[string]*big.Int)

	for _, market := range marketParams {
		pair := market.Pair
		shouldWritePrice, price := pa.ShouldWritePriceToStore(ctx, prices, market)
		if !shouldWritePrice {
			continue
		}

		newPrice := pricestypes.MarketPriceUpdates_MarketPriceUpdate{
			MarketId: market.Id,
			Price:    price.Uint64(),
		}

		if err := pa.pricesKeeper.UpdateMarketPrice(ctx, &newPrice); err != nil {
			pa.logger.Error(
				"failed to set price for currency pair",
				"currency_pair", pair,
				"err", err,
			)

			return nil, err
		}

		pa.logger.Info(
			"set price for currency pair",
			"currency_pair", pair,
			"quote_price", newPrice.Price,
		)

		pricesWrittenToState[pair] = price
	}
	return pricesWrittenToState, nil
}

func (pa *PriceApplier) ShouldWritePriceToStore(
	ctx sdk.Context,
	prices map[string]*big.Int,
	marketToUpdate pricestypes.MarketParam,
) (bool, *big.Int) {
	marketPairToUpdate := marketToUpdate.Pair
	price, ok := prices[marketPairToUpdate]
	if !ok || price == nil {
		pa.logger.Debug(
			"no price for currency pair",
			"currency_pair", marketPairToUpdate,
		)
		return false, nil
	}

	if price.Sign() == -1 {
		pa.logger.Error(
			"price is negative",
			"currency_pair", marketPairToUpdate,
			"price", price.String(),
		)

		return false, nil
	}
	priceUpdate := pricestypes.MarketPriceUpdates{
		MarketPriceUpdates: []*pricestypes.MarketPriceUpdates_MarketPriceUpdate{
			{
				MarketId: marketToUpdate.Id,
				Price:    price.Uint64(),
			},
		},
	}

	if pa.pricesKeeper.PerformStatefulPriceUpdateValidation(ctx, &priceUpdate, false) != nil {
		pa.logger.Error(
			"price update validation failed",
			"currency_pair", marketPairToUpdate,
			"price", price.String(),
		)

		return false, nil
	}

	return true, price
}
