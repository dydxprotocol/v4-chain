package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateLiquidationsConfig{}

// GetSigners requires that the MsgUpdateLiquidationsConfig message is signed by the gov module.
func (msg *MsgUpdateLiquidationsConfig) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic validates the message's LiquidationConfig.
func (msg *MsgUpdateLiquidationsConfig) ValidateBasic() error {
	return msg.LiquidationsConfig.Validate()
}
