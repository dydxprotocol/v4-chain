package process

import moderrors "cosmossdk.io/errors"

// DONTCOVER

const (
	ModuleName = "process_proposal"
)

var (
	// 1 - 99: Default.
	ErrDecodingTxBytes   = moderrors.Register(ModuleName, 1, "Decoding tx bytes failed")
	ErrMsgValidateBasic  = moderrors.Register(ModuleName, 2, "ValidateBasic failed on msg")
	ErrUnexpectedNumMsgs = moderrors.Register(ModuleName, 3, "Unexpected num of msgs")
	ErrUnexpectedMsgType = moderrors.Register(ModuleName, 4, "Unexpected msg type")
)
