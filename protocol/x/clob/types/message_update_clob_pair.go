package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateClobPair{}

// GetSigners requires that the MsgUpdateClobPair message is signed by the gov module.
func (msg *MsgUpdateClobPair) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic validates that the message's ClobPair status is a supported status.
func (msg *MsgUpdateClobPair) ValidateBasic() error {
	return msg.ClobPair.Validate()
}
