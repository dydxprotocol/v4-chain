package prices

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/dydxprotocol/v4-chain/protocol/app/process/errors"
)

var (
	ErrProposedPriceValidation = errorsmod.Register(errors.ModuleName, 5, "Validation of proposed MsgUpdateMarketPrices failed")
)

func IncorrectNumberUpdatesError(expected, actual int) error {
	return errorsmod.Wrapf(
		ErrProposedPriceValidation,
		"incorrect number of price-updates, expected: %d, actual: %d",
		expected,
		actual,
	)
}

func MissingPriceUpdateForMarket(marketID uint32) error {
	return errorsmod.Wrapf(
		ErrProposedPriceValidation,
		"missing price-update for market: %d",
		marketID,
	)
}

func IncorrectPriceUpdateForMarket(marketID uint32, expected, actual uint64) error {
	return errorsmod.Wrapf(
		ErrProposedPriceValidation,
		"incorrect price-update for market: %d, expected %d, actual %d",
		marketID,
		expected,
		actual,
	)
}
