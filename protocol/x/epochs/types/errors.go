package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/epochs module sentinel errors
var (
	ErrEmptyEpochInfoName = moderrors.Register(
		ModuleName,
		2,
		"EpochInfo name is empty",
	)
	ErrDurationIsZero = moderrors.Register(
		ModuleName,
		3,
		"Duration is zero",
	)
	ErrEpochInfoAlreadyExists = moderrors.Register(
		ModuleName,
		4,
		"EpochInfo name already exists",
	)
	ErrEpochInfoNotFound = moderrors.Register(
		ModuleName,
		5,
		"EpochInfo name not found",
	)
	ErrInvalidCurrentEpochAndCurrentEpochStartBlockTuple = moderrors.Register(
		ModuleName,
		6,
		"Invalid CurrentEpoch and CurrentEpochStartBlock tuple: CurrentEpoch should"+
			" be zero if and only if CurrentEpochStartBlock is zero",
	)
)
