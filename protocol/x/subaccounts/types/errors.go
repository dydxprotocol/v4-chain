package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/subaccounts module sentinel errors
var (
	// 0 - 99: generic.
	ErrIntegerOverflow = sdkerrors.Register(ModuleName, 0, "integer overflow")

	// 100 - 199: update related.
	ErrNonUniqueUpdatesPosition = sdkerrors.Register(
		ModuleName, 100, "multiple updates were specified for the same position id")
	ErrNonUniqueUpdatesSubaccount = sdkerrors.Register(
		ModuleName, 101, "multiple updates were specified for the same subaccountId")
	ErrFailedToUpdateSubaccounts = sdkerrors.Register(ModuleName, 102, "failed to apply subaccount updates")

	// 200 - 299: subaccount id related.
	ErrInvalidSubaccountIdNumber = sdkerrors.Register(ModuleName, 200, "subaccount id number cannot exceed 127")
	ErrInvalidSubaccountIdOwner  = sdkerrors.Register(ModuleName, 201, "subaccount id owner is an invalid address")
	ErrDuplicateSubaccountIds    = sdkerrors.Register(ModuleName, 202, "duplicate subaccount id found in genesis")

	// 300 - 399: asset position related.
	ErrAssetPositionsOutOfOrder       = sdkerrors.Register(ModuleName, 300, "asset positions are out of order")
	ErrAssetPositionZeroQuantum       = sdkerrors.Register(ModuleName, 301, "asset position's quantum cannot be zero")
	ErrAssetPositionNotSupported      = sdkerrors.Register(ModuleName, 302, "asset position is not supported")
	ErrMultAssetPositionsNotSupported = sdkerrors.Register(
		ModuleName, 303, "having multiple asset positions is not supported")

	// 400 - 499: perpetual position related.
	ErrPerpPositionsOutOfOrder = sdkerrors.Register(ModuleName, 400, "perpetual positions are out of order")
	ErrPerpPositionZeroQuantum = sdkerrors.Register(ModuleName, 401, "perpetual position's quantum cannot be zero")

	// 500 - 599: transfer related.
	ErrAssetTransferQuantumsNotPositive = sdkerrors.Register(
		ModuleName, 500, "asset transfer quantums is not positive")
	ErrAssetTransferThroughBankNotImplemented = sdkerrors.Register(
		ModuleName, 501, "asset transfer (other than USDC) through the bank module is not implemented")
)
