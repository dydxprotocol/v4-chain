package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/sending module sentinel errors
var (
	ErrSenderSameAsRecipient = errorsmod.Register(ModuleName, 1, "Sender is the same as recipient")
	ErrInvalidTransferAmount = errorsmod.Register(ModuleName, 2, "Invalid transfer amount")
	ErrDuplicatedTransfer    = errorsmod.Register(ModuleName, 3, "Duplicated transfer")
	ErrTransferNotFound      = errorsmod.Register(ModuleName, 4, "Transfer not found")
	ErrMissingFields         = errorsmod.Register(
		ModuleName,
		5,
		"Transfer does not contain all required fields",
	)
	ErrInvalidAccountAddress              = errorsmod.Register(ModuleName, 6, "Account address is invalid")
	ErrEmptyModuleName                    = errorsmod.Register(ModuleName, 7, "Module name is empty")
	ErrInvalidAuthority                   = errorsmod.Register(ModuleName, 8, "Authority is invalid")
	ErrNonUsdcAssetTransferNotImplemented = errorsmod.Register(
		ModuleName,
		1101,
		"Non-USDC asset transfer not implemented",
	)
)
