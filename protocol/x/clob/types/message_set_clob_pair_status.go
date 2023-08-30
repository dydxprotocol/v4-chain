package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetClobPairStatus{}

// GetSigners requires that the MsgSetClobPairStatus message is signed by the gov module.
func (msg *MsgSetClobPairStatus) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic validates that the message's ClobPair status is a supported status.
func (msg *MsgSetClobPairStatus) ValidateBasic() error {
	if !IsSupportedClobPairStatus(ClobPair_Status(msg.ClobPairStatus)) {
		return sdkerrors.Wrapf(
			ErrInvalidMsgSetClobPairStatus,
			"Unsupported ClobPair status: %+v",
			msg.ClobPairStatus,
		)
	}
	return nil
}
