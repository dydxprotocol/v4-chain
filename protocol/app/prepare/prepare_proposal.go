package prepare

import (
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4/lib/metrics"
)

var (
	EmptyResponse = abci.ResponsePrepareProposal{Txs: [][]byte{}}
)

// PricesTxResponse represents a response for creating `UpdateMarketPrices` tx.
type PricesTxResponse struct {
	Tx         []byte
	Err        error
	NumMarkets int
}

// FundingTxResponse represents a response for creating `AddPremiumVotes` tx.
type FundingTxResponse struct {
	Tx       []byte
	Err      error
	NumVotes int
}

// OperationTxResponse represents a response for creating 'ProposedOperations' tx
type OperationsTxResponse struct {
	Tx            []byte
	Err           error
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
	pricesKeeper PreparePricesKeeper,
	perpetualKeeper PreparePerpetualsKeeper,
) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
		defer telemetry.ModuleMeasureSince(
			ModuleName,
			time.Now(),
			ModuleName, // purposely repeated to add the module name to the metric key.
			metrics.Handler,
			metrics.Latency,
		)

		txs, err := NewPrepareProposalTxs(req)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("NewPrepareProposalTxs error: %v", err))
			recordErrorMetricsWithLabel(metrics.PrepareProposalTxs)
			return EmptyResponse
		}

		// Gather "FixedSize" group messages.
		pricesTxResp := GetUpdateMarketPricesTx(ctx, txConfig, req.ProposerAddress, pricesKeeper)
		if pricesTxResp.Err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetUpdateMarketPricesTx error: %v", pricesTxResp.Err))
			recordErrorMetricsWithLabel(metrics.PricesTx)
			return EmptyResponse
		}
		err = txs.SetUpdateMarketPricesTx(pricesTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetUpdateMarketPricesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.PricesTx)
			return EmptyResponse
		}

		fundingTxResp := GetAddPremiumVotesTx(ctx, txConfig, perpetualKeeper)
		if fundingTxResp.Err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetAddPremiumVotesTx error: %v", fundingTxResp.Err))
			recordErrorMetricsWithLabel(metrics.FundingTx)
			return EmptyResponse
		}
		err = txs.SetAddPremiumVotesTx(fundingTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetAddPremiumVotesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.FundingTx)
			return EmptyResponse
		}

		// Gather "Other" group messages.
		otherBytesAllocated := txs.GetAvailableBytes() / int64(4) // ~25% of the remainder.
		// filter out txs that have disallow messages.
		txsWithoutDisallowMsgs := RemoveDisallowMsgs(ctx, txConfig.TxDecoder(), req.Txs)
		otherTxsToInclude, otherTxsRemainder := GetGroupMsgOther(txsWithoutDisallowMsgs, otherBytesAllocated)
		if len(otherTxsToInclude) > 0 {
			err := txs.AddOtherTxs(otherTxsToInclude)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("AddOtherTxs error: %v", err))
				recordErrorMetricsWithLabel(metrics.OtherTxs)
				return EmptyResponse
			}
		}

		// Gather "OperationsRelated" group messages.
		// TODO(DEC-1237): ensure ProposedOperations is within a certain size.
		operationsTxResp := GetProposedOperationsTx(ctx, txConfig, clobKeeper)
		if operationsTxResp.Err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetProposedOperationsTx error: %v", operationsTxResp.Err))
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return EmptyResponse
		}
		err = txs.SetProposedOperationsTx(operationsTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetProposedOperationsTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return EmptyResponse
		}

		// Try to pack in more "Other" txs.
		availableBytes := txs.GetAvailableBytes()
		if availableBytes > 0 && len(otherTxsRemainder) > 0 {
			moreOtherTxsToInclude, _ := GetGroupMsgOther(otherTxsRemainder, availableBytes)
			if len(moreOtherTxsToInclude) > 0 {
				err := txs.AddOtherTxs(moreOtherTxsToInclude)
				if err != nil {
					ctx.Logger().Error(fmt.Sprintf("AddOtherTxs (additional) error: %v", err))
					recordErrorMetricsWithLabel(metrics.OtherTxs)
					return EmptyResponse
				}
			}
		}

		txsToReturn, err := txs.GetTxsInOrder()
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetTxsInOrder error: %v", err))
			recordErrorMetricsWithLabel(metrics.GetTxsInOrder)
			return EmptyResponse
		}

		// Record a success metric.
		recordSuccessMetrics(
			successMetricParams{
				txs:                 txs,
				pricesTx:            pricesTxResp,
				fundingTx:           fundingTxResp,
				operationsTx:        operationsTxResp,
				numTxsToReturn:      len(txsToReturn),
				numTxsInOriginalReq: len(req.Txs),
			},
		)

		return abci.ResponsePrepareProposal{Txs: txsToReturn}
	}
}

// GetUpdateMarketPricesTx returns a tx containing `MsgUpdateMarketPrices`.
// The response contains an error if encoding fails.
func GetUpdateMarketPricesTx(
	ctx sdk.Context,
	txConfig client.TxConfig,
	proposerAddress []byte,
	pricesKeeper PreparePricesKeeper,
) PricesTxResponse {
	// Get prices to update.
	msgUpdateMarketPrices := pricesKeeper.GetValidMarketPriceUpdates(ctx)
	if msgUpdateMarketPrices == nil {
		return PricesTxResponse{Err: fmt.Errorf("MsgUpdateMarketPrices cannot be nil")}
	}

	tx, err := EncodeMsgsIntoTxBytes(txConfig, msgUpdateMarketPrices)
	if err != nil {
		return PricesTxResponse{Err: err}
	}
	if len(tx) == 0 {
		return PricesTxResponse{Err: fmt.Errorf("Invalid tx: %v", tx)}
	}

	return PricesTxResponse{
		Tx:         tx,
		NumMarkets: len(msgUpdateMarketPrices.MarketPriceUpdates),
	}
}

// GetAddPremiumVotesTx returns a tx containing `MsgAddPremiumVotes`.
// The response contains an error if encoding fails.
func GetAddPremiumVotesTx(
	ctx sdk.Context,
	txConfig client.TxConfig,
	perpetualsKeeper PreparePerpetualsKeeper,
) FundingTxResponse {
	// Get premium votes.
	msgAddPremiumVotes := perpetualsKeeper.GetAddPremiumVotes(ctx)
	if msgAddPremiumVotes == nil {
		return FundingTxResponse{Err: fmt.Errorf("MsgAddPremiumVotes cannot be nil")}
	}

	tx, err := EncodeMsgsIntoTxBytes(txConfig, msgAddPremiumVotes)
	if err != nil {
		return FundingTxResponse{Err: err}
	}
	if len(tx) == 0 {
		return FundingTxResponse{Err: fmt.Errorf("Invalid tx: %v", tx)}
	}

	return FundingTxResponse{
		Tx:       tx,
		NumVotes: len(msgAddPremiumVotes.Votes),
	}
}

// GetProposedOperationsTx returns a tx containing `MsgProposedOperations`.
// The response contains an error if encoding fails.
func GetProposedOperationsTx(
	ctx sdk.Context,
	txConfig client.TxConfig,
	clobKeeper PrepareClobKeeper,
) OperationsTxResponse {
	// Get the order and fill messages from the CLOB keeper.
	msgOperations := clobKeeper.GetOperations(ctx)
	if msgOperations == nil {
		return OperationsTxResponse{Err: fmt.Errorf("MsgProposedOperations cannot be nil")}
	}

	tx, err := EncodeMsgsIntoTxBytes(txConfig, msgOperations)
	if err != nil {
		return OperationsTxResponse{Err: err}
	}
	if len(tx) == 0 {
		return OperationsTxResponse{Err: fmt.Errorf("Invalid tx: %v", tx)}
	}

	return OperationsTxResponse{Tx: tx, NumOperations: len(msgOperations.GetOperationsQueue())}
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
