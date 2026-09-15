package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateControls{}

func signer(address string) []sdk.AccAddress {
	parsed, _ := sdk.AccAddressFromBech32(address)
	return []sdk.AccAddress{parsed}
}

func validateSigner(address string) error {
	if _, err := sdk.AccAddressFromBech32(address); err != nil {
		return fmt.Errorf("%w: invalid signer: %v", ErrUnauthorized, err)
	}
	return nil
}

func (msg *MsgUpdateControls) GetSigners() []sdk.AccAddress { return signer(msg.Authority) }

func (msg *MsgUpdateControls) ValidateBasic() error {
	if err := validateSigner(msg.Authority); err != nil {
		return err
	}
	return msg.Controls.Validate()
}
