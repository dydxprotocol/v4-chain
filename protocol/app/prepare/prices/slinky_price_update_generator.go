package prices

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/slinky/abci/strategies/aggregator"
	"github.com/dydxprotocol/slinky/abci/strategies/codec"
	"github.com/dydxprotocol/slinky/abci/strategies/currencypair"
	"github.com/dydxprotocol/slinky/abci/ve"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// SlinkyPriceUpdateGenerator is an implementation of the PriceUpdateGenerator interface. This implementation
// retrieves the MsgUpdateMarketPricesTx by aggregating over all VoteExtensions from set of PreCommits on
// the last block (these commits are local to the proposer).
type SlinkyPriceUpdateGenerator struct {
	// VoteAggregator is responsible for reading all votes in the extended-commit and unmarshalling
	// them into a set of prices for the MsgUpdateMarketPricesTx.
	agg aggregator.VoteAggregator

	// extCommitCodec is responsible for unmarshalling the extended-commit from the proposal.
	extCommitCodec codec.ExtendedCommitCodec

	// veCodec is responsible for unmarshalling each vote-extension from the extended-commit
	veCodec codec.VoteExtensionCodec

	// currencyPairStrategy is responsible for mapping the currency-pairs to market-ids in the MsgUpdatemarketPricesTx
	currencyPairStrategy currencypair.CurrencyPairStrategy
}

// NewSlinkyPriceUpdateGenerator returns a new SlinkyPriceUpdateGenerator
func NewSlinkyPriceUpdateGenerator(
	agg aggregator.VoteAggregator,
	extCommitCodec codec.ExtendedCommitCodec,
	veCodec codec.VoteExtensionCodec,
	currencyPairStrategy currencypair.CurrencyPairStrategy,
) *SlinkyPriceUpdateGenerator {
	return &SlinkyPriceUpdateGenerator{
		agg:                  agg,
		extCommitCodec:       extCommitCodec,
		veCodec:              veCodec,
		currencyPairStrategy: currencyPairStrategy,
	}
}

func (pug *SlinkyPriceUpdateGenerator) GetValidMarketPriceUpdates(
	ctx sdk.Context, extCommitBz []byte) (*pricestypes.MsgUpdateMarketPrices, error) {
	// check whether VEs are enabled
	if !ve.VoteExtensionsEnabled(ctx) {
		// return a nil MsgUpdateMarketPricesTx w/ no updates
		return &pricestypes.MsgUpdateMarketPrices{}, nil
	}

	// unmarshal the injected extended-commit
	votes, err := aggregator.GetOracleVotes(
		[][]byte{extCommitBz},
		pug.veCodec,
		pug.extCommitCodec,
	)
	if err != nil {
		return nil, err
	}

	// aggregate the votes into a MsgUpdateMarketPricesTx
	prices, err := pug.agg.AggregateOracleVotes(ctx, votes)
	if err != nil {
		return nil, err
	}

	// create the update-market prices tx
	msg := &pricestypes.MsgUpdateMarketPrices{}

	// map the currency-pairs to market-ids
	for cp, price := range prices {
		marketID, err := pug.currencyPairStrategy.ID(ctx, cp)
		if err != nil {
			return nil, err
		}

		if !price.IsUint64() {
			return nil, &InvalidPriceError{
				MarketID: marketID,
				Reason:   "price is not a uint64",
			}
		}

		// add the price to the update-market prices tx
		msg.MarketPriceUpdates = append(msg.MarketPriceUpdates, &pricestypes.MsgUpdateMarketPrices_MarketPrice{
			MarketId: uint32(marketID),
			Price:    price.Uint64(),
		})
	}

	// sort the market-price updates
	sort.Slice(msg.MarketPriceUpdates, func(i, j int) bool {
		return msg.MarketPriceUpdates[i].MarketId < msg.MarketPriceUpdates[j].MarketId
	})

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	return msg, nil
}
