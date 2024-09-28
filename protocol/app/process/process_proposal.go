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

var (
	acceptResponse = &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}
	rejectResponse = &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}
)

// ProcessProposalHandler is responsible for ensuring that the list of txs in the proposed block are valid.
// Specifically, this validates:
//   - The number of txs in the proposal is at least `MinTxsCount`.
//   - VEs are valid .
//   - Tx bytes can be decoded to a valid tx.
//   - Txs are ordered correctly.
//   - Required "app-injected message" txs are included.
//   - No duplicate "app-injected message" txs are present (i.e. no "app-injected msg" in "other" txs).
//   - All messages are "valid" (i.e. `Msg.ValidateBasic` does not return errors).
//

// Note: stakingKeeper and perpetualKeeper are only needed for MEV calculations.
func ProcessProposalHandler(
	txConfig client.TxConfig,
	clobKeeper ProcessClobKeeper,
	perpetualKeeper ProcessPerpetualKeeper,
	pricesKeeper ve.PreBlockExecPricesKeeper,
	ratelimitKeeper ve.VoteExtensionRateLimitKeeper,
	extCodec codec.ExtendedCommitCodec,
	veCodec codec.VoteExtensionCodec,
	veApplier ProcessProposalVEApplier,
	validateVoteExtensionFn ve.ValidateVEConsensusInfoFn,
) sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, request *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
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
		if err := pricesKeeper.UpdateSmoothedSpotPrices(ctx, lib.Uint64LinearInterpolate); err != nil {
			recordErrorMetricsWithLabel(metrics.UpdateSmoothedPrices)
			error_lib.LogErrorWithOptionalContext(ctx, "UpdateSmoothedPrices failed", err)
		}

		if veutils.AreVEEnabled(ctx) {
			if len(request.Txs) < constants.MinTxsCountWithVE {
				ctx.Logger().Error("failed to process proposal: missing commit info", "num_txs", len(request.Txs))
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}

			extCommitBz := request.Txs[constants.DaemonInfoIndex]

			defer func() {
				// re-append the ve into the transactions to return to comet BFT
				// this gets removed to initially to not obstruct logic in validating
				// all other non-ve transactions
				// FLOW: validate ve -> remove ve -> validate all other txs -> re-append ve
				request.Txs = append([][]byte{extCommitBz}, request.Txs...)
			}()

			if err := DecodeValidateAndCacheVE(
				ctx,
				request,
				veApplier,
				extCommitBz,
				validateVoteExtensionFn,
				pricesKeeper,
				ratelimitKeeper,
				veCodec,
				extCodec,
			); err != nil {
				ctx.Logger().Error("failed to decode and validate ve", "err", err)
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}
		}

		txs, err := DecodeProcessProposalTxs(txConfig.TxDecoder(), request, pricesKeeper)
		if err != nil {
			error_lib.LogErrorWithOptionalContext(ctx, "DecodeProcessProposalTxs failed", err)
			recordErrorMetricsWithLabel(metrics.Decode)
			return rejectResponse, nil
		}

		err = txs.Validate()

		if err != nil {
			error_lib.LogErrorWithOptionalContext(ctx, "DecodeProcessProposalTxs.Validate failed", err)
			recordErrorMetricsWithLabel(metrics.Validate)
			return rejectResponse, nil
		}

		// Measure MEV metrics if enabled.
		// if clobKeeper.RecordMevMetricsIsEnabled() {
		// 	clobKeeper.RecordMevMetrics(ctx, stakingKeeper, perpetualKeeper, txs.ProposedOperationsTx.msg)
		// }

		// Record a success metric.
		recordSuccessMetrics(ctx, txs, len(request.Txs))

		return acceptResponse, nil
	}
}

func DecodeValidateAndCacheVE(
	ctx sdk.Context,
	request *abci.RequestProcessProposal,
	veApplier ProcessProposalVEApplier,
	extCommitBz []byte,
	validateVoteExtensionFn ve.ValidateVEConsensusInfoFn,
	pricesKeeper ve.PreBlockExecPricesKeeper,
	ratelimitKeeper ve.VoteExtensionRateLimitKeeper,
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
		request.Height,
		extInfo,
		voteCodec,
		pricesKeeper,
		ratelimitKeeper,
		validateVoteExtensionFn,
	); err != nil {
		return err
	}

	reqFinalizeBlock := &abci.RequestFinalizeBlock{
		Txs:    request.Txs,
		Height: request.Height,
		DecidedLastCommit: abci.CommitInfo{
			Round: ctx.CometInfo().GetLastCommit().Round(),
			Votes: []abci.VoteInfo{},
		},
	}

	if err := veApplier.ApplyVE(ctx, reqFinalizeBlock, true); err != nil {
		ctx.Logger().Error("failed to cache VE prices", "err", err)
	}
	request.Txs = request.Txs[1:]
	return nil
}
