package types

import (
	fmt "fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Subaccounts module event types.
const (
	EventTypeSettledFunding = "settled_funding"

	AttributeKeySubaccount       = "subaccount"
	AttributeKeySubaccountNumber = "subaccount_number"
	AttributeKeyPerpetualId      = "perpetual_id"
	AttributeKeyFundingPaid      = "funding_paid_quote_quantums"
)

// NewCreateSettledFundingEvent constructs a new funding sdk.Event. Note that `fundingPaid` is positive
// if the subaccount paid funding, negative if the subaccount received funding.
func NewCreateSettledFundingEvent(
	subaccount SubaccountId,
	perpetualId uint32,
	fundingPaid *big.Int,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeSettledFunding,
		sdk.NewAttribute(AttributeKeySubaccount, subaccount.Owner),
		sdk.NewAttribute(AttributeKeySubaccount, fmt.Sprint(subaccount.Number)),
		sdk.NewAttribute(AttributeKeyPerpetualId, fmt.Sprint(perpetualId)),
		sdk.NewAttribute(AttributeKeyFundingPaid, fundingPaid.String()),
	)
}
