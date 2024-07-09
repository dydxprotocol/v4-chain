package preblocker

import (
	"fmt"

	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
)

// PreBlockHandler is responsible for aggregating daemon data from each
// validator and writing the prices data into the store before any transactions
// are executed/finalized for a given block.
type PreBlockHandler struct { //golint:ignore
	logger log.Logger

	// keeper is the keeper for the prices module. This is utilized to write
	// daemon price data to state.
	keeper pk.Keeper

	// price applier writes the aggregated prices to state.
	priceApplier veaggregator.PriceApplier
}

// NewOraclePreBlockHandler returns a new PreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewDaemonPreBlockHandler(
	logger log.Logger,
	indexPriceCache *pricefeedtypes.MarketToExchangePrices,
	pk pk.Keeper,
	priceApplier veaggregator.PriceApplier,
) *PreBlockHandler {

	return &PreBlockHandler{
		logger:       logger,
		keeper:       pk,
		priceApplier: priceApplier,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating price daemon data from each validator
// and writing to the prices module store.

func (pbh *PreBlockHandler) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) (resp *sdk.ResponsePreBlock, err error) {

	if req == nil {
		ctx.Logger().Error(
			"received nil RequestFinalizeBlock in prices PreBlocker",
			"height", ctx.BlockHeight(),
		)

		return &sdk.ResponsePreBlock{}, fmt.Errorf("received nil RequestFinalizeBlock in prices preblocker: height %d", ctx.BlockHeight())
	}

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

	_, err = pbh.priceApplier.ApplyPricesFromVoteExtensions(ctx, req)
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
