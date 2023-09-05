package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/vest module sentinel errors
var (
	ErrInvalidVesterAccount    = moderrors.Register(ModuleName, 1001, "invalid vester account")
	ErrInvalidTreasuryAccount  = moderrors.Register(ModuleName, 1002, "invalid treasury account")
	ErrInvalidDenom            = moderrors.Register(ModuleName, 1003, "invalid denom")
	ErrVestEntryNotFound       = moderrors.Register(ModuleName, 1004, "account is not associated with a vest entry")
	ErrInvalidStartAndEndTimes = moderrors.Register(ModuleName, 1005, "start_time must be before end_time")
	ErrInvalidTimeZone         = moderrors.Register(ModuleName, 1006, "timestamp must be in UTC")
)
