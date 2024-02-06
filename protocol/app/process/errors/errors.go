package errors

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

const (
	ModuleName = "process_proposal"
)

var (
	// 1 - 99: Default.
	ErrDecodingTxBytes   = errorsmod.Register(ModuleName, 1, "Decoding tx bytes failed")
	ErrMsgValidateBasic  = errorsmod.Register(ModuleName, 2, "ValidateBasic failed on msg")
	ErrUnexpectedNumMsgs = errorsmod.Register(ModuleName, 3, "Unexpected num of msgs")
	ErrUnexpectedMsgType = errorsmod.Register(ModuleName, 4, "Unexpected msg type")
)

// GetValidateBasicError returns a sdk error for `Msg.ValidateBasic` failure.
func GetValidateBasicError(msg sdk.Msg, err error) error {
	return errorsmod.Wrapf(ErrMsgValidateBasic, "Msg Type: %T, Error: %+v", msg, err)
}

// GetDecodingError returns a sdk error for tx decoding failure.
func GetDecodingError(msgType reflect.Type, err error) error {
	return errorsmod.Wrapf(ErrDecodingTxBytes, "Msg Type: %s, Error: %+v", msgType, err)
}

// GetUnexpectedNumMsgsError returns a sdk error for having unexpected num of msgs in the tx.
func GetUnexpectedNumMsgsError(msgType reflect.Type, expectedNum int, actualNum int) error {
	return errorsmod.Wrapf(
		ErrUnexpectedNumMsgs,
		"Msg Type: %s, Expected %d num of msgs, but got %d",
		msgType,
		expectedNum,
		actualNum,
	)
}

// GetUnexpectedMsgTypeError returns a sdk error for having unexpected msg type in the tx.
func GetUnexpectedMsgTypeError(expectedMsgType reflect.Type, actualMsg sdk.Msg) error {
	return errorsmod.Wrapf(
		ErrUnexpectedMsgType, "Expected MsgType %s, but got %T", expectedMsgType, actualMsg,
	)
}
