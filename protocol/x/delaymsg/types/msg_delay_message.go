package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgDelayMessage) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic performs basic validation on the message.
func (msg *MsgDelayMessage) ValidateBasic() error {
	// Perform basic checks that the encoded message was set.
	if msg.Msg == nil || len(msg.Msg) == 0 {
		return ErrMsgIsNil
	}

	return nil
}
