package types

import errorsmod "cosmossdk.io/errors"

// x/delaymsg module sentinel errors
var (
	ErrInvalidInput     = errorsmod.Register(ModuleName, 1, "Invalid input")
	ErrMsgIsNil         = errorsmod.Register(ModuleName, 2, "Delayed msg is nil")
	ErrMsgIsUnroutable  = errorsmod.Register(ModuleName, 3, "Message not recognized by router")
	ErrInvalidSigner    = errorsmod.Register(ModuleName, 4, "Invalid signer")
	ErrInvalidAuthority = errorsmod.Register(ModuleName, 5, "Invalid authority")

	ErrInvalidGenesisState = errorsmod.Register(ModuleName, 10, "Invalid genesis state")
)
