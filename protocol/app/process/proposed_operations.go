package process

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process/errors"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var (
	msgProposedOperationsType = reflect.TypeOf(types.MsgProposedOperations{})
)

// ProposedOperationsTx represents `MsgProposedOperations` tx that can be validated.
type ProposedOperationsTx struct {
	msg *types.MsgProposedOperations
}

// DecodeProposedOperationsTx returns a new `ProposedOperationsTx` after validating the following:
//   - decodes the given tx bytes
//   - checks the num of msgs in the tx matches expectations
//   - checks the msg is of expected type
//
// If error occurs during any of the checks, returns error.
func DecodeProposedOperationsTx(decoder sdk.TxDecoder, txBytes []byte) (*ProposedOperationsTx, error) {
	// Decode.
	tx, err := decoder(txBytes)
	if err != nil {
		return nil, errors.GetDecodingError(msgProposedOperationsType, err)
	}

	// Check msg length.
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return nil, errors.GetUnexpectedNumMsgsError(msgProposedOperationsType, 1, len(msgs))
	}

	// Check msg type.
	proposedOperations, ok := msgs[0].(*types.MsgProposedOperations)
	if !ok {
		return nil, errors.GetUnexpectedMsgTypeError(msgProposedOperationsType, msgs[0])
	}

	return &ProposedOperationsTx{msg: proposedOperations}, nil
}

// Validate returns an error if the underlying msg fails `ValidateBasic`.
func (pmot *ProposedOperationsTx) Validate() error {
	if err := pmot.msg.ValidateBasic(); err != nil {
		return errors.GetValidateBasicError(pmot.msg, err)
	}
	return nil
}

// GetMsg returns the underlying `MsgProposedOperations`.
func (pmot *ProposedOperationsTx) GetMsg() sdk.Msg {
	return pmot.msg
}
