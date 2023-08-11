package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateParams = "update_params"

var _ sdk.Msg = &MsgUpdateParams{}

func NewMsgUpdateParams(params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Params: params,
	}
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	return msg.Params.Validate()
}
