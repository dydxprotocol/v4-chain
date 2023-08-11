package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/epochs module sentinel errors
var (
	ErrEmptyEpochInfoName = sdkerrors.Register(
		ModuleName,
		2,
		"EpochInfo name is empty",
	)
	ErrDurationIsZero = sdkerrors.Register(
		ModuleName,
		3,
		"Duration is zero",
	)
	ErrEpochInfoAlreadyExists = sdkerrors.Register(
		ModuleName,
		4,
		"EpochInfo name already exists",
	)
	ErrEpochInfoNotFound = sdkerrors.Register(
		ModuleName,
		5,
		"EpochInfo name not found",
	)
	ErrInvalidCurrentEpochAndCurrentEpochStartBlockTuple = sdkerrors.Register(
		ModuleName,
		6,
		"Invalid CurrentEpoch and CurrentEpochStartBlock tuple: CurrentEpoch should"+
			" be zero if and only if CurrentEpochStartBlock is zero",
	)
)
