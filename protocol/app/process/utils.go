package process

import (
	errorsmod "cosmossdk.io/errors"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// getValidateBasicError returns a sdk error for `Msg.ValidateBasic` failure.
func getValidateBasicError(msg sdk.Msg, err error) error {
	return errorsmod.Wrapf(ErrMsgValidateBasic, "Msg Type: %T, Error: %+v", msg, err)
}

// getDecodingError returns a sdk error for tx decoding failure.
func getDecodingError(msgType reflect.Type, err error) error {
	return errorsmod.Wrapf(ErrDecodingTxBytes, "Msg Type: %s, Error: %+v", msgType, err)
}

// getUnexpectedNumMsgsError returns a sdk error for having unexpected num of msgs in the tx.
func getUnexpectedNumMsgsError(msgType reflect.Type, expectedNum int, actualNum int) error {
	return errorsmod.Wrapf(
		ErrUnexpectedNumMsgs,
		"Msg Type: %s, Expected %d num of msgs, but got %d",
		msgType,
		expectedNum,
		actualNum,
	)
}

// getUnexpectedMsgTypeError returns a sdk error for having unexpected msg type in the tx.
func getUnexpectedMsgTypeError(expectedMsgType reflect.Type, actualMsg sdk.Msg) error {
	return errorsmod.Wrapf(
		ErrUnexpectedMsgType, "Expected MsgType %s, but got %T", expectedMsgType, actualMsg,
	)
}
