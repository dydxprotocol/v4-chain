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
