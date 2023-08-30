package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgUpdateDowntimeParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateDowntimeParams) ValidateBasic() error {
	return msg.Params.Validate()
}

func (msg *MsgIsDelayedBlock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgIsDelayedBlock) ValidateBasic() error {
	return nil
}
