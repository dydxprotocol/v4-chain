package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (msg *MsgUpdateEquityTierLimitConfiguration) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateEquityTierLimitConfiguration) ValidateBasic() error {
	return msg.EquityTierLimitConfig.Validate()
}
