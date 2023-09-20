package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSetLiquidityTier{}

func (msg *MsgSetLiquidityTier) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgSetLiquidityTier) ValidateBasic() error {
	if msg.Authority == "" {
		return errorsmod.Wrap(ErrInvalidAuthority, "authority cannot be empty")
	}
	return msg.LiquidityTier.Validate()
}
