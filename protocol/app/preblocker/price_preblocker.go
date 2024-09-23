package preblocker

import (
	"fmt"

	"cosmossdk.io/log"
	priceapplier "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/applier"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PreBlockHandler is responsible for aggregating daemon data from each
// validator and writing the prices data into the store before any transactions
// are executed/finalized for a given block.
type PreBlockHandler struct { //golint:ignore
	logger log.Logger

	// price applier writes the aggregated prices to state.
	priceApplier *priceapplier.PriceApplier
}

func NewDaemonPreBlockHandler(
	logger log.Logger,
	priceApplier *priceapplier.PriceApplier,
) *PreBlockHandler {
	return &PreBlockHandler{
		logger:       logger,
		priceApplier: priceApplier,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating price daemon data from each validator
// and writing to the prices module store.
func (pbh *PreBlockHandler) PreBlocker(ctx sdk.Context, request *abci.RequestFinalizeBlock) (resp *sdk.ResponsePreBlock, err error) {
	if request == nil {
		return &sdk.ResponsePreBlock{}, fmt.Errorf(
			"received nil RequestFinalizeBlock in prices preblocker: height %d",
			ctx.BlockHeight(),
		)
	}

	if !veutils.AreVEEnabled(ctx) {
		pbh.logger.Info(
			"vote extensions are not enabled, skipping prices pre-blocker",
			"height", ctx.BlockHeight(),
		)
		return &sdk.ResponsePreBlock{}, nil
	}

	err = pbh.priceApplier.ApplyPricesFromVE(ctx, request, true)
	if err != nil {
		pbh.logger.Error(
			"failed to apply prices from vote extensions",
			"height", request.Height,
			"err", err,
		)

		return &sdk.ResponsePreBlock{}, err
	}

	return &sdk.ResponsePreBlock{}, nil
}
