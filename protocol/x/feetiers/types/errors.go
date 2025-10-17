package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrNoTiersExist = errorsmod.Register(
		ModuleName,
		400,
		"Must have at least one fee tier",
	)
	ErrInvalidFirstTierRequirements = errorsmod.Register(
		ModuleName,
		401,
		"First fee tier must not have volume requirements",
	)
	ErrTiersOutOfOrder = errorsmod.Register(
		ModuleName,
		402,
		"Fee tiers must have ascending requirements",
	)
	ErrInvalidFee = errorsmod.Register(
		ModuleName,
		403,
		"No maker and taker fee combination should result in a net rebate",
	)
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		404,
		"Authority is invalid",
	)
	ErrRevShareSafetyViolation = errorsmod.Register(
		ModuleName,
		405,
		"Rev share safety violation",
	)
	ErrInvalidTimeRange = errorsmod.Register(
		ModuleName,
		406,
		"Invalid time range for market fee discount",
	)
	ErrMarketFeeDiscountNotFound = errorsmod.Register(
		ModuleName,
		407,
		"Market fee discount not found for the specified CLOB pair",
	)
	ErrInvalidChargePpm = errorsmod.Register(
		ModuleName,
		408,
		"Invalid charge PPM value",
	)
)
