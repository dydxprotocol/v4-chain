package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

var (
	ErrNoTiersExist = moderrors.Register(
		ModuleName,
		400,
		"Must have at least one fee tier",
	)
	ErrInvalidFirstTierRequirements = moderrors.Register(
		ModuleName,
		401,
		"First fee tier must not have volume requirements",
	)
	ErrTiersOutOfOrder = moderrors.Register(
		ModuleName,
		402,
		"Fee tiers must have ascending requirements",
	)
	ErrInvalidFee = moderrors.Register(
		ModuleName,
		403,
		"No maker and taker fee combination should result in a net rebate",
	)
)
