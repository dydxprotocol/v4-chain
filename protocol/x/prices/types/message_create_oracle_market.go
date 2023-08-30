package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateOracleMarket{}

func (msg *MsgCreateOracleMarket) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgCreateOracleMarket) ValidateBasic() error {
	// TODO(CORE-504): Implement message validation.
	return errors.New("not implemented")
}
