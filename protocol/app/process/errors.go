package process

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

const (
	ModuleName = "process_proposal"
)

var (
	// 1 - 99: Default.
	ErrDecodingTxBytes         = errorsmod.Register(ModuleName, 1, "Decoding tx bytes failed")
	ErrMsgValidateBasic        = errorsmod.Register(ModuleName, 2, "ValidateBasic failed on msg")
	ErrUnexpectedNumMsgs       = errorsmod.Register(ModuleName, 3, "Unexpected num of msgs")
	ErrUnexpectedMsgType       = errorsmod.Register(ModuleName, 4, "Unexpected msg type")
	ErrProposedPriceValidation = errorsmod.Register(ModuleName, 5, "Validation of proposed MsgUpdateMarketPrices failed")
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

func InvalidMarketPriceUpdateError(err error) error {
	return errorsmod.Wrap(ErrProposedPriceValidation, err.Error())
}
