package types

import moderrors "cosmossdk.io/errors"

// x/delaymsg module sentinel errors
var (
	ErrInvalidInput    = moderrors.Register(ModuleName, 1, "Invalid input")
	ErrMsgIsNil        = moderrors.Register(ModuleName, 2, "Delayed msg is nil")
	ErrMsgIsUnroutable = moderrors.Register(ModuleName, 3, "Message not recognized by router")
	ErrInvalidSigner   = moderrors.Register(ModuleName, 4, "Invalid signer")

	ErrInvalidGenesisState = moderrors.Register(ModuleName, 10, "Invalid genesis state")
)
