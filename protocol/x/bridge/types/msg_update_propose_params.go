package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgUpdateProposeParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateProposeParams) ValidateBasic() error {
	return nil
}
