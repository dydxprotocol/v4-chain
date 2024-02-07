package process

import (
	gometrics "github.com/hashicorp/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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

	// Bridge tx.
	msgAcknowledgeBridges, ok := txs.AcknowledgeBridgesTx.GetMsg().(*bridgetypes.MsgAcknowledgeBridges)
	if ok {
		telemetry.IncrCounter(
			float32(len(msgAcknowledgeBridges.Events)),
			ModuleName,
			metrics.NumBridges,
		)
	} else {
		ctx.Logger().Error("ProcessProposal: expected MsgAcknowledgeBridges")
	}

	// Order tx.
	msgProposedOperations, ok := txs.ProposedOperationsTx.GetMsg().(*clobtypes.MsgProposedOperations)
	if ok {
		recordMsgProposedOperationsMetrics(ctx, msgProposedOperations)
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

// recordMsgProposedOperationsMetrics reports metrics on a `MsgProposedOperations`
// object. It is used in the process_proposal module.
func recordMsgProposedOperationsMetrics(ctx sdk.Context, msg *clobtypes.MsgProposedOperations) {
	operationsStats := clobtypes.StatMsgProposedOperations(msg.GetOperationsQueue())
	operationsStats.EmitStats(metrics.ProcessProposal)
}
