package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/assets module sentinel errors
var (
	ErrAssetDoesNotExist            = moderrors.Register(ModuleName, 1, "Asset does not exist")
	ErrNegativeLongInterest         = moderrors.Register(ModuleName, 2, "LongInterest cannot be negative")
	ErrNoAssetWithDenom             = moderrors.Register(ModuleName, 3, "No asset found associated with given denom")
	ErrAssetDenomAlreadyExists      = moderrors.Register(ModuleName, 4, "Existing asset found with the same denom")
	ErrAssetIdAlreadyExists         = moderrors.Register(ModuleName, 5, "Existing asset found with the same asset id")
	ErrGapFoundInAssetId            = moderrors.Register(ModuleName, 6, "Found gap in asset Id")
	ErrAssetZeroNotUsdc             = moderrors.Register(ModuleName, 7, "First asset is not USDC")
	ErrNoAssetInGenesis             = moderrors.Register(ModuleName, 8, "No asset found in genesis state")
	ErrInvalidMarketId              = moderrors.Register(ModuleName, 9, "Found market id for asset without market")
	ErrInvalidAssetAtomicResolution = moderrors.Register(ModuleName, 10, "Invalid asset atomic resolution")
	ErrInvalidDenomExponent         = moderrors.Register(ModuleName, 11, "Invalid denom exponent")

	// Errors for Not Implemented
	ErrNotImplementedMulticollateral = moderrors.Register(ModuleName, 401, "Not Implemented: Multi-Collateral")
	ErrNotImplementedMargin          = moderrors.Register(ModuleName, 402, "Not Implemented: Margin-Trading of Assets")
)
