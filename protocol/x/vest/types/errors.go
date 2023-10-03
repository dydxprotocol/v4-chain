package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/vest module sentinel errors
var (
	ErrInvalidAuthority        = errorsmod.Register(ModuleName, 1000, "invalid authority")
	ErrInvalidVesterAccount    = errorsmod.Register(ModuleName, 1001, "invalid vester account")
	ErrInvalidTreasuryAccount  = errorsmod.Register(ModuleName, 1002, "invalid treasury account")
	ErrInvalidDenom            = errorsmod.Register(ModuleName, 1003, "invalid denom")
	ErrVestEntryNotFound       = errorsmod.Register(ModuleName, 1004, "account is not associated with a vest entry")
	ErrInvalidStartAndEndTimes = errorsmod.Register(ModuleName, 1005, "start_time must be before end_time")
	ErrInvalidTimeZone         = errorsmod.Register(ModuleName, 1006, "timestamp must be in UTC")
)
