package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateLiquidationsConfig{}

// GetSigners requires that the MsgUpdateLiquidationsConfig message is signed by the gov module.
func (msg *MsgUpdateLiquidationsConfig) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic validates the message's LiquidationConfig. Returns an error if the authority
// is empty or if the LiquidationsConfig is invalid.
func (msg *MsgUpdateLiquidationsConfig) ValidateBasic() error {
	if msg.Authority == "" {
		return errorsmod.Wrap(ErrInvalidAuthority, "authority cannot be empty")
	}

	return msg.LiquidationsConfig.Validate()
}
