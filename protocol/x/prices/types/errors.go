package types

import moderrors "cosmossdk.io/errors"

// DONTCOVER

var (
	// 1 - 99: Default.
	ErrInvalidInput = moderrors.Register(ModuleName, 1, "Invalid input")

	// 100 - 199: Exchange related errors.
	ErrExchangeDoesNotExist = moderrors.Register(ModuleName, 100, "Exchange does not exist")
	ErrZeroMinExchanges     = moderrors.Register(ModuleName, 101, "Min exchanges must be greater than zero")
	ErrTooFewExchanges      = moderrors.Register(ModuleName, 102, "Exchanges is fewer than minExchanges")
	ErrDuplicateExchanges   = moderrors.Register(
		ModuleName,
		103,
		"Exchanges must not contain duplicates and must be provided in ascending order",
	)

	// 200 - 299: Market related errors.
	ErrMarketParamDoesNotExist        = moderrors.Register(ModuleName, 200, "Market param does not exist")
	ErrMarketPriceDoesNotExist        = moderrors.Register(ModuleName, 201, "Market price does not exist")
	ErrMarketExponentCannotBeUpdated  = moderrors.Register(ModuleName, 202, "Market exponent cannot be updated")
	ErrMarketPricesAndParamsDontMatch = moderrors.Register(ModuleName, 203, "Market prices and params don't match")

	// 300 - 399: Price related errors.
	ErrIndexPriceNotAvailable = moderrors.Register(ModuleName, 300, "Index price is not available")

	// 400 - 499: Market price update related errors.
	ErrInvalidMarketPriceUpdateStateless = moderrors.Register(
		ModuleName, 400, "Market price update is invalid: stateless.")
	ErrInvalidMarketPriceUpdateDeterministic = moderrors.Register(
		ModuleName, 401, "Market price update is invalid: deterministic.")
	ErrInvalidMarketPriceUpdateNonDeterministic = moderrors.Register(
		ModuleName, 402, "Market price update is invalid: non-deterministic.")
)
