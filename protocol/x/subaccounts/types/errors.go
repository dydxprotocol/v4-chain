package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/subaccounts module sentinel errors
var (
	// 0 - 99: generic.
	ErrIntegerOverflow = moderrors.Register(ModuleName, 0, "integer overflow")

	// 100 - 199: update related.
	ErrNonUniqueUpdatesPosition = moderrors.Register(
		ModuleName, 100, "multiple updates were specified for the same position id")
	ErrNonUniqueUpdatesSubaccount = moderrors.Register(
		ModuleName, 101, "multiple updates were specified for the same subaccountId")
	ErrFailedToUpdateSubaccounts = moderrors.Register(ModuleName, 102, "failed to apply subaccount updates")

	// 200 - 299: subaccount id related.
	ErrInvalidSubaccountIdNumber = moderrors.Register(ModuleName, 200, "subaccount id number cannot exceed 127")
	ErrInvalidSubaccountIdOwner  = moderrors.Register(ModuleName, 201, "subaccount id owner is an invalid address")
	ErrDuplicateSubaccountIds    = moderrors.Register(ModuleName, 202, "duplicate subaccount id found in genesis")

	// 300 - 399: asset position related.
	ErrAssetPositionsOutOfOrder       = moderrors.Register(ModuleName, 300, "asset positions are out of order")
	ErrAssetPositionZeroQuantum       = moderrors.Register(ModuleName, 301, "asset position's quantum cannot be zero")
	ErrAssetPositionNotSupported      = moderrors.Register(ModuleName, 302, "asset position is not supported")
	ErrMultAssetPositionsNotSupported = moderrors.Register(
		ModuleName, 303, "having multiple asset positions is not supported")

	// 400 - 499: perpetual position related.
	ErrPerpPositionsOutOfOrder = moderrors.Register(ModuleName, 400, "perpetual positions are out of order")
	ErrPerpPositionZeroQuantum = moderrors.Register(ModuleName, 401, "perpetual position's quantum cannot be zero")

	// 500 - 599: transfer related.
	ErrAssetTransferQuantumsNotPositive = moderrors.Register(
		ModuleName, 500, "asset transfer quantums is not positive")
	ErrAssetTransferThroughBankNotImplemented = moderrors.Register(
		ModuleName, 501, "asset transfer (other than USDC) through the bank module is not implemented")
)
