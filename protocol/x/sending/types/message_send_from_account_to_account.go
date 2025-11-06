package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSendFromAccountToAccount{}

// NewMsgSendFromAccountToAccount constructs a `MsgSendFromAccountToAccount` from an
// authority, sender address, recipient address, and a coin.
func NewMsgSendFromAccountToAccount(
	authority string,
	sender string,
	recipient string,
	coin sdk.Coin,
) *MsgSendFromAccountToAccount {
	return &MsgSendFromAccountToAccount{
		Authority: authority,
		Sender:    sender,
		Recipient: recipient,
		Coin:      coin,
	}
}

// ValidateBasic runs validation on the fields of a MsgSendFromAccountToAccount.
func (msg *MsgSendFromAccountToAccount) ValidateBasic() error {
	// Validate authority.
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(
			ErrInvalidAuthority,
			fmt.Sprintf(
				"authority '%s' must be a valid bech32 address, but got error '%v'",
				msg.Authority,
				err.Error(),
			),
		)
	}

	// Validate sender address.
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return ErrInvalidAccountAddress
	}

	// Validate recipient address.
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return ErrInvalidAccountAddress
	}

	// Validate coin.
	if err := msg.Coin.Validate(); err != nil {
		return err
	}

	return nil
}
