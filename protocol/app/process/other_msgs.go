package process

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	libprocess "github.com/dydxprotocol/v4-chain/protocol/lib/process"
)

// OtherMsgsTx represents tx msgs in the "other" category that can be validated.
type OtherMsgsTx struct {
	msgs []sdk.Msg
}

// DecodeOtherMsgsTx returns a new `OtherMsgsTx` after validating the following:
//   - decodes the given tx bytes
//   - checks the num of msgs in the tx is not 0
//   - checks the msgs do not contain "app-injected msgs"
//
// If error occurs during any of the checks, returns error.
func DecodeOtherMsgsTx(decoder sdk.TxDecoder, txBytes []byte) (*OtherMsgsTx, error) {
	// Decode.
	tx, err := decoder(txBytes)
	if err != nil {
		return nil, errorsmod.Wrapf(ErrDecodingTxBytes, "OtherMsgsTx Error: %+v", err)
	}

	// Check msg length.
	allMsgs := tx.GetMsgs()
	if len(allMsgs) == 0 {
		return nil, errorsmod.Wrapf(ErrUnexpectedNumMsgs, "OtherMsgs len cannot be zero")
	}

	// Check msg type.
	for _, msg := range allMsgs {
		if ante.IsDisallowExternalSubmitMsg(msg) {
			return nil,
				errorsmod.Wrapf(
					ErrUnexpectedMsgType,
					"Invalid msg type or content in OtherTxs %T",
					msg,
				)
		}

		if libprocess.IsDisallowTopLevelMsgInOtherTxs(msg) {
			return nil,
				errorsmod.Wrapf(
					ErrUnexpectedMsgType,
					"Msg type %T is not allowed in OtherTxs",
					msg,
				)
		}
	}

	return &OtherMsgsTx{msgs: allMsgs}, nil
}

// Validate returns an error if one of the underlying msgs fails `ValidateBasic`.
func (omt *OtherMsgsTx) Validate() error {
	for _, msg := range omt.msgs {
		if err := msg.ValidateBasic(); err != nil {
			return getValidateBasicError(msg, err)
		}
	}
	return nil
}

// GetMsgs returns the underlying msgs in the tx.
func (omt *OtherMsgsTx) GetMsgs() []sdk.Msg {
	return omt.msgs
}
