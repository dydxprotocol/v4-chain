package process

import (
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"

	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	error_lib "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/error"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ConsensusRound = sdk.ContextKey("consensus_round")

// ProcessProposalHandler is responsible for ensuring that the list of txs in the proposed block are valid.
// Specifically, this validates:
//   - Tx bytes can be decoded to a valid tx.
//   - Txs are ordered correctly.
//   - Required "app-injected message" txs are included.
//   - No duplicate "app-injected message" txs are present (i.e. no "app-injected msg" in "other" txs).
//   - All messages are "valid" (i.e. `Msg.ValidateBasic` does not return errors).
//   - All proposed prices within `MsgUpdateMarketPrices` are valid according to non-deterministic validation.
//
// Note: `MsgUpdateMarketPrices` is an exception to only doing stateless validation. In order for this msg
// to be valid, the proposed price update values are compared against the local index price. Because the
// outcome depends on the local index price, this validation is dependent on "in-memory state"; therefore,
// this check is NOT stateless.
// Note: stakingKeeper and perpetualKeeper are only needed for MEV calculations.
func ProcessProposalHandler(
	txConfig client.TxConfig,
	clobKeeper ProcessClobKeeper,
	perpetualKeeper ProcessPerpetualKeeper,
	pricesKeeper ProcessPricesKeeper,
	extCodec codec.ExtendedCommitCodec,
	veCodec codec.VoteExtensionCodec,
	validateVoteExtensionFn func(ctx sdk.Context, extCommitInfo abci.ExtendedCommitInfo) error,
) sdk.ProcessProposalHandler {
	// Keep track of the current block height and consensus round.
	currentBlockHeight := int64(0)
	currentConsensusRound := int64(0)

	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		defer telemetry.ModuleMeasureSince(
			ModuleName,
			time.Now(),
			ModuleName, // purposely repeated to add the module name to the metric key.
			metrics.Handler,
			metrics.Latency,
		)

		// Perform the update of smoothed prices here to ensure that smoothed prices are updated even if a block is later
		// rejected by consensus. We want smoothed prices to be updated on fixed cadence, and we are piggybacking on
		// consensus round to do so.
		if err := pricesKeeper.UpdateSmoothedPrices(ctx, lib.Uint64LinearInterpolate); err != nil {
			recordErrorMetricsWithLabel(metrics.UpdateSmoothedPrices)
			error_lib.LogErrorWithOptionalContext(ctx, "UpdateSmoothedPrices failed", err)
		}

		veEnabled := ve.AreVoteExtensionsEnabled(ctx)

		if veEnabled {

			if len(req.Txs) < constants.MinTxsCount+1 {
				ctx.Logger().Error("failed to process proposal: missing commit info", "num_txs", len(req.Txs))

				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}

			extCommitBz := req.Txs[constants.DaemonInfoIndex]

			defer func() {
				// re-append the ve into the transactions to return to comet BFT
				// this gets removed to initially to not obstruct logic in validating
				// all other non-ve transactions
				// FLOW: validate ve -> remove ve -> validate all other txs -> re-append ve
				req.Txs = append([][]byte{extCommitBz}, req.Txs...)
			}()

			var extInfo abci.ExtendedCommitInfo
			extInfo, err := extCodec.Decode(extCommitBz)
			if err != nil {
				ctx.Logger().Error("failed to decode extended commit info", "err", err)

				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}
			// TODO: PreparePricesKeeper and ProcessPricesKeeper are the same (clean up)
			if err := ve.ValidateExtendedCommitInfo(ctx, req.Height, extInfo, veCodec, pricesKeeper.(ve.PreparePricesKeeper), validateVoteExtensionFn); err != nil {

				ctx.Logger().Error(
					"failed to validate extended commit info",
					"height", req.Height,
					"commit_info", extInfo,
					"err", err,
				)
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}
			req.Txs = req.Txs[1:]
		}

		// Update the current block height and consensus round.
		if ctx.BlockHeight() != currentBlockHeight {
			currentBlockHeight = ctx.BlockHeight()
			currentConsensusRound = 0
		} else {
			currentConsensusRound += 1
		}
		ctx = ctx.WithValue(ConsensusRound, currentConsensusRound)
		ctx = log.AddPersistentTagsToLogger(
			ctx,
			log.Module, ModuleName,
		)

		txs, err := DecodeProcessProposalTxs(ctx, txConfig.TxDecoder(), req, pricesKeeper)
		if err != nil {
			error_lib.LogErrorWithOptionalContext(ctx, "DecodeProcessProposalTxs failed", err)
			recordErrorMetricsWithLabel(metrics.Decode)
			return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
		}

		err = txs.Validate()
		if err != nil {
			error_lib.LogErrorWithOptionalContext(ctx, "DecodeProcessProposalTxs.Validate failed", err)
			recordErrorMetricsWithLabel(metrics.Validate)
			return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
		}

		// Measure MEV metrics if enabled.
		// if clobKeeper.RecordMevMetricsIsEnabled() {
		// 	clobKeeper.RecordMevMetrics(ctx, stakingKeeper, perpetualKeeper, txs.ProposedOperationsTx.msg)
		// }

		// Record a success metric.
		recordSuccessMetrics(ctx, txs, len(req.Txs))

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}
