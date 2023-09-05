package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/vest module sentinel errors
var (
	ErrInvalidVesterAccount    = sdkerrors.Register(ModuleName, 1001, "invalid vester account")
	ErrInvalidTreasuryAccount  = sdkerrors.Register(ModuleName, 1002, "invalid treasury account")
	ErrInvalidDenom            = sdkerrors.Register(ModuleName, 1003, "invalid denom")
	ErrVestEntryNotFound       = sdkerrors.Register(ModuleName, 1004, "account is not associated with a vest entry")
	ErrInvalidStartAndEndTimes = sdkerrors.Register(ModuleName, 1005, "start_time must be before end_time")
	ErrInvalidTimeZone         = sdkerrors.Register(ModuleName, 1006, "timestamp must be in UTC")
)
