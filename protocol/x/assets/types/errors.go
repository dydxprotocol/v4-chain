package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/assets module sentinel errors
var (
	ErrAssetDoesNotExist            = sdkerrors.Register(ModuleName, 1, "Asset does not exist")
	ErrNegativeLongInterest         = sdkerrors.Register(ModuleName, 2, "LongInterest cannot be negative")
	ErrNoAssetWithDenom             = sdkerrors.Register(ModuleName, 3, "No asset found associated with given denom")
	ErrAssetDenomAlreadyExists      = sdkerrors.Register(ModuleName, 4, "Existing asset found with the same denom")
	ErrAssetIdAlreadyExists         = sdkerrors.Register(ModuleName, 5, "Existing asset found with the same asset id")
	ErrGapFoundInAssetId            = sdkerrors.Register(ModuleName, 6, "Found gap in asset Id")
	ErrAssetZeroNotUsdc             = sdkerrors.Register(ModuleName, 7, "First asset is not USDC")
	ErrNoAssetInGenesis             = sdkerrors.Register(ModuleName, 8, "No asset found in genesis state")
	ErrInvalidMarketId              = sdkerrors.Register(ModuleName, 9, "Found market id for asset without market")
	ErrInvalidAssetAtomicResolution = sdkerrors.Register(ModuleName, 10, "Invalid asset atomic resolution")
	ErrInvalidDenomExponent         = sdkerrors.Register(ModuleName, 11, "Invalid denom exponent")

	// Errors for Not Implemented
	ErrNotImplementedMulticollateral = sdkerrors.Register(ModuleName, 401, "Not Implemented: Multi-Collateral")
	ErrNotImplementedMargin          = sdkerrors.Register(ModuleName, 402, "Not Implemented: Margin-Trading of Assets")
)
