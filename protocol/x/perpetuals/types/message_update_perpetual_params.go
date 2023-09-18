package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ sdk.Msg = &MsgUpdatePerpetualParams{}

func (msg *MsgUpdatePerpetualParams) GetSigners() []sdk.AccAddress {
	// TODO(CORE-562): implement this method
	return []sdk.AccAddress{}
}

func (msg *MsgUpdatePerpetualParams) ValidateBasic() error {
	// TODO(CORE-562): implement this method
	return status.Errorf(codes.Unimplemented, "ValidateBasic not implemented")
}
