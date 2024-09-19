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
	ErrMarketParamPairAlreadyExists   = errorsmod.Register(ModuleName, 205, "Market params pair already exists")
	ErrMarketPairConversionFailed     = errorsmod.Register(
		ModuleName,
		206,
		"Market pair conversion to currency pair failed",
	)
	ErrTickerNotFoundInMarketMap  = errorsmod.Register(ModuleName, 207, "Ticker not found in market map")
	ErrMarketCouldNotBeDisabled   = errorsmod.Register(ModuleName, 208, "Market could not be disabled")
	ErrMarketCouldNotBeEnabled    = errorsmod.Register(ModuleName, 209, "Market could not be enabled")
	ErrInvalidMarketPriceExponent = errorsmod.Register(
		ModuleName,
		210,
		"Market price exponent does not match the negation of the Decimals value in the market map",
	)

	// 300 - 399: Price related errors.
	ErrIndexPriceNotAvailable = errorsmod.Register(ModuleName, 300, "Index price is not available")

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
	ErrUnsafeMarketUpdate = errorsmod.Register(
		ModuleName,
		501,
		"Market update is unsafe",
	)

	ErrMarketUpdateChangesMarketMapEnabledValue = errorsmod.Register(
		ModuleName,
		502,
		"Market update changes market map enabled value",
	)

	ErrMarketDoesNotExistInMarketMap = errorsmod.Register(
		ModuleName,
		503,
		"Market does not exist in market map",
	)

	ErrAdditionOfEnabledMarket = errorsmod.Register(
		ModuleName,
		504,
		"Newly added markets must be disabled",
	)
)
