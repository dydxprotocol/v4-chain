package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	return msg.Params.Validate()
}
