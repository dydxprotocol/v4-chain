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
