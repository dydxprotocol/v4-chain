package types

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = &MsgWithdrawFromMegavault{}

// ValidateBasic performs stateless validation on a MsgWithdrawFromMegavault.
func (msg *MsgWithdrawFromMegavault) ValidateBasic() error {
	// Validate subaccount to withdraw to.
	if err := msg.SubaccountId.Validate(); err != nil {
		return err
	}

	// Validate that shares is positive.
	if msg.Shares.NumShares.Sign() <= 0 {
		return ErrNonPositiveShares
	}

	// Validate that min quote quantums is non-negative and an uint64.
	quoteQuantums := msg.MinQuoteQuantums.BigInt()
	if quoteQuantums.Sign() < 0 || !quoteQuantums.IsUint64() {
		return errors.Wrap(
			ErrInvalidQuoteQuantums,
			"min quote quantums must be non-negative and less than 2^64",
		)
	}

	return nil
}
