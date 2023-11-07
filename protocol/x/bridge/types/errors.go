package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/bridge module sentinel errors
var (
	ErrBridgeIdNotRecognized = errorsmod.Register(
		ModuleName,
		1,
		"Bridge event ID is not recognized",
	)
	ErrBridgeIdNotNextToAcknowledge = errorsmod.Register(
		ModuleName,
		2,
		"Bridge event ID is not the ID to be next acknowledged",
	)
	ErrBridgeIdsNotConsecutive = errorsmod.Register(
		ModuleName,
		3,
		"Bridge event IDs are not consecutive",
	)
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		4,
		"Authority is invalid",
	)
	ErrBridgeEventNotFound = errorsmod.Register(
		ModuleName,
		5,
		"Bridge event not found",
	)
	ErrBridgeEventContentMismatch = errorsmod.Register(
		ModuleName,
		6,
		"Bridge event content mismatch",
	)
	ErrInvalidEthAddress = errorsmod.Register(
		ModuleName,
		7,
		"Invalid Ethereum address",
	)

	ErrNegativeDuration = errorsmod.Register(
		ModuleName,
		400,
		"Duration is negative",
	)
	ErrRateOutOfBounds = errorsmod.Register(
		ModuleName,
		401,
		"Rate is out of bounds",
	)
	ErrBridgingDisabled = errorsmod.Register(
		ModuleName,
		402,
		"Bridging is disabled",
	)
)
