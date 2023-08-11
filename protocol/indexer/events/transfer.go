package events

import (
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// NewTransferEvent creates a TransferEvent representing a transfer of an asset between a sender
// and recipient subaccount.
func NewTransferEvent(
	senderSubaccountId satypes.SubaccountId,
	recipientSubaccountId satypes.SubaccountId,
	assetId uint32,
	amount satypes.BaseQuantums,
) *TransferEvent {
	return &TransferEvent{
		SenderSubaccountId:    senderSubaccountId,
		RecipientSubaccountId: recipientSubaccountId,
		AssetId:               assetId,
		Amount:                amount.ToUint64(),
	}
}
