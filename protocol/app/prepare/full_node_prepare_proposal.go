package prepare

import (
	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

// FullNodePrepareProposalHandler returns an EmptyResponse and logs an error
// if a node running in `--non-validating-full-node` mode attempts to run PrepareProposal.
func FullNodePrepareProposalHandler() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		ctx.Logger().Error(`
        Full nodes do not support PrepareProposal.
        This validator may be incorrectly running in full-node mode!
        Please check your configuration.
      `)
		recordErrorMetricsWithLabel(metrics.PrepareProposalTxs)

		// Return an empty response if the node is running in full-node mode so that the proposal fails.
		return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
	}
}
