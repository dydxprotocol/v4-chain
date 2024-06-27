package deamon

import (
	"fmt"
	"math/big"

	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
)

// PreBlockHandler is responsible for aggregating deamon data from each
// validator and writing the prices data into the store before any transactions
// are executed/finalized for a given block.
type PreBlockHandler struct { //golint:ignore
	logger log.Logger

	// keeper is the keeper for the prices module. This is utilized to write
	// deamon price data to state.
	keeper pk.Keeper

	// price applier writes the aggregated prices to state.
	pa veaggregator.PriceApplier
}

// NewOraclePreBlockHandler returns a new PreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewOraclePreBlockHandler(
	logger log.Logger,
	aggregateFn func(ctx sdk.Context) (map[string]*big.Int, error),
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	pk pk.Keeper,
	veCodec codec.VoteExtensionCodec,
	ecCodec codec.ExtendedCommitCodec,
) *PreBlockHandler {

	aggregator := veaggregator.NewMedianAggregator(
		logger,
		indexPriceCache,
		aggregateFn,
	)

	priceApplier := veaggregator.NewPriceWriter(
		aggregator,
		pk,
		veCodec,
		ecCodec,
		logger,
	)

	return &PreBlockHandler{
		logger: logger,
		keeper: pk,
		pa:     priceApplier,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating price deamon data from each validator
// and writing to the prices module store.

func (pbh *PreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(ctx sdk.Context, req *abci.RequestFinalizeBlock) (_ *sdk.ResponsePreBlock, err error) {
		if req == nil {
			ctx.Logger().Error(
				"received nil RequestFinalizeBlock in prices PreBlocker",
				"height", ctx.BlockHeight(),
			)

			return &sdk.ResponsePreBlock{}, fmt.Errorf("received nil RequestFinalizeBlock in prices preblocker: height %d", ctx.BlockHeight())
		}
		defer func() {
			// only measure latency in Finalize
			if ctx.ExecMode() == sdk.ExecModeFinalize {
				pbh.logger.Debug(
					"finished executing the pre-block hook",
					"height", ctx.BlockHeight(),
				)
			}
		}()

		if !ve.AreVoteExtensionsEnabled(ctx) {
			pbh.logger.Info(
				"vote extensions are not enabled, skipping prices pre-blocker",
				"height", ctx.BlockHeight(),
			)
			return &sdk.ResponsePreBlock{}, nil
		}
		pbh.logger.Debug(
			"executing the prices pre-block hook",
			"height", req.Height,
		)

		_, err = pbh.pa.ApplyPricesFromVoteExtensions(ctx, req)
		if err != nil {
			pbh.logger.Error(
				"failed to apply prices from vote extensions",
				"height", req.Height,
				"err", err,
			)

			return &sdk.ResponsePreBlock{}, err
		}

		return &sdk.ResponsePreBlock{}, nil
	}
}
