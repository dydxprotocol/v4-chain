package process

import (
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"

	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	error_lib "github.com/StreamFinance-Protocol/stream-chain/protocol/lib/error"
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

		if veutils.AreVEEnabled(ctx) {

			if len(req.Txs) < constants.MinTxsCountWithVE {
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

			if err := DecodeAndValidateVE(
				ctx,
				req,
				extCommitBz,
				validateVoteExtensionFn,
				pricesKeeper,
				veCodec,
				extCodec,
			); err != nil {
				ctx.Logger().Error("failed to decode and validate ve", "err", err)
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}

		}

		txs, err := DecodeProcessProposalTxs(txConfig.TxDecoder(), req, pricesKeeper)
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

func DecodeAndValidateVE(
	ctx sdk.Context,
	req *abci.RequestProcessProposal,
	extCommitBz []byte,
	validateVoteExtensionFn func(ctx sdk.Context, extCommitInfo abci.ExtendedCommitInfo) error,
	pricesKeeper ProcessPricesKeeper,
	voteCodec codec.VoteExtensionCodec,
	extCodec codec.ExtendedCommitCodec,

) error {
	var extInfo abci.ExtendedCommitInfo
	extInfo, err := extCodec.Decode(extCommitBz)
	if err != nil {
		return err
	}
	if err := ve.ValidateExtendedCommitInfo(
		ctx,
		req.Height,
		extInfo,
		voteCodec,
		pricesKeeper.(ve.PreparePricesKeeper),
		validateVoteExtensionFn,
	); err != nil {
		return err
	}
	// should this happend even if it fails?
	req.Txs = req.Txs[1:]
	return nil
}
