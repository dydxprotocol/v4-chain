package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSendFromModuleToAccount{}

// NewMsgWithdrawFromSubaccount constructs a `MsgWithdrawFromSubaccount` from an
// `x/subaccounts` subaccount sender, an `x/bank` account recipient, an asset ID,
// and a number of quantums.
func NewMsgSendFromModuleToAccount(
	authority string,
	senderModuleName string,
	recipient string,
	coin sdk.Coin,
) *MsgSendFromModuleToAccount {
	return &MsgSendFromModuleToAccount{
		Authority:        authority,
		SenderModuleName: senderModuleName,
		Recipient:        recipient,
		Coin:             coin,
	}
}

// ValidateBasic runs validation on the fields of a MsgSendFromModuleToAccount.
func (msg *MsgSendFromModuleToAccount) ValidateBasic() error {
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

	// Validate sender module name is non-empty.
	if len(msg.SenderModuleName) == 0 {
		return ErrEmptyModuleName
	}

	// Validate account recipient.
	_, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return ErrInvalidAccountAddress
	}

	// Validate coin.
	if err := msg.Coin.Validate(); err != nil {
		return err
	}

	return nil
}
