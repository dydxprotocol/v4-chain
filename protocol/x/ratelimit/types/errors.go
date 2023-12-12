package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/ratelimit module sentinel errors
var (
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		1001,
		"Authority is invalid",
	)
	ErrWithdrawalExceedsCapacity = errorsmod.Register(
		ModuleName,
		1002,
		"withdrawal amount would exceed rate-limit capacity",
	)
)
