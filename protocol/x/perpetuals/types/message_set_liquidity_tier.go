package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ sdk.Msg = &MsgSetLiquidityTier{}

func (msg *MsgSetLiquidityTier) GetSigners() []sdk.AccAddress {
	// TODO(CORE-563): implement this method
	return []sdk.AccAddress{}
}

func (msg *MsgSetLiquidityTier) ValidateBasic() error {
	// TODO(CORE-563): implement this method
	return status.Errorf(codes.Unimplemented, "ValidateBasic not implemented")
}
