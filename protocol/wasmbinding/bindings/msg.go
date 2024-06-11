package bindings

import (
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"

	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

type SendingMsg struct {
	CreateTransfer         *sendingtypes.MsgCreateTransfer         `json:"create_transfer,omitempty"`
	DepositToSubaccount    *sendingtypes.MsgDepositToSubaccount    `json:"deposit_to_subaccount,omitempty"`
	WithdrawFromSubaccount *sendingtypes.MsgWithdrawFromSubaccount `json:"withdraw_from_subaccount,omitempty"`
	PlaceOrder             *clobtypes.MsgPlaceOrder                `json:"place_order,omitempty"`
	CancelOrder            *clobtypes.MsgCancelOrder               `json:"cancel_order,omitempty"`
}
