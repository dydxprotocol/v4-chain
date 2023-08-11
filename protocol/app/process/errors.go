package process

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ModuleName = "process_proposal"
)

var (
	// 1 - 99: Default.
	ErrDecodingTxBytes   = sdkerrors.Register(ModuleName, 1, "Decoding tx bytes failed")
	ErrMsgValidateBasic  = sdkerrors.Register(ModuleName, 2, "ValidateBasic failed on msg")
	ErrUnexpectedNumMsgs = sdkerrors.Register(ModuleName, 3, "Unexpected num of msgs")
	ErrUnexpectedMsgType = sdkerrors.Register(ModuleName, 4, "Unexpected msg type")
)
