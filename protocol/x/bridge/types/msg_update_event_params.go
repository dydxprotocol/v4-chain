package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgUpdateEventParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateEventParams) ValidateBasic() error {
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
	return nil
}
