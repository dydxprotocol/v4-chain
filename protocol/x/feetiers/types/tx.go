package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgUpdatePerpetualFeeParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdatePerpetualFeeParams) ValidateBasic() error {
	return msg.Params.Validate()
}
