package prepare

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	gometrics "github.com/hashicorp/go-metrics"

	"github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// GetGroupMsgOther returns two separate slices of byte txs given a single slice of byte txs and max bytes.
// The first slice contains the first N txs where the total bytes of the N txs is <= max bytes.
// The second slice contains the rest of txs, if any.
func GetGroupMsgOther(availableTxs [][]byte, maxBytes uint64) ([][]byte, [][]byte) {
	var (
		txsToInclude [][]byte
		txsRemainder [][]byte
		byteCount    uint64
	)

	for _, tx := range availableTxs {
		byteCount += uint64(len(tx))
		if byteCount <= maxBytes {
			txsToInclude = append(txsToInclude, tx)
		} else {
			txsRemainder = append(txsRemainder, tx)
		}
	}

	return txsToInclude, txsRemainder
}

// RemoveDisallowMsgs removes any txs that contain a disallowed msg.
func RemoveDisallowMsgs(
	ctx sdk.Context,
	decoder sdk.TxDecoder,
	txs [][]byte,
) [][]byte {
	defer telemetry.ModuleMeasureSince(
		ModuleName,
		time.Now(),
		metrics.RemoveDisallowMsgs,
		metrics.Latency,
	)

	var filteredTxs [][]byte
	for i, txBytes := range txs {
		// Decode tx so we can read msgs.
		tx, err := decoder(txBytes)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("RemoveDisallowMsgs: failed to decode tx (index %v of %v txs): %v", i, len(txs), err))
			continue // continue to next tx.
		}

		// For each msg in tx, check if it is disallowed.
		containsDisallowMsg := false
		for _, msg := range tx.GetMsgs() {
			if ante.IsDisallowExternalSubmitMsg(msg) {
				telemetry.IncrCounterWithLabels(
					[]string{ModuleName, metrics.RemoveDisallowMsgs, metrics.DisallowMsg, metrics.Count},
					1,
					[]gometrics.Label{metrics.GetLabelForStringValue(metrics.Detail, proto.MessageName(msg))},
				)
				containsDisallowMsg = true
				break // break out of loop over msgs.
			}
		}

		// If tx contains disallowed msg, skip it.
		if containsDisallowMsg {
			ctx.Logger().Error(
				fmt.Sprintf("RemoveDisallowMsgs: skipping tx with disallowed msg. Size: %d", len(txBytes)))
			continue // continue to next tx.
		}

		// Otherwise, add tx to filtered txs.
		filteredTxs = append(filteredTxs, txBytes)
	}

	return filteredTxs
}

// ReorderClobCancelsFirst returns a new slice with CLOB cancel-only txs placed
// before all other txs, preserving relative order within each bucket.
// It also returns the number of cancel-only txs (the first N in the returned slice).
func ReorderClobCancelsFirst(
	ctx sdk.Context,
	decoder sdk.TxDecoder,
	txs [][]byte,
) (reordered [][]byte, numCancelOnly int) {
	if len(txs) == 0 {
		return txs, 0
	}

	cancelTxs := make([][]byte, 0, len(txs))
	otherTxs := make([][]byte, 0, len(txs))

	for i, txBytes := range txs {
		tx, err := decoder(txBytes)
		if err != nil {
			// Preserve tx ordering by treating decode failures as non-cancel txs.
			ctx.Logger().Error(
				fmt.Sprintf("ReorderClobCancelsFirst: failed to decode tx (index %v of %v txs): %v", i, len(txs), err))
			otherTxs = append(otherTxs, txBytes)
			continue
		}

		msgs := tx.GetMsgs()
		if len(msgs) == 0 {
			otherTxs = append(otherTxs, txBytes)
			continue
		}

		allCancels := true
		for _, msg := range msgs {
			switch msg.(type) {
			case *clobtypes.MsgCancelOrder, *clobtypes.MsgBatchCancel:
				// continue til we get non-cancel msg.
				continue
			default:
				allCancels = false
			}
		}

		if allCancels {
			cancelTxs = append(cancelTxs, txBytes)
			continue
		}

		otherTxs = append(otherTxs, txBytes)
	}

	return append(cancelTxs, otherTxs...), len(cancelTxs)
}

// IsTxCancelOnly returns true if the tx contains only CLOB cancel messages (MsgCancelOrder or MsgBatchCancel).
// It is used by tests to assert that PrepareProposal places cancel-only txs before other tx types.
func IsTxCancelOnly(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	if len(msgs) == 0 {
		return false
	}
	for _, msg := range msgs {
		switch msg.(type) {
		case *clobtypes.MsgCancelOrder, *clobtypes.MsgBatchCancel:
			continue
		default:
			return false
		}
	}
	return true
}
