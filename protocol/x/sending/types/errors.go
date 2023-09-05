package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

// x/sending module sentinel errors
var (
	ErrSenderSameAsRecipient       = moderrors.Register(ModuleName, 1, "Sender is the same as recipient")
	ErrInvalidTransferAmount       = moderrors.Register(ModuleName, 2, "Invalid transfer amount")
	ErrDuplicatedTransfer          = moderrors.Register(ModuleName, 3, "Duplicated transfer")
	ErrTransferNotFound            = moderrors.Register(ModuleName, 4, "Transfer not found")
	ErrMissingFields               = moderrors.Register(ModuleName, 5, "Transfer does not contain all required fields")
	ErrInvalidAccountAddress       = moderrors.Register(ModuleName, 6, "Account address is invalid")
	ErrKeeperMethodsNotImplemented = moderrors.Register(
		ModuleName,
		1100,
		"Sending module keeper method not implemented",
	)
	ErrNonUsdcAssetTransferNotImplemented = moderrors.Register(
		ModuleName,
		1101,
		"Non-USDC asset transfer not implemented",
	)
)
