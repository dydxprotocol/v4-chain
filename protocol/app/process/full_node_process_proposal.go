package process

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FullNodeProcessProposalHandler is the `ProcessProposal` implementation for full-nodes.
// This implementation calculates and reports MEV metrics and always returns `abci.ResponseProcessProposal_ACCEPT`.
// Validators within the validator set should never use this implementation.
func FullNodeProcessProposalHandler(
	txConfig client.TxConfig,
	bridgeKeeepr ProcessBridgeKeeper,
	clobKeeper ProcessClobKeeper,
	stakingKeeper ProcessStakingKeeper,
	perpetualKeeper ProcessPerpetualKeeper,
	pricesTxDecoder UpdateMarketPriceTxDecoder,
) sdk.ProcessProposalHandler {
	// Keep track of the current block height and consensus round.
	currentBlockHeight := int64(0)
	currentConsensusRound := int64(0)

	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		// Always return `abci.ResponseProcessProposal_ACCEPT`
		response := &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}

		// Update the current block height and consensus round.
		if ctx.BlockHeight() != currentBlockHeight {
			currentBlockHeight = ctx.BlockHeight()
			currentConsensusRound = 0
		} else {
			currentConsensusRound += 1
		}
		ctx = ctx.WithValue(ConsensusRound, currentConsensusRound)

		txs, err := DecodeProcessProposalTxs(ctx, txConfig.TxDecoder(), req, bridgeKeeepr, pricesTxDecoder)
		if err != nil {
			return response, nil
		}

		// Only validate the `ProposedOperationsTx` since full nodes don't have
		// pricefeed enabled by default and therefore, stateful validation of `UpdateMarketPricesTx`
		// would fail due to missing index prices.
		err = txs.ProposedOperationsTx.Validate()
		if err != nil {
			return response, nil
		}

		// Measure MEV metrics if enabled.
		if clobKeeper.RecordMevMetricsIsEnabled() {
			clobKeeper.RecordMevMetrics(ctx, stakingKeeper, perpetualKeeper, txs.ProposedOperationsTx.msg)
		}

		return response, nil
	}
}
