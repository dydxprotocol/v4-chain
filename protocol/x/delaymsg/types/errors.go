package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// x/delaymsg module sentinel errors
var (
	ErrInvalidInput    = sdkerrors.Register(ModuleName, 1, "Invalid input")
	ErrMsgIsNil        = sdkerrors.Register(ModuleName, 2, "Delayed msg is nil")
	ErrMsgIsUnroutable = sdkerrors.Register(ModuleName, 3, "Message not recognized by router")

	ErrInvalidGenesisState = sdkerrors.Register(ModuleName, 10, "Invalid genesis state")
)
