package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// ShouldSkipSequenceValidation returns whether sequence validation can be skipped for a given list of messages.
// If the transaction consists of only messages which use `GoodTilBlock` to prevent transaction replay,
// the sequence numbers for this transaction won't get incremented and verified.
//
// Important:
// ALL messages from a transaction must use `GoodTilBlock` in order to skip sequence number validation.
// Otherwise, attackers can create transactions with a single `GoodTilBlock` message, followed by any number of messages
// that they wish to be replayed which normally use sequence numbers. This would cause the sequence validation
// to be skipped for all of those messages and this transaction could be replayed.
func ShouldSkipSequenceValidation(msgs []sdk.Msg) (shouldSkipValidation bool) {
	for _, msg := range msgs {
		switch typedMsg := msg.(type) {
		case
			*clobtypes.MsgPlaceOrder:
			// Stateful orders need to use sequence numbers for replay prevention.
			orderId := typedMsg.GetOrder().OrderId
			if orderId.IsStatefulOrder() {
				return false
			}
			// This is a short term order, continue to check the next message.
			continue
		case
			*clobtypes.MsgCancelOrder:
			orderId := typedMsg.GetOrderId()
			if orderId.IsStatefulOrder() {
				return false
			}
			// This is a `GoodTilBlock` message, continue to check the next message.
			continue
		case
			*clobtypes.MsgBatchCancel:
			// MsgBatchCancel only supports short term orders.
			continue
		default:
			// Early return for messages that require sequence number validation.
			return false
		}
	}
	// All messages use `GoodTilBlock`.
	return true
}
