package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ sdk.Msg = &MsgUpdateParams{}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	// TODO(CORE-560): implement this method
	return []sdk.AccAddress{}
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	// TODO(CORE-560): implement this method
	return status.Errorf(codes.Unimplemented, "ValidateBasic not implemented")
}
