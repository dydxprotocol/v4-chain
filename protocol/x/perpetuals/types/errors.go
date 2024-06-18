package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/perpetuals module sentinel errors
var (
	ErrPerpetualDoesNotExist = errorsmod.Register(
		ModuleName,
		1,
		"Perpetual does not exist",
	)
	ErrMarketDoesNotExist = errorsmod.Register(
		ModuleName,
		2,
		"MarketId on perpetual does not exist",
	)
	ErrInitialMarginPpmExceedsMax = errorsmod.Register(
		ModuleName,
		3,
		"InitialMarginPpm exceeds maximum value of 1e6",
	)
	ErrMaintenanceFractionPpmExceedsMax = errorsmod.Register(
		ModuleName,
		4,
		"MaintenanceFractionPpm exceeds maximum value of 1e6",
	)
	ErrDefaultFundingPpmMagnitudeExceedsMax = errorsmod.Register(
		ModuleName,
		5,
		"DefaultFundingPpm magnitude exceeds maximum value of 1e6",
	)
	ErrTickerEmptyString = errorsmod.Register(
		ModuleName,
		6,
		"Ticker must be non-empty string",
	)
	ErrNoNewPremiumVotes = errorsmod.Register(
		ModuleName,
		7,
		"No new premium votes were collected",
	)
	ErrMoreFundingSamplesThanExpected = errorsmod.Register(
		ModuleName,
		8,
		"Recorded more than expected funding samples in the past funding-tick epoch",
	)
	ErrInvalidAddPremiumVotes = errorsmod.Register(ModuleName, 9, "MsgAddPremiumVotes is invalid")
	ErrPremiumVoteNotClamped  = errorsmod.Register(
		ModuleName,
		10,
		"Premium vote value is not clamped by MaxAbsPremiumVotePpm",
	)
	ErrFundingRateInt32Overflow = errorsmod.Register(
		ModuleName,
		11,
		"Funding rate int32 overflow",
	)
	ErrLiquidityTierDoesNotExist = errorsmod.Register(
		ModuleName,
		12,
		"Liquidity Tier does not exist",
	)
	ErrFundingRateClampFactorPpmIsZero = errorsmod.Register(
		ModuleName,
		14,
		"Funding rate clamp factor ppm is zero",
	)
	ErrPremiumVoteClampFactorPpmIsZero = errorsmod.Register(
		ModuleName,
		15,
		"Premium vote clamp factor ppm is zero",
	)
	ErrImpactNotionalIsZero = errorsmod.Register(
		ModuleName,
		16,
		"Impact notional is zero",
	)
	ErrPerpetualAlreadyExists = errorsmod.Register(
		ModuleName,
		17,
		"Perpetual already exists",
	)
	ErrPremiumVoteForNonActiveMarket = errorsmod.Register(
		ModuleName,
		18,
		"Premium votes are disallowed for non active markets",
	)
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		19,
		"Authority is invalid",
	)
	ErrLiquidityTierAlreadyExists = errorsmod.Register(
		ModuleName,
		20,
		"Liquidity tier already exists",
	)
	ErrMaintenanceMarginLargerThanInitialMargin = errorsmod.Register(
		ModuleName,
		21,
		"Maintenance margin fraction is larger than initial margin fraction",
	)
	ErrMinNumVotesPerSampleIsZero = errorsmod.Register(
		ModuleName,
		22,
		"MinNumVotesPerSample is zero",
	)
	ErrInvalidMarketType = errorsmod.Register(
		ModuleName,
		23,
		"Market type is invalid",
	)
	ErrOpenInterestLowerCapLargerThanUpperCap = errorsmod.Register(
		ModuleName,
		24,
		"open interest lower cap is larger than upper cap",
	)
	ErrOpenInterestWouldBecomeNegative = errorsmod.Register(
		ModuleName,
		25,
		"open interest would become negative after update",
	)
	ErrPerpetualInfoDoesNotExist = errorsmod.Register(
		ModuleName,
		26,
		"PerpetualInfo does not exist",
	)

	// Errors for Not Implemented
	ErrNotImplementedFunding = errorsmod.Register(ModuleName, 1001, "Not Implemented: Perpetuals Funding")
)
