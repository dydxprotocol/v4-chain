package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// Panic strings
const (
	ErrMatchUpdatesMustHaveTwoUpdates = "internalCanUpdateSubaccounts: MATCH subaccount updates must consist of " +
		"exactly 2 updates, got settledUpdates: %+v"
	ErrMatchUpdatesMustUpdateOnePerp = "internalCanUpdateSubaccounts: MATCH subaccount updates must each have " +
		"exactly 1 PerpetualUpdate, got settledUpdates: %+v"
	ErrMatchUpdatesMustBeSamePerpId = "internalCanUpdateSubaccounts: MATCH subaccount updates must consists of two " +
		"updates on same perpetual Id, got settledUpdates: %+v"
	ErrMatchUpdatesInvalidSize = "internalCanUpdateSubaccounts: MATCH subaccount updates must consists of two " +
		"updates of equal absolute base quantums and opposite sign: %+v"
)

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
	ErrInvalidSubaccountIdNumber = errorsmod.Register(
		ModuleName,
		200,
		"subaccount id number cannot exceed "+lib.IntToString(MaxSubaccountIdNumber),
	)
	ErrInvalidSubaccountIdOwner = errorsmod.Register(ModuleName, 201, "subaccount id owner is an invalid address")
	ErrDuplicateSubaccountIds   = errorsmod.Register(ModuleName, 202, "duplicate subaccount id found in genesis")

	// 300 - 399: asset position related.
	ErrAssetPositionsOutOfOrder       = errorsmod.Register(ModuleName, 300, "asset positions are out of order")
	ErrAssetPositionZeroQuantum       = errorsmod.Register(ModuleName, 301, "asset position's quantum cannot be zero")
	ErrAssetPositionNotSupported      = errorsmod.Register(ModuleName, 302, "asset position is not supported")
	ErrMultAssetPositionsNotSupported = errorsmod.Register(
		ModuleName, 303, "having multiple asset positions is not supported")

	// 400 - 499: perpetual position related.
	ErrPerpPositionsOutOfOrder = errorsmod.Register(ModuleName, 400, "perpetual positions are out of order")
	ErrPerpPositionZeroQuantum = errorsmod.Register(
		ModuleName,
		401,
		"perpetual position's quantum cannot be zero",
	)
	ErrCannotModifyPerpOpenInterestForOIMF = errorsmod.Register(
		ModuleName,
		402,
		"cannot modify perpetual open interest for OIMF calculation",
	)
	ErrCannotRevertPerpOpenInterestForOIMF = errorsmod.Register(
		ModuleName,
		403,
		"cannot revert perpetual open interest for OIMF calculation",
	)
	ErrPerpetualInfoDoesNotExist = errorsmod.Register(
		ModuleName,
		404,
		"PerpetualInfo does not exist in map",
	)

	// 500 - 599: transfer related.
	ErrAssetTransferQuantumsNotPositive = errorsmod.Register(
		ModuleName, 500, "asset transfer quantums is not positive")
	ErrAssetTransferThroughBankNotImplemented = errorsmod.Register(
		ModuleName, 501, "asset transfer (other than USDC) through the bank module is not implemented")
)
