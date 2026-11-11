package types

import (
	"time"
)

// ShortBlockWindow represents the maximum number of blocks past the current block height that a
// `MsgPlaceOrder` or `MsgCancelOrder` message will be considered valid by the validator.
const ShortBlockWindow uint32 = 40

// MaxMsgBatchCancelBatchSize represents the maximum number of cancels that a MsgBatchCancel
// can have in one Msg.
const MaxMsgBatchCancelBatchSize uint32 = 100

// MaxShortTermClobTxBytes bounds the size of a short-term CLOB tx accepted into the local mempool.
// Short-term CLOB txs are fee- and gas-free and their raw bytes are re-emitted verbatim in the
// operations tx, so padding inflates the proposed block. A legitimate single-signer short-term tx is
// well under 4 KiB, so 8 KiB is a generous node-local CheckTx bound, far below max_tx_bytes.
const MaxShortTermClobTxBytes = 8 * 1024

// StatefulOrderTimeWindow represents the maximum amount of time in seconds past the current block time that a
// long-term/conditional `MsgPlaceOrder` message will be considered valid by the validator.
const StatefulOrderTimeWindow time.Duration = 95 * 24 * time.Hour // 95 days.

// ConditionalOrderTriggerMultiplier represents the multiplier used to calculate the upper and lower bounds of
// the trigger price for a conditional order.
// The upper bound is calculated as:
//
//	upper_bound = (1 + min_price_change_ppm / 1_000_000 * conditional_order_trigger_multiplier) * oracle_price
//
// The lower bound is calculated as:
//
//	lower_bound = (1 - min_price_change_ppm / 1_000_000 * conditional_order_trigger_multiplier) * oracle_price
const ConditionalOrderTriggerMultiplier uint64 = 5
