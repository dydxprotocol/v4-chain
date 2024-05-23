package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgDepositToVault{}

// ValidateBasic performs stateless validation on a MsgDepositToVault.
func (msg *MsgDepositToVault) ValidateBasic() error {
	// Validate subaccount to deposit from.
	if err := msg.SubaccountId.Validate(); err != nil {
		return err
	}

	// Validate that quote quantums is positive and an uint64.
	quoteQuantums := msg.QuoteQuantums.BigInt()
	if quoteQuantums.Sign() <= 0 || !quoteQuantums.IsUint64() {
		return errorsmod.Wrap(ErrInvalidDepositAmount, "quote quantums must be strictly positive and less than 2^64")
	}

	return nil
}
