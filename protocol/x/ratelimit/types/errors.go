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
	ErrInvalidSender = errorsmod.Register(
		ModuleName,
		1008,
		"Sender address is invalid",
	)
	ErrUnableToDecodeBigInt = errorsmod.Register(
		ModuleName,
		1009,
		"Unable to decode bigint",
	)
	ErrValueIsNegative = errorsmod.Register(
		ModuleName,
		1110,
		"Value is negative",
	)
	ErrInvalidSDAIConversionRate = errorsmod.Register(
		ModuleName,
		1111,
		"Proposed SDAI conversion rate is invalid",
	)
	ErrSDAIConversionRateNotInitisialised = errorsmod.Register(
		ModuleName,
		1112,
		"The SDAI rate has not been initialised",
	)
	ErrEpochNotStored = errorsmod.Register(
		ModuleName,
		1113,
		"Epoch info is not stored",
	)
	ErrEpochNotRetrieved = errorsmod.Register(
		ModuleName,
		1114,
		"Epoch info could not be retrieved from store",
	)
	ErrFailedSDaiToTDaiConversion = errorsmod.Register(
		ModuleName,
		1115,
		"Failed to convert sDai amount to corresponding TDai Amount",
	)
)
