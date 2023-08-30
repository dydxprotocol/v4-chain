package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/sending module sentinel errors
var (
	ErrSenderSameAsRecipient       = sdkerrors.Register(ModuleName, 1, "Sender is the same as recipient")
	ErrInvalidTransferAmount       = sdkerrors.Register(ModuleName, 2, "Invalid transfer amount")
	ErrDuplicatedTransfer          = sdkerrors.Register(ModuleName, 3, "Duplicated transfer")
	ErrTransferNotFound            = sdkerrors.Register(ModuleName, 4, "Transfer not found")
	ErrMissingFields               = sdkerrors.Register(ModuleName, 5, "Transfer does not contain all required fields")
	ErrInvalidAccountAddress       = sdkerrors.Register(ModuleName, 6, "Account address is invalid")
	ErrKeeperMethodsNotImplemented = sdkerrors.Register(
		ModuleName,
		1100,
		"Sending module keeper method not implemented",
	)
	ErrNonUsdcAssetTransferNotImplemented = sdkerrors.Register(
		ModuleName,
		1101,
		"Non-USDC asset transfer not implemented",
	)
)
