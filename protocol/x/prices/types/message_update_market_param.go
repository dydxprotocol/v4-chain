package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateMarketParam{}

func (msg *MsgUpdateMarketParam) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgUpdateMarketParam) ValidateBasic() error {
	if msg.Authority == "" {
		return errorsmod.Wrapf(ErrInvalidAuthority, "authority cannot be empty")
	}
	return msg.MarketParam.Validate()
}
