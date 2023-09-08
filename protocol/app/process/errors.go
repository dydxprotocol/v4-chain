package process

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

const (
	ModuleName = "process_proposal"
)

var (
	// 1 - 99: Default.
	ErrDecodingTxBytes   = errorsmod.Register(ModuleName, 1, "Decoding tx bytes failed")
	ErrMsgValidateBasic  = errorsmod.Register(ModuleName, 2, "ValidateBasic failed on msg")
	ErrUnexpectedNumMsgs = errorsmod.Register(ModuleName, 3, "Unexpected num of msgs")
	ErrUnexpectedMsgType = errorsmod.Register(ModuleName, 4, "Unexpected msg type")
)
