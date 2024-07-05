package prepare

import (
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
)

var (
	EmptyResponse = abci.ResponsePrepareProposal{Txs: [][]byte{}}
)

// PricesTxResponse represents a response for creating `UpdateMarketPrices` tx.
type PricesTxResponse struct {
	Tx         []byte
	NumMarkets int
}

// FundingTxResponse represents a response for creating `AddPremiumVotes` tx.
type FundingTxResponse struct {
	Tx       []byte
	NumVotes int
}

// OperationTxResponse represents a response for creating 'ProposedOperations' tx
type OperationsTxResponse struct {
	Tx            []byte
	NumOperations int
}

// PrepareProposalHandler is responsible for preparing a block proposal that's returned to Tendermint via ABCI++.
//
// The returned txs are gathered in the following way to fit within the given request's max bytes:
//   - "Fixed" Group: Bytes=unbound. Includes price updates and premium votes.
//   - "Others" Group: Bytes=25% of max bytes minus "Fixed" Group size. Includes txs in the request.
//   - "Order" Group: Bytes=75% of max bytes minus "Fixed" Group size. Includes order matches.
//   - If there are extra available bytes and there are more txs in "Other" group, add more txs from this group.
func PrepareProposalHandler(
	txConfig client.TxConfig,
	clobKeeper PrepareClobKeeper,
	perpetualKeeper PreparePerpetualsKeeper,
	pricesKeeper PreparePricesKeeper,
	veCodec codec.VoteExtensionCodec,
	extCommitCodec codec.ExtendedCommitCodec,
	validateVoteExtensionFn func(ctx sdk.Context, extCommitInfo abci.ExtendedCommitInfo) error,
) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (resp *abci.ResponsePrepareProposal, err error) {

		var (
			extInfoBz []byte
			finalTxs  [][]byte
		)

		defer telemetry.ModuleMeasureSince(
			ModuleName,
			time.Now(),
			ModuleName, // purposely repeated to add the module name to the metric key.
			metrics.Handler,
			metrics.Latency,
		)

		if req == nil {
			ctx.Logger().Error("PrepareProposalHandler received a nil request")
			return &EmptyResponse, err
		}

		txs, err := NewPrepareProposalTxs(req)

		voteExtensionsEnabled := ve.AreVoteExtensionsEnabled(ctx)
		if voteExtensionsEnabled {
			ctx.Logger().Info(
				"Providing oracle data using vote extensions",
				"height", req.Height,
			)

			// Get the vote extnesions

			extCommitInfo, err := ve.PruneAndValidateExtendedCommitInfo(
				ctx,
				req.LocalLastCommit,
				veCodec,
				pricesKeeper,
				validateVoteExtensionFn,
			)

			if err != nil {
				ctx.Logger().Error(
					"failed to prune extended commit info",
					"height", req.Height,
					"local_last_commit", req.LocalLastCommit,
					"err", err,
				)

				return &abci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, fmt.Errorf("failed to prune extended commit info: %w", err)
			}

			// Create the vote extension injection data which will be injected into the proposal. These contain the
			// oracle data for the current block which will be committed to state in PreBlock.
			extInfoBz, err = extCommitCodec.Encode(extCommitInfo)
			if err != nil {
				ctx.Logger().Error(
					"failed to extended commit info",
					"commit_info", extCommitInfo,
					"err", err,
				)

				return &abci.ResponsePrepareProposal{Txs: make([][]byte, 0)}, fmt.Errorf("failed to encode extended commit info: %w", err)
			}

			err = txs.SetExtInfoBz(extInfoBz)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("SetExtInfoBz error: %v", err))
				recordErrorMetricsWithLabel(metrics.FundingTx)
				return &EmptyResponse, nil
			}
		} else {
			// set empty VE's on first block to maintain minTxs invariant within block
			err := txs.SetExtInfoBz([]byte{})
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("SetExtInfoBz (empty) error: %v", err))
				recordErrorMetricsWithLabel(metrics.FundingTx)
				return &EmptyResponse, nil
			}
		}

		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("NewPrepareProposalTxs error: %v", err))
			recordErrorMetricsWithLabel(metrics.PrepareProposalTxs)
			return &EmptyResponse, nil
		}

		fundingTxResp, err := GetAddPremiumVotesTx(ctx, txConfig, perpetualKeeper)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetAddPremiumVotesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.FundingTx)
			return &EmptyResponse, nil
		}
		err = txs.SetAddPremiumVotesTx(fundingTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetAddPremiumVotesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.FundingTx)
			return &EmptyResponse, nil
		}

		// Gather "Other" group messages.
		otherBytesAllocated := txs.GetAvailableBytes() / 4 // ~25% of the remainder.
		// filter out txs that have disallow messages.
		txsWithoutDisallowMsgs := RemoveDisallowMsgs(ctx, txConfig.TxDecoder(), req.Txs)
		otherTxsToInclude, otherTxsRemainder := GetGroupMsgOther(txsWithoutDisallowMsgs, otherBytesAllocated)
		if len(otherTxsToInclude) > 0 {
			err := txs.AddOtherTxs(otherTxsToInclude)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("AddOtherTxs error: %v", err))
				recordErrorMetricsWithLabel(metrics.OtherTxs)
				return &EmptyResponse, nil
			}
		}

		// Gather "OperationsRelated" group messages.
		// TODO(DEC-1237): ensure ProposedOperations is within a certain size.
		operationsTxResp, err := GetProposedOperationsTx(ctx, txConfig, clobKeeper)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetProposedOperationsTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return &EmptyResponse, nil
		}
		err = txs.SetProposedOperationsTx(operationsTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetProposedOperationsTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return &EmptyResponse, nil
		}

		// Try to pack in more "Other" txs.
		availableBytes := txs.GetAvailableBytes()
		if availableBytes > 0 && len(otherTxsRemainder) > 0 {
			moreOtherTxsToInclude, _ := GetGroupMsgOther(otherTxsRemainder, availableBytes)
			if len(moreOtherTxsToInclude) > 0 {
				err = txs.AddOtherTxs(moreOtherTxsToInclude)
				if err != nil {
					ctx.Logger().Error(fmt.Sprintf("AddOtherTxs (additional) error: %v", err))
					recordErrorMetricsWithLabel(metrics.OtherTxs)
					return &EmptyResponse, nil
				}
			}
		}

		if voteExtensionsEnabled {
			finalTxs, err = txs.GetTxsInOrder(true)
		} else {
			finalTxs, err = txs.GetTxsInOrder(false)
		}

		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetTxsInOrder error: %v", err))
			recordErrorMetricsWithLabel(metrics.GetTxsInOrder)
			return &EmptyResponse, nil
		}

		// Record a success metric.
		recordSuccessMetrics(
			successMetricParams{
				txs:                 txs,
				fundingTx:           fundingTxResp,
				operationsTx:        operationsTxResp,
				numTxsToReturn:      len(finalTxs),
				numTxsInOriginalReq: len(req.Txs),
			},
		)

		return &abci.ResponsePrepareProposal{Txs: finalTxs}, nil
	}
}

// GetAddPremiumVotesTx returns a tx containing `MsgAddPremiumVotes`.
func GetAddPremiumVotesTx(
	ctx sdk.Context,
	txConfig client.TxConfig,
	perpetualsKeeper PreparePerpetualsKeeper,
) (FundingTxResponse, error) {
	// Get premium votes.
	msgAddPremiumVotes := perpetualsKeeper.GetAddPremiumVotes(ctx)
	if msgAddPremiumVotes == nil {
		return FundingTxResponse{}, fmt.Errorf("MsgAddPremiumVotes cannot be nil")
	}

	tx, err := EncodeMsgsIntoTxBytes(txConfig, msgAddPremiumVotes)
	if err != nil {
		return FundingTxResponse{}, err
	}
	if len(tx) == 0 {
		return FundingTxResponse{}, fmt.Errorf("Invalid tx: %v", tx)
	}

	return FundingTxResponse{
		Tx:       tx,
		NumVotes: len(msgAddPremiumVotes.Votes),
	}, nil
}

// GetProposedOperationsTx returns a tx containing `MsgProposedOperations`.
func GetProposedOperationsTx(
	ctx sdk.Context,
	txConfig client.TxConfig,
	clobKeeper PrepareClobKeeper,
) (OperationsTxResponse, error) {
	// Get the order and fill messages from the CLOB keeper.
	msgOperations := clobKeeper.GetOperations(ctx)
	if msgOperations == nil {
		return OperationsTxResponse{}, fmt.Errorf("MsgProposedOperations cannot be nil")
	}

	tx, err := EncodeMsgsIntoTxBytes(txConfig, msgOperations)
	if err != nil {
		return OperationsTxResponse{}, err
	}
	if len(tx) == 0 {
		return OperationsTxResponse{}, fmt.Errorf("Invalid tx: %v", tx)
	}

	return OperationsTxResponse{
		Tx:            tx,
		NumOperations: len(msgOperations.GetOperationsQueue()),
	}, nil
}

// EncodeMsgsIntoTxBytes encodes the given msgs into a single transaction.
func EncodeMsgsIntoTxBytes(txConfig client.TxConfig, msgs ...sdk.Msg) ([]byte, error) {
	txBuilder := txConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}
