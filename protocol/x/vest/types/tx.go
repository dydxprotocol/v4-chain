package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetVestEntry = "set_vest_entry"

var _ sdk.Msg = &MsgSetVestEntry{}

func NewMsgSetVestEntry(authority string, entry VestEntry) *MsgSetVestEntry {
	return &MsgSetVestEntry{
		Authority: authority,
		Entry:     entry,
	}
}

func (msg *MsgSetVestEntry) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrapf(
			ErrInvalidAuthority,
			"authority address is invalid: %v, err = %v",
			msg.Authority,
			err,
		)
	}

	return msg.Entry.Validate()
}

func NewMsgDeleteVestEntry(authority string, vesterAccount string) *MsgDeleteVestEntry {
	return &MsgDeleteVestEntry{
		Authority:     authority,
		VesterAccount: vesterAccount,
	}
}

func (msg *MsgDeleteVestEntry) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrapf(
			ErrInvalidAuthority,
			"authority address is invalid: %v, err = %v",
			msg.Authority,
			err,
		)
	}

	if msg.VesterAccount == "" {
		return errorsmod.Wrapf(
			ErrInvalidVesterAccount,
			"vester account cannot be empty",
		)
	}

	return nil
}
