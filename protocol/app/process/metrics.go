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
	operationQueue := msg.GetOperationsQueue()

	numMatchTakerOrders := 0
	numFills := 0
	numMatchPerpLiquidationsOperations := 0
	numMatchPerpDeleveragingOperations := 0
	numMatchedShortTermOrders := 0
	numMatchedStatefulOrders := 0
	numOffsettingSubaccountsForDeleveraging := 0

	for _, operation := range operationQueue {
		if matchOperation := operation.GetMatch(); matchOperation != nil {
			switch match := matchOperation.Match.(type) {
			case *clobtypes.ClobMatch_MatchOrders:
				matchOrders := match.MatchOrders
				numMatchTakerOrders += 1
				numFills += len(matchOrders.GetFills())

				takerOrderId := matchOrders.TakerOrderId
				if takerOrderId.IsStatefulOrder() {
					numMatchedStatefulOrders += 1
				} else {
					numMatchedShortTermOrders += 1
				}

				for _, fill := range matchOrders.GetFills() {
					if fill.MakerOrderId.IsStatefulOrder() {
						numMatchedStatefulOrders += 1
					} else {
						numMatchedShortTermOrders += 1
					}
				}
			case *clobtypes.ClobMatch_MatchPerpetualLiquidation:
				numMatchPerpLiquidationsOperations += 1
				perpLiquidation := matchOperation.GetMatchPerpetualLiquidation()
				fills := perpLiquidation.GetFills()
				for _, fill := range fills {
					numFills += 1
					if fill.MakerOrderId.IsStatefulOrder() {
						numMatchedStatefulOrders += 1
					} else {
						numMatchedShortTermOrders += 1
					}
				}
			case *clobtypes.ClobMatch_MatchPerpetualDeleveraging:
				numMatchPerpDeleveragingOperations += 1
				numOffsettingSubaccountsForDeleveraging += len(match.MatchPerpetualDeleveraging.GetFills())
			}
		}
	}

	// Report the number of matches.
	// This is equivalent to the number of fills.
	telemetry.SetGauge(
		float32(numFills),
		ModuleName,
		metrics.NumFills,
	)

	// Report the number of match order operations.
	// This is the number of taker orders that generated fills.
	telemetry.SetGauge(
		float32(numMatchTakerOrders),
		ModuleName,
		metrics.NumMatchTakerOrders,
	)

	// Report the number of matched stateful orders.
	telemetry.SetGauge(
		float32(numMatchedStatefulOrders),
		ModuleName,
		metrics.NumMatchStatefulOrders,
	)

	// Report the number of matched short term orders.
	telemetry.SetGauge(
		float32(numMatchedShortTermOrders),
		ModuleName,
		metrics.NumMatchedShortTermOrders,
	)

	// Report the number of match perp liquidation operations.
	telemetry.SetGauge(
		float32(numMatchPerpLiquidationsOperations),
		ModuleName,
		metrics.NumMatchPerpLiquidationsOperations,
	)

	// Report the number of match perp deleveraging operations.
	telemetry.SetGauge(
		float32(numMatchPerpDeleveragingOperations),
		ModuleName,
		metrics.NumMatchPerpDeleveragingOperations,
	)

	// Report the number of offsetting subaccounts for perp deleveraging.
	telemetry.SetGauge(
		float32(numOffsettingSubaccountsForDeleveraging),
		ModuleName,
		metrics.NumOffsettingSubaccountsForDeleveraging,
	)

	// Report length of operations queue in process proposal.
	telemetry.SetGauge(
		float32(len(operationQueue)),
		ModuleName,
		metrics.NumProposedOperations,
	)
}
