package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNoTiersExist = sdkerrors.Register(
		ModuleName,
		400,
		"Must have at least one fee tier",
	)
	ErrInvalidFirstTierRequirements = sdkerrors.Register(
		ModuleName,
		401,
		"First fee tier must not have volume requirements",
	)
	ErrTiersOutOfOrder = sdkerrors.Register(
		ModuleName,
		402,
		"Fee tiers must have ascending requirements",
	)
	ErrInvalidFee = sdkerrors.Register(
		ModuleName,
		403,
		"No maker and taker fee combination should result in a net rebate",
	)
)
