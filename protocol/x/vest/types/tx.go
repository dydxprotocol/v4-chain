package types

import (
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

func (msg *MsgSetVestEntry) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgSetVestEntry) ValidateBasic() error {
	return msg.Entry.Validate()
}

func NewMsgDeleteVestEntry(authority string, vesterAccount string) *MsgDeleteVestEntry {
	return &MsgDeleteVestEntry{
		Authority:     authority,
		VesterAccount: vesterAccount,
	}
}

func (msg *MsgDeleteVestEntry) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

func (msg *MsgDeleteVestEntry) ValidateBasic() error {
	// TODO
	return nil
}
