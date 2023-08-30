package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// sending module event types
const (
	EventTypeCreateTransfer         = "create_transfer"
	EventTypeDepositToSubaccount    = "deposit_to_subaccount"
	EventTypeWithdrawFromSubaccount = "withdraw_from_subaccount"

	AttributeKeySender          = "sender"
	AttributeKeySenderNumber    = "sender_number"
	AttributeKeyRecipient       = "recipient"
	AttributeKeyRecipientNumber = "recipient_number"
	AttributeKeyQuantums        = "quantums"
	AttributeKeyAssetId         = "asset_id"
)

// NewCreateTransferEvent constructs a new create_transfer sdk.Event
func NewCreateTransferEvent(
	sender satypes.SubaccountId,
	recipient satypes.SubaccountId,
	assetId uint32,
	quantums uint64,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeCreateTransfer,
		sdk.NewAttribute(AttributeKeySender, sender.Owner),
		sdk.NewAttribute(AttributeKeySenderNumber, fmt.Sprintf("%d", sender.Number)),
		sdk.NewAttribute(AttributeKeyRecipient, recipient.Owner),
		sdk.NewAttribute(AttributeKeyRecipientNumber, fmt.Sprintf("%d", recipient.Number)),
		sdk.NewAttribute(AttributeKeyAssetId, fmt.Sprintf("%d", assetId)),
		sdk.NewAttribute(AttributeKeyQuantums, fmt.Sprintf("%d", quantums)),
	)
}

// NewDepositToSubaccountEvent a new deposit_to_subaccount sdk.Event
func NewDepositToSubaccountEvent(
	sender sdk.Address,
	recipient satypes.SubaccountId,
	assetId uint32,
	quantums uint64,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeDepositToSubaccount,
		sdk.NewAttribute(AttributeKeySender, sender.String()),
		sdk.NewAttribute(AttributeKeyRecipient, recipient.Owner),
		sdk.NewAttribute(AttributeKeyRecipientNumber, fmt.Sprintf("%d", recipient.Number)),
		sdk.NewAttribute(AttributeKeyAssetId, fmt.Sprintf("%d", assetId)),
		sdk.NewAttribute(AttributeKeyQuantums, fmt.Sprintf("%d", quantums)),
	)
}

// NewWithdrawFromSubaccount constructs a new withdraw_from_subaccount sdk.Event
func NewWithdrawFromSubaccountEvent(
	sender satypes.SubaccountId,
	recipient sdk.Address,
	assetId uint32,
	quantums uint64,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeWithdrawFromSubaccount,
		sdk.NewAttribute(AttributeKeySender, sender.Owner),
		sdk.NewAttribute(AttributeKeySenderNumber, fmt.Sprintf("%d", sender.Number)),
		sdk.NewAttribute(AttributeKeyRecipient, recipient.String()),
		sdk.NewAttribute(AttributeKeyAssetId, fmt.Sprintf("%d", assetId)),
		sdk.NewAttribute(AttributeKeyQuantums, fmt.Sprintf("%d", quantums)),
	)
}
