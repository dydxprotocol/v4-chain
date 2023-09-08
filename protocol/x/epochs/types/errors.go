package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/epochs module sentinel errors
var (
	ErrEmptyEpochInfoName = errorsmod.Register(
		ModuleName,
		2,
		"EpochInfo name is empty",
	)
	ErrDurationIsZero = errorsmod.Register(
		ModuleName,
		3,
		"Duration is zero",
	)
	ErrEpochInfoAlreadyExists = errorsmod.Register(
		ModuleName,
		4,
		"EpochInfo name already exists",
	)
	ErrEpochInfoNotFound = errorsmod.Register(
		ModuleName,
		5,
		"EpochInfo name not found",
	)
	ErrInvalidCurrentEpochAndCurrentEpochStartBlockTuple = errorsmod.Register(
		ModuleName,
		6,
		"Invalid CurrentEpoch and CurrentEpochStartBlock tuple: CurrentEpoch should"+
			" be zero if and only if CurrentEpochStartBlock is zero",
	)
)
