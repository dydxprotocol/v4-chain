package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ sdk.Msg = &MsgUpdateMarketParam{}

func (msg *MsgUpdateMarketParam) GetSigners() []sdk.AccAddress {
	// TODO(CORE-564): implement this method
	return []sdk.AccAddress{}
}

func (msg *MsgUpdateMarketParam) ValidateBasic() error {
	// TODO(CORE-564): implement this method
	return status.Errorf(codes.Unimplemented, "ValidateBasic not implemented")
}
