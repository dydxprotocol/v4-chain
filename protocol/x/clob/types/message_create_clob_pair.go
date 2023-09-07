package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateClobPair{}

func (msg *MsgCreateClobPair) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgCreateClobPair) ValidateBasic() error {
	if msg.Authority == "" {
		return sdkerrors.Wrap(ErrInvalidAuthority, "authority cannot be empty")
	}

	return msg.ClobPair.Validate()
}
