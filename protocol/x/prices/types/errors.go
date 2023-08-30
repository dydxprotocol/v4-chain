package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// 1 - 99: Default.
	ErrInvalidInput = sdkerrors.Register(ModuleName, 1, "Invalid input")

	// 100 - 199: Exchange related errors.
	ErrExchangeDoesNotExist = sdkerrors.Register(ModuleName, 100, "Exchange does not exist")
	ErrZeroMinExchanges     = sdkerrors.Register(ModuleName, 101, "Min exchanges must be greater than zero")
	ErrTooFewExchanges      = sdkerrors.Register(ModuleName, 102, "Exchanges is fewer than minExchanges")
	ErrDuplicateExchanges   = sdkerrors.Register(
		ModuleName,
		103,
		"Exchanges must not contain duplicates and must be provided in ascending order",
	)

	// 200 - 299: Market related errors.
	ErrMarketParamDoesNotExist        = sdkerrors.Register(ModuleName, 200, "Market param does not exist")
	ErrMarketPriceDoesNotExist        = sdkerrors.Register(ModuleName, 201, "Market price does not exist")
	ErrMarketExponentCannotBeUpdated  = sdkerrors.Register(ModuleName, 202, "Market exponent cannot be updated")
	ErrMarketPricesAndParamsDontMatch = sdkerrors.Register(ModuleName, 203, "Market prices and params don't match")

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
