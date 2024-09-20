package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = &MsgAllocateToVault{}

// ValidateBasic performs stateless validation on a MsgAllocateToVault.
func (msg *MsgAllocateToVault) ValidateBasic() error {
	if msg.Authority == "" {
		return ErrInvalidAuthority
	}

	// Validate that quote quantums is positive and an uint64.
	quoteQuantums := msg.QuoteQuantums.BigInt()
	if quoteQuantums.Sign() <= 0 || !quoteQuantums.IsUint64() {
		return ErrInvalidQuoteQuantums
	}

	return nil
}
