package types

import (
	"time"
)

// ShortBlockWindow represents the maximum number of blocks past the current block height that a
// `MsgPlaceOrder` or `MsgCancelOrder` message will be considered valid by the validator.
const ShortBlockWindow uint32 = 20

// StatefulOrderTimeWindow represents the maximum amount of time in seconds past the current block time that a
// long-term/conditional `MsgPlaceOrder` message will be considered valid by the validator.
const StatefulOrderTimeWindow time.Duration = 95 * 24 * time.Hour // 95 days.

// ConditionalOrderTriggerMultiplier represents the multiplier used to calculate the upper and lower bounds of
// the trigger price for a conditional order.
const ConditionalOrderTriggerMultiplier uint64 = 5
