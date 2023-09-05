package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/perpetuals module sentinel errors
var (
	ErrPerpetualDoesNotExist = moderrors.Register(
		ModuleName,
		1,
		"Perpetual does not exist",
	)
	ErrMarketDoesNotExist = moderrors.Register(
		ModuleName,
		2,
		"MarketId on perpetual does not exist",
	)
	ErrInitialMarginPpmExceedsMax = moderrors.Register(
		ModuleName,
		3,
		"InitialMarginPpm exceeds maximum value of 1e6",
	)
	ErrMaintenanceFractionPpmExceedsMax = moderrors.Register(
		ModuleName,
		4,
		"MaintenanceFractionPpm exceeds maximum value of 1e6",
	)
	ErrDefaultFundingPpmMagnitudeExceedsMax = moderrors.Register(
		ModuleName,
		5,
		"DefaultFundingPpm magnitude exceeds maximum value of 1e6",
	)
	ErrTickerEmptyString = moderrors.Register(
		ModuleName,
		6,
		"Ticker must be non-empty string",
	)
	ErrNoNewPremiumVotes = moderrors.Register(
		ModuleName,
		7,
		"No new premium votes were collected",
	)
	ErrMoreFundingSamplesThanExpected = moderrors.Register(
		ModuleName,
		8,
		"Recorded more than expected funding samples in the past funding-tick epoch",
	)
	ErrInvalidAddPremiumVotes = moderrors.Register(ModuleName, 9, "MsgAddPremiumVotes is invalid")
	ErrPremiumVoteNotClamped  = moderrors.Register(
		ModuleName,
		10,
		"Premium vote value is not clamped by MaxAbsPremiumVotePpm",
	)
	ErrFundingRateInt32Overflow = moderrors.Register(
		ModuleName,
		11,
		"Funding rate int32 overflow",
	)
	ErrLiquidityTierDoesNotExist = moderrors.Register(
		ModuleName,
		12,
		"Liquidity Tier does not exist",
	)
	ErrBasePositionNotionalIsZero = moderrors.Register(
		ModuleName,
		13,
		"Base position notional is zero",
	)
	ErrFundingRateClampFactorPpmIsZero = moderrors.Register(
		ModuleName,
		14,
		"Funding rate clamp factor ppm is zero",
	)
	ErrPremiumVoteClampFactorPpmIsZero = moderrors.Register(
		ModuleName,
		15,
		"Premium vote clamp factor ppm is zero",
	)
	ErrImpactNotionalIsZero = moderrors.Register(
		ModuleName,
		16,
		"Impact notional is zero",
	)
	ErrPerpetualAlreadyExists = moderrors.Register(
		ModuleName,
		17,
		"Perpetual already exists",
	)
	ErrPremiumVoteForNonActiveMarket = moderrors.Register(
		ModuleName,
		18,
		"Premium votes are disallowed for non active markets",
	)

	// Errors for Not Implemented
	ErrNotImplementedFunding      = moderrors.Register(ModuleName, 1001, "Not Implemented: Perpetuals Funding")
	ErrNotImplementedOpenInterest = moderrors.Register(ModuleName, 1002, "Not Implemented: Perpetuals Open Interest")
)
