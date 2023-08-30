package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/perpetuals module sentinel errors
var (
	ErrPerpetualDoesNotExist = sdkerrors.Register(
		ModuleName,
		1,
		"Perpetual does not exist",
	)
	ErrMarketDoesNotExist = sdkerrors.Register(
		ModuleName,
		2,
		"MarketId on perpetual does not exist",
	)
	ErrInitialMarginPpmExceedsMax = sdkerrors.Register(
		ModuleName,
		3,
		"InitialMarginPpm exceeds maximum value of 1e6",
	)
	ErrMaintenanceFractionPpmExceedsMax = sdkerrors.Register(
		ModuleName,
		4,
		"MaintenanceFractionPpm exceeds maximum value of 1e6",
	)
	ErrDefaultFundingPpmMagnitudeExceedsMax = sdkerrors.Register(
		ModuleName,
		5,
		"DefaultFundingPpm magnitude exceeds maximum value of 1e6",
	)
	ErrTickerEmptyString = sdkerrors.Register(
		ModuleName,
		6,
		"Ticker must be non-empty string",
	)
	ErrNoNewPremiumVotes = sdkerrors.Register(
		ModuleName,
		7,
		"No new premium votes were collected",
	)
	ErrMoreFundingSamplesThanExpected = sdkerrors.Register(
		ModuleName,
		8,
		"Recorded more than expected funding samples in the past funding-tick epoch",
	)
	ErrInvalidAddPremiumVotes = sdkerrors.Register(ModuleName, 9, "MsgAddPremiumVotes is invalid")
	ErrPremiumVoteNotClamped  = sdkerrors.Register(
		ModuleName,
		10,
		"Premium vote value is not clamped by MaxAbsPremiumVotePpm",
	)
	ErrFundingRateInt32Overflow = sdkerrors.Register(
		ModuleName,
		11,
		"Funding rate int32 overflow",
	)
	ErrLiquidityTierDoesNotExist = sdkerrors.Register(
		ModuleName,
		12,
		"Liquidity Tier does not exist",
	)
	ErrBasePositionNotionalIsZero = sdkerrors.Register(
		ModuleName,
		13,
		"Base position notional is zero",
	)
	ErrFundingRateClampFactorPpmIsZero = sdkerrors.Register(
		ModuleName,
		14,
		"Funding rate clamp factor ppm is zero",
	)
	ErrPremiumVoteClampFactorPpmIsZero = sdkerrors.Register(
		ModuleName,
		15,
		"Premium vote clamp factor ppm is zero",
	)
	ErrImpactNotionalIsZero = sdkerrors.Register(
		ModuleName,
		16,
		"Impact notional is zero",
	)
	ErrPerpetualAlreadyExists = sdkerrors.Register(
		ModuleName,
		17,
		"Perpetual already exists",
	)
	ErrPremiumVoteForInitializingMarket = sdkerrors.Register(
		ModuleName,
		18,
		"Premium votes are disallowed for initializing markets",
	)

	// Errors for Not Implemented
	ErrNotImplementedFunding      = sdkerrors.Register(ModuleName, 1001, "Not Implemented: Perpetuals Funding")
	ErrNotImplementedOpenInterest = sdkerrors.Register(ModuleName, 1002, "Not Implemented: Perpetuals Open Interest")
)
