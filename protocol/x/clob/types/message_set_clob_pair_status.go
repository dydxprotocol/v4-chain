package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var _ sdk.Msg = &MsgSetClobPairStatus{}

// GetSigners requires that the MsgSetClobPairStatus message is signed by the gov module.
func (msg *MsgSetClobPairStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		authtypes.NewModuleAddress(govtypes.ModuleName),
	}
}

func (msg *MsgSetClobPairStatus) ValidateBasic() error {
	if !IsSupportedClobPairStatus(msg.Status) {
		return sdkerrors.Wrapf(
			ErrInvalidMsgSetClobPairStatus,
			"Unsupported ClobPair status: %+v",
			msg.Status,
		)
	}
	return nil
}
