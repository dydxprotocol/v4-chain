package types

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = &MsgDepositToVault{}

// ValidateBasic performs stateless validation on a MsgDepositToVault.
func (msg *MsgDepositToVault) ValidateBasic() error {
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
		return errors.Wrap(ErrInvalidDepositAmount, "quote quantums must be strictly positive and less than 2^64")
	}

	return nil
}
