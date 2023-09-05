package types

import (
	fmt "fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Subaccounts module event types.
const (
	EventTypeFunding = "funding"

	AttributeKeySubaccount       = "subaccount"
	AttributeKeySubaccountNumber = "subaccount_number"
	AttributeKeyPerpetualId      = "perpetual_id"
	AttributeKeyFundingPaid      = "funding_received_quote_quantums"
)

// NewCreateFundingEvent constructs a new funding sdk.Event. Note that `fundingPaid` is positive
// if the subaccount paid funding, negative if the subaccount received funding.
func NewCreateFundingEvent(
	subaccount SubaccountId,
	perpetualId uint32,
	fundingPaid *big.Int,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeFunding,
		sdk.NewAttribute(AttributeKeySubaccount, subaccount.Owner),
		sdk.NewAttribute(AttributeKeySubaccount, fmt.Sprint(subaccount.Number)),
		sdk.NewAttribute(AttributeKeyPerpetualId, fmt.Sprint(perpetualId)),
		sdk.NewAttribute(AttributeKeyFundingPaid, fundingPaid.String()),
	)
}
