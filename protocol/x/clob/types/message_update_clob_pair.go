package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateClobPair{}

// GetSigners requires that the MsgUpdateClobPair message is signed by the gov module.
func (msg *MsgUpdateClobPair) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic validates that the message's ClobPair status is a supported status.
func (msg *MsgUpdateClobPair) ValidateBasic() error {
	// TODO(CORE-504): Implement message validation, copy from MsgCreateClobPair.

	if !IsSupportedClobPairStatus(msg.ClobPair.Status) {
		return sdkerrors.Wrapf(
			ErrInvalidMsgUpdateClobPair,
			"Cannot set status for ClobPair with id %d to unsupported ClobPair status %s",
			msg.ClobPair.Id,
			msg.ClobPair.Status,
		)
	}
	return nil
}
