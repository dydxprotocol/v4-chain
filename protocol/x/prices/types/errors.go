package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	// 1 - 99: Default.
	ErrInvalidInput = errorsmod.Register(ModuleName, 1, "Invalid input")

	// 100 - 199: Exchange related errors.
	ErrExchangeDoesNotExist = errorsmod.Register(ModuleName, 100, "Exchange does not exist")
	ErrZeroMinExchanges     = errorsmod.Register(ModuleName, 101, "Min exchanges must be greater than zero")
	ErrTooFewExchanges      = errorsmod.Register(ModuleName, 102, "Exchanges is fewer than minExchanges")
	ErrDuplicateExchanges   = errorsmod.Register(
		ModuleName,
		103,
		"Exchanges must not contain duplicates and must be provided in ascending order",
	)

	// 200 - 299: Market related errors.
	ErrMarketParamDoesNotExist        = errorsmod.Register(ModuleName, 200, "Market param does not exist")
	ErrMarketPriceDoesNotExist        = errorsmod.Register(ModuleName, 201, "Market price does not exist")
	ErrMarketExponentCannotBeUpdated  = errorsmod.Register(ModuleName, 202, "Market exponent cannot be updated")
	ErrMarketPricesAndParamsDontMatch = errorsmod.Register(ModuleName, 203, "Market prices and params don't match")
	ErrMarketParamAlreadyExists       = errorsmod.Register(ModuleName, 204, "Market params already exists")

	// 300 - 399: Price related errors.
	ErrDaemonPriceNotAvailable = errorsmod.Register(ModuleName, 300, "daemon price is not available")

	// 400 - 499: Market price update related errors.
	ErrInvalidMarketPriceUpdateStateless = errorsmod.Register(
		ModuleName, 400, "Market price update is invalid: stateless.")
	ErrInvalidMarketPriceUpdateDeterministic = errorsmod.Register(
		ModuleName, 401, "Market price update is invalid: deterministic.")
	ErrInvalidMarketPriceUpdateNonDeterministic = errorsmod.Register(
		ModuleName, 402, "Market price update is invalid: non-deterministic.")

	// 500 - 599: sdk.Msg related errors.
	ErrInvalidAuthority = errorsmod.Register(
		ModuleName,
		500,
		"Authority is invalid",
	)
)
