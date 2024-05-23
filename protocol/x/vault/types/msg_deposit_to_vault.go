package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
)

var _ sdk.Msg = &MsgDepositToVault{}

// ValidateBasic performs stateless validation on a MsgDepositToVault.
func (msg *MsgDepositToVault) ValidateBasic() error {
	// Validate subaccount to deposit from.
	if err := msg.SubaccountId.Validate(); err != nil {
		return err
	}

	// Validate that quote quantums is positive and an uint64.
	if msg.QuoteQuantums.Cmp(dtypes.NewInt(0)) <= 0 || !msg.QuoteQuantums.BigInt().IsUint64() {
		return errorsmod.Wrap(ErrInvalidDepositAmount, "quote quantums must be strictly positive and less than 2^64")
	}

	return nil
}
