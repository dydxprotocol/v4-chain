package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrAuthenticatorNotFound = errorsmod.Register(
		ModuleName,
		1,
		"Authenticator is not found",
	)
	ErrInvalidAccountAddress = errorsmod.Register(
		ModuleName,
		2,
		"Invalid account address",
	)
	ErrAuthenticatorDataExceedsMaximumLength = errorsmod.Register(
		ModuleName,
		3,
		"Authenticator data exceeds maximum length",
	)
	ErrSmartAccountNotActive = errorsmod.Register(
		ModuleName,
		4,
		"Smart account is not active",
	)
	ErrInitializingAuthenticator = errorsmod.Register(
		ModuleName,
		5,
		"Error initializing authenticator",
	)
	ErrTxnHasMultipleSigners = errorsmod.Register(
		ModuleName,
		6,
		"The transaction has multiple signers",
	)

	// Errors for failing authenticator validation
	ErrSignatureVerification = errorsmod.Register(
		ModuleName,
		100,
		"Signature verification failed",
	)
	ErrMessageTypeVerification = errorsmod.Register(
		ModuleName,
		101,
		"Message type verification failed",
	)
	ErrClobPairIdVerification = errorsmod.Register(
		ModuleName,
		102,
		"Clob pair id verification failed",
	)
	ErrSubaccountVerification = errorsmod.Register(
		ModuleName,
		103,
		"Subaccount verification failed",
	)
	ErrAllOfVerification = errorsmod.Register(
		ModuleName,
		104,
		"AllOf verification failed",
	)
	ErrAnyOfVerification = errorsmod.Register(
		ModuleName,
		105,
		"AnyOf verification failed",
	)
)
