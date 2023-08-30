package events

import (
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// NewTransferEvent creates a TransferEvent representing a transfer of an asset between a sender
// and recipient subaccount.
func NewTransferEvent(
	senderSubaccountId satypes.SubaccountId,
	recipientSubaccountId satypes.SubaccountId,
	assetId uint32,
	amount satypes.BaseQuantums,
) *TransferEventV1 {
	indexerSenderSubaccountId := v1.SubaccountIdToIndexerSubaccountId(senderSubaccountId)
	indexerRecipientSubaccountId := v1.SubaccountIdToIndexerSubaccountId(recipientSubaccountId)
	return &TransferEventV1{
		SenderSubaccountId:    &indexerSenderSubaccountId,
		RecipientSubaccountId: &indexerRecipientSubaccountId,
		Sender: &SourceOfFunds{
			Source: &SourceOfFunds_SubaccountId{
				SubaccountId: &indexerSenderSubaccountId,
			},
		},
		Recipient: &SourceOfFunds{
			Source: &SourceOfFunds_SubaccountId{
				SubaccountId: &indexerRecipientSubaccountId,
			},
		},
		AssetId: assetId,
		Amount:  amount.ToUint64(),
	}
}

// NewDepositEvent creates a DepositEvent representing a deposit of an asset from a sender
// wallet address to a recipient subaccount.
func NewDepositEvent(
	senderAddress string,
	recipientSubaccountId satypes.SubaccountId,
	assetId uint32,
	amount satypes.BaseQuantums,
) *TransferEventV1 {
	indexerRecipientSubaccountId := v1.SubaccountIdToIndexerSubaccountId(recipientSubaccountId)
	return &TransferEventV1{
		AssetId: assetId,
		Amount:  amount.ToUint64(),
		Sender: &SourceOfFunds{
			Source: &SourceOfFunds_Address{
				Address: senderAddress,
			},
		},
		Recipient: &SourceOfFunds{
			Source: &SourceOfFunds_SubaccountId{
				SubaccountId: &indexerRecipientSubaccountId,
			},
		},
	}
}

// NewWithdrawEvent creates a WithdrawEvent representing a withdrawal of an asset from a sender
// subaccount to a recipient wallet address.
func NewWithdrawEvent(
	senderSubaccountId satypes.SubaccountId,
	recipientAddress string,
	assetId uint32,
	amount satypes.BaseQuantums,
) *TransferEventV1 {
	indexerSenderSubaccountId := v1.SubaccountIdToIndexerSubaccountId(senderSubaccountId)
	return &TransferEventV1{
		AssetId: assetId,
		Amount:  amount.ToUint64(),
		Sender: &SourceOfFunds{
			Source: &SourceOfFunds_SubaccountId{
				SubaccountId: &indexerSenderSubaccountId,
			},
		},
		Recipient: &SourceOfFunds{
			Source: &SourceOfFunds_Address{
				Address: recipientAddress,
			},
		},
	}
}
