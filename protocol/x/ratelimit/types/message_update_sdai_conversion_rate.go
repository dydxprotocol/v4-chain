package types

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateSDAIConversionRate{}

// ValidateBasic runs validation on the fields of a MsgUpdateSDAIConversionRate.
func (msg *MsgUpdateSDAIConversionRate) ValidateBasic() error {

	// Validate account sender.
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(
			ErrInvalidSender,
			fmt.Sprintf(
				"authority '%s' must be a valid bech32 address, but got error '%v'",
				msg.Sender,
				err.Error(),
			),
		)
	}

	bigConversionRate, ok := new(big.Int).SetString(msg.ConversionRate, 10)
	if !ok {
		return errorsmod.Wrap(
			ErrUnableToDecodeBigInt,
			fmt.Sprintf(
				"Unable to convert the sDAI conversion rate to a big int",
			),
		)
	}

	// Validate that the conversion rate is positive // BigInt().Sign()
	if bigConversionRate.Sign() <= 0 {
		return errorsmod.Wrap(
			ErrValueIsNegative,
			fmt.Sprintf(
				"Invalid sDAI conversion rate",
			),
		)
	}

	bigEthereumBlockNumber, ok := new(big.Int).SetString(msg.ConversionRate, 10)
	if !ok {
		return errorsmod.Wrap(
			ErrUnableToDecodeBigInt,
			fmt.Sprintf(
				"Unable to convert the ethereum block number to a big int",
			),
		)
	}

	// Validate that eth ethereum block height is positive
	if bigEthereumBlockNumber.Sign() <= 0 {
		return errorsmod.Wrap(
			ErrValueIsNegative,
			fmt.Sprintf(
				"Invalid ethereum block number",
			),
		)
	}

	return nil
}
