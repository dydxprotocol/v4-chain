package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/subaccounts module sentinel errors
var (
	// 0 - 99: generic.
	ErrIntegerOverflow = errorsmod.Register(ModuleName, 0, "integer overflow")

	// 100 - 199: update related.
	ErrNonUniqueUpdatesPosition = errorsmod.Register(
		ModuleName, 100, "multiple updates were specified for the same position id")
	ErrNonUniqueUpdatesSubaccount = errorsmod.Register(
		ModuleName, 101, "multiple updates were specified for the same subaccountId")
	ErrFailedToUpdateSubaccounts   = errorsmod.Register(ModuleName, 102, "failed to apply subaccount updates")
	ErrProductPositionNotUpdatable = errorsmod.Register(ModuleName, 103, "product position is not updatable")

	// 200 - 299: subaccount id related.
	ErrInvalidSubaccountIdNumber = errorsmod.Register(ModuleName, 200, "subaccount id number cannot exceed 127")
	ErrInvalidSubaccountIdOwner  = errorsmod.Register(ModuleName, 201, "subaccount id owner is an invalid address")
	ErrDuplicateSubaccountIds    = errorsmod.Register(ModuleName, 202, "duplicate subaccount id found in genesis")

	// 300 - 399: asset position related.
	ErrAssetPositionsOutOfOrder       = errorsmod.Register(ModuleName, 300, "asset positions are out of order")
	ErrAssetPositionZeroQuantum       = errorsmod.Register(ModuleName, 301, "asset position's quantum cannot be zero")
	ErrAssetPositionNotSupported      = errorsmod.Register(ModuleName, 302, "asset position is not supported")
	ErrMultAssetPositionsNotSupported = errorsmod.Register(
		ModuleName, 303, "having multiple asset positions is not supported")

	// 400 - 499: perpetual position related.
	ErrPerpPositionsOutOfOrder = errorsmod.Register(ModuleName, 400, "perpetual positions are out of order")
	ErrPerpPositionZeroQuantum = errorsmod.Register(ModuleName, 401, "perpetual position's quantum cannot be zero")

	// 500 - 599: transfer related.
	ErrAssetTransferQuantumsNotPositive = errorsmod.Register(
		ModuleName, 500, "asset transfer quantums is not positive")
	ErrAssetTransferThroughBankNotImplemented = errorsmod.Register(
		ModuleName, 501, "asset transfer (other than USDC) through the bank module is not implemented")
)
