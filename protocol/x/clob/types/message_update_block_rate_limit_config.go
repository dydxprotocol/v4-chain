package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (msg *MsgUpdateBlockRateLimitConfiguration) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateBlockRateLimitConfiguration) ValidateBasic() error {
	return msg.BlockRateLimitConfig.Validate()
}
