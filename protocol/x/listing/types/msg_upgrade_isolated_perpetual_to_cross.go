package types

import (
	fmt "fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpgradeIsolatedPerpetualToCross{}

func (msg *MsgUpgradeIsolatedPerpetualToCross) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpgradeIsolatedPerpetualToCross) ValidateBasic() error {
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
	if msg.PerpetualId == 0 {
		return ErrInvalidPerpetualId
	}
	return nil
}