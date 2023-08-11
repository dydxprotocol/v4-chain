package events

import (
	"github.com/dydxprotocol/v4/indexer/protocol/v1"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// NewTransferEvent creates a TransferEvent representing a transfer of an asset between a sender
// and recipient subaccount.
func NewTransferEvent(
	senderSubaccountId satypes.SubaccountId,
	recipientSubaccountId satypes.SubaccountId,
	assetId uint32,
	amount satypes.BaseQuantums,
) *TransferEventV1 {
	return &TransferEventV1{
		SenderSubaccountId:    v1.SubaccountIdToIndexerSubaccountId(senderSubaccountId),
		RecipientSubaccountId: v1.SubaccountIdToIndexerSubaccountId(recipientSubaccountId),
		AssetId:               assetId,
		Amount:                amount.ToUint64(),
	}
}
