package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// x/delaymsg module sentinel errors
var (
	ErrInvalidInput        = sdkerrors.Register(ModuleName, 1, "Invalid input")
	ErrMsgIsNil            = sdkerrors.Register(ModuleName, 2, "Delayed msg is nil")
	ErrInvalidGenesisState = sdkerrors.Register(ModuleName, 3, "Invalid genesis state")
)
