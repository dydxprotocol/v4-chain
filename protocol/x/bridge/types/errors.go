package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/bridge module sentinel errors
var (
	ErrBridgeIdNotRecognized = sdkerrors.Register(
		ModuleName,
		1,
		"Bridge event ID is not recognized",
	)
	ErrBridgeIdNotNextToAcknowledge = sdkerrors.Register(
		ModuleName,
		2,
		"Bridge event ID is not the ID to be next acknowledged",
	)
	ErrBridgeIdsNotConsecutive = sdkerrors.Register(
		ModuleName,
		3,
		"Bridge event IDs are not consecutive",
	)
	ErrInvalidAuthority = sdkerrors.Register(
		ModuleName,
		4,
		"Authority is invalid",
	)
	ErrBridgeEventNotFound = sdkerrors.Register(
		ModuleName,
		5,
		"Bridge event not found",
	)
	ErrBridgeEventContentMismatch = sdkerrors.Register(
		ModuleName,
		6,
		"Bridge event content mismatch",
	)

	ErrNegativeDuration = sdkerrors.Register(
		ModuleName,
		400,
		"Duration is negative",
	)
	ErrRateOutOfBounds = sdkerrors.Register(
		ModuleName,
		401,
		"Rate is out of bounds",
	)
)
