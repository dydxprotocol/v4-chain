package prepare

import (
	"fmt"
	"strings"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/slinky/abci/ve"
	"github.com/dydxprotocol/v4-chain/protocol/app/constants"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare/prices"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricetypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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

// BridgeTxResponse represents a response for creating 'AcknowledgeBridges' tx
type BridgeTxResponse struct {
	Tx         []byte
	NumBridges int
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
	bridgeKeeper PrepareBridgeKeeper,
	clobKeeper PrepareClobKeeper,
	perpetualKeeper PreparePerpetualsKeeper,
	priceUpdateGenerator prices.PriceUpdateGenerator,
) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
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
			return &EmptyResponse, nil
		}

		// Grab the injected VEs from the previous block.
		var extCommitBzTx []byte
		// Sanity check to ensure that there is at least 1 tx. This should never return false unless
		// before VE are enabled, there are no tx in the block.
		if len(req.Txs) >= constants.OracleVEInjectedTxs {
			extCommitBzTx = req.Txs[constants.OracleInfoIndex]
		}

		// get the update market prices tx
		msg, err := priceUpdateGenerator.GetValidMarketPriceUpdates(ctx, extCommitBzTx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetValidMarketPriceUpdates error: %v", err))
			recordErrorMetricsWithLabel(metrics.PricesTx)
			return &EmptyResponse, nil
		}

		// Gather "FixedSize" group messages.
		pricesTxResp, err := EncodeMarketPriceUpdates(txConfig, msg)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetUpdateMarketPricesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.PricesTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}
		err = txs.SetUpdateMarketPricesTx(pricesTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetUpdateMarketPricesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.PricesTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}

		fundingTxResp, err := GetAddPremiumVotesTx(ctx, txConfig, perpetualKeeper)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetAddPremiumVotesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.FundingTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}
		err = txs.SetAddPremiumVotesTx(fundingTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetAddPremiumVotesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.FundingTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}

		acknowledgeBridgesTxResp, err := GetAcknowledgeBridgesTx(ctx, txConfig, bridgeKeeper)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetAcknowledgeBridgesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.AcknowledgeBridgesTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}
		// Set AcknowledgeBridgesTx whether there are bridge events or not to ensure
		// consistent ordering of txs received by ProcessProposal.
		err = txs.SetAcknowledgeBridgesTx(acknowledgeBridgesTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetAcknowledgeBridgesTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.AcknowledgeBridgesTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}

		// Gather "Other" group messages.
		otherBytesAllocated := txs.GetAvailableBytes() / 4 // ~25% of the remainder.
		// filter out txs that have disallow messages.
		var txsWithoutDisallowMsgs [][]byte
		if ve.VoteExtensionsEnabled(ctx) {
			txsWithoutDisallowMsgs = RemoveDisallowMsgs(ctx, txConfig.TxDecoder(), req.Txs[1:])
		} else {
			txsWithoutDisallowMsgs = RemoveDisallowMsgs(ctx, txConfig.TxDecoder(), req.Txs)
		}
		txsWithoutDisallowMsgs, numCancelOnlyTxs := ReorderClobCancelsFirst(ctx, txConfig.TxDecoder(), txsWithoutDisallowMsgs)
		ctx.Logger().Info("PrepareProposal: reordered txs for CLOB cancel-first",
			"height", req.Height,
			"total_txs", len(txsWithoutDisallowMsgs),
			"cancel_only_txs_first", numCancelOnlyTxs,
			"other_txs_after", len(txsWithoutDisallowMsgs)-numCancelOnlyTxs,
		)

		otherTxsToInclude, otherTxsRemainder := GetGroupMsgOther(txsWithoutDisallowMsgs, otherBytesAllocated)
		ctx.Logger().Info("PrepareProposal: Other group allocation",
			"height", req.Height,
			"other_bytes_allocated", otherBytesAllocated,
			"other_txs_included", len(otherTxsToInclude),
			"other_txs_remainder", len(otherTxsRemainder),
		)
		if len(otherTxsToInclude) > 0 {
			err := txs.AddOtherTxs(otherTxsToInclude)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("AddOtherTxs error: %v", err))
				recordErrorMetricsWithLabel(metrics.OtherTxs)
				return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
			}
		}

		// Gather "OperationsRelated" group messages.
		// TODO(DEC-1237): ensure ProposedOperations is within a certain size.
		memClobCtx, _ := ctx.CacheContext()
		memClobCtx = memClobCtx.WithIsCheckTx(true)

		var operationsTxResp OperationsTxResponse
		err = ApplyClobMsgsToMemClob(memClobCtx, txConfig.TxDecoder(), clobKeeper, txsWithoutDisallowMsgs)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("ApplyClobMsgsToMemClob error (using empty ops): %v", err))
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}

		msgOperations := clobKeeper.GetOperations(memClobCtx)
		if msgOperations == nil {
			ctx.Logger().Error("GetOperations returned nil msg")
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}
		// Log all proposed operations in detail.
		logProposedOperations(ctx, req.Height, msgOperations)
		operationsTxResp, err = EncodeProposedOperationsTx(txConfig, msgOperations)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("EncodeProposedOperationsTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}
		err = txs.SetProposedOperationsTx(operationsTxResp.Tx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("SetProposedOperationsTx error: %v", err))
			recordErrorMetricsWithLabel(metrics.OperationsTx)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
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
					return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
				}
			}
		}

		txsToReturn, err := txs.GetTxsInOrder()
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("GetTxsInOrder error: %v", err))
			recordErrorMetricsWithLabel(metrics.GetTxsInOrder)
			return &abci.ResponsePrepareProposal{Txs: [][]byte{}}, nil
		}
		ctx.Logger().Info("PrepareProposal: final txs to return",
			"height", req.Height,
			"num_txs", len(txsToReturn),
			"num_original_req_txs", len(req.Txs),
		)

		// Record a success metric.
		recordSuccessMetrics(
			successMetricParams{
				txs:                 txs,
				pricesTx:            pricesTxResp,
				fundingTx:           fundingTxResp,
				bridgeTx:            acknowledgeBridgesTxResp,
				operationsTx:        operationsTxResp,
				numTxsToReturn:      len(txsToReturn),
				numTxsInOriginalReq: len(req.Txs),
			},
		)

		return &abci.ResponsePrepareProposal{Txs: txsToReturn}, nil
	}
}

// logProposedOperations logs each operation in the proposed operations queue for the given height.
func logProposedOperations(ctx sdk.Context, height int64, msg *clobtypes.MsgProposedOperations) {
	if msg == nil {
		return
	}
	queue := msg.GetOperationsQueue()
	ctx.Logger().Info("PrepareProposal: proposed operations",
		"height", height,
		"num_operations", len(queue),
	)
	for i, op := range queue {
		ctx.Logger().Info("PrepareProposal: proposed operation",
			"height", height,
			"index", i,
			"operation", formatOperationForLog(op),
		)
	}
}

// formatOperationForLog returns a short, human-readable description of an OperationRaw for logging.
func formatOperationForLog(op clobtypes.OperationRaw) string {
	switch o := op.Operation.(type) {
	case *clobtypes.OperationRaw_Match:
		return formatMatchForLog(o.Match)
	case *clobtypes.OperationRaw_ShortTermOrderPlacement:
		return fmt.Sprintf("ShortTermOrderPlacement(len=%d)", len(o.ShortTermOrderPlacement))
	case *clobtypes.OperationRaw_OrderRemoval:
		if o.OrderRemoval != nil {
			return fmt.Sprintf("OrderRemoval orderId=%v reason=%s",
				o.OrderRemoval.GetOrderId(),
				o.OrderRemoval.GetRemovalReason().String())
		}
		return "OrderRemoval(nil)"
	default:
		return "UnknownOperation"
	}
}

// formatMatchForLog returns a short description of a ClobMatch for logging.
func formatMatchForLog(match *clobtypes.ClobMatch) string {
	if match == nil {
		return "Match(nil)"
	}
	switch m := match.Match.(type) {
	case *clobtypes.ClobMatch_MatchOrders:
		if mo := m.MatchOrders; mo != nil {
			taker := mo.GetTakerOrderId()
			fills := mo.GetFills()
			makerDesc := fmt.Sprintf("%d maker(s)", len(fills))
			if len(fills) > 0 {
				makerIds := make([]string, 0, len(fills))
				for _, f := range fills {
					makerIds = append(makerIds, fmt.Sprintf("%v", f.GetMakerOrderId()))
				}
				makerDesc = fmt.Sprintf("makers=[%s]", strings.Join(makerIds, ","))
			}
			return fmt.Sprintf("MatchOrders taker=%v %s", taker, makerDesc)
		}
		return "MatchOrders(nil)"
	case *clobtypes.ClobMatch_MatchPerpetualLiquidation:
		if liq := m.MatchPerpetualLiquidation; liq != nil {
			return fmt.Sprintf("MatchPerpetualLiquidation liquidated=%v clobPairId=%d perpetualId=%d totalSize=%d",
				liq.GetLiquidated(), liq.GetClobPairId(), liq.GetPerpetualId(), liq.GetTotalSize())
		}
		return "MatchPerpetualLiquidation(nil)"
	case *clobtypes.ClobMatch_MatchPerpetualDeleveraging:
		if d := m.MatchPerpetualDeleveraging; d != nil {
			return fmt.Sprintf("MatchPerpetualDeleveraging liquidated=%v perpetualId=%d num_fills=%d",
				d.GetLiquidated(), d.GetPerpetualId(), len(d.GetFills()))
		}
		return "MatchPerpetualDeleveraging(nil)"
	default:
		return "Match(unknown type)"
	}
}

// EncodeMarketPriceUpdates returns a tx containing `MsgUpdateMarketPrices`.
func EncodeMarketPriceUpdates(
	txConfig client.TxConfig,
	msg *pricetypes.MsgUpdateMarketPrices,
) (PricesTxResponse, error) {
	tx, err := EncodeMsgsIntoTxBytes(txConfig, msg)
	if err != nil {
		return PricesTxResponse{}, err
	}
	if len(tx) == 0 {
		return PricesTxResponse{}, fmt.Errorf("Invalid tx: %v", tx)
	}

	return PricesTxResponse{
		Tx:         tx,
		NumMarkets: len(msg.MarketPriceUpdates),
	}, nil
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
// GetAcknowledgeBridgeTx returns a tx containing a list of `MsgAcknowledgeBridge`.
func GetAcknowledgeBridgesTx(
	ctx sdk.Context,
	txConfig client.TxConfig,
	bridgeKeeper PrepareBridgeKeeper,
) (BridgeTxResponse, error) {
	msgAcknowledgeBridges := bridgeKeeper.GetAcknowledgeBridges(ctx, ctx.BlockTime())

	tx, err := EncodeMsgsIntoTxBytes(txConfig, msgAcknowledgeBridges)
	if err != nil {
		return BridgeTxResponse{}, err
	}
	if len(tx) == 0 {
		return BridgeTxResponse{}, fmt.Errorf("Invalid tx: %v", tx)
	}

	return BridgeTxResponse{
		Tx:         tx,
		NumBridges: len(msgAcknowledgeBridges.Events),
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

// EncodeProposedOperationsTx encodes a MsgProposedOperations into a tx.
func EncodeProposedOperationsTx(
	txConfig client.TxConfig,
	msgOperations *clobtypes.MsgProposedOperations,
) (OperationsTxResponse, error) {
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

// GetProposedOperationsTx returns a tx containing `MsgProposedOperations`.
func GetProposedOperationsTx(
	ctx sdk.Context,
	txConfig client.TxConfig,
	clobKeeper PrepareClobKeeper,
) (OperationsTxResponse, error) {
	msgOperations := clobKeeper.GetOperations(ctx)
	if msgOperations == nil {
		return OperationsTxResponse{}, fmt.Errorf("MsgProposedOperations cannot be nil")
	}

	return EncodeProposedOperationsTx(txConfig, msgOperations)
}

// ApplyClobMsgsToMemClob replays CLOB cancel/place messages on a check-tx cache context so the memclob
// reflects those changes before computing proposed operations.
func ApplyClobMsgsToMemClob(
	ctx sdk.Context,
	decoder sdk.TxDecoder,
	clobKeeper PrepareClobKeeper,
	txs [][]byte,
) error {
	for i, txBytes := range txs {
		tx, err := decoder(txBytes)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("ApplyCancelsToMemClob: failed to decode tx (index %v of %v): %v", i, len(txs), err))
			continue
		}

		txCtx := ctx.WithTxBytes(txBytes)

		for _, msg := range tx.GetMsgs() {
			var applyErr error
			switch m := msg.(type) {
			case *clobtypes.MsgCancelOrder:
				if m.OrderId.IsStatefulOrder() {
					applyErr = clobKeeper.CancelStatefulOrder(txCtx, m)
				} else {
					applyErr = clobKeeper.CancelShortTermOrder(txCtx, m)
				}
			case *clobtypes.MsgBatchCancel:
				_, _, applyErr = clobKeeper.BatchCancelShortTermOrder(txCtx, m)
			case *clobtypes.MsgPlaceOrder:
				if m.Order.OrderId.IsStatefulOrder() {
					applyErr = clobKeeper.PlaceStatefulOrder(txCtx, m, false)
				} else {
					_, _, applyErr = clobKeeper.PlaceShortTermOrder(txCtx, m)
				}
			default:
				continue
			}

			if applyErr != nil {
				return applyErr
			}
		}
	}

	return nil
}

// ExtractCancelledOrderIds builds a set of order IDs cancelled by the provided txs.
func ExtractCancelledOrderIds(
	ctx sdk.Context,
	decoder sdk.TxDecoder,
	txs [][]byte,
) map[clobtypes.OrderId]struct{} {
	cancelled := make(map[clobtypes.OrderId]struct{})
	for i, txBytes := range txs {
		tx, err := decoder(txBytes)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("ExtractCancelledOrderIds: failed to decode tx (index %v of %v): %v", i, len(txs), err))
			continue
		}
		for _, msg := range tx.GetMsgs() {
			switch m := msg.(type) {
			case *clobtypes.MsgCancelOrder:
				cancelled[m.GetOrderId()] = struct{}{}
			case *clobtypes.MsgBatchCancel:
				subaccount := m.GetSubaccountId()
				for _, batch := range m.GetShortTermCancels() {
					clobPairId := batch.GetClobPairId()
					for _, clientId := range batch.GetClientIds() {
						orderId := clobtypes.OrderId{
							SubaccountId: subaccount,
							ClientId:     clientId,
							ClobPairId:   clobPairId,
							OrderFlags:   clobtypes.OrderIdFlags_ShortTerm,
						}
						cancelled[orderId] = struct{}{}
					}
				}
			}
		}
	}
	return cancelled
}

// FilterMatchesByCancelledOrders drops match operations that reference any cancelled order ID.
func FilterMatchesByCancelledOrders(
	msgOperations *clobtypes.MsgProposedOperations,
	cancelled map[clobtypes.OrderId]struct{},
) *clobtypes.MsgProposedOperations {
	if len(cancelled) == 0 || msgOperations == nil {
		return msgOperations
	}

	filtered := &clobtypes.MsgProposedOperations{
		OperationsQueue: make([]clobtypes.OperationRaw, 0, len(msgOperations.GetOperationsQueue())),
	}

	for _, op := range msgOperations.GetOperationsQueue() {
		match := op.GetMatch()
		if match != nil && matchContainsCancelledOrder(match, cancelled) {
			continue
		}
		filtered.OperationsQueue = append(filtered.OperationsQueue, op)
	}

	return filtered
}

func matchContainsCancelledOrder(match *clobtypes.ClobMatch, cancelled map[clobtypes.OrderId]struct{}) bool {
	switch m := match.Match.(type) {
	case *clobtypes.ClobMatch_MatchOrders:
		mo := m.MatchOrders
		if mo == nil {
			return false
		}
		taker := mo.GetTakerOrderId()
		if _, ok := cancelled[taker]; ok {
			return true
		}
		for _, fill := range mo.GetFills() {
			maker := fill.GetMakerOrderId()
			if _, ok := cancelled[maker]; ok {
				return true
			}
		}
	default:
		// Other match types are not filtered.
	}
	return false
}
