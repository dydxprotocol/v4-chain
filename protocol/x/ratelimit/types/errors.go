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
	ErrMismatchedCapacityLimitersLength = errorsmod.Register(
		ModuleName,
		1003,
		"capacity list length does not match number of limiters",
	)
	ErrInvalidRateLimitPeriod = errorsmod.Register(
		ModuleName,
		1004,
		"rate limit period should be positive",
	)
	ErrInvalidBaselineMinimum = errorsmod.Register(
		ModuleName,
		1005,
		"baseline_minimum should be positive",
	)
	ErrInvalidBaselineTvlPpm = errorsmod.Register(
		ModuleName,
		1006,
		"Baseline_tvl_ppm must in the range (0, 1_000_000)",
	)
	ErrInvalidInput = errorsmod.Register(
		ModuleName,
		1007,
		"Invalid input",
	)
)
