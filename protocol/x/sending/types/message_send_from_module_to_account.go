package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ sdk.Msg = &MsgSendFromModuleToAccount{}

// NewMsgWithdrawFromSubaccount constructs a `MsgWithdrawFromSubaccount` from an
// `x/subaccounts` subaccount sender, an `x/bank` account recipient, an asset ID,
// and a number of quantums.
func NewMsgSendFromModuleToAccount(
	authority string,
	senderModuleName string,
	recipient string,
	coin sdk.Coin,
) *MsgSendFromModuleToAccount {
	return &MsgSendFromModuleToAccount{
		Authority:        authority,
		SenderModuleName: senderModuleName,
		Recipient:        recipient,
		Coin:             coin,
	}
}

func (msg *MsgSendFromModuleToAccount) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs validation on the fields of a MsgSendFromModuleToAccount.
func (msg *MsgSendFromModuleToAccount) ValidateBasic() error {
	return status.Errorf(codes.Unimplemented, "ValidateBasic not implemented")
}
