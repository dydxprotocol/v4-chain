package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateOracleMarket{}

func (msg *MsgCreateOracleMarket) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgCreateOracleMarket) ValidateBasic() error {
	if msg.Authority == "" {
		return sdkerrors.Wrap(ErrInvalidAuthority, "authority cannot be empty")
	}
	return msg.Params.Validate()
}
