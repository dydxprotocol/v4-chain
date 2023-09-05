package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/bridge module sentinel errors
var (
	ErrBridgeIdNotRecognized = moderrors.Register(
		ModuleName,
		1,
		"Bridge event ID is not recognized",
	)
	ErrBridgeIdNotNextToAcknowledge = moderrors.Register(
		ModuleName,
		2,
		"Bridge event ID is not the ID to be next acknowledged",
	)
	ErrBridgeIdsNotConsecutive = moderrors.Register(
		ModuleName,
		3,
		"Bridge event IDs are not consecutive",
	)
	ErrInvalidAuthority = moderrors.Register(
		ModuleName,
		4,
		"Authority is invalid",
	)
	ErrBridgeEventNotFound = moderrors.Register(
		ModuleName,
		5,
		"Bridge event not found",
	)
	ErrBridgeEventContentMismatch = moderrors.Register(
		ModuleName,
		6,
		"Bridge event content mismatch",
	)

	ErrNegativeDuration = moderrors.Register(
		ModuleName,
		400,
		"Duration is negative",
	)
	ErrRateOutOfBounds = moderrors.Register(
		ModuleName,
		401,
		"Rate is out of bounds",
	)
)
