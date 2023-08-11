package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgUpdateEventParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateEventParams) ValidateBasic() error {
	return nil
}
