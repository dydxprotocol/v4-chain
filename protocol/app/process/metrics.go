package process

import (
	gometrics "github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
)

// recordErrorMetricsWithLabel records an error metric in `ProcessProposalHandler` with a label.
func recordErrorMetricsWithLabel(label string) {
	telemetry.IncrCounterWithLabels(
		[]string{ModuleName, metrics.Handler, metrics.Error, metrics.Count},
		1,
		[]gometrics.Label{metrics.GetLabelForStringValue(metrics.Detail, label)},
	)
}

// recordSuccessMetrics records a success metric details for `ProcessProposalHandler`.
func recordSuccessMetrics(ctx sdk.Context, txs *ProcessProposalTxs, totalNumTxs int) {
	// Record success.
	telemetry.IncrCounter(
		1,
		ModuleName,
		metrics.Handler,
		metrics.Success,
		metrics.Count,
	)

	// Prices tx.
	updateMarketPricesMsg, ok := txs.UpdateMarketPricesTx.GetMsg().(*pricestypes.MsgUpdateMarketPrices)
	if ok {
		telemetry.SetGauge(
			float32(len(updateMarketPricesMsg.MarketPriceUpdates)),
			ModuleName,
			metrics.NumMarketPricesToUpdate,
		)
	} else {
		ctx.Logger().Error("ProcessProposal: expected MsgUpdateMarketPrices")
	}

	// Funding tx.
	// TODO(DEC-1254): add more metrics for Funding tx.
	addPremiumVotesMsg, ok := txs.AddPremiumVotesTx.GetMsg().(*perptypes.MsgAddPremiumVotes)
	if ok {
		telemetry.SetGauge(
			float32(len(addPremiumVotesMsg.Votes)),
			ModuleName,
			metrics.NumPremiumVotes,
		)
	} else {
		ctx.Logger().Error("ProcessProposal: expected MsgAddPremiumVotes")
	}

	// Order tx.
	// TODO(DEC-1406): add more metrics for Order tx (ie LT Place/Cancels), also update current metric.
	_, ok = txs.ProposedOperationsTx.GetMsg().(*clobtypes.MsgProposedOperations)
	if ok {
		// telemetry.SetGauge(  TODO(CLOB-279) - replace with metric for operations queue
		// 	float32(len(operationsMsg.GetOperationsQueue().MatchOrders)),
		// 	ModuleName,
		// 	metrics.Fills,
		// )
	} else {
		ctx.Logger().Error("ProcessProposal: expected MsgProposedOperations")
	}

	// Other txs.
	telemetry.SetGauge(
		float32(len(txs.OtherTxs)),
		ModuleName,
		metrics.NumOtherTxs,
	)

	// Total # of txs in the new proposal.
	telemetry.SetGauge(
		float32(totalNumTxs),
		ModuleName,
		metrics.TotalNumTxs,
	)
}
