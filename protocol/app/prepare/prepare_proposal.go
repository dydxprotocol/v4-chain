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
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
)

var (
	EmptyPrepareProposalResponse = abci.ResponsePrepareProposal{Txs: make([][]byte, 0)}
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

// common params between tx setters
type TxSetterUtils struct {
	Ctx      sdk.Context
	TxConfig client.TxConfig
	Txs      *PrepareProposalTxs
	Request  *abci.RequestPrepareProposal
}

// PrepareProposalHandler is responsible for preparing a block proposal that's returned to Tendermint via ABCI++.
//
// The returned txs are gathered in the following way to fit within the given request's max bytes:
//   - "Fixed" Group: Bytes=unbound. Includes price updates and premium votes and VEs.
//   - "Others" Group: Bytes=25% of max bytes minus "Fixed" Group size. Includes txs in the request.
//   - "Order" Group: Bytes=75% of max bytes minus "Fixed" Group size. Includes order matches.
//   - If there are extra available bytes and there are more txs in "Other" group, add more txs from this group.
func PrepareProposalHandler(
	txConfig client.TxConfig,
	clobKeeper PrepareClobKeeper,
	perpetualKeeper PreparePerpetualsKeeper,
	pricesKeeper ve.PreBlockExecPricesKeeper,
	ratelimitKeeper ve.VoteExtensionRateLimitKeeper,
	veCodec codec.VoteExtensionCodec,
	extCommitCodec codec.ExtendedCommitCodec,
) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, request *abci.RequestPrepareProposal) (resp *abci.ResponsePrepareProposal, err error) {
		var finalTxs [][]byte

		defer telemetry.ModuleMeasureSince(
			ModuleName,
			time.Now(),
			ModuleName, // purposely repeated to add the module name to the metric key.
			metrics.Handler,
			metrics.Latency,
		)

		if request == nil {
			ctx.Logger().Error("PrepareProposalHandler received a nil request")
			return &EmptyPrepareProposalResponse, nil
		}

		txs, err := NewPrepareProposalTxs(request)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("NewPrepareProposalTxs error: %v", err))
			recordErrorMetricsWithLabel(metrics.PrepareProposalTxs)
			return &EmptyPrepareProposalResponse, nil
		}

		txSetterUtils := TxSetterUtils{
			Ctx:      ctx,
			TxConfig: txConfig,
			Txs:      &txs,
			Request:  request,
		}

		//------------------------ VOTE EXTENSIONS ------------------------
		if err := SetVE(
			txSetterUtils,
			pricesKeeper,
			ratelimitKeeper,
			veCodec,
			extCommitCodec,
		); err != nil {
			ctx.Logger().Error(
				"failed to inject vote extensions into block",
				"height", request.Height,
				"err", err,
			)
			return &EmptyPrepareProposalResponse, nil
		}

		//------------------------ PREMIUM VOTES ------------------------
		fundingTxResp, err := SetPremiumVotesTx(
			txSetterUtils,
			perpetualKeeper,
		)

		if err != nil {
			ctx.Logger().Error(
				"failed to inject premium votes into block",
				"height", request.Height,
				"err", err,
			)
			recordErrorMetricsWithLabel(metrics.FundingTx)
			return &EmptyPrepareProposalResponse, nil
		}

		//------------------------ OTHER TXS ------------------------
		otherTxsRemainder, err := SetOneFourthOtherTxsAndGetRemainder(
			txSetterUtils,
		)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("AddOtherTxs error: %v", err))
			recordErrorMetricsWithLabel(metrics.OtherTxs)
			return &EmptyPrepareProposalResponse, nil
		}

		//------------------------ PROPOSED OPERATIONS ------------------------
		operationsTxResp, err := SetProposedOperationsTx(
			txSetterUtils,
			clobKeeper,
		)

		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetProposedOperationsTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return &EmptyPrepareProposalResponse, nil
		}

		//------------------------ REMAINDER TXS ------------------------
		if err := FillRemainderWithOtherTxs(
			txSetterUtils,
			otherTxsRemainder,
		); err != nil {
			ctx.Logger().Error(fmt.Sprintf("AddOtherTxs (additional) error: %v", err))
			recordErrorMetricsWithLabel(metrics.OtherTxs)
			return &EmptyPrepareProposalResponse, nil
		}

		finalTxs, err = GetFinalTxs(ctx, txs)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetTxsInOrder error: %v", err))
			recordErrorMetricsWithLabel(metrics.GetTxsInOrder)
			return &EmptyPrepareProposalResponse, nil
		}

		// Record a success metric.
		recordSuccessMetrics(
			successMetricParams{
				txs:                 txs,
				fundingTx:           fundingTxResp,
				operationsTx:        operationsTxResp,
				numTxsToReturn:      len(finalTxs),
				numTxsInOriginalReq: len(request.Txs),
			},
		)

		return &abci.ResponsePrepareProposal{Txs: finalTxs}, nil
	}
}

func SetVE(
	txSetterUtils TxSetterUtils,
	pricesKeeper ve.PreBlockExecPricesKeeper,
	ratelimitKeeper ve.VoteExtensionRateLimitKeeper,
	voteCodec codec.VoteExtensionCodec,
	extCodec codec.ExtendedCommitCodec,
) error {
	if !veutils.AreVEEnabled(txSetterUtils.Ctx) {
		return nil
	}

	txSetterUtils.Ctx.Logger().Info(
		"Providing oracle data using vote extensions",
		"height", txSetterUtils.Request.Height,
	)

	cleanExtCommitInfo, err := ve.CleanAndValidateExtCommitInfo(
		txSetterUtils.Ctx,
		txSetterUtils.Request.LocalLastCommit,
		voteCodec,
		pricesKeeper,
		ratelimitKeeper,
	)

	if err != nil {
		return err
	}
	// Create the vote extension injection data which will be injected into the proposal. These contain the
	// oracle data for the current block along with the sDAI conversion rate at the appropriate heights
	// which will be committed to state in PreBlock.
	extInfoBz, err := extCodec.Encode(cleanExtCommitInfo)
	if err != nil {
		return err
	}

	err = txSetterUtils.Txs.SetExtInfoBz(extInfoBz)
	if err != nil {
		return err
	}

	return nil
}

func FillRemainderWithOtherTxs(
	txSetterUtils TxSetterUtils,
	otherTxsRemainder [][]byte,
) error {
	// Try to pack in more "Other" txs.
	availableBytes := txSetterUtils.Txs.GetAvailableBytes()
	if availableBytes > 0 && len(otherTxsRemainder) > 0 {
		moreOtherTxsToInclude, _ := GetGroupMsgOther(otherTxsRemainder, availableBytes)
		if len(moreOtherTxsToInclude) > 0 {
			if err := txSetterUtils.Txs.AddOtherTxs(moreOtherTxsToInclude); err != nil {
				return err
			}
		}
	}
	return nil
}

func SetPremiumVotesTx(
	txSetterUtils TxSetterUtils,
	perpetualKeeper PreparePerpetualsKeeper,
) (FundingTxResponse, error) {
	fundingTxResp, err := GetAddPremiumVotesTx(
		txSetterUtils,
		perpetualKeeper,
	)
	if err != nil {
		return fundingTxResp, err
	}

	if err := txSetterUtils.Txs.SetAddPremiumVotesTx(fundingTxResp.Tx); err != nil {
		return fundingTxResp, err
	}

	return fundingTxResp, nil
}

func SetProposedOperationsTx(
	txSetterUtils TxSetterUtils,
	clobKeeper PrepareClobKeeper,
) (OperationsTxResponse, error) {
	// Gather "OperationsRelated" group messages.
	// TODO(DEC-1237): ensure ProposedOperations is within a certain size.
	operationsTxResp, err := GetProposedOperationsTx(
		txSetterUtils,
		clobKeeper,
	)
	if err != nil {
		return operationsTxResp, err
	}
	if err := txSetterUtils.Txs.SetProposedOperationsTx(operationsTxResp.Tx); err != nil {
		return operationsTxResp, err
	}

	return operationsTxResp, nil
}

func SetOneFourthOtherTxsAndGetRemainder(
	txSetterUtils TxSetterUtils,
) ([][]byte, error) {
	// Gather "Other" group messages.
	otherBytesAllocated := txSetterUtils.Txs.GetAvailableBytes() / 4 // ~25% of the remainder.
	// filter out txs that have disallow messages.
	txsWithoutDisallowMsgs := RemoveDisallowMsgs(
		txSetterUtils.Ctx,
		txSetterUtils.TxConfig.TxDecoder(),
		txSetterUtils.Request.Txs,
	)
	otherTxsToInclude, otherTxsRemainder := GetGroupMsgOther(txsWithoutDisallowMsgs, otherBytesAllocated)
	if len(otherTxsToInclude) > 0 {
		err := txSetterUtils.Txs.AddOtherTxs(otherTxsToInclude)
		if err != nil {
			return nil, err
		}
	}
	return otherTxsRemainder, nil
}

// GetAddPremiumVotesTx returns a tx containing `MsgAddPremiumVotes`.
func GetAddPremiumVotesTx(
	txSetterUtils TxSetterUtils,
	perpetualsKeeper PreparePerpetualsKeeper,
) (FundingTxResponse, error) {
	// Get premium votes.
	msgAddPremiumVotes := perpetualsKeeper.GetAddPremiumVotes(txSetterUtils.Ctx)
	if msgAddPremiumVotes == nil {
		return FundingTxResponse{}, fmt.Errorf("MsgAddPremiumVotes cannot be nil")
	}

	tx, err := EncodeMsgsIntoTxBytes(txSetterUtils.TxConfig, msgAddPremiumVotes)
	if err != nil {
		return FundingTxResponse{}, err
	}
	if len(tx) == 0 {
		return FundingTxResponse{}, fmt.Errorf("invalid tx: %v", tx)
	}

	return FundingTxResponse{
		Tx:       tx,
		NumVotes: len(msgAddPremiumVotes.Votes),
	}, nil
}

// GetProposedOperationsTx returns a tx containing `MsgProposedOperations`.
func GetProposedOperationsTx(
	txSetterUtils TxSetterUtils,
	clobKeeper PrepareClobKeeper,
) (OperationsTxResponse, error) {
	// Get the order and fill messages from the CLOB keeper.
	msgOperations := clobKeeper.GetOperations(txSetterUtils.Ctx)
	if msgOperations == nil {
		return OperationsTxResponse{}, fmt.Errorf("MsgProposedOperations cannot be nil")
	}

	tx, err := EncodeMsgsIntoTxBytes(txSetterUtils.TxConfig, msgOperations)
	if err != nil {
		return OperationsTxResponse{}, err
	}
	if len(tx) == 0 {
		return OperationsTxResponse{}, fmt.Errorf("invalid tx: %v", tx)
	}

	return OperationsTxResponse{
		Tx:            tx,
		NumOperations: len(msgOperations.GetOperationsQueue()),
	}, nil
}

func GetFinalTxs(ctx sdk.Context, txs PrepareProposalTxs) ([][]byte, error) {
	if veutils.AreVEEnabled(ctx) {
		return txs.GetTxsInOrder(true)
	}
	return txs.GetTxsInOrder(false)
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
