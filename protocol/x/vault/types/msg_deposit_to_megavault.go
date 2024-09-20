package types

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = &MsgDepositToMegavault{}

// ValidateBasic performs stateless validation on a MsgDepositToMegavault.
func (msg *MsgDepositToMegavault) ValidateBasic() error {
	// Note: msg signer must be the owner of the subaccount.
	// This is enforced by the following notatino on the msg proto:
	//    option (cosmos.msg.v1.signer) = "subaccount_id"

	// Validate subaccount to deposit from.
	if err := msg.SubaccountId.Validate(); err != nil {
		return err
	}

	// Validate that quote quantums is positive and an uint64.
	quoteQuantums := msg.QuoteQuantums.BigInt()
	if quoteQuantums.Sign() <= 0 || !quoteQuantums.IsUint64() {
		return errors.Wrap(
			ErrInvalidQuoteQuantums,
			"quote quantums must be strictly positive and less than 2^64",
		)
	}

	return nil
}
