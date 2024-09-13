package process

import (
	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FullNodeProcessProposalHandler is the `ProcessProposal` implementation for full-nodes.
// This implementation calculates and reports MEV metrics and always returns `abci.ResponseProcessProposal_ACCEPT`.
// Validators within the validator set should never use this implementation.
func FullNodeProcessProposalHandler(
	txConfig client.TxConfig,
	clobKeeper ProcessClobKeeper,
	perpetualKeeper ProcessPerpetualKeeper,
	pricesKeeper ve.PreBlockExecPricesKeeper,
) sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		// Always return `abci.ResponseProcessProposal_ACCEPT`
		response := &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}

		txs, err := DecodeProcessProposalTxs(txConfig.TxDecoder(), req, pricesKeeper)
		if err != nil {
			return response, nil
		}

		// Only validate the `ProposedOperationsTx` since full nodes don't have
		// pricefeed enabled by default and therefore, stateful validation of `UpdateMarketPricesTx`
		// would fail due to missing daemon prices.
		err = txs.ProposedOperationsTx.Validate()
		if err != nil {
			return response, nil
		}

		// // Measure MEV metrics if enabled.
		// if clobKeeper.RecordMevMetricsIsEnabled() {
		// 	clobKeeper.RecordMevMetrics(ctx, stakingKeeper, perpetualKeeper, txs.ProposedOperationsTx.msg)
		// }

		return response, nil
	}
}
