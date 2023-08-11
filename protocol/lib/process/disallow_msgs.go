package process

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
)

// IsDisallowTopLevelMsgInOtherTxs returns true if the given msg type is not allowed
// to be in a proposed block as part of OtherTxs. Otherwise, returns false.
func IsDisallowTopLevelMsgInOtherTxs(targetMsg sdk.Msg) bool {
	switch targetMsg.(type) {
	case *clobtypes.MsgCancelOrder, *clobtypes.MsgPlaceOrder:
		// We do not expect PlaceOrder/CancelOrder to be in OtherTxs as a top-level msg.
		// These should be app-injected as part of order tx.
		return true
	}
	return false
}
