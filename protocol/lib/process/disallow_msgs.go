package process

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// IsDisallowTopLevelMsgInOtherTxs returns true if the given msg type is not allowed
// to be in a proposed block as part of OtherTxs. Otherwise, returns false.
func IsDisallowTopLevelMsgInOtherTxs(targetMsg sdk.Msg) bool {
	switch msg := targetMsg.(type) {
	// Non-stateful cancel and place orders are not allowed in the proposed blocks.
	// These should be app-injected as part of MsgProposedOperation tx.
	case *clobtypes.MsgCancelOrder:
		orderId := msg.GetOrderId()
		return !orderId.IsStatefulOrder()
	case *clobtypes.MsgPlaceOrder:
		order := msg.GetOrder()
		orderId := order.GetOrderId()
		return !orderId.IsStatefulOrder()
	}
	return false
}
