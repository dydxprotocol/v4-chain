package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// 1 - 99: Default.
	ErrInvalidInput = sdkerrors.Register(ModuleName, 1, "Invalid input")

	// 100 - 199: Exchange related errors.
	ErrExchangeFeedDoesNotExist = sdkerrors.Register(ModuleName, 100, "ExchangeFeed does not exist")
	ErrZeroMinExchanges         = sdkerrors.Register(ModuleName, 101, "Min exchanges must be greater than zero")
	ErrTooFewExchanges          = sdkerrors.Register(ModuleName, 102, "Exchanges is fewer than minExchanges")
	ErrDuplicateExchanges       = sdkerrors.Register(
		ModuleName,
		103,
		"Exchanges must not contain duplicates and must be provided in ascending order",
	)

	// 200 - 299: Market related errors.
	ErrMarketDoesNotExist = sdkerrors.Register(ModuleName, 200, "Market does not exist")

	// 300 - 399: Price related errors.
	ErrIndexPriceNotAvailable = sdkerrors.Register(ModuleName, 300, "Index price is not available")

	// 400 - 499: Market price update related errors.
	ErrInvalidMarketPriceUpdateStateless = sdkerrors.Register(
		ModuleName, 400, "Market price update is invalid: stateless.")
	ErrInvalidMarketPriceUpdateDeterministic = sdkerrors.Register(
		ModuleName, 401, "Market price update is invalid: deterministic.")
	ErrInvalidMarketPriceUpdateNonDeterministic = sdkerrors.Register(
		ModuleName, 402, "Market price update is invalid: non-deterministic.")
)
